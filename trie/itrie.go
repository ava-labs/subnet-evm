package trie

import (
	"github.com/ava-labs/subnet-evm/ethdb"
	"github.com/ava-labs/subnet-evm/trie/trienode"
	"github.com/ethereum/go-ethereum/common"
)

type ITrie interface {
	Hash() common.Hash
	Commit(collectLeaf bool) (common.Hash, *trienode.NodeSet)
	Delete(key []byte) error
	MustDelete(key []byte)
	Update(key, value []byte) error
	MustUpdate(key, value []byte)
	Get(key []byte) ([]byte, error)
	MustGet(key []byte) []byte

	ICopy() ITrie

	NodeIterator(start []byte) NodeIterator
	Prove(key []byte, fromLevel uint, proofDb ethdb.KeyValueWriter) error
}

func (t *Trie) ICopy() ITrie {
	return t.Copy()
}
