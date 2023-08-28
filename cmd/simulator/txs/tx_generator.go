// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/ava-labs/subnet-evm/core/types"
)

var _ TxSequence[*types.Transaction] = (*txSequence)(nil)

type CreateTx func(key *ecdsa.PrivateKey, nonce uint64) (*types.Transaction, error)

// GenerateTxSequence fetches the current nonce of key and calls [generator]
// [numTxs] times sequentially to generate a sequence of transactions.
func GenerateTxSequence(
	ctx context.Context,
	generator CreateTx,
	key *ecdsa.PrivateKey,
	startingNonce uint64,
	numTxs uint64,
) (TxSequence[*types.Transaction], error) {
	txs := make([]*types.Transaction, 0, numTxs)
	for i := uint64(0); i < numTxs; i++ {
		tx, err := generator(key, startingNonce+i)
		if err != nil {
			return nil, fmt.Errorf("failed to sign tx at index %d: %w", i, err)
		}
		txs = append(txs, tx)
	}
	return ConvertTxSliceToSequence(txs), nil
}

func GenerateTxSequences(
	ctx context.Context,
	generator CreateTx,
	keys []*ecdsa.PrivateKey,
	startingNonces []uint64,
	txsPerKey uint64,
) ([]TxSequence[*types.Transaction], error) {
	txSequences := make([]TxSequence[*types.Transaction], len(keys))
	for i, key := range keys {
		txs, err := GenerateTxSequence(ctx, generator, key, startingNonces[i], txsPerKey)
		if err != nil {
			return nil, fmt.Errorf("failed to generate tx sequence at index %d: %w", i, err)
		}
		txSequences[i] = txs
	}
	return txSequences, nil
}

type txSequence struct {
	txChan chan *types.Transaction
}

func ConvertTxSliceToSequence(txs []*types.Transaction) TxSequence[*types.Transaction] {
	txChan := make(chan *types.Transaction, len(txs))
	for _, tx := range txs {
		txChan <- tx
	}
	close(txChan)

	return &txSequence{
		txChan: txChan,
	}
}

func (t *txSequence) Chan() <-chan *types.Transaction {
	return t.txChan
}
