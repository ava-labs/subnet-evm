package chain

import (
	"github.com/ava-labs/coreth/core/types"
	"github.com/ethereum/go-ethereum/trie"
)

type TrieWriter interface {
	InsertTrie(block *types.Block) error // Handle inserted trie reference of [root]
	AcceptTrie(block *types.Block) error // Mark [root] as part of an accepted block
	RejectTrie(block *types.Block) error // Notify TrieWriter that the block containing [root] has been rejected
	Shutdown() error
}

type stateManager struct {
	tdb *trie.Database
}

func NewStateManager(tdb *trie.Database) *stateManager {
	return &stateManager{
		tdb: tdb,
	}
}

func (sm *stateManager) AcceptTrie(block *types.Block) error {
	return sm.tdb.Commit(block.Root(), false)
}

func (sm *stateManager) InsertTrie(block *types.Block) error { return nil }
func (sm *stateManager) RejectTrie(block *types.Block) error { return nil }
func (sm *stateManager) Shutdown() error                     { return nil }
