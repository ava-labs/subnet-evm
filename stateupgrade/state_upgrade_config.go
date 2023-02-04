// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package stateupgrade

import (
	"fmt"
	"math/big"
)

// StateUpgradeConfig defines the interface for a state upgrade configuration.
type StateUpgradeConfig interface {
	// Timestamp returns the timestamp at which this stateful precompile should be enabled.
	// 1) 0 indicates that the state upgrade should be enabled from genesis.
	// 2) n indicates that the state upgrade should be enabled in the first block with timestamp >= [n].
	// 3) nil indicates that the state upgrade is never enabled.
	Timestamp() *big.Int
	// IsDisabled returns true if this network upgrade should disable the precompile.
	IsDisabled() bool
	// Equal returns true if the provided argument configures the same precompile with the same parameters.
	Equal(config StateUpgradeConfig) bool
	// RunUpgrade is called on the first block where the stateful precompile should be enabled.
	// This allows the stateful precompile to configure its own state via [StateDB] and [BlockContext] as necessary.
	// This function must be deterministic since it will impact the EVM state. If a change to the
	// config causes a change to the state modifications made in Configure, then it cannot be safely
	// made to the config after the network upgrade has gone into effect.
	//
	// RunUpgrade is called on the first block where the stateful precompile should be enabled. This
	// provides the config the ability to set its initial state and should only modify the state within
	// its own address space.
	RunUpgrade(ChainConfig, StateDB, BlockContext)
	// Verify is called on startup and an error is treated as fatal. Configure can assume the Config has passed verification.
	Verify() error

	fmt.Stringer
}

// RunUpgrade calls the RunUpgrade method on [config] if it is non-nil.
// Assumes that [stateUpgradeConfig] is non-nil.
func RunUpgrade(chainConfig ChainConfig, blockContext BlockContext, stateUpgradeConfig StateUpgradeConfig, state StateDB) {
	stateUpgradeConfig.RunUpgrade(chainConfig, state, blockContext)
}
