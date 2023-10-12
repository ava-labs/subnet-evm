package hubbleutils

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

var hState = &HubbleState{
	Assets: []Collateral{
		{
			Price:    big.NewInt(1.01 * 1e6), // 1.01
			Weight:   big.NewInt(1e6),        // 1
			Decimals: 6,
		},
		{
			Price:    big.NewInt(54.36 * 1e6), // 54.36
			Weight:   big.NewInt(0.7 * 1e6),   // 0.7
			Decimals: 6,
		},
	},
	MidPrices: map[Market]*big.Int{
		0: big.NewInt(1544.21 * 1e6), // 1544.21
		1: big.NewInt(19.5 * 1e6),    // 19.5
	},
	OraclePrices: map[Market]*big.Int{
		0: big.NewInt(1503.21 * 1e6),
		1: big.NewInt(17.5 * 1e6),
	},
	ActiveMarkets: []Market{
		0, 1,
	},
	MinAllowableMargin: big.NewInt(100000), // 0.1
	MaintenanceMargin:  big.NewInt(200000), // 0.2
}

var userState = &UserState{
	Positions: map[Market]*Position{
		0: {
			Size:         big.NewInt(0.582 * 1e18), // 0.0582
			OpenNotional: big.NewInt(875 * 1e6),    // 87.5, openPrice = 1503.43
		},
		1: {
			Size:         Scale(big.NewInt(-101), 18), // -101
			OpenNotional: big.NewInt(1767.5 * 1e6),    // 1767.5, openPrice = 17.5
		},
	},
	Margins: []*big.Int{
		big.NewInt(30.5 * 1e6), // 30.5
		big.NewInt(14 * 1e6),   // 14
	},
	PendingFunding: big.NewInt(0),
	ReservedMargin: big.NewInt(0),
}

func TestWeightedAndSpotCollateral(t *testing.T) {
	assets := hState.Assets
	margins := userState.Margins
	expectedWeighted := Unscale(Mul(Mul(margins[0], assets[0].Price), assets[0].Weight), assets[0].Decimals+6)
	expectedWeighted.Add(expectedWeighted, Unscale(Mul(Mul(margins[1], assets[1].Price), assets[1].Weight), assets[1].Decimals+6))

	expectedSpot := Unscale(Mul(margins[0], assets[0].Price), assets[0].Decimals)
	expectedSpot.Add(expectedSpot, Unscale(Mul(margins[1], assets[1].Price), assets[1].Decimals))

	resultWeighted, resultSpot := WeightedAndSpotCollateral(assets, margins)
	fmt.Println(resultWeighted, resultSpot)
	assert.Equal(t, expectedWeighted, resultWeighted)
	assert.Equal(t, expectedSpot, resultSpot)

	normalisedMargin := GetNormalizedMargin(assets, margins)
	assert.Equal(t, expectedWeighted, normalisedMargin)

}

func TestGetNotionalPosition(t *testing.T) {
	price := Scale(big.NewInt(1200), 6)
	size := Scale(big.NewInt(5), 18)
	expected := Scale(big.NewInt(6000), 6)

	result := GetNotionalPosition(price, size)

	assert.Equal(t, expected, result)
}

func TestGetPositionMetadata(t *testing.T) {
	price := big.NewInt(20250000)        // 20.25
	openNotional := big.NewInt(75369000) // 75.369 (size * 18.5)
	size := Scale(big.NewInt(40740), 14) // 4.074
	margin := big.NewInt(20000000)       // 20

	notionalPosition, unrealisedPnl, marginFraction := GetPositionMetadata(price, openNotional, size, margin)

	expectedNotionalPosition := big.NewInt(82498500) // 82.4985
	expectedUnrealisedPnl := big.NewInt(7129500)     // 7.1295
	expectedMarginFraction := big.NewInt(328848)     // 0.328848

	assert.Equal(t, expectedNotionalPosition, notionalPosition)
	assert.Equal(t, expectedUnrealisedPnl, unrealisedPnl)
	assert.Equal(t, expectedMarginFraction, marginFraction)

	// ------ when size is negative ------
	size = Scale(big.NewInt(-40740), 14) // -4.074
	openNotional = big.NewInt(75369000)  // 75.369 (size * 18.5)
	notionalPosition, unrealisedPnl, marginFraction = GetPositionMetadata(price, openNotional, size, margin)
	fmt.Println("notionalPosition", notionalPosition, "unrealisedPnl", unrealisedPnl, "marginFraction", marginFraction)

	expectedNotionalPosition = big.NewInt(82498500) // 82.4985
	expectedUnrealisedPnl = big.NewInt(-7129500)    // -7.1295
	expectedMarginFraction = big.NewInt(156008)     // 0.156008

	assert.Equal(t, expectedNotionalPosition, notionalPosition)
	assert.Equal(t, expectedUnrealisedPnl, unrealisedPnl)
	assert.Equal(t, expectedMarginFraction, marginFraction)
}

func TestGetOptimalPnl(t *testing.T) {
	margin := big.NewInt(20 * 1e6) // 20
	market := 0
	position := userState.Positions[market]
	marginMode := Maintenance_Margin

	notionalPosition, uPnL := getOptimalPnl(hState, position, margin, market, marginMode, 0)

	// mid price pnl is more than oracle price pnl
	expectedNotionalPosition := Unscale(Mul(position.Size, hState.MidPrices[market]), 18)
	expectedUPnL := Sub(expectedNotionalPosition, position.OpenNotional)
	fmt.Println("Maintenace_Margin_Mode", "notionalPosition", notionalPosition, "uPnL", uPnL)

	assert.Equal(t, expectedNotionalPosition, notionalPosition)
	assert.Equal(t, expectedUPnL, uPnL)

	// ------ when marginMode is Min_Allowable_Margin ------

	marginMode = Min_Allowable_Margin
	notionalPosition, uPnL = getOptimalPnl(hState, position, margin, market, marginMode, 0)

	expectedNotionalPosition = Unscale(Mul(position.Size, hState.OraclePrices[market]), 18)
	expectedUPnL = Sub(expectedNotionalPosition, position.OpenNotional)
	fmt.Println("Min_Allowable_Margin_Mode", "notionalPosition", notionalPosition, "uPnL", uPnL)

	assert.Equal(t, expectedNotionalPosition, notionalPosition)
	assert.Equal(t, expectedUPnL, uPnL)
}

func TestGetOptimalPnlDeprecated(t *testing.T) {
	margin := big.NewInt(20 * 1e6) // 20
	market := 0
	position := userState.Positions[market]
	marginMode := Maintenance_Margin

	notionalPosition, uPnL := getOptimalPnl(hState, position, margin, market, marginMode, 1)

	// mid price pnl is more than oracle price pnl
	expectedNotionalPosition := Unscale(Mul(position.Size, hState.MidPrices[market]), 18)
	expectedUPnL := Sub(expectedNotionalPosition, position.OpenNotional)
	fmt.Println("Maintenace_Margin_Mode", "notionalPosition", notionalPosition, "uPnL", uPnL)

	assert.Equal(t, expectedNotionalPosition, notionalPosition)
	assert.Equal(t, expectedUPnL, uPnL)

	// ------ when marginMode is Min_Allowable_Margin ------

	marginMode = Min_Allowable_Margin
	notionalPosition, uPnL = getOptimalPnl(hState, position, margin, market, marginMode, 1)

	expectedNotionalPosition = Unscale(Mul(position.Size, hState.OraclePrices[market]), 18)
	expectedUPnL = Sub(expectedNotionalPosition, position.OpenNotional)
	fmt.Println("Min_Allowable_Margin_Mode", "notionalPosition", notionalPosition, "uPnL", uPnL)

	assert.Equal(t, expectedNotionalPosition, notionalPosition)
	assert.Equal(t, expectedUPnL, uPnL)
}

func TestGetTotalNotionalPositionAndUnrealizedPnl(t *testing.T) {
	margin := GetNormalizedMargin(hState.Assets, userState.Margins)
	// margin := big.NewInt(2000 * 1e6) // 50
	fmt.Println("margin = ", margin) // 563.533
	marginMode := Maintenance_Margin
	fmt.Println("availableMargin = ", GetAvailableMargin(hState, userState))
	fmt.Println("marginFraction = ", GetMarginFraction(hState, userState))

	notionalPosition, uPnL := GetTotalNotionalPositionAndUnrealizedPnl(hState, userState, margin, marginMode, 0)
	fmt.Println("Maintenace_Margin_Mode ", "notionalPosition = ", notionalPosition, "uPnL = ", uPnL)
	_, pnl := getOptimalPnl(hState, userState.Positions[0], margin, 0, marginMode, 0)
	fmt.Println("best pnl market 0 =", pnl)
	_, pnl = getOptimalPnl(hState, userState.Positions[1], margin, 1, marginMode, 0)
	fmt.Println("best pnl market 1 =", pnl)

	// mid price pnl is more than oracle price pnl for long position
	expectedNotionalPosition := Unscale(Mul(userState.Positions[0].Size, hState.MidPrices[0]), 18)
	expectedUPnL := Sub(expectedNotionalPosition, userState.Positions[0].OpenNotional)
	// oracle price pnl is more than mid price pnl for short position
	expectedNotional2 := Abs(Unscale(Mul(userState.Positions[1].Size, hState.OraclePrices[1]), 18))
	expectedNotionalPosition.Add(expectedNotionalPosition, expectedNotional2)
	expectedUPnL.Add(expectedUPnL, Sub(userState.Positions[1].OpenNotional, expectedNotional2))

	assert.Equal(t, expectedNotionalPosition, notionalPosition)
	assert.Equal(t, expectedUPnL, uPnL)

	// ------ when marginMode is Min_Allowable_Margin ------

	marginMode = Min_Allowable_Margin
	notionalPosition, uPnL = GetTotalNotionalPositionAndUnrealizedPnl(hState, userState, margin, marginMode, 0)
	fmt.Println("Min_Allowable_Margin_Mode ", "notionalPosition = ", notionalPosition, "uPnL = ", uPnL)

	_, pnl = getOptimalPnl(hState, userState.Positions[0], margin, 0, marginMode, 0)
	fmt.Println("worst pnl market 0 =", pnl)
	_, pnl = getOptimalPnl(hState, userState.Positions[1], margin, 1, marginMode, 0)
	fmt.Println("worst pnl market 1 =", pnl)

	expectedNotionalPosition = Unscale(Mul(userState.Positions[0].Size, hState.OraclePrices[0]), 18)
	expectedUPnL = Sub(expectedNotionalPosition, userState.Positions[0].OpenNotional)
	expectedNotional2 = Abs(Unscale(Mul(userState.Positions[1].Size, hState.MidPrices[1]), 18))
	expectedNotionalPosition.Add(expectedNotionalPosition, expectedNotional2)
	expectedUPnL.Add(expectedUPnL, Sub(userState.Positions[1].OpenNotional, expectedNotional2))

	assert.Equal(t, expectedNotionalPosition, notionalPosition)
	assert.Equal(t, expectedUPnL, uPnL)
}
