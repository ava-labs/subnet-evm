// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package allowlist

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

// Role mirrors the Solidity enum of the same name.
type Role uint64

// Role constants; they MUST be incremented and MUST NOT be reordered.
const (
	// NoRole - this is equivalent to common.Hash{} and deletes the key from the DB when set
	NoRole Role = iota
	// EnabledRole - allowed to call the precompile
	EnabledRole
	// AdminRole - allowed to both modify the allowlist and call the precompile
	AdminRole
	// ManagerRole - allowed to add and remove only enabled addresses and also call the precompile. (only after Durango)
	ManagerRole

	// firstInvalidRole is useful for loops and other upper-bound checks of
	// Role values. See [Role.IsValid].
	firstInvalidRole
)

// ErrInvalidRole is returned if an invalid role is encountered. If this error
// is returned along with a [Role] value, said value MUST NOT be used; this
// allows functions to return `Role(0), ErrInvalidRole` without worrying about
// misinterpretation as [NoRole].
var ErrInvalidRole = errors.New("invalid role")

// IsValid returns whether or not `r` is a valid value.
func (r Role) IsValid() bool {
	return r < firstInvalidRole
}

// IsEnabled returns true if [r] indicates that it has permission to access the resource.
func (r Role) IsEnabled() bool {
	return r != NoRole && r.IsValid()
}

func (r Role) CanModify(from, target Role) bool {
	switch r {
	case AdminRole:
		return true
	case ManagerRole:
		return from.canBeManaged() && target.canBeManaged()
	default:
		return false
	}
}

// canBeManaged returns whether an account with the `Manager` Role can modify
// another account's Role to or from `r`.
func (r Role) canBeManaged() bool {
	return r == EnabledRole || r == NoRole
}

func (r Role) uint256() *uint256.Int {
	return uint256.NewInt(uint64(r))
}

func (r Role) Bytes() []byte {
	b := r.uint256().Bytes32()
	return b[:]
}

func (r Role) Big() *big.Int {
	return r.uint256().ToBig()
}

func (r Role) Hash() common.Hash {
	return common.Hash(r.uint256().Bytes32())
}

func (r Role) GetSetterFunctionName() (string, error) {
	switch r {
	case AdminRole:
		return "setAdmin", nil
	case ManagerRole:
		return "setManager", nil
	case EnabledRole:
		return "setEnabled", nil
	case NoRole:
		return "setNone", nil
	default:
		return "", ErrInvalidRole
	}
}

// String returns a string representation of [r].
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
		return fmt.Sprintf("UnknownRole[%d]", r)
	}
}

// A Big value may or may not be representable as a uint64. Such types include
// [big.Int] and [uint256.Int].
type Big interface {
	Uint64() uint64
	IsUint64() bool
}

// RoleFromBig converts `u.Uint64()` into a [Role] if `u.IsUint64()` is true.
func RoleFromBig(b Big) (Role, error) {
	r := Role(b.Uint64())
	if !b.IsUint64() || !r.IsValid() {
		return r, ErrInvalidRole
	}
	return r, nil
}

// RoleFromHash converts `h` into a [Role].
func RoleFromHash(h common.Hash) (Role, error) {
	return RoleFromBig(new(uint256.Int).SetBytes32(h[:]))
}
