// (c) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package message

import (
	"context"

	"github.com/ethereum/go-ethereum/log"

	"github.com/ava-labs/avalanchego/ids"
)

var _ GossipHandler = NoopMempoolGossipHandler{}

// GossipHandler handles incoming gossip messages
type GossipHandler interface {
	HandleTxs(nodeID ids.NodeID, msg TxsGossip) error
}

type NoopMempoolGossipHandler struct{}

func (NoopMempoolGossipHandler) HandleTxs(nodeID ids.NodeID, _ TxsGossip) error {
	log.Debug("dropping unexpected Txs message", "peerID", nodeID)
	return nil
}

// RequestHandler interface handles incoming requests from peers
// Must have methods in format of handleType(context.Context, ids.NodeID, uint32, request Type) error
// so that the Request object of relevant Type can invoke its respective handle method
// on this struct.
// Also see GossipHandler for implementation style.
type RequestHandler interface {
	HandleTrieLeafsRequest(ctx context.Context, nodeID ids.NodeID, requestID uint32, leafsRequest LeafsRequest) ([]byte, error)
	HandleBlockRequest(ctx context.Context, nodeID ids.NodeID, requestID uint32, request BlockRequest) ([]byte, error)
	HandleCodeRequest(ctx context.Context, nodeID ids.NodeID, requestID uint32, codeRequest CodeRequest) ([]byte, error)
}

// ResponseHandler handles response for a sent request
// Only one of OnResponse or OnFailure is called for a given requestID, not both
type ResponseHandler interface {
	// OnResponse is invoked when the peer responded to a request
	OnResponse(nodeID ids.NodeID, requestID uint32, response []byte) error
	// OnFailure is invoked when there was a failure in processing a request
	// The FailureReason outlines the underlying cause.
	OnFailure(nodeID ids.NodeID, requestID uint32) error
}
