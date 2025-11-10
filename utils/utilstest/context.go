// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utilstest

import (
	"testing"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/snowtest"
)

// NewTestSnowContext returns a snow.Context with validator state properly configured for testing.
// This wraps snowtest.Context and sets the validator state to avoid the missing GetValidatorSetF issue.
//
// Usage example:
//
//	// Instead of:
//	// snowCtx := utilstest.NewTestSnowContext(t, snowtest.CChainID)
//	// validatorState := utils.NewTestValidatorState()
//	// snowCtx.ValidatorState = validatorState
//
//	// Use:
//	snowCtx := utils.NewTestSnowContext(t)
//
// This function ensures that the snow context has a properly configured validator state
// that includes the GetValidatorSetF function, which is required by many tests.
func NewTestSnowContext(t testing.TB) *snow.Context {
	snowCtx := snowtest.Context(t, SubnetEVMTestChainID)
	snowCtx.ValidatorState = NewTestValidatorState()
	return snowCtx
}

// NewTestSnowContextWithChainID returns a snow.Context with validator state properly configured for testing
// with a specific chain ID. This is provided for backward compatibility when a specific chain ID is needed.
func NewTestSnowContextWithChainID(t testing.TB, chainID ids.ID) *snow.Context {
	snowCtx := snowtest.Context(t, chainID)
	snowCtx.ValidatorState = NewTestValidatorState()
	return snowCtx
}
