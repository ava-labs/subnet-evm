package evm

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"math/big"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/txpool"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/eth"
	"github.com/ava-labs/subnet-evm/eth/filters"
	"github.com/ava-labs/subnet-evm/metrics"
	"github.com/ava-labs/subnet-evm/plugin/evm/orderbook"
	"github.com/ava-labs/subnet-evm/utils"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

const (
	memoryDBSnapshotKey string = "memoryDBSnapshot"
	snapshotInterval    uint64 = 1000 // save snapshot every 1000 blocks
)

type LimitOrderProcesser interface {
	ListenAndProcessTransactions(blockBuilder *blockBuilder)
	GetOrderBookAPI() *orderbook.OrderBookAPI
	GetTradingAPI() *orderbook.TradingAPI
}

type limitOrderProcesser struct {
	ctx                    *snow.Context
	mu                     *sync.Mutex
	txPool                 *txpool.TxPool
	shutdownChan           <-chan struct{}
	shutdownWg             *sync.WaitGroup
	backend                *eth.EthAPIBackend
	blockChain             *core.BlockChain
	memoryDb               orderbook.LimitOrderDatabase
	limitOrderTxProcessor  orderbook.LimitOrderTxProcessor
	contractEventProcessor *orderbook.ContractEventsProcessor
	matchingPipeline       *orderbook.MatchingPipeline
	filterAPI              *filters.FilterAPI
	hubbleDB               database.Database
	configService          orderbook.IConfigService
	blockBuilder           *blockBuilder
}

func NewLimitOrderProcesser(ctx *snow.Context, txPool *txpool.TxPool, shutdownChan <-chan struct{}, shutdownWg *sync.WaitGroup, backend *eth.EthAPIBackend, blockChain *core.BlockChain, hubbleDB database.Database, validatorPrivateKey string) LimitOrderProcesser {
	log.Info("**** NewLimitOrderProcesser")
	configService := orderbook.NewConfigService(blockChain)
	memoryDb := orderbook.NewInMemoryDatabase(configService)
	lotp := orderbook.NewLimitOrderTxProcessor(txPool, memoryDb, backend, validatorPrivateKey)
	contractEventProcessor := orderbook.NewContractEventsProcessor(memoryDb)
	matchingPipeline := orderbook.NewMatchingPipeline(memoryDb, lotp, configService)
	filterSystem := filters.NewFilterSystem(backend, filters.Config{})
	filterAPI := filters.NewFilterAPI(filterSystem)

	// need to register the types for gob encoding because memory DB has an interface field(ContractOrder)
	gob.Register(&orderbook.LimitOrder{})
	gob.Register(&orderbook.IOCOrder{})
	return &limitOrderProcesser{
		ctx:                    ctx,
		mu:                     &sync.Mutex{},
		txPool:                 txPool,
		shutdownChan:           shutdownChan,
		shutdownWg:             shutdownWg,
		backend:                backend,
		memoryDb:               memoryDb,
		hubbleDB:               hubbleDB,
		blockChain:             blockChain,
		limitOrderTxProcessor:  lotp,
		contractEventProcessor: contractEventProcessor,
		matchingPipeline:       matchingPipeline,
		filterAPI:              filterAPI,
		configService:          configService,
	}
}

func (lop *limitOrderProcesser) ListenAndProcessTransactions(blockBuilder *blockBuilder) {
	lop.mu.Lock()

	lastAccepted := lop.blockChain.LastAcceptedBlock()
	lastAcceptedBlockNumber := lastAccepted.Number()
	if lastAcceptedBlockNumber.Sign() > 0 {
		fromBlock := big.NewInt(0)

		// first load the last snapshot containing finalised data till block x and query the logs of [x+1, latest]
		acceptedBlockNumber, err := lop.loadMemoryDBSnapshot()
		if err != nil {
			log.Error("ListenAndProcessTransactions - error in loading snapshot", "err", err)
		} else {
			if acceptedBlockNumber > 0 {
				fromBlock = big.NewInt(int64(acceptedBlockNumber) + 1)
				log.Info("ListenAndProcessTransactions - memory DB snapshot loaded", "acceptedBlockNumber", acceptedBlockNumber)
			} else {
				// not an error, but unlikely after the blockchain is running for some time
				log.Warn("ListenAndProcessTransactions - no snapshot found")
			}
		}

		logHandler := log.Root().GetHandler()
		log.Info("ListenAndProcessTransactions - beginning sync", " till block number", lastAcceptedBlockNumber)
		JUMP := big.NewInt(3999)
		toBlock := utils.BigIntMin(lastAcceptedBlockNumber, big.NewInt(0).Add(fromBlock, JUMP))
		for toBlock.Cmp(fromBlock) > 0 {
			logs := lop.getLogs(fromBlock, toBlock)
			// set the log handler to discard logs so that the ProcessEvents doesn't spam the logs
			log.Root().SetHandler(log.DiscardHandler())
			lop.contractEventProcessor.ProcessEvents(logs)
			lop.contractEventProcessor.ProcessAcceptedEvents(logs, true)
			lop.memoryDb.Accept(toBlock.Uint64(), 0) // will delete stale orders from the memorydb
			log.Root().SetHandler(logHandler)
			log.Info("ListenAndProcessTransactions - processed log chunk", "fromBlock", fromBlock.String(), "toBlock", toBlock.String(), "number of logs", len(logs))

			fromBlock = fromBlock.Add(toBlock, big.NewInt(1))
			toBlock = utils.BigIntMin(lastAcceptedBlockNumber, big.NewInt(0).Add(fromBlock, JUMP))
		}
		lop.memoryDb.Accept(lastAcceptedBlockNumber.Uint64(), lastAccepted.Time()) // will delete stale orders from the memorydb
		log.Root().SetHandler(logHandler)

		// needs to be run everytime as long as the db.UpdatePosition uses configService.GetCumulativePremiumFraction
		lop.UpdateLastPremiumFractionFromStorage()
	}

	lop.mu.Unlock()

	lop.blockBuilder = blockBuilder
	lop.runMatchingTimer()
	lop.listenAndStoreLimitOrderTransactions()
}

func (lop *limitOrderProcesser) RunMatchingPipeline() {
	executeFuncAndRecoverPanic(func() {
		matchesFound := lop.matchingPipeline.Run(new(big.Int).Add(lop.blockChain.CurrentBlock().Number, big.NewInt(1)))
		if matchesFound {
			lop.blockBuilder.signalTxsReady()
		}
	}, orderbook.RunMatchingPipelinePanicMessage, orderbook.RunMatchingPipelinePanicsCounter)
}

func (lop *limitOrderProcesser) GetOrderBookAPI() *orderbook.OrderBookAPI {
	return orderbook.NewOrderBookAPI(lop.memoryDb, lop.backend, lop.configService)
}

func (lop *limitOrderProcesser) GetTradingAPI() *orderbook.TradingAPI {
	return orderbook.NewTradingAPI(lop.memoryDb, lop.backend, lop.configService)
}

func (lop *limitOrderProcesser) listenAndStoreLimitOrderTransactions() {
	logsCh := make(chan []*types.Log)
	logsSubscription := lop.backend.SubscribeHubbleLogsEvent(logsCh)
	lop.shutdownWg.Add(1)
	go func() {
		defer lop.shutdownWg.Done()
		defer logsSubscription.Unsubscribe()
		for {
			select {
			case logs := <-logsCh:
				executeFuncAndRecoverPanic(func() {
					lop.mu.Lock()
					defer lop.mu.Unlock()
					lop.contractEventProcessor.ProcessEvents(logs)
					go lop.contractEventProcessor.PushtoTraderFeed(logs, orderbook.ConfirmationLevelHead)
				}, orderbook.HandleHubbleFeedLogsPanicMessage, orderbook.HandleHubbleFeedLogsPanicsCounter)

				lop.RunMatchingPipeline()

			case <-lop.shutdownChan:
				return
			}
		}
	}()

	acceptedLogsCh := make(chan []*types.Log)
	acceptedLogsSubscription := lop.backend.SubscribeAcceptedLogsEvent(acceptedLogsCh)
	lop.shutdownWg.Add(1)
	go func() {
		defer lop.shutdownWg.Done()
		defer acceptedLogsSubscription.Unsubscribe()

		for {
			select {
			case logs := <-acceptedLogsCh:
				executeFuncAndRecoverPanic(func() {
					lop.mu.Lock()
					defer lop.mu.Unlock()
					lop.contractEventProcessor.ProcessAcceptedEvents(logs, false)
					go lop.contractEventProcessor.PushtoTraderFeed(logs, orderbook.ConfirmationLevelAccepted)
				}, orderbook.HandleChainAcceptedLogsPanicMessage, orderbook.HandleChainAcceptedLogsPanicsCounter)
			case <-lop.shutdownChan:
				return
			}
		}
	}()

	chainAcceptedEventCh := make(chan core.ChainEvent)
	chainAcceptedEventSubscription := lop.backend.SubscribeChainAcceptedEvent(chainAcceptedEventCh)
	lop.shutdownWg.Add(1)
	go func() {
		defer lop.shutdownWg.Done()
		defer chainAcceptedEventSubscription.Unsubscribe()

		for {
			select {
			case chainAcceptedEvent := <-chainAcceptedEventCh:
				executeFuncAndRecoverPanic(func() {
					lop.handleChainAcceptedEvent(chainAcceptedEvent)
				}, orderbook.HandleChainAcceptedEventPanicMessage, orderbook.HandleChainAcceptedEventPanicsCounter)
			case <-lop.shutdownChan:
				return
			}
		}
	}()
}

// executes the matching pipeline periodically
func (lop *limitOrderProcesser) runMatchingTimer() {
	lop.shutdownWg.Add(1)
	go executeFuncAndRecoverPanic(func() {
		defer lop.shutdownWg.Done()

		for {
			select {
			case <-lop.matchingPipeline.MatchingTicker.C:
				lop.RunMatchingPipeline()

			case <-lop.shutdownChan:
				lop.matchingPipeline.MatchingTicker.Stop()
				return
			}
		}
	}, orderbook.RunMatchingPipelinePanicMessage, orderbook.RunMatchingPipelinePanicsCounter)
}

func (lop *limitOrderProcesser) handleChainAcceptedEvent(event core.ChainEvent) {
	lop.mu.Lock()
	defer lop.mu.Unlock()
	block := event.Block
	log.Info("received ChainAcceptedEvent", "number", block.NumberU64(), "hash", block.Hash().String())
	lop.memoryDb.Accept(block.NumberU64(), block.Time())

	// update metrics asynchronously
	go lop.limitOrderTxProcessor.UpdateMetrics(block)
	if block.NumberU64()%snapshotInterval == 0 {
		err := lop.saveMemoryDBSnapshot(block.Number())
		if err != nil {
			log.Error("Error in saving memory DB snapshot", "err", err)
		}
	}
}

func (lop *limitOrderProcesser) loadMemoryDBSnapshot() (acceptedBlockNumber uint64, err error) {
	snapshotFound, err := lop.hubbleDB.Has([]byte(memoryDBSnapshotKey))
	if err != nil {
		return acceptedBlockNumber, fmt.Errorf("Error in checking snapshot in hubbleDB: err=%v", err)
	}

	if !snapshotFound {
		return acceptedBlockNumber, nil
	}

	memorySnapshotBytes, err := lop.hubbleDB.Get([]byte(memoryDBSnapshotKey))
	if err != nil {
		return acceptedBlockNumber, fmt.Errorf("Error in fetching snapshot from hubbleDB; err=%v", err)
	}

	buf := bytes.NewBuffer(memorySnapshotBytes)
	var snapshot orderbook.Snapshot
	err = gob.NewDecoder(buf).Decode(&snapshot)
	if err != nil {
		return acceptedBlockNumber, fmt.Errorf("Error in snapshot parsing; err=%v", err)
	}

	if snapshot.AcceptedBlockNumber != nil && snapshot.AcceptedBlockNumber.Uint64() > 0 {
		err = lop.memoryDb.LoadFromSnapshot(snapshot)
		if err != nil {
			return acceptedBlockNumber, fmt.Errorf("Error in loading from snapshot: err=%v", err)
		}

		return snapshot.AcceptedBlockNumber.Uint64(), nil
	} else {
		return acceptedBlockNumber, nil
	}
}

// assumes that memory DB lock is held
func (lop *limitOrderProcesser) saveMemoryDBSnapshot(acceptedBlockNumber *big.Int) error {
	currentHeadBlock := lop.blockChain.CurrentBlock()

	memoryDBCopy, err := lop.memoryDb.GetOrderBookDataCopy()
	if err != nil {
		return fmt.Errorf("Error in getting memory DB copy: err=%v", err)
	}
	if currentHeadBlock.Number.Cmp(acceptedBlockNumber) == 1 {
		// if current head is ahead of the accepted block, then certain events(OrderBook)
		// need to be removed from the saved state
		logsToRemove := []*types.Log{}
		for {
			logs := lop.blockChain.GetLogs(currentHeadBlock.Hash(), currentHeadBlock.Number.Uint64())
			flattenedLogs := types.FlattenLogs(logs)
			logsToRemove = append(logsToRemove, flattenedLogs...)

			currentHeadBlock = lop.blockChain.GetHeaderByHash(currentHeadBlock.ParentHash)
			if currentHeadBlock.Number.Cmp(acceptedBlockNumber) == 0 {
				break
			}
		}

		for i := 0; i < len(logsToRemove); i++ {
			logsToRemove[i].Removed = true
		}

		cev := orderbook.NewContractEventsProcessor(memoryDBCopy)
		cev.ProcessEvents(logsToRemove)
	}

	snapshot := orderbook.Snapshot{
		Data:                memoryDBCopy,
		AcceptedBlockNumber: acceptedBlockNumber,
	}

	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(&snapshot)
	if err != nil {
		return fmt.Errorf("error in gob encoding: err=%v", err)
	}

	err = lop.hubbleDB.Put([]byte(memoryDBSnapshotKey), buf.Bytes())
	if err != nil {
		return fmt.Errorf("Error in saving to DB: err=%v", err)
	}

	log.Info("Saved memory DB snapshot successfully", "accepted block", acceptedBlockNumber, "head block number", currentHeadBlock.Number, "head block hash", currentHeadBlock.Hash())

	return nil
}

func (lop *limitOrderProcesser) getLogs(fromBlock, toBlock *big.Int) []*types.Log {
	ctx := context.Background()
	logs, err := lop.filterAPI.GetLogs(ctx, filters.FilterCriteria{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: []common.Address{orderbook.OrderBookContractAddress, orderbook.ClearingHouseContractAddress, orderbook.MarginAccountContractAddress},
	})

	if err != nil {
		log.Error("ListenAndProcessTransactions - GetLogs failed", "err", err)
		panic(err)
	}
	return logs
}

func (lop *limitOrderProcesser) UpdateLastPremiumFractionFromStorage() {
	// This is to fix the bug that was causing the LastPremiumFraction to be set to 0 in the snapshot whenever a trader's position was updated
	traderMap := lop.memoryDb.GetOrderBookData().TraderMap
	count := 0
	start := time.Now()
	for traderAddr, trader := range traderMap {
		for market := range trader.Positions {
			lastPremiumFraction := lop.configService.GetLastPremiumFraction(market, &traderAddr)
			cumulativePremiumFraction := lop.configService.GetCumulativePremiumFraction(market)
			lop.memoryDb.UpdateLastPremiumFraction(market, traderAddr, lastPremiumFraction, cumulativePremiumFraction)
			count++
		}
	}
	log.Info("@@@@ UpdateLastPremiumFractionFromStorage - update complete", "count", count, "time taken", time.Since(start))
}

func executeFuncAndRecoverPanic(fn func(), panicMessage string, panicCounter metrics.Counter) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			var errorMessage string
			switch panicInfo := panicInfo.(type) {
			case string:
				errorMessage = fmt.Sprintf("recovered (string) panic: %s", panicInfo)
			case runtime.Error:
				errorMessage = fmt.Sprintf("recovered (runtime.Error) panic: %s", panicInfo.Error())
			case error:
				errorMessage = fmt.Sprintf("recovered (error) panic: %s", panicInfo.Error())
			default:
				errorMessage = fmt.Sprintf("recovered (default) panic: %v", panicInfo)
			}

			log.Error(panicMessage, "errorMessage", errorMessage, "stack_trace", string(debug.Stack()))
			panicCounter.Inc(1)
		}
	}()
	fn()
}
