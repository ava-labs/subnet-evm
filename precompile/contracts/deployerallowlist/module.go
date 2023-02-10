// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package deployerallowlist

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
)

var _ contract.Module = Module{}

// ConfigKey is the key used in json config files to specify this precompile config.
// must be unique across all precompiles.
const ConfigKey = "contractDeployerAllowListConfig"

var ContractAddress = common.HexToAddress("0x0200000000000000000000000000000000000000")

type Module struct{}

func (Module) Key() string {
	return ConfigKey
}

// Address returns the address of the contract deployer allow list.
func (Module) Address() common.Address {
	return ContractAddress
}

func (Module) NewConfig() config.Config {
	return &ContractDeployerAllowListConfig{}
}

// Configure configures [state] with the desired admins based on [cfg].
func (Module) Configure(_ contract.ChainConfig, cfg config.Config, state contract.StateDB, _ contract.BlockContext) error {
	config, ok := cfg.(*ContractDeployerAllowListConfig)
	if !ok {
		return fmt.Errorf("incorrect config %T: %v", config, config)
	}
	return config.AllowListConfig.Configure(state, ContractAddress)
}

func (Module) Contract() contract.StatefulPrecompiledContract {
	return ContractDeployerAllowListPrecompile
}
