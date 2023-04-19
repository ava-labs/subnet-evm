// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package load

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"

	"github.com/ava-labs/subnet-evm/cmd/simulator/config"
	"github.com/ava-labs/subnet-evm/cmd/simulator/key"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

func CreateLoader(ctx context.Context, config config.Config) (*WorkerGroup, error) {
	// Construct the arguments for the load simulator
	switch {
	case len(config.Endpoints) == 0:
		fmt.Printf("Must specify at least one clientURI\n")
		os.Exit(1)
	case len(config.Endpoints) < config.Workers:
		// Ensure there are at least [config.Workers] config.Endpoints by creating
		// duplicates as needed.
		for i := 0; len(config.Endpoints) < config.Workers; i++ {
			config.Endpoints = append(config.Endpoints, config.Endpoints[i])
		}
	}
	clients := make([]ethclient.Client, 0, len(config.Endpoints))
	for _, clientURI := range config.Endpoints {
		client, err := ethclient.Dial(clientURI)
		if err != nil {
			return nil, fmt.Errorf("failed to dial client at %s: %w", clientURI, err)
		}
		clients = append(clients, client)
	}

	keys, err := key.LoadAll(ctx, config.KeyDir)
	if err != nil {
		return nil, err
	}
	// Ensure there are at least [config.Workers] keys and save any newly generated ones.
	if len(keys) < config.Workers {
		for i := 0; len(keys) < config.Workers; i++ {
			newKey, err := key.Generate()
			if err != nil {
				return nil, fmt.Errorf("failed to generate %d new key: %w", i, err)
			}
			if err := newKey.Save(config.KeyDir); err != nil {
				return nil, fmt.Errorf("failed to save %d new key: %w", i, err)
			}
			keys = append(keys, newKey)
		}
	}

	maxFeeCap := new(big.Int).Mul(big.NewInt(params.GWei), big.NewInt(config.MaxFeeCap))
	minFundsPerAddr := new(big.Int).Mul(maxFeeCap, big.NewInt(int64(config.TxsPerWorker*params.TxGas)))
	log.Info("Distributing funds", "numTxsPerWorker", config.TxsPerWorker, "minFunds", minFundsPerAddr)
	keys, err = DistributeFunds(ctx, clients[0], keys, config.Workers, minFundsPerAddr)
	if err != nil {
		return nil, err
	}

	pks := make([]*ecdsa.PrivateKey, 0, len(keys))
	senders := make([]common.Address, 0, len(keys))
	for _, key := range keys {
		pks = append(pks, key.PrivKey)
		senders = append(senders, key.Address)
	}

	bigGwei := big.NewInt(params.GWei)
	gasTipCap := new(big.Int).Mul(bigGwei, big.NewInt(config.MaxTipCap))
	gasFeeCap := new(big.Int).Mul(bigGwei, big.NewInt(config.MaxFeeCap))

	txSequences, err := GenerateTxSequences(ctx, clients[0], pks, gasFeeCap, gasTipCap, config.TxsPerWorker)
	if err != nil {
		return nil, err
	}
	if len(clients) < config.Workers {
		return nil, fmt.Errorf("less clients %d than requested workers %d", len(clients), config.Workers)
	}
	if len(senders) < config.Workers {
		return nil, fmt.Errorf("less senders %d than requested workers %d", len(senders), config.Workers)
	}
	if len(txSequences) < config.Workers {
		return nil, fmt.Errorf("less txSequences %d than requested workers %d", len(txSequences), config.Workers)
	}

	wg := NewWorkerGroup(clients[:config.Workers], senders[:config.Workers], txSequences[:config.Workers])
	return wg, nil
}

func ExecuteLoader(ctx context.Context, config config.Config) error {
	ctx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	loader, err := CreateLoader(ctx, config)
	if err != nil {
		return err
	}
	return loader.Execute(ctx)
}
