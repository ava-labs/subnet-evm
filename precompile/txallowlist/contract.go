// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txallowlist

import (
	"errors"

	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ precompile.StatefulPrecompileConfig = &TxAllowListConfig{}
	// Singleton StatefulPrecompiledContract for W/R access to the contract deployer allow list.
	TxAllowListPrecompile precompile.StatefulPrecompiledContract = allowlist.CreateAllowListPrecompile(precompile.TxAllowListAddress)

	ErrSenderAddressNotAllowListed = errors.New("cannot issue transaction from non-allow listed address")
)

// GetTxAllowListStatus returns the role of [address] for the contract deployer
// allow list.
func GetTxAllowListStatus(stateDB precompile.StateDB, address common.Address) allowlist.AllowListRole {
	return allowlist.GetAllowListStatus(stateDB, precompile.TxAllowListAddress, address)
}

// SetTxAllowListStatus sets the permissions of [address] to [role] for the
// tx allow list.
// assumes [role] has already been verified as valid.
func SetTxAllowListStatus(stateDB precompile.StateDB, address common.Address, role allowlist.AllowListRole) {
	allowlist.SetAllowListRole(stateDB, precompile.TxAllowListAddress, address, role)
}
