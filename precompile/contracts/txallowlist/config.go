// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txallowlist

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ethereum/go-ethereum/common"
)

var _ config.StatefulPrecompileConfig = &Config{}

// Config implements the StatefulPrecompileConfig interface while adding in the
// TxAllowList specific precompile config.
type Config struct {
	allowlist.AllowListConfig
	config.Upgrade
}

// NewConfig returns a config for a network upgrade at [blockTimestamp] that enables
// TxAllowList with the given [admins] and [enableds] as members of the allowlist.
func NewConfig(blockTimestamp *big.Int, admins []common.Address, enableds []common.Address) *Config {
	return &Config{
		AllowListConfig: allowlist.AllowListConfig{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
		},
		Upgrade: config.Upgrade{BlockTimestamp: blockTimestamp},
	}
}

// NewDisableConfig returns config for a network upgrade at [blockTimestamp]
// that disables TxAllowList.
func NewDisableConfig(blockTimestamp *big.Int) *Config {
	return &Config{
		Upgrade: config.Upgrade{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

func (c *Config) Key() string {
	return ConfigKey
}

// Equal returns true if [cfg] is a [*TxAllowListConfig] and it has been configured identical to [c].
func (c *Config) Equal(cfg config.StatefulPrecompileConfig) bool {
	// typecast before comparison
	other, ok := (cfg).(*Config)
	if !ok {
		return false
	}
	return c.Upgrade.Equal(&other.Upgrade) && c.AllowListConfig.Equal(&other.AllowListConfig)
}
