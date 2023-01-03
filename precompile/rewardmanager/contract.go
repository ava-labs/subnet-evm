// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated
// This file is a generated precompile contract with stubbed abstract functions.

package rewardmanager

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/constants"
	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/vmerrs"

	_ "embed"

	"github.com/ethereum/go-ethereum/common"
)

const (
	AllowFeeRecipientsGasCost      uint64 = (precompile.WriteGasCostPerSlot) + allowlist.ReadAllowListGasCost // write 1 slot + read allow list
	AreFeeRecipientsAllowedGasCost uint64 = allowlist.ReadAllowListGasCost
	CurrentRewardAddressGasCost    uint64 = allowlist.ReadAllowListGasCost
	DisableRewardsGasCost          uint64 = (precompile.WriteGasCostPerSlot) + allowlist.ReadAllowListGasCost // write 1 slot + read allow list
	SetRewardAddressGasCost        uint64 = (precompile.WriteGasCostPerSlot) + allowlist.ReadAllowListGasCost // write 1 slot + read allow list
)

// Singleton StatefulPrecompiledContract and signatures.
var (
	ErrCannotAllowFeeRecipients      = errors.New("non-enabled cannot call allowFeeRecipients")
	ErrCannotAreFeeRecipientsAllowed = errors.New("non-enabled cannot call areFeeRecipientsAllowed")
	ErrCannotCurrentRewardAddress    = errors.New("non-enabled cannot call currentRewardAddress")
	ErrCannotDisableRewards          = errors.New("non-enabled cannot call disableRewards")
	ErrCannotSetRewardAddress        = errors.New("non-enabled cannot call setRewardAddress")

	ErrCannotEnableBothRewards = errors.New("cannot enable both fee recipients and reward address at the same time")
	ErrEmptyRewardAddress      = errors.New("reward address cannot be empty")

	// RewardManagerRawABI contains the raw ABI of RewardManager contract.
	//go:embed contract.abi
	RewardManagerRawABI string

	RewardManagerABI        abi.ABI                                // will be initialized by init function
	RewardManagerPrecompile precompile.StatefulPrecompiledContract // will be initialized by init function

	rewardAddressStorageKey        = common.Hash{'r', 'a', 's', 'k'}
	allowFeeRecipientsAddressValue = common.Hash{'a', 'f', 'r', 'a', 'v'}
)

func init() {
	parsed, err := abi.JSON(strings.NewReader(RewardManagerRawABI))
	if err != nil {
		panic(err)
	}
	RewardManagerABI = parsed
	RewardManagerPrecompile, err = createRewardManagerPrecompile(precompile.RewardManagerAddress)
	if err != nil {
		panic(err)
	}
}

// GetRewardManagerAllowListStatus returns the role of [address] for the RewardManager list.
func GetRewardManagerAllowListStatus(stateDB precompile.StateDB, address common.Address) allowlist.AllowListRole {
	return allowlist.GetAllowListStatus(stateDB, precompile.RewardManagerAddress, address)
}

// SetRewardManagerAllowListStatus sets the permissions of [address] to [role] for the
// RewardManager list. Assumes [role] has already been verified as valid.
func SetRewardManagerAllowListStatus(stateDB precompile.StateDB, address common.Address, role allowlist.AllowListRole) {
	allowlist.SetAllowListRole(stateDB, precompile.RewardManagerAddress, address, role)
}

// PackAllowFeeRecipients packs the function selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackAllowFeeRecipients() ([]byte, error) {
	return RewardManagerABI.Pack("allowFeeRecipients")
}

// EnableAllowFeeRecipients enables fee recipients.
func EnableAllowFeeRecipients(stateDB precompile.StateDB) {
	stateDB.SetState(precompile.RewardManagerAddress, rewardAddressStorageKey, allowFeeRecipientsAddressValue)
}

// DisableRewardAddress disables rewards and burns them by sending to Blackhole Address.
func DisableFeeRewards(stateDB precompile.StateDB) {
	stateDB.SetState(precompile.RewardManagerAddress, rewardAddressStorageKey, constants.BlackholeAddr.Hash())
}

func allowFeeRecipients(accessibleState precompile.PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = precompile.DeductGas(suppliedGas, AllowFeeRecipientsGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// no input provided for this function

	// Allow list is enabled and AllowFeeRecipients is a state-changer function.
	// This part of the code restricts the function to be called only by enabled/admin addresses in the allow list.
	// You can modify/delete this code if you don't want this function to be restricted by the allow list.
	stateDB := accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to modify it
	callerStatus := allowlist.GetAllowListStatus(stateDB, precompile.RewardManagerAddress, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("%w: %s", ErrCannotAllowFeeRecipients, caller)
	}
	// allow list code ends here.

	// this function does not return an output, leave this one as is
	EnableAllowFeeRecipients(stateDB)
	packedOutput := []byte{}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// PackAreFeeRecipientsAllowed packs the include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackAreFeeRecipientsAllowed() ([]byte, error) {
	return RewardManagerABI.Pack("areFeeRecipientsAllowed")
}

// PackAreFeeRecipientsAllowedOutput attempts to pack given isAllowed of type bool
// to conform the ABI outputs.
func PackAreFeeRecipientsAllowedOutput(isAllowed bool) ([]byte, error) {
	return RewardManagerABI.PackOutput("areFeeRecipientsAllowed", isAllowed)
}

func areFeeRecipientsAllowed(accessibleState precompile.PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = precompile.DeductGas(suppliedGas, AreFeeRecipientsAllowedGasCost); err != nil {
		return nil, 0, err
	}
	// no input provided for this function

	stateDB := accessibleState.GetStateDB()
	var output bool
	_, output = GetStoredRewardAddress(stateDB)

	packedOutput, err := PackAreFeeRecipientsAllowedOutput(output)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// PackCurrentRewardAddress packs the include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackCurrentRewardAddress() ([]byte, error) {
	return RewardManagerABI.Pack("currentRewardAddress")
}

// PackCurrentRewardAddressOutput attempts to pack given rewardAddress of type common.Address
// to conform the ABI outputs.
func PackCurrentRewardAddressOutput(rewardAddress common.Address) ([]byte, error) {
	return RewardManagerABI.PackOutput("currentRewardAddress", rewardAddress)
}

// GetStoredRewardAddress returns the current value of the address stored under rewardAddressStorageKey.
// Returns an empty address and true if allow fee recipients is enabled, otherwise returns current reward address and false.
func GetStoredRewardAddress(stateDB precompile.StateDB) (common.Address, bool) {
	val := stateDB.GetState(precompile.RewardManagerAddress, rewardAddressStorageKey)
	return common.BytesToAddress(val.Bytes()), val == allowFeeRecipientsAddressValue
}

// StoredRewardAddress stores the given [val] under rewardAddressStorageKey.
func StoreRewardAddress(stateDB precompile.StateDB, val common.Address) error {
	// if input is empty, return an error
	if val == (common.Address{}) {
		return ErrEmptyRewardAddress
	}
	stateDB.SetState(precompile.RewardManagerAddress, rewardAddressStorageKey, val.Hash())
	return nil
}

// PackSetRewardAddress packs [addr] of type common.Address into the appropriate arguments for setRewardAddress.
// the packed bytes include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackSetRewardAddress(addr common.Address) ([]byte, error) {
	return RewardManagerABI.Pack("setRewardAddress", addr)
}

// UnpackSetRewardAddressInput attempts to unpack [input] into the common.Address type argument
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackSetRewardAddressInput(input []byte) (common.Address, error) {
	res, err := RewardManagerABI.UnpackInput("setRewardAddress", input)
	if err != nil {
		return common.Address{}, err
	}
	unpacked := *abi.ConvertType(res[0], new(common.Address)).(*common.Address)
	return unpacked, nil
}

func setRewardAddress(accessibleState precompile.PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = precompile.DeductGas(suppliedGas, SetRewardAddressGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the SetRewardAddressInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackSetRewardAddressInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// Allow list is enabled and SetRewardAddress is a state-changer function.
	// This part of the code restricts the function to be called only by enabled/admin addresses in the allow list.
	// You can modify/delete this code if you don't want this function to be restricted by the allow list.
	stateDB := accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to modify it
	callerStatus := allowlist.GetAllowListStatus(stateDB, precompile.RewardManagerAddress, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("%w: %s", ErrCannotSetRewardAddress, caller)
	}
	// allow list code ends here.

	if err := StoreRewardAddress(stateDB, inputStruct); err != nil {
		return nil, remainingGas, err
	}
	// this function does not return an output, leave this one as is
	packedOutput := []byte{}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

func currentRewardAddress(accessibleState precompile.PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = precompile.DeductGas(suppliedGas, CurrentRewardAddressGasCost); err != nil {
		return nil, 0, err
	}

	// no input provided for this function
	stateDB := accessibleState.GetStateDB()
	output, _ := GetStoredRewardAddress(stateDB)
	packedOutput, err := PackCurrentRewardAddressOutput(output)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// PackDisableRewards packs the include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackDisableRewards() ([]byte, error) {
	return RewardManagerABI.Pack("disableRewards")
}

func disableRewards(accessibleState precompile.PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = precompile.DeductGas(suppliedGas, DisableRewardsGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// no input provided for this function

	// Allow list is enabled and DisableRewards is a state-changer function.
	// This part of the code restricts the function to be called only by enabled/admin addresses in the allow list.
	// You can modify/delete this code if you don't want this function to be restricted by the allow list.
	stateDB := accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to modify it
	callerStatus := allowlist.GetAllowListStatus(stateDB, precompile.RewardManagerAddress, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("%w: %s", ErrCannotDisableRewards, caller)
	}
	// allow list code ends here.
	DisableFeeRewards(stateDB)
	// this function does not return an output, leave this one as is
	packedOutput := []byte{}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// createRewardManagerPrecompile returns a StatefulPrecompiledContract with getters and setters for the precompile.
// Access to the getters/setters is controlled by an allow list for [precompileAddr].
func createRewardManagerPrecompile(precompileAddr common.Address) (precompile.StatefulPrecompiledContract, error) {
	var functions []*precompile.StatefulPrecompileFunction
	functions = append(functions, allowlist.CreateAllowListFunctions(precompileAddr)...)
	abiFunctionMap := map[string]precompile.RunStatefulPrecompileFunc{
		"allowFeeRecipients":      allowFeeRecipients,
		"areFeeRecipientsAllowed": areFeeRecipientsAllowed,
		"currentRewardAddress":    currentRewardAddress,
		"disableRewards":          disableRewards,
		"setRewardAddress":        setRewardAddress,
	}

	for name, function := range abiFunctionMap {
		method, ok := RewardManagerABI.Methods[name]
		if !ok {
			return nil, fmt.Errorf("given method (%s) does not exist in the ABI", name)
		}
		functions = append(functions, precompile.NewStatefulPrecompileFunction(method.ID, function))
	}

	// Construct the contract with no fallback function.
	return precompile.NewStatefulPrecompileContract(nil, functions)
}
