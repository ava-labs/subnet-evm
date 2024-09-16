// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"math/big"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/uptime"
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

type UptimeAPI struct {
	ctx        *snow.Context
	calculator uptime.LockedCalculator
}

// TODO: add StartTime
type GetUptimeResponse struct {
	NodeID ids.NodeID    `json:"nodeID"`
	Uptime time.Duration `json:"uptime"`
}

// GetUptime returns the uptime of the node
func (api *UptimeAPI) GetUptime(ctx context.Context, nodeID ids.NodeID) (*GetUptimeResponse, error) {
	uptime, _, err := api.calculator.CalculateUptime(nodeID)
	if err != nil {
		return nil, err
	}

	return &GetUptimeResponse{
		NodeID: nodeID,
		Uptime: uptime,
	}, nil
}

// TODO: add GetUptime for currently tracked peers
