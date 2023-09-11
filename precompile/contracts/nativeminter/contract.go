// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nativeminter

import (
	_ "embed"
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
)

const (
	mintInputAddressSlot = iota
	mintInputAmountSlot

	mintInputLen = common.HashLength + common.HashLength

	MintGasCost = 30_000
)

type MintNativeCoinInput struct {
	Addr   common.Address
	Amount *big.Int
}

var (
	// Singleton StatefulPrecompiledContract for minting native assets by permissioned callers.
	ContractNativeMinterPrecompile contract.StatefulPrecompiledContract = createNativeMinterPrecompile()

	mintSignature = contract.CalculateFunctionSelector("mintNativeCoin(address,uint256)") // address, amount
	ErrCannotMint = errors.New("non-enabled cannot mint")
	ErrInvalidLen = errors.New("invalid input length for minting")

	// NativeMinterV2RawABI contains the raw ABI of NativeMinterV2 contract.
	//go:embed contract.abi
	NativeMinterV2RawABI string

	NativeMinterV2ABI = contract.ParseABI(NativeMinterV2RawABI)
)

// GetContractNativeMinterStatus returns the role of [address] for the minter list.
func GetContractNativeMinterStatus(stateDB contract.StateDB, address common.Address) allowlist.Role {
	return allowlist.GetAllowListStatus(stateDB, ContractAddress, address)
}

// SetContractNativeMinterStatus sets the permissions of [address] to [role] for the
// minter list. assumes [role] has already been verified as valid.
func SetContractNativeMinterStatus(stateDB contract.StateDB, address common.Address, role allowlist.Role) {
	allowlist.SetAllowListRole(stateDB, ContractAddress, address, role)
}

// PackMintInput packs [address] and [amount] into the appropriate arguments for minting operation.
// Assumes that [amount] can be represented by 32 bytes.
func PackMintInput(address common.Address, amount *big.Int) ([]byte, error) {
	// function selector (4 bytes) + input(hash for address + hash for amount)
	res := make([]byte, contract.SelectorLen+mintInputLen)
	err := contract.PackOrderedHashesWithSelector(res, mintSignature, []common.Hash{
		address.Hash(),
		common.BigToHash(amount),
	})

	return res, err
}

// UnpackMintInput attempts to unpack [input] into the arguments to the mint precompile
// assumes that [input] does not include selector (omits first 4 bytes in PackMintInput)
func UnpackMintInput(input []byte) (common.Address, *big.Int, error) {
	if len(input) != mintInputLen {
		return common.Address{}, nil, fmt.Errorf("%w: %d", ErrInvalidLen, len(input))
	}
	to := common.BytesToAddress(contract.PackedHash(input, mintInputAddressSlot))
	assetAmount := new(big.Int).SetBytes(contract.PackedHash(input, mintInputAmountSlot))
	return to, assetAmount, nil
}

// mintNativeCoin checks if the caller is permissioned for minting operation.
// The execution function parses the [input] into native coin amount and receiver address.
func mintNativeCoin(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, MintGasCost); err != nil {
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
	// Verify that the caller is in the allow list and therefore has the right to call this function.
	callerStatus := allowlist.GetAllowListStatus(stateDB, ContractAddress, caller)
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

// createNativeMinterPrecompile returns a StatefulPrecompiledContract for native coin minting. The precompile
// is accessed controlled by an allow list at [precompileAddr].
func createNativeMinterPrecompile() contract.StatefulPrecompiledContract {
	enabledFuncs := allowlist.CreateAllowListFunctions(ContractAddress)

	mintFunc := contract.NewStatefulPrecompileFunction(mintSignature, mintNativeCoin)

	enabledFuncs = append(enabledFuncs, mintFunc)
	// Construct the contract with no fallback function.
	contract, err := contract.NewStatefulPrecompileContract(nil, enabledFuncs)
	// TODO: Change this to be returned as an error after refactoring this precompile
	// to use the new precompile template.
	if err != nil {
		panic(err)
	}
	return contract
}

// UnpackMintNativeCoinInput attempts to unpack [input] as MintNativeCoinInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
// if [doLenCheck] is true, it will return an error if the length of [input] is not [mintInputLen]
func UnpackMintNativeCoinV2Input(input []byte, doLenCheck bool) (MintNativeCoinInput, error) {
	// Initially we had this check to ensure that the input was the correct length.
	// However solidity does not always pack the input to the correct length, and allows
	// for extra padding bytes to be added to the end of the input. Therefore, we have removed
	// this check with the DUpgrade. We still need to keep this check for backwards compatibility.
	if doLenCheck && len(input) != mintInputLen {
		return MintNativeCoinInput{}, fmt.Errorf("%w: %d", ErrInvalidLen, len(input))
	}
	inputStruct := MintNativeCoinInput{}
	err := NativeMinterV2ABI.UnpackInputIntoInterface(&inputStruct, "mintNativeCoin", input)

	return inputStruct, err
}

// PackMintNativeCoin packs [inputStruct] of type MintNativeCoinInput into the appropriate arguments for mintNativeCoin.
func PackMintNativeCoinV2(inputStruct MintNativeCoinInput) ([]byte, error) {
	return NativeMinterV2ABI.Pack("mintNativeCoin", inputStruct.Addr, inputStruct.Amount)
}
