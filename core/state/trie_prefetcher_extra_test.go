// (c) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package state

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state/snapshot"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/metrics"
	"github.com/ava-labs/subnet-evm/triedb"
	"github.com/ava-labs/subnet-evm/triedb/hashdb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/stretchr/testify/require"
)

const prefix = "chain"

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
	address1 := common.Address{42}
	address2 := common.Address{43}
	addBlock := func(db Database, snaps *snapshot.Tree, kvsPerBlock int, prefetchers int) (statedb *StateDB) {
		statedb, root, err = addKVs(db, snaps, address1, address2, root, block, kvsPerBlock, prefetchers)
		require.NoError(err)
		count += uint64(kvsPerBlock)
		block++

		return statedb
	}

	lastCommit := block
	commit := func(levelDB ethdb.Database, snaps *snapshot.Tree, db Database) {
		// update the tracking keys
		err = levelDB.Put(rootKey, root.Bytes())
		require.NoError(err)
		err = database.PutUInt64(levelDB, blockKey, block)
		require.NoError(err)
		err = database.PutUInt64(levelDB, countKey, count)
		require.NoError(err)

		require.NoError(db.TrieDB().Commit(root, false))

		for i := lastCommit + 1; i <= block; i++ {
			require.NoError(snaps.Flatten(fakeHash(i)))
		}
		lastCommit = block
	}

	tdbConfig := &triedb.Config{
		HashDB: &hashdb.Config{
			CleanCacheSize: 3 * 1024 * 1024 * 1024,
		},
	}
	db := NewDatabaseWithConfig(levelDB, tdbConfig)
	snaps := snapshot.NewTestTree(levelDB, fakeHash(block), root)
	for count < uint64(wantKVs) {
		previous := root
		_ = addBlock(db, snaps, 100_000, 0) // Note this updates root and count
		b.Logf("Root: %v, kvs: %d, block: %d", root, count, block)

		// Commit every 10 blocks or on the last iteration
		if block%10 == 0 || count >= uint64(wantKVs) {
			commit(levelDB, snaps, db)
			b.Logf("Root: %v, kvs: %d, block: %d (committed)", root, count, block)
		}
		if previous != root {
			require.NoError(db.TrieDB().Dereference(previous))
		} else {
			panic("root did not change")
		}
	}
	require.NoError(levelDB.Close())
	b.Logf("Starting benchmarks")
	b.Logf("Root: %v, kvs: %d, block: %d", root, count, block)
	for _, updates := range []int{100, 200, 500, 1_000, 10_000, 100_000} {
		for _, prefetchers := range []int{0, 1, 4, 16} {
			b.Run(fmt.Sprintf("updates_%d_prefetchers_%d", updates, prefetchers), func(b *testing.B) {
				startRoot, startBlock, startCount := root, block, count
				defer func() { root, block, count = startRoot, startBlock, startCount }()

				levelDB, err := rawdb.NewLevelDBDatabase(path.Join(dir, "level.db"), 0, 0, "", false)
				require.NoError(err)
				snaps := snapshot.NewTestTree(levelDB, fakeHash(block), root)
				db := NewDatabaseWithConfig(levelDB, tdbConfig)
				storage := int64(0)
				for i := 0; i < b.N; i++ {
					_ = addBlock(db, snaps, updates, prefetchers)
					meter := metrics.GetOrRegisterMeter(prefix+"/storage/skip", nil)
					storage += meter.Snapshot().Count()
				}
				require.NoError(levelDB.Close())
				b.ReportMetric(float64(storage)/float64(b.N), "storage")
			})
		}
	}
}

func fakeHash(block uint64) common.Hash {
	return common.BytesToHash(binary.BigEndian.AppendUint64(nil, block))
}

func addKVs(
	db Database, snaps *snapshot.Tree,
	address1, address2 common.Address, root common.Hash, block uint64,
	count int, prefetchers int,
) (*StateDB, common.Hash, error) {
	snap := snaps.Snapshot(root)
	if snap == nil {
		return nil, common.Hash{}, fmt.Errorf("snapshot not found")
	}
	statedb, err := NewWithSnapshot(root, db, snap)
	if err != nil {
		return nil, common.Hash{}, err
	}
	if prefetchers > 0 {
		statedb.StartPrefetcher(prefix, prefetchers)
		// defer statedb.StopPrefetcher()
	}
	statedb.SetNonce(address1, 1)
	for i := 0; i < count/2; i++ {
		key := make([]byte, 32)
		value := make([]byte, 32)
		rand.Read(key)
		rand.Read(value)

		statedb.SetState(address1, common.BytesToHash(key), common.BytesToHash(value))
	}
	statedb.SetNonce(address2, 1)
	for i := 0; i < count/2; i++ {
		key := make([]byte, 32)
		value := make([]byte, 32)
		rand.Read(key)
		rand.Read(value)

		statedb.SetState(address2, common.BytesToHash(key), common.BytesToHash(value))
	}
	root, err = statedb.CommitWithSnap(block+1, true, snaps, fakeHash(block+1), fakeHash(block), false)
	if err != nil {
		return nil, common.Hash{}, err
	}
	return statedb, root, err
}
