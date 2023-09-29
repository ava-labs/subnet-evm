package hubbleutils

import (
	"math/big"
)

type HubbleState struct {
	Assets             []Collateral
	OraclePrices       map[Market]*big.Int
	LastPrices         map[Market]*big.Int
	ActiveMarkets      []Market
	MinAllowableMargin *big.Int
}

type UserState struct {
	Positions      map[Market]*Position
	Margins        []*big.Int
	PendingFunding *big.Int
	ReservedMargin *big.Int
}

func GetAvailableMargin(hState *HubbleState, userState *UserState) *big.Int {
	notionalPosition, margin := GetNotionalPositionAndMargin(hState, userState, Min_Allowable_Margin)
	return GetAvailableMargin_(notionalPosition, margin, userState.ReservedMargin, hState.MinAllowableMargin)
}

func GetAvailableMargin_(notionalPosition, margin, reservedMargin, minAllowableMargin *big.Int) *big.Int {
	utilisedMargin := Div1e6(Mul(notionalPosition, minAllowableMargin))
	return Sub(Sub(margin, utilisedMargin), reservedMargin)
}

func GetNotionalPositionAndMargin(hState *HubbleState, userState *UserState, marginMode MarginMode) (*big.Int, *big.Int) {
	margin := Sub(GetNormalizedMargin(hState.Assets, userState.Margins), userState.PendingFunding)
	notionalPosition, unrealizedPnl := GetTotalNotionalPositionAndUnrealizedPnl(hState, userState, margin, marginMode)
	return notionalPosition, Add(margin, unrealizedPnl)
}

func GetTotalNotionalPositionAndUnrealizedPnl(hState *HubbleState, userState *UserState, margin *big.Int, marginMode MarginMode) (*big.Int, *big.Int) {
	notionalPosition := big.NewInt(0)
	unrealizedPnl := big.NewInt(0)
	for _, market := range hState.ActiveMarkets {
		_notionalPosition, _unrealizedPnl := GetOptimalPnl(hState, userState.Positions[market], margin, market, marginMode)
		notionalPosition.Add(notionalPosition, _notionalPosition)
		unrealizedPnl.Add(unrealizedPnl, _unrealizedPnl)
	}
	return notionalPosition, unrealizedPnl
}

func GetOptimalPnl(hState *HubbleState, position *Position, margin *big.Int, market Market, marginMode MarginMode) (notionalPosition *big.Int, uPnL *big.Int) {
	if position == nil || position.Size.Sign() == 0 {
		return big.NewInt(0), big.NewInt(0)
	}

	// based on last price
	notionalPosition, unrealizedPnl, lastPriceBasedMF := GetPositionMetadata(
		hState.LastPrices[market],
		position.OpenNotional,
		position.Size,
		margin,
	)

	// based on oracle price
	oracleBasedNotional, oracleBasedUnrealizedPnl, oracleBasedMF := GetPositionMetadata(
		hState.OraclePrices[market],
		position.OpenNotional,
		position.Size,
		margin,
	)

	if (marginMode == Maintenance_Margin && oracleBasedMF.Cmp(lastPriceBasedMF) == 1) || // for liquidations
		(marginMode == Min_Allowable_Margin && oracleBasedMF.Cmp(lastPriceBasedMF) == -1) { // for increasing leverage
		return oracleBasedNotional, oracleBasedUnrealizedPnl
	}
	return notionalPosition, unrealizedPnl
}

func GetPositionMetadata(price *big.Int, openNotional *big.Int, size *big.Int, margin *big.Int) (notionalPosition *big.Int, unrealisedPnl *big.Int, marginFraction *big.Int) {
	notionalPosition = GetNotionalPosition(price, size)
	uPnL := new(big.Int)
	if notionalPosition.Sign() == 0 {
		return big.NewInt(0), big.NewInt(0), big.NewInt(0)
	}
	if size.Cmp(big.NewInt(0)) > 0 {
		uPnL = Sub(notionalPosition, openNotional)
	} else {
		uPnL = Sub(openNotional, notionalPosition)
	}
	mf := Div(Mul1e6(Add(margin, uPnL)), notionalPosition)
	return notionalPosition, uPnL, mf
}

func GetNotionalPosition(price *big.Int, size *big.Int) *big.Int {
	return big.NewInt(0).Abs(Div1e18(Mul(size, price)))
}

func GetNormalizedMargin(assets []Collateral, margins []*big.Int) *big.Int {
	weighted, _ := WeightedAndSpotCollateral(assets, margins)
	return weighted
}

func WeightedAndSpotCollateral(assets []Collateral, margins []*big.Int) (weighted, spot *big.Int) {
	weighted = big.NewInt(0)
	spot = big.NewInt(0)
	for i, asset := range assets {
		if margins[i] == nil || margins[i].Sign() == 0 {
			continue
		}
		numerator := Mul(margins[i], asset.Price) // margin[i] is scaled by asset.Decimal
		spot.Add(spot, Unscale(numerator, asset.Decimals))
		weighted.Add(weighted, Unscale(Mul(numerator, asset.Weight), asset.Decimals+6))
	}
	return weighted, spot
}
