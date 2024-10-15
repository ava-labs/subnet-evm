package validators

import (
	"sync"

	ids "github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/set"
)

type lockedStateReader struct {
	lock sync.Locker
	s    StateReader
}

func NewLockedStateReader(lock sync.Locker, s State) StateReader {
	return &lockedStateReader{
		lock: lock,
		s:    s,
	}
}

func (s *lockedStateReader) GetStatus(vID ids.ID) (bool, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.s.GetStatus(vID)
}

func (s *lockedStateReader) GetValidationIDs() set.Set[ids.ID] {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.s.GetValidationIDs()
}

func (s *lockedStateReader) GetNodeIDs() set.Set[ids.NodeID] {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.s.GetNodeIDs()
}

func (s *lockedStateReader) GetValidator(nodeID ids.NodeID) (*ValidatorOutput, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.s.GetValidator(nodeID)
}

func (s *lockedStateReader) GetNodeID(vID ids.ID) (ids.NodeID, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.s.GetNodeID(vID)
}
