// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package mdb

import (
	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/subnet-evm/trie"
	"github.com/ethereum/go-ethereum/common"
)

// Just to support this code in dumping:
//
// func (it *Iterator) Next() bool {
// 	for it.nodeIt.Next(true) {
// 		if it.nodeIt.Leaf() {
// 			it.Key = it.nodeIt.LeafKey()
// 			it.Value = it.nodeIt.LeafBlob()
// 			return true
// 		}
// 	}
// 	it.Key = nil
// 	it.Value = nil
// 	it.Err = it.nodeIt.Error()
// 	return false
// }

var _ trie.NodeIterator = &mdbIterator{}

type mdbIterator struct {
	prefix      []byte
	it          database.Iterator
	expectedLen int
}

func (t *merkleDBTrie) NodeIterator(start []byte) trie.NodeIterator {
	prefix := PrefixBytes(t.owner, start)
	startPos := PrefixBytes(t.owner, append(start, 0)) // start here to skip the root node

	if !t.hashed {
		if err := t.hash(); err != nil {
			return &mdbIterator{
				it: &database.IteratorError{Err: err},
			}
		}
	}
	return &mdbIterator{
		prefix:      prefix,
		it:          t.tv.NewIteratorWithStartAndPrefix(startPos, prefix),
		expectedLen: common.HashLength,
	}
}

func (it *mdbIterator) Next(bool) bool {
	for it.it.Next() {
		if len(it.LeafKey()) == it.expectedLen {
			return true
		}
	}
	it.it.Release()
	return false
}

func (it *mdbIterator) Error() error { return it.it.Error() }
func (it *mdbIterator) Leaf() bool   { return true }

func (it *mdbIterator) LeafBlob() []byte { return it.it.Value() }
func (it *mdbIterator) LeafKey() []byte  { return it.it.Key()[len(it.prefix):] }

// unimplemented
func (it *mdbIterator) Hash() common.Hash                      { return common.Hash{} }
func (it *mdbIterator) Parent() common.Hash                    { panic("implement me") }
func (it *mdbIterator) Path() []byte                           { panic("implement me") }
func (it *mdbIterator) NodeBlob() []byte                       { panic("implement me") }
func (it *mdbIterator) LeafProof() [][]byte                    { panic("implement me") }
func (it *mdbIterator) AddResolver(resolver trie.NodeResolver) { panic("implement me") }
