// (c) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package state

import (
	"crypto/rand"
	"encoding/binary"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/triedb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// Write a test to add 100m kvs to a leveldb so that we can test the prefetcher
// performance.

func BenchmarkPrefetcherDatabase(b *testing.B) {
	require := require.New(b)

	dir := b.TempDir()
	if env := os.Getenv("TEST_DB_DIR"); env != "" {
		dir = env
	}
	wantKVs := 100_000
	if env := os.Getenv("TEST_DB_KVS"); env != "" {
		var err error
		wantKVs, err = strconv.Atoi(env)
		require.NoError(err)
	}

	levelDB, err := rawdb.NewLevelDBDatabase(path.Join(dir, "level.db"), 0, 0, "", false)
	require.NoError(err)
	defer levelDB.Close()

	root := types.EmptyRootHash
	count := uint64(0)
	block := uint64(0)

	rootKey := []byte("root")
	countKey := []byte("count")
	blockKey := []byte("block")
	got, err := levelDB.Get(rootKey)
	if err == nil {
		root = common.BytesToHash(got)
	}
	got, err = levelDB.Get(countKey)
	if err == nil {
		count = binary.BigEndian.Uint64(got)
	}
	got, err = levelDB.Get(blockKey)
	if err == nil {
		block = binary.BigEndian.Uint64(got)
	}

	// Make a trie on the levelDB
	address := common.Address{42}
	addBlock := func(db Database, kvsPerBlock int) {
		root, err = addKVs(db, address, root, block, kvsPerBlock)
		require.NoError(err)
		count += uint64(kvsPerBlock)
		block++
		b.Logf("Root: %v, kvs: %d", root, count)

		// update the tracking keys
		err = levelDB.Put(rootKey, root.Bytes())
		require.NoError(err)
		err = database.PutUInt64(levelDB, blockKey, block)
		require.NoError(err)
		err = database.PutUInt64(levelDB, countKey, count)
		require.NoError(err)
	}

	db := NewDatabaseWithConfig(levelDB, triedb.HashDefaults)
	for count < uint64(wantKVs) {
		addBlock(db, 100_000)
	}

	b.ResetTimer()
	db = NewDatabaseWithConfig(levelDB, triedb.HashDefaults)
	for i := 0; i < b.N; i++ {
		addBlock(db, 1_000)
	}
}

func addKVs(db Database, address common.Address, root common.Hash, block uint64, count int) (common.Hash, error) {
	statedb, err := New(root, db, nil)
	if err != nil {
		return common.Hash{}, err
	}
	statedb.SetNonce(address, 1)
	for i := 0; i < count; i++ {
		key := make([]byte, 32)
		value := make([]byte, 32)
		rand.Read(key)
		rand.Read(value)

		statedb.SetState(address, common.BytesToHash(key), common.BytesToHash(value))
	}
	root, err = statedb.Commit(block+1, true, false)
	if err != nil {
		return common.Hash{}, err
	}
	err = statedb.db.TrieDB().Commit(root, false)
	return root, err
}
