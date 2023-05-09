package limitorders

import (
	"context"
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/eth"
	"github.com/stretchr/testify/assert"
)

func TestAggregatedOrderBook(t *testing.T) {
	t.Run("it aggregates long and short orders by price and returns aggregated data in json format with blockNumber", func(t *testing.T) {
		db := NewInMemoryDatabase()
		service := NewOrderBookAPI(db, &eth.EthAPIBackend{})

		longOrder1 := getLongOrder()
		db.Add(getIdFromLimitOrder(longOrder1), &longOrder1)

		longOrder2 := getLongOrder()
		longOrder2.Salt.Add(longOrder2.Salt, big.NewInt(100))
		longOrder2.Price.Mul(longOrder2.Price, big.NewInt(2))
		db.Add(getIdFromLimitOrder(longOrder2), &longOrder2)

		shortOrder1 := getShortOrder()
		shortOrder1.Salt.Add(shortOrder1.Salt, big.NewInt(200))
		db.Add(getIdFromLimitOrder(shortOrder1), &shortOrder1)

		shortOrder2 := getShortOrder()
		shortOrder2.Salt.Add(shortOrder1.Salt, big.NewInt(300))
		shortOrder2.Price.Mul(shortOrder2.Price, big.NewInt(2))
		db.Add(getIdFromLimitOrder(shortOrder2), &shortOrder2)

		ctx := context.TODO()
		response := service.GetDepthForMarket(ctx, int(AvaxPerp))
		expectedAggregatedOrderBookState := MarketDepth{
			Market: AvaxPerp,
			Longs: map[string]string{
				longOrder1.Price.String(): longOrder1.BaseAssetQuantity.String(),
				longOrder2.Price.String(): longOrder2.BaseAssetQuantity.String(),
			},
			Shorts: map[string]string{
				shortOrder1.Price.String(): shortOrder1.BaseAssetQuantity.String(),
				shortOrder2.Price.String(): shortOrder2.BaseAssetQuantity.String(),
			},
		}
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
