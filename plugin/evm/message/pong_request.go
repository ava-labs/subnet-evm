// (c) 2021-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package message

import (
	"context"
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
)

var _ CrossChainRequest = PongRequest{}

type PongRequest struct {
	RequestBytes []byte `serialize:"true"`
}

type PongResponse struct {
	ResponseBytes []byte `serialize:"true"`
}

func (p PongRequest) String() string {
	return fmt.Sprintf("%#v", p)
}

// Handle returns the encoded EthCallResponse by executing EVM call with the given EthCallRequest
func (p PongRequest) Handle(ctx context.Context, requestingChainID ids.ID, requestID uint32, handler CrossChainRequestHandler) ([]byte, error) {
	return nil, nil
}
