// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package mdb

import (
	"context"
	"errors"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/utils/maybe"
	"github.com/ava-labs/avalanchego/x/merkledb"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/ethdb"
	"github.com/ava-labs/subnet-evm/trie"
	"github.com/ava-labs/subnet-evm/trie/trienode"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

var (
	_ trie.ITrie          = &merkleDBTrie{}
	_ state.HashChainTrie = &merkleDBTrie{}
)

type merkleDBTrie struct {
	vc     merkledb.ViewChanges
	parent merkledb.Trie
	tv     merkledb.TrieView

	stateRoot common.Hash
	owner     common.Hash
	hashed    bool

	db         *WithMerkleDB
	hashParent *merkleDBTrie
}

func (t *merkleDBTrie) SetLastHashed(tr trie.ITrie) {
	if tr == nil {
		return
	}
	t.hashParent = tr.(*merkleDBTrie)
}

func (t *merkleDBTrie) initialize() {
	t.vc.MapOps = make(map[string]maybe.Maybe[[]byte])
}

func (t *merkleDBTrie) Get(k []byte) ([]byte, error) {
	key := PrefixBytes(t.owner, k)
	val, ok := t.vc.MapOps[string(key)]
	if ok {
		return val.Value(), nil
	}
	trieVal, err := t.parent.GetValue(context.Background(), key)
	switch {
	case errors.Is(err, database.ErrNotFound):
		return nil, nil
	case err != nil:
		return nil, err
	}
	return trieVal, nil
}

func (t *merkleDBTrie) Update(k, value []byte) error {
	k = common.CopyBytes(k)
	t.hashed = false
	key := PrefixBytes(t.owner, k)
	val := maybe.Nothing[[]byte]()
	if len(value) > 0 {
		val = maybe.Some(value)
	}

	t.vc.MapOps[string(key)] = val
	return nil
}

func (t *merkleDBTrie) Delete(k []byte) error {
	t.hashed = false
	key := PrefixBytes(t.owner, k)
	t.vc.MapOps[string(key)] = maybe.Nothing[[]byte]()
	return nil
}

func (t *merkleDBTrie) MustGet(key []byte) []byte {
	res, err := t.Get(key)
	if err != nil {
		log.Error("Unhandled trie error in merkledbTrie.Get", "err", err)
	}
	return res
}

func (t *merkleDBTrie) MustUpdate(key, value []byte) {
	if err := t.Update(key, value); err != nil {
		log.Error("Unhandled trie error in merkledbTrie.Update", "err", err)
	}
}

func (t *merkleDBTrie) MustDelete(key []byte) {
	if err := t.Delete(key); err != nil {
		log.Error("Unhandled trie error in merkledbTrie.Delete", "err", err)
	}
}

func (t *merkleDBTrie) Commit(collectLeaf bool) (common.Hash, *trienode.NodeSet) {
	root := t.Hash()
	log.Debug("mtree commit", "root", root, "owner", t.owner)
	nodeSet := trienode.NewNodeSet(t.owner)
	nodeSet.Commit = t
	return root, nodeSet
}

func (t *merkleDBTrie) Hash() common.Hash {
	if !t.hashed {
		if err := t.hash(); err != nil {
			panic(err)
		}
	}
	id, err := t.tv.GetAltMerkleRoot(context.Background())
	if err != nil {
		panic(err)
	}
	t.hashed = true
	hash := common.BytesToHash(id[:])
	log.Debug("mtree hash", "root", hash, "owner", t.owner)
	return hash
}

func (t *merkleDBTrie) hash() error {
	rootPrefix := merkledb.ToKey(PrefixBytes(t.owner, nil))
	parent := t.parent
	if t.hashParent != nil {
		parent = t.hashParent.tv
	}

	tv, err := parent.NewViewWithRootPrefix(context.Background(), t.vc, rootPrefix)
	if err != nil {
		return err
	}
	t.tv = tv
	return nil
}

func PrefixBytes(owner common.Hash, key []byte) []byte {
	if owner == (common.Hash{}) {
		return key
	}
	return append(append(owner[:], []byte{0}...), key...)
}

func (t *merkleDBTrie) ICopy() trie.ITrie {
	vc := merkledb.ViewChanges{
		BatchOps:     slices.Clone(t.vc.BatchOps),
		MapOps:       maps.Clone(t.vc.MapOps),
		ConsumeBytes: t.vc.ConsumeBytes,
	}
	return &merkleDBTrie{
		vc:        vc,
		parent:    t.parent,
		owner:     t.owner,
		db:        t.db,
		stateRoot: t.stateRoot,

		// Note we don't copy the id or hashed fields
		// this forces a rehash on the copy
	}
}

func (t *merkleDBTrie) Prove(key []byte, fromLevel uint, proofDb ethdb.KeyValueWriter) error {
	panic("implement me")
}
