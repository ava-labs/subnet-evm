// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/subnet-evm/plugin/evm/message"
	"github.com/ethereum/go-ethereum/log"
)

// WarpAPI introduces snowman specific functionality to the evm
type WarpAPI struct{ vm *VM }

// GetSignature returns the BLS signature associated with a messageID. In the raw request, [messageID] should be cb58 encoded
func (api *WarpAPI) GetSignature(ctx context.Context, messageID ids.ID) (*message.SignatureResponse, error) {
	signature, err := api.vm.backend.GetSignature(ctx, messageID)
	if err != nil {
		log.Debug("Unknown warp signature requested", "messageID", messageID)
		return nil, nil
	}

	response := message.SignatureResponse{
		Signature: signature,
	}
	return &response, nil
}
