// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package deployerallowlist

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ethereum/go-ethereum/common"
)

var _ config.Config = &Config{}

// Config contains the configuration for the ContractDeployerAllowList precompile,
// consisting of the initial allowlist and the timestamp for the network upgrade.
type Config struct {
	allowlist.Config
	config.Upgrade
}

// NewConfig returns a config for a network upgrade at [blockTimestamp] that enables
// ContractDeployerAllowList with [admins] and [enableds] as members of the allowlist.
func NewConfig(blockTimestamp *big.Int, admins []common.Address, enableds []common.Address) *Config {
	return &Config{
		Config: allowlist.Config{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
		},
		Upgrade: config.Upgrade{BlockTimestamp: blockTimestamp},
	}
}

// NewDisableConfig returns config for a network upgrade at [blockTimestamp]
// that disables ContractDeployerAllowList.
func NewDisableConfig(blockTimestamp *big.Int) *Config {
	return &Config{
		Upgrade: config.Upgrade{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

func (Config) Key() string { return ConfigKey }

// Equal returns true if [cfg] is a [*ContractDeployerAllowListConfig] and it has been configured identical to [c].
func (c *Config) Equal(cfg config.Config) bool {
	// typecast before comparison
	other, ok := (cfg).(*Config)
	if !ok {
		return false
	}
	return c.Upgrade.Equal(&other.Upgrade) && c.Config.Equal(&other.Config)
}

func (c *Config) Verify() error { return c.Config.Verify() }
