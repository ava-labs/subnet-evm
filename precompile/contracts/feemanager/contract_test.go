// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package feemanager

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/accounts/abi"
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
	setFeeConfigSignature              = contract.CalculateFunctionSelector("setFeeConfig(uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)")
	getFeeConfigSignature              = contract.CalculateFunctionSelector("getFeeConfig()")
	getFeeConfigLastChangedAtSignature = contract.CalculateFunctionSelector("getFeeConfigLastChangedAt()")

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

func FuzzPackGetFeeConfigOutputEqualTest(f *testing.F) {
	f.Add([]byte{}, uint64(0))
	f.Add(big.NewInt(0).Bytes(), uint64(0))
	f.Add(big.NewInt(1).Bytes(), uint64(math.MaxUint64))
	f.Add(math.MaxBig256.Bytes(), uint64(0))
	f.Add(math.MaxBig256.Sub(math.MaxBig256, common.Big1).Bytes(), uint64(0))
	f.Add(math.MaxBig256.Add(math.MaxBig256, common.Big1).Bytes(), uint64(0))
	f.Fuzz(func(t *testing.T, bigIntBytes []byte, blockRate uint64) {
		bigIntVal := new(big.Int).SetBytes(bigIntBytes)
		feeConfig := commontype.FeeConfig{
			GasLimit:                 bigIntVal,
			TargetBlockRate:          blockRate,
			MinBaseFee:               bigIntVal,
			TargetGas:                bigIntVal,
			BaseFeeChangeDenominator: bigIntVal,
			MinBlockGasCost:          bigIntVal,
			MaxBlockGasCost:          bigIntVal,
			BlockGasCostStep:         bigIntVal,
		}
		doCheckOutputs := true
		// we can only check if outputs are correct if the value is less than MaxUint256
		// otherwise the value will be truncated when packed,
		// and thus unpacked output will not be equal to the value
		if bigIntVal.Cmp(abi.MaxUint256) > 0 {
			doCheckOutputs = false
		}
		testOldPackGetFeeConfigOutputEqual(t, feeConfig, doCheckOutputs)
	})
}

func TestPackUnpackGetFeeConfigOutputEdgeCases(t *testing.T) {
	testOldPackGetFeeConfigOutputEqual(t, testFeeConfig, true)
	// These should panic
	require.Panics(t, func() {
		_, _ = OldPackFeeConfig(commontype.FeeConfig{})
	})
	require.Panics(t, func() {
		_, _ = PackGetFeeConfigOutput(commontype.FeeConfig{})
	})

	unpacked, err := OldUnpackFeeConfig([]byte{})
	require.ErrorIs(t, err, ErrInvalidLen)
	unpacked2, err := UnpackGetFeeConfigOutput([]byte{}, false)
	require.ErrorIs(t, err, ErrInvalidLen)
	require.Equal(t, unpacked, unpacked2)

	_, err = UnpackGetFeeConfigOutput([]byte{}, true)
	require.Error(t, err)

	// Test for extra padded bytes
	input, err := PackGetFeeConfigOutput(testFeeConfig)
	require.NoError(t, err)
	// add extra padded bytes
	input = append(input, make([]byte, 32)...)
	_, err = OldUnpackFeeConfig(input)
	require.ErrorIs(t, err, ErrInvalidLen)
	_, err = UnpackGetFeeConfigOutput([]byte{}, false)
	require.ErrorIs(t, err, ErrInvalidLen)

	_, err = UnpackGetFeeConfigOutput(input, true)
	require.NoError(t, err)

	// now it's now divisible by 32
	input = append(input, make([]byte, 1)...)
	_, err = UnpackGetFeeConfigOutput(input, true)
	require.Error(t, err)
}

func TestGetFeeConfig(t *testing.T) {
	// Compare OldPackGetFeeConfigInput vs PackGetFeeConfig
	// to see if they are equivalent
	input := OldPackGetFeeConfigInput()

	input2, err := PackGetFeeConfig()
	require.NoError(t, err)

	require.Equal(t, input, input2)
}

func TestGetLastChangedAtInput(t *testing.T) {
	// Compare OldPackGetFeeConfigInput vs PackGetFeeConfigLastChangedAt
	// to see if they are equivalent

	input := OldPackGetLastChangedAtInput()

	input2, err := PackGetFeeConfigLastChangedAt()
	require.NoError(t, err)

	require.Equal(t, input, input2)
}

func FuzzPackGetLastChangedAtOutput(f *testing.F) {
	f.Add([]byte{})
	f.Add(big.NewInt(0).Bytes())
	f.Add(big.NewInt(1).Bytes())
	f.Add(math.MaxBig256.Bytes())
	f.Add(math.MaxBig256.Sub(math.MaxBig256, common.Big1).Bytes())
	f.Add(math.MaxBig256.Add(math.MaxBig256, common.Big1).Bytes())
	f.Fuzz(func(t *testing.T, bigIntBytes []byte) {
		bigIntVal := new(big.Int).SetBytes(bigIntBytes)
		doCheckOutputs := true
		// we can only check if outputs are correct if the value is less than MaxUint256
		// otherwise the value will be truncated when packed,
		// and thus unpacked output will not be equal to the value
		if bigIntVal.Cmp(abi.MaxUint256) > 0 {
			doCheckOutputs = false
		}
		testOldPackGetLastChangedAtOutputEqual(t, bigIntVal, doCheckOutputs)
	})
}

func testOldPackGetFeeConfigOutputEqual(t *testing.T, feeConfig commontype.FeeConfig, checkOutputs bool) {
	t.Helper()
	t.Run(fmt.Sprintf("TestGetFeeConfigOutput, feeConfig %v", feeConfig), func(t *testing.T) {
		input, err := OldPackFeeConfig(feeConfig)
		input2, err2 := PackGetFeeConfigOutput(feeConfig)
		if err != nil {
			require.ErrorContains(t, err2, err.Error())
			return
		}
		require.NoError(t, err2)
		require.Equal(t, input, input2)

		config, err := OldUnpackFeeConfig(input)
		unpacked, err2 := UnpackGetFeeConfigOutput(input, false)
		if err != nil {
			require.ErrorContains(t, err2, err.Error())
			return
		}
		require.NoError(t, err2)
		require.True(t, config.Equal(&unpacked), "not equal: config %v, unpacked %v", feeConfig, unpacked)
		if checkOutputs {
			require.True(t, feeConfig.Equal(&unpacked), "not equal: feeConfig %v, unpacked %v", feeConfig, unpacked)
		}
	})
}

func testOldPackGetLastChangedAtOutputEqual(t *testing.T, blockNumber *big.Int, checkOutputs bool) {
	t.Helper()
	t.Run(fmt.Sprintf("TestGetLastChangedAtOutput, blockNumber %v", blockNumber), func(t *testing.T) {
		input := OldPackGetLastChangedAtOutput(blockNumber)
		input2, err2 := PackGetFeeConfigLastChangedAtOutput(blockNumber)
		require.NoError(t, err2)
		require.Equal(t, input, input2)

		value, err := OldUnpackGetLastChangedAtOutput(input)
		unpacked, err2 := UnpackGetFeeConfigLastChangedAtOutput(input)
		if err != nil {
			require.ErrorContains(t, err2, err.Error())
			return
		}
		require.NoError(t, err2)
		require.True(t, value.Cmp(unpacked) == 0, "not equal: value %v, unpacked %v", value, unpacked)
		if checkOutputs {
			require.True(t, blockNumber.Cmp(unpacked) == 0, "not equal: blockNumber %v, unpacked %v", blockNumber, unpacked)
		}
	})
}

func FuzzPackSetFeeConfigEqualTest(f *testing.F) {
	f.Add([]byte{}, uint64(0))
	f.Add(big.NewInt(0).Bytes(), uint64(0))
	f.Add(big.NewInt(1).Bytes(), uint64(math.MaxUint64))
	f.Add(math.MaxBig256.Bytes(), uint64(0))
	f.Add(math.MaxBig256.Sub(math.MaxBig256, common.Big1).Bytes(), uint64(0))
	f.Add(math.MaxBig256.Add(math.MaxBig256, common.Big1).Bytes(), uint64(0))
	f.Fuzz(func(t *testing.T, bigIntBytes []byte, blockRate uint64) {
		bigIntVal := new(big.Int).SetBytes(bigIntBytes)
		feeConfig := commontype.FeeConfig{
			GasLimit:                 bigIntVal,
			TargetBlockRate:          blockRate,
			MinBaseFee:               bigIntVal,
			TargetGas:                bigIntVal,
			BaseFeeChangeDenominator: bigIntVal,
			MinBlockGasCost:          bigIntVal,
			MaxBlockGasCost:          bigIntVal,
			BlockGasCostStep:         bigIntVal,
		}
		doCheckOutputs := true
		// we can only check if outputs are correct if the value is less than MaxUint256
		// otherwise the value will be truncated when packed,
		// and thus unpacked output will not be equal to the value
		if bigIntVal.Cmp(abi.MaxUint256) > 0 {
			doCheckOutputs = false
		}
		testOldPackSetFeeConfigInputEqual(t, feeConfig, doCheckOutputs)
	})
}

func TestPackSetFeeConfigInputEdgeCases(t *testing.T) {
	// Some edge cases
	testOldPackSetFeeConfigInputEqual(t, testFeeConfig, true)
	// These should panic
	require.Panics(t, func() {
		_, _ = OldPackSetFeeConfig(commontype.FeeConfig{})
	})
	require.Panics(t, func() {
		_, _ = PackSetFeeConfig(commontype.FeeConfig{})
	})
	// These should err
	_, err := UnpackSetFeeConfigInput([]byte{123}, false)
	require.ErrorIs(t, err, ErrInvalidLen)

	_, err = UnpackSetFeeConfigInput([]byte{123}, true)
	require.ErrorContains(t, err, "abi: improperly formatted input")

	_, err = OldUnpackFeeConfig([]byte{123})
	require.ErrorIs(t, err, ErrInvalidLen)

	// Test for extra padded bytes
	input, err := PackSetFeeConfig(testFeeConfig)
	require.NoError(t, err)
	// exclude 4 bytes for function selector
	input = input[4:]
	// add extra padded bytes
	input = append(input, make([]byte, 32)...)
	_, err = OldUnpackFeeConfig(input)
	require.ErrorIs(t, err, ErrInvalidLen)
	_, err = UnpackSetFeeConfigInput(input, false)
	require.ErrorIs(t, err, ErrInvalidLen)

	unpacked, err := UnpackSetFeeConfigInput(input, true)
	require.NoError(t, err)
	require.True(t, testFeeConfig.Equal(&unpacked))
}

func TestFunctionSignatures(t *testing.T) {
	abiSetFeeConfig := FeeManagerABI.Methods["setFeeConfig"]
	require.Equal(t, setFeeConfigSignature, abiSetFeeConfig.ID)

	abiGetFeeConfig := FeeManagerABI.Methods["getFeeConfig"]
	require.Equal(t, getFeeConfigSignature, abiGetFeeConfig.ID)

	abiGetFeeConfigLastChangedAt := FeeManagerABI.Methods["getFeeConfigLastChangedAt"]
	require.Equal(t, getFeeConfigLastChangedAtSignature, abiGetFeeConfigLastChangedAt.ID)
}

func testOldPackSetFeeConfigInputEqual(t *testing.T, feeConfig commontype.FeeConfig, checkOutputs bool) {
	t.Helper()
	t.Run(fmt.Sprintf("TestSetFeeConfigInput, feeConfig %v", feeConfig), func(t *testing.T) {
		input, err := OldPackSetFeeConfig(feeConfig)
		input2, err2 := PackSetFeeConfig(feeConfig)
		if err != nil {
			require.ErrorContains(t, err2, err.Error())
			return
		}
		require.NoError(t, err2)
		require.Equal(t, input, input2)

		value, err := OldUnpackFeeConfig(input)
		unpacked, err2 := UnpackSetFeeConfigInput(input, false)
		if err != nil {
			require.ErrorContains(t, err2, err.Error())
			return
		}
		require.NoError(t, err2)
		require.True(t, value.Equal(&unpacked), "not equal: value %v, unpacked %v", value, unpacked)
		if checkOutputs {
			require.True(t, feeConfig.Equal(&unpacked), "not equal: feeConfig %v, unpacked %v", feeConfig, unpacked)
		}
	})
}

func OldPackFeeConfig(feeConfig commontype.FeeConfig) ([]byte, error) {
	return packFeeConfigHelper(feeConfig, false)
}

func OldUnpackFeeConfig(input []byte) (commontype.FeeConfig, error) {
	if len(input) != feeConfigInputLen {
		return commontype.FeeConfig{}, fmt.Errorf("%w: %d", ErrInvalidLen, len(input))
	}
	feeConfig := commontype.FeeConfig{}
	for i := minFeeConfigFieldKey; i <= numFeeConfigField; i++ {
		listIndex := i - 1
		packedElement := contract.PackedHash(input, listIndex)
		switch i {
		case gasLimitKey:
			feeConfig.GasLimit = new(big.Int).SetBytes(packedElement)
		case targetBlockRateKey:
			feeConfig.TargetBlockRate = new(big.Int).SetBytes(packedElement).Uint64()
		case minBaseFeeKey:
			feeConfig.MinBaseFee = new(big.Int).SetBytes(packedElement)
		case targetGasKey:
			feeConfig.TargetGas = new(big.Int).SetBytes(packedElement)
		case baseFeeChangeDenominatorKey:
			feeConfig.BaseFeeChangeDenominator = new(big.Int).SetBytes(packedElement)
		case minBlockGasCostKey:
			feeConfig.MinBlockGasCost = new(big.Int).SetBytes(packedElement)
		case maxBlockGasCostKey:
			feeConfig.MaxBlockGasCost = new(big.Int).SetBytes(packedElement)
		case blockGasCostStepKey:
			feeConfig.BlockGasCostStep = new(big.Int).SetBytes(packedElement)
		default:
			// This should never encounter an unknown fee config key
			panic(fmt.Sprintf("unknown fee config key: %d", i))
		}
	}
	return feeConfig, nil
}

func packFeeConfigHelper(feeConfig commontype.FeeConfig, useSelector bool) ([]byte, error) {
	hashes := []common.Hash{
		common.BigToHash(feeConfig.GasLimit),
		common.BigToHash(new(big.Int).SetUint64(feeConfig.TargetBlockRate)),
		common.BigToHash(feeConfig.MinBaseFee),
		common.BigToHash(feeConfig.TargetGas),
		common.BigToHash(feeConfig.BaseFeeChangeDenominator),
		common.BigToHash(feeConfig.MinBlockGasCost),
		common.BigToHash(feeConfig.MaxBlockGasCost),
		common.BigToHash(feeConfig.BlockGasCostStep),
	}

	if useSelector {
		res := make([]byte, len(setFeeConfigSignature)+feeConfigInputLen)
		err := contract.PackOrderedHashesWithSelector(res, setFeeConfigSignature, hashes)
		return res, err
	}

	res := make([]byte, len(hashes)*common.HashLength)
	err := contract.PackOrderedHashes(res, hashes)
	return res, err
}

// PackGetFeeConfigInput packs the getFeeConfig signature
func OldPackGetFeeConfigInput() []byte {
	return getFeeConfigSignature
}

// PackGetLastChangedAtInput packs the getFeeConfigLastChangedAt signature
func OldPackGetLastChangedAtInput() []byte {
	return getFeeConfigLastChangedAtSignature
}

func OldPackGetLastChangedAtOutput(lastChangedAt *big.Int) []byte {
	return common.BigToHash(lastChangedAt).Bytes()
}

func OldUnpackGetLastChangedAtOutput(input []byte) (*big.Int, error) {
	return new(big.Int).SetBytes(input), nil
}

func OldPackSetFeeConfig(feeConfig commontype.FeeConfig) ([]byte, error) {
	// function selector (4 bytes) + input(feeConfig)
	return packFeeConfigHelper(feeConfig, true)
}
