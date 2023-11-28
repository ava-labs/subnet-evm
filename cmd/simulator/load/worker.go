// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package load

import (
	"context"
	"fmt"
	"time"

	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/interfaces"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type ethereumTxWorker struct {
	client ethclient.Client

	acceptedNonce uint64
	address       common.Address

	sub      interfaces.Subscription
	newHeads chan *types.Header
}

// NewSingleAddressTxWorker creates and returns a new ethereumTxWorker that confirms transactions by checking the latest
// nonce of [address] and assuming any transaction with a lower nonce was already accepted.
func NewSingleAddressTxWorker(ctx context.Context, client ethclient.Client, address common.Address) *ethereumTxWorker {
	newHeads := make(chan *types.Header)
	tw := &ethereumTxWorker{
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

// NewTxReceiptWorker creates and returns a new ethereumTxWorker that confirms transactions by checking for the
// corresponding transaction receipt.
func NewTxReceiptWorker(ctx context.Context, client ethclient.Client) *ethereumTxWorker {
	newHeads := make(chan *types.Header)
	tw := &ethereumTxWorker{
		client:   client,
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

func (tw *ethereumTxWorker) IssueTx(ctx context.Context, tx *types.Transaction) error {
	return tw.client.SendTransaction(ctx, tx)
}

func (tw *ethereumTxWorker) ConfirmTx(ctx context.Context, tx *types.Transaction) error {
	if tw.address == (common.Address{}) {
		return tw.confirmTxByReceipt(ctx, tx)
	}
	return tw.confirmTxByNonce(ctx, tx)
}

func (tw *ethereumTxWorker) confirmTxByNonce(ctx context.Context, tx *types.Transaction) error {
	txNonce := tx.Nonce()

	for {
		// Update the worker's accepted nonce, so we can check on the next iteration
		// if the transaction has been accepted.
		acceptedNonce, err := tw.client.NonceAt(ctx, tw.address, nil)
		if err != nil {
			return fmt.Errorf("failed to await tx %s nonce %d: %w", tx.Hash(), txNonce, err)
		}
		tw.acceptedNonce = acceptedNonce

		log.Info("trying to confirm tx", "txHash", tx.Hash(), "txNonce", txNonce, "acceptedNonce", tw.acceptedNonce)
		// If the is less than what has already been accepted, the transaction is confirmed
		if txNonce < tw.acceptedNonce {
			return nil
		}

		select {
		case <-tw.newHeads:
		case <-time.After(time.Second):
		case <-ctx.Done():
			return fmt.Errorf("failed to await tx %s nonce %d: %w", tx.Hash(), txNonce, ctx.Err())
		}
	}
}

func (tw *ethereumTxWorker) confirmTxByReceipt(ctx context.Context, tx *types.Transaction) error {
	for {
		_, err := tw.client.TransactionReceipt(ctx, tx.Hash())
		if err == nil {
			return nil
		}
		log.Debug("no tx receipt", "txHash", tx.Hash(), "nonce", tx.Nonce(), "err", err)

		select {
		case <-tw.newHeads:
		case <-time.After(time.Second):
		case <-ctx.Done():
			return fmt.Errorf("failed to await tx %s nonce %d: %w", tx.Hash(), tx.Nonce(), ctx.Err())
		}
	}
}

func (tw *ethereumTxWorker) LatestHeight(ctx context.Context) (uint64, error) {
	return tw.client.BlockNumber(ctx)
}

func (tw *ethereumTxWorker) Close(ctx context.Context) error {
	if tw.sub != nil {
		tw.sub.Unsubscribe()
	}
	close(tw.newHeads)
	return nil
}
