// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package sharedmemory

import (
	"github.com/ava-labs/avalanchego/chains/atomic"
	"github.com/ava-labs/avalanchego/codec"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/wrappers"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
)

const (
	syncKeyPrefix   = 's'
	spentUTXOPrefix = 'u'
)

type trie interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
}

type stateTrie struct {
	s contract.StateDB
}

func (s *stateTrie) Get(key []byte) ([]byte, error) {
	return []byte(s.s.GetStateVariableLength(ContractAddress, string(key))), nil
}

func (s *stateTrie) Put(key, value []byte) error {
	s.s.SetStateVariableLength(ContractAddress, string(key), string(value))
	return nil
}

func isSpent(codec codec.Manager, utxo ids.ID, trie trie) (bool, error) {
	key := mkSpentUTXOKey(utxo)
	data, err := trie.Get(key)
	if err != nil {
		return false, err
	}
	return len(data) > 0, nil
}

func markSpent(codec codec.Manager, utxo ids.ID, trie trie) error {
	key := mkSpentUTXOKey(utxo)
	return trie.Put(key, []byte{0})
}

func addAtomicOpsToSyncRecord(codec codec.Manager, height uint64, chainID ids.ID, requests *atomic.Requests, trie trie) error {
	key := mkSyncKey(height, chainID)

	// First, get any existing sync record from the trie
	data, err := trie.Get(key)
	if err != nil {
		return err
	}

	var syncRecord *atomic.Requests
	if len(data) > 0 {
		// If there is an existing sync record, unmarshal it
		if _, err := codec.Unmarshal(data, syncRecord); err != nil {
			return err
		}
	} else {
		// If there is no existing sync record, create a new one
		syncRecord = &atomic.Requests{}
	}

	// Add the atomic ops to the sync record
	syncRecord.PutRequests = append(syncRecord.PutRequests, requests.PutRequests...)
	syncRecord.RemoveRequests = append(syncRecord.RemoveRequests, requests.RemoveRequests...)

	// Marshal the sync record
	data, err = codec.Marshal(0, syncRecord)
	if err != nil {
		return err
	}

	// Put the marshalled sync record in the trie
	if err := trie.Put(key, data); err != nil {
		return err
	}

	return nil
}

func mkSyncKey(height uint64, blockchainID ids.ID) []byte {
	packer := wrappers.Packer{Bytes: make([]byte, 0, 1+wrappers.LongLen+common.HashLength)}
	packer.PackByte(syncKeyPrefix)
	packer.PackLong(height)
	packer.PackFixedBytes(blockchainID[:])
	return packer.Bytes
}

func mkSpentUTXOKey(utxo ids.ID) []byte {
	packer := wrappers.Packer{Bytes: make([]byte, 0, 1+common.HashLength)}
	packer.PackByte(spentUTXOPrefix)
	packer.PackFixedBytes(utxo[:])
	return packer.Bytes
}
