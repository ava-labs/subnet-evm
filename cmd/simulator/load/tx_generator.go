// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package load

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

type CreateTxData func() *TxData

type TxData struct {
	To    *common.Address
	Data  []byte
	Value *big.Int
	Gas   uint64
}

func GenerateFundDistributionTxSequence(key *ecdsa.PrivateKey, chainID *big.Int, signer types.Signer, startingNonce uint64, gasFeeCap *big.Int, gasTipCap *big.Int, fundAddrs []common.Address, value *big.Int) ([]*types.Transaction, error) {
	// Create a closure to use in GenerateTxSequence
	addrs := make([]common.Address, len(fundAddrs))
	copy(addrs, fundAddrs)
	i := 0
	return GenerateTxSequence(
		key,
		chainID,
		signer,
		startingNonce,
		gasFeeCap,
		gasTipCap,
		func() *TxData {
			data := &TxData{
				To:    &addrs[i],
				Data:  nil,
				Value: new(big.Int).Add(common.Big0, value),
				Gas:   params.TxGas,
			}
			i++
			return data
		},
		uint64(len(addrs)),
	)
}

func GenerateTxSequence(key *ecdsa.PrivateKey, chainID *big.Int, signer types.Signer, startingNonce uint64, gasFeeCap *big.Int, gasTipCap *big.Int, generator CreateTxData, numTxs uint64) ([]*types.Transaction, error) {
	txs := make([]*types.Transaction, 0, numTxs)
	for i := uint64(0); i < numTxs; i++ {
		txData := generator()
		tx, err := types.SignNewTx(key, signer, &types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     startingNonce + i,
			GasTipCap: gasTipCap,
			GasFeeCap: gasFeeCap,
			Gas:       txData.Gas,
			To:        txData.To,
			Data:      txData.Data,
			Value:     txData.Value,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to sign tx at index %d: %w", i, err)
		}
		txs = append(txs, tx)
	}
	return txs, nil
}

func GenerateNoopTxSequence(key *ecdsa.PrivateKey, chainID *big.Int, signer types.Signer, startingNonce uint64, gasFeeCap *big.Int, gasTipCap *big.Int, address common.Address, numTxs uint64) ([]*types.Transaction, error) {
	return GenerateTxSequence(
		key,
		chainID,
		signer,
		startingNonce,
		gasFeeCap,
		gasTipCap,
		func() *TxData {
			return &TxData{
				To:    &address,
				Data:  nil,
				Value: common.Big0,
				Gas:   params.TxGas,
			}
		},
		numTxs,
	)
}

func GenerateTxSequences(ctx context.Context, client ethclient.Client, keys []*ecdsa.PrivateKey, gasFeeCap *big.Int, gasTipCap *big.Int, txsPerKey uint64) ([]types.Transactions, error) {
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chainID: %w", err)
	}
	signer := types.LatestSignerForChainID(chainID)

	txSequences := make([]types.Transactions, len(keys))
	for i, key := range keys {
		address := ethcrypto.PubkeyToAddress(key.PublicKey)
		startingNonce, err := client.NonceAt(ctx, address, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch nonce for address %s at index %d: %w", address, i, err)
		}
		txs, err := GenerateNoopTxSequence(key, chainID, signer, startingNonce, gasFeeCap, gasTipCap, address, txsPerKey)
		if err != nil {
			return nil, fmt.Errorf("failed to generate tx sequence at index %d: %w", i, err)
		}
		txSequences[i] = txs
	}
	return txSequences, nil
}
