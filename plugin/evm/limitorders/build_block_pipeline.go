package limitorders

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type BuildBlockPipeline struct {
	db   LimitOrderDatabase
	lotp LimitOrderTxProcessor
}

func NewBuildBlockPipeline(db LimitOrderDatabase, lotp LimitOrderTxProcessor) *BuildBlockPipeline {
	return &BuildBlockPipeline{
		db:   db,
		lotp: lotp,
	}
}

func (pipeline *BuildBlockPipeline) Run(lastBlockTime uint64) {
	pipeline.lotp.PurgeLocalTx()
	if isFundingPaymentTime(lastBlockTime, pipeline.db) {
		log.Info("BuildBlockPipeline:isFundingPaymentTime")
		// just execute the funding payment and skip running the matching engine
		err := executeFundingPayment(pipeline.lotp)
		if err != nil {
			log.Error("Funding payment job failed", "err", err)
		}
	} else {
		// fetch the underlying price and run the matching engine
		underlyingPrices, err := pipeline.lotp.GetUnderlyingPrice()
		if err != nil {
			log.Error("could not fetch underlying price", "err", err)
		} else {
			cancellableOrderIds := pipeline.cancelOrders(underlyingPrices)
			for _, market := range GetActiveMarkets() {
				pipeline.runLiquidationsAndMatchingForMarket(market, underlyingPrices[market], cancellableOrderIds)
			}
		}
	}
}

func (pipeline *BuildBlockPipeline) runLiquidationsAndMatchingForMarket(market Market, underlyingPrice *big.Int, cancellableOrderIds map[common.Hash]struct{}) {
	log.Info("BuildBlockPipeline:runLiquidationsAndMatchingForMarket", "underlyingPrice", prettifyScaledBigInt(underlyingPrice, 6))

	// 1. Get long orders
	longCutOffPrice := divideByBasePrecision(big.NewInt(0).Mul(
		underlyingPrice,
		big.NewInt(0).Add(_1e6, spreadRatioThreshold),
	),
	)
	longOrders := pipeline.db.GetLongOrders(market, longCutOffPrice)

	// 2. Get short orders
	shortCutOffPrice := big.NewInt(0)
	if _1e6.Cmp(spreadRatioThreshold) > 0 {
		shortCutOffPrice = divideByBasePrecision(big.NewInt(0).Mul(
			underlyingPrice,
			big.NewInt(0).Sub(_1e6, spreadRatioThreshold),
		),
		)
	}
	shortOrders := pipeline.db.GetShortOrders(market, shortCutOffPrice)

	// 3. Remove orders that were just cancelled
	longOrders = removeOrdersWithIds(longOrders, cancellableOrderIds)
	shortOrders = removeOrdersWithIds(shortOrders, cancellableOrderIds)

	// 4. Run liquidations
	modifiedLongOrders, modifiedShortOrders := pipeline.runLiquidations(market, longOrders, shortOrders)

	// 5. Run matching engine
	runMatchingEngine(pipeline.lotp, modifiedLongOrders, modifiedShortOrders)
}

func (pipeline *BuildBlockPipeline) cancelOrders(oraclePrices map[Market]*big.Int) map[common.Hash]struct{} {
	cancellableOrders := pipeline.db.GetOrdersToCancel(oraclePrices)
	cancellableOrderIds := map[common.Hash]struct{}{}
	// @todo: if there are too many cancellable orders, they might not fit in a single block. Need to adjust for that.
	for _, orderIds := range cancellableOrders {
		err := pipeline.lotp.ExecuteOrderCancel(orderIds)
		if err != nil {
			log.Error("Error in ExecuteOrderCancel", "orderIds", formatHashSlice(orderIds), "err", err)
		} else {
			for _, orderId := range orderIds {
				cancellableOrderIds[orderId] = struct{}{}
			}
		}
	}

	return cancellableOrderIds
}

func (pipeline *BuildBlockPipeline) runLiquidations(market Market, longOrders []LimitOrder, shortOrders []LimitOrder) (filteredLongOrder []LimitOrder, filteredShortOrder []LimitOrder) {
	if len(longOrders) == 0 && len(shortOrders) == 0 {
		return
	}
	oraclePrice := big.NewInt(20 * 10e6) // @todo: get it from the oracle

	liquidablePositions := GetLiquidableTraders(pipeline.db.GetAllTraders(), market, pipeline.db.GetLastPrice(market), oraclePrice)
	log.Info("liquidablePositions", "liquidablePositions", liquidablePositions)

	for i, liquidable := range liquidablePositions {
		var oppositeOrders []LimitOrder
		switch liquidable.PositionType {
		case "long":
			oppositeOrders = longOrders
		case "short":
			oppositeOrders = shortOrders
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
			pipeline.lotp.ExecuteLiquidation(liquidable.Address, oppositeOrder, fillAmount)

			switch liquidable.PositionType {
			case "long":
				oppositeOrders[j].FilledBaseAssetQuantity = big.NewInt(0).Add(oppositeOrders[j].FilledBaseAssetQuantity, fillAmount)
				liquidablePositions[i].FilledSize.Add(liquidablePositions[i].FilledSize, fillAmount)
			case "short":
				oppositeOrders[j].FilledBaseAssetQuantity = big.NewInt(0).Sub(oppositeOrders[j].FilledBaseAssetQuantity, fillAmount)
				liquidablePositions[i].FilledSize.Sub(liquidablePositions[i].FilledSize, fillAmount)
			}
		}
	}
	return longOrders, shortOrders
}

func runMatchingEngine(lotp LimitOrderTxProcessor, longOrders []LimitOrder, shortOrders []LimitOrder) {
	if len(longOrders) == 0 || len(shortOrders) == 0 {
		log.Info("BuildBlockPipeline - either no long or no short orders", "long", len(longOrders), "short", len(shortOrders))
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
			longOrders[i], shortOrders[j], ordersMatched = matchLongAndShortOrder(lotp, longOrders[i], shortOrders[j])
			if !ordersMatched {
				i = len(longOrders)
				break
			}
		}
	}
}

func matchLongAndShortOrder(lotp LimitOrderTxProcessor, longOrder LimitOrder, shortOrder LimitOrder) (LimitOrder, LimitOrder, bool) {
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

func isFundingPaymentTime(lastBlockTime uint64, db LimitOrderDatabase) bool {
	if db.GetNextFundingTime() == 0 {
		return false
	}
	return lastBlockTime >= db.GetNextFundingTime()
}

func executeFundingPayment(lotp LimitOrderTxProcessor) error {
	// @todo get index twap for each market with warp msging

	return lotp.ExecuteFundingPaymentTx()
}

func removeOrdersWithIds(orders []LimitOrder, orderIds map[common.Hash]struct{}) []LimitOrder {
	var filteredOrders []LimitOrder
	for _, order := range orders {
		if _, ok := orderIds[order.Id]; !ok {
			filteredOrders = append(filteredOrders, order)
		}
	}
	return filteredOrders
}

func formatHashSlice(hashes []common.Hash) []string {
	var formattedHashes []string
	for _, hash := range hashes {
		formattedHashes = append(formattedHashes, hash.String())
	}
	return formattedHashes
}
