// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package aggregator

import (
	"context"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/validators"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	warpPrecompile "github.com/ava-labs/subnet-evm/precompile/contracts/warp"
	wrappedValidators "github.com/ava-labs/subnet-evm/warp/validators"
)

type Aggregator struct {
	chainContext *snow.Context
	client       ClientBackend
	state        validators.State
}

func NewAggregator(chainContext *snow.Context, client ClientBackend) *Aggregator {
	return &Aggregator{
		chainContext: chainContext,
		client:       client,
		state:        wrappedValidators.NewState(chainContext),
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
		a.chainContext.SubnetID,
		quorumNum,
		quorumNum,
		warpPrecompile.DefaultQuorumNumerator,
		a.state,
		unsignedMessage,
	)

	return job.Execute(ctx)
}
