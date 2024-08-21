// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package validators

import (
	"context"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/validators"
	"github.com/ava-labs/avalanchego/utils/constants"
)

var _ validators.State = (*State)(nil)

// State provides a special case used to handle Avalanche Warp Message verification for messages sent
// from the Primary Network. Subnets have strictly fewer validators than the Primary Network, so we require
// signatures from a threshold of the RECEIVING subnet validator set rather than the full Primary Network
// since the receiving subnet already relies on a majority of its validators being correct.
type State struct {
	validators.State
	chainContext                 *snow.Context
	requirePrimaryNetworkSigners func() bool
}

// NewState returns a wrapper of [validators.State] which special cases the handling of the Primary Network.
//
// The wrapped state will return the chainContext's Subnet validator set instead of the Primary Network when
// the Primary Network SubnetID is passed in.
// Additionally, it will return the chainContext's Subnet instead of the P-Chain, so that messages from the
// P-Chains are verified against the Subnet's validator set.
func NewState(chainContext *snow.Context, requirePrimaryNetworkSigners func() bool) *State {
	return &State{
		State:                        chainContext.ValidatorState,
		chainContext:                 chainContext,
		requirePrimaryNetworkSigners: requirePrimaryNetworkSigners,
	}
}

func (s *State) GetSubnetID(ctx context.Context, chainID ids.ID) (ids.ID, error) {
	// Messages from the P-Chain should be verified against the Subnet's validator set
	if chainID == constants.PlatformChainID {
		return s.chainContext.SubnetID, nil
	}

	return s.State.GetSubnetID(ctx, chainID)
}

func (s *State) GetValidatorSet(
	ctx context.Context,
	height uint64,
	subnetID ids.ID,
) (map[ids.NodeID]*validators.GetValidatorOutput, error) {
	// If the subnetID is anything other than the Primary Network, or Primary
	// Network signers are required, this is a direct passthrough.
	if s.requirePrimaryNetworkSigners() || subnetID != constants.PrimaryNetworkID {
		return s.State.GetValidatorSet(ctx, height, subnetID)
	}

	// If the requested subnet is the primary network, then we return the validator
	// set for the Subnet that is receiving the message instead.
	return s.State.GetValidatorSet(ctx, height, s.chainContext.SubnetID)
}
