package hubbleutils

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeightedAndSpotCollateral(t *testing.T) {
	assets := []Collateral{
		{
			Price:    big.NewInt(80500000), // 80.5
			Weight:   big.NewInt(800000),   // 0.8
			Decimals: 6,
		},
		{
			Price:    big.NewInt(410000), // 0.41
			Weight:   big.NewInt(900000), // 0.9
			Decimals: 6,
		},
	}
	margins := []*big.Int{
		big.NewInt(3500000),    // 3.5
		big.NewInt(1040000000), // 1040
	}
	expectedWeighted := big.NewInt(609160000) // 609.16
	expectedSpot := big.NewInt(708150000)     // 708.15
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
	hState := &HubbleState{
		Assets: []Collateral{
			{
				Price:    big.NewInt(101000000), // 101
				Weight:   big.NewInt(900000),    // 0.9
				Decimals: 6,
			},
			{
				Price:    big.NewInt(54360000), // 54.36
				Weight:   big.NewInt(700000),   // 0.7
				Decimals: 6,
			},
		},
		MidPrices: map[Market]*big.Int{
			0: big.NewInt(1545340000), // 1545.34
		},
		OraclePrices: map[Market]*big.Int{
			0: big.NewInt(1545210000), // 1545.21
		},
		ActiveMarkets: []Market{
			0,
		},
		MinAllowableMargin: big.NewInt(100000), // 0.1
		MaintenanceMargin:  big.NewInt(200000), // 0.2
	}
	position := &Position{
		Size:         Scale(big.NewInt(582), 14), // 0.0582
		OpenNotional: big.NewInt(87500000),       // 87.5
	}
	margin := big.NewInt(20000000) // 20
	market := 0
	marginMode := Maintenance_Margin

	notionalPosition, uPnL := GetOptimalPnl(hState, position, margin, market, marginMode)

	expectedNotionalPosition := big.NewInt(89938788)
	expectedUPnL := big.NewInt(2438788)

	assert.Equal(t, expectedNotionalPosition, notionalPosition)
	assert.Equal(t, expectedUPnL, uPnL)

	// ------ when marginMode is Min_Allowable_Margin ------

	marginMode = Min_Allowable_Margin

	notionalPosition, uPnL = GetOptimalPnl(hState, position, margin, market, marginMode)

	expectedNotionalPosition = big.NewInt(89931222)
	expectedUPnL = big.NewInt(2431222)

	assert.Equal(t, expectedNotionalPosition, notionalPosition)
	assert.Equal(t, expectedUPnL, uPnL)
}
