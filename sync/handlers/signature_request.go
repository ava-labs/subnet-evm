// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package handlers

import (
	"context"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/subnet-evm/plugin/evm"
	"github.com/ava-labs/subnet-evm/plugin/evm/message"
)

type SignatureRequestHandler struct {
	backend evm.WarpBackend
}

func NewSignatureRequestHandler(backend evm.WarpBackend) *SignatureRequestHandler {
	return &SignatureRequestHandler{
		backend: backend,
	}
}

func (s *SignatureRequestHandler) OnSignatureRequest(ctx context.Context, nodeID ids.NodeID, requestID uint32, signatureRequest message.SignatureRequest) ([]byte, error) {
	return s.backend.GetSignature(ctx, signatureRequest.MessageID)
}
