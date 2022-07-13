// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var (
	_                            StatefulPrecompileConfig    = &HelloWorldConfig{}
	ContractHelloWorldPrecompile StatefulPrecompiledContract = createHelloWorldPrecompile()

	helloWorldSignature             = CalculateFunctionSelector("helloWorld()")
	setHelloWorldRecipientSignature = CalculateFunctionSelector("setRecipient(string)")

	nameKey      = common.BytesToHash([]byte("recipient"))
	initialValue = common.BytesToHash([]byte("world!"))
)

type HelloWorldConfig struct {
	BlockTimestamp *big.Int `json:"helloWorldTimestamp"`
}

// Address returns the address of the precompile
func (h *HelloWorldConfig) Address() common.Address { return HelloWorldAddress }

// Return the timestamp at which the precompile is enabled or nil, if it is never enabled
func (h *HelloWorldConfig) Timestamp() *big.Int { return h.BlockTimestamp }

func (h *HelloWorldConfig) Configure(stateDB StateDB) {
	// This will be called in the first block where HelloWorld stateful precompile is enabled.
	// 1) If BlockTimestamp is nil, this will not be called
	// 2) If BlockTimestamp is 0, this will be called while setting up the genesis block
	// 3) If BlockTimestamp is 1000, this will be called while processing the first block whose timestamp is >= 1000
	//
	// Set the initial value under [nameKey] to "world!"
	stateDB.SetState(HelloWorldAddress, nameKey, initialValue)
}

// Return the precompile contract
func (h *HelloWorldConfig) Contract() StatefulPrecompiledContract {
	return ContractHelloWorldPrecompile
}

// Arguments are passed in to functions according to the ABI specification: https://docs.soliditylang.org/en/latest/abi-spec.html.
// Therefore, we maintain compatibility with Solidity by following the same specification while encoding/decoding arguments.
func PackHelloWorldInput(name string) ([]byte, error) {
	byteStr := []byte(name)
	if len(byteStr) > common.HashLength {
		return nil, fmt.Errorf("cannot pack hello world input with string: %s", name)
	}

	input := make([]byte, common.HashLength+len(byteStr))
	strLength := new(big.Int).SetUint64(uint64(len(byteStr)))
	strLengthBytes := strLength.Bytes()
	copy(input[:common.HashLength], strLengthBytes)
	copy(input[common.HashLength:], byteStr)

	return input, nil
}

// UnpackHelloWorldInput unpacks the recipient string from the hello world input
func UnpackHelloWorldInput(input []byte) (string, error) {
	if len(input) < common.HashLength {
		return "", fmt.Errorf("cannot unpack hello world input with length: %d", len(input))
	}

	strLengthBig := new(big.Int).SetBytes(input[:common.HashLength])
	if !strLengthBig.IsUint64() {
		return "", fmt.Errorf("cannot unpack hello world input with stated length that is non-uint64")
	}

	strLength := strLengthBig.Uint64()
	if strLength > common.HashLength {
		return "", fmt.Errorf("cannot unpack hello world string with length: %d", strLength)
	}

	if len(input) != common.HashLength+int(strLength) {
		return "", fmt.Errorf("input had unexpected length %d with string length defined as %d", len(input), strLength)
	}

	str := string(input[common.HashLength:])
	return str, nil
}

func GetReceipient(state StateDB) string {
	value := state.GetState(HelloWorldAddress, nameKey)
	b := value.Bytes()
	trimmedbytes := common.TrimLeftZeroes(b)
	return string(trimmedbytes)
}

// SetRecipient sets the recipient for the hello world precompile
func SetRecipient(state StateDB, recipient string) {
	state.SetState(HelloWorldAddress, nameKey, common.BytesToHash([]byte(recipient)))
}

// sayHello is the execution function of "sayHello()"
func sayHello(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if len(input) != 0 {
		return nil, 0, fmt.Errorf("fuck")
	}
	remainingGas, err = deductGas(suppliedGas, HelloWorldGasCost)
	if err != nil {
		return nil, 0, err
	}

	recipient := GetReceipient(accessibleState.GetStateDB())
	return []byte(fmt.Sprintf("Hello %s!", recipient)), suppliedGas - SetRecipientGasCost, nil
}

// setRecipient is the execution function of "setRecipient(name string)" and sets the recipient in the string returned by hello world
func setRecipient(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	recipient, err := UnpackHelloWorldInput(input)
	if err != nil {
		return nil, 0, err
	}
	remainingGas, err = deductGas(suppliedGas, SetRecipientGasCost)
	if err != nil {
		return nil, 0, err
	}

	SetRecipient(accessibleState.GetStateDB(), recipient)
	return []byte{}, remainingGas, nil
}

// createHelloWorldPrecompile returns the StatefulPrecompile contract that implements the HelloWorld interface from solidity
func createHelloWorldPrecompile() StatefulPrecompiledContract {
	return newStatefulPrecompileWithFunctionSelectors(nil, []*statefulPrecompileFunction{
		newStatefulPrecompileFunction(helloWorldSignature, sayHello),
		newStatefulPrecompileFunction(setHelloWorldRecipientSignature, setRecipient),
	})
}
