package message

import (
	"context"
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
)

var _ CrossChainRequest = EthCallRequest{}

// EthCallRequest has the JSON Data necessary to execute a new EVM call on the blockchain
type EthCallRequest struct {
	RequestArgs []byte `serialize:"true"`
}

// EthCallResponse represents the JSON return value of the executed EVM call
type EthCallResponse struct {
	ExecutionResult []byte `serialize:"true"`
}

// String converts EthCallRequest to a string
func (e EthCallRequest) String() string {
	return fmt.Sprintf("%#v", e)
}

// Handle returns the encoded EthCallResponse by executing EVM call with the given EthCallRequest
func (e EthCallRequest) Handle(ctx context.Context, requestingChainID ids.ID, requestID uint32, handler CrossChainRequestHandler) ([]byte, error) {
	return nil, nil
}
