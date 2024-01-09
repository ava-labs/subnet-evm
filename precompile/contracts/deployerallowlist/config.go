// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package deployerallowlist

import (
	"github.com/ava-labs/avalanchego/utils/wrappers"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/docker/docker/pkg/units"
	"github.com/ethereum/go-ethereum/common"
)

var _ precompileconfig.Config = &Config{}

// Config contains the configuration for the ContractDeployerAllowList precompile,
// consisting of the initial allowlist and the timestamp for the network upgrade.
type Config struct {
	allowlist.AllowListConfig
	precompileconfig.Upgrade
}

// NewConfig returns a config for a network upgrade at [blockTimestamp] that enables
// ContractDeployerAllowList with [admins], [enableds] and [managers] as members of the allowlist.
func NewConfig(blockTimestamp *uint64, admins []common.Address, enableds []common.Address, managers []common.Address) *Config {
	return &Config{
		AllowListConfig: allowlist.AllowListConfig{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
			ManagerAddresses: managers,
		},
		Upgrade: precompileconfig.Upgrade{BlockTimestamp: blockTimestamp},
	}
}

// NewDisableConfig returns config for a network upgrade at [blockTimestamp]
// that disables ContractDeployerAllowList.
func NewDisableConfig(blockTimestamp *uint64) *Config {
	return &Config{
		Upgrade: precompileconfig.Upgrade{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

func (*Config) Key() string { return ConfigKey }

// Equal returns true if [cfg] is a [*ContractDeployerAllowListConfig] and it has been configured identical to [c].
func (c *Config) Equal(cfg precompileconfig.Config) bool {
	// typecast before comparison
	other, ok := (cfg).(*Config)
	if !ok {
		return false
	}
	return c.Upgrade.Equal(&other.Upgrade) && c.AllowListConfig.Equal(&other.AllowListConfig)
}

func (c *Config) Verify(chainConfig precompileconfig.ChainConfig) error {
	return c.AllowListConfig.Verify(chainConfig, c.Upgrade)
}

func (c *Config) MarshalBinary() ([]byte, error) {
	p := wrappers.Packer{
		Bytes:   []byte{},
		MaxSize: 1 * units.MiB,
	}

	bytes, err := c.AllowListConfig.MarshalBinary()
	if err != nil {
		return nil, err
	}

	p.PackBytes(bytes)
	if p.Err != nil {
		return nil, p.Err
	}

	bytes, err = c.Upgrade.MarshalBinary()
	if err != nil {
		return nil, err
	}

	p.PackBytes(bytes)
	if p.Err != nil {
		return nil, p.Err
	}

	return p.Bytes, nil
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
	return nil
}
