// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated
// This file is a generated precompile contract with stubbed abstract functions.

package precompile

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/vmerrs"

	"github.com/ethereum/go-ethereum/common"
)

const (
	AllowFeeRecipientsGasCost      uint64 = (writeGasCostPerSlot * 2) + ReadAllowListGasCost // 2 slots + read allow list
	AreFeeRecipientsAllowedGasCost uint64 = readGasCostPerSlot
	CurrentRewardAddressGasCost    uint64 = readGasCostPerSlot
	DisableRewardsGasCost          uint64 = (writeGasCostPerSlot * 2) + ReadAllowListGasCost // 2 slots + read allow list
	SetRewardAddressGasCost        uint64 = (writeGasCostPerSlot * 2) + ReadAllowListGasCost // 2 slots + read allow list

	// RewardManagerRawABI contains the raw ABI of RewardManager contract.
	RewardManagerRawABI = "[{\"inputs\":[],\"name\":\"allowFeeRecipients\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"areFeeRecipientsAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isAllowed\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentRewardAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"rewardAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"disableRewards\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"readAllowList\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setEnabled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setNone\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setRewardAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
)

// Singleton StatefulPrecompiledContract and signatures.
var (
	_ StatefulPrecompileConfig = &RewardManagerConfig{}

	ErrCannotAllowFeeRecipients      = errors.New("non-enabled cannot allowFeeRecipients")
	ErrCannotAreFeeRecipientsAllowed = errors.New("non-enabled cannot areFeeRecipientsAllowed")
	ErrCannotCurrentRewardAddress    = errors.New("non-enabled cannot currentRewardAddress")
	ErrCannotDisableRewards          = errors.New("non-enabled cannot disableRewards")
	ErrCannotSetRewardAddress        = errors.New("non-enabled cannot setRewardAddress")

	ErrCannotEnableBothRewards = errors.New("cannot enable both fee recipients and reward address at the same time")
	ErrEmptyRewardAddress      = errors.New("reward address cannot be empty")

	RewardManagerABI        abi.ABI                     // will be initialized by init function
	RewardManagerPrecompile StatefulPrecompiledContract // will be initialized by init function

	allowFeeRecipientsStorageKey = common.Hash{'a', 'f', 'r', 's', 'k'}
	rewardAddressStorageKey      = common.Hash{'r', 'a', 's', 'k'}
)

type InitialRewardConfig struct {
	AllowFeeRecipients bool           `json:"allowFeeRecipients"`
	RewardAddress      common.Address `json:"rewardAddress,omitempty"`
}

func (i *InitialRewardConfig) Verify() error {
	switch {
	case i.AllowFeeRecipients && i.RewardAddress != (common.Address{}):
		return ErrCannotEnableBothRewards
		// shall we also check blackhole address here?
	default:
		return nil
	}
}

func (c *InitialRewardConfig) Equal(other *InitialRewardConfig) bool {
	if other == nil {
		return false
	}

	return c.AllowFeeRecipients == other.AllowFeeRecipients && c.RewardAddress == other.RewardAddress
}

// RewardManagerConfig implements the StatefulPrecompileConfig
// interface while adding in the RewardManager specific precompile address.
type RewardManagerConfig struct {
	AllowListConfig
	UpgradeableConfig
	InitialRewardConfig *InitialRewardConfig `json:"initialRewardConfig,omitempty"`
}

func init() {
	parsed, err := abi.JSON(strings.NewReader(RewardManagerRawABI))
	if err != nil {
		panic(err)
	}
	RewardManagerABI = parsed
	RewardManagerPrecompile = createRewardManagerPrecompile(RewardManagerAddress)
}

// NewRewardManagerConfig returns a config for a network upgrade at [blockTimestamp] that enables
// RewardManager with the given [admins] and [enableds] as members of the allowlist with [initialConfig] as initial rewards config if specified.
func NewRewardManagerConfig(blockTimestamp *big.Int, admins []common.Address, enableds []common.Address, initialConfig *InitialRewardConfig) *RewardManagerConfig {
	return &RewardManagerConfig{
		AllowListConfig: AllowListConfig{
			AllowListAdmins:  admins,
			EnabledAddresses: enableds,
		},
		UpgradeableConfig:   UpgradeableConfig{BlockTimestamp: blockTimestamp},
		InitialRewardConfig: initialConfig,
	}
}

// NewDisableRewardManagerConfig returns config for a network upgrade at [blockTimestamp]
// that disables RewardManager.
func NewDisableRewardManagerConfig(blockTimestamp *big.Int) *RewardManagerConfig {
	return &RewardManagerConfig{
		UpgradeableConfig: UpgradeableConfig{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

// Equal returns true if [s] is a [*RewardManagerConfig] and it has been configured identical to [c].
func (c *RewardManagerConfig) Equal(s StatefulPrecompileConfig) bool {
	// typecast before comparison
	other, ok := (s).(*RewardManagerConfig)
	if !ok {
		return false
	}
	// CUSTOM CODE STARTS HERE
	// modify this boolean accordingly with your custom RewardManagerConfig, to check if [other] and the current [c] are equal
	// if RewardManagerConfig contains only UpgradeableConfig  and AllowListConfig  you can skip modifying it.
	equals := c.UpgradeableConfig.Equal(&other.UpgradeableConfig) && c.AllowListConfig.Equal(&other.AllowListConfig)
	if !equals {
		return false
	}

	if c.InitialRewardConfig == nil {
		return other.InitialRewardConfig == nil
	}

	return c.InitialRewardConfig.Equal(other.InitialRewardConfig)
}

// Address returns the address of the RewardManager. Addresses reside under the precompile/params.go
// Select a non-conflicting address and set it in the params.go.
func (c *RewardManagerConfig) Address() common.Address {
	return RewardManagerAddress
}

// Configure configures [state] with the initial configuration.
func (c *RewardManagerConfig) Configure(_ ChainConfig, state StateDB, _ BlockContext) {
	c.AllowListConfig.Configure(state, RewardManagerAddress)
	// CUSTOM CODE STARTS HERE
	// configure the RewardManager with the initial configuration
	if c.InitialRewardConfig != nil {
		// set the initial reward config
		if c.InitialRewardConfig.AllowFeeRecipients {
			StoreAllowFeeRecipients(state, true)
		} else if c.InitialRewardConfig.RewardAddress != (common.Address{}) {
			StoreRewardAddress(state, c.InitialRewardConfig.RewardAddress)
		}
	}
}

// Contract returns the singleton stateful precompiled contract to be used for RewardManager.
func (c *RewardManagerConfig) Contract() StatefulPrecompiledContract {
	return RewardManagerPrecompile
}

func (c *RewardManagerConfig) Verify() error {
	if err := c.AllowListConfig.Verify(); err != nil {
		return err
	}
	if c.InitialRewardConfig != nil {
		return c.InitialRewardConfig.Verify()
	}
	return nil
}

// GetRewardManagerAllowListStatus returns the role of [address] for the RewardManager list.
func GetRewardManagerAllowListStatus(stateDB StateDB, address common.Address) AllowListRole {
	return getAllowListStatus(stateDB, RewardManagerAddress, address)
}

// SetRewardManagerAllowListStatus sets the permissions of [address] to [role] for the
// RewardManager list. Assumes [role] has already been verified as valid.
func SetRewardManagerAllowListStatus(stateDB StateDB, address common.Address, role AllowListRole) {
	setAllowListRole(stateDB, RewardManagerAddress, address, role)
}

// PackAllowFeeRecipients packs the include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackAllowFeeRecipients() ([]byte, error) {
	return RewardManagerABI.Pack("allowFeeRecipients")
}

// GetStoredAllowFeeRecipients returns the current value of the stored allowFeeRecipients flag.
func GetStoredAllowFeeRecipients(stateDB StateDB) bool {
	val := stateDB.GetState(RewardManagerAddress, allowFeeRecipientsStorageKey)
	return hashToBool(val)
}

// StoreAllowFeeRecipients stores the given [val] under allowFeeRecipientsStoragekey.
func StoreAllowFeeRecipients(stateDB StateDB, val bool) {
	stateDB.SetState(RewardManagerAddress, allowFeeRecipientsStorageKey, boolToHash(val))
}

func allowFeeRecipients(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = deductGas(suppliedGas, AllowFeeRecipientsGasCost); err != nil {
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
	callerStatus := getAllowListStatus(stateDB, RewardManagerAddress, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("%w: %s", ErrCannotAllowFeeRecipients, caller)
	}
	// allow list code ends here.

	// CUSTOM CODE STARTS HERE
	// this function does not return an output, leave this one as is
	if GetStoredRewardAddress(stateDB) != (common.Address{}) {
		// reset stored reward address first
		StoreRewardAddress(stateDB, common.Address{})
	}
	StoreAllowFeeRecipients(stateDB, true)
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

func areFeeRecipientsAllowed(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = deductGas(suppliedGas, AreFeeRecipientsAllowedGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// no input provided for this function

	// Allow list is enabled and AreFeeRecipientsAllowed is a state-changer function.
	// This part of the code restricts the function to be called only by enabled/admin addresses in the allow list.
	// You can modify/delete this code if you don't want this function to be restricted by the allow list.
	stateDB := accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to modify it
	callerStatus := getAllowListStatus(stateDB, RewardManagerAddress, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("%w: %s", ErrCannotAreFeeRecipientsAllowed, caller)
	}
	// allow list code ends here.

	// CUSTOM CODE STARTS HERE
	var output bool // CUSTOM CODE FOR AN OUTPUT
	output = GetStoredAllowFeeRecipients(stateDB)
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
func GetStoredRewardAddress(stateDB StateDB) common.Address {
	val := stateDB.GetState(RewardManagerAddress, rewardAddressStorageKey)
	return common.BytesToAddress(val.Bytes())
}

// StoredRewardAddress stores the given [val] under rewardAddressStorageKey.
func StoreRewardAddress(stateDB StateDB, val common.Address) {
	stateDB.SetState(RewardManagerAddress, rewardAddressStorageKey, val.Hash())
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

func setRewardAddress(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = deductGas(suppliedGas, SetRewardAddressGasCost); err != nil {
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
	callerStatus := getAllowListStatus(stateDB, RewardManagerAddress, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("%w: %s", ErrCannotSetRewardAddress, caller)
	}
	// allow list code ends here.

	// CUSTOM CODE STARTS HERE
	// if input is empty, return an error
	if inputStruct == (common.Address{}) {
		return nil, remainingGas, ErrEmptyRewardAddress
	}
	// reset stored allow fee recipients flag only if it's already set to true
	if GetStoredAllowFeeRecipients(stateDB) {
		StoreAllowFeeRecipients(stateDB, false)
	}

	StoreRewardAddress(stateDB, inputStruct)
	// this function does not return an output, leave this one as is
	packedOutput := []byte{}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

func currentRewardAddress(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = deductGas(suppliedGas, CurrentRewardAddressGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// no input provided for this function

	// Allow list is enabled and CurrentRewardAddress is a state-changer function.
	// This part of the code restricts the function to be called only by enabled/admin addresses in the allow list.
	// You can modify/delete this code if you don't want this function to be restricted by the allow list.
	stateDB := accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to modify it
	callerStatus := getAllowListStatus(stateDB, RewardManagerAddress, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("%w: %s", ErrCannotCurrentRewardAddress, caller)
	}
	// allow list code ends here.

	// CUSTOM CODE STARTS HERE
	output := GetStoredRewardAddress(stateDB)
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

func disableRewards(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = deductGas(suppliedGas, DisableRewardsGasCost); err != nil {
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
	callerStatus := getAllowListStatus(stateDB, RewardManagerAddress, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("%w: %s", ErrCannotDisableRewards, caller)
	}
	// allow list code ends here.

	// CUSTOM CODE STARTS HERE
	// reset stored allow fee recipients flag only if it's already set to true
	if GetStoredAllowFeeRecipients(stateDB) {
		StoreAllowFeeRecipients(stateDB, false)
	}

	// reset stored reward address only if it's already set to non empty address
	if GetStoredRewardAddress(stateDB) != (common.Address{}) {
		StoreRewardAddress(stateDB, common.Address{})
	}
	// this function does not return an output, leave this one as is
	packedOutput := []byte{}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// createRewardManagerPrecompile returns a StatefulPrecompiledContract with getters and setters for the precompile.
// Access to the getters/setters is controlled by an allow list for [precompileAddr].
func createRewardManagerPrecompile(precompileAddr common.Address) StatefulPrecompiledContract {
	var functions []*statefulPrecompileFunction
	functions = append(functions, createAllowListFunctions(precompileAddr)...)

	methodAllowFeeRecipients, ok := RewardManagerABI.Methods["allowFeeRecipients"]
	if !ok {
		panic("given method does not exist in the ABI")
	}
	functions = append(functions, newStatefulPrecompileFunction(methodAllowFeeRecipients.ID, allowFeeRecipients))

	methodAreFeeRecipientsAllowed, ok := RewardManagerABI.Methods["areFeeRecipientsAllowed"]
	if !ok {
		panic("given method does not exist in the ABI")
	}
	functions = append(functions, newStatefulPrecompileFunction(methodAreFeeRecipientsAllowed.ID, areFeeRecipientsAllowed))

	methodCurrentRewardAddress, ok := RewardManagerABI.Methods["currentRewardAddress"]
	if !ok {
		panic("given method does not exist in the ABI")
	}
	functions = append(functions, newStatefulPrecompileFunction(methodCurrentRewardAddress.ID, currentRewardAddress))

	methodDisableRewards, ok := RewardManagerABI.Methods["disableRewards"]
	if !ok {
		panic("given method does not exist in the ABI")
	}
	functions = append(functions, newStatefulPrecompileFunction(methodDisableRewards.ID, disableRewards))

	methodSetRewardAddress, ok := RewardManagerABI.Methods["setRewardAddress"]
	if !ok {
		panic("given method does not exist in the ABI")
	}
	functions = append(functions, newStatefulPrecompileFunction(methodSetRewardAddress.ID, setRewardAddress))

	// Construct the contract with no fallback function.
	contract := newStatefulPrecompileWithFunctionSelectors(nil, functions)
	return contract
}
