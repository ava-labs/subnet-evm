// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"

	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/subnet-evm/plugin/evm/message"
	"github.com/ethereum/go-ethereum/log"
)

// SnowmanAPI introduces snowman specific functionality to the evm
type WarpAPI struct{ vm *VM }

func (api *WarpAPI) GetSignature(ctx context.Context, signatureRequest *message.SignatureRequest) (*message.SignatureResponse, error) {
	sig, err := api.vm.backend.GetSignature(ctx, signatureRequest.MessageID)
	if err != nil {
		log.Debug("Unknown warp signature requested", "messageID", signatureRequest.MessageID)
		return nil, nil
	}
	var signature [bls.SignatureLen]byte
	copy(signature[:], sig)
	response := message.SignatureResponse{
		Signature: signature,
	}
	return &response, nil
}
