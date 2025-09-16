package acp224feemanagertest

import (
	"math/big"
	"testing"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/vm"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/core/extstate"
	"github.com/ava-labs/subnet-evm/precompile/allowlist/allowlisttest"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/contracts/acp224feemanager"
	"github.com/ava-labs/subnet-evm/precompile/precompiletest"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	_ = vm.ErrOutOfGas
	_ = big.NewInt
	_ = common.Big0
	_ = require.New
)

// These tests are run against the precompile contract directly with
// the given input and expected output. They're just a guide to
// help you write your own tests. These tests are for general cases like
// allowlist, readOnly behaviour, and gas cost. You should write your own
// tests for specific cases.
var (
	testBlockNumber             = big.NewInt(7)
	testDefaultUpdatedFeeConfig = commontype.ACP224FeeConfig{
		TargetGas:          big.NewInt(10_000_000),
		MinGasPrice:        common.Big1,
		TimeToFillCapacity: big.NewInt(5),
		TimeToDouble:       big.NewInt(60),
	}
	tests = map[string]precompiletest.PrecompileTest{
		"calling getFeeConfig from NoRole should succeed": {
			Caller: allowlisttest.TestNoRoleAddr,
			BeforeHook: func(t testing.TB, state *extstate.StateDB) {
				blockContext := contract.NewMockBlockContext(gomock.NewController(t))
				blockContext.EXPECT().Number().Return(big.NewInt(6)).Times(1)
				allowlisttest.SetDefaultRoles(acp224feemanager.Module.Address)(t, state)
				require.NoError(t, acp224feemanager.StoreFeeConfig(state, testDefaultUpdatedFeeConfig, blockContext))
			},
			InputFn: func(t testing.TB) []byte {
				input, err := acp224feemanager.PackGetFeeConfig()
				require.NoError(t, err)
				return input
			},
			ExpectedRes: func() []byte {
				packedOutput, err := acp224feemanager.PackGetFeeConfigOutput(testDefaultUpdatedFeeConfig)
				if err != nil {
					panic(err)
				}
				return packedOutput
			}(),
			SuppliedGas: acp224feemanager.GetFeeConfigGasCost,
			ReadOnly:    false,
			ExpectedErr: "",
		},
		"calling getFeeConfig from Enabled should succeed": {
			Caller: allowlisttest.TestEnabledAddr,
			BeforeHook: func(t testing.TB, state *extstate.StateDB) {
				blockContext := contract.NewMockBlockContext(gomock.NewController(t))
				blockContext.EXPECT().Number().Return(big.NewInt(6)).Times(1)
				allowlisttest.SetDefaultRoles(acp224feemanager.Module.Address)(t, state)
				require.NoError(t, acp224feemanager.StoreFeeConfig(state, testDefaultUpdatedFeeConfig, blockContext))
			},
			InputFn: func(t testing.TB) []byte {
				input, err := acp224feemanager.PackGetFeeConfig()
				require.NoError(t, err)
				return input
			},
			ExpectedRes: func() []byte {
				packedOutput, err := acp224feemanager.PackGetFeeConfigOutput(testDefaultUpdatedFeeConfig)
				if err != nil {
					panic(err)
				}
				return packedOutput
			}(),
			SuppliedGas: acp224feemanager.GetFeeConfigGasCost,
			ReadOnly:    false,
			ExpectedErr: "",
		},
		"calling getFeeConfig from Manager should succeed": {
			Caller: allowlisttest.TestManagerAddr,
			BeforeHook: func(t testing.TB, state *extstate.StateDB) {
				blockContext := contract.NewMockBlockContext(gomock.NewController(t))
				blockContext.EXPECT().Number().Return(big.NewInt(6)).Times(1)
				allowlisttest.SetDefaultRoles(acp224feemanager.Module.Address)(t, state)
				require.NoError(t, acp224feemanager.StoreFeeConfig(state, testDefaultUpdatedFeeConfig, blockContext))
			},
			InputFn: func(t testing.TB) []byte {
				input, err := acp224feemanager.PackGetFeeConfig()
				require.NoError(t, err)
				return input
			},
			ExpectedRes: func() []byte {
				packedOutput, err := acp224feemanager.PackGetFeeConfigOutput(testDefaultUpdatedFeeConfig)
				if err != nil {
					panic(err)
				}
				return packedOutput
			}(),
			SuppliedGas: acp224feemanager.GetFeeConfigGasCost,
			ReadOnly:    false,
			ExpectedErr: "",
		},
		"calling getFeeConfig from Admin should succeed": {
			Caller: allowlisttest.TestAdminAddr,
			BeforeHook: func(t testing.TB, state *extstate.StateDB) {
				blockContext := contract.NewMockBlockContext(gomock.NewController(t))
				blockContext.EXPECT().Number().Return(big.NewInt(6)).Times(1)
				allowlisttest.SetDefaultRoles(acp224feemanager.Module.Address)(t, state)
				require.NoError(t, acp224feemanager.StoreFeeConfig(state, testDefaultUpdatedFeeConfig, blockContext))
			},
			InputFn: func(t testing.TB) []byte {
				input, err := acp224feemanager.PackGetFeeConfig()
				require.NoError(t, err)
				return input
			},
			ExpectedRes: func() []byte {
				packedOutput, err := acp224feemanager.PackGetFeeConfigOutput(testDefaultUpdatedFeeConfig)
				if err != nil {
					panic(err)
				}
				return packedOutput
			}(),
			SuppliedGas: acp224feemanager.GetFeeConfigGasCost,
			ReadOnly:    false,
			ExpectedErr: "",
		},
		"insufficient gas for getFeeConfig should fail": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				input, err := acp224feemanager.PackGetFeeConfig()
				require.NoError(t, err)
				return input
			},
			SuppliedGas: acp224feemanager.GetFeeConfigGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vm.ErrOutOfGas.Error(),
		},
		"calling getFeeConfigLastChangedAt from NoRole should succeed": {
			Caller: allowlisttest.TestNoRoleAddr,
			BeforeHook: func(t testing.TB, state *extstate.StateDB) {
				blockContext := contract.NewMockBlockContext(gomock.NewController(t))
				blockContext.EXPECT().Number().Return(testBlockNumber).Times(1)
				allowlisttest.SetDefaultRoles(acp224feemanager.Module.Address)(t, state)
				require.NoError(t, acp224feemanager.StoreFeeConfig(state, testDefaultUpdatedFeeConfig, blockContext))
			},
			InputFn: func(t testing.TB) []byte {
				input, err := acp224feemanager.PackGetFeeConfigLastChangedAt()
				require.NoError(t, err)
				return input
			},
			ExpectedRes: func() []byte {
				packedOutput, err := acp224feemanager.PackGetFeeConfigLastChangedAtOutput(testBlockNumber)
				if err != nil {
					panic(err)
				}
				return packedOutput
			}(),
			AfterHook: func(t testing.TB, state *extstate.StateDB) {
				feeConfig := acp224feemanager.GetStoredFeeConfig(state)
				require.Equal(t, testDefaultUpdatedFeeConfig, feeConfig)
				lastChangedAt := acp224feemanager.GetFeeConfigLastChangedAt(state)
				require.Equal(t, testBlockNumber, lastChangedAt)
			},
			SuppliedGas: acp224feemanager.GetFeeConfigLastChangedAtGasCost,
			ReadOnly:    false,
			ExpectedErr: "",
		},
		"calling getFeeConfigLastChangedAt from Enabled should succeed": {
			Caller: allowlisttest.TestEnabledAddr,
			BeforeHook: func(t testing.TB, state *extstate.StateDB) {
				blockContext := contract.NewMockBlockContext(gomock.NewController(t))
				blockContext.EXPECT().Number().Return(testBlockNumber).Times(1)
				allowlisttest.SetDefaultRoles(acp224feemanager.Module.Address)(t, state)
				require.NoError(t, acp224feemanager.StoreFeeConfig(state, testDefaultUpdatedFeeConfig, blockContext))
			},
			InputFn: func(t testing.TB) []byte {
				input, err := acp224feemanager.PackGetFeeConfigLastChangedAt()
				require.NoError(t, err)
				return input
			},
			ExpectedRes: func() []byte {
				packedOutput, err := acp224feemanager.PackGetFeeConfigLastChangedAtOutput(testBlockNumber)
				if err != nil {
					panic(err)
				}
				return packedOutput
			}(),
			AfterHook: func(t testing.TB, state *extstate.StateDB) {
				feeConfig := acp224feemanager.GetStoredFeeConfig(state)
				require.Equal(t, testDefaultUpdatedFeeConfig, feeConfig)
				lastChangedAt := acp224feemanager.GetFeeConfigLastChangedAt(state)
				require.Equal(t, testBlockNumber, lastChangedAt)
			},
			SuppliedGas: acp224feemanager.GetFeeConfigLastChangedAtGasCost,
			ReadOnly:    false,
			ExpectedErr: "",
		},
		"calling getFeeConfigLastChangedAt from Manager should succeed": {
			Caller: allowlisttest.TestManagerAddr,
			BeforeHook: func(t testing.TB, state *extstate.StateDB) {
				blockContext := contract.NewMockBlockContext(gomock.NewController(t))
				blockContext.EXPECT().Number().Return(testBlockNumber).Times(1)
				allowlisttest.SetDefaultRoles(acp224feemanager.Module.Address)(t, state)
				require.NoError(t, acp224feemanager.StoreFeeConfig(state, testDefaultUpdatedFeeConfig, blockContext))
			},
			InputFn: func(t testing.TB) []byte {
				input, err := acp224feemanager.PackGetFeeConfigLastChangedAt()
				require.NoError(t, err)
				return input
			},
			ExpectedRes: func() []byte {
				packedOutput, err := acp224feemanager.PackGetFeeConfigLastChangedAtOutput(testBlockNumber)
				if err != nil {
					panic(err)
				}
				return packedOutput
			}(),
			AfterHook: func(t testing.TB, state *extstate.StateDB) {
				feeConfig := acp224feemanager.GetStoredFeeConfig(state)
				require.Equal(t, testDefaultUpdatedFeeConfig, feeConfig)
				lastChangedAt := acp224feemanager.GetFeeConfigLastChangedAt(state)
				require.Equal(t, testBlockNumber, lastChangedAt)
			},
			SuppliedGas: acp224feemanager.GetFeeConfigLastChangedAtGasCost,
			ReadOnly:    false,
			ExpectedErr: "",
		},
		"calling getFeeConfigLastChangedAt from Admin should succeed": {
			Caller: allowlisttest.TestAdminAddr,
			BeforeHook: func(t testing.TB, state *extstate.StateDB) {
				blockContext := contract.NewMockBlockContext(gomock.NewController(t))
				blockContext.EXPECT().Number().Return(testBlockNumber).Times(1)
				allowlisttest.SetDefaultRoles(acp224feemanager.Module.Address)(t, state)
				require.NoError(t, acp224feemanager.StoreFeeConfig(state, testDefaultUpdatedFeeConfig, blockContext))
			},
			InputFn: func(t testing.TB) []byte {
				input, err := acp224feemanager.PackGetFeeConfigLastChangedAt()
				require.NoError(t, err)
				return input
			},
			ExpectedRes: func() []byte {
				packedOutput, err := acp224feemanager.PackGetFeeConfigLastChangedAtOutput(testBlockNumber)
				if err != nil {
					panic(err)
				}
				return packedOutput
			}(),
			AfterHook: func(t testing.TB, state *extstate.StateDB) {
				feeConfig := acp224feemanager.GetStoredFeeConfig(state)
				require.Equal(t, testDefaultUpdatedFeeConfig, feeConfig)
				lastChangedAt := acp224feemanager.GetFeeConfigLastChangedAt(state)
				require.Equal(t, testBlockNumber, lastChangedAt)
			},
			SuppliedGas: acp224feemanager.GetFeeConfigLastChangedAtGasCost,
			ReadOnly:    false,
			ExpectedErr: "",
		},
		"insufficient gas for getFeeConfigLastChangedAt should fail": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				input, err := acp224feemanager.PackGetFeeConfigLastChangedAt()
				require.NoError(t, err)
				return input
			},
			SuppliedGas: acp224feemanager.GetFeeConfigLastChangedAtGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vm.ErrOutOfGas.Error(),
		},
		"calling setFeeConfig from NoRole should fail": {
			Caller:     allowlisttest.TestNoRoleAddr,
			BeforeHook: allowlisttest.SetDefaultRoles(acp224feemanager.Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := acp224feemanager.PackSetFeeConfig(testDefaultUpdatedFeeConfig)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: acp224feemanager.SetFeeConfigGasCost,
			ReadOnly:    false,
			ExpectedErr: acp224feemanager.ErrCannotSetFeeConfig.Error(),
		},
		"calling setFeeConfig from Enabled should succeed": {
			Caller:     allowlisttest.TestEnabledAddr,
			BeforeHook: allowlisttest.SetDefaultRoles(acp224feemanager.Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := acp224feemanager.PackSetFeeConfig(testDefaultUpdatedFeeConfig)
				require.NoError(t, err)
				return input
			},
			ExpectedRes: func() []byte {
				return []byte{}
			}(),
			SuppliedGas: acp224feemanager.SetFeeConfigGasCost,
			ReadOnly:    false,
			ExpectedErr: "",
		},
		"calling setFeeConfig from Manager should succeed": {
			Caller:     allowlisttest.TestManagerAddr,
			BeforeHook: allowlisttest.SetDefaultRoles(acp224feemanager.Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := acp224feemanager.PackSetFeeConfig(testDefaultUpdatedFeeConfig)
				require.NoError(t, err)
				return input
			},
			ExpectedRes: func() []byte {
				return []byte{}
			}(),
			SuppliedGas: acp224feemanager.SetFeeConfigGasCost,
			ReadOnly:    false,
			ExpectedErr: "",
		},
		"calling setFeeConfig from Admin should succeed": {
			Caller:     allowlisttest.TestAdminAddr,
			BeforeHook: allowlisttest.SetDefaultRoles(acp224feemanager.Module.Address),
			InputFn: func(t testing.TB) []byte {
				input, err := acp224feemanager.PackSetFeeConfig(testDefaultUpdatedFeeConfig)
				require.NoError(t, err)
				return input
			},
			ExpectedRes: func() []byte {
				return []byte{}
			}(),
			SuppliedGas: acp224feemanager.SetFeeConfigGasCost,
			ReadOnly:    false,
			ExpectedErr: "",
		},
		"readOnly setFeeConfig should fail": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				testInput := commontype.ACP224FeeConfig{
					TargetGas:          big.NewInt(10_000_000),
					MinGasPrice:        common.Big1,
					TimeToFillCapacity: big.NewInt(5),
					TimeToDouble:       big.NewInt(60),
				}
				input, err := acp224feemanager.PackSetFeeConfig(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: acp224feemanager.SetFeeConfigGasCost,
			ReadOnly:    true,
			ExpectedErr: vm.ErrWriteProtection.Error(),
		},
		"insufficient gas for setFeeConfig should fail": {
			Caller: common.Address{1},
			InputFn: func(t testing.TB) []byte {
				testInput := commontype.ACP224FeeConfig{
					TargetGas:          big.NewInt(10_000_000),
					MinGasPrice:        common.Big1,
					TimeToFillCapacity: big.NewInt(5),
					TimeToDouble:       big.NewInt(60),
				}
				input, err := acp224feemanager.PackSetFeeConfig(testInput)
				require.NoError(t, err)
				return input
			},
			SuppliedGas: acp224feemanager.SetFeeConfigGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vm.ErrOutOfGas.Error(),
		},
	}
)

// TestACP224FeeManagerRun tests the Run function of the precompile contract.
func TestACP224FeeManagerRun(t *testing.T) {
	allowlisttest.RunPrecompileWithAllowListTests(t, acp224feemanager.Module, tests)
}

// TestPackUnpackFeeConfigUpdatedEventData tests the Pack/UnpackFeeConfigUpdatedEventData.
func TestPackUnpackFeeConfigUpdatedEventData(t *testing.T) {
	var senderInput common.Address = acp224feemanager.ContractAddress
	oldFeeConfig := commontype.ValidTestACP224FeeConfig
	newFeeConfig := commontype.ACP224FeeConfig{
		TargetGas:          big.NewInt(42_000_000),
		MinGasPrice:        big.NewInt(42),
		TimeToFillCapacity: big.NewInt(42),
		TimeToDouble:       big.NewInt(42),
	}

	dataInput := acp224feemanager.FeeConfigUpdatedEventData{
		OldFeeConfig: oldFeeConfig,
		NewFeeConfig: newFeeConfig,
	}

	_, data, err := acp224feemanager.PackFeeConfigUpdatedEvent(
		senderInput,
		dataInput,
	)
	require.NoError(t, err)

	unpacked, err := acp224feemanager.UnpackFeeConfigUpdatedEventData(data)
	require.NoError(t, err)
	require.Equal(t, dataInput, unpacked)
}
