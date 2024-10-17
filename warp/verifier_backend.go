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
	if err := b.verifyMessage(unsignedMessage); err != nil {
		return err
	}
	return nil
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
			b.stats.IncAddressedCallSignatureValidationFail()
			return apperr
		}
	case *payload.Hash:
		apperr := b.verifyBlockMessage(p)
		if apperr != nil {
			b.stats.IncBlockSignatureValidationFail()
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

// verifyBlockMessage verifies the block message (payload.Hash)
func (b *backend) verifyBlockMessage(blockHashPayload *payload.Hash) *common.AppError {
	blockID := blockHashPayload.Hash
	_, err := b.blockClient.GetAcceptedBlock(context.TODO(), blockID)
	if err != nil {
		return &common.AppError{
			Code:    VerifyErrCode,
			Message: fmt.Sprintf("failed to get block %s: %s", blockID, err.Error()),
		}
	}

	return nil
}

// verifyAddressedCall verifies the addressed call message
func (b *backend) verifyAddressedCall(addressedCall *payload.AddressedCall) *common.AppError {
	// Further, parse the payload to see if it is a known type.
	parsed, err := messages.Parse(addressedCall.Payload)
	if err != nil {
		return &common.AppError{
			Code:    ParseErrCode,
			Message: "failed to parse addressed call message: " + err.Error(),
		}
	}

	switch p := parsed.(type) {
	case *messages.ValidatorUptime:
		return b.verifyUptimeMessage(p)
	default:
		return &common.AppError{
			Code:    ParseErrCode,
			Message: fmt.Sprintf("unknown message type: %T", p),
		}
	}
}

func (b *backend) verifyUptimeMessage(uptimeMsg *messages.ValidatorUptime) *common.AppError {
	// first get the validator's nodeID
	nodeID, err := b.validatorState.GetNodeID(uptimeMsg.ValidationID)
	if err != nil {
		return &common.AppError{
			Code:    VerifyErrCode,
			Message: fmt.Sprintf("failed to get validator for validationID %s: %s", uptimeMsg.ValidationID, err.Error()),
		}
	}

	// then get the current uptime
	currentUptime, _, err := b.uptimeCalculator.CalculateUptime(nodeID)
	if err != nil {
		return &common.AppError{
			Code:    VerifyErrCode,
			Message: fmt.Sprintf("failed to calculate uptime for nodeID %s: %s", nodeID, err.Error()),
		}
	}

	currentUptimeSeconds := uint64(currentUptime.Seconds())
	// verify the current uptime against the total uptime in the message
	if currentUptimeSeconds < uptimeMsg.TotalUptime {
		return &common.AppError{
			Code:    VerifyErrCode,
			Message: fmt.Sprintf("current uptime %d is less than queried uptime %d for nodeID %s", currentUptimeSeconds, uptimeMsg.TotalUptime, nodeID),
		}
	}

	return nil
}
