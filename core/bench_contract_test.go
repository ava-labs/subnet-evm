// (c) 2020-2021, Ava Labs, Inc.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2014 The go-ethereum Authors
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
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	_ "embed"
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/consensus/dummy"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/core/vm"
	"github.com/ava-labs/subnet-evm/ethdb"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
)

var (
	//go:embed TrieStressTest.bin
	stressBinStr string
	//go:embed TrieStressTest.abi
	stressABIStr string
)

func BenchmarkTrie(t *testing.B) {
	benchTrieInserts(t, "test", 5000, false)
}

func generateTx(elements int64) func(int, *BlockGen) {
	return func(i int, gen *BlockGen) {
		gasPrice := big.NewInt(225000000000)
		tx := types.NewContractCreation(gen.TxNonce(benchRootAddr), big.NewInt(0), 3000000, gasPrice, common.FromHex(stressBinStr))
		tx, _ = types.SignTx(tx, signer, testKey)
		addr := gen.AddTx(tx).ContractAddress

		stressABI := contract.ParseABI(stressABIStr)
		txPayload, _ := stressABI.Pack(
			"writeValues",
			big.NewInt(elements),
		)
		tx = types.NewTransaction(gen.TxNonce(benchRootAddr), addr, big.NewInt(0), 3000000, gasPrice, txPayload)
		tx, _ = types.SignTx(tx, signer, testKey)
		gen.AddTx(tx)
	}
}

func benchTrieInserts(b *testing.B, name string, elements int64, disk bool) {
	// Create the database in memory or in a temporary directory.
	var db ethdb.Database
	var err error
	if !disk {
		db = rawdb.NewMemoryDatabase()
	} else {
		dir := b.TempDir()
		db, err = rawdb.NewLevelDBDatabase(dir, 128, 128, "", false)
		if err != nil {
			b.Fatalf("cannot create temporary database: %v", err)
		}
		defer db.Close()
	}

	gspec := &Genesis{
		Config: params.TestChainConfig,
		Alloc:  GenesisAlloc{benchRootAddr: {Balance: benchRootFunds}},
	}

	_, chain, _, _ := GenerateChainWithGenesis(gspec, dummy.NewCoinbaseFaker(), b.N, 10, generateTx(elements))

	// Time the insertion of the new chain.
	// State and blocks are stored in the same DB.
	chainman, _ := NewBlockChain(db, DefaultCacheConfig, gspec, dummy.NewCoinbaseFaker(), vm.Config{}, common.Hash{}, false)
	defer chainman.Stop()

	if i, err := chainman.InsertChain(chain); err != nil {
		b.Fatalf("insert error (block %d): %v\n", i, err)
	}
}
