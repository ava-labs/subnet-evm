package chain

import (
	"context"
	"errors"
	"math/big"

	"github.com/ava-labs/coreth/commontype"
	"github.com/ava-labs/coreth/consensus/dummy"
	"github.com/ava-labs/coreth/core"
	"github.com/ava-labs/coreth/core/types"
	"github.com/ava-labs/coreth/params"
	"github.com/ava-labs/coreth/rpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
)

var ErrUnfinalizedData = errors.New("cannot query unfinalized data")

type gpoBackend struct {
	chain                   BlockChain
	allowUnfinalizedQueries bool
}

func NewGPOBackend(chain BlockChain, allowUnfinalizedQueries bool) *gpoBackend {
	return &gpoBackend{
		chain:                   chain,
		allowUnfinalizedQueries: allowUnfinalizedQueries,
	}
}

func (b *gpoBackend) IsAllowUnfinalizedQueries() bool {
	return b.allowUnfinalizedQueries
}

func (b *gpoBackend) isLatestAndAllowed(number rpc.BlockNumber) bool {
	return number.IsLatest() && b.IsAllowUnfinalizedQueries()
}

func (b *gpoBackend) HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	// Treat requests for the pending, latest, or accepted block
	// identically.
	acceptedBlock := b.LastAcceptedBlock()
	if number.IsAccepted() {
		if b.isLatestAndAllowed(number) {
			return b.chain.CurrentHeader(), nil
		}
		return acceptedBlock.Header(), nil
	}

	if !b.IsAllowUnfinalizedQueries() && acceptedBlock != nil {
		if number.Int64() > acceptedBlock.Number().Int64() {
			return nil, ErrUnfinalizedData
		}
	}

	return b.chain.GetHeaderByNumber(uint64(number)), nil
}

func (b *gpoBackend) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	// Treat requests for the pending, latest, or accepted block
	// identically.
	acceptedBlock := b.LastAcceptedBlock()
	if number.IsAccepted() {
		if b.isLatestAndAllowed(number) {
			header := b.chain.CurrentBlock()
			return b.chain.GetBlock(header.Hash(), header.Number.Uint64()), nil
		}
		return acceptedBlock, nil
	}

	if !b.IsAllowUnfinalizedQueries() && acceptedBlock != nil {
		if number.Int64() > acceptedBlock.Number().Int64() {
			return nil, ErrUnfinalizedData
		}
	}

	return b.chain.GetBlockByNumber(uint64(number)), nil
}

func (b *gpoBackend) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return b.chain.GetReceiptsByHash(hash), nil
}

func (b *gpoBackend) MinRequiredTip(ctx context.Context, header *types.Header) (*big.Int, error) {
	return dummy.MinRequiredTip(b.ChainConfig(), header)
}

func (g *gpoBackend) ChainConfig() *params.ChainConfig {
	return g.chain.Config()
}

func (g *gpoBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return g.chain.SubscribeChainHeadEvent(ch)
}

func (g *gpoBackend) SubscribeChainAcceptedEvent(ch chan<- core.ChainEvent) event.Subscription {
	return g.chain.SubscribeChainAcceptedEvent(ch)
}

func (g *gpoBackend) LastAcceptedBlock() *types.Block {
	return g.chain.LastAcceptedBlock()
}

func (g *gpoBackend) GetFeeConfigAt(parent *types.Header) (commontype.FeeConfig, *big.Int, error) {
	return g.chain.GetFeeConfigAt(parent)
}
