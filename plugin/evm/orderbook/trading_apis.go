// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package orderbook

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ava-labs/subnet-evm/eth"
	hu "github.com/ava-labs/subnet-evm/plugin/evm/orderbook/hubbleutils"
	"github.com/ava-labs/subnet-evm/rpc"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
)

var traderFeed event.Feed
var marketFeed event.Feed
var MakerbookDatabaseFile = ""

type TradingAPI struct {
	db            LimitOrderDatabase
	backend       *eth.EthAPIBackend
	configService IConfigService
}

func NewTradingAPI(database LimitOrderDatabase, backend *eth.EthAPIBackend, configService IConfigService) *TradingAPI {
	return &TradingAPI{
		db:            database,
		backend:       backend,
		configService: configService,
	}
}

type TradingOrderBookDepthResponse struct {
	LastUpdateID int        `json:"lastUpdateId"`
	E            int64      `json:"E"`
	T            int64      `json:"T"`
	Symbol       int64      `json:"symbol"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

type TradingOrderBookDepthUpdateResponse struct {
	T      int64      `json:"T"`
	Symbol int64      `json:"s"`
	Bids   [][]string `json:"b"`
	Asks   [][]string `json:"a"`
}

// found at https://binance-docs.github.io/apidocs/futures/en/#query-order-user_data
// commented field values are from the binance docs
type OrderStatusResponse struct {
	ExecutedQty  string   `json:"executedQty"`  // "0"
	OrderID      string   `json:"orderId"`      // 1917641
	OrigQty      string   `json:"origQty"`      // "0.40"
	Price        string   `json:"price"`        // "0"
	ReduceOnly   bool     `json:"reduceOnly"`   // false
	PostOnly     bool     `json:"postOnly"`     // false
	PositionSide string   `json:"positionSide"` // "SHORT"
	Status       string   `json:"status"`       // "NEW"
	Symbol       int64    `json:"symbol"`       // "BTCUSDT"
	Time         int64    `json:"time"`         // 1579276756075
	Type         string   `json:"type"`         // "LIMIT"
	UpdateTime   int64    `json:"updateTime"`   // 1579276756075
	Salt         *big.Int `json:"salt"`
}

type TraderPosition struct {
	Market               Market `json:"market"`
	OpenNotional         string `json:"openNotional"`
	Size                 string `json:"size"`
	UnrealisedFunding    string `json:"unrealisedFunding"`
	LiquidationThreshold string `json:"liquidationThreshold"`
	NotionalPosition     string `json:"notionalPosition"`
	UnrealisedProfit     string `json:"unrealisedProfit"`
	MarginFraction       string `json:"marginFraction"`
	LiquidationPrice     string `json:"liquidationPrice"`
	MarkPrice            string `json:"markPrice"`
}

type GetPositionsResponse struct {
	Margin         string           `json:"margin"`
	ReservedMargin string           `json:"reservedMargin"`
	Positions      []TraderPosition `json:"positions"`
}

var mapStatus = map[Status]string{
	Placed:           "NEW",
	FulFilled:        "FILLED",
	Cancelled:        "CANCELED",
	Execution_Failed: "REJECTED",
}

func (api *TradingAPI) GetTradingOrderBookDepth(ctx context.Context, market int8) TradingOrderBookDepthResponse {
	response := TradingOrderBookDepthResponse{
		Asks: [][]string{},
		Bids: [][]string{},
	}
	depth := getDepthForMarket(api.db, Market(market))

	response = transformMarketDepth(depth)
	response.T = time.Now().Unix()

	return response
}

func (api *TradingAPI) GetOrderStatus(ctx context.Context, orderId common.Hash) (OrderStatusResponse, error) {
	response := OrderStatusResponse{}

	limitOrder := api.db.GetOrderById(orderId)
	if limitOrder == nil {
		return response, fmt.Errorf("order not found")
	}

	status := mapStatus[limitOrder.getOrderStatus().Status]
	if limitOrder.FilledBaseAssetQuantity.Sign() != 0 {
		status = "PARTIALLY_FILLED"
	}

	limitOrder.BaseAssetQuantity.Abs(limitOrder.BaseAssetQuantity)
	limitOrder.FilledBaseAssetQuantity.Abs(limitOrder.FilledBaseAssetQuantity)

	var positionSide string
	switch limitOrder.PositionType {
	case LONG:
		positionSide = "LONG"
	case SHORT:
		positionSide = "SHORT"
	}

	var time, updateTime int64
	placedBlock, err := api.backend.BlockByNumber(ctx, rpc.BlockNumber(limitOrder.BlockNumber.Int64()))
	if err == nil {
		time = int64(placedBlock.Time())
	}

	updateBlock, err := api.backend.BlockByNumber(ctx, rpc.BlockNumber(limitOrder.getOrderStatus().BlockNumber))
	if err == nil {
		updateTime = int64(updateBlock.Time())
	}

	response = OrderStatusResponse{
		ExecutedQty:  utils.BigIntToDecimal(limitOrder.FilledBaseAssetQuantity, 18, 8),
		OrderID:      limitOrder.Id.String(),
		OrigQty:      utils.BigIntToDecimal(limitOrder.BaseAssetQuantity, 18, 8),
		Price:        utils.BigIntToDecimal(limitOrder.Price, 6, 8),
		ReduceOnly:   limitOrder.ReduceOnly,
		PostOnly:     limitOrder.isPostOnly(),
		PositionSide: positionSide,
		Status:       status,
		Symbol:       int64(limitOrder.Market),
		Time:         time,
		Type:         "LIMIT_ORDER",
		UpdateTime:   updateTime,
		Salt:         limitOrder.Salt,
	}

	return response, nil
}

func (api *TradingAPI) GetMarginAndPositions(ctx context.Context, trader string) (GetPositionsResponse, error) {
	response := GetPositionsResponse{Positions: []TraderPosition{}}

	traderAddr := common.HexToAddress(trader)

	traderInfo := api.db.GetTraderInfo(traderAddr)
	if traderInfo == nil {
		return response, fmt.Errorf("trader not found")
	}

	count := api.configService.GetActiveMarketsCount()
	markets := make([]Market, count)
	for i := int64(0); i < count; i++ {
		markets[i] = Market(i)
	}

	assets := api.configService.GetCollaterals()
	pendingFunding := getTotalFunding(traderInfo, markets)
	margin := new(big.Int).Sub(getNormalisedMargin(traderInfo, assets), pendingFunding)
	response.Margin = utils.BigIntToDecimal(margin, 6, 8)
	response.ReservedMargin = utils.BigIntToDecimal(traderInfo.Margin.Reserved, 6, 8)

	for market, position := range traderInfo.Positions {
		midPrice := api.configService.GetMidPrices()[market]
		notionalPosition, uPnL, mf := getPositionMetadata(midPrice, position.OpenNotional, position.Size, margin)

		response.Positions = append(response.Positions, TraderPosition{
			Market:               market,
			OpenNotional:         utils.BigIntToDecimal(position.OpenNotional, 6, 8),
			Size:                 utils.BigIntToDecimal(position.Size, 18, 8),
			UnrealisedFunding:    utils.BigIntToDecimal(position.UnrealisedFunding, 6, 8),
			LiquidationThreshold: utils.BigIntToDecimal(position.LiquidationThreshold, 6, 8),
			UnrealisedProfit:     utils.BigIntToDecimal(uPnL, 6, 8),
			MarginFraction:       utils.BigIntToDecimal(mf, 6, 8),
			NotionalPosition:     utils.BigIntToDecimal(notionalPosition, 6, 8),
			LiquidationPrice:     "0", // todo: calculate
			MarkPrice:            utils.BigIntToDecimal(midPrice, 6, 8),
		})
	}

	return response, nil
}

// used by the sdk
func (api *TradingAPI) StreamDepthUpdateForMarket(ctx context.Context, market int) (*rpc.Subscription, error) {
	notifier, _ := rpc.NotifierFromContext(ctx)
	rpcSub := notifier.CreateSubscription()

	ticker := time.NewTicker(1 * time.Second)

	var oldMarketDepth = &MarketDepth{}

	go executeFuncAndRecoverPanic(func() {
		for {
			select {
			case <-ticker.C:
				newMarketDepth := getDepthForMarket(api.db, Market(market))
				depthUpdate := getUpdateInDepth(newMarketDepth, oldMarketDepth)
				transformedDepthUpdate := transformMarketDepth(depthUpdate)
				response := TradingOrderBookDepthUpdateResponse{
					T:      time.Now().Unix(),
					Symbol: int64(market),
					Bids:   transformedDepthUpdate.Bids,
					Asks:   transformedDepthUpdate.Asks,
				}
				notifier.Notify(rpcSub.ID, response)
				oldMarketDepth = newMarketDepth
			case <-notifier.Closed():
				ticker.Stop()
				return
			}
		}
	}, "panic in StreamDepthUpdateForMarket", RPCPanicsCounter)

	return rpcSub, nil
}

func transformMarketDepth(depth *MarketDepth) TradingOrderBookDepthResponse {
	response := TradingOrderBookDepthResponse{}
	for price, quantity := range depth.Longs {
		bigPrice, _ := big.NewInt(0).SetString(price, 10)
		bigQuantity, _ := big.NewInt(0).SetString(quantity, 10)
		response.Bids = append(response.Bids, []string{
			utils.BigIntToDecimal(bigPrice, 6, 8),
			utils.BigIntToDecimal(bigQuantity, 18, 8),
		})
	}

	for price, quantity := range depth.Shorts {
		bigPrice, _ := big.NewInt(0).SetString(price, 10)
		bigQuantity, _ := big.NewInt(0).SetString(quantity, 10)
		response.Asks = append(response.Asks, []string{
			utils.BigIntToDecimal(bigPrice, 6, 8),
			utils.BigIntToDecimal(bigQuantity, 18, 8),
		})
	}

	return response
}

func (api *TradingAPI) StreamTraderUpdates(ctx context.Context, trader string, blockStatus string) (*rpc.Subscription, error) {
	notifier, _ := rpc.NotifierFromContext(ctx)
	rpcSub := notifier.CreateSubscription()
	confirmationLevel := BlockConfirmationLevel(blockStatus)

	traderFeedCh := make(chan TraderEvent)
	traderFeedSubscription := traderFeed.Subscribe(traderFeedCh)
	go executeFuncAndRecoverPanic(func() {
		defer traderFeedSubscription.Unsubscribe()

		for {
			select {
			case event := <-traderFeedCh:
				if strings.EqualFold(event.Trader.String(), trader) && event.BlockStatus == confirmationLevel {
					notifier.Notify(rpcSub.ID, event)
				}
			case <-notifier.Closed():
				return
			}
		}
	}, "panic in StreamTraderUpdates", RPCPanicsCounter)

	return rpcSub, nil
}

func (api *TradingAPI) StreamMarketTrades(ctx context.Context, market Market, blockStatus string) (*rpc.Subscription, error) {
	notifier, _ := rpc.NotifierFromContext(ctx)
	rpcSub := notifier.CreateSubscription()
	confirmationLevel := BlockConfirmationLevel(blockStatus)

	marketFeedCh := make(chan MarketFeedEvent)
	acceptedLogsSubscription := marketFeed.Subscribe(marketFeedCh)
	go executeFuncAndRecoverPanic(func() {
		defer acceptedLogsSubscription.Unsubscribe()

		for {
			select {
			case event := <-marketFeedCh:
				if event.Market == market && event.BlockStatus == confirmationLevel {
					notifier.Notify(rpcSub.ID, event)
				}
			case <-notifier.Closed():
				return
			}
		}
	}, "panic in StreamMarketTrades", RPCPanicsCounter)

	return rpcSub, nil
}

// @todo cache api.configService values to avoid db lookups on every order placement
func (api *TradingAPI) PlaceOrder(order *hu.SignedOrder) (common.Hash, error) {
	orderId, err := order.Hash()
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to hash order: %s", err)
	}
	fields := api.db.GetOrderValidationFields(orderId, order)
	// P1. Order is not already in memdb
	if fields.Exists {
		return orderId, hu.ErrOrderAlreadyExists
	}
	marketId := int(order.AmmIndex.Int64())
	trader, signer, err := hu.ValidateSignedOrder(
		order,
		hu.SignedOrderValidationFields{
			OrderHash:          orderId,
			Now:                uint64(time.Now().Unix()),
			ActiveMarketsCount: api.configService.GetActiveMarketsCount(),
			MinSize:            api.configService.getMinSizeRequirement(marketId),
			PriceMultiplier:    api.configService.GetPriceMultiplier(marketId),
			Status:             api.configService.GetSignedOrderStatus(orderId),
		},
	)
	if err != nil {
		return orderId, err
	}
	if trader != signer && !api.configService.IsTradingAuthority(trader, signer) {
		log.Error("not trading authority", "trader", trader.String(), "signer", signer.String())
		return orderId, hu.ErrNoTradingAuthority
	}

	requiredMargin := big.NewInt(0)
	if !order.ReduceOnly {
		// P2. Margin is available for non-reduce only orders
		minAllowableMargin := api.configService.GetMinAllowableMargin()
		// even tho order might be matched at a different price, we reserve margin at the price the order was placed at to keep it simple
		requiredMargin = hu.GetRequiredMargin(order.Price, hu.Abs(order.BaseAssetQuantity), minAllowableMargin, big.NewInt(0))
		if fields.AvailableMargin.Cmp(requiredMargin) == -1 {
			return orderId, hu.ErrInsufficientMargin
		}
	} else {
		// @todo P3. Sum of all reduce only orders should not exceed the total position size
		return orderId, errors.New("reduce only orders via makerbook are not supported yet")
	}

	// P4. Post only order shouldn't cross the market
	if order.PostOnly {
		orderSide := hu.Side(hu.Long)
		if order.BaseAssetQuantity.Sign() == -1 {
			orderSide = hu.Side(hu.Short)
		}
		asksHead := fields.AsksHead
		bidsHead := fields.BidsHead
		if (orderSide == hu.Side(hu.Short) && bidsHead.Sign() != 0 && order.Price.Cmp(bidsHead) != 1) || (orderSide == hu.Side(hu.Long) && asksHead.Sign() != 0 && order.Price.Cmp(asksHead) != -1) {
			return orderId, hu.ErrCrossingMarket
		}
	}

	// P5. HasReferrer
	if !api.configService.HasReferrer(order.Trader) {
		return orderId, hu.ErrNoReferrer
	}

	// validations passed, add to db
	signedOrder := &Order{
		Id:                      orderId,
		Market:                  Market(order.AmmIndex.Int64()),
		PositionType:            getPositionTypeBasedOnBaseAssetQuantity(order.BaseAssetQuantity),
		Trader:                  trader,
		BaseAssetQuantity:       order.BaseAssetQuantity,
		FilledBaseAssetQuantity: big.NewInt(0),
		Price:                   order.Price,
		Salt:                    order.Salt,
		ReduceOnly:              order.ReduceOnly,
		BlockNumber:             big.NewInt(0),
		RawOrder:                order,
		OrderType:               Signed,
	}
	placeSignedOrderCounter.Inc(1)
	api.db.AddSignedOrder(signedOrder, requiredMargin)

	if len(MakerbookDatabaseFile) > 0 {
		go writeOrderToFile(order, orderId)
	}

	// send to trader feed - both for head and accepted block
	go func() {
		orderMap := order.Map()
		orderMap["orderType"] = "signed"
		orderMap["expireAt"] = order.ExpireAt.String()
		args := map[string]interface{}{
			"order": orderMap,
		}

		traderEvent := TraderEvent{
			Trader:      trader,
			Removed:     false,
			EventName:   "OrderAccepted",
			Args:        args,
			BlockStatus: ConfirmationLevelHead,
			OrderId:     orderId,
			OrderType:   Signed.String(),
			Timestamp:   big.NewInt(time.Now().Unix()),
		}

		traderFeed.Send(traderEvent)

		traderEvent.BlockStatus = ConfirmationLevelAccepted
		traderFeed.Send(traderEvent)
	}()

	return orderId, nil
}

func writeOrderToFile(order *hu.SignedOrder, orderId common.Hash) {
	orderMap := order.Map()
	orderMap["orderType"] = "signed"
	orderMap["expireAt"] = order.ExpireAt.String()
	doc := map[string]interface{}{
		"type":      "OrderAccepted",
		"timestamp": time.Now().Unix(),
		"trader":    order.Trader.String(),
		"orderHash": strings.ToLower(orderId.String()),
		"orderType": "signed",
		"order": map[string]interface{}{
			"orderType":         2,
			"expireAt":          order.ExpireAt.Uint64(),
			"ammIndex":          order.AmmIndex.Uint64(),
			"trader":            order.Trader.String(),
			"baseAssetQuantity": utils.BigIntToFloat(order.BaseAssetQuantity, 18),
			"price":             utils.BigIntToFloat(order.Price, 6),
			"salt":              order.Salt.Int64(),
			"reduceOnly":        order.ReduceOnly,
		},
	}
	jsonDoc, err := json.Marshal(doc)
	if err != nil {
		log.Error("writeOrderToFile: failed to marshal order", "err", err)
		makerBookWriteFailuresCounter.Inc(1)
		return
	}
	err = utils.AppendToFile(MakerbookDatabaseFile, jsonDoc)
	if err != nil {
		log.Error("writeOrderToFile: failed to write order to file", "err", err)
		makerBookWriteFailuresCounter.Inc(1)
	}
}
