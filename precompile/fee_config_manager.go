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

	setFeeConfigSignature              = CalculateFunctionSelector("setFeeConfig(uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)")
	getFeeConfigSignature              = CalculateFunctionSelector("getFeeConfig()")
	getFeeConfigLastChangedAtSignature = CalculateFunctionSelector("getFeeConfigLastChangedAt()")

	// 8 fields in FeeConfig struct
	feeConfigInputLen = common.HashLength * numFeeConfigField

	feeConfigLastChangedAtKey = common.Hash{'l', 'c', 'a'}
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
		panic(fmt.Sprintf("fee config should have been verified in genesis: %s", err))
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

// PackGetLastChangedAtInput packs the getFeeConfigLastChangedAt signature
func PackGetLastChangedAtInput() []byte {
	return getFeeConfigLastChangedAtSignature
}

// PackFeeConfig packs [feeConfig] without the selector into the appropriate arguments for fee config operations.
func PackFeeConfig(feeConfig commontype.FeeConfig) ([]byte, error) {
	//  input(feeConfig)
	return packFeeConfigHelper(feeConfig, false), nil
}

// PackSetFeeConfig packs [feeConfig] with the selector into the appropriate arguments for setting fee config operations.
func PackSetFeeConfig(feeConfig commontype.FeeConfig) ([]byte, error) {
	// function selector (4 bytes) + input(feeConfig)
	return packFeeConfigHelper(feeConfig, true), nil
}

func packFeeConfigHelper(feeConfig commontype.FeeConfig, useSelector bool) []byte {
	hashes := []common.Hash{
		common.BigToHash(feeConfig.GasLimit),
		common.BigToHash(new(big.Int).SetUint64(feeConfig.TargetBlockRate)),
		common.BigToHash(feeConfig.MinBaseFee),
		common.BigToHash(feeConfig.TargetGas),
		common.BigToHash(feeConfig.BaseFeeChangeDenominator),
		common.BigToHash(feeConfig.MinBlockGasCost),
		common.BigToHash(feeConfig.MaxBlockGasCost),
		common.BigToHash(feeConfig.BlockGasCostStep),
	}

	if useSelector {
		res := make([]byte, len(setFeeConfigSignature)+len(hashes)*common.HashLength)
		packOrderedHashesWithSelector(res, setFeeConfigSignature, hashes)
		return res
	}

	res := make([]byte, len(hashes)*common.HashLength)
	packOrderedHashes(res, hashes)
	return res
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
		packedElement := returnPackedHash(input, listIndex)
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
			panic(fmt.Sprintf("unknown fee config key: %d", i))
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
			panic(fmt.Sprintf("unknown fee config key: %d", i))
		}
	}
	return feeConfig
}

func GetFeeConfigLastUpdatedAt(stateDB StateDB) *big.Int {
	val := stateDB.GetState(FeeConfigManagerAddress, feeConfigLastChangedAtKey)
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
	stateDB.SetState(FeeConfigManagerAddress, feeConfigLastChangedAtKey, common.BigToHash(blockNumber))

	for i := minFeeConfigFieldKey; i <= numFeeConfigField; i++ {
		var input *big.Int
		switch i {
		case gasLimitKey:
			input = feeConfig.GasLimit
		case targetBlockRateKey:
			input = new(big.Int).SetUint64(feeConfig.TargetBlockRate)
		case minBaseFeeKey:
			input = feeConfig.MinBaseFee
		case targetGasKey:
			input = feeConfig.TargetGas
		case baseFeeChangeDenominatorKey:
			input = feeConfig.BaseFeeChangeDenominator
		case minBlockGasCostKey:
			input = feeConfig.MinBlockGasCost
		case maxBlockGasCostKey:
			input = feeConfig.MaxBlockGasCost
		case blockGasCostStepKey:
			input = feeConfig.BlockGasCostStep
		default:
			panic(fmt.Sprintf("unknown fee config key: %d", i))
		}
		stateDB.SetState(FeeConfigManagerAddress, common.Hash{byte(i)}, common.BigToHash(input))
	}

	return nil
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

// getFeeConfigLastChangedAt returns the block number that fee config was last changed in.
// The execution function reads the contract state for the stored block number and returns the output accordingly.
func getFeeConfigLastChangedAt(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = deductGas(suppliedGas, GetLastChangedAtGasCost); err != nil {
		return nil, 0, err
	}

	lastChangedAt := GetFeeConfigLastUpdatedAt(accessibleState.GetStateDB())

	// Return an empty output and the remaining gas
	return common.BigToHash(lastChangedAt).Bytes(), remainingGas, err
}

// createFeeConfigManagerPrecompile returns a StatefulPrecompiledContract with R/W control of an allow list at [precompileAddr],
// and a storage setter for fee configs.
func createFeeConfigManagerPrecompile(precompileAddr common.Address) StatefulPrecompiledContract {
	allowListFuncs := createAllowListFunctions(precompileAddr)

	setFeeConfigFunc := newStatefulPrecompileFunction(setFeeConfigSignature, setFeeConfig)
	getFeeConfigFunc := newStatefulPrecompileFunction(getFeeConfigSignature, getFeeConfig)
	getFeeConfigLastChangedAtFunc := newStatefulPrecompileFunction(getFeeConfigLastChangedAtSignature, getFeeConfigLastChangedAt)

	// Construct the contract with no fallback function.
	contract := newStatefulPrecompileWithFunctionSelectors(nil, append(allowListFuncs, setFeeConfigFunc, getFeeConfigFunc, getFeeConfigLastChangedAtFunc))
	return contract
}
