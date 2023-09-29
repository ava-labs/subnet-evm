package ticks

import (
	"fmt"
	"math/big"

	"testing"

	hu "github.com/ava-labs/subnet-evm/plugin/evm/orderbook/hubbleutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	b "github.com/ava-labs/subnet-evm/precompile/contracts/bibliophile"
	gomock "github.com/golang/mock/gomock"
)

func TestGetPrevTick(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBibliophile := b.NewMockBibliophileClient(ctrl)
	ammAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	t.Run("when input tick price is 0", func(t *testing.T) {
		t.Run("For a bid", func(t *testing.T) {
			input := GetPrevTickInput{
				Amm:   ammAddress,
				Tick:  big.NewInt(0),
				IsBid: true,
			}
			output, err := GetPrevTick(mockBibliophile, input)
			assert.Equal(t, "tick price cannot be zero", err.Error())
			var expectedPrevTick *big.Int = nil
			assert.Equal(t, expectedPrevTick, output)
		})
		t.Run("For an ask", func(t *testing.T) {
			input := GetPrevTickInput{
				Amm:   ammAddress,
				Tick:  big.NewInt(0),
				IsBid: false,
			}
			output, err := GetPrevTick(mockBibliophile, input)
			assert.Equal(t, "tick price cannot be zero", err.Error())
			var expectedPrevTick *big.Int = nil
			assert.Equal(t, expectedPrevTick, output)
		})
	})
	t.Run("when input tick price > 0", func(t *testing.T) {
		t.Run("For a bid", func(t *testing.T) {
			bidsHead := big.NewInt(10000000) // 10
			t.Run("when bid price >= bidsHead", func(t *testing.T) {
				//covers bidsHead == 0
				t.Run("it returns error when bid price == bidsHead", func(t *testing.T) {
					input := GetPrevTickInput{
						Amm:   ammAddress,
						Tick:  big.NewInt(0).Set(bidsHead),
						IsBid: true,
					}
					mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
					prevTick, err := GetPrevTick(mockBibliophile, input)
					assert.Equal(t, fmt.Sprintf("tick %v is greater than or equal to bidsHead %v", input.Tick, bidsHead), err.Error())
					var expectedPrevTick *big.Int = nil
					assert.Equal(t, expectedPrevTick, prevTick)
				})
				t.Run("it returns error when bid price > bidsHead", func(t *testing.T) {
					input := GetPrevTickInput{
						Amm:   ammAddress,
						Tick:  big.NewInt(0).Add(bidsHead, big.NewInt(1)),
						IsBid: true,
					}
					mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
					prevTick, err := GetPrevTick(mockBibliophile, input)
					assert.Equal(t, fmt.Sprintf("tick %v is greater than or equal to bidsHead %v", input.Tick, bidsHead), err.Error())
					var expectedPrevTick *big.Int = nil
					assert.Equal(t, expectedPrevTick, prevTick)
				})
			})
			t.Run("when bid price < bidsHead", func(t *testing.T) {
				t.Run("when there is only 1 bid in orderbook", func(t *testing.T) {
					t.Run("it returns bidsHead as prevTick", func(t *testing.T) {
						input := GetPrevTickInput{
							Amm:   ammAddress,
							Tick:  big.NewInt(0).Div(bidsHead, big.NewInt(2)),
							IsBid: true,
						}
						mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
						mockBibliophile.EXPECT().GetNextBidPrice(input.Amm, bidsHead).Return(big.NewInt(0)).Times(1)
						prevTick, err := GetPrevTick(mockBibliophile, input)
						assert.Equal(t, nil, err)
						assert.Equal(t, bidsHead, prevTick)
					})
				})
				t.Run("when there are more than 1 bids in orderbook", func(t *testing.T) {
					bids := []*big.Int{big.NewInt(10000000), big.NewInt(9000000), big.NewInt(8000000), big.NewInt(7000000)}
					t.Run("when bid price does not match any bids in orderbook", func(t *testing.T) {
						t.Run("it returns prevTick when bid price falls between bids in orderbook", func(t *testing.T) {
							input := GetPrevTickInput{
								Amm:   ammAddress,
								Tick:  big.NewInt(8100000),
								IsBid: true,
							}
							mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
							mockBibliophile.EXPECT().GetNextBidPrice(input.Amm, bids[0]).Return(bids[1]).Times(1)
							mockBibliophile.EXPECT().GetNextBidPrice(input.Amm, bids[1]).Return(bids[2]).Times(1)
							prevTick, err := GetPrevTick(mockBibliophile, input)
							assert.Equal(t, nil, err)
							assert.Equal(t, bids[1], prevTick)
						})
						t.Run("it returns prevTick when bid price is lowest in orderbook", func(t *testing.T) {
							input := GetPrevTickInput{
								Amm:   ammAddress,
								Tick:  big.NewInt(400000),
								IsBid: true,
							}
							mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
							for i := 0; i < len(bids)-1; i++ {
								mockBibliophile.EXPECT().GetNextBidPrice(input.Amm, bids[i]).Return(bids[i+1]).Times(1)
							}
							mockBibliophile.EXPECT().GetNextBidPrice(input.Amm, bids[len(bids)-1]).Return(big.NewInt(0)).Times(1)
							prevTick, err := GetPrevTick(mockBibliophile, input)
							assert.Equal(t, nil, err)
							assert.Equal(t, bids[len(bids)-1], prevTick)
						})
					})
					t.Run("when bid price matches another bid's price in orderbook", func(t *testing.T) {
						t.Run("it returns prevTick", func(t *testing.T) {
							input := GetPrevTickInput{
								Amm:   ammAddress,
								Tick:  bids[2],
								IsBid: true,
							}
							mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
							mockBibliophile.EXPECT().GetNextBidPrice(input.Amm, bids[0]).Return(bids[1]).Times(1)
							mockBibliophile.EXPECT().GetNextBidPrice(input.Amm, bids[1]).Return(bids[2]).Times(1)
							prevTick, err := GetPrevTick(mockBibliophile, input)
							assert.Equal(t, nil, err)
							assert.Equal(t, bids[1], prevTick)
						})
					})
				})
			})
		})
		t.Run("For an ask", func(t *testing.T) {
			t.Run("when asksHead is 0", func(t *testing.T) {
				t.Run("it returns error", func(t *testing.T) {
					asksHead := big.NewInt(0)

					mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
					input := GetPrevTickInput{
						Amm:   ammAddress,
						Tick:  big.NewInt(10),
						IsBid: false,
					}
					prevTick, err := GetPrevTick(mockBibliophile, input)
					assert.Equal(t, "asksHead is zero", err.Error())
					var expectedPrevTick *big.Int = nil
					assert.Equal(t, expectedPrevTick, prevTick)
				})
			})
			t.Run("when asksHead > 0", func(t *testing.T) {
				asksHead := big.NewInt(10000000)
				t.Run("it returns error when ask price == asksHead", func(t *testing.T) {
					mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
					input := GetPrevTickInput{
						Amm:   ammAddress,
						Tick:  big.NewInt(0).Set(asksHead),
						IsBid: false,
					}
					prevTick, err := GetPrevTick(mockBibliophile, input)
					assert.Equal(t, fmt.Sprintf("tick %d is less than or equal to asksHead %d", input.Tick, asksHead), err.Error())
					var expectedPrevTick *big.Int = nil
					assert.Equal(t, expectedPrevTick, prevTick)
				})
				t.Run("it returns error when ask price < asksHead", func(t *testing.T) {
					mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
					input := GetPrevTickInput{
						Amm:   ammAddress,
						Tick:  big.NewInt(0).Sub(asksHead, big.NewInt(1)),
						IsBid: false,
					}
					prevTick, err := GetPrevTick(mockBibliophile, input)
					assert.Equal(t, fmt.Sprintf("tick %d is less than or equal to asksHead %d", input.Tick, asksHead), err.Error())
					var expectedPrevTick *big.Int = nil
					assert.Equal(t, expectedPrevTick, prevTick)
				})
				t.Run("when ask price > asksHead", func(t *testing.T) {
					t.Run("when there is only one ask in orderbook", func(t *testing.T) {
						t.Run("it returns asksHead as prevTick", func(t *testing.T) {
							mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
							input := GetPrevTickInput{
								Amm:   ammAddress,
								Tick:  big.NewInt(0).Add(asksHead, big.NewInt(1)),
								IsBid: false,
							}
							mockBibliophile.EXPECT().GetNextAskPrice(input.Amm, asksHead).Return(big.NewInt(0)).Times(1)
							prevTick, err := GetPrevTick(mockBibliophile, input)
							assert.Equal(t, nil, err)
							var expectedPrevTick *big.Int = asksHead
							assert.Equal(t, expectedPrevTick, prevTick)
						})
					})
					t.Run("when there are multiple asks in orderbook", func(t *testing.T) {
						asks := []*big.Int{asksHead, big.NewInt(11000000), big.NewInt(12000000), big.NewInt(13000000)}
						t.Run("when ask price does not match any asks in orderbook", func(t *testing.T) {
							t.Run("it returns prevTick when ask price falls between asks in orderbook", func(t *testing.T) {
								mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
								askPrice := big.NewInt(11500000)
								input := GetPrevTickInput{
									Amm:   ammAddress,
									Tick:  askPrice,
									IsBid: false,
								}
								mockBibliophile.EXPECT().GetNextAskPrice(input.Amm, asksHead).Return(asks[0]).Times(1)
								mockBibliophile.EXPECT().GetNextAskPrice(input.Amm, asks[0]).Return(asks[1]).Times(1)
								mockBibliophile.EXPECT().GetNextAskPrice(input.Amm, asks[1]).Return(asks[2]).Times(1)
								prevTick, err := GetPrevTick(mockBibliophile, input)
								assert.Equal(t, nil, err)
								var expectedPrevTick *big.Int = asks[1]
								assert.Equal(t, expectedPrevTick, prevTick)
							})
							t.Run("it returns prevTick when ask price is highest in orderbook", func(t *testing.T) {
								mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
								askPrice := big.NewInt(0).Add(asks[len(asks)-1], big.NewInt(1))
								input := GetPrevTickInput{
									Amm:   ammAddress,
									Tick:  askPrice,
									IsBid: false,
								}
								mockBibliophile.EXPECT().GetNextAskPrice(input.Amm, asksHead).Return(asks[0]).Times(1)
								mockBibliophile.EXPECT().GetNextAskPrice(input.Amm, asks[0]).Return(asks[1]).Times(1)
								mockBibliophile.EXPECT().GetNextAskPrice(input.Amm, asks[1]).Return(asks[2]).Times(1)
								mockBibliophile.EXPECT().GetNextAskPrice(input.Amm, asks[2]).Return(big.NewInt(0)).Times(1)
								prevTick, err := GetPrevTick(mockBibliophile, input)
								assert.Equal(t, nil, err)
								var expectedPrevTick *big.Int = asks[2]
								assert.Equal(t, expectedPrevTick, prevTick)
							})
						})
						t.Run("when ask price matches another ask's price in orderbook", func(t *testing.T) {
							t.Run("it returns prevTick", func(t *testing.T) {
								mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
								askPrice := asks[1]
								input := GetPrevTickInput{
									Amm:   ammAddress,
									Tick:  askPrice,
									IsBid: false,
								}
								mockBibliophile.EXPECT().GetNextAskPrice(input.Amm, asksHead).Return(asks[0]).Times(1)
								mockBibliophile.EXPECT().GetNextAskPrice(input.Amm, asks[0]).Return(asks[1]).Times(1)
								prevTick, err := GetPrevTick(mockBibliophile, input)
								assert.Equal(t, nil, err)
								var expectedPrevTick *big.Int = asks[0]
								assert.Equal(t, expectedPrevTick, prevTick)
							})
						})
					})
				})
			})
		})
	})
}

func TestSampleImpactBid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBibliophile := b.NewMockBibliophileClient(ctrl)
	ammAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	t.Run("when impactMarginNotional is zero", func(t *testing.T) {
		mockBibliophile.EXPECT().GetImpactMarginNotional(ammAddress).Return(big.NewInt(0)).Times(1)
		output := SampleImpactBid(mockBibliophile, ammAddress)
		assert.Equal(t, big.NewInt(0), output)
	})
	t.Run("when impactMarginNotional is > zero", func(t *testing.T) {
		impactMarginNotional := big.NewInt(4000000000) // 4000 units
		t.Run("when bidsHead is 0", func(t *testing.T) {
			mockBibliophile.EXPECT().GetImpactMarginNotional(ammAddress).Return(impactMarginNotional).Times(1)
			mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(big.NewInt(0)).Times(1)
			output := SampleImpactBid(mockBibliophile, ammAddress)
			assert.Equal(t, big.NewInt(0), output)
		})
		t.Run("when bidsHead > 0", func(t *testing.T) {
			bidsHead := big.NewInt(20000000) // 20 units
			t.Run("when bids in orderbook are not enough to cover impactMarginNotional", func(t *testing.T) {
				t.Run("when there is only one bid in orderbook it returns 0", func(t *testing.T) {
					mockBibliophile.EXPECT().GetImpactMarginNotional(ammAddress).Return(impactMarginNotional).Times(1)
					mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
					mockBibliophile.EXPECT().GetBidSize(ammAddress, bidsHead).Return(big.NewInt(1e18)).Times(1)
					mockBibliophile.EXPECT().GetNextBidPrice(ammAddress, bidsHead).Return(big.NewInt(0)).Times(1)
					output := SampleImpactBid(mockBibliophile, ammAddress)
					assert.Equal(t, big.NewInt(0), output)
				})
				t.Run("when there are multiple bids", func(t *testing.T) {
					bids := []*big.Int{bidsHead, big.NewInt(2100000), big.NewInt(2200000), big.NewInt(2300000)}
					size := big.NewInt(1e18) // 1 ether
					mockBibliophile.EXPECT().GetImpactMarginNotional(ammAddress).Return(impactMarginNotional).Times(1)
					mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
					for i := 0; i < len(bids); i++ {
						mockBibliophile.EXPECT().GetBidSize(ammAddress, bids[i]).Return(size).Times(1)
						if i != len(bids)-1 {
							mockBibliophile.EXPECT().GetNextBidPrice(ammAddress, bids[i]).Return(bids[i+1]).Times(1)
						} else {
							mockBibliophile.EXPECT().GetNextBidPrice(ammAddress, bids[i]).Return(big.NewInt(0)).Times(1)
						}
					}

					accumulatedMarginNotional := big.NewInt(0)
					for i := 0; i < len(bids); i++ {
						accumulatedMarginNotional.Add(accumulatedMarginNotional, hu.Div(hu.Mul(bids[i], size), big.NewInt(1e18)))
					}
					//asserting to check if testing conditions are setup correctly
					assert.Equal(t, -1, accumulatedMarginNotional.Cmp(impactMarginNotional))
					// accBaseQ := big.NewInt(0).Mul(size, big.NewInt(int64(len(bids))))
					// expectedSampleImpactBid := hu.Div(hu.Mul(accumulatedMarginNotional, big.NewInt(1e18)), accBaseQ)
					output := SampleImpactBid(mockBibliophile, ammAddress)
					assert.Equal(t, big.NewInt(0), output)
					// assert.Equal(t, expectedSampleImpactBid, output)
				})
			})
			t.Run("when bids in orderbook are enough to cover impactMarginNotional", func(t *testing.T) {
				t.Run("when there is only one bid in orderbook it returns bidsHead", func(t *testing.T) {
					bidsHead := impactMarginNotional
					mockBibliophile.EXPECT().GetImpactMarginNotional(ammAddress).Return(impactMarginNotional).Times(1)
					mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
					mockBibliophile.EXPECT().GetBidSize(ammAddress, bidsHead).Return(big.NewInt(1e18)).Times(1)
					output := SampleImpactBid(mockBibliophile, ammAddress)
					assert.Equal(t, bidsHead, output)
				})
				t.Run("when there are multiple bids, it tries to fill with available bids and average price is returned for rest", func(t *testing.T) {
					bidsHead := big.NewInt(2000000000) // 2000 units
					bids := []*big.Int{bidsHead}
					for i := int64(1); i < 6; i++ {
						bids = append(bids, big.NewInt(0).Sub(bidsHead, big.NewInt(i)))
					}
					size := big.NewInt(6e17) // 0.6 ether
					mockBibliophile.EXPECT().GetImpactMarginNotional(ammAddress).Return(impactMarginNotional).Times(1)
					mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
					mockBibliophile.EXPECT().GetNextBidPrice(ammAddress, bids[0]).Return(bids[1]).Times(1)
					mockBibliophile.EXPECT().GetNextBidPrice(ammAddress, bids[1]).Return(bids[2]).Times(1)
					mockBibliophile.EXPECT().GetNextBidPrice(ammAddress, bids[2]).Return(bids[3]).Times(1)
					mockBibliophile.EXPECT().GetBidSize(ammAddress, bids[0]).Return(size).Times(1)
					mockBibliophile.EXPECT().GetBidSize(ammAddress, bids[1]).Return(size).Times(1)
					mockBibliophile.EXPECT().GetBidSize(ammAddress, bids[2]).Return(size).Times(1)
					mockBibliophile.EXPECT().GetBidSize(ammAddress, bids[3]).Return(size).Times(1)

					output := SampleImpactBid(mockBibliophile, ammAddress)
					// 3 bids are filled and 3 are left
					totalBaseQ := big.NewInt(0).Mul(size, big.NewInt(3))
					filledQuote := big.NewInt(0)
					for i := 0; i < 3; i++ {
						filledQuote.Add(filledQuote, (hu.Div(hu.Mul(bids[i], size), big.NewInt(1e18))))
					}
					unfulFilledQuote := big.NewInt(0).Sub(impactMarginNotional, filledQuote)
					// as quantity is in 1e18 baseQ = price * 1e18 / price
					baseQAtTick := big.NewInt(0).Div(big.NewInt(0).Mul(unfulFilledQuote, big.NewInt(1e18)), bids[3])
					expectedOutput := big.NewInt(0).Div(big.NewInt(0).Mul(impactMarginNotional, big.NewInt(1e18)), big.NewInt(0).Add(totalBaseQ, baseQAtTick))
					assert.Equal(t, expectedOutput, output)
				})
			})
		})
	})
}

func TestSampleImpactAsk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBibliophile := b.NewMockBibliophileClient(ctrl)
	ammAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	t.Run("when impactMarginNotional is zero", func(t *testing.T) {
		mockBibliophile.EXPECT().GetImpactMarginNotional(ammAddress).Return(big.NewInt(0)).Times(1)
		output := SampleImpactAsk(mockBibliophile, ammAddress)
		assert.Equal(t, big.NewInt(0), output)
	})
	t.Run("when impactMarginNotional is > zero", func(t *testing.T) {
		impactMarginNotional := big.NewInt(4000000000) // 4000 units
		t.Run("when asksHead is 0", func(t *testing.T) {
			mockBibliophile.EXPECT().GetImpactMarginNotional(ammAddress).Return(impactMarginNotional).Times(1)
			mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(big.NewInt(0)).Times(1)
			output := SampleImpactAsk(mockBibliophile, ammAddress)
			assert.Equal(t, big.NewInt(0), output)
		})
		t.Run("when asksHead > 0", func(t *testing.T) {
			asksHead := big.NewInt(20000000) // 20 units
			t.Run("when asks in orderbook are not enough to cover impactMarginNotional", func(t *testing.T) {
				t.Run("when there is only one ask in orderbook it returns asksHead", func(t *testing.T) {
					mockBibliophile.EXPECT().GetImpactMarginNotional(ammAddress).Return(impactMarginNotional).Times(1)
					mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
					mockBibliophile.EXPECT().GetAskSize(ammAddress, asksHead).Return(big.NewInt(1e18)).Times(1)
					mockBibliophile.EXPECT().GetNextAskPrice(ammAddress, asksHead).Return(big.NewInt(0)).Times(1)
					output := SampleImpactAsk(mockBibliophile, ammAddress)
					assert.Equal(t, big.NewInt(0), output)
				})
				t.Run("when there are multiple asks", func(t *testing.T) {
					asks := []*big.Int{asksHead, big.NewInt(2100000), big.NewInt(2200000), big.NewInt(2300000)}
					size := big.NewInt(1e18) // 1 ether
					mockBibliophile.EXPECT().GetImpactMarginNotional(ammAddress).Return(impactMarginNotional).Times(1)
					mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
					for i := 0; i < len(asks); i++ {
						mockBibliophile.EXPECT().GetAskSize(ammAddress, asks[i]).Return(size).Times(1)
						if i != len(asks)-1 {
							mockBibliophile.EXPECT().GetNextAskPrice(ammAddress, asks[i]).Return(asks[i+1]).Times(1)
						} else {
							mockBibliophile.EXPECT().GetNextAskPrice(ammAddress, asks[i]).Return(big.NewInt(0)).Times(1)
						}
					}

					accumulatedMarginNotional := big.NewInt(0)
					for i := 0; i < len(asks); i++ {
						accumulatedMarginNotional.Add(accumulatedMarginNotional, hu.Div(hu.Mul(asks[i], size), big.NewInt(1e18)))
					}
					//asserting to check if testing conditions are setup correctly
					assert.Equal(t, -1, accumulatedMarginNotional.Cmp(impactMarginNotional))
					output := SampleImpactAsk(mockBibliophile, ammAddress)
					assert.Equal(t, big.NewInt(0), output)
				})
			})
			t.Run("when asks in orderbook are enough to cover impactMarginNotional", func(t *testing.T) {
				t.Run("when there is only one ask in orderbook it returns asksHead", func(t *testing.T) {
					newAsksHead := impactMarginNotional
					mockBibliophile.EXPECT().GetImpactMarginNotional(ammAddress).Return(impactMarginNotional).Times(1)
					mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(newAsksHead).Times(1)
					mockBibliophile.EXPECT().GetAskSize(ammAddress, newAsksHead).Return(big.NewInt(1e18)).Times(1)
					output := SampleImpactAsk(mockBibliophile, ammAddress)
					assert.Equal(t, newAsksHead, output)
				})
				t.Run("when there are multiple asks, it tries to fill with available asks and average price is returned for rest", func(t *testing.T) {
					newAsksHead := big.NewInt(2000000000) // 2000 units
					asks := []*big.Int{newAsksHead}
					for i := int64(1); i < 6; i++ {
						asks = append(asks, big.NewInt(0).Add(newAsksHead, big.NewInt(i)))
					}
					size := big.NewInt(6e17) // 0.6 ether
					mockBibliophile.EXPECT().GetImpactMarginNotional(ammAddress).Return(impactMarginNotional).Times(1)
					mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(newAsksHead).Times(1)
					mockBibliophile.EXPECT().GetNextAskPrice(ammAddress, asks[0]).Return(asks[1]).Times(1)
					mockBibliophile.EXPECT().GetNextAskPrice(ammAddress, asks[1]).Return(asks[2]).Times(1)
					mockBibliophile.EXPECT().GetNextAskPrice(ammAddress, asks[2]).Return(asks[3]).Times(1)
					mockBibliophile.EXPECT().GetAskSize(ammAddress, asks[0]).Return(size).Times(1)
					mockBibliophile.EXPECT().GetAskSize(ammAddress, asks[1]).Return(size).Times(1)
					mockBibliophile.EXPECT().GetAskSize(ammAddress, asks[2]).Return(size).Times(1)
					mockBibliophile.EXPECT().GetAskSize(ammAddress, asks[3]).Return(size).Times(1)

					// 2000 * .6 + 2001 * .6 + 2002 * .6 = 3,601.8
					// 3 asks are filled and 3 are left
					accBaseQ := big.NewInt(0).Mul(size, big.NewInt(3))
					filledQuote := big.NewInt(0)
					for i := 0; i < 3; i++ {
						filledQuote.Add(filledQuote, hu.Div1e6(big.NewInt(0).Mul(asks[i], size)))
					}
					_impactMarginNotional := new(big.Int).Mul(impactMarginNotional, big.NewInt(1e12))
					baseQAtTick := new(big.Int).Div(hu.Mul1e6(big.NewInt(0).Sub(_impactMarginNotional, filledQuote)), asks[3])
					expectedOutput := new(big.Int).Div(hu.Mul1e6(_impactMarginNotional), new(big.Int).Add(baseQAtTick, accBaseQ))
					assert.Equal(t, expectedOutput, SampleImpactAsk(mockBibliophile, ammAddress))
				})
			})
		})
	})
}

func TestSampleBid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBibliophile := b.NewMockBibliophileClient(ctrl)
	ammAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	bidsHead := big.NewInt(20 * 1e6)      // $20
	baseAssetQuantity := big.NewInt(1e18) // 1 ether
	t.Run("when bidsHead is 0", func(t *testing.T) {
		mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(big.NewInt(0)).Times(1)
		output := _sampleBid(mockBibliophile, ammAddress, baseAssetQuantity)
		assert.Equal(t, big.NewInt(0), output)
	})
	t.Run("when bidsHead > 0", func(t *testing.T) {
		t.Run("when bids in orderbook are not enough to cover baseAssetQuantity", func(t *testing.T) {
			t.Run("when there is only one bid in orderbook it returns 0", func(t *testing.T) {
				mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
				mockBibliophile.EXPECT().GetBidSize(ammAddress, bidsHead).Return(hu.Sub(baseAssetQuantity, big.NewInt(1))).Times(1)
				mockBibliophile.EXPECT().GetNextBidPrice(ammAddress, bidsHead).Return(big.NewInt(0)).Times(1)
				output := _sampleBid(mockBibliophile, ammAddress, baseAssetQuantity)
				assert.Equal(t, big.NewInt(0), output)
			})
			t.Run("when there are multiple bids", func(t *testing.T) {
				bids := []*big.Int{bidsHead, big.NewInt(2100000), big.NewInt(2200000), big.NewInt(2300000)}
				size := big.NewInt(24 * 1e16) // 0.24 ether
				mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
				for i := 0; i < len(bids); i++ {
					mockBibliophile.EXPECT().GetBidSize(ammAddress, bids[i]).Return(size).Times(1)
					if i != len(bids)-1 {
						mockBibliophile.EXPECT().GetNextBidPrice(ammAddress, bids[i]).Return(bids[i+1]).Times(1)
					} else {
						mockBibliophile.EXPECT().GetNextBidPrice(ammAddress, bids[i]).Return(big.NewInt(0)).Times(1)
					}
				}
				output := _sampleBid(mockBibliophile, ammAddress, baseAssetQuantity)
				assert.Equal(t, big.NewInt(0), output)
			})
		})
		t.Run("when bids in orderbook are enough to cover baseAssetQuantity", func(t *testing.T) {
			t.Run("when there is only one bid in orderbook it returns bidsHead", func(t *testing.T) {
				mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
				mockBibliophile.EXPECT().GetBidSize(ammAddress, bidsHead).Return(baseAssetQuantity).Times(1)
				output := _sampleBid(mockBibliophile, ammAddress, baseAssetQuantity)
				assert.Equal(t, bidsHead, output)
			})
			t.Run("when there are multiple bids, it tries to fill with available bids and average price is returned for rest", func(t *testing.T) {
				bids := []*big.Int{bidsHead}
				for i := int64(1); i < 6; i++ {
					bids = append(bids, hu.Sub(bidsHead, big.NewInt(i)))
				}
				size := big.NewInt(3e17) // 0.3 ether
				mockBibliophile.EXPECT().GetBidsHead(ammAddress).Return(bidsHead).Times(1)
				mockBibliophile.EXPECT().GetNextBidPrice(ammAddress, bids[0]).Return(bids[1]).Times(1)
				mockBibliophile.EXPECT().GetNextBidPrice(ammAddress, bids[1]).Return(bids[2]).Times(1)
				mockBibliophile.EXPECT().GetNextBidPrice(ammAddress, bids[2]).Return(bids[3]).Times(1)
				mockBibliophile.EXPECT().GetBidSize(ammAddress, bids[0]).Return(size).Times(1)
				mockBibliophile.EXPECT().GetBidSize(ammAddress, bids[1]).Return(size).Times(1)
				mockBibliophile.EXPECT().GetBidSize(ammAddress, bids[2]).Return(size).Times(1)
				mockBibliophile.EXPECT().GetBidSize(ammAddress, bids[3]).Return(size).Times(1)

				output := _sampleBid(mockBibliophile, ammAddress, baseAssetQuantity)
				accBaseQ := hu.Mul(size, big.NewInt(3))
				accNotional := big.NewInt(0)
				for i := 0; i < 3; i++ {
					accNotional.Add(accNotional, (hu.Div1e6(hu.Mul(bids[i], size))))
				}
				notionalAtTick := hu.Div1e6(hu.Mul(hu.Sub(baseAssetQuantity, accBaseQ), bids[3]))
				expectedOutput := hu.Div(hu.Mul1e6(hu.Add(accNotional, notionalAtTick)), baseAssetQuantity)
				assert.Equal(t, expectedOutput, output)
			})
		})
	})
}

func TestSampleAsk(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBibliophile := b.NewMockBibliophileClient(ctrl)
	ammAddress := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	asksHead := big.NewInt(20 * 1e6)      // $20
	baseAssetQuantity := big.NewInt(1e18) // 1 ether
	t.Run("when asksHead is 0", func(t *testing.T) {
		mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(big.NewInt(0)).Times(1)
		output := _sampleAsk(mockBibliophile, ammAddress, baseAssetQuantity)
		assert.Equal(t, big.NewInt(0), output)
	})
	t.Run("when asksHead > 0", func(t *testing.T) {
		t.Run("when asks in orderbook are not enough to cover baseAssetQuantity", func(t *testing.T) {
			t.Run("when there is only one ask in orderbook it returns 0", func(t *testing.T) {
				mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
				mockBibliophile.EXPECT().GetAskSize(ammAddress, asksHead).Return(hu.Sub(baseAssetQuantity, big.NewInt(1))).Times(1)
				mockBibliophile.EXPECT().GetNextAskPrice(ammAddress, asksHead).Return(big.NewInt(0)).Times(1)
				output := _sampleAsk(mockBibliophile, ammAddress, baseAssetQuantity)
				assert.Equal(t, big.NewInt(0), output)
			})
			t.Run("when there are multiple asks, it tries to fill with available asks", func(t *testing.T) {
				asks := []*big.Int{asksHead, big.NewInt(2100000), big.NewInt(2200000), big.NewInt(2300000)}
				size := big.NewInt(24 * 1e16) // 0.24 ether
				mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
				for i := 0; i < len(asks); i++ {
					mockBibliophile.EXPECT().GetAskSize(ammAddress, asks[i]).Return(size).Times(1)
					if i != len(asks)-1 {
						mockBibliophile.EXPECT().GetNextAskPrice(ammAddress, asks[i]).Return(asks[i+1]).Times(1)
					} else {
						mockBibliophile.EXPECT().GetNextAskPrice(ammAddress, asks[i]).Return(big.NewInt(0)).Times(1)
					}
				}
				output := _sampleAsk(mockBibliophile, ammAddress, baseAssetQuantity)
				assert.Equal(t, big.NewInt(0), output)
			})
		})
		t.Run("when asks in orderbook are enough to cover baseAssetQuantity", func(t *testing.T) {
			t.Run("when there is only one ask in orderbook it returns asksHead", func(t *testing.T) {
				mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
				mockBibliophile.EXPECT().GetAskSize(ammAddress, asksHead).Return(baseAssetQuantity).Times(1)
				output := _sampleAsk(mockBibliophile, ammAddress, baseAssetQuantity)
				assert.Equal(t, asksHead, output)
			})
			t.Run("when there are multiple asks, it tries to fill with available asks and average price is returned for rest", func(t *testing.T) {
				asks := []*big.Int{asksHead}
				for i := int64(1); i < 6; i++ {
					asks = append(asks, hu.Sub(asksHead, big.NewInt(i)))
				}
				size := big.NewInt(31e16) // 0.31 ether
				mockBibliophile.EXPECT().GetAsksHead(ammAddress).Return(asksHead).Times(1)
				mockBibliophile.EXPECT().GetNextAskPrice(ammAddress, asks[0]).Return(asks[1]).Times(1)
				mockBibliophile.EXPECT().GetNextAskPrice(ammAddress, asks[1]).Return(asks[2]).Times(1)
				mockBibliophile.EXPECT().GetNextAskPrice(ammAddress, asks[2]).Return(asks[3]).Times(1)
				mockBibliophile.EXPECT().GetAskSize(ammAddress, asks[0]).Return(size).Times(1)
				mockBibliophile.EXPECT().GetAskSize(ammAddress, asks[1]).Return(size).Times(1)
				mockBibliophile.EXPECT().GetAskSize(ammAddress, asks[2]).Return(size).Times(1)
				mockBibliophile.EXPECT().GetAskSize(ammAddress, asks[3]).Return(size).Times(1)

				output := _sampleAsk(mockBibliophile, ammAddress, baseAssetQuantity)
				accBaseQ := hu.Mul(size, big.NewInt(3))
				accNotional := big.NewInt(0)
				for i := 0; i < 3; i++ {
					accNotional.Add(accNotional, (hu.Div1e6(hu.Mul(asks[i], size))))
				}
				notionalAtTick := hu.Div1e6(hu.Mul(hu.Sub(baseAssetQuantity, accBaseQ), asks[3]))
				expectedOutput := hu.Div(hu.Mul1e6(hu.Add(accNotional, notionalAtTick)), baseAssetQuantity)
				assert.Equal(t, expectedOutput, output)
			})
		})
	})
}
