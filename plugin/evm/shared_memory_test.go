// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"math/big"
	"testing"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	engCommon "github.com/ava-labs/avalanchego/snow/engine/common"
	"github.com/ava-labs/avalanchego/vms/components/avax"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contracts/sharedmemory"
	"github.com/ava-labs/subnet-evm/utils/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestSharedMemory(t *testing.T) {
	// create genesis with the shared memory precompile enabled
	genesis := &core.Genesis{}
	require.NoError(t, genesis.UnmarshalJSON([]byte(genesisJSONSubnetEVM)))
	genesis.Config.GenesisPrecompiles = params.Precompiles{
		sharedmemory.ConfigKey: sharedmemory.NewConfig(common.Big0),
	}
	genesisJSON, err := genesis.MarshalJSON()
	require.NoError(t, err)

	ctx, dbManager, genesisBytes, issuer, atomicMemory := setupGenesis(t, string(genesisJSON))

	// Find the X Chain's shared memory
	xChainSharedMemory := atomicMemory.NewSharedMemory(testXChainID)

	// initialize and configure the VM
	vm := &VM{}
	appSender := &engCommon.SenderTest{T: t}
	appSender.CantSendAppGossip = true
	appSender.SendAppGossipF = func(context.Context, []byte) error { return nil }
	err = vm.Initialize(
		context.Background(),
		ctx,
		dbManager,
		genesisBytes,
		nil,
		nil,
		issuer,
		[]*engCommon.Fx{},
		appSender,
	)
	require.NoError(t, err)
	defer func() { require.NoError(t, vm.Shutdown(context.Background())) }()

	require.NoError(t, vm.SetState(context.Background(), snow.Bootstrapping))
	require.NoError(t, vm.SetState(context.Background(), snow.NormalOp))

	newTxPoolHeadChan := make(chan core.NewTxPoolReorgEvent, 1)
	vm.txPool.SubscribeNewReorgEvent(newTxPoolHeadChan)

	// Note: the addresses specified here must correspond to the account that should control the exported funds on the recipient chain
	// If this were to be imported to the X-Chain, we would need to use an address derived from an X-Chain private key
	data, err := sharedmemory.PackExportAVAX(sharedmemory.ExportAVAXInput{
		DestinationChainID: testXChainID,
		Locktime:           uint64(0),
		Threshold:          uint64(1),
		Addrs:              []common.Address{testEthAddrs[1]},
	})
	require.NoError(t, err)

	tx := types.NewTransaction(uint64(0), sharedmemory.ContractAddress, big.NewInt(params.Ether), 200_000, big.NewInt(testMinGasPrice), data)
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

	values, _, _, err := xChainSharedMemory.Indexed(ctx.ChainID, [][]byte{testEthAddrs[1][:]}, nil, nil, 100)
	require.NoError(t, err)
	require.Len(t, values, 1)
	utxo := &avax.UTXO{}
	v, err := codec.Codec.Unmarshal(values[0], utxo)
	require.NoError(t, err)
	require.Equal(t, uint16(0), v)

	expectedUTXO := &avax.UTXO{
		// Derive unique UTXOID from txHash and log index
		UTXOID: avax.UTXOID{
			TxID:        ids.ID(signedTx.Hash()),
			OutputIndex: uint32(0),
		},
		Asset: avax.Asset{ID: ctx.AVAXAssetID},
		Out: &secp256k1fx.TransferOutput{
			Amt: 1_000_000_000,
			OutputOwners: secp256k1fx.OutputOwners{
				Locktime:  0,
				Threshold: uint32(1),
				Addrs:     []ids.ShortID{ids.ShortID(testEthAddrs[1])},
			},
		},
	}

	expectedUTXOBytes, err := codec.Codec.Marshal(uint16(0), expectedUTXO)
	require.NoError(t, err)
	require.Equal(t, expectedUTXOBytes, values[0])
}
