package orderbook

import (
	"math/big"
	"testing"
	"time"

	hu "github.com/ava-labs/subnet-evm/plugin/evm/orderbook/hubbleutils"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRunLiquidations(t *testing.T) {
	traderAddress := common.HexToAddress("0x710bf5f942331874dcbc7783319123679033b63b")
	traderAddress1 := common.HexToAddress("0x376c47978271565f56DEB45495afa69E59c16Ab2")
	market := Market(0)
	liqUpperBound := big.NewInt(22)
	liqLowerBound := big.NewInt(18)

	t.Run("when there are no liquidable positions", func(t *testing.T) {
		_, lotp, pipeline, underlyingPrices, _ := setupDependencies(t)
		longOrders := []Order{getLongOrder()}
		shortOrders := []Order{getShortOrder()}

		orderMap := map[Market]*Orders{market: {longOrders, shortOrders}}
		pipeline.runLiquidations([]LiquidablePosition{}, orderMap, underlyingPrices, map[common.Address]*big.Int{})
		assert.Equal(t, longOrders, orderMap[market].longOrders)
		assert.Equal(t, shortOrders, orderMap[market].shortOrders)
		lotp.AssertNotCalled(t, "ExecuteLiquidation", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("when liquidable position is long", func(t *testing.T) {
		t.Run("when no long orders are present in database for matching", func(t *testing.T) {
			_, lotp, pipeline, underlyingPrices, cs := setupDependencies(t)
			longOrders := []Order{}
			shortOrders := []Order{getShortOrder()}

			orderMap := map[Market]*Orders{market: {longOrders, shortOrders}}

			cs.On("GetAcceptableBoundsForLiquidation", market).Return(liqUpperBound, liqLowerBound)
			cs.On("getMinAllowableMargin").Return(big.NewInt(1e5))
			cs.On("GetTakerFee").Return(big.NewInt(1e5))
			pipeline.runLiquidations([]LiquidablePosition{getLiquidablePos(traderAddress, LONG, 7)}, orderMap, underlyingPrices, map[common.Address]*big.Int{})
			assert.Equal(t, longOrders, orderMap[market].longOrders)
			assert.Equal(t, shortOrders, orderMap[market].shortOrders)
			lotp.AssertNotCalled(t, "ExecuteLiquidation", mock.Anything, mock.Anything, mock.Anything)
			cs.AssertCalled(t, "GetAcceptableBoundsForLiquidation", market)
		})
		t.Run("when long orders are present in database for matching", func(t *testing.T) {
			liquidablePositions := []LiquidablePosition{getLiquidablePos(traderAddress, LONG, 7)}
			_, lotp, pipeline, underlyingPrices, cs := setupDependencies(t)
			longOrder := getLongOrder()
			shortOrder := getShortOrder()
			expectedFillAmount := utils.BigIntMinAbs(longOrder.BaseAssetQuantity, liquidablePositions[0].Size)
			cs.On("GetAcceptableBoundsForLiquidation", market).Return(liqUpperBound, liqLowerBound)
			cs.On("getMinAllowableMargin").Return(big.NewInt(1e5))
			cs.On("GetTakerFee").Return(big.NewInt(1e5))
			lotp.On("ExecuteLiquidation", traderAddress, longOrder, expectedFillAmount).Return(nil)

			orderMap := map[Market]*Orders{market: {[]Order{longOrder}, []Order{shortOrder}}}

			pipeline.runLiquidations(liquidablePositions, orderMap, underlyingPrices, map[common.Address]*big.Int{})

			lotp.AssertCalled(t, "ExecuteLiquidation", traderAddress, longOrder, expectedFillAmount)
			cs.AssertCalled(t, "GetAcceptableBoundsForLiquidation", market)
			assert.Equal(t, shortOrder, orderMap[market].shortOrders[0])
			assert.Equal(t, expectedFillAmount.Uint64(), orderMap[market].longOrders[0].FilledBaseAssetQuantity.Uint64())
		})
		t.Run("2nd long order < liqLowerBound", func(t *testing.T) {
			liquidablePositions := []LiquidablePosition{getLiquidablePos(traderAddress, LONG, 7)}
			_, lotp, pipeline, underlyingPrices, cs := setupDependencies(t)
			longOrder := getLongOrder()
			longOrder.BaseAssetQuantity = big.NewInt(5) // 5 < liquidable.Size (7)

			longOrder2 := getLongOrder()
			longOrder2.Price = big.NewInt(17) // 17 < lower bound (18)

			expectedFillAmount := utils.BigIntMinAbs(longOrder.BaseAssetQuantity, liquidablePositions[0].Size)
			cs.On("GetAcceptableBoundsForLiquidation", market).Return(liqUpperBound, liqLowerBound)
			cs.On("getMinAllowableMargin").Return(big.NewInt(1e5))
			cs.On("GetTakerFee").Return(big.NewInt(1e5))
			lotp.On("ExecuteLiquidation", traderAddress, longOrder, expectedFillAmount).Return(nil)

			orderMap := map[Market]*Orders{market: {[]Order{longOrder, longOrder2}, []Order{}}}

			pipeline.runLiquidations(liquidablePositions, orderMap, underlyingPrices, map[common.Address]*big.Int{})

			lotp.AssertCalled(t, "ExecuteLiquidation", traderAddress, longOrder, expectedFillAmount)
			cs.AssertCalled(t, "GetAcceptableBoundsForLiquidation", market)
			assert.Equal(t, 1, len(orderMap[market].longOrders))              // 0th order was consumed
			assert.Equal(t, longOrder2, orderMap[market].longOrders[0])       // untouched
			assert.Equal(t, big.NewInt(5), liquidablePositions[0].FilledSize) // 7 - 5
		})

		t.Run("4 liquidable positions", func(t *testing.T) {
			liquidablePositions := []LiquidablePosition{getLiquidablePos(traderAddress, LONG, 7), getLiquidablePos(traderAddress, SHORT, -8), getLiquidablePos(traderAddress1, LONG, 9), getLiquidablePos(traderAddress1, SHORT, -2)}
			_, lotp, pipeline, underlyingPrices, cs := setupDependencies(t)
			longOrder0 := buildLongOrder(20, 5)
			longOrder1 := buildLongOrder(19, 12)

			shortOrder0 := buildShortOrder(19, -9)
			shortOrder1 := buildShortOrder(liqLowerBound.Int64()-1, -8)
			orderMap := map[Market]*Orders{market: {[]Order{longOrder0, longOrder1}, []Order{shortOrder0, shortOrder1}}}

			cs.On("GetAcceptableBoundsForLiquidation", market).Return(liqUpperBound, liqLowerBound)
			cs.On("getMinAllowableMargin").Return(big.NewInt(1e5))
			cs.On("GetTakerFee").Return(big.NewInt(1e5))
			lotp.On("ExecuteLiquidation", traderAddress, orderMap[market].longOrders[0], big.NewInt(5)).Return(nil)
			lotp.On("ExecuteLiquidation", traderAddress, orderMap[market].longOrders[1], big.NewInt(2)).Return(nil)
			lotp.On("ExecuteLiquidation", traderAddress1, orderMap[market].longOrders[1], big.NewInt(9)).Return(nil)
			lotp.On("ExecuteLiquidation", traderAddress, orderMap[market].shortOrders[0], big.NewInt(8)).Return(nil)
			lotp.On("ExecuteLiquidation", traderAddress1, orderMap[market].shortOrders[0], big.NewInt(1)).Return(nil)
			lotp.On("ExecuteLiquidation", traderAddress1, orderMap[market].shortOrders[1], big.NewInt(1)).Return(nil)

			pipeline.runLiquidations(liquidablePositions, orderMap, underlyingPrices, map[common.Address]*big.Int{})
			cs.AssertCalled(t, "GetAcceptableBoundsForLiquidation", market)

			lotp.AssertCalled(t, "ExecuteLiquidation", traderAddress, longOrder0, big.NewInt(5))
			lotp.AssertCalled(t, "ExecuteLiquidation", traderAddress, longOrder0, big.NewInt(5))
			lotp.AssertCalled(t, "ExecuteLiquidation", traderAddress1, longOrder1, big.NewInt(9))
			lotp.AssertCalled(t, "ExecuteLiquidation", traderAddress, shortOrder0, big.NewInt(8))
			lotp.AssertCalled(t, "ExecuteLiquidation", traderAddress1, shortOrder0, big.NewInt(1))
			lotp.AssertCalled(t, "ExecuteLiquidation", traderAddress1, shortOrder1, big.NewInt(1))

			assert.Equal(t, 1, len(orderMap[market].longOrders)) // 0th order was consumed
			assert.Equal(t, big.NewInt(11), orderMap[market].longOrders[0].FilledBaseAssetQuantity)
			assert.Equal(t, big.NewInt(7), liquidablePositions[0].FilledSize)
			assert.Equal(t, big.NewInt(9), liquidablePositions[2].FilledSize)

			assert.Equal(t, 1, len(orderMap[market].shortOrders))
			assert.Equal(t, big.NewInt(-1), orderMap[market].shortOrders[0].FilledBaseAssetQuantity)
			assert.Equal(t, big.NewInt(-8), liquidablePositions[1].FilledSize)
			assert.Equal(t, big.NewInt(-2), liquidablePositions[3].FilledSize)
		})
	})

	t.Run("when liquidable position is short", func(t *testing.T) {
		liquidablePositions := []LiquidablePosition{{
			Address:      traderAddress,
			Market:       market,
			PositionType: SHORT,
			Size:         big.NewInt(-7),
			FilledSize:   big.NewInt(0),
		}}
		t.Run("when no short orders are present in database for matching", func(t *testing.T) {
			_, lotp, pipeline, underlyingPrices, cs := setupDependencies(t)
			shortOrders := []Order{}
			longOrders := []Order{getLongOrder()}

			orderMap := map[Market]*Orders{market: {longOrders, shortOrders}}

			cs.On("GetAcceptableBoundsForLiquidation", market).Return(liqUpperBound, liqLowerBound)
			cs.On("getMinAllowableMargin").Return(big.NewInt(1e5))
			cs.On("GetTakerFee").Return(big.NewInt(1e5))
			pipeline.runLiquidations(liquidablePositions, orderMap, underlyingPrices, map[common.Address]*big.Int{})
			assert.Equal(t, longOrders, orderMap[market].longOrders)
			assert.Equal(t, shortOrders, orderMap[market].shortOrders)
			lotp.AssertNotCalled(t, "ExecuteLiquidation", mock.Anything, mock.Anything, mock.Anything)
			cs.AssertCalled(t, "GetAcceptableBoundsForLiquidation", market)
		})
		t.Run("when short orders are present in database for matching", func(t *testing.T) {
			_, lotp, pipeline, underlyingPrices, cs := setupDependencies(t)
			longOrder := getLongOrder()
			shortOrder := getShortOrder()
			expectedFillAmount := utils.BigIntMinAbs(shortOrder.BaseAssetQuantity, liquidablePositions[0].Size)
			lotp.On("ExecuteLiquidation", traderAddress, shortOrder, expectedFillAmount).Return(nil)
			cs.On("GetAcceptableBoundsForLiquidation", market).Return(liqUpperBound, liqLowerBound)
			cs.On("getMinAllowableMargin").Return(big.NewInt(1e5))
			cs.On("GetTakerFee").Return(big.NewInt(1e5))

			orderMap := map[Market]*Orders{market: {[]Order{longOrder}, []Order{shortOrder}}}
			pipeline.runLiquidations(liquidablePositions, orderMap, underlyingPrices, map[common.Address]*big.Int{})

			lotp.AssertCalled(t, "ExecuteLiquidation", traderAddress, shortOrder, expectedFillAmount)
			cs.AssertCalled(t, "GetAcceptableBoundsForLiquidation", market)

			assert.Equal(t, longOrder, orderMap[market].longOrders[0])
			assert.Equal(t, expectedFillAmount.Uint64(), orderMap[market].shortOrders[0].FilledBaseAssetQuantity.Uint64())
		})
	})
}

func TestRunMatchingEngine(t *testing.T) {
	minAllowableMargin := big.NewInt(1e6)
	takerFee := big.NewInt(1e6)
	upperBound := big.NewInt(22)
	t.Run("when either long or short orders are not present in memorydb", func(t *testing.T) {
		t.Run("when no short and long orders are present", func(t *testing.T) {
			_, lotp, pipeline, _, _ := setupDependencies(t)
			longOrders := make([]Order, 0)
			shortOrders := make([]Order, 0)
			pipeline.runMatchingEngine(lotp, longOrders, shortOrders, map[common.Address]*big.Int{}, minAllowableMargin, takerFee, upperBound)
			lotp.AssertNotCalled(t, "ExecuteMatchedOrdersTx", mock.Anything, mock.Anything, mock.Anything)
		})
		t.Run("when longOrders are not present but short orders are present", func(t *testing.T) {
			_, lotp, pipeline, _, _ := setupDependencies(t)
			longOrders := make([]Order, 0)
			shortOrders := []Order{getShortOrder()}
			pipeline.runMatchingEngine(lotp, longOrders, shortOrders, map[common.Address]*big.Int{}, minAllowableMargin, takerFee, upperBound)
			lotp.AssertNotCalled(t, "ExecuteMatchedOrdersTx", mock.Anything, mock.Anything, mock.Anything)
		})
		t.Run("when short orders are not present but long orders are present", func(t *testing.T) {
			db, lotp, pipeline, _, _ := setupDependencies(t)
			longOrders := make([]Order, 0)
			shortOrders := make([]Order, 0)
			longOrder := getLongOrder()
			longOrders = append(longOrders, longOrder)
			db.On("GetLongOrders").Return(longOrders)
			db.On("GetShortOrders").Return(shortOrders)
			lotp.On("PurgeLocalTx").Return(nil)
			pipeline.runMatchingEngine(lotp, longOrders, shortOrders, map[common.Address]*big.Int{}, minAllowableMargin, takerFee, upperBound)
			lotp.AssertNotCalled(t, "ExecuteMatchedOrdersTx", mock.Anything, mock.Anything, mock.Anything)
		})
	})
	t.Run("When both long and short orders are present in db", func(t *testing.T) {
		t.Run("when longOrder.Price < shortOrder.Price", func(t *testing.T) {
			_, lotp, pipeline, _, _ := setupDependencies(t)
			shortOrder := getShortOrder()
			longOrder := getLongOrder()
			longOrder.Price.Sub(shortOrder.Price, big.NewInt(1))

			pipeline.runMatchingEngine(lotp, []Order{longOrder}, []Order{shortOrder}, map[common.Address]*big.Int{}, minAllowableMargin, takerFee, upperBound)
			lotp.AssertNotCalled(t, "ExecuteMatchedOrdersTx", mock.Anything, mock.Anything, mock.Anything)
		})
		t.Run("When longOrder.Price >= shortOrder.Price same", func(t *testing.T) {
			t.Run("When long order and short order's unfulfilled quantity is same", func(t *testing.T) {
				t.Run("When long order and short order's base asset quantity is same", func(t *testing.T) {
					//Add 2 long orders
					_, lotp, pipeline, _, _ := setupDependencies(t)
					longOrders := make([]Order, 0)
					longOrder1 := getLongOrder()
					longOrders = append(longOrders, longOrder1)
					longOrder2 := getLongOrder()
					longOrders = append(longOrders, longOrder2)

					// Add 2 short orders
					shortOrder1 := getShortOrder()
					shortOrders := make([]Order, 0)
					shortOrders = append(shortOrders, shortOrder1)
					shortOrder2 := getShortOrder()
					shortOrders = append(shortOrders, shortOrder2)

					fillAmount1 := longOrder1.BaseAssetQuantity
					fillAmount2 := longOrder2.BaseAssetQuantity
					marginMap := map[common.Address]*big.Int{
						common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa"): big.NewInt(1e9), // $1000
					}
					lotp.On("ExecuteMatchedOrdersTx", longOrder1, shortOrder1, fillAmount1).Return(nil)
					lotp.On("ExecuteMatchedOrdersTx", longOrder2, shortOrder2, fillAmount2).Return(nil)
					pipeline.runMatchingEngine(lotp, longOrders, shortOrders, marginMap, minAllowableMargin, takerFee, upperBound)
					lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder1, shortOrder1, fillAmount1)
					lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder2, shortOrder2, fillAmount2)
				})
				t.Run("When long order and short order's base asset quantity is different", func(t *testing.T) {
					db, lotp, pipeline, _, _ := setupDependencies(t)
					//Add 2 long orders with half base asset quantity of short order
					longOrders := make([]Order, 0)
					longOrder := getLongOrder()
					longOrder.BaseAssetQuantity = big.NewInt(20)
					longOrder.FilledBaseAssetQuantity = big.NewInt(5)
					longOrders = append(longOrders, longOrder)

					// Add 2 short orders
					shortOrder := getShortOrder()
					shortOrder.BaseAssetQuantity = big.NewInt(-30)
					shortOrder.FilledBaseAssetQuantity = big.NewInt(-15)
					shortOrders := make([]Order, 0)
					shortOrders = append(shortOrders, shortOrder)

					fillAmount := big.NewInt(0).Sub(longOrder.BaseAssetQuantity, longOrder.FilledBaseAssetQuantity)
					db.On("GetLongOrders").Return(longOrders)
					db.On("GetShortOrders").Return(shortOrders)
					lotp.On("PurgeLocalTx").Return(nil)
					lotp.On("ExecuteMatchedOrdersTx", longOrder, shortOrder, fillAmount).Return(nil)
					marginMap := map[common.Address]*big.Int{
						common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa"): big.NewInt(1e9), // $1000
					}
					pipeline.runMatchingEngine(lotp, longOrders, shortOrders, marginMap, minAllowableMargin, takerFee, upperBound)
					lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder, shortOrder, fillAmount)
				})
			})
			t.Run("When long order and short order's unfulfilled quantity is not same", func(t *testing.T) {
				db, lotp, pipeline, _, _ := setupDependencies(t)
				longOrders := make([]Order, 0)
				longOrder1 := getLongOrder()
				longOrder1.BaseAssetQuantity = big.NewInt(20)
				longOrder1.FilledBaseAssetQuantity = big.NewInt(5)
				longOrder2 := getLongOrder()
				longOrder2.BaseAssetQuantity = big.NewInt(40)
				longOrder2.FilledBaseAssetQuantity = big.NewInt(0)
				longOrder3 := getLongOrder()
				longOrder3.BaseAssetQuantity = big.NewInt(10)
				longOrder3.FilledBaseAssetQuantity = big.NewInt(3)
				longOrders = append(longOrders, longOrder1, longOrder2, longOrder3)

				// Add 2 short orders
				shortOrders := make([]Order, 0)
				shortOrder1 := getShortOrder()
				shortOrder1.BaseAssetQuantity = big.NewInt(-30)
				shortOrder1.FilledBaseAssetQuantity = big.NewInt(-2)
				shortOrder2 := getShortOrder()
				shortOrder2.BaseAssetQuantity = big.NewInt(-50)
				shortOrder2.FilledBaseAssetQuantity = big.NewInt(-20)
				shortOrder3 := getShortOrder()
				shortOrder3.BaseAssetQuantity = big.NewInt(-20)
				shortOrder3.FilledBaseAssetQuantity = big.NewInt(-10)
				shortOrders = append(shortOrders, shortOrder1, shortOrder2, shortOrder3)

				lotp.On("ExecuteMatchedOrdersTx", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(5)

				db.On("GetLongOrders").Return(longOrders)
				db.On("GetShortOrders").Return(shortOrders)
				lotp.On("PurgeLocalTx").Return(nil)
				log.Info("longOrder1", "longOrder1", longOrder1)
				marginMap := map[common.Address]*big.Int{
					common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa"): big.NewInt(1e9), // $1000
				}
				pipeline.runMatchingEngine(lotp, longOrders, shortOrders, marginMap, minAllowableMargin, takerFee, upperBound)
				log.Info("longOrder1", "longOrder1", longOrder1)

				//During 1st  matching iteration
				longOrder1UnfulfilledQuantity := longOrder1.GetUnFilledBaseAssetQuantity()
				shortOrder1UnfulfilledQuantity := shortOrder1.GetUnFilledBaseAssetQuantity()
				fillAmount := utils.BigIntMinAbs(longOrder1UnfulfilledQuantity, shortOrder1UnfulfilledQuantity)
				lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder1, shortOrder1, fillAmount)
				//After 1st matching iteration longOrder1 has been matched fully but shortOrder1 has not
				longOrder1.FilledBaseAssetQuantity.Add(longOrder1.FilledBaseAssetQuantity, fillAmount)
				shortOrder1.FilledBaseAssetQuantity.Sub(shortOrder1.FilledBaseAssetQuantity, fillAmount)

				//During 2nd iteration longOrder2 with be matched with shortOrder1
				longOrder2UnfulfilledQuantity := longOrder2.GetUnFilledBaseAssetQuantity()
				shortOrder1UnfulfilledQuantity = shortOrder1.GetUnFilledBaseAssetQuantity()
				fillAmount = utils.BigIntMinAbs(longOrder2UnfulfilledQuantity, shortOrder1UnfulfilledQuantity)
				lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder2, shortOrder1, fillAmount)
				//After 2nd matching iteration shortOrder1 has been matched fully but longOrder2 has not
				longOrder2.FilledBaseAssetQuantity.Add(longOrder2.FilledBaseAssetQuantity, fillAmount)
				shortOrder1.FilledBaseAssetQuantity.Sub(longOrder2.FilledBaseAssetQuantity, fillAmount)

				//During 3rd iteration longOrder2 with be matched with shortOrder2
				longOrder2UnfulfilledQuantity = longOrder2.GetUnFilledBaseAssetQuantity()
				shortOrder2UnfulfilledQuantity := shortOrder2.GetUnFilledBaseAssetQuantity()
				fillAmount = utils.BigIntMinAbs(longOrder2UnfulfilledQuantity, shortOrder2UnfulfilledQuantity)
				lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder2, shortOrder2, fillAmount)
				//After 3rd matching iteration longOrder2 has been matched fully but shortOrder2 has not
				longOrder2.FilledBaseAssetQuantity.Add(longOrder2.FilledBaseAssetQuantity, fillAmount)
				shortOrder2.FilledBaseAssetQuantity.Sub(shortOrder2.FilledBaseAssetQuantity, fillAmount)

				//So during 4th iteration longOrder3 with be matched with shortOrder2
				longOrder3UnfulfilledQuantity := longOrder3.GetUnFilledBaseAssetQuantity()
				shortOrder2UnfulfilledQuantity = shortOrder2.GetUnFilledBaseAssetQuantity()
				fillAmount = utils.BigIntMinAbs(longOrder3UnfulfilledQuantity, shortOrder2UnfulfilledQuantity)
				lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder3, shortOrder2, fillAmount)
				//After 4rd matching iteration shortOrder2 has been matched fully but longOrder3 has not
				longOrder3.FilledBaseAssetQuantity.Add(longOrder3.FilledBaseAssetQuantity, fillAmount)
				shortOrder2.FilledBaseAssetQuantity.Sub(shortOrder2.FilledBaseAssetQuantity, fillAmount)

				//So during 5th iteration longOrder3 with be matched with shortOrder3
				longOrder3UnfulfilledQuantity = longOrder3.GetUnFilledBaseAssetQuantity()
				shortOrder3UnfulfilledQuantity := shortOrder3.GetUnFilledBaseAssetQuantity()
				fillAmount = utils.BigIntMinAbs(longOrder3UnfulfilledQuantity, shortOrder3UnfulfilledQuantity)
				lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder3, shortOrder3, fillAmount)
			})
		})
	})
}

func TestMatchLongAndShortOrder(t *testing.T) {
	t.Run("When longPrice is less than shortPrice ,it returns orders unchanged and ordersMatched=false", func(t *testing.T) {
		_, lotp, _, _, _ := setupDependencies(t)
		longOrder := getLongOrder()
		shortOrder := getShortOrder()
		longOrder.Price.Sub(shortOrder.Price, big.NewInt(1))
		changedLongOrder, changedShortOrder, ordersMatched := matchLongAndShortOrder(lotp, longOrder, shortOrder)
		lotp.AssertNotCalled(t, "ExecuteMatchedOrdersTx", mock.Anything, mock.Anything, mock.Anything)
		assert.Equal(t, longOrder, changedLongOrder)
		assert.Equal(t, shortOrder, changedShortOrder)
		assert.Equal(t, false, ordersMatched)
	})
	t.Run("When longPrice is >= shortPrice", func(t *testing.T) {
		t.Run("When either longOrder or/and shortOrder is fully filled ", func(t *testing.T) {
			t.Run("When longOrder is fully filled, it returns orders unchanged and ordersMatched=false", func(t *testing.T) {
				_, lotp, _, _, _ := setupDependencies(t)
				longOrder := getLongOrder()
				longOrder.FilledBaseAssetQuantity = longOrder.BaseAssetQuantity
				shortOrder := getShortOrder()
				longOrder.Price.Add(shortOrder.Price, big.NewInt(1))
				changedLongOrder, changedShortOrder, ordersMatched := matchLongAndShortOrder(lotp, longOrder, shortOrder)
				lotp.AssertNotCalled(t, "ExecuteMatchedOrdersTx", mock.Anything, mock.Anything, mock.Anything)
				assert.Equal(t, longOrder, changedLongOrder)
				assert.Equal(t, shortOrder, changedShortOrder)
				assert.Equal(t, false, ordersMatched)
			})
			t.Run("When shortOrder is fully filled, it returns orders unchanged and ordersMatched=false", func(t *testing.T) {
				_, lotp, _, _, _ := setupDependencies(t)
				longOrder := getLongOrder()
				shortOrder := getShortOrder()
				longOrder.Price.Add(shortOrder.Price, big.NewInt(1))
				shortOrder.FilledBaseAssetQuantity = shortOrder.BaseAssetQuantity
				changedLongOrder, changedShortOrder, ordersMatched := matchLongAndShortOrder(lotp, longOrder, shortOrder)
				lotp.AssertNotCalled(t, "ExecuteMatchedOrdersTx", mock.Anything, mock.Anything, mock.Anything)
				assert.Equal(t, longOrder, changedLongOrder)
				assert.Equal(t, shortOrder, changedShortOrder)
				assert.Equal(t, false, ordersMatched)
			})
			t.Run("When longOrder and shortOrder are fully filled, it returns orders unchanged and ordersMatched=false", func(t *testing.T) {
				_, lotp, _, _, _ := setupDependencies(t)
				longOrder := getLongOrder()
				longOrder.FilledBaseAssetQuantity = longOrder.BaseAssetQuantity
				shortOrder := getShortOrder()
				longOrder.Price.Add(shortOrder.Price, big.NewInt(1))
				shortOrder.FilledBaseAssetQuantity = shortOrder.BaseAssetQuantity
				changedLongOrder, changedShortOrder, ordersMatched := matchLongAndShortOrder(lotp, longOrder, shortOrder)
				lotp.AssertNotCalled(t, "ExecuteMatchedOrdersTx", mock.Anything, mock.Anything, mock.Anything)
				assert.Equal(t, longOrder, changedLongOrder)
				assert.Equal(t, shortOrder, changedShortOrder)
				assert.Equal(t, false, ordersMatched)
			})
		})
		t.Run("when both long and short order are not fully filled", func(t *testing.T) {
			t.Run("when unfilled is same for longOrder and shortOrder", func(t *testing.T) {
				t.Run("When filled is zero for long and short order, it returns fully filled longOrder and shortOrder and ordersMatched=true", func(t *testing.T) {
					_, lotp, _, _, _ := setupDependencies(t)
					longOrder := getLongOrder()
					longOrder.FilledBaseAssetQuantity = big.NewInt(0)
					shortOrder := getShortOrder()
					longOrder.Price.Add(shortOrder.Price, big.NewInt(1))
					shortOrder.FilledBaseAssetQuantity = big.NewInt(0)
					shortOrder.BaseAssetQuantity = big.NewInt(0).Neg(longOrder.BaseAssetQuantity)

					expectedFillAmount := longOrder.BaseAssetQuantity
					lotp.On("ExecuteMatchedOrdersTx", longOrder, shortOrder, expectedFillAmount).Return(nil)

					changedLongOrder, changedShortOrder, ordersMatched := matchLongAndShortOrder(lotp, longOrder, shortOrder)
					lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder, shortOrder, expectedFillAmount)

					//setting this to test if returned order is same as original except for FilledBaseAssetQuantity
					longOrder.FilledBaseAssetQuantity = longOrder.BaseAssetQuantity
					shortOrder.FilledBaseAssetQuantity = shortOrder.BaseAssetQuantity
					assert.Equal(t, longOrder, changedLongOrder)
					assert.Equal(t, shortOrder, changedShortOrder)
					assert.Equal(t, true, ordersMatched)
				})
				t.Run("When filled is non zero for long and short order, it returns fully filled longOrder and shortOrder and ordersMatched=true", func(t *testing.T) {
					_, lotp, _, _, _ := setupDependencies(t)
					longOrder := getLongOrder()
					longOrder.BaseAssetQuantity = big.NewInt(20)
					longOrder.FilledBaseAssetQuantity = big.NewInt(5)
					shortOrder := getShortOrder()
					longOrder.Price.Add(shortOrder.Price, big.NewInt(1))
					shortOrder.BaseAssetQuantity = big.NewInt(-30)
					shortOrder.FilledBaseAssetQuantity = big.NewInt(-15)

					expectedFillAmount := big.NewInt(0).Sub(longOrder.BaseAssetQuantity, longOrder.FilledBaseAssetQuantity)
					lotp.On("ExecuteMatchedOrdersTx", longOrder, shortOrder, expectedFillAmount).Return(nil)
					changedLongOrder, changedShortOrder, ordersMatched := matchLongAndShortOrder(lotp, longOrder, shortOrder)

					lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder, shortOrder, expectedFillAmount)
					//setting this to test if returned order is same as original except for FilledBaseAssetQuantity
					longOrder.FilledBaseAssetQuantity = longOrder.BaseAssetQuantity
					shortOrder.FilledBaseAssetQuantity = shortOrder.BaseAssetQuantity
					assert.Equal(t, longOrder, changedLongOrder)
					assert.Equal(t, shortOrder, changedShortOrder)
					assert.Equal(t, true, ordersMatched)
				})
			})
			t.Run("when unfilled(amount x) is less for longOrder, it returns fully filled longOrder and adds fillAmount(x) to shortOrder with and ordersMatched=true", func(t *testing.T) {
				_, lotp, _, _, _ := setupDependencies(t)
				longOrder := getLongOrder()
				longOrder.BaseAssetQuantity = big.NewInt(20)
				longOrder.FilledBaseAssetQuantity = big.NewInt(15)
				shortOrder := getShortOrder()
				longOrder.Price.Add(shortOrder.Price, big.NewInt(1))
				shortOrder.BaseAssetQuantity = big.NewInt(-30)
				shortOrder.FilledBaseAssetQuantity = big.NewInt(-15)

				expectedFillAmount := big.NewInt(0).Sub(longOrder.BaseAssetQuantity, longOrder.FilledBaseAssetQuantity)
				lotp.On("ExecuteMatchedOrdersTx", longOrder, shortOrder, expectedFillAmount).Return(nil)
				changedLongOrder, changedShortOrder, ordersMatched := matchLongAndShortOrder(lotp, longOrder, shortOrder)
				lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder, shortOrder, expectedFillAmount)

				expectedShortOrderFilled := big.NewInt(0).Sub(shortOrder.FilledBaseAssetQuantity, expectedFillAmount)
				//setting this to test if returned order is same as original except for FilledBaseAssetQuantity
				longOrder.FilledBaseAssetQuantity = longOrder.BaseAssetQuantity
				shortOrder.FilledBaseAssetQuantity = expectedShortOrderFilled
				assert.Equal(t, longOrder, changedLongOrder)
				assert.Equal(t, shortOrder, changedShortOrder)
				assert.Equal(t, true, ordersMatched)
			})
			t.Run("when unfilled(amount x) is less for shortOrder, it returns fully filled shortOrder and adds fillAmount(x) to longOrder and ordersMatched=true", func(t *testing.T) {
				_, lotp, _, _, _ := setupDependencies(t)
				longOrder := getLongOrder()
				longOrder.BaseAssetQuantity = big.NewInt(20)
				longOrder.FilledBaseAssetQuantity = big.NewInt(5)
				shortOrder := getShortOrder()
				longOrder.Price.Add(shortOrder.Price, big.NewInt(1))
				shortOrder.BaseAssetQuantity = big.NewInt(-30)
				shortOrder.FilledBaseAssetQuantity = big.NewInt(-25)

				expectedFillAmount := big.NewInt(0).Neg(big.NewInt(0).Sub(shortOrder.BaseAssetQuantity, shortOrder.FilledBaseAssetQuantity))
				lotp.On("ExecuteMatchedOrdersTx", longOrder, shortOrder, expectedFillAmount).Return(nil)
				changedLongOrder, changedShortOrder, ordersMatched := matchLongAndShortOrder(lotp, longOrder, shortOrder)
				lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder, shortOrder, expectedFillAmount)

				expectedLongOrderFilled := big.NewInt(0).Add(longOrder.FilledBaseAssetQuantity, expectedFillAmount)
				//setting this to test if returned order is same as original except for FilledBaseAssetQuantity
				longOrder.FilledBaseAssetQuantity = expectedLongOrderFilled
				shortOrder.FilledBaseAssetQuantity = shortOrder.BaseAssetQuantity
				assert.Equal(t, longOrder, changedLongOrder)
				assert.Equal(t, shortOrder, changedShortOrder)
				assert.Equal(t, true, ordersMatched)
			})
		})
	})
}

func TestAreMatchingOrders(t *testing.T) {
	minAllowableMargin := big.NewInt(1e6)
	takerFee := big.NewInt(1e6)
	upperBound := big.NewInt(22)

	trader := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
	longOrder_ := Order{
		Market:                  1,
		PositionType:            LONG,
		BaseAssetQuantity:       big.NewInt(10),
		Trader:                  trader,
		FilledBaseAssetQuantity: big.NewInt(0),
		Salt:                    big.NewInt(1),
		Price:                   big.NewInt(100),
		ReduceOnly:              false,
		LifecycleList:           []Lifecycle{Lifecycle{}},
		BlockNumber:             big.NewInt(21),
		RawOrder: &LimitOrder{
			BaseOrder: hu.BaseOrder{
				AmmIndex:          big.NewInt(1),
				Trader:            trader,
				BaseAssetQuantity: big.NewInt(10),
				Price:             big.NewInt(100),
				Salt:              big.NewInt(1),
				ReduceOnly:        false,
			},
			PostOnly: false,
		},
		OrderType: Limit,
	}
	shortOrder_ := Order{
		Market:                  1,
		PositionType:            SHORT,
		BaseAssetQuantity:       big.NewInt(-10),
		Trader:                  trader,
		FilledBaseAssetQuantity: big.NewInt(0),
		Salt:                    big.NewInt(2),
		Price:                   big.NewInt(100),
		ReduceOnly:              false,
		LifecycleList:           []Lifecycle{Lifecycle{}},
		BlockNumber:             big.NewInt(21),
		RawOrder: &LimitOrder{
			BaseOrder: hu.BaseOrder{
				AmmIndex:          big.NewInt(1),
				Trader:            trader,
				BaseAssetQuantity: big.NewInt(-10),
				Price:             big.NewInt(100),
				Salt:              big.NewInt(2),
				ReduceOnly:        false,
			},
			PostOnly: false,
		},
		OrderType: Limit,
	}

	t.Run("longOrder's price < shortOrder's price", func(t *testing.T) {
		longOrder := deepCopyOrder(&longOrder_)
		shortOrder := deepCopyOrder(&shortOrder_)

		longOrder.Price = big.NewInt(80)
		marginMap := map[common.Address]*big.Int{
			common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa"): big.NewInt(1e9), // $1000
		}
		actualFillAmount := areMatchingOrders(longOrder, shortOrder, marginMap, minAllowableMargin, takerFee, upperBound)

		assert.Nil(t, actualFillAmount)
	})

	t.Run("longOrder was placed first", func(t *testing.T) {
		longOrder := deepCopyOrder(&longOrder_)
		shortOrder := deepCopyOrder(&shortOrder_)
		longOrder.BlockNumber = big.NewInt(20)
		shortOrder.BlockNumber = big.NewInt(21)
		t.Run("longOrder is IOC", func(t *testing.T) {
			longOrder.OrderType = IOC
			rawOrder := longOrder.RawOrder.(*LimitOrder)
			longOrder.RawOrder = &IOCOrder{
				BaseOrder: rawOrder.BaseOrder,
				OrderType: 1,
				ExpireAt:  big.NewInt(0),
			}
			marginMap := map[common.Address]*big.Int{
				common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa"): big.NewInt(1e9), // $1000
			}
			actualFillAmount := areMatchingOrders(longOrder, shortOrder, marginMap, minAllowableMargin, takerFee, upperBound)
			assert.Nil(t, actualFillAmount)
		})
		t.Run("short order is post only", func(t *testing.T) {
			longOrder := deepCopyOrder(&longOrder_)
			longOrder.BlockNumber = big.NewInt(20)

			shortOrder.RawOrder.(*LimitOrder).PostOnly = true
			marginMap := map[common.Address]*big.Int{
				common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa"): big.NewInt(1e9), // $1000
			}
			actualFillAmount := areMatchingOrders(longOrder, shortOrder, marginMap, minAllowableMargin, takerFee, upperBound)
			assert.Nil(t, actualFillAmount)
		})
	})

	t.Run("shortOrder was placed first", func(t *testing.T) {
		longOrder := deepCopyOrder(&longOrder_)
		shortOrder := deepCopyOrder(&shortOrder_)
		longOrder.BlockNumber = big.NewInt(21)
		shortOrder.BlockNumber = big.NewInt(20)
		t.Run("shortOrder is IOC", func(t *testing.T) {
			shortOrder.OrderType = IOC
			rawOrder := shortOrder.RawOrder.(*LimitOrder)
			shortOrder.RawOrder = &IOCOrder{
				BaseOrder: rawOrder.BaseOrder,
				OrderType: 1,
				ExpireAt:  big.NewInt(0),
			}
			marginMap := map[common.Address]*big.Int{
				common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa"): big.NewInt(1e9), // $1000
			}
			actualFillAmount := areMatchingOrders(longOrder, shortOrder, marginMap, minAllowableMargin, takerFee, upperBound)
			assert.Nil(t, actualFillAmount)
		})
		t.Run("longOrder is post only", func(t *testing.T) {
			longOrder := deepCopyOrder(&longOrder_)
			longOrder.BlockNumber = big.NewInt(21)

			longOrder.RawOrder.(*LimitOrder).PostOnly = true
			marginMap := map[common.Address]*big.Int{
				common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa"): big.NewInt(1e9), // $1000
			}
			actualFillAmount := areMatchingOrders(longOrder, shortOrder, marginMap, minAllowableMargin, takerFee, upperBound)
			assert.Nil(t, actualFillAmount)
		})
	})

	t.Run("one of the orders is fully filled", func(t *testing.T) {
		longOrder := deepCopyOrder(&longOrder_)
		shortOrder := deepCopyOrder(&shortOrder_)

		longOrder.FilledBaseAssetQuantity = longOrder.BaseAssetQuantity
		marginMap := map[common.Address]*big.Int{
			common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa"): big.NewInt(1e9), // $1000
		}
		actualFillAmount := areMatchingOrders(longOrder, shortOrder, marginMap, minAllowableMargin, takerFee, upperBound)
		assert.Nil(t, actualFillAmount)
	})

	t.Run("success", func(t *testing.T) {
		longOrder := deepCopyOrder(&longOrder_)
		shortOrder := deepCopyOrder(&shortOrder_)

		longOrder.FilledBaseAssetQuantity = big.NewInt(5)
		marginMap := map[common.Address]*big.Int{
			common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa"): big.NewInt(1e9), // $1000
		}
		actualFillAmount := areMatchingOrders(longOrder, shortOrder, marginMap, minAllowableMargin, takerFee, upperBound)
		assert.Equal(t, big.NewInt(5), actualFillAmount)
	})
}

func getShortOrder() Order {
	salt := big.NewInt(time.Now().Unix())
	shortOrder := createLimitOrder(SHORT, "0x22Bb736b64A0b4D4081E103f83bccF864F0404aa", big.NewInt(-10), big.NewInt(20.0), Placed, big.NewInt(2), salt)
	return shortOrder
}

func getLongOrder() Order {
	salt := big.NewInt(time.Now().Unix())
	longOrder := createLimitOrder(LONG, "0x22Bb736b64A0b4D4081E103f83bccF864F0404aa", big.NewInt(10), big.NewInt(20.0), Placed, big.NewInt(2), salt)
	return longOrder
}

func buildLongOrder(price, q int64) Order {
	salt := big.NewInt(time.Now().Unix())
	longOrder := createLimitOrder(LONG, "0x22Bb736b64A0b4D4081E103f83bccF864F0404aa", big.NewInt(q), big.NewInt(price), Placed, big.NewInt(2), salt)
	return longOrder
}

func buildShortOrder(price, q int64) Order {
	salt := big.NewInt(time.Now().Unix())
	order := createLimitOrder(SHORT, "0x22Bb736b64A0b4D4081E103f83bccF864F0404aa", big.NewInt(q), big.NewInt(price), Placed, big.NewInt(2), salt)
	return order
}

func getLiquidablePos(address common.Address, posType PositionType, size int64) LiquidablePosition {
	return LiquidablePosition{
		Address:      address,
		Market:       market,
		PositionType: posType,
		Size:         big.NewInt(size),
		FilledSize:   big.NewInt(0),
	}
}

func setupDependencies(t *testing.T) (*MockLimitOrderDatabase, *MockLimitOrderTxProcessor, *MatchingPipeline, map[Market]*big.Int, *MockConfigService) {
	db := NewMockLimitOrderDatabase()
	lotp := NewMockLimitOrderTxProcessor()
	cs := NewMockConfigService()
	pipeline := NewMatchingPipeline(db, lotp, cs)
	underlyingPrices := make(map[Market]*big.Int)
	underlyingPrices[market] = big.NewInt(20.0)
	return db, lotp, pipeline, underlyingPrices, cs
}
