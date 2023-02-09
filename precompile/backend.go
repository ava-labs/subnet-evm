// (c) 2022-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"context"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/platformvm/teleporter"
	"github.com/ava-labs/subnet-evm/plugin/evm/warp"
)

var _ Backend = (*noopBackend)(nil)

// Backend defines the interface for precompiles to interact with vm backends.
type Backend interface {
	warp.WarpBackend
}

type noopBackend struct{}

func NewNoopBackend() Backend {
	return &noopBackend{}
}

func (n noopBackend) AddMessage(ctx context.Context, unsignedMessage *teleporter.UnsignedMessage) error {
	return nil
}

func (n noopBackend) GetSignature(ctx context.Context, messageHash ids.ID) ([]byte, error) {
	return nil, nil
}
