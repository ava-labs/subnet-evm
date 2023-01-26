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
	ContractAddress = common.HexToAddress("0x0200000000000000000000000000000000000002")

	// Singleton StatefulPrecompiledContract for W/R access to the tx allow list.
	TxAllowListPrecompile precompile.StatefulPrecompiledContract = allowlist.CreateAllowListPrecompile(ContractAddress)

	ErrSenderAddressNotAllowListed = errors.New("cannot issue transaction from non-allow listed address")
)

// GetTxAllowListStatus returns the role of [address] for the tx allow list.
func GetTxAllowListStatus(stateDB precompile.StateDB, address common.Address) allowlist.AllowListRole {
	return allowlist.GetAllowListStatus(stateDB, ContractAddress, address)
}

// SetTxAllowListStatus sets the permissions of [address] to [role] for the
// tx allow list.
// assumes [role] has already been verified as valid.
func SetTxAllowListStatus(stateDB precompile.StateDB, address common.Address, role allowlist.AllowListRole) {
	allowlist.SetAllowListRole(stateDB, ContractAddress, address, role)
}
