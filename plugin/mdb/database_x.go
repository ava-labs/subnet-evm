// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package mdb

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/trace"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/avalanchego/x/merkledb"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethdb"
	"github.com/ava-labs/subnet-evm/trie"
	"github.com/ava-labs/subnet-evm/trie/trienode"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/prometheus/client_golang/prometheus"
)

const MerkleDBScheme = "merkleDBScheme"

var (
	_ ethdb.Database  = &WithMerkleDB{}
	_ trie.TrieOpener = &WithMerkleDB{}
	_ trie.Backender  = &WithMerkleDB{}

	ErrUnknownRoot = errors.New("unknown root")
)

type commit struct {
	stack  []*merkleDBTrie
	parent common.Hash
}

type WithMerkleDB struct {
	ethdb.Database
	merkleDB  merkledb.MerkleDB
	archiveDB ArchiveDB

	lock           sync.RWMutex
	pendingCommits map[common.Hash][]commit
	refCount       map[common.Hash]int
}

type backend WithMerkleDB

func toHash(id ids.ID) common.Hash {
	return common.BytesToHash(id[:])
}

func NewWithMerkleDB(db ethdb.Database, merkleDB merkledb.MerkleDB, archiveDB ArchiveDB) *WithMerkleDB {
	return &WithMerkleDB{
		Database:       db,
		merkleDB:       merkleDB,
		archiveDB:      archiveDB,
		pendingCommits: make(map[common.Hash][]commit),
		refCount:       make(map[common.Hash]int),
	}
}

func (db *WithMerkleDB) MerkleDB() merkledb.MerkleDB {
	return db.merkleDB
}

func (db *WithMerkleDB) GetAltMerkleRoot(ctx context.Context) (common.Hash, error) {
	id, err := db.merkleDB.GetAltMerkleRoot(ctx)
	if err != nil {
		return common.Hash{}, err
	}
	return toHash(id), nil
}

func (db *WithMerkleDB) Backend() trie.Backend {
	return (*backend)(db)
}

func (db *WithMerkleDB) getParent(root common.Hash) (merkledb.Trie, error) {
	pending, ok := db.pendingCommits[root]
	if ok {
		return pending[0].stack[0].tv, nil
	}
	ctx := context.TODO()
	hash, err := db.GetAltMerkleRoot(ctx)
	if err != nil {
		return nil, err
	}
	if hash == root {
		return db.merkleDB, nil
	}
	return nil, fmt.Errorf("%w: %x", ErrUnknownRoot, root)
}

func (db *WithMerkleDB) OpenTrie(root common.Hash) (trie.ITrie, error) {
	return db.OpenStorageTrie(root, common.Hash{}, root)
}

func (db *WithMerkleDB) OpenStorageTrie(stateRoot common.Hash, addrHash, root common.Hash) (trie.ITrie, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	parent, err := db.getParent(stateRoot)
	if errors.Is(err, ErrUnknownRoot) && db.archiveDB != nil {
		// try to open from archive
		return db.archiveDB.OpenTrie(stateRoot, addrHash, root)
	} else if err != nil {
		return nil, err
	}
	tr := &merkleDBTrie{
		parent:       parent,
		stateRoot:    stateRoot,
		owner:        addrHash,
		db:           db,
		originalRoot: root,
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

func (db *backend) UpdateAndReferenceRoot(root common.Hash, parent common.Hash, nodes *trienode.MergedNodeSet) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if err := db.update(root, parent, nodes); err != nil {
		return err
	}
	db.refCount[root]++
	return nil
}

// Update performs a state transition by committing dirty nodes contained
// in the given set in order to update state from the specified parent to
// the specified root.
func (db *backend) Update(root common.Hash, parent common.Hash, nodes *trienode.MergedNodeSet) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	return db.update(root, parent, nodes)
}

func (db *backend) update(root common.Hash, parent common.Hash, nodes *trienode.MergedNodeSet) error {
	t := nodes.Sets[common.Hash{}].Commit.(*merkleDBTrie)
	tvsToCommit := make([]*merkleDBTrie, 0, len(nodes.Sets))
	log.Debug("Update", "root", root, "parent", parent, "numTries", len(nodes.Sets))

	for t != nil {
		tvsToCommit = append(tvsToCommit, t)
		t = t.hashParent
	}
	if len(tvsToCommit) != len(nodes.Sets) {
		return fmt.Errorf("TrieView chain length does not match tries to commit (%d != %d)", len(tvsToCommit), len(nodes.Sets))
	}

	db.pendingCommits[root] = append(db.pendingCommits[root], commit{
		stack:  tvsToCommit,
		parent: parent,
	})
	db.refCount[parent]++ // keep the parent around
	return nil
}

// Commit writes all relevant trie nodes belonging to the specified state
// to disk. Report specifies whether logs will be displayed in info level.
func (db *backend) Commit(root common.Hash, report bool) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	log.Debug("Commit", "root", root)

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

		changes := make([]MapOps, 0, len(commit.stack))
		log.Debug("commit <--", "root", root, "parent", commit.parent)
		for i := len(commit.stack) - 1; i >= 0; i-- {
			change := commit.stack[i]
			if err := change.tv.CommitToDB(ctx); err != nil {
				return false, err
			}
			changes = append(changes, change.vc.MapOps)
		}
		if db.archiveDB != nil {
			// commit to archive
			if err := db.archiveDB.Commit(root, changes); err != nil {
				return false, err
			}
		}

		for _, commit := range db.pendingCommits[root] {
			db.dereference(commit.parent)
		}
		delete(db.pendingCommits, root)
		delete(db.refCount, root)
		return true, nil
	}
	return false, nil
}

// Scheme returns the identifier of used storage scheme.
func (db *backend) Scheme() string {
	return MerkleDBScheme
}

// Close closes the trie database backend and releases all held resources.
func (db *backend) Close() error {
	err := db.merkleDB.Close()
	log.Info("Closing merkleDB", "err", err)
	return err
}

func (db *backend) Dereference(root common.Hash) {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.dereference(root)
}

func (db *backend) dereference(root common.Hash) {
	if _, ok := db.refCount[root]; !ok {
		return
	}

	db.refCount[root]--
	if db.refCount[root] == 0 {
		delete(db.refCount, root)
		for _, commit := range db.pendingCommits[root] {
			db.dereference(commit.parent)
		}
		delete(db.pendingCommits, root)
	}
}

func (db *backend) Reference(root common.Hash, parent common.Hash) {
	db.lock.Lock()
	defer db.lock.Unlock()

	// only care about trie roots, which will reference the meta root aka empty
	if parent != (common.Hash{}) {
		return
	}
	db.refCount[root]++
}

func (db *backend) Cap(limit common.StorageSize) error {
	panic("implement me")
}

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
