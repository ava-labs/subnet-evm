// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package aggregator

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ava-labs/subnet-evm/params"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/validators"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/set"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
)

var (
	errNoValidators       = errors.New("cannot aggregate signatures from subnet with no validators")
	errInsufficientWeight = errors.New("verification failed with insufficient weight")
)

// SignatureGetter defines the minimum network interface to perform signature aggregation
type SignatureGetter interface {
	// GetSignature attempts to fetch a BLS Signature from [nodeID] for [unsignedWarpMessage]
	GetSignature(ctx context.Context, nodeID ids.NodeID, unsignedWarpMessage *avalancheWarp.UnsignedMessage) (*bls.Signature, error)
}

type AggregateSignatureResult struct {
	// Weight of validators included in the aggregate signature.
	SignatureWeight uint64
	// Total weight of all validators in the subnet.
	TotalWeight uint64
	// The message with the aggregate signature.
	Message *avalancheWarp.Message
}

// Aggregator requests signatures from validators and
// aggregates them into a single signature.
type Aggregator struct {
	// Aggregating signatures for a chain validated by this subnet.
	subnetID ids.ID
	// Fetches signatures from validators.
	client SignatureGetter
	// Validator state for this chain.
	state validators.State
}

// New returns a signature aggregator for the chain with the given [state] on the
// given [subnet], and where [client] can be used to fetch signatures from validators.
func New(subnetID ids.ID, state validators.State, client SignatureGetter) *Aggregator {
	return &Aggregator{
		subnetID: subnetID,
		client:   client,
		state:    state,
	}
}

func (a *Aggregator) AggregateSignatures(ctx context.Context, unsignedMessage *avalancheWarp.UnsignedMessage, quorumNum uint64) (*AggregateSignatureResult, error) {
	// Note: we use the current height as a best guess of the canonical validator set when the aggregated signature will be verified
	// by the recipient chain. If the validator set changes from [pChainHeight] to the P-Chain height that is actually specified by the
	// ProposerVM header when this message is verified, then the aggregate signature could become outdated and require re-aggregation.
	pChainHeight, err := a.state.GetCurrentHeight(ctx)
	if err != nil {
		return nil, err
	}

	log.Debug("Fetching signature",
		"a.subnetID", a.subnetID,
		"height", pChainHeight,
	)
	validators, totalWeight, err := avalancheWarp.GetCanonicalValidatorSet(ctx, a.state, pChainHeight, a.subnetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get validator set: %w", err)
	}
	if len(validators) == 0 {
		return nil, fmt.Errorf("%w (SubnetID: %s, Height: %d)", errNoValidators, a.subnetID, pChainHeight)
	}

	return aggregateSignatures(ctx, a.client, unsignedMessage, validators, quorumNum, totalWeight)
}

func aggregateSignatures(
	ctx context.Context,
	client SignatureGetter,
	unsignedMessage *avalancheWarp.UnsignedMessage,
	validators []*avalancheWarp.Validator,
	quorumNum uint64,
	totalWeight uint64,
) (*AggregateSignatureResult, error) {
	var (
		// [signatureLock] must be held when accessing [blsSignatures],
		// [signersBitset], or [signatureWeight] in the goroutine below.
		signatureLock   = sync.Mutex{}
		signatures      = make([]*bls.Signature, 0, len(validators))
		signersBitset   = set.NewBits()
		signatureWeight = uint64(0)
	)

	// Create a child context to cancel signature fetching if we reach [maxNeededQuorumNum] threshold
	signatureFetchCtx, signatureFetchCancel := context.WithCancel(ctx)
	defer signatureFetchCancel()

	wg := sync.WaitGroup{}
	wg.Add(len(validators))
	for i, validator := range validators {
		var (
			i         = i
			validator = validator
			// TODO: update from a single nodeID to the original slice and use extra nodeIDs as backup.
			nodeID = validator.NodeIDs[0]
		)
		go func() {
			defer wg.Done()

			log.Debug("Fetching warp signature",
				"nodeID", nodeID,
				"index", i,
			)

			signature, err := client.GetSignature(signatureFetchCtx, nodeID, unsignedMessage)
			if err != nil {
				log.Debug("Failed to fetch warp signature",
					"nodeID", nodeID,
					"index", i,
					"err", err,
				)
				return
			}

			sigHex := hexutil.Bytes(bls.SignatureToBytes(signature))
			log.Debug("Retrieved warp signature",
				"nodeID", nodeID,
				"index", i,
				"signature", sigHex,
			)

			if !bls.Verify(validator.PublicKey, signature, unsignedMessage.Bytes()) {
				log.Debug("Failed to verify warp signature",
					"nodeID", nodeID,
					"index", i,
					"signature", sigHex,
					"msgID", unsignedMessage.ID(),
				)
				return
			}

			// Add the signature and check if we've reached the requested threshold
			signatureLock.Lock()
			defer signatureLock.Unlock()

			signatures = append(signatures, signature)
			signersBitset.Add(i)
			signatureWeight += validator.Weight
			log.Debug("Updated weight",
				"totalWeight", signatureWeight,
				"addedWeight", validator.Weight,
			)

			// If the signature weight meets the requested threshold, cancel signature fetching
			if err := avalancheWarp.VerifyWeight(signatureWeight, totalWeight, quorumNum, params.WarpQuorumDenominator); err == nil {
				log.Debug("Verify weight passed, exiting aggregation early",
					"maxNeededQuorumNum", quorumNum,
					"totalWeight", totalWeight,
					"signatureWeight", signatureWeight,
				)
				signatureFetchCancel()
			}
		}()
	}
	wg.Wait()

	// If I failed to fetch sufficient signature stake, return an error
	if err := avalancheWarp.VerifyWeight(signatureWeight, totalWeight, quorumNum, params.WarpQuorumDenominator); err != nil {
		return nil, fmt.Errorf("%w: %w", errInsufficientWeight, err)
	}

	// Otherwise, return the aggregate signature
	aggregateSignature, err := bls.AggregateSignatures(signatures)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate BLS signatures: %w", err)
	}

	warpSignature := &avalancheWarp.BitSetSignature{
		Signers: signersBitset.Bytes(),
	}
	copy(warpSignature.Signature[:], bls.SignatureToBytes(aggregateSignature))

	msg, err := avalancheWarp.NewMessage(unsignedMessage, warpSignature)
	if err != nil {
		return nil, fmt.Errorf("failed to construct warp message: %w", err)
	}

	return &AggregateSignatureResult{
		Message:         msg,
		SignatureWeight: signatureWeight,
		TotalWeight:     totalWeight,
	}, nil
}
