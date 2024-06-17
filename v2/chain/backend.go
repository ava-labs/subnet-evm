package chain

import (
	"context"
	"math/big"

	"github.com/ava-labs/avalanchego/utils/timer/mockable"
	"github.com/ava-labs/coreth/consensus"
	"github.com/ava-labs/coreth/core/txpool"
	"github.com/ava-labs/coreth/core/txpool/legacypool"
	"github.com/ava-labs/coreth/core/types"
	"github.com/ava-labs/coreth/miner"
	"github.com/ava-labs/coreth/rpc"
	"github.com/ethereum/go-ethereum/common"
)

// legacyBackend attaches legacy backend components (txPool, miner) to the chain
// so it can be used by the plugin interface.
type legacyBackend struct {
	chain  BlockChain
	txPool TxPool
	miner  *miner.Miner
	engine consensus.Engine
}

func NewLegacyBackend(
	chain BlockChain,
	poolConfig legacypool.Config,
	minerConfig *miner.Config,
	clock *mockable.Clock,
) (*legacyBackend, error) {
	legacyPool := legacypool.New(poolConfig, chain)
	txPool, err := txpool.New(new(big.Int).SetUint64(poolConfig.PriceLimit), chain, []txpool.SubPool{legacyPool}) // Note: blobpool omitted
	if err != nil {
		return nil, err
	}

	engine := chain.Engine()
	miner := miner.New(chain, txPool, minerConfig, chain.Config(), engine, clock)
	return &legacyBackend{
		chain:  chain,
		txPool: txPool,
		miner:  miner,
		engine: engine,
	}, nil
}

func (b *legacyBackend) Start() {}
func (b *legacyBackend) Stop() error {
	b.txPool.Close()
	b.chain.Stop()
	b.engine.Close()
	return nil
}

func (b *legacyBackend) EstimateBaseFee(context.Context) (*big.Int, error) { panic("unimplemented") }
func (b *legacyBackend) SetEtherbase(common.Address)                       { panic("unimplemented") }
func (b *legacyBackend) ResetToStateSyncedBlock(*types.Block) error {
	panic("unimplemented")
}

func (b *legacyBackend) BlockChain() BlockChain { return b.chain }
func (b *legacyBackend) TxPool() TxPool         { return b.txPool }
func (b *legacyBackend) Miner() *miner.Miner    { return b.miner }
func (b *legacyBackend) APIs() []rpc.API        { return nil }
