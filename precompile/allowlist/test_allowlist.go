// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package allowlist

import (
	"testing"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	TestAdminAddr   = common.HexToAddress("0x0000000000000000000000000000000000000011")
	TestEnabledAddr = common.HexToAddress("0x0000000000000000000000000000000000000022")
	TestNoRoleAddr  = common.HexToAddress("0x0000000000000000000000000000000000000033")
	TestManagerAddr = common.HexToAddress("0x0000000000000000000000000000000000000044")
)

func AllowListTests(t testing.TB, module modules.Module) map[string]testutils.PrecompileTest {
	contractAddress := module.Address
	return map[string]testutils.PrecompileTest{
		"admin set admin": {
			Caller:     TestAdminAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestNoRoleAddr, AdminRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost + AllowListEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, state contract.StateDB) {
				res := GetAllowListStatus(state, contractAddress, TestNoRoleAddr)
				require.Equal(t, AdminRole, res)
			},
		},
		"admin set enabled": {
			Caller:     TestAdminAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestNoRoleAddr, EnabledRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost + AllowListEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, state contract.StateDB) {
				res := GetAllowListStatus(state, contractAddress, TestNoRoleAddr)
				require.Equal(t, EnabledRole, res)
			},
		},
		"admin set no role": {
			Caller:     TestAdminAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestEnabledAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost + AllowListEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, state contract.StateDB) {
				res := GetAllowListStatus(state, contractAddress, TestEnabledAddr)
				require.Equal(t, NoRole, res)
			},
		},
		"set no role from no role": {
			Caller:     TestNoRoleAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestEnabledAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set enabled from no role": {
			Caller:     TestNoRoleAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestNoRoleAddr, EnabledRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set admin from no role": {
			Caller:     TestNoRoleAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestEnabledAddr, AdminRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set no role from enabled": {
			Caller:     TestEnabledAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestAdminAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set enabled from enabled": {
			Caller:     TestEnabledAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestNoRoleAddr, EnabledRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set admin from enabled": {
			Caller:     TestEnabledAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestNoRoleAddr, AdminRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set manager from no role before activation": {
			Caller:     TestNoRoleAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDUpgrade(gomock.Any()).Return(false).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestNoRoleAddr, ManagerRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: 0,
			ReadOnly:    false,
			ExpectedErr: "invalid non-activated function selector",
		},
		"set manager from no role after activation": {
			Caller:     TestNoRoleAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDUpgrade(gomock.Any()).Return(true).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestNoRoleAddr, ManagerRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set manager from enabled role before activation": {
			Caller:     TestEnabledAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDUpgrade(gomock.Any()).Return(false).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestNoRoleAddr, ManagerRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: 0,
			ReadOnly:    false,
			ExpectedErr: "invalid non-activated function selector",
		},
		"set manager from enabled after activation": {
			Caller:     TestNoRoleAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDUpgrade(gomock.Any()).Return(true).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestNoRoleAddr, ManagerRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set manager from admin before activation": {
			Caller:     TestAdminAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestNoRoleAddr, ManagerRole)
				require.NoError(t, err)

				return input
			},
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDUpgrade(gomock.Any()).Return(false).AnyTimes()
				return config
			},
			SuppliedGas: 0,
			ReadOnly:    false,
			ExpectedErr: "invalid non-activated function selector",
		},
		"set manager from admin after activation": {
			Caller:     TestAdminAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestNoRoleAddr, ManagerRole)
				require.NoError(t, err)

				return input
			},
			ExpectedRes: []byte{},
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDUpgrade(gomock.Any()).Return(true).AnyTimes()
				return config
			},
			SuppliedGas: ModifyAllowListGasCost + AllowListEventGasCost,
			ReadOnly:    false,
			AfterHook: func(t testing.TB, state contract.StateDB) {
				res := GetAllowListStatus(state, contractAddress, TestNoRoleAddr)
				require.Equal(t, ManagerRole, res)
			},
		},
		"set no role to no role from manager after activation": {
			Caller:     TestManagerAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestNoRoleAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost + AllowListEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			ExpectedErr: "",
			AfterHook: func(t testing.TB, state contract.StateDB) {
				res := GetAllowListStatus(state, contractAddress, TestNoRoleAddr)
				require.Equal(t, NoRole, res)
			},
		},
		"set no role to enabled from manager after activation": {
			Caller:     TestManagerAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestNoRoleAddr, EnabledRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost + AllowListEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			ExpectedErr: "",
			AfterHook: func(t testing.TB, state contract.StateDB) {
				res := GetAllowListStatus(state, contractAddress, TestNoRoleAddr)
				require.Equal(t, EnabledRole, res)
			},
		},
		"set no role to manager from manager after activation": {
			Caller:     TestManagerAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDUpgrade(gomock.Any()).Return(true).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestNoRoleAddr, ManagerRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set no role to admin from manager after activation": {
			Caller:     TestManagerAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestNoRoleAddr, AdminRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set enabled role to admin from manager after activation": {
			Caller:     TestManagerAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestEnabledAddr, AdminRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set enabled role to manager from manager after activation": {
			Caller:     TestManagerAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDUpgrade(gomock.Any()).Return(true).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestEnabledAddr, ManagerRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set enabled role to no role from manager after activation": {
			Caller:     TestManagerAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestEnabledAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost + AllowListEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, state contract.StateDB) {
				res := GetAllowListStatus(state, contractAddress, TestNoRoleAddr)
				require.Equal(t, NoRole, res)
			},
		},
		"set admin to no role from manager after activation": {
			Caller:     TestManagerAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestAdminAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set admin role to enabled from manager after activation": {
			Caller:     TestManagerAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestAdminAddr, EnabledRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set admin to manager from manager after activation": {
			Caller:     TestManagerAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDUpgrade(gomock.Any()).Return(true).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestAdminAddr, ManagerRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set manager role to no role from manager after activation": {
			Caller:     TestManagerAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestManagerAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set no role with readOnly enabled": {
			Caller:     TestAdminAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestEnabledAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"set no role insufficient gas": {
			Caller:     TestAdminAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestEnabledAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"read allow list no role": {
			Caller:     TestNoRoleAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackReadAllowList(TestNoRoleAddr)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ReadAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: common.Hash(NoRole).Bytes(),
		},
		"read allow list admin role": {
			Caller:     TestAdminAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackReadAllowList(TestAdminAddr)
				require.NoError(t, err)

				return input
			}, SuppliedGas: ReadAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: common.Hash(AdminRole).Bytes(),
		},
		"read allow list with readOnly enabled": {
			Caller:     TestAdminAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackReadAllowList(TestNoRoleAddr)
				require.NoError(t, err)

				return input
			}, SuppliedGas: ReadAllowListGasCost,
			ReadOnly:    true,
			ExpectedRes: common.Hash(NoRole).Bytes(),
		},
		"read allow list out of gas": {
			Caller:     TestAdminAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t testing.TB) []byte {
				input, err := PackReadAllowList(TestNoRoleAddr)
				require.NoError(t, err)

				return input
			}, SuppliedGas: ReadAllowListGasCost - 1,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"initial config sets admins": {
			Config: mkConfigWithAllowList(
				module,
				&AllowListConfig{
					AdminAddresses: []common.Address{TestNoRoleAddr, TestEnabledAddr},
				},
			),
			SuppliedGas: 0,
			ReadOnly:    false,
			AfterHook: func(t testing.TB, state contract.StateDB) {
				require.Equal(t, AdminRole, GetAllowListStatus(state, contractAddress, TestNoRoleAddr))
				require.Equal(t, AdminRole, GetAllowListStatus(state, contractAddress, TestEnabledAddr))
			},
		},
		"initial config sets managers": {
			Config: mkConfigWithAllowList(
				module,
				&AllowListConfig{
					ManagerAddresses: []common.Address{TestNoRoleAddr, TestEnabledAddr},
				},
			),
			SuppliedGas: 0,
			ReadOnly:    false,
			AfterHook: func(t testing.TB, state contract.StateDB) {
				require.Equal(t, ManagerRole, GetAllowListStatus(state, contractAddress, TestNoRoleAddr))
				require.Equal(t, ManagerRole, GetAllowListStatus(state, contractAddress, TestEnabledAddr))
			},
		},
		"initial config sets enabled": {
			Config: mkConfigWithAllowList(
				module,
				&AllowListConfig{
					EnabledAddresses: []common.Address{TestNoRoleAddr, TestAdminAddr},
				},
			),
			SuppliedGas: 0,
			ReadOnly:    false,
			AfterHook: func(t testing.TB, state contract.StateDB) {
				require.Equal(t, EnabledRole, GetAllowListStatus(state, contractAddress, TestAdminAddr))
				require.Equal(t, EnabledRole, GetAllowListStatus(state, contractAddress, TestNoRoleAddr))
			},
		},
		"allowList admin event doesn't log if D fork is not active": {
			Caller:     TestAdminAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDUpgrade(gomock.Any()).Return(false).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestEnabledAddr, AdminRole)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, stateDB contract.StateDB) {
				// Check no logs are stored in state
				topics, data := stateDB.GetLogData()
				require.Len(t, topics, 0)
				require.Len(t, data, 0)
			},
		},
		"allowList enabled event doesn't log if D fork is not active": {
			Caller:     TestAdminAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDUpgrade(gomock.Any()).Return(false).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestNoRoleAddr, EnabledRole)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, stateDB contract.StateDB) {
				// Check no logs are stored in state
				topics, data := stateDB.GetLogData()
				require.Len(t, topics, 0)
				require.Len(t, data, 0)
			},
		},
		"allowList role remove event doesn't log if D fork is not active": {
			Caller:     TestAdminAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDUpgrade(gomock.Any()).Return(false).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestEnabledAddr, NoRole)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, stateDB contract.StateDB) {
				// Check no logs are stored in state
				topics, data := stateDB.GetLogData()
				require.Len(t, topics, 0)
				require.Len(t, data, 0)
			},
		},
		"allowList admin event does log if D fork is active": {
			Caller:     TestAdminAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDUpgrade(gomock.Any()).Return(true).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestEnabledAddr, AdminRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost + AllowListEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, stateDB contract.StateDB) {
				logsTopics, logsData := stateDB.GetLogData()
				require.Len(t, logsTopics, 1)
				require.Len(t, logsData, 1)
				topics := logsTopics[0]
				require.Len(t, topics, 3)
				require.Equal(t, topics[0], AllowListABI.Events["AdminAdded"].ID)
				require.Equal(t, topics[1], TestAdminAddr.Hash())
				require.Equal(t, topics[2], TestEnabledAddr.Hash())
				require.Len(t, logsData[0], 0)
			},
		},
		"allowList manager event does log if D fork is active": {
			Caller:     TestAdminAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDUpgrade(gomock.Any()).Return(true).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestEnabledAddr, ManagerRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost + AllowListEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, stateDB contract.StateDB) {
				logsTopics, logsData := stateDB.GetLogData()
				require.Len(t, logsTopics, 1)
				require.Len(t, logsData, 1)
				topics := logsTopics[0]
				require.Len(t, topics, 3)
				require.Equal(t, topics[0], AllowListABI.Events["ManagerAdded"].ID)
				require.Equal(t, topics[1], TestAdminAddr.Hash())
				require.Equal(t, topics[2], TestEnabledAddr.Hash())
				require.Len(t, logsData[0], 0)
			},
		},
		"allowList enabled event does log if D fork is active": {
			Caller:     TestAdminAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDUpgrade(gomock.Any()).Return(true).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestNoRoleAddr, EnabledRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost + AllowListEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, stateDB contract.StateDB) {
				logsTopics, logsData := stateDB.GetLogData()
				require.Len(t, logsTopics, 1)
				require.Len(t, logsData, 1)
				topics := logsTopics[0]
				require.Len(t, topics, 3)
				require.Equal(t, topics[0], AllowListABI.Events["EnabledAdded"].ID)
				require.Equal(t, topics[1], TestAdminAddr.Hash())
				require.Equal(t, topics[2], TestNoRoleAddr.Hash())
				require.Len(t, logsData[0], 0)
			},
		},
		"allowList role remove event does log if D fork is active": {
			Caller:     TestAdminAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			ChainConfigFn: func(ctrl *gomock.Controller) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(ctrl)
				config.EXPECT().IsDUpgrade(gomock.Any()).Return(true).AnyTimes()
				return config
			},
			InputFn: func(t testing.TB) []byte {
				input, err := PackModifyAllowList(TestEnabledAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost + AllowListEventGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, stateDB contract.StateDB) {
				logsTopics, logsData := stateDB.GetLogData()
				require.Len(t, logsTopics, 1)
				require.Len(t, logsData, 1)
				topics := logsTopics[0]
				require.Len(t, topics, 3)
				require.Equal(t, topics[0], AllowListABI.Events["RoleRemoved"].ID)
				require.Equal(t, topics[1], TestAdminAddr.Hash())
				require.Equal(t, topics[2], TestEnabledAddr.Hash())
				require.Len(t, logsData[0], 0)
			},
		},
	}
}

// SetDefaultRoles returns a BeforeHook that sets roles TestAdminAddr and TestEnabledAddr
// to have the AdminRole and EnabledRole respectively.
func SetDefaultRoles(contractAddress common.Address) func(t testing.TB, state contract.StateDB) {
	return func(t testing.TB, state contract.StateDB) {
		SetAllowListRole(state, contractAddress, TestAdminAddr, AdminRole)
		SetAllowListRole(state, contractAddress, TestManagerAddr, ManagerRole)
		SetAllowListRole(state, contractAddress, TestEnabledAddr, EnabledRole)
		require.Equal(t, AdminRole, GetAllowListStatus(state, contractAddress, TestAdminAddr))
		require.Equal(t, ManagerRole, GetAllowListStatus(state, contractAddress, TestManagerAddr))
		require.Equal(t, EnabledRole, GetAllowListStatus(state, contractAddress, TestEnabledAddr))
		require.Equal(t, NoRole, GetAllowListStatus(state, contractAddress, TestNoRoleAddr))
	}
}

func RunPrecompileWithAllowListTests(t *testing.T, module modules.Module, newStateDB func(t testing.TB) contract.StateDB, contractTests map[string]testutils.PrecompileTest) {
	t.Helper()
	tests := AllowListTests(t, module)
	// Add the contract specific tests to the map of tests to run.
	for name, test := range contractTests {
		if _, exists := tests[name]; exists {
			t.Fatalf("duplicate test name: %s", name)
		}
		tests[name] = test
	}

	testutils.RunPrecompileTests(t, module, newStateDB, tests)
}

func BenchPrecompileWithAllowList(b *testing.B, module modules.Module, newStateDB func(t testing.TB) contract.StateDB, contractTests map[string]testutils.PrecompileTest) {
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
			test.Bench(b, module, newStateDB(b))
		})
	}
}
