// (c) 2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package header

import (
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/stretchr/testify/assert"
)

func TestBlockGasCost(t *testing.T) {
	testFeeConfig := commontype.FeeConfig{
		MinBlockGasCost:  big.NewInt(0),
		MaxBlockGasCost:  big.NewInt(1_000_000),
		TargetBlockRate:  2,
		BlockGasCostStep: big.NewInt(50_000),
	}
	BlockGasCostTest(t, testFeeConfig)

	testFeeConfigDouble := commontype.FeeConfig{
		MinBlockGasCost:  big.NewInt(2),
		MaxBlockGasCost:  big.NewInt(2_000_000),
		TargetBlockRate:  4,
		BlockGasCostStep: big.NewInt(100_000),
	}
	BlockGasCostTest(t, testFeeConfigDouble)
}

func BlockGasCostTest(t *testing.T, testFeeConfig commontype.FeeConfig) {
	maxBlockGasCostBig := testFeeConfig.MaxBlockGasCost
	maxBlockGasCost := testFeeConfig.MaxBlockGasCost.Uint64()
	blockGasCostStep := testFeeConfig.BlockGasCostStep.Uint64()
	minBlockGasCost := testFeeConfig.MinBlockGasCost.Uint64()
	targetBlockRate := testFeeConfig.TargetBlockRate

	tests := []struct {
		name       string
		parentTime uint64
		parentCost *big.Int
		timestamp  uint64
		expected   uint64
	}{
		{
			name:       "normal",
			parentTime: 10,
			parentCost: maxBlockGasCostBig,
			timestamp:  10 + targetBlockRate + 1,
			expected:   maxBlockGasCost - blockGasCostStep,
		},
		{
			name:       "negative_time_elapsed",
			parentTime: 10,
			parentCost: maxBlockGasCostBig,
			timestamp:  9,
			expected:   minBlockGasCost + blockGasCostStep*targetBlockRate,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parent := &types.Header{
				Time:         test.parentTime,
				BlockGasCost: test.parentCost,
			}

			assert.Equal(t, test.expected, BlockGasCost(
				testFeeConfig,
				parent,
				test.timestamp,
			))
		})
	}
}

func TestBlockGasCostWithStep(t *testing.T) {
	testFeeConfig := commontype.FeeConfig{
		MinBlockGasCost:  big.NewInt(0),
		MaxBlockGasCost:  big.NewInt(1_000_000),
		TargetBlockRate:  2,
		BlockGasCostStep: big.NewInt(50_000),
	}
	BlockGasCostWithStepTest(t, testFeeConfig)

	testFeeConfigDouble := commontype.FeeConfig{
		MinBlockGasCost:  big.NewInt(2),
		MaxBlockGasCost:  big.NewInt(2_000_000),
		TargetBlockRate:  4,
		BlockGasCostStep: big.NewInt(100_000),
	}
	BlockGasCostWithStepTest(t, testFeeConfigDouble)
}

func BlockGasCostWithStepTest(t *testing.T, testFeeConfig commontype.FeeConfig) {
	minBlockGasCost := testFeeConfig.MinBlockGasCost.Uint64()
	blockGasCostStep := testFeeConfig.BlockGasCostStep.Uint64()
	targetBlockRate := testFeeConfig.TargetBlockRate
	bigMaxBlockGasCost := testFeeConfig.MaxBlockGasCost
	maxBlockGasCost := bigMaxBlockGasCost.Uint64()
	tests := []struct {
		name        string
		parentCost  *big.Int
		timeElapsed uint64
		expected    uint64
	}{
		{
			name:        "nil_parentCost",
			parentCost:  nil,
			timeElapsed: 0,
			expected:    minBlockGasCost,
		},
		{
			name:        "timeElapsed_0",
			parentCost:  big.NewInt(0),
			timeElapsed: 0,
			expected:    targetBlockRate * blockGasCostStep,
		},
		{
			name:        "timeElapsed_1",
			parentCost:  big.NewInt(0),
			timeElapsed: 1,
			expected:    (targetBlockRate - 1) * blockGasCostStep,
		},
		{
			name:        "timeElapsed_0_with_parentCost",
			parentCost:  big.NewInt(50_000),
			timeElapsed: 0,
			expected:    50_000 + targetBlockRate*blockGasCostStep,
		},
		{
			name:        "timeElapsed_0_with_max_parentCost",
			parentCost:  bigMaxBlockGasCost,
			timeElapsed: 0,
			expected:    maxBlockGasCost,
		},
		{
			name:        "timeElapsed_1_with_max_parentCost",
			parentCost:  bigMaxBlockGasCost,
			timeElapsed: 1,
			expected:    maxBlockGasCost,
		},
		{
			name:        "timeElapsed_at_target",
			parentCost:  big.NewInt(900_000),
			timeElapsed: targetBlockRate,
			expected:    900_000,
		},
		{
			name:        "timeElapsed_over_target_3",
			parentCost:  bigMaxBlockGasCost,
			timeElapsed: 3,
			expected:    maxBlockGasCost - (3-targetBlockRate)*blockGasCostStep,
		},
		{
			name:        "timeElapsed_over_target_10",
			parentCost:  bigMaxBlockGasCost,
			timeElapsed: 10,
			expected:    maxBlockGasCost - (10-targetBlockRate)*blockGasCostStep,
		},
		{
			name:        "timeElapsed_over_target_20",
			parentCost:  bigMaxBlockGasCost,
			timeElapsed: 20,
			expected:    maxBlockGasCost - (20-targetBlockRate)*blockGasCostStep,
		},
		{
			name:        "timeElapsed_over_target_22",
			parentCost:  bigMaxBlockGasCost,
			timeElapsed: 22,
			expected:    maxBlockGasCost - (22-targetBlockRate)*blockGasCostStep,
		},
		{
			name:        "timeElapsed_large_clamped_to_0",
			parentCost:  bigMaxBlockGasCost,
			timeElapsed: 23,
			expected:    0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, BlockGasCostWithStep(
				testFeeConfig,
				test.parentCost,
				blockGasCostStep,
				test.timeElapsed,
			))
		})
	}
}
