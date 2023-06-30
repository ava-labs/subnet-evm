// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package load

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"regexp"
	"time"

	"github.com/ava-labs/subnet-evm/cmd/simulator/config"
	"github.com/ava-labs/subnet-evm/cmd/simulator/key"
	"github.com/ava-labs/subnet-evm/cmd/simulator/txs"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/sync/errgroup"
)

// ExecuteLoader creates txSequences from [config] and has txAgents execute the specified simulation.
func ExecuteLoader(ctx context.Context, config config.Config) error {
	if config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, config.Timeout)
		defer cancel()
	}

	// Construct the arguments for the load simulator
	clients := make([]ethclient.Client, 0, len(config.Endpoints))
	// Extract blockchainIDStr from the clientURI
	re := regexp.MustCompile(`bc\/(.*)\/`)
	matches := re.FindStringSubmatch(config.Endpoints[0])
	if len(matches) < 1 {
		return fmt.Errorf("failed to get blockchainIDStr from the clientURI %s", config.Endpoints[0])
	}
	// Get the last element in matches
	blockchainIDStr := matches[len(matches)-1]
	log.Info("Extracted blockchainIDStr from the clientURI", "blockchainIDStr", blockchainIDStr)

	re = regexp.MustCompile(`127.0.0.1:(.*)/ext/bc`)
	matches = re.FindStringSubmatch(config.Endpoints[0])
	if len(matches) < 1 {
		return fmt.Errorf("failed to get endpoint from the clientURI %s", config.Endpoints[0])
	}
	// Get the last element in matches
	endpoint := matches[len(matches)-1]
	log.Info("Extracted endpoint from the clientURI", "endpoint", endpoint)

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
	log.Info("Distributing funds", "numTxsPerWorker", config.TxsPerWorker, "minFunds", minFundsPerAddr)
	keys, err = DistributeFunds(ctx, clients[0], keys, config.Workers, minFundsPerAddr)
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

	bigGwei := big.NewInt(params.GWei)
	gasTipCap := new(big.Int).Mul(bigGwei, big.NewInt(config.MaxTipCap))
	gasFeeCap := new(big.Int).Mul(bigGwei, big.NewInt(config.MaxFeeCap))
	client := clients[0]
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch chainID: %w", err)
	}
	signer := types.LatestSignerForChainID(chainID)

	log.Info("Creating transaction sequences...")
	txGenerator := func(key *ecdsa.PrivateKey, nonce uint64) (*types.Transaction, error) {
		addr := ethcrypto.PubkeyToAddress(key.PublicKey)
		tx, err := types.SignNewTx(key, signer, &types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			GasTipCap: gasTipCap,
			GasFeeCap: gasFeeCap,
			Gas:       params.TxGas,
			To:        &addr,
			Data:      nil,
			Value:     common.Big0,
		})
		if err != nil {
			return nil, err
		}
		return tx, nil
	}
	txSequences, err := txs.GenerateTxSequences(ctx, txGenerator, clients[0], pks, config.TxsPerWorker)
	if err != nil {
		return err
	}

	log.Info("Constructing tx agents...", "numAgents", config.Workers)
	agents := make([]txs.Agent[*types.Transaction], 0, config.Workers)
	for i := 0; i < config.Workers; i++ {
		agents = append(agents, txs.NewIssueNAgent[*types.Transaction](txSequences[i], NewSingleAddressTxWorker(ctx, clients[i], senders[i]), config.BatchSize))
	}

	log.Info("Starting tx agents...")
	eg := errgroup.Group{}
	for _, agent := range agents {
		agent := agent
		eg.Go(func() error {
			return agent.Execute(ctx)
		})
	}

	log.Info("Waiting for tx agents...")
	if err := eg.Wait(); err != nil {
		return err
	}
	log.Info("Tx agents completed successfully.")

	logOtherMetrics(blockchainIDStr, endpoint)
	return nil
}

func logOtherMetrics(blockchainIDStr string, endpoint string) error {
	getCallStart := time.Now()
	metricsAPI := fmt.Sprintf("http://127.0.0.1:%s/ext/metrics", endpoint)
	resp, err := http.Get(metricsAPI)
	getCallDuration := time.Since(getCallStart)

	log.Info("GET Metrics API Data", "time", getCallDuration.Seconds())
	if err != nil {
		return fmt.Errorf("failed getting metrics: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed reading response body of metrics: %w", err)
	}

	bodyString := string(body)
	re := regexp.MustCompile(fmt.Sprintf(".*avalanche_%s_vm_metervm_build_block_sum.*", blockchainIDStr))
	matches := re.FindAllStringSubmatch(bodyString, -1)
	if len(matches) < 1 {
		log.Info("No build_block metrics found from metrics API for blockchainIDStr", "blockchainIDStr", blockchainIDStr, "metricsAPI", metricsAPI)
		return nil
	}
	log.Info("Sum of time (in ns) of a build_block", "time", matches[len(matches)-1])

	re = regexp.MustCompile(fmt.Sprintf(".*avalanche_%s_blks_accepted_sum.*", blockchainIDStr))
	matches = re.FindAllStringSubmatch(bodyString, -1)
	if len(matches) < 1 {
		log.Info("No accepted_block metrics found from metrics API for blockchainIDStr", "blockchainIDStr", blockchainIDStr, "metricsAPI", metricsAPI)
		return nil
	}
	log.Info("Sum of time (in ns) from issuance of a block(s) to its acceptance", "time", matches[len(matches)-1])

	re = regexp.MustCompile(fmt.Sprintf(".*avalanche_%s_vm_metervm_verify_sum.*", blockchainIDStr))
	matches = re.FindAllStringSubmatch(bodyString, -1)
	if len(matches) < 1 {
		log.Info("No verify metrics found from metrics API for blockchainIDStr", "blockchainIDStr", blockchainIDStr, "metricsAPI", metricsAPI)
		return nil
	}
	log.Info("Sum of time (in ns) of a verify", "time", matches[len(matches)-1])

	return nil
}
