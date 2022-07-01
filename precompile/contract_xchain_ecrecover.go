// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

var (
	_ StatefulPrecompileConfig = &ContractXChainECRecoverConfig{}
	// Singleton StatefulPrecompiledContract for XChain ECRecover.
	ContractXChainECRecoverPrecompile StatefulPrecompiledContract = createXChainECRecoverPrecompile(ContractXchainECRecoverAddress)

	xChainECRecoverSignature = CalculateFunctionSelector("xChainECRecover(string)") // address, amount
	xChainECRecoverReadSignature = CalculateFunctionSelector("getXChainECRecover(string)")
)

// ContractXChainECRecoverConfig uses it to implement the StatefulPrecompileConfig
type ContractXChainECRecoverConfig struct {
	BlockTimestamp *big.Int `json:"blockTimestamp"`
}

// Address returns the address of the XChain ECRecover contract.
func (c *ContractXChainECRecoverConfig) Address() common.Address {
	return ContractXchainECRecoverAddress
}

// Contract returns the singleton stateful precompiled contract to be used for the XChain ECRecover.
func (c *ContractXChainECRecoverConfig) Contract() StatefulPrecompiledContract {
	return ContractXChainECRecoverPrecompile
}

// Configure configures [state] with the desired admins based on [c].
func (c *ContractXChainECRecoverConfig) Configure(state StateDB) {
	
}

func (c *ContractXChainECRecoverConfig) Timestamp() *big.Int { return c.BlockTimestamp }

// getXChainECRecover returns an execution function that reads the input and return the input from the given [precompileAddr].
// The execution function parses the input into a string and returns the string
func getXChainECRecover(precompileAddr common.Address) RunStatefulPrecompileFunc {
	log.Info("Reached 2 1");
	return func(evm PrecompileAccessibleState, callerAddr common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
		if remainingGas, err = deductGas(suppliedGas, XChainECRecoverCost); err != nil {
			return nil, 0, err
		}
		log.Info("Reached 2 2");
		log.Info(string(input[:]));

		out := []byte(string(input[:]))
		return out, remainingGas, nil
	}
}

// createXChainECRecoverPrecompile returns a StatefulPrecompiledContract with R/W control of an allow list at [precompileAddr] and a native coin minter.
func createXChainECRecoverPrecompile(precompileAddr common.Address) StatefulPrecompiledContract {
	log.Info("Reached 1");
	funcGetXChainECRecover := newStatefulPrecompileFunction(xChainECRecoverReadSignature, getXChainECRecover(precompileAddr))

	// Construct the contract with no fallback function.
	contract := newStatefulPrecompileWithFunctionSelectors(nil, []*statefulPrecompileFunction{xChainECRecover,funcGetXChainECRecover})
	return contract
}
