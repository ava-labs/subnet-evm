// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package rewardmanager

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/constants"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	"github.com/ava-labs/subnet-evm/vmerrs"
)

var (
	rewardAddress = common.HexToAddress("0x0123")
	tests         = map[string]testutils.PrecompileTest{
		"set allow fee recipients from no role fails": {
			Caller:     allowlist.TestNoRoleAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackAllowFeeRecipients()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: AllowFeeRecipientsGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotAllowFeeRecipients.Error(),
		},
		"set reward address from no role fails": {
			Caller:     allowlist.TestNoRoleAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackSetRewardAddress(rewardAddress)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetRewardAddressGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotSetRewardAddress.Error(),
		},
		"disable rewards from no role fails": {
			Caller:     allowlist.TestNoRoleAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackDisableRewards()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: DisableRewardsGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotDisableRewards.Error(),
		},
		"set allow fee recipients from enabled succeeds": {
			Caller:     allowlist.TestEnabledAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackAllowFeeRecipients()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: AllowFeeRecipientsGasCost + FeeRecipientsAllowedEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, state contract.StateDB) {
				_, isFeeRecipients := GetStoredRewardAddress(state)
				require.True(t, isFeeRecipients)

				logsTopics, logsData := state.GetLogData()
				assertFeeRecipientsAllowed(t, logsTopics, logsData, allowlist.TestEnabledAddr)
			},
		},
		"set fee recipients should not emit events pre-Durango": {
			Caller:     allowlist.TestEnabledAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackAllowFeeRecipients()
				require.NoError(t, err)

				return input
			},
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				mockChainConfig := precompileconfig.NewMockChainConfig(ctrl)
				mockChainConfig.EXPECT().GetFeeConfig().AnyTimes().Return(commontype.ValidTestFeeConfig)
				mockChainConfig.EXPECT().AllowedFeeRecipients().AnyTimes().Return(false)
				mockChainConfig.EXPECT().IsDurango(gomock.Any()).AnyTimes().Return(false)
				return mockChainConfig
			},
			SuppliedGas: AllowFeeRecipientsGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, stateDB contract.StateDB) {
				// Check no logs are stored in state
				logsTopics, logsData := stateDB.GetLogData()
				require.Len(t, logsTopics, 0)
				require.Len(t, logsData, 0)
			},
		},
		"set reward address from enabled succeeds": {
			Caller:     allowlist.TestEnabledAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackSetRewardAddress(rewardAddress)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetRewardAddressGasCost + RewardAddressChangedEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, state contract.StateDB) {
				address, isFeeRecipients := GetStoredRewardAddress(state)
				require.Equal(t, rewardAddress, address)
				require.False(t, isFeeRecipients)

				logsTopics, logsData := state.GetLogData()
				assertRewardAddressChanged(t, logsTopics, logsData, allowlist.TestEnabledAddr, common.Address{}, rewardAddress)
			},
		},
		"set allow fee recipients from manager succeeds": {
			Caller:     allowlist.TestManagerAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackAllowFeeRecipients()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: AllowFeeRecipientsGasCost + FeeRecipientsAllowedEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, state contract.StateDB) {
				_, isFeeRecipients := GetStoredRewardAddress(state)
				require.True(t, isFeeRecipients)

				logsTopics, logsData := state.GetLogData()
				assertFeeRecipientsAllowed(t, logsTopics, logsData, allowlist.TestManagerAddr)
			},
		},
		"set reward address from manager succeeds": {
			Caller:     allowlist.TestManagerAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackSetRewardAddress(rewardAddress)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetRewardAddressGasCost + RewardAddressChangedEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, state contract.StateDB) {
				address, isFeeRecipients := GetStoredRewardAddress(state)
				require.Equal(t, rewardAddress, address)
				require.False(t, isFeeRecipients)

				logsTopics, logsData := state.GetLogData()
				assertRewardAddressChanged(t, logsTopics, logsData, allowlist.TestManagerAddr, common.Address{}, rewardAddress)
			},
		},
		"change reward address should not emit events pre-Durango": {
			Caller:     allowlist.TestManagerAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackSetRewardAddress(rewardAddress)
				require.NoError(t, err)

				return input
			},
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				mockChainConfig := precompileconfig.NewMockChainConfig(ctrl)
				mockChainConfig.EXPECT().GetFeeConfig().AnyTimes().Return(commontype.ValidTestFeeConfig)
				mockChainConfig.EXPECT().AllowedFeeRecipients().AnyTimes().Return(false)
				mockChainConfig.EXPECT().IsDurango(gomock.Any()).AnyTimes().Return(false)
				return mockChainConfig
			},
			SuppliedGas: SetRewardAddressGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, stateDB contract.StateDB) {
				// Check no logs are stored in state
				logsTopics, logsData := stateDB.GetLogData()
				require.Len(t, logsTopics, 0)
				require.Len(t, logsData, 0)
			},
		},
		"disable rewards from manager succeeds": {
			Caller:     allowlist.TestManagerAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackDisableRewards()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: DisableRewardsGasCost + RewardsDisabledEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, state contract.StateDB) {
				address, isFeeRecipients := GetStoredRewardAddress(state)
				require.False(t, isFeeRecipients)
				require.Equal(t, constants.BlackholeAddr, address)

				logsTopics, logsData := state.GetLogData()
				assertRewardsDisabled(t, logsTopics, logsData, allowlist.TestManagerAddr)
			},
		},
		"disable rewards from enabled succeeds": {
			Caller:     allowlist.TestEnabledAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackDisableRewards()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: DisableRewardsGasCost + RewardsDisabledEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, state contract.StateDB) {
				address, isFeeRecipients := GetStoredRewardAddress(state)
				require.False(t, isFeeRecipients)
				require.Equal(t, constants.BlackholeAddr, address)

				logsTopics, logsData := state.GetLogData()
				assertRewardsDisabled(t, logsTopics, logsData, allowlist.TestEnabledAddr)
			},
		},
		"disable rewards should not emit event pre-Durango": {
			Caller:     allowlist.TestManagerAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackDisableRewards()
				require.NoError(t, err)

				return input
			},
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				mockChainConfig := precompileconfig.NewMockChainConfig(ctrl)
				mockChainConfig.EXPECT().GetFeeConfig().AnyTimes().Return(commontype.ValidTestFeeConfig)
				mockChainConfig.EXPECT().AllowedFeeRecipients().AnyTimes().Return(false)
				mockChainConfig.EXPECT().IsDurango(gomock.Any()).AnyTimes().Return(false)
				return mockChainConfig
			},
			SuppliedGas: SetRewardAddressGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, stateDB contract.StateDB) {
				// Check logs are not stored in state
				topics, data := stateDB.GetLogData()
				require.Len(t, topics, 0)
				require.Len(t, data, 0)
			},
		},
		"get current reward address from no role succeeds": {
			Caller: allowlist.TestNoRoleAddr,
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				allowlist.SetDefaultRoles(Module.Address)(t, state)
				StoreRewardAddress(state, rewardAddress)
			},
			InputFn: func(t testing.TB) []byte {
				input, err := PackCurrentRewardAddress()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: CurrentRewardAddressGasCost,
			ReadOnly:    false,
			ExpectedRes: func() []byte {
				res, err := PackCurrentRewardAddressOutput(rewardAddress)
				if err != nil {
					panic(err)
				}
				return res
			}(),
		},
		"get are fee recipients allowed from no role succeeds": {
			Caller: allowlist.TestNoRoleAddr,
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				allowlist.SetDefaultRoles(Module.Address)(t, state)
				EnableAllowFeeRecipients(state)
			},
			InputFn: func(t testing.TB) []byte {
				input, err := PackAreFeeRecipientsAllowed()
				require.NoError(t, err)
				return input
			},
			SuppliedGas: AreFeeRecipientsAllowedGasCost,
			ReadOnly:    false,
			ExpectedRes: func() []byte {
				res, err := PackAreFeeRecipientsAllowedOutput(true)
				if err != nil {
					panic(err)
				}
				return res
			}(),
		},
		"get initial config with address": {
			Caller:     allowlist.TestNoRoleAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackCurrentRewardAddress()
				require.NoError(t, err)
				return input
			},
			SuppliedGas: CurrentRewardAddressGasCost,
			Config: &Config{
				InitialRewardConfig: &InitialRewardConfig{
					RewardAddress: rewardAddress,
				},
			},
			ReadOnly: false,
			ExpectedRes: func() []byte {
				res, err := PackCurrentRewardAddressOutput(rewardAddress)
				if err != nil {
					panic(err)
				}
				return res
			}(),
		},
		"get initial config with allow fee recipients enabled": {
			Caller:     allowlist.TestNoRoleAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
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
				if err != nil {
					panic(err)
				}
				return res
			}(),
		},
		"readOnly allow fee recipients with allowed role fails": {
			Caller:     allowlist.TestEnabledAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackAllowFeeRecipients()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: AllowFeeRecipientsGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly set reward addresss with allowed role fails": {
			Caller:     allowlist.TestEnabledAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackSetRewardAddress(rewardAddress)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetRewardAddressGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"insufficient gas set reward address from allowed role": {
			Caller:     allowlist.TestEnabledAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackSetRewardAddress(rewardAddress)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetRewardAddressGasCost + RewardAddressChangedEventGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas allow fee recipients from allowed role": {
			Caller:     allowlist.TestEnabledAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackAllowFeeRecipients()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: AllowFeeRecipientsGasCost + FeeRecipientsAllowedEventGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas read current reward address from allowed role": {
			Caller:     allowlist.TestEnabledAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackCurrentRewardAddress()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: CurrentRewardAddressGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"insufficient gas are fee recipients allowed from allowed role": {
			Caller:     allowlist.TestEnabledAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackAreFeeRecipientsAllowed()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: AreFeeRecipientsAllowedGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
	}
)

func TestRewardManagerRun(t *testing.T) {
	allowlist.RunPrecompileWithAllowListTests(t, Module, state.NewTestStateDB, tests)
}

func BenchmarkRewardManager(b *testing.B) {
	allowlist.BenchPrecompileWithAllowList(b, Module, state.NewTestStateDB, tests)
}

func assertRewardAddressChanged(
	t testing.TB,
	logsTopics [][]common.Hash,
	logsData [][]byte,
	caller,
	oldAddress,
	newAddress common.Address) {
	require.Len(t, logsTopics, 1)
	require.Len(t, logsData, 1)
	topics := logsTopics[0]
	require.Len(t, topics, 4)
	require.Equal(t, RewardManagerABI.Events["RewardAddressChanged"].ID, topics[0])
	require.Equal(t, caller.Hash(), topics[1])
	require.Equal(t, oldAddress.Hash(), topics[2])
	require.Equal(t, newAddress.Hash(), topics[3])
	require.Len(t, logsData[0], 0)
}

func assertRewardsDisabled(
	t testing.TB,
	logsTopics [][]common.Hash,
	logsData [][]byte,
	caller common.Address) {
	require.Len(t, logsTopics, 1)
	require.Len(t, logsData, 1)
	topics := logsTopics[0]
	require.Len(t, topics, 2)
	require.Equal(t, RewardManagerABI.Events["RewardsDisabled"].ID, topics[0])
	require.Equal(t, caller.Hash(), topics[1])
	require.Len(t, logsData[0], 0)
}

func assertFeeRecipientsAllowed(
	t testing.TB,
	logsTopics [][]common.Hash,
	logsData [][]byte,
	caller common.Address) {
	require.Len(t, logsTopics, 1)
	require.Len(t, logsData, 1)
	topics := logsTopics[0]
	require.Len(t, topics, 2)
	require.Equal(t, RewardManagerABI.Events["FeeRecipientsAllowed"].ID, topics[0])
	require.Equal(t, caller.Hash(), topics[1])
	require.Len(t, logsData[0], 0)
}
