// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"encoding/binary"
	"errors"

	"github.com/ava-labs/avalanchego/chains/atomic"
	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/database/versiondb"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/wrappers"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/contracts/sharedmemory"
	"github.com/ethereum/go-ethereum/common"
)

const maxOpsPerBatch = 10_000

var (
	lastAppliedKey = []byte("lastApplied")
)

type stateProvider interface {
	StateAt(root common.Hash) (*state.StateDB, error)
}

type SharedMemorySyncer interface {
	SyncSharedMemoryToState(stateRoot common.Hash) error
	IncLastApplied(numLogs int) error
}

type sharedMemorySyncer struct {
	metadataDB    database.Database
	versionDB     *versiondb.Database
	stateProvider stateProvider
	sharedMemory  atomic.SharedMemory

	lastApplied uint64
}

func newSharedMemorySyncer(
	metadataDB database.Database, versionDB *versiondb.Database, stateProvider stateProvider, sharedMemory atomic.SharedMemory,
) (*sharedMemorySyncer, error) {
	s := &sharedMemorySyncer{
		metadataDB:    metadataDB,
		stateProvider: stateProvider,
		sharedMemory:  sharedMemory,
		versionDB:     versionDB,
	}
	return s, s.initialize()
}

func (s *sharedMemorySyncer) initialize() error {
	// Read the last applied block from the database
	lastAppliedBytes, err := s.metadataDB.Get(lastAppliedKey)
	if errors.Is(err, database.ErrNotFound) {
		return nil
	} else if err != nil {
		return err
	}

	s.lastApplied = binary.BigEndian.Uint64(lastAppliedBytes)
	return nil
}

func (s *sharedMemorySyncer) SyncSharedMemoryToState(stateRoot common.Hash) error {
	// Get the state at [stateRoot]
	stateDB, err := s.stateProvider.StateAt(stateRoot)
	if err != nil {
		return err
	}
	trie := &sharedmemory.StateTrie{StateDB: stateDB}
	lastSerialNumber, err := sharedmemory.GetSerialNumber(trie)
	if err != nil {
		return err
	}
	writer := NewSharedMemoryWriter()
	for i := s.lastApplied + 1; i <= lastSerialNumber; i++ {
		// Get the sync record
		syncRecord, err := sharedmemory.GetSyncRecord(i, trie)
		if err != nil {
			return err
		}
		// Apply the operations to shared memory
		writer.AddSharedMemoryRequests(syncRecord.ChainID, syncRecord.Requests)

		// Keep the number of operations per batch under [maxOpsPerBatch]
		if writer.requestsLen >= maxOpsPerBatch {
			if err := s.applyOps(i, writer.requests); err != nil {
				return err
			}
			writer = NewSharedMemoryWriter()
		}
	}

	// Apply the remaining operations (if any) and update the lastApplied serial number
	return s.applyOps(lastSerialNumber, writer.requests)
}

func (s *sharedMemorySyncer) putLastApplied(serialNumber uint64) error {
	lastAppliedBytes := make([]byte, wrappers.LongLen)
	binary.BigEndian.PutUint64(lastAppliedBytes, serialNumber)
	return s.metadataDB.Put(lastAppliedKey, lastAppliedBytes)
}

func (s *sharedMemorySyncer) applyOps(serialNumber uint64, ops map[ids.ID]*atomic.Requests) error {
	if err := s.putLastApplied(serialNumber); err != nil {
		return err
	}
	vdbBatch, err := s.versionDB.CommitBatch()
	if err != nil {
		return err
	}

	// Atomically apply the ops and update the lastApplied serial number
	s.sharedMemory.Apply(ops, vdbBatch)
	s.lastApplied = serialNumber
	return nil
}

func (s *sharedMemorySyncer) IncLastApplied(numLogs int) error {
	s.lastApplied += uint64(numLogs)
	return s.putLastApplied(s.lastApplied)
}
