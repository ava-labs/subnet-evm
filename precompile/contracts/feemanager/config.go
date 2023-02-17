// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package feemanager

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	precompileConfig "github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ethereum/go-ethereum/common"
)

var _ precompileConfig.Config = &Config{}

// Config implements the StatefulPrecompileConfig interface while adding in the
// FeeManager specific precompile config.
type Config struct {
	allowlist.AllowList // Config for the fee config manager allow list
	precompileConfig.Upgrade
	InitialFeeConfig *commontype.FeeConfig `json:"initialFeeConfig,omitempty"` // initial fee config to be immediately activated
}

// NewConfig returns a config for a network upgrade at [blockTimestamp] that enables
// FeeManager with the given [admins] and [enableds] as members of the
// allowlist with [initialConfig] as initial fee config if specified.
func NewConfig(blockTimestamp *big.Int, admins []common.Address, enableds []common.Address, initialConfig *commontype.FeeConfig) *Config {
	return &Config{
		AllowList: allowlist.AllowList{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
		},
		Upgrade:          precompileConfig.Upgrade{BlockTimestamp: blockTimestamp},
		InitialFeeConfig: initialConfig,
	}
}

// NewDisableConfig returns config for a network upgrade at [blockTimestamp]
// that disables FeeManager.
func NewDisableConfig(blockTimestamp *big.Int) *Config {
	return &Config{
		Upgrade: precompileConfig.Upgrade{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

func (*Config) Key() string { return ConfigKey }

// Equal returns true if [cfg] is a [*FeeManagerConfig] and it has been configured identical to [c].
func (c *Config) Equal(cfg precompileConfig.Config) bool {
	// typecast before comparison
	other, ok := (cfg).(*Config)
	if !ok {
		return false
	}
	eq := c.Upgrade.Equal(&other.Upgrade) && c.AllowList.Equal(&other.AllowList)
	if !eq {
		return false
	}

	if c.InitialFeeConfig == nil {
		return other.InitialFeeConfig == nil
	}

	return c.InitialFeeConfig.Equal(other.InitialFeeConfig)
}

func (c *Config) Verify() error {
	if err := c.AllowList.Verify(); err != nil {
		return err
	}
	if c.InitialFeeConfig == nil {
		return nil
	}

	return c.InitialFeeConfig.Verify()
}
