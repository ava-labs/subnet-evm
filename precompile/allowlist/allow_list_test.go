// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package allowlist

import (
	"testing"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestAllowListRun(t *testing.T) {
	type test struct {
		caller      common.Address
		input       func() []byte
		suppliedGas uint64
		readOnly    bool

		expectedRes []byte
		expectedErr string

		config *AllowList

		assertState func(t *testing.T, state *state.StateDB)
	}

	adminAddr := common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")
	enabledAddr := common.HexToAddress("0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B")
	noRoleAddr := common.HexToAddress("0xF60C45c607D0f41687c94C314d300f483661E13a")
	dummyContractAddr := common.HexToAddress("0x0000000000000000000000000000000000000000")
	testAllowListPrecompile := CreateAllowListPrecompile(dummyContractAddr)

	for name, test := range map[string]test{
		"set admin": {
			caller: adminAddr,
			input: func() []byte {
				input, err := PackModifyAllowList(noRoleAddr, AdminRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := GetAllowListStatus(state, dummyContractAddr, noRoleAddr)
				require.Equal(t, AdminRole, res)
			},
		},
		"set enabled": {
			caller: adminAddr,
			input: func() []byte {
				input, err := PackModifyAllowList(noRoleAddr, EnabledRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := GetAllowListStatus(state, dummyContractAddr, noRoleAddr)
				require.Equal(t, EnabledRole, res)
			},
		},
		"set no role": {
			caller: adminAddr,
			input: func() []byte {
				input, err := PackModifyAllowList(enabledAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := GetAllowListStatus(state, dummyContractAddr, enabledAddr)
				require.Equal(t, NoRole, res)
			},
		},
		"set no role from no role": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := PackModifyAllowList(enabledAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set enabled from no role": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := PackModifyAllowList(noRoleAddr, EnabledRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set admin from no role": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := PackModifyAllowList(enabledAddr, AdminRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set no role from enabled": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := PackModifyAllowList(adminAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set enabled from enabled": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := PackModifyAllowList(noRoleAddr, EnabledRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set admin from enabled": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := PackModifyAllowList(noRoleAddr, AdminRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: ErrCannotModifyAllowList.Error(),
		},
		"set no role with readOnly enabled": {
			caller: adminAddr,
			input: func() []byte {
				input, err := PackModifyAllowList(enabledAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: ModifyAllowListGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"set no role insufficient gas": {
			caller: adminAddr,
			input: func() []byte {
				input, err := PackModifyAllowList(enabledAddr, NoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: ModifyAllowListGasCost - 1,
			readOnly:    false,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"read allow list no role": {
			caller: noRoleAddr,
			input: func() []byte {
				return PackReadAllowList(noRoleAddr)
			},
			suppliedGas: ReadAllowListGasCost,
			readOnly:    false,
			expectedRes: common.Hash(NoRole).Bytes(),
			assertState: nil,
		},
		"read allow list admin role": {
			caller: adminAddr,
			input: func() []byte {
				return PackReadAllowList(adminAddr)
			},
			suppliedGas: ReadAllowListGasCost,
			readOnly:    false,
			expectedRes: common.Hash(AdminRole).Bytes(),
			assertState: nil,
		},
		"read allow list with readOnly enabled": {
			caller: adminAddr,
			input: func() []byte {
				return PackReadAllowList(noRoleAddr)
			},
			suppliedGas: ReadAllowListGasCost,
			readOnly:    true,
			expectedRes: common.Hash(NoRole).Bytes(),
			assertState: nil,
		},
		"read allow list out of gas": {
			caller: adminAddr,
			input: func() []byte {
				return PackReadAllowList(noRoleAddr)
			},
			suppliedGas: ReadAllowListGasCost - 1,
			readOnly:    true,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"initial config sets admins": {
			config: &AllowList{
				AdminAddresses: []common.Address{noRoleAddr, enabledAddr},
			},
			suppliedGas: 0,
			readOnly:    false,
			expectedErr: "",
			assertState: func(t *testing.T, state *state.StateDB) {
				require.Equal(t, AdminRole, GetAllowListStatus(state, dummyContractAddr, noRoleAddr))
				require.Equal(t, AdminRole, GetAllowListStatus(state, dummyContractAddr, enabledAddr))
			},
		},
		"initial config sets enabled": {
			config: &AllowList{
				EnabledAddresses: []common.Address{noRoleAddr, adminAddr},
			},
			suppliedGas: 0,
			readOnly:    false,
			expectedErr: "",
			assertState: func(t *testing.T, state *state.StateDB) {
				require.Equal(t, EnabledRole, GetAllowListStatus(state, dummyContractAddr, adminAddr))
				require.Equal(t, EnabledRole, GetAllowListStatus(state, dummyContractAddr, noRoleAddr))
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			db := rawdb.NewMemoryDatabase()
			state, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
			require.NoError(t, err)

			// Set up the state so that each address has the expected permissions at the start.
			SetAllowListRole(state, dummyContractAddr, adminAddr, AdminRole)
			SetAllowListRole(state, dummyContractAddr, enabledAddr, EnabledRole)
			require.Equal(t, AdminRole, GetAllowListStatus(state, dummyContractAddr, adminAddr))
			require.Equal(t, EnabledRole, GetAllowListStatus(state, dummyContractAddr, enabledAddr))

			if test.config != nil {
				test.config.Configure(state, dummyContractAddr)
			}

			blockContext := contract.NewMockBlockContext(common.Big0, 0)
			accesibleState := contract.NewMockAccessibleState(state, blockContext, snow.DefaultContextTest())
			if test.input != nil {
				ret, remainingGas, err := testAllowListPrecompile.Run(accesibleState, test.caller, dummyContractAddr, test.input(), test.suppliedGas, test.readOnly)

				if len(test.expectedErr) != 0 {
					require.ErrorContains(t, err, test.expectedErr)
				} else {
					require.NoError(t, err)
				}

				require.Equal(t, uint64(0), remainingGas)
				require.Equal(t, test.expectedRes, ret)
			}

			if test.assertState != nil {
				test.assertState(t, state)
			}
		})
	}
}
