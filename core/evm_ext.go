// (c) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package core

import (
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/subnet-evm/constants"
	"github.com/ava-labs/subnet-evm/core/vm"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/contracts/deployerallowlist"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

var defaultEVMFactory vm.EvmFactory

func init() {
	defaultEVMFactory = vm.DefaultEVMFactory()
	vm.SetDefaultEVMFactory(&evmFactory{})
}

// IsProhibited returns true if [addr] is in the prohibited list of addresses which should
// not be allowed as an EOA or newly created contract address.
func IsProhibited(addr common.Address) bool {
	if addr == constants.BlackholeAddr {
		return true
	}

	return modules.ReservedAddress(addr)
}

type evmFactory struct{}

type EVM struct {
	*vm.EVM

	chainConfig *params.ChainConfig
}

func (*evmFactory) NewEVM(
	blockCtx vm.BlockContext, txCtx vm.TxContext, statedb vm.StateDB,
	chainConfig vm.ChainConfig, config vm.Config,
) *vm.EVM {
	customChainConfig, ok := chainConfig.(*params.ChainConfig)
	if !ok {
		// If the chainConfig is not a params.ChainConfig, then we can't use the custom
		// EVM implementation, so we fall back to the default implementation.
		log.Warn("ChainConfig is not a *params.ChainConfig, falling back to default EVM")
		return defaultEVMFactory.NewEVM(blockCtx, txCtx, statedb, chainConfig, config)
	}
	evm := &EVM{
		chainConfig: customChainConfig,
	}
	config.IsProhibited = IsProhibited
	config.DeployerAllowed = evm.DeployerAllowed
	config.CustomPrecompiles = evm.CustomPrecompiles

	evm.EVM = defaultEVMFactory.NewEVM(blockCtx, txCtx, statedb, chainConfig, config)
	return evm.EVM
}

func (evm *EVM) GetBlockContext() contract.BlockContext {
	return &evm.EVM.Context
}

func (evm *EVM) GetStateDB() contract.StateDB {
	return evm.StateDB
}

func (evm *EVM) GetChainConfig() precompileconfig.ChainConfig {
	return evm.chainConfig
}

func (evm *EVM) GetSnowContext() *snow.Context {
	return evm.chainConfig.SnowCtx
}

func (evm *EVM) DeployerAllowed(addr common.Address) bool {
	rules := evm.chainConfig.Rules(evm.Context.BlockNumber, evm.Context.Time)
	if rules.IsPrecompileEnabled(deployerallowlist.ContractAddress) {
		allowListRole := deployerallowlist.GetContractDeployerAllowListStatus(evm.StateDB, evm.TxContext.Origin)
		if !allowListRole.IsEnabled() {
			return false
		}
	}
	return true
}

func (evm *EVM) CustomPrecompiles(addr common.Address) (vm.RunFunc, bool) {
	rules := evm.chainConfig.Rules(evm.Context.BlockNumber, evm.Context.Time)
	if _, ok := rules.ActivePrecompiles[addr]; !ok {
		return nil, false
	}
	module, ok := modules.GetPrecompileModuleByAddress(addr)
	if !ok {
		return nil, false
	}

	return func(
		caller common.Address, input []byte, suppliedGas uint64, readOnly bool,
	) (ret []byte, remainingGas uint64, err error) {
		return module.Contract.Run(evm, caller, addr, input, suppliedGas, readOnly)
	}, true
}
