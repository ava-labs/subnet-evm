// (c) 2022-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"context"
	"time"

	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
)

var (
	_ Backend = (*NoopBackend)(nil)
	_ Backend = (*MockBackend)(nil)
)

// Backend defines the interface for precompiles to interact with vm backends.
type Backend interface {
	// AddMessage adds an unsigned message to the warp backend database
	AddMessage(ctx context.Context, unsignedMessage *warp.UnsignedMessage) error
}

type MockBackend struct {
	Stats MockHandlerStats
}

func (m *MockBackend) AddMessage(ctx context.Context, unsignedMessage *warp.UnsignedMessage) error {
	startTime := time.Now()
	m.Stats.IncAddSignature()

	// Always report signature request time
	defer func() {
		m.Stats.UpdateAddSignatureProcessingTime(time.Since(startTime))
	}()

	return nil
}

type NoopBackend struct{}

func (n NoopBackend) AddMessage(ctx context.Context, unsignedMessage *warp.UnsignedMessage) error {
	return nil
}
