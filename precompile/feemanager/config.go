// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package feemanager

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ethereum/go-ethereum/common"
)

var _ config.Config = &FeeManagerConfig{}

// FeeManagerConfig wraps [AllowListConfig] and uses it to implement the StatefulPrecompileConfig
// interface while adding in the FeeManager specific precompile address.
type FeeManagerConfig struct {
	allowlist.AllowListConfig // Config for the fee config manager allow list
	config.UpgradeableConfig
	InitialFeeConfig *commontype.FeeConfig `json:"initialFeeConfig,omitempty"` // initial fee config to be immediately activated
}

// NewFeeManagerConfig returns a config for a network upgrade at [blockTimestamp] that enables
// FeeManager with the given [admins] and [enableds] as members of the allowlist with [initialConfig] as initial fee config if specified.
func NewFeeManagerConfig(blockTimestamp *big.Int, admins []common.Address, enableds []common.Address, initialConfig *commontype.FeeConfig) *FeeManagerConfig {
	return &FeeManagerConfig{
		AllowListConfig: allowlist.AllowListConfig{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
		},
		UpgradeableConfig: config.UpgradeableConfig{BlockTimestamp: blockTimestamp},
		InitialFeeConfig:  initialConfig,
	}
}

// NewDisableFeeManagerConfig returns config for a network upgrade at [blockTimestamp]
// that disables FeeManager.
func NewDisableFeeManagerConfig(blockTimestamp *big.Int) *FeeManagerConfig {
	return &FeeManagerConfig{
		UpgradeableConfig: config.UpgradeableConfig{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

func (FeeManagerConfig) Key() string             { return ConfigKey }
func (FeeManagerConfig) Address() common.Address { return ContractAddress }

// Equal returns true if [cfg] is a [*FeeManagerConfig] and it has been configured identical to [c].
func (c *FeeManagerConfig) Equal(cfg config.Config) bool {
	// typecast before comparison
	other, ok := (cfg).(*FeeManagerConfig)
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

func (c *FeeManagerConfig) Verify() error {
	if err := c.AllowListConfig.Verify(); err != nil {
		return err
	}
	if c.InitialFeeConfig == nil {
		return nil
	}

	return c.InitialFeeConfig.Verify()
}
