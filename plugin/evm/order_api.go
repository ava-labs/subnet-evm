// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"encoding/hex"
	"strings"

	"github.com/ava-labs/subnet-evm/plugin/evm/orderbook"
	hu "github.com/ava-labs/subnet-evm/plugin/evm/orderbook/hubbleutils"
)

type OrderAPI struct {
	tradingAPI *orderbook.TradingAPI
	vm         *VM
}

func NewOrderAPI(tradingAPI *orderbook.TradingAPI, vm *VM) *OrderAPI {
	return &OrderAPI{
		tradingAPI: tradingAPI,
		vm:         vm,
	}
}

type PlaceOrderResponse struct {
	Success bool `json:"success"`
}

func (api *OrderAPI) PlaceSignedOrder(ctx context.Context, rawOrder string) (PlaceOrderResponse, error) {
	testData, err := hex.DecodeString(strings.TrimPrefix(rawOrder, "0x"))
	if err != nil {
		return PlaceOrderResponse{Success: false}, err
	}
	order, err := hu.DecodeSignedOrder(testData)
	if err != nil {
		return PlaceOrderResponse{Success: false}, err
	}

	err = api.tradingAPI.PlaceOrder(order)
	if err != nil {
		return PlaceOrderResponse{Success: false}, err
	}
	api.vm.gossiper.GossipSignedOrders([]*hu.SignedOrder{order})

	return PlaceOrderResponse{Success: true}, nil
}
