// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package message

import (
	"context"

	"github.com/ava-labs/avalanchego/utils/crypto/bls"

	"github.com/ava-labs/avalanchego/ids"
)

var _ Request = SignatureRequest{}

// SignatureRequest is a request for the BLS signature for the Teleporter message identified by the MessageID/DestinationChainID pair
type SignatureRequest struct {
	UnsignedMessageID ids.ID `serialize:"true"`
}

func (s SignatureRequest) String() string {
	//TODO implement me
	panic("implement me")
}

func (s SignatureRequest) Handle(ctx context.Context, nodeID ids.NodeID, requestID uint32, handler RequestHandler) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

type SignatureResponse struct {
	Signature [bls.SignatureLen]byte `serialize:"true"`
}
