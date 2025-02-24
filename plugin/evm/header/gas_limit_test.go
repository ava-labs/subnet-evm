// (c) 2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package header

import (
	"testing"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/stretchr/testify/require"
)

func TestGasLimit(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		GasLimitTest(t, testFeeConfig)
	})
	t.Run("double", func(t *testing.T) {
		GasLimitTest(t, testFeeConfigDouble)
	})
}

func GasLimitTest(t *testing.T, feeConfig commontype.FeeConfig) {
	tests := []struct {
		name      string
		upgrades  params.NetworkUpgrades
		calcGas   CalculateGasLimitFunc
		parent    *types.Header
		timestamp uint64
		want      uint64
	}{
		{
			name:     "subnet_evm",
			upgrades: params.TestSubnetEVMChainConfig.NetworkUpgrades,
			want:     feeConfig.GasLimit.Uint64(),
		},
		{
			name:     "pre_subnet_evm",
			upgrades: params.TestPreSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				GasLimit: 1,
			},
			calcGas: func(parentGasUsed, parentGasLimit, gasFloor, gasCeil uint64) uint64 {
				return 1
			},
			want: 1, // Same as parent
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := &params.ChainConfig{
				NetworkUpgrades: test.upgrades,
			}
			got := GasLimit(config, feeConfig, test.parent, test.calcGas, test.timestamp)
			require.Equal(t, test.want, got)
		})
	}
}

func TestVerifyGasLimit(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		VerifyGasLimitTest(t, testFeeConfig)
	})
	t.Run("double", func(t *testing.T) {
		VerifyGasLimitTest(t, testFeeConfigDouble)
	})
}

func VerifyGasLimitTest(t *testing.T, feeConfig commontype.FeeConfig) {
	tests := []struct {
		name     string
		upgrades params.NetworkUpgrades
		parent   *types.Header
		header   *types.Header
		want     error
	}{
		{
			name:     "subnet_evm_valid",
			upgrades: params.TestSubnetEVMChainConfig.NetworkUpgrades,
			header: &types.Header{
				GasLimit: feeConfig.GasLimit.Uint64(),
			},
		},
		{
			name:     "subnet_evm_invalid",
			upgrades: params.TestSubnetEVMChainConfig.NetworkUpgrades,
			header: &types.Header{
				GasLimit: feeConfig.GasLimit.Uint64() + 1,
			},
			want: errInvalidGasLimit,
		},
		{
			name:     "pre_subnet_evm_valid",
			upgrades: params.TestPreSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				GasLimit: 50_000,
			},
			header: &types.Header{
				GasLimit: 50_001, // Gas limit is allowed to change by 1/1024
			},
		},
		{
			name:     "pre_subnet_evm_too_low",
			upgrades: params.TestPreSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				GasLimit: params.MinGasLimit,
			},
			header: &types.Header{
				GasLimit: params.MinGasLimit - 1,
			},
			want: errInvalidGasLimit,
		},
		{
			name:     "pre_subnet_evm_too_high",
			upgrades: params.TestPreSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				GasLimit: params.MaxGasLimit,
			},
			header: &types.Header{
				GasLimit: params.MaxGasLimit + 1,
			},
			want: errInvalidGasLimit,
		},
		{
			name:     "pre_subnet_evm_too_large",
			upgrades: params.TestPreSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				GasLimit: params.MinGasLimit,
			},
			header: &types.Header{
				GasLimit: params.MaxGasLimit,
			},
			want: errInvalidGasLimit,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := &params.ChainConfig{
				NetworkUpgrades: test.upgrades,
			}
			err := VerifyGasLimit(config, feeConfig, test.parent, test.header)
			require.ErrorIs(t, err, test.want)
		})
	}
}
