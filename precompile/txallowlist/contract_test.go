// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txallowlist

import (
	"testing"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestTxAllowListRun(t *testing.T) {
	type test struct {
		caller         common.Address
		precompileAddr common.Address
		input          func() []byte
		suppliedGas    uint64
		readOnly       bool

		expectedRes []byte
		expectedErr string

		assertState func(t *testing.T, state *state.StateDB)
	}

	adminAddr := common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")
	noRoleAddr := common.HexToAddress("0xF60C45c607D0f41687c94C314d300f483661E13a")

	for name, test := range map[string]test{
		"set admin": {
			caller: adminAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(noRoleAddr, precompile.AllowListAdmin)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := GetTxAllowListStatus(state, noRoleAddr)
				require.Equal(t, precompile.AllowListAdmin, res)
			},
		},
		"set allowed": {
			caller: adminAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(noRoleAddr, precompile.AllowListEnabled)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := GetTxAllowListStatus(state, noRoleAddr)
				require.Equal(t, precompile.AllowListEnabled, res)
			},
		},
		"set no role": {
			caller: adminAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(adminAddr, precompile.AllowListNoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := GetTxAllowListStatus(state, adminAddr)
				require.Equal(t, precompile.AllowListNoRole, res)
			},
		},
		"set no role from non-admin": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(adminAddr, precompile.AllowListNoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: precompile.ErrCannotModifyAllowList.Error(),
		},
		"set allowed from non-admin": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(adminAddr, precompile.AllowListEnabled)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: precompile.ErrCannotModifyAllowList.Error(),
		},
		"set admin from non-admin": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(adminAddr, precompile.AllowListAdmin)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: precompile.ErrCannotModifyAllowList.Error(),
		},
		"set no role with readOnly enabled": {
			caller:         adminAddr,
			precompileAddr: Address,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(adminAddr, precompile.AllowListNoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"set no role insufficient gas": {
			caller: adminAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(adminAddr, precompile.AllowListNoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost - 1,
			readOnly:    false,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"read allow list no role": {
			caller: noRoleAddr,
			input: func() []byte {
				return precompile.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost,
			readOnly:    false,
			expectedRes: common.Hash(precompile.AllowListNoRole).Bytes(),
			assertState: nil,
		},
		"read allow list admin role": {
			caller: adminAddr,
			input: func() []byte {
				return precompile.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost,
			readOnly:    false,
			expectedRes: common.Hash(precompile.AllowListNoRole).Bytes(),
			assertState: nil,
		},
		"read allow list with readOnly enabled": {
			caller: adminAddr,
			input: func() []byte {
				return precompile.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost,
			readOnly:    true,
			expectedRes: common.Hash(precompile.AllowListNoRole).Bytes(),
			assertState: nil,
		},
		"read allow list out of gas": {
			caller: adminAddr,
			input: func() []byte {
				return precompile.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost - 1,
			readOnly:    true,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			db := rawdb.NewMemoryDatabase()
			state, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
			require.NoError(t, err)

			// Set up the state so that each address has the expected permissions at the start.
			SetTxAllowListStatus(state, adminAddr, precompile.AllowListAdmin)
			require.Equal(t, precompile.AllowListAdmin, GetTxAllowListStatus(state, adminAddr))

			blockContext := precompile.NewMockBlockContext(common.Big0, 0)
			accesibleState := precompile.NewMockAccessibleState(state, blockContext, snow.DefaultContextTest())
			ret, remainingGas, err := TxAllowListPrecompile.Run(accesibleState, test.caller, Address, test.input(), test.suppliedGas, test.readOnly)
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
