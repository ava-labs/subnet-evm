// (c) 2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// blockgascost implements the block gas cost logic
package blockgascost

import (
	"math"

	safemath "github.com/ava-labs/avalanchego/utils/math"
	"github.com/ava-labs/subnet-evm/commontype"
)

// BlockGasCost calculates the required block gas cost.
//
// cost = parentCost + step * (TargetBlockRate - timeElapsed)
//
// The returned cost is clamped to [MinBlockGasCost, MaxBlockGasCost].
func BlockGasCost(
	feeConfig commontype.FeeConfig,
	parentCost uint64,
	step uint64,
	timeElapsed uint64,
) uint64 {
	deviation := safemath.AbsDiff(feeConfig.TargetBlockRate, timeElapsed)
	change, err := safemath.Mul(step, deviation)
	if err != nil {
		change = math.MaxUint64
	}

	var (
		minBlockGasCost uint64 = feeConfig.MinBlockGasCost.Uint64()
		maxBlockGasCost uint64 = feeConfig.MaxBlockGasCost.Uint64()
		op                     = safemath.Add[uint64]
		defaultCost     uint64 = feeConfig.MaxBlockGasCost.Uint64()
	)
	if timeElapsed > feeConfig.TargetBlockRate {
		op = safemath.Sub
		defaultCost = minBlockGasCost
	}

	cost, err := op(parentCost, change)
	if err != nil {
		cost = defaultCost
	}

	switch {
	case cost < minBlockGasCost:
		// This is technically dead code because [MinBlockGasCost] is 0, but it
		// makes the code more clear.
		return minBlockGasCost
	case cost > maxBlockGasCost:
		return maxBlockGasCost
	default:
		return cost
	}
}
