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
		err := executeFundingPayment(pipeline.lotp)
		if err != nil {
			log.Error("Funding payment job failed", "err", err)
		}
	}

	// fetch the underlying price and run the matching engine
	underlyingPrices := pipeline.GetUnderlyingPrices()

	// build trader map
	liquidablePositions, ordersToCancel := pipeline.db.GetNaughtyTraders(underlyingPrices, markets)
	cancellableOrderIds := pipeline.cancelOrders(ordersToCancel)
	orderMap := make(map[Market]*Orders)
	for _, market := range markets {
		orderMap[market] = pipeline.fetchOrders(market, underlyingPrices[market], cancellableOrderIds)
	}
	pipeline.runLiquidations(liquidablePositions, orderMap, underlyingPrices)
	for _, market := range markets {
		// @todo should we prioritize matching in any particular market?
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

func (pipeline *BuildBlockPipeline) GetUnderlyingPrices() map[Market]*big.Int {
	prices := pipeline.configService.GetUnderlyingPrices()
	log.Info("GetUnderlyingPrices", "prices", prices)
	underlyingPrices := make(map[Market]*big.Int)
	for market, price := range prices {
		underlyingPrices[Market(market)] = price
	}
	return underlyingPrices
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
	_, lowerbound := pipeline.configService.GetAcceptableBounds(market)

	// any long orders below the permissible lowerbound are irrelevant, because they won't be matched no matter what.
	// this assumes that all above cancelOrder transactions got executed successfully
	longOrders := removeOrdersWithIds(pipeline.db.GetLongOrders(market, lowerbound), cancellableOrderIds)

	var shortOrders []LimitOrder
	// all short orders above price of the highest long order are irrelevant
	if len(longOrders) > 0 {
		shortOrders = removeOrdersWithIds(pipeline.db.GetShortOrders(market, longOrders[0].Price /* upperbound */), cancellableOrderIds)
	}
	return &Orders{longOrders, shortOrders}
}

func (pipeline *BuildBlockPipeline) runLiquidations(liquidablePositions []LiquidablePosition, orderMap map[Market]*Orders, underlyingPrices map[Market]*big.Int) {
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

func (pipeline *BuildBlockPipeline) runMatchingEngine(lotp LimitOrderTxProcessor, longOrders []LimitOrder, shortOrders []LimitOrder) {
	if len(longOrders) == 0 || len(shortOrders) == 0 {
		return
	}

	matchingComplete := false
	for i := 0; i < len(longOrders); i++ {
		numOrdersExhausted := 0
		for j := 0; j < len(shortOrders); j++ {
			var ordersMatched bool
			longOrders[i], shortOrders[j], ordersMatched = matchLongAndShortOrder(lotp, longOrders[i], shortOrders[j])
			if !ordersMatched {
				matchingComplete = true
				break

			}
			if shortOrders[j].GetUnFilledBaseAssetQuantity().Sign() == 0 {
				numOrdersExhausted++
			}
			if longOrders[i].GetUnFilledBaseAssetQuantity().Sign() == 0 {
				break
			}
		}
		if matchingComplete {
			break
		}
		shortOrders = shortOrders[numOrdersExhausted:]
	}
}

func matchLongAndShortOrder(lotp LimitOrderTxProcessor, longOrder, shortOrder LimitOrder) (LimitOrder, LimitOrder, bool) {
	fillAmount := utils.BigIntMinAbs(longOrder.GetUnFilledBaseAssetQuantity(), shortOrder.GetUnFilledBaseAssetQuantity())
	if longOrder.Price.Cmp(shortOrder.Price) == -1 || fillAmount.Sign() == 0 {
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
