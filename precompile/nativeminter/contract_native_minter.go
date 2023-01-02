// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nativeminter

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
)

const (
	mintInputAddressSlot = iota
	mintInputAmountSlot

	mintInputLen = common.HashLength + common.HashLength

	MintGasCost = 30_000
)

var (
	// Singleton StatefulPrecompiledContract for minting native assets by permissioned callers.
	ContractNativeMinterPrecompile precompile.StatefulPrecompiledContract = createNativeMinterPrecompile(Address)

	mintSignature = precompile.CalculateFunctionSelector("mintNativeCoin(address,uint256)") // address, amount
	ErrCannotMint = errors.New("non-enabled cannot mint")
)

// GetContractNativeMinterStatus returns the role of [address] for the minter list.
func GetContractNativeMinterStatus(stateDB precompile.StateDB, address common.Address) precompile.AllowListRole {
	return precompile.GetAllowListStatus(stateDB, Address, address)
}

// SetContractNativeMinterStatus sets the permissions of [address] to [role] for the
// minter list. assumes [role] has already been verified as valid.
func SetContractNativeMinterStatus(stateDB precompile.StateDB, address common.Address, role precompile.AllowListRole) {
	precompile.SetAllowListRole(stateDB, Address, address, role)
}

// PackMintInput packs [address] and [amount] into the appropriate arguments for minting operation.
// Assumes that [amount] can be represented by 32 bytes.
func PackMintInput(address common.Address, amount *big.Int) ([]byte, error) {
	// function selector (4 bytes) + input(hash for address + hash for amount)
	res := make([]byte, precompile.SelectorLen+mintInputLen)
	err := precompile.PackOrderedHashesWithSelector(res, mintSignature, []common.Hash{
		address.Hash(),
		common.BigToHash(amount),
	})

	return res, err
}

// UnpackMintInput attempts to unpack [input] into the arguments to the mint precompile
// assumes that [input] does not include selector (omits first 4 bytes in PackMintInput)
func UnpackMintInput(input []byte) (common.Address, *big.Int, error) {
	if len(input) != mintInputLen {
		return common.Address{}, nil, fmt.Errorf("invalid input length for minting: %d", len(input))
	}
	to := common.BytesToAddress(precompile.PackedHash(input, mintInputAddressSlot))
	assetAmount := new(big.Int).SetBytes(precompile.PackedHash(input, mintInputAmountSlot))
	return to, assetAmount, nil
}

// mintNativeCoin checks if the caller is permissioned for minting operation.
// The execution function parses the [input] into native coin amount and receiver address.
func mintNativeCoin(accessibleState precompile.PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = precompile.DeductGas(suppliedGas, MintGasCost); err != nil {
		return nil, 0, err
	}

	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}

	to, amount, err := UnpackMintInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	stateDB := accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to modify it
	callerStatus := precompile.GetAllowListStatus(stateDB, Address, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("%w: %s", ErrCannotMint, caller)
	}

	// if there is no address in the state, create one.
	if !stateDB.Exist(to) {
		stateDB.CreateAccount(to)
	}

	stateDB.AddBalance(to, amount)
	// Return an empty output and the remaining gas
	return []byte{}, remainingGas, nil
}

// createNativeMinterPrecompile returns a StatefulPrecompiledContract with R/W control of an allow list at [precompileAddr] and a native coin minter.
func createNativeMinterPrecompile(precompileAddr common.Address) precompile.StatefulPrecompiledContract {
	enabledFuncs := precompile.CreateAllowListFunctions(precompileAddr)

	mintFunc := precompile.NewStatefulPrecompileFunction(mintSignature, mintNativeCoin)

	enabledFuncs = append(enabledFuncs, mintFunc)
	// Construct the contract with no fallback function.
	contract, err := precompile.NewStatefulPrecompileContract(nil, enabledFuncs)
	// Change this to be returned as an error after refactoring this precompile
	// to use the new precompile template.
	if err != nil {
		panic(err)
	}
	return contract
}
