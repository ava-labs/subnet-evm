package limitorders

import (
	"math"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

var maintenanceMargin = big.NewInt(1e5)
var spreadRatioThreshold = big.NewInt(1e6)

var BASE_PRECISION = big.NewInt(1e6)
var SIZE_BASE_PRECISION = big.NewInt(1e18)

type LiquidablePosition struct {
	Address        common.Address
	Size           *big.Int
	MarginFraction *big.Int
	FilledSize     *big.Int
	PositionType   string
}

func (liq LiquidablePosition) GetUnfilledSize() *big.Int {
	return big.NewInt(0).Sub(liq.Size, liq.FilledSize)
}

func GetLiquidableTraders(traderMap map[common.Address]Trader, market Market, lastPrice *big.Int, oraclePrice *big.Int) []LiquidablePosition {
	liquidablePositions := []LiquidablePosition{}
	markPrice := lastPrice

	overSpreadLimit := isOverSpreadLimit(markPrice, oraclePrice)
	log.Info("GetLiquidableTraders:", "markPrice", markPrice, "oraclePrice", oraclePrice, "overSpreadLimit", overSpreadLimit)

	for addr, trader := range traderMap {
		position := trader.Positions[market]
		if position != nil && position.Size.Sign() != 0 {
			margin := getMarginForTrader(trader, market)
			marginFraction := getMarginFraction(margin, markPrice, position)

			log.Info("GetLiquidableTraders", "trader", addr.String(), "traderInfo", trader, "marginFraction", marginFraction, "margin", margin.Uint64())
			if overSpreadLimit {
				oracleBasedMarginFraction := getMarginFraction(margin, oraclePrice, position)
				if oracleBasedMarginFraction.Cmp(marginFraction) == 1 {
					marginFraction = oracleBasedMarginFraction
				}
				log.Info("GetLiquidableTraders", "trader", addr.String(), "oracleBasedMarginFraction", oracleBasedMarginFraction)
			}

			if marginFraction.Cmp(maintenanceMargin) == -1 {
				log.Info("GetLiquidableTraders - below maintenanceMargin", "trader", addr.String(), "marginFraction", marginFraction)
				liquidable := LiquidablePosition{
					Address:        addr,
					Size:           position.LiquidationThreshold,
					MarginFraction: marginFraction,
					FilledSize:     big.NewInt(0),
				}
				if position.Size.Sign() == -1 {
					liquidable.PositionType = "short"
				} else {
					liquidable.PositionType = "long"
				}
				liquidablePositions = append(liquidablePositions, liquidable)
			}
		}
	}

	// lower margin fraction positions should be liquidated first
	sortLiquidableSliceByMarginFraction(liquidablePositions)
	return liquidablePositions
}

func sortLiquidableSliceByMarginFraction(positions []LiquidablePosition) []LiquidablePosition {
	sort.SliceStable(positions, func(i, j int) bool {
		return positions[i].MarginFraction.Cmp(positions[j].MarginFraction) == -1
	})
	return positions
}

func isOverSpreadLimit(markPrice *big.Int, oraclePrice *big.Int) bool {
	// diff := abs(markPrice - oraclePrice)
	diff := multiplyBasePrecision(big.NewInt(0).Abs(big.NewInt(0).Sub(markPrice, oraclePrice)))
	// spreadRatioAbs := diff * 100 / oraclePrice
	spreadRatioAbs := big.NewInt(0).Div(diff, oraclePrice)
	if spreadRatioAbs.Cmp(spreadRatioThreshold) >= 0 {
		return true
	} else {
		return false
	}
}

func getNormalisedMargin(trader Trader) *big.Int {
	return trader.Margins[HUSD]

	// this will change after multi collateral
	// var normalisedMargin *big.Int
	// for coll, margin := range trader.Margins {
	// 	normalisedMargin += margin * priceMap[coll] * collateralWeightMap[coll]
	// }

	// return normalisedMargin
}

func getMarginForTrader(trader Trader, market Market) *big.Int {
	if position, ok := trader.Positions[market]; ok {
		if position.UnrealisedFunding != nil {
			return big.NewInt(0).Sub(getNormalisedMargin(trader), position.UnrealisedFunding)
		}
	}
	return getNormalisedMargin(trader)
}

func getNotionalPosition(price *big.Int, size *big.Int) *big.Int {
	//notional position is base precision 1e6
	return big.NewInt(0).Abs(dividePrecisionSize(big.NewInt(0).Mul(size, price)))
}

func getUnrealisedPnl(price *big.Int, position *Position, notionalPosition *big.Int) *big.Int {
	if position.Size.Sign() == 1 {
		return big.NewInt(0).Sub(notionalPosition, position.OpenNotional)
	} else {
		return big.NewInt(0).Sub(position.OpenNotional, notionalPosition)
	}
}

func getMarginFraction(margin *big.Int, price *big.Int, position *Position) *big.Int {
	notionalPosition := getNotionalPosition(price, position.Size)
	unrealisedPnl := getUnrealisedPnl(price, position, notionalPosition)
	log.Info("getMarginFraction:", "notionalPosition", notionalPosition, "unrealisedPnl", unrealisedPnl)
	effectionMargin := big.NewInt(0).Add(margin, unrealisedPnl)
	mf := big.NewInt(0).Div(multiplyBasePrecision(effectionMargin), notionalPosition)
	if mf.Sign() == -1 {
		return big.NewInt(0) // why?
	}
	return mf
}

func multiplyBasePrecision(number *big.Int) *big.Int {
	return big.NewInt(0).Mul(number, BASE_PRECISION)
}

func multiplyPrecisionSize(number *big.Int) *big.Int {
	return big.NewInt(0).Mul(number, SIZE_BASE_PRECISION)
}

func dividePrecisionSize(number *big.Int) *big.Int {
	return big.NewInt(0).Div(number, SIZE_BASE_PRECISION)
}

func divideByBasePrecision(number *big.Int) *big.Int {
	return big.NewInt(0).Div(number, BASE_PRECISION)
}

func prettifyScaledBigInt(number *big.Int, precision int8) string {
	return new(big.Float).Quo(new(big.Float).SetInt(number), big.NewFloat(math.Pow10(int(precision)))).String()
}
