// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package orderbook

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"math/big"
	"os"

	"github.com/ava-labs/subnet-evm/eth"
	"github.com/ava-labs/subnet-evm/precompile/contracts/bibliophile"
	"github.com/ava-labs/subnet-evm/rpc"
	"github.com/ethereum/go-ethereum/common"
)

type TestingAPI struct {
	db               LimitOrderDatabase
	backend          *eth.EthAPIBackend
	configService    IConfigService
	snapshotFilePath string
}

func NewTestingAPI(database LimitOrderDatabase, backend *eth.EthAPIBackend, configService IConfigService, snapshotFilePath string) *TestingAPI {
	return &TestingAPI{
		db:               database,
		backend:          backend,
		configService:    configService,
		snapshotFilePath: snapshotFilePath,
	}
}

func (api *TestingAPI) GetClearingHouseVars(ctx context.Context, trader common.Address) bibliophile.VariablesReadFromClearingHouseSlots {
	stateDB, _, _ := api.backend.StateAndHeaderByNumber(ctx, rpc.BlockNumber(getCurrentBlockNumber(api.backend)))
	return bibliophile.GetClearingHouseVariables(stateDB, trader)
}

func (api *TestingAPI) GetMarginAccountVars(ctx context.Context, collateralIdx *big.Int, traderAddress string) bibliophile.VariablesReadFromMarginAccountSlots {
	stateDB, _, _ := api.backend.StateAndHeaderByNumber(ctx, rpc.BlockNumber(getCurrentBlockNumber(api.backend)))
	return bibliophile.GetMarginAccountVariables(stateDB, collateralIdx, common.HexToAddress(traderAddress))
}

func (api *TestingAPI) GetAMMVars(ctx context.Context, ammAddress string, ammIndex int, traderAddress string) bibliophile.VariablesReadFromAMMSlots {
	stateDB, _, _ := api.backend.StateAndHeaderByNumber(ctx, rpc.BlockNumber(getCurrentBlockNumber(api.backend)))
	return bibliophile.GetAMMVariables(stateDB, common.HexToAddress(ammAddress), int64(ammIndex), common.HexToAddress(traderAddress))
}

func (api *TestingAPI) GetIOCOrdersVars(ctx context.Context, orderHash common.Hash) bibliophile.VariablesReadFromIOCOrdersSlots {
	stateDB, _, _ := api.backend.StateAndHeaderByNumber(ctx, rpc.BlockNumber(getCurrentBlockNumber(api.backend)))
	return bibliophile.GetIOCOrdersVariables(stateDB, orderHash)
}

func (api *TestingAPI) GetOrderBookVars(ctx context.Context, traderAddress string, senderAddress string, orderHash common.Hash) bibliophile.VariablesReadFromOrderbookSlots {
	stateDB, _, _ := api.backend.StateAndHeaderByNumber(ctx, rpc.BlockNumber(getCurrentBlockNumber(api.backend)))
	return bibliophile.GetOrderBookVariables(stateDB, traderAddress, senderAddress, orderHash)
}

func (api *TestingAPI) GetSnapshot(ctx context.Context) (Snapshot, error) {
	var snapshot Snapshot

	memorySnapshotBytes, err := os.ReadFile(api.snapshotFilePath)
	if err != nil {
		return snapshot, fmt.Errorf("error in fetching snapshot from hubbleDB; err=%v", err)
	}

	buf := bytes.NewBuffer(memorySnapshotBytes)
	err = gob.NewDecoder(buf).Decode(&snapshot)
	if err != nil {
		return snapshot, fmt.Errorf("error in snapshot parsing; err=%v", err)
	}

	return snapshot, nil
}

func getCurrentBlockNumber(backend *eth.EthAPIBackend) uint64 {
	return backend.CurrentHeader().Number.Uint64()
}
