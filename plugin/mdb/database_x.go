// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package mdb

import (
	"context"
	"fmt"
	"sync"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/trace"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/avalanchego/x/merkledb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethdb"
	"github.com/ava-labs/subnet-evm/trie"
	"github.com/ava-labs/subnet-evm/trie/trienode"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/prometheus/client_golang/prometheus"
)

const merkleDBScheme = "merkleDBScheme"

var (
	_ ethdb.Database   = &WithMerkleDB{}
	_ state.TrieOpener = &WithMerkleDB{}
	_ trie.Backender   = &WithMerkleDB{}
)

type commit struct {
	stack  []merkledb.TrieView
	parent common.Hash
}

type WithMerkleDB struct {
	ethdb.Database
	merkleDB merkledb.MerkleDB

	lock           sync.RWMutex
	pendingCommits map[common.Hash][]commit
}

type backend WithMerkleDB

func toHash(id ids.ID) common.Hash {
	return common.BytesToHash(id[:])
}

func NewWithMerkleDB(db ethdb.Database, merkleDB merkledb.MerkleDB) *WithMerkleDB {
	return &WithMerkleDB{
		Database:       db,
		merkleDB:       merkleDB,
		pendingCommits: make(map[common.Hash][]commit),
	}
}

func (db *WithMerkleDB) Backend() trie.Backend {
	return (*backend)(db)
}

func (db *WithMerkleDB) getParent(root common.Hash) (merkledb.Trie, error) {
	pending, ok := db.pendingCommits[root]
	if ok {
		return pending[0].stack[0], nil
	}
	ctx := context.TODO()
	id, err := db.merkleDB.GetAltMerkleRoot(ctx)
	if err != nil {
		return nil, err
	}
	hash := toHash(id)
	if hash == root {
		return db.merkleDB, nil
	}
	return nil, fmt.Errorf("unknown root %x", root)
}

func (db *WithMerkleDB) OpenTrie(root common.Hash) (trie.ITrie, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	parent, err := db.getParent(root)
	if err != nil {
		return nil, err
	}
	tr := &merkleDBTrie{
		parent:    parent,
		stateRoot: root,
		db:        db,
	}
	tr.initialize()
	return tr, nil
}

func (db *WithMerkleDB) OpenStorageTrie(stateRoot common.Hash, addrHash, prev common.Hash) (trie.ITrie, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	parent, err := db.getParent(stateRoot)
	if err != nil {
		return nil, err
	}
	tr := &merkleDBTrie{
		parent:    parent,
		stateRoot: stateRoot,
		owner:     addrHash,
		db:        db,
	}
	tr.initialize()
	return tr, nil
}

// Initialized returns an indicator if the state data is already initialized
// according to the state scheme.
func (db *backend) Initialized(genesisRoot common.Hash) bool {
	rootID, err := db.merkleDB.GetAltMerkleRoot(context.Background())
	if err != nil {
		panic(err)
	}
	return toHash(rootID) != types.EmptyRootHash
}

// Size returns the current storage size of the memory cache in front of the
// persistent database layer.
func (db *backend) Size() common.StorageSize {
	return 0
	// panic("implement me")
}

// Update performs a state transition by committing dirty nodes contained
// in the given set in order to update state from the specified parent to
// the specified root.
func (db *backend) Update(root common.Hash, parent common.Hash, nodes *trienode.MergedNodeSet) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	t := nodes.Sets[common.Hash{}].Commit.(*merkleDBTrie)
	tvsToCommit := make([]merkledb.TrieView, 0, len(nodes.Sets))
	log.Info("Update", "root", root, "parent", parent, "numTries", len(nodes.Sets))

	for t != nil {
		tvsToCommit = append(tvsToCommit, t.tv)
		t = t.hashParent
	}
	if len(tvsToCommit) != len(nodes.Sets) {
		return fmt.Errorf("TrieView chain length does not match tries to commit (%d != %d)", len(tvsToCommit), len(nodes.Sets))
	}

	db.pendingCommits[root] = append(db.pendingCommits[root], commit{
		stack:  tvsToCommit,
		parent: parent,
	})
	return nil
}

// Commit writes all relevant trie nodes belonging to the specified state
// to disk. Report specifies whether logs will be displayed in info level.
func (db *backend) Commit(root common.Hash, report bool) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	log.Info("Commit", "root", root)

	ctx := context.TODO()
	dbRootID, err := db.merkleDB.GetAltMerkleRoot(ctx)
	if err != nil {
		return err
	}
	dbRoot := toHash(dbRootID)
	success, err := db.commit(ctx, root, dbRoot)
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("failed to commit root %x", root)
	}
	return nil
}

func (db *backend) commit(ctx context.Context, root common.Hash, dbRoot common.Hash) (bool, error) {
	if root == dbRoot {
		return true, nil
	}

	for _, commit := range db.pendingCommits[root] {
		// try committing through this path
		success, err := db.commit(ctx, commit.parent, dbRoot)
		if err != nil {
			return false, err
		}
		if !success {
			continue
		}

		log.Info("commit <--", "root", root, "parent", commit.parent)
		for i := len(commit.stack) - 1; i >= 0; i-- {
			tv := commit.stack[i]
			if err := tv.CommitToDB(ctx); err != nil {
				return false, err
			}
		}
		delete(db.pendingCommits, root)
		return true, nil
	}
	return false, nil
}

// Scheme returns the identifier of used storage scheme.
func (db *backend) Scheme() string {
	return merkleDBScheme
}

func (db *backend) UpdateAndReferenceRoot(root common.Hash, parent common.Hash, nodes *trienode.MergedNodeSet) error {
	return db.Update(root, parent, nodes)
}

// Close closes the trie database backend and releases all held resources.
func (db *backend) Close() error {
	return db.merkleDB.Close()
}

func (db *backend) Dereference(root common.Hash) {
}

func (db *backend) Cap(limit common.StorageSize) error {
	panic("implement me")
}

func (db *backend) Reference(root common.Hash, parent common.Hash) {}

func NewBasicConfig() merkledb.Config {
	return merkledb.Config{
		EvictionBatchSize:         10,
		HistoryLength:             300,
		ValueNodeCacheSize:        units.MiB,
		IntermediateNodeCacheSize: units.MiB,
		Reg:                       prometheus.NewRegistry(),
		Tracer:                    trace.Noop,
		BranchFactor:              merkledb.BranchFactor16,
	}
}
