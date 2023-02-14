// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package allowlist

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
)

// Config specifies the initial set of allow list admins.
type Config struct {
	AdminAddresses   []common.Address `json:"adminAddresses,omitempty"`
	EnabledAddresses []common.Address `json:"enabledAddresses,omitempty"` // initial enabled addresses
}

// Configure initializes the address space of [precompileAddr] by initializing the role of each of
// the addresses in [AllowListAdmins].
func (c *Config) Configure(state contract.StateDB, precompileAddr common.Address) error {
	for _, enabledAddr := range c.EnabledAddresses {
		SetAllowListRole(state, precompileAddr, enabledAddr, EnabledRole)
	}
	for _, adminAddr := range c.AdminAddresses {
		SetAllowListRole(state, precompileAddr, adminAddr, AdminRole)
	}
	return nil
}

// Equal returns true iff [other] has the same admins in the same order in its allow list.
func (c *Config) Equal(other *Config) bool {
	if other == nil {
		return false
	}
	if !areEqualAddressLists(c.AdminAddresses, other.AdminAddresses) {
		return false
	}

	return areEqualAddressLists(c.EnabledAddresses, other.EnabledAddresses)
}

// areEqualAddressLists returns true iff [a] and [b] have the same addresses in the same order.
func areEqualAddressLists(current []common.Address, other []common.Address) bool {
	if len(current) != len(other) {
		return false
	}
	for i, address := range current {
		if address != other[i] {
			return false
		}
	}
	return true
}

// Verify returns an error if there is an overlapping address between admin and enabled roles
func (c *Config) Verify() error {
	// return early if either list is empty
	if len(c.EnabledAddresses) == 0 || len(c.AdminAddresses) == 0 {
		return nil
	}

	addressMap := make(map[common.Address]bool)
	for _, enabledAddr := range c.EnabledAddresses {
		// check for duplicates
		if _, ok := addressMap[enabledAddr]; ok {
			return fmt.Errorf("duplicate address %s in enabled list", enabledAddr)
		}
		addressMap[enabledAddr] = false
	}

	for _, adminAddr := range c.AdminAddresses {
		// check for overlap between enabled and admin lists
		if inAdmin, ok := addressMap[adminAddr]; ok {
			if inAdmin {
				return fmt.Errorf("duplicate address %s in admin list", adminAddr)
			} else {
				return fmt.Errorf("cannot set address %s as both admin and enabled", adminAddr)
			}
		}
		addressMap[adminAddr] = true
	}

	return nil
}
