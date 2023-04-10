// (c) 2021-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package handlers

import (
	"time"

	"github.com/ava-labs/subnet-evm/core/state/snapshot"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethdb"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ BlockProvider    = &TestBlockProvider{}
	_ SnapshotProvider = &TestSnapshotProvider{}
)

type TestBlockProvider struct {
	GetBlockFn func(common.Hash, uint64) *types.Block
}

func (t *TestBlockProvider) GetBlock(hash common.Hash, number uint64) *types.Block {
	return t.GetBlockFn(hash, number)
}

type TestSnapshotProvider struct {
	Snapshot *snapshot.Tree
}

func (t *TestSnapshotProvider) Snapshots() *snapshot.Tree {
	return t.Snapshot
}

type delayedReader struct {
	ethdb.KeyValueStore
	delay time.Duration
}

func (d *delayedReader) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	return &delayedIterator{
		Iterator:      d.KeyValueStore.NewIterator(prefix, start),
		delayedReader: d,
	}
}

type delayedIterator struct {
	ethdb.Iterator
	delayedReader *delayedReader
}

func (d *delayedIterator) Next() bool {
	time.Sleep(d.delayedReader.delay)
	return d.Iterator.Next()
}
