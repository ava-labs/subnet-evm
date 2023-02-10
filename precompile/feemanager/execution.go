// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package feemanager

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/precompile/contract"
)

type Executor struct{}

// Configure configures [state] with the desired admins based on [configIface].
func (Executor) Configure(chainConfig contract.ChainConfig, cfg config.Config, state contract.StateDB, blockContext contract.BlockContext) error {
	config, ok := cfg.(*FeeManagerConfig)
	if !ok {
		return fmt.Errorf("incorrect config %T: %v", config, config)
	}
	// Store the initial fee config into the state when the fee manager activates.
	if config.InitialFeeConfig != nil {
		if err := StoreFeeConfig(state, *config.InitialFeeConfig, blockContext); err != nil {
			// This should not happen since we already checked this config with Verify()
			return fmt.Errorf("cannot configure given initial fee config: %w", err)
		}
	} else {
		if err := StoreFeeConfig(state, chainConfig.GetFeeConfig(), blockContext); err != nil {
			// This should not happen since we already checked the chain config in the genesis creation.
			return fmt.Errorf("cannot configure fee config in chain config: %w", err)
		}
	}
	return config.AllowListConfig.Configure(state, ContractAddress)
}

func (Executor) Contract() contract.StatefulPrecompiledContract {
	return FeeManagerPrecompile
}
