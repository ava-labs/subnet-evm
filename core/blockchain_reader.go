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
// Copyright 2021 The go-ethereum Authors
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

package core

import (
	"math/big"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/rawdb"
	"github.com/ava-labs/libevm/core/types"
	"github.com/ava-labs/libevm/core/vm"
	"github.com/ava-labs/libevm/event"
	"github.com/ava-labs/libevm/triedb"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/consensus"
	"github.com/ava-labs/subnet-evm/constants"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/core/state/snapshot"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contracts/feemanager"
	"github.com/ava-labs/subnet-evm/precompile/contracts/rewardmanager"
)

// CurrentHeader retrieves the current head header of the canonical chain. The
// header is retrieved from the HeaderChain's internal cache.
func (bc *BlockChain) CurrentHeader() *types.Header {
	return bc.hc.CurrentHeader()
}

// CurrentBlock retrieves the current head block of the canonical chain. The
// block is retrieved from the blockchain's internal cache.
func (bc *BlockChain) CurrentBlock() *types.Header {
	return bc.currentBlock.Load()
}

// HasHeader checks if a block header is present in the database or not, caching
// it if present.
func (bc *BlockChain) HasHeader(hash common.Hash, number uint64) bool {
	return bc.hc.HasHeader(hash, number)
}

// GetHeader retrieves a block header from the database by hash and number,
// caching it if found.
func (bc *BlockChain) GetHeader(hash common.Hash, number uint64) *types.Header {
	return bc.hc.GetHeader(hash, number)
}

// GetHeaderByHash retrieves a block header from the database by hash, caching it if
// found.
func (bc *BlockChain) GetHeaderByHash(hash common.Hash) *types.Header {
	return bc.hc.GetHeaderByHash(hash)
}

// GetHeaderByNumber retrieves a block header from the database by number,
// caching it (associated with its hash) if found.
func (bc *BlockChain) GetHeaderByNumber(number uint64) *types.Header {
	return bc.hc.GetHeaderByNumber(number)
}

// GetBody retrieves a block body (transactions and uncles) from the database by
// hash, caching it if found.
func (bc *BlockChain) GetBody(hash common.Hash) *types.Body {
	// Short circuit if the body's already in the cache, retrieve otherwise
	if cached, ok := bc.bodyCache.Get(hash); ok {
		return cached
	}
	number := bc.hc.GetBlockNumber(hash)
	if number == nil {
		return nil
	}
	body := rawdb.ReadBody(bc.db, hash, *number)
	if body == nil {
		return nil
	}
	// Cache the found body for next time and return
	bc.bodyCache.Add(hash, body)
	return body
}

// HasBlock checks if a block is fully present in the database or not.
func (bc *BlockChain) HasBlock(hash common.Hash, number uint64) bool {
	if bc.blockCache.Contains(hash) {
		return true
	}
	if !bc.HasHeader(hash, number) {
		return false
	}
	return rawdb.HasBody(bc.db, hash, number)
}

// HasFastBlock checks if a fast block is fully present in the database or not.
func (bc *BlockChain) HasFastBlock(hash common.Hash, number uint64) bool {
	if !bc.HasBlock(hash, number) {
		return false
	}
	if bc.receiptsCache.Contains(hash) {
		return true
	}
	return rawdb.HasReceipts(bc.db, hash, number)
}

// GetBlock retrieves a block from the database by hash and number,
// caching it if found.
func (bc *BlockChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	// Short circuit if the block's already in the cache, retrieve otherwise
	if block, ok := bc.blockCache.Get(hash); ok {
		return block
	}
	block := rawdb.ReadBlock(bc.db, hash, number)
	if block == nil {
		return nil
	}
	// Cache the found block for next time and return
	bc.blockCache.Add(block.Hash(), block)
	return block
}

// GetBlockByHash retrieves a block from the database by hash, caching it if found.
func (bc *BlockChain) GetBlockByHash(hash common.Hash) *types.Block {
	number := bc.hc.GetBlockNumber(hash)
	if number == nil {
		return nil
	}
	return bc.GetBlock(hash, *number)
}

// GetBlockByNumber retrieves a block from the database by number, caching it
// (associated with its hash) if found.
func (bc *BlockChain) GetBlockByNumber(number uint64) *types.Block {
	hash := rawdb.ReadCanonicalHash(bc.db, number)
	if hash == (common.Hash{}) {
		return nil
	}
	return bc.GetBlock(hash, number)
}

// GetBlocksFromHash returns the block corresponding to hash and up to n-1 ancestors.
// [deprecated by eth/62]
func (bc *BlockChain) GetBlocksFromHash(hash common.Hash, n int) (blocks []*types.Block) {
	number := bc.hc.GetBlockNumber(hash)
	if number == nil {
		return nil
	}
	for i := 0; i < n; i++ {
		block := bc.GetBlock(hash, *number)
		if block == nil {
			break
		}
		blocks = append(blocks, block)
		hash = block.ParentHash()
		*number--
	}
	return
}

// GetReceiptsByHash retrieves the receipts for all transactions in a given block.
func (bc *BlockChain) GetReceiptsByHash(hash common.Hash) types.Receipts {
	if receipts, ok := bc.receiptsCache.Get(hash); ok {
		return receipts
	}
	number := rawdb.ReadHeaderNumber(bc.db, hash)
	if number == nil {
		return nil
	}
	header := bc.GetHeader(hash, *number)
	if header == nil {
		return nil
	}
	receipts := rawdb.ReadReceipts(bc.db, hash, *number, header.Time, bc.chainConfig)
	if receipts == nil {
		return nil
	}
	bc.receiptsCache.Add(hash, receipts)
	return receipts
}

// GetCanonicalHash returns the canonical hash for a given block number
func (bc *BlockChain) GetCanonicalHash(number uint64) common.Hash {
	return bc.hc.GetCanonicalHash(number)
}

// GetTransactionLookup retrieves the lookup along with the transaction
// itself associate with the given transaction hash.
//
// An error will be returned if the transaction is not found, and background
// indexing for transactions is still in progress. The transaction might be
// reachable shortly once it's indexed.
//
// A null will be returned in the transaction is not found and background
// transaction indexing is already finished. The transaction is not existent
// from the node's perspective.
func (bc *BlockChain) GetTransactionLookup(hash common.Hash) (*rawdb.LegacyTxLookupEntry, *types.Transaction, error) {
	// Short circuit if the txlookup already in the cache, retrieve otherwise
	if item, exist := bc.txLookupCache.Get(hash); exist {
		return item.lookup, item.transaction, nil
	}
	tx, blockHash, blockNumber, txIndex := rawdb.ReadTransaction(bc.db, hash)
	if tx == nil {
		// The transaction is already indexed, the transaction is either
		// not existent or not in the range of index, returning null.
		return nil, nil, nil
	}
	lookup := &rawdb.LegacyTxLookupEntry{
		BlockHash:  blockHash,
		BlockIndex: blockNumber,
		Index:      txIndex,
	}
	bc.txLookupCache.Add(hash, txLookup{
		lookup:      lookup,
		transaction: tx,
	})
	return lookup, tx, nil
}

// HasState checks if state trie is fully present in the database or not.
func (bc *BlockChain) HasState(hash common.Hash) bool {
	_, err := bc.stateCache.OpenTrie(hash)
	return err == nil
}

// HasBlockAndState checks if a block and associated state trie is fully present
// in the database or not, caching it if present.
func (bc *BlockChain) HasBlockAndState(hash common.Hash, number uint64) bool {
	// Check first that the block itself is known
	block := bc.GetBlock(hash, number)
	if block == nil {
		return false
	}
	return bc.HasState(block.Root())
}

// State returns a new mutable state based on the current HEAD block.
func (bc *BlockChain) State() (*state.StateDB, error) {
	return bc.StateAt(bc.CurrentBlock().Root)
}

// StateAt returns a new mutable state based on a particular point in time.
func (bc *BlockChain) StateAt(root common.Hash) (*state.StateDB, error) {
	return state.New(root, bc.stateCache, bc.snaps)
}

// Config retrieves the chain's fork configuration.
func (bc *BlockChain) Config() *params.ChainConfig { return bc.chainConfig }

// Engine retrieves the blockchain's consensus engine.
func (bc *BlockChain) Engine() consensus.Engine { return bc.engine }

// Snapshots returns the blockchain snapshot tree.
func (bc *BlockChain) Snapshots() *snapshot.Tree {
	return bc.snaps
}

// Validator returns the current validator.
func (bc *BlockChain) Validator() Validator {
	return bc.validator
}

// Processor returns the current processor.
func (bc *BlockChain) Processor() Processor {
	return bc.processor
}

// StateCache returns the caching database underpinning the blockchain instance.
func (bc *BlockChain) StateCache() state.Database {
	return bc.stateCache
}

// GasLimit returns the gas limit of the current HEAD block.
func (bc *BlockChain) GasLimit() uint64 {
	return bc.CurrentBlock().GasLimit
}

// Genesis retrieves the chain's genesis block.
func (bc *BlockChain) Genesis() *types.Block {
	return bc.genesisBlock
}

// GetVMConfig returns the block chain VM config.
func (bc *BlockChain) GetVMConfig() *vm.Config {
	return &bc.vmConfig
}

// TrieDB retrieves the low level trie database used for data storage.
func (bc *BlockChain) TrieDB() *triedb.Database {
	return bc.triedb
}

// HeaderChain returns the underlying header chain.
func (bc *BlockChain) HeaderChain() *HeaderChain {
	return bc.hc
}

// SubscribeRemovedLogsEvent registers a subscription of RemovedLogsEvent.
func (bc *BlockChain) SubscribeRemovedLogsEvent(ch chan<- RemovedLogsEvent) event.Subscription {
	return bc.scope.Track(bc.rmLogsFeed.Subscribe(ch))
}

// SubscribeChainEvent registers a subscription of ChainEvent.
func (bc *BlockChain) SubscribeChainEvent(ch chan<- ChainEvent) event.Subscription {
	return bc.scope.Track(bc.chainFeed.Subscribe(ch))
}

// SubscribeChainHeadEvent registers a subscription of ChainHeadEvent.
func (bc *BlockChain) SubscribeChainHeadEvent(ch chan<- ChainHeadEvent) event.Subscription {
	return bc.scope.Track(bc.chainHeadFeed.Subscribe(ch))
}

// SubscribeChainSideEvent registers a subscription of ChainSideEvent.
func (bc *BlockChain) SubscribeChainSideEvent(ch chan<- ChainSideEvent) event.Subscription {
	return bc.scope.Track(bc.chainSideFeed.Subscribe(ch))
}

// SubscribeLogsEvent registers a subscription of []*types.Log.
func (bc *BlockChain) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return bc.scope.Track(bc.logsFeed.Subscribe(ch))
}

// SubscribeBlockProcessingEvent registers a subscription of bool where true means
// block processing has started while false means it has stopped.
func (bc *BlockChain) SubscribeBlockProcessingEvent(ch chan<- bool) event.Subscription {
	return bc.scope.Track(bc.blockProcFeed.Subscribe(ch))
}

// SubscribeChainAcceptedEvent registers a subscription of ChainEvent.
func (bc *BlockChain) SubscribeChainAcceptedEvent(ch chan<- ChainEvent) event.Subscription {
	return bc.scope.Track(bc.chainAcceptedFeed.Subscribe(ch))
}

// SubscribeAcceptedLogsEvent registers a subscription of accepted []*types.Log.
func (bc *BlockChain) SubscribeAcceptedLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return bc.scope.Track(bc.logsAcceptedFeed.Subscribe(ch))
}

// SubscribeAcceptedTransactionEvent registers a subscription of accepted transactions
func (bc *BlockChain) SubscribeAcceptedTransactionEvent(ch chan<- NewTxsEvent) event.Subscription {
	return bc.scope.Track(bc.txAcceptedFeed.Subscribe(ch))
}

// GetFeeConfigAt returns the fee configuration and the last changed block number at [parent].
// If Subnet-EVM is not activated, returns default fee config and nil block number.
// If FeeManager is activated at [parent], returns the fee config in the precompile contract state.
// Otherwise returns the fee config in the chain config.
// Assumes that a valid configuration is stored when the precompile is activated.
func (bc *BlockChain) GetFeeConfigAt(parent *types.Header) (commontype.FeeConfig, *big.Int, error) {
	config := params.GetExtra(bc.Config())
	if !config.IsSubnetEVM(parent.Time) {
		return params.DefaultFeeConfig, nil, nil
	}
	if !config.IsPrecompileEnabled(feemanager.ContractAddress, parent.Time) {
		return config.FeeConfig, common.Big0, nil
	}

	// try to return it from the cache
	if cached, hit := bc.feeConfigCache.Get(parent.Root); hit {
		return cached.feeConfig, cached.lastChangedAt, nil
	}

	stateDB, err := bc.StateAt(parent.Root)
	if err != nil {
		return commontype.EmptyFeeConfig, nil, err
	}

	storedFeeConfig := feemanager.GetStoredFeeConfig(stateDB)
	// this should not return an invalid fee config since it's assumed that
	// StoreFeeConfig returns an error when an invalid fee config is attempted to be stored.
	// However an external stateDB call can modify the contract state.
	// This check is added to add a defense in-depth.
	if err := storedFeeConfig.Verify(); err != nil {
		return commontype.EmptyFeeConfig, nil, err
	}
	lastChangedAt := feemanager.GetFeeConfigLastChangedAt(stateDB)
	cacheable := &cacheableFeeConfig{feeConfig: storedFeeConfig, lastChangedAt: lastChangedAt}
	// add it to the cache
	bc.feeConfigCache.Add(parent.Root, cacheable)
	return storedFeeConfig, lastChangedAt, nil
}

// GetCoinbaseAt returns the configured coinbase address at [parent].
// If RewardManager is activated at [parent], returns the reward manager config in the precompile contract state.
// If fee recipients are allowed, returns true in the second return value.
func (bc *BlockChain) GetCoinbaseAt(parent *types.Header) (common.Address, bool, error) {
	configExtra := params.GetExtra(bc.Config())
	if !configExtra.IsSubnetEVM(parent.Time) {
		return constants.BlackholeAddr, false, nil
	}

	if !configExtra.IsPrecompileEnabled(rewardmanager.ContractAddress, parent.Time) {
		if configExtra.AllowFeeRecipients {
			return common.Address{}, true, nil
		} else {
			return constants.BlackholeAddr, false, nil
		}
	}

	// try to return it from the cache
	if cached, hit := bc.coinbaseConfigCache.Get(parent.Root); hit {
		return cached.coinbaseAddress, cached.allowFeeRecipients, nil
	}

	stateDB, err := bc.StateAt(parent.Root)
	if err != nil {
		return common.Address{}, false, err
	}
	rewardAddress, feeRecipients := rewardmanager.GetStoredRewardAddress(stateDB)

	cacheable := &cacheableCoinbaseConfig{coinbaseAddress: rewardAddress, allowFeeRecipients: feeRecipients}
	bc.coinbaseConfigCache.Add(parent.Root, cacheable)
	return rewardAddress, feeRecipients, nil
}

// GetLogs fetches all logs from a given block.
func (bc *BlockChain) GetLogs(hash common.Hash, number uint64) [][]*types.Log {
	logs, ok := bc.acceptedLogsCache.Get(hash) // this cache is thread-safe
	if ok {
		return logs
	}
	block := bc.GetBlockByHash(hash)
	if block == nil {
		return nil
	}
	logs = bc.collectUnflattenedLogs(block, false)
	return logs
}
