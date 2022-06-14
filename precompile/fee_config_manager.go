// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// must preserve order of these fields
	gasLimitKey = iota
	targetBlockRateKey
	minBaseFeeKey
	targetGasKey
	baseFeeChangeDenominatorKey
	minBlockGasCostKey
	maxBlockGasCostKey
	blockGasCostStepKey

	minFeeConfigFieldKey = gasLimitKey
	maxFeeConfigFieldKey = blockGasCostStepKey

	numFeeConfigField = (maxFeeConfigFieldKey - minFeeConfigFieldKey + 1)
)

var (
	_ StatefulPrecompileConfig = &FeeConfigManagerConfig{}

	// Singleton StatefulPrecompiledContract for setting fee configs by permissioned callers.
	FeeConfigManagerPrecompile StatefulPrecompiledContract = createFeeConfigManagerPrecompile(FeeConfigManagerAddress)

	setFeeConfigSignature = CalculateFunctionSelector("setFeeConfig(uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)")
	getFeeConfigSignature = CalculateFunctionSelector("getFeeConfig()")

	ErrCannotChangeFee = errors.New("non-enabled cannot change fee config")

	// 8 fields in FeeConfig struct
	feeConfigInputLen = common.HashLength * numFeeConfigField
)

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

// Contract returns the singleton stateful precompiled contract to be used for the fee manager.
func (c *FeeConfigManagerConfig) Contract() StatefulPrecompiledContract {
	return FeeConfigManagerPrecompile
}

// GetFeeConfigManagerStatus returns the role of [address] for the fee config manager list.
func GetFeeConfigManagerStatus(stateDB StateDB, address common.Address) AllowListRole {
	return getAllowListStatus(stateDB, FeeConfigManagerAddress, address)
}

// SetFeeConfigManagerStatus sets the permissions of [address] to [role] for the
// fee config manager list. assumes [role] has already been verified as valid.
func SetFeeConfigManagerStatus(stateDB StateDB, address common.Address, role AllowListRole) {
	setAllowListRole(stateDB, FeeConfigManagerAddress, address, role)
}

// PackGetFeeConfigInput packs the getFeeConfig signature
func PackGetFeeConfigInput() []byte {
	return getFeeConfigSignature
}

// PackFeeConfig packs [feeConfig] without the selector into the appropriate arguments for fee config operations.
func PackFeeConfig(feeConfig commontype.FeeConfig) ([]byte, error) {
	//  input(feeConfig)
	return packHelper(feeConfig, false)
}

// PackSetFeeConfig packs [feeConfig] with the selector into the appropriate arguments for setting fee config operations.
func PackSetFeeConfig(feeConfig commontype.FeeConfig) ([]byte, error) {
	// function selector (4 bytes) + input(feeConfig)
	return packHelper(feeConfig, true)
}

func packHelper(feeConfig commontype.FeeConfig, useSelector bool) ([]byte, error) {
	fullLen := feeConfigInputLen
	packed := [][]byte{
		feeConfig.GasLimit.FillBytes(make([]byte, 32)),
		new(big.Int).SetUint64(feeConfig.TargetBlockRate).FillBytes(make([]byte, 32)),
		feeConfig.MinBaseFee.FillBytes(make([]byte, 32)),
		feeConfig.TargetGas.FillBytes(make([]byte, 32)),
		feeConfig.BaseFeeChangeDenominator.FillBytes(make([]byte, 32)),
		feeConfig.MinBlockGasCost.FillBytes(make([]byte, 32)),
		feeConfig.MaxBlockGasCost.FillBytes(make([]byte, 32)),
		feeConfig.BlockGasCostStep.FillBytes(make([]byte, 32)),
	}
	if useSelector {
		packed = append([][]byte{setFeeConfigSignature}, packed...)
		return packOrderedHashesWithSelector(packed, fullLen+selectorLen)
	}

	return packOrderedHashes(packed, fullLen)
}

// UnpackFeeConfigInput attempts to unpack [input] into the arguments to the fee config precompile
// assumes that [input] does not include selector (omits first 4 bytes in PackSetFeeConfigInput)
func UnpackFeeConfigInput(input []byte) (commontype.FeeConfig, error) {
	if len(input) != feeConfigInputLen {
		return commontype.FeeConfig{}, fmt.Errorf("invalid input length for fee config input: %d", len(input))
	}
	return commontype.FeeConfig{
		GasLimit:                 new(big.Int).SetBytes(returnPackedElement(input, 0)),
		TargetBlockRate:          new(big.Int).SetBytes(returnPackedElement(input, 1)).Uint64(),
		MinBaseFee:               new(big.Int).SetBytes(returnPackedElement(input, 2)),
		TargetGas:                new(big.Int).SetBytes(returnPackedElement(input, 3)),
		BaseFeeChangeDenominator: new(big.Int).SetBytes(returnPackedElement(input, 4)),
		MinBlockGasCost:          new(big.Int).SetBytes(returnPackedElement(input, 5)),
		MaxBlockGasCost:          new(big.Int).SetBytes(returnPackedElement(input, 6)),
		BlockGasCostStep:         new(big.Int).SetBytes(returnPackedElement(input, 7)),
	}, nil
}

// GetStoredFeeConfig returns fee config from contract storage in given state
func GetStoredFeeConfig(stateDB StateDB) (commontype.FeeConfig, bool) {
	if !stateDB.Exist(FeeConfigManagerAddress) {
		return commontype.FeeConfig{}, false
	}
	feeConfig := commontype.FeeConfig{}
	isAllZeros := true
	for i := minFeeConfigFieldKey; i <= maxFeeConfigFieldKey; i++ {
		val := stateDB.GetState(FeeConfigManagerAddress, common.Hash{byte(i)})
		// even if one of these vals are non-zero means that this state exist
		if isAllZeros && val != (common.Hash{0}) {
			isAllZeros = false
		}
		switch i {
		case gasLimitKey:
			feeConfig.GasLimit = new(big.Int).Set(val.Big())
		case targetBlockRateKey:
			feeConfig.TargetBlockRate = val.Big().Uint64()
		case minBaseFeeKey:
			feeConfig.MinBaseFee = new(big.Int).Set(val.Big())
		case targetGasKey:
			feeConfig.TargetGas = new(big.Int).Set(val.Big())
		case baseFeeChangeDenominatorKey:
			feeConfig.BaseFeeChangeDenominator = new(big.Int).Set(val.Big())
		case minBlockGasCostKey:
			feeConfig.MinBlockGasCost = new(big.Int).Set(val.Big())
		case maxBlockGasCostKey:
			feeConfig.MaxBlockGasCost = new(big.Int).Set(val.Big())
		case blockGasCostStepKey:
			feeConfig.BlockGasCostStep = new(big.Int).Set(val.Big())
		default:
			panic("unknown key")
		}
	}
	if isAllZeros {
		return commontype.EmptyFeeConfig, false
	}
	return feeConfig, true
}

func StoreFeeConfig(stateDB StateDB, feeConfig commontype.FeeConfig) {
	for i := minFeeConfigFieldKey; i <= maxFeeConfigFieldKey; i++ {
		var hashInput common.Hash
		switch i {
		case gasLimitKey:
			hashInput = common.BigToHash(feeConfig.GasLimit)
		case targetBlockRateKey:
			hashInput = common.BigToHash(new(big.Int).SetUint64(feeConfig.TargetBlockRate))
		case minBaseFeeKey:
			hashInput = common.BigToHash(feeConfig.MinBaseFee)
		case targetGasKey:
			hashInput = common.BigToHash(feeConfig.TargetGas)
		case baseFeeChangeDenominatorKey:
			hashInput = common.BigToHash(feeConfig.BaseFeeChangeDenominator)
		case minBlockGasCostKey:
			hashInput = common.BigToHash(feeConfig.MinBlockGasCost)
		case maxBlockGasCostKey:
			hashInput = common.BigToHash(feeConfig.MaxBlockGasCost)
		case blockGasCostStepKey:
			hashInput = common.BigToHash(feeConfig.BlockGasCostStep)
		default:
			panic("unknown key")
		}
		stateDB.SetState(FeeConfigManagerAddress, common.Hash{byte(i)}, hashInput)
	}
}

// setFeeConfig checks if the caller is permissioned for setting fee config operation.
// The execution function parses the [input] into FeeConfig structure and sets contract storage accordingly.
func setFeeConfig(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
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

	StoreFeeConfig(stateDB, feeConfig)

	// Return an empty output and the remaining gas
	return []byte{}, remainingGas, nil
}

// getFeeConfig returns the stored fee config as an output
func getFeeConfig(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = deductGas(suppliedGas, GetFeeConfigGasCost); err != nil {
		return nil, 0, err
	}

	feeConfig, ok := GetStoredFeeConfig(accessibleState.GetStateDB())
	if !ok {
		// we cannot pack empty fee config, instead use zero fee config
		output, err := PackFeeConfig(commontype.ZeroFeeConfig)
		// should never fail to packing an empty fee config
		if err != nil {
			panic(err)
		}
		return output, remainingGas, nil
	}

	output, err := PackFeeConfig(feeConfig)

	// Return an empty output and the remaining gas
	return output, remainingGas, err
}

// createFeeConfigManagerPrecompile returns a StatefulPrecompiledContract with R/W control of an allow list at [precompileAddr],
// and a storage setter for fee configs.
func createFeeConfigManagerPrecompile(precompileAddr common.Address) StatefulPrecompiledContract {
	setAdmin := newStatefulPrecompileFunction(setAdminSignature, createAllowListRoleSetter(precompileAddr, AllowListAdmin))
	setEnabled := newStatefulPrecompileFunction(setEnabledSignature, createAllowListRoleSetter(precompileAddr, AllowListEnabled))
	setNone := newStatefulPrecompileFunction(setNoneSignature, createAllowListRoleSetter(precompileAddr, AllowListNoRole))
	read := newStatefulPrecompileFunction(readAllowListSignature, createReadAllowList(precompileAddr))

	setFeeConfig := newStatefulPrecompileFunction(setFeeConfigSignature, setFeeConfig)
	getFeeConfig := newStatefulPrecompileFunction(getFeeConfigSignature, getFeeConfig)

	// Construct the contract with no fallback function.
	contract := newStatefulPrecompileWithFunctionSelectors(nil, []*statefulPrecompileFunction{setAdmin, setEnabled, setNone, read, setFeeConfig, getFeeConfig})
	return contract
}
