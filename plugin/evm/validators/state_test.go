// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package validators

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

func TestState(t *testing.T) {
	require := require.New(t)
	db := memdb.New()
	state, err := NewState(db)
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

	// add new validator
	state.AddNewValidator(vID, nodeID, uint64(startTime.Unix()), true)

	// adding the same validator should fail
	err = state.AddNewValidator(vID, nodeID, uint64(startTime.Unix()), true)
	require.Error(err)
	// adding the same nodeID should fail
	err = state.AddNewValidator(ids.GenerateTestID(), nodeID, uint64(startTime.Unix()), true)
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
	state.DeleteValidator(vID)

	// get deleted uptime
	_, _, err = state.GetUptime(nodeID)
	require.ErrorIs(err, database.ErrNotFound)
}

func TestWriteValidator(t *testing.T) {
	require := require.New(t)
	db := memdb.New()
	state, err := NewState(db)
	require.NoError(err)
	// write empty uptimes
	require.NoError(state.WriteState())

	// load uptime
	nodeID := ids.GenerateTestNodeID()
	vID := ids.GenerateTestID()
	startTime := time.Now()
	state.AddNewValidator(vID, nodeID, uint64(startTime.Unix()), true)

	// write state, should reflect to DB
	require.NoError(state.WriteState())
	require.True(db.Has(vID[:]))

	// set uptime
	newUpDuration := 2 * time.Minute
	newLastUpdated := startTime.Add(time.Hour)
	require.NoError(state.SetUptime(nodeID, newUpDuration, newLastUpdated))
	require.NoError(state.WriteState())

	// refresh state, should load from DB
	state, err = NewState(db)
	require.NoError(err)

	// get uptime
	upDuration, lastUpdated, err := state.GetUptime(nodeID)
	require.NoError(err)
	require.Equal(newUpDuration, upDuration)
	require.Equal(newLastUpdated.Unix(), lastUpdated.Unix())

	// delete
	state.DeleteValidator(vID)

	// write state, should reflect to DB
	require.NoError(state.WriteState())
	require.False(db.Has(vID[:]))
}

func TestParseValidator(t *testing.T) {
	testNodeID, err := ids.NodeIDFromString("NodeID-CaBYJ9kzHvrQFiYWowMkJGAQKGMJqZoat")
	require.NoError(t, err)
	type test struct {
		name        string
		bytes       []byte
		expected    *validatorData
		expectedErr error
	}
	tests := []test{
		{
			name:  "nil",
			bytes: nil,
			expected: &validatorData{
				lastUpdated: time.Unix(0, 0),
				startTime:   time.Unix(0, 0),
			},
			expectedErr: nil,
		},
		{
			name:  "empty",
			bytes: []byte{},
			expected: &validatorData{
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
			expected: &validatorData{
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
			var data validatorData
			err := parseValidatorData(tt.bytes, &data)
			require.ErrorIs(err, tt.expectedErr)
			if tt.expectedErr != nil {
				return
			}
			require.Equal(tt.expected, &data)
		})
	}
}

func TestStateListener(t *testing.T) {
	require := require.New(t)
	db := memdb.New()
	state, err := NewState(db)
	require.NoError(err)

	expectedvID := ids.GenerateTestID()
	expectedNodeID := ids.GenerateTestNodeID()
	expectedStartTime := time.Now()

	// add listener
	listener := &testCallbackListener{
		t: t,
		onAdd: func(vID ids.ID, nodeID ids.NodeID, startTime uint64, isActive bool) {
			require.Equal(expectedvID, vID)
			require.Equal(expectedNodeID, nodeID)
			require.Equal(uint64(expectedStartTime.Unix()), startTime)
			require.True(isActive)
		},
		onRemove: func(vID ids.ID, nodeID ids.NodeID) {
			require.Equal(expectedvID, vID)
			require.Equal(expectedNodeID, nodeID)
		},
		onStatusUpdate: func(vID ids.ID, nodeID ids.NodeID, isActive bool) {
			require.Equal(expectedvID, vID)
			require.Equal(expectedNodeID, nodeID)
			require.False(isActive)
		},
	}
	state.RegisterListener(listener)

	// add new validator
	state.AddNewValidator(expectedvID, expectedNodeID, uint64(expectedStartTime.Unix()), true)

	// set status
	require.NoError(state.SetStatus(expectedvID, false))

	// remove validator
	state.DeleteValidator(expectedvID)

	require.Equal(3, listener.called)

	// test case: check initial trigger when registering listener
	// add new validator
	state.AddNewValidator(expectedvID, expectedNodeID, uint64(expectedStartTime.Unix()), true)
	newListener := &testCallbackListener{
		t: t,
		onAdd: func(vID ids.ID, nodeID ids.NodeID, startTime uint64, isActive bool) {
			require.Equal(expectedvID, vID)
			require.Equal(expectedNodeID, nodeID)
			require.Equal(uint64(expectedStartTime.Unix()), startTime)
			require.True(isActive)
		},
	}
	state.RegisterListener(newListener)
	require.Equal(1, newListener.called)
}

var _ StateCallbackListener = (*testCallbackListener)(nil)

type testCallbackListener struct {
	t              *testing.T
	called         int
	onAdd          func(vID ids.ID, nodeID ids.NodeID, startTime uint64, isActive bool)
	onRemove       func(ids.ID, ids.NodeID)
	onStatusUpdate func(ids.ID, ids.NodeID, bool)
}

func (t *testCallbackListener) OnValidatorAdded(vID ids.ID, nodeID ids.NodeID, startTime uint64, isActive bool) {
	t.called++
	if t.onAdd != nil {
		t.onAdd(vID, nodeID, startTime, isActive)
	} else {
		t.t.Fail()
	}
}

func (t *testCallbackListener) OnValidatorRemoved(vID ids.ID, nodeID ids.NodeID) {
	t.called++
	if t.onRemove != nil {
		t.onRemove(vID, nodeID)
	} else {
		t.t.Fail()
	}
}

func (t *testCallbackListener) OnValidatorStatusUpdated(vID ids.ID, nodeID ids.NodeID, isActive bool) {
	t.called++
	if t.onStatusUpdate != nil {
		t.onStatusUpdate(vID, nodeID, isActive)
	} else {
		t.t.Fail()
	}
}
