package state

import (
	"github.com/ava-labs/coreth/core/types"
	"github.com/ethereum/go-ethereum/common"
)

func (s *StateDB) AccessList() types.AccessList {
	return s._accessList
}

// Warning: Test Only
func (s *StateDB) SetAccessList(list types.AccessList) {
	s._accessList = list
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
