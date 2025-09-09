// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package header

import (
	"math/big"
	"testing"

	"github.com/ava-labs/avalanchego/vms/components/gas"
	"github.com/ava-labs/avalanchego/vms/evm/upgrade/acp176"
	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/types"
	"github.com/stretchr/testify/require"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/params/extras"
	"github.com/ava-labs/subnet-evm/plugin/evm/upgrade/subnetevm"
	"github.com/ava-labs/subnet-evm/utils"
)

const (
	maxBaseFee = 225 * utils.GWei
)

func TestBaseFee(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		BaseFeeTest(t, testFeeConfig, testACP224FeeConfig)
	})
	t.Run("double", func(t *testing.T) {
		BaseFeeTest(t, testFeeConfigDouble, testACP224FeeConfigDouble)
	})
}

func BaseFeeTest(t *testing.T, feeConfig commontype.FeeConfig, acp224FeeConfig commontype.ACP224FeeConfig) {
	tests := []struct {
		name      string
		upgrades  extras.NetworkUpgrades
		parent    *types.Header
		timestamp uint64
		want      *big.Int
		wantErr   error
	}{
		{
			name:     "pre_subnet_evm",
			upgrades: extras.TestPreSubnetEVMChainConfig.NetworkUpgrades,
			want:     nil,
			wantErr:  nil,
		},
		{
			name: "subnet_evm_first_block",
			upgrades: extras.NetworkUpgrades{
				SubnetEVMTimestamp: utils.NewUint64(1),
			},
			parent: &types.Header{
				Number: big.NewInt(1),
			},
			timestamp: 1,
			want:      big.NewInt(feeConfig.MinBaseFee.Int64()),
		},
		{
			name:     "subnet_evm_genesis_block",
			upgrades: extras.TestSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number: big.NewInt(0),
			},
			want: big.NewInt(feeConfig.MinBaseFee.Int64()),
		},
		{
			name:     "subnet_evm_invalid_fee_window",
			upgrades: extras.TestSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number: big.NewInt(1),
			},
			wantErr: subnetevm.ErrWindowInsufficientLength,
		},
		{
			name:     "subnet_evm_invalid_timestamp",
			upgrades: extras.TestSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number: big.NewInt(1),
				Time:   1,
				Extra:  (&subnetevm.Window{}).Bytes(),
			},
			timestamp: 0,
			wantErr:   errInvalidTimestamp,
		},
		{
			name:     "subnet_evm_no_change",
			upgrades: extras.TestSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number:  big.NewInt(1),
				GasUsed: feeConfig.TargetGas.Uint64(),
				Time:    1,
				Extra:   (&subnetevm.Window{}).Bytes(),
				BaseFee: big.NewInt(feeConfig.MinBaseFee.Int64() + 1),
			},
			timestamp: 1,
			want:      big.NewInt(feeConfig.MinBaseFee.Int64() + 1),
		},
		{
			name:     "subnet_evm_small_decrease",
			upgrades: extras.TestSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number:  big.NewInt(1),
				Extra:   (&subnetevm.Window{}).Bytes(),
				BaseFee: big.NewInt(maxBaseFee),
			},
			timestamp: 1,
			want: func() *big.Int {
				var (
					gasTarget                  = feeConfig.TargetGas.Int64()
					gasUsed                    = int64(0)
					amountUnderTarget          = gasTarget - gasUsed
					parentBaseFee              = int64(maxBaseFee)
					smoothingFactor            = feeConfig.BaseFeeChangeDenominator.Int64()
					baseFeeFractionUnderTarget = amountUnderTarget * parentBaseFee / gasTarget
					delta                      = baseFeeFractionUnderTarget / smoothingFactor
					baseFee                    = parentBaseFee - delta
				)
				return big.NewInt(baseFee)
			}(),
		},
		{
			name:     "subnet_evm_large_decrease",
			upgrades: extras.TestSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number:  big.NewInt(1),
				Extra:   (&subnetevm.Window{}).Bytes(),
				BaseFee: big.NewInt(maxBaseFee),
			},
			timestamp: 2 * subnetevm.WindowLen,
			want: func() *big.Int {
				var (
					gasTarget                  = feeConfig.TargetGas.Int64()
					gasUsed                    = int64(0)
					amountUnderTarget          = gasTarget - gasUsed
					parentBaseFee              = int64(maxBaseFee)
					smoothingFactor            = feeConfig.BaseFeeChangeDenominator.Int64()
					baseFeeFractionUnderTarget = amountUnderTarget * parentBaseFee / gasTarget
					windowsElapsed             = int64(2)
					delta                      = windowsElapsed * baseFeeFractionUnderTarget / smoothingFactor
					baseFee                    = parentBaseFee - delta
				)
				return big.NewInt(baseFee)
			}(),
		},
		{
			name:     "subnet_evm_increase",
			upgrades: extras.TestSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number:  big.NewInt(1),
				GasUsed: 2 * feeConfig.TargetGas.Uint64(),
				Extra:   (&subnetevm.Window{}).Bytes(),
				BaseFee: big.NewInt(feeConfig.MinBaseFee.Int64()),
			},
			timestamp: 1,
			want: func() *big.Int {
				var (
					gasTarget                 = feeConfig.TargetGas.Int64()
					gasUsed                   = 2 * gasTarget
					amountOverTarget          = gasUsed - gasTarget
					parentBaseFee             = feeConfig.MinBaseFee.Int64()
					smoothingFactor           = feeConfig.BaseFeeChangeDenominator.Int64()
					baseFeeFractionOverTarget = amountOverTarget * parentBaseFee / gasTarget
					delta                     = baseFeeFractionOverTarget / smoothingFactor
					baseFee                   = parentBaseFee + delta
				)
				return big.NewInt(baseFee)
			}(),
		},
		{
			name:     "subnet_evm_big_1_not_modified",
			upgrades: extras.TestSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number:  big.NewInt(1),
				GasUsed: 1,
				Extra:   (&subnetevm.Window{}).Bytes(),
				BaseFee: big.NewInt(1),
			},
			timestamp: 2 * subnetevm.WindowLen,
			want:      big.NewInt(feeConfig.MinBaseFee.Int64()),
		},
		{
			name:     "fortuna_invalid_timestamp",
			upgrades: extras.TestFortunaChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number: big.NewInt(1),
				Time:   1,
				Extra:  (&acp176.State{}).Bytes(),
			},
			timestamp: 0,
			wantErr:   errInvalidTimestamp,
		},
		{
			name: "fortuna_first_block",
			upgrades: extras.NetworkUpgrades{
				FortunaTimestamp: utils.NewUint64(1),
			},
			parent: &types.Header{
				Number: big.NewInt(1),
			},
			timestamp: 1,
			want:      big.NewInt(acp176.MinGasPrice),
		},
		{
			name:     "fortuna_genesis_block",
			upgrades: extras.TestFortunaChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number: big.NewInt(0),
			},
			want: big.NewInt(acp176.MinGasPrice),
		},
		{
			name:     "fortuna_invalid_fee_state",
			upgrades: extras.TestFortunaChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number: big.NewInt(1),
				Extra:  make([]byte, acp176.StateSize-1),
			},
			wantErr: acp176.ErrStateInsufficientLength,
		},
		{
			name:     "fortuna_current",
			upgrades: extras.TestFortunaChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number: big.NewInt(1),
				Extra: (&acp176.State{
					Gas: gas.State{
						Excess: 2_704_386_192, // 1_500_000 * ln(nAVAX) * [acp176.TargetToPriceUpdateConversion]
					},
					TargetExcess: 13_605_152, // 2^25 * ln(1.5)
				}).Bytes(),
			},
			want: big.NewInt(1_000_000_002), // nAVAX + 2 due to rounding
		},
		{
			name:     "fortuna_decrease",
			upgrades: extras.TestFortunaChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number: big.NewInt(1),
				Extra: (&acp176.State{
					Gas: gas.State{
						Excess: 2_704_386_192, // 1_500_000 * ln(nAVAX) * [acp176.TargetToPriceUpdateConversion]
					},
					TargetExcess: 13_605_152, // 2^25 * ln(1.5)
				}).Bytes(),
			},
			timestamp: 1,
			want:      big.NewInt(988_571_555), // e^((2_704_386_192 - 1_500_000) / 1_500_000 / [acp176.TargetToPriceUpdateConversion])
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)

			config := &extras.ChainConfig{
				NetworkUpgrades: test.upgrades,
			}
			got, err := BaseFee(config, feeConfig, acp224FeeConfig, test.parent, test.timestamp)
			require.ErrorIs(err, test.wantErr)
			require.Equal(test.want, got)

			// Verify that [common.Big1] is not modified by [BaseFee].
			require.Equal(big.NewInt(1), common.Big1)
		})
	}
}

func TestEstimateNextBaseFee(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		EstimateNextBaseFeeTest(t, testFeeConfig, testACP224FeeConfig)
	})
	t.Run("double", func(t *testing.T) {
		EstimateNextBaseFeeTest(t, testFeeConfigDouble, testACP224FeeConfigDouble)
	})
}

func EstimateNextBaseFeeTest(t *testing.T, feeConfig commontype.FeeConfig, acp224FeeConfig commontype.ACP224FeeConfig) {
	testBaseFee := uint64(225 * utils.GWei)
	nilUpgrade := extras.NetworkUpgrades{}
	tests := []struct {
		name      string
		upgrades  extras.NetworkUpgrades
		parent    *types.Header
		timestamp uint64
		want      *big.Int
		wantErr   error
	}{
		{
			name:     "activated",
			upgrades: extras.TestSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number:  big.NewInt(1),
				Extra:   (&subnetevm.Window{}).Bytes(),
				BaseFee: new(big.Int).SetUint64(testBaseFee),
			},
			timestamp: 1,
			want: func() *big.Int {
				var (
					gasTarget                  = feeConfig.TargetGas.Uint64()
					gasUsed                    = uint64(0)
					amountUnderTarget          = gasTarget - gasUsed
					parentBaseFee              = testBaseFee
					smoothingFactor            = feeConfig.BaseFeeChangeDenominator.Uint64()
					baseFeeFractionUnderTarget = amountUnderTarget * parentBaseFee / gasTarget
					delta                      = baseFeeFractionUnderTarget / smoothingFactor
					baseFee                    = parentBaseFee - delta
				)
				return new(big.Int).SetUint64(baseFee)
			}(),
		},
		{
			name:     "not_scheduled",
			upgrades: nilUpgrade,
			wantErr:  errEstimateBaseFeeWithoutActivation,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)

			config := &extras.ChainConfig{
				NetworkUpgrades: test.upgrades,
			}
			got, err := EstimateNextBaseFee(config, feeConfig, acp224FeeConfig, test.parent, test.timestamp)
			require.ErrorIs(err, test.wantErr)
			require.Equal(test.want, got)
		})
	}
}
