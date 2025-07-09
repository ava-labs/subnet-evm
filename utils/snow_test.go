// (c) 2025 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"testing"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/snowtest"
	"github.com/ava-labs/subnet-evm/utils/utilstest"
	"github.com/stretchr/testify/require"
)

func TestNewTestSnowContext(t *testing.T) {
	// Test that NewTestSnowContext creates a context with validator state
	snowCtx := utilstest.NewTestSnowContext(t, snowtest.CChainID)
	require.NotNil(t, snowCtx.ValidatorState)

	// Test that the validator state has the required functions
	validatorState := snowCtx.ValidatorState
	require.NotNil(t, validatorState)

	// Test that we can call GetValidatorSetF without panicking
	validators, err := validatorState.GetValidatorSet(nil, 0, ids.Empty)
	require.NoError(t, err)
	require.NotNil(t, validators)

	// Test that we can call GetCurrentValidatorSetF without panicking
	currentValidators, height, err := validatorState.GetCurrentValidatorSet(nil, ids.Empty)
	require.NoError(t, err)
	require.NotNil(t, currentValidators)
	require.Equal(t, uint64(0), height)
}
