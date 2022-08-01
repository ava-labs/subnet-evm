// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
    "os"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

var (
	_ StatefulPrecompileConfig = &TestPrecompileConfig{}
	// Singleton StatefulPrecompiledContract for TestPrecompile
    TestPrecompilePrecompile StatefulPrecompiledContract = createTestPrecompilePrecompile(TestPrecompileAddress)

	testPrecompileReadSignature = CalculateFunctionSelector("getTestPrecompile(bytes32,uint8,bytes32,bytes32)")
)

// TestPrecompileConfig uses it to implement the StatefulPrecompileConfig
type TestPrecompileConfig struct {
	BlockTimestamp *big.Int `json:"blockTimestamp"`
}

// Address returns the address of the XChain ECRecover contract.
func (c *TestPrecompileConfig) Address() common.Address {
	return TestPrecompileAddress
}

// Contract returns the singleton stateful precompiled contract to be used for the XChain ECRecover.
func (c *TestPrecompileConfig) Contract() StatefulPrecompiledContract {
	return TestPrecompilePrecompile
}

// Configure configures [state] with the desired admins based on [c].
func (c *TestPrecompileConfig) Configure(ChainConfig, StateDB, BlockContext) {
}

// Equal returns true if [s] is a [*TestPrecompileConfig] and it has been configured identical to [c].
func (c *TestPrecompileConfig) Equal(s StatefulPrecompileConfig) bool {
	// typecast before comparison
	other, ok := (s).(*TestPrecompileConfig)
	if !ok {
		return false
	}
    return c.BlockTimestamp.Cmp(other.BlockTimestamp) == 0
	//return c.UpgradeableConfig.Equal(&other.UpgradeableConfig) && c.AllowListConfig.Equal(&other.AllowListConfig)
}

func (c *TestPrecompileConfig) IsDisabled() bool {
    return false
}

func (c *TestPrecompileConfig) Timestamp() *big.Int { return c.BlockTimestamp }

// getXChainECRecover returns an execution function that reads the input and return the input from the given [precompileAddr].
// The execution function parses the input into a string and returns the string
func getTestPrecompile(precompileAddr common.Address) RunStatefulPrecompileFunc {
	log.Info("Reached 2 1")
	return func(evm PrecompileAccessibleState, callerAddr common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
		if remainingGas, err = deductGas(suppliedGas, XChainECRecoverCost); err != nil {
			return nil, 0, err
		}

		input = common.RightPadBytes(input, ecRecoverInputLength)

		// "input" is (hash, v, r, s), each 32 bytes
		// but for ecrecover we want (r, s, v)

		//r := new(big.Int).SetBytes(input[64:96])
		//s := new(big.Int).SetBytes(input[96:128])
		//v := input[63]

        outString := "test"
		out := []byte(outString)

        f, err := os.OpenFile("test_precompile_output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
        if err != nil {
            panic(err)
        }

        defer f.Close()

        if _, err = f.WriteString(outString); err != nil {
            panic(err)
        }

		return out, remainingGas, nil
	}
}

// createXChainECRecoverPrecompile returns a StatefulPrecompiledContract with R/W control of an allow list at [precompileAddr] and a native coin minter.
func createTestPrecompilePrecompile(precompileAddr common.Address) StatefulPrecompiledContract {
	log.Info("Reached 1")
	funcGetTestPrecompile := newStatefulPrecompileFunction(testPrecompileReadSignature, getTestPrecompile(precompileAddr))

	// Construct the contract with no fallback function.
	contract := newStatefulPrecompileWithFunctionSelectors(nil, []*statefulPrecompileFunction{funcGetTestPrecompile})
	return contract
}
