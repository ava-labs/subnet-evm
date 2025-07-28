// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package allowlisttest

import (
	"testing"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/vm"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/precompile/precompiletest"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	TestAdminAddr   = common.HexToAddress("0x0000000000000000000000000000000000000011")
	TestEnabledAddr = common.HexToAddress("0x0000000000000000000000000000000000000022")
	TestNoRoleAddr  = common.HexToAddress("0x0000000000000000000000000000000000000033")
	TestManagerAddr = common.HexToAddress("0x0000000000000000000000000000000000000044")
)

func AllowListTests(t testing.TB, module modules.Module) map[string]precompiletest.PrecompileTest {
	contractAddress := module.Address
	return map[string]precompiletest.PrecompileTest{
		"admin set admin": {
			Caller: TestAdminAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestNoRoleAddr, allowlist.AdminRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost + allowlist.AllowListEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, state contract.StateDB) {
				res := allowlist.GetAllowListStatus(state, contractAddress, TestNoRoleAddr)
				require.Equal(t, allowlist.AdminRole, res)
				// Check logs are stored in state
				logsTopics, logsData := state.GetLogData()
				assertSetRoleEvent(t, logsTopics, logsData, allowlist.AdminRole, TestNoRoleAddr, TestAdminAddr, allowlist.NoRole)
			},
		},
		"admin set enabled": {
			Caller: TestAdminAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestNoRoleAddr, allowlist.EnabledRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost + allowlist.AllowListEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, state contract.StateDB) {
				res := allowlist.GetAllowListStatus(state, contractAddress, TestNoRoleAddr)
				require.Equal(t, allowlist.EnabledRole, res)
				// Check logs are stored in state
				logsTopics, logsData := state.GetLogData()
				assertSetRoleEvent(t, logsTopics, logsData, allowlist.EnabledRole, TestNoRoleAddr, TestAdminAddr, allowlist.NoRole)
			},
		},
		"admin set no role": {
			Caller: TestAdminAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestEnabledAddr, allowlist.NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost + allowlist.AllowListEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, state contract.StateDB) {
				res := allowlist.GetAllowListStatus(state, contractAddress, TestEnabledAddr)
				require.Equal(t, allowlist.NoRole, res)
				// Check logs are stored in state
				logsTopics, logsData := state.GetLogData()
				assertSetRoleEvent(t, logsTopics, logsData, allowlist.NoRole, TestEnabledAddr, TestAdminAddr, allowlist.EnabledRole)
			},
		},
		"no role set no role": {
			Caller: TestNoRoleAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestEnabledAddr, allowlist.NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"no role set enabled": {
			Caller: TestNoRoleAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestNoRoleAddr, allowlist.EnabledRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"no role set admin": {
			Caller: TestNoRoleAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestEnabledAddr, allowlist.AdminRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"enabled set no role": {
			Caller: TestEnabledAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestAdminAddr, allowlist.NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"enabled set enabled": {
			Caller: TestEnabledAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestNoRoleAddr, allowlist.EnabledRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"enabled set admin": {
			Caller: TestEnabledAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestNoRoleAddr, allowlist.AdminRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"no role set manager pre-Durango": {
			Caller: TestNoRoleAddr,
			Config: DefaultAllowListConfig(module),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDurango(gomock.Any()).Return(false).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestNoRoleAddr, allowlist.ManagerRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: 0,
			ReadOnly:    false,
			ExpectedErr: "invalid non-activated function selector",
		},
		"no role set manager": {
			Caller: TestNoRoleAddr,
			Config: DefaultAllowListConfig(module),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDurango(gomock.Any()).Return(true).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestNoRoleAddr, allowlist.ManagerRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"enabled role set manager pre-Durango": {
			Caller: TestEnabledAddr,
			Config: DefaultAllowListConfig(module),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDurango(gomock.Any()).Return(false).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestNoRoleAddr, allowlist.ManagerRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: 0,
			ReadOnly:    false,
			ExpectedErr: "invalid non-activated function selector",
		},
		"enabled set manager": {
			Caller: TestNoRoleAddr,
			Config: DefaultAllowListConfig(module),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDurango(gomock.Any()).Return(true).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestNoRoleAddr, allowlist.ManagerRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"admin set manager pre-DUpgarde": {
			Caller: TestAdminAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestNoRoleAddr, allowlist.ManagerRole)
				require.NoError(t, err)

				return input
			},
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDurango(gomock.Any()).Return(false).AnyTimes()
				return config
			},
			SuppliedGas: 0,
			ReadOnly:    false,
			ExpectedErr: "invalid non-activated function selector",
		},
		"admin set manager": {
			Caller: TestAdminAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestNoRoleAddr, allowlist.ManagerRole)
				require.NoError(t, err)

				return input
			},
			ExpectedRes: []byte{},
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDurango(gomock.Any()).Return(true).AnyTimes()
				return config
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost + allowlist.AllowListEventGasCost,
			ReadOnly:    false,
			AfterHook: func(t testing.TB, state contract.StateDB) {
				res := allowlist.GetAllowListStatus(state, contractAddress, TestNoRoleAddr)
				require.Equal(t, allowlist.ManagerRole, res)
				// Check logs are stored in state
				logsTopics, logsData := state.GetLogData()
				assertSetRoleEvent(t, logsTopics, logsData, allowlist.ManagerRole, TestNoRoleAddr, TestAdminAddr, allowlist.NoRole)
			},
		},
		"manager set no role to no role": {
			Caller: TestManagerAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestNoRoleAddr, allowlist.NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost + allowlist.AllowListEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			ExpectedErr: "",
			AfterHook: func(t testing.TB, state contract.StateDB) {
				res := allowlist.GetAllowListStatus(state, contractAddress, TestNoRoleAddr)
				require.Equal(t, allowlist.NoRole, res)
				// Check logs are stored in state
				logsTopics, logsData := state.GetLogData()
				assertSetRoleEvent(t, logsTopics, logsData, allowlist.NoRole, TestNoRoleAddr, TestManagerAddr, allowlist.NoRole)
			},
		},
		"manager set no role to enabled": {
			Caller: TestManagerAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestNoRoleAddr, allowlist.EnabledRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost + allowlist.AllowListEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			ExpectedErr: "",
			AfterHook: func(t testing.TB, state contract.StateDB) {
				res := allowlist.GetAllowListStatus(state, contractAddress, TestNoRoleAddr)
				require.Equal(t, allowlist.EnabledRole, res)

				// Check logs are stored in state
				logsTopics, logsData := state.GetLogData()
				assertSetRoleEvent(t, logsTopics, logsData, allowlist.EnabledRole, TestNoRoleAddr, TestManagerAddr, allowlist.NoRole)
			},
		},
		"manager set no role to manager": {
			Caller: TestManagerAddr,
			Config: DefaultAllowListConfig(module),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDurango(gomock.Any()).Return(true).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestNoRoleAddr, allowlist.ManagerRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"manager set no role to admin": {
			Caller: TestManagerAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestNoRoleAddr, allowlist.AdminRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"manager set enabled to admin": {
			Caller: TestManagerAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestEnabledAddr, allowlist.AdminRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"manager set enabled role to manager": {
			Caller: TestManagerAddr,
			Config: DefaultAllowListConfig(module),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDurango(gomock.Any()).Return(true).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestEnabledAddr, allowlist.ManagerRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"manager set enabled role to no role": {
			Caller: TestManagerAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestEnabledAddr, allowlist.NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost + allowlist.AllowListEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, state contract.StateDB) {
				res := allowlist.GetAllowListStatus(state, contractAddress, TestNoRoleAddr)
				require.Equal(t, allowlist.NoRole, res)

				// Check logs are stored in state
				logsTopics, logsData := state.GetLogData()
				assertSetRoleEvent(t, logsTopics, logsData, allowlist.NoRole, TestEnabledAddr, TestManagerAddr, allowlist.EnabledRole)
			},
		},
		"manager set admin to no role": {
			Caller: TestManagerAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestAdminAddr, allowlist.NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"manager set admin role to enabled": {
			Caller: TestManagerAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestAdminAddr, allowlist.EnabledRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"manager set admin to manager": {
			Caller: TestManagerAddr,
			Config: DefaultAllowListConfig(module),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDurango(gomock.Any()).Return(true).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestAdminAddr, allowlist.ManagerRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"manager set manager to no role": {
			Caller: TestManagerAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestManagerAddr, allowlist.NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"admin set no role with readOnly enabled": {
			Caller: TestAdminAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestEnabledAddr, allowlist.NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    true,
			ExpectedErr: vm.ErrWriteProtection.Error(),
		},
		"admin set no role insufficient gas": {
			Caller: TestAdminAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestEnabledAddr, allowlist.NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vm.ErrOutOfGas.Error(),
		},
		"no role read allow list": {
			Caller: TestNoRoleAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackReadAllowList(TestNoRoleAddr)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: allowlist.ReadAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: common.Hash(allowlist.NoRole).Bytes(),
		},
		"admin role read allow list": {
			Caller: TestAdminAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackReadAllowList(TestAdminAddr)
				require.NoError(t, err)

				return input
			}, SuppliedGas: allowlist.ReadAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: common.Hash(allowlist.AdminRole).Bytes(),
		},
		"admin read allow list with readOnly enabled": {
			Caller: TestAdminAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackReadAllowList(TestNoRoleAddr)
				require.NoError(t, err)

				return input
			}, SuppliedGas: allowlist.ReadAllowListGasCost,
			ReadOnly:    true,
			ExpectedRes: common.Hash(allowlist.NoRole).Bytes(),
		},
		"radmin read allow list out of gas": {
			Caller: TestAdminAddr,
			Config: DefaultAllowListConfig(module),
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackReadAllowList(TestNoRoleAddr)
				require.NoError(t, err)

				return input
			}, SuppliedGas: allowlist.ReadAllowListGasCost - 1,
			ReadOnly:    true,
			ExpectedErr: vm.ErrOutOfGas.Error(),
		},
		"initial config sets admins": {
			Config: mkConfigWithAllowList(
				module,
				&allowlist.AllowListConfig{
					AdminAddresses: []common.Address{TestNoRoleAddr, TestEnabledAddr},
				},
			),
			SuppliedGas: 0,
			ReadOnly:    false,
			AfterHook: func(t testing.TB, state contract.StateDB) {
				require.Equal(t, allowlist.AdminRole, allowlist.GetAllowListStatus(state, contractAddress, TestNoRoleAddr))
				require.Equal(t, allowlist.AdminRole, allowlist.GetAllowListStatus(state, contractAddress, TestEnabledAddr))
			},
		},
		"initial config sets managers": {
			Config: mkConfigWithAllowList(
				module,
				&allowlist.AllowListConfig{
					ManagerAddresses: []common.Address{TestNoRoleAddr, TestEnabledAddr},
				},
			),
			SuppliedGas: 0,
			ReadOnly:    false,
			AfterHook: func(t testing.TB, state contract.StateDB) {
				require.Equal(t, allowlist.ManagerRole, allowlist.GetAllowListStatus(state, contractAddress, TestNoRoleAddr))
				require.Equal(t, allowlist.ManagerRole, allowlist.GetAllowListStatus(state, contractAddress, TestEnabledAddr))
			},
		},
		"initial config sets enabled": {
			Config: mkConfigWithAllowList(
				module,
				&allowlist.AllowListConfig{
					EnabledAddresses: []common.Address{TestNoRoleAddr, TestAdminAddr},
				},
			),
			SuppliedGas: 0,
			ReadOnly:    false,
			AfterHook: func(t testing.TB, state contract.StateDB) {
				require.Equal(t, allowlist.EnabledRole, allowlist.GetAllowListStatus(state, contractAddress, TestAdminAddr))
				require.Equal(t, allowlist.EnabledRole, allowlist.GetAllowListStatus(state, contractAddress, TestNoRoleAddr))
			},
		},
		"admin set admin pre-Durango": {
			Caller: TestAdminAddr,
			Config: DefaultAllowListConfig(module),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDurango(gomock.Any()).Return(false).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestNoRoleAddr, allowlist.AdminRole)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, stateDB contract.StateDB) {
				// Check no logs are stored in state
				topics, data := stateDB.GetLogData()
				require.Len(t, topics, 0)
				require.Len(t, data, 0)
			},
		},
		"admin set enabled pre-Durango": {
			Caller: TestAdminAddr,
			Config: DefaultAllowListConfig(module),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDurango(gomock.Any()).Return(false).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestNoRoleAddr, allowlist.EnabledRole)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, stateDB contract.StateDB) {
				// Check no logs are stored in state
				topics, data := stateDB.GetLogData()
				require.Len(t, topics, 0)
				require.Len(t, data, 0)
			},
		},
		"admin set no role pre-Durango": {
			Caller: TestAdminAddr,
			Config: DefaultAllowListConfig(module),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDurango(gomock.Any()).Return(false).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := allowlist.PackModifyAllowList(TestEnabledAddr, allowlist.NoRole)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, stateDB contract.StateDB) {
				// Check no logs are stored in state
				topics, data := stateDB.GetLogData()
				require.Len(t, topics, 0)
				require.Len(t, data, 0)
			},
		},
	}
}

// DefaultAllowListConfig returns the default allowlist configuration with Admin, Enabled, and Manager roles
func DefaultAllowListConfig(module modules.Module) precompileconfig.Config {
	return mkConfigWithAllowList(
		module,
		&allowlist.AllowListConfig{
			AdminAddresses:   []common.Address{TestAdminAddr},
			EnabledAddresses: []common.Address{TestEnabledAddr},
			ManagerAddresses: []common.Address{TestManagerAddr},
		},
	)
}

func RunPrecompileWithAllowListTests(t *testing.T, module modules.Module, contractTests map[string]precompiletest.PrecompileTest) {
	t.Helper()
	tests := AllowListTests(t, module)
	// Add the contract specific tests to the map of tests to run.
	for name, test := range contractTests {
		if _, exists := tests[name]; exists {
			t.Fatalf("duplicate test name: %s", name)
		}
		tests[name] = test
	}

	precompiletest.RunPrecompileTests(t, module, tests)
}

func BenchPrecompileWithAllowList(b *testing.B, module modules.Module, contractTests map[string]precompiletest.PrecompileTest) {
	b.Helper()

	tests := AllowListTests(b, module)
	// Add the contract specific tests to the map of tests to run.
	for name, test := range contractTests {
		if _, exists := tests[name]; exists {
			b.Fatalf("duplicate bench name: %s", name)
		}
		tests[name] = test
	}

	for name, test := range tests {
		b.Run(name, func(b *testing.B) {
			test.Bench(b, module)
		})
	}
}

func assertSetRoleEvent(t testing.TB, logsTopics [][]common.Hash, logsData [][]byte, role allowlist.Role, addr common.Address, caller common.Address, oldRole allowlist.Role) {
	require.Len(t, logsTopics, 1)
	require.Len(t, logsData, 1)
	topics := logsTopics[0]
	require.Len(t, topics, 4)
	require.Equal(t, allowlist.AllowListABI.Events["RoleSet"].ID, topics[0])
	require.Equal(t, role.Hash(), topics[1])
	require.Equal(t, common.BytesToHash(addr[:]), topics[2])
	require.Equal(t, common.BytesToHash(caller[:]), topics[3])
	data := logsData[0]
	require.Equal(t, oldRole.Bytes(), data)
}
