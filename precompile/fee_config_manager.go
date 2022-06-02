// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
)

// TODO: edit comments
var (
	_ StatefulPrecompileConfig = &FeeConfigManagerConfig{}

	// Singleton StatefulPrecompiledContract for minting native assets by permissioned callers.
	FeeConfigManagerPrecompile StatefulPrecompiledContract = createFeeConfigManagerPrecompile(ContractNativeMinterAddress)

	setFeeConfigSignature = CalculateFunctionSelector("setFeeConfig(uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)")
	// getFeeConfigSignature = CalculateFunctionSelector("getFeeConfig()")

	storageKey = common.Hash{}

	ErrCannotChangeFee = errors.New("non-enabled cannot change fee config")

	// 8 fields in FeeConfig struct
	feeConfigInputLen = common.HashLength * 8
)

// TODO: find a common place with this and params.FeeConfig
type FeeConfig struct {
	GasLimit *big.Int `json:"gasLimit,omitempty"`
	// TODO: make this uint64?
	TargetBlockRate *big.Int `json:"targetBlockRate,omitempty"`

	MinBaseFee               *big.Int `json:"minBaseFee,omitempty"`
	TargetGas                *big.Int `json:"targetGas,omitempty"`
	BaseFeeChangeDenominator *big.Int `json:"baseFeeChangeDenominator,omitempty"`

	MinBlockGasCost  *big.Int `json:"minBlockGasCost,omitempty"`
	MaxBlockGasCost  *big.Int `json:"maxBlockGasCost,omitempty"`
	BlockGasCostStep *big.Int `json:"blockGasCostStep,omitempty"`
}

// FeeConfigManagerConfig wraps [AllowListConfig] and uses it to implement the StatefulPrecompileConfig
// interface while adding in the contract deployer specific precompile address.
type FeeConfigManagerConfig struct {
	AllowListConfig
}

// Address returns the address of the fee config manager contract.
func (c *FeeConfigManagerConfig) Address() common.Address {
	return FeeConfigManagerAddress
}

// Configure configures [state] with the desired admins based on [c].
func (c *FeeConfigManagerConfig) Configure(state StateDB) {
	c.AllowListConfig.Configure(state, FeeConfigManagerAddress)
}

// Contract returns the singleton stateful precompiled contract to be used for the native minter.
func (c *FeeConfigManagerConfig) Contract() StatefulPrecompiledContract {
	return FeeConfigManagerPrecompile
}

// GetContractNativeMinterStatus returns the role of [address] for the minter list.
func GetFeeConfigManagerStatus(stateDB StateDB, address common.Address) AllowListRole {
	return getAllowListStatus(stateDB, FeeConfigManagerAddress, address)
}

// SetContractNativeMinterStatus sets the permissions of [address] to [role] for the
// minter list. assumes [role] has already been verified as valid.
func SetFeeConfigManagerStatus(stateDB StateDB, address common.Address, role AllowListRole) {
	setAllowListRole(stateDB, FeeConfigManagerAddress, address, role)
}

// PackMintInput packs [address] and [amount] into the appropriate arguments for minting operation.
func PackFeeConfigInput(feeConfig FeeConfig) ([]byte, error) {
	// function selector (4 bytes) + input(hash for address + hash for amount)
	fullLen := selectorLen + feeConfigInputLen
	packed := [][]byte{
		setFeeConfigSignature,
		feeConfig.GasLimit.FillBytes(make([]byte, 32)),
		feeConfig.TargetBlockRate.FillBytes(make([]byte, 32)),
		feeConfig.MinBaseFee.FillBytes(make([]byte, 32)),
		feeConfig.TargetGas.FillBytes(make([]byte, 32)),
		feeConfig.BaseFeeChangeDenominator.FillBytes(make([]byte, 32)),
		feeConfig.MinBlockGasCost.FillBytes(make([]byte, 32)),
		feeConfig.MaxBlockGasCost.FillBytes(make([]byte, 32)),
		feeConfig.BlockGasCostStep.FillBytes(make([]byte, 32)),
	}
	return inputPackOrdered(packed, fullLen)
}

// UnpackMintInput attempts to unpack [input] into the arguments to the mint precompile
// assumes that [input] does not include selector (omits first 4 bytes in PackMintInput)
func UnpackFeeConfigInput(input []byte) (FeeConfig, error) {
	if len(input) != feeConfigInputLen {
		return FeeConfig{}, fmt.Errorf("invalid input length for fee config input: %d", len(input))
	}
	return FeeConfig{
		GasLimit:                 new(big.Int).SetBytes(returnPackedElement(input, 0)),
		TargetBlockRate:          new(big.Int).SetBytes(returnPackedElement(input, 1)),
		MinBaseFee:               new(big.Int).SetBytes(returnPackedElement(input, 2)),
		TargetGas:                new(big.Int).SetBytes(returnPackedElement(input, 3)),
		BaseFeeChangeDenominator: new(big.Int).SetBytes(returnPackedElement(input, 4)),
		MinBlockGasCost:          new(big.Int).SetBytes(returnPackedElement(input, 5)),
		MaxBlockGasCost:          new(big.Int).SetBytes(returnPackedElement(input, 6)),
		BlockGasCostStep:         new(big.Int).SetBytes(returnPackedElement(input, 7)),
	}, nil
}

func getFeeConfig(state StateDB) (FeeConfig, error) {
	// Generate the state key for [address]
	feeHash := state.GetState(FeeConfigManagerAddress, storageKey)
	if (feeHash == common.Hash{}) {
		return FeeConfig{}, nil
	}
	v := FeeConfig{}
	err := json.Unmarshal(feeHash.Bytes(), &v)
	return v, err
}

func setFeeConfig(stateDB StateDB, feeConfig FeeConfig) error {
	// TODO: use a better serilization?
	feeBytes, err := json.Marshal(feeConfig)
	if err != nil {
		return err
	}
	// Assign [role] to the address
	stateDB.SetState(FeeConfigManagerAddress, storageKey, common.BytesToHash(feeBytes))
	return nil
}

// createMintNativeCoin checks if the caller is permissioned for minting operation.
// The execution function parses the [input] into native coin amount and receiver address.
func createSetFeeConfig(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = deductGas(suppliedGas, SetFeeConfigGasCost); err != nil {
		return nil, 0, err
	}

	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}

	feeConfig, err := UnpackFeeConfigInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	stateDB := accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to modify it
	callerStatus := getAllowListStatus(stateDB, FeeConfigManagerAddress, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("%w: %s", ErrCannotChangeFee, caller)
	}

	setFeeConfig(accessibleState.GetStateDB(), feeConfig)

	// Return an empty output and the remaining gas
	return []byte{}, remainingGas, nil
}

// createNativeMinterPrecompile returns a StatefulPrecompiledContract with R/W control of an allow list at [precompileAddr] and a native coin minter.
func createFeeConfigManagerPrecompile(precompileAddr common.Address) StatefulPrecompiledContract {
	setAdmin := newStatefulPrecompileFunction(setAdminSignature, createAllowListRoleSetter(precompileAddr, AllowListAdmin))
	setEnabled := newStatefulPrecompileFunction(setEnabledSignature, createAllowListRoleSetter(precompileAddr, AllowListEnabled))
	setNone := newStatefulPrecompileFunction(setNoneSignature, createAllowListRoleSetter(precompileAddr, AllowListNoRole))
	read := newStatefulPrecompileFunction(readAllowListSignature, createReadAllowList(precompileAddr))

	setFeeConfig := newStatefulPrecompileFunction(setFeeConfigSignature, createSetFeeConfig)

	// Construct the contract with no fallback function.
	contract := newStatefulPrecompileWithFunctionSelectors(nil, []*statefulPrecompileFunction{setAdmin, setEnabled, setNone, read, setFeeConfig})
	return contract
}
