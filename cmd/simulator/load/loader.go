// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package load

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/ava-labs/subnet-evm/cmd/simulator/config"
	"github.com/ava-labs/subnet-evm/cmd/simulator/key"
	"github.com/ava-labs/subnet-evm/cmd/simulator/metrics"
	"github.com/ava-labs/subnet-evm/cmd/simulator/txs"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
)

const (
	MetricsEndpoint = "/metrics" // Endpoint for the Prometheus Metrics Server
)

// ExecuteLoader creates txSequences from [config] and has txAgents execute the specified simulation.
func ExecuteLoader(ctx context.Context, config config.Config) error {
	if config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, config.Timeout)
		defer cancel()
	}

	// Create buffered sigChan to receive SIGINT notifications
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	// Create context with cancel
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		// Blocks until we receive a SIGINT notification or if parent context is done
		select {
		case <-sigChan:
		case <-ctx.Done():
		}

		// Cancel the child context and end all processes
		cancel()
	}()

	// Construct the arguments for the load simulator
	clients := make([]ethclient.Client, 0, len(config.Endpoints))
	for i := 0; i < config.Workers; i++ {
		clientURI := config.Endpoints[i%len(config.Endpoints)]
		client, err := ethclient.Dial(clientURI)
		if err != nil {
			return fmt.Errorf("failed to dial client at %s: %w", clientURI, err)
		}
		clients = append(clients, client)
	}

	keys, err := key.LoadAll(ctx, config.KeyDir)
	if err != nil {
		return err
	}
	// Ensure there are at least [config.Workers] keys and save any newly generated ones.
	if len(keys) < config.Workers {
		for i := 0; len(keys) < config.Workers; i++ {
			newKey, err := key.Generate()
			if err != nil {
				return fmt.Errorf("failed to generate %d new key: %w", i, err)
			}
			if err := newKey.Save(config.KeyDir); err != nil {
				return fmt.Errorf("failed to save %d new key: %w", i, err)
			}
			keys = append(keys, newKey)
		}
	}

	// Each address needs: params.GWei * MaxFeeCap * params.TxGas * TxsPerWorker total wei
	// to fund gas for all of their transactions.
	maxFeeCap := new(big.Int).Mul(big.NewInt(params.GWei), big.NewInt(config.MaxFeeCap))
	minFundsPerAddr := new(big.Int).Mul(maxFeeCap, big.NewInt(int64(config.TxsPerWorker*params.TxGas)))

	// Create metrics
	reg := prometheus.NewRegistry()
	m := metrics.NewMetrics(reg)
	metricsPort := strconv.Itoa(int(config.MetricsPort))

	log.Info("Distributing funds", "numTxsPerWorker", config.TxsPerWorker, "minFunds", minFundsPerAddr)
	keys, err = DistributeFunds(ctx, clients[0], keys, config.Workers, minFundsPerAddr, m)
	if err != nil {
		return err
	}
	log.Info("Distributed funds successfully")

	pks := make([]*ecdsa.PrivateKey, 0, len(keys))
	senders := make([]common.Address, 0, len(keys))
	for _, key := range keys {
		pks = append(pks, key.PrivKey)
		senders = append(senders, key.Address)
	}
	client := clients[0]
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch chainID: %w", err)
	}
	log.Info("Constructing tx agents...", "numAgents", config.Workers)
	agentBuilder, err := NewTransferTxAgentBuilder(ctx, config, chainID, pks, client, m)
	if err != nil {
		return err
	}
	agents := make([]txs.Agent, 0, config.Workers)
	for i := 0; i < config.Workers; i++ {
		agent, err := agentBuilder.NewAgent(ctx, i, clients[i], senders[i])
		if err != nil {
			return err
		}
		agents = append(agents, agent)
	}

	log.Info("Starting tx agents...")
	eg := errgroup.Group{}
	for _, agent := range agents {
		agent := agent
		eg.Go(func() error {
			return agent.Execute(ctx)
		})
	}

	go startMetricsServer(ctx, metricsPort, reg)

	log.Info("Waiting for tx agents...")
	if err := eg.Wait(); err != nil {
		return err
	}
	log.Info("Tx agents completed successfully.")

	printOutputFromMetricsServer(metricsPort)
	return nil
}

func startMetricsServer(ctx context.Context, metricsPort string, reg *prometheus.Registry) {
	// Create a prometheus server to expose individual tx metrics
	server := &http.Server{
		Addr: fmt.Sprintf(":%s", metricsPort),
	}

	// Start up go routine to listen for SIGINT notifications to gracefully shut down server
	go func() {
		// Blocks until signal is received
		<-ctx.Done()

		if err := server.Shutdown(ctx); err != nil {
			log.Error("Metrics server error: %v", err)
		}
		log.Info("Received a SIGINT signal: Gracefully shutting down metrics server")
	}()

	// Start metrics server
	http.Handle(MetricsEndpoint, promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Info(fmt.Sprintf("Metrics Server: localhost:%s%s", metricsPort, MetricsEndpoint))
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Error("Metrics server error: %v", err)
	}
}

func printOutputFromMetricsServer(metricsPort string) {
	// Get response from server
	resp, err := http.Get(fmt.Sprintf("http://localhost:%s%s", metricsPort, MetricsEndpoint))
	if err != nil {
		log.Error("cannot get response from metrics servers", "err", err)
		return
	}
	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("cannot read response body", "err", err)
		return
	}
	// Print out formatted individual metrics
	parts := strings.Split(string(respBody), "\n")
	for _, s := range parts {
		fmt.Printf("       \t\t\t%s\n", s)
	}
}
