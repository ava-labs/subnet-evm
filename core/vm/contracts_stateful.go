// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package vm

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
// FeeManagerAddress         				= common.HexToAddress("0x0200000000000000000000000000000000000003")
// RewardManagerAddress             = common.HexToAddress("0x0200000000000000000000000000000000000004")
// ADD YOUR PRECOMPILE HERE
// {YourPrecompile}Address       = common.HexToAddress("0x03000000000000000000000000000000000000??")
func registerStatefulPrecompiles() {
	precompile.RegisterModule(deployerallowlist.ContractDeployerAllowListConfig{})
	precompile.RegisterModule(nativeminter.ContractNativeMinterConfig{})
	precompile.RegisterModule(txallowlist.TxAllowListConfig{})
	precompile.RegisterModule(feemanager.FeeConfigManagerConfig{})
	precompile.RegisterModule(rewardmanager.RewardManagerConfig{})
	// ADD YOUR PRECOMPILE HERE
	// RegisterModule("{yourPrecompile}Config", {YourPrecompile}Address, {YourPrecompile}.Create{YourPrecompile}Precompile)
}

// wrappedPrecompiledContract implements StatefulPrecompiledContract by wrapping stateless native precompiled contracts
// in Ethereum.
type wrappedPrecompiledContract struct {
	p PrecompiledContract
}

// newWrappedPrecompiledContract returns a wrapped version of [PrecompiledContract] to be executed according to the StatefulPrecompiledContract
// interface.
func newWrappedPrecompiledContract(p PrecompiledContract) precompile.StatefulPrecompiledContract {
	return &wrappedPrecompiledContract{p: p}
}

// Run implements the StatefulPrecompiledContract interface
func (w *wrappedPrecompiledContract) Run(accessibleState precompile.PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	return RunPrecompiledContract(w.p, input, suppliedGas)
}

// RunStatefulPrecompiledContract confirms runs [precompile] with the specified parameters.
func RunStatefulPrecompiledContract(precompile precompile.StatefulPrecompiledContract, accessibleState precompile.PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	return precompile.Run(accessibleState, caller, addr, input, suppliedGas, readOnly)
}
