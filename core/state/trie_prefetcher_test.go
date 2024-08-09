// (c) 2019-2020, Ava Labs, Inc.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2021 The go-ethereum Authors
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

package state

import (
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/tracing"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

const maxConcurrency = 4

func filledStateDB() *StateDB {
	state, _ := New(types.EmptyRootHash, NewDatabase(rawdb.NewMemoryDatabase()), nil)

	// Create an account and check if the retrieved balance is correct
	addr := common.HexToAddress("0xaffeaffeaffeaffeaffeaffeaffeaffeaffeaffe")
	skey := common.HexToHash("aaa")
	sval := common.HexToHash("bbb")

	state.SetBalance(addr, uint256.NewInt(42), tracing.BalanceChangeUnspecified) // Change the account trie
	state.SetCode(addr, []byte("hello"))                                         // Change an external metadata
	state.SetState(addr, skey, sval)                                             // Change the storage trie
	for i := 0; i < 100; i++ {
		sk := common.BigToHash(big.NewInt(int64(i)))
		state.SetState(addr, sk, sk) // Change the storage trie
	}
	return state
}

func TestUseAfterTerminate(t *testing.T) {
	db := filledStateDB()
	prefetcher := newTriePrefetcher(db.db, db.originalRoot, "", true, 1)
	skey := common.HexToHash("aaa")

	if err := prefetcher.prefetch(common.Hash{}, db.originalRoot, common.Address{}, [][]byte{skey.Bytes()}, false); err != nil {
		t.Errorf("Prefetch failed before terminate: %v", err)
	}
	prefetcher.terminate(false)

	if err := prefetcher.prefetch(common.Hash{}, db.originalRoot, common.Address{}, [][]byte{skey.Bytes()}, false); err == nil {
		t.Errorf("Prefetch succeeded after terminate: %v", err)
	}
	if tr := prefetcher.trie(common.Hash{}, db.originalRoot); tr == nil {
		t.Errorf("Prefetcher returned nil trie after terminate")
	}
}
