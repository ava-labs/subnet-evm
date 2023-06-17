// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"context"
	"fmt"
)

type TxSequence[T any] interface {
	Chan() <-chan T
}

type Worker[T any] interface {
	IssueTx(ctx context.Context, tx T) error
	ConfirmTx(ctx context.Context, tx T) error
	Close(ctx context.Context) error
}

type Agent[T any] interface {
	Execute(ctx context.Context) error
}

type issueAllAgent[T any] struct {
	sequence TxSequence[T]
	worker   Worker[T]
}

func NewIssueAllAgent[T any](sequence TxSequence[T], worker Worker[T]) Agent[T] {
	return &issueAllAgent[T]{
		sequence: sequence,
		worker:   worker,
	}
}

func (a issueAllAgent[T]) Execute(ctx context.Context) error {
	txChan := a.sequence.Chan()

	txs := make([]T, 0, 100)
	for tx := range txChan {
		if err := a.worker.IssueTx(ctx, tx); err != nil {
			return fmt.Errorf("failed to issue transaction %d: %w", len(txs), err)
		}
		txs = append(txs, tx)
	}

	for i, tx := range txs {
		if err := a.worker.ConfirmTx(ctx, tx); err != nil {
			return fmt.Errorf("failed to await transaction %d: %w", i, err)
		}
	}

	return a.worker.Close(ctx)
}

type issueNAgent[T any] struct {
	sequence TxSequence[T]
	worker   Worker[T]
	n        int
}

func NewIssueNAgent[T any](sequence TxSequence[T], worker Worker[T], n int) Agent[T] {
	return &issueNAgent[T]{
		sequence: sequence,
		worker:   worker,
		n:        n,
	}
}

func (a issueNAgent[T]) Execute(ctx context.Context) error {
	txChan := a.sequence.Chan()

	for {
		var (
			txs  = make([]T, 0, a.n)
			tx   T
			done bool
		)
		for i := 0; i < a.n; i++ {
			select {
			case tx, done = <-txChan:
				if done {
					return a.worker.Close(ctx)
				}
			case <-ctx.Done():
				return ctx.Err()
			}

			if err := a.worker.IssueTx(ctx, tx); err != nil {
				return fmt.Errorf("failed to issue transaction %d: %w", len(txs), err)
			}
			txs = append(txs, tx)
		}

		for i, tx := range txs {
			if err := a.worker.ConfirmTx(ctx, tx); err != nil {
				return fmt.Errorf("failed to await transaction %d: %w", i, err)
			}
		}
	}
}
