// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package feemanager

import (
	"math/big"
	"testing"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var testFeeConfig = commontype.FeeConfig{
	GasLimit:        big.NewInt(8_000_000),
	TargetBlockRate: 2, // in seconds

	MinBaseFee:               big.NewInt(25_000_000_000),
	TargetGas:                big.NewInt(15_000_000),
	BaseFeeChangeDenominator: big.NewInt(36),

	MinBlockGasCost:  big.NewInt(0),
	MaxBlockGasCost:  big.NewInt(1_000_000),
	BlockGasCostStep: big.NewInt(200_000),
}

func TestFeeManagerRun(t *testing.T) {
	testBlockNumber := big.NewInt(7)

	adminAddr := common.BigToAddress(common.Big0)
	enabledAddr := common.BigToAddress(common.Big1)
	noRoleAddr := common.BigToAddress(common.Big2)

	for name, test := range map[string]contract.PrecompileTest{
		"set config from no role fails": {
			Caller: noRoleAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotChangeFee.Error(),
		},
		"set config from enabled address": {
			Caller: enabledAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				require.Equal(t, testFeeConfig, feeConfig)
			},
		},
		"set invalid config from enabled address": {
			Caller: enabledAddr,
			Input: func(tt *testing.T) []byte {
				feeConfig := testFeeConfig
				feeConfig.MinBlockGasCost = new(big.Int).Mul(feeConfig.MaxBlockGasCost, common.Big2)
				input, err := PackSetFeeConfig(feeConfig)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    false,
			Config: &Config{
				InitialFeeConfig: &testFeeConfig,
			},
			ExpectedErr: "cannot be greater than maxBlockGasCost",
			AfterHook: func(t *testing.T, state contract.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				require.Equal(t, testFeeConfig, feeConfig)
			},
		},
		"set config from admin address": {
			Caller: adminAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				require.Equal(t, testFeeConfig, feeConfig)
				lastChangedAt := GetFeeConfigLastChangedAt(state)
				require.EqualValues(t, testBlockNumber, lastChangedAt)
			},
		},
		"get fee config from non-enabled address": {
			Caller: noRoleAddr,
			BeforeHook: func(t *testing.T, state contract.StateDB) {
				err := StoreFeeConfig(state, testFeeConfig, contract.NewMockBlockContext(big.NewInt(6), 0))
				require.NoError(t, err)
			},
			Input:       PackGetFeeConfigInput(),
			SuppliedGas: GetFeeConfigGasCost,
			ReadOnly:    true,
			ExpectedRes: func() []byte {
				res, err := PackFeeConfig(testFeeConfig)
				require.NoError(t, err)
				return res
			}(),
			AfterHook: func(t *testing.T, state contract.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				lastChangedAt := GetFeeConfigLastChangedAt(state)
				require.Equal(t, testFeeConfig, feeConfig)
				require.EqualValues(t, big.NewInt(6), lastChangedAt)
			},
		},
		"get initial fee config": {
			Caller:      noRoleAddr,
			Input:       PackGetFeeConfigInput(),
			SuppliedGas: GetFeeConfigGasCost,
			Config: &Config{
				InitialFeeConfig: &testFeeConfig,
			},
			ReadOnly: true,
			ExpectedRes: func() []byte {
				res, err := PackFeeConfig(testFeeConfig)
				require.NoError(t, err)
				return res
			}(),
			AfterHook: func(t *testing.T, state contract.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				lastChangedAt := GetFeeConfigLastChangedAt(state)
				require.Equal(t, testFeeConfig, feeConfig)
				require.EqualValues(t, testBlockNumber, lastChangedAt)
			},
		},
		"get last changed at from non-enabled address": {
			Caller: noRoleAddr,
			BeforeHook: func(t *testing.T, state contract.StateDB) {
				err := StoreFeeConfig(state, testFeeConfig, contract.NewMockBlockContext(testBlockNumber, 0))
				require.NoError(t, err)
			},
			Input:       PackGetLastChangedAtInput(),
			SuppliedGas: GetLastChangedAtGasCost,
			ReadOnly:    true,
			ExpectedRes: common.BigToHash(testBlockNumber).Bytes(),
			AfterHook: func(t *testing.T, state contract.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				lastChangedAt := GetFeeConfigLastChangedAt(state)
				require.Equal(t, testFeeConfig, feeConfig)
				require.Equal(t, testBlockNumber, lastChangedAt)
			},
		},
		"readOnly setFeeConfig with noRole fails": {
			Caller: noRoleAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly setFeeConfig with allow role fails": {
			Caller: enabledAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly setFeeConfig with admin role fails": {
			Caller: adminAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"insufficient gas setFeeConfig from admin": {
			Caller: adminAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: SetFeeConfigGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			db := rawdb.NewMemoryDatabase()
			state, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
			require.NoError(t, err)

			// Set up the state so that each address has the expected permissions at the start.
			SetFeeManagerStatus(state, adminAddr, allowlist.AdminRole)
			SetFeeManagerStatus(state, enabledAddr, allowlist.EnabledRole)
			require.Equal(t, allowlist.AdminRole, GetFeeManagerStatus(state, adminAddr))
			require.Equal(t, allowlist.EnabledRole, GetFeeManagerStatus(state, enabledAddr))

			if test.BeforeHook != nil {
				test.BeforeHook(t, state)
			}
			blockContext := contract.NewMockBlockContext(testBlockNumber, 0)
			accesibleState := contract.NewMockAccessibleState(state, blockContext, snow.DefaultContextTest())
			if test.Config != nil {
				Module.Configure(nil, test.Config, state, blockContext)
			}
			ret, remainingGas, err := FeeManagerPrecompile.Run(accesibleState, test.Caller, ContractAddress, test.Input, test.SuppliedGas, test.ReadOnly)
			if len(test.ExpectedErr) != 0 {
				require.ErrorContains(t, err, test.ExpectedErr)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, uint64(0), remainingGas)
			require.Equal(t, test.ExpectedRes, ret)

			if test.AfterHook != nil {
				test.AfterHook(t, state)
			}
		})
	}
}
