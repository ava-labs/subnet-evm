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
	"github.com/ava-labs/subnet-evm/precompile/deployerallowlist"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestContractDeployerAllowListRun(t *testing.T) {
	adminAddr := common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")
	noRoleAddr := common.HexToAddress("0xF60C45c607D0f41687c94C314d300f483661E13a")

	for name, test := range map[string]precompileTest{
		"set admin": {
			caller: adminAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(noRoleAddr, allowlist.AllowListAdmin)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := deployerallowlist.GetContractDeployerAllowListStatus(state, noRoleAddr)
				require.Equal(t, allowlist.AllowListAdmin, res)
			},
		},
		"set deployer": {
			caller: adminAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(noRoleAddr, allowlist.AllowListEnabled)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := deployerallowlist.GetContractDeployerAllowListStatus(state, noRoleAddr)
				require.Equal(t, allowlist.AllowListEnabled, res)
			},
		},
		"set no role": {
			caller: adminAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(adminAddr, allowlist.AllowListEnabled)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := deployerallowlist.GetContractDeployerAllowListStatus(state, adminAddr)
				require.Equal(t, allowlist.AllowListEnabled, res)
			},
		},
		"set no role from non-admin": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(adminAddr, allowlist.AllowListEnabled)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"set deployer from non-admin": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(adminAddr, allowlist.AllowListEnabled)
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
				input, err := allowlist.PackModifyAllowList(adminAddr, allowlist.AllowListAdmin)
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
				input, err := allowlist.PackModifyAllowList(adminAddr, allowlist.AllowListEnabled)
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
				input, err := allowlist.PackModifyAllowList(adminAddr, allowlist.AllowListEnabled)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost - 1,
			readOnly:    false,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"read allow list no role": {
			caller: noRoleAddr,
			input: func() []byte {
				return allowlist.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: allowlist.ReadAllowListGasCost,
			readOnly:    false,
			expectedRes: common.Hash(allowlist.AllowListEnabled).Bytes(),
			assertState: nil,
		},
		"read allow list admin role": {
			caller: adminAddr,
			input: func() []byte {
				return allowlist.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: allowlist.ReadAllowListGasCost,
			readOnly:    false,
			expectedRes: common.Hash(allowlist.AllowListEnabled).Bytes(),
			assertState: nil,
		},
		"read allow list with readOnly enabled": {
			caller: adminAddr,
			input: func() []byte {
				return allowlist.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: allowlist.ReadAllowListGasCost,
			readOnly:    true,
			expectedRes: common.Hash(allowlist.AllowListEnabled).Bytes(),
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
			deployerallowlist.SetContractDeployerAllowListStatus(state, adminAddr, allowlist.AllowListAdmin)
			deployerallowlist.SetContractDeployerAllowListStatus(state, noRoleAddr, allowlist.AllowListEnabled)
			require.Equal(t, allowlist.AllowListAdmin, deployerallowlist.GetContractDeployerAllowListStatus(state, adminAddr))
			require.Equal(t, allowlist.AllowListEnabled, deployerallowlist.GetContractDeployerAllowListStatus(state, noRoleAddr))

			blockContext := precompile.NewMockBlockContext(common.Big0, 0)
			accesibleState := precompile.NewMockAccessibleState(state, blockContext, snow.DefaultContextTest())
			ret, remainingGas, err := deployerallowlist.ContractDeployerAllowListPrecompile.Run(accesibleState, test.caller, deployerallowlist.ContractAddress, test.input(), test.suppliedGas, test.readOnly)
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
