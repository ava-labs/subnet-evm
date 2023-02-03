// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/precompile/deployerallowlist"
	"github.com/ava-labs/subnet-evm/precompile/feemanager"
	"github.com/ava-labs/subnet-evm/precompile/nativeminter"
	"github.com/ava-labs/subnet-evm/precompile/rewardmanager"
	"github.com/ava-labs/subnet-evm/precompile/txallowlist"
	"github.com/ethereum/go-ethereum/common"
)

// This list is kept just for reference. The actual addresses defined in respective packages of precompiles.
// Note: it is important that none of these addresses conflict with each other or any other precompiles
// in core/vm/contracts.go.
// The first stateful precompiles were added in coreth to support nativeAssetCall and nativeAssetBalance. New stateful precompiles
// originating in coreth will continue at this prefix, so we reserve this range in subnet-evm so that they can be migrated into
// subnet-evm without issue.
// These start at the address: 0x0100000000000000000000000000000000000000 and will increment by 1.
// Optional precompiles implemented in subnet-evm start at 0x0200000000000000000000000000000000000000 and will increment by 1
// from here to reduce the risk of conflicts.
// For forks of subnet-evm, users should start at 0x0300000000000000000000000000000000000000 to ensure
// that their own modifications do not conflict with stateful precompiles that may be added to subnet-evm
// in the future.
// ContractDeployerAllowListAddress = common.HexToAddress("0x0200000000000000000000000000000000000000")
// ContractNativeMinterAddress      = common.HexToAddress("0x0200000000000000000000000000000000000001")
// TxAllowListAddress               = common.HexToAddress("0x0200000000000000000000000000000000000002")
// FeeManagerAddress                = common.HexToAddress("0x0200000000000000000000000000000000000003")
// RewardManagerAddress             = common.HexToAddress("0x0200000000000000000000000000000000000004")
// ADD YOUR PRECOMPILE HERE
// {YourPrecompile}Address          = common.HexToAddress("0x03000000000000000000000000000000000000??")
func init() {
	// Order is important here.
	// RegisterModule registers a precompile in the order it is registered.
	// The order of registration is important because it determines the configuration order
	// in the state.
	precompile.RegisterModule(deployerallowlist.NewModule())
	precompile.RegisterModule(nativeminter.NewModule())
	precompile.RegisterModule(txallowlist.InitModule(
		common.HexToAddress("0x0200000000000000000000000000000000000002"),
		"txAllowListConfig",
	))
	precompile.RegisterModule(feemanager.NewModule())
	precompile.RegisterModule(rewardmanager.NewModule())
	// ADD YOUR PRECOMPILE HERE
	// precompile.RegisterModule({yourPrecompilePackage}.{YourPrecompile}Config{})
}
