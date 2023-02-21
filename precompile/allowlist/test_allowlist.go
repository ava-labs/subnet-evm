// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package allowlist

import (
	"encoding/json"
	"testing"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var (
	TestAdminAddr   = common.HexToAddress("0x0000000000000000000000000000000000000011")
	TestEnabledAddr = common.HexToAddress("0x0000000000000000000000000000000000000022")
	TestNoRoleAddr  = common.HexToAddress("0x0000000000000000000000000000000000000033")
)

// mkConfigWithAllowList creates a new config with the correct type for [module]
// by marshalling [cfg] to JSON and then unmarshalling it into the config.
func mkConfigWithAllowList(module modules.Module, cfg *AllowListConfig) precompileconfig.Config {
	jsonBytes, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}

	moduleCfg := module.MakeConfig()
	err = json.Unmarshal(jsonBytes, moduleCfg)
	if err != nil {
		panic(err)
	}

	return moduleCfg
}

func AllowListTests(module modules.Module) map[string]testutils.PrecompileTest {
	contractAddress := module.Address
	return map[string]testutils.PrecompileTest{
		"set admin": {
			Caller:     TestAdminAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t *testing.T) []byte {
				input, err := PackModifyAllowList(TestNoRoleAddr, AdminRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				res := GetAllowListStatus(state, contractAddress, TestNoRoleAddr)
				require.Equal(t, AdminRole, res)
			},
		},
		"set enabled": {
			Caller:     TestAdminAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t *testing.T) []byte {
				input, err := PackModifyAllowList(TestNoRoleAddr, EnabledRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				res := GetAllowListStatus(state, contractAddress, TestNoRoleAddr)
				require.Equal(t, EnabledRole, res)
			},
		},
		"set no role": {
			Caller:     TestAdminAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t *testing.T) []byte {
				input, err := PackModifyAllowList(TestEnabledAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				res := GetAllowListStatus(state, contractAddress, TestEnabledAddr)
				require.Equal(t, NoRole, res)
			},
		},
		"set no role from no role": {
			Caller:     TestNoRoleAddr,
			BeforeHook: SetDefaultRoles(contractAddress),
			InputFn: func(t *testing.T) []byte {
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
			InputFn: func(t *testing.T) []byte {
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
			InputFn: func(t *testing.T) []byte {
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
			InputFn: func(t *testing.T) []byte {
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
			InputFn: func(t *testing.T) []byte {
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
			InputFn: func(t *testing.T) []byte {
				input, err := PackModifyAllowList(TestNoRoleAddr, AdminRole)
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
			InputFn: func(t *testing.T) []byte {
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
			InputFn: func(t *testing.T) []byte {
				input, err := PackModifyAllowList(TestEnabledAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"read allow list no role": {
			Caller:      TestNoRoleAddr,
			BeforeHook:  SetDefaultRoles(contractAddress),
			Input:       PackReadAllowList(TestNoRoleAddr),
			SuppliedGas: ReadAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: common.Hash(NoRole).Bytes(),
		},
		"read allow list admin role": {
			Caller:      TestAdminAddr,
			BeforeHook:  SetDefaultRoles(contractAddress),
			Input:       PackReadAllowList(TestAdminAddr),
			SuppliedGas: ReadAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: common.Hash(AdminRole).Bytes(),
		},
		"read allow list with readOnly enabled": {
			Caller:      TestAdminAddr,
			BeforeHook:  SetDefaultRoles(contractAddress),
			Input:       PackReadAllowList(TestNoRoleAddr),
			SuppliedGas: ReadAllowListGasCost,
			ReadOnly:    true,
			ExpectedRes: common.Hash(NoRole).Bytes(),
		},
		"read allow list out of gas": {
			Caller:      TestAdminAddr,
			BeforeHook:  SetDefaultRoles(contractAddress),
			Input:       PackReadAllowList(TestNoRoleAddr),
			SuppliedGas: ReadAllowListGasCost - 1,
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
			AfterHook: func(t *testing.T, state contract.StateDB) {
				require.Equal(t, AdminRole, GetAllowListStatus(state, contractAddress, TestNoRoleAddr))
				require.Equal(t, AdminRole, GetAllowListStatus(state, contractAddress, TestEnabledAddr))
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
			AfterHook: func(t *testing.T, state contract.StateDB) {
				require.Equal(t, EnabledRole, GetAllowListStatus(state, contractAddress, TestAdminAddr))
				require.Equal(t, EnabledRole, GetAllowListStatus(state, contractAddress, TestNoRoleAddr))
			},
		},
	}
}

// SetDefaultRoles returns a BeforeHook that sets roles TestAdminAddr and TestEnabledAddr
// to have the AdminRole and EnabledRole respectively.
func SetDefaultRoles(contractAddress common.Address) func(t *testing.T, state contract.StateDB) {
	return func(t *testing.T, state contract.StateDB) {
		SetAllowListRole(state, contractAddress, TestAdminAddr, AdminRole)
		SetAllowListRole(state, contractAddress, TestEnabledAddr, EnabledRole)
		require.Equal(t, AdminRole, GetAllowListStatus(state, contractAddress, TestAdminAddr))
		require.Equal(t, EnabledRole, GetAllowListStatus(state, contractAddress, TestEnabledAddr))
		require.Equal(t, NoRole, GetAllowListStatus(state, contractAddress, TestNoRoleAddr))
	}
}

func RunPrecompileWithAllowListTests(t *testing.T, module modules.Module, newStateDB func(t *testing.T) contract.StateDB, contractTests map[string]testutils.PrecompileTest) {
	t.Helper()
	tests := AllowListTests(module)
	// Add the contract specific tests to the map of tests to run.
	for name, test := range contractTests {
		if _, exists := tests[name]; exists {
			t.Fatalf("duplicate test name: %s", name)
		}
		tests[name] = test
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.Run(t, module, newStateDB(t))
		})
	}
}
