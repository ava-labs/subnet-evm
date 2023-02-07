// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package registration

// Force imports of each precompile to ensure each precompile's init function runs and registers itself
import (
	_ "github.com/ava-labs/subnet-evm/precompile/deployerallowlist"
	_ "github.com/ava-labs/subnet-evm/precompile/feemanager"
	_ "github.com/ava-labs/subnet-evm/precompile/nativeminter"
	_ "github.com/ava-labs/subnet-evm/precompile/rewardmanager"
	_ "github.com/ava-labs/subnet-evm/precompile/txallowlist"
)
