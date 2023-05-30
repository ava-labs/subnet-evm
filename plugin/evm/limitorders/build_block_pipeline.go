package limitorders

import (
	"math/big"
	"time"

	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type BuildBlockPipeline struct {
	db            LimitOrderDatabase
	lotp          LimitOrderTxProcessor
	configService IConfigService
}

func NewBuildBlockPipeline(db LimitOrderDatabase, lotp LimitOrderTxProcessor, configService IConfigService) *BuildBlockPipeline {
	return &BuildBlockPipeline{
		db:            db,
		lotp:          lotp,
		configService: configService,
	}
}

func (pipeline *BuildBlockPipeline) Run() {
	markets := pipeline.GetActiveMarkets()

	if len(markets) == 0 {
		return
	}

	pipeline.lotp.PurgeLocalTx()

	if isFundingPaymentTime(pipeline.db.GetNextFundingTime()) {
		log.Info("BuildBlockPipeline:isFundingPaymentTime")
		// just execute the funding payment and skip running the matching engine
		err := executeFundingPayment(pipeline.lotp)
		if err != nil {
			log.Error("Funding payment job failed", "err", err)
		}
		return
	}

	// fetch the underlying price and run the matching engine
	underlyingPrices, err := pipeline.lotp.GetUnderlyingPrice()
	if err != nil {
		log.Error("could not fetch underlying price", "err", err)
		return
	}

	// build trader map
	liquidablePositions, ordersToCancel := pipeline.db.GetNaughtyTraders(underlyingPrices, markets)
	cancellableOrderIds := pipeline.cancelOrders(ordersToCancel)
	orderMap := make(map[Market]*Orders)
	for _, market := range markets {
		orderMap[market] = pipeline.fetchOrders(market, underlyingPrices[market], cancellableOrderIds)
	}
	pipeline.runLiquidations(liquidablePositions, orderMap)
	for _, market := range markets {
		pipeline.runMatchingEngine(pipeline.lotp, orderMap[market].longOrders, orderMap[market].shortOrders)
	}
}

type Orders struct {
	longOrders  []LimitOrder
	shortOrders []LimitOrder
}

type Market int64

func (pipeline *BuildBlockPipeline) GetActiveMarkets() []Market {
	count := pipeline.configService.GetActiveMarketsCount()
	markets := make([]Market, count)
	for i := int64(0); i < count; i++ {
		markets[i] = Market(i)
	}
	return markets
}

func (pipeline *BuildBlockPipeline) cancelOrders(cancellableOrders map[common.Address][]LimitOrder) map[common.Hash]struct{} {
	cancellableOrderIds := map[common.Hash]struct{}{}
	// @todo: if there are too many cancellable orders, they might not fit in a single block. Need to adjust for that.
	for _, orders := range cancellableOrders {
		if len(orders) > 0 {
			rawOrders := make([]Order, len(orders))
			for i, order := range orders {
				rawOrders[i] = order.RawOrder
			}
			err := pipeline.lotp.ExecuteOrderCancel(rawOrders)
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

func (pipeline *BuildBlockPipeline) fetchOrders(market Market, underlyingPrice *big.Int, cancellableOrderIds map[common.Hash]struct{}) *Orders {
	spreadRatioThreshold := pipeline.configService.getSpreadRatioThreshold(market)
	// 1. Get long orders
	longCutOffPrice := divideByBasePrecision(big.NewInt(0).Mul(underlyingPrice, big.NewInt(0).Add(_1e6, spreadRatioThreshold)))
	longOrders := pipeline.db.GetLongOrders(market, longCutOffPrice)

	// 2. Get short orders
	shortCutOffPrice := big.NewInt(0)
	if _1e6.Cmp(spreadRatioThreshold) > 0 {
		shortCutOffPrice = divideByBasePrecision(big.NewInt(0).Mul(underlyingPrice, big.NewInt(0).Sub(_1e6, spreadRatioThreshold)))
	}
	shortOrders := pipeline.db.GetShortOrders(market, shortCutOffPrice)

	// 3. Remove orders that were just cancelled
	longOrders = removeOrdersWithIds(longOrders, cancellableOrderIds)
	shortOrders = removeOrdersWithIds(shortOrders, cancellableOrderIds)

	return &Orders{longOrders, shortOrders}
}

func (pipeline *BuildBlockPipeline) runLiquidations(liquidablePositions []LiquidablePosition, orderMap map[Market]*Orders) {
	if len(liquidablePositions) > 0 {
		log.Info("found positions to liquidate", "liquidablePositions", liquidablePositions)
	}

	for i, liquidable := range liquidablePositions {
		var oppositeOrders []LimitOrder
		switch liquidable.PositionType {
		case LONG:
			oppositeOrders = orderMap[liquidable.Market].longOrders
		case SHORT:
			oppositeOrders = orderMap[liquidable.Market].shortOrders
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
			case LONG:
				oppositeOrders[j].FilledBaseAssetQuantity.Add(oppositeOrders[j].FilledBaseAssetQuantity, fillAmount)
				liquidablePositions[i].FilledSize.Add(liquidablePositions[i].FilledSize, fillAmount)
			case SHORT:
				oppositeOrders[j].FilledBaseAssetQuantity.Sub(oppositeOrders[j].FilledBaseAssetQuantity, fillAmount)
				liquidablePositions[i].FilledSize.Sub(liquidablePositions[i].FilledSize, fillAmount)
			}
		}
	}
}

func (pipeline *BuildBlockPipeline) runMatchingEngine(lotp LimitOrderTxProcessor, longOrders []LimitOrder, shortOrders []LimitOrder) {
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

func isFundingPaymentTime(nextFundingTime uint64) bool {
	if nextFundingTime == 0 {
		return false
	}

	now := uint64(time.Now().Unix())
	return now >= nextFundingTime
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
