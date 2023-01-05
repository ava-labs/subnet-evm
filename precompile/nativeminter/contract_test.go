// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nativeminter

import (
	"testing"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"
)

func TestContractNativeMinterRun(t *testing.T) {
	type test struct {
		caller      common.Address
		input       func() []byte
		suppliedGas uint64
		readOnly    bool
		config      *ContractNativeMinterConfig

		expectedRes []byte
		expectedErr string

		assertState func(t *testing.T, state *state.StateDB)
	}

	adminAddr := common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")
	enabledAddr := common.HexToAddress("0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B")
	noRoleAddr := common.HexToAddress("0xF60C45c607D0f41687c94C314d300f483661E13a")
	testAddr := common.HexToAddress("0x123456789")

	for name, test := range map[string]test{
		"mint funds from no role fails": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := PackMintInput(noRoleAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			suppliedGas: MintGasCost,
			readOnly:    false,
			expectedErr: ErrCannotMint.Error(),
		},
		"mint funds from enabled address": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := PackMintInput(enabledAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			suppliedGas: MintGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				require.Equal(t, common.Big1, state.GetBalance(enabledAddr), "expected minted funds")
			},
		},
		"enabled role by config": {
			caller: noRoleAddr,
			input: func() []byte {
				return precompile.PackReadAllowList(testAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost,
			readOnly:    false,
			expectedRes: common.Hash(precompile.AllowListEnabled).Bytes(),
			assertState: func(t *testing.T, state *state.StateDB) {
				require.Equal(t, precompile.AllowListEnabled, GetContractNativeMinterStatus(state, testAddr))
			},
			config: &ContractNativeMinterConfig{
				AllowListConfig: precompile.AllowListConfig{EnabledAddresses: []common.Address{testAddr}},
			},
		},
		"initial mint funds": {
			caller: enabledAddr,
			config: &ContractNativeMinterConfig{
				InitialMint: map[common.Address]*math.HexOrDecimal256{
					enabledAddr: math.NewHexOrDecimal256(2),
				},
			},
			input: func() []byte {
				return precompile.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost,
			readOnly:    false,
			expectedRes: common.Hash(precompile.AllowListNoRole).Bytes(),
			assertState: func(t *testing.T, state *state.StateDB) {
				require.Equal(t, common.Big2, state.GetBalance(enabledAddr), "expected minted funds")
			},
		},
		"mint funds from admin address": {
			caller: adminAddr,
			input: func() []byte {
				input, err := PackMintInput(adminAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			suppliedGas: MintGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				require.Equal(t, common.Big1, state.GetBalance(adminAddr), "expected minted funds")
			},
		},
		"mint max big funds": {
			caller: adminAddr,
			input: func() []byte {
				input, err := PackMintInput(adminAddr, math.MaxBig256)
				require.NoError(t, err)

				return input
			},
			suppliedGas: MintGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				require.Equal(t, math.MaxBig256, state.GetBalance(adminAddr), "expected minted funds")
			},
		},
		"readOnly mint with noRole fails": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := PackMintInput(adminAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			suppliedGas: MintGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly mint with allow role fails": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := PackMintInput(enabledAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			suppliedGas: MintGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly mint with admin role fails": {
			caller: adminAddr,
			input: func() []byte {
				input, err := PackMintInput(adminAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			suppliedGas: MintGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"insufficient gas mint from admin": {
			caller: adminAddr,
			input: func() []byte {
				input, err := PackMintInput(enabledAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			suppliedGas: MintGasCost - 1,
			readOnly:    false,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"read from noRole address": {
			caller: noRoleAddr,
			input: func() []byte {
				return precompile.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost,
			readOnly:    false,
			expectedRes: common.Hash(precompile.AllowListNoRole).Bytes(),
			assertState: func(t *testing.T, state *state.StateDB) {},
		},
		"read from noRole address readOnly enabled": {
			caller: noRoleAddr,
			input: func() []byte {
				return precompile.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost,
			readOnly:    true,
			expectedRes: common.Hash(precompile.AllowListNoRole).Bytes(),
			assertState: func(t *testing.T, state *state.StateDB) {},
		},
		"read from noRole address with insufficient gas": {
			caller: noRoleAddr,
			input: func() []byte {
				return precompile.PackReadAllowList(noRoleAddr)
			},
			suppliedGas: precompile.ReadAllowListGasCost - 1,
			readOnly:    false,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"set allow role from admin": {
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
				res := GetContractNativeMinterStatus(state, noRoleAddr)
				require.Equal(t, precompile.AllowListEnabled, res)
			},
		},
		"set allow role from non-admin fails": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := precompile.PackModifyAllowList(noRoleAddr, precompile.AllowListEnabled)
				require.NoError(t, err)

				return input
			},
			suppliedGas: precompile.ModifyAllowListGasCost,
			readOnly:    false,
			expectedErr: precompile.ErrCannotModifyAllowList.Error(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			db := rawdb.NewMemoryDatabase()
			state, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
			require.NoError(t, err)

			// Set up the state so that each address has the expected permissions at the start.
			SetContractNativeMinterStatus(state, adminAddr, precompile.AllowListAdmin)
			SetContractNativeMinterStatus(state, enabledAddr, precompile.AllowListEnabled)
			SetContractNativeMinterStatus(state, noRoleAddr, precompile.AllowListNoRole)
			require.Equal(t, precompile.AllowListAdmin, GetContractNativeMinterStatus(state, adminAddr))
			require.Equal(t, precompile.AllowListEnabled, GetContractNativeMinterStatus(state, enabledAddr))
			require.Equal(t, precompile.AllowListNoRole, GetContractNativeMinterStatus(state, noRoleAddr))

			blockContext := precompile.NewMockBlockContext(common.Big0, 0)
			accesibleState := precompile.NewMockAccessibleState(state, blockContext, snow.DefaultContextTest())
			if test.config != nil {
				test.config.Configure(params.TestChainConfig, state, blockContext)
			}
			ret, remainingGas, err := ContractNativeMinterPrecompile.Run(accesibleState, test.caller, Address, test.input(), test.suppliedGas, test.readOnly)
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
