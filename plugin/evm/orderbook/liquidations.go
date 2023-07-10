package orderbook

import (
	"math"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/log"
)

var BASE_PRECISION = _1e6
var SIZE_BASE_PRECISION = _1e18

type LiquidablePosition struct {
	Address        common.Address
	Market         Market
	Size           *big.Int
	MarginFraction *big.Int
	FilledSize     *big.Int
	PositionType   PositionType
}

func (liq LiquidablePosition) GetUnfilledSize() *big.Int {
	return big.NewInt(0).Sub(liq.Size, liq.FilledSize)
}

// returns the max(oracle_mf, last_mf); hence should only be used to determine the margin fraction for liquidation and not to increase leverage
func calcMarginFraction(trader *Trader, pendingFunding *big.Int, oraclePrices map[Market]*big.Int, lastPrices map[Market]*big.Int, markets []Market) *big.Int {
	margin := new(big.Int).Sub(getNormalisedMargin(trader), pendingFunding)
	notionalPosition, unrealizePnL := getTotalNotionalPositionAndUnrealizedPnl(trader, margin, Maintenance_Margin, oraclePrices, lastPrices, markets)
	// log.Info("calcMarginFraction:M", "margin", margin, "notionalPosition", notionalPosition, "unrealizePnL", unrealizePnL)
	if notionalPosition.Sign() == 0 {
		return big.NewInt(math.MaxInt64)
	}
	margin.Add(margin, unrealizePnL)
	// log.Info("calcMarginFraction", "margin", margin, "notionalPosition", notionalPosition)
	return new(big.Int).Div(multiplyBasePrecision(margin), notionalPosition)
}

func sortLiquidableSliceByMarginFraction(positions []LiquidablePosition) []LiquidablePosition {
	sort.SliceStable(positions, func(i, j int) bool {
		return positions[i].MarginFraction.Cmp(positions[j].MarginFraction) == -1
	})
	return positions
}

func getNormalisedMargin(trader *Trader) *big.Int {
	return trader.Margin.Deposited[HUSD]
	// @todo: Write for multi-collateral
}

func getTotalFunding(trader *Trader, markets []Market) *big.Int {
	totalPendingFunding := big.NewInt(0)
	for _, market := range markets {
		if trader.Positions[market] != nil {
			totalPendingFunding.Add(totalPendingFunding, trader.Positions[market].UnrealisedFunding)
		}
	}
	return totalPendingFunding
}

func getNotionalPosition(price *big.Int, size *big.Int) *big.Int {
	return big.NewInt(0).Abs(dividePrecisionSize(big.NewInt(0).Mul(size, price)))
}

type MarginMode uint8

const (
	Maintenance_Margin MarginMode = iota
	Min_Allowable_Margin
)

func getTotalNotionalPositionAndUnrealizedPnl(trader *Trader, margin *big.Int, marginMode MarginMode, oraclePrices map[Market]*big.Int, lastPrices map[Market]*big.Int, markets []Market) (*big.Int, *big.Int) {
	notionalPosition := big.NewInt(0)
	unrealizedPnl := big.NewInt(0)
	for _, market := range markets {
		_notionalPosition, _unrealizedPnl := getOptimalPnl(market, oraclePrices[market], lastPrices[market], trader, margin, marginMode)
		notionalPosition.Add(notionalPosition, _notionalPosition)
		unrealizedPnl.Add(unrealizedPnl, _unrealizedPnl)
	}
	return notionalPosition, unrealizedPnl
}

func getOptimalPnl(market Market, oraclePrice *big.Int, lastPrice *big.Int, trader *Trader, margin *big.Int, marginMode MarginMode) (notionalPosition *big.Int, uPnL *big.Int) {
	position := trader.Positions[market]
	if position == nil || position.Size.Sign() == 0 {
		return big.NewInt(0), big.NewInt(0)
	}

	// based on last price
	notionalPosition, unrealizedPnl, lastPriceBasedMF := getPositionMetadata(
		lastPrice,
		position.OpenNotional,
		position.Size,
		margin,
	)
	// log.Info("in getOptimalPnl", "notionalPosition", notionalPosition, "unrealizedPnl", unrealizedPnl, "lastPriceBasedMF", lastPriceBasedMF)

	// based on oracle price
	oracleBasedNotional, oracleBasedUnrealizedPnl, oracleBasedMF := getPositionMetadata(
		oraclePrice,
		position.OpenNotional,
		position.Size,
		margin,
	)
	// log.Info("in getOptimalPnl", "oracleBasedNotional", oracleBasedNotional, "oracleBasedUnrealizedPnl", oracleBasedUnrealizedPnl, "oracleBasedMF", oracleBasedMF)

	if (marginMode == Maintenance_Margin && oracleBasedMF.Cmp(lastPriceBasedMF) == 1) || // for liquidations
		(marginMode == Min_Allowable_Margin && oracleBasedMF.Cmp(lastPriceBasedMF) == -1) { // for increasing leverage
		return oracleBasedNotional, oracleBasedUnrealizedPnl
	}
	return notionalPosition, unrealizedPnl
}

func getPositionMetadata(price *big.Int, openNotional *big.Int, size *big.Int, margin *big.Int) (notionalPosition *big.Int, unrealisedPnl *big.Int, marginFraction *big.Int) {
	// log.Info("in getPositionMetadata", "price", price, "openNotional", openNotional, "size", size, "margin", margin)
	notionalPosition = getNotionalPosition(price, size)
	uPnL := new(big.Int)
	if notionalPosition.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(0), big.NewInt(0), big.NewInt(0)
	}
	if size.Cmp(big.NewInt(0)) > 0 {
		uPnL = new(big.Int).Sub(notionalPosition, openNotional)
	} else {
		uPnL = new(big.Int).Sub(openNotional, notionalPosition)
	}
	mf := new(big.Int).Div(multiplyBasePrecision(new(big.Int).Add(margin, uPnL)), notionalPosition)
	return notionalPosition, uPnL, mf
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
