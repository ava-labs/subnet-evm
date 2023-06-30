// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package aggregator

import (
	"context"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/validators"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/params"
)

type Aggregator struct {
	subnetID ids.ID
	client   SignatureBackend
	state    validators.State
}

func NewAggregator(subnetID ids.ID, state validators.State, client SignatureBackend) *Aggregator {
	return &Aggregator{
		subnetID: subnetID,
		client:   client,
		state:    state,
	}
}

func (a *Aggregator) AggregateSignatures(ctx context.Context, unsignedMessage *avalancheWarp.UnsignedMessage, quorumNum uint64) (*AggregateSignatureResult, error) {
	pChainHeight, err := a.state.GetCurrentHeight(ctx)
	if err != nil {
		return nil, err
	}
	job := NewSignatureAggregationJob(
		a.client,
		pChainHeight,
		a.subnetID,
		quorumNum,
		quorumNum,
		params.WarpQuorumDenominator,
		a.state,
		unsignedMessage,
	)

	return job.Execute(ctx)
}
