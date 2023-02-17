// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package feemanager

import (
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/ethdb/memorydb"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/contracts/test_utils"
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

func newStateDB(t *testing.T) contract.StateDB {
	db := memorydb.New()
	stateDB, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
	require.NoError(t, err)
	return stateDB
}

func TestFeeManager(t *testing.T) {
	testBlockNumber := big.NewInt(7)
	tests := map[string]test_utils.PrecompileTest{
		"set config from no role fails": {
			Caller: allowlist.NoRoleAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotChangeFee.Error(),
		},
		"set config from enabled address": {
			Caller: allowlist.EnabledAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				require.Equal(t, testFeeConfig, feeConfig)
			},
		},
		"set invalid config from enabled address": {
			Caller: allowlist.EnabledAddr,
			InputFn: func(t *testing.T) []byte {
				feeConfig := testFeeConfig
				feeConfig.MinBlockGasCost = new(big.Int).Mul(feeConfig.MaxBlockGasCost, common.Big2)
				input, err := PackSetFeeConfig(feeConfig)
				require.NoError(t, err)

				return input
			},
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
			Caller: allowlist.AdminAddr,
			InputFn: func(tt *testing.T) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(tt, err)

				return input
			},
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			BlockNumber: testBlockNumber.Int64(),
			AfterHook: func(t *testing.T, state contract.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				require.Equal(t, testFeeConfig, feeConfig)
				lastChangedAt := GetFeeConfigLastChangedAt(state)
				require.EqualValues(t, testBlockNumber, lastChangedAt)
			},
		},
		"get fee config from non-enabled address": {
			Caller: allowlist.NoRoleAddr,
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
			Caller:      allowlist.NoRoleAddr,
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
			BlockNumber: testBlockNumber.Int64(),
			AfterHook: func(t *testing.T, state contract.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				lastChangedAt := GetFeeConfigLastChangedAt(state)
				require.Equal(t, testFeeConfig, feeConfig)
				require.EqualValues(t, testBlockNumber, lastChangedAt)
			},
		},
		"get last changed at from non-enabled address": {
			Caller: allowlist.NoRoleAddr,
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
			Caller: allowlist.NoRoleAddr,
			InputFn: func(tt *testing.T) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(tt, err)

				return input
			},
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly setFeeConfig with allow role fails": {
			Caller: allowlist.EnabledAddr,
			InputFn: func(tt *testing.T) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(tt, err)

				return input
			},
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly setFeeConfig with admin role fails": {
			Caller: allowlist.AdminAddr,
			InputFn: func(tt *testing.T) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(tt, err)

				return input
			},
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"insufficient gas setFeeConfig from admin": {
			Caller: allowlist.AdminAddr,
			InputFn: func(tt *testing.T) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(tt, err)

				return input
			},
			SuppliedGas: SetFeeConfigGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
	}

	allowlist.RunTestsWithAllowListSetup(t, Module, newStateDB, tests)
	allowlist.RunTestsWithAllowListSetup(t, Module, newStateDB, allowlist.AllowListTests(Module))
}
