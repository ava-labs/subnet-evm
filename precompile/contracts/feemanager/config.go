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

var _ config.Config = &Config{}

// Config wraps [AllowListConfig] and uses it to implement the StatefulPrecompileConfig
// interface while adding in the FeeManager specific precompile address.
type Config struct {
	allowlist.Config // Config for the fee config manager allow list
	config.Uprade
	InitialFeeConfig *commontype.FeeConfig `json:"initialFeeConfig,omitempty"` // initial fee config to be immediately activated
}

// NewConfig returns a config for a network upgrade at [blockTimestamp] that enables
// FeeManager with the given [admins] and [enableds] as members of the allowlist with [initialConfig] as initial fee config if specified.
func NewConfig(blockTimestamp *big.Int, admins []common.Address, enableds []common.Address, initialConfig *commontype.FeeConfig) *Config {
	return &Config{
		Config: allowlist.Config{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
		},
		Uprade:           config.Uprade{BlockTimestamp: blockTimestamp},
		InitialFeeConfig: initialConfig,
	}
}

// NewDisableConfig returns config for a network upgrade at [blockTimestamp]
// that disables FeeManager.
func NewDisableConfig(blockTimestamp *big.Int) *Config {
	return &Config{
		Uprade: config.Uprade{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

func (Config) Key() string { return ConfigKey }

// Equal returns true if [cfg] is a [*FeeManagerConfig] and it has been configured identical to [c].
func (c *Config) Equal(cfg config.Config) bool {
	// typecast before comparison
	other, ok := (cfg).(*Config)
	if !ok {
		return false
	}
	eq := c.Uprade.Equal(&other.Uprade) && c.Config.Equal(&other.Config)
	if !eq {
		return false
	}

	if c.InitialFeeConfig == nil {
		return other.InitialFeeConfig == nil
	}

	return c.InitialFeeConfig.Equal(other.InitialFeeConfig)
}

func (c *Config) Verify() error {
	if err := c.Config.Verify(); err != nil {
		return err
	}
	if c.InitialFeeConfig == nil {
		return nil
	}

	return c.InitialFeeConfig.Verify()
}
