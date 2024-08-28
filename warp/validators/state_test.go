// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package validators

import (
	"context"
	"testing"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/validators"
	"github.com/ava-labs/avalanchego/snow/validators/validatorstest"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/stretchr/testify/require"
)

func TestGetValidatorSetPrimaryNetwork(t *testing.T) {
	require := require.New(t)

	mySubnetID := ids.GenerateTestID()
	otherSubnetID := ids.GenerateTestID()

	snowCtx := utils.TestSnowContext()
	snowCtx.SubnetID = mySubnetID
	snowCtx.ValidatorState = &validatorstest.State{
		GetValidatorSetF: func(context.Context, uint64, ids.ID) (map[ids.NodeID]*validators.GetValidatorOutput, error) {
			return nil, nil
		},
	}
	state := NewState(snowCtx.ValidatorState, snowCtx.SubnetID, snowCtx.ChainID, false)
	// Expect that requesting my validator set returns my validator set
	output, err := state.GetValidatorSet(context.Background(), 10, mySubnetID)
	require.NoError(err)
	require.Len(output, 0)

	// Expect that requesting the Primary Network validator set overrides and returns my validator set
	output, err = state.GetValidatorSet(context.Background(), 10, constants.PrimaryNetworkID)
	require.NoError(err)
	require.Len(output, 0)

	// Expect that requesting other validator set returns that validator set
	output, err = state.GetValidatorSet(context.Background(), 10, otherSubnetID)
	require.NoError(err)
	require.Len(output, 0)
}
