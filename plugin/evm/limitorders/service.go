// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package limitorders

import (
	"context"
	"math/big"
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

type OrderMin struct {
	Market
	Price *big.Int
	Size  *big.Int
}

func (api *OrderBookAPI) GetDetailedOrderBookData(ctx context.Context) InMemoryDatabase {
	return api.db.GetOrderBookData()
}

func (api *OrderBookAPI) GetOrderBook(ctx context.Context) OrderBookResponse {
	allOrders := api.db.GetAllOrders()
	orders := []OrderMin{}

	for _, order := range allOrders {
		orders = append(orders, OrderMin{
			Market: order.Market,
			Price:  order.Price,
			Size:   order.GetUnFilledBaseAssetQuantity(),
		})
	}

	return OrderBookResponse{Orders: orders}
}
