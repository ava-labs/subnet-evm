// (c) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/codec"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/network/p2p/acp118"
	"github.com/ava-labs/avalanchego/snow/engine/common"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/payload"
	"github.com/ava-labs/subnet-evm/warp"
)

var (
	_ avalancheWarp.Signer = (*p2pSignerVerifier)(nil)
	_ acp118.Verifier      = (*p2pSignerVerifier)(nil)
)

const (
	ParseErrCode = iota
	GetSigErrCode
	MarshalErrCode
	ValidateErrCode
)

var (
	errUnknownPayloadType = fmt.Errorf("unknown payload type")
	errFailedToParse      = fmt.Errorf("failed to parse payload")
	errFailedToGetSig     = fmt.Errorf("failed to get signature")
)

type SignerVerifier interface {
	acp118.Verifier
	avalancheWarp.Signer
}

// p2pSignerVerifier serves warp signature requests using the p2p
// framework from avalanchego. It is a peer.RequestHandler for
// message.MessageSignatureRequest.
type p2pSignerVerifier struct {
	backend warp.Backend
	codec   codec.Manager
	stats   *handlerStats
}

func NewSignerVerifier(backend warp.Backend, codec codec.Manager) SignerVerifier {
	return &p2pSignerVerifier{
		backend: backend,
		codec:   codec,
		stats:   newStats(),
	}
}

func (s *p2pSignerVerifier) Verify(_ context.Context, unsignedMessage *avalancheWarp.UnsignedMessage, _ []byte) *common.AppError {
	parsed, err := payload.Parse(unsignedMessage.Payload)
	if err != nil {
		return &common.AppError{
			Code:    ParseErrCode,
			Message: "failed to parse payload: " + err.Error(),
		}
	}

	switch p := parsed.(type) {
	case *payload.AddressedCall:
		err = s.backend.ValidateMessage(unsignedMessage)
		if err != nil {
			s.stats.IncMessageSignatureValidationFail()
			return &common.AppError{
				Code:    ValidateErrCode,
				Message: "failed to validate message: " + err.Error(),
			}
		}
	case *payload.Hash:
		err = s.backend.ValidateBlockMessage(p.Hash)
		if err != nil {
			s.stats.IncBlockSignatureValidationFail()
			return &common.AppError{
				Code:    ValidateErrCode,
				Message: "failed to validate block message: " + err.Error(),
			}
		}
	default:
		return &common.AppError{
			Code:    ParseErrCode,
			Message: fmt.Sprintf("unknown payload type: %T", p),
		}
	}
	return nil
}

func (s *p2pSignerVerifier) Sign(unsignedMessage *avalancheWarp.UnsignedMessage) ([]byte, error) {
	parsed, err := payload.Parse(unsignedMessage.Payload)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errFailedToParse, err.Error())
	}

	var sig []byte
	switch p := parsed.(type) {
	case *payload.AddressedCall:
		sig, err = s.GetMessageSignature(unsignedMessage)
		if err != nil {
			s.stats.IncMessageSignatureMiss()
		} else {
			s.stats.IncMessageSignatureHit()
		}
	case *payload.Hash:
		sig, err = s.GetBlockSignature(p.Hash)
		if err != nil {
			s.stats.IncBlockSignatureMiss()
		} else {
			s.stats.IncBlockSignatureHit()
		}
	default:
		return nil, fmt.Errorf("%w: %T", errUnknownPayloadType, p)
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %s", errFailedToGetSig, err.Error())
	}

	return sig, nil
}

func (s *p2pSignerVerifier) GetMessageSignature(message *avalancheWarp.UnsignedMessage) ([]byte, error) {
	startTime := time.Now()
	s.stats.IncMessageSignatureRequest()

	// Always report signature request time
	defer func() {
		s.stats.UpdateMessageSignatureRequestTime(time.Since(startTime))
	}()

	// TODO: consider changing backend to return []byte
	sig, err := s.backend.GetMessageSignature(message)
	return sig[:], err
}

func (s *p2pSignerVerifier) GetBlockSignature(blockID ids.ID) ([]byte, error) {
	startTime := time.Now()
	s.stats.IncBlockSignatureRequest()

	// Always report signature request time
	defer func() {
		s.stats.UpdateBlockSignatureRequestTime(time.Since(startTime))
	}()

	// TODO: consider changing backend to return []byte
	sig, err := s.backend.GetBlockSignature(blockID)
	return sig[:], err
}
