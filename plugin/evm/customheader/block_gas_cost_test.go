// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package customheader

import (
	"math"
	"math/big"
	"testing"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/params/extras"
	"github.com/ava-labs/subnet-evm/plugin/evm/customtypes"
	"github.com/ava-labs/subnet-evm/utils"
)

var (
	testFeeConfig = commontype.FeeConfig{
		MinBlockGasCost:          big.NewInt(0),
		MaxBlockGasCost:          big.NewInt(1_000_000),
		TargetBlockRate:          2,
		BlockGasCostStep:         big.NewInt(50_000),
		TargetGas:                big.NewInt(10_000_000),
		BaseFeeChangeDenominator: big.NewInt(12),
		MinBaseFee:               big.NewInt(25 * utils.GWei),
		GasLimit:                 big.NewInt(12_000_000),
	}

	testFeeConfigDouble = commontype.FeeConfig{
		MinBlockGasCost:          big.NewInt(2),
		MaxBlockGasCost:          big.NewInt(2_000_000),
		TargetBlockRate:          4,
		BlockGasCostStep:         big.NewInt(100_000),
		TargetGas:                big.NewInt(20_000_000),
		BaseFeeChangeDenominator: big.NewInt(24),
		MinBaseFee:               big.NewInt(50 * utils.GWei),
		GasLimit:                 big.NewInt(24_000_000),
	}
)

func TestBlockGasCost(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		BlockGasCostTest(t, testFeeConfig)
	})
	t.Run("double", func(t *testing.T) {
		BlockGasCostTest(t, testFeeConfigDouble)
	})
}

func BlockGasCostTest(t *testing.T, feeConfig commontype.FeeConfig) {
	maxBlockGasCostBig := feeConfig.MaxBlockGasCost
	maxBlockGasCost := feeConfig.MaxBlockGasCost.Uint64()
	blockGasCostStep := feeConfig.BlockGasCostStep.Uint64()
	minBlockGasCost := feeConfig.MinBlockGasCost.Uint64()
	targetBlockRate := feeConfig.TargetBlockRate

	tests := []struct {
		name       string
		upgrades   extras.NetworkUpgrades
		parentTime uint64
		parentCost *big.Int
		timestamp  uint64
		expected   *big.Int
	}{
		{
			name:       "before_subnet_evm",
			parentTime: 10,
			upgrades:   extras.TestPreSubnetEVMChainConfig.NetworkUpgrades,
			parentCost: maxBlockGasCostBig,
			timestamp:  10 + targetBlockRate + 1,
			expected:   nil,
		},
		{
			name:       "normal",
			upgrades:   extras.TestChainConfig.NetworkUpgrades,
			parentTime: 10,
			parentCost: maxBlockGasCostBig,
			timestamp:  10 + targetBlockRate + 1,
			expected:   new(big.Int).SetUint64(maxBlockGasCost - blockGasCostStep),
		},
		{
			name:       "negative_time_elapsed",
			upgrades:   extras.TestChainConfig.NetworkUpgrades,
			parentTime: 10,
			parentCost: feeConfig.MinBlockGasCost,
			timestamp:  9,
			expected:   new(big.Int).SetUint64(minBlockGasCost + blockGasCostStep*targetBlockRate),
		},
		{
			name:       "granite_returns_zero",
			upgrades:   extras.TestGraniteChainConfig.NetworkUpgrades,
			parentTime: 10,
			parentCost: big.NewInt(int64(maxBlockGasCost)),
			timestamp:  10 + targetBlockRate + 1,
			expected:   big.NewInt(0),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := &extras.ChainConfig{
				NetworkUpgrades: test.upgrades,
			}
			parent := customtypes.WithHeaderExtra(
				&types.Header{
					Time: test.parentTime,
				},
				&customtypes.HeaderExtra{
					BlockGasCost: test.parentCost,
				},
			)

			assert.Equal(t, test.expected, BlockGasCost(
				config,
				feeConfig,
				parent,
				test.timestamp,
			))
		})
	}
}

func TestBlockGasCostWithStep(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		BlockGasCostWithStepTest(t, testFeeConfig)
	})
	t.Run("double", func(t *testing.T) {
		BlockGasCostWithStepTest(t, testFeeConfigDouble)
	})
}

func BlockGasCostWithStepTest(t *testing.T, feeConfig commontype.FeeConfig) {
	minBlockGasCost := feeConfig.MinBlockGasCost.Uint64()
	blockGasCostStep := feeConfig.BlockGasCostStep.Uint64()
	targetBlockRate := feeConfig.TargetBlockRate
	bigMaxBlockGasCost := feeConfig.MaxBlockGasCost
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
			name:        "timeElapsed_over_target_1",
			parentCost:  bigMaxBlockGasCost,
			timeElapsed: targetBlockRate + 1,
			expected:    maxBlockGasCost - blockGasCostStep,
		},
		{
			name:        "timeElapsed_over_target_10",
			parentCost:  bigMaxBlockGasCost,
			timeElapsed: targetBlockRate + 10,
			expected:    maxBlockGasCost - 10*blockGasCostStep,
		},
		{
			name:        "timeElapsed_over_target_15",
			parentCost:  bigMaxBlockGasCost,
			timeElapsed: targetBlockRate + 15,
			expected:    maxBlockGasCost - 15*blockGasCostStep,
		},
		{
			name:        "timeElapsed_large_clamped_to_0",
			parentCost:  bigMaxBlockGasCost,
			timeElapsed: targetBlockRate + 100,
			expected:    minBlockGasCost,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, BlockGasCostWithStep(
				feeConfig,
				test.parentCost,
				blockGasCostStep,
				test.timeElapsed,
			))
		})
	}
}

func TestEstimateRequiredTip(t *testing.T) {
	tests := []struct {
		name               string
		subnetEVMTimestamp *uint64
		header             *types.Header
		want               *big.Int
		wantErr            error
	}{
		{
			name:               "not_subnet_evm",
			subnetEVMTimestamp: utils.NewUint64(1),
			header:             &types.Header{},
		},
		{
			name:               "nil_base_fee",
			subnetEVMTimestamp: utils.NewUint64(0),
			header: customtypes.WithHeaderExtra(
				&types.Header{},
				&customtypes.HeaderExtra{
					BlockGasCost: big.NewInt(1),
				},
			),
			wantErr: errBaseFeeNil,
		},
		{
			name:               "nil_block_gas_cost",
			subnetEVMTimestamp: utils.NewUint64(0),
			header: &types.Header{
				BaseFee: big.NewInt(1),
			},
			wantErr: errBlockGasCostNil,
		},
		{
			name:               "no_gas_used",
			subnetEVMTimestamp: utils.NewUint64(0),
			header: customtypes.WithHeaderExtra(
				&types.Header{
					GasUsed: 0,
					BaseFee: big.NewInt(1),
				},
				&customtypes.HeaderExtra{
					BlockGasCost: big.NewInt(1),
				},
			),
			wantErr: errNoGasUsed,
		},
		{
			name:               "success",
			subnetEVMTimestamp: utils.NewUint64(0),
			header: customtypes.WithHeaderExtra(
				&types.Header{
					GasUsed: 912,
					BaseFee: big.NewInt(456),
				},
				&customtypes.HeaderExtra{
					BlockGasCost: big.NewInt(101112),
				},
			),
			// totalRequiredTips = BlockGasCost * BaseFee
			// estimatedTip = totalRequiredTips / GasUsed
			want: big.NewInt((101112 * 456) / (912)),
		},
		{
			name:               "success_rounds_up",
			subnetEVMTimestamp: utils.NewUint64(0),
			header: customtypes.WithHeaderExtra(
				&types.Header{
					GasUsed: 124,
					BaseFee: big.NewInt(456),
				},
				&customtypes.HeaderExtra{
					BlockGasCost: big.NewInt(101112),
				},
			),
			// totalGasUsed = GasUsed + ExtDataGasUsed
			// totalRequiredTips = BlockGasCost * BaseFee
			// estimatedTip = totalRequiredTips / totalGasUsed
			want: big.NewInt((101112*456)/(124) + 1), // +1 to round up
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)

			config := &extras.ChainConfig{
				NetworkUpgrades: extras.NetworkUpgrades{
					SubnetEVMTimestamp: test.subnetEVMTimestamp,
				},
			}
			requiredTip, err := EstimateRequiredTip(config, test.header)
			require.ErrorIs(err, test.wantErr)
			require.Equal(test.want, requiredTip)
		})
	}
}

func TestVerifyBlockFee(t *testing.T) {
	testAddr := common.HexToAddress("7ef5a6135f1fd6a02593eedc869c6d41d934aef8")
	tests := map[string]struct {
		baseFee                *big.Int
		parentBlockGasCost     *big.Int
		timeElapsed            uint64
		txs                    []*types.Transaction
		receipts               []*types.Receipt
		extraStateContribution *big.Int
		shouldErr              bool
	}{
		"tx only base fee": {
			baseFee:            big.NewInt(100),
			parentBlockGasCost: big.NewInt(0),
			timeElapsed:        0,
			txs: []*types.Transaction{
				types.NewTransaction(0, testAddr, big.NewInt(0), 100, big.NewInt(100), nil),
			},
			receipts: []*types.Receipt{
				{GasUsed: 1000},
			},
			extraStateContribution: nil,
			shouldErr:              true,
		},
		"tx covers exactly block fee": {
			baseFee:            big.NewInt(100),
			parentBlockGasCost: big.NewInt(0),
			timeElapsed:        0,
			txs: []*types.Transaction{
				types.NewTransaction(0, testAddr, big.NewInt(0), 100_000, big.NewInt(200), nil),
			},
			receipts: []*types.Receipt{
				{GasUsed: 100_000},
			},
			extraStateContribution: nil,
			shouldErr:              false,
		},
		"txs share block fee": {
			baseFee:            big.NewInt(100),
			parentBlockGasCost: big.NewInt(0),
			timeElapsed:        0,
			txs: []*types.Transaction{
				types.NewTransaction(0, testAddr, big.NewInt(0), 100_000, big.NewInt(200), nil),
				types.NewTransaction(1, testAddr, big.NewInt(0), 100_000, big.NewInt(100), nil),
			},
			receipts: []*types.Receipt{
				{GasUsed: 100_000},
				{GasUsed: 100_000},
			},
			extraStateContribution: nil,
			shouldErr:              false,
		},
		"txs split block fee": {
			baseFee:            big.NewInt(100),
			parentBlockGasCost: big.NewInt(0),
			timeElapsed:        0,
			txs: []*types.Transaction{
				types.NewTransaction(0, testAddr, big.NewInt(0), 100_000, big.NewInt(150), nil),
				types.NewTransaction(1, testAddr, big.NewInt(0), 100_000, big.NewInt(150), nil),
			},
			receipts: []*types.Receipt{
				{GasUsed: 100_000},
				{GasUsed: 100_000},
			},
			extraStateContribution: nil,
			shouldErr:              false,
		},
		"split block fee with extra state contribution": {
			baseFee:            big.NewInt(100),
			parentBlockGasCost: big.NewInt(0),
			timeElapsed:        0,
			txs: []*types.Transaction{
				types.NewTransaction(0, testAddr, big.NewInt(0), 100_000, big.NewInt(150), nil),
			},
			receipts: []*types.Receipt{
				{GasUsed: 100_000},
			},
			extraStateContribution: big.NewInt(5_000_000),
			shouldErr:              false,
		},
		"extra state contribution insufficient": {
			baseFee:                big.NewInt(100),
			parentBlockGasCost:     big.NewInt(0),
			timeElapsed:            0,
			txs:                    nil,
			receipts:               nil,
			extraStateContribution: big.NewInt(9_999_999),
			shouldErr:              true,
		},
		"negative extra state contribution": {
			baseFee:                big.NewInt(100),
			parentBlockGasCost:     big.NewInt(0),
			timeElapsed:            0,
			txs:                    nil,
			receipts:               nil,
			extraStateContribution: big.NewInt(-1),
			shouldErr:              true,
		},
		"extra state contribution covers block fee": {
			baseFee:                big.NewInt(100),
			parentBlockGasCost:     big.NewInt(0),
			timeElapsed:            0,
			txs:                    nil,
			receipts:               nil,
			extraStateContribution: big.NewInt(10_000_000),
			shouldErr:              false,
		},
		"extra state contribution covers more than block fee": {
			baseFee:                big.NewInt(100),
			parentBlockGasCost:     big.NewInt(0),
			timeElapsed:            0,
			txs:                    nil,
			receipts:               nil,
			extraStateContribution: big.NewInt(10_000_001),
			shouldErr:              false,
		},
		"tx only base fee after full time window": {
			baseFee:            big.NewInt(100),
			parentBlockGasCost: big.NewInt(500_000),
			timeElapsed:        testFeeConfig.TargetBlockRate + 10,
			txs: []*types.Transaction{
				types.NewTransaction(0, testAddr, big.NewInt(0), 100, big.NewInt(100), nil),
			},
			receipts: []*types.Receipt{
				{GasUsed: 1000},
			},
			extraStateContribution: nil,
			shouldErr:              false,
		},
		"tx only base fee after large time window": {
			baseFee:            big.NewInt(100),
			parentBlockGasCost: big.NewInt(100_000),
			timeElapsed:        math.MaxUint64,
			txs: []*types.Transaction{
				types.NewTransaction(0, testAddr, big.NewInt(0), 100, big.NewInt(100), nil),
			},
			receipts: []*types.Receipt{
				{GasUsed: 1000},
			},
			extraStateContribution: nil,
			shouldErr:              false,
		},
		"zero block gas cost": {
			baseFee:                big.NewInt(100),
			parentBlockGasCost:     big.NewInt(0),
			timeElapsed:            testFeeConfig.TargetBlockRate + 1,
			txs:                    nil,
			receipts:               nil,
			extraStateContribution: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			blockGasCost := BlockGasCostWithStep(
				testFeeConfig,
				test.parentBlockGasCost,
				testFeeConfig.BlockGasCostStep.Uint64(),
				test.timeElapsed,
			)
			bigBlockGasCost := new(big.Int).SetUint64(blockGasCost)

			if err := VerifyBlockFee(test.baseFee, bigBlockGasCost, test.txs, test.receipts, test.extraStateContribution); err != nil {
				if !test.shouldErr {
					t.Fatalf("Unexpected error: %s", err)
				}
			} else {
				if test.shouldErr {
					t.Fatal("Should have failed verification")
				}
			}
		})
	}
}
