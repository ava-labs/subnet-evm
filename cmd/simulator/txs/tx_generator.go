// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/ava-labs/subnet-evm/core/types"
)

var _ TxSequence[*types.Transaction] = (*txSequence[*types.Transaction])(nil)

type CreateTx[T any] func(key *ecdsa.PrivateKey, nonce uint64) (T, error)

// GenerateTxSequence fetches the current nonce of key and calls [generator]
// [numTxs] times sequentially to generate a sequence of transactions.
func GenerateTxSequence[T any](
	ctx context.Context, generator CreateTx[T],
	key *ecdsa.PrivateKey, startingNonce uint64, numTxs uint64,
) (TxSequence[T], error) {
	txs := make([]T, 0, numTxs)
	for i := uint64(0); i < numTxs; i++ {
		tx, err := generator(key, startingNonce+i)
		if err != nil {
			return nil, fmt.Errorf("failed to sign tx at index %d: %w", i, err)
		}
		txs = append(txs, tx)
	}
	return ConvertTxSliceToSequence(txs), nil
}

func GenerateTxSequences[T any](
	ctx context.Context, generator CreateTx[T],
	keys []*ecdsa.PrivateKey, startingNonces []uint64, txsPerKey uint64,
) ([]TxSequence[T], error) {
	txSequences := make([]TxSequence[T], len(keys))
	for i, key := range keys {
		txs, err := GenerateTxSequence(ctx, generator, key, startingNonces[i], txsPerKey)
		if err != nil {
			return nil, fmt.Errorf("failed to generate tx sequence at index %d: %w", i, err)
		}
		txSequences[i] = txs
	}
	return txSequences, nil
}

type txSequence[T any] struct {
	txChan chan T
}

func ConvertTxSliceToSequence[T any](txs []T) TxSequence[T] {
	txChan := make(chan T, len(txs))
	for _, tx := range txs {
		txChan <- tx
	}
	close(txChan)

	return &txSequence[T]{
		txChan: txChan,
	}
}

func (t *txSequence[T]) Chan() <-chan T {
	return t.txChan
}
