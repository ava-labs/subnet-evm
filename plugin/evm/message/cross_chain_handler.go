// (c) 2021-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package message

import (
	"context"

	"github.com/ava-labs/avalanchego/codec"
	"github.com/ava-labs/avalanchego/ids"

	"github.com/ava-labs/coreth/internal/ethapi"

	"github.com/ethereum/go-ethereum/log"
)

var _ CrossChainRequestHandler = &crossChainHandler{}

// crossChainHandler implements the CrossChainRequestHandler interface
type crossChainHandler struct {
	backend         ethapi.Backend
	crossChainCodec codec.Manager
}

// NewCrossChainHandler creates and returns a new instance of CrossChainRequestHandler
func NewCrossChainHandler(b ethapi.Backend, codec codec.Manager) CrossChainRequestHandler {
	return &crossChainHandler{
		backend:         b,
		crossChainCodec: codec,
	}
}

// HandleEthCallRequests returns an encoded EthCallResponse to the given [ethCallRequest]
// This function executes EVM Call against the state associated with [rpc.AcceptedBlockNumber] with the given
// transaction call object [ethCallRequest].
// This function does not return an error as errors are treated as FATAL to the node.
func (c *crossChainHandler) HandleEthCallRequest(ctx context.Context, requestingChainID ids.ID, requestID uint32, ethCallRequest EthCallRequest) ([]byte, error) {
	// XXX: Don't care about this for now
	response := EthCallResponse{}

	responseBytes, err := c.crossChainCodec.Marshal(Version, response)
	if err != nil {
		log.Error("error occurred with marshalling EthCallResponse", "err", err, "EthCallResponse", response)
		return nil, nil
	}

	return responseBytes, nil
}
