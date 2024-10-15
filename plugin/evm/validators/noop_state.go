package validators

import (
	"time"

	ids "github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/set"
)

var NoOpState State = &noOpState{}

type noOpState struct{}

func (n *noOpState) GetStatus(vID ids.ID) (bool, error) { return false, nil }

func (n *noOpState) GetValidationIDs() set.Set[ids.ID] { return set.NewSet[ids.ID](0) }

func (n *noOpState) GetNodeIDs() set.Set[ids.NodeID] { return set.NewSet[ids.NodeID](0) }

func (n *noOpState) GetValidator(nodeID ids.NodeID) (*ValidatorOutput, error) {
	return nil, nil
}

func (n *noOpState) GetNodeID(vID ids.ID) (ids.NodeID, error) { return ids.NodeID{}, nil }

func (n *noOpState) GetSubnetID(chainID ids.ID) (ids.ID, error) { return ids.ID{}, nil }

func (n *noOpState) GetUptime(
	nodeID ids.NodeID,
) (upDuration time.Duration, lastUpdated time.Time, err error) {
	return 0, time.Time{}, nil
}

func (n *noOpState) SetUptime(
	nodeID ids.NodeID,
	upDuration time.Duration,
	lastUpdated time.Time,
) error {
	return nil
}

func (n *noOpState) GetStartTime(
	nodeID ids.NodeID,
) (startTime time.Time, err error) {
	return time.Time{}, nil
}

func (n *noOpState) AddValidator(vID ids.ID, nodeID ids.NodeID, startTimestamp uint64, isActive bool) error {
	return nil
}

func (n *noOpState) DeleteValidator(vID ids.ID) error {
	return nil
}
func (n *noOpState) WriteState() error { return nil }

func (n *noOpState) SetStatus(vID ids.ID, isActive bool) error { return nil }

func (n *noOpState) RegisterListener(StateCallbackListener) {}
