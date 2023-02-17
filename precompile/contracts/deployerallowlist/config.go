// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package deployerallowlist

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	precompileConfig "github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ethereum/go-ethereum/common"
)

var _ precompileConfig.Config = &Config{}

// Config contains the configuration for the ContractDeployerAllowList precompile,
// consisting of the initial allowlist and the timestamp for the network upgrade.
type Config struct {
	allowlist.AllowList
	precompileConfig.Upgrade
}

// NewConfig returns a config for a network upgrade at [blockTimestamp] that enables
// ContractDeployerAllowList with [admins] and [enableds] as members of the allowlist.
func NewConfig(blockTimestamp *big.Int, admins []common.Address, enableds []common.Address) *Config {
	return &Config{
		AllowList: allowlist.AllowList{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
		},
		Upgrade: precompileConfig.Upgrade{BlockTimestamp: blockTimestamp},
	}
}

// NewDisableConfig returns config for a network upgrade at [blockTimestamp]
// that disables ContractDeployerAllowList.
func NewDisableConfig(blockTimestamp *big.Int) *Config {
	return &Config{
		Upgrade: precompileConfig.Upgrade{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

func (*Config) Key() string { return ConfigKey }

// Equal returns true if [cfg] is a [*ContractDeployerAllowListConfig] and it has been configured identical to [c].
func (c *Config) Equal(cfg precompileConfig.Config) bool {
	// typecast before comparison
	other, ok := (cfg).(*Config)
	if !ok {
		return false
	}
	return c.Upgrade.Equal(&other.Upgrade) && c.AllowList.Equal(&other.AllowList)
}

func (c *Config) Verify() error { return c.AllowList.Verify() }
