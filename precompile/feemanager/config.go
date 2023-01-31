// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package feemanager

import (
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ethereum/go-ethereum/common"
)

var _ precompile.StatefulPrecompileConfig = &FeeManagerConfig{}

// ConfigKey is the key used in json config files to specify this precompile config.
// must be unique across all precompiles.
const ConfigKey = "feeManagerConfig"

// FeeManagerConfig wraps [AllowListConfig] and uses it to implement the StatefulPrecompileConfig
// interface while adding in the FeeManager specific precompile address.
type FeeManagerConfig struct {
	allowlist.AllowListConfig // Config for the fee config manager allow list
	precompile.UpgradeableConfig
	InitialFeeConfig *commontype.FeeConfig `json:"initialFeeConfig,omitempty"` // initial fee config to be immediately activated
}

// NewModule returns a new module for FeeManager.
func NewModule() precompile.StatefulPrecompileModule {
	return &FeeManagerConfig{}
}

// NewFeeManagerConfig returns a config for a network upgrade at [blockTimestamp] that enables
// FeeManager with the given [admins] and [enableds] as members of the allowlist with [initialConfig] as initial fee config if specified.
func NewFeeManagerConfig(blockTimestamp *big.Int, admins []common.Address, enableds []common.Address, initialConfig *commontype.FeeConfig) *FeeManagerConfig {
	return &FeeManagerConfig{
		AllowListConfig: allowlist.AllowListConfig{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
		},
		UpgradeableConfig: precompile.UpgradeableConfig{BlockTimestamp: blockTimestamp},
		InitialFeeConfig:  initialConfig,
	}
}

// NewDisableFeeManagerConfig returns config for a network upgrade at [blockTimestamp]
// that disables FeeManager.
func NewDisableFeeManagerConfig(blockTimestamp *big.Int) *FeeManagerConfig {
	return &FeeManagerConfig{
		UpgradeableConfig: precompile.UpgradeableConfig{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

// Address returns the address of the fee manager contract.
func (FeeManagerConfig) Address() common.Address {
	return ContractAddress
}

// Contract returns the singleton stateful precompiled contract to be used for the fee manager.
func (FeeManagerConfig) Contract() precompile.StatefulPrecompiledContract {
	return FeeManagerPrecompile
}

// Key returns the key used in json config files to specify this precompile config.
func (FeeManagerConfig) Key() string {
	return ConfigKey
}

// NewConfig returns a new instance of this config.
func (FeeManagerConfig) NewConfig() precompile.StatefulPrecompileConfig {
	return new(FeeManagerConfig)
}

// Equal returns true if [s] is a [*FeeManagerConfig] and it has been configured identical to [c].
func (c *FeeManagerConfig) Equal(s precompile.StatefulPrecompileConfig) bool {
	// typecast before comparison
	other, ok := (s).(*FeeManagerConfig)
	if !ok {
		return false
	}
	eq := c.UpgradeableConfig.Equal(&other.UpgradeableConfig) && c.AllowListConfig.Equal(&other.AllowListConfig)
	if !eq {
		return false
	}

	if c.InitialFeeConfig == nil {
		return other.InitialFeeConfig == nil
	}

	return c.InitialFeeConfig.Equal(other.InitialFeeConfig)
}

// Configure configures [state] with the desired admins based on [c].
func (c *FeeManagerConfig) Configure(chainConfig precompile.ChainConfig, state precompile.StateDB, blockContext precompile.BlockContext) error {
	// Store the initial fee config into the state when the fee manager activates.
	if c.InitialFeeConfig != nil {
		if err := StoreFeeConfig(state, *c.InitialFeeConfig, blockContext); err != nil {
			// This should not happen since we already checked this config with Verify()
			return fmt.Errorf("cannot configure given initial fee config: %w", err)
		}
	} else {
		if err := StoreFeeConfig(state, chainConfig.GetFeeConfig(), blockContext); err != nil {
			// This should not happen since we already checked the chain config in the genesis creation.
			return fmt.Errorf("cannot configure fee config in chain config: %w", err)
		}
	}
	return c.AllowListConfig.Configure(state, ContractAddress)
}

func (c *FeeManagerConfig) Verify() error {
	if err := c.AllowListConfig.Verify(); err != nil {
		return err
	}
	if c.InitialFeeConfig == nil {
		return nil
	}

	return c.InitialFeeConfig.Verify()
}
