// Code generated
// This file is a generated precompile contract config with stubbed abstract functions.
// The file is generated by a template. Please inspect every code and comment in this file before use.

package counter

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"

	"github.com/ethereum/go-ethereum/common"
)

var _ contract.Configurator = &configurator{}

// ConfigKey is the key used in json config files to specify this precompile precompileconfig.
// must be unique across all precompiles.
const ConfigKey = "counterConfig"

var ContractAddress = common.HexToAddress("0x0300000000000000000000000000000000000099") // SET A SUITABLE HEX ADDRESS HERE

// Module is the precompile module. It is used to register the precompile contract.
var Module = modules.Module{
	ConfigKey:    ConfigKey,
	Address:      ContractAddress,
	Contract:     ICounter,
	Configurator: &configurator{},
}

type configurator struct{}

func init() {
	// Register the precompile module.
	// Each precompile contract registers itself through [RegisterModule] function.
	if err := modules.RegisterModule(Module); err != nil {
		panic(err)
	}
}

// MakeConfig returns a new precompile config instance.
// This is required for Marshal/Unmarshal the precompile config.
func (*configurator) MakeConfig() precompileconfig.Config {
	return new(Config)
}

const defaultCounter = 99

func (*configurator) Configure(chainConfig contract.ChainConfig, cfg precompileconfig.Config, state contract.StateDB, _ contract.BlockContext) error {
	config, ok := cfg.(*Config)
	if !ok {
		return fmt.Errorf("incorrect config %T: %v", config, config)
	}

	// AllowList is activated for this precompile. Configuring allowlist addresses here.
	return config.AllowListConfig.Configure(state, ContractAddress)
}
