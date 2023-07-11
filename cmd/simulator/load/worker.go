// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package load

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ava-labs/subnet-evm/cmd/simulator/txs"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/interfaces"
	"github.com/aybabtme/uniplot/histogram"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type singleAddressTxWorker struct {
	client ethclient.Client

	acceptedNonce uint64
	address       common.Address

	sub      interfaces.Subscription
	newHeads chan *types.Header

	issuanceToConfirmationHistogram []float64
	confirmationHistogram           []float64
	issuanceHistogram               []float64
}

// NewSingleAddressTxWorker creates and returns a singleAddressTxWorker
func NewSingleAddressTxWorker(ctx context.Context, client ethclient.Client, address common.Address) *singleAddressTxWorker {
	newHeads := make(chan *types.Header)
	tw := &singleAddressTxWorker{
		client:   client,
		address:  address,
		newHeads: newHeads,
	}

	sub, err := client.SubscribeNewHead(ctx, newHeads)
	if err != nil {
		log.Debug("failed to subscribe new heads, falling back to polling", "err", err)
	} else {
		tw.sub = sub
	}

	return tw
}

func (tw *singleAddressTxWorker) IssueTx(ctx context.Context, timedTx txs.TimedTx) error {
	start := time.Now()
	err := tw.client.SendTransaction(ctx, timedTx.Tx)
	if err != nil {
		return err
	}
	issuanceDuration := time.Since(start)

	timedTx.IssuanceStart = start
	timedTx.IssuanceDuration = issuanceDuration
	tw.issuanceHistogram = append(tw.issuanceHistogram, timedTx.IssuanceDuration.Seconds())
	return nil
}

func (tw *singleAddressTxWorker) ConfirmTx(ctx context.Context, timedTx txs.TimedTx) error {
	tx := timedTx.Tx
	txNonce := tx.Nonce()

	start := time.Now()
	for {
		// If the is less than what has already been accepted, the transaction is confirmed
		if txNonce < tw.acceptedNonce {
			confirmationEnd := time.Now()
			timedTx.ConfirmationDuration = confirmationEnd.Sub(start)
			timedTx.IssuanceToConfirmationDuration = confirmationEnd.Sub(timedTx.IssuanceStart)
			tw.issuanceToConfirmationHistogram = append(tw.issuanceToConfirmationHistogram, timedTx.IssuanceToConfirmationDuration.Seconds())
			tw.confirmationHistogram = append(tw.confirmationHistogram, timedTx.ConfirmationDuration.Seconds())

			return nil
		}

		select {
		case <-tw.newHeads:
		case <-time.After(time.Second):
		case <-ctx.Done():
			return fmt.Errorf("failed to await tx %s nonce %d: %w", tx.Hash(), txNonce, ctx.Err())
		}

		// Update the worker's accepted nonce, so we can check on the next iteration
		// if the transaction has been accepted.
		acceptedNonce, err := tw.client.NonceAt(ctx, tw.address, nil)
		if err != nil {
			return fmt.Errorf("failed to await tx %s nonce %d: %w", tx.Hash(), txNonce, err)
		}
		tw.acceptedNonce = acceptedNonce
	}
}

func (tw *singleAddressTxWorker) Close(ctx context.Context) error {
	if tw.sub != nil {
		tw.sub.Unsubscribe()
	}
	close(tw.newHeads)
	return nil
}

func (tw *singleAddressTxWorker) CollectMetrics(ctx context.Context) error {
	log.Info("Individual Tx Issuance to Confirmation Duration Histogram (s)")
	hist := histogram.Hist(10, tw.issuanceToConfirmationHistogram)
	err := histogram.Fprint(os.Stdout, hist, histogram.Linear(5))
	if err != nil {
		return err
	}

	log.Info("Individual Tx Issuance Histogram (s)")
	hist = histogram.Hist(10, tw.issuanceHistogram)
	err = histogram.Fprint(os.Stdout, hist, histogram.Linear(5))
	if err != nil {
		return err
	}

	log.Info("Individual Tx Confirmation (s)")
	hist = histogram.Hist(10, tw.confirmationHistogram)
	err = histogram.Fprint(os.Stdout, hist, histogram.Linear(5))
	if err != nil {
		return err
	}
	return nil
}
