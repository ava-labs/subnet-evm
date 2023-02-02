// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package statefulprecompiles

import (
	"testing"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/txallowlist"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestTxAllowListRun(t *testing.T) {
	adminAddr := common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")
	noRoleAddr := common.HexToAddress("0xF60C45c607D0f41687c94C314d300f483661E13a")

	for name, test := range map[string]precompileTest{
		"set admin": {
			caller: adminAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(noRoleAddr, allowlist.AdminRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := txallowlist.GetTxAllowListStatus(state, noRoleAddr)
				require.Equal(t, allowlist.AdminRole, res)
			},
		},
		"set allowed": {
			caller: adminAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(noRoleAddr, allowlist.EnabledRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := txallowlist.GetTxAllowListStatus(state, noRoleAddr)
				require.Equal(t, allowlist.EnabledRole, res)
			},
		},
		"set no role": {
			caller: adminAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(adminAddr, allowlist.NoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := txallowlist.GetTxAllowListStatus(state, adminAddr)
				require.Equal(t, allowlist.NoRole, res)
			},
		},
		"set no role from non-admin": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(adminAddr, allowlist.NoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"set allowed from non-admin": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(adminAddr, allowlist.EnabledRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"set admin from non-admin": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(adminAddr, allowlist.AdminRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"set no role with readOnly enabled": {
			caller: adminAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(adminAddr, allowlist.EnabledRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"set no role insufficient gas": {
			caller: adminAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(adminAddr, allowlist.EnabledRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost - 1,
			readOnly:    false,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"read allow listwith no role": {
			caller: noRoleAddr,
			input: func() []byte {
				return allowlist.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: allowlist.ReadAllowListGasCost,
			readOnly:    false,
			expectedRes: common.Hash(allowlist.NoRole).Bytes(),
			assertState: nil,
		},
		"read allow list with admin role": {
			caller: adminAddr,
			input: func() []byte {
				return allowlist.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: allowlist.ReadAllowListGasCost,
			readOnly:    false,
			expectedRes: common.Hash(allowlist.NoRole).Bytes(),
			assertState: nil,
		},
		"read allow list with readOnly enabled": {
			caller: adminAddr,
			input: func() []byte {
				return allowlist.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: allowlist.ReadAllowListGasCost,
			readOnly:    true,
			expectedRes: common.Hash(allowlist.NoRole).Bytes(),
			assertState: nil,
		},
		"read allow list out of gas": {
			caller: adminAddr,
			input: func() []byte {
				return allowlist.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: allowlist.ReadAllowListGasCost - 1,
			readOnly:    true,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			db := rawdb.NewMemoryDatabase()
			state, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
			require.NoError(t, err)

			// Set up the state so that each address has the expected permissions at the start.
			txallowlist.SetTxAllowListStatus(state, adminAddr, allowlist.AdminRole)
			require.Equal(t, allowlist.AdminRole, txallowlist.GetTxAllowListStatus(state, adminAddr))

			blockContext := precompile.NewMockBlockContext(common.Big0, 0)
			accesibleState := precompile.NewMockAccessibleState(state, blockContext, snow.DefaultContextTest())
			ret, remainingGas, err := txallowlist.TxAllowListPrecompile.Run(accesibleState, test.caller, txallowlist.ContractAddress, test.input(), test.suppliedGas, test.readOnly)
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
