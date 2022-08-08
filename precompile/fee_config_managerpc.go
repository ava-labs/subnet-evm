// Code generated
// This file is a generated precompile with stubbed abstract functions.

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
	GetFeeConfigPCGasCost              uint64 = GetFeeConfigGasCost     // SET A GAS COST HERE
	GetFeeConfigPCLastChangedAtGasCost uint64 = GetLastChangedAtGasCost // SET A GAS COST HERE
	SetFeeConfigPCGasCost              uint64 = SetFeeConfigGasCost     // SET A GAS COST HERE
	GetTestPCGasCost                   uint64 = 1                       // SET A GAS COST HERE

	// FeeConfigManagerPCRawABI contains the raw ABI of FeeConfigManagerPC contract.
	FeeConfigManagerPCRawABI = "[{\"inputs\":[],\"name\":\"getFeeConfigPC\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"targetBlockRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"targetGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseFeeChangeDenominator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBlockGasCost\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxBlockGasCost\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockGasCostStep\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeeConfigPCLastChangedAt\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"testt\",\"type\":\"uint256\"}],\"name\":\"getTestPC\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"readAllowList\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setEnabled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"targetBlockRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBaseFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"targetGas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseFeeChangeDenominator\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBlockGasCost\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxBlockGasCost\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockGasCostStep\",\"type\":\"uint256\"}],\"name\":\"setFeeConfigPC\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setNone\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
)

// Singleton StatefulPrecompiledContract and signatures.
var (
	_ StatefulPrecompileConfig = &FeeConfigManagerPCConfig{}

	FeeConfigManagerPCPrecompile         StatefulPrecompiledContract = createFeeConfigManagerPCPrecompile(FeeConfigManagerPCAddress)
	getFeeConfigPCSignature                                          = CalculateFunctionSelector("getFeeConfigPC()")
	getFeeConfigPCLastChangedAtSignature                             = CalculateFunctionSelector("getFeeConfigPCLastChangedAt()")
	getTestPCSignature                                               = CalculateFunctionSelector("getTestPC(uint256)")
	setFeeConfigPCSignature                                          = CalculateFunctionSelector("setFeeConfigPC(uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)")

	ErrCannotSetFeeConfigPC = errors.New("non-enabled cannot setFeeConfigPC")

	FeeConfigManagerPCABI abi.ABI // will be filled by init func
)

// FeeConfigManagerPCConfig wraps [AllowListConfig] and uses it to implement  the StatefulPrecompileConfig
// interface while adding in the FeeConfigManagerPC specific precompile address.
type FeeConfigManagerPCConfig struct {
	AllowListConfig
	UpgradeableConfig
}

type GetFeeConfigPCOutput struct {
	GasLimit                 *big.Int
	TargetBlockRate          *big.Int
	MinBaseFee               *big.Int
	TargetGas                *big.Int
	BaseFeeChangeDenominator *big.Int
	MinBlockGasCost          *big.Int
	MaxBlockGasCost          *big.Int
	BlockGasCostStep         *big.Int
}

type SetFeeConfigPCInput struct {
	GasLimit                 *big.Int
	TargetBlockRate          *big.Int
	MinBaseFee               *big.Int
	TargetGas                *big.Int
	BaseFeeChangeDenominator *big.Int
	MinBlockGasCost          *big.Int
	MaxBlockGasCost          *big.Int
	BlockGasCostStep         *big.Int
}

func init() {
	parsed, err := abi.JSON(strings.NewReader(FeeConfigManagerPCRawABI))
	if err != nil {
		panic(err)
	}
	FeeConfigManagerPCABI = parsed
}

// NewFeeConfigManagerPCConfig returns a config for a network upgrade at [blockTimestamp] that enables
// FeeConfigManagerPC  with the given [admins] as members of the allowlist .
func NewFeeConfigManagerPCConfig(blockTimestamp *big.Int, admins []common.Address) *FeeConfigManagerPCConfig {
	return &FeeConfigManagerPCConfig{
		AllowListConfig:   AllowListConfig{AllowListAdmins: admins},
		UpgradeableConfig: UpgradeableConfig{BlockTimestamp: blockTimestamp},
	}
}

// NewDisableFeeConfigManagerPCConfig returns config for a network upgrade at [blockTimestamp]
// that disables FeeConfigManagerPC.
func NewDisableFeeConfigManagerPCConfig(blockTimestamp *big.Int) *FeeConfigManagerPCConfig {
	return &FeeConfigManagerPCConfig{
		UpgradeableConfig: UpgradeableConfig{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

// Equal returns true if [s] is a [*FeeConfigManagerPCConfig] and it has been configured identical to [c].
func (c *FeeConfigManagerPCConfig) Equal(s StatefulPrecompileConfig) bool {
	// typecast before comparison
	other, ok := (s).(*FeeConfigManagerPCConfig)
	if !ok {
		return false
	}
	return c.UpgradeableConfig.Equal(&other.UpgradeableConfig) && c.AllowListConfig.Equal(&other.AllowListConfig)
}

// Address returns the address of the FeeConfigManagerPC. Addresses reside under the precompile/params.go
// Select a non-conflicting address and set it in the params.go.
func (c *FeeConfigManagerPCConfig) Address() common.Address {
	return FeeConfigManagerPCAddress
}

// Configure configures [state] with the initial configuration.
func (c *FeeConfigManagerPCConfig) Configure(_ ChainConfig, state StateDB, _ BlockContext) {
	c.AllowListConfig.Configure(state, FeeConfigManagerPCAddress)
	// CUSTOM CODE STARTS HERE
}

// Contract returns the singleton stateful precompiled contract to be used for FeeConfigManagerPC.
func (c *FeeConfigManagerPCConfig) Contract() StatefulPrecompiledContract {
	return FeeConfigManagerPCPrecompile
}

// GetFeeConfigManagerPCStatus returns the role of [address] for the FeeConfigManagerPC list.
func GetFeeConfigManagerPCStatus(stateDB StateDB, address common.Address) AllowListRole {
	return getAllowListStatus(stateDB, FeeConfigManagerPCAddress, address)
}

// SetFeeConfigManagerPCStatus sets the permissions of [address] to [role] for the
// FeeConfigManagerPC list. Assumes [role] has already been verified as valid.
func SetFeeConfigManagerPCStatus(stateDB StateDB, address common.Address, role AllowListRole) {
	setAllowListRole(stateDB, FeeConfigManagerPCAddress, address, role)
}

// CUSTOM CODE ADDED HERE FOR TEST
// GetStoredFeeConfig returns fee config from contract storage in given state
func GetStoredFeeConfigPC(stateDB StateDB) GetFeeConfigPCOutput {
	feeConfig := GetFeeConfigPCOutput{}
	for i := minFeeConfigFieldKey; i <= numFeeConfigField; i++ {
		val := stateDB.GetState(FeeConfigManagerAddress, common.Hash{byte(i)})
		switch i {
		case gasLimitKey:
			feeConfig.GasLimit = new(big.Int).Set(val.Big())
		case targetBlockRateKey:
			feeConfig.TargetBlockRate = new(big.Int).Set(val.Big())
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

// CUSTOM CODE ADDED HERE FOR TEST
func StoreFeeConfigPC(stateDB StateDB, feeConfig SetFeeConfigPCInput, blockContext BlockContext) error {
	for i := minFeeConfigFieldKey; i <= numFeeConfigField; i++ {
		var input common.Hash
		switch i {
		case gasLimitKey:
			input = common.BigToHash(feeConfig.GasLimit)
		case targetBlockRateKey:
			input = common.BigToHash(feeConfig.TargetBlockRate)
		case minBaseFeeKey:
			input = common.BigToHash(feeConfig.MinBaseFee)
		case targetGasKey:
			input = common.BigToHash(feeConfig.TargetGas)
		case baseFeeChangeDenominatorKey:
			input = common.BigToHash(feeConfig.BaseFeeChangeDenominator)
		case minBlockGasCostKey:
			input = common.BigToHash(feeConfig.MinBlockGasCost)
		case maxBlockGasCostKey:
			input = common.BigToHash(feeConfig.MaxBlockGasCost)
		case blockGasCostStepKey:
			input = common.BigToHash(feeConfig.BlockGasCostStep)
		default:
			panic(fmt.Sprintf("unknown fee config key: %d", i))
		}
		stateDB.SetState(FeeConfigManagerAddress, common.Hash{byte(i)}, input)
	}

	blockNumber := blockContext.Number()
	if blockNumber == nil {
		return fmt.Errorf("blockNumber cannot be nil")
	}
	stateDB.SetState(FeeConfigManagerAddress, feeConfigLastChangedAtKey, common.BigToHash(blockNumber))

	return nil
}

// PackGetFeeConfigPC packs the include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackGetFeeConfigPC() ([]byte, error) {
	return FeeConfigManagerPCABI.Pack("getFeeConfigPC")
}

// PackGetFeeConfigPCOutput attempts to pack given [outputStruct] of type GetFeeConfigPCOutput
// to conform the ABI outputs.
func PackGetFeeConfigPCOutput(outputStruct GetFeeConfigPCOutput) ([]byte, error) {
	return FeeConfigManagerPCABI.PackOutput("getFeeConfigPC",
		outputStruct.GasLimit,
		outputStruct.TargetBlockRate,
		outputStruct.MinBaseFee,
		outputStruct.TargetGas,
		outputStruct.BaseFeeChangeDenominator,
		outputStruct.MinBlockGasCost,
		outputStruct.MaxBlockGasCost,
		outputStruct.BlockGasCostStep,
	)
}

func getFeeConfigPC(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = deductGas(suppliedGas, GetFeeConfigPCGasCost); err != nil {
		return nil, 0, err
	}
	// no input provided for this function

	// CUSTOM CODE STARTS HERE
	// use inputStruct ...
	feeConfig := GetStoredFeeConfigPC(accessibleState.GetStateDB())

	packedOutput, err := PackGetFeeConfigPCOutput(feeConfig)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return an empty output and the remaining gas
	return packedOutput, remainingGas, nil
}

// CUSTOM CODE ADDED FOR TEST
func GetFeeConfigLastChangedAtPC(stateDB StateDB) *big.Int {
	val := stateDB.GetState(FeeConfigManagerAddress, feeConfigLastChangedAtKey)
	return val.Big()
}

// PackGetFeeConfigPCLastChangedAt packs the include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackGetFeeConfigPCLastChangedAt() ([]byte, error) {
	return FeeConfigManagerPCABI.Pack("getFeeConfigPCLastChangedAt")
}

// PackGetFeeConfigPCLastChangedAtOutput attempts to pack given blockNumber of type *big.Int
// to conform the ABI outputs.
func PackGetFeeConfigPCLastChangedAtOutput(blockNumber *big.Int) ([]byte, error) {
	return FeeConfigManagerPCABI.PackOutput("getFeeConfigPCLastChangedAt", blockNumber)
}

func getFeeConfigPCLastChangedAt(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = deductGas(suppliedGas, GetFeeConfigPCLastChangedAtGasCost); err != nil {
		return nil, 0, err
	}
	// no input provided for this function

	// CUSTOM CODE STARTS HERE
	// use inputStruct ...
	outputBlockNumber := GetFeeConfigLastChangedAtPC(accessibleState.GetStateDB())
	packedOutput, err := PackGetFeeConfigPCLastChangedAtOutput(outputBlockNumber)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return an empty output and the remaining gas
	return packedOutput, remainingGas, nil
}

// UnpackGetTestPCInput attempts to unpack [input] into the *big.Int type argument
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackGetTestPCInput(input []byte) (*big.Int, error) {
	res, err := FeeConfigManagerPCABI.UnpackInput("getTestPC", input)
	if err != nil {
		return nil, err
	}
	unpacked := *abi.ConvertType(res[0], new(*big.Int)).(**big.Int)
	return unpacked, nil
}

// PackGetTestPC packs [testt] of type *big.Int into the appropriate arguments for getTestPC.
// the packed bytes include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackGetTestPC(testt *big.Int) ([]byte, error) {
	return FeeConfigManagerPCABI.Pack("getTestPC", testt)
}

// PackGetTestPCOutput attempts to pack given blockNumber of type *big.Int
// to conform the ABI outputs.
func PackGetTestPCOutput(blockNumber *big.Int) ([]byte, error) {
	return FeeConfigManagerPCABI.PackOutput("getTestPC", blockNumber)
}

func getTestPC(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = deductGas(suppliedGas, GetTestPCGasCost); err != nil {
		return nil, 0, err
	}
	// attempts to unpack [input] into the arguments to the GetTestPCInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackGetTestPCInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	// use inputStruct ...
	packedOutput, err := PackGetTestPCOutput(inputStruct)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return an empty output and the remaining gas
	return packedOutput, remainingGas, nil
}

// UnpackSetFeeConfigPCInput attempts to unpack [input] into the arguments for the SetFeeConfigPCInput{}
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackSetFeeConfigPCInput(input []byte) (SetFeeConfigPCInput, error) {
	inputStruct := SetFeeConfigPCInput{}
	err := FeeConfigManagerPCABI.UnpackInputIntoInterface(&inputStruct, "setFeeConfigPC", input)

	return inputStruct, err
}

// PackSetFeeConfigPC packs [inputStruct] of type SetFeeConfigPCInput into the appropriate arguments for setFeeConfigPC.
func PackSetFeeConfigPC(inputStruct SetFeeConfigPCInput) ([]byte, error) {
	return FeeConfigManagerPCABI.Pack("setFeeConfigPC", inputStruct.GasLimit, inputStruct.TargetBlockRate, inputStruct.MinBaseFee, inputStruct.TargetGas, inputStruct.BaseFeeChangeDenominator, inputStruct.MinBlockGasCost, inputStruct.MaxBlockGasCost, inputStruct.BlockGasCostStep)
}

func setFeeConfigPC(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = deductGas(suppliedGas, SetFeeConfigPCGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// attempts to unpack [input] into the arguments to the SetFeeConfigPCInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackSetFeeConfigPCInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	stateDB := accessibleState.GetStateDB()
	// Verify that the caller is in the allow list and therefore has the right to modify it
	callerStatus := getAllowListStatus(stateDB, FeeConfigManagerPCAddress, caller)
	if !callerStatus.IsEnabled() {
		return nil, remainingGas, fmt.Errorf("%w: %s", ErrCannotSetFeeConfigPC, caller)
	}

	// CUSTOM CODE STARTS HERE
	// use inputStruct ...
	if err := StoreFeeConfigPC(stateDB, inputStruct, accessibleState.GetBlockContext()); err != nil {
		return nil, remainingGas, err
	}

	// this function does not return an output, leave this one as is
	packedOutput := []byte{}

	// Return an empty output and the remaining gas
	return packedOutput, remainingGas, nil
}

// createFeeConfigManagerPCPrecompile returns a StatefulPrecompiledContract
// with getters and setters for the precompile.
// Access to the getters/setters is controlled by an allow list for [precompileAddr].
func createFeeConfigManagerPCPrecompile(precompileAddr common.Address) StatefulPrecompiledContract {
	var functions []*statefulPrecompileFunction
	functions = append(functions, createAllowListFunctions(precompileAddr)...)
	functions = append(functions, newStatefulPrecompileFunction(getFeeConfigPCSignature, getFeeConfigPC))
	functions = append(functions, newStatefulPrecompileFunction(getFeeConfigPCLastChangedAtSignature, getFeeConfigPCLastChangedAt))
	functions = append(functions, newStatefulPrecompileFunction(getTestPCSignature, getTestPC))
	functions = append(functions, newStatefulPrecompileFunction(setFeeConfigPCSignature, setFeeConfigPC))

	// Construct the contract with no fallback function.
	contract := newStatefulPrecompileWithFunctionSelectors(nil, functions)
	return contract
}
