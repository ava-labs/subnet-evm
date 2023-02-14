// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Defines the stateless interface for unmarshalling an arbitrary config of a precompile
package registry

// Force imports of each precompile to ensure each precompile's init function runs and registers itself
// with the registry.
//
import (
	_ "github.com/ava-labs/subnet-evm/precompile/contracts/deployerallowlist"
	_ "github.com/ava-labs/subnet-evm/precompile/contracts/feemanager"
	_ "github.com/ava-labs/subnet-evm/precompile/contracts/nativeminter"
	_ "github.com/ava-labs/subnet-evm/precompile/contracts/rewardmanager"
	_ "github.com/ava-labs/subnet-evm/precompile/contracts/txallowlist"
	// ADD YOUR PRECOMPILE HERE
	// _ "github.com/ava-labs/subnet-evm/precompile/contracts/yourprecompile"
)
