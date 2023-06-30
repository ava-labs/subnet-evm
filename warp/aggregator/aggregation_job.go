// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package aggregator

import (
	"context"
	"fmt"
	"sync"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/validators"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/set"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
)

type signatureAggregationJob struct {
	// SignatureBackend is assumed to be thread-safe and may be used by multiple signature aggregation jobs concurrently
	client   SignatureBackend
	height   uint64
	subnetID ids.ID

	quorumNum       uint64 // Minimum threshold at which to bother returning the resulting signature
	cancelQuorumNum uint64 // Threshold at which to cancel further signature fetching
	quorumDen       uint64 // Denominator to use when checking if we've reached the threshold
	state           validators.State
	msg             *avalancheWarp.UnsignedMessage
}

type AggregateSignatureResult struct {
	SignatureWeight uint64
	TotalWeight     uint64
	Message         *avalancheWarp.Message
}

func NewSignatureAggregationJob(
	client SignatureBackend,
	height uint64,
	subnetID ids.ID,
	quorumNum uint64,
	cancelQuorumNum uint64,
	quorumDen uint64,
	state validators.State,
	msg *avalancheWarp.UnsignedMessage,
) *signatureAggregationJob {
	return &signatureAggregationJob{
		client:          client,
		height:          height,
		subnetID:        subnetID,
		quorumNum:       quorumNum,
		cancelQuorumNum: cancelQuorumNum,
		quorumDen:       quorumDen,
		state:           state,
		msg:             msg,
	}
}

// Execute aggregates signatures for the requested message
func (a *signatureAggregationJob) Execute(ctx context.Context) (*AggregateSignatureResult, error) {
	validators, totalWeight, err := avalancheWarp.GetCanonicalValidatorSet(ctx, a.state, a.height, a.subnetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get validator set: %w", err)
	}
	signatureJobs := make([]*signatureJob, 0, len(validators))
	for _, validator := range validators {
		signatureJobs = append(signatureJobs, newSignatureJob(a.client, validator, a.msg))
	}

	// signatureLock is used to access any of the signature attributes in the goroutines created below
	signatureLock := sync.Mutex{}
	blsSignatures := make([]*bls.Signature, 0, len(signatureJobs))
	bitSet := set.NewBits()
	signatureWeight := uint64(0)

	// Create a child context to cancel signature fetching if we reach [cancelQuorumNum] threshold
	signatureFetchCtx, signatureFetchCancel := context.WithCancel(ctx)
	defer signatureFetchCancel()

	wg := sync.WaitGroup{}
	for i, signatureJob := range signatureJobs {
		i := i
		signatureJob := signatureJob
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Info("Fetching warp signature", "nodeID", signatureJob.nodeID, "index", i)
			blsSignature, err := signatureJob.Execute(signatureFetchCtx)
			if err != nil {
				log.Info("Failed to fetch signature at index %d: %s", i, signatureJob)
				return
			}
			log.Info("Retrieved warp signature", "nodeID", signatureJob.nodeID, "index", i, "signature", hexutil.Bytes(bls.SignatureToBytes(blsSignature)))
			// Add the signature and check if we've reached the requested threshold
			signatureLock.Lock()
			defer signatureLock.Unlock()

			blsSignatures = append(blsSignatures, blsSignature)
			bitSet.Add(i)
			log.Info("Updated weight", "totalWeight", signatureWeight+signatureJob.weight, "addedWeight", signatureJob.weight)
			signatureWeight += signatureJob.weight
			// If the signature weight meets the requested threshold, cancel signature fetching
			if err := avalancheWarp.VerifyWeight(signatureWeight, totalWeight, a.cancelQuorumNum, a.quorumDen); err == nil {
				log.Info("Verify weight passed, exiting aggregation early", "cancelQuorumNum", a.cancelQuorumNum, "totalWeight", totalWeight, "signatureWeight", signatureWeight)
				signatureFetchCancel()
			}
		}()
	}
	wg.Wait()

	// If I failed to fetch sufficient signature stake, return an error
	if err := avalancheWarp.VerifyWeight(signatureWeight, totalWeight, a.quorumNum, a.quorumDen); err != nil {
		return nil, fmt.Errorf("failed to aggregate signature: %w", err)
	}
	// Otherwise, return the aggregate signature
	aggregateSignature, err := bls.AggregateSignatures(blsSignatures)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate BLS signatures: %w", err)
	}
	warpSignature := &avalancheWarp.BitSetSignature{
		Signers: bitSet.Bytes(),
	}
	copy(warpSignature.Signature[:], bls.SignatureToBytes(aggregateSignature))
	msg, err := avalancheWarp.NewMessage(a.msg, warpSignature)
	if err != nil {
		return nil, fmt.Errorf("failed to construct warp message: %w", err)
	}
	return &AggregateSignatureResult{
		Message:         msg,
		SignatureWeight: signatureWeight,
		TotalWeight:     totalWeight,
	}, nil
}
