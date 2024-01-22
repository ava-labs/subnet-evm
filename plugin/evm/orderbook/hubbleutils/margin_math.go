package hubbleutils

import (
	"math"
	"math/big"
)

type UpgradeVersion uint8

const (
	V0 UpgradeVersion = iota
	V1
	V2
)

const V1ActivationTime = uint64(1697129100) // Thursday, 12 October 2023 16:45:00 GMT
type HubbleState struct {
	Assets             []Collateral
	OraclePrices       map[Market]*big.Int
	MidPrices          map[Market]*big.Int
	ActiveMarkets      []Market
	MinAllowableMargin *big.Int
	MaintenanceMargin  *big.Int
	TakerFee           *big.Int
	UpgradeVersion     UpgradeVersion
}

type UserState struct {
	Positions         map[Market]*Position
	ReduceOnlyAmounts []*big.Int
	Margins           []*big.Int
	PendingFunding    *big.Int
	ReservedMargin    *big.Int
}

func UpgradeVersionV0orV1(blockTimestamp uint64) UpgradeVersion {
	if blockTimestamp >= V1ActivationTime {
		return V1
	}
	return V0
}

func GetAvailableMargin(hState *HubbleState, userState *UserState) *big.Int {
	notionalPosition, margin := GetNotionalPositionAndMargin(hState, userState, Min_Allowable_Margin)
	return GetAvailableMargin_(notionalPosition, margin, userState.ReservedMargin, hState.MinAllowableMargin)
}

func GetAvailableMargin_(notionalPosition, margin, reservedMargin, minAllowableMargin *big.Int) *big.Int {
	utilisedMargin := Div1e6(Mul(notionalPosition, minAllowableMargin))
	return Sub(Sub(margin, utilisedMargin), reservedMargin)
}

func GetMarginFraction(hState *HubbleState, userState *UserState) *big.Int {
	notionalPosition, margin := GetNotionalPositionAndMargin(hState, userState, Maintenance_Margin)
	if notionalPosition.Sign() == 0 {
		return big.NewInt(math.MaxInt64)
	}
	return Div(Mul1e6(margin), notionalPosition)
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
		_notionalPosition, _unrealizedPnl := getOptimalPnl(hState, userState.Positions[market], margin, market, marginMode)
		notionalPosition.Add(notionalPosition, _notionalPosition)
		unrealizedPnl.Add(unrealizedPnl, _unrealizedPnl)
	}
	return notionalPosition, unrealizedPnl
}

func getOptimalPnl(hState *HubbleState, position *Position, margin *big.Int, market Market, marginMode MarginMode) (notionalPosition *big.Int, uPnL *big.Int) {
	if position == nil || position.Size.Sign() == 0 {
		return big.NewInt(0), big.NewInt(0)
	}

	// based on oracle price
	oracleBasedNotional, oracleBasedUnrealizedPnl, oracleBasedMF := GetPositionMetadata(
		hState.OraclePrices[market],
		position.OpenNotional,
		position.Size,
		margin,
	)

	// convert to uint8 so that it auto-applies to future version upgrades that may touch unrelated parts of the code
	if uint8(hState.UpgradeVersion) >= uint8(V2) {
		return oracleBasedNotional, oracleBasedUnrealizedPnl
	}

	// based on last price
	notionalPosition, unrealizedPnl, midPriceBasedMF := GetPositionMetadata(
		hState.MidPrices[market],
		position.OpenNotional,
		position.Size,
		margin,
	)

	if hState.UpgradeVersion == V1 {
		if (marginMode == Maintenance_Margin && oracleBasedUnrealizedPnl.Cmp(unrealizedPnl) == 1) || // for liquidations
			(marginMode == Min_Allowable_Margin && oracleBasedUnrealizedPnl.Cmp(unrealizedPnl) == -1) { // for increasing leverage
			return oracleBasedNotional, oracleBasedUnrealizedPnl
		}
		return notionalPosition, unrealizedPnl
	}

	// use V0 logic
	if (marginMode == Maintenance_Margin && oracleBasedMF.Cmp(midPriceBasedMF) == 1) || // for liquidations
		(marginMode == Min_Allowable_Margin && oracleBasedMF.Cmp(midPriceBasedMF) == -1) { // for increasing leverage
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
	if size.Sign() > 0 {
		uPnL = Sub(notionalPosition, openNotional)
	} else {
		uPnL = Sub(openNotional, notionalPosition)
	}
	mf := Div(Mul1e6(Add(margin, uPnL)), notionalPosition)
	return notionalPosition, uPnL, mf
}

func GetNotionalPosition(price *big.Int, size *big.Int) *big.Int {
	return big.NewInt(0).Abs(Div1e18(Mul(price, size)))
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

func GetRequiredMargin(price, fillAmount, minAllowableMargin, takerFee *big.Int) *big.Int {
	quoteAsset := Div1e18(Mul(fillAmount, price))
	return Add(Div1e6(Mul(quoteAsset, minAllowableMargin)), Div1e6(Mul(quoteAsset, takerFee)))
}
