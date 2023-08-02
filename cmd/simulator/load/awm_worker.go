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
	"github.com/ethereum/go-ethereum/log"
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
}

func (aw *awmWorker) IssueTx(ctx context.Context, tx *AwmTx) error {
	if aw.onIssued != nil {
		aw.onIssued(tx.AwmID)
	}
	log.Info("awm worker issuing tx", "hash", tx.Tx.Hash())
	return aw.worker.IssueTx(ctx, tx.Tx)
}

func (aw *awmWorker) ConfirmTx(ctx context.Context, tx *AwmTx) error {
	log.Info("awm worker confirming tx", "hash", tx.Tx.Hash())
	if err := aw.worker.ConfirmTx(ctx, tx.Tx); err != nil {
		return err
	}
	if aw.onConfirmed != nil {
		aw.onConfirmed(tx.AwmID)
	}
	return nil
}

func (aw *awmWorker) Close(ctx context.Context) error {
	return aw.worker.Close(ctx)
}

type timeTracker struct {
	lock     sync.Mutex
	issued   map[common.Hash]time.Time
	observer func(float64)
}

func newTimeTracker(observer func(float64)) *timeTracker {
	return &timeTracker{
		issued:   make(map[common.Hash]time.Time),
		observer: observer,
	}
}

func (tt *timeTracker) IssueTx(id common.Hash) {
	tt.lock.Lock()
	defer tt.lock.Unlock()

	tt.issued[id] = time.Now()
	log.Info("awm time tracker issued tx", "id", id)
}

func (tt *timeTracker) ConfirmTx(id common.Hash) {
	tt.lock.Lock()
	defer tt.lock.Unlock()

	start, ok := tt.issued[id]
	if !ok {
		panic("unexpected confirm " + id.Hex())
	}
	delete(tt.issued, id)
	duration := time.Since(start)
	tt.observer(duration.Seconds())
	log.Info("awm time tracker confirmed tx", "id", id, "duration", duration)
}
