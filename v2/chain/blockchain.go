package chain

import (
	"sync"

	"github.com/ava-labs/coreth/consensus"
	"github.com/ava-labs/coreth/core"
	"github.com/ava-labs/coreth/core/state"
	"github.com/ava-labs/coreth/core/types"
	"github.com/ava-labs/coreth/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/event"
)

var _ BlockChain = (*blockChain)(nil)

type blockChain struct {
	lock sync.RWMutex
}

func (b *blockChain) Accept(*types.Block) error                               { panic("unimplemented") }
func (b *blockChain) Reject(*types.Block) error                               { panic("unimplemented") }
func (b *blockChain) CurrentBlock() *types.Header                             { panic("unimplemented") }
func (b *blockChain) CurrentHeader() *types.Header                            { panic("unimplemented") }
func (b *blockChain) LastAcceptedBlock() *types.Block                         { panic("unimplemented") }
func (b *blockChain) LastConsensusAcceptedBlock() *types.Block                { panic("unimplemented") }
func (b *blockChain) GetBlockByNumber(uint64) *types.Block                    { panic("unimplemented") }
func (b *blockChain) HasBlock(common.Hash, uint64) bool                       { panic("unimplemented") }
func (b *blockChain) GetBlock(common.Hash, uint64) *types.Block               { panic("unimplemented") }
func (b *blockChain) HasState(common.Hash) bool                               { panic("unimplemented") }
func (b *blockChain) State() (*state.StateDB, error)                          { panic("unimplemented") }
func (b *blockChain) StateAt(common.Hash) (*state.StateDB, error)             { panic("unimplemented") }
func (b *blockChain) InsertBlockManual(*types.Block, bool) error              { panic("unimplemented") }
func (b *blockChain) GetBlockByHash(common.Hash) *types.Block                 { panic("unimplemented") }
func (b *blockChain) SetPreference(*types.Block) error                        { panic("unimplemented") }
func (b *blockChain) GetHeader(hash common.Hash, number uint64) *types.Header { panic("unimplemented") }
func (b *blockChain) GetHeaderByHash(hash common.Hash) *types.Header          { panic("unimplemented") }
func (b *blockChain) GetHeaderByNumber(number uint64) *types.Header           { panic("unimplemented") }
func (b *blockChain) SenderCacher() *core.TxSenderCacher                      { panic("unimplemented") }

func (b *blockChain) CacheConfig() *core.CacheConfig { panic("unimplemented") }
func (b *blockChain) Config() *params.ChainConfig    { panic("unimplemented") }
func (b *blockChain) GetVMConfig() *vm.Config        { panic("unimplemented") }
func (b *blockChain) Engine() consensus.Engine       { panic("unimplemented") }
func (b *blockChain) Stop()                          { panic("unimplemented") }

// Subscriptions
func (b *blockChain) SubscribeAcceptedLogsEvent(ch chan<- []*types.Log) event.Subscription {
	panic("unimplemented")
}

func (b *blockChain) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	panic("unimplemented")
}

// No-ops
func (b *blockChain) DrainAcceptorQueue()           {}
func (b *blockChain) InitializeSnapshots()          {}
func (b *blockChain) ValidateCanonicalChain() error { return nil }
