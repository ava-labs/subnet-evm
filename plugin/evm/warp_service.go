// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/subnet-evm/warp"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type SignatureResponse struct {
	Signature hexutil.Bytes
}

// WarpAPI introduces snowman specific functionality to the evm
type WarpAPI struct {
	backend warp.WarpBackend
}

// GetSignature returns the BLS signature associated with a messageID. In the raw request, [messageID] should be cb58 encoded
func (api *WarpAPI) GetSignature(ctx context.Context, messageID ids.ID) (*SignatureResponse, error) {
	signature, err := api.backend.GetSignature(ctx, messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get signature for with error %w", err)
	}
	sigBytes := (hexutil.Bytes)(signature[:])

	response := SignatureResponse{
		Signature: sigBytes,
	}
	return &response, nil
}
