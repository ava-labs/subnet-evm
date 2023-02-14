// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package deployerallowlist

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ethereum/go-ethereum/common"
)

// ConfigKey is the key used in json config files to specify this precompile config.
// must be unique across all precompiles.
const ConfigKey = "contractDeployerAllowListConfig"

var ContractAddress = common.HexToAddress("0x0200000000000000000000000000000000000000")

var Module = modules.Module{
	ConfigKey:    ConfigKey,
	Address:      ContractAddress,
	Contract:     ContractDeployerAllowListPrecompile,
	Configurator: &configuror{},
}

type configuror struct{}

func init() {
	modules.RegisterModule(Module)
}

func (*configuror) NewConfig() config.Config {
	return &ContractDeployerAllowListConfig{}
}

// Configure configures [state] with the desired admins based on [cfg].
func (c *configuror) Configure(_ contract.ChainConfig, cfg config.Config, state contract.StateDB, _ contract.BlockContext) error {
	config, ok := cfg.(*ContractDeployerAllowListConfig)
	if !ok {
		return fmt.Errorf("incorrect config %T: %v", config, config)
	}
	return config.AllowListConfig.Configure(state, ContractAddress)
}
