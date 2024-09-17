// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"fmt"
	"math/big"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/subnet-evm/constants"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/contracts/deployerallowlist"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/vmerrs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/libevm"
	gethparams "github.com/ethereum/go-ethereum/params"
)

var PredicateParser = func(extra []byte) (PredicateResults, error) {
	return nil, nil
}

func (r RulesExtra) JumpTable() interface{} {
	// XXX: This does not account for the any possible differences in EIP-3529
	// Do not merge without verifying.
	return nil
}

func (r RulesExtra) CanCreateContract(ac *libevm.AddressContext, gas uint64, state libevm.StateReader) (uint64, error) {
	// IsProhibited
	if ac.Self == constants.BlackholeAddr || modules.ReservedAddress(ac.Self) {
		return gas, vmerrs.ErrAddrProhibited
	}

	// If the allow list is enabled, check that [ac.Origin] has permission to deploy a contract.
	if r.IsPrecompileEnabled(deployerallowlist.ContractAddress) {
		allowListRole := deployerallowlist.GetContractDeployerAllowListStatus(state, ac.Origin)
		if !allowListRole.IsEnabled() {
			gas = 0
			return gas, fmt.Errorf("tx.origin %s is not authorized to deploy a contract", ac.Origin)
		}
	}

	return gas, nil
}

func (r RulesExtra) PrecompileOverride(addr common.Address) (libevm.PrecompiledContract, bool) {
	if _, ok := r.ActivePrecompiles[addr]; !ok {
		return nil, false
	}
	module, ok := modules.GetPrecompileModuleByAddress(addr)
	if !ok {
		return nil, false
	}

	precompile := func(env vm.PrecompileEnvironment, input []byte, suppliedGas uint64) ([]byte, uint64, error) {
		header, err := env.BlockHeader()
		if err != nil {
			panic(err) // Should never happen
		}
		predicateResults, err := PredicateParser(header.Extra)
		if err != nil {
			panic(err) // Should never happen, because predicates are parsed in NewEVMBlockContext.
		}
		// XXX: this should be moved to the precompiles
		var state libevm.StateReader
		if env.ReadOnly() {
			state = env.ReadOnlyState()
		} else {
			state = env.StateDB()
		}
		accessableState := accessableState{
			StateReader: state,
			chainConfig: GetRulesExtra(env.Rules()).chainConfig,
			blockContext: &BlockContext{
				number:           env.BlockNumber(),
				time:             env.BlockTime(),
				predicateResults: predicateResults,
			}}
		return module.Contract.Run(accessableState, env.Addresses().Caller, env.Addresses().Self, input, suppliedGas, env.ReadOnly())
	}
	return vm.NewStatefulPrecompile(precompile), true
}

type accessableState struct {
	libevm.StateReader
	chainConfig  *gethparams.ChainConfig
	blockContext *BlockContext
}

func (a accessableState) GetStateDB() contract.StateDB {
	// XXX: Whoa, this is a hack
	return a.StateReader.(contract.StateDB)
}

func (a accessableState) GetBlockContext() contract.BlockContext {
	return a.blockContext
}

func (a accessableState) GetChainConfig() precompileconfig.ChainConfig {
	extra := GetExtra(a.chainConfig)
	return extra
}

func (a accessableState) GetSnowContext() *snow.Context {
	return GetExtra(a.chainConfig).SnowCtx
}

type PredicateResults interface {
	GetPredicateResults(txHash common.Hash, address common.Address) []byte
}

type BlockContext struct {
	number           *big.Int
	time             uint64
	predicateResults PredicateResults
}

func NewBlockContext(number *big.Int, time uint64, predicateResults PredicateResults) *BlockContext {
	return &BlockContext{
		number:           number,
		time:             time,
		predicateResults: predicateResults,
	}
}

func (b *BlockContext) Number() *big.Int {
	return b.number
}

func (b *BlockContext) Timestamp() uint64 {
	return b.time
}

func (b *BlockContext) GetPredicateResults(txHash common.Hash, address common.Address) []byte {
	if b.predicateResults == nil {
		return nil
	}
	return b.predicateResults.GetPredicateResults(txHash, address)
}
