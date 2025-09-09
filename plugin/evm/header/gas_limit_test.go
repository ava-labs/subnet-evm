// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package header

import (
	"math/big"
	"testing"

	"github.com/ava-labs/avalanchego/utils/math"
	"github.com/ava-labs/avalanchego/vms/evm/upgrade/acp176"

	"github.com/ava-labs/libevm/core/types"
	"github.com/stretchr/testify/require"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/params/extras"

	ethparams "github.com/ava-labs/libevm/params"
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
		upgrades  extras.NetworkUpgrades
		parent    *types.Header
		timestamp uint64
		want      uint64
		wantErr   error
	}{
		{
			name:     "subnet_evm",
			upgrades: extras.TestSubnetEVMChainConfig.NetworkUpgrades,
			want:     feeConfig.GasLimit.Uint64(),
		},
		{
			name:     "fortuna_invalid_parent_header",
			upgrades: extras.TestFortunaChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number: big.NewInt(1),
			},
			wantErr: acp176.ErrStateInsufficientLength,
		},
		{
			name:     "fortuna_initial_max_capacity",
			upgrades: extras.TestFortunaChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number: big.NewInt(0),
			},
			want: acp176.MinMaxCapacity,
		},
		{
			name:     "pre_subnet_evm",
			upgrades: extras.TestPreSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				GasLimit: 1,
			},
			want: 1, // Same as parent
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)

			config := &extras.ChainConfig{
				NetworkUpgrades: test.upgrades,
			}
			got, err := GasLimit(config, feeConfig, test.parent, test.timestamp)
			require.ErrorIs(err, test.wantErr)
			require.Equal(test.want, got)
		})
	}
}

func TestVerifyGasUsed(t *testing.T) {
	tests := []struct {
		name            string
		feeConfig       commontype.FeeConfig
		acp224FeeConfig commontype.ACP224FeeConfig
		upgrades        extras.NetworkUpgrades
		parent          *types.Header
		header          *types.Header
		want            error
	}{
		{
			name:     "fortuna_gas_used_overflow",
			upgrades: extras.TestFortunaChainConfig.NetworkUpgrades,
			header: &types.Header{
				GasUsed: math.MaxUint[uint64](),
			},
			want: math.ErrOverflow,
		},
		{
			name:     "fortuna_invalid_capacity",
			upgrades: extras.TestFortunaChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number: big.NewInt(1),
			},
			header: &types.Header{},
			want:   acp176.ErrStateInsufficientLength,
		},
		{
			name:     "fortuna_invalid_usage",
			upgrades: extras.TestFortunaChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number: big.NewInt(0),
			},
			header: &types.Header{
				Time: 1,
				// The maximum allowed gas used is:
				// (header.Time - parent.Time) * [acp176.MinMaxPerSecond]
				// which is equal to [acp176.MinMaxPerSecond].
				GasUsed: acp176.MinMaxPerSecond + 1,
			},
			want: errInvalidGasUsed,
		},
		{
			name:     "fortuna_max_consumption",
			upgrades: extras.TestFortunaChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number: big.NewInt(0),
			},
			header: &types.Header{
				Time:    1,
				GasUsed: acp176.MinMaxPerSecond,
			},
			want: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := &extras.ChainConfig{
				NetworkUpgrades: test.upgrades,
			}
			err := VerifyGasUsed(config, test.feeConfig, test.parent, test.header)
			require.ErrorIs(t, err, test.want)
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
		upgrades extras.NetworkUpgrades
		parent   *types.Header
		header   *types.Header
		want     error
	}{
		{
			name:     "fortuna_invalid_header",
			upgrades: extras.TestFortunaChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number: big.NewInt(1),
			},
			header: &types.Header{},
			want:   acp176.ErrStateInsufficientLength,
		},
		{
			name:     "fortuna_invalid",
			upgrades: extras.TestFortunaChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number: big.NewInt(0),
			},
			header: &types.Header{
				GasLimit: acp176.MinMaxCapacity + 1,
			},
			want: errInvalidGasLimit,
		},
		{
			name:     "fortuna_valid",
			upgrades: extras.TestFortunaChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number: big.NewInt(0),
			},
			header: &types.Header{
				GasLimit: acp176.MinMaxCapacity,
			},
		},
		{
			name:     "subnet_evm_valid",
			upgrades: extras.TestSubnetEVMChainConfig.NetworkUpgrades,
			header: &types.Header{
				GasLimit: feeConfig.GasLimit.Uint64(),
			},
		},
		{
			name:     "subnet_evm_invalid",
			upgrades: extras.TestSubnetEVMChainConfig.NetworkUpgrades,
			header: &types.Header{
				GasLimit: feeConfig.GasLimit.Uint64() + 1,
			},
			want: errInvalidGasLimit,
		},
		{
			name:     "pre_subnet_evm_valid",
			upgrades: extras.TestPreSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				GasLimit: 50_000,
			},
			header: &types.Header{
				GasLimit: 50_001, // Gas limit is allowed to change by 1/1024
			},
		},
		{
			name:     "pre_subnet_evm_too_low",
			upgrades: extras.TestPreSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				GasLimit: ethparams.MinGasLimit,
			},
			header: &types.Header{
				GasLimit: ethparams.MinGasLimit - 1,
			},
			want: errInvalidGasLimit,
		},
		{
			name:     "pre_subnet_evm_too_high",
			upgrades: extras.TestPreSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				GasLimit: ethparams.MaxGasLimit,
			},
			header: &types.Header{
				GasLimit: ethparams.MaxGasLimit + 1,
			},
			want: errInvalidGasLimit,
		},
		{
			name:     "pre_subnet_evm_too_large",
			upgrades: extras.TestPreSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				GasLimit: ethparams.MinGasLimit,
			},
			header: &types.Header{
				GasLimit: ethparams.MaxGasLimit,
			},
			want: errInvalidGasLimit,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := &extras.ChainConfig{
				NetworkUpgrades: test.upgrades,
			}
			err := VerifyGasLimit(config, feeConfig, test.parent, test.header)
			require.ErrorIs(t, err, test.want)
		})
	}
}

func TestGasCapacity(t *testing.T) {
	tests := []struct {
		name      string
		feeConfig commontype.FeeConfig
		upgrades  extras.NetworkUpgrades
		parent    *types.Header
		timestamp uint64
		want      uint64
		wantErr   error
	}{
		{
			name:     "subnet_evm",
			upgrades: extras.TestSubnetEVMChainConfig.NetworkUpgrades,
			want:     0, // TODO: XXX Handle feeConfig with Fortuna here
		},
		{
			name:     "fortuna_invalid_header",
			upgrades: extras.TestFortunaChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number: big.NewInt(1),
			},
			wantErr: acp176.ErrStateInsufficientLength,
		},
		{
			name:     "fortuna_after_1s",
			upgrades: extras.TestFortunaChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number: big.NewInt(0),
			},
			timestamp: 1,
			want:      acp176.MinMaxPerSecond,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)

			config := &extras.ChainConfig{
				NetworkUpgrades: test.upgrades,
			}
			got, err := GasCapacity(config, test.feeConfig, test.parent, test.timestamp)
			require.ErrorIs(err, test.wantErr)
			require.Equal(test.want, got)
		})
	}
}
