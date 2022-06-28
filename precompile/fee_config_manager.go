// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ethereum/go-ethereum/common"
)

const (
	minFeeConfigFieldKey = iota + 1
	// add new fields below this
	// must preserve order of these fields
	gasLimitKey = iota
	targetBlockRateKey
	minBaseFeeKey
	targetGasKey
	baseFeeChangeDenominatorKey
	minBlockGasCostKey
	maxBlockGasCostKey
	blockGasCostStepKey
	// add new fields above this
	numFeeConfigField = iota - 1
)

var (
	_ StatefulPrecompileConfig = &FeeConfigManagerConfig{}

	// Singleton StatefulPrecompiledContract for setting fee configs by permissioned callers.
	FeeConfigManagerPrecompile StatefulPrecompiledContract = createFeeConfigManagerPrecompile(FeeConfigManagerAddress)

	setFeeConfigSignature     = CalculateFunctionSelector("setFeeConfig(uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)")
	getFeeConfigSignature     = CalculateFunctionSelector("getFeeConfig()")
	getLastChangedAtSignature = CalculateFunctionSelector("getLastChangedAt()")

	// 8 fields in FeeConfig struct
	feeConfigInputLen = common.HashLength * numFeeConfigField

	lastChangedAtKey = common.Hash{'l', 'c', 'a'}
)

// FeeConfigManagerConfig wraps [AllowListConfig] and uses it to implement the StatefulPrecompileConfig
// interface while adding in the contract deployer specific precompile address.
type FeeConfigManagerConfig struct {
	AllowListConfig
	commontype.FeeConfig
}

// Address returns the address of the fee config manager contract.
func (c *FeeConfigManagerConfig) Address() common.Address {
	return FeeConfigManagerAddress
}

// Configure configures [state] with the desired admins based on [c].
func (c *FeeConfigManagerConfig) Configure(state StateDB, blockContext BlockContext) {
	if err := StoreFeeConfig(state, c.FeeConfig, blockContext); err != nil {
		panic(err) // this should be already verified in genesis
	}
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

// PackGetLastChangedAtInput packs the getLastChangedAt signature
func PackGetLastChangedAtInput() []byte {
	return getLastChangedAtSignature
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
		feeConfig.GasLimit.FillBytes(make([]byte, common.HashLength)),
		new(big.Int).SetUint64(feeConfig.TargetBlockRate).FillBytes(make([]byte, common.HashLength)),
		feeConfig.MinBaseFee.FillBytes(make([]byte, common.HashLength)),
		feeConfig.TargetGas.FillBytes(make([]byte, common.HashLength)),
		feeConfig.BaseFeeChangeDenominator.FillBytes(make([]byte, common.HashLength)),
		feeConfig.MinBlockGasCost.FillBytes(make([]byte, common.HashLength)),
		feeConfig.MaxBlockGasCost.FillBytes(make([]byte, common.HashLength)),
		feeConfig.BlockGasCostStep.FillBytes(make([]byte, common.HashLength)),
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
	feeConfig := commontype.FeeConfig{}
	for i := minFeeConfigFieldKey; i <= numFeeConfigField; i++ {
		listIndex := i - 1
		packedElement := returnPackedElement(input, listIndex)
		switch i {
		case gasLimitKey:
			feeConfig.GasLimit = new(big.Int).SetBytes(packedElement)
		case targetBlockRateKey:
			feeConfig.TargetBlockRate = new(big.Int).SetBytes(packedElement).Uint64()
		case minBaseFeeKey:
			feeConfig.MinBaseFee = new(big.Int).SetBytes(packedElement)
		case targetGasKey:
			feeConfig.TargetGas = new(big.Int).SetBytes(packedElement)
		case baseFeeChangeDenominatorKey:
			feeConfig.BaseFeeChangeDenominator = new(big.Int).SetBytes(packedElement)
		case minBlockGasCostKey:
			feeConfig.MinBlockGasCost = new(big.Int).SetBytes(packedElement)
		case maxBlockGasCostKey:
			feeConfig.MaxBlockGasCost = new(big.Int).SetBytes(packedElement)
		case blockGasCostStepKey:
			feeConfig.BlockGasCostStep = new(big.Int).SetBytes(packedElement)
		default:
			panic("unknown key")
		}
	}
	return feeConfig, nil
}

// GetStoredFeeConfig returns fee config from contract storage in given state
func GetStoredFeeConfig(stateDB StateDB) commontype.FeeConfig {
	feeConfig := commontype.FeeConfig{}
	for i := minFeeConfigFieldKey; i <= numFeeConfigField; i++ {
		val := stateDB.GetState(FeeConfigManagerAddress, common.Hash{byte(i)})
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
	return feeConfig
}

func GetStoredLastChangedAt(stateDB StateDB) *big.Int {
	val := stateDB.GetState(FeeConfigManagerAddress, lastChangedAtKey)
	return val.Big()
}

// StoreFeeConfig stores given [feeConfig] and block number in the [blockContext] to the [stateDB].
// A validation on [feeConfig] is done before storing.
func StoreFeeConfig(stateDB StateDB, feeConfig commontype.FeeConfig, blockContext BlockContext) error {
	if err := feeConfig.Verify(); err != nil {
		return err
	}

	blockNumber := blockContext.Number()
	if blockNumber == nil {
		return fmt.Errorf("blockNumber cannot be nil")
	}

	hashes, err := getFeeConfigHashes(feeConfig)
	if err != nil {
		return err
	}
	for i := minFeeConfigFieldKey; i <= numFeeConfigField; i++ {
		stateDB.SetState(FeeConfigManagerAddress, common.Hash{byte(i)}, hashes[i])
	}

	stateDB.SetState(FeeConfigManagerAddress, lastChangedAtKey, common.BigToHash(blockNumber))
	return nil
}

// getFeeConfigHashes takes [feeConfig] and converts them to an array of hashes, with ordered with key indexes.
func getFeeConfigHashes(feeConfig commontype.FeeConfig) ([]common.Hash, error) {
	res := make([]common.Hash, minFeeConfigFieldKey+numFeeConfigField)
	for i := minFeeConfigFieldKey; i <= numFeeConfigField; i++ {
		var hashInput common.Hash
		var err error
		switch i {
		case gasLimitKey:
			hashInput, err = bigToHashSafe(feeConfig.GasLimit)
		case targetBlockRateKey:
			hashInput, err = bigToHashSafe(new(big.Int).SetUint64(feeConfig.TargetBlockRate))
		case minBaseFeeKey:
			hashInput, err = bigToHashSafe(feeConfig.MinBaseFee)
		case targetGasKey:
			hashInput, err = bigToHashSafe(feeConfig.TargetGas)
		case baseFeeChangeDenominatorKey:
			hashInput, err = bigToHashSafe(feeConfig.BaseFeeChangeDenominator)
		case minBlockGasCostKey:
			hashInput, err = bigToHashSafe(feeConfig.MinBlockGasCost)
		case maxBlockGasCostKey:
			hashInput, err = bigToHashSafe(feeConfig.MaxBlockGasCost)
		case blockGasCostStepKey:
			hashInput, err = bigToHashSafe(feeConfig.BlockGasCostStep)
		default:
			panic("unknown key")
		}
		if err != nil {
			return nil, err
		}
		// omits first slot in order to normalize indexes with keys
		res[i] = hashInput
	}
	return res, nil
}

// setFeeConfig checks if the caller is permissioned for setting fee config operation.
// The execution function parses the [input] into FeeConfig structure and sets contract storage accordingly.
func setFeeConfig(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	stateDB := accessibleState.GetStateDB()

	if remainingGas, err = allowListSetterEnabledCheck(FeeConfigManagerAddress, caller, suppliedGas, SetFeeConfigGasCost, readOnly, stateDB); err != nil {
		return nil, remainingGas, err
	}

	feeConfig, err := UnpackFeeConfigInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	if err := StoreFeeConfig(stateDB, feeConfig, accessibleState.GetBlockContext()); err != nil {
		return nil, remainingGas, err
	}

	// Return an empty output and the remaining gas
	return []byte{}, remainingGas, nil
}

// getFeeConfig returns the stored fee config as an output.
// The execution function reads the contract state for the stored fee config and returns the output accordingly.
func getFeeConfig(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = deductGas(suppliedGas, GetFeeConfigGasCost); err != nil {
		return nil, 0, err
	}

	feeConfig := GetStoredFeeConfig(accessibleState.GetStateDB())

	output, err := PackFeeConfig(feeConfig)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return an empty output and the remaining gas
	return output, remainingGas, err
}

// getLastChangedAt returns the block number that fee config was last changed in.
// The execution function reads the contract state for the stored block number and returns the output accordingly.
func getLastChangedAt(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = deductGas(suppliedGas, GetLastChangedAtGasCost); err != nil {
		return nil, 0, err
	}

	lastChangedAt := GetStoredLastChangedAt(accessibleState.GetStateDB())

	// Return an empty output and the remaining gas
	return common.BigToHash(lastChangedAt).Bytes(), remainingGas, err
}

// createFeeConfigManagerPrecompile returns a StatefulPrecompiledContract with R/W control of an allow list at [precompileAddr],
// and a storage setter for fee configs.
func createFeeConfigManagerPrecompile(precompileAddr common.Address) StatefulPrecompiledContract {
	allowListFuncs := createAllowListFunctions(precompileAddr)

	setFeeConfigFunc := newStatefulPrecompileFunction(setFeeConfigSignature, setFeeConfig)
	getFeeConfigFunc := newStatefulPrecompileFunction(getFeeConfigSignature, getFeeConfig)
	getLastChangedAtFunc := newStatefulPrecompileFunction(getLastChangedAtSignature, getLastChangedAt)

	enabledFuncs := append(allowListFuncs, setFeeConfigFunc, getFeeConfigFunc, getLastChangedAtFunc)
	// Construct the contract with no fallback function.
	contract := newStatefulPrecompileWithFunctionSelectors(nil, enabledFuncs)
	return contract
}
