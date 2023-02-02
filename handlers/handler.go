// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package handlers

import (
	"context"

	"github.com/ava-labs/avalanchego/codec"
	"github.com/ava-labs/subnet-evm/handlers/stats"
	"github.com/ava-labs/subnet-evm/metrics"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/subnet-evm/plugin/evm/message"
	"github.com/ava-labs/subnet-evm/sync/handlers"
	syncStats "github.com/ava-labs/subnet-evm/sync/handlers/stats"
	"github.com/ava-labs/subnet-evm/trie"
)

var _ message.RequestHandler = &networkHandler{}

type networkHandler struct {
	stateTrieLeafsRequestHandler *handlers.LeafsRequestHandler
	blockRequestHandler          *handlers.BlockRequestHandler
	codeRequestHandler           *handlers.CodeRequestHandler
	signatureRequestHandler      SignatureRequestHandler
}

func NewNetworkHandler(
	provider handlers.SyncDataProvider,
	evmTrieDB *trie.Database,
	networkCodec codec.Manager,
) message.RequestHandler {
	syncStats := syncStats.NewHandlerStats(metrics.Enabled)
	handlerStats := stats.NewHandlerStats(metrics.Enabled, syncStats)
	return &networkHandler{
		stateTrieLeafsRequestHandler: handlers.NewLeafsRequestHandler(evmTrieDB, provider, networkCodec, handlerStats),
		blockRequestHandler:          handlers.NewBlockRequestHandler(provider, networkCodec, handlerStats),
		codeRequestHandler:           handlers.NewCodeRequestHandler(evmTrieDB.DiskDB(), networkCodec, handlerStats),

		// TODO: initialize actual signature request handler when warp is ready
		signatureRequestHandler: &NoopSignatureRequestHandler{},
	}
}

func (n networkHandler) HandleTrieLeafsRequest(ctx context.Context, nodeID ids.NodeID, requestID uint32, leafsRequest message.LeafsRequest) ([]byte, error) {
	return n.stateTrieLeafsRequestHandler.OnLeafsRequest(ctx, nodeID, requestID, leafsRequest)
}

func (n networkHandler) HandleBlockRequest(ctx context.Context, nodeID ids.NodeID, requestID uint32, blockRequest message.BlockRequest) ([]byte, error) {
	return n.blockRequestHandler.OnBlockRequest(ctx, nodeID, requestID, blockRequest)
}

func (n networkHandler) HandleCodeRequest(ctx context.Context, nodeID ids.NodeID, requestID uint32, codeRequest message.CodeRequest) ([]byte, error) {
	return n.codeRequestHandler.OnCodeRequest(ctx, nodeID, requestID, codeRequest)
}

func (n networkHandler) HandleSignatureRequest(ctx context.Context, nodeID ids.NodeID, requestID uint32, signatureRequest message.SignatureRequest) ([]byte, error) {
	return n.signatureRequestHandler.OnSignatureRequest(ctx, nodeID, requestID, signatureRequest)
}
