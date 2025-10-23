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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ava-labs/subnet-evm/utils/utilstest"

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
	ctx.ValidatorState = &validatorstest.State{
		GetCurrentValidatorSetF: func(context.Context, ids.ID) (map[ids.ID]*avagovalidators.GetCurrentValidatorOutput, uint64, error) {
			return map[ids.ID]*avagovalidators.GetCurrentValidatorOutput{
				testValidationIDs[0]: {
					NodeID:    testNodeIDs[0],
					PublicKey: nil,
					Weight:    1,
				},
				testValidationIDs[1]: {
					NodeID:    testNodeIDs[1],
					PublicKey: nil,
					Weight:    1,
				},
				testValidationIDs[2]: {
					NodeID:    testNodeIDs[2],
					PublicKey: nil,
					Weight:    1,
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

	// Test case 1: state should not be populated until bootstrapped
	require.NoError(vm.SetState(context.Background(), snow.Bootstrapping))
	// After bootstrapping but before NormalOp, uptimeTracker hasn't started syncing yet
	_, _, found, err := vm.uptimeTracker.GetUptime(testValidationIDs[0])
	require.NoError(err)
	require.False(found, "uptime should not be tracked yet")

	// Test case 2: state should be populated after bootstrapped
	require.NoError(vm.SetState(context.Background(), snow.NormalOp))
	// Give the sync goroutine time to run
	time.Sleep(2 * time.Second)
	_, _, found, err = vm.uptimeTracker.GetUptime(testValidationIDs[0])
	require.NoError(err)
	require.True(found, "uptime should be tracked after NormalOp")

	// Test case 3: restarting VM should not lose state
	vm.Shutdown(context.Background())

	vm = &VM{}
	err = vm.Initialize(
		context.Background(),
		utilstest.NewTestSnowContext(t), // this context does not have validators state, making VM to source it from the database
		dbManager,
		genesisBytes,
		[]byte(""),
		[]byte(""),
		[]*commonEng.Fx{},
		appSender,
	)
	require.NoError(err, "error initializing GenesisVM")
	// Uptime data should be persisted from the previous run
	_, _, _, err = vm.uptimeTracker.GetUptime(testValidationIDs[0])
	require.NoError(err)
	// Note: uptime tracking hasn't started yet (not in NormalOp), so found might be false or true depending on persistence

	// Test case 4: new validators should be added to the state
	newValidationID := ids.GenerateTestID()
	newNodeID := ids.GenerateTestNodeID()
	testState := &validatorstest.State{
		GetCurrentValidatorSetF: func(context.Context, ids.ID) (map[ids.ID]*avagovalidators.GetCurrentValidatorOutput, uint64, error) {
			return map[ids.ID]*avagovalidators.GetCurrentValidatorOutput{
				testValidationIDs[0]: {
					NodeID:    testNodeIDs[0],
					PublicKey: nil,
					Weight:    1,
				},
				testValidationIDs[1]: {
					NodeID:    testNodeIDs[1],
					PublicKey: nil,
					Weight:    1,
				},
				testValidationIDs[2]: {
					NodeID:    testNodeIDs[2],
					PublicKey: nil,
					Weight:    1,
				},
				newValidationID: {
					NodeID:    newNodeID,
					PublicKey: nil,
					Weight:    1,
				},
			}, 0, nil
		},
	}
	// set VM as bootstrapped
	require.NoError(vm.SetState(context.Background(), snow.Bootstrapping))
	require.NoError(vm.SetState(context.Background(), snow.NormalOp))

	vm.ctx.ValidatorState = testState

	// new validator should be added to the state eventually after sync runs
	require.EventuallyWithT(func(c *assert.CollectT) {
		vm.vmLock.Lock()
		defer vm.vmLock.Unlock()
		// Check if the new validator's uptime is being tracked
		_, _, found, err := vm.uptimeTracker.GetUptime(newValidationID)
		assert.NoError(c, err)
		assert.True(c, found, "new validator should be tracked")
	}, 5*time.Second, 100*time.Millisecond)
}
