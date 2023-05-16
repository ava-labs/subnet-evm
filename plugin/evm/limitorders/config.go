package limitorders

import "math/big"

var (
	minAllowableMargin   = big.NewInt(2 * 1e5) // 5x
	maintenanceMargin    = big.NewInt(1e5)
	spreadRatioThreshold = big.NewInt(1e6)
	maxLiquidationRatio  = big.NewInt(25 * 1e4) // 25%
	minSizeRequirement   = big.NewInt(1e16)
)
