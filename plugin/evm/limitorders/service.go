// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package limitorders

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type OrderBookAPI struct {
	db LimitOrderDatabase
}

func NewOrderBookAPI(database LimitOrderDatabase) *OrderBookAPI {
	return &OrderBookAPI{
		db: database,
	}
}

type OrderBookResponse struct {
	Orders []OrderMin
}

type OpenOrdersResponse struct {
	Orders []OrderForOpenOrders
}

type OrderMin struct {
	Market
	Price          string
	Size           string
	InprogressSize string
	Hash           string
}

type OrderForOpenOrders struct {
	Market
	Price      string
	Size       string
	FilledSize string
	Timestamp  uint64
	Salt       string
	Hash       string
}

func (api *OrderBookAPI) GetDetailedOrderBookData(ctx context.Context) InMemoryDatabase {
	return api.db.GetOrderBookData()
}

func (api *OrderBookAPI) GetOrderBook(ctx context.Context, marketStr string) (*OrderBookResponse, error) {
	// market is a string cuz it's an optional param
	allOrders := api.db.GetOrderBookData().OrderMap
	orders := []OrderMin{}

	if len(marketStr) > 0 {
		market, err := strconv.Atoi(marketStr)
		if err != nil {
			return nil, fmt.Errorf("invalid market")
		}
		marketOrders := map[common.Hash]*LimitOrder{}
		for hash, order := range allOrders {
			if order.Market == Market(market) {
				marketOrders[hash] = order
			}
		}
		allOrders = marketOrders
	}

	for hash, order := range allOrders {
		orders = append(orders, OrderMin{
			Market: order.Market,
			Price:  order.Price.String(),
			Size:   order.GetUnFilledBaseAssetQuantity().String(),
			Hash:   hash.String(),
		})
	}

	return &OrderBookResponse{Orders: orders}, nil
}

func (api *OrderBookAPI) GetOpenOrders(ctx context.Context, trader string) OpenOrdersResponse {
	traderOrders := []OrderForOpenOrders{}
	orderMap := api.db.GetOrderBookData().OrderMap
	for hash, order := range orderMap {
		if strings.EqualFold(order.UserAddress, trader) {
			traderOrders = append(traderOrders, OrderForOpenOrders{
				Market:     order.Market,
				Price:      order.Price.String(),
				Size:       order.BaseAssetQuantity.String(),
				FilledSize: order.FilledBaseAssetQuantity.String(),
				Salt:       getOrderFromRawOrder(order.RawOrder).Salt.String(),
				Hash:       hash.String(),
			})
		}
	}

	return OpenOrdersResponse{Orders: traderOrders}
}
