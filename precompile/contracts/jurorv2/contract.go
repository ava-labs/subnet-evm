// Code generated
// This file is a generated precompile contract config with stubbed abstract functions.
// The file is generated by a template. Please inspect every code and comment in this file before use.

package jurorv2

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/contracts/bibliophile"

	_ "embed"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

const (
	// Gas costs for each function. These are set to 1 by default.
	// You should set a gas cost for each function in your contract.
	// Generally, you should not set gas costs very low as this may cause your network to be vulnerable to DoS attacks.
	// There are some predefined gas costs in contract/utils.go that you can use.
	GetNotionalPositionAndMarginGasCost                  uint64 = 69
	ValidateLiquidationOrderAndDetermineFillPriceGasCost uint64 = 69
	ValidateOrdersAndDetermineFillPriceGasCost           uint64 = 69
)

// CUSTOM CODE STARTS HERE
// Reference imports to suppress errors from unused imports. This code and any unnecessary imports can be removed.
var (
	_ = abi.JSON
	_ = errors.New
	_ = big.NewInt
)

// Singleton StatefulPrecompiledContract and signatures.
var (

	// JurorRawABI contains the raw ABI of Juror contract.
	//go:embed contract.abi
	JurorRawABI string

	JurorABI = contract.ParseABI(JurorRawABI)

	JurorPrecompile = createJurorPrecompile()
)

// IClearingHouseInstruction is an auto generated low-level Go binding around an user-defined struct.
type IClearingHouseInstruction struct {
	AmmIndex  *big.Int
	Trader    common.Address
	OrderHash [32]byte
	Mode      uint8
}

// ILimitOrderBookOrder is an auto generated low-level Go binding around an user-defined struct.
type ILimitOrderBookOrder struct {
	AmmIndex          *big.Int
	Trader            common.Address
	BaseAssetQuantity *big.Int
	Price             *big.Int
	Salt              *big.Int
	ReduceOnly        bool
	PostOnly          bool
}

// IOrderHandlerLiquidationMatchingValidationRes is an auto generated low-level Go binding around an user-defined struct.
type IOrderHandlerLiquidationMatchingValidationRes struct {
	Instruction  IClearingHouseInstruction
	OrderType    uint8
	EncodedOrder []byte
	FillPrice    *big.Int
	FillAmount   *big.Int
}

// IOrderHandlerMatchingValidationRes is an auto generated low-level Go binding around an user-defined struct.
type IOrderHandlerMatchingValidationRes struct {
	Instructions  [2]IClearingHouseInstruction
	OrderTypes    [2]uint8
	EncodedOrders [2][]byte
	FillPrice     *big.Int
}

type GetNotionalPositionAndMarginInput struct {
	Trader                 common.Address
	IncludeFundingPayments bool
	Mode                   uint8
}

type GetNotionalPositionAndMarginOutput struct {
	NotionalPosition *big.Int
	Margin           *big.Int
}

type ValidateLiquidationOrderAndDetermineFillPriceInput struct {
	Data              []byte
	LiquidationAmount *big.Int
}

type ValidateLiquidationOrderAndDetermineFillPriceOutput struct {
	Err     string
	Element uint8
	Res     IOrderHandlerLiquidationMatchingValidationRes
}

type ValidateOrdersAndDetermineFillPriceInput struct {
	Data       [2][]byte
	FillAmount *big.Int
}

type ValidateOrdersAndDetermineFillPriceOutput struct {
	Err     string
	Element uint8
	Res     IOrderHandlerMatchingValidationRes
}

// UnpackGetNotionalPositionAndMarginInput attempts to unpack [input] as GetNotionalPositionAndMarginInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackGetNotionalPositionAndMarginInput(input []byte) (GetNotionalPositionAndMarginInput, error) {
	inputStruct := GetNotionalPositionAndMarginInput{}
	err := JurorABI.UnpackInputIntoInterface(&inputStruct, "getNotionalPositionAndMargin", input)

	return inputStruct, err
}

// PackGetNotionalPositionAndMargin packs [inputStruct] of type GetNotionalPositionAndMarginInput into the appropriate arguments for getNotionalPositionAndMargin.
func PackGetNotionalPositionAndMargin(inputStruct GetNotionalPositionAndMarginInput) ([]byte, error) {
	return JurorABI.Pack("getNotionalPositionAndMargin", inputStruct.Trader, inputStruct.IncludeFundingPayments, inputStruct.Mode)
}

// PackGetNotionalPositionAndMarginOutput attempts to pack given [outputStruct] of type GetNotionalPositionAndMarginOutput
// to conform the ABI outputs.
func PackGetNotionalPositionAndMarginOutput(outputStruct GetNotionalPositionAndMarginOutput) ([]byte, error) {
	return JurorABI.PackOutput("getNotionalPositionAndMargin",
		outputStruct.NotionalPosition,
		outputStruct.Margin,
	)
}

func getNotionalPositionAndMargin(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, GetNotionalPositionAndMarginGasCost); err != nil {
		return nil, 0, err
	}
	// attempts to unpack [input] into the arguments to the GetNotionalPositionAndMarginInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackGetNotionalPositionAndMarginInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	bibliophile := bibliophile.NewBibliophileClient(accessibleState)
	log.Info("getNotionalPositionAndMargin", "accessibleState.GetSnowContext().ChainID", accessibleState.GetSnowContext().ChainID.String())
	output := GetNotionalPositionAndMargin(bibliophile, &inputStruct)
	packedOutput, err := PackGetNotionalPositionAndMarginOutput(output)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// UnpackValidateLiquidationOrderAndDetermineFillPriceInput attempts to unpack [input] as ValidateLiquidationOrderAndDetermineFillPriceInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackValidateLiquidationOrderAndDetermineFillPriceInput(input []byte) (ValidateLiquidationOrderAndDetermineFillPriceInput, error) {
	inputStruct := ValidateLiquidationOrderAndDetermineFillPriceInput{}
	err := JurorABI.UnpackInputIntoInterface(&inputStruct, "validateLiquidationOrderAndDetermineFillPrice", input)

	return inputStruct, err
}

// PackValidateLiquidationOrderAndDetermineFillPrice packs [inputStruct] of type ValidateLiquidationOrderAndDetermineFillPriceInput into the appropriate arguments for validateLiquidationOrderAndDetermineFillPrice.
func PackValidateLiquidationOrderAndDetermineFillPrice(inputStruct ValidateLiquidationOrderAndDetermineFillPriceInput) ([]byte, error) {
	return JurorABI.Pack("validateLiquidationOrderAndDetermineFillPrice", inputStruct.Data, inputStruct.LiquidationAmount)
}

// PackValidateLiquidationOrderAndDetermineFillPriceOutput attempts to pack given [outputStruct] of type ValidateLiquidationOrderAndDetermineFillPriceOutput
// to conform the ABI outputs.
func PackValidateLiquidationOrderAndDetermineFillPriceOutput(outputStruct ValidateLiquidationOrderAndDetermineFillPriceOutput) ([]byte, error) {
	return JurorABI.PackOutput("validateLiquidationOrderAndDetermineFillPrice",
		outputStruct.Err,
		outputStruct.Element,
		outputStruct.Res,
	)
}

func validateLiquidationOrderAndDetermineFillPrice(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, ValidateLiquidationOrderAndDetermineFillPriceGasCost); err != nil {
		return nil, 0, err
	}
	// attempts to unpack [input] into the arguments to the ValidateLiquidationOrderAndDetermineFillPriceInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackValidateLiquidationOrderAndDetermineFillPriceInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	bibliophile := bibliophile.NewBibliophileClient(accessibleState)
	output := ValidateLiquidationOrderAndDetermineFillPrice(bibliophile, &inputStruct)
	packedOutput, err := PackValidateLiquidationOrderAndDetermineFillPriceOutput(output)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// UnpackValidateOrdersAndDetermineFillPriceInput attempts to unpack [input] as ValidateOrdersAndDetermineFillPriceInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackValidateOrdersAndDetermineFillPriceInput(input []byte) (ValidateOrdersAndDetermineFillPriceInput, error) {
	inputStruct := ValidateOrdersAndDetermineFillPriceInput{}
	err := JurorABI.UnpackInputIntoInterface(&inputStruct, "validateOrdersAndDetermineFillPrice", input)

	return inputStruct, err
}

// PackValidateOrdersAndDetermineFillPrice packs [inputStruct] of type ValidateOrdersAndDetermineFillPriceInput into the appropriate arguments for validateOrdersAndDetermineFillPrice.
func PackValidateOrdersAndDetermineFillPrice(inputStruct ValidateOrdersAndDetermineFillPriceInput) ([]byte, error) {
	return JurorABI.Pack("validateOrdersAndDetermineFillPrice", inputStruct.Data, inputStruct.FillAmount)
}

// PackValidateOrdersAndDetermineFillPriceOutput attempts to pack given [outputStruct] of type ValidateOrdersAndDetermineFillPriceOutput
// to conform the ABI outputs.
func PackValidateOrdersAndDetermineFillPriceOutput(outputStruct ValidateOrdersAndDetermineFillPriceOutput) ([]byte, error) {
	return JurorABI.PackOutput("validateOrdersAndDetermineFillPrice",
		outputStruct.Err,
		outputStruct.Element,
		outputStruct.Res,
	)
}

func validateOrdersAndDetermineFillPrice(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, ValidateOrdersAndDetermineFillPriceGasCost); err != nil {
		return nil, 0, err
	}
	// attempts to unpack [input] into the arguments to the ValidateOrdersAndDetermineFillPriceInput.
	// Assumes that [input] does not include selector
	// You can use unpacked [inputStruct] variable in your code
	inputStruct, err := UnpackValidateOrdersAndDetermineFillPriceInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	// CUSTOM CODE STARTS HERE
	log.Info("validateOrdersAndDetermineFillPrice", "inputStruct", inputStruct)
	bibliophile := bibliophile.NewBibliophileClient(accessibleState)
	output := ValidateOrdersAndDetermineFillPrice(bibliophile, &inputStruct)
	packedOutput, err := PackValidateOrdersAndDetermineFillPriceOutput(output)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// createJurorPrecompile returns a StatefulPrecompiledContract with getters and setters for the precompile.

func createJurorPrecompile() contract.StatefulPrecompiledContract {
	var functions []*contract.StatefulPrecompileFunction

	abiFunctionMap := map[string]contract.RunStatefulPrecompileFunc{
		"getNotionalPositionAndMargin":                  getNotionalPositionAndMargin,
		"validateLiquidationOrderAndDetermineFillPrice": validateLiquidationOrderAndDetermineFillPrice,
		"validateOrdersAndDetermineFillPrice":           validateOrdersAndDetermineFillPrice,
	}

	for name, function := range abiFunctionMap {
		method, ok := JurorABI.Methods[name]
		if !ok {
			panic(fmt.Errorf("given method (%s) does not exist in the ABI", name))
		}
		functions = append(functions, contract.NewStatefulPrecompileFunction(method.ID, function))
	}
	// Construct the contract with no fallback function.
	statefulContract, err := contract.NewStatefulPrecompileContract(nil, functions)
	if err != nil {
		panic(err)
	}
	return statefulContract
}
