// (c) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/network/p2p/gossip"
	"github.com/ava-labs/subnet-evm/consensus/dummy"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/rawdb"
	"github.com/ava-labs/subnet-evm/core/txpool"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/core/vm"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestGossipEthTxMarshaller(t *testing.T) {
	require := require.New(t)

	blobTx := &types.BlobTx{}
	want := &GossipEthTx{Tx: types.NewTx(blobTx)}
	marshaller := GossipEthTxMarshaller{}

	bytes, err := marshaller.MarshalGossip(want)
	require.NoError(err)

	got, err := marshaller.UnmarshalGossip(bytes)
	require.NoError(err)
	require.Equal(want.GossipID(), got.GossipID())
}

func TestGossipSubscribe(t *testing.T) {
	require := require.New(t)
	key, err := crypto.GenerateKey()
	require.NoError(err)
	addr := crypto.PubkeyToAddress(key.PublicKey)

	require.NoError(err)
	txPool := setupPoolWithConfig(t, params.TestChainConfig, addr)
	defer txPool.Stop()
	txPool.SetGasPrice(common.Big1)
	txPool.SetMinFee(common.Big0)

	gossipTxPool, err := NewGossipEthTxPool(txPool)
	require.NoError(err)

	// use a custom bloom filter to test the bloom filter reset
	gossipTxPool.bloom, err = gossip.NewBloomFilter(1, 0.01, 0.0000000000000001) // maxCount =1
	require.NoError(err)
	ctx, cancel := context.WithCancel(context.TODO())
	go func() {
		gossipTxPool.Subscribe(ctx)
	}()

	// create first eth txes
	ethTxs := getValidEthTxs(key, 10, big.NewInt(226*params.GWei))

	// Notify VM about first eth txs batch
	batch1 := ethTxs[:5]
	errs := txPool.AddRemotesSync(batch1)
	for _, err := range errs {
		require.NoError(err, "failed adding subnet-evm tx to remote mempool")
	}

	// Notify VM about second eth txs batch
	batch2 := ethTxs[5:]
	errs = txPool.AddRemotesSync(batch2)
	for _, err := range errs {
		require.NoError(err, "failed adding subnet-evm tx to remote mempool")
	}
	time.Sleep(50 * time.Millisecond)
	cancel()
	for i, tx := range ethTxs {
		gossipable := &GossipEthTx{Tx: tx}
		require.Truef(gossipTxPool.bloom.Has(gossipable), "expected tx to be in bloom filter: index %d", i)
	}
}

func setupPoolWithConfig(t *testing.T, config *params.ChainConfig, fundedAddress common.Address) *txpool.TxPool {
	diskdb := rawdb.NewMemoryDatabase()
	engine := dummy.NewETHFaker()

	var gspec = &core.Genesis{
		Config: config,
		Alloc:  core.GenesisAlloc{fundedAddress: core.GenesisAccount{Balance: big.NewInt(1000000000000000000)}},
	}
	chain, err := core.NewBlockChain(diskdb, core.DefaultCacheConfig, gspec, engine, vm.Config{}, common.Hash{}, false)
	require.NoError(t, err)
	testTxPoolConfig := txpool.DefaultConfig
	testTxPoolConfig.Journal = ""
	pool := txpool.NewTxPool(testTxPoolConfig, config, chain)

	return pool
}
