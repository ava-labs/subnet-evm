// (c) 2019-2020, Ava Labs, Inc.
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
	"github.com/ava-labs/subnet-evm/consensus"
	"github.com/ava-labs/subnet-evm/consensus/dummy"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
)

const (
	targetTxsSize = 1800 * units.KiB
)

// environment is the worker's current environment and holds all of the current state information.
type environment struct {
	signer types.Signer

	state   *state.StateDB // apply state changes here
	tcount  int            // tx count in cycle
	gasPool *core.GasPool  // available gas used to pack transactions

	parent   *types.Header
	header   *types.Header
	txs      []*types.Transaction
	receipts []*types.Receipt
	size     common.StorageSize

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
	mux      *event.TypeMux // TODO replace
	mu       sync.RWMutex   // The lock used to protect the coinbase and extra fields
	coinbase common.Address
	clock    *mockable.Clock // Allows us mock the clock for testing
}

func newWorker(config *Config, chainConfig *params.ChainConfig, engine consensus.Engine, eth Backend, mux *event.TypeMux, clock *mockable.Clock) *worker {
	worker := &worker{
		config:      config,
		chainConfig: chainConfig,
		engine:      engine,
		eth:         eth,
		mux:         mux,
		chain:       eth.BlockChain(),
		clock:       clock,
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
func (w *worker) commitNewWork(predicateContext *precompileconfig.ProposerPredicateContext) (*types.Block, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	tstart := w.clock.Time()
	timestamp := uint64(tstart.Unix())
	parent := w.chain.CurrentBlock()
	// Note: in order to support asynchronous block production, blocks are allowed to have
	// the same timestamp as their parent. This allows more than one block to be produced
	// per second.
	if parent.Time() >= timestamp {
		timestamp = parent.Time()
	}

	bigTimestamp := new(big.Int).SetUint64(timestamp)
	var gasLimit uint64
	// The fee manager relies on the state of the parent block to set the fee config
	// because the fee config may be changed by the current block.
	feeConfig, _, err := w.chain.GetFeeConfigAt(parent.Header())
	if err != nil {
		return nil, err
	}
	configuredGasLimit := feeConfig.GasLimit.Uint64()
	if w.chainConfig.IsSubnetEVM(bigTimestamp) {
		gasLimit = configuredGasLimit
	} else {
		// The gas limit is set in SubnetEVMGasLimit because the ceiling and floor were set to the same value
		// such that the gas limit converged to it. Since this is hardbaked now, we remove the ability to configure it.
		gasLimit = core.CalcGasLimit(parent.GasUsed(), parent.GasLimit(), configuredGasLimit, configuredGasLimit)
	}
	num := parent.Number()
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     num.Add(num, common.Big1),
		GasLimit:   gasLimit,
		Extra:      nil,
		Time:       timestamp,
	}

	if w.chainConfig.IsSubnetEVM(bigTimestamp) {
		var err error
		header.Extra, header.BaseFee, err = dummy.CalcBaseFee(w.chainConfig, feeConfig, parent.Header(), timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate new base fee: %w", err)
		}
	}

	if w.coinbase == (common.Address{}) {
		return nil, errors.New("cannot mine without etherbase")
	}
	header.Coinbase = w.coinbase

	configuredCoinbase, isAllowFeeRecipient, err := w.chain.GetCoinbaseAt(parent.Header())
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

	env, err := w.createCurrentEnvironment(parent, header, tstart)
	if err != nil {
		return nil, fmt.Errorf("failed to create new current environment: %w", err)
	}
	// Configure any upgrades that should go into effect during this block.
	err = core.ApplyUpgrades(w.chainConfig, new(big.Int).SetUint64(parent.Time()), types.NewBlockWithHeader(header), env.state)
	if err != nil {
		log.Error("failed to configure precompiles mining new block", "parent", parent.Hash(), "number", header.Number, "timestamp", header.Time, "err", err)
		return nil, err
	}

	// Get the pending txs from TxPool
	pending := w.eth.TxPool().Pending(true)
	// Filter out transactions that don't satisfy predicateContext and remove them from TxPool
	rules := w.chainConfig.AvalancheRules(header.Number, new(big.Int).SetUint64(header.Time))
	pending = w.enforcePredicates(rules, predicateContext, pending)

	// Split the pending transactions into locals and remotes
	localTxs := make(map[common.Address]types.Transactions)
	remoteTxs := pending
	for _, account := range w.eth.TxPool().Locals() {
		if txs := remoteTxs[account]; len(txs) > 0 {
			delete(remoteTxs, account)
			localTxs[account] = txs
		}
	}
	if len(localTxs) > 0 {
		txs := types.NewTransactionsByPriceAndNonce(env.signer, localTxs, header.BaseFee)
		w.commitTransactions(env, txs, header.Coinbase)
	}
	if len(remoteTxs) > 0 {
		txs := types.NewTransactionsByPriceAndNonce(env.signer, remoteTxs, header.BaseFee)
		w.commitTransactions(env, txs, header.Coinbase)
	}

	return w.commit(env)
}

func (w *worker) createCurrentEnvironment(parent *types.Block, header *types.Header, tstart time.Time) (*environment, error) {
	state, err := w.chain.StateAt(parent.Root())
	if err != nil {
		return nil, err
	}
	return &environment{
		signer:  types.MakeSigner(w.chainConfig, header.Number, new(big.Int).SetUint64(header.Time)),
		state:   state,
		parent:  parent.Header(),
		header:  header,
		tcount:  0,
		gasPool: new(core.GasPool).AddGas(header.GasLimit),
		start:   tstart,
	}, nil
}

func (w *worker) commitTransaction(env *environment, tx *types.Transaction, coinbase common.Address) ([]*types.Log, error) {
	snap := env.state.Snapshot()

	receipt, err := core.ApplyTransaction(w.chainConfig, w.chain, &coinbase, env.gasPool, env.state, env.header, tx, &env.header.GasUsed, *w.chain.GetVMConfig())
	if err != nil {
		env.state.RevertToSnapshot(snap)
		return nil, err
	}
	env.txs = append(env.txs, tx)
	env.receipts = append(env.receipts, receipt)
	env.size += tx.Size()

	return receipt.Logs, nil
}

func (w *worker) commitTransactions(env *environment, txs *types.TransactionsByPriceAndNonce, coinbase common.Address) {
	for {
		// If we don't have enough gas for any further transactions then we're done
		if env.gasPool.Gas() < params.TxGas {
			log.Trace("Not enough gas for further transactions", "have", env.gasPool, "want", params.TxGas)
			break
		}
		// Retrieve the next transaction and abort if all done
		tx := txs.Peek()
		if tx == nil {
			break
		}
		// Abort transaction if it won't fit in the block and continue to search for a smaller
		// transction that will fit.
		if totalTxsSize := env.size + tx.Size(); totalTxsSize > targetTxsSize {
			log.Trace("Skipping transaction that would exceed target size", "hash", tx.Hash(), "totalTxsSize", totalTxsSize, "txSize", tx.Size())

			txs.Pop()
			continue
		}
		// Error may be ignored here. The error has already been checked
		// during transaction acceptance is the transaction pool.
		//
		// We use the eip155 signer regardless of the current hf.
		from, _ := types.Sender(env.signer, tx)
		// Check whether the tx is replay protected. If we're not in the EIP155 hf
		// phase, start ignoring the sender until we do.
		if tx.Protected() && !w.chainConfig.IsEIP155(env.header.Number) {
			log.Trace("Ignoring reply protected transaction", "hash", tx.Hash(), "eip155", w.chainConfig.EIP155Block)

			txs.Pop()
			continue
		}
		// Start executing the transaction
		env.state.Prepare(tx.Hash(), env.tcount)

		_, err := w.commitTransaction(env, tx, coinbase)
		switch {
		case errors.Is(err, core.ErrGasLimitReached):
			// Pop the current out-of-gas transaction without shifting in the next from the account
			log.Trace("Gas limit exceeded for current block", "sender", from)
			txs.Pop()

		case errors.Is(err, core.ErrNonceTooLow):
			// New head notification data race between the transaction pool and miner, shift
			log.Trace("Skipping transaction with low nonce", "sender", from, "nonce", tx.Nonce())
			txs.Shift()

		case errors.Is(err, core.ErrNonceTooHigh):
			// Reorg notification data race between the transaction pool and miner, skip account =
			log.Trace("Skipping account with high nonce", "sender", from, "nonce", tx.Nonce())
			txs.Pop()

		case errors.Is(err, nil):
			env.tcount++
			txs.Shift()

		case errors.Is(err, core.ErrTxTypeNotSupported):
			// Pop the unsupported transaction without shifting in the next from the account
			log.Trace("Skipping unsupported transaction type", "sender", from, "type", tx.Type())
			txs.Pop()

		default:
			// Strange error, discard the transaction and get the next in line (note, the
			// nonce-too-high clause will prevent us from executing in vain).
			log.Debug("Transaction failed, account skipped", "hash", tx.Hash(), "err", err)
			txs.Shift()
		}
	}
}

// commit runs any post-transaction state modifications, assembles the final block
// and commits new work if consensus engine is running.
func (w *worker) commit(env *environment) (*types.Block, error) {
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
	if w.chain.HasBlock(block.Hash(), block.NumberU64()) {
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

	log.Info("Commit new mining work", "number", block.Number(), "hash", hash, "timestamp", block.Time(), "uncles", 0, "txs", env.tcount,
		"gas", block.GasUsed(), "fees", totalFees(block, receipts), "elapsed", common.PrettyDuration(time.Since(env.start)))

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

// totalFees computes total consumed fees in ETH. Block transactions and receipts have to have the same order.
func totalFees(block *types.Block, receipts []*types.Receipt) *big.Float {
	feesWei := new(big.Int)
	for i, tx := range block.Transactions() {
		feesWei.Add(feesWei, new(big.Int).Mul(new(big.Int).SetUint64(receipts[i].GasUsed), tx.GasPrice()))
	}
	return new(big.Float).Quo(new(big.Float).SetInt(feesWei), new(big.Float).SetInt(big.NewInt(params.Ether)))
}

// enforcePredicates takes a set of pending transactions (grouped by sender, and ordered by nonce)
// and returns the subset of those transactions (following the same grouping) that satisfy predicateContext.
// Any transaction that fails predicate verification will be removed from the tx pool and excluded
// from the return value.
// Transactions with a nonce that follows a removed transaction will be added back to the future
// queue of the tx pool.
func (w *worker) enforcePredicates(
	rules params.Rules,
	predicateContext *precompileconfig.ProposerPredicateContext,
	pending map[common.Address]types.Transactions,
) map[common.Address]types.Transactions {
	// Short circuit early if there are no precompile predicates to verify and return the
	// unmodified pending transactions.
	if !rules.PredicatesExist() {
		return pending
	}
	result := make(map[common.Address]types.Transactions, len(pending))
	for addr, txs := range pending {
		for i, tx := range txs {
			if err := core.CheckPredicates(rules, predicateContext, tx); err != nil {
				log.Debug("Transaction predicate failed verification in miner", "sender", addr, "err", err)
				// If the transaction fails the predicate check, we remove the transaction from the mempool
				// and move all transactions from the same address with a subsequent nonce back to the
				// future queue of the transaction pool.
				w.eth.TxPool().RemoveTx(tx.Hash())
				txs = txs[:i] // Cut off any transactions past the failed predicate in the return value
				break
			}
		}
		if len(txs) > 0 {
			result[addr] = txs
		}
	}
	return result
}
