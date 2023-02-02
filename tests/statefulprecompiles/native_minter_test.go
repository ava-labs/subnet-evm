// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package statefulprecompiles

import (
	"testing"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/nativeminter"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"
)

func TestContractNativeMinterRun(t *testing.T) {
	adminAddr := common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")
	enabledAddr := common.HexToAddress("0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B")
	noRoleAddr := common.HexToAddress("0xF60C45c607D0f41687c94C314d300f483661E13a")

	for name, test := range map[string]precompileTest{
		"mint funds from no role fails": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := nativeminter.PackMintInput(noRoleAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			suppliedGas: nativeminter.MintGasCost,
			readOnly:    false,
			expectedErr: nativeminter.ErrCannotMint.Error(),
		},
		"mint funds from enabled address": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := nativeminter.PackMintInput(enabledAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			suppliedGas: nativeminter.MintGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				require.Equal(t, common.Big1, state.GetBalance(enabledAddr), "expected minted funds")
			},
		},
		"initial mint funds": {
			caller: enabledAddr,
			config: &nativeminter.ContractNativeMinterConfig{
				InitialMint: map[common.Address]*math.HexOrDecimal256{
					enabledAddr: math.NewHexOrDecimal256(2),
				},
			},
			assertState: func(t *testing.T, state *state.StateDB) {
				require.Equal(t, common.Big2, state.GetBalance(enabledAddr), "expected minted funds")
			},
		},
		"mint funds from admin address": {
			caller: adminAddr,
			input: func() []byte {
				input, err := nativeminter.PackMintInput(adminAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			suppliedGas: nativeminter.MintGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				require.Equal(t, common.Big1, state.GetBalance(adminAddr), "expected minted funds")
			},
		},
		"mint max big funds": {
			caller: adminAddr,
			input: func() []byte {
				input, err := nativeminter.PackMintInput(adminAddr, math.MaxBig256)
				require.NoError(t, err)

				return input
			},
			suppliedGas: nativeminter.MintGasCost,
			readOnly:    false,
			expectedRes: []byte{},
			assertState: func(t *testing.T, state *state.StateDB) {
				require.Equal(t, math.MaxBig256, state.GetBalance(adminAddr), "expected minted funds")
			},
		},
		"readOnly mint with noRole fails": {
			caller: noRoleAddr,
			input: func() []byte {
				input, err := nativeminter.PackMintInput(adminAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			suppliedGas: nativeminter.MintGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly mint with allow role fails": {
			caller: enabledAddr,
			input: func() []byte {
				input, err := nativeminter.PackMintInput(enabledAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			suppliedGas: nativeminter.MintGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly mint with admin role fails": {
			caller: adminAddr,
			input: func() []byte {
				input, err := nativeminter.PackMintInput(adminAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			suppliedGas: nativeminter.MintGasCost,
			readOnly:    true,
			expectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"insufficient gas mint from admin": {
			caller: adminAddr,
			input: func() []byte {
				input, err := nativeminter.PackMintInput(enabledAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			suppliedGas: nativeminter.MintGasCost - 1,
			readOnly:    false,
			expectedErr: vmerrs.ErrOutOfGas.Error(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			db := rawdb.NewMemoryDatabase()
			state, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
			require.NoError(t, err)

			// Set up the state so that each address has the expected permissions at the start.
			nativeminter.SetContractNativeMinterStatus(state, adminAddr, allowlist.AdminRole)
			nativeminter.SetContractNativeMinterStatus(state, enabledAddr, allowlist.EnabledRole)
			require.Equal(t, allowlist.AdminRole, nativeminter.GetContractNativeMinterStatus(state, adminAddr))
			require.Equal(t, allowlist.EnabledRole, nativeminter.GetContractNativeMinterStatus(state, enabledAddr))

			blockContext := precompile.NewMockBlockContext(common.Big0, 0)
			accesibleState := precompile.NewMockAccessibleState(state, blockContext, snow.DefaultContextTest())
			if test.config != nil {
				test.config.Configure(params.TestChainConfig, state, blockContext)
			}
			if test.input != nil {
				ret, remainingGas, err := nativeminter.ContractNativeMinterPrecompile.Run(accesibleState, test.caller, nativeminter.ContractAddress, test.input(), test.suppliedGas, test.readOnly)
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
