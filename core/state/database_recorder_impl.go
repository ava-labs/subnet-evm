// (c) 2020-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package state

import (
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethdb"
	"github.com/ava-labs/subnet-evm/trie"
	"github.com/ava-labs/subnet-evm/trie/trienode"
	"github.com/ethereum/go-ethereum/common"
)

type recordingTrie struct {
	Trie
	r *recording
}

type RecordingDatabase struct {
	Database
	r *recording
}

func NewRecordingDatabase(db Database, r *recording) *RecordingDatabase {
	return &RecordingDatabase{Database: db, r: r}
}

func (db *RecordingDatabase) OpenTrie(root common.Hash) (Trie, error) {
	tr, err := db.Database.OpenTrie(root)
	tr = &recordingTrie{Trie: tr, r: db.r}
	db.r.Record(db, "OpenTrie", root).Returns(tr, err)
	return tr, err
}

func (db *RecordingDatabase) OpenStorageTrie(stateRoot common.Hash, addrHash, root common.Hash) (Trie, error) {
	tr, err := db.Database.OpenStorageTrie(stateRoot, addrHash, root)
	tr = &recordingTrie{Trie: tr, r: db.r}
	db.r.Record(db, "OpenStorageTrie", stateRoot, addrHash, root).Returns(tr, err)
	return tr, err
}

func (db *RecordingDatabase) CopyTrie(trie Trie) Trie {
	if r, ok := trie.(*recordingTrie); ok {
		trie = r.Trie
	}
	tr := db.Database.CopyTrie(trie)
	tr = &recordingTrie{Trie: tr, r: db.r}
	db.r.Record(db, "CopyTrie", trie).Returns(tr)
	return tr
}

func (t *recordingTrie) GetKey(key []byte) (ret0 []byte) {
	defer t.r.Record(t, "GetKey", key).Returns(ret0)
	return t.Trie.GetKey(key)
}

func (t *recordingTrie) GetStorage(addr common.Address, key []byte) ([]byte, error) {
	ret0, err := t.Trie.GetStorage(addr, key)
	t.r.Record(t, "GetStorage", addr, key).Returns(ret0, err)
	return ret0, err
}

func (t *recordingTrie) GetAccount(address common.Address) (*types.StateAccount, error) {
	ret0, err := t.Trie.GetAccount(address)
	t.r.Record(t, "GetAccount", address).Returns(ret0, err)
	return ret0, err
}

func (t *recordingTrie) UpdateStorage(addr common.Address, key, value []byte) error {
	err := t.Trie.UpdateStorage(addr, key, value)
	t.r.Record(t, "UpdateStorage", addr, common.CopyBytes(key), value).Returns(err)
	return err
}

func (t *recordingTrie) UpdateAccount(address common.Address, account *types.StateAccount) error {
	err := t.Trie.UpdateAccount(address, account)
	t.r.Record(t, "UpdateAccount", address, account).Returns(err)
	return err
}

func (t *recordingTrie) DeleteStorage(addr common.Address, key []byte) error {
	err := t.Trie.DeleteStorage(addr, key)
	t.r.Record(t, "DeleteStorage", addr, key).Returns(err)
	return err
}

func (t *recordingTrie) DeleteAccount(address common.Address) error {
	err := t.Trie.DeleteAccount(address)
	t.r.Record(t, "DeleteAccount", address).Returns(err)
	return err
}

func (t *recordingTrie) Hash() common.Hash {
	op := t.r.Record(t, "Hash")
	ret0 := t.Trie.Hash()
	op.Returns(ret0)
	return ret0
}

func (t *recordingTrie) Commit(collectLeaf bool) (common.Hash, *trienode.NodeSet) {
	ret0, ret1 := t.Trie.Commit(collectLeaf)
	t.r.Record(t, "Commit", collectLeaf).Returns(ret0, ret1)
	return ret0, ret1
}

func (t *recordingTrie) NodeIterator(startKey []byte) trie.NodeIterator {
	panic("not implemented")
}

func (t *recordingTrie) Prove(key []byte, fromLevel uint, proofDb ethdb.KeyValueWriter) error {
	panic("not implemented")
}

func (t *recordingTrie) ITrie() trie.ITrie {
	return t.Trie.(*trie.StateTrie).ITrie()
}
