// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nativeminter

import (
	"math/big"
	"testing"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/common/math"
	"github.com/ava-labs/libevm/core/vm"
	"github.com/ava-labs/subnet-evm/core/extstate/extstatetest"
	"github.com/ava-labs/subnet-evm/precompile/allowlist/allowlisttest"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/precompile/precompiletest"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	tests = map[string]precompiletest.PrecompileTest{
		"calling mintNativeCoin from NoRole should fail": {
			Caller:     allowlisttest.TestNoRoleAddr,
			BeforeHook: allowlisttest.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackMintNativeCoin(allowlisttest.TestNoRoleAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: MintGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotMint.Error(),
		},
		"calling mintNativeCoin from Enabled should succeed": {
			Caller:     allowlisttest.TestEnabledAddr,
			BeforeHook: allowlisttest.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackMintNativeCoin(allowlisttest.TestEnabledAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: MintGasCost + NativeCoinMintedEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, stateDB contract.StateDB) {
				expected := uint256.MustFromBig(common.Big1)
				require.Equal(t, expected, stateDB.GetBalance(allowlisttest.TestEnabledAddr), "expected minted funds")

				logsTopics, logsData := stateDB.GetLogData()
				assertNativeCoinMintedEvent(t, logsTopics, logsData, allowlisttest.TestEnabledAddr, allowlisttest.TestEnabledAddr, common.Big1)
			},
		},
		"initial mint funds": {
			Caller:     allowlisttest.TestEnabledAddr,
			BeforeHook: allowlisttest.SetDefaultRoles(Module.Address),
			Config: &Config{
				InitialMint: map[common.Address]*math.HexOrDecimal256{
					allowlisttest.TestEnabledAddr: math.NewHexOrDecimal256(2),
				},
			},
			AfterHook: func(t testing.TB, stateDB contract.StateDB) {
				expected := uint256.MustFromBig(common.Big2)
				require.Equal(t, expected, stateDB.GetBalance(allowlisttest.TestEnabledAddr), "expected minted funds")
			},
		},
		"calling mintNativeCoin from Manager should succeed": {
			Caller:     allowlisttest.TestManagerAddr,
			BeforeHook: allowlisttest.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackMintNativeCoin(allowlisttest.TestEnabledAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: MintGasCost + NativeCoinMintedEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, stateDB contract.StateDB) {
				expected := uint256.MustFromBig(common.Big1)
				require.Equal(t, expected, stateDB.GetBalance(allowlisttest.TestEnabledAddr), "expected minted funds")

				logsTopics, logsData := stateDB.GetLogData()
				assertNativeCoinMintedEvent(t, logsTopics, logsData, allowlisttest.TestManagerAddr, allowlisttest.TestEnabledAddr, common.Big1)
			},
		},
		"mint funds from admin address": {
			Caller:     allowlisttest.TestAdminAddr,
			BeforeHook: allowlisttest.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackMintNativeCoin(allowlisttest.TestAdminAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: MintGasCost + NativeCoinMintedEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, stateDB contract.StateDB) {
				expected := uint256.MustFromBig(common.Big1)
				require.Equal(t, expected, stateDB.GetBalance(allowlisttest.TestAdminAddr), "expected minted funds")

				logsTopics, logsData := stateDB.GetLogData()
				assertNativeCoinMintedEvent(t, logsTopics, logsData, allowlisttest.TestAdminAddr, allowlisttest.TestAdminAddr, common.Big1)
			},
		},
		"mint max big funds": {
			Caller:     allowlisttest.TestAdminAddr,
			BeforeHook: allowlisttest.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackMintNativeCoin(allowlisttest.TestAdminAddr, math.MaxBig256)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: MintGasCost + NativeCoinMintedEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, stateDB contract.StateDB) {
				expected := uint256.MustFromBig(math.MaxBig256)
				require.Equal(t, expected, stateDB.GetBalance(allowlisttest.TestAdminAddr), "expected minted funds")

				logsTopics, logsData := stateDB.GetLogData()
				assertNativeCoinMintedEvent(t, logsTopics, logsData, allowlisttest.TestAdminAddr, allowlisttest.TestAdminAddr, math.MaxBig256)
			},
		},
		"readOnly mint with noRole fails": {
			Caller:     allowlisttest.TestNoRoleAddr,
			BeforeHook: allowlisttest.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackMintNativeCoin(allowlisttest.TestAdminAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: MintGasCost,
			ReadOnly:    true,
			ExpectedErr: vm.ErrWriteProtection.Error(),
		},
		"readOnly mint with allow role fails": {
			Caller:     allowlisttest.TestEnabledAddr,
			BeforeHook: allowlisttest.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackMintNativeCoin(allowlisttest.TestEnabledAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: MintGasCost,
			ReadOnly:    true,
			ExpectedErr: vm.ErrWriteProtection.Error(),
		},
		"readOnly mint with admin role fails": {
			Caller:     allowlisttest.TestAdminAddr,
			BeforeHook: allowlisttest.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackMintNativeCoin(allowlisttest.TestAdminAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: MintGasCost,
			ReadOnly:    true,
			ExpectedErr: vm.ErrWriteProtection.Error(),
		},
		"insufficient gas mint from admin": {
			Caller:     allowlisttest.TestAdminAddr,
			BeforeHook: allowlisttest.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackMintNativeCoin(allowlisttest.TestEnabledAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: MintGasCost + NativeCoinMintedEventGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vm.ErrOutOfGas.Error(),
		},
		"mint doesn't log pre-Durango": {
			Caller:     allowlisttest.TestEnabledAddr,
			BeforeHook: allowlisttest.SetDefaultRoles(Module.Address),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDurango(gomock.Any()).Return(false).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := PackMintNativeCoin(allowlisttest.TestEnabledAddr, common.Big1)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: MintGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, stateDB contract.StateDB) {
				// Check no logs are stored in state
				logsTopics, logsData := stateDB.GetLogData()
				require.Len(t, logsTopics, 0)
				require.Len(t, logsData, 0)
			},
		},
		"mint with extra padded bytes should fail pre-Durango": {
			Caller:     allowlisttest.TestEnabledAddr,
			BeforeHook: allowlisttest.SetDefaultRoles(Module.Address),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDurango(gomock.Any()).Return(false).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := PackMintNativeCoin(allowlisttest.TestEnabledAddr, common.Big1)
				require.NoError(t, err)

				// Add extra bytes to the end of the input
				input = append(input, make([]byte, 32)...)

				return input
			},
			SuppliedGas: MintGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrInvalidLen.Error(),
		},
		"mint with extra padded bytes should succeed with Durango": {
			Caller:     allowlisttest.TestEnabledAddr,
			BeforeHook: allowlisttest.SetDefaultRoles(Module.Address),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDurango(gomock.Any()).Return(true).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := PackMintNativeCoin(allowlisttest.TestEnabledAddr, common.Big1)
				require.NoError(t, err)

				// Add extra bytes to the end of the input
				input = append(input, make([]byte, 32)...)

				return input
			},
			ExpectedRes: []byte{},
			SuppliedGas: MintGasCost + NativeCoinMintedEventGasCost,
			ReadOnly:    false,
			AfterHook: func(t testing.TB, state contract.StateDB) {
				expected := uint256.MustFromBig(common.Big1)
				require.Equal(t, expected, state.GetBalance(allowlisttest.TestEnabledAddr), "expected minted funds")

				logsTopics, logsData := state.GetLogData()
				assertNativeCoinMintedEvent(t, logsTopics, logsData, allowlisttest.TestEnabledAddr, allowlisttest.TestEnabledAddr, common.Big1)
			},
		},
	}
)

func TestContractNativeMinterRun(t *testing.T) {
	allowlisttest.RunPrecompileWithAllowListTests(t, Module, extstatetest.NewTestStateDB, tests)
}

func BenchmarkContractNativeMinter(b *testing.B) {
	allowlisttest.BenchPrecompileWithAllowList(b, Module, extstatetest.NewTestStateDB, tests)
}

func assertNativeCoinMintedEvent(t testing.TB,
	logsTopics [][]common.Hash,
	logsData [][]byte,
	expectedSender common.Address,
	expectedRecipient common.Address,
	expectedAmount *big.Int,
) {
	require.Len(t, logsTopics, 1)
	require.Len(t, logsData, 1)
	topics := logsTopics[0]
	require.Len(t, topics, 3)
	require.Equal(t, NativeMinterABI.Events["NativeCoinMinted"].ID, topics[0])
	require.Equal(t, common.BytesToHash(expectedSender[:]), topics[1])
	require.Equal(t, common.BytesToHash(expectedRecipient[:]), topics[2])
	require.NotEmpty(t, logsData[0])
	amount, err := UnpackNativeCoinMintedEventData(logsData[0])
	require.NoError(t, err)
	require.True(t, expectedAmount.Cmp(amount) == 0, "expected", expectedAmount, "got", amount)
}
