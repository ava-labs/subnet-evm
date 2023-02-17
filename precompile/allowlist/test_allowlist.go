// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package allowlist

import (
	"testing"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/contracts/test_utils"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var (
	AdminAddr   = common.BigToAddress(common.Big1)
	EnabledAddr = common.BigToAddress(common.Big2)
	NoRoleAddr  = common.BigToAddress(common.Big3)
)

func AllowListTests(module modules.Module) map[string]test_utils.PrecompileTest {
	contractAddress := module.Address
	return map[string]test_utils.PrecompileTest{
		"set admin": {
			Caller: AdminAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackModifyAllowList(NoRoleAddr, AdminRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				res := GetAllowListStatus(state, contractAddress, NoRoleAddr)
				require.Equal(t, AdminRole, res)
			},
		},
		"set enabled": {
			Caller: AdminAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackModifyAllowList(NoRoleAddr, EnabledRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				res := GetAllowListStatus(state, contractAddress, NoRoleAddr)
				require.Equal(t, EnabledRole, res)
			},
		},
		"set no role": {
			Caller: AdminAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackModifyAllowList(EnabledAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				res := GetAllowListStatus(state, contractAddress, EnabledAddr)
				require.Equal(t, NoRole, res)
			},
		},
		"set no role from no role": {
			Caller: NoRoleAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackModifyAllowList(EnabledAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set enabled from no role": {
			Caller: NoRoleAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackModifyAllowList(NoRoleAddr, EnabledRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set admin from no role": {
			Caller: NoRoleAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackModifyAllowList(EnabledAddr, AdminRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set no role from enabled": {
			Caller: EnabledAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackModifyAllowList(AdminAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set enabled from enabled": {
			Caller: EnabledAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackModifyAllowList(NoRoleAddr, EnabledRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set admin from enabled": {
			Caller: EnabledAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackModifyAllowList(NoRoleAddr, AdminRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set no role with readOnly enabled": {
			Caller: AdminAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackModifyAllowList(EnabledAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"set no role insufficient gas": {
			Caller: AdminAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackModifyAllowList(EnabledAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: ModifyAllowListGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"read allow list no role": {
			Caller:      NoRoleAddr,
			Input:       PackReadAllowList(NoRoleAddr),
			SuppliedGas: ReadAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: common.Hash(NoRole).Bytes(),
		},
		"read allow list admin role": {
			Caller:      AdminAddr,
			Input:       PackReadAllowList(AdminAddr),
			SuppliedGas: ReadAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: common.Hash(AdminRole).Bytes(),
		},
		"read allow list with readOnly enabled": {
			Caller:      AdminAddr,
			Input:       PackReadAllowList(NoRoleAddr),
			SuppliedGas: ReadAllowListGasCost,
			ReadOnly:    true,
			ExpectedRes: common.Hash(NoRole).Bytes(),
		},
		"read allow list out of gas": {
			Caller: AdminAddr,
			InputFn: func(t *testing.T) []byte {
				return PackReadAllowList(NoRoleAddr)
			},
			SuppliedGas: ReadAllowListGasCost - 1,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
	}
}

func RunTestsWithAllowListSetup(t *testing.T, module modules.Module, newStateDB func(t *testing.T) contract.StateDB, tests map[string]test_utils.PrecompileTest) {
	contractAddress := module.Address
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			state := newStateDB(t)

			// Set up the state so that each address has the expected permissions at the start.
			SetAllowListRole(state, contractAddress, AdminAddr, AdminRole)
			SetAllowListRole(state, contractAddress, EnabledAddr, EnabledRole)
			require.Equal(t, AdminRole, GetAllowListStatus(state, contractAddress, AdminAddr))
			require.Equal(t, EnabledRole, GetAllowListStatus(state, contractAddress, EnabledAddr))
			require.Equal(t, NoRole, GetAllowListStatus(state, contractAddress, NoRoleAddr))

			// Run the test
			test.Run(t, module, state)
		})
	}
}
