// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package aggregator

import (
	"context"
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/plugin/evm/message"
)

const (
	initialRetryFetchSignatureDelay = 100 * time.Millisecond
	retryBackoffFactor              = 2
)

var _ SignatureGetter = (*NetworkSignatureGetter)(nil)

// SignatureGetter defines the minimum network interface to perform signature aggregation
type SignatureGetter interface {
	// GetSignature attempts to fetch a BLS Signature from [nodeID] for [unsignedWarpMessage]
	GetSignature(ctx context.Context, nodeID ids.NodeID, unsignedWarpMessage *avalancheWarp.UnsignedMessage) (*bls.Signature, error)
}

type NetworkClient interface {
	SendAppRequest(nodeID ids.NodeID, message []byte) ([]byte, error)
}

// NetworkSignatureGetter fetches warp signatures on behalf of the
// aggregator using VM App-Specific Messaging
type NetworkSignatureGetter struct {
	Client NetworkClient
}

// GetSignature attempts to fetch a BLS Signature of [unsignedWarpMessage] from [nodeID] until it succeeds or receives an invalid response
//
// Note: this function will continue attempting to fetch the signature from [nodeID] until it receives an invalid value or [ctx] is cancelled.
// The caller is responsible to cancel [ctx] if it no longer needs to fetch this signature.
func (s *NetworkSignatureGetter) GetSignature(ctx context.Context, nodeID ids.NodeID, unsignedWarpMessage *avalancheWarp.UnsignedMessage) (*bls.Signature, error) {
	signatureReq := message.SignatureRequest{
		MessageID: unsignedWarpMessage.ID(),
	}
	signatureReqBytes, err := message.RequestToBytes(message.Codec, signatureReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal signature request: %w", err)
	}

	delay := initialRetryFetchSignatureDelay
	timer := time.NewTimer(delay)
	defer timer.Stop()
	for {
		signatureRes, err := s.Client.SendAppRequest(nodeID, signatureReqBytes)
		if err != nil {
			// Wait until the retry delay has elapsed before retrying.
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(delay)

			select {
			case <-ctx.Done():
				return nil, err
			case <-timer.C:
			}

			// Exponential backoff.
			delay *= retryBackoffFactor
			continue
		}

		var response message.SignatureResponse
		if _, err := message.Codec.Unmarshal(signatureRes, &response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal signature res: %w", err)
		}

		blsSignature, err := bls.SignatureFromBytes(response.Signature[:])
		if err != nil {
			return nil, fmt.Errorf("failed to parse signature from res: %w", err)
		}
		return blsSignature, nil
	}
}
