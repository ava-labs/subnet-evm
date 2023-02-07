// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package stats

import (
	"sync"

	warpStats "github.com/ava-labs/subnet-evm/handlers/warp/stats"
	syncStats "github.com/ava-labs/subnet-evm/sync/handlers/stats"
)

var _ HandlerStats = &MockHandlerStats{}

// MockHandlerStats is mock for capturing and asserting on handler metrics in test
type MockHandlerStats struct {
	lock sync.Mutex

	syncStats.MockHandlerStats
	warpStats.MockSignatureRequestHandlerStats
}

func (m *MockHandlerStats) Reset() {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.MockHandlerStats.Reset()
	m.MockSignatureRequestHandlerStats.Reset()
}
