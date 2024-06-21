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
// Copyright 2014 The go-ethereum Authors
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

// Package core implements the Ethereum consensus protocol.
package chain

import (
	"errors"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"time"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/consensus"
	"github.com/ava-labs/subnet-evm/consensus/misc/eip4844"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/metrics"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/lru"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/trie"
)

var _ BlockChain = (*blockChain)(nil)

var (
	accountReadTimer   = metrics.NewRegisteredCounter("chain/account/reads", nil)
	accountHashTimer   = metrics.NewRegisteredCounter("chain/account/hashes", nil)
	accountUpdateTimer = metrics.NewRegisteredCounter("chain/account/updates", nil)
	accountCommitTimer = metrics.NewRegisteredCounter("chain/account/commits", nil)
	storageReadTimer   = metrics.NewRegisteredCounter("chain/storage/reads", nil)
	storageHashTimer   = metrics.NewRegisteredCounter("chain/storage/hashes", nil)
	storageUpdateTimer = metrics.NewRegisteredCounter("chain/storage/updates", nil)
	storageCommitTimer = metrics.NewRegisteredCounter("chain/storage/commits", nil)
	triedbCommitTimer  = metrics.NewRegisteredCounter("chain/triedb/commits", nil)

	blockInsertTimer            = metrics.NewRegisteredCounter("chain/block/inserts", nil)
	blockInsertCount            = metrics.NewRegisteredCounter("chain/block/inserts/count", nil)
	blockContentValidationTimer = metrics.NewRegisteredCounter("chain/block/validations/content", nil)
	blockStateInitTimer         = metrics.NewRegisteredCounter("chain/block/inits/state", nil)
	blockExecutionTimer         = metrics.NewRegisteredCounter("chain/block/executions", nil)
	blockTrieOpsTimer           = metrics.NewRegisteredCounter("chain/block/trie", nil)
	blockValidationTimer        = metrics.NewRegisteredCounter("chain/block/validations/state", nil)
	blockWriteTimer             = metrics.NewRegisteredCounter("chain/block/writes", nil)

	acceptorWorkTimer            = metrics.NewRegisteredCounter("chain/acceptor/work", nil)
	acceptorWorkCount            = metrics.NewRegisteredCounter("chain/acceptor/work/count", nil)
	processedBlockGasUsedCounter = metrics.NewRegisteredCounter("chain/block/gas/used/processed", nil)
	acceptedBlockGasUsedCounter  = metrics.NewRegisteredCounter("chain/block/gas/used/accepted", nil)
	badBlockCounter              = metrics.NewRegisteredCounter("chain/block/bad/count", nil)

	acceptedTxsCounter  = metrics.NewRegisteredCounter("chain/txs/accepted", nil)
	processedTxsCounter = metrics.NewRegisteredCounter("chain/txs/processed", nil)

	acceptedLogsCounter  = metrics.NewRegisteredCounter("chain/logs/accepted", nil)
	processedLogsCounter = metrics.NewRegisteredCounter("chain/logs/processed", nil)

	errFutureBlockUnsupported  = errors.New("future block insertion not supported")
	errCacheConfigNotSpecified = errors.New("must specify cache config")
	errInvalidOldChain         = errors.New("invalid old chain")
	errInvalidNewChain         = errors.New("invalid new chain")
)

const (
	blockCacheLimit          = 256
	receiptsCacheLimit       = 32
	feeConfigCacheLimit      = 256
	coinbaseConfigCacheLimit = 256
)

// cacheableFeeConfig encapsulates fee configuration itself and the block number that it has changed at,
// in order to cache them together.
type cacheableFeeConfig struct {
	feeConfig     commontype.FeeConfig
	lastChangedAt *big.Int
}

// cacheableCoinbaseConfig encapsulates coinbase address itself and allowFeeRecipient flag,
// in order to cache them together.
type cacheableCoinbaseConfig struct {
	coinbaseAddress    common.Address
	allowFeeRecipients bool
}

type blockChain struct {
	chainmu sync.RWMutex

	db                  ethdb.Database
	state               committableStateDB
	senderCacher        *core.TxSenderCacher
	hc                  *core.HeaderChain
	blockCache          *lru.Cache[common.Hash, *types.Block]             // Cache for the most recent entire blocks
	receiptsCache       *lru.Cache[common.Hash, []*types.Receipt]         // Cache for the most recent receipts per block
	feeConfigCache      *lru.Cache[common.Hash, *cacheableFeeConfig]      // Cache for the most recent feeConfig lookup data.
	coinbaseConfigCache *lru.Cache[common.Hash, *cacheableCoinbaseConfig] // Cache for the most recent coinbaseConfig lookup data.
	lastAccepted        *types.Block                                      // Prevents reorgs past this height
	validator           core.Validator                                    // Block and state validator interface
	processor           core.Processor                                    // Block transaction processor interface
	genesisBlock        *types.Block

	// TODO: should make a config struct?
	config      *params.ChainConfig
	cacheConfig *core.CacheConfig
	vmConfig    vm.Config
	engine      consensus.Engine

	scope             event.SubscriptionScope
	chainHeadFeed     event.Feed
	chainAcceptedFeed event.Feed
	logsAcceptedFeed  event.Feed
}

func NewBlockChain(
	chaindb ethdb.Database,
	committable committableStateDB,
	cacheConfig *core.CacheConfig,
	genesis *core.Genesis,
	engine consensus.Engine,
	vmConfig vm.Config,
	lastAcceptedHash common.Hash,
	skipChainConfigCheckCompatible bool,
) (*blockChain, error) {
	if cacheConfig == nil {
		return nil, errCacheConfigNotSpecified
	}
	// Setup the genesis block, commit the provided genesis specification
	// to database if the genesis block is not present yet, or load the
	// stored one from database.
	// Note: In go-ethereum, the code rewinds the chain on an incompatible config upgrade.
	// We don't do this and expect the node operator to always update their node's configuration
	// before network upgrades take effect.
	config, _, err := core.SetupGenesisBlockWithCommitable(
		chaindb, committable, genesis, lastAcceptedHash, skipChainConfigCheckCompatible,
	)
	if err != nil {
		return nil, err
	}
	bc := &blockChain{
		db:                  chaindb,
		state:               committable,
		config:              config,
		cacheConfig:         cacheConfig,
		receiptsCache:       lru.NewCache[common.Hash, []*types.Receipt](receiptsCacheLimit),
		blockCache:          lru.NewCache[common.Hash, *types.Block](blockCacheLimit),
		feeConfigCache:      lru.NewCache[common.Hash, *cacheableFeeConfig](feeConfigCacheLimit),
		coinbaseConfigCache: lru.NewCache[common.Hash, *cacheableCoinbaseConfig](coinbaseConfigCacheLimit),
		engine:              engine,
		vmConfig:            vmConfig,
		senderCacher:        core.NewTxSenderCacher(runtime.NumCPU()),
	}
	bc.validator = core.NewBlockValidator(config, bc, engine)
	bc.processor = core.NewStateProcessor(config, bc, engine)

	bc.hc, err = core.NewHeaderChain(chaindb, config, cacheConfig, engine)
	if err != nil {
		return nil, err
	}
	bc.genesisBlock = bc.GetBlockByNumber(0)
	if bc.genesisBlock == nil {
		return nil, core.ErrNoGenesis
	}

	// Re-generate current block state if it is missing
	if err := bc.loadLastState(lastAcceptedHash); err != nil {
		return nil, err
	}
	// Make sure the state associated with the block is available
	head := bc.CurrentBlock()
	if !bc.HasState(head.Root) {
		return nil, fmt.Errorf("head state missing %d:%s", head.Number, head.Hash())
	}

	bc.warmAcceptedCaches()
	return bc, nil
}

func (bc *blockChain) Stop() {
	bc.senderCacher.Shutdown()
	bc.scope.Close()
	if err := bc.state.Close(bc.CurrentBlock().Root); err != nil {
		log.Error("Failed to close state database", "err", err)
	}
}

func (bc *blockChain) InsertBlockManual(block *types.Block, writes bool) error {
	bc.chainmu.Lock()
	defer bc.chainmu.Unlock()

	return bc.insertBlock(block, writes)
}

func (bc *blockChain) insertBlock(block *types.Block, writes bool) error {
	start := time.Now()
	bc.senderCacher.Recover(types.MakeSigner(bc.config, block.Number(), block.Time()), block.Transactions())

	substart := time.Now()
	err := bc.engine.VerifyHeader(bc, block.Header())
	if err == nil {
		err = bc.validator.ValidateBody(block)
	}

	switch {
	case errors.Is(err, core.ErrKnownBlock):
		// even if the block is already known, we still need to generate the
		// snapshot layer and add a reference to the triedb, so we re-execute
		// the block. Note that insertBlock should only be called on a block
		// once if it returns nil
		if bc.newTip(block) {
			log.Debug("Setting head to be known block", "number", block.Number(), "hash", block.Hash())
		} else {
			log.Debug("Reprocessing already known block", "number", block.Number(), "hash", block.Hash())
		}

	// If an ancestor has been pruned, then this block cannot be acceptable.
	case errors.Is(err, consensus.ErrPrunedAncestor):
		return errors.New("side chain insertion is not supported")

	// Future blocks are not supported, but should not be reported, so we return an error
	// early here
	case errors.Is(err, consensus.ErrFutureBlock):
		return errFutureBlockUnsupported

	// Some other error occurred, abort
	case err != nil:
		bc.reportBlock(block, nil, err)
		return err
	}
	blockContentValidationTimer.Inc(time.Since(substart).Milliseconds())

	// No validation errors for the block
	var activeState *state.StateDB
	defer func() {
		// The chain importer is starting and stopping trie prefetchers. If a bad
		// block or other error is hit however, an early return may not properly
		// terminate the background threads. This defer ensures that we clean up
		// and dangling prefetcher, without defering each and holding on live refs.
		if activeState != nil {
			activeState.StopPrefetcher()
		}
	}()

	// Retrieve the parent block to determine which root to build state on
	substart = time.Now()
	parent := bc.GetHeader(block.ParentHash(), block.NumberU64()-1)

	// Instantiate the statedb to use for processing transactions
	statedb, err := state.New(parent.Root, bc.state, nil)
	if err != nil {
		return err
	}
	blockStateInitTimer.Inc(time.Since(substart).Milliseconds())

	// Enable prefetching to pull in trie node paths while processing transactions
	statedb.StartPrefetcher("chain", bc.cacheConfig.TriePrefetcherParallelism)
	activeState = statedb

	// Process block using the parent state as reference point
	pstart := time.Now()
	receipts, logs, usedGas, err := bc.processor.Process(block, parent, statedb, bc.vmConfig)
	if serr := statedb.Error(); serr != nil {
		log.Error("statedb error encountered", "err", serr, "number", block.Number(), "hash", block.Hash())
	}
	if err != nil {
		bc.reportBlock(block, receipts, err)
		return err
	}
	ptime := time.Since(pstart)

	// Validate the state using the default validator
	vstart := time.Now()
	if err := bc.validator.ValidateState(block, statedb, receipts, usedGas); err != nil {
		bc.reportBlock(block, receipts, err)
		return err
	}
	vtime := time.Since(vstart)

	// Update the metrics touched during block processing and validation
	accountReadTimer.Inc(statedb.AccountReads.Milliseconds())                  // Account reads are complete(in processing)
	storageReadTimer.Inc(statedb.StorageReads.Milliseconds())                  // Storage reads are complete(in processing)
	accountUpdateTimer.Inc(statedb.AccountUpdates.Milliseconds())              // Account updates are complete(in validation)
	storageUpdateTimer.Inc(statedb.StorageUpdates.Milliseconds())              // Storage updates are complete(in validation)
	accountHashTimer.Inc(statedb.AccountHashes.Milliseconds())                 // Account hashes are complete(in validation)
	storageHashTimer.Inc(statedb.StorageHashes.Milliseconds())                 // Storage hashes are complete(in validation)
	triehash := statedb.AccountHashes + statedb.StorageHashes                  // The time spent on tries hashing
	trieUpdate := statedb.AccountUpdates + statedb.StorageUpdates              // The time spent on tries update
	trieRead := statedb.SnapshotAccountReads + statedb.AccountReads            // The time spent on account read
	trieRead += statedb.SnapshotStorageReads + statedb.StorageReads            // The time spent on storage read
	blockExecutionTimer.Inc((ptime - trieRead).Milliseconds())                 // The time spent on EVM processing
	blockValidationTimer.Inc((vtime - (triehash + trieUpdate)).Milliseconds()) // The time spent on block validation
	blockTrieOpsTimer.Inc((triehash + trieUpdate + trieRead).Milliseconds())   // The time spent on trie operations

	// If [writes] are disabled, skip [writeBlockWithState] so that we do not write the block
	// or the state trie to disk.
	// Note: in pruning mode, this prevents us from generating a reference to the state root.
	if !writes {
		return nil
	}

	// Write the block to the chain and get the status.
	// writeBlockWithState (called within writeBlockAndSethead) creates a reference that
	// will be cleaned up in Accept/Reject so we need to ensure an error cannot occur
	// later in verification, since that would cause the referenced root to never be dereferenced.
	wstart := time.Now()
	if err := bc.writeBlockAndSetHead(block, receipts, logs, statedb); err != nil {
		return err
	}
	// Update the metrics touched during block commit
	accountCommitTimer.Inc(statedb.AccountCommits.Milliseconds()) // Account commits are complete, we can mark them
	storageCommitTimer.Inc(statedb.StorageCommits.Milliseconds()) // Storage commits are complete, we can mark them
	triedbCommitTimer.Inc(statedb.TrieDBCommits.Milliseconds())   // Trie database commits are complete, we can mark them
	blockWriteTimer.Inc((time.Since(wstart) - statedb.AccountCommits - statedb.StorageCommits - statedb.SnapshotCommits - statedb.TrieDBCommits).Milliseconds())
	blockInsertTimer.Inc(time.Since(start).Milliseconds())

	log.Debug("Inserted new block", "number", block.Number(), "hash", block.Hash(),
		"parentHash", block.ParentHash(),
		"uncles", len(block.Uncles()), "txs", len(block.Transactions()), "gas", block.GasUsed(),
		"elapsed", common.PrettyDuration(time.Since(start)),
		"root", block.Root(), "baseFeePerGas", block.BaseFee(), "blockGasCost", block.BlockGasCost(),
	)

	processedBlockGasUsedCounter.Inc(int64(block.GasUsed()))
	processedTxsCounter.Inc(int64(block.Transactions().Len()))
	processedLogsCounter.Inc(int64(len(logs)))
	blockInsertCount.Inc(1)
	return nil
}

// newTip returns a boolean indicating if the block should be appended to
// the canonical chain.
func (bc *blockChain) newTip(block *types.Block) bool {
	return block.ParentHash() == bc.CurrentBlock().Hash()
}

// writeBlockAndSetHead persists the block and associated state to the database
// and optimistically updates the canonical chain if [block] extends the current
// canonical chain.
// writeBlockAndSetHead expects to be the last verification step during InsertBlock
// since it creates a reference that will only be cleaned up by Accept/Reject.
func (bc *blockChain) writeBlockAndSetHead(block *types.Block, receipts []*types.Receipt, logs []*types.Log, state *state.StateDB) error {
	if err := bc.writeBlockWithState(block, receipts, state); err != nil {
		return err
	}

	// If [block] represents a new tip of the canonical chain, we optimistically add it before
	// setPreference is called. Otherwise, we consider it a side chain block.
	if bc.newTip(block) {
		bc.writeCanonicalBlockWithLogs(block, logs)
	}

	return nil
}

// writeBlockWithState writes the block and all associated state to the database,
// but it expects the chain mutex to be held.
func (bc *blockChain) writeBlockWithState(block *types.Block, receipts []*types.Receipt, state *state.StateDB) error {
	// Irrelevant of the canonical status, write the block itself to the database.
	//
	// Note all the components of block(hash->number map, header, body, receipts)
	// should be written atomically. BlockBatch is used for containing all components.
	blockBatch := bc.db.NewBatch()
	rawdb.WriteBlock(blockBatch, block)
	rawdb.WriteReceipts(blockBatch, block.Hash(), block.NumberU64(), receipts)
	rawdb.WritePreimages(blockBatch, state.Preimages())
	if err := blockBatch.Write(); err != nil {
		return err
	}

	// Commit all cached state changes into underlying memory database.
	// If snapshots are enabled, call CommitWithSnaps to explicitly create a snapshot
	// diff layer for the block.
	_, err := state.CommitWithBlockHash(block.NumberU64(), bc.config.IsEIP158(block.Number()), block.Hash(), block.ParentHash())
	if err != nil {
		return err
	}
	return nil
}

// writeCanonicalBlockWithLogs writes the new head [block] and emits events
// for the new head block.
func (bc *blockChain) writeCanonicalBlockWithLogs(block *types.Block, logs []*types.Log) {
	bc.writeHeadBlock(block)
	// bc.chainFeed.Send(ChainEvent{Block: block, Hash: block.Hash(), Logs: logs})
	// if len(logs) > 0 {
	// 	bc.logsFeed.Send(logs)
	// }
	bc.chainHeadFeed.Send(core.ChainHeadEvent{Block: block})
}

// writeHeadBlock injects a new head block into the current block chain. This method
// assumes that the block is indeed a true head. It will also reset the head
// header to this very same block if they are older or if they are on a different side chain.
//
// Note, this function assumes that the `mu` mutex is held!
func (bc *blockChain) writeHeadBlock(block *types.Block) {
	// If the block is on a side chain or an unknown one, force other heads onto it too
	// Add the block to the canonical chain number scheme and mark as the head
	batch := bc.db.NewBatch()
	rawdb.WriteCanonicalHash(batch, block.Hash(), block.NumberU64())

	rawdb.WriteHeadBlockHash(batch, block.Hash())
	rawdb.WriteHeadHeaderHash(batch, block.Hash())

	// Flush the whole batch into the disk, exit the node if failed
	if err := batch.Write(); err != nil {
		log.Crit("Failed to update chain indexes and markers", "err", err)
	}
	// Update all in-memory chain markers in the last step
	bc.hc.SetCurrentHeader(block.Header())
}

func (bc *blockChain) Accept(block *types.Block) error {
	bc.chainmu.Lock()
	defer bc.chainmu.Unlock()

	// The parent of [block] must be the last accepted block.
	if bc.lastAccepted.Hash() != block.ParentHash() {
		return fmt.Errorf(
			"expected accepted block to have parent %s:%d but got %s:%d",
			bc.lastAccepted.Hash().Hex(),
			bc.lastAccepted.NumberU64(),
			block.ParentHash().Hex(),
			block.NumberU64()-1,
		)
	}

	// If the canonical hash at the block height does not match the block we are
	// accepting, we need to trigger a reorg.
	canonical := bc.hc.GetCanonicalHash(block.NumberU64())
	if canonical != block.Hash() {
		log.Debug("Accepting block in non-canonical chain", "number", block.Number(), "hash", block.Hash())
		if err := bc.setPreference(block); err != nil {
			return fmt.Errorf("could not set new preferred block %d:%s as preferred: %w", block.Number(), block.Hash(), err)
		}
	}

	bc.lastAccepted = block
	acceptedBlockGasUsedCounter.Inc(int64(block.GasUsed()))
	acceptedTxsCounter.Inc(int64(len(block.Transactions())))
	return bc.accept(block)
}

// accept processes a block that has been verified and updates the snapshot
// and indexes.
func (bc *blockChain) accept(next *types.Block) error {
	start := time.Now()

	if err := bc.state.Commit(next.Root(), false); err != nil {
		return fmt.Errorf("unable to accept trie: %w", err)
	}

	// Update last processed and transaction lookup index
	// if err := bc.writeBlockAcceptedIndices(next); err != nil {
	// 	return fmt.Errorf("unable to write block accepted indices: %w", err)
	// }

	// Ensure [hc.acceptedNumberCache] and [acceptedLogsCache] have latest content
	bc.hc.PutAcceptedHeader(next.NumberU64(), next.Header())
	logs := bc.collectUnflattenedLogs(next, false)

	// Update accepted feeds
	flattenedLogs := types.FlattenLogs(logs)
	bc.chainAcceptedFeed.Send(core.ChainEvent{Block: next, Hash: next.Hash(), Logs: flattenedLogs})
	if len(flattenedLogs) > 0 {
		bc.logsAcceptedFeed.Send(flattenedLogs)
	}

	acceptorWorkTimer.Inc(time.Since(start).Milliseconds())
	acceptorWorkCount.Inc(1)
	acceptedLogsCounter.Inc(int64(len(logs)))
	return nil
}

func (bc *blockChain) Reject(block *types.Block) error {
	bc.chainmu.Lock()
	defer bc.chainmu.Unlock()

	// Remove the block since its data is no longer needed
	batch := bc.db.NewBatch()
	rawdb.DeleteBlock(batch, block.Hash(), block.NumberU64())
	if err := batch.Write(); err != nil {
		return fmt.Errorf("failed to write delete block batch: %w", err)
	}

	return nil
}

// SetPreference attempts to update the head block to be the provided block and
// emits a ChainHeadEvent if successful. This function will handle all reorg
// side effects, if necessary.
//
// Note: This function should ONLY be called on blocks that have already been
// inserted into the chain.
func (bc *blockChain) SetPreference(block *types.Block) error {
	bc.chainmu.Lock()
	defer bc.chainmu.Unlock()

	return bc.setPreference(block)
}

// setPreference attempts to update the head block to be the provided block and
// emits a ChainHeadEvent if successful. This function will handle all reorg
// side effects, if necessary.
func (bc *blockChain) setPreference(block *types.Block) error {
	current := bc.CurrentBlock()

	// Return early if the current block is already the block
	// we are trying to write.
	if current.Hash() == block.Hash() {
		return nil
	}

	log.Debug("Setting preference", "number", block.Number(), "hash", block.Hash())

	if block.ParentHash() != current.Hash() {
		if err := bc.reorg(current, block); err != nil {
			return err
		}
	}
	bc.writeHeadBlock(block)

	bc.chainHeadFeed.Send(core.ChainHeadEvent{Block: block})
	return nil
}

// collectUnflattenedLogs collects the logs that were generated or removed during
// the processing of a block.
func (bc *blockChain) collectUnflattenedLogs(b *types.Block, removed bool) [][]*types.Log {
	var blobGasPrice *big.Int
	excessBlobGas := b.ExcessBlobGas()
	if excessBlobGas != nil {
		blobGasPrice = eip4844.CalcBlobFee(*excessBlobGas)
	}
	receipts := rawdb.ReadRawReceipts(bc.db, b.Hash(), b.NumberU64())
	if err := receipts.DeriveFields(bc.config, b.Hash(), b.NumberU64(), b.Time(), b.BaseFee(), blobGasPrice, b.Transactions()); err != nil {
		log.Error("Failed to derive block receipts fields", "hash", b.Hash(), "number", b.NumberU64(), "err", err)
	}

	// Note: gross but this needs to be initialized here because returning nil will be treated specially as an incorrect
	// error case downstream.
	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		receiptLogs := make([]*types.Log, len(receipt.Logs))
		for i, log := range receipt.Logs {
			if removed {
				log.Removed = true
			}
			receiptLogs[i] = log
		}
		logs[i] = receiptLogs
	}
	return logs
}

func (bc *blockChain) reportBlock(block *types.Block, receipts types.Receipts, err error) {
	reason := &core.BadBlockReason{
		ChainConfig: bc.config,
		Receipts:    receipts,
		Number:      block.NumberU64(),
		Hash:        block.Hash(),
		Error:       err.Error(),
	}

	badBlockCounter.Inc(1)
	log.Debug(reason.String())

	// TODO: remove this log once we have a better way to report bad blocks
	log.Error("Invalid block detected", "number", block.Number(), "hash", block.Hash(), "err", err)
}

// loadLastState loads the last known chain state from the database. This method
// assumes that the chain manager mutex is held.
func (bc *blockChain) loadLastState(lastAcceptedHash common.Hash) error {
	// Initialize genesis state
	if lastAcceptedHash == (common.Hash{}) {
		return bc.loadGenesisState()
	}

	// Restore the last known head block
	head := rawdb.ReadHeadBlockHash(bc.db)
	if head == (common.Hash{}) {
		return errors.New("could not read head block hash")
	}
	// Make sure the entire head block is available
	headBlock := bc.GetBlockByHash(head)
	if headBlock == nil {
		return fmt.Errorf("could not load head block %s", head.Hex())
	}

	// Restore the last known head header
	currentHeader := headBlock.Header()
	if head := rawdb.ReadHeadHeaderHash(bc.db); head != (common.Hash{}) {
		if header := bc.GetHeaderByHash(head); header != nil {
			currentHeader = header
		}
	}
	// Everything seems to be fine, set as the head block
	bc.hc.SetCurrentHeader(currentHeader)

	log.Info("Loaded most recent local header", "number", currentHeader.Number, "hash", currentHeader.Hash(), "age", common.PrettyAge(time.Unix(int64(currentHeader.Time), 0)))
	log.Info("Loaded most recent local full block", "number", headBlock.Number(), "hash", headBlock.Hash(), "age", common.PrettyAge(time.Unix(int64(headBlock.Time()), 0)))

	// Otherwise, set the last accepted block and perform a re-org.
	bc.lastAccepted = bc.GetBlockByHash(lastAcceptedHash)
	if bc.lastAccepted == nil {
		return fmt.Errorf("could not load last accepted block")
	}
	// reprocessState is necessary to ensure that the last accepted state is
	// available. The state may not be available if it was not committed due
	// to an unclean shutdown.
	reprocessBlocks := uint64(128)
	if err := bc.reprocessState(bc.lastAccepted, reprocessBlocks); err != nil {
		return fmt.Errorf("failed to reprocess state for last accepted block: %w", err)
	}

	// This ensures that the head block is updated to the last accepted block on startup
	if err := bc.setPreference(bc.lastAccepted); err != nil {
		return fmt.Errorf("failed to set preference to last accepted block while loading last state: %w", err)
	}
	return nil
}

// reprocessState reprocesses the state up to [block], iterating through its ancestors until
// it reaches a block with a state committed to the database. reprocessState does not use
// snapshots since the disk layer for snapshots will most likely be above the last committed
// state that reprocessing will start from.
func (bc *blockChain) reprocessState(current *types.Block, reexec uint64) error {
	origin := current.NumberU64()

	// If the state is already available, skip re-processing.
	if bc.HasState(current.Root()) {
		log.Info("Skipping state reprocessing", "root", current.Root())
		return nil
	}

	var err error
	for i := 0; i < int(reexec); i++ {
		// TODO: handle canceled context

		if current.NumberU64() == 0 {
			return errors.New("genesis state is missing")
		}
		parent := bc.GetBlock(current.ParentHash(), current.NumberU64()-1)
		if parent == nil {
			return fmt.Errorf("missing block %s:%d", current.ParentHash().Hex(), current.NumberU64()-1)
		}
		current = parent
		_, err = bc.state.OpenTrie(current.Root())
		if err == nil {
			break
		}
	}
	if err != nil {
		switch err.(type) {
		case *trie.MissingNodeError:
			return fmt.Errorf("required historical state unavailable (reexec=%d)", reexec)
		default:
			return err
		}
	}

	// State was available at historical point, regenerate
	var (
		start  = time.Now()
		logged time.Time
	)
	// Note: we add 1 since in each iteration, we attempt to re-execute the next block.
	log.Info("Re-executing blocks to generate state for last accepted block", "from", current.NumberU64()+1, "to", origin)
	for current.NumberU64() < origin {
		// TODO: handle canceled context

		// Print progress logs if long enough time elapsed
		if time.Since(logged) > 8*time.Second {
			log.Info("Regenerating historical state", "block", current.NumberU64()+1, "target", origin, "remaining", origin-current.NumberU64(), "elapsed", time.Since(start))
			logged = time.Now()
		}

		// Retrieve the next block to regenerate and process it
		parent := current
		next := current.NumberU64() + 1
		if current = bc.GetBlockByNumber(next); current == nil {
			return fmt.Errorf("failed to retrieve block %d while re-generating state", next)
		}

		// Reprocess next block using previously fetched data
		_, err := bc.reprocessBlock(parent, current)
		if err != nil {
			return err
		}
	}

	log.Info("Historical state regenerated", "block", current.NumberU64(), "elapsed", time.Since(start))
	bc.writeHeadBlock(current)
	return nil
}

// reprocessBlock reprocesses a previously accepted block. This is often used
// to regenerate previously pruned state tries.
func (bc *blockChain) reprocessBlock(parent *types.Block, current *types.Block) (common.Hash, error) {
	// Retrieve the parent block and its state to execute block
	var (
		statedb    *state.StateDB
		err        error
		parentRoot = parent.Root()
	)
	statedb, err = state.New(parentRoot, bc.state, nil)
	if err != nil {
		return common.Hash{}, fmt.Errorf("could not fetch state for (%s: %d): %v", parent.Hash().Hex(), parent.NumberU64(), err)
	}

	// Enable prefetching to pull in trie node paths while processing transactions
	statedb.StartPrefetcher("chain", bc.cacheConfig.TriePrefetcherParallelism)
	defer func() {
		statedb.StopPrefetcher()
	}()

	// Process previously stored block
	receipts, _, usedGas, err := bc.processor.Process(current, parent.Header(), statedb, vm.Config{})
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to re-process block (%s: %d): %v", current.Hash().Hex(), current.NumberU64(), err)
	}

	// Validate the state using the default validator
	if err := bc.validator.ValidateState(current, statedb, receipts, usedGas); err != nil {
		return common.Hash{}, fmt.Errorf("failed to validate state while re-processing block (%s: %d): %v", current.Hash().Hex(), current.NumberU64(), err)
	}
	log.Debug("Processed block", "block", current.Hash(), "number", current.NumberU64())

	// Commit all cached state changes into underlying memory database.
	return statedb.CommitWithBlockHash(current.NumberU64(), bc.config.IsEIP158(current.Number()), current.Hash(), current.ParentHash())
}

func (bc *blockChain) loadGenesisState() error {
	// Prepare the genesis block and reinitialise the chain
	batch := bc.db.NewBatch()
	rawdb.WriteBlock(batch, bc.genesisBlock)
	if err := batch.Write(); err != nil {
		log.Crit("Failed to write genesis block", "err", err)
	}
	bc.writeHeadBlock(bc.genesisBlock)

	// Last update all in-memory chain markers
	bc.lastAccepted = bc.genesisBlock
	bc.hc.SetGenesis(bc.genesisBlock.Header())
	bc.hc.SetCurrentHeader(bc.genesisBlock.Header())
	return nil
}

// Getters
func (bc *blockChain) LastAcceptedBlock() *types.Block {
	bc.chainmu.RLock()
	defer bc.chainmu.RUnlock()

	return bc.lastAccepted
}

func (bc *blockChain) LastConsensusAcceptedBlock() *types.Block { return bc.LastAcceptedBlock() }
func (bc *blockChain) CurrentBlock() *types.Header              { return bc.hc.CurrentHeader() }
func (bc *blockChain) CurrentHeader() *types.Header             { return bc.hc.CurrentHeader() }
func (bc *blockChain) SenderCacher() *core.TxSenderCacher       { return bc.senderCacher }
func (bc *blockChain) CacheConfig() *core.CacheConfig           { return bc.cacheConfig }
func (bc *blockChain) Config() *params.ChainConfig              { return bc.config }
func (bc *blockChain) GetVMConfig() *vm.Config                  { return &bc.vmConfig }
func (bc *blockChain) Engine() consensus.Engine                 { return bc.engine }

// No-ops
func (bc *blockChain) DrainAcceptorQueue()           {}
func (bc *blockChain) InitializeSnapshots()          {}
func (bc *blockChain) ValidateCanonicalChain() error { return nil }
