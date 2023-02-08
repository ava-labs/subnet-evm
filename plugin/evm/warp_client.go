// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"fmt"

	"github.com/ava-labs/subnet-evm/rpc"

	"github.com/ava-labs/avalanchego/utils/cb58"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/subnet-evm/plugin/evm/message"
)

// Interface compliance
var _ WarpClient = (*warpClient)(nil)

type WarpClient interface {
	GetSignature(ctx context.Context, signatureRequest message.SignatureRequest) (*[bls.SignatureLen]byte, error)
}

// warpClient implementation for interacting with EVM [chain]
type warpClient struct {
	client *rpc.Client
}

// NewWarpClient returns a WarpClient for interacting with EVM [chain]
func NewWarpClient(uri, chain string) (WarpClient, error) {
	client, err := rpc.Dial(fmt.Sprintf("%s/ext/bc/%s/rpc", uri, chain))
	if err != nil {
		return nil, fmt.Errorf("failed to dial client. err: %w", err)
	}
	return &warpClient{
		client: client,
	}, nil
}

// GetSignature requests the BLS signature associated with a messageID
func (c *warpClient) GetSignature(ctx context.Context, signatureRequest message.SignatureRequest) (*[bls.SignatureLen]byte, error) {
	req, err := cb58.Encode(signatureRequest.MessageID[:])
	if err != nil {
		return nil, fmt.Errorf("failed to base58 encode the request. messageID: %s, error: %w", signatureRequest.MessageID, err)
	}

	var res message.SignatureResponse
	err = c.client.CallContext(ctx, &res, "warp_getSignature", req)
	if err != nil {
		return nil, fmt.Errorf("call to warp_getSignature failed. err: %w", err)
	}
	return &res.Signature, err
}
