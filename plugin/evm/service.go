// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"math/big"
	"time"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/json"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/subnet-evm/plugin/evm/validators"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

// SnowmanAPI introduces snowman specific functionality to the evm
type SnowmanAPI struct{ vm *VM }

// GetAcceptedFrontReply defines the reply that will be sent from the
// GetAcceptedFront API call
type GetAcceptedFrontReply struct {
	Hash   common.Hash `json:"hash"`
	Number *big.Int    `json:"number"`
}

// GetAcceptedFront returns the last accepted block's hash and height
func (api *SnowmanAPI) GetAcceptedFront(ctx context.Context) (*GetAcceptedFrontReply, error) {
	blk := api.vm.blockChain.LastConsensusAcceptedBlock()
	return &GetAcceptedFrontReply{
		Hash:   blk.Hash(),
		Number: blk.Number(),
	}, nil
}

// IssueBlock to the chain
func (api *SnowmanAPI) IssueBlock(ctx context.Context) error {
	log.Info("Issuing a new block")
	api.vm.builder.signalTxsReady()
	return nil
}

type ValidatorsAPI struct {
	vm *VM
}

type GetCurrentValidatorResponse struct {
	ValidationID     ids.ID        `json:"validationID"`
	NodeID           ids.NodeID    `json:"nodeID"`
	StartTime        time.Time     `json:"startTime"`
	IsActive         bool          `json:"isActive"`
	IsConnected      bool          `json:"isConnected"`
	UptimePercentage *json.Float32 `json:"uptimePercentage"`
	Uptime           time.Duration `json:"uptime"`
}

// GetUptime returns the uptime of the node
func (api *ValidatorsAPI) GetCurrentValidators(ctx context.Context, nodeIDsArg *[]ids.NodeID) ([]GetCurrentValidatorResponse, error) {
	api.vm.ctx.Lock.Lock()
	defer api.vm.ctx.Lock.Unlock()
	var nodeIDs set.Set[ids.NodeID]
	if nodeIDsArg == nil || len(*nodeIDsArg) == 0 {
		nodeIDs = api.vm.validatorState.GetValidatorIDs()
	} else {
		nodeIDs = set.Of(*nodeIDsArg...)
	}

	responses := make([]GetCurrentValidatorResponse, 0, nodeIDs.Len())

	for _, nodeID := range nodeIDs.List() {
		validator, err := api.vm.validatorState.GetValidator(nodeID)
		switch {
		case err == database.ErrNotFound:
			continue
		case err != nil:
			return nil, err
		}
		uptimePerc, err := api.getAPIUptimePerc(validator)
		if err != nil {
			return nil, err
		}
		isConnected := api.vm.uptimeManager.IsConnected(nodeID)

		uptime, _, err := api.vm.uptimeManager.CalculateUptime(nodeID)
		if err != nil {
			return nil, err
		}

		responses = append(responses, GetCurrentValidatorResponse{
			ValidationID:     validator.ValidationID,
			NodeID:           nodeID,
			StartTime:        validator.StartTime,
			IsActive:         validator.IsActive,
			UptimePercentage: uptimePerc,
			IsConnected:      isConnected,
			Uptime:           time.Duration(uptime.Seconds()),
		})
	}
	return responses, nil
}

func (api *ValidatorsAPI) getAPIUptimePerc(validator *validators.ValidatorOutput) (*json.Float32, error) {
	rawUptime, err := api.vm.uptimeManager.CalculateUptimePercentFrom(validator.NodeID, validator.StartTime)
	if err != nil {
		return nil, err
	}
	// Transform this to a percentage (0-100)
	uptime := json.Float32(rawUptime * 100)
	return &uptime, nil
}
