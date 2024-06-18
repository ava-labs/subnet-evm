// (c) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"math/big"

	"github.com/ava-labs/coreth/commontype"
	"github.com/ava-labs/coreth/core"
	"github.com/ava-labs/coreth/core/types"
	"github.com/ava-labs/coreth/eth"
	"github.com/ava-labs/coreth/miner"
	"github.com/ava-labs/coreth/params"
	"github.com/ava-labs/coreth/precompile/contract"
	"github.com/ava-labs/coreth/rpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
)

type BlockChain interface {
	Accept(*types.Block) error
	Reject(*types.Block) error
	CurrentBlock() *types.Header
	CurrentHeader() *types.Header
	LastAcceptedBlock() *types.Block
	LastConsensusAcceptedBlock() *types.Block
	GetBlockByNumber(uint64) *types.Block
	InitializeSnapshots()
	HasBlock(common.Hash, uint64) bool
	DrainAcceptorQueue()
	HasState(common.Hash) bool
	State() (StateDB, error)
	StateAt(common.Hash) (StateDB, error)
	ValidateCanonicalChain() error
	InsertBlockManual(*types.Block, bool) error
	GetBlockByHash(common.Hash) *types.Block
	SetPreference(*types.Block) error
	GetFeeConfigAt(parent *types.Header) (commontype.FeeConfig, *big.Int, error)
	SubscribeAcceptedLogsEvent(ch chan<- []*types.Log) event.Subscription
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

type Backend interface {
	BlockChain() BlockChain
	TxPool() TxPool
	Miner() *miner.Miner
	EstimateBaseFee(context.Context) (*big.Int, error)
	Start()
	Stop() error
	SetEtherbase(common.Address)
	ResetToStateSyncedBlock(*types.Block) error
	APIs() []rpc.API
}

type ethBackender struct {
	*eth.Ethereum
}

func (e *ethBackender) BlockChain() BlockChain {
	return &ethBlockChainer{e.Ethereum.BlockChain()}
}

func (e *ethBackender) TxPool() TxPool {
	return e.Ethereum.TxPool()
}

func (e *ethBackender) EstimateBaseFee(ctx context.Context) (*big.Int, error) {
	// Note: this is cheating a little, but it's only used to estimate
	// fees, and in principle we can fix the gpo to not depend on
	// the APIBackend (see OracleBackend).
	return e.Ethereum.APIBackend.EstimateBaseFee(ctx)
}

func (e *ethBackender) ResetToStateSyncedBlock(block *types.Block) error {
	// BloomIndexer needs to know that some parts of the chain are not available
	// and cannot be indexed. This is done by calling [AddCheckpoint] here.
	// Since the indexer uses sections of size [params.BloomBitsBlocks] (= 4096),
	// each block is indexed in section number [blockNumber/params.BloomBitsBlocks].
	// To allow the indexer to start with the block we just synced to,
	// we create a checkpoint for its parent.
	// Note: This requires assuming the synced block height is divisible
	// by [params.BloomBitsBlocks].
	parentHeight := block.NumberU64() - 1
	parentHash := block.ParentHash()
	e.Ethereum.BloomIndexer().AddCheckpoint(parentHeight/params.BloomBitsBlocks, parentHash)

	return e.Ethereum.BlockChain().ResetToStateSyncedBlock(block)
}

func (e *ethBackender) APIs() []rpc.API {
	return nil // deliberately turn off the APIs
}
