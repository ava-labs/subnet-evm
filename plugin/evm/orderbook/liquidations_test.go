package orderbook

import (
	"math/big"
	"testing"

	hu "github.com/ava-labs/subnet-evm/plugin/evm/orderbook/hubbleutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestGetLiquidableTraders(t *testing.T) {
	var market Market = Market(0)
	collateral := HUSD
	assets := []hu.Collateral{{Price: big.NewInt(1e6), Weight: big.NewInt(1e6), Decimals: 6}}
	t.Run("When no trader exist", func(t *testing.T) {
		db := getDatabase()
		hState := &hu.HubbleState{
			Assets:        assets,
			OraclePrices:  map[Market]*big.Int{market: hu.Mul1e6(big.NewInt(110))},
			ActiveMarkets: []hu.Market{market},
		}
		liquidablePositions, _, _ := db.GetNaughtyTraders(hState)
		assert.Equal(t, 0, len(liquidablePositions))
	})

	t.Run("When no trader has any positions", func(t *testing.T) {
		longTraderAddress := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
		margin := big.NewInt(10000000000)
		db := getDatabase()
		db.TraderMap = map[common.Address]*Trader{
			longTraderAddress: {
				Margin: Margin{
					Reserved:  big.NewInt(0),
					Deposited: map[Collateral]*big.Int{collateral: margin},
				},
				Positions: map[Market]*Position{},
			},
		}
		hState := &hu.HubbleState{
			Assets:            assets,
			OraclePrices:      map[Market]*big.Int{market: hu.Mul1e6(big.NewInt(110))},
			MidPrices:         map[Market]*big.Int{market: hu.Mul1e6(big.NewInt(100))},
			ActiveMarkets:     []hu.Market{market},
			MaintenanceMargin: db.configService.getMaintenanceMargin(),
		}
		liquidablePositions, _, _ := db.GetNaughtyTraders(hState)
		assert.Equal(t, 0, len(liquidablePositions))
	})

	t.Run("long trader", func(t *testing.T) {
		longTraderAddress := common.HexToAddress("0x22Bb736b64A0b4D4081E103f83bccF864F0404aa")
		marginLong := hu.Mul1e6(big.NewInt(500))
		longSize := hu.Mul1e18(big.NewInt(10))
		longEntryPrice := hu.Mul1e6(big.NewInt(90))
		openNotionalLong := hu.Div1e18(big.NewInt(0).Mul(longEntryPrice, longSize))
		pendingFundingLong := hu.Mul1e6(big.NewInt(42))
		t.Run("is saved from liquidation zone by mark price", func(t *testing.T) {
			// setup db
			db := getDatabase()
			longTrader := Trader{
				Margin: Margin{
					Reserved:  big.NewInt(0),
					Deposited: map[Collateral]*big.Int{collateral: marginLong},
				},
				Positions: map[Market]*Position{
					market: getPosition(market, openNotionalLong, longSize, pendingFundingLong, big.NewInt(0), big.NewInt(0), db.configService.getMaxLiquidationRatio(market), db.configService.getMinSizeRequirement(market)),
				},
			}
			db.TraderMap = map[common.Address]*Trader{
				longTraderAddress: &longTrader,
			}

			hState := &hu.HubbleState{
				Assets:             assets,
				OraclePrices:       map[Market]*big.Int{market: hu.Mul1e6(big.NewInt(49))},
				MidPrices:          map[Market]*big.Int{market: hu.Mul1e6(big.NewInt(50))},
				ActiveMarkets:      []hu.Market{market},
				MinAllowableMargin: db.configService.getMinAllowableMargin(),
				MaintenanceMargin:  db.configService.getMaintenanceMargin(),
			}
			oraclePrices := hState.OraclePrices
			// assertions begin
			// for long trader
			_trader := &longTrader
			assert.Equal(t, marginLong, getNormalisedMargin(_trader, assets))
			assert.Equal(t, pendingFundingLong, getTotalFunding(_trader, []Market{market}))

			// open notional = 90 * 10 = 900
			// last price: notional = 50 * 10 = 500, pnl = 500-900 = -400, mf = (500-42-400)/500 = 0.116
			// oracle price: notional = 49 * 10 = 490, pnl = 490-900 = -410, mf = (500-42-410)/490 = 0.097

			// for hu.Min_Allowable_Margin we select the min of 2 hence orale_mf
			notionalPosition, unrealizePnL := getTotalNotionalPositionAndUnrealizedPnl(_trader, new(big.Int).Add(marginLong, pendingFundingLong), hu.Min_Allowable_Margin, oraclePrices, hState.MidPrices, []Market{market})
			assert.Equal(t, hu.Mul1e6(big.NewInt(490)), notionalPosition)
			assert.Equal(t, hu.Mul1e6(big.NewInt(-410)), unrealizePnL)

			hState.MinAllowableMargin = db.configService.getMinAllowableMargin()
			availableMargin := getAvailableMargin(_trader, hState)
			// availableMargin = 500 - 42 (pendingFundingLong) - 410 (uPnL) - 490/5 = -50
			assert.Equal(t, hu.Mul1e6(big.NewInt(-50)), availableMargin)

			// for hu.Maintenance_Margin we select the max of 2 hence, last_mf
			notionalPosition, unrealizePnL = getTotalNotionalPositionAndUnrealizedPnl(_trader, new(big.Int).Add(marginLong, pendingFundingLong), hu.Maintenance_Margin, oraclePrices, hState.MidPrices, []Market{market})
			assert.Equal(t, hu.Mul1e6(big.NewInt(500)), notionalPosition)
			assert.Equal(t, hu.Mul1e6(big.NewInt(-400)), unrealizePnL)

			// marginFraction := calcMarginFraction(_trader, pendingFundingLong, assets, oraclePrices, hState.MidPrices, []Market{market})
			marginFraction := calcMarginFraction(_trader, hState)
			assert.Equal(t, new(big.Int).Div(hu.Mul1e6(new(big.Int).Add(new(big.Int).Sub(marginLong, pendingFundingLong), unrealizePnL)), notionalPosition), marginFraction)

			liquidablePositions, _, _ := db.GetNaughtyTraders(hState)
			assert.Equal(t, 0, len(liquidablePositions))
		})

		t.Run("is saved from liquidation zone by oracle price", func(t *testing.T) {
			// setup db
			db := getDatabase()
			longTrader := Trader{
				Margin: Margin{
					Reserved:  big.NewInt(0),
					Deposited: map[Collateral]*big.Int{collateral: marginLong},
				},
				Positions: map[Market]*Position{
					market: getPosition(market, openNotionalLong, longSize, pendingFundingLong, big.NewInt(0), big.NewInt(0), db.configService.getMaxLiquidationRatio(market), db.configService.getMinSizeRequirement(market)),
				},
			}
			db.TraderMap = map[common.Address]*Trader{
				longTraderAddress: &longTrader,
			}
			hState := &hu.HubbleState{
				Assets:             assets,
				OraclePrices:       map[Market]*big.Int{market: hu.Mul1e6(big.NewInt(50))},
				MidPrices:          map[Market]*big.Int{market: hu.Mul1e6(big.NewInt(49))},
				ActiveMarkets:      []hu.Market{market},
				MinAllowableMargin: db.configService.getMinAllowableMargin(),
				MaintenanceMargin:  db.configService.getMaintenanceMargin(),
			}
			oraclePrices := hState.OraclePrices

			// assertions begin
			// for long trader
			_trader := &longTrader
			assert.Equal(t, marginLong, getNormalisedMargin(_trader, assets))
			assert.Equal(t, pendingFundingLong, getTotalFunding(_trader, []Market{market}))

			// open notional = 90 * 10 = 900
			// last price: notional = 49 * 10 = 490, pnl = 490-900 = -410, mf = (500-42-410)/490 = 0.097
			// oracle price: notional = 50 * 10 = 500, pnl = 500-900 = -400, mf = (500-42-400)/500 = 0.116

			// for hu.Min_Allowable_Margin we select the min of 2 hence last_mf
			notionalPosition, unrealizePnL := getTotalNotionalPositionAndUnrealizedPnl(_trader, new(big.Int).Add(marginLong, pendingFundingLong), hu.Min_Allowable_Margin, oraclePrices, hState.MidPrices, []Market{market})
			assert.Equal(t, hu.Mul1e6(big.NewInt(490)), notionalPosition)
			assert.Equal(t, hu.Mul1e6(big.NewInt(-410)), unrealizePnL)

			availableMargin := getAvailableMargin(_trader, hState)
			// availableMargin = 500 - 42 (pendingFundingLong) - 410 (uPnL) - 490/5 = -50
			assert.Equal(t, hu.Mul1e6(big.NewInt(-50)), availableMargin)

			// for hu.Maintenance_Margin we select the max of 2 hence, oracle_mf
			notionalPosition, unrealizePnL = getTotalNotionalPositionAndUnrealizedPnl(_trader, new(big.Int).Add(marginLong, pendingFundingLong), hu.Maintenance_Margin, oraclePrices, hState.MidPrices, []Market{market})
			assert.Equal(t, hu.Mul1e6(big.NewInt(500)), notionalPosition)
			assert.Equal(t, hu.Mul1e6(big.NewInt(-400)), unrealizePnL)

			marginFraction := calcMarginFraction(_trader, hState)
			assert.Equal(t, new(big.Int).Div(hu.Mul1e6(new(big.Int).Add(new(big.Int).Sub(marginLong, pendingFundingLong), unrealizePnL)), notionalPosition), marginFraction)

			liquidablePositions, _, _ := db.GetNaughtyTraders(hState)
			assert.Equal(t, 0, len(liquidablePositions))
		})
	})

	t.Run("short trader is saved from liquidation zone by mark price", func(t *testing.T) {
		shortTraderAddress := common.HexToAddress("0x710bf5F942331874dcBC7783319123679033b63b")
		marginShort := hu.Mul1e6(big.NewInt(1000))
		shortSize := hu.Mul1e18(big.NewInt(-20))
		shortEntryPrice := hu.Mul1e6(big.NewInt(105))
		openNotionalShort := hu.Div1e18(big.NewInt(0).Abs(big.NewInt(0).Mul(shortEntryPrice, shortSize)))
		pendingFundingShort := hu.Mul1e6(big.NewInt(-37))
		t.Run("is saved from liquidation zone by mark price", func(t *testing.T) {
			// setup db
			db := getDatabase()
			shortTrader := Trader{
				Margin: Margin{
					Reserved:  big.NewInt(0),
					Deposited: map[Collateral]*big.Int{collateral: marginShort},
				},
				Positions: map[Market]*Position{
					market: getPosition(market, openNotionalShort, shortSize, pendingFundingShort, big.NewInt(0), big.NewInt(0), db.configService.getMaxLiquidationRatio(market), db.configService.getMinSizeRequirement(market)),
				},
			}
			db.TraderMap = map[common.Address]*Trader{
				shortTraderAddress: &shortTrader,
			}
			hState := &hu.HubbleState{
				Assets:             assets,
				OraclePrices:       map[Market]*big.Int{market: hu.Mul1e6(big.NewInt(143))},
				MidPrices:          map[Market]*big.Int{market: hu.Mul1e6(big.NewInt(142))},
				ActiveMarkets:      []hu.Market{market},
				MinAllowableMargin: db.configService.getMinAllowableMargin(),
				MaintenanceMargin:  db.configService.getMaintenanceMargin(),
			}
			oraclePrices := hState.OraclePrices

			// assertions begin
			_trader := &shortTrader
			assert.Equal(t, marginShort, getNormalisedMargin(_trader, assets))
			assert.Equal(t, pendingFundingShort, getTotalFunding(_trader, []Market{market}))

			// open notional = 105 * 20 = 2100
			// last price: notional = 142 * 20 = 2840, pnl = 2100-2840 = -740, mf = (1000+37-740)/2840 = 0.104
			// oracle price based notional = 143 * 20 = 2860, pnl = 2100-2860 = -760, mf = (1000+37-760)/2860 = 0.096

			// for hu.Min_Allowable_Margin we select the min of 2 hence, oracle_mf
			notionalPosition, unrealizePnL := getTotalNotionalPositionAndUnrealizedPnl(_trader, new(big.Int).Add(marginShort, pendingFundingShort), hu.Min_Allowable_Margin, oraclePrices, hState.MidPrices, []Market{market})
			assert.Equal(t, hu.Mul1e6(big.NewInt(2860)), notionalPosition)
			assert.Equal(t, hu.Mul1e6(big.NewInt(-760)), unrealizePnL)

			availableMargin := getAvailableMargin(_trader, hState)
			// availableMargin = 1000 + 37 (pendingFundingShort) -760 (uPnL) - 2860/5 = -295
			assert.Equal(t, hu.Mul1e6(big.NewInt(-295)), availableMargin)

			// for hu.Maintenance_Margin we select the max of 2 hence, last_mf
			notionalPosition, unrealizePnL = getTotalNotionalPositionAndUnrealizedPnl(_trader, new(big.Int).Add(marginShort, pendingFundingShort), hu.Maintenance_Margin, oraclePrices, hState.MidPrices, []Market{market})
			assert.Equal(t, hu.Mul1e6(big.NewInt(2840)), notionalPosition)
			assert.Equal(t, hu.Mul1e6(big.NewInt(-740)), unrealizePnL)

			marginFraction := calcMarginFraction(_trader, hState)
			assert.Equal(t, new(big.Int).Div(hu.Mul1e6(new(big.Int).Add(new(big.Int).Sub(marginShort, pendingFundingShort), unrealizePnL)), notionalPosition), marginFraction)

			liquidablePositions, _, _ := db.GetNaughtyTraders(hState)
			assert.Equal(t, 0, len(liquidablePositions))
		})

		t.Run("is saved from liquidation zone by oracle price", func(t *testing.T) {
			// setup db
			db := getDatabase()
			shortTrader := Trader{
				Margin: Margin{
					Reserved:  big.NewInt(0),
					Deposited: map[Collateral]*big.Int{collateral: marginShort},
				},
				Positions: map[Market]*Position{
					market: getPosition(market, openNotionalShort, shortSize, pendingFundingShort, big.NewInt(0), big.NewInt(0), db.configService.getMaxLiquidationRatio(market), db.configService.getMinSizeRequirement(market)),
				},
			}
			db.TraderMap = map[common.Address]*Trader{
				shortTraderAddress: &shortTrader,
			}
			hState := &hu.HubbleState{
				Assets:             assets,
				OraclePrices:       map[Market]*big.Int{market: hu.Mul1e6(big.NewInt(142))},
				MidPrices:          map[Market]*big.Int{market: hu.Mul1e6(big.NewInt(143))},
				ActiveMarkets:      []hu.Market{market},
				MinAllowableMargin: db.configService.getMinAllowableMargin(),
				MaintenanceMargin:  db.configService.getMaintenanceMargin(),
			}
			oraclePrices := hState.OraclePrices

			// assertions begin
			_trader := &shortTrader
			assert.Equal(t, marginShort, getNormalisedMargin(_trader, assets))
			assert.Equal(t, pendingFundingShort, getTotalFunding(_trader, []Market{market}))

			// open notional = 105 * 20 = 2100
			// last price: = 143 * 20 = 2860, pnl = 2100-2860 = -760, mf = (1000+37-760)/2860 = 0.096
			// oracle price: notional = 142 * 20 = 2840, pnl = 2100-2840 = -740, mf = (1000+37-740)/2840 = 0.104

			// for hu.Min_Allowable_Margin we select the min of 2 hence, last_mf
			notionalPosition, unrealizePnL := getTotalNotionalPositionAndUnrealizedPnl(_trader, new(big.Int).Add(marginShort, pendingFundingShort), hu.Min_Allowable_Margin, oraclePrices, hState.MidPrices, []Market{market})
			assert.Equal(t, hu.Mul1e6(big.NewInt(2860)), notionalPosition)
			assert.Equal(t, hu.Mul1e6(big.NewInt(-760)), unrealizePnL)

			availableMargin := getAvailableMargin(_trader, hState)
			// availableMargin = 1000 + 37 (pendingFundingShort) - 760 (uPnL) - 2860/5 = -295
			assert.Equal(t, hu.Mul1e6(big.NewInt(-295)), availableMargin)

			// for hu.Maintenance_Margin we select the max of 2 hence, oracle_mf
			notionalPosition, unrealizePnL = getTotalNotionalPositionAndUnrealizedPnl(_trader, new(big.Int).Add(marginShort, pendingFundingShort), hu.Maintenance_Margin, oraclePrices, hState.MidPrices, []Market{market})
			assert.Equal(t, hu.Mul1e6(big.NewInt(2840)), notionalPosition)
			assert.Equal(t, hu.Mul1e6(big.NewInt(-740)), unrealizePnL)

			marginFraction := calcMarginFraction(_trader, hState)
			assert.Equal(t, new(big.Int).Div(hu.Mul1e6(new(big.Int).Add(new(big.Int).Sub(marginShort, pendingFundingShort), unrealizePnL)), notionalPosition), marginFraction)

			liquidablePositions, _, _ := db.GetNaughtyTraders(hState)
			assert.Equal(t, 0, len(liquidablePositions))
		})
	})
}

func TestGetNormalisedMargin(t *testing.T) {
	assets := []hu.Collateral{{Price: big.NewInt(1e6), Weight: big.NewInt(1e6), Decimals: 6}}
	t.Run("When trader has no margin", func(t *testing.T) {
		trader := Trader{}
		assert.Equal(t, big.NewInt(0), getNormalisedMargin(&trader, assets))
	})
	t.Run("When trader has margin in HUSD", func(t *testing.T) {
		margin := hu.Mul1e6(big.NewInt(10))
		trader := Trader{
			Margin: Margin{Deposited: map[Collateral]*big.Int{
				HUSD: margin,
			}},
		}
		assert.Equal(t, margin, getNormalisedMargin(&trader, assets))
	})
}

func TestGetNotionalPosition(t *testing.T) {
	t.Run("When size is positive, it return abs value", func(t *testing.T) {
		price := hu.Mul1e6(big.NewInt(10))
		size := hu.Mul1e18(big.NewInt(20))
		expectedNotionalPosition := hu.Div1e18(big.NewInt(0).Mul(price, size))
		assert.Equal(t, expectedNotionalPosition, hu.GetNotionalPosition(price, size))
	})
	t.Run("When size is negative, it return abs value", func(t *testing.T) {
		price := hu.Mul1e6(big.NewInt(10))
		size := hu.Mul1e18(big.NewInt(-20))
		expectedNotionalPosition := hu.Div1e18(big.NewInt(0).Abs(big.NewInt(0).Mul(price, size)))
		assert.Equal(t, expectedNotionalPosition, hu.GetNotionalPosition(price, size))
	})
}

func TestGetPositionMetadata(t *testing.T) {
	t.Run("When newPrice is > entryPrice", func(t *testing.T) {
		t.Run("When size is positive", func(t *testing.T) {
			size := hu.Mul1e18(big.NewInt(10))
			entryPrice := hu.Mul1e6(big.NewInt(10))
			newPrice := hu.Mul1e6(big.NewInt(15))
			position := &Position{
				Position: hu.Position{OpenNotional: hu.GetNotionalPosition(entryPrice, size), Size: size},
			}

			arbitaryMarginValue := hu.Mul1e6(big.NewInt(69))
			notionalPosition, uPnL, mf := getPositionMetadata(newPrice, position.OpenNotional, position.Size, arbitaryMarginValue)
			assert.Equal(t, hu.GetNotionalPosition(newPrice, size), notionalPosition)
			expectedPnl := hu.Div1e18(big.NewInt(0).Mul(big.NewInt(0).Sub(newPrice, entryPrice), size))
			assert.Equal(t, expectedPnl, uPnL)
			assert.Equal(t, new(big.Int).Div(hu.Mul1e6(new(big.Int).Add(arbitaryMarginValue, uPnL)), notionalPosition), mf)
		})
		t.Run("When size is negative", func(t *testing.T) {
			size := hu.Mul1e18(big.NewInt(-10))
			entryPrice := hu.Mul1e6(big.NewInt(10))
			newPrice := hu.Mul1e6(big.NewInt(15))
			position := &Position{
				Position: hu.Position{OpenNotional: hu.GetNotionalPosition(entryPrice, size), Size: size},
			}

			notionalPosition, uPnL, _ := getPositionMetadata(newPrice, position.OpenNotional, position.Size, big.NewInt(0))
			assert.Equal(t, hu.GetNotionalPosition(newPrice, size), notionalPosition)
			expectedPnl := hu.Div1e18(big.NewInt(0).Mul(big.NewInt(0).Sub(newPrice, entryPrice), size))
			assert.Equal(t, expectedPnl, uPnL)
		})
	})
	t.Run("When newPrice is < entryPrice", func(t *testing.T) {
		t.Run("When size is positive", func(t *testing.T) {
			size := hu.Mul1e18(big.NewInt(10))
			entryPrice := hu.Mul1e6(big.NewInt(10))
			newPrice := hu.Mul1e6(big.NewInt(5))
			position := &Position{
				Position: hu.Position{OpenNotional: hu.GetNotionalPosition(entryPrice, size), Size: size},
			}

			notionalPosition, uPnL, _ := getPositionMetadata(newPrice, position.OpenNotional, position.Size, big.NewInt(0))
			assert.Equal(t, hu.GetNotionalPosition(newPrice, size), notionalPosition)
			expectedPnl := hu.Div1e18(big.NewInt(0).Mul(big.NewInt(0).Sub(newPrice, entryPrice), size))
			assert.Equal(t, expectedPnl, uPnL)
		})
		t.Run("When size is negative", func(t *testing.T) {
			size := hu.Mul1e18(big.NewInt(-10))
			entryPrice := hu.Mul1e6(big.NewInt(10))
			newPrice := hu.Mul1e6(big.NewInt(5))
			position := &Position{
				Position: hu.Position{OpenNotional: hu.GetNotionalPosition(entryPrice, size), Size: size},
			}
			notionalPosition, uPnL, _ := getPositionMetadata(newPrice, position.OpenNotional, position.Size, big.NewInt(0))
			assert.Equal(t, hu.GetNotionalPosition(newPrice, size), notionalPosition)
			expectedPnl := hu.Div1e18(big.NewInt(0).Mul(big.NewInt(0).Sub(newPrice, entryPrice), size))
			assert.Equal(t, expectedPnl, uPnL)
		})
	})
}

func getPosition(market Market, openNotional *big.Int, size *big.Int, unrealizedFunding *big.Int, lastPremiumFraction *big.Int, liquidationThreshold *big.Int, maxLiquidationRatio *big.Int, minSizeRequirement *big.Int) *Position {
	if liquidationThreshold.Sign() == 0 {
		liquidationThreshold = getLiquidationThreshold(maxLiquidationRatio, minSizeRequirement, size)
	}
	return &Position{
		Position:             hu.Position{OpenNotional: openNotional, Size: size},
		UnrealisedFunding:    unrealizedFunding,
		LastPremiumFraction:  lastPremiumFraction,
		LiquidationThreshold: liquidationThreshold,
	}
}

func getDatabase() *InMemoryDatabase {
	configService := NewMockConfigService()
	configService.Mock.On("getMaintenanceMargin").Return(big.NewInt(1e5))
	configService.Mock.On("getMinAllowableMargin").Return(big.NewInt(2e5))
	configService.Mock.On("getMaxLiquidationRatio").Return(big.NewInt(1e6))
	configService.Mock.On("getMinSizeRequirement").Return(big.NewInt(1e16))

	return NewInMemoryDatabase(configService)
}
