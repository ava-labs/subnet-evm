// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package uptime

import (
	"errors"

	"github.com/ava-labs/subnet-evm/plugin/evm/validators"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/uptime"
	"github.com/ava-labs/avalanchego/utils/set"
)

var _ validators.StateCallbackListener = &pausableManager{}

var ErrPausedDc = errors.New("paused node cannot be disconnected")

type PausableManager interface {
	uptime.Manager
	validators.StateCallbackListener
	IsPaused(nodeID ids.NodeID) bool
}

type pausableManager struct {
	uptime.Manager
	pausedVdrs set.Set[ids.NodeID]
	// connectedVdrs is a set of nodes that are connected to the manager.
	// This is used to keep track of the nodes that are connected to the manager
	// but are paused.
	connectedVdrs set.Set[ids.NodeID]
}

// NewPausableManager takes an uptime.Manager and returns a PausableManager
func NewPausableManager(manager uptime.Manager) PausableManager {
	return &pausableManager{
		pausedVdrs:    make(set.Set[ids.NodeID]),
		connectedVdrs: make(set.Set[ids.NodeID]),
		Manager:       manager,
	}
}

// Connect connects the node with the given ID to the uptime.Manager
// If the node is paused, it will not be connected
func (p *pausableManager) Connect(nodeID ids.NodeID) error {
	p.connectedVdrs.Add(nodeID)
	if !p.IsPaused(nodeID) && !p.Manager.IsConnected(nodeID) {
		return p.Manager.Connect(nodeID)
	}
	return nil
}

// Disconnect disconnects the node with the given ID from the uptime.Manager
// If the node is paused, it will not be disconnected
// Invariant: we should never have a connected paused node that is disconnecting
func (p *pausableManager) Disconnect(nodeID ids.NodeID) error {
	p.connectedVdrs.Remove(nodeID)
	if p.Manager.IsConnected(nodeID) {
		if p.IsPaused(nodeID) {
			// We should never see this case
			return ErrPausedDc
		}
		return p.Manager.Disconnect(nodeID)
	}
	return nil
}

// StartTracking starts tracking uptime for the nodes with the given IDs
// If a node is paused, it will not be tracked
func (p *pausableManager) StartTracking(nodeIDs []ids.NodeID) error {
	var activeNodeIDs []ids.NodeID
	for _, nodeID := range nodeIDs {
		if !p.IsPaused(nodeID) {
			activeNodeIDs = append(activeNodeIDs, nodeID)
		}
	}
	return p.Manager.StartTracking(activeNodeIDs)
}

// OnValidatorAdded is called when a validator is added.
// If the node is inactive, it will be paused.
func (p *pausableManager) OnValidatorAdded(vID ids.ID, nodeID ids.NodeID, startTime uint64, isActive bool) {
	if !isActive {
		p.pause(nodeID)
	}
}

// OnValidatorRemoved is called when a validator is removed.
// If the node is already paused, it will be resumed.
func (p *pausableManager) OnValidatorRemoved(vID ids.ID, nodeID ids.NodeID) {
	if p.IsPaused(nodeID) {
		p.resume(nodeID)
	}
}

// OnValidatorStatusUpdated is called when the status of a validator is updated.
// If the node is active, it will be resumed. If the node is inactive, it will be paused.
func (p *pausableManager) OnValidatorStatusUpdated(vID ids.ID, nodeID ids.NodeID, isActive bool) {
	if isActive {
		p.resume(nodeID)
	} else {
		p.pause(nodeID)
	}
}

// IsPaused returns true if the node with the given ID is paused.
func (p *pausableManager) IsPaused(nodeID ids.NodeID) bool {
	return p.pausedVdrs.Contains(nodeID)
}

// pause pauses uptime tracking for the node with the given ID
// pause can disconnect the node from the uptime.Manager if it is connected.
// Returns an error if the node is already paused.
func (p *pausableManager) pause(nodeID ids.NodeID) error {
	p.pausedVdrs.Add(nodeID)
	if p.Manager.IsConnected(nodeID) {
		// If the node is connected, then we need to disconnect it from
		// manager
		// This should be fine in case tracking has not started yet since
		// the inner manager should handle disconnects accordingly
		return p.Manager.Disconnect(nodeID)
	}
	return nil
}

// resume resumes uptime tracking for the node with the given ID
// resume can connect the node to the uptime.Manager if it was connected.
// Returns an error if the node is not paused.
func (p *pausableManager) resume(nodeID ids.NodeID) error {
	p.pausedVdrs.Remove(nodeID)
	if p.connectedVdrs.Contains(nodeID) {
		return p.Manager.Connect(nodeID)
	}
	return nil
}
