// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Defines the interface for the configuration and execution of a precompile contract
package execution

import (
	"math/big"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/core/state"
	_ "github.com/ava-labs/subnet-evm/params"
	precompileConfig "github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ethereum/go-ethereum/common"
)

type Contract interface {
	// Run executes the precompiled contract.
	Run(accessibleState AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error)
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

// PrecompileAccessibleState defines the interface exposed to stateful precompile contracts
type AccessibleState interface {
	GetStateDB() *state.StateDB
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

type Execution interface {
	Configure(chainConfig ChainConfig, precompileConfig precompileConfig.Config, state *state.StateDB, blockContext BlockContext) error
	Contract() Contract
}
