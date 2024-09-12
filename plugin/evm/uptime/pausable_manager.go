// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package uptime

import (
	"errors"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/uptime"
	"github.com/ava-labs/avalanchego/utils/set"
)

var (
	ErrPaused    = errors.New("node is paused")
	ErrNotPaused = errors.New("node is not paused")
	ErrPausedDc  = errors.New("paused node cannot be disconnected")
)

// Pausable is an interface that allows pausing and resuming uptime tracking
// for a node
type Pausable interface {
	// Pause pauses uptime tracking the node with the given ID
	Pause(ids.NodeID) error
	// Resume resumes uptime tracking for the node with the given ID
	Resume(ids.NodeID) error
	// IsPaused returns true if the node with the given ID is paused
	IsPaused(ids.NodeID) bool
}

// PausableManager is an interface that extends the uptime.Manager interface
// with the ability to pause and resume uptime tracking for a node
type PausableManager interface {
	Pausable
	uptime.Manager
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

// Pause pauses uptime tracking for the node with the given ID
// Pause can disconnect the node from the uptime.Manager if it is connected.
// Returns an error if the node is already paused.
func (p *pausableManager) Pause(nodeID ids.NodeID) error {
	if p.IsPaused(nodeID) {
		return ErrPaused
	}

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

// Resume resumes uptime tracking for the node with the given ID
// Resume can connect the node to the uptime.Manager if it was connected.
// Returns an error if the node is not paused.
func (p *pausableManager) Resume(nodeID ids.NodeID) error {
	if !p.IsPaused(nodeID) {
		return ErrNotPaused
	}
	p.pausedVdrs.Remove(nodeID)
	if p.connectedVdrs.Contains(nodeID) {
		return p.Manager.Connect(nodeID)
	}
	return nil
}

// IsPaused returns true if the node with the given ID is paused
func (p *pausableManager) IsPaused(nodeID ids.NodeID) bool {
	return p.pausedVdrs.Contains(nodeID)
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
