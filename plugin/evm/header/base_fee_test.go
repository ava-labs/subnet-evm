// (c) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package header

import (
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/plugin/evm/upgrade/subnetevm"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const (
	maxBaseFee = 225 * utils.GWei
)

func TestBaseFee(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		BaseFeeTest(t, testFeeConfig)
	})
	t.Run("double", func(t *testing.T) {
		BaseFeeTest(t, testFeeConfigDouble)
	})
}

func BaseFeeTest(t *testing.T, feeConfig commontype.FeeConfig) {
	tests := []struct {
		name      string
		upgrades  params.NetworkUpgrades
		parent    *types.Header
		timestamp uint64
		want      *big.Int
		wantErr   error
	}{
		{
			name:     "pre_subnet_evm",
			upgrades: params.TestPreSubnetEVMChainConfig.NetworkUpgrades,
			want:     nil,
			wantErr:  nil,
		},
		{
			name: "subnet_evm_first_block",
			upgrades: params.NetworkUpgrades{
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
			upgrades: params.TestSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number: big.NewInt(0),
			},
			want: big.NewInt(feeConfig.MinBaseFee.Int64()),
		},
		{
			name:     "subnet_evm_invalid_fee_window",
			upgrades: params.TestSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number: big.NewInt(1),
			},
			wantErr: subnetevm.ErrWindowInsufficientLength,
		},
		{
			name:     "subnet_evm_invalid_timestamp",
			upgrades: params.TestSubnetEVMChainConfig.NetworkUpgrades,
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
			upgrades: params.TestSubnetEVMChainConfig.NetworkUpgrades,
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
			upgrades: params.TestSubnetEVMChainConfig.NetworkUpgrades,
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
			upgrades: params.TestSubnetEVMChainConfig.NetworkUpgrades,
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
			upgrades: params.TestSubnetEVMChainConfig.NetworkUpgrades,
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
			upgrades: params.TestSubnetEVMChainConfig.NetworkUpgrades,
			parent: &types.Header{
				Number:  big.NewInt(1),
				GasUsed: 1,
				Extra:   (&subnetevm.Window{}).Bytes(),
				BaseFee: big.NewInt(1),
			},
			timestamp: 2 * subnetevm.WindowLen,
			want:      big.NewInt(feeConfig.MinBaseFee.Int64()),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)

			config := &params.ChainConfig{
				NetworkUpgrades: test.upgrades,
			}
			got, err := BaseFee(config, feeConfig, test.parent, test.timestamp)
			require.ErrorIs(err, test.wantErr)
			require.Equal(test.want, got)

			// Verify that [common.Big1] is not modified by [BaseFee].
			require.Equal(big.NewInt(1), common.Big1)
		})
	}
}

func TestEstimateNextBaseFee(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		BlockGasCostTest(t, testFeeConfig)
	})
	t.Run("double", func(t *testing.T) {
		BlockGasCostTest(t, testFeeConfigDouble)
	})
}

func EstimateNextBaseFeeTest(t *testing.T, feeConfig commontype.FeeConfig) {
	testBaseFee := uint64(225 * utils.GWei)
	nilUpgrade := params.NetworkUpgrades{}
	tests := []struct {
		name      string
		upgrades  params.NetworkUpgrades
		parent    *types.Header
		timestamp uint64
		want      *big.Int
		wantErr   error
	}{
		{
			name:     "activated",
			upgrades: params.TestSubnetEVMChainConfig.NetworkUpgrades,
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

			config := &params.ChainConfig{
				NetworkUpgrades: test.upgrades,
			}
			got, err := EstimateNextBaseFee(config, feeConfig, test.parent, test.timestamp)
			require.ErrorIs(err, test.wantErr)
			require.Equal(test.want, got)
		})
	}
}
