package limitorders

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestGetLiquidableTraders(t *testing.T) {
	spreadRatioThreshold = big.NewInt(2e5) // this assumption has been made in the test cases
	t.Run("When no trader exist", func(t *testing.T) {
		var market Market = 1
		traderMap := map[common.Address]Trader{}
		markPrice := multiplyBasePrecision(big.NewInt(100))
		oraclePrice := multiplyBasePrecision(big.NewInt(110))
		liquidablePositions := GetLiquidableTraders(traderMap, market, markPrice, oraclePrice)
		assert.Equal(t, 0, len(liquidablePositions))
	})

	t.Run("When traders exist", func(t *testing.T) {
		t.Run("When no trader has any positions", func(t *testing.T) {
			var market Market = 1
			longTraderAddress := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
			collateral := HUSD
			margin := big.NewInt(10000000000)
			traderMap := map[common.Address]Trader{
				longTraderAddress: Trader{
					Margin: Margin{Deposited: map[Collateral]*big.Int{
						collateral: margin,
					}},
					Positions: map[Market]*Position{},
				},
			}
			markPrice := multiplyBasePrecision(big.NewInt(100))
			oraclePrice := multiplyBasePrecision(big.NewInt(110))
			liquidablePositions := GetLiquidableTraders(traderMap, market, markPrice, oraclePrice)
			assert.Equal(t, 0, len(liquidablePositions))
		})
		t.Run("When traders have positions", func(t *testing.T) {
			t.Run("When mark price is within 20% of oracle price, it uses mark price for calculating margin fraction", func(t *testing.T) {
				markPrice := multiplyBasePrecision(big.NewInt(100))
				oraclePrice := multiplyBasePrecision(big.NewInt(110))
				t.Run("When traders margin fraction is >= than maintenance margin, GetLiquidableTraders returns empty array", func(t *testing.T) {
					var market Market = 1
					collateral := HUSD

					//long trader
					longTraderAddress := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
					marginLong := multiplyBasePrecision(big.NewInt(500))
					longSize := multiplyPrecisionSize(big.NewInt(10))
					longEntryPrice := multiplyBasePrecision(big.NewInt(90))
					openNotionalLong := dividePrecisionSize(big.NewInt(0).Mul(longEntryPrice, longSize))
					longTrader := Trader{
						Margin: Margin{Deposited: map[Collateral]*big.Int{
							collateral: marginLong,
						}},
						Positions: map[Market]*Position{
							market: getPosition(market, openNotionalLong, longSize, big.NewInt(0), big.NewInt(0), big.NewInt(0)),
						},
					}

					//short trader
					shortTraderAddress := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
					marginShort := multiplyBasePrecision(big.NewInt(1000))
					shortSize := multiplyPrecisionSize(big.NewInt(-20))
					shortEntryPrice := multiplyBasePrecision(big.NewInt(105))
					openNotionalShort := dividePrecisionSize(big.NewInt(0).Abs(big.NewInt(0).Mul(shortEntryPrice, shortSize)))
					shortTrader := Trader{
						Margin: Margin{Deposited: map[Collateral]*big.Int{
							collateral: marginShort,
						}},
						Positions: map[Market]*Position{
							market: getPosition(market, openNotionalShort, shortSize, big.NewInt(0), big.NewInt(0), big.NewInt(0)),
						},
					}
					traderMap := map[common.Address]Trader{
						longTraderAddress:  longTrader,
						shortTraderAddress: shortTrader,
					}

					//long margin fraction - ((500 +(100-90)*10)*1e6/(100*10) = 600000 > maintenanceMargin(1e5)
					//short margin fraction - ((1000 + (105-100)*20)*1e6/(20*100) = 550000 > maintenanceMargin(1e5)
					expectedLiquidablePositionsCount := 0
					liquidablePositions := GetLiquidableTraders(traderMap, market, markPrice, oraclePrice)
					assert.Equal(t, expectedLiquidablePositionsCount, len(liquidablePositions))
				})
				t.Run("When trader margin fraction is < than maintenance margin, it returns trader's info in GetLiquidableTraders sorted by marginFraction", func(t *testing.T) {
					var market Market = 1
					collateral := HUSD

					//long trader
					longTraderAddress := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
					marginLong := multiplyBasePrecision(big.NewInt(500))
					longSize := multiplyPrecisionSize(big.NewInt(10))
					longEntryPrice := multiplyBasePrecision(big.NewInt(145))
					openNotionalLong := dividePrecisionSize(big.NewInt(0).Mul(longEntryPrice, longSize))
					longTrader := Trader{
						Margin: Margin{Deposited: map[Collateral]*big.Int{
							collateral: marginLong,
						}},
						Positions: map[Market]*Position{
							market: getPosition(market, openNotionalLong, longSize, big.NewInt(0), big.NewInt(0), big.NewInt(0)),
						},
					}

					//short trader
					shortTraderAddress := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
					marginShort := multiplyBasePrecision(big.NewInt(500))
					shortSize := multiplyPrecisionSize(big.NewInt(-20))
					shortEntryPrice := multiplyBasePrecision(big.NewInt(80))
					openNotionalShort := dividePrecisionSize(big.NewInt(0).Abs(big.NewInt(0).Mul(shortEntryPrice, shortSize)))
					shortTrader := Trader{
						Margin: Margin{Deposited: map[Collateral]*big.Int{
							collateral: marginShort,
						}},
						Positions: map[Market]*Position{
							market: getPosition(market, openNotionalShort, shortSize, big.NewInt(0), big.NewInt(0), big.NewInt(0)),
						},
					}
					traderMap := map[common.Address]Trader{
						longTraderAddress:  longTrader,
						shortTraderAddress: shortTrader,
					}

					liquidablePositions := GetLiquidableTraders(traderMap, market, markPrice, oraclePrice)

					assert.Equal(t, 2, len(liquidablePositions))
					//oldNotional := 1450 * 1e6-> (longEntryPrice * longSize)/1e18
					//unrealizedPnl :=  -450 * 1e6 -> longSize(markPrice - longEntryPrice)/1e18
					//effectiveMarginLong1 := 500 - 450 = 50 * 1e6 -> margin + unrealizedPnl
					//newNotional := 1000 * 1e6 -> (markPrice * longSize1)/1e18
					//expectedMarginFractionLong1 = effectiveMarginLong1*1e6/newNotional
					expectedMarginFractionLong := big.NewInt(50000)
					//(1e6(margin*1e18 + shortSize.Abs*(shortEntryPrice-markPrice)))/(shortSize*markPrice)

					// Add short position 1
					//oldNotional := 1600 * 1e6 -> (ShortEntryPrice * ShortSize)/1e18
					//unrealizedPnl := -400 * 1e6 -> ShortSize1(markPrice - ShortEntryPrice1)/1e18
					//effectiveMarginShort1 := 100 * 1e6 -> margin + unrealizedPnl
					//newNotional := 2000 * 1e6 -> (markPrice * ShortSize1)/1e18
					//expectedMarginFractionShort1 = effectiveMarginShort1*1e6/newNotional
					expectedMarginFractionShort := big.NewInt(50000)

					//both mfs are same so liquidable position order will same as traderMap so long liquidable comes first
					assert.Equal(t, longTraderAddress, liquidablePositions[0].Address)
					assert.Equal(t, getLiquidationThreshold(longSize), liquidablePositions[0].Size)
					assert.Equal(t, expectedMarginFractionLong, liquidablePositions[0].MarginFraction)
					assert.Equal(t, shortTraderAddress, liquidablePositions[1].Address)
					assert.Equal(t, getLiquidationThreshold(shortSize), liquidablePositions[1].Size)
					assert.Equal(t, expectedMarginFractionShort, liquidablePositions[1].MarginFraction)
				})
			})
			t.Run("When mark price is outside of 20% of oracle price, it also uses oracle price for calculating margin fraction", func(t *testing.T) {
				t.Run("When trader margin fraction is >= than maintenance margin", func(t *testing.T) {
					markPrice := multiplyBasePrecision(big.NewInt(75))
					oraclePrice := multiplyBasePrecision(big.NewInt(100))
					var market Market = 1
					collateral := HUSD

					//long position for trader 1
					longTraderAddress := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
					marginLong := multiplyBasePrecision(big.NewInt(500))
					longSize := multiplyPrecisionSize(big.NewInt(10))
					longEntryPrice := multiplyBasePrecision(big.NewInt(90))
					openNotionalLong := dividePrecisionSize(big.NewInt(0).Mul(longEntryPrice, longSize))
					longTrader := Trader{
						Margin: Margin{Deposited: map[Collateral]*big.Int{
							collateral: marginLong,
						}},
						Positions: map[Market]*Position{
							market: getPosition(market, openNotionalLong, longSize, big.NewInt(0), big.NewInt(0), big.NewInt(0)),
						},
					}

					//short Position for trader 2
					shortTraderAddress := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
					marginShort := multiplyBasePrecision(big.NewInt(1000))
					// open price for short is 2100/20= 105 so trader 2 is in loss
					shortSize := multiplyPrecisionSize(big.NewInt(-20))
					shortEntryPrice := multiplyBasePrecision(big.NewInt(105))
					openNotionalShort := dividePrecisionSize(big.NewInt(0).Abs(big.NewInt(0).Mul(shortEntryPrice, shortSize)))
					shortTrader := Trader{
						Margin: Margin{Deposited: map[Collateral]*big.Int{
							collateral: marginShort,
						}},
						Positions: map[Market]*Position{
							market: getPosition(market, openNotionalShort, shortSize, big.NewInt(0), big.NewInt(0), big.NewInt(0)),
						},
					}
					traderMap := map[common.Address]Trader{
						longTraderAddress:  longTrader,
						shortTraderAddress: shortTrader,
					}

					//long margin fraction - ((500 +(100-90)*10)*1e6/(100*10) = 600000 > maintenanceMargin(1e5)
					//short margin fraction - ((1000 + (105-75)*20)*1e6/(20*75) = 1066666 > maintenanceMargin(1e5)
					liquidablePositions := GetLiquidableTraders(traderMap, market, markPrice, oraclePrice)
					assert.Equal(t, 0, len(liquidablePositions))
				})
				t.Run("When trader margin fraction is < than maintenance margin, it returns trader's info in GetLiquidableTraders", func(t *testing.T) {
					t.Run("When mf-markPrice > mf-oraclePrice, it uses mf with mark price", func(t *testing.T) {
						t.Run("For long order", func(t *testing.T) {
							// for both long mf-markPrice will > mf-oraclePrice
							markPrice := multiplyBasePrecision(big.NewInt(140))
							oraclePrice := multiplyBasePrecision(big.NewInt(110))
							var market Market = 1
							collateral := HUSD

							//long trader 1
							longTraderAddress1 := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
							marginLong1 := multiplyBasePrecision(big.NewInt(500))
							longSize1 := multiplyPrecisionSize(big.NewInt(10))
							longEntryPrice1 := multiplyBasePrecision(big.NewInt(180))
							openNotionalLong1 := dividePrecisionSize(big.NewInt(0).Mul(longEntryPrice1, longSize1))
							longTrader1 := Trader{
								Margin: Margin{Deposited: map[Collateral]*big.Int{
									collateral: marginLong1,
								}},
								Positions: map[Market]*Position{
									market: getPosition(market, openNotionalLong1, longSize1, big.NewInt(0), big.NewInt(0), big.NewInt(0)),
								},
							}

							//long trader 2
							longTraderAddress2 := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
							marginLong2 := multiplyBasePrecision(big.NewInt(500))
							longSize2 := multiplyPrecisionSize(big.NewInt(10))
							longEntryPrice2 := multiplyBasePrecision(big.NewInt(145))
							openNotionalLong2 := dividePrecisionSize(big.NewInt(0).Mul(longEntryPrice2, longSize2))
							longTrader2 := Trader{
								Margin: Margin{Deposited: map[Collateral]*big.Int{
									collateral: marginLong2,
								}},
								Positions: map[Market]*Position{
									market: getPosition(market, openNotionalLong2, longSize2, big.NewInt(0), big.NewInt(0), big.NewInt(0)),
								},
							}

							traderMap := map[common.Address]Trader{
								longTraderAddress1: longTrader1,
								longTraderAddress2: longTrader2,
							}

							liquidablePositions := GetLiquidableTraders(traderMap, market, markPrice, oraclePrice)
							assert.Equal(t, 1, len(liquidablePositions))

							//long trader 1 mf-markPrice < maintenanceMargin so it is liquidated
							//long trader 2 mf-markPrice > maintenanceMargin so it is not liquidated

							//oldNotional := 1800000000 -> (longEntryPrice1 * longSize1)/1e18
							//unrealizedPnl := -400000000 -> longSize1(markPrice - longEntryPrice1)/1e18
							//effectiveMarginLong1 := 100000000 -> margin + unrealizedPnl
							//newNotional := 1400000000 -> (markPrice * longSize1)/1e18
							//expectedMarginFractionLong1 = effectiveMarginLong1*1e6/newNotional
							expectedMarginFractionLong1 := big.NewInt(71428)
							assert.Equal(t, longTraderAddress1, liquidablePositions[0].Address)
							assert.Equal(t, getLiquidationThreshold(longSize1), liquidablePositions[0].Size)
							assert.Equal(t, expectedMarginFractionLong1, liquidablePositions[0].MarginFraction)
						})
						t.Run("For short order", func(t *testing.T) {
							markPrice := multiplyBasePrecision(big.NewInt(110))
							oraclePrice := multiplyBasePrecision(big.NewInt(140))
							var market Market = 1
							collateral := HUSD

							//short trader 1
							shortTraderAddress1 := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
							marginShort1 := multiplyBasePrecision(big.NewInt(500))
							// Add short position 1
							shortSize1 := multiplyPrecisionSize(big.NewInt(-20))
							shortEntryPrice1 := multiplyBasePrecision(big.NewInt(80))
							openNotionalShort1 := dividePrecisionSize(big.NewInt(0).Abs(big.NewInt(0).Mul(shortEntryPrice1, shortSize1)))
							shortTrader1 := Trader{
								Margin: Margin{Deposited: map[Collateral]*big.Int{
									collateral: marginShort1,
								}},
								Positions: map[Market]*Position{
									market: getPosition(market, openNotionalShort1, shortSize1, big.NewInt(0), big.NewInt(0), big.NewInt(0)),
								},
							}
							//short trader 2
							shortTraderAddress2 := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
							marginShort2 := multiplyBasePrecision(big.NewInt(500))
							shortSize2 := multiplyPrecisionSize(big.NewInt(-20))
							shortEntryPrice2 := multiplyBasePrecision(big.NewInt(100))
							openNotionalShort2 := dividePrecisionSize(big.NewInt(0).Abs(big.NewInt(0).Mul(shortEntryPrice2, shortSize2)))
							shortTrader2 := Trader{
								Margin: Margin{Deposited: map[Collateral]*big.Int{
									collateral: marginShort2,
								}},
								Positions: map[Market]*Position{
									market: getPosition(market, openNotionalShort2, shortSize2, big.NewInt(0), big.NewInt(0), big.NewInt(0)),
								},
							}
							traderMap := map[common.Address]Trader{
								shortTraderAddress1: shortTrader1,
								shortTraderAddress2: shortTrader2,
							}

							liquidablePositions := GetLiquidableTraders(traderMap, market, markPrice, oraclePrice)

							//Short trader 1 mf-markPrice < maintenanceMargin so it is liquidated
							//Short trader 2 mf-markPrice > maintenanceMargin so it is notliquidated

							//oldNotional := 1600000000 -> (ShortEntryPrice1 * ShortSize1)/1e18
							//unrealizedPnl := -600000000 -> ShortSize1(markPrice - ShortEntryPrice1)/1e18
							//effectiveMarginShort1 := -10000000 -> margin + unrealizedPnl
							//newNotional := 2800000000 -> (markPrice * ShortSize1)/1e18
							//expectedMarginFractionShort1 = effectiveMarginShort1*1e6/newNotional
							expectedMarginFractionShort1 := big.NewInt(0)

							assert.Equal(t, 1, len(liquidablePositions))
							assert.Equal(t, shortTraderAddress1, liquidablePositions[0].Address)
							assert.Equal(t, getLiquidationThreshold(shortSize1), liquidablePositions[0].Size)
							assert.Equal(t, expectedMarginFractionShort1, liquidablePositions[0].MarginFraction)
						})
					})
					t.Run("When mf-markPrice < mf-oraclePrice, it uses mf with oracle price", func(t *testing.T) {
						t.Run("For long order", func(t *testing.T) {
							// for both long mf-markPrice will > mf-oraclePrice
							markPrice := multiplyBasePrecision(big.NewInt(110))
							oraclePrice := multiplyBasePrecision(big.NewInt(140))
							db := NewInMemoryDatabase()
							var market Market = 1
							db.LastPrice[market] = markPrice
							collateral := HUSD

							//long trader 1
							longTraderAddress1 := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
							marginLong1 := multiplyBasePrecision(big.NewInt(500))
							longSize1 := multiplyPrecisionSize(big.NewInt(10))
							longEntryPrice1 := multiplyBasePrecision(big.NewInt(180))
							openNotionalLong1 := dividePrecisionSize(big.NewInt(0).Mul(longEntryPrice1, longSize1))
							longTrader1 := Trader{
								Margin: Margin{Deposited: map[Collateral]*big.Int{
									collateral: marginLong1,
								}},
								Positions: map[Market]*Position{
									market: getPosition(market, openNotionalLong1, longSize1, big.NewInt(0), big.NewInt(0), big.NewInt(0)),
								},
							}

							//long trader 2
							longTraderAddress2 := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
							marginLong2 := multiplyBasePrecision(big.NewInt(500))
							longSize2 := multiplyPrecisionSize(big.NewInt(10))
							longEntryPrice2 := multiplyBasePrecision(big.NewInt(145))
							openNotionalLong2 := dividePrecisionSize(big.NewInt(0).Mul(longEntryPrice2, longSize2))
							longTrader2 := Trader{
								Margin: Margin{Deposited: map[Collateral]*big.Int{
									collateral: marginLong2,
								}},
								Positions: map[Market]*Position{
									market: getPosition(market, openNotionalLong2, longSize2, big.NewInt(0), big.NewInt(0), big.NewInt(0)),
								},
							}
							traderMap := map[common.Address]Trader{
								longTraderAddress1: longTrader1,
								longTraderAddress2: longTrader2,
							}

							liquidablePositions := GetLiquidableTraders(traderMap, market, markPrice, oraclePrice)
							assert.Equal(t, 1, len(liquidablePositions))

							//long trader 1 mf-oraclePrice < maintenanceMargin so it is liquidated
							//long trader 2 mf-oraclePrice > maintenanceMargin so it is notliquidated
							//oldNotional := 1800000000 -> (longEntryPrice1 * longSize1)/1e18
							//unrealizedPnl := -400000000 -> longSize1(markPrice - longEntryPrice1)/1e18
							//effectiveMarginLong1 := 100000000 -> margin + unrealizedPnl
							//newNotional := 1400000000 -> (markPrice * longSize1)/1e18
							//expectedMarginFractionLong1 = effectiveMarginLong1*1e6/newNotional
							expectedMarginFractionLong1 := big.NewInt(71428)
							assert.Equal(t, longTraderAddress1, liquidablePositions[0].Address)
							assert.Equal(t, getLiquidationThreshold(longSize1), liquidablePositions[0].Size)
							assert.Equal(t, expectedMarginFractionLong1, liquidablePositions[0].MarginFraction)
						})
						t.Run("For short order", func(t *testing.T) {
							markPrice := multiplyBasePrecision(big.NewInt(140))
							oraclePrice := multiplyBasePrecision(big.NewInt(110))
							db := NewInMemoryDatabase()
							var market Market = 1
							db.LastPrice[market] = markPrice
							collateral := HUSD

							//short trader 1
							shortTraderAddress1 := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
							marginShort1 := multiplyBasePrecision(big.NewInt(500))
							shortSize1 := multiplyPrecisionSize(big.NewInt(-20))
							shortEntryPrice1 := multiplyBasePrecision(big.NewInt(80))
							openNotionalShort1 := dividePrecisionSize(big.NewInt(0).Abs(big.NewInt(0).Mul(shortEntryPrice1, shortSize1)))
							shortTrader1 := Trader{
								Margin: Margin{Deposited: map[Collateral]*big.Int{
									collateral: marginShort1,
								}},
								Positions: map[Market]*Position{
									market: getPosition(market, openNotionalShort1, shortSize1, big.NewInt(0), big.NewInt(0), big.NewInt(0)),
								},
							}

							//short trader 2
							shortTraderAddress2 := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
							marginShort2 := multiplyBasePrecision(big.NewInt(500))
							shortSize2 := multiplyPrecisionSize(big.NewInt(-20))
							shortEntryPrice2 := multiplyBasePrecision(big.NewInt(100))
							openNotionalShort2 := dividePrecisionSize(big.NewInt(0).Abs(big.NewInt(0).Mul(shortEntryPrice2, shortSize2)))
							shortTrader2 := Trader{
								Margin: Margin{Deposited: map[Collateral]*big.Int{
									collateral: marginShort2,
								}},
								Positions: map[Market]*Position{
									market: getPosition(market, openNotionalShort2, shortSize2, big.NewInt(0), big.NewInt(0), big.NewInt(0)),
								},
							}
							traderMap := map[common.Address]Trader{
								shortTraderAddress1: shortTrader1,
								shortTraderAddress2: shortTrader2,
							}

							liquidablePositions := GetLiquidableTraders(traderMap, market, markPrice, oraclePrice)

							//oldNotional := 1600000000 -> (ShortEntryPrice1 * ShortSize1)/1e18
							//unrealizedPnl := -600000000 -> ShortSize1(markPrice - ShortEntryPrice1)/1e18
							//effectiveMarginShort1 := -10000000 -> margin + unrealizedPnl
							//newNotional := 2800000000 -> (markPrice * ShortSize1)/1e18
							expectedMarginFractionShort1 := big.NewInt(0)
							assert.Equal(t, 1, len(liquidablePositions))
							assert.Equal(t, shortTraderAddress1, liquidablePositions[0].Address)
							assert.Equal(t, getLiquidationThreshold(shortSize1), liquidablePositions[0].Size)
							assert.Equal(t, expectedMarginFractionShort1, liquidablePositions[0].MarginFraction)
						})
					})
				})
			})
		})
	})
}

func TestGetNormalisedMargin(t *testing.T) {
	t.Run("When trader has no margin", func(t *testing.T) {
		trader := Trader{}
		assert.Equal(t, trader.Margin.Deposited[HUSD], getNormalisedMargin(trader))
	})
	t.Run("When trader has margin in HUSD", func(t *testing.T) {
		margin := multiplyBasePrecision(big.NewInt(10))
		trader := Trader{
			Margin: Margin{Deposited: map[Collateral]*big.Int{
				HUSD: margin,
			}},
		}
		assert.Equal(t, margin, getNormalisedMargin(trader))
	})
}

func TestGetMarginForTrader(t *testing.T) {
	margin := multiplyBasePrecision(big.NewInt(10))
	trader := Trader{
		Margin: Margin{Deposited: map[Collateral]*big.Int{
			HUSD: margin,
		}},
	}
	t.Run("when trader has no positions for a market, it returns output of getNormalized margin", func(t *testing.T) {
		var market Market = 1
		assert.Equal(t, margin, getMarginForTrader(trader, market))
	})
	t.Run("when trader has positions for a market, it subtracts unrealized funding from margin", func(t *testing.T) {
		var market Market = 1
		unrealizedFunding := multiplyBasePrecision(big.NewInt(5))
		position := &Position{UnrealisedFunding: unrealizedFunding}
		trader.Positions = map[Market]*Position{market: position}
		expectedMargin := big.NewInt(0).Sub(margin, unrealizedFunding)
		assert.Equal(t, expectedMargin, getMarginForTrader(trader, market))
	})
}

func TestGetNotionalPosition(t *testing.T) {
	t.Run("When size is positive, it return abs value", func(t *testing.T) {
		price := multiplyBasePrecision(big.NewInt(10))
		size := multiplyPrecisionSize(big.NewInt(20))
		expectedNotionalPosition := dividePrecisionSize(big.NewInt(0).Mul(price, size))
		assert.Equal(t, expectedNotionalPosition, getNotionalPosition(price, size))
	})
	t.Run("When size is negative, it return abs value", func(t *testing.T) {
		price := multiplyBasePrecision(big.NewInt(10))
		size := multiplyPrecisionSize(big.NewInt(-20))
		expectedNotionalPosition := dividePrecisionSize(big.NewInt(0).Abs(big.NewInt(0).Mul(price, size)))
		assert.Equal(t, expectedNotionalPosition, getNotionalPosition(price, size))
	})
}

func TestGetUnrealizedPnl(t *testing.T) {
	t.Run("When newPrice is > entryPrice", func(t *testing.T) {
		t.Run("When size is positive", func(t *testing.T) {
			size := multiplyPrecisionSize(big.NewInt(10))
			entryPrice := multiplyBasePrecision(big.NewInt(10))
			newPrice := multiplyBasePrecision(big.NewInt(15))
			position := &Position{
				Size:         size,
				OpenNotional: getNotionalPosition(entryPrice, size),
			}
			expectedPnl := dividePrecisionSize(big.NewInt(0).Mul(big.NewInt(0).Sub(newPrice, entryPrice), size))
			newNotional := getNotionalPosition(newPrice, size)
			assert.Equal(t, expectedPnl, getUnrealisedPnl(newPrice, position, newNotional))
		})
		t.Run("When size is negative", func(t *testing.T) {
			size := multiplyPrecisionSize(big.NewInt(-10))
			entryPrice := multiplyBasePrecision(big.NewInt(10))
			newPrice := multiplyBasePrecision(big.NewInt(15))
			position := &Position{
				Size:         size,
				OpenNotional: getNotionalPosition(entryPrice, size),
			}
			expectedPnl := dividePrecisionSize(big.NewInt(0).Mul(big.NewInt(0).Sub(newPrice, entryPrice), size))
			newNotional := getNotionalPosition(newPrice, size)
			assert.Equal(t, expectedPnl, getUnrealisedPnl(newPrice, position, newNotional))
		})
	})
	t.Run("When newPrice is < entryPrice", func(t *testing.T) {
		t.Run("When size is positive", func(t *testing.T) {
			size := multiplyPrecisionSize(big.NewInt(10))
			entryPrice := multiplyBasePrecision(big.NewInt(10))
			newPrice := multiplyBasePrecision(big.NewInt(5))
			position := &Position{
				Size:         size,
				OpenNotional: getNotionalPosition(entryPrice, size),
			}
			expectedPnl := dividePrecisionSize(big.NewInt(0).Mul(big.NewInt(0).Sub(newPrice, entryPrice), size))
			newNotional := getNotionalPosition(newPrice, size)
			assert.Equal(t, expectedPnl, getUnrealisedPnl(newPrice, position, newNotional))
		})
		t.Run("When size is negative", func(t *testing.T) {
			size := multiplyPrecisionSize(big.NewInt(-10))
			entryPrice := multiplyBasePrecision(big.NewInt(10))
			newPrice := multiplyBasePrecision(big.NewInt(5))
			position := &Position{
				Size:         size,
				OpenNotional: getNotionalPosition(entryPrice, size),
			}
			expectedPnl := dividePrecisionSize(big.NewInt(0).Mul(big.NewInt(0).Sub(newPrice, entryPrice), size))
			newNotional := getNotionalPosition(newPrice, size)
			assert.Equal(t, expectedPnl, getUnrealisedPnl(newPrice, position, newNotional))
		})
	})
}

func TestGetMarginFraction(t *testing.T) {
	t.Run("If margin + unrealized pnl < 0, it returns 0", func(t *testing.T) {
		margin := multiplyBasePrecision(big.NewInt(5))
		size := multiplyPrecisionSize(big.NewInt(10))
		entryPrice := multiplyBasePrecision(big.NewInt(10))
		newPrice := multiplyBasePrecision(big.NewInt(4))
		position := &Position{
			Size:         size,
			OpenNotional: getNotionalPosition(entryPrice, size),
		}
		assert.Equal(t, big.NewInt(0), getMarginFraction(margin, newPrice, position))
	})
	t.Run("If margin + unrealized pnl > 0, it returns calculated mf", func(t *testing.T) {
		margin := multiplyBasePrecision(big.NewInt(50))
		size := multiplyPrecisionSize(big.NewInt(10))
		entryPrice := multiplyBasePrecision(big.NewInt(10))
		newPrice := multiplyBasePrecision(big.NewInt(6))
		position := &Position{
			Size:         size,
			OpenNotional: getNotionalPosition(entryPrice, size),
		}
		newNotional := getNotionalPosition(newPrice, size)
		expectedMarginFraction := big.NewInt(0).Div(multiplyBasePrecision(big.NewInt(0).Add(margin, getUnrealisedPnl(newPrice, position, newNotional))), getNotionalPosition(newPrice, size))
		assert.Equal(t, expectedMarginFraction, getMarginFraction(margin, newPrice, position))
	})
}

func getPosition(market Market, openNotional *big.Int, size *big.Int, unrealizedFunding *big.Int, lastPremiumFraction *big.Int, liquidationThreshold *big.Int) *Position {
	if liquidationThreshold.Sign() == 0 {
		liquidationThreshold = getLiquidationThreshold(size)
	}
	return &Position{
		OpenNotional:         openNotional,
		Size:                 size,
		UnrealisedFunding:    unrealizedFunding,
		LastPremiumFraction:  lastPremiumFraction,
		LiquidationThreshold: liquidationThreshold,
	}
}
