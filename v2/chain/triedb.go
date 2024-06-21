package chain

import (
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/trie/triedb/pathdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"
)

const commitLag = 32

func NewTrieDB(db ethdb.Database, config *core.CacheConfig) *trie.Database {
	pCfg := &pathdb.Config{
		StateHistory:   config.StateHistory,
		CleanCacheSize: config.TrieCleanLimit * 1024 * 1024,
		DirtyCacheSize: config.TrieDirtyLimit * 1024 * 1024,
		CommitLag:      commitLag,
	}
	return trie.NewDatabase(db, &trie.Config{PathDB: pCfg})
}
