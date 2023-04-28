package limitorders

import (
	"math/big"
	"testing"
	"time"

	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRunLiquidations(t *testing.T) {
	market := AvaxPerp
	t.Run("when there are no liquidable positions", func(t *testing.T) {
		db, lotp, pipeline := setupDependencies(t)
		traderMap := map[common.Address]Trader{}
		lastPrice := big.NewInt(20 * 10e6)
		longOrders := []LimitOrder{getLongOrder()}
		shortOrders := []LimitOrder{getShortOrder()}
		db.On("GetAllTraders").Return(traderMap)
		db.On("GetLastPrice").Return(lastPrice)
		modifiedLongOrders, modifiedShortOrders := pipeline.runLiquidations(market, longOrders, shortOrders)
		assert.Equal(t, longOrders, modifiedLongOrders)
		assert.Equal(t, shortOrders, modifiedShortOrders)
		lotp.AssertNotCalled(t, "ExecuteLiquidation", mock.Anything, mock.Anything, mock.Anything)
	})
	t.Run("when there are liquidable positions", func(t *testing.T) {
		t.Run("when liquidable position is long", func(t *testing.T) {
			longTraderAddress := common.HexToAddress("0x710bf5f942331874dcbc7783319123679033b63b")
			collateral := HUSD
			market := AvaxPerp
			marginLong := multiplyBasePrecision(big.NewInt(500))
			longSize := multiplyPrecisionSize(big.NewInt(10))
			longEntryPrice := multiplyBasePrecision(big.NewInt(300))
			openNotionalLong := dividePrecisionSize(big.NewInt(0).Mul(longEntryPrice, longSize))
			trader := Trader{
				Margin: Margin{
					Deposited: map[Collateral]*big.Int{
						collateral: marginLong,
					}},
				Positions: map[Market]*Position{
					market: getPosition(market, openNotionalLong, longSize, big.NewInt(0), big.NewInt(0), big.NewInt(0)),
				},
			}
			traderMap := map[common.Address]Trader{longTraderAddress: trader}
			lastPrice := big.NewInt(20 * 10e6)
			t.Run("when no long orders are present in database for matching", func(t *testing.T) {
				db, lotp, pipeline := setupDependencies(t)
				db.On("GetAllTraders").Return(traderMap)
				db.On("GetLastPrice").Return(lastPrice)
				longOrders := []LimitOrder{}
				shortOrders := []LimitOrder{getShortOrder()}
				modifiedLongOrders, modifiedShortOrders := pipeline.runLiquidations(market, longOrders, shortOrders)
				assert.Equal(t, longOrders, modifiedLongOrders)
				assert.Equal(t, shortOrders, modifiedShortOrders)
				lotp.AssertNotCalled(t, "ExecuteLiquidation", mock.Anything, mock.Anything, mock.Anything)
			})
			t.Run("when long orders are present in database for matching", func(t *testing.T) {
				db, lotp, pipeline := setupDependencies(t)
				db.On("GetAllTraders").Return(traderMap)
				db.On("GetLastPrice").Return(lastPrice)
				longOrder := getLongOrder()
				longOrders := []LimitOrder{longOrder}
				shortOrder := getShortOrder()
				shortOrders := []LimitOrder{shortOrder}
				expectedFillAmount := big.NewInt(0).Abs(longOrder.BaseAssetQuantity)
				lotp.On("ExecuteLiquidation", longTraderAddress, longOrder, expectedFillAmount).Return(nil)
				modifiedLongOrders, modifiedShortOrders := pipeline.runLiquidations(market, longOrders, shortOrders)

				lotp.AssertCalled(t, "ExecuteLiquidation", longTraderAddress, longOrder, expectedFillAmount)

				assert.Equal(t, shortOrder, modifiedShortOrders[0])
				assert.Equal(t, modifiedLongOrders[0].FilledBaseAssetQuantity.Uint64(), expectedFillAmount.Uint64())
			})
		})
		t.Run("when liquidable position is short", func(t *testing.T) {
			shortTraderAddress := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
			collateral := HUSD
			market := AvaxPerp
			marginShort := multiplyBasePrecision(big.NewInt(500))
			shortSize := multiplyPrecisionSize(big.NewInt(-20))
			shortEntryPrice := multiplyBasePrecision(big.NewInt(100))
			openNotionalShort := dividePrecisionSize(big.NewInt(0).Abs(big.NewInt(0).Mul(shortEntryPrice, shortSize)))
			trader := Trader{
				Margin: Margin{
					Deposited: map[Collateral]*big.Int{
						collateral: marginShort,
					}},
				Positions: map[Market]*Position{
					market: getPosition(market, openNotionalShort, shortSize, big.NewInt(0), big.NewInt(0), big.NewInt(0)),
				},
			}
			traderMap := map[common.Address]Trader{shortTraderAddress: trader}
			lastPrice := big.NewInt(20 * 10e6)
			t.Run("when no short orders are present in database for matching", func(t *testing.T) {
				db, lotp, pipeline := setupDependencies(t)
				db.On("GetAllTraders").Return(traderMap)
				db.On("GetLastPrice").Return(lastPrice)
				longOrders := []LimitOrder{getLongOrder()}
				shortOrders := []LimitOrder{}
				modifiedLongOrders, modifiedShortOrders := pipeline.runLiquidations(market, longOrders, shortOrders)
				assert.Equal(t, longOrders[0], modifiedLongOrders[0])
				assert.Equal(t, 0, len(modifiedShortOrders))
				lotp.AssertNotCalled(t, "ExecuteLiquidation", mock.Anything, mock.Anything, mock.Anything)
			})
			t.Run("when short orders are present in database for matching", func(t *testing.T) {
				db, lotp, pipeline := setupDependencies(t)
				db.On("GetAllTraders").Return(traderMap)
				db.On("GetLastPrice").Return(lastPrice)
				longOrder := getLongOrder()
				longOrders := []LimitOrder{longOrder}
				shortOrder := getShortOrder()
				shortOrders := []LimitOrder{shortOrder}
				expectedFillAmount := shortOrder.BaseAssetQuantity.Abs(shortOrder.BaseAssetQuantity)
				lotp.On("ExecuteLiquidation", shortTraderAddress, shortOrder, expectedFillAmount).Return(nil)
				modifiedLongOrders, modifiedShortOrders := pipeline.runLiquidations(market, longOrders, shortOrders)
				assert.Equal(t, longOrders[0], modifiedLongOrders[0])

				lotp.AssertCalled(t, "ExecuteLiquidation", shortTraderAddress, shortOrder, expectedFillAmount)

				assert.NotEqual(t, shortOrder, modifiedShortOrders[0])
				assert.Equal(t, expectedFillAmount, big.NewInt(0).Abs(shortOrders[0].FilledBaseAssetQuantity))
				shortOrder.FilledBaseAssetQuantity = expectedFillAmount
				assert.Equal(t, longOrders, modifiedLongOrders)
			})
		})
	})
}
func TestRunMatchingEngine(t *testing.T) {
	t.Run("when either long or short orders are not present in memorydb", func(t *testing.T) {
		t.Run("when no short and long orders are present", func(t *testing.T) {
			db, lotp, _ := setupDependencies(t)
			longOrders := make([]LimitOrder, 0)
			shortOrders := make([]LimitOrder, 0)
			db.On("GetLongOrders").Return(longOrders)
			db.On("GetShortOrders").Return(shortOrders)
			lotp.On("PurgeLocalTx").Return(nil)
			runMatchingEngine(lotp, longOrders, shortOrders)
			lotp.AssertNotCalled(t, "ExecuteMatchedOrdersTx", mock.Anything, mock.Anything, mock.Anything)
		})
		t.Run("when longOrders are not present but short orders are present", func(t *testing.T) {
			db, lotp, _ := setupDependencies(t)
			longOrders := make([]LimitOrder, 0)
			shortOrders := make([]LimitOrder, 0)
			shortOrders = append(shortOrders, getShortOrder())
			db.On("GetLongOrders").Return(longOrders)
			db.On("GetShortOrders").Return(shortOrders)
			lotp.On("PurgeLocalTx").Return(nil)
			runMatchingEngine(lotp, longOrders, shortOrders)
			lotp.AssertNotCalled(t, "ExecuteMatchedOrdersTx", mock.Anything, mock.Anything, mock.Anything)
		})
		t.Run("when short orders are not present but long orders are present", func(t *testing.T) {
			db, lotp, _ := setupDependencies(t)
			longOrders := make([]LimitOrder, 0)
			shortOrders := make([]LimitOrder, 0)
			longOrder := getLongOrder()
			longOrders = append(longOrders, longOrder)
			db.On("GetLongOrders").Return(longOrders)
			db.On("GetShortOrders").Return(shortOrders)
			lotp.On("PurgeLocalTx").Return(nil)
			runMatchingEngine(lotp, longOrders, shortOrders)
			lotp.AssertNotCalled(t, "ExecuteMatchedOrdersTx", mock.Anything, mock.Anything, mock.Anything)
		})
	})
	t.Run("When both long and short orders are present in db", func(t *testing.T) {
		t.Run("when longOrder.Price < shortOrder.Price", func(t *testing.T) {
			db, lotp, _ := setupDependencies(t)
			longOrders := make([]LimitOrder, 0)
			shortOrders := make([]LimitOrder, 0)
			longOrder := getLongOrder()
			longOrders = append(longOrders, longOrder)
			shortOrder := getShortOrder()
			shortOrder.Price = big.NewInt(0).Add(shortOrder.Price, big.NewInt(2))
			shortOrders = append(shortOrders, shortOrder)
			db.On("GetLongOrders").Return(longOrders)
			db.On("GetShortOrders").Return(shortOrders)
			lotp.On("PurgeLocalTx").Return(nil)
			runMatchingEngine(lotp, longOrders, shortOrders)
			lotp.AssertNotCalled(t, "ExecuteMatchedOrdersTx", mock.Anything, mock.Anything, mock.Anything)
		})
		t.Run("When longOrder.Price >= shortOrder.Price same", func(t *testing.T) {
			t.Run("When long order and short order's unfulfilled quantity is same", func(t *testing.T) {
				t.Run("When long order and short order's base asset quantity is same", func(t *testing.T) {
					//Add 2 long orders
					db, lotp, _ := setupDependencies(t)
					longOrders := make([]LimitOrder, 0)
					longOrder1 := getLongOrder()
					longOrders = append(longOrders, longOrder1)
					longOrder2 := getLongOrder()
					longOrder2.Signature = []byte("Here is a 2nd long order")
					longOrders = append(longOrders, longOrder2)

					// Add 2 short orders
					shortOrder1 := getShortOrder()
					shortOrders := make([]LimitOrder, 0)
					shortOrders = append(shortOrders, shortOrder1)
					shortOrder2 := getShortOrder()
					shortOrder2.Signature = []byte("Here is a 2nd short order")
					shortOrders = append(shortOrders, shortOrder2)

					db.On("GetLongOrders").Return(longOrders)
					db.On("GetShortOrders").Return(shortOrders)
					lotp.On("PurgeLocalTx").Return(nil)
					fillAmount1 := longOrder1.BaseAssetQuantity
					fillAmount2 := longOrder2.BaseAssetQuantity
					lotp.On("ExecuteMatchedOrdersTx", longOrder1, shortOrder1, fillAmount1).Return(nil)
					lotp.On("ExecuteMatchedOrdersTx", longOrder2, shortOrder2, fillAmount2).Return(nil)
					runMatchingEngine(lotp, longOrders, shortOrders)
					lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder1, shortOrder1, fillAmount1)
					lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder2, shortOrder2, fillAmount2)
				})
				t.Run("When long order and short order's base asset quantity is different", func(t *testing.T) {
					db, lotp, _ := setupDependencies(t)
					//Add 2 long orders with half base asset quantity of short order
					longOrders := make([]LimitOrder, 0)
					longOrder := getLongOrder()
					longOrder.BaseAssetQuantity = big.NewInt(20)
					longOrder.FilledBaseAssetQuantity = big.NewInt(5)
					longOrders = append(longOrders, longOrder)

					// Add 2 short orders
					shortOrder := getShortOrder()
					shortOrder.BaseAssetQuantity = big.NewInt(-30)
					shortOrder.FilledBaseAssetQuantity = big.NewInt(-15)
					shortOrders := make([]LimitOrder, 0)
					shortOrders = append(shortOrders, shortOrder)

					fillAmount := big.NewInt(0).Sub(longOrder.BaseAssetQuantity, longOrder.FilledBaseAssetQuantity)
					db.On("GetLongOrders").Return(longOrders)
					db.On("GetShortOrders").Return(shortOrders)
					lotp.On("PurgeLocalTx").Return(nil)
					lotp.On("ExecuteMatchedOrdersTx", longOrder, shortOrder, fillAmount).Return(nil)
					runMatchingEngine(lotp, longOrders, shortOrders)
					lotp.AssertCalled(t, "ExecuteMatchedOrdersTx", longOrder, shortOrder, fillAmount)
				})
			})
			t.Run("When long order and short order's unfulfilled quantity is not same", func(t *testing.T) {
				db, lotp, _ := setupDependencies(t)
				longOrders := make([]LimitOrder, 0)
				longOrder1 := getLongOrder()
				longOrder1.BaseAssetQuantity = big.NewInt(20)
				longOrder1.FilledBaseAssetQuantity = big.NewInt(5)
				longOrder2 := getLongOrder()
				longOrder2.BaseAssetQuantity = big.NewInt(40)
				longOrder2.FilledBaseAssetQuantity = big.NewInt(0)
				longOrder2.Signature = []byte("Here is a 2nd long order")
				longOrder3 := getLongOrder()
				longOrder3.BaseAssetQuantity = big.NewInt(10)
				longOrder3.FilledBaseAssetQuantity = big.NewInt(3)
				longOrder3.Signature = []byte("Here is a 3rd long order")
				longOrders = append(longOrders, longOrder1, longOrder2, longOrder3)

				// Add 2 short orders
				shortOrders := make([]LimitOrder, 0)
				shortOrder1 := getShortOrder()
				shortOrder1.BaseAssetQuantity = big.NewInt(-30)
				shortOrder1.FilledBaseAssetQuantity = big.NewInt(-2)
				shortOrder2 := getShortOrder()
				shortOrder2.BaseAssetQuantity = big.NewInt(-50)
				shortOrder2.FilledBaseAssetQuantity = big.NewInt(-20)
				shortOrder2.Signature = []byte("Here is a 2nd short order")
				shortOrder3 := getShortOrder()
				shortOrder3.BaseAssetQuantity = big.NewInt(-20)
				shortOrder3.FilledBaseAssetQuantity = big.NewInt(-10)
				shortOrder3.Signature = []byte("Here is a 3rd short order")
				shortOrders = append(shortOrders, shortOrder1, shortOrder2, shortOrder3)

				lotp.On("ExecuteMatchedOrdersTx", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(5)

				db.On("GetLongOrders").Return(longOrders)
				db.On("GetShortOrders").Return(shortOrders)
				lotp.On("PurgeLocalTx").Return(nil)
				runMatchingEngine(lotp, longOrders, shortOrders)

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
		_, lotp, _ := setupDependencies(t)
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
				_, lotp, _ := setupDependencies(t)
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
				_, lotp, _ := setupDependencies(t)
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
				_, lotp, _ := setupDependencies(t)
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
					_, lotp, _ := setupDependencies(t)
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
					_, lotp, _ := setupDependencies(t)
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
				_, lotp, _ := setupDependencies(t)
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
				_, lotp, _ := setupDependencies(t)
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

func getShortOrder() LimitOrder {
	signature := []byte("Here is a short order")
	salt := big.NewInt(time.Now().Unix())
	shortOrder, _ := createLimitOrder("short", "0x22Bb736b64A0b4D4081E103f83bccF864F0404aa", big.NewInt(-10), big.NewInt(20.0), Placed, signature, big.NewInt(2), salt)
	return shortOrder
}

func getLongOrder() LimitOrder {
	signature := []byte("Here is a long order")
	salt := big.NewInt(time.Now().Unix())
	longOrder, _ := createLimitOrder("long", "0x22Bb736b64A0b4D4081E103f83bccF864F0404aa", big.NewInt(10), big.NewInt(20.0), Placed, signature, big.NewInt(2), salt)
	return longOrder
}

func setupDependencies(t *testing.T) (*MockLimitOrderDatabase, *MockLimitOrderTxProcessor, *BuildBlockPipeline) {
	db := NewMockLimitOrderDatabase()
	lotp := NewMockLimitOrderTxProcessor()
	pipeline := NewBuildBlockPipeline(db, lotp)
	return db, lotp, pipeline
}
