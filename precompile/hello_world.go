// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	hello "github.com/ava-labs/subnet-evm/precompile/hello/contracts"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_                            StatefulPrecompileConfig    = &HelloWorldConfig{}
	ContractHelloWorldPrecompile StatefulPrecompiledContract = createHelloWorldPrecompile()

	helloWorldSignature             = CalculateFunctionSelector("sayHello()")
	setHelloWorldRecipientSignature = CalculateFunctionSelector("setGreeting(string)")

	nameKey       = common.BytesToHash([]byte("recipient"))
	helloWorldStr = "Hello World!"

	ErrInvalidGreeting = errors.New("invalid input length to say hello")

	helloABI abi.ABI // The ABI for the hello world interface
)

func init() {
	parsed, err := abi.JSON(strings.NewReader(hello.HelloWorldABI))
	if err != nil {
		panic(err)
	}

	helloABI = parsed
}

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
	// Set the initial value under [nameKey] to "Hello World!"
	SetGreeting(stateDB, helloWorldStr)
}

// Return the precompile contract
func (h *HelloWorldConfig) Contract() StatefulPrecompiledContract {
	return ContractHelloWorldPrecompile
}

// PackSayHelloInput returns the calldata necessary to call HelloWorld's sayHello
func PackSayHelloInput() []byte {
	return common.CopyBytes(helloWorldSignature)
}

// Arguments are passed in to functions according to the ABI specification: https://docs.soliditylang.org/en/latest/abi-spec.html.
// Therefore, we maintain compatibility with Solidity by following the same specification while encoding/decoding arguments.
func PackHelloWorldSetGreetingInput(name string) ([]byte, error) {
	if len([]byte(name)) > common.HashLength {
		return nil, fmt.Errorf("cannot pack hello world input with string: %s", name)
	}
	return helloABI.Pack("setGreeting", name)
}

// UnpackHelloWorldInput unpacks the recipient string from the hello world input
func UnpackHelloWorldSetGreetingInput(input []byte) (string, error) {
	res, err := helloABI.Methods["setGreeting"].Inputs.Unpack(input)
	if err != nil {
		return "", err
	}

	if len(res) != 1 {
		return "", fmt.Errorf("unexpected response length: %d", len(res))
	}
	str, ok := res[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected response type: %T of %v", res[0], res[0])
	}

	byteStr := []byte(str)
	if len(byteStr) > common.HashLength {
		return "", fmt.Errorf("cannot unpack string of byte length %d", len(byteStr))
	}

	return str, nil
}

func GetGreeting(state StateDB) string {
	value := state.GetState(HelloWorldAddress, nameKey)
	b := value.Bytes()
	trimmedbytes := common.TrimLeftZeroes(b)
	return string(trimmedbytes)
}

// SetGreeting sets the recipient for the hello world precompile
func SetGreeting(state StateDB, recipient string) {
	res := common.LeftPadBytes([]byte(recipient), common.HashLength)
	state.SetState(HelloWorldAddress, nameKey, common.BytesToHash(res))
}

// sayHello is the execution function of "sayHello()"
func sayHello(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if len(input) != 0 {
		return nil, 0, fmt.Errorf("%w: %d", ErrInvalidGreeting, len(input))
	}
	remainingGas, err = deductGas(suppliedGas, HelloWorldGasCost)
	if err != nil {
		return nil, 0, err
	}

	recipient := GetGreeting(accessibleState.GetStateDB())
	return []byte(recipient), remainingGas, nil
}

// setGreeting is the execution function of "SetGreeting(name string)" and sets the recipient in the string returned by hello world
func setGreeting(accessibleState PrecompileAccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	recipient, err := UnpackHelloWorldSetGreetingInput(input)
	if err != nil {
		return nil, 0, err
	}
	remainingGas, err = deductGas(suppliedGas, SetGreetingGasCost)
	if err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, 0, vmerrs.ErrWriteProtection
	}

	SetGreeting(accessibleState.GetStateDB(), recipient)
	return []byte{}, remainingGas, nil
}

// createHelloWorldPrecompile returns the StatefulPrecompile contract that implements the HelloWorld interface from solidity
func createHelloWorldPrecompile() StatefulPrecompiledContract {
	return newStatefulPrecompileWithFunctionSelectors(nil, []*statefulPrecompileFunction{
		newStatefulPrecompileFunction(helloWorldSignature, sayHello),
		newStatefulPrecompileFunction(setHelloWorldRecipientSignature, setGreeting),
	})
}
