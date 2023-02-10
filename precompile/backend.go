// (c) 2022-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"context"

	"github.com/ava-labs/avalanchego/vms/platformvm/teleporter"
)

var _ Backend = (*noopBackend)(nil)

// Backend defines the interface for precompiles to interact with vm backends.
type Backend interface {
	// AddMessage adds an unsigned message to the warp backend database
	AddMessage(ctx context.Context, unsignedMessage *teleporter.UnsignedMessage) error
}

type noopBackend struct{}

func NewNoopBackend() Backend {
	return &noopBackend{}
}

func (n noopBackend) AddMessage(ctx context.Context, unsignedMessage *teleporter.UnsignedMessage) error {
	return nil
}
