package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/commontype"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

var _ consensus.ChainHeaderReader = (*WithFeeConfig)(nil)

type chainHeaderReader interface {
	Config() *params.ChainConfig
	CurrentHeader() *types.Header
	GetHeader(hash common.Hash, number uint64) *types.Header
	GetHeaderByNumber(number uint64) *types.Header
	GetHeaderByHash(hash common.Hash) *types.Header
	Engine() consensus.Engine
}

type WithFeeConfig struct {
	chainHeaderReader
	FeeConfig commontype.FeeConfig
	Coinbase  common.Address
	Modified  *big.Int
}

func (w *WithFeeConfig) GetFeeConfigAt(parent *types.Header) (commontype.FeeConfig, *big.Int, error) {
	return w.FeeConfig, w.Modified, nil
}

func (w *WithFeeConfig) GetCoinbaseAt(parent *types.Header) (common.Address, bool, error) {
	return w.Coinbase, false, nil
}
