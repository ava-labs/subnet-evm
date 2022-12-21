// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"math/big"
	"testing"

	"github.com/ava-labs/avalanchego/snow"
	engCommon "github.com/ava-labs/avalanchego/snow/engine/common"
	"github.com/ava-labs/avalanchego/vms/components/avax"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/utils/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestSharedMemory(t *testing.T) {
	genesis := &core.Genesis{}
	require.NoError(t, genesis.UnmarshalJSON([]byte(genesisJSONSubnetEVM)))

	genesis.Config.SharedMemoryConfig = precompile.NewSharedMemoryConfig(common.Big0)
	genesisJSON, err := genesis.MarshalJSON()
	require.NoError(t, err)

	vm := &VM{}
	ctx, dbManager, genesisBytes, issuer, m := setupGenesis(t, string(genesisJSON))
	appSender := &engCommon.SenderTest{T: t}
	appSender.CantSendAppGossip = true
	appSender.SendAppGossipF = func(context.Context, []byte) error { return nil }
	if err := vm.Initialize(
		context.Background(),
		ctx,
		dbManager,
		genesisBytes,
		nil,
		nil,
		issuer,
		[]*engCommon.Fx{},
		appSender,
	); err != nil {
		t.Fatal(err)
	}

	require.NoError(t, vm.SetState(context.Background(), snow.Bootstrapping))
	require.NoError(t, vm.SetState(context.Background(), snow.NormalOp))

	defer func() {
		require.NoError(t, vm.Shutdown(context.Background()))
	}()

	newTxPoolHeadChan := make(chan core.NewTxPoolReorgEvent, 1)
	vm.txPool.SubscribeNewReorgEvent(newTxPoolHeadChan)

	data, err := precompile.PackExportAVAX(testXChainID, uint64(0), uint64(1), []common.Address{testEthAddrs[0]})
	require.NoError(t, err)

	tx := types.NewTransaction(uint64(0), precompile.SharedMemoryAddress, big.NewInt(1), 200_000, big.NewInt(testMinGasPrice), data)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(vm.chainConfig.ChainID), testKeys[0])
	require.NoError(t, err)

	txErrors := vm.txPool.AddRemotesSync([]*types.Transaction{signedTx})
	for _, err := range txErrors {
		require.NoError(t, err)
	}

	blk := issueAndAccept(t, issuer, vm)
	newHead := <-newTxPoolHeadChan
	require.Equal(t, newHead.Head.Hash(), common.Hash(blk.ID()))

	// Drain the acceptor queue so that we finish processing the atomic operations
	vm.blockChain.DrainAcceptorQueue()

	xChainSharedMemory := m.NewSharedMemory(testXChainID)
	values, _, _, err := xChainSharedMemory.Indexed(ctx.ChainID, [][]byte{testEthAddrs[0][:]}, nil, nil, 100)
	require.NoError(t, err)
	require.Len(t, values, 1)
	utxo := &avax.UTXO{}
	v, err := codec.Codec.Unmarshal(values[0], utxo)
	require.NoError(t, err)
	require.Equal(t, uint16(0), v)
}
