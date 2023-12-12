// (c) 2020-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package mdb

import (
	"errors"
	"fmt"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/database/prefixdb"
	"github.com/ava-labs/avalanchego/utils/maybe"
	"github.com/ava-labs/avalanchego/x/archivedb"
	"github.com/ava-labs/subnet-evm/ethdb"
	"github.com/ava-labs/subnet-evm/trie"
	"github.com/ava-labs/subnet-evm/trie/trienode"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

var (
	_              trie.ITrie = &archiveTrie{}
	metadataPrefix            = []byte("metadata")
)

type MapOps map[string]maybe.Maybe[[]byte]

type ArchiveDB interface {
	NewBatch(height uint64) database.Batch
	OpenTrie(stateRoot common.Hash, owner common.Hash, root common.Hash) (trie.ITrie, error)
	Commit(root common.Hash, changes []MapOps) error
}

type archiveDB struct {
	archive *archivedb.Database
	db      database.Database
	meta    database.Database
}

func NewArchiveDB(db database.Database) ArchiveDB {
	return &archiveDB{
		archive: archivedb.New(db),
		meta:    prefixdb.New(metadataPrefix, db),
		db:      db,
	}
}

func (db *archiveDB) NewBatch(height uint64) database.Batch {
	return db.archive.NewBatch(height)
}

func (db *archiveDB) Commit(root common.Hash, changes []MapOps) error {
	last, err := db.archive.Height()
	if errors.Is(err, database.ErrNotFound) {
		last = 0 // first commit
	} else if err != nil {
		return err
	}
	b := db.NewBatch(last + 1)
	for _, change := range changes {
		for k, v := range change {
			if v.HasValue() {
				b.Put([]byte(k), v.Value())
			} else {
				b.Delete([]byte(k))
			}
		}
	}
	if err := b.Write(); err != nil {
		return err
	}
	return database.PutUInt64(db.meta, root[:], last+1)
}

func (db *archiveDB) OpenTrie(stateRoot common.Hash, owner common.Hash, root common.Hash) (trie.ITrie, error) {
	height, err := database.GetUInt64(db.meta, stateRoot[:])
	if errors.Is(err, database.ErrNotFound) {
		return nil, fmt.Errorf("%w %x", ErrUnknownRoot, stateRoot)
	} else if err != nil {
		return nil, err
	}
	return &archiveTrie{
		reader: db.archive.Open(height),
		owner:  owner,
	}, nil
}

type archiveTrie struct {
	reader database.KeyValueReader
	owner  common.Hash
}

func (a *archiveTrie) prefixBytes(key []byte) []byte {
	if a.owner == (common.Hash{}) {
		return key
	}
	return append(append(a.owner[:], []byte{0}...), key...)
}

func (a *archiveTrie) Get(k []byte) ([]byte, error) {
	key := a.prefixBytes(k)
	return a.reader.Get(key)
}

func (a *archiveTrie) MustGet(key []byte) []byte {
	got, err := a.Get(key)
	if err != nil {
		log.Error("Unhandled trie error in archiveTrie.Get", "err", err)
	}
	return got
}

func (a *archiveTrie) ICopy() trie.ITrie {
	return &archiveTrie{
		reader: a.reader,
		owner:  a.owner,
	}
}

func (a *archiveTrie) Hash() common.Hash { panic("implement me") }
func (a *archiveTrie) Commit(collectLeaf bool) (common.Hash, *trienode.NodeSet) {
	panic("implement me")
}
func (a *archiveTrie) Delete(key []byte) error                     { panic("implement me") }
func (a *archiveTrie) MustDelete(key []byte)                       { panic("implement me") }
func (a *archiveTrie) Update(key, value []byte) error              { panic("implement me") }
func (a *archiveTrie) MustUpdate(key, value []byte)                { panic("implement me") }
func (a *archiveTrie) NodeIterator(start []byte) trie.NodeIterator { panic("implement me") }
func (a *archiveTrie) Prove(key []byte, fromLevel uint, proofDb ethdb.KeyValueWriter) error {
	panic("implement me")
}
