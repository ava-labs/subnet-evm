// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package allowlist

import (
	_ "embed"
	"errors"
	"fmt"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
)

// AllowList is an abstraction that allows other precompiles to manage
// which addresses can call the precompile by maintaining an allowlist
// in the storage trie.

const (
	ModifyAllowListGasCost = contract.WriteGasCostPerSlot
	ReadAllowListGasCost   = contract.ReadGasCostPerSlot

	allowListInputLen = common.HashLength
)

var (
	// AllowList function signatures
	setAdminSignature   = contract.CalculateFunctionSelector("setAdmin(address)")
	setManagerSignature = contract.CalculateFunctionSelector("setManager(address)")
	setEnabledSignature = contract.CalculateFunctionSelector("setEnabled(address)")
	setNoneSignature    = contract.CalculateFunctionSelector("setNone(address)")
	// Error returned when an invalid write is attempted
	ErrCannotModifyAllowList = errors.New("cannot modify allow list")

	// AllowListRawABI contains the raw ABI of AllowList library interface.
	//go:embed allowlist.abi
	AllowListRawABI string

	AllowListABI = contract.ParseABI(AllowListRawABI)
)

// GetAllowListStatus returns the allow list role of [address] for the precompile
// at [precompileAddr]
func GetAllowListStatus(state contract.StateDB, precompileAddr common.Address, address common.Address) Role {
	// Generate the state key for [address]
	addressKey := address.Hash()
	return Role(state.GetState(precompileAddr, addressKey))
}

// SetAllowListRole sets the permissions of [address] to [role] for the precompile
// at [precompileAddr].
// assumes [role] has already been verified as valid.
func SetAllowListRole(stateDB contract.StateDB, precompileAddr, address common.Address, role Role) {
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
func PackModifyAllowList(address common.Address, role Role) ([]byte, error) {
	// function selector (4 bytes) + hash for address
	input := make([]byte, 0, contract.SelectorLen+common.HashLength)

	switch role {
	case AdminRole:
		input = append(input, setAdminSignature...)
	case ManagerRole:
		input = append(input, setManagerSignature...)
	case EnabledRole:
		input = append(input, setEnabledSignature...)
	case NoRole:
		input = append(input, setNoneSignature...)
	default:
		return nil, fmt.Errorf("cannot pack modify list input with invalid role: %s", role)
	}

	input = append(input, address.Hash().Bytes()...)
	return input, nil
}

func getAllowListFunctionSelector(role Role) (string, error) {
	switch role {
	case AdminRole:
		return "setAdmin", nil
	case ManagerRole:
		return "setManager", nil
	case EnabledRole:
		return "setEnabled", nil
	case NoRole:
		return "setNone", nil
	default:
		return "", fmt.Errorf("unknown role: %s", role)
	}
}

func PackModifyAllowListV2(address common.Address, role Role) ([]byte, error) {
	funcName, err := getAllowListFunctionSelector(role)
	if err != nil {
		return nil, fmt.Errorf("cannot pack modify list input with invalid role: %s", role)
	}

	return AllowListABI.Pack(funcName, address)
}

func UnpackModifyAllowListInput(input []byte, r Role) (common.Address, error) {
	if len(input) != allowListInputLen {
		return common.Address{}, fmt.Errorf("invalid input length for modifying allow list: %d", len(input))
	}

	funcName, err := getAllowListFunctionSelector(r)
	if err != nil {
		return common.Address{}, err
	}

	var modifyAddress common.Address
	err = AllowListABI.UnpackInputIntoInterface(&modifyAddress, funcName, input)
	return modifyAddress, err
}

// createAllowListRoleSetter returns an execution function for setting the allow list status of the input address argument to [role].
// This execution function is speciifc to [precompileAddr].
func createAllowListRoleSetter(precompileAddr common.Address, role Role) contract.RunStatefulPrecompileFunc {
	return func(evm contract.AccessibleState, callerAddr, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
		if remainingGas, err = contract.DeductGas(suppliedGas, ModifyAllowListGasCost); err != nil {
			return nil, 0, err
		}

		modifyAddress, err := UnpackModifyAllowListInput(input, role)

		if err != nil {
			return nil, remainingGas, err
		}

		if readOnly {
			return nil, remainingGas, vmerrs.ErrWriteProtection
		}

		stateDB := evm.GetStateDB()

		// Verify that the caller is an admin with permission to modify the allow list
		callerStatus := GetAllowListStatus(stateDB, precompileAddr, callerAddr)
		// Verify that the address we are trying to modify has a status that allows it to be modified
		modifyStatus := GetAllowListStatus(stateDB, precompileAddr, modifyAddress)
		if !callerStatus.CanModify(modifyStatus, role) {
			return nil, remainingGas, fmt.Errorf("%w: modify address: %s, from role: %s, to role: %s", ErrCannotModifyAllowList, callerAddr, modifyStatus, role)
		}
		SetAllowListRole(stateDB, precompileAddr, modifyAddress, role)
		// Return an empty output and the remaining gas
		return []byte{}, remainingGas, nil
	}
}

// PackReadAllowList packs [address] into the input data to the read allow list function
func PackReadAllowList(address common.Address) ([]byte, error) {
	return AllowListABI.Pack("readAllowList", address)
}

func UnpackReadAllowListInput(input []byte, skipLenCheck bool) (common.Address, error) {
	if !skipLenCheck && len(input) != allowListInputLen {
		return common.Address{}, fmt.Errorf("invalid input length for read allow list: %d", len(input))
	}

	var modifyAddress common.Address
	err := AllowListABI.UnpackInputIntoInterface(&modifyAddress, "readAllowList", input)
	return modifyAddress, err
}

// createReadAllowList returns an execution function that reads the allow list for the given [precompileAddr].
// The execution function parses the input into a single address and returns the 32 byte hash that specifies the
// designated role of that address
func createReadAllowList(precompileAddr common.Address) contract.RunStatefulPrecompileFunc {
	return func(evm contract.AccessibleState, callerAddr common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
		if remainingGas, err = contract.DeductGas(suppliedGas, ReadAllowListGasCost); err != nil {
			return nil, 0, err
		}

		skipLenCheck := contract.IsDUpgradeActivated(evm)
		readAddress, err := UnpackReadAllowListInput(input, skipLenCheck)
		if err != nil {
			return nil, remainingGas, err
		}

		role := GetAllowListStatus(evm.GetStateDB(), precompileAddr, readAddress)
		return role.Bytes(), remainingGas, nil
	}
}

// CreateAllowListPrecompile returns a StatefulPrecompiledContract with R/W control of an allow list at [precompileAddr]
func CreateAllowListPrecompile(precompileAddr common.Address) contract.StatefulPrecompiledContract {
	// Construct the contract with no fallback function.
	allowListFuncs := CreateAllowListFunctions(precompileAddr)
	contract, err := contract.NewStatefulPrecompileContract(nil, allowListFuncs)
	if err != nil {
		panic(err)
	}
	return contract
}

func CreateAllowListFunctions(precompileAddr common.Address) []*contract.StatefulPrecompileFunction {
	var functions []*contract.StatefulPrecompileFunction

	type precompileFn struct {
		fn        contract.RunStatefulPrecompileFunc
		activator contract.ActivationFunc
	}

	abiFunctionMap := map[string]precompileFn{
		"setAdmin": {
			fn: createAllowListRoleSetter(precompileAddr, AdminRole),
		},
		"setEnabled": {
			fn: createAllowListRoleSetter(precompileAddr, EnabledRole),
		},
		"setNone": {
			fn: createAllowListRoleSetter(precompileAddr, NoRole),
		},
		"readAllowList": {
			fn: createReadAllowList(precompileAddr),
		},
		"setManager": {
			fn:        createAllowListRoleSetter(precompileAddr, ManagerRole),
			activator: contract.IsDUpgradeActivated,
		},
	}

	for name, function := range abiFunctionMap {
		method, ok := AllowListABI.Methods[name]
		if !ok {
			panic(fmt.Errorf("given method (%s) does not exist in the ABI", name))
		}
		var spFn *contract.StatefulPrecompileFunction
		if function.activator != nil {
			spFn = contract.NewStatefulPrecompileFunctionWithActivator(method.ID, function.fn, function.activator)
		} else {
			spFn = contract.NewStatefulPrecompileFunction(method.ID, function.fn)
		}
		functions = append(functions, spFn)
	}

	return functions
}
