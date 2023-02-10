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
	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var testBlockNumber = big.NewInt(7)

type precompileTest struct {
	caller      common.Address
	input       func() []byte
	suppliedGas uint64
	readOnly    bool

	config config.Config

	preCondition func(t *testing.T, state *state.StateDB)
	assertState  func(t *testing.T, state *state.StateDB)

	expectedRes []byte
	expectedErr string
}

func TestRewardManagerRun(t *testing.T) {
	adminAddr := common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")
	enabledAddr := common.HexToAddress("0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B")
	noRoleAddr := common.HexToAddress("0xF60C45c607D0f41687c94C314d300f483661E13a")
	testAddr := common.HexToAddress("0x0123")

	for name, test := range map[string]precompileTest{
		"set allow fee recipients from no role fails": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := PackAllowFeeRecipients()
				require.NoError(t, err)

				return input
			},
			suppliedGas: AllowFeeRecipientsGasCost,
			readOnly:    false,
			expectedErr: ErrCannotAllowFeeRecipients.Error(),
		},
		"set reward address from no role fails": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := PackSetRewardAddress(testAddr)
				require.NoError(t, err)

				return input
			},
			suppliedGas: SetRewardAddressGasCost,
			readOnly:    false,
			expectedErr: ErrCannotSetRewardAddress.Error(),
		},
		"disable rewards from no role fails": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := PackDisableRewards()
				require.NoError(t, err)

				return input
			},
			suppliedGas: DisableRewardsGasCost,
			readOnly:    false,
			expectedErr: ErrCannotDisableRewards.Error(),
		},
		"set allow fee recipients from enabled succeeds": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := PackAllowFeeRecipients()
				require.NoError(t, err)

				return input
			},
			suppliedGas: AllowFeeRecipientsGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				_, isFeeRecipients := GetStoredRewardAddress(state)
				require.True(t, isFeeRecipients)
			},
		},
		"set reward address from enabled succeeds": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := PackSetRewardAddress(testAddr)
				require.NoError(t, err)

				return input
			},
			suppliedGas: SetRewardAddressGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				address, isFeeRecipients := GetStoredRewardAddress(state)
				require.Equal(t, testAddr, address)
				require.False(t, isFeeRecipients)
			},
		},
		"disable rewards from enabled succeeds": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := PackDisableRewards()
				require.NoError(t, err)

				return input
			},
			suppliedGas: DisableRewardsGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				address, isFeeRecipients := GetStoredRewardAddress(state)
				require.False(t, isFeeRecipients)
				require.Equal(t, constants.BlackholeAddr, address)
			},
		},
		"get current reward address from no role succeeds": {
			caller: noRoleAddr,
			preCondition: func(t *testing.T, state *state.StateDB) {
				StoreRewardAddress(state, testAddr)
			},
			input: func() []byte {
				input, err := PackCurrentRewardAddress()
				require.NoError(t, err)

				return input
			},
			suppliedGas: CurrentRewardAddressGasCost,
			readOnly:    false,
			expectedRes: func() []byte {
				res, err := PackCurrentRewardAddressOutput(testAddr)
				require.NoError(t, err)
				return res
			}(),
		},
		"get are fee recipients allowed from no role succeeds": {
			caller: noRoleAddr,
			preCondition: func(t *testing.T, state *state.StateDB) {
				EnableAllowFeeRecipients(state)
			},
			input: func() []byte {
				input, err := PackAreFeeRecipientsAllowed()
				require.NoError(t, err)
				return input
			},
			suppliedGas: AreFeeRecipientsAllowedGasCost,
			readOnly:    false,
			expectedRes: func() []byte {
				res, err := PackAreFeeRecipientsAllowedOutput(true)
				require.NoError(t, err)
				return res
			}(),
		},
		"get initial config with address": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := PackCurrentRewardAddress()
				require.NoError(t, err)
				return input
			},
			suppliedGas: CurrentRewardAddressGasCost,
			config: &RewardManagerConfig{
				InitialRewardConfig: &InitialRewardConfig{
					RewardAddress: testAddr,
				},
			},
			readOnly: false,
			expectedRes: func() []byte {
				res, err := PackCurrentRewardAddressOutput(testAddr)
				require.NoError(t, err)
				return res
			}(),
		},
		"get initial config with allow fee recipients enabled": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := PackAreFeeRecipientsAllowed()
				require.NoError(t, err)
				return input
			},
			suppliedGas: AreFeeRecipientsAllowedGasCost,
			config: &RewardManagerConfig{
				InitialRewardConfig: &InitialRewardConfig{
					AllowFeeRecipients: true,
				},
			},
			readOnly: false,
			expectedRes: func() []byte {
				res, err := PackAreFeeRecipientsAllowedOutput(true)
				require.NoError(t, err)
				return res
			}(),
		},
		"readOnly allow fee recipients with allowed role fails": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := PackAllowFeeRecipients()
				require.NoError(t, err)

				return input
			},
			suppliedGas: AllowFeeRecipientsGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly set reward addresss with allowed role fails": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := PackSetRewardAddress(testAddr)
				require.NoError(t, err)

				return input
			},
			suppliedGas: SetRewardAddressGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"insufficient gas set reward address from allowed role": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := PackSetRewardAddress(testAddr)
				require.NoError(t, err)

				return input
			},
			suppliedGas: SetRewardAddressGasCost - 1,
			readOnly:    false,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas allow fee recipients from allowed role": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := PackAllowFeeRecipients()
				require.NoError(t, err)

				return input
			},
			suppliedGas: AllowFeeRecipientsGasCost - 1,
			readOnly:    false,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas read current reward address from allowed role": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := PackCurrentRewardAddress()
				require.NoError(t, err)

				return input
			},
			suppliedGas: CurrentRewardAddressGasCost - 1,
			readOnly:    false,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas are fee recipients allowed from allowed role": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := PackAreFeeRecipientsAllowed()
				require.NoError(t, err)

				return input
			},
			suppliedGas: AreFeeRecipientsAllowedGasCost - 1,
			readOnly:    false,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
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

			if test.preCondition != nil {
				test.preCondition(t, state)
			}

			blockContext := contract.NewMockBlockContext(testBlockNumber, 0)
			accesibleState := contract.NewMockAccessibleState(state, blockContext, snow.DefaultContextTest())

			if test.config != nil {
				Module{}.Configure(nil, test.config, state, blockContext)
			}
			ret, remainingGas, err := RewardManagerPrecompile.Run(accesibleState, test.caller, ContractAddress, test.input(), test.suppliedGas, test.readOnly)
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
