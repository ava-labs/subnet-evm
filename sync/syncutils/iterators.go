// (c) 2021-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package syncutils

import (
	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/subnet-evm/core/state/snapshot"
	"github.com/ava-labs/subnet-evm/ethdb"
	"github.com/ava-labs/subnet-evm/plugin/mdb"
)

var (
	_ ethdb.Iterator = &AccountIterator{}
	_ ethdb.Iterator = &StorageIterator{}
)

// AccountIterator wraps a [snapshot.AccountIterator] to conform to [ethdb.Iterator]
// accounts will be returned in consensus (FullRLP) format for compatibility with trie data.
type AccountIterator struct {
	snapshot.AccountIterator
	err error
	val []byte
}

func (it *AccountIterator) Next() bool {
	if it.err != nil {
		return false
	}
	for it.AccountIterator.Next() {
		it.val, it.err = snapshot.FullAccountRLP(it.Account())
		return it.err == nil
	}
	it.val = nil
	return false
}

func (it *AccountIterator) Key() []byte {
	if it.err != nil {
		return nil
	}
	return it.Hash().Bytes()
}

func (it *AccountIterator) Value() []byte {
	if it.err != nil {
		return nil
	}
	return it.val
}

func (it *AccountIterator) Error() error {
	if it.err != nil {
		return it.err
	}
	return it.AccountIterator.Error()
}

// StorageIterator wraps a [snapshot.StorageIterator] to conform to [ethdb.Iterator]
type StorageIterator struct {
	snapshot.StorageIterator
}

func (it *StorageIterator) Key() []byte {
	return it.Hash().Bytes()
}

func (it *StorageIterator) Value() []byte {
	return it.Slot()
}

func ClearPartialDB(db ethdb.Database) error {
	if wmdb, ok := db.(*mdb.WithMerkleDB); ok {
		if err := database.Clear(wmdb.MerkleDB(), units.MiB); err != nil {
			return err
		}
	}

	// Wipe the snapshot completely if we are not resuming from an existing sync, so that we do not
	// use a corrupted snapshot.
	// Note: this assumes that when the node is started with state sync disabled, the in-progress state
	// sync marker will be wiped, so we do not accidentally resume progress from an incorrect version
	// of the snapshot. (if switching between versions that come before this change and back this could
	// lead to the snapshot not being cleaned up correctly)
	<-snapshot.WipeSnapshot(db, true)
	// Reset the snapshot generator here so that when state sync completes, snapshots will not attempt to read an
	// invalid generator.
	// Note: this must be called after WipeSnapshot is called so that we do not invalidate a partially generated snapshot.
	snapshot.ResetSnapshotGeneration(db)
	return nil
}
