// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package stats

import (
	"sync"
	"time"

	syncStats "github.com/ava-labs/subnet-evm/sync/handlers/stats"
)

var _ HandlerStats = &MockHandlerStats{}

// MockHandlerStats is mock for capturing and asserting on handler metrics in test
type MockHandlerStats struct {
	lock sync.Mutex

	syncHandlerStats syncStats.MockHandlerStats

	SignatureRequestCount,
	SignatureRequestHit,
	SignatureRequestMiss uint32
	SignatureRequestDuration time.Duration
}

func (m *MockHandlerStats) Reset() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.syncHandlerStats.Reset()

	m.SignatureRequestCount = 0
	m.SignatureRequestHit = 0
	m.SignatureRequestMiss = 0
	m.SignatureRequestDuration = 0
}

func (m *MockHandlerStats) IncBlockRequest() {
	m.syncHandlerStats.IncBlockRequest()
}

func (m *MockHandlerStats) IncMissingBlockHash() {
	m.syncHandlerStats.IncMissingBlockHash()
}

func (m *MockHandlerStats) UpdateBlocksReturned(num uint16) {
	m.syncHandlerStats.UpdateBlocksReturned(num)
}

func (m *MockHandlerStats) UpdateBlockRequestProcessingTime(duration time.Duration) {
	m.syncHandlerStats.UpdateBlockRequestProcessingTime(duration)
}

func (m *MockHandlerStats) IncCodeRequest() {
	m.syncHandlerStats.IncCodeRequest()
}

func (m *MockHandlerStats) IncMissingCodeHash() {
	m.syncHandlerStats.IncMissingCodeHash()
}

func (m *MockHandlerStats) IncTooManyHashesRequested() {
	m.syncHandlerStats.IncTooManyHashesRequested()
}

func (m *MockHandlerStats) IncDuplicateHashesRequested() {
	m.syncHandlerStats.IncDuplicateHashesRequested()
}

func (m *MockHandlerStats) UpdateCodeReadTime(duration time.Duration) {
	m.syncHandlerStats.UpdateCodeReadTime(duration)
}

func (m *MockHandlerStats) UpdateCodeBytesReturned(bytes uint32) {
	m.syncHandlerStats.UpdateCodeBytesReturned(bytes)
}

func (m *MockHandlerStats) IncLeafsRequest() {
	m.syncHandlerStats.IncLeafsRequest()
}

func (m *MockHandlerStats) IncInvalidLeafsRequest() {
	m.syncHandlerStats.IncInvalidLeafsRequest()
}

func (m *MockHandlerStats) UpdateLeafsReturned(numLeafs uint16) {
	m.syncHandlerStats.UpdateLeafsReturned(numLeafs)
}

func (m *MockHandlerStats) UpdateLeafsRequestProcessingTime(duration time.Duration) {
	m.syncHandlerStats.UpdateLeafsRequestProcessingTime(duration)
}

func (m *MockHandlerStats) UpdateReadLeafsTime(duration time.Duration) {
	m.syncHandlerStats.UpdateReadLeafsTime(duration)
}

func (m *MockHandlerStats) UpdateSnapshotReadTime(duration time.Duration) {
	m.syncHandlerStats.UpdateSnapshotReadTime(duration)
}

func (m *MockHandlerStats) UpdateGenerateRangeProofTime(duration time.Duration) {
	m.syncHandlerStats.UpdateGenerateRangeProofTime(duration)
}

func (m *MockHandlerStats) UpdateRangeProofValsReturned(numProofVals int64) {
	m.syncHandlerStats.UpdateRangeProofValsReturned(numProofVals)
}

func (m *MockHandlerStats) IncMissingRoot() {
	m.syncHandlerStats.IncMissingRoot()
}

func (m *MockHandlerStats) IncTrieError() {
	m.syncHandlerStats.IncTrieError()
}

func (m *MockHandlerStats) IncProofError() {
	m.syncHandlerStats.IncProofError()
}

func (m *MockHandlerStats) IncSnapshotReadError() {
	m.syncHandlerStats.IncSnapshotReadError()
}

func (m *MockHandlerStats) IncSnapshotReadAttempt() {
	m.syncHandlerStats.IncSnapshotReadAttempt()
}

func (m *MockHandlerStats) IncSnapshotReadSuccess() {
	m.syncHandlerStats.IncSnapshotReadSuccess()
}

func (m *MockHandlerStats) IncSnapshotSegmentValid() {
	m.syncHandlerStats.IncSnapshotSegmentValid()
}

func (m *MockHandlerStats) IncSnapshotSegmentInvalid() {
	m.syncHandlerStats.IncSnapshotSegmentInvalid()
}

func (m *MockHandlerStats) IncSignatureRequest() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.SignatureRequestCount++
}

func (m *MockHandlerStats) IncSignatureHit() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.SignatureRequestHit++
}

func (m *MockHandlerStats) IncSignatureMiss() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.SignatureRequestMiss++
}

func (m *MockHandlerStats) UpdateSignatureRequestTime(duration time.Duration) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.SignatureRequestDuration += duration
}
