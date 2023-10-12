package orderbook

import (
	"math"
	"math/big"
	"sync"
	"time"

	hu "github.com/ava-labs/subnet-evm/plugin/evm/orderbook/hubbleutils"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

const (
	// ticker frequency for calling signalTxsReady
	matchingTickerDuration = 5 * time.Second
)

type MatchingPipeline struct {
	mu             sync.Mutex
	db             LimitOrderDatabase
	lotp           LimitOrderTxProcessor
	configService  IConfigService
	MatchingTicker *time.Ticker
}

func NewMatchingPipeline(
	db LimitOrderDatabase,
	lotp LimitOrderTxProcessor,
	configService IConfigService) *MatchingPipeline {

	return &MatchingPipeline{
		db:             db,
		lotp:           lotp,
		configService:  configService,
		MatchingTicker: time.NewTicker(matchingTickerDuration),
	}
}

func (pipeline *MatchingPipeline) Run(blockNumber *big.Int) bool {
	pipeline.mu.Lock()
	defer pipeline.mu.Unlock()

	// reset ticker
	pipeline.MatchingTicker.Reset(matchingTickerDuration)
	markets := pipeline.GetActiveMarkets()
	log.Info("MatchingPipeline:Run", "blockNumber", blockNumber)

	if len(markets) == 0 {
		return false
	}

	// start fresh and purge all local transactions
	pipeline.lotp.PurgeOrderBookTxs()

	if isFundingPaymentTime(pipeline.db.GetNextFundingTime()) {
		log.Info("MatchingPipeline:isFundingPaymentTime")
		err := executeFundingPayment(pipeline.lotp)
		if err != nil {
			log.Error("Funding payment job failed", "err", err)
		}
	}

	// check nextSamplePITime
	if isSamplePITime(pipeline.db.GetNextSamplePITime(), pipeline.db.GetSamplePIAttemptedTime()) {
		log.Info("MatchingPipeline:isSamplePITime")
		err := pipeline.lotp.ExecuteSamplePITx()
		if err != nil {
			log.Error("Sample PI job failed", "err", err)
		}
	}

	// fetch the underlying price and run the matching engine
	hState := &hu.HubbleState{
		Assets:             pipeline.GetCollaterals(),
		OraclePrices:       pipeline.GetUnderlyingPrices(),
		MidPrices:          pipeline.GetMidPrices(),
		ActiveMarkets:      markets,
		MinAllowableMargin: pipeline.configService.getMinAllowableMargin(),
		MaintenanceMargin:  pipeline.configService.getMaintenanceMargin(),
	}

	// build trader map
	liquidablePositions, ordersToCancel := pipeline.db.GetNaughtyTraders(hState)
	cancellableOrderIds := pipeline.cancelLimitOrders(ordersToCancel)
	orderMap := make(map[Market]*Orders)
	for _, market := range markets {
		orderMap[market] = pipeline.fetchOrders(market, hState.OraclePrices[market], cancellableOrderIds, blockNumber)
	}
	pipeline.runLiquidations(liquidablePositions, orderMap, hState.OraclePrices)
	for _, market := range markets {
		// @todo should we prioritize matching in any particular market?
		pipeline.runMatchingEngine(pipeline.lotp, orderMap[market].longOrders, orderMap[market].shortOrders)
	}

	orderBookTxsCount := pipeline.lotp.GetOrderBookTxsCount()
	if orderBookTxsCount > 0 {
		pipeline.lotp.SetOrderBookTxsBlockNumber(blockNumber.Uint64())
		return true
	}
	return false
}

type Orders struct {
	longOrders  []Order
	shortOrders []Order
}

func (pipeline *MatchingPipeline) GetActiveMarkets() []Market {
	count := pipeline.configService.GetActiveMarketsCount()
	markets := make([]Market, count)
	for i := int64(0); i < count; i++ {
		markets[i] = Market(i)
	}
	return markets
}

func (pipeline *MatchingPipeline) GetUnderlyingPrices() map[Market]*big.Int {
	prices := pipeline.configService.GetUnderlyingPrices()
	log.Info("GetUnderlyingPrices", "prices", prices)
	underlyingPrices := make(map[Market]*big.Int)
	for market, price := range prices {
		underlyingPrices[Market(market)] = price
	}
	return underlyingPrices
}

func (pipeline *MatchingPipeline) GetMidPrices() map[Market]*big.Int {
	prices := pipeline.configService.GetMidPrices()
	log.Info("GetMidPrices", "prices", prices)
	midPrices := make(map[Market]*big.Int)
	for market, price := range prices {
		midPrices[Market(market)] = price
	}
	return midPrices
}

func (pipeline *MatchingPipeline) GetCollaterals() []hu.Collateral {
	return pipeline.configService.GetCollaterals()
}

func (pipeline *MatchingPipeline) cancelLimitOrders(cancellableOrders map[common.Address][]Order) map[common.Hash]struct{} {
	cancellableOrderIds := map[common.Hash]struct{}{}
	// @todo: if there are too many cancellable orders, they might not fit in a single block. Need to adjust for that.
	for _, orders := range cancellableOrders {
		if len(orders) > 0 {
			rawOrders := make([]LimitOrder, len(orders))
			for i, order := range orders {
				rawOrder := order.RawOrder.(*LimitOrder)
				rawOrders[i] = *rawOrder // @todo: make sure only limit orders reach here
			}
			log.Info("orders to cancel", "num", len(orders))
			// cancel max of 30 orders
			err := pipeline.lotp.ExecuteLimitOrderCancel(rawOrders[0:int(math.Min(float64(len(rawOrders)), 30))]) // change this if the tx gas limit (1.5m) is changed
			if err != nil {
				log.Error("Error in ExecuteOrderCancel", "orders", orders, "err", err)
			} else {
				for _, order := range orders {
					cancellableOrderIds[order.Id] = struct{}{}
				}
			}
		}
	}
	return cancellableOrderIds
}

func (pipeline *MatchingPipeline) fetchOrders(market Market, underlyingPrice *big.Int, cancellableOrderIds map[common.Hash]struct{}, blockNumber *big.Int) *Orders {
	_, lowerBoundForLongs := pipeline.configService.GetAcceptableBounds(market)
	// any long orders below the permissible lowerbound are irrelevant, because they won't be matched no matter what.
	// this assumes that all above cancelOrder transactions got executed successfully (or atleast they are not meant to be executed anyway if they passed the cancellation criteria)
	longOrders := removeOrdersWithIds(pipeline.db.GetLongOrders(market, lowerBoundForLongs, blockNumber), cancellableOrderIds)

	// say if there were no long orders, then shord orders above liquidation upper bound are irrelevant, because they won't be matched no matter what
	// note that this assumes that permissible liquidation spread <= oracle spread
	upperBoundforShorts, _ := pipeline.configService.GetAcceptableBoundsForLiquidation(market)

	// however, if long orders exist, then
	if len(longOrders) != 0 {
		oracleUpperBound, _ := pipeline.configService.GetAcceptableBounds(market)
		// take the max of price of highest long and liq upper bound. But say longOrders[0].Price > oracleUpperBound ? - then we discard orders above oracleUpperBound, because they won't be matched no matter what
		upperBoundforShorts = utils.BigIntMin(utils.BigIntMax(longOrders[0].Price, upperBoundforShorts), oracleUpperBound)
	}
	shortOrders := removeOrdersWithIds(pipeline.db.GetShortOrders(market, upperBoundforShorts, blockNumber), cancellableOrderIds)
	return &Orders{longOrders, shortOrders}
}

func (pipeline *MatchingPipeline) runLiquidations(liquidablePositions []LiquidablePosition, orderMap map[Market]*Orders, underlyingPrices map[Market]*big.Int) {
	if len(liquidablePositions) == 0 {
		return
	}

	log.Info("found positions to liquidate", "num", len(liquidablePositions))

	// we need to retreive permissible bounds for liquidations in each market
	markets := pipeline.GetActiveMarkets()
	type S struct {
		Upperbound *big.Int
		Lowerbound *big.Int
	}
	liquidationBounds := make([]S, len(markets))
	for _, market := range markets {
		upperbound, lowerbound := pipeline.configService.GetAcceptableBoundsForLiquidation(market)
		liquidationBounds[market] = S{Upperbound: upperbound, Lowerbound: lowerbound}
	}

	for _, liquidable := range liquidablePositions {
		market := liquidable.Market
		numOrdersExhausted := 0
		switch liquidable.PositionType {
		case LONG:
			for _, order := range orderMap[market].longOrders {
				if order.Price.Cmp(liquidationBounds[market].Lowerbound) == -1 {
					// further orders are not not eligible to liquidate with
					break
				}
				fillAmount := utils.BigIntMinAbs(liquidable.GetUnfilledSize(), order.GetUnFilledBaseAssetQuantity())
				pipeline.lotp.ExecuteLiquidation(liquidable.Address, order, fillAmount)
				order.FilledBaseAssetQuantity.Add(order.FilledBaseAssetQuantity, fillAmount)
				liquidable.FilledSize.Add(liquidable.FilledSize, fillAmount)
				if order.GetUnFilledBaseAssetQuantity().Sign() == 0 {
					numOrdersExhausted++
				}
				if liquidable.GetUnfilledSize().Sign() == 0 {
					break // partial/full liquidation for this position slated for this run is complete
				}
			}
			orderMap[market].longOrders = orderMap[market].longOrders[numOrdersExhausted:]
		case SHORT:
			for _, order := range orderMap[market].shortOrders {
				if order.Price.Cmp(liquidationBounds[market].Upperbound) == 1 {
					// further orders are not not eligible to liquidate with
					break
				}
				fillAmount := utils.BigIntMinAbs(liquidable.GetUnfilledSize(), order.GetUnFilledBaseAssetQuantity())
				pipeline.lotp.ExecuteLiquidation(liquidable.Address, order, fillAmount)
				order.FilledBaseAssetQuantity.Sub(order.FilledBaseAssetQuantity, fillAmount)
				liquidable.FilledSize.Sub(liquidable.FilledSize, fillAmount)
				if order.GetUnFilledBaseAssetQuantity().Sign() == 0 {
					numOrdersExhausted++
				}
				if liquidable.GetUnfilledSize().Sign() == 0 {
					break // partial/full liquidation for this position slated for this run is complete
				}
			}
			orderMap[market].shortOrders = orderMap[market].shortOrders[numOrdersExhausted:]
		}
		if liquidable.GetUnfilledSize().Sign() != 0 {
			log.Info("unquenched liquidation", "liquidable", liquidable)
		}
	}
}

func (pipeline *MatchingPipeline) runMatchingEngine(lotp LimitOrderTxProcessor, longOrders []Order, shortOrders []Order) {
	for i := 0; i < len(longOrders); i++ {
		// if there are no short orders or if the price of the first long order is < the price of the first short order, then we can stop matching
		if len(shortOrders) == 0 || longOrders[i].Price.Cmp(shortOrders[0].Price) == -1 {
			break
		}
		numOrdersExhausted := 0
		for j := 0; j < len(shortOrders); j++ {
			fillAmount := areMatchingOrders(longOrders[i], shortOrders[j])
			if fillAmount == nil {
				continue
			}
			longOrders[i], shortOrders[j] = ExecuteMatchedOrders(lotp, longOrders[i], shortOrders[j], fillAmount)
			if shortOrders[j].GetUnFilledBaseAssetQuantity().Sign() == 0 {
				numOrdersExhausted++
			}
			if longOrders[i].GetUnFilledBaseAssetQuantity().Sign() == 0 {
				break
			}
		}
		shortOrders = shortOrders[numOrdersExhausted:]
	}
}

func areMatchingOrders(longOrder, shortOrder Order) *big.Int {
	if longOrder.Price.Cmp(shortOrder.Price) == -1 {
		return nil
	}
	blockDiff := longOrder.BlockNumber.Cmp(shortOrder.BlockNumber)
	if blockDiff == -1 && (longOrder.OrderType == IOC || shortOrder.isPostOnly()) ||
		blockDiff == 1 && (shortOrder.OrderType == IOC || longOrder.isPostOnly()) {
		return nil
	}
	fillAmount := utils.BigIntMinAbs(longOrder.GetUnFilledBaseAssetQuantity(), shortOrder.GetUnFilledBaseAssetQuantity())
	if fillAmount.Sign() == 0 {
		return nil
	}
	return fillAmount
}

func ExecuteMatchedOrders(lotp LimitOrderTxProcessor, longOrder, shortOrder Order, fillAmount *big.Int) (Order, Order) {
	lotp.ExecuteMatchedOrdersTx(longOrder, shortOrder, fillAmount)
	longOrder.FilledBaseAssetQuantity = big.NewInt(0).Add(longOrder.FilledBaseAssetQuantity, fillAmount)
	shortOrder.FilledBaseAssetQuantity = big.NewInt(0).Sub(shortOrder.FilledBaseAssetQuantity, fillAmount)
	return longOrder, shortOrder
}

func matchLongAndShortOrder(lotp LimitOrderTxProcessor, longOrder, shortOrder Order) (Order, Order, bool) {
	fillAmount := utils.BigIntMinAbs(longOrder.GetUnFilledBaseAssetQuantity(), shortOrder.GetUnFilledBaseAssetQuantity())
	if longOrder.Price.Cmp(shortOrder.Price) == -1 || fillAmount.Sign() == 0 {
		return longOrder, shortOrder, false
	}
	if longOrder.BlockNumber.Cmp(shortOrder.BlockNumber) > 0 && longOrder.isPostOnly() {
		log.Warn("post only long order matched with a resting order", "longOrder", longOrder, "shortOrder", shortOrder)
		return longOrder, shortOrder, false
	}
	if shortOrder.BlockNumber.Cmp(longOrder.BlockNumber) > 0 && shortOrder.isPostOnly() {
		log.Warn("post only short order matched with a resting order", "longOrder", longOrder, "shortOrder", shortOrder)
		return longOrder, shortOrder, false
	}
	if err := lotp.ExecuteMatchedOrdersTx(longOrder, shortOrder, fillAmount); err != nil {
		return longOrder, shortOrder, false
	}
	longOrder.FilledBaseAssetQuantity = big.NewInt(0).Add(longOrder.FilledBaseAssetQuantity, fillAmount)
	shortOrder.FilledBaseAssetQuantity = big.NewInt(0).Sub(shortOrder.FilledBaseAssetQuantity, fillAmount)
	return longOrder, shortOrder, true
}

func isFundingPaymentTime(nextFundingTime uint64) bool {
	if nextFundingTime == 0 {
		return false
	}

	now := uint64(time.Now().Unix())
	return now >= nextFundingTime
}

func isSamplePITime(nextSamplePITime, lastAttempt uint64) bool {
	if nextSamplePITime == 0 {
		return false
	}

	now := uint64(time.Now().Unix())
	return now >= nextSamplePITime && now >= lastAttempt+5 // give 5 secs for the tx to be mined
}

func executeFundingPayment(lotp LimitOrderTxProcessor) error {
	// @todo get index twap for each market with warp msging

	return lotp.ExecuteFundingPaymentTx()
}

func removeOrdersWithIds(orders []Order, orderIds map[common.Hash]struct{}) []Order {
	var filteredOrders []Order
	for _, order := range orders {
		if _, ok := orderIds[order.Id]; !ok {
			filteredOrders = append(filteredOrders, order)
		}
	}
	return filteredOrders
}
