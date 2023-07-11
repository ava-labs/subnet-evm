// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

var _ TxSequence[TimedTx] = (*txSequence)(nil)

type CreateTx func(key *ecdsa.PrivateKey, nonce uint64) (*types.Transaction, error)

type TimedTx struct {
	Tx                             *types.Transaction
	IssuanceDuration               time.Duration
	ConfirmationDuration           time.Duration
	IssuanceStart                  time.Time
	IssuanceToConfirmationDuration time.Duration
}

// GenerateTxSequence fetches the current nonce of key and calls [generator] [numTxs] times sequentially to generate a sequence of transactions.
func GenerateTxSequence(ctx context.Context, generator CreateTx, client ethclient.Client, key *ecdsa.PrivateKey, numTxs uint64) (TxSequence[TimedTx], error) {
	address := ethcrypto.PubkeyToAddress(key.PublicKey)
	startingNonce, err := client.NonceAt(ctx, address, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch nonce for address %s: %w", address, err)
	}
	txs := make([]TimedTx, 0, numTxs)
	for i := uint64(0); i < numTxs; i++ {
		tx, err := generator(key, startingNonce+i)
		if err != nil {
			return nil, fmt.Errorf("failed to sign tx at index %d: %w", i, err)
		}
		timedTx := TimedTx{Tx: tx}
		txs = append(txs, timedTx)
	}
	return ConvertTxSliceToSequence(txs), nil
}

func GenerateTxSequences(ctx context.Context, generator CreateTx, client ethclient.Client, keys []*ecdsa.PrivateKey, txsPerKey uint64) ([]TxSequence[TimedTx], error) {
	txSequences := make([]TxSequence[TimedTx], len(keys))
	for i, key := range keys {
		txs, err := GenerateTxSequence(ctx, generator, client, key, txsPerKey)
		if err != nil {
			return nil, fmt.Errorf("failed to generate tx sequence at index %d: %w", i, err)
		}
		txSequences[i] = txs
	}
	return txSequences, nil
}

type txSequence struct {
	txChan chan TimedTx
}

func ConvertTxSliceToSequence(txs []TimedTx) TxSequence[TimedTx] {
	txChan := make(chan TimedTx, len(txs))
	for _, tx := range txs {
		txChan <- tx
	}
	close(txChan)

	return &txSequence{
		txChan: txChan,
	}
}

func (t *txSequence) Chan() <-chan TimedTx {
	return t.txChan
}
