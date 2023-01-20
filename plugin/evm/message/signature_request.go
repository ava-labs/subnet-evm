// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package message

import (
	"context"
	"fmt"

	"github.com/ava-labs/avalanchego/utils/crypto/bls"

	"github.com/ava-labs/avalanchego/ids"
)

var _ Request = SignatureRequest{}

// SignatureRequest is used to request a warp message's signature.
type SignatureRequest struct {
	MessageID ids.ID `serialize:"true"`
}

func (s SignatureRequest) String() string {
	return fmt.Sprintf("SignatureRequest(MessageID=%s)", s.MessageID.String())
}

func (s SignatureRequest) Handle(ctx context.Context, nodeID ids.NodeID, requestID uint32, handler RequestHandler) ([]byte, error) {
	return handler.HandleSignatureRequest(ctx, nodeID, requestID, s)
}

type SignatureResponse struct {
	Signature [bls.SignatureLen]byte `serialize:"true"`
}
