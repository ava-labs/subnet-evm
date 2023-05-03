// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package sharedmemory

import (
	"encoding/binary"

	"github.com/ava-labs/avalanchego/chains/atomic"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/wrappers"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/utils/codec"
	"github.com/ethereum/go-ethereum/common"
)

const (
	syncKeyPrefix   = 's' // latest serial number is stored at this key
	spentUTXOPrefix = 'u'
)

type trie interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
}

type StateTrie struct {
	StateDB contract.StateDB
}

func (s *StateTrie) Get(key []byte) ([]byte, error) {
	return []byte(s.StateDB.GetStateVariableLength(ContractAddress, string(key))), nil
}

func (s *StateTrie) Put(key, value []byte) error {
	s.StateDB.SetStateVariableLength(ContractAddress, string(key), string(value))
	return nil
}

func IsSpent(utxo ids.ID, state contract.StateDB) (bool, error) {
	return isSpent(utxo, &StateTrie{state})
}

func isSpent(utxo ids.ID, trie trie) (bool, error) {
	key := mkSpentUTXOKey(utxo)
	data, err := trie.Get(key)
	if err != nil {
		return false, err
	}
	return len(data) > 0, nil
}

func markSpent(utxo ids.ID, trie trie) error {
	key := mkSpentUTXOKey(utxo)
	return trie.Put(key, []byte{0})
}

type SyncRecord struct {
	// TODO: maybe we want height here too?
	ChainID  ids.ID           `serialize:"true"`
	Requests *atomic.Requests `serialize:"true"`
}

// addAtomicOpsToSyncRecord adds the atomic ops for [chainID] represented by
// [requests] to the [trie] provided.
func addAtomicOpsToSyncRecord(height uint64, chainID ids.ID, requests *atomic.Requests, trie trie) error {
	// Get the key to store the next sync record at
	key, err := nextSyncKey(trie)
	if err != nil {
		return err
	}

	// Create the sync record
	syncRecord := &SyncRecord{
		ChainID:  chainID,
		Requests: requests,
	}

	// Marshal the sync record
	data, err := codec.Codec.Marshal(codec.CodecVersion, syncRecord)
	if err != nil {
		return err
	}

	// Put the marshalled sync record in the trie
	if err := trie.Put(key, data); err != nil {
		return err
	}

	return nil
}

func nextSyncKey(trie trie) ([]byte, error) {
	// Get the last used serial number from the trie
	serialNumber, err := GetSerialNumber(trie)
	if err != nil {
		return nil, err
	}

	// Create the key and persist the new serial number
	// to the trie
	key := mkSyncKey(serialNumber + 1)
	if err := trie.Put([]byte{syncKeyPrefix}, key[1:]); err != nil {
		return nil, err
	}

	return key, nil
}

func mkSyncKey(serialNumber uint64) []byte {
	key := make([]byte, 1+wrappers.LongLen)
	key[0] = syncKeyPrefix
	binary.BigEndian.PutUint64(key[1:], serialNumber)
	return key
}

func GetSyncRecord(serialNumber uint64, trie trie) (SyncRecord, error) {
	key := mkSyncKey(serialNumber)
	data, err := trie.Get(key)
	if err != nil {
		return SyncRecord{}, err
	}

	var syncRecord SyncRecord
	if _, err := codec.Codec.Unmarshal(data, &syncRecord); err != nil {
		return SyncRecord{}, err
	}

	return syncRecord, nil
}

func GetSerialNumber(trie trie) (uint64, error) {
	data, err := trie.Get([]byte{syncKeyPrefix})
	if err != nil {
		return 0, err
	}
	if len(data) == 0 {
		return 0, nil
	}
	return binary.BigEndian.Uint64(data), nil
}

func mkSpentUTXOKey(utxo ids.ID) []byte {
	packer := wrappers.Packer{Bytes: make([]byte, 1+common.HashLength)}
	packer.PackByte(spentUTXOPrefix)
	packer.PackFixedBytes(utxo[:])
	return packer.Bytes
}
