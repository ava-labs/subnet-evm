package chain

import (
	"math/big"

	"github.com/ava-labs/coreth/consensus"
	"github.com/ava-labs/coreth/core"
	"github.com/ava-labs/coreth/core/state"
	"github.com/ava-labs/coreth/core/txpool"
	"github.com/ava-labs/coreth/core/txpool/legacypool"
	"github.com/ava-labs/coreth/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
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
	State() (*state.StateDB, error)
	StateAt(common.Hash) (*state.StateDB, error)
	ValidateCanonicalChain() error
	InsertBlockManual(*types.Block, bool) error
	GetBlockByHash(common.Hash) *types.Block
	SetPreference(*types.Block) error
	SubscribeAcceptedLogsEvent(ch chan<- []*types.Log) event.Subscription
	SubscribeChainAcceptedEvent(ch chan<- core.ChainEvent) event.Subscription

	GetReceiptsByHash(common.Hash) types.Receipts
	ResetToStateSyncedBlock(block *types.Block) error
	Stop()

	// used by miner
	consensus.ChainHeaderReader
	Engine() consensus.Engine
	CacheConfig() *core.CacheConfig
	GetVMConfig() *vm.Config

	// used by txpool
	legacypool.BlockChain
	txpool.BlockChain
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

	Locals() []common.Address
	PendingWithBaseFee(enforceTips bool, baseFee *big.Int) map[common.Address][]*txpool.LazyTransaction

	Close() error
}

type committableStateDB interface {
	state.Database
	Commit(root common.Hash, report bool) error
	Initialized(root common.Hash) bool
	Close(lastBlockRoot common.Hash) error
}
