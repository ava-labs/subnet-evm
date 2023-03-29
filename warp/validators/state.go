// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package validators

import (
	"context"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/validators"
)

var _ validators.State = (*State)(nil)

type State struct {
	chainContext *snow.Context
	validators.State
}

// NewState returns a wrapper of [validators.State] which special cases the handling of the Primary Network.
//
// The wrapped state will return the chainContext's Subnet validator set instead of the Primary Network when
// the Primary Network SubnetID is passed in.
func NewState(chainContext *snow.Context) *State {
	return &State{
		chainContext: chainContext,
		State:        chainContext.ValidatorState,
	}
}

func (s *State) GetValidatorSet(
	ctx context.Context,
	height uint64,
	subnetID ids.ID,
) (map[ids.NodeID]*validators.GetValidatorOutput, error) {
	// If the subnetID is anything other than the Primary Network, this is a direct
	// passthrough
	if subnetID != ids.Empty {
		return s.State.GetValidatorSet(ctx, height, subnetID)
	}

	// If the requested subnet is the primary network, then we return the validator
	// set for the Subnet that is receiving the message instead.
	return s.State.GetValidatorSet(ctx, height, s.chainContext.SubnetID)
}
