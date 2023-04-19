// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package load

import (
	"context"
	"time"

	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/sync/errgroup"
)

type WorkerGroup struct {
	Workers []*Worker
}

func NewWorkerGroup(clients []ethclient.Client, senders []common.Address, txSequences []types.Transactions) *WorkerGroup {
	workers := make([]*Worker, len(clients))
	for i, client := range clients {
		workers[i] = NewWorker(client, senders[i], txSequences[i])
	}

	return &WorkerGroup{
		Workers: workers,
	}
}

func (wg *WorkerGroup) executeTask(ctx context.Context, f func(ctx context.Context, w *Worker) error) error {
	eg := errgroup.Group{}
	for _, worker := range wg.Workers {
		worker := worker
		eg.Go(func() error {
			return f(ctx, worker)
		})
	}

	return eg.Wait()
}

func (wg *WorkerGroup) IssueTxs(ctx context.Context) error {
	return wg.executeTask(ctx, func(ctx context.Context, w *Worker) error {
		return w.ExecuteTxsFromAddress(ctx)
	})
}

func (wg *WorkerGroup) AwaitTxs(ctx context.Context) error {
	return wg.executeTask(ctx, func(ctx context.Context, w *Worker) error {
		return w.AwaitTxs(ctx)
	})
}

func (wg *WorkerGroup) ConfirmAllTransactions(ctx context.Context) error {
	return wg.executeTask(ctx, func(ctx context.Context, w *Worker) error {
		return w.ConfirmAllTransactions(ctx)
	})
}

func (wg *WorkerGroup) Execute(ctx context.Context) error {
	start := time.Now()
	defer func() {
		log.Info("Completed execution", "totalTime", time.Since(start))
	}()

	executionStart := time.Now()
	log.Info("Executing transactions", "startTime", executionStart)
	if err := wg.IssueTxs(ctx); err != nil {
		return err
	}
	awaitStart := time.Now()
	log.Info("Awaiting transactions", "startTime", awaitStart, "executionElapsed", awaitStart.Sub(executionStart))
	if err := wg.AwaitTxs(ctx); err != nil {
		return err
	}

	confirmationStart := time.Now()
	log.Info("Confirming transactions", "startTime", confirmationStart, "awaitElapsed", confirmationStart.Sub(awaitStart))
	if err := wg.ConfirmAllTransactions(ctx); err != nil {
		return err
	}
	log.Info("Transaction confirmation completed", "confirmationElapsed", time.Since(confirmationStart))

	return nil
}
