package core

import (
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/trie"
)

type commitableStateDB interface {
	state.Database
	Commit(root common.Hash, report bool) error
	Initialized(root common.Hash) bool
	Close(lastBlockRoot common.Hash) error
}

type withCommit struct {
	state.Database
}

func (w withCommit) Commit(root common.Hash, report bool) error {
	return w.Database.TrieDB().Commit(root, report)
}

func (w withCommit) Initialized(root common.Hash) bool {
	return w.Database.TrieDB().Initialized(root)
}

func (w withCommit) Close(lastBlockRoot common.Hash) error {
	triedb := w.Database.TrieDB()
	if triedb.Scheme() == rawdb.PathScheme {
		// Ensure that the in-memory trie nodes are journaled to disk properly.
		if err := triedb.Journal(lastBlockRoot); err != nil {
			log.Info("Failed to journal in-memory trie nodes", "err", err)
		}
	}
	// Close the trie database, release all the held resources as the last step.
	if err := triedb.Close(); err != nil {
		log.Error("Failed to close trie database", "err", err)
	}
	return nil
}

func AsCommittable(db ethdb.Database, tdb *trie.Database) commitableStateDB {
	return withCommit{state.NewDatabaseWithNodeDB(db, tdb)}
}
