// (c) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package handlers

import (
	"context"
	"time"

	"github.com/ava-labs/avalanchego/codec"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/network/p2p"
	"github.com/ava-labs/avalanchego/proto/pb/sdk"
	"github.com/ava-labs/avalanchego/snow/engine/common"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/warp"
	"google.golang.org/protobuf/proto"
)

var _ p2p.Handler = (*SignatureRequestHandlerP2P)(nil)

const (
	ErrFailedToParse = iota
	ErrFailedToGetSig
	ErrFailedToMarshal
)

// SignatureRequestHandlerP2P serves warp signature requests using the p2p
// framework from avalanchego. It is a peer.RequestHandler for
// message.MessageSignatureRequest.
type SignatureRequestHandlerP2P struct {
	backend warp.Backend
	codec   codec.Manager
	stats   *handlerStats
}

func NewSignatureRequestHandlerP2P(backend warp.Backend, codec codec.Manager) *SignatureRequestHandlerP2P {
	return &SignatureRequestHandlerP2P{
		backend: backend,
		codec:   codec,
		stats:   newStats(),
	}
}

func (s *SignatureRequestHandlerP2P) AppRequest(
	ctx context.Context,
	nodeID ids.NodeID,
	deadline time.Time,
	requestBytes []byte,
) ([]byte, *common.AppError) {
	// Per ACP-118, the requestBytes are the serialized form of
	// sdk.SignatureRequest.
	req := new(sdk.SignatureRequest)
	if err := proto.Unmarshal(requestBytes, req); err != nil {
		return nil, &common.AppError{
			Code:    ErrFailedToParse,
			Message: "failed to unmarshal request: " + err.Error(),
		}
	}

	var (
		sig [bls.SignatureLen]byte
		err error
	)
	if len(req.Message) == ids.IDLen {
		// Treat requests with 32-byte [message] as a block hash.
		var blockID ids.ID
		copy(blockID[:], req.Message)
		sig, err = s.GetBlockSignature(blockID)
		if err != nil {
			s.stats.IncBlockSignatureMiss()
		} else {
			s.stats.IncBlockSignatureHit()
		}
	} else {
		// Otherwise, treat the request as an unsigned warp message.
		sig, err = s.GetMessageSignature(req.Message)
		if err != nil {
			s.stats.IncMessageSignatureMiss()
		} else {
			s.stats.IncMessageSignatureHit()
		}
	}
	if err != nil {
		return nil, &common.AppError{
			Code:    ErrFailedToGetSig,
			Message: "failed to get signature: " + err.Error(),
		}
	}

	// Per ACP-118, the responseBytes are the serialized form of
	// sdk.SignatureResponse.
	resp := &sdk.SignatureResponse{Signature: sig[:]}
	respBytes, err := proto.Marshal(resp)
	if err != nil {
		return nil, &common.AppError{
			Code:    ErrFailedToMarshal,
			Message: "failed to marshal response: " + err.Error(),
		}
	}
	return respBytes, nil
}

func (s *SignatureRequestHandlerP2P) GetMessageSignature(message []byte) ([bls.SignatureLen]byte, error) {
	startTime := time.Now()
	s.stats.IncMessageSignatureRequest()

	// Always report signature request time
	defer func() {
		s.stats.UpdateMessageSignatureRequestTime(time.Since(startTime))
	}()

	unsignedMessage, err := avalancheWarp.ParseUnsignedMessage(message)
	if err != nil {
		return [bls.SignatureLen]byte{}, err
	}

	messageID := unsignedMessage.ID()
	return s.backend.GetMessageSignature(messageID)
}

func (s *SignatureRequestHandlerP2P) GetBlockSignature(blockID ids.ID) ([bls.SignatureLen]byte, error) {
	startTime := time.Now()
	s.stats.IncBlockSignatureRequest()

	// Always report signature request time
	defer func() {
		s.stats.UpdateBlockSignatureRequestTime(time.Since(startTime))
	}()

	return s.backend.GetBlockSignature(blockID)
}

func (s *SignatureRequestHandlerP2P) CrossChainAppRequest(
	ctx context.Context,
	chainID ids.ID,
	deadline time.Time,
	requestBytes []byte,
) ([]byte, error) {
	return nil, nil
}

func (s *SignatureRequestHandlerP2P) AppGossip(
	ctx context.Context, nodeID ids.NodeID, gossipBytes []byte) {
}
