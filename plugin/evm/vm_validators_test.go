// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/engine/enginetest"
	"github.com/ava-labs/avalanchego/snow/validators/validatorstest"
	"github.com/ava-labs/avalanchego/upgrade/upgradetest"
	"github.com/stretchr/testify/require"

	commonEng "github.com/ava-labs/avalanchego/snow/engine/common"
	avagovalidators "github.com/ava-labs/avalanchego/snow/validators"
)

func TestValidatorState(t *testing.T) {
	require := require.New(t)
	ctx, dbManager, genesisBytes := setupGenesis(t, upgradetest.Latest)

	vm := &VM{}

	appSender := &enginetest.Sender{T: t}
	appSender.CantSendAppGossip = true
	testNodeIDs := []ids.NodeID{
		ids.GenerateTestNodeID(),
		ids.GenerateTestNodeID(),
		ids.GenerateTestNodeID(),
	}
	testValidationIDs := []ids.ID{
		ids.GenerateTestID(),
		ids.GenerateTestID(),
		ids.GenerateTestID(),
	}
	startTime := uint64(time.Now().Unix())
	ctx.ValidatorState = &validatorstest.State{
		GetCurrentValidatorSetF: func(context.Context, ids.ID) (map[ids.ID]*avagovalidators.GetCurrentValidatorOutput, uint64, error) {
			return map[ids.ID]*avagovalidators.GetCurrentValidatorOutput{
				testValidationIDs[0]: {
					NodeID:    testNodeIDs[0],
					PublicKey: nil,
					Weight:    1,
					StartTime: startTime,
					IsActive:  true,
				},
				testValidationIDs[1]: {
					NodeID:    testNodeIDs[1],
					PublicKey: nil,
					Weight:    1,
					StartTime: startTime,
					IsActive:  true,
				},
				testValidationIDs[2]: {
					NodeID:    testNodeIDs[2],
					PublicKey: nil,
					Weight:    1,
					StartTime: startTime,
					IsActive:  true,
				},
			}, 0, nil
		},
	}
	appSender.SendAppGossipF = func(context.Context, commonEng.SendConfig, []byte) error { return nil }
	err := vm.Initialize(
		context.Background(),
		ctx,
		dbManager,
		genesisBytes,
		[]byte(""),
		[]byte(""),
		[]*commonEng.Fx{},
		appSender,
	)
	require.NoError(err, "error initializing GenesisVM")

	// Test case 1: uptime should not be tracked until NormalOp
	require.NoError(vm.SetState(context.Background(), snow.Bootstrapping))
	// After bootstrapping but before NormalOp, uptimeTracker hasn't started syncing yet
	_, _, found, err := vm.uptimeTracker.GetUptime(testValidationIDs[0])
	require.NoError(err)
	require.False(found, "uptime should not be tracked yet")

	// Test case 2: uptime should be tracked after NormalOp
	require.NoError(vm.SetState(context.Background(), snow.NormalOp))
	// Give the sync goroutine time to run at least once
	time.Sleep(2 * time.Second)
	_, _, found, err = vm.uptimeTracker.GetUptime(testValidationIDs[0])
	require.NoError(err)
	require.True(found, "uptime should be tracked after NormalOp")

	// Test case 3: uptime data should be persisted across restarts
	require.NoError(vm.Shutdown(context.Background()))

	// Create a new context for the restarted VM to avoid metric conflicts
	ctx2, _, _ := setupGenesis(t, upgradetest.Latest)
	ctx2.ValidatorState = ctx.ValidatorState // Reuse the same validator state

	vm = &VM{}
	err = vm.Initialize(
		context.Background(),
		ctx2,
		dbManager,
		genesisBytes,
		[]byte(""),
		[]byte(""),
		[]*commonEng.Fx{},
		appSender,
	)
	require.NoError(err, "error initializing GenesisVM after restart")

	// Uptime data should be persisted from the previous run - the state is \
	// persisted in the database and will be loaded on initialization
	require.NoError(vm.SetState(context.Background(), snow.Bootstrapping))
	require.NoError(vm.SetState(context.Background(), snow.NormalOp))
	time.Sleep(2 * time.Second)

	_, _, found, err = vm.uptimeTracker.GetUptime(testValidationIDs[0])
	require.NoError(err)
	require.True(found, "uptime should be tracked after restart and NormalOp")

	require.NoError(vm.Shutdown(context.Background()))
}
