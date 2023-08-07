// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package allowlist

import "github.com/ethereum/go-ethereum/common"

// 1. NoRole - this is equivalent to common.Hash{} and deletes the key from the DB when set
// 2. EnabledRole - allowed to call the precompile
// 3. Admin - allowed to both modify the allowlist and call the precompile
// 4. Manager - allowed to add and remove only enabled address (only after DUpgrade), with enabled permissions.
var (
	// NoRole - this is equivalent to common.Hash{} and deletes the key from the DB when set.
	NoRole = Role(common.BigToHash(common.Big0))
	// EnabledRole - allowed to call the precompile.
	EnabledRole = Role(common.BigToHash(common.Big1))
	// Admin - allowed to both modify the allowlist and call the precompile.
	AdminRole = Role(common.BigToHash(common.Big2))
	// Manager - allowed to add and remove only enabled addresses, with being able to call the precompile (i.e enabled). Activated only after DUpgrade.
	ManagerRole = Role(common.BigToHash(common.Big3))
)

// Enum constants for valid Role
type Role common.Hash

// IsNoRole returns true if [s] indicates no specific role.
func (r Role) IsNoRole(isManagerActivated bool) bool {
	switch r {
	case NoRole:
		return true
	case ManagerRole:
		return !isManagerActivated
	default:
		return false
	}
}

// IsAdmin returns true if [s] indicates the permission to modify the allow list.
func (r Role) IsAdmin() bool {
	switch r {
	case AdminRole:
		return true
	default:
		return false
	}
}

// IsEnabled returns true if [s] indicates that it has permission to access the resource.
func (r Role) IsEnabled(isManagerActivated bool) bool {
	switch r {
	case AdminRole, EnabledRole:
		return true
	case ManagerRole:
		return isManagerActivated
	default:
		return false
	}
}

// IsManager returns true if [s] indicates that it has permission to add and remove
// addresses from the allow list.
func (r Role) IsManager(isManagerActivated bool) bool {
	switch r {
	case AdminRole:
		return true
	case ManagerRole:
		return isManagerActivated
	default:
		return false
	}
}

func (r Role) CanModify(isManagerActivated bool, from, target Role) bool {
	switch r {
	case AdminRole:
		return true
	case ManagerRole:
		return isManagerActivated &&
			(from == EnabledRole && target == NoRole) || (from == NoRole && target == EnabledRole)
	default:
		return false
	}
}

// String returns a string representation of [s].
func (r Role) String() string {
	switch r {
	case NoRole:
		return "NoRole"
	case EnabledRole:
		return "EnabledRole"
	case ManagerRole:
		return "ManagerRole"
	case AdminRole:
		return "AdminRole"
	default:
		return "UnknownRole"
	}
}
