// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package orderbook

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ava-labs/subnet-evm/eth"
	"github.com/ava-labs/subnet-evm/rpc"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
)

var traderFeed event.Feed
var marketFeed event.Feed

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

	pendingFunding := getTotalFunding(traderInfo, markets)
	margin := new(big.Int).Sub(getNormalisedMargin(traderInfo), pendingFunding)
	response.Margin = utils.BigIntToDecimal(margin, 6, 8)
	response.ReservedMargin = utils.BigIntToDecimal(traderInfo.Margin.Reserved, 6, 8)

	for market, position := range traderInfo.Positions {
		lastPrice := api.db.GetLastPrice(market)
		notionalPosition, uPnL, mf := getPositionMetadata(lastPrice, position.OpenNotional, position.Size, margin)

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
			MarkPrice:            utils.BigIntToDecimal(lastPrice, 6, 8),
		})
	}

	return response, nil
}

func (api *TradingAPI) StreamDepthUpdateForMarket(ctx context.Context, market int) (*rpc.Subscription, error) {
	notifier, _ := rpc.NotifierFromContext(ctx)
	rpcSub := notifier.CreateSubscription()

	ticker := time.NewTicker(1 * time.Second)

	var oldMarketDepth = &MarketDepth{}

	go func() {
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
	}()

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
	go func() {
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
	}()

	return rpcSub, nil
}

func (api *TradingAPI) StreamMarketTrades(ctx context.Context, market Market, blockStatus string) (*rpc.Subscription, error) {
	notifier, _ := rpc.NotifierFromContext(ctx)
	rpcSub := notifier.CreateSubscription()
	confirmationLevel := BlockConfirmationLevel(blockStatus)

	marketFeedCh := make(chan MarketFeedEvent)
	acceptedLogsSubscription := marketFeed.Subscribe(marketFeedCh)
	go func() {
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
	}()

	return rpcSub, nil
}
