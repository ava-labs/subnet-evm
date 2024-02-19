// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package rewardmanager

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"

	"github.com/ethereum/go-ethereum/common"
)

var _ contract.Configurator = &configurator{}

// ConfigKey is the key used in json config files to specify this precompile config.
// must be unique across all precompiles.
const ConfigKey = "rewardManagerConfig"

var ContractAddress = common.HexToAddress("0x0200000000000000000000000000000000000004")

// Module is the precompile module. It is used to register the precompile contract.
var Module = modules.Module{
	ConfigKey:    ConfigKey,
	Address:      ContractAddress,
	Contract:     RewardManagerPrecompile,
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

// Configure configures [state] with the given [cfg] precompileconfig.
// This function is called by the EVM once per precompile contract activation.
// You can use this function to set up your precompile contract's initial state,
// by using the [cfg] config and [state] stateDB.
func (*configurator) Configure(chainConfig precompileconfig.ChainConfig, cfg precompileconfig.Config, state contract.StateDB, blockContext contract.ConfigurationBlockContext) error {
	config, ok := cfg.(*Config)
	if !ok {
		return fmt.Errorf("expected config type %T, got %T: %v", &Config{}, cfg, cfg)
	}
	// configure the RewardManager with the given initial configuration
	if config.InitialRewardConfig != nil {
		config.InitialRewardConfig.Configure(state)
	} else if chainConfig.AllowedFeeRecipients() {
		// configure the RewardManager according to chainConfig
		EnableAllowFeeRecipients(state)
	} else {
		// chainConfig does not have any reward address
		// if chainConfig does not enable fee recipients
		// default to disabling rewards
		DisableFeeRewards(state)
	}
	return config.AllowListConfig.Configure(chainConfig, ContractAddress, state, blockContext)
}
