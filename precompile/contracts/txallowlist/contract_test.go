// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txallowlist

import (
	"testing"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestTxAllowListRun(t *testing.T) {
	adminAddr := common.BigToAddress(common.Big0)
	noRoleAddr := common.BigToAddress(common.Big2)

	for name, test := range map[string]contract.PrecompileTest{
		"set admin": {
			Caller: adminAddr,
			Input: func(tt *testing.T) []byte {
				input, err := allowlist.PackModifyAllowList(noRoleAddr, allowlist.AdminRole)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				res := GetTxAllowListStatus(state, noRoleAddr)
				require.Equal(t, allowlist.AdminRole, res)
			},
		},
		"set allowed": {
			Caller: adminAddr,
			Input: func(tt *testing.T) []byte {
				input, err := allowlist.PackModifyAllowList(noRoleAddr, allowlist.EnabledRole)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				res := GetTxAllowListStatus(state, noRoleAddr)
				require.Equal(t, allowlist.EnabledRole, res)
			},
		},
		"set no role": {
			Caller: adminAddr,
			Input: func(tt *testing.T) []byte {
				input, err := allowlist.PackModifyAllowList(adminAddr, allowlist.NoRole)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				res := GetTxAllowListStatus(state, adminAddr)
				require.Equal(t, allowlist.NoRole, res)
			},
		},
		"set no role from non-admin": {
			Caller: noRoleAddr,
			Input: func(tt *testing.T) []byte {
				input, err := allowlist.PackModifyAllowList(adminAddr, allowlist.NoRole)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"set allowed from non-admin": {
			Caller: noRoleAddr,
			Input: func(tt *testing.T) []byte {
				input, err := allowlist.PackModifyAllowList(adminAddr, allowlist.EnabledRole)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"set admin from non-admin": {
			Caller: noRoleAddr,
			Input: func(tt *testing.T) []byte {
				input, err := allowlist.PackModifyAllowList(adminAddr, allowlist.AdminRole)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    false,
			ExpectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"set no role with readOnly enabled": {
			Caller: adminAddr,
			Input: func(tt *testing.T) []byte {
				input, err := allowlist.PackModifyAllowList(adminAddr, allowlist.EnabledRole)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: allowlist.ModifyAllowListGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"set no role insufficient gas": {
			Caller: adminAddr,
			Input: func(tt *testing.T) []byte {
				input, err := allowlist.PackModifyAllowList(adminAddr, allowlist.EnabledRole)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: allowlist.ModifyAllowListGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"read allow list with no role": {
			Caller:      noRoleAddr,
			Input:       allowlist.PackReadAllowList(noRoleAddr),
			SuppliedGas: allowlist.ReadAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: common.Hash(allowlist.NoRole).Bytes(),
			AfterHook:   nil,
		},
		"read allow list with admin role": {
			Caller:      adminAddr,
			Input:       allowlist.PackReadAllowList(noRoleAddr),
			SuppliedGas: allowlist.ReadAllowListGasCost,
			ReadOnly:    false,
			ExpectedRes: common.Hash(allowlist.NoRole).Bytes(),
			AfterHook:   nil,
		},
		"read allow list with readOnly enabled": {
			Caller:      adminAddr,
			Input:       allowlist.PackReadAllowList(noRoleAddr),
			SuppliedGas: allowlist.ReadAllowListGasCost,
			ReadOnly:    true,
			ExpectedRes: common.Hash(allowlist.NoRole).Bytes(),
			AfterHook:   nil,
		},
		"read allow list out of gas": {
			Caller:      adminAddr,
			Input:       allowlist.PackReadAllowList(noRoleAddr),
			SuppliedGas: allowlist.ReadAllowListGasCost - 1,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			db := rawdb.NewMemoryDatabase()
			state, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
			require.NoError(t, err)

			// Set up the state so that each address has the expected permissions at the start.
			SetTxAllowListStatus(state, adminAddr, allowlist.AdminRole)
			require.Equal(t, allowlist.AdminRole, GetTxAllowListStatus(state, adminAddr))

			blockContext := contract.NewMockBlockContext(common.Big0, 0)
			accesibleState := contract.NewMockAccessibleState(state, blockContext, snow.DefaultContextTest())
			ret, remainingGas, err := TxAllowListPrecompile.Run(accesibleState, test.Caller, ContractAddress, test.Input, test.SuppliedGas, test.ReadOnly)
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
