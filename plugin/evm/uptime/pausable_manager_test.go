// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package uptime

import (
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/uptime"
	"github.com/ava-labs/avalanchego/utils/timer/mockable"
	"github.com/stretchr/testify/require"
)

func TestPausableManager(t *testing.T) {
	nodeID0 := ids.GenerateTestNodeID()
	startTime := time.Now()
	require := require.New(t)

	// Connect before pause before tracking
	{
		up, clk, _ := setupTestEnv(nodeID0, startTime)

		// Connect before tracking
		require.NoError(up.Connect(nodeID0))
		addTime(clk, time.Second)

		// Pause before tracking
		require.NoError(up.Pause(nodeID0))

		// Elapse Time
		addTime(clk, time.Second)

		// Start tracking
		require.NoError(up.StartTracking([]ids.NodeID{nodeID0}))
		currentTime := addTime(clk, time.Second)
		// Uptime should not have increased since the node was paused
		checkUptime(t, up, nodeID0, 0*time.Second, currentTime)

		// Disconnect
		require.NoError(up.Disconnect(nodeID0))
		// Uptime should not have increased
		checkUptime(t, up, nodeID0, 0*time.Second, currentTime)
	}

	// Paused after tracking resumed after tracking
	{
		up, clk, _ := setupTestEnv(nodeID0, startTime)

		// Start tracking
		require.NoError(up.StartTracking([]ids.NodeID{nodeID0}))

		// Connect
		addTime(clk, time.Second)
		require.NoError(up.Connect(nodeID0))

		// Pause
		addTime(clk, time.Second)
		require.NoError(up.Pause(nodeID0))

		// Elapse time
		currentTime := addTime(clk, 2*time.Second)
		// Uptime should be 1 second since the node was paused after 1 sec
		checkUptime(t, up, nodeID0, 1*time.Second, currentTime)

		// Disconnect and check uptime
		currentTime = addTime(clk, 3*time.Second)
		require.NoError(up.Disconnect(nodeID0))
		// Uptime should not have increased since the node was paused
		checkUptime(t, up, nodeID0, 1*time.Second, currentTime)

		// Connect again and check uptime
		addTime(clk, 4*time.Second)
		require.NoError(up.Connect(nodeID0))
		currentTime = addTime(clk, 5*time.Second)
		// Uptime should not have increased since the node was paused
		checkUptime(t, up, nodeID0, 1*time.Second, currentTime)

		// Resume and check uptime
		currentTime = addTime(clk, 6*time.Second)
		require.NoError(up.Resume(nodeID0))
		// Uptime should not have increased since the node was paused
		// and we just resumed it
		checkUptime(t, up, nodeID0, 1*time.Second, currentTime)

		// Elapsed time check
		currentTime = addTime(clk, 7*time.Second)
		// Uptime should increase by 7 seconds above since the node was resumed
		checkUptime(t, up, nodeID0, 8*time.Second, currentTime)
	}

	// Paused before tracking resumed after tracking
	{
		up, clk, _ := setupTestEnv(nodeID0, startTime)

		// Pause before tracking
		require.NoError(up.Pause(nodeID0))

		// Start tracking
		addTime(clk, time.Second)
		require.NoError(up.StartTracking([]ids.NodeID{nodeID0}))

		// Connect and check uptime
		addTime(clk, 1*time.Second)
		require.NoError(up.Connect(nodeID0))

		currentTime := addTime(clk, 2*time.Second)
		// Uptime should not have increased since the node was paused
		checkUptime(t, up, nodeID0, 0*time.Second, currentTime)

		// Disconnect and check uptime
		currentTime = addTime(clk, 3*time.Second)
		require.NoError(up.Disconnect(nodeID0))
		// Uptime should not have increased since the node was paused
		checkUptime(t, up, nodeID0, 0*time.Second, currentTime)

		// Connect again and resume
		addTime(clk, 4*time.Second)
		require.NoError(up.Connect(nodeID0))
		addTime(clk, 5*time.Second)
		require.NoError(up.Resume(nodeID0))

		// Check uptime after resume
		currentTime = addTime(clk, 6*time.Second)
		// Uptime should have increased by 6 seconds since the node was resumed
		checkUptime(t, up, nodeID0, 6*time.Second, currentTime)
	}

	// Paused after tracking resumed before tracking
	{
		up, clk, s := setupTestEnv(nodeID0, startTime)

		// Start tracking and connect
		require.NoError(up.StartTracking([]ids.NodeID{nodeID0}))
		addTime(clk, time.Second)
		require.NoError(up.Connect(nodeID0))

		// Pause and check uptime
		currentTime := addTime(clk, 2*time.Second)
		require.NoError(up.Pause(nodeID0))
		// Uptime should be 2 seconds since the node was paused after 2 seconds
		checkUptime(t, up, nodeID0, 2*time.Second, currentTime)

		// Stop tracking and reinitialize manager
		currentTime = addTime(clk, 3*time.Second)
		require.NoError(up.StopTracking([]ids.NodeID{nodeID0}))
		up = NewPausableManager(uptime.NewManager(s, clk))

		// Uptime should not have increased since the node was paused
		// and we have not started tracking again
		checkUptime(t, up, nodeID0, 2*time.Second, currentTime)

		// Pause and check uptime
		require.NoError(up.Pause(nodeID0))
		// Uptime should not have increased since the node was paused
		checkUptime(t, up, nodeID0, 2*time.Second, currentTime)

		// Resume and check uptime
		currentTime = addTime(clk, 5*time.Second)
		require.NoError(up.Resume(nodeID0))
		// Uptime should have increased by 5 seconds since the node was resumed
		checkUptime(t, up, nodeID0, 7*time.Second, currentTime)

		// Start tracking and check elapsed time
		currentTime = addTime(clk, 6*time.Second)
		require.NoError(up.StartTracking([]ids.NodeID{nodeID0}))
		// Uptime should have increased by 6 seconds since we started tracking
		// and node was resumed (we assume the node was online until we started tracking)
		checkUptime(t, up, nodeID0, 13*time.Second, currentTime)

		// Elapsed time
		currentTime = addTime(clk, 7*time.Second)
		// Uptime should not have increased since the node was not connected
		checkUptime(t, up, nodeID0, 13*time.Second, currentTime)

		// Connect and final uptime check
		require.NoError(up.Connect(nodeID0))
		currentTime = addTime(clk, 8*time.Second)
		// Uptime should have increased by 8 seconds since the node was connected
		checkUptime(t, up, nodeID0, 21*time.Second, currentTime)
	}
}

func setupTestEnv(nodeID ids.NodeID, startTime time.Time) (PausableManager, *mockable.Clock, uptime.State) {
	clk := mockable.Clock{}
	clk.Set(startTime)
	s := uptime.NewTestState()
	s.AddNode(nodeID, startTime)
	up := NewPausableManager(uptime.NewManager(s, &clk))
	return up, &clk, s
}

func addTime(clk *mockable.Clock, duration time.Duration) time.Time {
	clk.Set(clk.Time().Add(duration))
	return clk.Time()
}

func checkUptime(t *testing.T, up PausableManager, nodeID ids.NodeID, expectedUptime time.Duration, expectedLastUpdate time.Time) {
	uptime, lastUpdated, err := up.CalculateUptime(nodeID)
	require.NoError(t, err)
	require.Equal(t, expectedLastUpdate.Unix(), lastUpdated.Unix())
	require.Equal(t, expectedUptime, uptime)
}
