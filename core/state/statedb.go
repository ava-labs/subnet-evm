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

// Package state provides a caching layer atop the Ethereum state trie.
package state

import (
	"github.com/ava-labs/coreth/utils"
	"github.com/ava-labs/libevm/common"
	ethstate "github.com/ava-labs/libevm/core/state"
)

// StateDB structs within the ethereum protocol are used to store anything
// within the merkle trie. StateDBs take care of caching and storing
// nested states. It's the general query interface to retrieve:
//
// * Contracts
// * Accounts
//
// Once the state is committed, tries cached in stateDB (including account
// trie, storage tries) will no longer be functional. A new state instance
// must be created with new root and updated database for accessing post-
// commit states.
type StateDB struct {
	*ethstate.StateDB

	// The tx context
	thash   common.Hash
	txIndex int

	// Some fields remembered as they are used in tests
	db    Database
	snaps ethstate.SnapshotTree
}

// New creates a new state from a given trie.
func New(root common.Hash, db Database, snaps ethstate.SnapshotTree) (*StateDB, error) {
	stateDB, err := ethstate.New(root, db, snaps)
	if err != nil {
		return nil, err
	}
	return &StateDB{
		StateDB: stateDB,
		db:      db,
		snaps:   snaps,
	}, nil
}

type workerPool struct {
	*utils.BoundedWorkers
}

func (wp *workerPool) Done() {
	// Done is guaranteed to only be called after all work is already complete,
	// so we call Wait for goroutines to finish before returning.
	wp.BoundedWorkers.Wait()
}

func withConcurrentWorkers(prefetchers int) ethstate.PrefetcherOption {
	pool := &workerPool{
		BoundedWorkers: utils.NewBoundedWorkers(prefetchers),
	}
	return ethstate.WithWorkerPools(func() ethstate.WorkerPool { return pool })
}

// StartPrefetcher initializes a new trie prefetcher to pull in nodes from the
// state trie concurrently while the state is mutated so that when we reach the
// commit phase, most of the needed data is already hot.
func (s *StateDB) StartPrefetcher(namespace string, maxConcurrency int) {
	s.StateDB.StartPrefetcher(namespace, withConcurrentWorkers(maxConcurrency))
}

// GetState retrieves a value from the given account's storage trie.
func (s *StateDB) GetState(addr common.Address, hash common.Hash) common.Hash {
	return s.StateDB.GetState(addr, hash)
}

func (s *StateDB) SetState(addr common.Address, key, value common.Hash) {
	s.StateDB.SetState(addr, key, value)
}

// SetTxContext sets the current transaction hash and index which are
// used when the EVM emits new state logs. It should be invoked before
// transaction execution.
func (s *StateDB) SetTxContext(thash common.Hash, ti int) {
	s.thash = thash
	s.txIndex = ti
	s.StateDB.SetTxContext(thash, ti)
}

// GetTxHash returns the current tx hash on the StateDB set by SetTxContext.
func (s *StateDB) GetTxHash() common.Hash {
	return s.thash
}

func (s *StateDB) Copy() *StateDB {
	return &StateDB{
		StateDB: s.StateDB.Copy(),
		db:      s.db,
		snaps:   s.snaps,
		thash:   s.thash,
		txIndex: s.txIndex,
	}
}

// NormalizeCoinID ORs the 0th bit of the first byte in
// [coinID], which ensures this bit will be 1 and all other
// bits are left the same.
// This partitions multicoin storage from normal state storage.
func NormalizeCoinID(coinID *common.Hash) {
	coinID[0] |= 0x01
}

// NormalizeStateKey ANDs the 0th bit of the first byte in
// [key], which ensures this bit will be 0 and all other bits
// are left the same.
// This partitions normal state storage from multicoin storage.
func NormalizeStateKey(key *common.Hash) {
	key[0] &= 0xfe
}
