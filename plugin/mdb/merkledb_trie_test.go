// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package mdb

import (
	"context"
	"fmt"
	"testing"

	"github.com/ava-labs/avalanchego/database/memdb"
	"github.com/ava-labs/avalanchego/x/merkledb"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/trie"
	"github.com/ava-labs/subnet-evm/trie/trienode"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestMerkleDBTries(t *testing.T) {
	require := require.New(t)
	memDB := memdb.New()
	db, err := merkledb.New(context.Background(), memDB, NewBasicConfig())
	require.NoError(err)

	chaindb := rawdb.NewMemoryDatabase()
	trieDB := NewWithMerkleDB(chaindb, db)
	// Open at empty
	stateRoot := types.EmptyRootHash
	st, err := trieDB.OpenTrie(stateRoot)
	require.NoError(err)

	var (
		prefixes   = []common.Hash{common.Hash{1}, common.Hash{2}, common.Hash{3}}
		storages   []trie.ITrie
		lastHashed trie.ITrie
	)
	numStorageKeys := 10
	for _, prefix := range prefixes {
		// Get some storage tries
		storage, err := trieDB.OpenStorageTrie(stateRoot, prefix, common.Hash{})
		require.NoError(err)

		// Insert some values
		for i := 0; i < numStorageKeys; i++ {
			stKey := common.BytesToHash([]byte(fmt.Sprintf("k-%d", i)))
			stVal := common.BytesToHash([]byte(fmt.Sprintf("v-%d", i)))
			err := storage.Update(stKey[:], stVal[:])
			require.NoError(err)
		}

		// Next, we will hash the storage tries.
		storage.(state.HashChainTrie).SetLastHashed(lastHashed)
		root := storage.Hash()
		require.NotZero(root)
		t.Logf("storage trie hash: %x", root)
		storages = append(storages, storage)
		lastHashed = storage
	}

	// Then, update the state trie with some info from the storage tries.
	for i, storage := range storages {
		err := st.Update(prefixes[i][:], storage.Hash().Bytes())
		require.NoError(err)
	}

	// Next, hash the state trie
	st.(state.HashChainTrie).SetLastHashed(lastHashed)
	nextRoot := st.Hash()
	require.NotZero(nextRoot)
	t.Logf("state root hash: %x", nextRoot)

	// Next, commit the storage tries
	merged := trienode.NewMergedNodeSet()
	for _, storage := range storages {
		_, nodes := storage.Commit(false)
		merged.Merge(nodes)
	}

	// Then, commit the account trie
	_, nodes := st.Commit(false)
	merged.Merge(nodes)

	// Perform an update
	backend := (*backend)(trieDB)
	err = backend.Update(nextRoot, stateRoot, merged)
	require.NoError(err)

	// Let's open the state trie from memory
	reopened, err := trieDB.OpenTrie(nextRoot)
	require.NoError(err)
	for i, storage := range storages {
		got, err := reopened.Get(prefixes[i][:])
		require.NoError(err)
		require.Equal(storage.Hash().Bytes(), got)

		// Check that the storage trie is still intact
		storageTrie, err := trieDB.OpenStorageTrie(nextRoot, prefixes[i], common.Hash{})
		require.NoError(err)
		for i := 0; i < numStorageKeys; i++ {
			stKey := common.BytesToHash([]byte(fmt.Sprintf("k-%d", i)))
			stVal := common.BytesToHash([]byte(fmt.Sprintf("v-%d", i)))
			got, err := storageTrie.Get(stKey[:])
			require.NoError(err)
			require.Equal(stVal[:], got)
		}
	}

	// Write to disk
	err = backend.Commit(nextRoot, false)
	require.NoError(err)
}
