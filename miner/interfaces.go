// (c) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package miner

import (
	"math/big"

	"github.com/ava-labs/coreth/consensus"
	"github.com/ava-labs/coreth/core"
	"github.com/ava-labs/coreth/core/state"
	"github.com/ava-labs/coreth/core/txpool"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

type TxPool interface {
	Locals() []common.Address
	PendingWithBaseFee(enforceTips bool, baseFee *big.Int) map[common.Address][]*txpool.LazyTransaction
}

type BlockChain interface {
	consensus.ChainHeaderReader
	Engine() consensus.Engine
	HasBlock(common.Hash, uint64) bool
	CacheConfig() *core.CacheConfig
	GetVMConfig() *vm.Config
	StateAt(common.Hash) (*state.StateDB, error)
}
