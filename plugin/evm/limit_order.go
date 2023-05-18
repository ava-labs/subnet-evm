package evm

import (
	"context"
	"math/big"
	"sync"

	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/eth"
	"github.com/ava-labs/subnet-evm/eth/filters"
	"github.com/ava-labs/subnet-evm/plugin/evm/limitorders"
	"github.com/ava-labs/subnet-evm/utils"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type LimitOrderProcesser interface {
	ListenAndProcessTransactions()
	RunBuildBlockPipeline()
	GetOrderBookAPI() *limitorders.OrderBookAPI
}

type limitOrderProcesser struct {
	ctx                    *snow.Context
	txPool                 *core.TxPool
	shutdownChan           <-chan struct{}
	shutdownWg             *sync.WaitGroup
	backend                *eth.EthAPIBackend
	blockChain             *core.BlockChain
	memoryDb               limitorders.LimitOrderDatabase
	limitOrderTxProcessor  limitorders.LimitOrderTxProcessor
	contractEventProcessor *limitorders.ContractEventsProcessor
	buildBlockPipeline     *limitorders.BuildBlockPipeline
	filterAPI              *filters.FilterAPI
}

func NewLimitOrderProcesser(ctx *snow.Context, txPool *core.TxPool, shutdownChan <-chan struct{}, shutdownWg *sync.WaitGroup, backend *eth.EthAPIBackend, blockChain *core.BlockChain) LimitOrderProcesser {
	log.Info("**** NewLimitOrderProcesser")
	configService := limitorders.NewConfigService(blockChain)
	memoryDb := limitorders.NewInMemoryDatabase(configService)
	lotp := limitorders.NewLimitOrderTxProcessor(txPool, memoryDb, backend)
	contractEventProcessor := limitorders.NewContractEventsProcessor(memoryDb)
	buildBlockPipeline := limitorders.NewBuildBlockPipeline(memoryDb, lotp, configService)
	filterSystem := filters.NewFilterSystem(backend, filters.Config{})
	filterAPI := filters.NewFilterAPI(filterSystem, true)
	return &limitOrderProcesser{
		ctx:                    ctx,
		txPool:                 txPool,
		shutdownChan:           shutdownChan,
		shutdownWg:             shutdownWg,
		backend:                backend,
		memoryDb:               memoryDb,
		blockChain:             blockChain,
		limitOrderTxProcessor:  lotp,
		contractEventProcessor: contractEventProcessor,
		buildBlockPipeline:     buildBlockPipeline,
		filterAPI:              filterAPI,
	}
}

func (lop *limitOrderProcesser) ListenAndProcessTransactions() {
	lastAccepted := lop.blockChain.LastAcceptedBlock().Number()
	if lastAccepted.Sign() > 0 {
		log.Info("ListenAndProcessTransactions - beginning sync", " till block number", lastAccepted)
		ctx := context.Background()

		var fromBlock, toBlock *big.Int
		fromBlock = big.NewInt(0)
		toBlock = utils.BigIntMin(lastAccepted, big.NewInt(0).Add(fromBlock, big.NewInt(10000)))
		for toBlock.Cmp(fromBlock) >= 0 {
			logs, err := lop.filterAPI.GetLogs(ctx, filters.FilterCriteria{
				FromBlock: fromBlock,
				ToBlock:   toBlock, // check that this is inclusive...
				Addresses: []common.Address{limitorders.OrderBookContractAddress, limitorders.ClearingHouseContractAddress, limitorders.MarginAccountContractAddress},
			})
			log.Info("ListenAndProcessTransactions", "fromBlock", fromBlock.String(), "toBlock", toBlock.String(), "number of logs", len(logs), "err", err)
			if err != nil {
				log.Error("ListenAndProcessTransactions - GetLogs failed", "err", err)
				panic(err)
			}
			lop.contractEventProcessor.ProcessEvents(logs)
			lop.contractEventProcessor.ProcessAcceptedEvents(logs)

			fromBlock = fromBlock.Add(toBlock, big.NewInt(1))
			toBlock = utils.BigIntMin(lastAccepted, big.NewInt(0).Add(fromBlock, big.NewInt(10000)))
		}
		lop.memoryDb.Accept(lastAccepted.Uint64()) // will delete stale orders from the memorydb
	}

	lop.listenAndStoreLimitOrderTransactions()
}

func (lop *limitOrderProcesser) RunBuildBlockPipeline() {
	lop.buildBlockPipeline.Run()
}

func (lop *limitOrderProcesser) GetOrderBookAPI() *limitorders.OrderBookAPI {
	return limitorders.NewOrderBookAPI(lop.memoryDb, lop.backend)
}

func (lop *limitOrderProcesser) listenAndStoreLimitOrderTransactions() {
	logsCh := make(chan []*types.Log)
	logsSubscription := lop.backend.SubscribeHubbleLogsEvent(logsCh)
	lop.shutdownWg.Add(1)
	go lop.ctx.Log.RecoverAndPanic(func() {
		defer lop.shutdownWg.Done()
		defer logsSubscription.Unsubscribe()

		for {
			select {
			case logs := <-logsCh:
				lop.contractEventProcessor.ProcessEvents(logs)
			case <-lop.shutdownChan:
				return
			}
		}
	})

	acceptedLogsCh := make(chan []*types.Log)
	acceptedLogsSubscription := lop.backend.SubscribeAcceptedLogsEvent(acceptedLogsCh)
	lop.shutdownWg.Add(1)
	go lop.ctx.Log.RecoverAndPanic(func() {
		defer lop.shutdownWg.Done()
		defer acceptedLogsSubscription.Unsubscribe()

		for {
			select {
			case logs := <-acceptedLogsCh:
				lop.contractEventProcessor.ProcessAcceptedEvents(logs)
			case <-lop.shutdownChan:
				return
			}
		}
	})

	chainAcceptedEventCh := make(chan core.ChainEvent)
	chainAcceptedEventSubscription := lop.backend.SubscribeChainAcceptedEvent(chainAcceptedEventCh)
	lop.shutdownWg.Add(1)
	go lop.ctx.Log.RecoverAndPanic(func() {
		defer lop.shutdownWg.Done()
		defer chainAcceptedEventSubscription.Unsubscribe()

		for {
			select {
			case chainAcceptedEvent := <-chainAcceptedEventCh:
				lop.handleChainAcceptedEvent(chainAcceptedEvent)
			case <-lop.shutdownChan:
				return
			}
		}
	})
}

func (lop *limitOrderProcesser) handleChainAcceptedEvent(event core.ChainEvent) {
	block := event.Block
	log.Info("#### received ChainAcceptedEvent", "number", block.NumberU64(), "hash", block.Hash().String())
	lop.memoryDb.Accept(block.NumberU64())
}
