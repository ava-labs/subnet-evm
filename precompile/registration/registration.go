// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package registration

import (
	"github.com/ava-labs/subnet-evm/precompile/deployerallowlist"
	"github.com/ava-labs/subnet-evm/precompile/modules"
)

// This init function is defined as a convenience location to register precompile modules
func init() {
	modules.RegisterModule(deployerallowlist.NewModule())
	// modules.RegisterModule(nativeminter.NewModule())
	// modules.RegisterModule(txallowlist.NewModule())
	// modules.RegisterModule(feemanager.NewModule())
	// modules.RegisterModule(rewardmanager.NewModule())
	// ADD YOUR PRECOMPILE HERE
	// mdoules.RegisterModule({yourPrecompilePackage}.{YourPrecompile}Config{})
}
