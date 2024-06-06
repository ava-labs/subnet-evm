package eth

import (
	"github.com/ava-labs/coreth/consensus"
	"github.com/ava-labs/coreth/core"
	"github.com/ava-labs/coreth/core/state"
	"github.com/ava-labs/coreth/core/txpool"
	"github.com/ava-labs/coreth/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

type BlockChain interface {
	HasBlock(common.Hash, uint64) bool
	GetBlock(common.Hash, uint64) *types.Block
	LastAcceptedBlock() *types.Block

	consensus.ChainHeaderReader
	Engine() consensus.Engine
	CacheConfig() *core.CacheConfig
	GetVMConfig() *vm.Config
	StateAt(common.Hash) (*state.StateDB, error)

	txpool.BlockChain
	SenderCacher() *core.TxSenderCacher
}
