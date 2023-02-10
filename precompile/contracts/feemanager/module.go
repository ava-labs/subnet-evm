// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package feemanager

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
)

var _ contract.Module = Module{}

// ConfigKey is the key used in json config files to specify this precompile config.
// must be unique across all precompiles.
const ConfigKey = "feeManagerConfig"

var ContractAddress = common.HexToAddress("0x0200000000000000000000000000000000000003")

type Module struct{}

func (Module) Key() string {
	return ConfigKey
}

// Address returns the address of the fee manager.
func (Module) Address() common.Address {
	return ContractAddress
}

func (Module) NewConfig() config.Config {
	return &FeeManagerConfig{}
}

// Configure configures [state] with the desired admins based on [configIface].
func (Module) Configure(chainConfig contract.ChainConfig, cfg config.Config, state contract.StateDB, blockContext contract.BlockContext) error {
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

func (Module) Contract() contract.StatefulPrecompiledContract {
	return FeeManagerPrecompile
}
