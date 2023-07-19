package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	// Summary of the quantiles of Individual Confirmation Tx Times
	ConfirmationTxTimes prometheus.Summary
	// Summary of the quantiles of Individual Issuance Tx Times
	IssuanceTxTimes prometheus.Summary
}

// NewMetrics creates and returns a metrics and registers it with a Collector
func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		ConfirmationTxTimes: prometheus.NewSummary(prometheus.SummaryOpts{
			Name:       "confirmation_tx_times",
			Help:       "Individual Tx Confirmation Times for a Load Test",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
		IssuanceTxTimes: prometheus.NewSummary(prometheus.SummaryOpts{
			Name:       "issuance_tx_times",
			Help:       "Individual Tx Issuance Times for a Load Test",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
	}
	reg.MustRegister(m.ConfirmationTxTimes)
	reg.MustRegister(m.IssuanceTxTimes)
	return m
}
