// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package metastate

import (
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/uptime"
	"github.com/ava-labs/avalanchego/utils/set"
)

var _ uptime.State = &metastate{}

type ValidatorMetastate interface {
	uptime.State
	// AddNewValidatorMetadata adds a new validator metadata to the metastate
	AddNewValidatorMetadata(vID ids.ID, nodeID ids.NodeID, startTimestamp uint64, isActive bool) error
	// DeleteValidatorMetadata deletes the validator metadata from the metastate
	DeleteValidatorMetadata(vID ids.ID) error
	// WriteValidatorMetadata writes the metastate to the disk
	WriteValidatorMetadata() error

	// SetStatus sets the active status of the validator with the given vID
	SetStatus(vID ids.ID, isActive bool) error
	// GetStatus returns the active status of the validator with the given vID
	GetStatus(vID ids.ID) (bool, error)

	// GetValidationIDs returns the validation IDs in the metastate
	GetValidationIDs() set.Set[ids.ID]
	// GetValidatorIDs returns the validator node IDs in the metastate
	GetValidatorIDs() set.Set[ids.NodeID]
}

type validatorMetadata struct {
	UpDuration  time.Duration `serialize:"true"`
	LastUpdated uint64        `serialize:"true"`
	NodeID      ids.NodeID    `serialize:"true"`
	StartTime   uint64        `serialize:"true"`
	IsActive    bool          `serialize:"true"`

	validationID ids.ID // database key
	lastUpdated  time.Time
	startTime    time.Time
}

type metastate struct {
	data  map[ids.ID]*validatorMetadata // vID -> metadata
	index map[ids.NodeID]ids.ID         // nodeID -> vID
	// updatedMetadata tracks the updates since las WriteValidatorMetadata was called
	updatedMetadata map[ids.ID]bool // vID -> true(updated)/false(deleted)
	db              database.Database
}

// NewValidatorMetaState creates a new ValidatorMetastate, it also loads the metadata from the disk
func NewValidatorMetaState(db database.Database) (ValidatorMetastate, error) {
	m := &metastate{
		index:           make(map[ids.NodeID]ids.ID),
		data:            make(map[ids.ID]*validatorMetadata),
		updatedMetadata: make(map[ids.ID]bool),
		db:              db,
	}
	if err := m.loadFromDisk(); err != nil {
		return nil, fmt.Errorf("failed to load metadata from disk: %w", err)
	}
	return m, nil
}

// GetUptime returns the uptime of the validator with the given nodeID
func (m *metastate) GetUptime(
	nodeID ids.NodeID,
) (time.Duration, time.Time, error) {
	metadata, err := m.getMetadata(nodeID)
	if err != nil {
		return 0, time.Time{}, err
	}
	return metadata.UpDuration, metadata.lastUpdated, nil
}

// SetUptime sets the uptime of the validator with the given nodeID
func (m *metastate) SetUptime(
	nodeID ids.NodeID,
	upDuration time.Duration,
	lastUpdated time.Time,
) error {
	metadata, err := m.getMetadata(nodeID)
	if err != nil {
		return err
	}
	metadata.UpDuration = upDuration
	metadata.lastUpdated = lastUpdated

	m.updatedMetadata[metadata.validationID] = true
	return nil
}

// GetStartTime returns the start time of the validator with the given nodeID
func (m *metastate) GetStartTime(nodeID ids.NodeID) (time.Time, error) {
	metadata, err := m.getMetadata(nodeID)
	if err != nil {
		return time.Time{}, err
	}
	return metadata.startTime, nil
}

// AddNewValidatorMetadata adds a new validator metadata to the metastate
// the new metadata is marked as updated and will be written to the disk when WriteValidatorMetadata is called
func (m *metastate) AddNewValidatorMetadata(vID ids.ID, nodeID ids.NodeID, startTimestamp uint64, isActive bool) error {
	startTimeUnix := time.Unix(int64(startTimestamp), 0)

	metadata := &validatorMetadata{
		NodeID:       nodeID,
		validationID: vID,
		IsActive:     isActive,
		StartTime:    startTimestamp,
		UpDuration:   0,
		LastUpdated:  startTimestamp,
		lastUpdated:  startTimeUnix,
		startTime:    startTimeUnix,
	}
	if err := m.putMetadata(vID, metadata); err != nil {
		return err
	}

	m.updatedMetadata[vID] = true
	return nil
}

// DeleteValidatorMetadata marks the validator metadata as deleted
// marked metadata will be deleted when WriteValidatorMetadata is called
func (m *metastate) DeleteValidatorMetadata(vID ids.ID) error {
	metadata, exists := m.data[vID]
	if !exists {
		return database.ErrNotFound
	}
	delete(m.data, metadata.validationID)
	delete(m.index, metadata.NodeID)

	// mark as deleted for WriteValidatorMetadata
	m.updatedMetadata[metadata.validationID] = false
	return nil
}

// WriteValidatorMetadata writes the updated metastate to the disk
func (m *metastate) WriteValidatorMetadata() error {
	// TODO: add batch size
	batch := m.db.NewBatch()
	for vID, updated := range m.updatedMetadata {
		if updated {
			metadata := m.data[vID]
			metadata.LastUpdated = uint64(metadata.lastUpdated.Unix())
			// should never change but in case
			metadata.StartTime = uint64(metadata.startTime.Unix())

			metadataBytes, err := metadataCodec.Marshal(codecVersion, metadata)
			if err != nil {
				return err
			}
			if err := batch.Put(vID[:], metadataBytes); err != nil {
				return err
			}
		} else { // deleted
			if err := batch.Delete(vID[:]); err != nil {
				return err
			}
		}
		// we're done, remove the updated marker
		delete(m.updatedMetadata, vID)
	}
	return batch.Write()
}

// SetStatus sets the active status of the validator with the given vID
func (m *metastate) SetStatus(vID ids.ID, isActive bool) error {
	metadata, exists := m.data[vID]
	if !exists {
		return database.ErrNotFound
	}
	metadata.IsActive = isActive
	m.updatedMetadata[vID] = true
	return nil
}

// GetStatus returns the active status of the validator with the given vID
func (m *metastate) GetStatus(vID ids.ID) (bool, error) {
	metadata, exists := m.data[vID]
	if !exists {
		return false, database.ErrNotFound
	}
	return metadata.IsActive, nil
}

// GetValidationIDs returns the validation IDs in the metastate
func (m *metastate) GetValidationIDs() set.Set[ids.ID] {
	ids := set.NewSet[ids.ID](len(m.data))
	for vID := range m.data {
		ids.Add(vID)
	}
	return ids
}

// GetValidatorIDs returns the validator IDs in the metastate
func (m *metastate) GetValidatorIDs() set.Set[ids.NodeID] {
	ids := set.NewSet[ids.NodeID](len(m.index))
	for nodeID := range m.index {
		ids.Add(nodeID)
	}
	return ids
}

// parseValidatorMetadata parses the metadata from the bytes and returns the metadata
func parseValidatorMetadata(bytes []byte, metadata *validatorMetadata) error {
	if len(bytes) != 0 {
		if _, err := metadataCodec.Unmarshal(bytes, metadata); err != nil {
			return err
		}
	}
	metadata.lastUpdated = time.Unix(int64(metadata.LastUpdated), 0)
	metadata.startTime = time.Unix(int64(metadata.StartTime), 0)
	return nil
}

// Load the metadata from the disk
func (m *metastate) loadFromDisk() error {
	it := m.db.NewIterator()
	defer it.Release()
	for it.Next() {
		vIDBytes := it.Key()
		vID, err := ids.ToID(vIDBytes)
		if err != nil {
			return fmt.Errorf("failed to parse validator ID: %w", err)
		}
		metadata := &validatorMetadata{
			validationID: vID,
		}
		if err := parseValidatorMetadata(it.Value(), metadata); err != nil {
			return fmt.Errorf("failed to parse validator metadata: %w", err)
		}
		if err := m.putMetadata(vID, metadata); err != nil {
			return err
		}
	}
	return it.Error()
}

func (m *metastate) putMetadata(vID ids.ID, metadata *validatorMetadata) error {
	if _, exists := m.data[vID]; exists {
		return fmt.Errorf("validator metadata already exists for %s", vID)
	}
	// should never happen
	if _, exists := m.index[metadata.NodeID]; exists {
		return fmt.Errorf("validator metadata already exists for %s", metadata.NodeID)
	}

	m.data[vID] = metadata
	m.index[metadata.NodeID] = vID
	return nil
}

// getMetadata returns the metadata for the validator with the given nodeID
// returns ErrNotFound if the metadata does not exist
func (m *metastate) getMetadata(nodeID ids.NodeID) (*validatorMetadata, error) {
	vID, exists := m.index[nodeID]
	if !exists {
		return nil, database.ErrNotFound
	}
	metadata, exists := m.data[vID]
	if !exists {
		return nil, database.ErrNotFound
	}
	return metadata, nil
}
