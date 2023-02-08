// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"fmt"
	"math/big"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/engine/snowman/block"
	"github.com/ava-labs/subnet-evm/plugin/evm"
	"github.com/ethereum/go-ethereum/common"
)

// PredicateContext provides context to stateful precompile predicates
type PredicateContext struct {
	SnowCtx            *snow.Context
	ProposerVMBlockCtx *block.Context
}

// PredicateFunc is the function type for validating that an access list tuple touching a precompile follows the predicate
// It is possible [blockContext] to be nil if the proposer VM is not fully activated, so predicates should explicitly check this.
type PredicateFunc func(predicateContext *PredicateContext, storageSlots []byte) error

// OnAcceptFunc is called on any log produced in a block where the address matches the precompile address
type OnAcceptFunc func(backend evm.Backend, txHash common.Hash, logIndex int, topics []common.Hash, logData []byte) error

// StatefulPrecompileConfig defines the interface for a stateful precompile to
type StatefulPrecompileConfig interface {
	// Address returns the address where the stateful precompile is accessible.
	Address() common.Address

	// Timestamp returns the timestamp at which this stateful precompile should be enabled.
	// 1) 0 indicates that the precompile should be enabled from genesis.
	// 2) n indicates that the precompile should be enabled in the first block with timestamp >= [n].
	// 3) nil indicates that the precompile is never enabled.
	Timestamp() *big.Int

	// IsDisabled returns true if this network upgrade should disable the precompile.
	IsDisabled() bool

	// Equal returns true if the provided argument configures the same precompile with the same parameters.
	Equal(StatefulPrecompileConfig) bool

	// Configure is called on the first block where the stateful precompile should be enabled.
	// This allows the stateful precompile to configure its own state via [StateDB] and [BlockContext] as necessary.
	// This function must be deterministic since it will impact the EVM state. If a change to the
	// config causes a change to the state modifications made in Configure, then it cannot be safely
	// made to the config after the network upgrade has gone into effect.
	//
	// Configure is called on the first block where the stateful precompile should be enabled. This
	// provides the config the ability to set its initial state and should only modify the state within
	// its own address space.
	Configure(ChainConfig, StateDB, BlockContext)

	// Contract returns a thread-safe singleton that can be used as the StatefulPrecompiledContract when
	// this config is enabled.
	Contract() StatefulPrecompiledContract

	// Predicate returns an optional function which is called a predicate on every transaction's access list
	// for any access list tuple whose address matches the address of the precompile.
	// This allows the precompile to enforce a predicate on the transaction itself for it to be considered valid
	// to be included in a block.
	Predicate() PredicateFunc

	// OnAccept returns an optional function which is called when a block gets accepted on each log whose
	// address matches the address of the precompile.
	// This can be used to perform precompile specific logic on acceptance, so that the precompile can emit
	// events and perform logic on those events only after the block has been accepted.
	OnAccept() OnAcceptFunc

	// Verify is called on startup and an error is treated as fatal. Configure can assume the Config has passed verification.
	Verify() error

	fmt.Stringer
}

// Configure sets the nonce and code to non-empty values then calls Configure on [precompileConfig] to make the necessary
// state update to enable the StatefulPrecompile.
// Assumes that [precompileConfig] is non-nil.
func Configure(chainConfig ChainConfig, blockContext BlockContext, precompileConfig StatefulPrecompileConfig, state StateDB) {
	// Set the nonce of the precompile's address (as is done when a contract is created) to ensure
	// that it is marked as non-empty and will not be cleaned up when the statedb is finalized.
	state.SetNonce(precompileConfig.Address(), 1)
	// Set the code of the precompile's address to a non-zero length byte slice to ensure that the precompile
	// can be called from within Solidity contracts. Solidity adds a check before invoking a contract to ensure
	// that it does not attempt to invoke a non-existent contract.
	state.SetCode(precompileConfig.Address(), []byte{0x1})
	precompileConfig.Configure(chainConfig, state, blockContext)
}
