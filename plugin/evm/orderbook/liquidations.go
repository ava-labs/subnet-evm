package orderbook

import (
	"math"
	"math/big"
	"sort"

	hu "github.com/ava-labs/subnet-evm/plugin/evm/orderbook/hubbleutils"
	"github.com/ethereum/go-ethereum/common"
)

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

func calcMarginFraction(trader *Trader, pendingFunding *big.Int, assets []hu.Collateral, oraclePrices map[Market]*big.Int, lastPrices map[Market]*big.Int, markets []Market) *big.Int {
	margin := new(big.Int).Sub(getNormalisedMargin(trader, assets), pendingFunding)
	notionalPosition, unrealizePnL := getTotalNotionalPositionAndUnrealizedPnl(trader, margin, hu.Maintenance_Margin, oraclePrices, lastPrices, markets)
	if notionalPosition.Sign() == 0 {
		return big.NewInt(math.MaxInt64)
	}
	margin.Add(margin, unrealizePnL)
	return new(big.Int).Div(hu.Mul1e6(margin), notionalPosition)
}

func sortLiquidableSliceByMarginFraction(positions []LiquidablePosition) []LiquidablePosition {
	sort.SliceStable(positions, func(i, j int) bool {
		return positions[i].MarginFraction.Cmp(positions[j].MarginFraction) == -1
	})
	return positions
}

func getNormalisedMargin(trader *Trader, assets []hu.Collateral) *big.Int {
	return hu.GetNormalizedMargin(assets, getMargins(trader, len(assets)))
}

func getMargins(trader *Trader, numAssets int) []*big.Int {
	margin := make([]*big.Int, numAssets)
	for i := 0; i < numAssets; i++ {
		margin[i] = trader.Margin.Deposited[Collateral(i)]
	}
	return margin
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

type MarginMode = hu.MarginMode

func getTotalNotionalPositionAndUnrealizedPnl(trader *Trader, margin *big.Int, marginMode MarginMode, oraclePrices map[Market]*big.Int, lastPrices map[Market]*big.Int, markets []Market) (*big.Int, *big.Int) {
	return hu.GetTotalNotionalPositionAndUnrealizedPnl(
		&hu.HubbleState{
			OraclePrices:  oraclePrices,
			LastPrices:    lastPrices,
			ActiveMarkets: markets,
		},
		&hu.UserState{
			Positions: translatePositions(trader.Positions),
		},
		margin,
		marginMode,
	)
}

func getPositionMetadata(price *big.Int, openNotional *big.Int, size *big.Int, margin *big.Int) (notionalPosition *big.Int, unrealisedPnl *big.Int, marginFraction *big.Int) {
	return hu.GetPositionMetadata(price, openNotional, size, margin)
}

func prettifyScaledBigInt(number *big.Int, precision int8) string {
	return new(big.Float).Quo(new(big.Float).SetInt(number), big.NewFloat(math.Pow10(int(precision)))).String()
}

func translatePositions(positions map[int]*Position) map[int]*hu.Position {
	huPositions := make(map[int]*hu.Position)
	for key, value := range positions {
		huPositions[key] = &value.Position
	}
	return huPositions
}
