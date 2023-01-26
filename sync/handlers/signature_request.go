// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package handlers

import (
	"context"
	"time"

	"github.com/ava-labs/avalanchego/codec"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/subnet-evm/plugin/evm/message"
	"github.com/ava-labs/subnet-evm/plugin/evm/warp"
	"github.com/ava-labs/subnet-evm/sync/handlers/stats"
	"github.com/ethereum/go-ethereum/log"
)

// SignatureRequestHandler is a peer.RequestHandler for message.SignatureRequest
// serving requested BLS signature data
type SignatureRequestHandler struct {
	backend warp.WarpBackend
	codec   codec.Manager
	stats   stats.SignatureRequestHandlerStats
}

func NewSignatureRequestHandler(backend warp.WarpBackend, codec codec.Manager, stats stats.SignatureRequestHandlerStats) *SignatureRequestHandler {
	return &SignatureRequestHandler{
		backend: backend,
		codec:   codec,
		stats:   stats,
	}
}

// OnSignatureRequest handles message.SignatureRequest, and retrieves a warp signature for the requested message ID.
// Never returns an error
// Expects returned errors to be treated as FATAL
// Returns empty response if signature is not found
// Assumes ctx is active
func (s *SignatureRequestHandler) OnSignatureRequest(ctx context.Context, nodeID ids.NodeID, requestID uint32, signatureRequest message.SignatureRequest) ([]byte, error) {
	startTime := time.Now()
	s.stats.IncSignatureRequest()

	// always report signature request time
	defer func() {
		s.stats.UpdateSignatureRequestTime(time.Since(startTime))
	}()

	var signature [bls.SignatureLen]byte
	sig, err := s.backend.GetSignature(ctx, signatureRequest.MessageID)
	if err != nil {
		log.Debug("Unknown warp signature requested", "messageID", signatureRequest.MessageID)
		s.stats.IncSignatureMiss()
		return nil, nil
	}

	s.stats.IncSignatureHit()
	copy(signature[:], sig)
	response := message.SignatureResponse{Signature: signature}
	responseBytes, err := s.codec.Marshal(message.Version, response)
	if err != nil {
		log.Warn("could not marshal SignatureResponse, dropping request", "nodeID", nodeID, "requestID", requestID, "err", err)
		return nil, nil
	}

	return responseBytes, nil
}
