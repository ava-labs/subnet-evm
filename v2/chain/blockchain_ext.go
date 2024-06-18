package chain

import (
	"fmt"
	"time"

	"github.com/ava-labs/coreth/core/rawdb"
	"github.com/ava-labs/coreth/core/state"
	"github.com/ava-labs/coreth/core/types"
	"github.com/ethereum/go-ethereum/log"
)

// warmAcceptedCaches fetches previously accepted headers and logs from disk to
// pre-populate [hc.acceptedNumberCache] and [acceptedLogsCache].
func (bc *blockChain) warmAcceptedCaches() {
	var (
		startTime       = time.Now()
		lastAccepted    = bc.LastAcceptedBlock().NumberU64()
		startIndex      = uint64(1)
		targetCacheSize = uint64(bc.cacheConfig.AcceptedCacheSize)
	)
	if targetCacheSize == 0 {
		log.Info("Not warming accepted cache because disabled")
		return
	}
	if lastAccepted < startIndex {
		// This could occur if we haven't accepted any blocks yet
		log.Info("Not warming accepted cache because there are no accepted blocks")
		return
	}
	cacheDiff := targetCacheSize - 1 // last accepted lookback is inclusive, so we reduce size by 1
	if cacheDiff < lastAccepted {
		startIndex = lastAccepted - cacheDiff
	}
	for i := startIndex; i <= lastAccepted; i++ {
		block := bc.GetBlockByNumber(i)
		if block == nil {
			// This could happen if a node state-synced
			log.Info("Exiting accepted cache warming early because header is nil", "height", i, "t", time.Since(startTime))
			break
		}
		// TODO: handle blocks written to disk during state sync
		bc.hc.PutAcceptedHeader(block.NumberU64(), block.Header())
	}
	log.Info("Warmed accepted caches", "start", startIndex, "end", lastAccepted, "t", time.Since(startTime))
}

// ResetToStateSyncedBlock reinitializes the state of the blockchain
// to the trie represented by [block.Root()] after updating
// in-memory and on disk current block pointers to [block].
// Only should be called after state sync has completed.
func (bc *blockChain) ResetToStateSyncedBlock(block *types.Block) error {
	bc.chainmu.Lock()
	defer bc.chainmu.Unlock()

	// Update head block and snapshot pointers on disk
	batch := bc.db.NewBatch()
	// if err := bc.batchBlockAcceptedIndices(batch, block); err != nil {
	// 	return err
	// }
	rawdb.WriteHeadBlockHash(batch, block.Hash())
	rawdb.WriteHeadHeaderHash(batch, block.Hash())
	if err := rawdb.WriteSyncPerformed(batch, block.NumberU64()); err != nil {
		return err
	}

	if err := batch.Write(); err != nil {
		return err
	}

	// Update all in-memory chain markers
	bc.lastAccepted = block
	bc.hc.SetCurrentHeader(block.Header())

	lastAcceptedHash := block.Hash()
	bc.state = state.NewDatabaseWithNodeDB(bc.db, bc.triedb)

	if err := bc.loadLastState(lastAcceptedHash); err != nil {
		return err
	}
	// Create the state manager
	bc.stateManager = NewStateManager(bc.triedb)

	// Make sure the state associated with the block is available
	head := bc.CurrentBlock()
	if !bc.HasState(head.Root) {
		return fmt.Errorf("head state missing %d:%s", head.Number, head.Hash())
	}
	return nil
}
