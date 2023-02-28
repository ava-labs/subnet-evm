// (c) 2022-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"sync"
	"time"
)

var _ HandlerStats = (*MockHandlerStats)(nil)

type HandlerStats interface {
	IncAddSignature()
	UpdateAddSignatureProcessingTime(duration time.Duration)
}

type MockHandlerStats struct {
	lock sync.Mutex

	AddSignatureCount             uint32
	AddSignatureProcessingTimeSum time.Duration
}

func (m *MockHandlerStats) Reset() {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.AddSignatureCount = 0
	m.AddSignatureProcessingTimeSum = 0
}

func (m *MockHandlerStats) IncAddSignature() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.AddSignatureCount++
}

func (m *MockHandlerStats) UpdateAddSignatureProcessingTime(duration time.Duration) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.AddSignatureProcessingTimeSum += duration
}
