// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"encoding/hex"

	"github.com/ava-labs/subnet-evm/plugin/evm/message"
	"github.com/ethereum/go-ethereum/log"
)

const (
	MessageIDLength = 32
)

// WarpAPI introduces snowman specific functionality to the evm
type WarpAPI struct{ vm *VM }

type SignatureRequest struct {
	MessageID string `json:"messageID"`
}

func (api *WarpAPI) GetSignature(ctx context.Context, signatureRequest *SignatureRequest) (*message.SignatureResponse, error) {
	sigReqBytes, err := hex.DecodeString(signatureRequest.MessageID)
	if err != nil || len(sigReqBytes) != MessageIDLength {
		log.Info("Invalid messageID hex in signature request")
		return nil, err
	}

	signature, err := api.vm.backend.GetSignature(ctx, *(*[32]byte)(sigReqBytes))
	if err != nil {
		log.Debug("Unknown warp signature requested", "messageID", signatureRequest.MessageID)
		return nil, nil
	}

	response := message.SignatureResponse{
		Signature: signature,
	}
	return &response, nil
}
