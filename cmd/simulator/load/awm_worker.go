// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package load

import (
	"context"
	"sync"
	"time"

	"github.com/ava-labs/subnet-evm/cmd/simulator/txs"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ethereum/go-ethereum/common"
)

var _ txs.Worker[*AwmTx] = &awmWorker{}

type AwmTx struct {
	AwmID common.Hash // a hash of the payload used to uniquely identify this awm message
	Tx    *types.Transaction
}

func (a *AwmTx) Hash() common.Hash {
	return a.Tx.Hash() // this is used by the batch worker to track tx times
}

type awmWorker struct {
	worker      *singleAddressTxWorker
	onIssued    func(common.Hash)
	onConfirmed func(common.Hash)
	onClosed    func()
}

func (aw *awmWorker) IssueTx(ctx context.Context, tx *AwmTx) error {
	if aw.onIssued != nil {
		aw.onIssued(tx.AwmID)
	}
	return aw.worker.IssueTx(ctx, tx.Tx)
}

func (aw *awmWorker) ConfirmTx(ctx context.Context, tx *AwmTx) error {
	if err := aw.worker.ConfirmTx(ctx, tx.Tx); err != nil {
		return err
	}
	if aw.onConfirmed != nil {
		aw.onConfirmed(tx.AwmID)
	}
	return nil
}

func (aw *awmWorker) Close(ctx context.Context) error {
	if aw.onClosed != nil {
		aw.onClosed()
	}
	return aw.worker.Close(ctx)
}

type txTracker struct {
	lock   sync.Mutex
	closed bool
	done   chan struct{}

	issued   map[common.Hash]time.Time
	observer func(float64)
}

func newTxTracker(observer func(float64)) *txTracker {
	return &txTracker{
		issued:   make(map[common.Hash]time.Time),
		observer: observer,
		done:     make(chan struct{}),
	}
}

func (tt *txTracker) IssueTx(id common.Hash) {
	tt.lock.Lock()
	defer tt.lock.Unlock()

	tt.issued[id] = time.Now()
}

func (tt *txTracker) ConfirmTx(id common.Hash) {
	tt.lock.Lock()
	defer tt.lock.Unlock()

	start, ok := tt.issued[id]
	if !ok {
		panic("unexpected confirm " + id.Hex())
	}
	duration := time.Since(start)
	tt.observer(duration.Seconds())

	delete(tt.issued, id)
	tt.checkDone()
}

func (tt *txTracker) Close() {
	tt.lock.Lock()
	defer tt.lock.Unlock()

	tt.closed = true
	tt.checkDone()
}

// assumes lock is held
func (tt *txTracker) checkDone() {
	if !tt.closed || len(tt.issued) > 0 {
		return
	}
	close(tt.done)
}
