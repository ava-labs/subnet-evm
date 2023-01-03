// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package deployerallowlist

import (
	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ethereum/go-ethereum/common"
)

// Singleton StatefulPrecompiledContract for W/R access to the contract deployer allow list.
var ContractDeployerAllowListPrecompile precompile.StatefulPrecompiledContract = allowlist.CreateAllowListPrecompile(precompile.ContractDeployerAllowListAddress)

// GetContractDeployerAllowListStatus returns the role of [address] for the contract deployer
// allow list.
func GetContractDeployerAllowListStatus(stateDB precompile.StateDB, address common.Address) allowlist.AllowListRole {
	return allowlist.GetAllowListStatus(stateDB, precompile.ContractDeployerAllowListAddress, address)
}

// SetContractDeployerAllowListStatus sets the permissions of [address] to [role] for the
// contract deployer allow list.
// assumes [role] has already been verified as valid.
func SetContractDeployerAllowListStatus(stateDB precompile.StateDB, address common.Address, role allowlist.AllowListRole) {
	allowlist.SetAllowListRole(stateDB, precompile.ContractDeployerAllowListAddress, address, role)
}
