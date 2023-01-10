// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package feemanager

import (
	"math/big"
	"testing"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var (
	testFeeConfig = commontype.FeeConfig{
		GasLimit:        big.NewInt(8_000_000),
		TargetBlockRate: 2, // in seconds

		MinBaseFee:               big.NewInt(25_000_000_000),
		TargetGas:                big.NewInt(15_000_000),
		BaseFeeChangeDenominator: big.NewInt(36),

		MinBlockGasCost:  big.NewInt(0),
		MaxBlockGasCost:  big.NewInt(1_000_000),
		BlockGasCostStep: big.NewInt(200_000),
	}

	testBlockNumber = big.NewInt(7)
)

func TestFeeConfigManagerRun(t *testing.T) {
	type test struct {
		caller       common.Address
		preCondition func(t *testing.T, state *state.StateDB)
		input        func() []byte
		suppliedGas  uint64
		readOnly     bool
		config       *FeeConfigManagerConfig

		expectedRes []byte
		expectedErr string

		assertState func(t *testing.T, state *state.StateDB)
	}

	adminAddr := common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")
	enabledAddr := common.HexToAddress("0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B")
	noRoleAddr := common.HexToAddress("0xF60C45c607D0f41687c94C314d300f483661E13a")

	for name, test := range map[string]test{
		"set config from no role fails": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			suppliedGas: SetFeeConfigGasCost,
			readOnly:    false,
			expectedErr: ErrCannotChangeFee.Error(),
		},
		"set config from enabled address": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			suppliedGas: SetFeeConfigGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				require.Equal(t, testFeeConfig, feeConfig)
			},
		},
		"set invalid config from enabled address": {
			caller: enabledAddr,
			input: func() []byte {
				feeConfig := testFeeConfig
				feeConfig.MinBlockGasCost = new(big.Int).Mul(feeConfig.MaxBlockGasCost, common.Big2)
				input, err := PackSetFeeConfig(feeConfig)
				require.NoError(t, err)

				return input
			},
			suppliedGas: SetFeeConfigGasCost,
			readOnly:    false,
			expectedRes: nil,
			config: &FeeConfigManagerConfig{
				InitialFeeConfig: &testFeeConfig,
			},
			expectedErr: "cannot be greater than maxBlockGasCost",
			assertState: func(t *testing.T, state *state.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				require.Equal(t, testFeeConfig, feeConfig)
			},
		},
		"set config from admin address": {
			caller: adminAddr,
			input: func() []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			suppliedGas: SetFeeConfigGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				require.Equal(t, testFeeConfig, feeConfig)
				lastChangedAt := GetFeeConfigLastChangedAt(state)
				require.EqualValues(t, testBlockNumber, lastChangedAt)
			},
		},
		"get fee config from non-enabled address": {
			caller: noRoleAddr,
			preCondition: func(t *testing.T, state *state.StateDB) {
				err := StoreFeeConfig(state, testFeeConfig, precompile.NewMockBlockContext(big.NewInt(6), 0))
				require.NoError(t, err)
			},
			input: func() []byte {
				return PackGetFeeConfigInput()
			},
			suppliedGas: GetFeeConfigGasCost,
			readOnly:    true,
			expectedRes: func() []byte {
				res, err := PackFeeConfig(testFeeConfig)
				require.NoError(t, err)
				return res
			}(),
			assertState: func(t *testing.T, state *state.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				lastChangedAt := GetFeeConfigLastChangedAt(state)
				require.Equal(t, testFeeConfig, feeConfig)
				require.EqualValues(t, big.NewInt(6), lastChangedAt)
			},
		},
		"get initial fee config": {
			caller: noRoleAddr,
			input: func() []byte {
				return PackGetFeeConfigInput()
			},
			suppliedGas: GetFeeConfigGasCost,
			config: &FeeConfigManagerConfig{
				InitialFeeConfig: &testFeeConfig,
			},
			readOnly: true,
			expectedRes: func() []byte {
				res, err := PackFeeConfig(testFeeConfig)
				require.NoError(t, err)
				return res
			}(),
			assertState: func(t *testing.T, state *state.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				lastChangedAt := GetFeeConfigLastChangedAt(state)
				require.Equal(t, testFeeConfig, feeConfig)
				require.EqualValues(t, testBlockNumber, lastChangedAt)
			},
		},
		"get last changed at from non-enabled address": {
			caller: noRoleAddr,
			preCondition: func(t *testing.T, state *state.StateDB) {
				err := StoreFeeConfig(state, testFeeConfig, precompile.NewMockBlockContext(testBlockNumber, 0))
				require.NoError(t, err)
			},
			input: func() []byte {
				return PackGetLastChangedAtInput()
			},
			suppliedGas: GetLastChangedAtGasCost,
			readOnly:    true,
			expectedRes: common.BigToHash(testBlockNumber).Bytes(),
			assertState: func(t *testing.T, state *state.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				lastChangedAt := GetFeeConfigLastChangedAt(state)
				require.Equal(t, testFeeConfig, feeConfig)
				require.Equal(t, testBlockNumber, lastChangedAt)
			},
		},
		"readOnly setFeeConfig with noRole fails": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			suppliedGas: SetFeeConfigGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly setFeeConfig with allow role fails": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			suppliedGas: SetFeeConfigGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly setFeeConfig with admin role fails": {
			caller: adminAddr,
			input: func() []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			suppliedGas: SetFeeConfigGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"insufficient gas setFeeConfig from admin": {
			caller: adminAddr,
			input: func() []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			suppliedGas: SetFeeConfigGasCost - 1,
			readOnly:    false,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"set allow role from admin": {
			caller: adminAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(noRoleAddr, allowlist.AllowListEnabled)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := GetFeeConfigManagerStatus(state, noRoleAddr)
				require.Equal(t, allowlist.AllowListEnabled, res)
			},
		},
		"set allow role from non-admin fails": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(noRoleAddr, allowlist.AllowListEnabled)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			db := rawdb.NewMemoryDatabase()
			state, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
			require.NoError(t, err)

			// Set up the state so that each address has the expected permissions at the start.
			SetFeeConfigManagerStatus(state, adminAddr, allowlist.AllowListAdmin)
			SetFeeConfigManagerStatus(state, enabledAddr, allowlist.AllowListEnabled)
			SetFeeConfigManagerStatus(state, noRoleAddr, allowlist.AllowListEnabled)

			if test.preCondition != nil {
				test.preCondition(t, state)
			}
			blockContext := precompile.NewMockBlockContext(testBlockNumber, 0)
			accesibleState := precompile.NewMockAccessibleState(state, blockContext, snow.DefaultContextTest())
			if test.config != nil {
				test.config.Configure(nil, state, blockContext)
			}
			ret, remainingGas, err := FeeConfigManagerPrecompile.Run(accesibleState, test.caller, ContractAddress, test.input(), test.suppliedGas, test.readOnly)
			if len(test.expectedErr) != 0 {
				require.ErrorContains(t, err, test.expectedErr)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, uint64(0), remainingGas)
			require.Equal(t, test.expectedRes, ret)

			if test.assertState != nil {
				test.assertState(t, state)
			}
		})
	}
}
