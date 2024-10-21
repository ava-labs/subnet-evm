// (c) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"context"
	"fmt"

	"github.com/ava-labs/subnet-evm/warp/messages"

	"github.com/ava-labs/avalanchego/snow/engine/common"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/payload"
)

const (
	ParseErrCode = iota + 1
	VerifyErrCode
)

// Verify implements the acp118.Verifier interface
func (b *backend) Verify(_ context.Context, unsignedMessage *avalancheWarp.UnsignedMessage, _ []byte) *common.AppError {
	return b.verifyMessage(unsignedMessage)
}

// verifyMessage verifies the signature of the message
// This is moved to a separate function to avoid having to use a context.Context
func (b *backend) verifyMessage(unsignedMessage *avalancheWarp.UnsignedMessage) *common.AppError {
	messageID := unsignedMessage.ID()
	// Known on-chain messages should be signed
	if _, err := b.GetMessage(messageID); err == nil {
		return nil
	}

	parsed, err := payload.Parse(unsignedMessage.Payload)
	if err != nil {
		b.stats.IncMessageParseFail()
		return &common.AppError{
			Code:    ParseErrCode,
			Message: "failed to parse payload: " + err.Error(),
		}
	}

	switch p := parsed.(type) {
	case *payload.AddressedCall:
		apperr := b.verifyAddressedCall(p)
		if apperr != nil {
			return apperr
		}
	case *payload.Hash:
		apperr := b.verifyBlockMessage(p)
		if apperr != nil {
			return apperr
		}
	default:
		b.stats.IncMessageParseFail()
		return &common.AppError{
			Code:    ParseErrCode,
			Message: fmt.Sprintf("unknown payload type: %T", p),
		}
	}
	return nil
}

// verifyBlockMessage returns nil if blockHashPayload contains the ID
// of an accepted block indicating it should be signed by the VM.
func (b *backend) verifyBlockMessage(blockHashPayload *payload.Hash) *common.AppError {
	blockID := blockHashPayload.Hash
	_, err := b.blockClient.GetAcceptedBlock(context.TODO(), blockID)
	if err != nil {
		b.stats.IncBlockSignatureValidationFail()
		return &common.AppError{
			Code:    VerifyErrCode,
			Message: fmt.Sprintf("failed to get block %s: %s", blockID, err.Error()),
		}
	}

	return nil
}

// verifyAddressedCall returns nil if addressedCall is parseable to a known payload type and
// passes type specific validation, indicating it should be signed by the VM.
// Note currently there are no valid payload types so this call always returns common.AppError
// with ParseErrCode.
func (b *backend) verifyAddressedCall(addressedCall *payload.AddressedCall) *common.AppError {
	// Parse the payload to see if it is a known type.
	parsed, err := messages.Parse(addressedCall.Payload)
	if err != nil {
		b.stats.IncMessageParseFail()
		return &common.AppError{
			Code:    ParseErrCode,
			Message: "failed to parse addressed call message: " + err.Error(),
		}
	}

	switch p := parsed.(type) {
	default:
		b.stats.IncMessageParseFail()
		return &common.AppError{
			Code:    ParseErrCode,
			Message: fmt.Sprintf("unknown message type: %T", p),
		}
	}
}
