// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Defines the stateless interface for unmarshalling an arbitrary config of a precompile
package precompileconfig

import (
	"math/big"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/engine/snowman/block"
	"github.com/ethereum/go-ethereum/common"
)

// StatefulPrecompileConfig defines the interface for a stateful precompile to
// be enabled via a network upgrade.
type Config interface {
	// Key returns the unique key for the stateful precompile.
	Key() string
	// Timestamp returns the timestamp at which this stateful precompile should be enabled.
	// 1) 0 indicates that the precompile should be enabled from genesis.
	// 2) n indicates that the precompile should be enabled in the first block with timestamp >= [n].
	// 3) nil indicates that the precompile is never enabled.
	Timestamp() *big.Int
	// IsDisabled returns true if this network upgrade should disable the precompile.
	IsDisabled() bool
	// Equal returns true if the provided argument configures the same precompile with the same parameters.
	Equal(Config) bool
	// Verify is called on startup and an error is treated as fatal. Configure can assume the Config has passed verification.
	Verify() error
}

// PredicateContext is the context passed in to the Predicater interface.
type PrecompilePredicateContext struct {
	SnowCtx *snow.Context
}

// Predicater is an optional interface for StatefulPrecompileContracts to implement.
// If implemented, the predicate will be enforced on every transaction in a block, prior to
// the blcok's execution.
// If VerifyPredicate returns an error, the block will fail verification with no further processing.
// WARNING: If you are implementing a custom precompile, beware that subnet-evm
// will not maintain backwards compatibility of this interface and your code should not
// rely on this. Designed for use only by precompiles that ship with subnet-evm.
type PrecompilePredicater interface {
	VerifyPredicate(predicateContext *PrecompilePredicateContext, storageSlots []byte) error
}


// ProposerPredicateContext is the context passed in to the ProposerPredicater interface to verify
// a precompile predicate within a specific ProposerVM wrapper.
type ProposerPredicateContext struct {
	// Embed the general precompile predicate context in the proposer predicate context for the
	// context provided to precompiles that DO NOT require the ProposerVM Block Context.
	PrecompilePredicateContext
	// Note: ProposerVMBlockCtx may be nil if the Snowman Consensus Engine calls BuildBlock or Verify
	// instead of BuildBlockWithContext or VerifyWithContext.
	// In this case, it is up to the precompile to determine if a nil ProposerVMBlockCtx is valid.
	ProposerVMBlockCtx *block.Context
}

// EnginePredicater is an optional interface for StatefulPrecompiledContracts to implement.
// If implemented, the predicate will be enforced on every transaction in a block, prior to
// the block's execution.
// If VerifyPredicate returns an error, the block will fail verification with no further processing.
// Note: ProposerVMBlockCtx may be nil if the engine does not specify it. In this case,
// it's up to the precompile to determine if a nil ProposerVMBlockCtx is valid.
// WARNING: If you are implementing a custom precompile, beware that subnet-evm
// will not maintain backwards compatibility of this interface and your code should not
// rely on this. Designed for use only by precompiles that ship with subnet-evm.
type ProposerPredicater interface {
	VerifyPredicate(proposerPredicateContext *ProposerPredicateContext, storageSlots []byte) error
}

// Accepter is an optional interface for StatefulPrecompiledContracts to implement.
// If implemented, Accept will be called for every log with the address of the precompile when the block is accepted.
// WARNING: If you are implementing a custom precompile, beware that subnet-evm
// will not maintain backwards compatibility of this interface and your code should not
// rely on this. Designed for use only by precompiles that ship with subnet-evm.
type Accepter interface {
	Accept(txHash common.Hash, logIndex int, topics []common.Hash, logData []byte) error
}
