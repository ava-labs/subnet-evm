// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nativeminter

import (
	"testing"

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"
)

func TestContractNativeMinterRun(t *testing.T) {
	tests := map[string]testutils.PrecompileTest{
		"mint funds from no role fails": {
			Caller: allowlist.NoRoleAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackMintInput(allowlist.NoRoleAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: MintGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotMint.Error(),
		},
		"mint funds from enabled address": {
			Caller: allowlist.EnabledAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackMintInput(allowlist.EnabledAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: MintGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				require.Equal(t, common.Big1, state.GetBalance(allowlist.EnabledAddr), "expected minted funds")
			},
		},
		"initial mint funds": {
			Caller: allowlist.EnabledAddr,
			Config: &Config{
				InitialMint: map[common.Address]*math.HexOrDecimal256{
					allowlist.EnabledAddr: math.NewHexOrDecimal256(2),
				},
			},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				require.Equal(t, common.Big2, state.GetBalance(allowlist.EnabledAddr), "expected minted funds")
			},
		},
		"mint funds from admin address": {
			Caller: allowlist.AdminAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackMintInput(allowlist.AdminAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: MintGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				require.Equal(t, common.Big1, state.GetBalance(allowlist.AdminAddr), "expected minted funds")
			},
		},
		"mint max big funds": {
			Caller: allowlist.AdminAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackMintInput(allowlist.AdminAddr, math.MaxBig256)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: MintGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				require.Equal(t, math.MaxBig256, state.GetBalance(allowlist.AdminAddr), "expected minted funds")
			},
		},
		"readOnly mint with noRole fails": {
			Caller: allowlist.NoRoleAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackMintInput(allowlist.AdminAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: MintGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly mint with allow role fails": {
			Caller: allowlist.EnabledAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackMintInput(allowlist.EnabledAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: MintGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly mint with admin role fails": {
			Caller: allowlist.AdminAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackMintInput(allowlist.AdminAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: MintGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"insufficient gas mint from admin": {
			Caller: allowlist.AdminAddr,
			InputFn: func(t *testing.T) []byte {
				input, err := PackMintInput(allowlist.EnabledAddr, common.Big1)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: MintGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
	}

	allowlist.RunTestsWithAllowListSetup(t, Module, state.NewTestStateDB, tests)
	allowlist.RunTestsWithAllowListSetup(t, Module, state.NewTestStateDB, allowlist.AllowListTests(Module))
}
