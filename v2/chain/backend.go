package chain

import (
	"context"
	"math/big"

	"github.com/ava-labs/avalanchego/utils/timer/mockable"
	"github.com/ava-labs/subnet-evm/consensus"
	"github.com/ava-labs/subnet-evm/core/txpool"
	"github.com/ava-labs/subnet-evm/core/txpool/legacypool"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/eth/gasprice"
	"github.com/ava-labs/subnet-evm/miner"
	"github.com/ava-labs/subnet-evm/rpc"
	"github.com/ethereum/go-ethereum/common"
)

// legacyBackend attaches legacy backend components (txPool, miner, and gas
// price oracle) to the chain so it can be used by the plugin interface.
type legacyBackend struct {
	chain  BlockChain
	txPool TxPool
	miner  *miner.Miner
	engine consensus.Engine
	gpo    *gasprice.Oracle
}

func NewLegacyBackend(
	chain BlockChain,
	poolConfig legacypool.Config,
	minerConfig *miner.Config,
	clock *mockable.Clock,
	gasPriceConfig gasprice.Config,
	allowUnfinalizedQueries bool,
) (*legacyBackend, error) {
	legacyPool := legacypool.New(poolConfig, chain)
	txPool, err := txpool.New(new(big.Int).SetUint64(poolConfig.PriceLimit), chain, []txpool.SubPool{legacyPool}) // Note: blobpool omitted
	if err != nil {
		return nil, err
	}

	engine := chain.Engine()
	miner := miner.New(chain, txPool, minerConfig, chain.Config(), engine, clock)
	gpoBackend := NewGPOBackend(chain, allowUnfinalizedQueries)
	gpo, err := gasprice.NewOracle(gpoBackend, gasPriceConfig)
	if err != nil {
		return nil, err
	}

	return &legacyBackend{
		chain:  chain,
		txPool: txPool,
		miner:  miner,
		engine: engine,
		gpo:    gpo,
	}, nil
}

func (b *legacyBackend) Start() {}
func (b *legacyBackend) Stop() error {
	b.txPool.Close()
	b.chain.Stop()
	b.engine.Close()
	return nil
}

func (b *legacyBackend) EstimateBaseFee(ctx context.Context) (*big.Int, error) {
	return b.gpo.EstimateBaseFee(ctx)
}

func (b *legacyBackend) SetEtherbase(etherbase common.Address) {
	b.miner.SetEtherbase(etherbase)
}

func (b *legacyBackend) ResetToStateSyncedBlock(block *types.Block) error {
	return b.chain.ResetToStateSyncedBlock(block)
}

func (b *legacyBackend) BlockChain() BlockChain { return b.chain }
func (b *legacyBackend) TxPool() TxPool         { return b.txPool }
func (b *legacyBackend) Miner() *miner.Miner    { return b.miner }
func (b *legacyBackend) APIs() []rpc.API        { return nil }
