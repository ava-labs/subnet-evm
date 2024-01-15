// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"encoding/hex"
	"encoding/json"
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
	OrderId string `json:"orderId,omitempty"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type PlaceSignedOrdersResponse struct {
	Orders []PlaceOrderResponse `json:"orders"`
}

func (api *OrderAPI) PlaceSignedOrders(ctx context.Context, input string) (PlaceSignedOrdersResponse, error) {
	// input is a json encoded array of signed orders
	var rawOrders []string
	err := json.Unmarshal([]byte(input), &rawOrders)
	if err != nil {
		return PlaceSignedOrdersResponse{}, err
	}

	ordersToGossip := []*hu.SignedOrder{}
	response := []PlaceOrderResponse{}
	for _, rawOrder := range rawOrders {
		orderResponse := PlaceOrderResponse{Success: false}
		testData, err := hex.DecodeString(strings.TrimPrefix(rawOrder, "0x"))
		if err != nil {
			orderResponse.Error = err.Error()
			response = append(response, orderResponse)
			continue
		}
		order, err := hu.DecodeSignedOrder(testData)
		if err != nil {
			orderResponse.Error = err.Error()
			response = append(response, orderResponse)
			continue
		}

		orderId, err := api.tradingAPI.PlaceOrder(order)
		orderResponse.OrderId = orderId.String()
		if err != nil {
			orderResponse.Error = err.Error()
			response = append(response, orderResponse)
			continue
		}
		orderResponse.Success = true
		response = append(response, orderResponse)
		ordersToGossip = append(ordersToGossip, order)
	}

	api.vm.gossiper.GossipSignedOrders(ordersToGossip)

	return PlaceSignedOrdersResponse{Orders: response}, nil
}
