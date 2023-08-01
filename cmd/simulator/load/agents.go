// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package load

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ava-labs/subnet-evm/cmd/simulator/config"
	"github.com/ava-labs/subnet-evm/cmd/simulator/metrics"
	"github.com/ava-labs/subnet-evm/cmd/simulator/txs"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

type mkAgentBuilder func(
	ctx context.Context, config config.Config, chainID *big.Int,
	pks []*ecdsa.PrivateKey, client ethclient.Client, metrics *metrics.Metrics,
) (AgentBuilder, error)

type AgentBuilder interface {
	NewAgent(ctx context.Context, idx int, client ethclient.Client, sender common.Address) (txs.Agent, error)
}

type transferTxAgentBuilder struct {
	txSequences []txs.TxSequence[*types.Transaction]
	batchSize   uint64
	metrics     *metrics.Metrics
}

func NewTransferTxAgentBuilder(
	ctx context.Context, config config.Config, chainID *big.Int,
	pks []*ecdsa.PrivateKey, client ethclient.Client, metrics *metrics.Metrics,
) (AgentBuilder, error) {
	log.Info("Creating transaction sequences...")
	bigGwei := big.NewInt(params.GWei)
	gasTipCap := new(big.Int).Mul(bigGwei, big.NewInt(config.MaxTipCap))
	gasFeeCap := new(big.Int).Mul(bigGwei, big.NewInt(config.MaxFeeCap))
	signer := types.LatestSignerForChainID(chainID)

	txGenerator := func(key *ecdsa.PrivateKey, nonce uint64) (*types.Transaction, error) {
		addr := ethcrypto.PubkeyToAddress(key.PublicKey)
		tx, err := types.SignNewTx(key, signer, &types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			GasTipCap: gasTipCap,
			GasFeeCap: gasFeeCap,
			Gas:       params.TxGas,
			To:        &addr,
			Data:      nil,
			Value:     common.Big0,
		})
		if err != nil {
			return nil, err
		}
		return tx, nil
	}
	txSequences, err := txs.GenerateTxSequences(ctx, txGenerator, client, pks, config.TxsPerWorker)
	if err != nil {
		return nil, err
	}

	return &transferTxAgentBuilder{
		txSequences: txSequences,
		batchSize:   config.BatchSize,
		metrics:     metrics,
	}, nil
}

func (t *transferTxAgentBuilder) NewAgent(
	ctx context.Context, idx int, client ethclient.Client, sender common.Address,
) (txs.Agent, error) {
	worker := NewSingleAddressTxWorker(ctx, client, sender)
	return txs.NewIssueNAgent[*types.Transaction](
		t.txSequences[idx], worker, t.batchSize, t.metrics), nil
}

type warpSendTxAgentBuilder struct {
	txSequences []txs.TxSequence[*AwmTx]
	batchSize   uint64
	metrics     *metrics.Metrics
	timeTracker *timeTracker
}

func NewWarpSendTxAgentBuilder(
	ctx context.Context, config config.Config, chainID *big.Int,
	pks []*ecdsa.PrivateKey, client ethclient.Client, metrics *metrics.Metrics,
	timeTracker *timeTracker,
) (AgentBuilder, error) {
	log.Info("Creating warp transaction sequences...")
	txSequences, err := GetWarpSendTxSequences(ctx, config, chainID, pks, client)
	if err != nil {
		return nil, err
	}
	return &warpSendTxAgentBuilder{
		txSequences: txSequences,
		batchSize:   config.BatchSize,
		metrics:     metrics,
		timeTracker: timeTracker,
	}, nil
}

func (w *warpSendTxAgentBuilder) NewAgent(
	ctx context.Context, idx int, client ethclient.Client, sender common.Address,
) (txs.Agent, error) {
	worker := NewSingleAddressTxWorker(ctx, client, sender)
	awmWorker := &awmWorker{
		worker:   worker,
		onIssued: w.timeTracker.IssueTx,
	}
	return txs.NewIssueNAgent[*AwmTx](
		w.txSequences[idx], awmWorker, w.batchSize, w.metrics), nil
}

type warpReceiveTxAgentBuilder struct {
	txSequences []txs.TxSequence[*AwmTx]
	batchSize   uint64
	metrics     *metrics.Metrics
	timeTracker *timeTracker
}

func NewWarpReceiveTxAgentBuilder(
	ctx context.Context, config config.Config, chainID *big.Int,
	pks []*ecdsa.PrivateKey, client ethclient.Client, metrics *metrics.Metrics,
	timeTracker *timeTracker,
) (AgentBuilder, error) {
	log.Info("Creating warp transaction sequences...")
	txSequences, err := GetWarpReceiveTxSequences(ctx, config, chainID, pks, client)
	if err != nil {
		return nil, err
	}
	return &warpReceiveTxAgentBuilder{
		txSequences: txSequences,
		batchSize:   config.BatchSize,
		metrics:     metrics,
		timeTracker: timeTracker,
	}, nil
}

func (w *warpReceiveTxAgentBuilder) NewAgent(
	ctx context.Context, idx int, client ethclient.Client, sender common.Address,
) (txs.Agent, error) {
	worker := NewSingleAddressTxWorker(ctx, client, sender)
	awmWorker := &awmWorker{
		worker:      worker,
		onConfirmed: w.timeTracker.ConfirmTx,
	}
	return txs.NewIssueNAgent[*AwmTx](
		w.txSequences[idx], awmWorker, w.batchSize, w.metrics), nil
}
