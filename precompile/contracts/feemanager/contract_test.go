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
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
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
				mbc.EXPECT().Timestamp().Return(uint64(0)).AnyTimes()
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
			InputFn: func(t testing.TB) []byte {
				input, err := PackGetFeeConfig()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: GetFeeConfigGasCost,
			ReadOnly:    true,
			ExpectedRes: func() []byte {
				res, err := PackGetFeeConfigOutput(testFeeConfig)
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
			Caller:     allowlist.TestNoRoleAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackGetFeeConfig()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: GetFeeConfigGasCost,
			Config: &Config{
				InitialFeeConfig: &testFeeConfig,
			},
			ReadOnly: true,
			ExpectedRes: func() []byte {
				res, err := PackGetFeeConfigOutput(testFeeConfig)
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
			InputFn: func(t testing.TB) []byte {
				input, err := PackGetFeeConfigLastChangedAt()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: GetLastChangedAtGasCost,
			ReadOnly:    true,
			ExpectedRes: func() []byte {
				res, err := PackGetFeeConfigLastChangedAtOutput(testBlockNumber)
				if err != nil {
					panic(err)
				}
				return res
			}(),
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
		"set config with extra padded bytes should fail before DUpgrade": {
			Caller:     allowlist.TestEnabledAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				input = append(input, make([]byte, 32)...)
				return input
			},
			ChainConfigFn: func(t testing.TB) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(gomock.NewController(t))
				config.EXPECT().IsDUpgrade(gomock.Any()).Return(false).AnyTimes()
				return config
			},
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrInvalidLen.Error(),
			SetupBlockContext: func(mbc *contract.MockBlockContext) {
				mbc.EXPECT().Number().Return(testBlockNumber).AnyTimes()
				mbc.EXPECT().Timestamp().Return(uint64(0)).AnyTimes()
			},
		},
		"set config with extra padded bytes should succeed with DUpgrade": {
			Caller:     allowlist.TestEnabledAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := PackSetFeeConfig(testFeeConfig)
				require.NoError(t, err)

				input = append(input, make([]byte, 32)...)
				return input
			},
			ChainConfigFn: func(t testing.TB) precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(gomock.NewController(t))
				config.EXPECT().IsDUpgrade(gomock.Any()).Return(true).AnyTimes()
				return config
			},
			SuppliedGas: SetFeeConfigGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			SetupBlockContext: func(mbc *contract.MockBlockContext) {
				mbc.EXPECT().Number().Return(testBlockNumber).AnyTimes()
				mbc.EXPECT().Timestamp().Return(uint64(0)).AnyTimes()
			},
			AfterHook: func(t testing.TB, state contract.StateDB) {
				feeConfig := GetStoredFeeConfig(state)
				require.Equal(t, testFeeConfig, feeConfig)
				lastChangedAt := GetFeeConfigLastChangedAt(state)
				require.EqualValues(t, testBlockNumber, lastChangedAt)
			},
		},
	}
)

func TestFeeManager(t *testing.T) {
	allowlist.RunPrecompileWithAllowListTests(t, Module, state.NewTestStateDB, tests)
}

func BenchmarkFeeManager(b *testing.B) {
	allowlist.BenchPrecompileWithAllowList(b, Module, state.NewTestStateDB, tests)
}

func TestPackUnpackGetFeeConfigOutput(t *testing.T) {
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
	// should panic
	require.Panics(t, func() {
		_, _ = PackGetFeeConfigOutput(commontype.FeeConfig{})
	})

	_, err := UnpackGetFeeConfigOutput([]byte{})
	require.Error(t, err)
}

func testGetFeeConfigOutput(t *testing.T, feeConfig commontype.FeeConfig) {
	t.Helper()
	t.Run(fmt.Sprintf("TestGetFeeConfigOutput, feeConfig %v", feeConfig), func(t *testing.T) {
		input, err := PackGetFeeConfigOutput(feeConfig)
		require.NoError(t, err)

		unpacked, err := UnpackGetFeeConfigOutput(input)
		require.NoError(t, err)

		require.True(t, feeConfig.Equal(&unpacked), "not equal: feeConfig %v, unpacked %v", feeConfig, unpacked)
	})
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
		input, err := PackGetFeeConfigLastChangedAtOutput(lastChangedAt)
		require.NoError(t, err)

		unpacked, err := UnpackGetFeeConfigLastChangedAtOutput(input)
		require.NoError(t, err)

		require.Zero(t, lastChangedAt.Cmp(unpacked), "not equal: lastChangedAt %v, unpacked %v", lastChangedAt, unpacked)
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
		_, _ = PackSetFeeConfig(commontype.FeeConfig{})
	})

	// These should err
	_, err := UnpackSetFeeConfigInput([]byte{123}, false)
	require.ErrorIs(t, err, ErrInvalidLen)

	_, err = UnpackSetFeeConfigInput([]byte{123}, true)
	require.ErrorContains(t, err, "abi: improperly formatted input")

	// Test for extra padded bytes
	input, err := PackSetFeeConfig(testFeeConfig)
	require.NoError(t, err)
	// exclude 4 bytes for function selector
	input = input[4:]
	// add extra padded bytes
	input = append(input, make([]byte, 32)...)
	_, err = UnpackSetFeeConfigInput(input, false)
	require.ErrorIs(t, err, ErrInvalidLen)

	unpacked, err := UnpackSetFeeConfigInput(input, true)
	require.NoError(t, err)
	require.True(t, testFeeConfig.Equal(&unpacked))
}

func TestSignatures(t *testing.T) {
	setFeeConfigSignature := contract.CalculateFunctionSelector("setFeeConfig(uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)")
	abiSetFeeConfig := FeeManagerABI.Methods["setFeeConfig"]
	require.Equal(t, setFeeConfigSignature, abiSetFeeConfig.ID)

	getFeeConfigSignature := contract.CalculateFunctionSelector("getFeeConfig()")
	abiGetFeeConfig := FeeManagerABI.Methods["getFeeConfig"]
	require.Equal(t, getFeeConfigSignature, abiGetFeeConfig.ID)

	getFeeConfigLastChangedAtSignature := contract.CalculateFunctionSelector("getFeeConfigLastChangedAt()")
	abiGetFeeConfigLastChangedAt := FeeManagerABI.Methods["getFeeConfigLastChangedAt"]
	require.Equal(t, getFeeConfigLastChangedAtSignature, abiGetFeeConfigLastChangedAt.ID)
}

func testPackSetFeeConfigInput(t *testing.T, feeConfig commontype.FeeConfig) {
	t.Helper()
	t.Run(fmt.Sprintf("TestPackSetFeeConfigInput, feeConfig %v", feeConfig), func(t *testing.T) {
		input, err := PackSetFeeConfig(feeConfig)
		require.NoError(t, err)
		// exclude 4 bytes for function selector
		input = input[4:]

		unpacked, err := UnpackSetFeeConfigInput(input, true)
		require.NoError(t, err)

		require.True(t, feeConfig.Equal(&unpacked), "not equal: feeConfig %v, unpacked %v", feeConfig, unpacked)
	})
}
