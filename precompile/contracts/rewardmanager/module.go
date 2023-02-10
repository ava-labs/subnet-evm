// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package rewardmanager

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
)

// ConfigKey is the key used in json config files to specify this precompile config.
// must be unique across all precompiles.
const ConfigKey = "rewardManagerConfig"

var ContractAddress = common.HexToAddress("0x0200000000000000000000000000000000000004")

type Module struct{}

func (Module) Key() string {
	return ConfigKey
}

// Address returns the address of the reward manager.
func (Module) Address() common.Address {
	return ContractAddress
}

func (Module) NewConfig() config.Config {
	return &RewardManagerConfig{}
}

// Configure configures [state] with the desired admins based on [cfg].
func (Module) Configure(chainConfig contract.ChainConfig, cfg config.Config, state contract.StateDB, _ contract.BlockContext) error {
	config, ok := cfg.(*RewardManagerConfig)
	if !ok {
		return fmt.Errorf("incorrect config %T: %v", config, config)
	}
	// TODO: should we move this to the end and return it for consistency with other precompiles?
	if err := config.AllowListConfig.Configure(state, ContractAddress); err != nil {
		return err
	}
	// configure the RewardManager with the given initial configuration
	if config.InitialRewardConfig != nil {
		return config.InitialRewardConfig.Configure(state)
	} else if chainConfig.AllowedFeeRecipients() {
		// configure the RewardManager according to chainConfig
		EnableAllowFeeRecipients(state)
	} else {
		// chainConfig does not have any reward address
		// if chainConfig does not enable fee recipients
		// default to disabling rewards
		DisableFeeRewards(state)
	}
	return nil
}

func (Module) Contract() contract.StatefulPrecompiledContract {
	return RewardManagerPrecompile
}
