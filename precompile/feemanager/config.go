// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package feemanager

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ethereum/go-ethereum/common"
)

// FeeConfigManagerConfig wraps [AllowListConfig] and uses it to implement the StatefulPrecompileConfig
// interface while adding in the FeeConfigManager specific precompile address.
type FeeConfigManagerConfig struct {
	allowlist.AllowListConfig // Config for the fee config manager allow list
	precompile.UpgradeableConfig
	InitialFeeConfig *commontype.FeeConfig `json:"initialFeeConfig,omitempty"` // initial fee config to be immediately activated
}

// NewFeeManagerConfig returns a config for a network upgrade at [blockTimestamp] that enables
// FeeConfigManager with the given [admins] and [enableds] as members of the allowlist with [initialConfig] as initial fee config if specified.
func NewFeeManagerConfig(blockTimestamp *big.Int, admins []common.Address, enableds []common.Address, initialConfig *commontype.FeeConfig) *FeeConfigManagerConfig {
	return &FeeConfigManagerConfig{
		AllowListConfig: allowlist.AllowListConfig{
			AllowListAdmins:  admins,
			EnabledAddresses: enableds,
		},
		UpgradeableConfig: precompile.UpgradeableConfig{BlockTimestamp: blockTimestamp},
		InitialFeeConfig:  initialConfig,
	}
}

// NewDisableFeeManagerConfig returns config for a network upgrade at [blockTimestamp]
// that disables FeeConfigManager.
func NewDisableFeeManagerConfig(blockTimestamp *big.Int) *FeeConfigManagerConfig {
	return &FeeConfigManagerConfig{
		UpgradeableConfig: precompile.UpgradeableConfig{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

// Address returns the address of the fee config manager contract.
func (c *FeeConfigManagerConfig) Address() common.Address {
	return precompile.FeeConfigManagerAddress
}

// Equal returns true if [s] is a [*FeeConfigManagerConfig] and it has been configured identical to [c].
func (c *FeeConfigManagerConfig) Equal(s precompile.StatefulPrecompileConfig) bool {
	// typecast before comparison
	other, ok := (s).(*FeeConfigManagerConfig)
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
func (c *FeeConfigManagerConfig) Configure(chainConfig precompile.ChainConfig, state precompile.StateDB, blockContext precompile.BlockContext) error {
	// Store the initial fee config into the state when the fee config manager activates.
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
	return c.AllowListConfig.Configure(state, precompile.FeeConfigManagerAddress)
}

// Contract returns the singleton stateful precompiled contract to be used for the fee manager.
func (c *FeeConfigManagerConfig) Contract() precompile.StatefulPrecompiledContract {
	return FeeConfigManagerPrecompile
}

func (c *FeeConfigManagerConfig) Verify() error {
	if err := c.AllowListConfig.Verify(); err != nil {
		return err
	}
	if c.InitialFeeConfig == nil {
		return nil
	}

	return c.InitialFeeConfig.Verify()
}

// String returns a string representation of the FeeConfigManagerConfig.
func (c *FeeConfigManagerConfig) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}
