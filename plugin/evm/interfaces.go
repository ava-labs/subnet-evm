// (c) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"math/big"

	"github.com/ava-labs/coreth/commontype"
	"github.com/ava-labs/coreth/core"
	"github.com/ava-labs/coreth/core/types"
	"github.com/ava-labs/coreth/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
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
	State() (StateDB, error)
	StateAt(common.Hash) (StateDB, error)
	ValidateCanonicalChain() error
	InsertBlockManual(*types.Block, bool) error
	GetBlockByHash(common.Hash) *types.Block
	SetPreference(*types.Block) error
	GetFeeConfigAt(parent *types.Header) (commontype.FeeConfig, *big.Int, error)
}

type ethBlockChainer struct {
	*core.BlockChain
}

func (e *ethBlockChainer) State() (StateDB, error) {
	return e.BlockChain.State()
}

func (e *ethBlockChainer) StateAt(root common.Hash) (StateDB, error) {
	return e.BlockChain.StateAt(root)
}

type StateDB interface {
	SubBalance(common.Address, *big.Int)
	AddBalance(common.Address, *big.Int)
	GetBalance(common.Address) *big.Int

	GetNonce(common.Address) uint64
	SetNonce(common.Address, uint64)

	GetCode(common.Address) []byte
	GetCodeHash(common.Address) common.Hash

	contract.StateDB
}

type TxPool interface {
	PendingSize(enforceTips bool) int
	IteratePending(f func(tx *types.Transaction) bool)

	SubscribeTransactions(ch chan<- core.NewTxsEvent, reorgs bool) event.Subscription
	SubscribeNewReorgEvent(ch chan<- core.NewTxPoolReorgEvent) event.Subscription

	Add(txs []*types.Transaction, local bool, sync bool) []error
	AddRemotesSync(txs []*types.Transaction) []error
	Has(hash common.Hash) bool

	SetMinFee(fee *big.Int)
	SetGasTip(tip *big.Int)
	GasTip() *big.Int
}
