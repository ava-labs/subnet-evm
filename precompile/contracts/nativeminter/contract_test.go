// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nativeminter

import (
	"testing"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"
)

func TestContractNativeMinterRun(t *testing.T) {
	adminAddr := common.BigToAddress(common.Big0)
	enabledAddr := common.BigToAddress(common.Big1)
	noRoleAddr := common.BigToAddress(common.Big2)

	for name, test := range map[string]contract.PrecompileTest{
		"mint funds from no role fails": {
			Caller: noRoleAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackMintInput(noRoleAddr, common.Big1)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: MintGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotMint.Error(),
		},
		"mint funds from enabled address": {
			Caller: enabledAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackMintInput(enabledAddr, common.Big1)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: MintGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				require.Equal(t, common.Big1, state.GetBalance(enabledAddr), "expected minted funds")
			},
		},
		"initial mint funds": {
			Caller: enabledAddr,
			Config: &Config{
				InitialMint: map[common.Address]*math.HexOrDecimal256{
					enabledAddr: math.NewHexOrDecimal256(2),
				},
			},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				require.Equal(t, common.Big2, state.GetBalance(enabledAddr), "expected minted funds")
			},
		},
		"mint funds from admin address": {
			Caller: adminAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackMintInput(adminAddr, common.Big1)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: MintGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				require.Equal(t, common.Big1, state.GetBalance(adminAddr), "expected minted funds")
			},
		},
		"mint max big funds": {
			Caller: adminAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackMintInput(adminAddr, math.MaxBig256)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: MintGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				require.Equal(t, math.MaxBig256, state.GetBalance(adminAddr), "expected minted funds")
			},
		},
		"readOnly mint with noRole fails": {
			Caller: noRoleAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackMintInput(adminAddr, common.Big1)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: MintGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly mint with allow role fails": {
			Caller: enabledAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackMintInput(enabledAddr, common.Big1)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: MintGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly mint with admin role fails": {
			Caller: adminAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackMintInput(adminAddr, common.Big1)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: MintGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"insufficient gas mint from admin": {
			Caller: adminAddr,
			Input: func(tt *testing.T) []byte {
				input, err := PackMintInput(enabledAddr, common.Big1)
				require.NoError(tt, err)

				return input
			}(t),
			SuppliedGas: MintGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			db := rawdb.NewMemoryDatabase()
			state, err := state.New(common.Hash{}, state.NewDatabase(db), nil)
			require.NoError(t, err)

			// Set up the state so that each address has the expected permissions at the start.
			SetContractNativeMinterStatus(state, adminAddr, allowlist.AdminRole)
			SetContractNativeMinterStatus(state, enabledAddr, allowlist.EnabledRole)
			require.Equal(t, allowlist.AdminRole, GetContractNativeMinterStatus(state, adminAddr))
			require.Equal(t, allowlist.EnabledRole, GetContractNativeMinterStatus(state, enabledAddr))

			blockContext := contract.NewMockBlockContext(common.Big0, 0)
			accesibleState := contract.NewMockAccessibleState(state, blockContext, snow.DefaultContextTest())
			if test.Config != nil {
				Module.Configure(params.TestChainConfig, test.Config, state, blockContext)
			}
			if test.Input != nil {
				ret, remainingGas, err := ContractNativeMinterPrecompile.Run(accesibleState, test.Caller, ContractAddress, test.Input, test.SuppliedGas, test.ReadOnly)
				if len(test.ExpectedErr) != 0 {
					require.ErrorContains(t, err, test.ExpectedErr)
				} else {
					require.NoError(t, err)
				}

				require.Equal(t, uint64(0), remainingGas)
				require.Equal(t, test.ExpectedRes, ret)
			}

			if test.AfterHook != nil {
				test.AfterHook(t, state)
			}
		})
	}
}
