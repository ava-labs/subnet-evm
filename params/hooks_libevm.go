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
	"github.com/ethereum/go-ethereum/libevm"
	gethparams "github.com/ethereum/go-ethereum/params"
)

func (r RulesExtra) JumpTable() interface{} {
	// XXX: This does not account for the any possible differences in EIP-3529
	// Do not merge without verifying.
	return nil
}

func (r RulesExtra) CanCreateContract(ac *libevm.AddressContext, state libevm.StateReader) error {
	// IsProhibited
	if ac.Self == constants.BlackholeAddr || modules.ReservedAddress(ac.Self) {
		return vmerrs.ErrAddrProhibited
	}

	// If the allow list is enabled, check that [ac.Origin] has permission to deploy a contract.
	if r.IsPrecompileEnabled(deployerallowlist.ContractAddress) {
		allowListRole := deployerallowlist.GetContractDeployerAllowListStatus(state, ac.Origin)
		if !allowListRole.IsEnabled() {
			ac.Gas = 0
			return fmt.Errorf("tx.origin %s is not authorized to deploy a contract", ac.Origin)
		}
	}

	return nil
}

func (r RulesExtra) PrecompileOverride(addr common.Address) (libevm.PrecompiledContract, bool) {
	if _, ok := r.ActivePrecompiles[addr]; !ok {
		return nil, false
	}
	module, ok := modules.GetPrecompileModuleByAddress(addr)
	if !ok {
		return nil, false
	}

	return libevmContract{module.Contract}, true
}

// XXX: This is a hack since we need the suppliedGas
// to determine the gas cost of the precompile
// also evm.interpreter.ReadOnly
type libevmContract struct {
	contract.StatefulPrecompiledContract
}

func (l libevmContract) RunExtra(
	chainConfig *gethparams.ChainConfig,
	blockNumber *big.Int, blockTime uint64, predicateResults libevm.PredicateResults,
	state libevm.StateDB, _ *gethparams.Rules, caller, self common.Address, input []byte, suppliedGas uint64, readOnly bool) ([]byte, uint64, error) {
	accessableState := accessableState{
		StateDB:     state,
		chainConfig: chainConfig,
		blockContext: &BlockContext{
			number:           blockNumber,
			time:             blockTime,
			predicateResults: predicateResults,
		}}
	return l.StatefulPrecompiledContract.Run(accessableState, caller, self, input, suppliedGas, readOnly)
}

func (libevmContract) Run(input []byte) ([]byte, error) {
	panic("implement me")
}

func (libevmContract) RequiredGas(input []byte) uint64 {
	panic("implement me")
}

type accessableState struct {
	libevm.StateDB
	chainConfig  *gethparams.ChainConfig
	blockContext *BlockContext
}

func (a accessableState) GetStateDB() contract.StateDB {
	// XXX: Whoa, this is a hack
	return a.StateDB.(contract.StateDB)
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

type BlockContext struct {
	number           *big.Int
	time             uint64
	predicateResults libevm.PredicateResults
}

func NewBlockContext(number *big.Int, time uint64, predicateResults libevm.PredicateResults) *BlockContext {
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
