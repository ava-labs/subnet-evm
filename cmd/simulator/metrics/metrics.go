// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	// Summary of the quantiles of Individual Issuance Tx Times
	IssuanceTxTimes prometheus.Summary
	// Summary of the quantiles of Individual Confirmation Tx Times
	ConfirmationTxTimes prometheus.Summary
	// Summary of the quantiles of Individual Issuance To Confirmation Tx Times
	IssuanceToConfirmationTxTimes prometheus.Summary
}

func NewSummary(name string, help string) prometheus.Summary {
	return prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       name,
		Help:       help,
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
}

// NewMetrics creates and returns a Metrics and registers it with a Collector
func NewMetrics(prefix string, reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		IssuanceTxTimes: NewSummary(
			fmt.Sprintf("%stx_issuance_time", prefix),
			"Individual Tx Issuance Times for a Load Test",
		),
		ConfirmationTxTimes: NewSummary(
			fmt.Sprintf("%stx_confirmation_time", prefix),
			"Individual Tx Confirmation Times for a Load Test",
		),
		IssuanceToConfirmationTxTimes: NewSummary(
			fmt.Sprintf("%stx_issuance_to_confirmation_time", prefix),
			"Individual Tx Issuance To Confirmation Times for a Load Test",
		),
	}
	reg.MustRegister(m.IssuanceTxTimes)
	reg.MustRegister(m.ConfirmationTxTimes)
	reg.MustRegister(m.IssuanceToConfirmationTxTimes)
	return m
}
