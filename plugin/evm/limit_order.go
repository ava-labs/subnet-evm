package evm

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"math/big"
	"os"
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
	hu "github.com/ava-labs/subnet-evm/plugin/evm/orderbook/hubbleutils"
	"github.com/ava-labs/subnet-evm/utils"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ethereum/go-ethereum/log"
)

const (
	snapshotInterval uint64 = 10 // save snapshot every 1000 blocks
)

type LimitOrderProcesser interface {
	ListenAndProcessTransactions(blockBuilder *blockBuilder)
	GetOrderBookAPI() *orderbook.OrderBookAPI
	GetTestingAPI() *orderbook.TestingAPI
	GetTradingAPI() *orderbook.TradingAPI
}

type limitOrderProcesser struct {
	ctx                      *snow.Context
	mu                       *sync.Mutex
	txPool                   *txpool.TxPool
	shutdownChan             <-chan struct{}
	shutdownWg               *sync.WaitGroup
	backend                  *eth.EthAPIBackend
	blockChain               *core.BlockChain
	memoryDb                 orderbook.LimitOrderDatabase
	limitOrderTxProcessor    orderbook.LimitOrderTxProcessor
	contractEventProcessor   *orderbook.ContractEventsProcessor
	matchingPipeline         *orderbook.MatchingPipeline
	filterAPI                *filters.FilterAPI
	configService            orderbook.IConfigService
	blockBuilder             *blockBuilder
	isValidator              bool
	tradingAPIEnabled        bool
	loadFromSnapshotEnabled  bool
	snapshotSavedBlockNumber uint64
	snapshotFilePath         string
	tradingAPI               *orderbook.TradingAPI
}

func NewLimitOrderProcesser(ctx *snow.Context, txPool *txpool.TxPool, shutdownChan <-chan struct{}, shutdownWg *sync.WaitGroup, backend *eth.EthAPIBackend, blockChain *core.BlockChain, validatorPrivateKey string, config Config) LimitOrderProcesser {
	log.Info("**** NewLimitOrderProcesser")
	configService := orderbook.NewConfigService(blockChain)
	memoryDb := orderbook.NewInMemoryDatabase(configService)
	lotp := orderbook.NewLimitOrderTxProcessor(txPool, memoryDb, backend, validatorPrivateKey)
	signedObAddy := configService.GetSignedOrderbookContract()
	contractEventProcessor := orderbook.NewContractEventsProcessor(memoryDb, signedObAddy)

	matchingPipeline := orderbook.NewMatchingPipeline(memoryDb, lotp, configService)
	// if any of the following values are changed, the nodes will need to be restarted
	hState := &hu.HubbleState{
		Assets:             matchingPipeline.GetCollaterals(),
		ActiveMarkets:      matchingPipeline.GetActiveMarkets(),
		MinAllowableMargin: configService.GetMinAllowableMargin(),
		MaintenanceMargin:  configService.GetMaintenanceMargin(),
		TakerFee:           configService.GetTakerFee(),
		UpgradeVersion:     configService.GetUpgradeVersion(),
	}
	hu.SetHubbleState(hState)
	hu.SetChainIdAndVerifyingSignedOrdersContract(backend.ChainConfig().ChainID.Int64(), signedObAddy.String())

	filterSystem := filters.NewFilterSystem(backend, filters.Config{})
	filterAPI := filters.NewFilterAPI(filterSystem)

	// need to register the types for gob encoding because memory DB has an interface field(ContractOrder)
	// naming hu.LimitOrder as orderbook.LimitOrder instead of hubbleutils.LimitOrder because of backward compatibility. DO NOT CHANGE
	// same for hu.IOCOrder
	gob.RegisterName("*orderbook.LimitOrder", &hu.LimitOrder{})
	gob.RegisterName("*orderbook.IOCOrder", &hu.IOCOrder{})
	gob.Register(&hu.SignedOrder{})
	return &limitOrderProcesser{
		ctx:                     ctx,
		mu:                      &sync.Mutex{},
		txPool:                  txPool,
		shutdownChan:            shutdownChan,
		shutdownWg:              shutdownWg,
		backend:                 backend,
		memoryDb:                memoryDb,
		blockChain:              blockChain,
		limitOrderTxProcessor:   lotp,
		contractEventProcessor:  contractEventProcessor,
		matchingPipeline:        matchingPipeline,
		filterAPI:               filterAPI,
		configService:           configService,
		isValidator:             config.IsValidator,
		tradingAPIEnabled:       config.TradingAPIEnabled,
		loadFromSnapshotEnabled: config.LoadFromSnapshotEnabled,
		snapshotFilePath:        config.SnapshotFilePath,
	}
}

func (lop *limitOrderProcesser) ListenAndProcessTransactions(blockBuilder *blockBuilder) {
	lop.mu.Lock()

	lastAccepted := lop.blockChain.LastAcceptedBlock()
	lastAcceptedBlockNumber := lastAccepted.Number()
	if lastAcceptedBlockNumber.Sign() > 0 {
		fromBlock := big.NewInt(0)

		if lop.loadFromSnapshotEnabled {
			// first load the last snapshot containing finalised data till block x and query the logs of [x+1, latest]
			acceptedBlockNumber, err := lop.loadMemoryDBSnapshotFromFile()
			if err != nil {
				log.Error("ListenAndProcessTransactions - error in loading snapshot", "err", err)
			} else {
				if acceptedBlockNumber > 0 {
					fromBlock = big.NewInt(int64(acceptedBlockNumber) + 1)
				} else {
					// not an error, but unlikely after the blockchain is running for some time
					log.Warn("ListenAndProcessTransactions - no snapshot found")
				}
			}
		} else {
			log.Info("ListenAndProcessTransactions - loading from snapshot is disabled")
		}

		logHandler := log.Root().GetHandler()
		errorOnlyHandler := ErrorOnlyHandler(logHandler)
		log.Info("ListenAndProcessTransactions - beginning sync", " till block number", lastAcceptedBlockNumber)
		JUMP := big.NewInt(3999)
		toBlock := utils.BigIntMin(lastAcceptedBlockNumber, big.NewInt(0).Add(fromBlock, JUMP))
		for toBlock.Cmp(fromBlock) > 0 {
			logs := lop.getLogs(fromBlock, toBlock)
			// set the log handler to discard logs so that the ProcessEvents doesn't spam the logs
			log.Root().SetHandler(errorOnlyHandler)
			lop.contractEventProcessor.ProcessEvents(logs)
			lop.contractEventProcessor.ProcessAcceptedEvents(logs, true)
			lop.memoryDb.Accept(toBlock.Uint64(), 0) // will delete stale orders from the memorydb
			log.Root().SetHandler(logHandler)
			log.Info("ListenAndProcessTransactions - processed log chunk", "fromBlock", fromBlock.String(), "toBlock", toBlock.String(), "number of logs", len(logs))

			fromBlock = fromBlock.Add(toBlock, big.NewInt(1))
			toBlock = utils.BigIntMin(lastAcceptedBlockNumber, big.NewInt(0).Add(fromBlock, JUMP))
		}
		lop.memoryDb.Accept(lastAcceptedBlockNumber.Uint64(), lastAccepted.Time()) // will delete stale orders from the memorydb
		lop.snapshotSavedBlockNumber = lastAcceptedBlockNumber.Uint64()
		log.Info("Set snapshotSavedBlockNumber", "snapshotSavedBlockNumber", lop.snapshotSavedBlockNumber)
		log.Root().SetHandler(logHandler)
	}

	lop.mu.Unlock()

	lop.blockBuilder = blockBuilder
	lop.runMatchingTimer()
	lop.listenAndStoreLimitOrderTransactions()
}

func (lop *limitOrderProcesser) RunMatchingPipeline() {
	if !lop.isValidator {
		return
	}
	executeFuncAndRecoverPanic(func() {
		matchesFound := lop.matchingPipeline.Run(new(big.Int).Add(lop.blockChain.CurrentBlock().Number, big.NewInt(1)))
		if matchesFound {
			lop.blockBuilder.signalTxsReady()
		}
	}, orderbook.RunMatchingPipelinePanicMessage, orderbook.RunMatchingPipelinePanicsCounter)
}

func (lop *limitOrderProcesser) RunSanitaryPipeline() {
	executeFuncAndRecoverPanic(func() {
		lop.matchingPipeline.RunSanitization()
	}, orderbook.RunSanitaryPipelinePanicMessage, orderbook.RunSanitaryPipelinePanicsCounter)
}

func (lop *limitOrderProcesser) GetOrderBookAPI() *orderbook.OrderBookAPI {
	return orderbook.NewOrderBookAPI(lop.memoryDb, lop.backend, lop.configService)
}

func (lop *limitOrderProcesser) GetTradingAPI() *orderbook.TradingAPI {
	if lop.tradingAPI == nil {
		lop.tradingAPI = orderbook.NewTradingAPI(lop.memoryDb, lop.backend, lop.configService)
	}
	return lop.tradingAPI
}

func (lop *limitOrderProcesser) GetTestingAPI() *orderbook.TestingAPI {
	return orderbook.NewTestingAPI(lop.memoryDb, lop.backend, lop.configService, lop.snapshotFilePath)
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
					if lop.tradingAPIEnabled {
						go lop.contractEventProcessor.PushToTraderFeed(logs, orderbook.ConfirmationLevelHead)
						go lop.contractEventProcessor.PushToMarketFeed(logs, orderbook.ConfirmationLevelHead)
					}
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

					if len(logs) == 0 {
						return
					}
					if lop.tradingAPIEnabled {
						go lop.contractEventProcessor.PushToTraderFeed(logs, orderbook.ConfirmationLevelAccepted)
						go lop.contractEventProcessor.PushToMarketFeed(logs, orderbook.ConfirmationLevelAccepted)
					}

					blockNumber := logs[0].BlockNumber
					block := lop.blockChain.GetBlockByHash(logs[0].BlockHash)

					// If n is the block at which snapshot should be saved(n is multiple of [snapshotInterval]), save the snapshot
					// when logs of block number >= n + 1 are received before applying them in memory db

					// snapshot should be saved at block number = blockNumber - 1 because Accepted logs
					// have been applied in memory DB at this point
					snapshotBlockNumber := blockNumber - 1
					blockNumberFloor := ((snapshotBlockNumber) / snapshotInterval) * snapshotInterval
					if blockNumberFloor > lop.snapshotSavedBlockNumber {
						log.Info("Saving memory DB snapshot", "snapshotBlockNumber", snapshotBlockNumber, "current blockNumber", blockNumber, "blockNumberFloor", blockNumberFloor)
						snapshotBlock := lop.blockChain.GetBlockByNumber(snapshotBlockNumber)
						lop.memoryDb.Accept(snapshotBlockNumber, snapshotBlock.Timestamp())
						err := lop.saveMemoryDBSnapshotToFile(big.NewInt(int64(snapshotBlockNumber)))
						if err != nil {
							orderbook.SnapshotWriteFailuresCounter.Inc(1)
							log.Error("Error in saving memory DB snapshot", "err", err, "snapshotBlockNumber", snapshotBlockNumber, "current blockNumber", blockNumber, "blockNumberFloor", blockNumberFloor)
						}
					}

					lop.contractEventProcessor.ProcessAcceptedEvents(logs, false)
					lop.memoryDb.Accept(blockNumber, block.Timestamp())
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
					lop.mu.Lock()
					defer lop.mu.Unlock()
					block := chainAcceptedEvent.Block
					log.Info("received ChainAcceptedEvent", "number", block.NumberU64(), "hash", block.Hash().String())

					// update metrics asynchronously
					go lop.limitOrderTxProcessor.UpdateMetrics(block)

				}, orderbook.HandleChainAcceptedEventPanicMessage, orderbook.HandleChainAcceptedEventPanicsCounter)
			case <-lop.shutdownChan:
				return
			}
		}
	}()

	chainHeadEventCh := make(chan core.ChainHeadEvent)
	chainHeadEventSubscription := lop.backend.SubscribeChainHeadEvent(chainHeadEventCh)
	lop.shutdownWg.Add(1)
	go func() {
		defer lop.shutdownWg.Done()
		defer chainHeadEventSubscription.Unsubscribe()

		for {
			select {
			case chainHeadEvent := <-chainHeadEventCh:
				block := chainHeadEvent.Block
				log.Info("received ChainHeadEvent", "number", block.NumberU64(), "hash", block.Hash().String())
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

			case <-lop.matchingPipeline.SanitaryTicker.C:
				lop.RunSanitaryPipeline()

			case <-lop.shutdownChan:
				lop.matchingPipeline.MatchingTicker.Stop()
				return
			}
		}
	}, orderbook.RunMatchingPipelinePanicMessage, orderbook.RunMatchingPipelinePanicsCounter)
}

func (lop *limitOrderProcesser) loadMemoryDBSnapshotFromFile() (acceptedBlockNumber uint64, err error) {
	if lop.snapshotFilePath == "" {
		return acceptedBlockNumber, fmt.Errorf("snapshot file path not set")
	}

	memorySnapshotBytes, err := os.ReadFile(lop.snapshotFilePath)
	if err != nil {
		return acceptedBlockNumber, fmt.Errorf("Error in reading snapshot file: err=%v", err)
	}

	buf := bytes.NewBuffer(memorySnapshotBytes)
	var snapshot orderbook.Snapshot
	err = gob.NewDecoder(buf).Decode(&snapshot)
	if err != nil {
		return acceptedBlockNumber, fmt.Errorf("Error in snapshot parsing from file; err=%v", err)
	}

	if snapshot.AcceptedBlockNumber != nil && snapshot.AcceptedBlockNumber.Uint64() > 0 {
		err = lop.memoryDb.LoadFromSnapshot(snapshot)
		if err != nil {
			return acceptedBlockNumber, fmt.Errorf("Error in loading snapshot from file: err=%v", err)
		} else {
			log.Info("memory DB snapshot loaded from file", "acceptedBlockNumber", snapshot.AcceptedBlockNumber)
		}

		return snapshot.AcceptedBlockNumber.Uint64(), nil
	} else {
		return acceptedBlockNumber, nil
	}
}

// assumes that memory DB lock is held
func (lop *limitOrderProcesser) saveMemoryDBSnapshotToFile(acceptedBlockNumber *big.Int) error {
	if lop.snapshotFilePath == "" {
		return fmt.Errorf("snapshot file path not set")
	}
	start := time.Now()
	currentHeadBlock := lop.blockChain.CurrentBlock()

	memoryDBCopy, err := lop.memoryDb.GetOrderBookDataCopy()
	if err != nil {
		return fmt.Errorf("error in getting memory DB copy: err=%v", err)
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

		cev := orderbook.NewContractEventsProcessor(memoryDBCopy, lop.configService.GetSignedOrderbookContract())
		cev.ProcessEvents(logsToRemove)
	}

	// these SHOULD be re-populated while loading from snapshot
	memoryDBCopy.LongOrders = nil
	memoryDBCopy.ShortOrders = nil
	snapshot := orderbook.Snapshot{
		Data:                memoryDBCopy,
		AcceptedBlockNumber: acceptedBlockNumber,
	}

	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(&snapshot)
	if err != nil {
		return fmt.Errorf("error in gob encoding: err=%v", err)
	}

	snapshotBytes := buf.Bytes()

	// write to snapshot file
	err = os.WriteFile(lop.snapshotFilePath, snapshotBytes, 0644)
	if err != nil {
		return fmt.Errorf("Error in writing to snapshot file: err=%v", err)
	}

	lop.snapshotSavedBlockNumber = acceptedBlockNumber.Uint64()
	log.Info("Saved memory DB snapshot successfully", "accepted block", acceptedBlockNumber, "head block number", currentHeadBlock.Number, "head block hash", currentHeadBlock.Hash(), "duration", time.Since(start))

	return nil
}

func (lop *limitOrderProcesser) getLogs(fromBlock, toBlock *big.Int) []*types.Log {
	ctx := context.Background()
	logs, err := lop.filterAPI.GetLogs(ctx, filters.FilterCriteria{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
	})

	if err != nil {
		log.Error("ListenAndProcessTransactions - GetLogs failed", "err", err)
		panic(err)
	}
	return logs
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
			orderbook.AllPanicsCounter.Inc(1)
		}
	}()
	fn()
}
