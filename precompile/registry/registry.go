// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Defines the stateless interface for unmarshalling an arbitrary config of a precompile
package registry

import (
	"github.com/ava-labs/avalanchego/utils/wrappers"
	"github.com/ava-labs/subnet-evm/precompile/contracts/deployerallowlist"
	"github.com/ava-labs/subnet-evm/precompile/contracts/feemanager"
	"github.com/ava-labs/subnet-evm/precompile/contracts/nativeminter"
	"github.com/ava-labs/subnet-evm/precompile/contracts/rewardmanager"
	"github.com/ava-labs/subnet-evm/precompile/contracts/txallowlist"
	"github.com/ava-labs/subnet-evm/precompile/registerer"
)

func RegisterPrecompileModules() error {
	errs := wrappers.Errs{}
	errs.Add(
		// Order is important here.
		// RegisterModule registers a precompile in the order it is registered.
		// The order of registration is important because it determines the configuration order
		// in the state.
		registerer.RegisterModule(deployerallowlist.Module{}),
		registerer.RegisterModule(nativeminter.Module{}),
		registerer.RegisterModule(txallowlist.Module{}),
		registerer.RegisterModule(feemanager.Module{}),
		registerer.RegisterModule(rewardmanager.Module{}),
	// ADD YOUR PRECOMPILE HERE
	// precompile.RegisterModule({yourPrecompilePackage}.{YourPrecompile}Config{}),
	)

	return errs.Err
}
