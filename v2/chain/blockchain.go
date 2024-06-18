package chain

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/ava-labs/coreth/consensus"
	"github.com/ava-labs/coreth/core"
	"github.com/ava-labs/coreth/core/state"
	"github.com/ava-labs/coreth/core/types"
	"github.com/ava-labs/coreth/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/lru"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
)

var _ BlockChain = (*blockChain)(nil)

type blockChain struct {
	lock sync.RWMutex

	blocksDb ethdb.Database // Block chain to store block headers and bodies
	state    state.Database

	// TODO: should make a config struct?
	config      *params.ChainConfig
	cacheConfig *core.CacheConfig
	vmConfig    vm.Config
	engine      consensus.Engine

	hc           *core.HeaderChain
	blockCache   *lru.Cache[common.Hash, *types.Block] // Cache for the most recent entire blocks
	lastAccepted atomic.Pointer[types.Block]           // Prevents reorgs past this height

	senderCacher *core.TxSenderCacher
}

func NewBlockChain(
	blocksDb ethdb.Database,
	config *params.ChainConfig,
	cacheConfig *core.CacheConfig,
	vmConfig vm.Config,
	engine consensus.Engine,
) (*blockChain, error) {
	hc, err := core.NewHeaderChain(blocksDb, config, cacheConfig, engine)
	if err != nil {
		return nil, err
	}

	return &blockChain{
		blocksDb:     blocksDb,
		config:       config,
		cacheConfig:  cacheConfig,
		engine:       engine,
		vmConfig:     vmConfig,
		senderCacher: core.NewTxSenderCacher(runtime.NumCPU()),
		hc:           hc,
	}, nil
}

func (bc *blockChain) Stop() {
	bc.senderCacher.Shutdown()
}

func (bc *blockChain) Accept(*types.Block) error                  { panic("unimplemented") }
func (bc *blockChain) Reject(*types.Block) error                  { panic("unimplemented") }
func (bc *blockChain) InsertBlockManual(*types.Block, bool) error { panic("unimplemented") }
func (bc *blockChain) SetPreference(*types.Block) error           { panic("unimplemented") }

// Subscriptions
func (bc *blockChain) SubscribeAcceptedLogsEvent(ch chan<- []*types.Log) event.Subscription {
	panic("unimplemented")
}

func (bc *blockChain) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	panic("unimplemented")
}

// Getters
func (bc *blockChain) CurrentBlock() *types.Header              { return bc.hc.CurrentHeader() }
func (bc *blockChain) CurrentHeader() *types.Header             { return bc.hc.CurrentHeader() }
func (bc *blockChain) LastAcceptedBlock() *types.Block          { return bc.lastAccepted.Load() }
func (bc *blockChain) LastConsensusAcceptedBlock() *types.Block { return bc.lastAccepted.Load() }
func (bc *blockChain) SenderCacher() *core.TxSenderCacher       { return bc.senderCacher }
func (bc *blockChain) CacheConfig() *core.CacheConfig           { return bc.cacheConfig }
func (bc *blockChain) Config() *params.ChainConfig              { return bc.config }
func (bc *blockChain) GetVMConfig() *vm.Config                  { return &bc.vmConfig }
func (bc *blockChain) Engine() consensus.Engine                 { return bc.engine }

func (bc *blockChain) HasState(hash common.Hash) bool {
	_, err := bc.state.OpenTrie(hash)
	return err == nil
}
func (bc *blockChain) StateAt(root common.Hash) (*state.StateDB, error) {
	return state.New(root, bc.state, nil)
}
func (bc *blockChain) State() (*state.StateDB, error) {
	return bc.StateAt(bc.CurrentHeader().Root)
}

// No-ops
func (bc *blockChain) DrainAcceptorQueue()           {}
func (bc *blockChain) InitializeSnapshots()          {}
func (bc *blockChain) ValidateCanonicalChain() error { return nil }
