// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txallowlist

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ethereum/go-ethereum/common"
)

var _ config.Config = &TxAllowListConfig{}

// TxAllowListConfig wraps [AllowListConfig] and uses it to implement the StatefulPrecompileConfig
// interface while adding in the TxAllowList specific precompile address.
type TxAllowListConfig struct {
	allowlist.Config
	config.Uprade
}

// NewTxAllowListConfig returns a config for a network upgrade at [blockTimestamp] that enables
// TxAllowList with the given [admins] and [enableds] as members of the allowlist.
func NewTxAllowListConfig(blockTimestamp *big.Int, admins []common.Address, enableds []common.Address) *TxAllowListConfig {
	return &TxAllowListConfig{
		Config: allowlist.Config{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
		},
		Uprade: config.Uprade{BlockTimestamp: blockTimestamp},
	}
}

// NewDisableTxAllowListConfig returns config for a network upgrade at [blockTimestamp]
// that disables TxAllowList.
func NewDisableTxAllowListConfig(blockTimestamp *big.Int) *TxAllowListConfig {
	return &TxAllowListConfig{
		Uprade: config.Uprade{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

func (c *TxAllowListConfig) Key() string { return ConfigKey }

// Equal returns true if [cfg] is a [*TxAllowListConfig] and it has been configured identical to [c].
func (c *TxAllowListConfig) Equal(cfg config.Config) bool {
	// typecast before comparison
	other, ok := (cfg).(*TxAllowListConfig)
	if !ok {
		return false
	}
	return c.Uprade.Equal(&other.Uprade) && c.Config.Equal(&other.Config)
}
