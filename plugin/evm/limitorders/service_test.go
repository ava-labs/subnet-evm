package limitorders

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/eth"
	"github.com/stretchr/testify/assert"
)

func TestAggregatedOrderBook(t *testing.T) {
	t.Run("it aggregates long and short orders by price and returns aggregated data in json format with blockNumber", func(t *testing.T) {
		db := getDatabase()
		service := NewOrderBookAPI(db, &eth.EthAPIBackend{}, db.configService)

		longOrder1 := getLongOrder()
		db.Add(&longOrder1)

		longOrder2 := getLongOrder()
		longOrder2.Salt.Add(longOrder2.Salt, big.NewInt(100))
		longOrder2.Price.Mul(longOrder2.Price, big.NewInt(2))
		longOrder2.Id = getIdFromLimitOrder(longOrder2)
		db.Add(&longOrder2)

		shortOrder1 := getShortOrder()
		shortOrder1.Salt.Add(shortOrder1.Salt, big.NewInt(200))
		shortOrder1.Id = getIdFromLimitOrder(shortOrder1)
		db.Add(&shortOrder1)

		shortOrder2 := getShortOrder()
		shortOrder2.Salt.Add(shortOrder1.Salt, big.NewInt(300))
		shortOrder2.Price.Mul(shortOrder2.Price, big.NewInt(2))
		shortOrder2.Id = getIdFromLimitOrder(shortOrder2)
		db.Add(&shortOrder2)

		ctx := context.TODO()
		response := service.GetDepthForMarket(ctx, int(Market(0)))
		expectedAggregatedOrderBookState := MarketDepth{
			Market: Market(0),
			Longs: map[string]string{
				longOrder1.Price.String(): longOrder1.BaseAssetQuantity.String(),
				longOrder2.Price.String(): longOrder2.BaseAssetQuantity.String(),
			},
			Shorts: map[string]string{
				shortOrder1.Price.String(): shortOrder1.BaseAssetQuantity.String(),
				shortOrder2.Price.String(): shortOrder2.BaseAssetQuantity.String(),
			},
		}
		fmt.Println(response)
		assert.Equal(t, expectedAggregatedOrderBookState, *response)

		orderbook, _ := service.GetOrderBook(ctx, "0")
		assert.Equal(t, 4, len(orderbook.Orders))
	})
	t.Run("when event is the first event after subscribe", func(t *testing.T) {
		t.Run("when orderbook has no orders", func(t *testing.T) {
		})
		t.Run("when orderbook has either long or short orders with same price", func(t *testing.T) {
			t.Run("when order is longOrder", func(t *testing.T) {
			})
			t.Run("when order is shortOrder", func(t *testing.T) {
			})
			t.Run("when order is one long and one short order", func(t *testing.T) {
			})
			t.Run("when orderbook has more than one order of long or short", func(t *testing.T) {
				t.Run("when orderbook has more than one long orders", func(t *testing.T) {
				})
				t.Run("when orderbook has more than one short orders", func(t *testing.T) {
				})
				t.Run("when orderbook has more than one long and short orders", func(t *testing.T) {
				})
			})
		})
		t.Run("when orderbook has orders of same type with different price", func(t *testing.T) {
			t.Run("when orderbook has long orders of different price", func(t *testing.T) {
			})
			t.Run("when orderbook has short orders of different price", func(t *testing.T) {
			})
			t.Run("when orderbook has long and short orders of different price", func(t *testing.T) {
			})
		})
	})
	t.Run("for all events received after first event after subscribe", func(t *testing.T) {
		t.Run("when there are no changes in orderbook like cancel/matching", func(t *testing.T) {
		})
		t.Run("if there are changes in orderbook like cancel/matching", func(t *testing.T) {
		})
	})
}
