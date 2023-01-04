package evm

import (
	"sync"

	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/eth"
	"github.com/ava-labs/subnet-evm/plugin/evm/limitorders"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ethereum/go-ethereum/log"
)

type LimitOrderProcesser interface {
	ListenAndProcessTransactions()
	RunMatchingEngine()
}

type limitOrderProcesser struct {
	ctx                   *snow.Context
	txPool                *core.TxPool
	shutdownChan          <-chan struct{}
	shutdownWg            *sync.WaitGroup
	backend               *eth.EthAPIBackend
	blockChain            *core.BlockChain
	memoryDb              *limitorders.InMemoryDatabase
	limitOrderTxProcessor *limitorders.LimitOrderTxProcessor
}

func NewLimitOrderProcesser(ctx *snow.Context, txPool *core.TxPool, shutdownChan <-chan struct{}, shutdownWg *sync.WaitGroup, backend *eth.EthAPIBackend, blockChain *core.BlockChain, memoryDb *limitorders.InMemoryDatabase, lotp *limitorders.LimitOrderTxProcessor) LimitOrderProcesser {
	log.Info("**** NewLimitOrderProcesser")
	return &limitOrderProcesser{
		ctx:                   ctx,
		txPool:                txPool,
		shutdownChan:          shutdownChan,
		shutdownWg:            shutdownWg,
		backend:               backend,
		memoryDb:              memoryDb,
		blockChain:            blockChain,
		limitOrderTxProcessor: lotp,
	}
}

func (lop *limitOrderProcesser) ListenAndProcessTransactions() {
	lastAccepted := lop.blockChain.LastAcceptedBlock().NumberU64()
	if lastAccepted > 0 {
		log.Info("ListenAndProcessTransactions - beginning sync", " till block number", lastAccepted)

		allTxs := types.Transactions{}
		for i := uint64(0); i <= lastAccepted; i++ {
			block := lop.blockChain.GetBlockByNumber(i)
			if block != nil {
				for _, tx := range block.Transactions() {
					if lop.limitOrderTxProcessor.CheckIfOrderBookContractCall(tx) {
						lop.limitOrderTxProcessor.HandleOrderBookTx(tx, i, *lop.backend)
					}
				}
			}
		}

		log.Info("ListenAndProcessTransactions - sync complete", "till block number", lastAccepted, "total transactions", len(allTxs))
	}

	lop.listenAndStoreLimitOrderTransactions()
}

func (lop *limitOrderProcesser) RunMatchingEngine() {
	lop.limitOrderTxProcessor.PurgeLocalTx()
	longOrders := lop.memoryDb.GetLongOrders()
	shortOrders := lop.memoryDb.GetShortOrders()
	if len(longOrders) == 0 || len(shortOrders) == 0 {
		return
	}
	for _, longOrder := range longOrders {
		for j, shortOrder := range shortOrders {
			if longOrder.Price == shortOrder.Price && longOrder.BaseAssetQuantity == (-shortOrder.BaseAssetQuantity) {
				err := lop.limitOrderTxProcessor.ExecuteMatchedOrdersTx(*longOrder, *shortOrder)
				if err == nil {
					shortOrders = append(shortOrders[:j], shortOrders[j+1:]...)
					break
				}
			}
		}
	}
}

func (lop *limitOrderProcesser) listenAndStoreLimitOrderTransactions() {
	newChainChan := make(chan core.ChainEvent)
	chainAcceptedEventSubscription := lop.backend.SubscribeChainAcceptedEvent(newChainChan)

	lop.shutdownWg.Add(1)
	go lop.ctx.Log.RecoverAndPanic(func() {
		defer lop.shutdownWg.Done()
		defer chainAcceptedEventSubscription.Unsubscribe()

		for {
			select {
			case newChainAcceptedEvent := <-newChainChan:
				tsHashes := []string{}
				blockNumber := newChainAcceptedEvent.Block.Number().Uint64()
				for _, tx := range newChainAcceptedEvent.Block.Transactions() {
					tsHashes = append(tsHashes, tx.Hash().String())
					if lop.limitOrderTxProcessor.CheckIfOrderBookContractCall(tx) {
						lop.limitOrderTxProcessor.HandleOrderBookTx(tx, blockNumber, *lop.backend)
					}
				}
				log.Info("$$$$$ New head event", "number", newChainAcceptedEvent.Block.Header().Number, "tx hashes", tsHashes,
					"miner", newChainAcceptedEvent.Block.Coinbase().String(),
					"root", newChainAcceptedEvent.Block.Header().Root.String(), "gas used", newChainAcceptedEvent.Block.Header().GasUsed,
					"nonce", newChainAcceptedEvent.Block.Header().Nonce)

			case <-lop.shutdownChan:
				return
			}
		}
	})
}
