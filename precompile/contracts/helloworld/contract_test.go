// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package helloworld

import (
	"testing"

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestHelloWorld(t *testing.T) {
	testGreeting := "test"
	tests := map[string]testutils.PrecompileTest{
		"set greeting from no role fails": {
			Caller:     allowlist.TestNoRoleAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t *testing.T) []byte {
				input, err := PackSetGreeting("test")
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetGreetingGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrCannotSetGreeting.Error(),
		},
		"set greeting from enabled address": {
			Caller:     allowlist.TestEnabledAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t *testing.T) []byte {
				input, err := PackSetGreeting(testGreeting)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetGreetingGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				greeting := GetGreeting(state)
				require.Equal(t, greeting, testGreeting)
			},
		},
		"set greeting from admin address": {
			Caller:     allowlist.TestAdminAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t *testing.T) []byte {
				input, err := PackSetGreeting(testGreeting)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetGreetingGasCost,
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t *testing.T, state contract.StateDB) {
				greeting := GetGreeting(state)
				require.Equal(t, greeting, testGreeting)
			},
		},
		"get default hello from non-enabled address": {
			Caller:     allowlist.TestNoRoleAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t *testing.T) []byte {
				input, err := PackSayHello()
				require.NoError(t, err)

				return input
			},
			Config:      NewConfig(common.Big0, nil, nil), // give a zero config for immediate activation
			SuppliedGas: SayHelloGasCost,
			ReadOnly:    true,
			ExpectedRes: func() []byte {
				res, err := PackSayHelloOutput(defaultGreeting)
				require.NoError(t, err)
				return res
			}(),
		},
		"store greeting then say hello from non-enabled address": {
			Caller: allowlist.TestNoRoleAddr,
			BeforeHook: func(t *testing.T, state contract.StateDB) {
				allowlist.SetDefaultRoles(Module.Address)(t, state)
				StoreGreeting(state, testGreeting)
			},
			InputFn: func(t *testing.T) []byte {
				input, err := PackSayHello()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SayHelloGasCost,
			ReadOnly:    true,
			ExpectedRes: func() []byte {
				res, err := PackSayHelloOutput(testGreeting)
				require.NoError(t, err)
				return res
			}(),
		},
		"set a very long greeting from enabled address": {
			Caller:     allowlist.TestEnabledAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t *testing.T) []byte {
				longString := "a very long string that is longer than 32 bytes and will cause an error"
				input, err := PackSetGreeting(longString)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetGreetingGasCost,
			ReadOnly:    false,
			ExpectedErr: ErrInputExceedsLimit.Error(),
		},
		"readOnly setFeeConfig with noRole fails": {
			Caller:     allowlist.TestNoRoleAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t *testing.T) []byte {
				input, err := PackSetGreeting(testGreeting)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetGreetingGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly setFeeConfig with enabled role fails": {
			Caller:     allowlist.TestEnabledAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t *testing.T) []byte {
				input, err := PackSetGreeting(testGreeting)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetGreetingGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"readOnly setFeeConfig with admin role fails": {
			Caller:     allowlist.TestAdminAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t *testing.T) []byte {
				input, err := PackSetGreeting(testGreeting)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetGreetingGasCost,
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"insufficient gas setFeeConfig from admin": {
			Caller:     allowlist.TestAdminAddr,
			BeforeHook: allowlist.SetDefaultRoles(Module.Address),
			InputFn: func(t *testing.T) []byte {
				input, err := PackSetGreeting(testGreeting)
				require.NoError(t, err)

				return input
			},
			SuppliedGas: SetGreetingGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
	}

	allowlist.RunPrecompileWithAllowListTests(t, Module, state.NewTestStateDB, tests)
}
