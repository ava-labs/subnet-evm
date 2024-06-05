// (c) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"math/big"

	"github.com/ava-labs/coreth/commontype"
	"github.com/ava-labs/coreth/core/state"
	"github.com/ava-labs/coreth/core/types"
	"github.com/ethereum/go-ethereum/common"
)

type BlockChain interface {
	Accept(*types.Block) error
	Reject(*types.Block) error
	CurrentBlock() *types.Header
	LastAcceptedBlock() *types.Block
	LastConsensusAcceptedBlock() *types.Block
	GetBlockByNumber(uint64) *types.Block
	InitializeSnapshots()
	HasBlock(common.Hash, uint64) bool
	GetBlock(common.Hash, uint64) *types.Block
	DrainAcceptorQueue()
	HasState(common.Hash) bool
	State() (*state.StateDB, error)
	StateAt(common.Hash) (*state.StateDB, error)
	ValidateCanonicalChain() error
	InsertBlockManual(*types.Block, bool) error
	GetBlockByHash(common.Hash) *types.Block
	SetPreference(*types.Block) error
	GetFeeConfigAt(parent *types.Header) (commontype.FeeConfig, *big.Int, error)
}
