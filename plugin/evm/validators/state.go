// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package validators

import (
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/uptime"
	"github.com/ava-labs/avalanchego/utils/set"
)

var _ uptime.State = &state{}

type State interface {
	uptime.State
	// AddNewValidator adds a new validator to the state
	AddNewValidator(vID ids.ID, nodeID ids.NodeID, startTimestamp uint64, isActive bool) error
	// DeleteValidator deletes the validator from the state
	DeleteValidator(vID ids.ID) error
	// WriteState writes the validator state to the disk
	WriteState() error

	// SetStatus sets the active status of the validator with the given vID
	SetStatus(vID ids.ID, isActive bool) error
	// GetStatus returns the active status of the validator with the given vID
	GetStatus(vID ids.ID) (bool, error)

	// GetValidationIDs returns the validation IDs in the state
	GetValidationIDs() set.Set[ids.ID]
	// GetValidatorIDs returns the validator node IDs in the state
	GetValidatorIDs() set.Set[ids.NodeID]

	// RegisterListener registers a listener to the state
	RegisterListener(StateCallbackListener)
}

// StateCallbackListener is a listener for the validator state
type StateCallbackListener interface {
	// OnValidatorAdded is called when a new validator is added
	OnValidatorAdded(vID ids.ID, nodeID ids.NodeID, startTime uint64, isActive bool)
	// OnValidatorRemoved is called when a validator is removed
	OnValidatorRemoved(vID ids.ID, nodeID ids.NodeID)
	// OnValidatorStatusUpdated is called when a validator status is updated
	OnValidatorStatusUpdated(vID ids.ID, nodeID ids.NodeID, isActive bool)
}

type validatorData struct {
	UpDuration  time.Duration `serialize:"true"`
	LastUpdated uint64        `serialize:"true"`
	NodeID      ids.NodeID    `serialize:"true"`
	StartTime   uint64        `serialize:"true"`
	IsActive    bool          `serialize:"true"`

	validationID ids.ID // database key
	lastUpdated  time.Time
	startTime    time.Time
}

type state struct {
	data  map[ids.ID]*validatorData // vID -> validatorData
	index map[ids.NodeID]ids.ID     // nodeID -> vID
	// updatedData tracks the updates since las WriteValidator was called
	updatedData map[ids.ID]bool // vID -> true(updated)/false(deleted)
	db          database.Database

	listeners []StateCallbackListener
}

// NewState creates a new State, it also loads the data from the disk
func NewState(db database.Database) (State, error) {
	s := &state{
		index:       make(map[ids.NodeID]ids.ID),
		data:        make(map[ids.ID]*validatorData),
		updatedData: make(map[ids.ID]bool),
		db:          db,
	}
	if err := s.loadFromDisk(); err != nil {
		return nil, fmt.Errorf("failed to load data from disk: %w", err)
	}
	return s, nil
}

// GetUptime returns the uptime of the validator with the given nodeID
func (s *state) GetUptime(
	nodeID ids.NodeID,
) (time.Duration, time.Time, error) {
	data, err := s.getData(nodeID)
	if err != nil {
		return 0, time.Time{}, err
	}
	return data.UpDuration, data.lastUpdated, nil
}

// SetUptime sets the uptime of the validator with the given nodeID
func (s *state) SetUptime(
	nodeID ids.NodeID,
	upDuration time.Duration,
	lastUpdated time.Time,
) error {
	data, err := s.getData(nodeID)
	if err != nil {
		return err
	}
	data.UpDuration = upDuration
	data.lastUpdated = lastUpdated

	s.updatedData[data.validationID] = true
	return nil
}

// GetStartTime returns the start time of the validator with the given nodeID
func (s *state) GetStartTime(nodeID ids.NodeID) (time.Time, error) {
	data, err := s.getData(nodeID)
	if err != nil {
		return time.Time{}, err
	}
	return data.startTime, nil
}

// AddNewValidator adds a new validator to the state
// the new validator is marked as updated and will be written to the disk when WriteState is called
func (s *state) AddNewValidator(vID ids.ID, nodeID ids.NodeID, startTimestamp uint64, isActive bool) error {
	startTimeUnix := time.Unix(int64(startTimestamp), 0)

	data := &validatorData{
		NodeID:       nodeID,
		validationID: vID,
		IsActive:     isActive,
		StartTime:    startTimestamp,
		UpDuration:   0,
		LastUpdated:  startTimestamp,
		lastUpdated:  startTimeUnix,
		startTime:    startTimeUnix,
	}
	if err := s.putData(vID, data); err != nil {
		return err
	}

	s.updatedData[vID] = true

	for _, listener := range s.listeners {
		listener.OnValidatorAdded(vID, nodeID, startTimestamp, isActive)
	}
	return nil
}

// DeleteValidator marks the validator as deleted
// marked validator will be deleted from disk when WriteState is called
func (s *state) DeleteValidator(vID ids.ID) error {
	data, exists := s.data[vID]
	if !exists {
		return database.ErrNotFound
	}
	delete(s.data, data.validationID)
	delete(s.index, data.NodeID)

	// mark as deleted for WriteValidator
	s.updatedData[data.validationID] = false

	for _, listener := range s.listeners {
		listener.OnValidatorRemoved(vID, data.NodeID)
	}
	return nil
}

// WriteState writes the updated state to the disk
func (s *state) WriteState() error {
	// TODO: consider adding batch size
	batch := s.db.NewBatch()
	for vID, updated := range s.updatedData {
		if updated {
			data := s.data[vID]
			data.LastUpdated = uint64(data.lastUpdated.Unix())
			// should never change but in case
			data.StartTime = uint64(data.startTime.Unix())

			dataBytes, err := vdrCodec.Marshal(codecVersion, data)
			if err != nil {
				return err
			}
			if err := batch.Put(vID[:], dataBytes); err != nil {
				return err
			}
		} else { // deleted
			if err := batch.Delete(vID[:]); err != nil {
				return err
			}
		}
		// we're done, remove the updated marker
		delete(s.updatedData, vID)
	}
	return batch.Write()
}

// SetStatus sets the active status of the validator with the given vID
func (s *state) SetStatus(vID ids.ID, isActive bool) error {
	data, exists := s.data[vID]
	if !exists {
		return database.ErrNotFound
	}
	data.IsActive = isActive
	s.updatedData[vID] = true

	for _, listener := range s.listeners {
		listener.OnValidatorStatusUpdated(vID, data.NodeID, isActive)
	}
	return nil
}

// GetStatus returns the active status of the validator with the given vID
func (s *state) GetStatus(vID ids.ID) (bool, error) {
	data, exists := s.data[vID]
	if !exists {
		return false, database.ErrNotFound
	}
	return data.IsActive, nil
}

// GetValidationIDs returns the validation IDs in the state
func (s *state) GetValidationIDs() set.Set[ids.ID] {
	ids := set.NewSet[ids.ID](len(s.data))
	for vID := range s.data {
		ids.Add(vID)
	}
	return ids
}

// GetValidatorIDs returns the validator IDs in the state
func (s *state) GetValidatorIDs() set.Set[ids.NodeID] {
	ids := set.NewSet[ids.NodeID](len(s.index))
	for nodeID := range s.index {
		ids.Add(nodeID)
	}
	return ids
}

// RegisterListener registers a listener to the state
// the listener will be notified of current validators via OnValidatorAdded
func (s *state) RegisterListener(listener StateCallbackListener) {
	s.listeners = append(s.listeners, listener)

	// notify the listener of the current state
	for vID, data := range s.data {
		listener.OnValidatorAdded(vID, data.NodeID, uint64(data.startTime.Unix()), data.IsActive)
	}
}

// parseValidatorData parses the data from the bytes into given validatorData
func parseValidatorData(bytes []byte, data *validatorData) error {
	if len(bytes) != 0 {
		if _, err := vdrCodec.Unmarshal(bytes, data); err != nil {
			return err
		}
	}
	data.lastUpdated = time.Unix(int64(data.LastUpdated), 0)
	data.startTime = time.Unix(int64(data.StartTime), 0)
	return nil
}

// Load the state from the disk
func (s *state) loadFromDisk() error {
	it := s.db.NewIterator()
	defer it.Release()
	for it.Next() {
		vIDBytes := it.Key()
		vID, err := ids.ToID(vIDBytes)
		if err != nil {
			return fmt.Errorf("failed to parse validator ID: %w", err)
		}
		vdr := &validatorData{
			validationID: vID,
		}
		if err := parseValidatorData(it.Value(), vdr); err != nil {
			return fmt.Errorf("failed to parse validator data: %w", err)
		}
		if err := s.putData(vID, vdr); err != nil {
			return err
		}
	}
	return it.Error()
}

func (s *state) putData(vID ids.ID, data *validatorData) error {
	if _, exists := s.data[vID]; exists {
		return fmt.Errorf("validator data already exists for %s", vID)
	}
	// should never happen
	if _, exists := s.index[data.NodeID]; exists {
		return fmt.Errorf("validator data already exists for %s", data.NodeID)
	}

	s.data[vID] = data
	s.index[data.NodeID] = vID
	return nil
}

// getData returns the data for the validator with the given nodeID
// returns ErrNotFound if the data does not exist
func (s *state) getData(nodeID ids.NodeID) (*validatorData, error) {
	vID, exists := s.index[nodeID]
	if !exists {
		return nil, database.ErrNotFound
	}
	data, exists := s.data[vID]
	if !exists {
		return nil, database.ErrNotFound
	}
	return data, nil
}
