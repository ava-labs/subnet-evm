// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package feemanager

import (
	"github.com/ava-labs/avalanchego/utils/wrappers"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/docker/docker/pkg/units"
	"github.com/ethereum/go-ethereum/common"
)

var _ precompileconfig.Config = &Config{}

// Config implements the StatefulPrecompileConfig interface while adding in the
// FeeManager specific precompile config.
type Config struct {
	allowlist.AllowListConfig
	precompileconfig.Upgrade
	InitialFeeConfig *commontype.FeeConfig `json:"initialFeeConfig,omitempty" ` // initial fee config to be immediately activated
}

// NewConfig returns a config for a network upgrade at [blockTimestamp] that enables
// FeeManager with the given [admins], [enableds] and [managers] as members of the
// allowlist with [initialConfig] as initial fee config if specified.
func NewConfig(blockTimestamp *uint64, admins []common.Address, enableds []common.Address, managers []common.Address, initialConfig *commontype.FeeConfig) *Config {
	return &Config{
		AllowListConfig: allowlist.AllowListConfig{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
			ManagerAddresses: managers,
		},
		Upgrade:          precompileconfig.Upgrade{BlockTimestamp: blockTimestamp},
		InitialFeeConfig: initialConfig,
	}
}

// NewDisableConfig returns config for a network upgrade at [blockTimestamp]
// that disables FeeManager.
func NewDisableConfig(blockTimestamp *uint64) *Config {
	return &Config{
		Upgrade: precompileconfig.Upgrade{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

func (*Config) Key() string { return ConfigKey }

// Equal returns true if [cfg] is a [*FeeManagerConfig] and it has been configured identical to [c].
func (c *Config) Equal(cfg precompileconfig.Config) bool {
	// typecast before comparison
	other, ok := (cfg).(*Config)
	if !ok {
		return false
	}
	eq := c.Upgrade.Equal(&other.Upgrade) && c.AllowListConfig.Equal(&other.AllowListConfig)
	if !eq {
		return false
	}

	if c.InitialFeeConfig == nil {
		return other.InitialFeeConfig == nil
	}

	return c.InitialFeeConfig.Equal(other.InitialFeeConfig)
}

func (c *Config) Verify(chainConfig precompileconfig.ChainConfig) error {
	if err := c.AllowListConfig.Verify(chainConfig, c.Upgrade); err != nil {
		return err
	}
	if c.InitialFeeConfig == nil {
		return nil
	}

	return c.InitialFeeConfig.Verify()
}

func (c *Config) MarshalBinary() ([]byte, error) {
	p := wrappers.Packer{
		Bytes:   []byte{},
		MaxSize: 1 * units.MiB,
	}

	allowList, err := c.AllowListConfig.MarshalBinary()
	if err != nil {
		return nil, err
	}

	upgrade, err := c.Upgrade.MarshalBinary()
	if err != nil {
		return nil, err
	}

	p.PackBytes(allowList)
	if p.Err != nil {
		return nil, p.Err
	}
	p.PackBytes(upgrade)
	if p.Err != nil {
		return nil, p.Err
	}

	if c.InitialFeeConfig == nil {
		p.PackBool(true)
		if p.Err != nil {
			return nil, p.Err
		}
	} else {
		p.PackBool(false)
		if p.Err != nil {
			return nil, p.Err
		}
		bytes, err := c.InitialFeeConfig.MarshalBinary()
		if err != nil {
			return nil, err
		}
		p.PackBytes(bytes)
	}

	return p.Bytes, p.Err
}

func (c *Config) UnmarshalBinary(bytes []byte) error {
	p := wrappers.Packer{
		Bytes: bytes,
	}
	allowList := p.UnpackBytes()
	if p.Err != nil {
		return p.Err
	}
	upgrade := p.UnpackBytes()
	if p.Err != nil {
		return p.Err
	}
	if err := c.AllowListConfig.UnmarshalBinary(allowList); err != nil {
		return err
	}
	if err := c.Upgrade.UnmarshalBinary(upgrade); err != nil {
		return err
	}
	isNil := p.UnpackBool()
	if !isNil {
		c.InitialFeeConfig = &commontype.FeeConfig{}
		bytes := p.UnpackBytes()
		if p.Err != nil {
			return p.Err
		}
		if err := c.InitialFeeConfig.UnmarshalBinary(bytes); err != nil {
			return err
		}
	}

	return nil
}
