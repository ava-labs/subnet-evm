// Copyright 2024 The go-ethereum Authors
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
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>

package core

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/dummy"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

// TODO: simplify the unindexer logic and this test.
// XXX: These tests are moved from blockchain_test.go here.
// Should we try to use the TestTxIndexer from upstream here instead
// or move this test to a new file eg, blockchain_extra_test.go?
func TestTransactionIndices(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		key2, _ = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		addr1   = crypto.PubkeyToAddress(key1.PublicKey)
		addr2   = crypto.PubkeyToAddress(key2.PublicKey)
		funds   = big.NewInt(10000000000000)
		gspec   = &Genesis{
			Config: &params.ChainConfig{HomesteadBlock: new(big.Int)},
			Alloc:  types.GenesisAlloc{addr1: {Balance: funds}},
		}
		signer = types.LatestSigner(gspec.Config)
	)
	genDb, blocks, _, err := GenerateChainWithGenesis(gspec, dummy.NewFaker(), 128, 10, func(i int, block *BlockGen) {
		tx, err := types.SignTx(types.NewTransaction(block.TxNonce(addr1), addr2, big.NewInt(10000), params.TxGas, nil, nil), signer, key1)
		require.NoError(t, err)
		block.AddTx(tx)
	})
	require.NoError(t, err)

	blocks2, _, err := GenerateChain(gspec.Config, blocks[len(blocks)-1], dummy.NewFaker(), genDb, 10, 10, func(i int, block *BlockGen) {
		tx, err := types.SignTx(types.NewTransaction(block.TxNonce(addr1), addr2, big.NewInt(10000), params.TxGas, nil, nil), signer, key1)
		require.NoError(t, err)
		block.AddTx(tx)
	})
	require.NoError(t, err)

	check := func(t *testing.T, tail *uint64, chain *BlockChain) {
		require := require.New(t)
		stored := rawdb.ReadTxIndexTail(chain.db)
		var tailValue uint64
		if tail == nil {
			require.Nil(stored)
			tailValue = 0
		} else {
			require.EqualValues(*tail, *stored, "expected tail %d, got %d", *tail, *stored)
			tailValue = *tail
		}

		for i := tailValue; i <= chain.CurrentBlock().Number.Uint64(); i++ {
			block := rawdb.ReadBlock(chain.db, rawdb.ReadCanonicalHash(chain.db, i), i)
			if block.Transactions().Len() == 0 {
				continue
			}
			for _, tx := range block.Transactions() {
				index := rawdb.ReadTxLookupEntry(chain.db, tx.Hash())
				require.NotNilf(index, "Miss transaction indices, number %d hash %s", i, tx.Hash().Hex())
			}
		}

		for i := uint64(0); i < tailValue; i++ {
			block := rawdb.ReadBlock(chain.db, rawdb.ReadCanonicalHash(chain.db, i), i)
			if block.Transactions().Len() == 0 {
				continue
			}
			for _, tx := range block.Transactions() {
				index := rawdb.ReadTxLookupEntry(chain.db, tx.Hash())
				require.Nilf(index, "Transaction indices should be deleted, number %d hash %s", i, tx.Hash().Hex())
			}
		}
	}

	conf := &CacheConfig{
		TrieCleanLimit:            256,
		TrieDirtyLimit:            256,
		TrieDirtyCommitTarget:     20,
		TriePrefetcherParallelism: 4,
		Pruning:                   true,
		CommitInterval:            4096,
		SnapshotLimit:             256,
		SnapshotNoBuild:           true, // Ensure the test errors if snapshot initialization fails
		AcceptorQueueLimit:        64,
	}

	// Init block chain and check all needed indices has been indexed.
	chainDB := rawdb.NewMemoryDatabase()
	chain, err := createBlockChain(chainDB, conf, gspec, common.Hash{})
	require.NoError(t, err)

	_, err = chain.InsertChain(blocks)
	require.NoError(t, err)

	for _, block := range blocks {
		err := chain.Accept(block)
		require.NoError(t, err)
	}
	chain.DrainAcceptorQueue()

	chain.Stop()
	check(t, nil, chain) // check all indices has been indexed

	lastAcceptedHash := chain.CurrentHeader().Hash()

	// Reconstruct a block chain which only reserves limited tx indices
	// 128 blocks were previously indexed. Now we add a new block at each test step.
	limits := []uint64{
		0,   /* tip: 129 reserve all (don't run) */
		131, /* tip: 130 reserve all */
		140, /* tip: 131 reserve all */
		64,  /* tip: 132, limit:64 */
		32,  /* tip: 133, limit:32  */
	}
	for i, l := range limits {
		t.Run(fmt.Sprintf("test-%d, limit: %d", i+1, l), func(t *testing.T) {
			conf.TxLookupLimit = l

			chain, err := createBlockChain(chainDB, conf, gspec, lastAcceptedHash)
			require.NoError(t, err)

			newBlks := blocks2[i : i+1]
			_, err = chain.InsertChain(newBlks) // Feed chain a higher block to trigger indices updater.
			require.NoError(t, err)

			err = chain.Accept(newBlks[0]) // Accept the block to trigger indices updater.
			require.NoError(t, err)

			chain.DrainAcceptorQueue()
			time.Sleep(50 * time.Millisecond) // Wait for indices initialisation

			chain.Stop()
			var tail *uint64
			if l == 0 {
				tail = nil
			} else {
				var tl uint64
				if chain.CurrentBlock().Number.Uint64() > l {
					// tail should be the first block number which is indexed
					// i.e the first block number that's in the lookup range
					tl = chain.CurrentBlock().Number.Uint64() - l + 1
				}
				tail = &tl
			}

			check(t, tail, chain)

			lastAcceptedHash = chain.CurrentHeader().Hash()
		})
	}
}

func TestTransactionSkipIndexing(t *testing.T) {
	// Configure and generate a sample block chain
	require := require.New(t)
	var (
		key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		key2, _ = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		addr1   = crypto.PubkeyToAddress(key1.PublicKey)
		addr2   = crypto.PubkeyToAddress(key2.PublicKey)
		funds   = big.NewInt(10000000000000)
		gspec   = &Genesis{
			Config: &params.ChainConfig{HomesteadBlock: new(big.Int)},
			Alloc:  types.GenesisAlloc{addr1: {Balance: funds}},
		}
		signer = types.LatestSigner(gspec.Config)
	)
	genDb, blocks, _, err := GenerateChainWithGenesis(gspec, dummy.NewCoinbaseFaker(), 5, 10, func(i int, block *BlockGen) {
		tx, err := types.SignTx(types.NewTransaction(block.TxNonce(addr1), addr2, big.NewInt(10000), params.TxGas, nil, nil), signer, key1)
		require.NoError(err)
		block.AddTx(tx)
	})
	require.NoError(err)

	blocks2, _, err := GenerateChain(gspec.Config, blocks[len(blocks)-1], dummy.NewCoinbaseFaker(), genDb, 5, 10, func(i int, block *BlockGen) {
		tx, err := types.SignTx(types.NewTransaction(block.TxNonce(addr1), addr2, big.NewInt(10000), params.TxGas, nil, nil), signer, key1)
		require.NoError(err)
		block.AddTx(tx)
	})
	require.NoError(err)

	checkRemoved := func(tail *uint64, to uint64, chain *BlockChain) {
		stored := rawdb.ReadTxIndexTail(chain.db)
		var tailValue uint64
		if tail == nil {
			require.Nil(stored)
			tailValue = 0
		} else {
			require.EqualValues(*tail, *stored, "expected tail %d, got %d", *tail, *stored)
			tailValue = *tail
		}

		for i := tailValue; i < to; i++ {
			block := rawdb.ReadBlock(chain.db, rawdb.ReadCanonicalHash(chain.db, i), i)
			if block.Transactions().Len() == 0 {
				continue
			}
			for _, tx := range block.Transactions() {
				index := rawdb.ReadTxLookupEntry(chain.db, tx.Hash())
				require.NotNilf(index, "Miss transaction indices, number %d hash %s", i, tx.Hash().Hex())
			}
		}

		for i := uint64(0); i < tailValue; i++ {
			block := rawdb.ReadBlock(chain.db, rawdb.ReadCanonicalHash(chain.db, i), i)
			if block.Transactions().Len() == 0 {
				continue
			}
			for _, tx := range block.Transactions() {
				index := rawdb.ReadTxLookupEntry(chain.db, tx.Hash())
				require.Nilf(index, "Transaction indices should be deleted, number %d hash %s", i, tx.Hash().Hex())
			}
		}

		for i := to; i <= chain.CurrentBlock().Number.Uint64(); i++ {
			block := rawdb.ReadBlock(chain.db, rawdb.ReadCanonicalHash(chain.db, i), i)
			if block.Transactions().Len() == 0 {
				continue
			}
			for _, tx := range block.Transactions() {
				index := rawdb.ReadTxLookupEntry(chain.db, tx.Hash())
				require.Nilf(index, "Transaction indices should be skipped, number %d hash %s", i, tx.Hash().Hex())
			}
		}
	}

	conf := &CacheConfig{
		TrieCleanLimit:            256,
		TrieDirtyLimit:            256,
		TrieDirtyCommitTarget:     20,
		TriePrefetcherParallelism: 4,
		Pruning:                   true,
		CommitInterval:            4096,
		SnapshotLimit:             256,
		SnapshotNoBuild:           true, // Ensure the test errors if snapshot initialization fails
		AcceptorQueueLimit:        64,
		SkipTxIndexing:            true,
	}

	// test1: Init block chain and check all indices has been skipped.
	chainDB := rawdb.NewMemoryDatabase()
	chain, err := createAndInsertChain(chainDB, conf, gspec, blocks, common.Hash{})
	require.NoError(err)
	checkRemoved(nil, 0, chain) // check all indices has been skipped

	// test2: specify lookuplimit with tx index skipping enabled. Blocks should not be indexed but tail should be updated.
	conf.TxLookupLimit = 2
	chain, err = createAndInsertChain(chainDB, conf, gspec, blocks2[0:1], chain.CurrentHeader().Hash())
	require.NoError(err)
	tail := chain.CurrentBlock().Number.Uint64() - conf.TxLookupLimit + 1
	checkRemoved(&tail, 0, chain)

	// test3: tx index skipping and unindexer disabled. Blocks should be indexed and tail should be updated.
	conf.TxLookupLimit = 0
	conf.SkipTxIndexing = false
	chainDB = rawdb.NewMemoryDatabase()
	chain, err = createAndInsertChain(chainDB, conf, gspec, blocks, common.Hash{})
	require.NoError(err)
	checkRemoved(nil, chain.CurrentBlock().Number.Uint64()+1, chain) // check all indices has been indexed

	// now change tx index skipping to true and check that the indices are skipped for the last block
	// and old indices are removed up to the tail, but [tail, current) indices are still there.
	conf.TxLookupLimit = 2
	conf.SkipTxIndexing = true
	chain, err = createAndInsertChain(chainDB, conf, gspec, blocks2[0:1], chain.CurrentHeader().Hash())
	require.NoError(err)
	tail = chain.CurrentBlock().Number.Uint64() - conf.TxLookupLimit + 1
	checkRemoved(&tail, chain.CurrentBlock().Number.Uint64(), chain)
}

func createAndInsertChain(db ethdb.Database, cacheConfig *CacheConfig, gspec *Genesis, blocks types.Blocks, lastAcceptedHash common.Hash) (*BlockChain, error) {
	chain, err := createBlockChain(db, cacheConfig, gspec, lastAcceptedHash)
	if err != nil {
		return nil, err
	}
	_, err = chain.InsertChain(blocks)
	if err != nil {
		return nil, err
	}
	for _, block := range blocks {
		err := chain.Accept(block)
		if err != nil {
			return nil, err
		}
	}

	chain.DrainAcceptorQueue()
	time.Sleep(1000 * time.Millisecond) // Wait for indices initialisation

	chain.Stop()
	return chain, nil
}
