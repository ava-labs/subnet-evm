// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Defines the interface for the configuration and execution of a precompile contract
package contract

import (
	"math/big"

	"github.com/ava-labs/avalanchego/chains/atomic"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/engine/snowman/block"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ethereum/go-ethereum/common"
)

// StatefulPrecompiledContract is the interface for executing a precompiled contract
type StatefulPrecompiledContract interface {
	// Run executes the precompiled contract.
	Run(accessibleState AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error)
}

// PredicateContext provides context to stateful precompile predicates
type PredicateContext struct {
	SnowCtx            *snow.Context
	ProposerVMBlockCtx *block.Context
}

// Predicater is an optional interface for StatefulPrecompiledContracts to implement.
// If implemented, the predicate will be enforced on every transaction in a block, prior to the block's execution.
// If VerifyPredicate returns an error, the block will fail verification with no further processing.
// WARNING: this is not intended to be used for custom precompiles. Backwards compatibility with custom precompiles that
// use the Predicater interface will not be supported.
type Predicater interface {
	VerifyPredicate(predicateContext *PredicateContext, storageSlots []byte) error
}

type SharedMemoryWriter interface {
	AddSharedMemoryRequests(chainID ids.ID, requests *atomic.Requests)
}

type AcceptContext struct {
	SnowCtx      *snow.Context
	SharedMemory SharedMemoryWriter
}

// Accepter is an optional interface for StatefulPrecompiledContracts to implement.
// If implemented, Accept will be called for every log with the address of the precompile when the block is accepted.
// WARNING: this is not intended to be used for custom precompiles. Backwards compatibility with custom precompiles that
// use the Accepter interface will not be supported.
type Accepter interface {
	Accept(acceptCtx *AcceptContext, txHash common.Hash, logIndex int, topics []common.Hash, logData []byte) error
}

// ChainContext defines an interface that provides information to a stateful precompile
// about the chain configuration. The precompile can access this information to initialize
// its state.
type ChainConfig interface {
	// GetFeeConfig returns the original FeeConfig that was set in the genesis.
	GetFeeConfig() commontype.FeeConfig
	// AllowedFeeRecipients returns true if fee recipients are allowed in the genesis.
	AllowedFeeRecipients() bool
}

// StateDB is the interface for accessing EVM state
type StateDB interface {
	GetState(common.Address, common.Hash) common.Hash
	SetState(common.Address, common.Hash, common.Hash)

	SetCode(common.Address, []byte)

	SetNonce(common.Address, uint64)
	GetNonce(common.Address) uint64

	GetBalance(common.Address) *big.Int
	AddBalance(common.Address, *big.Int)
	SubBalance(common.Address, *big.Int)

	CreateAccount(common.Address)
	Exist(common.Address) bool

	AddLog(addr common.Address, topics []common.Hash, data []byte, blockNumber uint64)
	GetPredicateStorageSlots(address common.Address) ([]byte, bool)

	Suicide(common.Address) bool
	Finalise(deleteEmptyObjects bool)
}

// AccessibleState defines the interface exposed to stateful precompile contracts
type AccessibleState interface {
	GetStateDB() StateDB
	GetBlockContext() BlockContext
	GetSnowContext() *snow.Context
	CallFromPrecompile(caller common.Address, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error)
}

// BlockContext defines an interface that provides information to a stateful precompile
// about the block that activates the upgrade. The precompile can access this information
// to initialize its state.
type BlockContext interface {
	Number() *big.Int
	Timestamp() *big.Int
}

type Configurator interface {
	MakeConfig() precompileconfig.Config
	Configure(
		chainConfig ChainConfig,
		precompileconfig precompileconfig.Config,
		state StateDB,
		blockContext BlockContext,
	) error
}
