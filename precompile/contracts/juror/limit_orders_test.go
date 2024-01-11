package juror

import (
	"encoding/hex"
	"math/big"

	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	hu "github.com/ava-labs/subnet-evm/plugin/evm/orderbook/hubbleutils"
	b "github.com/ava-labs/subnet-evm/precompile/contracts/bibliophile"
	gomock "github.com/golang/mock/gomock"
)

func TestValidatePlaceLimitOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBibliophile := b.NewMockBibliophileClient(ctrl)
	ammIndex := big.NewInt(0)
	longBaseAssetQuantity := big.NewInt(5000000000000000000)
	shortBaseAssetQuantity := big.NewInt(-5000000000000000000)
	price := big.NewInt(100000000)
	salt := big.NewInt(121)
	reduceOnly := false
	postOnly := false
	trader := common.HexToAddress("0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC")
	ammAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")

	t.Run("Basic Order checks", func(t *testing.T) {
		t.Run("when baseAssetQuantity is 0", func(t *testing.T) {
			newBaseAssetQuantity := big.NewInt(0)
			order := getOrder(ammIndex, trader, newBaseAssetQuantity, price, salt, reduceOnly, postOnly)

			mockBibliophile.EXPECT().GetMarketAddressFromMarketID(order.AmmIndex.Int64()).Return(ammAddress).Times(1)
			output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: order, Sender: trader})
			assert.Equal(t, ErrBaseAssetQuantityZero.Error(), output.Err)
			expectedOrderHash, _ := GetLimitOrderHashFromContractStruct(&order)
			assert.Equal(t, common.BytesToHash(output.Orderhash[:]), expectedOrderHash)
			assert.Equal(t, output.Res.Amm, ammAddress)
			assert.Equal(t, output.Res.ReserveAmount, big.NewInt(0))
		})
		t.Run("when baseAssetQuantity is not 0", func(t *testing.T) {
			t.Run("when sender is not the trader and is not trading authority, it returns error", func(t *testing.T) {
				sender := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C9")
				t.Run("it returns error for a long order", func(t *testing.T) {
					order := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, reduceOnly, postOnly)
					mockBibliophile.EXPECT().IsTradingAuthority(order.Trader, sender).Return(false).Times(1)
					output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: order, Sender: sender})
					assert.Equal(t, "de9b5c2bf047cda53602c6a3223cd4b84b2b659f2ad6bc4b3fb29aed156185bd", hex.EncodeToString(output.Orderhash[:]))
					assert.Equal(t, ErrNoTradingAuthority.Error(), output.Err)
				})
				t.Run("it returns error for a short order", func(t *testing.T) {
					order := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, reduceOnly, postOnly)
					mockBibliophile.EXPECT().IsTradingAuthority(order.Trader, sender).Return(false).Times(1)
					output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: order, Sender: sender})
					// fmt.Println("Orderhash", hex.EncodeToString(output.Orderhash[:]))
					assert.Equal(t, "8c9158cccd9795896fef87cc969deb425499f230ae9a4427d314f89ac76a0288", hex.EncodeToString(output.Orderhash[:]))
					assert.Equal(t, ErrNoTradingAuthority.Error(), output.Err)
				})
			})
			t.Run("when either sender is trader or a trading authority", func(t *testing.T) {
				t.Run("when baseAssetQuantity is not a multiple of minSizeRequirement", func(t *testing.T) {
					t.Run("when |baseAssetQuantity| is >0 but less than minSizeRequirement", func(t *testing.T) {
						t.Run("it returns error for a long Order", func(t *testing.T) {
							minSizeRequirement := big.NewInt(0).Add(longBaseAssetQuantity, big.NewInt(1))
							order := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, reduceOnly, postOnly)

							mockBibliophile.EXPECT().GetMarketAddressFromMarketID(order.AmmIndex.Int64()).Return(ammAddress).Times(1)
							mockBibliophile.EXPECT().GetMinSizeRequirement(order.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
							output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: order, Sender: trader})
							assert.Equal(t, ErrNotMultiple.Error(), output.Err)
							expectedOrderHash, _ := GetLimitOrderHashFromContractStruct(&order)
							assert.Equal(t, common.BytesToHash(output.Orderhash[:]), expectedOrderHash)
							assert.Equal(t, output.Res.Amm, ammAddress)
							assert.Equal(t, output.Res.ReserveAmount, big.NewInt(0))
						})
						t.Run("it returns error for a short Order", func(t *testing.T) {
							minSizeRequirement := big.NewInt(0).Sub(shortBaseAssetQuantity, big.NewInt(1))
							order := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, reduceOnly, postOnly)

							mockBibliophile.EXPECT().GetMarketAddressFromMarketID(order.AmmIndex.Int64()).Return(ammAddress).Times(1)
							mockBibliophile.EXPECT().GetMinSizeRequirement(order.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
							output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: order, Sender: trader})
							assert.Equal(t, ErrNotMultiple.Error(), output.Err)
							expectedOrderHash, _ := GetLimitOrderHashFromContractStruct(&order)
							assert.Equal(t, expectedOrderHash, common.BytesToHash(output.Orderhash[:]))
							assert.Equal(t, ammAddress, output.Res.Amm)
							assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
						})
					})
					t.Run("when |baseAssetQuantity| is > minSizeRequirement but not a multiple of minSizeRequirement", func(t *testing.T) {
						t.Run("it returns error for a long Order", func(t *testing.T) {
							minSizeRequirement := big.NewInt(0).Div(big.NewInt(0).Mul(longBaseAssetQuantity, big.NewInt(3)), big.NewInt(2))
							order := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, reduceOnly, postOnly)

							mockBibliophile.EXPECT().GetMarketAddressFromMarketID(order.AmmIndex.Int64()).Return(ammAddress).Times(1)
							mockBibliophile.EXPECT().GetMinSizeRequirement(order.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
							output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: order, Sender: trader})
							assert.Equal(t, ErrNotMultiple.Error(), output.Err)
							expectedOrderHash, _ := GetLimitOrderHashFromContractStruct(&order)
							assert.Equal(t, expectedOrderHash, common.BytesToHash(output.Orderhash[:]))
							assert.Equal(t, ammAddress, output.Res.Amm)
							assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
						})
						t.Run("it returns error for a short Order", func(t *testing.T) {
							minSizeRequirement := big.NewInt(0).Div(big.NewInt(0).Mul(longBaseAssetQuantity, big.NewInt(3)), big.NewInt(2))
							order := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, reduceOnly, postOnly)

							mockBibliophile.EXPECT().GetMarketAddressFromMarketID(order.AmmIndex.Int64()).Return(ammAddress).Times(1)
							mockBibliophile.EXPECT().GetMinSizeRequirement(order.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
							output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: order, Sender: trader})
							assert.Equal(t, ErrNotMultiple.Error(), output.Err)
							expectedOrderHash, _ := GetLimitOrderHashFromContractStruct(&order)
							assert.Equal(t, expectedOrderHash, common.BytesToHash(output.Orderhash[:]))
							assert.Equal(t, ammAddress, output.Res.Amm)
							assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
						})
					})
				})
				t.Run("when baseAssetQuantity is a multiple of minSizeRequirement", func(t *testing.T) {
					minSizeRequirement := big.NewInt(0).Div(longBaseAssetQuantity, big.NewInt(2))

					t.Run("when order was placed earlier", func(t *testing.T) {
						t.Run("when order status is placed", func(t *testing.T) {
							t.Run("it returns error for a longOrder", func(t *testing.T) {
								longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, reduceOnly, postOnly)

								mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
								mockBibliophile.EXPECT().GetMinSizeRequirement(longOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
								orderHash, err := GetLimitOrderHashFromContractStruct(&longOrder)
								if err != nil {
									panic("error in getting longOrder hash")
								}
								mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Placed)).Times(1)
								output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: longOrder, Sender: trader})
								assert.Equal(t, ErrOrderAlreadyExists.Error(), output.Err)
								assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
								assert.Equal(t, ammAddress, output.Res.Amm)
								assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
							})
							t.Run("it returns error for a shortOrder", func(t *testing.T) {
								shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, reduceOnly, postOnly)

								mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
								mockBibliophile.EXPECT().GetMinSizeRequirement(shortOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
								orderHash, err := GetLimitOrderHashFromContractStruct(&shortOrder)
								if err != nil {
									panic("error in getting shortOrder hash")
								}
								mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Placed)).Times(1)
								output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: shortOrder, Sender: trader})
								assert.Equal(t, ErrOrderAlreadyExists.Error(), output.Err)
								assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
								assert.Equal(t, ammAddress, output.Res.Amm)
								assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
							})
						})
						t.Run("when order status is filled", func(t *testing.T) {
							t.Run("it returns error for a longOrder", func(t *testing.T) {
								longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, reduceOnly, postOnly)

								mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
								mockBibliophile.EXPECT().GetMinSizeRequirement(longOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
								orderHash, err := GetLimitOrderHashFromContractStruct(&longOrder)
								if err != nil {
									panic("error in getting longOrder hash")
								}
								mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Filled)).Times(1)
								output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: longOrder, Sender: trader})
								assert.Equal(t, ErrOrderAlreadyExists.Error(), output.Err)
								assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
								assert.Equal(t, ammAddress, output.Res.Amm)
								assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
							})
							t.Run("it returns error for a shortOrder", func(t *testing.T) {
								shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, reduceOnly, postOnly)

								mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
								mockBibliophile.EXPECT().GetMinSizeRequirement(shortOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
								orderHash, err := GetLimitOrderHashFromContractStruct(&shortOrder)
								if err != nil {
									panic("error in getting shortOrder hash")
								}
								mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Filled)).Times(1)
								output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: shortOrder, Sender: trader})
								assert.Equal(t, ErrOrderAlreadyExists.Error(), output.Err)
								assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
								assert.Equal(t, ammAddress, output.Res.Amm)
								assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
							})
						})
						t.Run("when order status is cancelled", func(t *testing.T) {
							t.Run("it returns error for a longOrder", func(t *testing.T) {
								longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, reduceOnly, postOnly)

								mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
								mockBibliophile.EXPECT().GetMinSizeRequirement(longOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
								orderHash, err := GetLimitOrderHashFromContractStruct(&longOrder)
								if err != nil {
									panic("error in getting longOrder hash")
								}
								mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Cancelled)).Times(1)
								output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: longOrder, Sender: trader})
								assert.Equal(t, ErrOrderAlreadyExists.Error(), output.Err)
								assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
								assert.Equal(t, ammAddress, output.Res.Amm)
								assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
							})
							t.Run("it returns error for a shortOrder", func(t *testing.T) {
								shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, reduceOnly, postOnly)

								mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
								mockBibliophile.EXPECT().GetMinSizeRequirement(shortOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
								orderHash, err := GetLimitOrderHashFromContractStruct(&shortOrder)
								if err != nil {
									panic("error in getting shortOrder hash")
								}
								mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Cancelled)).Times(1)
								output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: shortOrder, Sender: trader})
								assert.Equal(t, ErrOrderAlreadyExists.Error(), output.Err)
								assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
								assert.Equal(t, ammAddress, output.Res.Amm)
								assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
							})
						})
					})
				})
			})
		})
	})
	t.Run("When basic order validations pass", func(t *testing.T) {
		minSizeRequirement := big.NewInt(0).Div(longBaseAssetQuantity, big.NewInt(2))
		t.Run("When order is reduceOnly order", func(t *testing.T) {
			t.Run("When reduceOnly does not reduce position", func(t *testing.T) {
				t.Run("when trader has longPosition", func(t *testing.T) {
					t.Run("it returns error when order is longOrder", func(t *testing.T) {
						positionSize := longBaseAssetQuantity
						longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, true, postOnly)

						mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
						mockBibliophile.EXPECT().GetMinSizeRequirement(longOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
						orderHash, err := GetLimitOrderHashFromContractStruct(&longOrder)
						if err != nil {
							panic("error in getting longOrder hash")
						}
						mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
						mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
						mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, longOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
						output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: longOrder, Sender: trader})
						assert.Equal(t, ErrReduceOnlyBaseAssetQuantityInvalid.Error(), output.Err)
						assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
						assert.Equal(t, ammAddress, output.Res.Amm)
						assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
					})
					t.Run("it returns error when order is shortOrder and |baseAssetQuantity| > |positionSize|", func(t *testing.T) {
						positionSize := big.NewInt(0).Abs(big.NewInt(0).Add(shortBaseAssetQuantity, big.NewInt(1)))
						shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, true, postOnly)

						mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
						mockBibliophile.EXPECT().GetMinSizeRequirement(shortOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
						orderHash, err := GetLimitOrderHashFromContractStruct(&shortOrder)
						if err != nil {
							panic("error in getting shortOrder hash")
						}
						mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
						mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
						mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, shortOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
						output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: shortOrder, Sender: trader})
						assert.Equal(t, ErrReduceOnlyBaseAssetQuantityInvalid.Error(), output.Err)
						assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
						assert.Equal(t, ammAddress, output.Res.Amm)
						assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
					})
				})
				t.Run("when trader has shortPosition", func(t *testing.T) {
					t.Run("it returns when order is shortOrder", func(t *testing.T) {
						positionSize := shortBaseAssetQuantity
						shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, true, postOnly)

						mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
						mockBibliophile.EXPECT().GetMinSizeRequirement(shortOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
						orderHash, err := GetLimitOrderHashFromContractStruct(&shortOrder)
						if err != nil {
							panic("error in getting shortOrder hash")
						}
						mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
						mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
						mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, shortOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
						output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: shortOrder, Sender: trader})
						assert.Equal(t, ErrReduceOnlyBaseAssetQuantityInvalid.Error(), output.Err)
						assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
						assert.Equal(t, ammAddress, output.Res.Amm)
						assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
					})
					t.Run("it returns error when order is longOrder and |baseAssetQuantity| > |positionSize|", func(t *testing.T) {
						positionSize := big.NewInt(0).Sub(longBaseAssetQuantity, big.NewInt(1))
						longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, true, postOnly)

						mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
						mockBibliophile.EXPECT().GetMinSizeRequirement(longOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
						orderHash, err := GetLimitOrderHashFromContractStruct(&longOrder)
						if err != nil {
							panic("error in getting longOrder hash")
						}
						mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
						mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
						mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, longOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
						output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: longOrder, Sender: trader})
						assert.Equal(t, ErrReduceOnlyBaseAssetQuantityInvalid.Error(), output.Err)
						assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
						assert.Equal(t, ammAddress, output.Res.Amm)
						assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
					})
				})
			})
			t.Run("When reduceOnly reduces position", func(t *testing.T) {
				t.Run("when there are non reduceOnly Orders in same direction", func(t *testing.T) {
					t.Run("for a short position", func(t *testing.T) {
						t.Run("it returns error if order is longOrder and there are open longOrders which are not reduceOnly", func(t *testing.T) {
							positionSize := shortBaseAssetQuantity
							longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, true, postOnly)

							longOpenOrdersAmount := big.NewInt(0).Div(positionSize, big.NewInt(4))

							mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
							mockBibliophile.EXPECT().GetMinSizeRequirement(longOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
							orderHash, err := GetLimitOrderHashFromContractStruct(&longOrder)
							if err != nil {
								panic("error in getting longOrder hash")
							}
							mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
							mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
							mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, longOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
							mockBibliophile.EXPECT().GetLongOpenOrdersAmount(trader, longOrder.AmmIndex).Return(longOpenOrdersAmount).Times(1)
							mockBibliophile.EXPECT().GetShortOpenOrdersAmount(trader, longOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
							output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: longOrder, Sender: trader})
							assert.Equal(t, ErrOpenOrders.Error(), output.Err)
							assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
							assert.Equal(t, ammAddress, output.Res.Amm)
							assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
						})
					})
					t.Run("for a long position", func(t *testing.T) {
						t.Run("it returns error if order is shortOrder and there are open shortOrders which are not reduceOnly", func(t *testing.T) {
							positionSize := longBaseAssetQuantity
							shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, true, postOnly)

							shortOpenOrdersAmount := big.NewInt(0).Div(longBaseAssetQuantity, big.NewInt(4))

							mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
							mockBibliophile.EXPECT().GetMinSizeRequirement(shortOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
							orderHash, err := GetLimitOrderHashFromContractStruct(&shortOrder)
							if err != nil {
								panic("error in getting shortOrder hash")
							}
							mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
							mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
							mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, shortOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
							mockBibliophile.EXPECT().GetLongOpenOrdersAmount(trader, shortOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
							mockBibliophile.EXPECT().GetShortOpenOrdersAmount(trader, shortOrder.AmmIndex).Return(shortOpenOrdersAmount).Times(1)
							output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: shortOrder, Sender: trader})
							assert.Equal(t, ErrOpenOrders.Error(), output.Err)
							assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
							assert.Equal(t, ammAddress, output.Res.Amm)
							assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
						})
					})
				})
				t.Run("when there are no non reduceOnly orders in same direction", func(t *testing.T) {
					t.Run("when current open reduceOnlyOrders plus currentOrder's baseAssetQuantity exceeds positionSize", func(t *testing.T) {
						t.Run("it returns error for a longOrder", func(t *testing.T) {
							positionSize := shortBaseAssetQuantity
							longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, true, postOnly)

							reduceOnlyAmount := big.NewInt(1)

							mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
							mockBibliophile.EXPECT().GetMinSizeRequirement(longOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
							orderHash, err := GetLimitOrderHashFromContractStruct(&longOrder)
							if err != nil {
								panic("error in getting longOrder hash")
							}
							mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
							mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
							mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, longOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
							mockBibliophile.EXPECT().GetLongOpenOrdersAmount(trader, longOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
							mockBibliophile.EXPECT().GetShortOpenOrdersAmount(trader, longOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
							output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: longOrder, Sender: trader})
							assert.Equal(t, ErrNetReduceOnlyAmountExceeded.Error(), output.Err)
							assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
							assert.Equal(t, ammAddress, output.Res.Amm)
							assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
						})
						t.Run("it returns error for a shortOrder", func(t *testing.T) {
							positionSize := longBaseAssetQuantity
							shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, true, postOnly)

							reduceOnlyAmount := big.NewInt(-1)

							mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
							mockBibliophile.EXPECT().GetMinSizeRequirement(shortOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
							orderHash, err := GetLimitOrderHashFromContractStruct(&shortOrder)
							if err != nil {
								panic("error in getting shortOrder hash")
							}
							mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
							mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
							mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, shortOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
							mockBibliophile.EXPECT().GetLongOpenOrdersAmount(trader, shortOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
							mockBibliophile.EXPECT().GetShortOpenOrdersAmount(trader, shortOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
							output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: shortOrder, Sender: trader})
							assert.Equal(t, ErrNetReduceOnlyAmountExceeded.Error(), output.Err)
							assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
							assert.Equal(t, ammAddress, output.Res.Amm)
							assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
						})
					})
					t.Run("when current open reduceOnlyOrders plus currentOrder's baseAssetQuantity <= positionSize", func(t *testing.T) {
						t.Run("when order is not postOnly order", func(t *testing.T) {
							t.Run("for a longOrder it returns no error and 0 as reserveAmount", func(t *testing.T) {
								positionSize := big.NewInt(0).Mul(longBaseAssetQuantity, big.NewInt(-2))
								longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, true, postOnly)

								reduceOnlyAmount := big.NewInt(1)

								mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
								mockBibliophile.EXPECT().GetMinSizeRequirement(longOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
								orderHash, err := GetLimitOrderHashFromContractStruct(&longOrder)
								if err != nil {
									panic("error in getting longOrder hash")
								}
								mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
								mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
								mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, longOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
								mockBibliophile.EXPECT().GetLongOpenOrdersAmount(trader, longOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
								mockBibliophile.EXPECT().GetShortOpenOrdersAmount(trader, longOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
								mockBibliophile.EXPECT().HasReferrer(trader).Return(true).Times(1)
								mockBibliophile.EXPECT().GetPriceMultiplier(ammAddress).Return(big.NewInt(1)).Times(1)
								output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: longOrder, Sender: trader})
								assert.Equal(t, "", output.Err)
								assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
								assert.Equal(t, ammAddress, output.Res.Amm)
								assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
							})
							t.Run("for a shortOrder it returns no error and 0 as reserveAmount", func(t *testing.T) {
								positionSize := big.NewInt(0).Mul(shortBaseAssetQuantity, big.NewInt(-2))
								shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, true, postOnly)

								reduceOnlyAmount := big.NewInt(-1)

								mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
								mockBibliophile.EXPECT().GetMinSizeRequirement(shortOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
								orderHash, err := GetLimitOrderHashFromContractStruct(&shortOrder)
								if err != nil {
									panic("error in getting shortOrder hash")
								}
								mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
								mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
								mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, shortOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
								mockBibliophile.EXPECT().GetLongOpenOrdersAmount(trader, shortOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
								mockBibliophile.EXPECT().GetShortOpenOrdersAmount(trader, shortOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
								mockBibliophile.EXPECT().HasReferrer(trader).Return(true).Times(1)
								mockBibliophile.EXPECT().GetPriceMultiplier(ammAddress).Return(big.NewInt(1)).Times(1)
								output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: shortOrder, Sender: trader})
								assert.Equal(t, "", output.Err)
								assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
								assert.Equal(t, ammAddress, output.Res.Amm)
								assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
							})
						})
						t.Run("when order is postOnly order", func(t *testing.T) {
							asksHead := big.NewInt(0).Sub(price, big.NewInt(1))
							bidsHead := big.NewInt(0).Add(price, big.NewInt(1))
							t.Run("when order crosses market", func(t *testing.T) {
								t.Run("it returns error if longOrder's price >= asksHead", func(t *testing.T) {
									positionSize := big.NewInt(0).Mul(longBaseAssetQuantity, big.NewInt(-1))
									reduceOnlyAmount := big.NewInt(0)

									t.Run("it returns error if longOrder's price = asksHead", func(t *testing.T) {
										longPrice := big.NewInt(0).Set(asksHead)
										longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, longPrice, salt, true, true)

										mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
										mockBibliophile.EXPECT().GetMinSizeRequirement(longOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
										orderHash, err := GetLimitOrderHashFromContractStruct(&longOrder)
										if err != nil {
											panic("error in getting longOrder hash")
										}
										mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
										mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
										mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, longOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
										mockBibliophile.EXPECT().GetLongOpenOrdersAmount(trader, longOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
										mockBibliophile.EXPECT().GetShortOpenOrdersAmount(trader, longOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
										mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
										mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
										output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: longOrder, Sender: trader})
										assert.Equal(t, ErrCrossingMarket.Error(), output.Err)
										assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
										assert.Equal(t, ammAddress, output.Res.Amm)
										assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
									})
									t.Run("it returns error if longOrder's price > asksHead", func(t *testing.T) {
										longPrice := big.NewInt(0).Add(asksHead, big.NewInt(1))
										longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, longPrice, salt, true, true)

										mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
										mockBibliophile.EXPECT().GetMinSizeRequirement(longOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
										orderHash, err := GetLimitOrderHashFromContractStruct(&longOrder)
										if err != nil {
											panic("error in getting longOrder hash")
										}
										mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
										mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
										mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, longOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
										mockBibliophile.EXPECT().GetLongOpenOrdersAmount(trader, longOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
										mockBibliophile.EXPECT().GetShortOpenOrdersAmount(trader, longOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
										mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
										mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
										output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: longOrder, Sender: trader})
										assert.Equal(t, ErrCrossingMarket.Error(), output.Err)
										assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
										assert.Equal(t, ammAddress, output.Res.Amm)
										assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
									})
								})
								t.Run("it returns error if shortOrder's price <= bidsHead", func(t *testing.T) {
									positionSize := big.NewInt(0).Mul(shortBaseAssetQuantity, big.NewInt(-1))
									reduceOnlyAmount := big.NewInt(0)

									t.Run("it returns error if shortOrder price = asksHead", func(t *testing.T) {
										shortOrderPrice := big.NewInt(0).Set(bidsHead)
										shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, shortOrderPrice, salt, true, true)
										mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
										mockBibliophile.EXPECT().GetMinSizeRequirement(shortOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
										orderHash, err := GetLimitOrderHashFromContractStruct(&shortOrder)
										if err != nil {
											panic("error in getting shortOrder hash")
										}
										mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
										mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
										mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, shortOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
										mockBibliophile.EXPECT().GetLongOpenOrdersAmount(trader, shortOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
										mockBibliophile.EXPECT().GetShortOpenOrdersAmount(trader, shortOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
										mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
										mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)

										output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: shortOrder, Sender: trader})
										assert.Equal(t, ErrCrossingMarket.Error(), output.Err)
										assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
										assert.Equal(t, ammAddress, output.Res.Amm)
										assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
									})
									t.Run("it returns error if shortOrder price < asksHead", func(t *testing.T) {
										shortOrderPrice := big.NewInt(0).Sub(bidsHead, big.NewInt(1))
										shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, shortOrderPrice, salt, true, true)
										mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
										mockBibliophile.EXPECT().GetMinSizeRequirement(shortOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
										orderHash, err := GetLimitOrderHashFromContractStruct(&shortOrder)
										if err != nil {
											panic("error in getting shortOrder hash")
										}
										mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
										mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
										mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, shortOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
										mockBibliophile.EXPECT().GetLongOpenOrdersAmount(trader, shortOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
										mockBibliophile.EXPECT().GetShortOpenOrdersAmount(trader, shortOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
										mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
										mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)

										output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: shortOrder, Sender: trader})
										assert.Equal(t, ErrCrossingMarket.Error(), output.Err)
										assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
										assert.Equal(t, ammAddress, output.Res.Amm)
										assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
									})
								})
							})
							t.Run("when order does not cross market", func(t *testing.T) {
								t.Run("for a longOrder it returns no error and 0 as reserveAmount", func(t *testing.T) {
									positionSize := big.NewInt(0).Mul(longBaseAssetQuantity, big.NewInt(-1))
									reduceOnlyAmount := big.NewInt(0)

									longPrice := big.NewInt(0).Sub(asksHead, big.NewInt(1))
									longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, longPrice, salt, true, true)

									mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
									mockBibliophile.EXPECT().GetMinSizeRequirement(longOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
									orderHash, err := GetLimitOrderHashFromContractStruct(&longOrder)
									if err != nil {
										panic("error in getting longOrder hash")
									}
									mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
									mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
									mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, longOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
									mockBibliophile.EXPECT().GetLongOpenOrdersAmount(trader, longOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
									mockBibliophile.EXPECT().GetShortOpenOrdersAmount(trader, longOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
									mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
									mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
									mockBibliophile.EXPECT().HasReferrer(trader).Return(true).Times(1)
									mockBibliophile.EXPECT().GetPriceMultiplier(ammAddress).Return(big.NewInt(1)).Times(1)
									output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: longOrder, Sender: trader})
									assert.Equal(t, "", output.Err)
									assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
									assert.Equal(t, ammAddress, output.Res.Amm)
									assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
								})
								t.Run("for a shortOrder it returns no error and 0 as reserveAmount", func(t *testing.T) {
									positionSize := big.NewInt(0).Mul(shortBaseAssetQuantity, big.NewInt(-1))
									reduceOnlyAmount := big.NewInt(0)

									shortOrderPrice := big.NewInt(0).Add(bidsHead, big.NewInt(1))
									shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, shortOrderPrice, salt, true, true)
									mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
									mockBibliophile.EXPECT().GetMinSizeRequirement(shortOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
									orderHash, err := GetLimitOrderHashFromContractStruct(&shortOrder)
									if err != nil {
										panic("error in getting shortOrder hash")
									}
									mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
									mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
									mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, shortOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
									mockBibliophile.EXPECT().GetLongOpenOrdersAmount(trader, shortOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
									mockBibliophile.EXPECT().GetShortOpenOrdersAmount(trader, shortOrder.AmmIndex).Return(big.NewInt(0)).Times(1)
									mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
									mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
									mockBibliophile.EXPECT().HasReferrer(trader).Return(true).Times(1)
									mockBibliophile.EXPECT().GetPriceMultiplier(ammAddress).Return(big.NewInt(1)).Times(1)
									output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: shortOrder, Sender: trader})
									assert.Equal(t, "", output.Err)
									assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
									assert.Equal(t, ammAddress, output.Res.Amm)
									assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
								})
							})
						})
					})
				})
			})
		})
		t.Run("when order is not reduceOnly order", func(t *testing.T) {
			t.Run("When order is in opposite direction of position and there are reduceOnly orders in orderbook", func(t *testing.T) {
				t.Run("it returns error for a long Order", func(t *testing.T) {
					longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, false, postOnly)
					positionSize := big.NewInt(0).Mul(longBaseAssetQuantity, big.NewInt(-3)) // short position
					reduceOnlyAmount := big.NewInt(0).Div(longBaseAssetQuantity, big.NewInt(2))

					mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
					mockBibliophile.EXPECT().GetMinSizeRequirement(longOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
					orderHash, err := GetLimitOrderHashFromContractStruct(&longOrder)
					if err != nil {
						panic("error in getting shortOrder hash")
					}
					mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
					mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
					mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, longOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)

					output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: longOrder, Sender: trader})
					assert.Equal(t, ErrOpenReduceOnlyOrders.Error(), output.Err)
					assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
					assert.Equal(t, ammAddress, output.Res.Amm)
					assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
				})
				t.Run("it returns error for a short Order", func(t *testing.T) {
					shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, false, postOnly)
					positionSize := big.NewInt(0).Mul(shortBaseAssetQuantity, big.NewInt(-3)) // long position
					reduceOnlyAmount := big.NewInt(0).Div(shortBaseAssetQuantity, big.NewInt(2))

					mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
					mockBibliophile.EXPECT().GetMinSizeRequirement(shortOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
					orderHash, err := GetLimitOrderHashFromContractStruct(&shortOrder)
					if err != nil {
						panic("error in getting shortOrder hash")
					}
					mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
					mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
					mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, shortOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)

					output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: shortOrder, Sender: trader})
					assert.Equal(t, ErrOpenReduceOnlyOrders.Error(), output.Err)
					assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
					assert.Equal(t, ammAddress, output.Res.Amm)
					assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
				})
			})
			//Using a bad description here. Not sure how to write it properly. I dont want to test so many branches
			t.Run("when above is not true", func(t *testing.T) {
				t.Run("when trader does not have available margin for order", func(t *testing.T) {
					t.Run("it returns error for a long Order", func(t *testing.T) {
						longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, false, postOnly)
						positionSize := big.NewInt(0).Mul(longBaseAssetQuantity, big.NewInt(-1)) // short position
						reduceOnlyAmount := big.NewInt(0)
						minAllowableMargin := big.NewInt(100000)
						takerFee := big.NewInt(5000)
						lowerBound := hu.Div(price, big.NewInt(2))
						upperBound := hu.Add(price, lowerBound)

						t.Run("when available margin is 0", func(t *testing.T) {
							mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
							mockBibliophile.EXPECT().GetMinSizeRequirement(longOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
							orderHash, err := GetLimitOrderHashFromContractStruct(&longOrder)
							if err != nil {
								panic("error in getting longOrder hash")
							}
							mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
							mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
							mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, longOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
							mockBibliophile.EXPECT().GetUpperAndLowerBoundForMarket(longOrder.AmmIndex.Int64()).Return(upperBound, lowerBound).Times(1)
							mockBibliophile.EXPECT().GetMinAllowableMargin().Return(minAllowableMargin).Times(1)
							mockBibliophile.EXPECT().GetTakerFee().Return(takerFee).Times(1)
							availableMargin := big.NewInt(0)
							mockBibliophile.EXPECT().GetTimeStamp().Return(hu.V1ActivationTime).Times(1)
							mockBibliophile.EXPECT().GetAvailableMargin(trader, hu.V1).Return(availableMargin).Times(1)

							output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: longOrder, Sender: trader})
							assert.Equal(t, ErrInsufficientMargin.Error(), output.Err)
							assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
							assert.Equal(t, ammAddress, output.Res.Amm)
							assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
						})
						t.Run("when available margin is one less than requiredMargin", func(t *testing.T) {
							mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
							mockBibliophile.EXPECT().GetMinSizeRequirement(longOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
							orderHash, err := GetLimitOrderHashFromContractStruct(&longOrder)
							if err != nil {
								panic("error in getting longOrder hash")
							}
							mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
							mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
							mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, longOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
							mockBibliophile.EXPECT().GetUpperAndLowerBoundForMarket(longOrder.AmmIndex.Int64()).Return(upperBound, lowerBound).Times(1)
							mockBibliophile.EXPECT().GetMinAllowableMargin().Return(minAllowableMargin).Times(1)
							mockBibliophile.EXPECT().GetTakerFee().Return(takerFee).Times(1)
							quoteAsset := big.NewInt(0).Abs(hu.Div(hu.Mul(longOrder.BaseAssetQuantity, longOrder.Price), big.NewInt(1e18)))
							requiredMargin := hu.Div(hu.Mul(hu.Add(takerFee, minAllowableMargin), quoteAsset), big.NewInt(1e6))
							availableMargin := hu.Sub(requiredMargin, big.NewInt(1))
							mockBibliophile.EXPECT().GetTimeStamp().Return(hu.V1ActivationTime).Times(1)
							mockBibliophile.EXPECT().GetAvailableMargin(trader, hu.V1).Return(availableMargin).Times(1)

							output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: longOrder, Sender: trader})
							assert.Equal(t, ErrInsufficientMargin.Error(), output.Err)
							assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
							assert.Equal(t, ammAddress, output.Res.Amm)
							assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
						})
					})
					t.Run("it returns error for a short Order", func(t *testing.T) {
						shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, false, postOnly)
						positionSize := big.NewInt(0).Mul(shortBaseAssetQuantity, big.NewInt(-1)) // short position
						reduceOnlyAmount := big.NewInt(0)
						minAllowableMargin := big.NewInt(100000)
						takerFee := big.NewInt(5000)
						lowerBound := hu.Div(price, big.NewInt(2))
						upperBound := hu.Add(price, lowerBound)

						t.Run("when available margin is 0", func(t *testing.T) {
							mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
							mockBibliophile.EXPECT().GetMinSizeRequirement(shortOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
							orderHash, err := GetLimitOrderHashFromContractStruct(&shortOrder)
							if err != nil {
								panic("error in getting shortOrder hash")
							}
							mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
							mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
							mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, shortOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
							mockBibliophile.EXPECT().GetUpperAndLowerBoundForMarket(shortOrder.AmmIndex.Int64()).Return(upperBound, lowerBound).Times(1)
							mockBibliophile.EXPECT().GetMinAllowableMargin().Return(minAllowableMargin).Times(1)
							mockBibliophile.EXPECT().GetTakerFee().Return(takerFee).Times(1)
							availableMargin := big.NewInt(0)
							mockBibliophile.EXPECT().GetTimeStamp().Return(hu.V1ActivationTime).Times(1)
							mockBibliophile.EXPECT().GetAvailableMargin(trader, hu.V1).Return(availableMargin).Times(1)

							output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: shortOrder, Sender: trader})
							assert.Equal(t, ErrInsufficientMargin.Error(), output.Err)
							assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
							assert.Equal(t, ammAddress, output.Res.Amm)
							assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
						})
						t.Run("when available margin is one less than requiredMargin", func(t *testing.T) {
							mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
							mockBibliophile.EXPECT().GetMinSizeRequirement(shortOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
							orderHash, err := GetLimitOrderHashFromContractStruct(&shortOrder)
							if err != nil {
								panic("error in getting shortOrder hash")
							}
							mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
							mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
							mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, shortOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
							mockBibliophile.EXPECT().GetUpperAndLowerBoundForMarket(shortOrder.AmmIndex.Int64()).Return(upperBound, lowerBound).Times(1)
							mockBibliophile.EXPECT().GetMinAllowableMargin().Return(minAllowableMargin).Times(1)
							mockBibliophile.EXPECT().GetTakerFee().Return(takerFee).Times(1)
							// use upperBound as price to calculate quoteAsset for short
							quoteAsset := big.NewInt(0).Abs(hu.Div(hu.Mul(shortOrder.BaseAssetQuantity, upperBound), big.NewInt(1e18)))
							requiredMargin := hu.Div(hu.Mul(hu.Add(takerFee, minAllowableMargin), quoteAsset), big.NewInt(1e6))
							availableMargin := hu.Sub(requiredMargin, big.NewInt(1))
							mockBibliophile.EXPECT().GetTimeStamp().Return(hu.V1ActivationTime).Times(1)
							mockBibliophile.EXPECT().GetAvailableMargin(trader, hu.V1).Return(availableMargin).Times(1)

							output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: shortOrder, Sender: trader})
							assert.Equal(t, ErrInsufficientMargin.Error(), output.Err)
							assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
							assert.Equal(t, ammAddress, output.Res.Amm)
							assert.Equal(t, big.NewInt(0), output.Res.ReserveAmount)
						})
					})
				})
				t.Run("when trader has available margin for order", func(t *testing.T) {
					t.Run("when order is not a postOnly order", func(t *testing.T) {
						minAllowableMargin := big.NewInt(100000)
						takerFee := big.NewInt(5000)
						reduceOnlyAmount := big.NewInt(0)
						t.Run("it returns nil error and reserverAmount when order is a long order", func(t *testing.T) {
							longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, false, false)
							positionSize := big.NewInt(0).Mul(longBaseAssetQuantity, big.NewInt(-1)) // short position
							quoteAsset := big.NewInt(0).Abs(hu.Div(hu.Mul(longOrder.BaseAssetQuantity, longOrder.Price), big.NewInt(1e18)))
							requiredMargin := hu.Div(hu.Mul(hu.Add(takerFee, minAllowableMargin), quoteAsset), big.NewInt(1e6))
							availableMargin := hu.Add(requiredMargin, big.NewInt(1))
							lowerBound := hu.Div(price, big.NewInt(2))
							upperBound := hu.Add(price, lowerBound)

							mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
							mockBibliophile.EXPECT().GetMinSizeRequirement(longOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
							orderHash, err := GetLimitOrderHashFromContractStruct(&longOrder)
							if err != nil {
								panic("error in getting longOrder hash")
							}
							mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
							mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
							mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, longOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
							mockBibliophile.EXPECT().GetUpperAndLowerBoundForMarket(longOrder.AmmIndex.Int64()).Return(upperBound, lowerBound).Times(1)
							mockBibliophile.EXPECT().GetMinAllowableMargin().Return(minAllowableMargin).Times(1)
							mockBibliophile.EXPECT().GetTakerFee().Return(takerFee).Times(1)
							mockBibliophile.EXPECT().GetTimeStamp().Return(hu.V1ActivationTime).Times(1)
							mockBibliophile.EXPECT().GetAvailableMargin(trader, hu.V1).Return(availableMargin).Times(1)
							mockBibliophile.EXPECT().HasReferrer(trader).Return(true).Times(1)
							mockBibliophile.EXPECT().GetPriceMultiplier(ammAddress).Return(big.NewInt(1)).Times(1)
							output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: longOrder, Sender: trader})
							assert.Equal(t, "", output.Err)
							assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
							assert.Equal(t, ammAddress, output.Res.Amm)
							assert.Equal(t, requiredMargin, output.Res.ReserveAmount)
						})
						t.Run("it returns nil error and reserverAmount when order is a short order", func(t *testing.T) {
							lowerBound := hu.Div(price, big.NewInt(2))
							upperBound := hu.Add(price, lowerBound)
							shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, false, false)
							positionSize := big.NewInt(0).Mul(shortBaseAssetQuantity, big.NewInt(-1)) // long position
							quoteAsset := big.NewInt(0).Abs(hu.Div(hu.Mul(shortOrder.BaseAssetQuantity, upperBound), big.NewInt(1e18)))
							requiredMargin := hu.Div(hu.Mul(hu.Add(takerFee, minAllowableMargin), quoteAsset), big.NewInt(1e6))
							availableMargin := hu.Add(requiredMargin, big.NewInt(1))

							mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
							mockBibliophile.EXPECT().GetMinSizeRequirement(shortOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
							orderHash, err := GetLimitOrderHashFromContractStruct(&shortOrder)
							if err != nil {
								panic("error in getting shortOrder hash")
							}
							mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
							mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
							mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, shortOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
							mockBibliophile.EXPECT().GetUpperAndLowerBoundForMarket(shortOrder.AmmIndex.Int64()).Return(upperBound, lowerBound).Times(1)
							mockBibliophile.EXPECT().GetMinAllowableMargin().Return(minAllowableMargin).Times(1)
							mockBibliophile.EXPECT().GetTakerFee().Return(takerFee).Times(1)
							mockBibliophile.EXPECT().GetTimeStamp().Return(hu.V1ActivationTime).Times(1)
							mockBibliophile.EXPECT().GetAvailableMargin(trader, hu.V1).Return(availableMargin).Times(1)
							mockBibliophile.EXPECT().HasReferrer(trader).Return(true).Times(1)
							mockBibliophile.EXPECT().GetPriceMultiplier(ammAddress).Return(big.NewInt(1)).Times(1)
							output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: shortOrder, Sender: trader})
							assert.Equal(t, "", output.Err)
							assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
							assert.Equal(t, ammAddress, output.Res.Amm)
							assert.Equal(t, requiredMargin, output.Res.ReserveAmount)
						})
					})
					t.Run("when order is a postOnly order", func(t *testing.T) {
						asksHead := big.NewInt(0).Add(price, big.NewInt(1))
						bidsHead := big.NewInt(0).Sub(price, big.NewInt(1))
						minAllowableMargin := big.NewInt(100000)
						takerFee := big.NewInt(5000)
						reduceOnlyAmount := big.NewInt(0)

						t.Run("when order crosses market", func(t *testing.T) {
							t.Run("it returns error if longOrder's price >= asksHead", func(t *testing.T) {
								t.Run("it returns error if longOrder's price = asksHead", func(t *testing.T) {
									longPrice := big.NewInt(0).Set(asksHead)
									longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, longPrice, salt, false, true)
									positionSize := big.NewInt(0).Mul(longBaseAssetQuantity, big.NewInt(-1)) // short position
									quoteAsset := big.NewInt(0).Abs(hu.Div(hu.Mul(longOrder.BaseAssetQuantity, longOrder.Price), big.NewInt(1e18)))
									requiredMargin := hu.Add(hu.Div(hu.Mul(minAllowableMargin, quoteAsset), big.NewInt(1e6)), hu.Div(hu.Mul(takerFee, quoteAsset), big.NewInt(1e6)))
									availableMargin := hu.Add(requiredMargin, big.NewInt(1))
									lowerBound := hu.Div(price, big.NewInt(2))
									upperBound := hu.Add(price, lowerBound)

									mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
									mockBibliophile.EXPECT().GetMinSizeRequirement(longOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
									orderHash, err := GetLimitOrderHashFromContractStruct(&longOrder)
									if err != nil {
										panic("error in getting longOrder hash")
									}
									mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
									mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
									mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, longOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
									mockBibliophile.EXPECT().GetUpperAndLowerBoundForMarket(longOrder.AmmIndex.Int64()).Return(upperBound, lowerBound).Times(1)
									mockBibliophile.EXPECT().GetMinAllowableMargin().Return(minAllowableMargin).Times(1)
									mockBibliophile.EXPECT().GetTakerFee().Return(takerFee).Times(1)
									mockBibliophile.EXPECT().GetTimeStamp().Return(hu.V1ActivationTime).Times(1)
									mockBibliophile.EXPECT().GetAvailableMargin(trader, hu.V1).Return(availableMargin).Times(1)
									mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
									mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
									output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: longOrder, Sender: trader})
									assert.Equal(t, ErrCrossingMarket.Error(), output.Err)
									assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
									assert.Equal(t, ammAddress, output.Res.Amm)
									assert.Equal(t, requiredMargin, output.Res.ReserveAmount)
								})
								t.Run("it returns error if longOrder's price > asksHead", func(t *testing.T) {
									longPrice := big.NewInt(0).Add(asksHead, big.NewInt(1))
									longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, longPrice, salt, false, true)
									positionSize := big.NewInt(0).Mul(longBaseAssetQuantity, big.NewInt(-1)) // short position
									quoteAsset := big.NewInt(0).Abs(hu.Div(hu.Mul(longOrder.BaseAssetQuantity, longOrder.Price), big.NewInt(1e18)))
									requiredMargin := hu.Div(hu.Mul(hu.Add(takerFee, minAllowableMargin), quoteAsset), big.NewInt(1e6))
									availableMargin := hu.Add(requiredMargin, big.NewInt(1))
									lowerBound := hu.Div(price, big.NewInt(2))
									upperBound := hu.Add(price, lowerBound)

									mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
									mockBibliophile.EXPECT().GetMinSizeRequirement(longOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
									orderHash, err := GetLimitOrderHashFromContractStruct(&longOrder)
									if err != nil {
										panic("error in getting longOrder hash")
									}
									mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
									mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
									mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, longOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
									mockBibliophile.EXPECT().GetUpperAndLowerBoundForMarket(longOrder.AmmIndex.Int64()).Return(upperBound, lowerBound).Times(1)
									mockBibliophile.EXPECT().GetMinAllowableMargin().Return(minAllowableMargin).Times(1)
									mockBibliophile.EXPECT().GetTakerFee().Return(takerFee).Times(1)
									mockBibliophile.EXPECT().GetTimeStamp().Return(hu.V1ActivationTime).Times(1)
									mockBibliophile.EXPECT().GetAvailableMargin(trader, hu.V1).Return(availableMargin).Times(1)
									mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
									mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
									output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: longOrder, Sender: trader})
									assert.Equal(t, ErrCrossingMarket.Error(), output.Err)
									assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
									assert.Equal(t, ammAddress, output.Res.Amm)
									assert.Equal(t, requiredMargin, output.Res.ReserveAmount)
								})
							})
							t.Run("it returns error if shortOrder's price <= bidsHead", func(t *testing.T) {
								positionSize := big.NewInt(0).Mul(shortBaseAssetQuantity, big.NewInt(-1))

								t.Run("it returns error if shortOrder price = asksHead", func(t *testing.T) {
									shortOrderPrice := big.NewInt(0).Set(bidsHead)
									lowerBound := hu.Div(price, big.NewInt(2))
									upperBound := hu.Add(price, lowerBound)
									shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, shortOrderPrice, salt, false, true)
									quoteAsset := big.NewInt(0).Abs(hu.Div(hu.Mul(shortOrder.BaseAssetQuantity, upperBound), big.NewInt(1e18)))
									requiredMargin := hu.Add(hu.Div(hu.Mul(minAllowableMargin, quoteAsset), big.NewInt(1e6)), hu.Div(hu.Mul(takerFee, quoteAsset), big.NewInt(1e6)))
									availableMargin := hu.Add(requiredMargin, big.NewInt(1))

									mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
									mockBibliophile.EXPECT().GetMinSizeRequirement(shortOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
									orderHash, err := GetLimitOrderHashFromContractStruct(&shortOrder)
									if err != nil {
										panic("error in getting shortOrder hash")
									}
									mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
									mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
									mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, shortOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
									mockBibliophile.EXPECT().GetUpperAndLowerBoundForMarket(shortOrder.AmmIndex.Int64()).Return(upperBound, lowerBound).Times(1)
									mockBibliophile.EXPECT().GetMinAllowableMargin().Return(minAllowableMargin).Times(1)
									mockBibliophile.EXPECT().GetTakerFee().Return(takerFee).Times(1)
									mockBibliophile.EXPECT().GetTimeStamp().Return(hu.V1ActivationTime).Times(1)
									mockBibliophile.EXPECT().GetAvailableMargin(trader, hu.V1).Return(availableMargin).Times(1)
									mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
									mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)

									output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: shortOrder, Sender: trader})
									assert.Equal(t, ErrCrossingMarket.Error(), output.Err)
									assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
									assert.Equal(t, ammAddress, output.Res.Amm)
									assert.Equal(t, requiredMargin, output.Res.ReserveAmount)
								})
								t.Run("it returns error if shortOrder price < asksHead", func(t *testing.T) {
									shortOrderPrice := big.NewInt(0).Sub(bidsHead, big.NewInt(1))
									lowerBound := hu.Div(price, big.NewInt(2))
									upperBound := hu.Add(price, lowerBound)
									shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, shortOrderPrice, salt, false, true)
									quoteAsset := big.NewInt(0).Abs(hu.Div(hu.Mul(shortOrder.BaseAssetQuantity, upperBound), big.NewInt(1e18)))
									requiredMargin := hu.Add(hu.Div(hu.Mul(minAllowableMargin, quoteAsset), big.NewInt(1e6)), hu.Div(hu.Mul(takerFee, quoteAsset), big.NewInt(1e6)))
									availableMargin := hu.Add(requiredMargin, big.NewInt(1))

									mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
									mockBibliophile.EXPECT().GetMinSizeRequirement(shortOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
									orderHash, err := GetLimitOrderHashFromContractStruct(&shortOrder)
									if err != nil {
										panic("error in getting shortOrder hash")
									}
									mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
									mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
									mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, shortOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
									mockBibliophile.EXPECT().GetUpperAndLowerBoundForMarket(shortOrder.AmmIndex.Int64()).Return(upperBound, lowerBound).Times(1)
									mockBibliophile.EXPECT().GetMinAllowableMargin().Return(minAllowableMargin).Times(1)
									mockBibliophile.EXPECT().GetTakerFee().Return(takerFee).Times(1)
									mockBibliophile.EXPECT().GetTimeStamp().Return(hu.V1ActivationTime).Times(1)
									mockBibliophile.EXPECT().GetAvailableMargin(trader, hu.V1).Return(availableMargin).Times(1)
									mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
									mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)

									output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: shortOrder, Sender: trader})
									assert.Equal(t, ErrCrossingMarket.Error(), output.Err)
									assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
									assert.Equal(t, ammAddress, output.Res.Amm)
									assert.Equal(t, requiredMargin, output.Res.ReserveAmount)
								})
							})
						})
						t.Run("when order does not cross market", func(t *testing.T) {
							t.Run("for a longOrder it returns no error and 0 as reserveAmount", func(t *testing.T) {
								longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, false, true)
								positionSize := big.NewInt(0)
								quoteAsset := big.NewInt(0).Abs(hu.Div(hu.Mul(longOrder.BaseAssetQuantity, longOrder.Price), big.NewInt(1e18)))
								requiredMargin := hu.Add(hu.Div(hu.Mul(minAllowableMargin, quoteAsset), big.NewInt(1e6)), hu.Div(hu.Mul(takerFee, quoteAsset), big.NewInt(1e6)))
								availableMargin := hu.Add(requiredMargin, big.NewInt(1))
								lowerBound := hu.Div(price, big.NewInt(2))
								upperBound := hu.Add(price, lowerBound)

								mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
								mockBibliophile.EXPECT().GetMinSizeRequirement(longOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
								orderHash, err := GetLimitOrderHashFromContractStruct(&longOrder)
								if err != nil {
									panic("error in getting longOrder hash")
								}
								mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
								mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
								mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, longOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
								mockBibliophile.EXPECT().GetUpperAndLowerBoundForMarket(longOrder.AmmIndex.Int64()).Return(upperBound, lowerBound).Times(1)
								mockBibliophile.EXPECT().GetMinAllowableMargin().Return(minAllowableMargin).Times(1)
								mockBibliophile.EXPECT().GetTakerFee().Return(takerFee).Times(1)
								mockBibliophile.EXPECT().GetTimeStamp().Return(hu.V1ActivationTime).Times(1)
								mockBibliophile.EXPECT().GetAvailableMargin(trader, hu.V1).Return(availableMargin).Times(1)
								mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
								mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
								mockBibliophile.EXPECT().HasReferrer(trader).Return(true).Times(1)
								mockBibliophile.EXPECT().GetPriceMultiplier(ammAddress).Return(big.NewInt(1)).Times(1)
								output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: longOrder, Sender: trader})
								assert.Equal(t, "", output.Err)
								assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
								assert.Equal(t, ammAddress, output.Res.Amm)
								assert.Equal(t, requiredMargin, output.Res.ReserveAmount)
							})
							t.Run("for a shortOrder it returns no error and 0 as reserveAmount", func(t *testing.T) {
								shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, false, true)
								positionSize := big.NewInt(0)
								lowerBound := hu.Div(price, big.NewInt(2))
								upperBound := hu.Add(price, lowerBound)
								quoteAsset := big.NewInt(0).Abs(hu.Div(hu.Mul(shortOrder.BaseAssetQuantity, upperBound), big.NewInt(1e18)))
								requiredMargin := hu.Add(hu.Div(hu.Mul(minAllowableMargin, quoteAsset), big.NewInt(1e6)), hu.Div(hu.Mul(takerFee, quoteAsset), big.NewInt(1e6)))
								availableMargin := hu.Add(requiredMargin, big.NewInt(1))

								mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
								mockBibliophile.EXPECT().GetMinSizeRequirement(shortOrder.AmmIndex.Int64()).Return(minSizeRequirement).Times(1)
								orderHash, err := GetLimitOrderHashFromContractStruct(&shortOrder)
								if err != nil {
									panic("error in getting shortOrder hash")
								}
								mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
								mockBibliophile.EXPECT().GetSize(ammAddress, &trader).Return(positionSize).Times(1)
								mockBibliophile.EXPECT().GetReduceOnlyAmount(trader, shortOrder.AmmIndex).Return(reduceOnlyAmount).Times(1)
								mockBibliophile.EXPECT().GetUpperAndLowerBoundForMarket(shortOrder.AmmIndex.Int64()).Return(upperBound, lowerBound).Times(1)
								mockBibliophile.EXPECT().GetMinAllowableMargin().Return(minAllowableMargin).Times(1)
								mockBibliophile.EXPECT().GetTakerFee().Return(takerFee).Times(1)
								mockBibliophile.EXPECT().GetTimeStamp().Return(hu.V1ActivationTime).Times(1)
								mockBibliophile.EXPECT().GetAvailableMargin(trader, hu.V1).Return(availableMargin).Times(1)
								mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
								mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
								mockBibliophile.EXPECT().HasReferrer(trader).Return(true).Times(1)
								mockBibliophile.EXPECT().GetPriceMultiplier(ammAddress).Return(big.NewInt(1)).Times(1)
								output := ValidatePlaceLimitOrder(mockBibliophile, &ValidatePlaceLimitOrderInput{Order: shortOrder, Sender: trader})
								assert.Equal(t, "", output.Err)
								assert.Equal(t, orderHash, common.BytesToHash(output.Orderhash[:]))
								assert.Equal(t, ammAddress, output.Res.Amm)
								assert.Equal(t, requiredMargin, output.Res.ReserveAmount)
							})
						})
					})
				})
			})
		})
	})
}

func TestValidateCancelLimitOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBibliophile := b.NewMockBibliophileClient(ctrl)
	ammIndex := big.NewInt(0)
	longBaseAssetQuantity := big.NewInt(5000000000000000000)
	shortBaseAssetQuantity := big.NewInt(-5000000000000000000)
	price := big.NewInt(100000000)
	salt := big.NewInt(121)
	reduceOnly := false
	postOnly := false
	trader := common.HexToAddress("0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC")
	ammAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	assertLowMargin := false

	t.Run("when sender is not the trader and is not trading authority, it returns error", func(t *testing.T) {
		sender := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C9")
		t.Run("it returns error for a long order", func(t *testing.T) {
			order := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, reduceOnly, postOnly)
			input := getValidateCancelLimitOrderInput(order, sender, assertLowMargin)
			mockBibliophile.EXPECT().IsTradingAuthority(order.Trader, sender).Return(false).Times(1)
			output := ValidateCancelLimitOrder(mockBibliophile, &input)
			assert.Equal(t, ErrNoTradingAuthority.Error(), output.Err)
		})
		t.Run("it returns error for a short order", func(t *testing.T) {
			order := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, reduceOnly, postOnly)
			input := getValidateCancelLimitOrderInput(order, sender, assertLowMargin)
			mockBibliophile.EXPECT().IsTradingAuthority(order.Trader, sender).Return(false).Times(1)
			output := ValidateCancelLimitOrder(mockBibliophile, &input)
			assert.Equal(t, ErrNoTradingAuthority.Error(), output.Err)
		})
	})
	t.Run("when either sender is trader or a trading authority", func(t *testing.T) {
		t.Run("When order status is not placed", func(t *testing.T) {
			t.Run("when order status was never placed", func(t *testing.T) {
				t.Run("it returns error for a longOrder", func(t *testing.T) {
					longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, reduceOnly, postOnly)
					orderHash := getOrderHash(longOrder)
					mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
					input := getValidateCancelLimitOrderInput(longOrder, trader, assertLowMargin)
					output := ValidateCancelLimitOrder(mockBibliophile, &input)
					assert.Equal(t, "Invalid", output.Err)
					assert.Equal(t, orderHash, common.BytesToHash(output.OrderHash[:]))
					assert.Equal(t, common.Address{}, output.Res.Amm)
					assert.Equal(t, big.NewInt(0), output.Res.UnfilledAmount)
				})
				t.Run("it returns error for a shortOrder", func(t *testing.T) {
					shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, reduceOnly, postOnly)
					orderHash := getOrderHash(shortOrder)
					mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Invalid)).Times(1)
					input := getValidateCancelLimitOrderInput(shortOrder, trader, assertLowMargin)
					output := ValidateCancelLimitOrder(mockBibliophile, &input)
					assert.Equal(t, "Invalid", output.Err)
					assert.Equal(t, orderHash, common.BytesToHash(output.OrderHash[:]))
					assert.Equal(t, common.Address{}, output.Res.Amm)
					assert.Equal(t, big.NewInt(0), output.Res.UnfilledAmount)
				})
			})
			t.Run("when order status is cancelled", func(t *testing.T) {
				t.Run("it returns error for a longOrder", func(t *testing.T) {
					longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, reduceOnly, postOnly)
					orderHash := getOrderHash(longOrder)
					mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Cancelled)).Times(1)
					input := getValidateCancelLimitOrderInput(longOrder, trader, assertLowMargin)
					output := ValidateCancelLimitOrder(mockBibliophile, &input)
					assert.Equal(t, "Cancelled", output.Err)
					assert.Equal(t, orderHash, common.BytesToHash(output.OrderHash[:]))
					assert.Equal(t, common.Address{}, output.Res.Amm)
					assert.Equal(t, big.NewInt(0), output.Res.UnfilledAmount)
				})
				t.Run("it returns error for a shortOrder", func(t *testing.T) {
					shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, reduceOnly, postOnly)
					orderHash := getOrderHash(shortOrder)
					mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Cancelled)).Times(1)
					input := getValidateCancelLimitOrderInput(shortOrder, trader, assertLowMargin)
					output := ValidateCancelLimitOrder(mockBibliophile, &input)
					assert.Equal(t, "Cancelled", output.Err)
					assert.Equal(t, orderHash, common.BytesToHash(output.OrderHash[:]))
					assert.Equal(t, common.Address{}, output.Res.Amm)
					assert.Equal(t, big.NewInt(0), output.Res.UnfilledAmount)
				})
			})
			t.Run("when order status is filled", func(t *testing.T) {
				t.Run("it returns error for a longOrder", func(t *testing.T) {
					longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, reduceOnly, postOnly)
					orderHash := getOrderHash(longOrder)
					mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Filled)).Times(1)
					input := getValidateCancelLimitOrderInput(longOrder, trader, assertLowMargin)
					output := ValidateCancelLimitOrder(mockBibliophile, &input)
					assert.Equal(t, "Filled", output.Err)
					assert.Equal(t, orderHash, common.BytesToHash(output.OrderHash[:]))
					assert.Equal(t, common.Address{}, output.Res.Amm)
					assert.Equal(t, big.NewInt(0), output.Res.UnfilledAmount)
				})
				t.Run("it returns error for a shortOrder", func(t *testing.T) {
					shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, reduceOnly, postOnly)
					orderHash := getOrderHash(shortOrder)
					mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Filled)).Times(1)
					input := getValidateCancelLimitOrderInput(shortOrder, trader, assertLowMargin)
					output := ValidateCancelLimitOrder(mockBibliophile, &input)
					assert.Equal(t, "Filled", output.Err)
					assert.Equal(t, orderHash, common.BytesToHash(output.OrderHash[:]))
					assert.Equal(t, common.Address{}, output.Res.Amm)
					assert.Equal(t, big.NewInt(0), output.Res.UnfilledAmount)
				})
			})
		})
		t.Run("When order status is placed", func(t *testing.T) {
			t.Run("when assertLowMargin is true", func(t *testing.T) {
				assertLowMargin := true
				t.Run("when availableMargin >= zero", func(t *testing.T) {
					t.Run("when availableMargin == 0 ", func(t *testing.T) {
						t.Run("it returns error for a longOrder", func(t *testing.T) {
							longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, reduceOnly, postOnly)
							orderHash := getOrderHash(longOrder)

							mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Placed)).Times(1)
							mockBibliophile.EXPECT().GetTimeStamp().Return(hu.V1ActivationTime).Times(1)
							mockBibliophile.EXPECT().GetAvailableMargin(longOrder.Trader, hu.V1).Return(big.NewInt(0)).Times(1)
							mockBibliophile.EXPECT().IsValidator(longOrder.Trader).Return(true).Times(1)
							input := getValidateCancelLimitOrderInput(longOrder, trader, assertLowMargin)
							output := ValidateCancelLimitOrder(mockBibliophile, &input)
							assert.Equal(t, "Not Low Margin", output.Err)
							assert.Equal(t, orderHash, common.BytesToHash(output.OrderHash[:]))
							assert.Equal(t, common.Address{}, output.Res.Amm)
							assert.Equal(t, big.NewInt(0), output.Res.UnfilledAmount)
						})
						t.Run("it returns error for a shortOrder", func(t *testing.T) {
							shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, reduceOnly, postOnly)
							orderHash := getOrderHash(shortOrder)

							mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Placed)).Times(1)
							mockBibliophile.EXPECT().GetTimeStamp().Return(hu.V1ActivationTime).Times(1)
							mockBibliophile.EXPECT().GetAvailableMargin(shortOrder.Trader, hu.V1).Return(big.NewInt(0)).Times(1)
							mockBibliophile.EXPECT().IsValidator(shortOrder.Trader).Return(true).Times(1)
							input := getValidateCancelLimitOrderInput(shortOrder, trader, assertLowMargin)
							output := ValidateCancelLimitOrder(mockBibliophile, &input)

							assert.Equal(t, "Not Low Margin", output.Err)
							assert.Equal(t, orderHash, common.BytesToHash(output.OrderHash[:]))
							assert.Equal(t, common.Address{}, output.Res.Amm)
							assert.Equal(t, big.NewInt(0), output.Res.UnfilledAmount)
						})
					})
					t.Run("when availableMargin > 0 ", func(t *testing.T) {
						newMargin := hu.Mul(price, longBaseAssetQuantity)
						t.Run("it returns error for a longOrder", func(t *testing.T) {
							longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, reduceOnly, postOnly)
							orderHash := getOrderHash(longOrder)

							mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Placed)).Times(1)
							mockBibliophile.EXPECT().GetTimeStamp().Return(hu.V1ActivationTime).Times(1)
							mockBibliophile.EXPECT().GetAvailableMargin(longOrder.Trader, hu.V1).Return(newMargin).Times(1)
							mockBibliophile.EXPECT().IsValidator(longOrder.Trader).Return(true).Times(1)
							input := getValidateCancelLimitOrderInput(longOrder, trader, assertLowMargin)
							output := ValidateCancelLimitOrder(mockBibliophile, &input)
							assert.Equal(t, "Not Low Margin", output.Err)
							assert.Equal(t, orderHash, common.BytesToHash(output.OrderHash[:]))
							assert.Equal(t, common.Address{}, output.Res.Amm)
							assert.Equal(t, big.NewInt(0), output.Res.UnfilledAmount)
						})
						t.Run("it returns error for a shortOrder", func(t *testing.T) {
							shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, reduceOnly, postOnly)
							orderHash := getOrderHash(shortOrder)

							mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Placed)).Times(1)
							mockBibliophile.EXPECT().GetTimeStamp().Return(hu.V1ActivationTime).Times(1)
							mockBibliophile.EXPECT().GetAvailableMargin(shortOrder.Trader, hu.V1).Return(newMargin).Times(1)
							mockBibliophile.EXPECT().IsValidator(shortOrder.Trader).Return(true).Times(1)
							input := getValidateCancelLimitOrderInput(shortOrder, trader, assertLowMargin)
							output := ValidateCancelLimitOrder(mockBibliophile, &input)
							assert.Equal(t, "Not Low Margin", output.Err)
							assert.Equal(t, orderHash, common.BytesToHash(output.OrderHash[:]))
							assert.Equal(t, common.Address{}, output.Res.Amm)
							assert.Equal(t, big.NewInt(0), output.Res.UnfilledAmount)
						})
					})
				})
				t.Run("when availableMargin < zero", func(t *testing.T) {
					t.Run("for an unfilled Order", func(t *testing.T) {
						t.Run("for a longOrder it returns err = nil, with ammAddress and unfilled amount of cancelled Order", func(t *testing.T) {
							longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, reduceOnly, postOnly)
							orderHash := getOrderHash(longOrder)

							mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Placed)).Times(1)
							mockBibliophile.EXPECT().GetTimeStamp().Return(hu.V1ActivationTime).Times(1)
							mockBibliophile.EXPECT().GetAvailableMargin(longOrder.Trader, hu.V1).Return(big.NewInt(-1)).Times(1)
							mockBibliophile.EXPECT().GetOrderFilledAmount(orderHash).Return(big.NewInt(0)).Times(1)
							mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
							mockBibliophile.EXPECT().IsValidator(longOrder.Trader).Return(true).Times(1)

							input := getValidateCancelLimitOrderInput(longOrder, trader, assertLowMargin)
							output := ValidateCancelLimitOrder(mockBibliophile, &input)
							assert.Equal(t, "", output.Err)
							assert.Equal(t, orderHash, common.BytesToHash(output.OrderHash[:]))
							assert.Equal(t, ammAddress, output.Res.Amm)
							assert.Equal(t, longOrder.BaseAssetQuantity, output.Res.UnfilledAmount)
						})
						t.Run("for a shortOrder it returns err = nil, with ammAddress and unfilled amount of cancelled Order", func(t *testing.T) {
							shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, reduceOnly, postOnly)
							orderHash := getOrderHash(shortOrder)

							mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Placed)).Times(1)
							mockBibliophile.EXPECT().GetTimeStamp().Return(hu.V1ActivationTime).Times(1)
							mockBibliophile.EXPECT().GetAvailableMargin(shortOrder.Trader, hu.V1).Return(big.NewInt(-1)).Times(1)
							mockBibliophile.EXPECT().GetOrderFilledAmount(orderHash).Return(big.NewInt(0)).Times(1)
							mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
							mockBibliophile.EXPECT().IsValidator(shortOrder.Trader).Return(true).Times(1)

							input := getValidateCancelLimitOrderInput(shortOrder, trader, assertLowMargin)
							output := ValidateCancelLimitOrder(mockBibliophile, &input)
							assert.Equal(t, "", output.Err)
							assert.Equal(t, orderHash, common.BytesToHash(output.OrderHash[:]))
							assert.Equal(t, ammAddress, output.Res.Amm)
							assert.Equal(t, shortOrder.BaseAssetQuantity, output.Res.UnfilledAmount)
						})
					})
					t.Run("for a partially filled Order", func(t *testing.T) {
						t.Run("for a longOrder it returns err = nil, with ammAddress and unfilled amount of cancelled Order", func(t *testing.T) {
							longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, reduceOnly, postOnly)
							orderHash := getOrderHash(longOrder)
							filledAmount := hu.Div(longOrder.BaseAssetQuantity, big.NewInt(2))

							mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Placed)).Times(1)
							mockBibliophile.EXPECT().GetTimeStamp().Return(hu.V1ActivationTime).Times(1)
							mockBibliophile.EXPECT().GetAvailableMargin(longOrder.Trader, hu.V1).Return(big.NewInt(-1)).Times(1)
							mockBibliophile.EXPECT().GetOrderFilledAmount(orderHash).Return(filledAmount).Times(1)
							mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
							mockBibliophile.EXPECT().IsValidator(longOrder.Trader).Return(true).Times(1)

							input := getValidateCancelLimitOrderInput(longOrder, trader, assertLowMargin)
							output := ValidateCancelLimitOrder(mockBibliophile, &input)
							assert.Equal(t, "", output.Err)
							assert.Equal(t, orderHash, common.BytesToHash(output.OrderHash[:]))
							assert.Equal(t, ammAddress, output.Res.Amm)
							expectedUnfilleAmount := hu.Sub(longOrder.BaseAssetQuantity, filledAmount)
							assert.Equal(t, expectedUnfilleAmount, output.Res.UnfilledAmount)
						})
						t.Run("for a shortOrder it returns err = nil, with ammAddress and unfilled amount of cancelled Order", func(t *testing.T) {
							shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, reduceOnly, postOnly)
							orderHash := getOrderHash(shortOrder)
							filledAmount := hu.Div(shortOrder.BaseAssetQuantity, big.NewInt(2))

							mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Placed)).Times(1)
							mockBibliophile.EXPECT().GetTimeStamp().Return(hu.V1ActivationTime).Times(1)
							mockBibliophile.EXPECT().GetAvailableMargin(shortOrder.Trader, hu.V1).Return(big.NewInt(-1)).Times(1)
							mockBibliophile.EXPECT().GetOrderFilledAmount(orderHash).Return(filledAmount).Times(1)
							mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)
							mockBibliophile.EXPECT().IsValidator(shortOrder.Trader).Return(true).Times(1)

							input := getValidateCancelLimitOrderInput(shortOrder, trader, assertLowMargin)
							output := ValidateCancelLimitOrder(mockBibliophile, &input)
							assert.Equal(t, "", output.Err)
							assert.Equal(t, orderHash, common.BytesToHash(output.OrderHash[:]))
							assert.Equal(t, ammAddress, output.Res.Amm)
							expectedUnfilleAmount := hu.Sub(shortOrder.BaseAssetQuantity, filledAmount)
							assert.Equal(t, expectedUnfilleAmount, output.Res.UnfilledAmount)
						})
					})
				})
			})
			t.Run("when assertLowMargin is false", func(t *testing.T) {
				assertLowMargin := false
				t.Run("for an unfilled Order", func(t *testing.T) {
					t.Run("for a longOrder it returns err = nil, with ammAddress and unfilled amount of cancelled Order", func(t *testing.T) {
						longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, reduceOnly, postOnly)
						orderHash := getOrderHash(longOrder)

						mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Placed)).Times(1)
						mockBibliophile.EXPECT().GetOrderFilledAmount(orderHash).Return(big.NewInt(0)).Times(1)
						mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)

						input := getValidateCancelLimitOrderInput(longOrder, trader, assertLowMargin)
						output := ValidateCancelLimitOrder(mockBibliophile, &input)
						assert.Equal(t, "", output.Err)
						assert.Equal(t, orderHash, common.BytesToHash(output.OrderHash[:]))
						assert.Equal(t, ammAddress, output.Res.Amm)
						assert.Equal(t, longOrder.BaseAssetQuantity, output.Res.UnfilledAmount)
					})
					t.Run("for a shortOrder it returns err = nil, with ammAddress and unfilled amount of cancelled Order", func(t *testing.T) {
						shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, reduceOnly, postOnly)
						orderHash := getOrderHash(shortOrder)

						mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Placed)).Times(1)
						mockBibliophile.EXPECT().GetOrderFilledAmount(orderHash).Return(big.NewInt(0)).Times(1)
						mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)

						input := getValidateCancelLimitOrderInput(shortOrder, trader, assertLowMargin)
						output := ValidateCancelLimitOrder(mockBibliophile, &input)
						assert.Equal(t, "", output.Err)
						assert.Equal(t, orderHash, common.BytesToHash(output.OrderHash[:]))
						assert.Equal(t, ammAddress, output.Res.Amm)
						assert.Equal(t, shortOrder.BaseAssetQuantity, output.Res.UnfilledAmount)
					})
				})
				t.Run("for a partially filled Order", func(t *testing.T) {
					t.Run("for a longOrder it returns err = nil, with ammAddress and unfilled amount of cancelled Order", func(t *testing.T) {
						longOrder := getOrder(ammIndex, trader, longBaseAssetQuantity, price, salt, reduceOnly, postOnly)
						orderHash := getOrderHash(longOrder)
						filledAmount := hu.Div(longOrder.BaseAssetQuantity, big.NewInt(2))

						mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Placed)).Times(1)
						mockBibliophile.EXPECT().GetOrderFilledAmount(orderHash).Return(filledAmount).Times(1)
						mockBibliophile.EXPECT().GetMarketAddressFromMarketID(longOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)

						input := getValidateCancelLimitOrderInput(longOrder, trader, assertLowMargin)
						output := ValidateCancelLimitOrder(mockBibliophile, &input)
						assert.Equal(t, "", output.Err)
						assert.Equal(t, orderHash, common.BytesToHash(output.OrderHash[:]))
						assert.Equal(t, ammAddress, output.Res.Amm)
						expectedUnfilleAmount := hu.Sub(longOrder.BaseAssetQuantity, filledAmount)
						assert.Equal(t, expectedUnfilleAmount, output.Res.UnfilledAmount)
					})
					t.Run("for a shortOrder it returns err = nil, with ammAddress and unfilled amount of cancelled Order", func(t *testing.T) {
						shortOrder := getOrder(ammIndex, trader, shortBaseAssetQuantity, price, salt, reduceOnly, postOnly)
						orderHash := getOrderHash(shortOrder)
						filledAmount := hu.Div(shortOrder.BaseAssetQuantity, big.NewInt(2))

						mockBibliophile.EXPECT().GetOrderStatus(orderHash).Return(int64(Placed)).Times(1)
						mockBibliophile.EXPECT().GetOrderFilledAmount(orderHash).Return(filledAmount).Times(1)
						mockBibliophile.EXPECT().GetMarketAddressFromMarketID(shortOrder.AmmIndex.Int64()).Return(ammAddress).Times(1)

						input := getValidateCancelLimitOrderInput(shortOrder, trader, assertLowMargin)
						output := ValidateCancelLimitOrder(mockBibliophile, &input)
						assert.Equal(t, "", output.Err)
						assert.Equal(t, orderHash, common.BytesToHash(output.OrderHash[:]))
						assert.Equal(t, ammAddress, output.Res.Amm)
						expectedUnfilleAmount := hu.Sub(shortOrder.BaseAssetQuantity, filledAmount)
						assert.Equal(t, expectedUnfilleAmount, output.Res.UnfilledAmount)
					})
				})
			})
		})
	})
}

func getValidateCancelLimitOrderInput(order ILimitOrderBookOrder, sender common.Address, assertLowMargin bool) ValidateCancelLimitOrderInput {
	return ValidateCancelLimitOrderInput{
		Order:           order,
		Sender:          sender,
		AssertLowMargin: assertLowMargin,
	}
}

func getOrder(ammIndex *big.Int, trader common.Address, baseAssetQuantity *big.Int, price *big.Int, salt *big.Int, reduceOnly bool, postOnly bool) ILimitOrderBookOrder {
	return ILimitOrderBookOrder{
		AmmIndex:          ammIndex,
		BaseAssetQuantity: baseAssetQuantity,
		Trader:            trader,
		Price:             price,
		Salt:              salt,
		ReduceOnly:        reduceOnly,
		PostOnly:          postOnly,
	}
}

func getOrderHash(order ILimitOrderBookOrder) common.Hash {
	orderHash, err := GetLimitOrderHashFromContractStruct(&order)
	if err != nil {
		panic("error in getting order hash")
	}
	return orderHash
}
