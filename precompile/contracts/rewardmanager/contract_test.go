// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package rewardmanager

import (
	"testing"

	"github.com/ava-labs/subnet-evm/constants"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/ethdb/memorydb"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/contracts/test_utils"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func newStateDB(t *testing.T) contract.StateDB {
	db := memorydb.New()
	stateDB, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
	require.NoError(t, err)
	return stateDB
}

func TestRewardManagerRun(t *testing.T) {
	testAddr := common.HexToAddress("0x0123")

	tests := map[string]test_utils.PrecompileTest{
		"set allow fee recipients from no role fails": {
			Caller: allowlist.NoRoleAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackAllowFeeRecipients()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: AllowFeeRecipientsGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotAllowFeeRecipients.Error(),
		},
		"set reward address from no role fails": {
			Caller: allowlist.NoRoleAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackSetRewardAddress(testAddr)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetRewardAddressGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotSetRewardAddress.Error(),
		},
		"disable rewards from no role fails": {
			Caller: allowlist.NoRoleAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackDisableRewards()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: DisableRewardsGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotDisableRewards.Error(),
		},
		"set allow fee recipients from enabled succeeds": {
			Caller: allowlist.EnabledAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackAllowFeeRecipients()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: AllowFeeRecipientsGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				_, isFeeRecipients := GetStoredRewardAddress(state)
				require.True(t, isFeeRecipients)
			},
		},
		"set reward address from enabled succeeds": {
			Caller: allowlist.EnabledAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackSetRewardAddress(testAddr)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetRewardAddressGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				address, isFeeRecipients := GetStoredRewardAddress(state)
				require.Equal(t, testAddr, address)
				require.False(t, isFeeRecipients)
			},
		},
		"disable rewards from enabled succeeds": {
			Caller: allowlist.EnabledAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackDisableRewards()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: DisableRewardsGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				address, isFeeRecipients := GetStoredRewardAddress(state)
				require.False(t, isFeeRecipients)
				require.Equal(t, constants.BlackholeAddr, address)
			},
		},
		"get current reward address from no role succeeds": {
			Caller: allowlist.NoRoleAddr,
			BeforeHook: func(t *testing.T, state contract.StateDB) {
				StoreRewardAddress(state, testAddr)
			},
			InputFn: func(t *testing.T) []byte {
				input, err := PackCurrentRewardAddress()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: CurrentRewardAddressGasCost,
			ReadOnly:    false,
			ExpectedRes: func() []byte {
				res, err := PackCurrentRewardAddressOutput(testAddr)
				require.NoError(t, err)
				return res
			}(),
		},
		"get are fee recipients allowed from no role succeeds": {
			Caller: allowlist.NoRoleAddr,
			BeforeHook: func(t *testing.T, state contract.StateDB) {
				EnableAllowFeeRecipients(state)
			},
			InputFn: func(t *testing.T) []byte {
				input, err := PackAreFeeRecipientsAllowed()
				require.NoError(t, err)
				return input
			},
			SuppliedGas: AreFeeRecipientsAllowedGasCost,
			ReadOnly:    false,
			ExpectedRes: func() []byte {
				res, err := PackAreFeeRecipientsAllowedOutput(true)
				require.NoError(t, err)
				return res
			}(),
		},
		"get initial config with address": {
			Caller: allowlist.NoRoleAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackCurrentRewardAddress()
				require.NoError(t, err)
				return input
			},
			SuppliedGas: CurrentRewardAddressGasCost,
			Config: &Config{
				InitialRewardConfig: &InitialRewardConfig{
					RewardAddress: testAddr,
				},
			},
			ReadOnly: false,
			ExpectedRes: func() []byte {
				res, err := PackCurrentRewardAddressOutput(testAddr)
				require.NoError(t, err)
				return res
			}(),
		},
		"get initial config with allow fee recipients enabled": {
			Caller: allowlist.NoRoleAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackAreFeeRecipientsAllowed()
				require.NoError(t, err)
				return input
			},
			SuppliedGas: AreFeeRecipientsAllowedGasCost,
			Config: &Config{
				InitialRewardConfig: &InitialRewardConfig{
					AllowFeeRecipients: true,
				},
			},
			ReadOnly: false,
			ExpectedRes: func() []byte {
				res, err := PackAreFeeRecipientsAllowedOutput(true)
				require.NoError(t, err)
				return res
			}(),
		},
		"readOnly allow fee recipients with allowed role fails": {
			Caller: allowlist.EnabledAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackAllowFeeRecipients()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: AllowFeeRecipientsGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly set reward addresss with allowed role fails": {
			Caller: allowlist.EnabledAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackSetRewardAddress(testAddr)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetRewardAddressGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"insufficient gas set reward address from allowed role": {
			Caller: allowlist.EnabledAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackSetRewardAddress(testAddr)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetRewardAddressGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas allow fee recipients from allowed role": {
			Caller: allowlist.EnabledAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackAllowFeeRecipients()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: AllowFeeRecipientsGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas read current reward address from allowed role": {
			Caller: allowlist.EnabledAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackCurrentRewardAddress()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: CurrentRewardAddressGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas are fee recipients allowed from allowed role": {
			Caller: allowlist.EnabledAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackAreFeeRecipientsAllowed()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: AreFeeRecipientsAllowedGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
	}

	allowlist.RunTestsWithAllowListSetup(t, Module, newStateDB, tests)
	allowlist.RunTestsWithAllowListSetup(t, Module, newStateDB, allowlist.AllowListTests(Module))
}
