package hubbleutils

import (
	// "encoding/json"
	"math/big"
)

type MarginMode = uint8

const (
	Maintenance_Margin MarginMode = iota
	Min_Allowable_Margin
)

type Collateral struct {
	Price    *big.Int // scaled by 1e6
	Weight   *big.Int // scaled by 1e6
	Decimals uint8
}

type Market = int

type Position struct {
	OpenNotional *big.Int `json:"open_notional"`
	Size         *big.Int `json:"size"`
	// UnrealisedFunding    *big.Int `json:"unrealised_funding"`
	// LastPremiumFraction  *big.Int `json:"last_premium_fraction"`
	// LiquidationThreshold *big.Int `json:"liquidation_threshold"`
}

type Trader struct {
	Positions map[Market]*Position `json:"positions"` // position for every market
	Margin    Margin               `json:"margin"`    // available margin/balance for every market
}

type Margin struct {
	Reserved  *big.Int                `json:"reserved"`
	Deposited map[Collateral]*big.Int `json:"deposited"`
}
