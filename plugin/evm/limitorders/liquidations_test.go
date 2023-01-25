package limitorders

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestGetLiquidableTraders(t *testing.T) {
	t.Run("When no trader exist", func(t *testing.T) {
		db := NewInMemoryDatabase()
		var market Market = 1
		markPrice := multiplyBasePrecision(big.NewInt(100))
		db.lastPrice[market] = markPrice
		oraclePrice := multiplyBasePrecision(big.NewInt(110))
		liquidablePositions := db.GetLiquidableTraders(market, oraclePrice)
		assert.Equal(t, 0, len(liquidablePositions))
	})

	t.Run("When traders exist", func(t *testing.T) {

		t.Run("When no trader has any positions", func(t *testing.T) {
			db := NewInMemoryDatabase()
			var market Market = 1
			longTraderAddress := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
			collateral := HUSD
			margin := multiplyBasePrecision(big.NewInt(100))
			db.UpdateMargin(longTraderAddress, collateral, margin)
			markPrice := multiplyBasePrecision(big.NewInt(100))
			db.lastPrice[market] = markPrice
			oraclePrice := multiplyBasePrecision(big.NewInt(110))
			liquidablePositions := db.GetLiquidableTraders(market, oraclePrice)
			assert.Equal(t, 0, len(liquidablePositions))
		})
		t.Run("When traders have positions", func(t *testing.T) {
			t.Run("When mark price is within 20% of oracle price, it uses mark price for calculating margin fraction", func(t *testing.T) {
				t.SkipNow()
				markPrice := multiplyBasePrecision(big.NewInt(100))
				oraclePrice := multiplyBasePrecision(big.NewInt(110))
				t.Run("When traders margin fraction is >= than maintenance margin, GetLiquidableTraders returns empty array", func(t *testing.T) {
					db := NewInMemoryDatabase()
					var market Market = 1
					db.lastPrice[market] = markPrice
					collateral := HUSD

					//long position for trader 1
					longTraderAddress := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
					marginLong := multiplyBasePrecision(big.NewInt(500))
					db.UpdateMargin(longTraderAddress, collateral, marginLong)

					longSize := multiplyPrecisionSize(big.NewInt(10))
					longEntryPrice := multiplyBasePrecision(big.NewInt(90))
					openNotionalLong := dividePrecisionSize(big.NewInt(0).Mul(longEntryPrice, longSize))
					addPosition(db, longTraderAddress, longSize, openNotionalLong, market)

					//short Position for trader 2
					shortTraderAddress := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
					marginShort := multiplyBasePrecision(big.NewInt(1000))
					db.UpdateMargin(shortTraderAddress, collateral, marginShort)
					// open price for short is 2100/20= 105 so trader 2 is in loss
					shortSize := multiplyPrecisionSize(big.NewInt(-20))
					shortEntryPrice := multiplyBasePrecision(big.NewInt(105))
					openNotionalShort := dividePrecisionSize(big.NewInt(0).Abs(big.NewInt(0).Mul(shortEntryPrice, shortSize)))
					addPosition(db, shortTraderAddress, shortSize, openNotionalShort, market)

					liquidablePositions := db.GetLiquidableTraders(market, oraclePrice)
					assert.Equal(t, 0, len(liquidablePositions))
				})
				t.Run("When trader margin fraction is < than maintenance margin, it returns trader's info in GetLiquidableTraders sorted by marginFraction", func(t *testing.T) {
					db := NewInMemoryDatabase()
					var market Market = 1
					db.lastPrice[market] = markPrice
					collateral := HUSD

					//long trader 1
					longTraderAddress := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
					marginLong := multiplyBasePrecision(big.NewInt(500))
					db.UpdateMargin(longTraderAddress, collateral, marginLong)
					// Add long position
					longSize := multiplyPrecisionSize(big.NewInt(10))
					longEntryPrice := multiplyBasePrecision(big.NewInt(145))
					openNotionalLong := dividePrecisionSize(big.NewInt(0).Mul(longEntryPrice, longSize))
					addPosition(db, longTraderAddress, longSize, openNotionalLong, market)
					positionLong := db.traderMap[longTraderAddress].Positions[market]
					expectedMarginFractionLong := getMarginFraction(marginLong, markPrice, positionLong)

					//short trader 1
					shortTraderAddress := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
					marginShort := multiplyBasePrecision(big.NewInt(500))
					db.UpdateMargin(shortTraderAddress, collateral, marginShort)
					shortSize := multiplyPrecisionSize(big.NewInt(-20))
					shortEntryPrice := multiplyBasePrecision(big.NewInt(80))
					openNotionalShort := dividePrecisionSize(big.NewInt(0).Abs(big.NewInt(0).Mul(shortEntryPrice, shortSize)))
					addPosition(db, shortTraderAddress, shortSize, openNotionalShort, market)
					positionShort := db.traderMap[shortTraderAddress].Positions[market]
					expectedMarginFractionShort := getMarginFraction(marginShort, markPrice, positionShort)

					liquidablePositions := db.GetLiquidableTraders(market, oraclePrice)
					assert.Equal(t, 0, len(liquidablePositions))
					assert.Equal(t, 2, len(liquidablePositions))
					assert.Equal(t, longTraderAddress, liquidablePositions[0].Address)
					assert.Equal(t, longSize, liquidablePositions[0].Size)
					assert.Equal(t, expectedMarginFractionLong, liquidablePositions[0].MarginFraction)
					assert.Equal(t, shortTraderAddress, liquidablePositions[1].Address)
					assert.Equal(t, shortSize, liquidablePositions[0].Size)
					assert.Equal(t, expectedMarginFractionShort, liquidablePositions[0].MarginFraction)
				})
			})
			t.Run("When mark price is outside of 20% of oracle price, it uses oracle price for calculating margin fraction", func(t *testing.T) {
				t.Run("When trader margin fraction is >= than maintenance margin", func(t *testing.T) {
					t.SkipNow()
					markPrice := multiplyBasePrecision(big.NewInt(75))
					oraclePrice := multiplyBasePrecision(big.NewInt(100))
					t.Run("When traders margin fraction is >= than maintenance margin, GetLiquidableTraders returns empty array", func(t *testing.T) {
						db := NewInMemoryDatabase()
						var market Market = 1
						db.lastPrice[market] = markPrice
						collateral := HUSD

						//long position for trader 1
						longTraderAddress := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
						marginLong := multiplyBasePrecision(big.NewInt(500))
						db.UpdateMargin(longTraderAddress, collateral, marginLong)

						longSize := multiplyPrecisionSize(big.NewInt(10))
						longEntryPrice := multiplyBasePrecision(big.NewInt(90))
						openNotionalLong := dividePrecisionSize(big.NewInt(0).Mul(longEntryPrice, longSize))
						addPosition(db, longTraderAddress, longSize, openNotionalLong, market)

						//short Position for trader 2
						shortTraderAddress := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
						marginShort := multiplyBasePrecision(big.NewInt(1000))
						db.UpdateMargin(shortTraderAddress, collateral, marginShort)
						// open price for short is 2100/20= 105 so trader 2 is in loss
						shortSize := multiplyPrecisionSize(big.NewInt(-20))
						shortEntryPrice := multiplyBasePrecision(big.NewInt(105))
						openNotionalShort := dividePrecisionSize(big.NewInt(0).Abs(big.NewInt(0).Mul(shortEntryPrice, shortSize)))
						addPosition(db, shortTraderAddress, shortSize, openNotionalShort, market)

						liquidablePositions := db.GetLiquidableTraders(market, oraclePrice)
						assert.Equal(t, 0, len(liquidablePositions))
					})
				})
				t.Run("When trader margin fraction is < than maintenance margin, it returns trader's info in GetLiquidableTraders", func(t *testing.T) {
					t.Run("When mf-markPrice > mf-oraclePrice, it uses mf with mark price", func(t *testing.T) {
						t.Run("For long order", func(t *testing.T) {
							// for both long mf-markPrice will > mf-oraclePrice
							markPrice := multiplyBasePrecision(big.NewInt(140))
							oraclePrice := multiplyBasePrecision(big.NewInt(110))
							db := NewInMemoryDatabase()
							var market Market = 1
							db.lastPrice[market] = markPrice
							collateral := HUSD

							//long trader 1
							longTraderAddress1 := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
							marginLong1 := multiplyBasePrecision(big.NewInt(500))
							db.UpdateMargin(longTraderAddress1, collateral, marginLong1)
							// Add long position 1
							longSize1 := multiplyPrecisionSize(big.NewInt(10))
							longEntryPrice1 := multiplyBasePrecision(big.NewInt(180))
							openNotionalLong1 := dividePrecisionSize(big.NewInt(0).Mul(longEntryPrice1, longSize1))
							addPosition(db, longTraderAddress1, longSize1, openNotionalLong1, market)
							position := db.traderMap[longTraderAddress1].Positions[market]
							expectedMarginFractionLong1 := getMarginFraction(marginLong1, markPrice, position)

							//long trader 2
							longTraderAddress2 := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
							marginLong2 := multiplyBasePrecision(big.NewInt(500))
							db.UpdateMargin(longTraderAddress2, collateral, marginLong2)
							// Add long position 2
							longSize2 := multiplyPrecisionSize(big.NewInt(10))
							longEntryPrice2 := multiplyBasePrecision(big.NewInt(145))
							openNotionalLong2 := dividePrecisionSize(big.NewInt(0).Mul(longEntryPrice2, longSize2))
							addPosition(db, longTraderAddress2, longSize2, openNotionalLong2, market)

							liquidablePositions := db.GetLiquidableTraders(market, oraclePrice)
							assert.Equal(t, 1, len(liquidablePositions))

							//long trader 1 mf-markPrice > maintenanceMargin so it is not liquidated
							//long trader 2 mf-markPrice < maintenanceMargin so it is liquidated
							assert.Equal(t, longTraderAddress1, liquidablePositions[0].Address)
							assert.Equal(t, getLiquidationThreshold(longSize1), liquidablePositions[0].Size)
							assert.Equal(t, expectedMarginFractionLong1, liquidablePositions[0].MarginFraction)
						})
						t.Run("For short order", func(t *testing.T) {
							markPrice := multiplyBasePrecision(big.NewInt(110))
							oraclePrice := multiplyBasePrecision(big.NewInt(140))
							db := NewInMemoryDatabase()
							var market Market = 1
							db.lastPrice[market] = markPrice
							collateral := HUSD

							//short trader 1
							shortTraderAddress1 := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
							marginShort1 := multiplyBasePrecision(big.NewInt(500))
							db.UpdateMargin(shortTraderAddress1, collateral, marginShort1)
							// Add short position 1
							shortSize1 := multiplyPrecisionSize(big.NewInt(-20))
							shortEntryPrice1 := multiplyBasePrecision(big.NewInt(80))
							openNotionalShort1 := dividePrecisionSize(big.NewInt(0).Abs(big.NewInt(0).Mul(shortEntryPrice1, shortSize1)))
							addPosition(db, shortTraderAddress1, shortSize1, openNotionalShort1, market)
							position := db.traderMap[shortTraderAddress1].Positions[market]
							expectedMarginFractionShort1 := getMarginFraction(marginShort1, markPrice, position)

							//short trader 2
							shortTraderAddress2 := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
							marginShort2 := multiplyBasePrecision(big.NewInt(500))
							db.UpdateMargin(shortTraderAddress2, collateral, marginShort2)
							// Add short position 2
							shortSize2 := multiplyPrecisionSize(big.NewInt(-20))
							shortEntryPrice2 := multiplyBasePrecision(big.NewInt(100))
							openNotionalShort2 := dividePrecisionSize(big.NewInt(0).Abs(big.NewInt(0).Mul(shortEntryPrice2, shortSize2)))
							addPosition(db, shortTraderAddress2, shortSize2, openNotionalShort2, market)

							liquidablePositions := db.GetLiquidableTraders(market, oraclePrice)
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
							db.lastPrice[market] = markPrice
							collateral := HUSD

							//long trader 1
							longTraderAddress1 := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
							marginLong1 := multiplyBasePrecision(big.NewInt(500))
							db.UpdateMargin(longTraderAddress1, collateral, marginLong1)
							// Add long position 1
							longSize1 := multiplyPrecisionSize(big.NewInt(10))
							longEntryPrice1 := multiplyBasePrecision(big.NewInt(180))
							openNotionalLong1 := dividePrecisionSize(big.NewInt(0).Mul(longEntryPrice1, longSize1))
							addPosition(db, longTraderAddress1, longSize1, openNotionalLong1, market)
							position := db.traderMap[longTraderAddress1].Positions[market]
							expectedMarginFractionLong1 := getMarginFraction(marginLong1, oraclePrice, position)

							//long trader 2
							longTraderAddress2 := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
							marginLong2 := multiplyBasePrecision(big.NewInt(500))
							db.UpdateMargin(longTraderAddress2, collateral, marginLong2)
							// Add long position 2
							longSize2 := multiplyPrecisionSize(big.NewInt(10))
							longEntryPrice2 := multiplyBasePrecision(big.NewInt(145))
							openNotionalLong2 := dividePrecisionSize(big.NewInt(0).Mul(longEntryPrice2, longSize2))
							addPosition(db, longTraderAddress2, longSize2, openNotionalLong2, market)

							liquidablePositions := db.GetLiquidableTraders(market, oraclePrice)
							assert.Equal(t, 1, len(liquidablePositions))

							//long trader 1 mf-markPrice > maintenanceMargin so it is not liquidated
							//long trader 2 mf-markPrice < maintenanceMargin so it is liquidated
							assert.Equal(t, longTraderAddress1, liquidablePositions[0].Address)
							assert.Equal(t, getLiquidationThreshold(longSize1), liquidablePositions[0].Size)
							assert.Equal(t, expectedMarginFractionLong1, liquidablePositions[0].MarginFraction)
						})
						t.Run("For short order", func(t *testing.T) {
							markPrice := multiplyBasePrecision(big.NewInt(140))
							oraclePrice := multiplyBasePrecision(big.NewInt(110))
							db := NewInMemoryDatabase()
							var market Market = 1
							db.lastPrice[market] = markPrice
							collateral := HUSD

							//short trader 1
							shortTraderAddress1 := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
							marginShort1 := multiplyBasePrecision(big.NewInt(500))
							db.UpdateMargin(shortTraderAddress1, collateral, marginShort1)
							// Add short position 1
							shortSize1 := multiplyPrecisionSize(big.NewInt(-20))
							shortEntryPrice1 := multiplyBasePrecision(big.NewInt(80))
							openNotionalShort1 := dividePrecisionSize(big.NewInt(0).Abs(big.NewInt(0).Mul(shortEntryPrice1, shortSize1)))
							addPosition(db, shortTraderAddress1, shortSize1, openNotionalShort1, market)
							position := db.traderMap[shortTraderAddress1].Positions[market]
							expectedMarginFractionShort1 := getMarginFraction(marginShort1, oraclePrice, position)

							//short trader 2
							shortTraderAddress2 := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
							marginShort2 := multiplyBasePrecision(big.NewInt(500))
							db.UpdateMargin(shortTraderAddress2, collateral, marginShort2)
							// Add short position 2
							shortSize2 := multiplyPrecisionSize(big.NewInt(-20))
							shortEntryPrice2 := multiplyBasePrecision(big.NewInt(100))
							openNotionalShort2 := dividePrecisionSize(big.NewInt(0).Abs(big.NewInt(0).Mul(shortEntryPrice2, shortSize2)))
							addPosition(db, shortTraderAddress2, shortSize2, openNotionalShort2, market)

							liquidablePositions := db.GetLiquidableTraders(market, oraclePrice)
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

func addPosition(db *InMemoryDatabase, address common.Address, size *big.Int, openNotional *big.Int, market Market) {
	db.UpdatePosition(address, market, size, openNotional, false)
}

func TestGetNormalisedMargin(t *testing.T) {
	t.Run("When trader has no margin", func(t *testing.T) {
		trader := &Trader{}
		assert.Equal(t, trader.Margins[HUSD], getNormalisedMargin(trader))
	})
	t.Run("When trader has margin in HUSD", func(t *testing.T) {
		margin := multiplyBasePrecision(big.NewInt(10))
		trader := &Trader{
			Margins: map[Collateral]*big.Int{
				HUSD: margin,
			},
		}
		assert.Equal(t, margin, getNormalisedMargin(trader))
	})
}

func TestGetMarginForTrader(t *testing.T) {
	margin := multiplyBasePrecision(big.NewInt(10))
	trader := &Trader{
		Margins: map[Collateral]*big.Int{
			HUSD: margin,
		},
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
