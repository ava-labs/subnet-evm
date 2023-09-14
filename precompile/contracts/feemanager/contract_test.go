// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package feemanager

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	testFeeConfig = commontype.FeeConfig{
		GasLimit:        big.NewInt(8_000_000),
		TargetBlockRate: 2, // in seconds

		MinBaseFee:               big.NewInt(25_000_000_000),
		TargetGas:                big.NewInt(15_000_000),
		BaseFeeChangeDenominator: big.NewInt(36),

		MinBlockGasCost:  big.NewInt(0),
		MaxBlockGasCost:  big.NewInt(1_000_000),
		BlockGasCostStep: big.NewInt(200_000),
	}
	testBlockNumber = big.NewInt(7)
	tests           = map[string]testutils.PrecompileTest{
		"set config from no role fails": {
			Caller:     allowlist.TestNoRoleAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotChangeFee.Error(),
		},
		"set config from enabled address": {
			Caller:     allowlist.TestEnabledAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, state contract.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				require.Equal(t, testFeeConfig, feeConfig)
			},
		},
		"set config from manager succeeds": {
			Caller:     allowlist.TestManagerAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, state contract.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				require.Equal(t, testFeeConfig, feeConfig)
			},
		},
		"set invalid config from enabled address": {
			Caller:     allowlist.TestEnabledAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				feeConfig := testFeeConfig
				feeConfig.MinBlockGasCost = new(big.Int).Mul(feeConfig.MaxBlockGasCost, common.Big2)
				input, err := PackSetFeeConfig(feeConfig)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    false,
			Config: &Config{
				InitialFeeConfig: &testFeeConfig,
			},
			ExpectedErr: "cannot be greater than maxBlockGasCost",
			AfterHook: func(t testing.TB, state contract.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				require.Equal(t, testFeeConfig, feeConfig)
			},
		},
		"set config from admin address": {
			Caller:     allowlist.TestAdminAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			SetupBlockContext: func(mbc *contract.MockBlockContext) {
				mbc.EXPECT().Number().Return(testBlockNumber).AnyTimes()
			},
			AfterHook: func(t testing.TB, state contract.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				require.Equal(t, testFeeConfig, feeConfig)
				lastChangedAt := GetFeeConfigLastChangedAt(state)
				require.EqualValues(t, testBlockNumber, lastChangedAt)
			},
		},
		"get fee config from non-enabled address": {
			Caller: allowlist.TestNoRoleAddr,
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				blockContext := contract.NewMockBlockContext(gomock.NewController(t))
				blockContext.EXPECT().Number().Return(big.NewInt(6)).Times(1)
				allowlist.SetDefaultRoles(Module.Address)(t, state)
				err := StoreFeeConfig(state, testFeeConfig, blockContext)
				require.NoError(t, err)
			},
			Input:       PackGetFeeConfigInput(),
			SuppliedGas: GetFeeConfigGasCost,
			ReadOnly:    true,
			ExpectedRes: func() []byte {
				res, err := PackFeeConfig(testFeeConfig)
				if err != nil {
					panic(err)
				}
				return res
			}(),
			AfterHook: func(t testing.TB, state contract.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				lastChangedAt := GetFeeConfigLastChangedAt(state)
				require.Equal(t, testFeeConfig, feeConfig)
				require.EqualValues(t, big.NewInt(6), lastChangedAt)
			},
		},
		"get initial fee config": {
			Caller:      allowlist.TestNoRoleAddr,
			BeforeHook:  allowlist.SetDefaultRoles(Module.Address),
			Input:       PackGetFeeConfigInput(),
			SuppliedGas: GetFeeConfigGasCost,
			Config: &Config{
				InitialFeeConfig: &testFeeConfig,
			},
			ReadOnly: true,
			ExpectedRes: func() []byte {
				res, err := PackFeeConfig(testFeeConfig)
				if err != nil {
					panic(err)
				}
				return res
			}(),
			SetupBlockContext: func(mbc *contract.MockBlockContext) {
				mbc.EXPECT().Number().Return(testBlockNumber)
			},
			AfterHook: func(t testing.TB, state contract.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				lastChangedAt := GetFeeConfigLastChangedAt(state)
				require.Equal(t, testFeeConfig, feeConfig)
				require.EqualValues(t, testBlockNumber, lastChangedAt)
			},
		},
		"get last changed at from non-enabled address": {
			Caller: allowlist.TestNoRoleAddr,
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				blockContext := contract.NewMockBlockContext(gomock.NewController(t))
				blockContext.EXPECT().Number().Return(testBlockNumber).Times(1)
				allowlist.SetDefaultRoles(Module.Address)(t, state)
				err := StoreFeeConfig(state, testFeeConfig, blockContext)
				require.NoError(t, err)
			},
			Input:       PackGetLastChangedAtInput(),
			SuppliedGas: GetLastChangedAtGasCost,
			ReadOnly:    true,
			ExpectedRes: common.BigToHash(testBlockNumber).Bytes(),
			AfterHook: func(t testing.TB, state contract.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				lastChangedAt := GetFeeConfigLastChangedAt(state)
				require.Equal(t, testFeeConfig, feeConfig)
				require.Equal(t, testBlockNumber, lastChangedAt)
			},
		},
		"readOnly setFeeConfig with noRole fails": {
			Caller:     allowlist.TestNoRoleAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly setFeeConfig with allow role fails": {
			Caller:     allowlist.TestEnabledAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly setFeeConfig with admin role fails": {
			Caller:     allowlist.TestAdminAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"insufficient gas setFeeConfig from admin": {
			Caller:     allowlist.TestAdminAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetFeeConfigGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
	}
)

func TestFeeManager(t *testing.T) {
	allowlist.RunPrecompileWithAllowListTests(t, Module, state.NewTestStateDB, tests)
}

func BenchmarkFeeManager(b *testing.B) {
	allowlist.BenchPrecompileWithAllowList(b, Module, state.NewTestStateDB, tests)
}
func TestGetFeeConfig(t *testing.T) {
	// Compare PackGetFeeConfigV2 vs PackGetFeeConfig
	// to see if they are equivalent

	input := PackGetFeeConfigInput()

	input2, err := PackGetFeeConfigV2()
	require.NoError(t, err)

	require.Equal(t, input, input2)
}

func TestGetFeeConfigOutput(t *testing.T) {
	// Compare PackFeeConfigV2 vs PackFeeConfig
	// to see if they are equivalent

	for i := 0; i < 1000; i++ {
		feeConfig := commontype.FeeConfig{
			GasLimit:        big.NewInt(rand.Int63()),
			TargetBlockRate: rand.Uint64(),

			MinBaseFee:               big.NewInt(rand.Int63()),
			TargetGas:                big.NewInt(rand.Int63()),
			BaseFeeChangeDenominator: big.NewInt(rand.Int63()),

			MinBlockGasCost:  big.NewInt(rand.Int63()),
			MaxBlockGasCost:  big.NewInt(rand.Int63()),
			BlockGasCostStep: big.NewInt(rand.Int63()),
		}

		testGetFeeConfigOutput(t, feeConfig)
	}
	// Some edge cases
	testGetFeeConfigOutput(t, testFeeConfig)
	// These should panic
	require.Panics(t, func() {
		_, _ = PackGetFeeConfigOutputV2(commontype.FeeConfig{})
	})
	require.Panics(t, func() {
		_, _ = PackFeeConfig(commontype.FeeConfig{})
	})

	// These should
	unpacked, err := UnpackGetFeeConfigOutputV2([]byte{})
	require.Error(t, err)

	unpacked2, err := UnpackFeeConfigInput([]byte{})
	require.ErrorIs(t, err, ErrInvalidLen)
	require.Equal(t, unpacked, unpacked2)

	// Test for extra padded bytes
	input, err := PackGetFeeConfigOutputV2(testFeeConfig)
	require.NoError(t, err)
	// exclude 4 bytes for function selector
	input = input[4:]
	// add extra padded bytes
	input = append(input, make([]byte, 32)...)
	_, err = UnpackFeeConfigInput(input)
	require.ErrorIs(t, err, ErrInvalidLen)

	_, err = UnpackGetFeeConfigOutputV2(input)
	require.Error(t, err)
}

func testGetFeeConfigOutput(t *testing.T, feeConfig commontype.FeeConfig) {
	t.Helper()
	t.Run(fmt.Sprintf("TestGetFeeConfigOutput, feeConfig %v", feeConfig), func(t *testing.T) {
		// Test PackGetFeeConfigOutputV2, UnpackGetFeeConfigOutputV2
		input, err := PackGetFeeConfigOutputV2(feeConfig)
		require.NoError(t, err)

		unpacked, err := UnpackGetFeeConfigOutputV2(input)
		require.NoError(t, err)

		require.True(t, feeConfig.Equal(&unpacked), "not equal: feeConfig %v, unpacked %v", feeConfig, unpacked)

		// Test PackGetFeeConfigOutput, UnpackGetFeeConfigOutput
		input, err = PackFeeConfig(feeConfig)
		require.NoError(t, err)

		unpacked, err = UnpackFeeConfigInput(input)
		require.NoError(t, err)

		require.True(t, feeConfig.Equal(&unpacked), "not equal: feeConfig %v, unpacked %v", feeConfig, unpacked)

		// // now mix and match
		// Test PackGetFeeConfigOutput, PackGetFeeConfigOutputV2
		input, err = PackGetFeeConfigOutputV2(feeConfig)
		require.NoError(t, err)
		input2, err := PackFeeConfig(feeConfig)
		require.NoError(t, err)
		require.Equal(t, input, input2)

		// // Test UnpackGetFeeConfigOutput, UnpackGetFeeConfigOutputV2
		unpacked, err = UnpackGetFeeConfigOutputV2(input2)
		require.NoError(t, err)
		unpacked2, err := UnpackFeeConfigInput(input)
		require.NoError(t, err)
		require.True(t, unpacked.Equal(&unpacked2), "not equal: unpacked %v, unpacked2 %v", unpacked, unpacked2)
	})
}

func TestGetLastChangedAtInput(t *testing.T) {
	// Compare PackGetFeeConfigLastChangedAtV2 vs PackGetLastChangedAtInput
	// to see if they are equivalent

	input := PackGetLastChangedAtInput()

	input2, err := PackGetFeeConfigLastChangedAtV2()
	require.NoError(t, err)

	require.Equal(t, input, input2)
}

func TestGetLastChangedAtOutput(t *testing.T) {
	// Compare PackGetFeeConfigLastChangedAtOutputV2 vs PackGetLastChangedAtOutput
	// to see if they are equivalent

	for i := 0; i < 1000; i++ {
		lastChangedAt := big.NewInt(rand.Int63())
		testGetLastChangedAtOutput(t, lastChangedAt)
	}
	// Some edge cases
	testGetLastChangedAtOutput(t, big.NewInt(0))
	testGetLastChangedAtOutput(t, big.NewInt(1))
	testGetLastChangedAtOutput(t, big.NewInt(2))
	testGetLastChangedAtOutput(t, math.MaxBig256)
	testGetLastChangedAtOutput(t, math.MaxBig256.Sub(math.MaxBig256, common.Big1))
	testGetLastChangedAtOutput(t, math.MaxBig256.Add(math.MaxBig256, common.Big1))
}

func testGetLastChangedAtOutput(t *testing.T, lastChangedAt *big.Int) {
	t.Helper()
	t.Run(fmt.Sprintf("TestGetLastChangedAtOutput, lastChangedAt %v", lastChangedAt), func(t *testing.T) {
		// Test PackGetFeeConfigLastChangedAtOutputV2, UnpackGetFeeConfigLastChangedAtOutputV2
		input, err := PackGetFeeConfigLastChangedAtOutputV2(lastChangedAt)
		require.NoError(t, err)

		unpacked, err := UnpackGetFeeConfigLastChangedAtOutputV2(input)
		require.NoError(t, err)

		require.Zero(t, lastChangedAt.Cmp(unpacked), "not equal: lastChangedAt %v, unpacked %v", lastChangedAt, unpacked)

		// Test PackGetLastChangedAtOutput, UnpackGetLastChangedAtOutput
		input = common.BigToHash(lastChangedAt).Bytes()

		unpacked = common.BytesToHash(input).Big()
		require.NoError(t, err)

		require.Zero(t, lastChangedAt.Cmp(unpacked), "not equal: lastChangedAt %v, unpacked %v", lastChangedAt)

		// now mix and match
		// Test PackGetLastChangedAtOutput, PackGetFeeConfigLastChangedAtOutputV2
		input = common.BigToHash(lastChangedAt).Bytes()
		require.NoError(t, err)
		input2, err := PackGetFeeConfigLastChangedAtOutputV2(lastChangedAt)
		require.NoError(t, err)
		require.Equal(t, input, input2)

		// Test UnpackGetLastChangedAtOutput, UnpackGetFeeConfigLastChangedAtOutputV2
		unpacked = common.BytesToHash(input).Big()
		require.NoError(t, err)
		unpacked2, err := UnpackGetFeeConfigLastChangedAtOutputV2(input)
		require.NoError(t, err)
		require.EqualValues(t, unpacked, unpacked2)
	})
}

func TestPackSetFeeConfigInput(t *testing.T) {
	// Compare PackSetFeeConfigV2 vs PackSetFeeConfig
	// to see if they are equivalent
	for i := 0; i < 1000; i++ {
		feeConfig := commontype.FeeConfig{
			GasLimit:        big.NewInt(rand.Int63()),
			TargetBlockRate: rand.Uint64(),

			MinBaseFee:               big.NewInt(rand.Int63()),
			TargetGas:                big.NewInt(rand.Int63()),
			BaseFeeChangeDenominator: big.NewInt(rand.Int63()),

			MinBlockGasCost:  big.NewInt(rand.Int63()),
			MaxBlockGasCost:  big.NewInt(rand.Int63()),
			BlockGasCostStep: big.NewInt(rand.Int63()),
		}

		testPackSetFeeConfigInput(t, feeConfig)
	}
	// Some edge cases
	// Some edge cases
	testPackSetFeeConfigInput(t, testFeeConfig)
	// These should panic
	require.Panics(t, func() {
		_, _ = PackSetFeeConfigV2(commontype.FeeConfig{})
	})
	require.Panics(t, func() {
		_, _ = PackSetFeeConfig(commontype.FeeConfig{})
	})

	// These should err
	_, err := UnpackSetFeeConfigInputV2([]byte{123}, true)
	require.ErrorIs(t, err, ErrInvalidLen)

	_, err = UnpackSetFeeConfigInputV2([]byte{123}, false)
	require.ErrorContains(t, err, "abi: improperly formatted input")

	_, err = UnpackFeeConfigInput([]byte{123})
	require.ErrorIs(t, err, ErrInvalidLen)

	// Test for extra padded bytes
	input, err := PackSetFeeConfigV2(testFeeConfig)
	require.NoError(t, err)
	// exclude 4 bytes for function selector
	input = input[4:]
	// add extra padded bytes
	input = append(input, make([]byte, 32)...)
	_, err = UnpackFeeConfigInput(input)
	require.ErrorIs(t, err, ErrInvalidLen)

	_, err = UnpackSetFeeConfigInputV2(input, true)
	require.ErrorIs(t, err, ErrInvalidLen)

	unpacked, err := UnpackSetFeeConfigInputV2(input, false)
	require.NoError(t, err)
	require.True(t, testFeeConfig.Equal(&unpacked))
}

func testPackSetFeeConfigInput(t *testing.T, feeConfig commontype.FeeConfig) {
	t.Helper()
	t.Run(fmt.Sprintf("TestPackSetFeeConfigInput, feeConfig %v", feeConfig), func(t *testing.T) {
		// Test PackSetFeeConfigV2, UnpackSetFeeConfigInputV2
		input, err := PackSetFeeConfigV2(feeConfig)
		require.NoError(t, err)
		// exclude 4 bytes for function selector
		input = input[4:]

		unpacked, err := UnpackSetFeeConfigInputV2(input, true)
		require.NoError(t, err)

		require.True(t, feeConfig.Equal(&unpacked), "not equal: feeConfig %v, unpacked %v", feeConfig, unpacked)

		// Test PackSetFeeConfig, UnpackFeeConfigInput
		input, err = PackSetFeeConfig(feeConfig)
		require.NoError(t, err)
		// exclude 4 bytes for function selector
		input = input[4:]

		unpacked, err = UnpackFeeConfigInput(input)
		require.NoError(t, err)

		require.True(t, feeConfig.Equal(&unpacked), "not equal: feeConfig %v, unpacked %v", feeConfig)

		// now mix and match
		// Test PackSetFeeConfig, PackSetFeeConfigV2
		input, err = PackSetFeeConfig(feeConfig)
		require.NoError(t, err)
		input2, err := PackSetFeeConfigV2(feeConfig)
		require.NoError(t, err)
		require.Equal(t, input, input2)
		// exclude 4 bytes for function selector
		input = input[4:]
		input2 = input2[4:]

		// Test UnpackSetFeeConfigInputV2, UnpackFeeConfigInput
		unpacked, err = UnpackSetFeeConfigInputV2(input2, true)
		require.NoError(t, err)
		unpacked2, err := UnpackFeeConfigInput(input)
		require.NoError(t, err)
		require.EqualValues(t, unpacked, unpacked2)
	})
}
