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

	makeValidatorState := func() *validatorstest.State {
		return &validatorstest.State{
			GetMinimumHeightF: func(context.Context) (uint64, error) {
				return 0, nil
			},
			GetCurrentHeightF: func(context.Context) (uint64, error) {
				return 0, nil
			},
			GetSubnetIDF: func(context.Context, ids.ID) (ids.ID, error) {
				return ids.Empty, nil
			},
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
	}

	t.Run("uptime not tracked before NormalOp", func(t *testing.T) {
		require := require.New(t)
		ctx, dbManager, genesisBytes := setupGenesis(t, upgradetest.Latest)
		ctx.ValidatorState = makeValidatorState()

		appSender := &enginetest.Sender{T: t}
		appSender.CantSendAppGossip = true
		appSender.SendAppGossipF = func(context.Context, commonEng.SendConfig, []byte) error { return nil }

		vm := &VM{}
		require.NoError(vm.Initialize(
			context.Background(),
			ctx,
			dbManager,
			genesisBytes,
			[]byte(""),
			[]byte(""),
			[]*commonEng.Fx{},
			appSender,
		))
		defer vm.Shutdown(context.Background())

		require.NoError(vm.SetState(context.Background(), snow.Bootstrapping))

		// After bootstrapping but before NormalOp, uptimeTracker hasn't started syncing yet
		_, _, found, err := vm.uptimeTracker.GetUptime(testValidationIDs[0])
		require.NoError(err)
		require.False(found, "uptime should not be tracked yet")
	})

	t.Run("uptime tracked after NormalOp and Connect", func(t *testing.T) {
		require := require.New(t)
		ctx, dbManager, genesisBytes := setupGenesis(t, upgradetest.Latest)
		ctx.ValidatorState = makeValidatorState()

		appSender := &enginetest.Sender{T: t}
		appSender.CantSendAppGossip = true
		appSender.SendAppGossipF = func(context.Context, commonEng.SendConfig, []byte) error { return nil }

		vm := &VM{}
		require.NoError(vm.Initialize(
			context.Background(),
			ctx,
			dbManager,
			genesisBytes,
			[]byte(""),
			[]byte(""),
			[]*commonEng.Fx{},
			appSender,
		))
		defer vm.Shutdown(context.Background())

		require.NoError(vm.SetState(context.Background(), snow.Bootstrapping))
		require.NoError(vm.SetState(context.Background(), snow.NormalOp))

		// Connect the validators to start tracking their uptime
		for _, nodeID := range testNodeIDs {
			require.NoError(vm.uptimeTracker.Connect(nodeID))
		}

		// Manually call Sync to ensure state is updated after Connect
		require.NoError(vm.uptimeTracker.Sync(context.Background()))

		_, _, found, err := vm.uptimeTracker.GetUptime(testValidationIDs[0])
		require.NoError(err)
		require.True(found, "uptime should be tracked after validators are connected")
	})
}
