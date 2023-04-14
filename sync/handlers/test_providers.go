// (c) 2021-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package handlers

import (
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

type blockingReader struct {
	ethdb.KeyValueStore
	blockChan <-chan struct{}
}

func (b *blockingReader) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	return &blockingIterator{
		Iterator:       b.KeyValueStore.NewIterator(prefix, start),
		blockingReader: b,
	}
}

type blockingIterator struct {
	ethdb.Iterator
	blockingReader *blockingReader
}

func (b *blockingIterator) Next() bool {
	if wait := b.blockingReader.blockChan; wait != nil {
		<-wait
	}
	return b.Iterator.Next()
}
