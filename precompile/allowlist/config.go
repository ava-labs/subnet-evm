// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package allowlist

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/ava-labs/avalanchego/utils/wrappers"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ethereum/go-ethereum/common"
)

var ErrCannotAddManagersBeforeDUpgrade = fmt.Errorf("cannot add managers before DUpgrade")

// AllowListConfig specifies the initial set of addresses with Admin or Enabled roles.
type AllowListConfig struct {
	AdminAddresses   []common.Address `json:"adminAddresses,omitempty" serialize:"true"`   // initial admin addresses
	ManagerAddresses []common.Address `json:"managerAddresses,omitempty" serialize:"true"` // initial manager addresses
	EnabledAddresses []common.Address `json:"enabledAddresses,omitempty" serialize:"true"` // initial enabled addresses
}

// Configure initializes the address space of [precompileAddr] by initializing the role of each of
// the addresses in [AllowListAdmins].
func (c *AllowListConfig) Configure(chainConfig precompileconfig.ChainConfig, precompileAddr common.Address, state contract.StateDB, blockContext contract.ConfigurationBlockContext) error {
	for _, enabledAddr := range c.EnabledAddresses {
		SetAllowListRole(state, precompileAddr, enabledAddr, EnabledRole)
	}
	for _, adminAddr := range c.AdminAddresses {
		SetAllowListRole(state, precompileAddr, adminAddr, AdminRole)
	}
	// Verify() should have been called before Configure()
	// so we know manager role is activated
	for _, managerAddr := range c.ManagerAddresses {
		SetAllowListRole(state, precompileAddr, managerAddr, ManagerRole)
	}
	return nil
}

func (c *AllowListConfig) packAddresses(addresses []common.Address, p *wrappers.Packer) error {
	sort.Slice(addresses, func(i, j int) bool {
		return bytes.Compare(addresses[i][:], addresses[j][:]) < 0
	})

	p.PackInt(uint32(len(addresses)))
	if p.Err != nil {
		return p.Err
	}

	for _, address := range addresses {
		p.PackBytes(address[:])
		if p.Err != nil {
			return p.Err
		}
	}
	return nil
}

func (c *AllowListConfig) unpackAddresses(p *wrappers.Packer) ([]common.Address, error) {
	length := p.UnpackInt()
	if p.Err != nil {
		return nil, p.Err
	}

	addresses := make([]common.Address, length)
	for i := uint32(0); i < length; i++ {
		bytes := p.UnpackBytes()
		addresses = append(addresses[:i], common.BytesToAddress(bytes))
		if p.Err != nil {
			return nil, p.Err
		}
	}

	return addresses, nil
}

// Equal returns true iff [other] has the same admins in the same order in its allow list.
func (c *AllowListConfig) Equal(other *AllowListConfig) bool {
	if other == nil {
		return false
	}

	return areEqualAddressLists(c.AdminAddresses, other.AdminAddresses) &&
		areEqualAddressLists(c.ManagerAddresses, other.ManagerAddresses) &&
		areEqualAddressLists(c.EnabledAddresses, other.EnabledAddresses)
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
func (c *AllowListConfig) Verify(chainConfig precompileconfig.ChainConfig, upgrade precompileconfig.Upgrade) error {
	addressMap := make(map[common.Address]Role) // tracks which addresses we have seen and their role

	// check for duplicates in enabled list
	for _, enabledAddr := range c.EnabledAddresses {
		if _, ok := addressMap[enabledAddr]; ok {
			return fmt.Errorf("duplicate address in enabled list: %s", enabledAddr)
		}
		addressMap[enabledAddr] = EnabledRole
	}

	// check for overlap between enabled and admin lists or duplicates in admin list
	for _, adminAddr := range c.AdminAddresses {
		if role, ok := addressMap[adminAddr]; ok {
			if role == AdminRole {
				return fmt.Errorf("duplicate address in admin list: %s", adminAddr)
			} else {
				return fmt.Errorf("cannot set address as both admin and enabled: %s", adminAddr)
			}
		}
		addressMap[adminAddr] = AdminRole
	}

	if len(c.ManagerAddresses) != 0 && upgrade.Timestamp() != nil {
		// If the config attempts to activate a manager before the DUpgrade, fail verification
		timestamp := *upgrade.Timestamp()
		if !chainConfig.IsDUpgrade(timestamp) {
			return ErrCannotAddManagersBeforeDUpgrade
		}
	}

	// check for overlap between admin and manager lists or duplicates in manager list
	for _, managerAddr := range c.ManagerAddresses {
		if role, ok := addressMap[managerAddr]; ok {
			switch role {
			case ManagerRole:
				return fmt.Errorf("duplicate address in manager list: %s", managerAddr)
			case AdminRole:
				return fmt.Errorf("cannot set address as both admin and manager: %s", managerAddr)
			case EnabledRole:
				return fmt.Errorf("cannot set address as both enabled and manager: %s", managerAddr)
			}
		}
		addressMap[managerAddr] = ManagerRole
	}

	return nil
}

func (c *AllowListConfig) ToBytesWithPacker(p *wrappers.Packer) error {
	if err := c.packAddresses(c.AdminAddresses, p); err != nil {
		return err
	}
	if err := c.packAddresses(c.ManagerAddresses, p); err != nil {
		return err
	}
	if err := c.packAddresses(c.EnabledAddresses, p); err != nil {
		return err
	}
	return nil
}

func (c *AllowListConfig) FromBytesWithPacker(p *wrappers.Packer) error {
	admins, err := c.unpackAddresses(p)
	if err != nil {
		return err
	}
	managers, err := c.unpackAddresses(p)
	if err != nil {
		return err
	}
	enableds, err := c.unpackAddresses(p)
	if err != nil {
		return err
	}

	c.AdminAddresses = admins
	c.ManagerAddresses = managers
	c.EnabledAddresses = enableds

	return nil
}
