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
	ErrSenderAddressNotAllowListed = errors.New("cannot issue transaction from non-allow listed address")
)

// GetTxAllowListStatus returns the role of [address] for the tx allow list.
func GetTxAllowListStatus(stateDB precompile.StateDB, address common.Address) allowlist.Role {
	return allowlist.GetAllowListStatus(stateDB, Module.Address, address)
}

// SetTxAllowListStatus sets the permissions of [address] to [role] for the
// tx allow list.
// assumes [role] has already been verified as valid.
func SetTxAllowListStatus(stateDB precompile.StateDB, address common.Address, role allowlist.Role) {
	allowlist.SetAllowListRole(stateDB, Module.Address, address, role)
}
