// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package limitorders

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/eth"
	"github.com/ava-labs/subnet-evm/rpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
)

type OrderBookAPI struct {
	db      LimitOrderDatabase
	backend *eth.EthAPIBackend
}

func NewOrderBookAPI(database LimitOrderDatabase, backend *eth.EthAPIBackend) *OrderBookAPI {
	return &OrderBookAPI{
		db:      database,
		backend: backend,
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
	Price   string
	Size    string
	Signer  string
	OrderId string
}

type OrderForOpenOrders struct {
	Market
	Price      string
	Size       string
	FilledSize string
	Timestamp  uint64
	Salt       string
	OrderId    string
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
				if order.PositionType == "short" /* || order.Price.Cmp(big.NewInt(20e6)) <= 0 */ {
					marketOrders[hash] = order
				}
			}
		}
		allOrders = marketOrders
	}

	for hash, order := range allOrders {
		orders = append(orders, OrderMin{
			Market:  order.Market,
			Price:   order.Price.String(),
			Size:    order.GetUnFilledBaseAssetQuantity().String(),
			Signer:  order.UserAddress,
			OrderId: hash.String(),
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
				OrderId:    hash.String(),
			})
		}
	}

	return OpenOrdersResponse{Orders: traderOrders}
}

// NewOrderBookState send a notification each time a new (header) block is appended to the chain.
func (api *OrderBookAPI) NewOrderBookState(ctx context.Context) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}

	rpcSub := notifier.CreateSubscription()

	go func() {
		var (
			headers    = make(chan core.ChainHeadEvent)
			headersSub event.Subscription
		)

		headersSub = api.backend.SubscribeChainHeadEvent(headers)
		defer headersSub.Unsubscribe()

		for {
			select {
			case <-headers:
				orderBookData := api.GetDetailedOrderBookData(ctx)
				notifier.Notify(rpcSub.ID, &orderBookData)
			case <-rpcSub.Err():
				headersSub.Unsubscribe()
				return
			case <-notifier.Closed():
				headersSub.Unsubscribe()
				return
			}
		}
	}()

	return rpcSub, nil
}
