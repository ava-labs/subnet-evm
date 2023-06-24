package counter

import (
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/vmerrs"

	_ "embed"

	"github.com/ethereum/go-ethereum/common"
)

const (
	// Define gas costs for the new functions
	IncrementByXGasCost   uint64 = contract.WriteGasCostPerSlot
	GetCounterGasCost     uint64 = contract.ReadGasCostPerSlot
)

// Gas costs for stateful precompiles
const (
    WriteGasCostPerSlot = 20_000
    ReadGasCostPerSlot  = 5_000
)

var (
	// Define a new storage key for the counter
	counterKeyHash = common.BytesToHash([]byte("counterKey"))
)

// Implement the IncrementByOne function
func incrementByOne(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	// Call incrementByX with 1 as the increment value
	return incrementByX(accessibleState, caller, addr, big.NewInt(1).Bytes(), suppliedGas, readOnly)
}

// Implement the IncrementByX function
func incrementByX(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, IncrementByXGasCost); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}

	x := new(big.Int).SetBytes(input)
	stateDB := accessibleState.GetStateDB()
	counter := new(big.Int).SetBytes(stateDB.GetState(ContractAddress, counterKeyHash).Bytes())
	counter.Add(counter, x)
	stateDB.SetState(ContractAddress, counterKeyHash, common.BigToHash(counter))

	return nil, remainingGas, nil
}

// Implement the getCounter function
func getCounter(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, GetCounterGasCost); err != nil {
		return nil, 0, err
	}

	stateDB := accessibleState.GetStateDB()
	counter := stateDB.GetState(ContractAddress, counterKeyHash).Bytes()

	return counter, remainingGas, nil
}

// Add the new functions to the contract
func createCounterPrecompile() contract.StatefulPrecompiledContract {
	var functions []*contract.StatefulPrecompileFunction
	functions = append(functions, allowlist.CreateAllowListFunctions(ContractAddress)...)

	abiFunctionMap := map[string]contract.RunStatefulPrecompileFunc{
		"IncrementByOne": incrementByOne,
		"IncrementByX":   incrementByX,
		"getCounter":     getCounter,
	}

	for name, function := range abiFunctionMap {
		method, ok := Counter.Methods[name]
		if !ok {
			panic(fmt.Errorf("given method (%s) does not exist in the ABI", name))
		}
		functions = append(functions, contract.NewStatefulPrecompileFunction(method.ID, function))
	}

}
