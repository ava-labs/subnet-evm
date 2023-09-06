// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nativeminter

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"github.com/ava-labs/subnet-evm/constants"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

var tests = map[string]testutils.PrecompileTest{
	"mint funds from no role fails": {
		Caller:     allowlist.TestNoRoleAddr,
		BeforeHook: allowlist.SetDefaultRoles(Module.Address),
		InputFn: func(t testing.TB) []byte {
			input, err := PackMintInput(allowlist.TestNoRoleAddr, common.Big1)
			require.NoError(t, err)

			return input
		},
		SuppliedGas: MintGasCost,
		ReadOnly:    false,
		ExpectedErr: ErrCannotMint.Error(),
	},
	"mint funds from enabled address": {
		Caller:     allowlist.TestEnabledAddr,
		BeforeHook: allowlist.SetDefaultRoles(Module.Address),
		InputFn: func(t testing.TB) []byte {
			input, err := PackMintInput(allowlist.TestEnabledAddr, common.Big1)
			require.NoError(t, err)

			return input
		},
		SuppliedGas: MintGasCost,
		ReadOnly:    false,
		ExpectedRes: []byte{},
		AfterHook: func(t testing.TB, state contract.StateDB) {
			require.Equal(t, common.Big1, state.GetBalance(allowlist.TestEnabledAddr), "expected minted funds")
		},
	},
	"initial mint funds": {
		Caller:     allowlist.TestEnabledAddr,
		BeforeHook: allowlist.SetDefaultRoles(Module.Address),
		Config: &Config{
			InitialMint: map[common.Address]*math.HexOrDecimal256{
				allowlist.TestEnabledAddr: math.NewHexOrDecimal256(2),
			},
		},
		AfterHook: func(t testing.TB, state contract.StateDB) {
			require.Equal(t, common.Big2, state.GetBalance(allowlist.TestEnabledAddr), "expected minted funds")
		},
	},
	"mint funds from admin address": {
		Caller:     allowlist.TestAdminAddr,
		BeforeHook: allowlist.SetDefaultRoles(Module.Address),
		InputFn: func(t testing.TB) []byte {
			input, err := PackMintInput(allowlist.TestAdminAddr, common.Big1)
			require.NoError(t, err)

			return input
		},
		SuppliedGas: MintGasCost,
		ReadOnly:    false,
		ExpectedRes: []byte{},
		AfterHook: func(t testing.TB, state contract.StateDB) {
			require.Equal(t, common.Big1, state.GetBalance(allowlist.TestAdminAddr), "expected minted funds")
		},
	},
	"mint max big funds": {
		Caller:     allowlist.TestAdminAddr,
		BeforeHook: allowlist.SetDefaultRoles(Module.Address),
		InputFn: func(t testing.TB) []byte {
			input, err := PackMintInput(allowlist.TestAdminAddr, math.MaxBig256)
			require.NoError(t, err)

			return input
		},
		SuppliedGas: MintGasCost,
		ReadOnly:    false,
		ExpectedRes: []byte{},
		AfterHook: func(t testing.TB, state contract.StateDB) {
			require.Equal(t, math.MaxBig256, state.GetBalance(allowlist.TestAdminAddr), "expected minted funds")
		},
	},
	"readOnly mint with noRole fails": {
		Caller:     allowlist.TestNoRoleAddr,
		BeforeHook: allowlist.SetDefaultRoles(Module.Address),
		InputFn: func(t testing.TB) []byte {
			input, err := PackMintInput(allowlist.TestAdminAddr, common.Big1)
			require.NoError(t, err)

			return input
		},
		SuppliedGas: MintGasCost,
		ReadOnly:    true,
		ExpectedErr: vmerrs.ErrWriteProtection.Error(),
	},
	"readOnly mint with allow role fails": {
		Caller:     allowlist.TestEnabledAddr,
		BeforeHook: allowlist.SetDefaultRoles(Module.Address),
		InputFn: func(t testing.TB) []byte {
			input, err := PackMintInput(allowlist.TestEnabledAddr, common.Big1)
			require.NoError(t, err)

			return input
		},
		SuppliedGas: MintGasCost,
		ReadOnly:    true,
		ExpectedErr: vmerrs.ErrWriteProtection.Error(),
	},
	"readOnly mint with admin role fails": {
		Caller:     allowlist.TestAdminAddr,
		BeforeHook: allowlist.SetDefaultRoles(Module.Address),
		InputFn: func(t testing.TB) []byte {
			input, err := PackMintInput(allowlist.TestAdminAddr, common.Big1)
			require.NoError(t, err)

			return input
		},
		SuppliedGas: MintGasCost,
		ReadOnly:    true,
		ExpectedErr: vmerrs.ErrWriteProtection.Error(),
	},
	"insufficient gas mint from admin": {
		Caller:     allowlist.TestAdminAddr,
		BeforeHook: allowlist.SetDefaultRoles(Module.Address),
		InputFn: func(t testing.TB) []byte {
			input, err := PackMintInput(allowlist.TestEnabledAddr, common.Big1)
			require.NoError(t, err)

			return input
		},
		SuppliedGas: MintGasCost - 1,
		ReadOnly:    false,
		ExpectedErr: vmerrs.ErrOutOfGas.Error(),
	},
}

func TestContractNativeMinterRun(t *testing.T) {
	allowlist.RunPrecompileWithAllowListTests(t, Module, state.NewTestStateDB, tests)
}

func BenchmarkContractNativeMinter(b *testing.B) {
	allowlist.BenchPrecompileWithAllowList(b, Module, state.NewTestStateDB, tests)
}

func TestUnpackAndPacks(t *testing.T) {
	// Compare UnpackMintNativeCoinV2Input, PackMintNativeCoinV2 vs
	// PackMintInput, UnpackMintInput to see if they are equivalent

	// Test UnpackMintNativeCoinV2Input, PackMintNativeCoinV2
	// against PackMintInput, UnpackMintInput
	// for 1000 random addresses and amounts
	for i := 0; i < 1000; i++ {
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)
		amount := new(big.Int).SetInt64(rand.Int63())
		testUnpackAndPacks(t, addr, amount)
	}

	// Some edge cases
	testUnpackAndPacks(t, common.Address{}, common.Big0)
	testUnpackAndPacks(t, common.Address{}, common.Big1)
	testUnpackAndPacks(t, common.Address{}, math.MaxBig256)
	testUnpackAndPacks(t, common.Address{}, math.MaxBig256.Sub(math.MaxBig256, common.Big1))
	testUnpackAndPacks(t, common.Address{}, math.MaxBig256.Add(math.MaxBig256, common.Big1))
	testUnpackAndPacks(t, constants.BlackholeAddr, common.Big2)
}

func testUnpackAndPacks(t *testing.T, addr common.Address, amount *big.Int) {
	// Test PackMintNativeCoinV2, UnpackMintNativeCoinV2Input
	t.Helper()
	t.Run(fmt.Sprintf("TestUnpackAndPacks, addr: %s, amount: %s", addr.String(), amount.String()), func(t *testing.T) {
		input, err := PackMintNativeCoinV2(MintNativeCoinInput{Addr: addr, Amount: amount})
		require.NoError(t, err)
		// exclude 4 bytes for function selector
		input = input[4:]

		unpacked, err := UnpackMintNativeCoinV2Input(input)
		require.NoError(t, err)

		require.EqualValues(t, addr, unpacked.Addr)
		require.Equal(t, amount.Bytes(), unpacked.Amount.Bytes())

		// Test PackMintInput, UnpackMintInput
		input, err = PackMintInput(addr, amount)
		require.NoError(t, err)
		// exclude 4 bytes for function selector
		input = input[4:]

		to, assetAmount, err := UnpackMintInput(input)
		require.NoError(t, err)

		require.Equal(t, addr, to)
		require.Equal(t, amount.Bytes(), assetAmount.Bytes())

		// now mix and match
		// Test PackMintInput, PackMintNativeCoinV2
		input, err = PackMintInput(addr, amount)
		// exclude 4 bytes for function selector
		input = input[4:]
		require.NoError(t, err)
		input2, err := PackMintNativeCoinV2(MintNativeCoinInput{Addr: addr, Amount: amount})
		// exclude 4 bytes for function selector
		input2 = input2[4:]
		require.NoError(t, err)
		require.Equal(t, input, input2)

		// Test UnpackMintInput, UnpackMintNativeCoinV2Input
		to, assetAmount, err = UnpackMintInput(input)
		require.NoError(t, err)
		unpacked, err = UnpackMintNativeCoinV2Input(input2)
		require.NoError(t, err)
		require.Equal(t, to, unpacked.Addr)
		require.Equal(t, assetAmount.Bytes(), unpacked.Amount.Bytes())
	})
}
