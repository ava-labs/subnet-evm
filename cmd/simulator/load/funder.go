// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package load

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/cmd/simulator/key"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

// DistributeFunds ensures that each address in keys has at least [minFundsPerAddr] by sending funds
// from the key with the highest starting balance.
// This function should never return a set of keys with length less than [numKeys]
func DistributeFunds(ctx context.Context, client ethclient.Client, keys []*key.Key, numKeys int, minFundsPerAddr *big.Int) ([]*key.Key, error) {
	if len(keys) < numKeys {
		return nil, fmt.Errorf("insufficient number of keys %d < %d", len(keys), numKeys)
	}
	fundedKeys := make([]*key.Key, 0, numKeys)
	// TODO: clean up fund distribution.
	needFundsKeys := make([]*key.Key, 0)
	needFundsAddrs := make([]common.Address, 0)

	maxFundsKey := keys[0]
	maxFundsBalance := common.Big0
	log.Info("Checking balance of each key to distribute funds")
	for _, key := range keys {
		balance, err := client.BalanceAt(ctx, key.Address, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch balance for addr %s: %w", key.Address, err)
		}

		if balance.Cmp(minFundsPerAddr) < 0 {
			needFundsKeys = append(needFundsKeys, key)
			needFundsAddrs = append(needFundsAddrs, key.Address)
		} else {
			fundedKeys = append(fundedKeys, key)
		}

		if balance.Cmp(maxFundsBalance) > 0 {
			maxFundsKey = key
			maxFundsBalance = balance
		}
	}
	requiredFunds := new(big.Int).Mul(minFundsPerAddr, big.NewInt(int64(numKeys)))
	if maxFundsBalance.Cmp(requiredFunds) < 0 {
		return nil, fmt.Errorf("insufficient funds to distribute %d < %d", maxFundsBalance, requiredFunds)
	}
	log.Info("Found max funded key", "address", maxFundsKey.Address, "balance", maxFundsBalance, "numFundAddrs", len(needFundsAddrs))
	if len(fundedKeys) >= numKeys {
		return fundedKeys[:numKeys], nil
	}

	// If there are not enough funded keys, cut [needFundsAddrs] to the number of keys that
	// must be funded to reach [numKeys] required.
	fundKeysCutLen := numKeys - len(fundedKeys)
	needFundsKeys = needFundsKeys[:fundKeysCutLen]
	needFundsAddrs = needFundsAddrs[:fundKeysCutLen]

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch chainID: %w", err)
	}
	nonce, err := client.NonceAt(ctx, maxFundsKey.Address, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch nonce of address %s: %w", maxFundsKey.Address, err)
	}
	gasFeeCap, err := client.EstimateBaseFee(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch estimated base fee: %w", err)
	}
	gasTipCap, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch suggested gas tip: %w", err)
	}
	signer := types.LatestSignerForChainID(chainID)

	// Generate a sequence of transactions to distribute the required funds.
	log.Info("Generating distribution transactions")
	txs, err := GenerateFundDistributionTxSequence(maxFundsKey.PrivKey, chainID, signer, nonce, gasFeeCap, gasTipCap, needFundsAddrs, minFundsPerAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to generate fund distribution sequence from %s of length %d", maxFundsKey.Address, len(needFundsAddrs))
	}

	log.Info("Executing distribution transactions...")
	worker := NewWorker(client, maxFundsKey.Address, txs)
	if err := worker.Execute(ctx); err != nil {
		return nil, err
	}

	for _, addr := range needFundsAddrs {
		balance, err := client.BalanceAt(ctx, addr, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch balance for addr %s: %w", addr, err)
		}
		log.Info("Funded address has balance", "balance", balance)
	}
	return needFundsKeys, nil
}
