// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/rpc"
	"github.com/ava-labs/subnet-evm/plugin/evm/message"
)

// Interface compliance
var _ WarpClient = (*warpClient)(nil)

type WarpClient interface {
	GetSignature(ctx context.Context, signatureRequest message.SignatureRequest) (*[bls.SignatureLen]byte, error)
}

// Client implementation for interacting with EVM [chain]
type warpClient struct {
	requester rpc.EndpointRequester
}

// NewClient returns a Client for interacting with EVM [chain]
func NewWarpClient(uri, chain string) WarpClient {
	return &warpClient{
		requester: rpc.NewEndpointRequester(fmt.Sprintf("%s/ext/bc/%s/rpc", uri, chain)),
	}
}

func (c *warpClient) GetSignature(ctx context.Context, signatureRequest message.SignatureRequest) (*[bls.SignatureLen]byte, error) {
	sigReqJson := SignatureRequest{
		MessageID: hex.EncodeToString(signatureRequest.MessageID[:]),
	}
	res := &message.SignatureResponse{}
	err := c.requester.SendRequest(ctx, "warp_getSignature", &sigReqJson, res)
	return &res.Signature, err
}
