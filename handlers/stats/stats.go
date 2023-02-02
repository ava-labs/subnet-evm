// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package stats

import (
	"time"

	"github.com/ava-labs/subnet-evm/metrics"
	syncStats "github.com/ava-labs/subnet-evm/sync/handlers/stats"
)

var (
	_ HandlerStats = &handlerStats{}
	_ HandlerStats = &noopHandlerStats{}
)

// HandlerStats reports prometheus metrics for the network handlers
type HandlerStats interface {
	SignatureRequestHandlerStats
	syncStats.HandlerStats
}

type SignatureRequestHandlerStats interface {
	IncSignatureRequest()
	IncSignatureHit()
	IncSignatureMiss()
	UpdateSignatureRequestTime(duration time.Duration)
}

type handlerStats struct {
	syncHandlerStats syncStats.HandlerStats

	// SignatureRequestHandler metrics
	signatureRequest        metrics.Counter
	signatureHit            metrics.Counter
	signatureMiss           metrics.Counter
	signatureProcessingTime metrics.Timer
}

func (h *handlerStats) IncBlockRequest() {
	h.syncHandlerStats.IncBlockRequest()
}

func (h *handlerStats) IncMissingBlockHash() {
	h.syncHandlerStats.IncMissingBlockHash()
}

func (h *handlerStats) UpdateBlocksReturned(num uint16) {
	h.syncHandlerStats.UpdateBlocksReturned(num)
}

func (h *handlerStats) UpdateBlockRequestProcessingTime(duration time.Duration) {
	h.syncHandlerStats.UpdateBlockRequestProcessingTime(duration)
}

func (h *handlerStats) IncCodeRequest() {
	h.syncHandlerStats.IncCodeRequest()
}

func (h *handlerStats) IncMissingCodeHash() {
	h.syncHandlerStats.IncMissingCodeHash()
}

func (h *handlerStats) IncTooManyHashesRequested() {
	h.syncHandlerStats.IncTooManyHashesRequested()
}

func (h *handlerStats) IncDuplicateHashesRequested() {
	h.syncHandlerStats.IncDuplicateHashesRequested()
}

func (h *handlerStats) UpdateCodeReadTime(duration time.Duration) {
	h.syncHandlerStats.UpdateCodeReadTime(duration)
}

func (h *handlerStats) UpdateCodeBytesReturned(bytes uint32) {
	h.syncHandlerStats.UpdateCodeBytesReturned(bytes)
}

func (h *handlerStats) IncLeafsRequest() {
	h.syncHandlerStats.IncLeafsRequest()
}

func (h *handlerStats) IncInvalidLeafsRequest() {
	h.syncHandlerStats.IncInvalidLeafsRequest()
}

func (h *handlerStats) UpdateLeafsReturned(numLeafs uint16) {
	h.syncHandlerStats.UpdateLeafsReturned(numLeafs)
}

func (h *handlerStats) UpdateLeafsRequestProcessingTime(duration time.Duration) {
	h.syncHandlerStats.UpdateLeafsRequestProcessingTime(duration)
}

func (h *handlerStats) UpdateReadLeafsTime(duration time.Duration) {
	h.syncHandlerStats.UpdateReadLeafsTime(duration)
}

func (h *handlerStats) UpdateSnapshotReadTime(duration time.Duration) {
	h.syncHandlerStats.UpdateSnapshotReadTime(duration)
}

func (h *handlerStats) UpdateGenerateRangeProofTime(duration time.Duration) {
	h.syncHandlerStats.UpdateGenerateRangeProofTime(duration)
}

func (h *handlerStats) UpdateRangeProofValsReturned(numProofVals int64) {
	h.syncHandlerStats.UpdateRangeProofValsReturned(numProofVals)
}

func (h *handlerStats) IncMissingRoot() {
	h.syncHandlerStats.IncMissingRoot()
}

func (h *handlerStats) IncTrieError() {
	h.syncHandlerStats.IncTrieError()
}

func (h *handlerStats) IncProofError() {
	h.syncHandlerStats.IncProofError()
}

func (h *handlerStats) IncSnapshotReadError() {
	h.syncHandlerStats.IncSnapshotReadError()
}

func (h *handlerStats) IncSnapshotReadAttempt() {
	h.syncHandlerStats.IncSnapshotReadAttempt()
}

func (h *handlerStats) IncSnapshotReadSuccess() {
	h.syncHandlerStats.IncSnapshotReadSuccess()
}

func (h *handlerStats) IncSnapshotSegmentValid() {
	h.syncHandlerStats.IncSnapshotSegmentValid()
}

func (h *handlerStats) IncSnapshotSegmentInvalid() {
	h.syncHandlerStats.IncSnapshotSegmentInvalid()
}

func (h *handlerStats) IncSignatureRequest() { h.signatureRequest.Inc(1) }
func (h *handlerStats) IncSignatureHit()     { h.signatureHit.Inc(1) }
func (h *handlerStats) IncSignatureMiss()    { h.signatureMiss.Inc(1) }
func (h *handlerStats) UpdateSignatureRequestTime(duration time.Duration) {
	h.signatureProcessingTime.Update(duration)
}

func NewHandlerStats(enabled bool) HandlerStats {
	if !enabled {
		return NewNoopHandlerStats()
	}
	return &handlerStats{
		syncHandlerStats: syncStats.NewHandlerStats(enabled),

		// initialize signature request stats
		signatureRequest:        metrics.GetOrRegisterCounter("signature_request_count", nil),
		signatureHit:            metrics.GetOrRegisterCounter("signature_request_hit", nil),
		signatureMiss:           metrics.GetOrRegisterCounter("signature_request_miss", nil),
		signatureProcessingTime: metrics.GetOrRegisterTimer("signature_request_duration", nil),
	}
}

// no op implementation
type noopHandlerStats struct {
	syncHandlerStats syncStats.HandlerStats
}

func NewNoopHandlerStats() HandlerStats {
	return &noopHandlerStats{
		syncHandlerStats: syncStats.NewNoopHandlerStats(),
	}
}

func (m *noopHandlerStats) IncBlockRequest() {
	m.syncHandlerStats.IncBlockRequest()
}

func (m *noopHandlerStats) IncMissingBlockHash() {
	m.syncHandlerStats.IncMissingBlockHash()
}

func (m *noopHandlerStats) UpdateBlocksReturned(num uint16) {
	m.syncHandlerStats.UpdateBlocksReturned(num)
}

func (m *noopHandlerStats) UpdateBlockRequestProcessingTime(duration time.Duration) {
	m.syncHandlerStats.UpdateBlockRequestProcessingTime(duration)
}

func (m *noopHandlerStats) IncCodeRequest() {
	m.syncHandlerStats.IncCodeRequest()
}

func (m *noopHandlerStats) IncMissingCodeHash() {
	m.syncHandlerStats.IncMissingCodeHash()
}

func (m *noopHandlerStats) IncTooManyHashesRequested() {
	m.syncHandlerStats.IncTooManyHashesRequested()
}

func (m *noopHandlerStats) IncDuplicateHashesRequested() {
	m.syncHandlerStats.IncDuplicateHashesRequested()
}

func (m *noopHandlerStats) UpdateCodeReadTime(duration time.Duration) {
	m.syncHandlerStats.UpdateCodeReadTime(duration)
}

func (m *noopHandlerStats) UpdateCodeBytesReturned(bytes uint32) {
	m.syncHandlerStats.UpdateCodeBytesReturned(bytes)
}

func (m *noopHandlerStats) IncLeafsRequest() {
	m.syncHandlerStats.IncLeafsRequest()
}

func (m *noopHandlerStats) IncInvalidLeafsRequest() {
	m.syncHandlerStats.IncInvalidLeafsRequest()
}

func (m *noopHandlerStats) UpdateLeafsReturned(numLeafs uint16) {
	m.syncHandlerStats.UpdateLeafsReturned(numLeafs)
}

func (m *noopHandlerStats) UpdateLeafsRequestProcessingTime(duration time.Duration) {
	m.syncHandlerStats.UpdateLeafsRequestProcessingTime(duration)
}

func (m *noopHandlerStats) UpdateReadLeafsTime(duration time.Duration) {
	m.syncHandlerStats.UpdateReadLeafsTime(duration)
}

func (m *noopHandlerStats) UpdateSnapshotReadTime(duration time.Duration) {
	m.syncHandlerStats.UpdateSnapshotReadTime(duration)
}

func (m *noopHandlerStats) UpdateGenerateRangeProofTime(duration time.Duration) {
	m.syncHandlerStats.UpdateGenerateRangeProofTime(duration)
}

func (m *noopHandlerStats) UpdateRangeProofValsReturned(numProofVals int64) {
	m.syncHandlerStats.UpdateRangeProofValsReturned(numProofVals)
}

func (m *noopHandlerStats) IncMissingRoot() {
	m.syncHandlerStats.IncMissingRoot()
}

func (m *noopHandlerStats) IncTrieError() {
	m.syncHandlerStats.IncTrieError()
}

func (m *noopHandlerStats) IncProofError() {
	m.syncHandlerStats.IncProofError()
}

func (m *noopHandlerStats) IncSnapshotReadError() {
	m.syncHandlerStats.IncSnapshotReadError()
}

func (m *noopHandlerStats) IncSnapshotReadAttempt() {
	m.syncHandlerStats.IncSnapshotReadAttempt()
}

func (m *noopHandlerStats) IncSnapshotReadSuccess() {
	m.syncHandlerStats.IncSnapshotReadSuccess()
}

func (m *noopHandlerStats) IncSnapshotSegmentValid() {
	m.syncHandlerStats.IncSnapshotSegmentValid()
}

func (m *noopHandlerStats) IncSnapshotSegmentInvalid() {
	m.syncHandlerStats.IncSnapshotSegmentInvalid()
}

func (m *noopHandlerStats) IncSignatureRequest()                              {}
func (m *noopHandlerStats) IncSignatureHit()                                  {}
func (m *noopHandlerStats) IncSignatureMiss()                                 {}
func (m *noopHandlerStats) UpdateSignatureRequestTime(duration time.Duration) {}
