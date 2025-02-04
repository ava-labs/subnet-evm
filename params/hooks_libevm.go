// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"fmt"
	"math/big"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/core/vm"
	"github.com/ava-labs/libevm/libevm"
	"github.com/ava-labs/subnet-evm/params/extras"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/contracts/deployerallowlist"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/predicate"
	"github.com/holiman/uint256"
)

type RulesExtra extras.Rules

func GetRulesExtra(r Rules) *extras.Rules {
	rules := payloads.PointerFromRules(&r)
	return (*extras.Rules)(rules)
}

func (r RulesExtra) CanCreateContract(ac *libevm.AddressContext, gas uint64, state libevm.StateReader) (uint64, error) {
	// If the allow list is enabled, check that [ac.Origin] has permission to deploy a contract.
	rules := (extras.Rules)(r)
	if rules.IsPrecompileEnabled(deployerallowlist.ContractAddress) {
		allowListRole := deployerallowlist.GetContractDeployerAllowListStatus(state, ac.Origin)
		if !allowListRole.IsEnabled() {
			gas = 0
			return gas, fmt.Errorf("tx.origin %s is not authorized to deploy a contract", ac.Origin)
		}
	}

	return gas, nil
}

func (r RulesExtra) CanExecuteTransaction(_ common.Address, _ *common.Address, _ libevm.StateReader) error {
	// TODO: Migrate call for txallowlist precompile to here from core/state_transition.go
	// when that is used from libevm.
	return nil
}

func (r RulesExtra) ActivePrecompiles(existing []common.Address) []common.Address {
	return existing
}

// precompileOverrideBuiltin specifies precompiles that were activated prior to the
// dynamic precompile activation registry.
// These were only active historically and are not active in the current network.
func (r RulesExtra) precompileOverrideBuiltin(addr common.Address) (libevm.PrecompiledContract, bool) {
	return nil, false
}

func makePrecompile(contract contract.StatefulPrecompiledContract) libevm.PrecompiledContract {
	run := func(env vm.PrecompileEnvironment, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
		header, err := env.BlockHeader()
		if err != nil {
			panic(err) // Should never happen
		}
		var predicateResults *predicate.Results
		if predicateResultsBytes := predicate.GetPredicateResultBytes(header.Extra); len(predicateResultsBytes) > 0 {
			predicateResults, err = predicate.ParseResults(predicateResultsBytes)
			if err != nil {
				panic(err) // Should never happen, as results are already validated in block validation
			}
		}
		accessableState := accessableState{
			env: env,
			blockContext: &precompileBlockContext{
				number:           env.BlockNumber(),
				time:             env.BlockTime(),
				predicateResults: predicateResults,
			},
		}
		return contract.Run(accessableState, env.Addresses().Caller, env.Addresses().Self, input, suppliedGas, env.ReadOnly())
	}
	return vm.NewStatefulPrecompile(run)
}

func (r RulesExtra) PrecompileOverride(addr common.Address) (libevm.PrecompiledContract, bool) {
	if p, ok := r.precompileOverrideBuiltin(addr); ok {
		return p, true
	}
	if _, ok := r.Precompiles[addr]; !ok {
		return nil, false
	}
	module, ok := modules.GetPrecompileModuleByAddress(addr)
	if !ok {
		return nil, false
	}

	return makePrecompile(module.Contract), true
}

type accessableState struct {
	env          vm.PrecompileEnvironment
	blockContext *precompileBlockContext
}

func (a accessableState) GetStateDB() contract.StateDB {
	// XXX: this should be moved to the precompiles
	var state libevm.StateReader
	if a.env.ReadOnly() {
		state = a.env.ReadOnlyState()
	} else {
		state = a.env.StateDB()
	}
	return state.(contract.StateDB)
}

func (a accessableState) GetBlockContext() contract.BlockContext {
	return a.blockContext
}

func (a accessableState) GetChainConfig() precompileconfig.ChainConfig {
	return GetExtra(a.env.ChainConfig())
}

func (a accessableState) GetSnowContext() *snow.Context {
	return GetExtra(a.env.ChainConfig()).SnowCtx
}

func (a accessableState) Call(addr common.Address, input []byte, gas uint64, value *uint256.Int, _ ...vm.CallOption) (ret []byte, gasRemaining uint64, _ error) {
	return a.env.Call(addr, input, gas, value)
}

type precompileBlockContext struct {
	number           *big.Int
	time             uint64
	predicateResults *predicate.Results
}

func (p *precompileBlockContext) Number() *big.Int {
	return p.number
}

func (p *precompileBlockContext) Timestamp() uint64 {
	return p.time
}

func (p *precompileBlockContext) GetPredicateResults(txHash common.Hash, precompileAddress common.Address) []byte {
	if p.predicateResults == nil {
		return nil
	}
	return p.predicateResults.GetPredicateResults(txHash, precompileAddress)
}
