// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txallowlist

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	precompileConfig "github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ethereum/go-ethereum/common"
)

var _ precompileConfig.Config = &Config{}

// Config implements the StatefulPrecompileConfig interface while adding in the
// TxAllowList specific precompile config.
type Config struct {
	allowlist.AllowList
	precompileConfig.Upgrade
}

// NewConfig returns a config for a network upgrade at [blockTimestamp] that enables
// TxAllowList with the given [admins] and [enableds] as members of the allowlist.
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
// that disables TxAllowList.
func NewDisableConfig(blockTimestamp *big.Int) *Config {
	return &Config{
		Upgrade: precompileConfig.Upgrade{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

func (c *Config) Key() string { return ConfigKey }

// Equal returns true if [cfg] is a [*TxAllowListConfig] and it has been configured identical to [c].
func (c *Config) Equal(cfg precompileConfig.Config) bool {
	// typecast before comparison
	other, ok := (cfg).(*Config)
	if !ok {
		return false
	}
	return c.Upgrade.Equal(&other.Upgrade) && c.AllowList.Equal(&other.AllowList)
}
