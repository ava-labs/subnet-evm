// (c) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// warptest exposes common functionality for testing the warp package.
package warptest

import (
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/subnet-evm/plugin/evm/validators/interfaces"
	stateinterfaces "github.com/ava-labs/subnet-evm/plugin/evm/validators/state/interfaces"
)

var _ interfaces.ValidatorReader = &NoOpValidatorReader{}

type NoOpValidatorReader struct{}

func (NoOpValidatorReader) CalculateUptime(nodeID ids.NodeID) (time.Duration, time.Time, error) {
	return 0, time.Time{}, nil
}

func (NoOpValidatorReader) CalculateUptimePercent(nodeID ids.NodeID) (float64, error) {
	return 0, nil
}

func (NoOpValidatorReader) CalculateUptimePercentFrom(nodeID ids.NodeID, startTime time.Time) (float64, error) {
	return 0, nil
}

func (NoOpValidatorReader) GetNodeIDs() set.Set[ids.NodeID] {
	return set.Set[ids.NodeID]{}
}

func (NoOpValidatorReader) GetValidationIDs() set.Set[ids.ID] {
	return set.Set[ids.ID]{}
}

func (NoOpValidatorReader) GetValidator(vID ids.ID) (stateinterfaces.Validator, error) {
	return stateinterfaces.Validator{}, nil
}
