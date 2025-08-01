// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// NOTE: this piece of code is modified by Ted Yin.
// The modification is also licensed under the same LGPL.

package miner

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ava-labs/avalanchego/utils/timer/mockable"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/types"
	"github.com/ava-labs/libevm/core/vm"
	"github.com/ava-labs/libevm/event"
	"github.com/ava-labs/libevm/log"
	ethparams "github.com/ava-labs/libevm/params"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/consensus"
	"github.com/ava-labs/subnet-evm/consensus/misc/eip4844"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/core/txpool"
	"github.com/ava-labs/subnet-evm/params"
	customheader "github.com/ava-labs/subnet-evm/plugin/evm/header"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/predicate"
	"github.com/holiman/uint256"
)

const (
	targetTxsSize = 1800 * units.KiB
)

// environment is the worker's current environment and holds all of the current state information.
type environment struct {
	signer  types.Signer
	state   *state.StateDB // apply state changes here
	tcount  int            // tx count in cycle
	gasPool *core.GasPool  // available gas used to pack transactions

	parent   *types.Header
	header   *types.Header
	txs      []*types.Transaction
	receipts []*types.Receipt
	sidecars []*types.BlobTxSidecar
	blobs    int
	size     uint64

	rules            params.Rules
	predicateContext *precompileconfig.PredicateContext
	// predicateResults contains the results of checking the predicates for each transaction in the miner.
	// The results are accumulated as transactions are executed by the miner and set on the BlockContext.
	// If a transaction is dropped, its results must explicitly be removed from predicateResults in the same
	// way that the gas pool and state is reset.
	predicateResults *predicate.Results

	start time.Time // Time that block building began
}

// worker is the main object which takes care of submitting new work to consensus engine
// and gathering the sealing result.
type worker struct {
	config      *Config
	chainConfig *params.ChainConfig
	engine      consensus.Engine
	eth         Backend
	chain       *core.BlockChain

	// Feeds
	// TODO remove since this will never be written to
	pendingLogsFeed event.Feed

	// Subscriptions
	mux        *event.TypeMux // TODO replace
	mu         sync.RWMutex   // The lock used to protect the coinbase and extra fields
	coinbase   common.Address
	clock      *mockable.Clock // Allows us mock the clock for testing
	beaconRoot *common.Hash    // TODO: set to empty hash, retained for upstream compatibility and future use
}

func newWorker(config *Config, chainConfig *params.ChainConfig, engine consensus.Engine, eth Backend, mux *event.TypeMux, clock *mockable.Clock) *worker {
	worker := &worker{
		config:      config,
		chainConfig: chainConfig,
		engine:      engine,
		eth:         eth,
		chain:       eth.BlockChain(),
		mux:         mux,
		coinbase:    config.Etherbase,
		clock:       clock,
		beaconRoot:  &common.Hash{},
	}

	return worker
}

// setEtherbase sets the etherbase used to initialize the block coinbase field.
func (w *worker) setEtherbase(addr common.Address) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.coinbase = addr
}

// commitNewWork generates several new sealing tasks based on the parent block.
func (w *worker) commitNewWork(predicateContext *precompileconfig.PredicateContext) (*types.Block, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	tstart := w.clock.Time()
	timestamp := uint64(tstart.Unix())
	parent := w.chain.CurrentBlock()
	// Note: in order to support asynchronous block production, blocks are allowed to have
	// the same timestamp as their parent. This allows more than one block to be produced
	// per second.
	if parent.Time >= timestamp {
		timestamp = parent.Time
	}

	// The fee manager relies on the state of the parent block to set the fee config
	// because the fee config may be changed by the current block.
	feeConfig, _, err := w.chain.GetFeeConfigAt(parent)
	if err != nil {
		return nil, err
	}
	chainConfig := params.GetExtra(w.chainConfig)
	gasLimit, err := customheader.GasLimit(chainConfig, feeConfig, parent, timestamp)
	if err != nil {
		return nil, fmt.Errorf("calculating new gas limit: %w", err)
	}
	baseFee, err := customheader.BaseFee(chainConfig, feeConfig, parent, timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate new base fee: %w", err)
	}

	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     new(big.Int).Add(parent.Number, common.Big1),
		GasLimit:   gasLimit,
		Time:       timestamp,
		BaseFee:    baseFee,
	}

	// Apply EIP-4844, EIP-4788.
	if w.chainConfig.IsCancun(header.Number, header.Time) {
		var excessBlobGas uint64
		if w.chainConfig.IsCancun(parent.Number, parent.Time) {
			excessBlobGas = eip4844.CalcExcessBlobGas(*parent.ExcessBlobGas, *parent.BlobGasUsed)
		} else {
			// For the first post-fork block, both parent.data_gas_used and parent.excess_data_gas are evaluated as 0
			excessBlobGas = eip4844.CalcExcessBlobGas(0, 0)
		}
		header.BlobGasUsed = new(uint64)
		header.ExcessBlobGas = &excessBlobGas
		header.ParentBeaconRoot = w.beaconRoot
	}

	if w.coinbase == (common.Address{}) {
		return nil, errors.New("cannot mine without etherbase")
	}
	header.Coinbase = w.coinbase

	configuredCoinbase, isAllowFeeRecipient, err := w.chain.GetCoinbaseAt(parent)
	if err != nil {
		return nil, fmt.Errorf("failed to get configured coinbase: %w", err)
	}

	// if fee recipients are not allowed, then the coinbase is the configured coinbase
	// don't set w.coinbase directly to the configured coinbase because that would override the
	// coinbase set by the user
	if !isAllowFeeRecipient && w.coinbase != configuredCoinbase {
		log.Info("fee recipients are not allowed, using required coinbase for the mining", "currentminer", w.coinbase, "required", configuredCoinbase)
		header.Coinbase = configuredCoinbase
	}

	if err := w.engine.Prepare(w.chain, header); err != nil {
		return nil, fmt.Errorf("failed to prepare header for mining: %w", err)
	}

	env, err := w.createCurrentEnvironment(predicateContext, parent, header, feeConfig, tstart)
	if err != nil {
		return nil, fmt.Errorf("failed to create new current environment: %w", err)
	}
	if header.ParentBeaconRoot != nil {
		context := core.NewEVMBlockContext(header, w.chain, nil)
		vmenv := vm.NewEVM(context, vm.TxContext{}, env.state, w.chainConfig, vm.Config{})
		core.ProcessBeaconBlockRoot(*header.ParentBeaconRoot, vmenv, env.state)
	}
	// Ensure we always stop prefetcher after block building is complete.
	defer func() {
		if env.state == nil {
			return
		}
		env.state.StopPrefetcher()
	}()
	// Configure any upgrades that should go into effect during this block.
	blockContext := core.NewBlockContext(header.Number, header.Time)
	err = core.ApplyUpgrades(w.chainConfig, &parent.Time, blockContext, env.state)
	if err != nil {
		log.Error("failed to configure precompiles mining new block", "parent", parent.Hash(), "number", header.Number, "timestamp", header.Time, "err", err)
		return nil, err
	}

	// Retrieve the pending transactions pre-filtered by the 1559/4844 dynamic fees
	filter := txpool.PendingFilter{
		MinTip: uint256.MustFromBig(w.eth.TxPool().GasTip()),
	}
	if env.header.BaseFee != nil {
		filter.BaseFee = uint256.MustFromBig(env.header.BaseFee)
	}
	if env.header.ExcessBlobGas != nil {
		filter.BlobFee = uint256.MustFromBig(eip4844.CalcBlobFee(*env.header.ExcessBlobGas))
	}
	filter.OnlyPlainTxs, filter.OnlyBlobTxs = true, false
	pendingPlainTxs := w.eth.TxPool().Pending(filter)

	filter.OnlyPlainTxs, filter.OnlyBlobTxs = false, true
	pendingBlobTxs := w.eth.TxPool().Pending(filter)

	// Split the pending transactions into locals and remotes.
	localPlainTxs, remotePlainTxs := make(map[common.Address][]*txpool.LazyTransaction), pendingPlainTxs
	localBlobTxs, remoteBlobTxs := make(map[common.Address][]*txpool.LazyTransaction), pendingBlobTxs
	for _, account := range w.eth.TxPool().Locals() {
		if txs := remotePlainTxs[account]; len(txs) > 0 {
			delete(remotePlainTxs, account)
			localPlainTxs[account] = txs
		}
		if txs := remoteBlobTxs[account]; len(txs) > 0 {
			delete(remoteBlobTxs, account)
			localBlobTxs[account] = txs
		}
	}
	// Fill the block with all available pending transactions.
	if len(localPlainTxs) > 0 || len(localBlobTxs) > 0 {
		plainTxs := newTransactionsByPriceAndNonce(env.signer, localPlainTxs, env.header.BaseFee)
		blobTxs := newTransactionsByPriceAndNonce(env.signer, localBlobTxs, env.header.BaseFee)

		w.commitTransactions(env, plainTxs, blobTxs, env.header.Coinbase)
	}
	if len(remotePlainTxs) > 0 || len(remoteBlobTxs) > 0 {
		plainTxs := newTransactionsByPriceAndNonce(env.signer, remotePlainTxs, env.header.BaseFee)
		blobTxs := newTransactionsByPriceAndNonce(env.signer, remoteBlobTxs, env.header.BaseFee)

		w.commitTransactions(env, plainTxs, blobTxs, env.header.Coinbase)
	}

	return w.commit(env)
}

func (w *worker) createCurrentEnvironment(predicateContext *precompileconfig.PredicateContext, parent *types.Header, header *types.Header, feeConfig commontype.FeeConfig, tstart time.Time) (*environment, error) {
	currentState, err := w.chain.StateAt(parent.Root)
	if err != nil {
		return nil, err
	}
	chainConfig := params.GetExtra(w.chainConfig)
	capacity, err := customheader.GasCapacity(chainConfig, feeConfig, parent, header.Time)
	if err != nil {
		return nil, fmt.Errorf("calculating gas capacity: %w", err)
	}
	numPrefetchers := w.chain.CacheConfig().TriePrefetcherParallelism
	currentState.StartPrefetcher("miner", state.WithConcurrentWorkers(numPrefetchers))
	return &environment{
		signer:           types.MakeSigner(w.chainConfig, header.Number, header.Time),
		state:            currentState,
		parent:           parent,
		header:           header,
		tcount:           0,
		gasPool:          new(core.GasPool).AddGas(capacity),
		rules:            w.chainConfig.Rules(header.Number, params.IsMergeTODO, header.Time),
		predicateContext: predicateContext,
		predicateResults: predicate.NewResults(),
		start:            tstart,
	}, nil
}

func (w *worker) commitTransaction(env *environment, tx *types.Transaction, coinbase common.Address) ([]*types.Log, error) {
	if tx.Type() == types.BlobTxType {
		return w.commitBlobTransaction(env, tx, coinbase)
	}
	receipt, err := w.applyTransaction(env, tx, coinbase)
	if err != nil {
		return nil, err
	}
	env.txs = append(env.txs, tx)
	env.receipts = append(env.receipts, receipt)
	env.size += tx.Size()
	return receipt.Logs, nil
}

func (w *worker) commitBlobTransaction(env *environment, tx *types.Transaction, coinbase common.Address) ([]*types.Log, error) {
	sc := tx.BlobTxSidecar()
	if sc == nil {
		panic("blob transaction without blobs in miner")
	}
	// Checking against blob gas limit: It's kind of ugly to perform this check here, but there
	// isn't really a better place right now. The blob gas limit is checked at block validation time
	// and not during execution. This means core.ApplyTransaction will not return an error if the
	// tx has too many blobs. So we have to explicitly check it here.
	if (env.blobs+len(sc.Blobs))*ethparams.BlobTxBlobGasPerBlob > ethparams.MaxBlobGasPerBlock {
		return nil, errors.New("max data blobs reached")
	}
	receipt, err := w.applyTransaction(env, tx, coinbase)
	if err != nil {
		return nil, err
	}
	env.txs = append(env.txs, tx.WithoutBlobTxSidecar())
	env.receipts = append(env.receipts, receipt)
	env.sidecars = append(env.sidecars, sc)
	env.blobs += len(sc.Blobs)
	*env.header.BlobGasUsed += receipt.BlobGasUsed
	return receipt.Logs, nil
}

// applyTransaction runs the transaction. If execution fails, state and gas pool are reverted.
func (w *worker) applyTransaction(env *environment, tx *types.Transaction, coinbase common.Address) (*types.Receipt, error) {
	var (
		snap         = env.state.Snapshot()
		gp           = env.gasPool.Gas()
		blockContext vm.BlockContext
	)

	if params.GetRulesExtra(env.rules).IsDurango {
		results, err := core.CheckPredicates(env.rules, env.predicateContext, tx)
		if err != nil {
			log.Debug("Transaction predicate failed verification in miner", "tx", tx.Hash(), "err", err)
			return nil, err
		}
		env.predicateResults.SetTxResults(tx.Hash(), results)

		predicateResultsBytes, err := env.predicateResults.Bytes()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal predicate results: %w", err)
		}
		blockContext = core.NewEVMBlockContextWithPredicateResults(env.header, w.chain, &coinbase, predicateResultsBytes)
	} else {
		blockContext = core.NewEVMBlockContext(env.header, w.chain, &coinbase)
	}

	receipt, err := core.ApplyTransaction(w.chainConfig, w.chain, blockContext, env.gasPool, env.state, env.header, tx, &env.header.GasUsed, *w.chain.GetVMConfig())
	if err != nil {
		env.state.RevertToSnapshot(snap)
		env.gasPool.SetGas(gp)
		env.predicateResults.DeleteTxResults(tx.Hash())
	}
	return receipt, err
}

func (w *worker) commitTransactions(env *environment, plainTxs, blobTxs *transactionsByPriceAndNonce, coinbase common.Address) {
	for {
		// If we don't have enough gas for any further transactions then we're done.
		if env.gasPool.Gas() < ethparams.TxGas {
			log.Trace("Not enough gas for further transactions", "have", env.gasPool, "want", ethparams.TxGas)
			break
		}
		// If we don't have enough blob space for any further blob transactions,
		// skip that list altogether
		if !blobTxs.Empty() && env.blobs*ethparams.BlobTxBlobGasPerBlob >= ethparams.MaxBlobGasPerBlock {
			log.Trace("Not enough blob space for further blob transactions")
			blobTxs.Clear()
			// Fall though to pick up any plain txs
		}
		// If we don't have enough blob space for any further blob transactions,
		// skip that list altogether
		if !blobTxs.Empty() && env.blobs*ethparams.BlobTxBlobGasPerBlob >= ethparams.MaxBlobGasPerBlock {
			log.Trace("Not enough blob space for further blob transactions")
			blobTxs.Clear()
			// Fall though to pick up any plain txs
		}
		// Retrieve the next transaction and abort if all done.
		var (
			ltx *txpool.LazyTransaction
			txs *transactionsByPriceAndNonce
		)
		pltx, ptip := plainTxs.Peek()
		bltx, btip := blobTxs.Peek()

		switch {
		case pltx == nil:
			txs, ltx = blobTxs, bltx
		case bltx == nil:
			txs, ltx = plainTxs, pltx
		default:
			if ptip.Lt(btip) {
				txs, ltx = blobTxs, bltx
			} else {
				txs, ltx = plainTxs, pltx
			}
		}
		if ltx == nil {
			break
		}
		// If we don't have enough space for the next transaction, skip the account.
		if env.gasPool.Gas() < ltx.Gas {
			log.Trace("Not enough gas left for transaction", "hash", ltx.Hash, "left", env.gasPool.Gas(), "needed", ltx.Gas)
			txs.Pop()
			continue
		}
		if left := uint64(ethparams.MaxBlobGasPerBlock - env.blobs*ethparams.BlobTxBlobGasPerBlob); left < ltx.BlobGas {
			log.Trace("Not enough blob gas left for transaction", "hash", ltx.Hash, "left", left, "needed", ltx.BlobGas)
			txs.Pop()
			continue
		}
		// Transaction seems to fit, pull it up from the pool
		tx := ltx.Resolve()
		if tx == nil {
			log.Trace("Ignoring evicted transaction", "hash", ltx.Hash)
			txs.Pop()
			continue
		}
		// Abort transaction if it won't fit in the block and continue to search for a smaller
		// transaction that will fit.
		if totalTxsSize := env.size + tx.Size(); totalTxsSize > targetTxsSize {
			log.Trace("Skipping transaction that would exceed target size", "hash", tx.Hash(), "totalTxsSize", totalTxsSize, "txSize", tx.Size())
			txs.Pop()
			continue
		}

		// Error may be ignored here. The error has already been checked
		// during transaction acceptance is the transaction pool.
		from, _ := types.Sender(env.signer, tx)

		// Check whether the tx is replay protected. If we're not in the EIP155 hf
		// phase, start ignoring the sender until we do.
		if tx.Protected() && !w.chainConfig.IsEIP155(env.header.Number) {
			log.Trace("Ignoring replay protected transaction", "hash", ltx.Hash, "eip155", w.chainConfig.EIP155Block)
			txs.Pop()
			continue
		}

		// Start executing the transaction
		env.state.SetTxContext(tx.Hash(), env.tcount)

		_, err := w.commitTransaction(env, tx, coinbase)
		switch {
		case errors.Is(err, core.ErrNonceTooLow):
			// New head notification data race between the transaction pool and miner, shift
			log.Trace("Skipping transaction with low nonce", "hash", ltx.Hash, "sender", from, "nonce", tx.Nonce())
			txs.Shift()

		case errors.Is(err, nil):
			env.tcount++
			txs.Shift()

		default:
			// Transaction is regarded as invalid, drop all consecutive transactions from
			// the same sender because of `nonce-too-high` clause.
			log.Debug("Transaction failed, account skipped", "hash", ltx.Hash, "err", err)
			txs.Pop()
		}
	}
}

// commit runs any post-transaction state modifications, assembles the final block
// and commits new work if consensus engine is running.
func (w *worker) commit(env *environment) (*types.Block, error) {
	if params.GetRulesExtra(env.rules).IsDurango {
		predicateResultsBytes, err := env.predicateResults.Bytes()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal predicate results: %w", err)
		}
		env.header.Extra = append(env.header.Extra, predicateResultsBytes...)
	}
	// Deep copy receipts here to avoid interaction between different tasks.
	receipts := copyReceipts(env.receipts)
	block, err := w.engine.FinalizeAndAssemble(w.chain, env.header, env.parent, env.state, env.txs, nil, receipts)
	if err != nil {
		return nil, err
	}

	return w.handleResult(env, block, time.Now(), receipts)
}

func (w *worker) handleResult(env *environment, block *types.Block, createdAt time.Time, unfinishedReceipts []*types.Receipt) (*types.Block, error) {
	// Short circuit when receiving duplicate result caused by resubmitting.
	if !w.config.TestOnlyAllowDuplicateBlocks && w.chain.HasBlock(block.Hash(), block.NumberU64()) {
		return nil, fmt.Errorf("produced duplicate block (Hash: %s, Number %d)", block.Hash(), block.NumberU64())
	}
	// Different block could share same sealhash, deep copy here to prevent write-write conflict.
	var (
		hash     = block.Hash()
		receipts = make([]*types.Receipt, len(unfinishedReceipts))
		logs     []*types.Log
	)
	for i, unfinishedReceipt := range unfinishedReceipts {
		receipt := new(types.Receipt)
		receipts[i] = receipt
		*receipt = *unfinishedReceipt

		// add block location fields
		receipt.BlockHash = hash
		receipt.BlockNumber = block.Number()
		receipt.TransactionIndex = uint(i)

		// Update the block hash in all logs since it is now available and not when the
		// receipt/log of individual transactions were created.
		receipt.Logs = make([]*types.Log, len(unfinishedReceipt.Logs))
		for j, unfinishedLog := range unfinishedReceipt.Logs {
			log := new(types.Log)
			receipt.Logs[j] = log
			*log = *unfinishedLog
			log.BlockHash = hash
		}
		logs = append(logs, receipt.Logs...)
	}

	feesInEther, err := core.TotalFeesFloat(block, receipts)
	if err != nil {
		log.Error("TotalFeesFloat error: %s", err)
	}
	log.Info("Commit new mining work", "number", block.Number(), "hash", hash,
		"uncles", 0, "txs", env.tcount,
		"gas", block.GasUsed(), "fees", feesInEther,
		"elapsed", common.PrettyDuration(time.Since(env.start)))

	// Note: the miner no longer emits a NewMinedBlock event. Instead the caller
	// is responsible for running any additional verification and then inserting
	// the block with InsertChain, which will also emit a new head event.
	return block, nil
}

// copyReceipts makes a deep copy of the given receipts.
func copyReceipts(receipts []*types.Receipt) []*types.Receipt {
	result := make([]*types.Receipt, len(receipts))
	for i, l := range receipts {
		cpy := *l
		result[i] = &cpy
	}
	return result
}
