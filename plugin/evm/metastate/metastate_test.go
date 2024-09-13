// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package metastate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ava-labs/avalanchego/codec"
	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/database/memdb"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/wrappers"
)

func TestValidatorUptimes(t *testing.T) {
	require := require.New(t)
	db := memdb.New()
	state, err := NewValidatorMetaState(db)
	require.NoError(err)

	// get non-existent uptime
	nodeID := ids.GenerateTestNodeID()
	vID := ids.GenerateTestID()
	_, _, err = state.GetUptime(nodeID)
	require.ErrorIs(err, database.ErrNotFound)

	// set non-existent uptime
	startTime := time.Now()
	err = state.SetUptime(nodeID, 1, startTime)
	require.ErrorIs(err, database.ErrNotFound)

	// add new validator metadata
	state.AddNewValidatorMetadata(vID, nodeID, uint64(startTime.Unix()), true)

	// adding the same validator should fail
	err = state.AddNewValidatorMetadata(vID, nodeID, uint64(startTime.Unix()), true)
	require.Error(err)
	// adding the same nodeID should fail
	err = state.AddNewValidatorMetadata(ids.GenerateTestID(), nodeID, uint64(startTime.Unix()), true)
	require.Error(err)

	// get uptime
	upDuration, lastUpdated, err := state.GetUptime(nodeID)
	require.NoError(err)
	require.Equal(time.Duration(0), upDuration)
	require.Equal(startTime.Unix(), lastUpdated.Unix())

	// set uptime
	newUpDuration := 2 * time.Minute
	newLastUpdated := lastUpdated.Add(time.Hour)
	require.NoError(state.SetUptime(nodeID, newUpDuration, newLastUpdated))
	// get new uptime
	upDuration, lastUpdated, err = state.GetUptime(nodeID)
	require.NoError(err)
	require.Equal(newUpDuration, upDuration)
	require.Equal(newLastUpdated, lastUpdated)

	// set status
	require.NoError(state.SetStatus(vID, false))
	// get status
	status, err := state.GetStatus(vID)
	require.NoError(err)
	require.False(status)

	// delete uptime
	state.DeleteValidatorMetadata(vID)

	// get deleted uptime
	_, _, err = state.GetUptime(nodeID)
	require.ErrorIs(err, database.ErrNotFound)
}

func TestWriteValidatorMetadata(t *testing.T) {
	require := require.New(t)
	db := memdb.New()
	state, err := NewValidatorMetaState(db)
	require.NoError(err)
	// write empty uptimes
	require.NoError(state.WriteValidatorMetadata())

	// load uptime
	nodeID := ids.GenerateTestNodeID()
	vID := ids.GenerateTestID()
	startTime := time.Now()
	state.AddNewValidatorMetadata(vID, nodeID, uint64(startTime.Unix()), true)

	// write state, should reflect to DB
	require.NoError(state.WriteValidatorMetadata())
	require.True(db.Has(vID[:]))

	// set uptime
	newUpDuration := 2 * time.Minute
	newLastUpdated := startTime.Add(time.Hour)
	require.NoError(state.SetUptime(nodeID, newUpDuration, newLastUpdated))
	require.NoError(state.WriteValidatorMetadata())

	// refresh state, should load from DB
	state, err = NewValidatorMetaState(db)
	require.NoError(err)

	// get uptime
	upDuration, lastUpdated, err := state.GetUptime(nodeID)
	require.NoError(err)
	require.Equal(newUpDuration, upDuration)
	require.Equal(newLastUpdated.Unix(), lastUpdated.Unix())

	// delete metadata
	state.DeleteValidatorMetadata(vID)

	// write state, should reflect to DB
	require.NoError(state.WriteValidatorMetadata())
	require.False(db.Has(vID[:]))
}

func TestParseValidatorMetadata(t *testing.T) {
	testNodeID, err := ids.NodeIDFromString("NodeID-CaBYJ9kzHvrQFiYWowMkJGAQKGMJqZoat")
	require.NoError(t, err)
	type test struct {
		name        string
		bytes       []byte
		expected    *validatorMetadata
		expectedErr error
	}
	tests := []test{
		{
			name:  "nil",
			bytes: nil,
			expected: &validatorMetadata{
				lastUpdated: time.Unix(0, 0),
				startTime:   time.Unix(0, 0),
			},
			expectedErr: nil,
		},
		{
			name:  "empty",
			bytes: []byte{},
			expected: &validatorMetadata{
				lastUpdated: time.Unix(0, 0),
				startTime:   time.Unix(0, 0),
			},
			expectedErr: nil,
		},
		{
			name: "valid",
			bytes: []byte{
				// codec version
				0x00, 0x00,
				// up duration
				0x00, 0x00, 0x00, 0x00, 0x00, 0x5B, 0x8D, 0x80,
				// last updated
				0x00, 0x00, 0x00, 0x00, 0x00, 0x0D, 0xBB, 0xA0,
				// node ID
				0x7e, 0xef, 0xe8, 0x8a, 0x45, 0xfb, 0x7a, 0xc4,
				0xb0, 0x59, 0xc9, 0x33, 0x71, 0x0a, 0x57, 0x33,
				0xff, 0x9f, 0x4b, 0xab,
				// start time
				0x00, 0x00, 0x00, 0x00, 0x00, 0x5B, 0x8D, 0x80,
				// status
				0x01,
			},
			expected: &validatorMetadata{
				UpDuration:  time.Duration(6000000),
				LastUpdated: 900000,
				lastUpdated: time.Unix(900000, 0),
				NodeID:      testNodeID,
				StartTime:   6000000,
				startTime:   time.Unix(6000000, 0),
				IsActive:    true,
			},
		},
		{
			name: "invalid codec version",
			bytes: []byte{
				// codec version
				0x00, 0x02,
				// up duration
				0x00, 0x00, 0x00, 0x00, 0x00, 0x5B, 0x8D, 0x80,
				// last updated
				0x00, 0x00, 0x00, 0x00, 0x00, 0x0D, 0xBB, 0xA0,
			},
			expected:    nil,
			expectedErr: codec.ErrUnknownVersion,
		},
		{
			name: "short byte len",
			bytes: []byte{
				// codec version
				0x00, 0x00,
				// up duration
				0x00, 0x00, 0x00, 0x00, 0x00, 0x5B, 0x8D, 0x80,
				// last updated
				0x00, 0x00, 0x00, 0x00, 0x00, 0x0D, 0xBB, 0xA0,
			},
			expected:    nil,
			expectedErr: wrappers.ErrInsufficientLength,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			var metadata validatorMetadata
			err := parseValidatorMetadata(tt.bytes, &metadata)
			require.ErrorIs(err, tt.expectedErr)
			if tt.expectedErr != nil {
				return
			}
			require.Equal(tt.expected, &metadata)
		})
	}
}
