// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package rewardmanager

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/precompile/execution"
)

type Executor struct{}

// Configure configures [state] with the desired admins based on [cfg].
func (Executor) Configure(chainConfig execution.ChainConfig, cfg config.Config, state execution.StateDB, _ execution.BlockContext) error {
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

func (Executor) Contract() execution.Contract {
	return RewardManagerPrecompile
}
