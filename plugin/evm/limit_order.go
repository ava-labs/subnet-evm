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
	RunLiquidationsAndMatching()
	IsFundingPaymentTime(lastBlockTime uint64) bool
	ExecuteFundingPayment() error
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
}

func NewLimitOrderProcesser(ctx *snow.Context, txPool *core.TxPool, shutdownChan <-chan struct{}, shutdownWg *sync.WaitGroup, backend *eth.EthAPIBackend, blockChain *core.BlockChain, memoryDb limitorders.LimitOrderDatabase, lotp limitorders.LimitOrderTxProcessor) LimitOrderProcesser {
	log.Info("**** NewLimitOrderProcesser")
	contractEventProcessor := limitorders.NewContractEventsProcessor(memoryDb)
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
			log.Info("ListenAndProcessTransactions", "number of logs", len(logs), "err", err)

			fromBlock = toBlock.Add(fromBlock, big.NewInt(1))
			toBlock = utils.BigIntMin(lastAccepted, big.NewInt(0).Add(fromBlock, big.NewInt(10000)))
		}
	}

	lop.listenAndStoreLimitOrderTransactions()
}

func (lop *limitOrderProcesser) IsFundingPaymentTime(lastBlockTime uint64) bool {
	if lop.memoryDb.GetNextFundingTime() == 0 {
		return false
	}
	return lastBlockTime >= lop.memoryDb.GetNextFundingTime()
}

func (lop *limitOrderProcesser) ExecuteFundingPayment() error {
	// @todo get index twap for each market with warp msging

	return lop.limitOrderTxProcessor.ExecuteFundingPaymentTx()
}

func (lop *limitOrderProcesser) RunLiquidationsAndMatching() {
	lop.limitOrderTxProcessor.PurgeLocalTx()
	for _, market := range limitorders.GetActiveMarkets() {
		longOrders := lop.memoryDb.GetLongOrders(market)
		shortOrders := lop.memoryDb.GetShortOrders(market)
		longOrders, shortOrders = lop.runLiquidations(market, longOrders, shortOrders)
		lop.runMatchingEngine(longOrders, shortOrders)
	}
}

func (lop *limitOrderProcesser) runMatchingEngine(longOrders []limitorders.LimitOrder, shortOrders []limitorders.LimitOrder) {

	if len(longOrders) == 0 || len(shortOrders) == 0 {
		return
	}
	for i := 0; i < len(longOrders); i++ {
		for j := 0; j < len(shortOrders); j++ {
			if longOrders[i].GetUnFilledBaseAssetQuantity().Sign() == 0 {
				break
			}
			if shortOrders[j].GetUnFilledBaseAssetQuantity().Sign() == 0 {
				continue
			}
			var ordersMatched bool
			longOrders[i], shortOrders[j], ordersMatched = matchLongAndShortOrder(lop.limitOrderTxProcessor, longOrders[i], shortOrders[j])
			if !ordersMatched {
				i = len(longOrders)
				break
			}
		}
	}
}

func (lop *limitOrderProcesser) runLiquidations(market limitorders.Market, longOrders []limitorders.LimitOrder, shortOrders []limitorders.LimitOrder) (filteredLongOrder []limitorders.LimitOrder, filteredShortOrder []limitorders.LimitOrder) {
	oraclePrice := big.NewInt(20 * 10e6) // @todo: get it from the oracle

	liquidablePositions := limitorders.GetLiquidableTraders(lop.memoryDb.GetAllTraders(), market, lop.memoryDb.GetLastPrice(market), oraclePrice)

	for i, liquidable := range liquidablePositions {
		var oppositeOrders []limitorders.LimitOrder
		switch liquidable.PositionType {
		case "long":
			oppositeOrders = shortOrders
		case "short":
			oppositeOrders = longOrders
		}
		if len(oppositeOrders) == 0 {
			log.Error("no matching order found for liquidation", "trader", liquidable.Address.String(), "size", liquidable.Size)
			continue // so that all other liquidable positions get logged
		}
		for j, oppositeOrder := range oppositeOrders {
			if liquidable.GetUnfilledSize().Sign() == 0 {
				break
			}
			// @todo: add a restriction on the price range that liquidation will occur on.
			// An illiquid market can be very adverse for trader being liquidated.
			fillAmount := utils.BigIntMinAbs(liquidable.GetUnfilledSize(), oppositeOrder.GetUnFilledBaseAssetQuantity())
			if fillAmount.Sign() == 0 {
				continue
			}
			lop.limitOrderTxProcessor.ExecuteLiquidation(liquidable.Address, oppositeOrder, fillAmount)

			switch liquidable.PositionType {
			case "long":
				oppositeOrders[j].FilledBaseAssetQuantity.Sub(oppositeOrders[j].FilledBaseAssetQuantity, fillAmount)
				liquidablePositions[i].FilledSize.Add(liquidablePositions[i].FilledSize, fillAmount)
			case "short":
				oppositeOrders[j].FilledBaseAssetQuantity.Add(oppositeOrders[j].FilledBaseAssetQuantity, fillAmount)
				liquidablePositions[i].FilledSize.Sub(liquidablePositions[i].FilledSize, fillAmount)
			}
		}
	}

	return longOrders, shortOrders
}

func matchLongAndShortOrder(lotp limitorders.LimitOrderTxProcessor, longOrder limitorders.LimitOrder, shortOrder limitorders.LimitOrder) (limitorders.LimitOrder, limitorders.LimitOrder, bool) {
	if longOrder.Price.Cmp(shortOrder.Price) >= 0 { // longOrder.Price >= shortOrder.Price
		fillAmount := utils.BigIntMinAbs(longOrder.GetUnFilledBaseAssetQuantity(), shortOrder.GetUnFilledBaseAssetQuantity())
		if fillAmount.Sign() != 0 {
			err := lotp.ExecuteMatchedOrdersTx(longOrder, shortOrder, fillAmount)
			if err == nil {
				longOrder.FilledBaseAssetQuantity = big.NewInt(0).Add(longOrder.FilledBaseAssetQuantity, fillAmount)
				shortOrder.FilledBaseAssetQuantity = big.NewInt(0).Sub(shortOrder.FilledBaseAssetQuantity, fillAmount)
				return longOrder, shortOrder, true
			}
		}
	}
	return longOrder, shortOrder, false
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

func (lop *limitOrderProcesser) GetOrderBookAPI() *limitorders.OrderBookAPI {
	return limitorders.NewOrderBookAPI(lop.memoryDb)
}
