// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package stats

import (
	"time"

	evmmetrics "github.com/ava-labs/libevm/metrics" // alias to avoid name collision with local 'metrics' in this package
)

// RequestHandlerStats provides the interface for metrics for app requests.
type RequestHandlerStats interface {
	UpdateTimeUntilDeadline(duration time.Duration)
	IncDeadlineDroppedRequest()
}

type requestHandlerStats struct {
	timeUntilDeadline evmmetrics.Timer
	droppedRequests   evmmetrics.Counter
}

func (h *requestHandlerStats) IncDeadlineDroppedRequest() {
	h.droppedRequests.Inc(1)
}

func (h *requestHandlerStats) UpdateTimeUntilDeadline(duration time.Duration) {
	h.timeUntilDeadline.Update(duration)
}

func NewRequestHandlerStats() RequestHandlerStats {
	return &requestHandlerStats{
		timeUntilDeadline: evmmetrics.GetOrRegisterTimer("net_req_time_until_deadline", nil),
		droppedRequests:   evmmetrics.GetOrRegisterCounter("net_req_deadline_dropped", nil),
	}
}
