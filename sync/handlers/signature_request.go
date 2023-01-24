// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package handlers

import (
	"context"

	"github.com/ava-labs/avalanchego/codec"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/subnet-evm/plugin/evm/message"
	"github.com/ava-labs/subnet-evm/plugin/evm/warp"
)

// SignatureRequestHandler is a peer.RequestHandler for message.SignatureRequest
// serving requested BLS signature data
type SignatureRequestHandler struct {
	backend warp.WarpBackend
	codec   codec.Manager
}

func NewSignatureRequestHandler(backend warp.WarpBackend, codec codec.Manager) *SignatureRequestHandler {
	return &SignatureRequestHandler{
		backend: backend,
		codec:   codec,
	}
}

// OnSignatureRequest handles message.SignatureRequest, and retrieves a warp signature for the requested message ID.
func (s *SignatureRequestHandler) OnSignatureRequest(ctx context.Context, nodeID ids.NodeID, requestID uint32, signatureRequest message.SignatureRequest) ([]byte, error) {
	return s.backend.GetSignature(ctx, signatureRequest.MessageID)
}
