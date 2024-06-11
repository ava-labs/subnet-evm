package state

import (
	"github.com/ava-labs/coreth/core/types"
	"github.com/ava-labs/coreth/predicate"
	"github.com/ava-labs/coreth/utils"
	"github.com/ethereum/go-ethereum/common"
)

// GetPredicateStorageSlots returns the storage slots associated with the address, index pair.
// A list of access tuples can be included within transaction types post EIP-2930. The address
// is declared directly on the access tuple and the index is the i'th occurrence of an access
// tuple with the specified address.
//
// Ex. AccessList[[AddrA, Predicate1], [AddrB, Predicate2], [AddrA, Predicate3]]
// In this case, the caller could retrieve predicates 1-3 with the following calls:
// GetPredicateStorageSlots(AddrA, 0) -> Predicate1
// GetPredicateStorageSlots(AddrB, 0) -> Predicate2
// GetPredicateStorageSlots(AddrA, 1) -> Predicate3
func (s *StateDB) GetPredicateStorageSlots(address common.Address, index int) ([]byte, bool) {
	predicates := predicate.GetPredicatesFromAccessList(s._accessList, address)
	if index >= len(predicates) {
		return nil, false
	}
	return predicates[index], true
}

// SetPredicateStorageSlots sets the predicate storage slots for the given address
// TODO: This test-only method can be replaced with setting the access list.
func (s *StateDB) SetPredicateStorageSlots(address common.Address, predicates [][]byte) {
	s._accessList = make(types.AccessList, 0, len(predicates))
	for _, predicateBytes := range predicates {
		s._accessList = append(s._accessList, types.AccessTuple{
			Address:     address,
			StorageKeys: utils.BytesToHashSlice(predicateBytes),
		})
	}
}

// GetTxHash returns the current tx hash on the StateDB set by SetTxContext.
func (s *StateDB) GetTxHash() common.Hash {
	return s.thash
}

// GetLogData returns the underlying topics and data from each log included in the StateDB
// Test helper function.
func (s *StateDB) GetLogData() ([][]common.Hash, [][]byte) {
	var logData [][]byte
	var topics [][]common.Hash
	for _, lgs := range s.logs {
		for _, log := range lgs {
			topics = append(topics, log.Topics)
			logData = append(logData, common.CopyBytes(log.Data))
		}
	}
	return topics, logData
}
