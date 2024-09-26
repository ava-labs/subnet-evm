// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"net/http"
	"time"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/set"
)

type ValidatorsAPI struct {
	vm *VM
}

type GetCurrentValidatorsRequest struct {
	NodeIDs []ids.NodeID `json:"nodeIDs"`
}

type GetCurrentValidatorsResponse struct {
	Validators []CurrentValidator `json:"validators"`
}

type CurrentValidator struct {
	ValidationID ids.ID        `json:"validationID"`
	NodeID       ids.NodeID    `json:"nodeID"`
	StartTime    time.Time     `json:"startTime"`
	IsActive     bool          `json:"isActive"`
	IsConnected  bool          `json:"isConnected"`
	Uptime       time.Duration `json:"uptime"`
}

// GetUptime returns the uptime of the node
func (api *ValidatorsAPI) GetCurrentValidators(_ *http.Request, args *GetCurrentValidatorsRequest, reply *GetCurrentValidatorsResponse) error {
	api.vm.ctx.Lock.RLock()
	defer api.vm.ctx.Lock.RUnlock()

	nodeIDs := set.Of(args.NodeIDs...)
	if nodeIDs.Len() == 0 {
		nodeIDs = api.vm.validatorState.GetValidatorIDs()
	}

	reply.Validators = make([]CurrentValidator, 0, nodeIDs.Len())

	for _, nodeID := range nodeIDs.List() {
		validator, err := api.vm.validatorState.GetValidator(nodeID)
		switch {
		case err == database.ErrNotFound:
			continue
		case err != nil:
			return err
		}

		isConnected := api.vm.uptimeManager.IsConnected(nodeID)

		uptime, _, err := api.vm.uptimeManager.CalculateUptime(nodeID)
		if err != nil {
			return err
		}

		reply.Validators = append(reply.Validators, CurrentValidator{
			ValidationID: validator.ValidationID,
			NodeID:       nodeID,
			StartTime:    validator.StartTime,
			IsActive:     validator.IsActive,
			IsConnected:  isConnected,
			Uptime:       time.Duration(uptime.Seconds()),
		})
	}
	return nil
}
