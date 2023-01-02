// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package deployerallowlist

import (
	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ethereum/go-ethereum/common"
)

var (
	// Singleton StatefulPrecompiledContract for W/R access to the contract deployer allow list.
	ContractDeployerAllowListPrecompile precompile.StatefulPrecompiledContract = precompile.CreateAllowListPrecompile(Address)
)

// GetContractDeployerAllowListStatus returns the role of [address] for the contract deployer
// allow list.
func GetContractDeployerAllowListStatus(stateDB precompile.StateDB, address common.Address) precompile.AllowListRole {
	return precompile.GetAllowListStatus(stateDB, Address, address)
}

// SetContractDeployerAllowListStatus sets the permissions of [address] to [role] for the
// contract deployer allow list.
// assumes [role] has already been verified as valid.
func SetContractDeployerAllowListStatus(stateDB precompile.StateDB, address common.Address, role precompile.AllowListRole) {
	precompile.SetAllowListRole(stateDB, Address, address, role)
}
