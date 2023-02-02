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

		config *allowlist.AllowListConfig

		assertState func(t *testing.T, state *state.StateDB)
	}

	adminAddr := common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")
	enabledAddr := common.HexToAddress("0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B")
	noRoleAddr := common.HexToAddress("0xF60C45c607D0f41687c94C314d300f483661E13a")
	dummyContractAddr := common.HexToAddress("0x0000000000000000000000000000000000000000")
	testAllowListPrecompile := allowlist.CreateAllowListPrecompile(dummyContractAddr)

	for name, test := range map[string]test{
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
				res := allowlist.GetAllowListStatus(state, dummyContractAddr, noRoleAddr)
				require.Equal(t, allowlist.AllowListAdmin, res)
			},
		},
		"set enabled": {
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
				res := allowlist.GetAllowListStatus(state, dummyContractAddr, noRoleAddr)
				require.Equal(t, allowlist.AllowListEnabled, res)
			},
		},
		"set no role": {
			caller: adminAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(enabledAddr, allowlist.AllowListNoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				res := allowlist.GetAllowListStatus(state, dummyContractAddr, enabledAddr)
				require.Equal(t, allowlist.AllowListNoRole, res)
			},
		},
		"set no role from no role": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(enabledAddr, allowlist.AllowListNoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"set enabled from no role": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(noRoleAddr, allowlist.AllowListEnabled)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"set admin from no role": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(enabledAddr, allowlist.AllowListAdmin)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"set no role from enabled": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(adminAddr, allowlist.AllowListNoRole)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"set enabled from enabled": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(noRoleAddr, allowlist.AllowListEnabled)
				require.NoError(t, err)

				return input
			},
			suppliedGas: allowlist.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: allowlist.ErrCannotModifyAllowList.Error(),
		},
		"set admin from enabled": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := allowlist.PackModifyAllowList(noRoleAddr, allowlist.AllowListAdmin)
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
				input, err := allowlist.PackModifyAllowList(enabledAddr, allowlist.AllowListNoRole)
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
				input, err := allowlist.PackModifyAllowList(enabledAddr, allowlist.AllowListNoRole)
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
			expectedRes: common.Hash(allowlist.AllowListNoRole).Bytes(),
			assertState: nil,
		},
		"read allow list admin role": {
			caller: adminAddr,
			input: func() []byte {
				return allowlist.PackReadAllowList(adminAddr)
			},
			suppliedGas: allowlist.ReadAllowListGasCost,
			readOnly:    false,
			expectedRes: common.Hash(allowlist.AllowListAdmin).Bytes(),
			assertState: nil,
		},
		"read allow list with readOnly enabled": {
			caller: adminAddr,
			input: func() []byte {
				return allowlist.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: allowlist.ReadAllowListGasCost,
			readOnly:    true,
			expectedRes: common.Hash(allowlist.AllowListNoRole).Bytes(),
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
		"initial config sets admins": {
			config: &allowlist.AllowListConfig{
				AdminAddresses: []common.Address{noRoleAddr, enabledAddr},
			},
			suppliedGas: 0,
			readOnly:    false,
			expectedErr: "",
			assertState: func(t *testing.T, state *state.StateDB) {
				require.Equal(t, allowlist.AllowListAdmin, allowlist.GetAllowListStatus(state, dummyContractAddr, noRoleAddr))
				require.Equal(t, allowlist.AllowListAdmin, allowlist.GetAllowListStatus(state, dummyContractAddr, enabledAddr))
			},
		},
		"initial config sets enabled": {
			config: &allowlist.AllowListConfig{
				EnabledAddresses: []common.Address{noRoleAddr, adminAddr},
			},
			suppliedGas: 0,
			readOnly:    false,
			expectedErr: "",
			assertState: func(t *testing.T, state *state.StateDB) {
				require.Equal(t, allowlist.AllowListEnabled, allowlist.GetAllowListStatus(state, dummyContractAddr, adminAddr))
				require.Equal(t, allowlist.AllowListEnabled, allowlist.GetAllowListStatus(state, dummyContractAddr, noRoleAddr))
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			db := rawdb.NewMemoryDatabase()
			state, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
			require.NoError(t, err)

			// Set up the state so that each address has the expected permissions at the start.
			allowlist.SetAllowListRole(state, dummyContractAddr, adminAddr, allowlist.AllowListAdmin)
			allowlist.SetAllowListRole(state, dummyContractAddr, enabledAddr, allowlist.AllowListEnabled)
			require.Equal(t, allowlist.AllowListAdmin, allowlist.GetAllowListStatus(state, dummyContractAddr, adminAddr))
			require.Equal(t, allowlist.AllowListEnabled, allowlist.GetAllowListStatus(state, dummyContractAddr, enabledAddr))

			if test.config != nil {
				test.config.Configure(state, dummyContractAddr)
			}

			blockContext := precompile.NewMockBlockContext(common.Big0, 0)
			accesibleState := precompile.NewMockAccessibleState(state, blockContext, snow.DefaultContextTest())
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
