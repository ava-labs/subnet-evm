// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
)

const (
	SetAdminFuncKey      = "setAdmin"
	SetEnabledFuncKey    = "setEnabled"
	SetNoneFuncKey       = "setNone"
	ReadAllowListFuncKey = "readAllowList"

	ModifyAllowListGasCost = WriteGasCostPerSlot
	ReadAllowListGasCost   = ReadGasCostPerSlot
)

var (
	AllowListNoRole  AllowListRole = AllowListRole(common.BigToHash(big.NewInt(0))) // No role assigned - this is equivalent to common.Hash{} and deletes the key from the DB when set
	AllowListEnabled AllowListRole = AllowListRole(common.BigToHash(big.NewInt(1))) // Deployers are allowed to create new contracts
	AllowListAdmin   AllowListRole = AllowListRole(common.BigToHash(big.NewInt(2))) // Admin - allowed to modify both the admin and deployer list as well as deploy contracts

	// AllowList function signatures
	setAdminSignature      = CalculateFunctionSelector("setAdmin(address)")
	setEnabledSignature    = CalculateFunctionSelector("setEnabled(address)")
	setNoneSignature       = CalculateFunctionSelector("setNone(address)")
	readAllowListSignature = CalculateFunctionSelector("readAllowList(address)")
	// Error returned when an invalid write is attempted
	ErrCannotModifyAllowList = errors.New("non-admin cannot modify allow list")

	allowListInputLen = common.HashLength
)

// AllowListConfig specifies the initial set of allow list admins.
type AllowListConfig struct {
	AllowListAdmins  []common.Address `json:"adminAddresses"`
	EnabledAddresses []common.Address `json:"enabledAddresses"` // initial enabled addresses
}

// Configure initializes the address space of [precompileAddr] by initializing the role of each of
// the addresses in [AllowListAdmins].
func (c *AllowListConfig) Configure(state StateDB, precompileAddr common.Address) error {
	for _, enabledAddr := range c.EnabledAddresses {
		SetAllowListRole(state, precompileAddr, enabledAddr, AllowListEnabled)
	}
	for _, adminAddr := range c.AllowListAdmins {
		SetAllowListRole(state, precompileAddr, adminAddr, AllowListAdmin)
	}
	return nil
}

// Equal returns true iff [other] has the same admins in the same order in its allow list.
func (c *AllowListConfig) Equal(other *AllowListConfig) bool {
	if other == nil {
		return false
	}
	if !areEqualAddressLists(c.AllowListAdmins, other.AllowListAdmins) {
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
func (c *AllowListConfig) Verify() error {
	// return early if either list is empty
	if len(c.EnabledAddresses) == 0 || len(c.AllowListAdmins) == 0 {
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

	for _, adminAddr := range c.AllowListAdmins {
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

// GetAllowListStatus returns the allow list role of [address] for the precompile
// at [precompileAddr]
func GetAllowListStatus(state StateDB, precompileAddr common.Address, address common.Address) AllowListRole {
	// Generate the state key for [address]
	addressKey := address.Hash()
	return AllowListRole(state.GetState(precompileAddr, addressKey))
}

// SetAllowListRole sets the permissions of [address] to [role] for the precompile
// at [precompileAddr].
// assumes [role] has already been verified as valid.
func SetAllowListRole(stateDB StateDB, precompileAddr, address common.Address, role AllowListRole) {
	// Generate the state key for [address]
	addressKey := address.Hash()
	// Assign [role] to the address
	// This stores the [role] in the contract storage with address [precompileAddr]
	// and [addressKey] hash. It means that any reusage of the [addressKey] for different value
	// conflicts with the same slot [role] is stored.
	// Precompile implementations must use a different key than [addressKey]
	stateDB.SetState(precompileAddr, addressKey, common.Hash(role))
}

// PackModifyAllowList packs [address] and [role] into the appropriate arguments for modifying the allow list.
// Note: [role] is not packed in the input value returned, but is instead used as a selector for the function
// selector that should be encoded in the input.
func PackModifyAllowList(address common.Address, role AllowListRole) ([]byte, error) {
	// function selector (4 bytes) + hash for address
	input := make([]byte, 0, selectorLen+common.HashLength)

	switch role {
	case AllowListAdmin:
		input = append(input, setAdminSignature...)
	case AllowListEnabled:
		input = append(input, setEnabledSignature...)
	case AllowListNoRole:
		input = append(input, setNoneSignature...)
	default:
		return nil, fmt.Errorf("cannot pack modify list input with invalid role: %s", role)
	}

	input = append(input, address.Hash().Bytes()...)
	return input, nil
}

// PackReadAllowList packs [address] into the input data to the read allow list function
func PackReadAllowList(address common.Address) []byte {
	input := make([]byte, 0, selectorLen+common.HashLength)
	input = append(input, readAllowListSignature...)
	input = append(input, address.Hash().Bytes()...)
	return input
}

// createAllowListRoleSetter returns an execution function for setting the allow list status of the input address argument to [role].
// This execution function is speciifc to [precompileAddr].
func createAllowListRoleSetter(precompileAddr common.Address, role AllowListRole) RunStatefulPrecompileFunc {
	return func(evm PrecompileAccessibleState, callerAddr, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
		if remainingGas, err = DeductGas(suppliedGas, ModifyAllowListGasCost); err != nil {
			return nil, 0, err
		}

		if len(input) != allowListInputLen {
			return nil, remainingGas, fmt.Errorf("invalid input length for modifying allow list: %d", len(input))
		}

		modifyAddress := common.BytesToAddress(input)

		if readOnly {
			return nil, remainingGas, vmerrs.ErrWriteProtection
		}

		stateDB := evm.GetStateDB()

		// Verify that the caller is in the allow list and therefore has the right to modify it
		callerStatus := GetAllowListStatus(stateDB, precompileAddr, callerAddr)
		if !callerStatus.IsAdmin() {
			return nil, remainingGas, fmt.Errorf("%w: %s", ErrCannotModifyAllowList, callerAddr)
		}

		SetAllowListRole(stateDB, precompileAddr, modifyAddress, role)
		// Return an empty output and the remaining gas
		return []byte{}, remainingGas, nil
	}
}

// createReadAllowList returns an execution function that reads the allow list for the given [precompileAddr].
// The execution function parses the input into a single address and returns the 32 byte hash that specifies the
// designated role of that address
func createReadAllowList(precompileAddr common.Address) RunStatefulPrecompileFunc {
	return func(evm PrecompileAccessibleState, callerAddr common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
		if remainingGas, err = DeductGas(suppliedGas, ReadAllowListGasCost); err != nil {
			return nil, 0, err
		}

		if len(input) != allowListInputLen {
			return nil, remainingGas, fmt.Errorf("invalid input length for read allow list: %d", len(input))
		}

		readAddress := common.BytesToAddress(input)
		role := GetAllowListStatus(evm.GetStateDB(), precompileAddr, readAddress)
		roleBytes := common.Hash(role).Bytes()
		return roleBytes, remainingGas, nil
	}
}

// createAllowListPrecompile returns a StatefulPrecompiledContract with R/W control of an allow list at [precompileAddr]
func createAllowListPrecompile(precompileAddr common.Address) StatefulPrecompiledContract {
	// Construct the contract with no fallback function.
	allowListFuncs := CreateAllowListFunctions(precompileAddr)
	contract, err := NewStatefulPrecompileContract(nil, allowListFuncs)
	// TODO Change this to be returned as an error after refactoring this precompile
	// to use the new precompile template.
	if err != nil {
		panic(err)
	}
	return contract
}

func CreateAllowListFunctions(precompileAddr common.Address) []*StatefulPrecompileFunction {
	setAdmin := NewStatefulPrecompileFunction(setAdminSignature, createAllowListRoleSetter(precompileAddr, AllowListAdmin))
	setEnabled := NewStatefulPrecompileFunction(setEnabledSignature, createAllowListRoleSetter(precompileAddr, AllowListEnabled))
	setNone := NewStatefulPrecompileFunction(setNoneSignature, createAllowListRoleSetter(precompileAddr, AllowListNoRole))
	read := NewStatefulPrecompileFunction(readAllowListSignature, createReadAllowList(precompileAddr))

	return []*StatefulPrecompileFunction{setAdmin, setEnabled, setNone, read}
}
