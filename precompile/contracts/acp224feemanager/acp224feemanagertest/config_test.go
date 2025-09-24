// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package acp224feemanagertest

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ava-labs/avalanchego/vms/evm/upgrade/acp176"
	"github.com/ava-labs/libevm/common"
	"go.uber.org/mock/gomock"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/precompile/allowlist/allowlisttest"
	"github.com/ava-labs/subnet-evm/precompile/contracts/acp224feemanager"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/precompile/precompiletest"
	"github.com/ava-labs/subnet-evm/utils"
)

// TestVerify tests the verification of Config.
func TestVerify(t *testing.T) {
	admins := []common.Address{allowlisttest.TestAdminAddr}
	enableds := []common.Address{allowlisttest.TestEnabledAddr}
	managers := []common.Address{allowlisttest.TestManagerAddr}
	tests := map[string]precompiletest.ConfigVerifyTest{
		"valid config": {
			Config: acp224feemanager.NewConfig(utils.NewUint64(3), admins, enableds, managers, &commontype.ValidTestACP224FeeConfig),
			ChainConfig: func() precompileconfig.ChainConfig {
				config := precompileconfig.NewMockChainConfig(gomock.NewController(t))
				config.EXPECT().IsDurango(gomock.Any()).Return(true).AnyTimes()
				return config
			}(),
			ExpectedError: "",
		},
		"invalid - nil TargetGas": {
			Config: acp224feemanager.NewConfig(utils.NewUint64(3), admins, enableds, managers, &commontype.ACP224FeeConfig{
				TargetGas:          nil,
				MinGasPrice:        common.Big1,
				TimeToFillCapacity: big.NewInt(5),
				TimeToDouble:       big.NewInt(60),
			}),
			ExpectedError: "targetGas cannot be nil",
		},
		"invalid - nil MinGasPrice": {
			Config: acp224feemanager.NewConfig(utils.NewUint64(3), admins, enableds, managers, &commontype.ACP224FeeConfig{
				TargetGas:          big.NewInt(10_000_000),
				MinGasPrice:        nil,
				TimeToFillCapacity: big.NewInt(5),
				TimeToDouble:       big.NewInt(60),
			}),
			ExpectedError: "minGasPrice cannot be nil",
		},
		"invalid - nil TimeToFillCapacity": {
			Config: acp224feemanager.NewConfig(utils.NewUint64(3), admins, enableds, managers, &commontype.ACP224FeeConfig{
				TargetGas:          big.NewInt(10_000_000),
				MinGasPrice:        common.Big1,
				TimeToFillCapacity: nil,
				TimeToDouble:       big.NewInt(60),
			}),
			ExpectedError: "timeToFillCapacity cannot be nil",
		},
		"invalid - nil TimeToDouble": {
			Config: acp224feemanager.NewConfig(utils.NewUint64(3), admins, enableds, managers, &commontype.ACP224FeeConfig{
				TargetGas:          big.NewInt(10_000_000),
				MinGasPrice:        common.Big1,
				TimeToFillCapacity: big.NewInt(5),
				TimeToDouble:       nil,
			}),
			ExpectedError: "timeToDouble cannot be nil",
		},
		"invalid - TargetGas <= 0": {
			Config: acp224feemanager.NewConfig(utils.NewUint64(3), admins, enableds, managers, &commontype.ACP224FeeConfig{
				TargetGas:          big.NewInt(0),
				MinGasPrice:        common.Big1,
				TimeToFillCapacity: big.NewInt(5),
				TimeToDouble:       big.NewInt(60),
			}),
			ExpectedError: "targetGas = 0 cannot be less than or equal to 0",
		},
		"invalid - MinGasPrice <= 0": {
			Config: acp224feemanager.NewConfig(utils.NewUint64(3), admins, enableds, managers, &commontype.ACP224FeeConfig{
				TargetGas:          big.NewInt(10_000_000),
				MinGasPrice:        big.NewInt(0),
				TimeToFillCapacity: big.NewInt(5),
				TimeToDouble:       big.NewInt(60),
			}),
			ExpectedError: "minGasPrice = 0 cannot be less than or equal to 0",
		},
		"invalid - TimeToFillCapacity < 0": {
			Config: acp224feemanager.NewConfig(utils.NewUint64(3), admins, enableds, managers, &commontype.ACP224FeeConfig{
				TargetGas:          big.NewInt(10_000_000),
				MinGasPrice:        common.Big1,
				TimeToFillCapacity: big.NewInt(-1),
				TimeToDouble:       big.NewInt(60),
			}),
			ExpectedError: "timeToFillCapacity = -1 cannot be less than 0",
		},
		"invalid - TimeToDouble < 0": {
			Config: acp224feemanager.NewConfig(utils.NewUint64(3), admins, enableds, managers, &commontype.ACP224FeeConfig{
				TargetGas:          big.NewInt(10_000_000),
				MinGasPrice:        common.Big1,
				TimeToFillCapacity: big.NewInt(5),
				TimeToDouble:       big.NewInt(-1),
			}),
			ExpectedError: "timeToDouble = -1 cannot be less than 0",
		},
		"invalid - TimeToFillCapacity > MaxTimeToFillCapacity": {
			Config: acp224feemanager.NewConfig(utils.NewUint64(3), admins, enableds, managers, &commontype.ACP224FeeConfig{
				TargetGas:          big.NewInt(10_000_000),
				MinGasPrice:        common.Big1,
				TimeToFillCapacity: big.NewInt(acp176.MaxTimeToFillCapacity + 1),
				TimeToDouble:       big.NewInt(60),
			}),
			ExpectedError: fmt.Sprintf("timeToFillCapacity = %d cannot be greater than %d", big.NewInt(acp176.MaxTimeToFillCapacity+1), acp176.MaxTimeToFillCapacity),
		},
		"invalid - TimeToDouble > MaxTimeToDouble": {
			Config: acp224feemanager.NewConfig(utils.NewUint64(3), admins, enableds, managers, &commontype.ACP224FeeConfig{
				TargetGas:          big.NewInt(10_000_000),
				MinGasPrice:        common.Big1,
				TimeToFillCapacity: big.NewInt(5),
				TimeToDouble:       big.NewInt(acp176.MaxTimeToDouble + 1),
			}),
			ExpectedError: fmt.Sprintf("timeToDouble = %d cannot be greater than %d", big.NewInt(acp176.MaxTimeToDouble+1), acp176.MaxTimeToDouble),
		},
	}
	// Verify the precompile with the allowlist.
	// This adds allowlist verify tests to your custom tests
	// and runs them all together.
	// Even if you don't add any custom tests, keep this. This will still
	// run the default allowlist verify tests.
	allowlisttest.VerifyPrecompileWithAllowListTests(t, acp224feemanager.Module, tests)
}

// TestEqual tests the equality of Config with other precompile configs.
func TestEqual(t *testing.T) {
	admins := []common.Address{allowlisttest.TestAdminAddr}
	enableds := []common.Address{allowlisttest.TestEnabledAddr}
	managers := []common.Address{allowlisttest.TestManagerAddr}
	tests := map[string]precompiletest.ConfigEqualTest{
		"non-nil config and nil other": {
			Config:   acp224feemanager.NewConfig(utils.NewUint64(3), admins, enableds, managers, &commontype.ACP224FeeConfig{}),
			Other:    nil,
			Expected: false,
		},
		"different type": {
			Config:   acp224feemanager.NewConfig(utils.NewUint64(3), admins, enableds, managers, &commontype.ACP224FeeConfig{}),
			Other:    precompileconfig.NewMockConfig(gomock.NewController(t)),
			Expected: false,
		},
		"different timestamp": {
			Config:   acp224feemanager.NewConfig(utils.NewUint64(3), admins, enableds, managers, &commontype.ACP224FeeConfig{}),
			Other:    acp224feemanager.NewConfig(utils.NewUint64(4), admins, enableds, managers, &commontype.ACP224FeeConfig{}),
			Expected: false,
		},
		"same config": {
			Config:   acp224feemanager.NewConfig(utils.NewUint64(3), admins, enableds, managers, &commontype.ACP224FeeConfig{}),
			Other:    acp224feemanager.NewConfig(utils.NewUint64(3), admins, enableds, managers, &commontype.ACP224FeeConfig{}),
			Expected: true,
		},
		// CUSTOM CODE STARTS HERE
		// Add your own Equal tests here
	}
	// Run allow list equal tests.
	// This adds allowlist equal tests to your custom tests
	// and runs them all together.
	// Even if you don't add any custom tests, keep this. This will still
	// run the default allowlist equal tests.
	allowlisttest.EqualPrecompileWithAllowListTests(t, acp224feemanager.Module, tests)
}
