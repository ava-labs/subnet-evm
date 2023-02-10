// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package deployerallowlist

import (
	"testing"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

type precompileTest struct {
	caller      common.Address
	input       func() []byte
	suppliedGas uint64
	readOnly    bool

	config config.Config

	preCondition func(t *testing.T, state *state.StateDB)
	assertState  func(t *testing.T, state *state.StateDB)

	expectedRes []byte
	expectedErr string
}

func TestContractDeployerAllowListRun(t *testing.T) {
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
				res := GetContractDeployerAllowListStatus(state, noRoleAddr)
				require.Equal(t, allowlist.AdminRole, res)
			},
		},
		"set deployer": {
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
				res := GetContractDeployerAllowListStatus(state, noRoleAddr)
				require.Equal(t, allowlist.EnabledRole, res)
			},
		},
		"set no role": {
			caller: adminAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(adminAddr, allowlist.EnabledRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := GetContractDeployerAllowListStatus(state, adminAddr)
				require.Equal(t, allowlist.EnabledRole, res)
			},
		},
		"set no role from non-admin": {
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
		"set deployer from non-admin": {
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
		"read allow list no role": {
			caller: noRoleAddr,
			input: func() []byte {
				return allowlist.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: allowlist.ReadAllowListGasCost,
			readOnly:    false,
			expectedRes: common.Hash(allowlist.EnabledRole).Bytes(),
			assertState: nil,
		},
		"read allow list admin role": {
			caller: adminAddr,
			input: func() []byte {
				return allowlist.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: allowlist.ReadAllowListGasCost,
			readOnly:    false,
			expectedRes: common.Hash(allowlist.EnabledRole).Bytes(),
			assertState: nil,
		},
		"read allow list with readOnly enabled": {
			caller: adminAddr,
			input: func() []byte {
				return allowlist.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: allowlist.ReadAllowListGasCost,
			readOnly:    true,
			expectedRes: common.Hash(allowlist.EnabledRole).Bytes(),
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
			SetContractDeployerAllowListStatus(state, adminAddr, allowlist.AdminRole)
			SetContractDeployerAllowListStatus(state, noRoleAddr, allowlist.EnabledRole)
			require.Equal(t, allowlist.AdminRole, GetContractDeployerAllowListStatus(state, adminAddr))
			require.Equal(t, allowlist.EnabledRole, GetContractDeployerAllowListStatus(state, noRoleAddr))

			blockContext := contract.NewMockBlockContext(common.Big0, 0)
			accesibleState := contract.NewMockAccessibleState(state, blockContext, snow.DefaultContextTest())
			ret, remainingGas, err := ContractDeployerAllowListPrecompile.Run(accesibleState, test.caller, ContractAddress, test.input(), test.suppliedGas, test.readOnly)
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
