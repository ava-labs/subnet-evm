// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package core

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ava-labs/avalanchego/database/memdb"
	"github.com/ava-labs/avalanchego/x/merkledb"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/plugin/mdb"
	"github.com/ava-labs/subnet-evm/trie"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestCompareStateDB(t *testing.T) {
	require := require.New(t)
	cacheConfig := &trie.Config{
		Cache:       256,
		StatsPrefix: trieCleanCacheStatsNamespace,
	}

	db1 := rawdb.NewMemoryDatabase()
	trieDB1 := trie.NewDatabaseWithConfig(db1, cacheConfig)
	stateDB1 := state.NewDatabaseWithNodeDB(db1, trieDB1)

	// Create a merkleDB
	db := rawdb.NewMemoryDatabase()
	memDB := memdb.New()
	merkleDB, err := merkledb.New(context.Background(), memDB, mdb.NewTestConfig())
	require.NoError(err)
	db2 := mdb.NewWithMerkleDB(db, merkleDB)
	trieDB2 := trie.NewDatabaseWithConfig(db2, cacheConfig)
	stateDB2 := state.NewDatabaseWithNodeDB(db2, trieDB2)

	addr1 := common.Address{0x01}
	addr2 := common.Address{0x02}
	_, _ = addr1, addr2
	ops := []op{
		setStorage{
			addrs:   []common.Address{addr1},
			numKeys: 1,
		},
		setStorage{
			addrs:       []common.Address{addr1},
			numKeys:     1,
			startOffset: 1,
		},
	}

	stateRoot := types.EmptyRootHash
	for i, op := range ops {
		sdb1, err := state.New(stateRoot, stateDB1, nil)
		require.NoError(err)
		sdb2, err := state.New(stateRoot, stateDB2, nil)
		require.NoError(err)

		op.apply(sdb1)
		op.apply(sdb2)

		r1 := sdb1.IntermediateRoot(false)
		r2 := sdb2.IntermediateRoot(false)
		fmt.Println(i)
		require.Equal(r1, r2)

		_, err = sdb1.Commit(true, true)
		require.NoError(err)
		_, err = sdb2.Commit(true, true)
		require.NoError(err)

		err = trieDB1.Commit(r1, false)
		require.NoError(err)
		err = trieDB2.Commit(r2, false)
		require.NoError(err)

		stateRoot = r1
	}

}

type op interface {
	apply(*state.StateDB)
}

type addBalance struct {
	addrs   []common.Address
	balance *big.Int
}

func (op addBalance) apply(sdb *state.StateDB) {
	for _, addr := range op.addrs {
		sdb.AddBalance(addr, op.balance)
	}
}

type setStorage struct {
	addrs       []common.Address
	numKeys     int
	startOffset int
}

func (op setStorage) apply(sdb *state.StateDB) {
	for _, addr := range op.addrs {
		for i := 0; i < op.numKeys; i++ {
			key := common.BigToHash(big.NewInt(int64(i)))
			value := common.BigToHash(big.NewInt(int64(i + op.startOffset)))
			sdb.SetState(addr, key, value)
		}
	}
}
