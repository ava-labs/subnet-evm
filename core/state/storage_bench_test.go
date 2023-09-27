package state

import (
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/utils"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/trie"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestGrowingStorageTrie(t *testing.T) {
	require := require.New(t)
	/// XXX: don't hardcode this
	dir := "/Users/darioush.jalali/tmpfoo"
	db, err := rawdb.NewLevelDBDatabase(dir, 128, 128, "", false)
	require.NoError(err)
	defer db.Close()

	testConfig := &trie.Config{
		Cache: 256,
	}
	triedb := trie.NewDatabaseWithConfig(db, testConfig)
	statedb := NewDatabaseWithNodeDB(db, triedb)

	addr := common.Address{2, 3, 4}
	numKeys := 60_000_000
	numKeysPerRoot := 512
	stateRoot := common.Hash{}
	commitInterval := 4096

	sampleTime := time.Duration(0)
	numRootsPerSample := 250

	for rootNum := 1; rootNum*numKeysPerRoot < numKeys; rootNum++ {
		state, err := New(stateRoot, statedb, nil)
		require.NoError(err)

		// In the beginning, we need to update the nonce, balance, or code.
		// Othewise committing will mark the object as deleted.
		// Alterntively, we can try passing false to delete empty objects
		// for this test.

		start := time.Now()
		for i := 0; i < numKeysPerRoot; i++ {
			k := common.BytesToHash(utils.RandomBytes(32))
			v := common.BytesToHash(utils.RandomBytes(32))
			state.SetState(addr, k, v)
		}
		stateRoot, err = state.Commit(false, false)
		require.NoError(err)

		// Ensure the state is not empty
		require.NotZero(stateRoot)
		require.NotEqual(types.EmptyRootHash, stateRoot)

		took := time.Since(start)
		sampleTime += took
		if rootNum%numRootsPerSample == 0 {
			avg := sampleTime / time.Duration(numRootsPerSample)
			t.Logf(
				"%d avg: %v (total keys: %dK)",
				rootNum, avg, (rootNum*numKeysPerRoot)/1000,
			)
			sampleTime = 0
		}

		if rootNum%commitInterval == 0 {
			start := time.Now()
			err := triedb.Commit(stateRoot, false)
			require.NoError(err)
			took := time.Since(start)
			t.Logf(
				"%d committed, took: %v (root %v)",
				rootNum, took, stateRoot,
			)
		}
	}
}
