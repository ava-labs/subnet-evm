// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package miner implements Ethereum block creation and mining.
package miner

import (
	"sync"

	"github.com/ava-labs/avalanchego/utils/timer/mockable"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/precompile/precompileconfig"
)

// Backend wraps all methods required for mining.
type Backend interface {
	BlockChain() *core.BlockChain
	TxPool() *txpool.TxPool
}

// Config is the configuration parameters of mining.
type Config struct {
	Etherbase            common.Address `toml:",omitempty"` // Public address for block mining rewards
	AllowDuplicateBlocks bool           // Allow mining of duplicate blocks (used in tests only)
}

// Miner is the main object which takes care of submitting new work to consensus
// engine and gathering the sealing result.
type Miner struct {
	confMu      sync.RWMutex // The lock used to protect the config fields: GasCeil, GasTip and Extradata
	config      *Config
	chainConfig *params.ChainConfig
	engine      consensus.Engine
	txpool      *txpool.TxPool
	chain       *core.BlockChain

	mu         sync.RWMutex
	coinbase   common.Address
	clock      *mockable.Clock // Allows us mock the clock for testing
	beaconRoot *common.Hash    // TODO: set to empty hash, retained for upstream compatibility and future use
}

func (miner *Miner) SetEtherbase(addr common.Address) {
	miner.mu.Lock()
	defer miner.mu.Unlock()
	miner.coinbase = addr
}

func (miner *Miner) GenerateBlock(predicateContext *precompileconfig.PredicateContext) (*types.Block, error) {
	return miner.commitNewWork(predicateContext)
}

// New creates a new miner with provided config.
func New(eth Backend, config Config, engine consensus.Engine, clock *mockable.Clock) *Miner {
	return &Miner{
		config:      &config,
		chainConfig: eth.BlockChain().Config(),
		engine:      engine,
		txpool:      eth.TxPool(),
		chain:       eth.BlockChain(),
		coinbase:    config.Etherbase,
		clock:       clock,
		beaconRoot:  &common.Hash{},
	}
}

// XXX
// // SetExtra sets the content used to initialize the block extra field.
// func (miner *Miner) SetExtra(extra []byte) error {
// if uint64(len(extra)) > params.MaximumExtraDataSize {
// return fmt.Errorf("extra exceeds max length. %d > %v", len(extra), params.MaximumExtraDataSize)
// }
// miner.confMu.Lock()
// miner.config.ExtraData = extra
// miner.confMu.Unlock()
// return nil
// }

// // SetGasCeil sets the gaslimit to strive for when mining blocks post 1559.
// // For pre-1559 blocks, it sets the ceiling.
// func (miner *Miner) SetGasCeil(ceil uint64) {
// miner.confMu.Lock()
// miner.config.GasCeil = ceil
// miner.confMu.Unlock()
// }

// // SetGasTip sets the minimum gas tip for inclusion.
// func (miner *Miner) SetGasTip(tip *big.Int) error {
// miner.confMu.Lock()
// miner.config.GasPrice = tip
// miner.confMu.Unlock()
// return nil
// }
