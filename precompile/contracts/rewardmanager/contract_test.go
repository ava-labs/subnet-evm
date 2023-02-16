// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package rewardmanager

import (
	"math/big"
	"testing"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/subnet-evm/constants"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var testBlockNumber = big.NewInt(7)

func TestRewardManagerRun(t *testing.T) {
	adminAddr := common.BigToAddress(common.Big0)
	enabledAddr := common.BigToAddress(common.Big1)
	noRoleAddr := common.BigToAddress(common.Big2)
	testAddr := common.HexToAddress("0x0123")

	for name, test := range map[string]contract.PrecompileTest{
		"set allow fee recipients from no role fails": {
			Caller: noRoleAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackAllowFeeRecipients()
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: AllowFeeRecipientsGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotAllowFeeRecipients.Error(),
		},
		"set reward address from no role fails": {
			Caller: noRoleAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackSetRewardAddress(testAddr)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: SetRewardAddressGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotSetRewardAddress.Error(),
		},
		"disable rewards from no role fails": {
			Caller: noRoleAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackDisableRewards()
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: DisableRewardsGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotDisableRewards.Error(),
		},
		"set allow fee recipients from enabled succeeds": {
			Caller: enabledAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackAllowFeeRecipients()
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: AllowFeeRecipientsGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				_, isFeeRecipients := GetStoredRewardAddress(state)
				require.True(t, isFeeRecipients)
			},
		},
		"set reward address from enabled succeeds": {
			Caller: enabledAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackSetRewardAddress(testAddr)
				require.NoError(tt, err)

				return input
			}(t),
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
			Caller: enabledAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackDisableRewards()
				require.NoError(tt, err)

				return input
			}(t),
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
			Caller: noRoleAddr,
			BeforeHook: func(t *testing.T, state contract.StateDB) {
				StoreRewardAddress(state, testAddr)
			},
			Input: func(tt *testing.T) []byte {
				input, err := PackCurrentRewardAddress()
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: CurrentRewardAddressGasCost,
			ReadOnly:    false,
			ExpectedRes: func() []byte {
				res, err := PackCurrentRewardAddressOutput(testAddr)
				require.NoError(t, err)
				return res
			}(),
		},
		"get are fee recipients allowed from no role succeeds": {
			Caller: noRoleAddr,
			BeforeHook: func(t *testing.T, state contract.StateDB) {
				EnableAllowFeeRecipients(state)
			},
			Input: func(tt *testing.T) []byte {
				input, err := PackAreFeeRecipientsAllowed()
				require.NoError(tt, err)
				return input
			}(t),
			SuppliedGas: AreFeeRecipientsAllowedGasCost,
			ReadOnly:    false,
			ExpectedRes: func() []byte {
				res, err := PackAreFeeRecipientsAllowedOutput(true)
				require.NoError(t, err)
				return res
			}(),
		},
		"get initial config with address": {
			Caller: noRoleAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackCurrentRewardAddress()
				require.NoError(tt, err)
				return input
			}(t),
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
			Caller: noRoleAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackAreFeeRecipientsAllowed()
				require.NoError(tt, err)
				return input
			}(t),
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
			Caller: enabledAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackAllowFeeRecipients()
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: AllowFeeRecipientsGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly set reward addresss with allowed role fails": {
			Caller: enabledAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackSetRewardAddress(testAddr)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: SetRewardAddressGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"insufficient gas set reward address from allowed role": {
			Caller: enabledAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackSetRewardAddress(testAddr)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: SetRewardAddressGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas allow fee recipients from allowed role": {
			Caller: enabledAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackAllowFeeRecipients()
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: AllowFeeRecipientsGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas read current reward address from allowed role": {
			Caller: enabledAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackCurrentRewardAddress()
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: CurrentRewardAddressGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas are fee recipients allowed from allowed role": {
			Caller: enabledAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackAreFeeRecipientsAllowed()
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: AreFeeRecipientsAllowedGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			db := rawdb.NewMemoryDatabase()
			state, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
			require.NoError(t, err)

			// Set up the state so that each address has the expected permissions at the start.
			SetRewardManagerAllowListStatus(state, adminAddr, allowlist.AdminRole)
			SetRewardManagerAllowListStatus(state, enabledAddr, allowlist.EnabledRole)
			require.Equal(t, allowlist.AdminRole, GetRewardManagerAllowListStatus(state, adminAddr))
			require.Equal(t, allowlist.EnabledRole, GetRewardManagerAllowListStatus(state, enabledAddr))

			if test.BeforeHook != nil {
				test.BeforeHook(t, state)
			}

			blockContext := contract.NewMockBlockContext(testBlockNumber, 0)
			accesibleState := contract.NewMockAccessibleState(state, blockContext, snow.DefaultContextTest())

			if test.Config != nil {
				Module.Configure(nil, test.Config, state, blockContext)
			}
			ret, remainingGas, err := RewardManagerPrecompile.Run(accesibleState, test.Caller, ContractAddress, test.Input, test.SuppliedGas, test.ReadOnly)
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
