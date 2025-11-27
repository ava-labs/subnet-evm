// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package metrics

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time" // Added for explicit timeout definitions

	"github.com/ava-labs/libevm/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ava-labs/subnet-evm/rpc"
)

// Metrics holds the Prometheus Registry and the specific metric collectors.
type Metrics struct {
	reg *prometheus.Registry
	// Summary of the quantiles of Individual Issuance Tx Times
	IssuanceTxTimes prometheus.Summary
	// Summary of the quantiles of Individual Confirmation Tx Times
	ConfirmationTxTimes prometheus.Summary
	// Summary of the quantiles of Individual Issuance To Confirmation Tx Times
	IssuanceToConfirmationTxTimes prometheus.Summary
}

// NewDefaultMetrics creates a standard Prometheus registry and initializes metrics with it.
func NewDefaultMetrics() *Metrics {
	registry := prometheus.NewRegistry()
	return NewMetrics(registry)
}

// NewMetrics creates and returns a Metrics object and registers all summary metrics
// with the provided Prometheus Collector registry.
func NewMetrics(reg *prometheus.Registry) *Metrics {
	m := &Metrics{
		reg: reg,
		IssuanceTxTimes: prometheus.NewSummary(prometheus.SummaryOpts{
			Name:       "tx_issuance_time",
			Help:       "Individual Tx Issuance Times for a Load Test",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}, // 50th, 90th, 99th percentile
		}),
		ConfirmationTxTimes: prometheus.NewSummary(prometheus.SummaryOpts{
			Name:       "tx_confirmation_time",
			Help:       "Individual Tx Confirmation Times for a Load Test",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
		IssuanceToConfirmationTxTimes: prometheus.NewSummary(prometheus.SummaryOpts{
			Name:       "tx_issuance_to_confirmation_time",
			Help:       "Individual Tx Issuance To Confirmation Times for a Load Test",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
	}
	// MustRegister panics if registration fails, ensuring metrics are correctly initialized.
	reg.MustRegister(m.IssuanceTxTimes)
	reg.MustRegister(m.ConfirmationTxTimes)
	reg.MustRegister(m.IssuanceToConfirmationTxTimes)
	return m
}

// MetricsServer encapsulates the control flow for the running HTTP server.
type MetricsServer struct {
	cancel context.CancelFunc // Function to signal server shutdown
	stopCh chan struct{}      // Channel closed when the server goroutine exits
}

// Serve starts the HTTP server to expose Prometheus metrics.
func (m *Metrics) Serve(ctx context.Context, metricsPort string, metricsEndpoint string) *MetricsServer {
	// Create a cancellable context for server control
	ctx, cancel := context.WithCancel(ctx)
	
	// Create the HTTP server instance with explicit timeouts for robustness and security.
	server := &http.Server{
		Addr:              ":" + metricsPort,
		ReadHeaderTimeout: rpc.DefaultHTTPTimeouts.ReadHeaderTimeout,
		WriteTimeout:      5 * time.Second, // Added: Timeout for response writes
		IdleTimeout:       30 * time.Second, // Added: Max time to wait for the next request
	}

	// Go routine to listen for cancellation (e.g., SIGINT or external call to Shutdown)
	go func() {
		// Blocks until the context is cancelled
		<-ctx.Done() 

		// Shutdown the server gracefully using the cancellation context
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second) // Added 5s timeout for shutdown
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			// Use the log package's formatting capability
			log.Error("Metrics server shutdown error: %v", err)
		}
		log.Info("Metrics server: Gracefully shutting down.")
	}()

	// Start the metrics server goroutine
	ms := &MetricsServer{
		stopCh: make(chan struct{}),
		cancel: cancel,
	}
	go func() {
		// Ensure the stop channel is closed when the goroutine finishes
		defer close(ms.stopCh)

		// Set up the Prometheus handler on the specified endpoint
		http.Handle(metricsEndpoint, promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})) // Removed redundant Registry option

		// Log the server startup info using format logging
		log.Info("Metrics Server started at: localhost:%s%s", metricsPort, metricsEndpoint)
		
		// ListenAndServe blocks. Check the error explicitly.
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// Log critical errors that are not due to normal shutdown
			log.Error("Metrics server critical error: %v", err)
		}
	}()

	return ms
}

// Shutdown signals the MetricsServer to stop and waits for the underlying goroutine to exit.
func (ms *MetricsServer) Shutdown() {
	// Signal the cancellation
	ms.cancel()
	// Wait for the server goroutine to exit (by reading from stopCh)
	<-ms.stopCh
}

// Print gathers all metrics and outputs them to stdout or a JSON file.
func (m *Metrics) Print(outputFile string) error {
	metrics, err := m.reg.Gather()
	if err != nil {
		return fmt.Errorf("failed to gather metrics: %w", err) // Use %w for error wrapping
	}

	if outputFile == "" {
		// Printout to stdout
		fmt.Println("*** Metrics Report (STDOUT) ***")
		for _, mf := range metrics {
			// Use fmt.Printf for clear output formatting
			fmt.Printf("Metric Name: %s (Type: %s)\n", mf.GetName(), mf.GetType().String())
			fmt.Printf("Help: %s\n", mf.GetHelp())
			
			for _, m := range mf.GetMetric() {
				// Use JSON encoding for consistent metric value representation
				metricJSON, err := json.MarshalIndent(m, "  ", "  ")
				if err == nil {
					fmt.Printf("  Values:\n%s\n", string(metricJSON))
				} else {
					fmt.Printf("  Values: %s (Error Marshaling)\n", m.String())
				}
			}
			fmt.Println("-------------------------------")
		}
		fmt.Println("*******************************")
	} else {
		// Printout to a JSON file
		jsonFile, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file %q: %w", outputFile, err)
		}
		// Defer closing the file
		defer jsonFile.Close() 

		// Encode the collected metrics to the file
		if err := json.NewEncoder(jsonFile).Encode(metrics); err != nil {
			return fmt.Errorf("failed to encode metrics to JSON: %w", err)
		}
	}

	return nil
}
