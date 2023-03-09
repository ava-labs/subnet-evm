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
	RunBuildBlockPipeline(lastBlockTime uint64)
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
}

func NewLimitOrderProcesser(ctx *snow.Context, txPool *core.TxPool, shutdownChan <-chan struct{}, shutdownWg *sync.WaitGroup, backend *eth.EthAPIBackend, blockChain *core.BlockChain, memoryDb limitorders.LimitOrderDatabase, lotp limitorders.LimitOrderTxProcessor) LimitOrderProcesser {
	log.Info("**** NewLimitOrderProcesser")
	contractEventProcessor := limitorders.NewContractEventsProcessor(memoryDb)
	buildBlockPipeline := limitorders.NewBuildBlockPipeline(memoryDb, lotp)
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
	}
}

func (lop *limitOrderProcesser) ListenAndProcessTransactions() {
	lastAccepted := lop.blockChain.LastAcceptedBlock().Number()
	if lastAccepted.Sign() > 0 {
		log.Info("ListenAndProcessTransactions - beginning sync", " till block number", lastAccepted)
		ctx := context.Background()

		filterSystem := filters.NewFilterSystem(lop.backend, filters.Config{})
		filterAPI := filters.NewFilterAPI(filterSystem, true)

		var fromBlock, toBlock *big.Int
		fromBlock = big.NewInt(0)
		toBlock = utils.BigIntMin(lastAccepted, big.NewInt(0).Add(fromBlock, big.NewInt(10000)))
		for toBlock.Cmp(fromBlock) >= 0 {
			logs, err := filterAPI.GetLogs(ctx, filters.FilterCriteria{
				FromBlock: fromBlock,
				ToBlock:   toBlock,
				Addresses: []common.Address{limitorders.OrderBookContractAddress, limitorders.ClearingHouseContractAddress, limitorders.MarginAccountContractAddress},
			})
			if err != nil {
				log.Error("ListenAndProcessTransactions - GetLogs failed", "err", err)
				panic(err)
			}
			lop.contractEventProcessor.ProcessEvents(logs)
			log.Info("ListenAndProcessTransactions", "fromBlock", fromBlock.String(), "toBlock", toBlock.String(), "number of logs", len(logs), "err", err)

			fromBlock = toBlock.Add(fromBlock, big.NewInt(1))
			toBlock = utils.BigIntMin(lastAccepted, big.NewInt(0).Add(fromBlock, big.NewInt(10000)))
		}
	}

	lop.listenAndStoreLimitOrderTransactions()
}

func (lop *limitOrderProcesser) RunBuildBlockPipeline(lastBlockTime uint64) {
	lop.buildBlockPipeline.Run(lastBlockTime)
}

func (lop *limitOrderProcesser) GetOrderBookAPI() *limitorders.OrderBookAPI {
	return limitorders.NewOrderBookAPI(lop.memoryDb)
}

func (lop *limitOrderProcesser) listenAndStoreLimitOrderTransactions() {
	logsCh := make(chan []*types.Log)
	logsSubscription := lop.backend.SubscribeAcceptedLogsEvent(logsCh)
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
}
