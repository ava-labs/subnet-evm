// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nativeminter

import (
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

var _ precompileconfig.Config = &Config{}

// Config implements the StatefulPrecompileConfig interface while adding in the
// ContractNativeMinter specific precompile config.
type Config struct {
	allowlist.AllowListConfig
	precompileconfig.Upgrade
	InitialMint map[common.Address]*math.HexOrDecimal256 `json:"initialMint,omitempty"` // addresses to receive the initial mint mapped to the amount to mint
}

// NewConfig returns a config for a network upgrade at [blockTimestamp] that enables
// ContractNativeMinter with the given [admins] and [enableds] as members of the allowlist. Also mints balances according to [initialMint] when the upgrade activates.
func NewConfig(blockTimestamp *big.Int, admins []common.Address, enableds []common.Address, initialMint map[common.Address]*math.HexOrDecimal256) *Config {
	return &Config{
		AllowListConfig: allowlist.AllowListConfig{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
		},
		Upgrade:     precompileconfig.Upgrade{BlockTimestamp: blockTimestamp},
		InitialMint: initialMint,
	}
}

// NewDisableConfig returns config for a network upgrade at [blockTimestamp]
// that disables ContractNativeMinter.
func NewDisableConfig(blockTimestamp *big.Int) *Config {
	return &Config{
		Upgrade: precompileconfig.Upgrade{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}
func (*Config) Key() string { return ConfigKey }

// Equal returns true if [cfg] is a [*ContractNativeMinterConfig] and it has been configured identical to [c].
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

	if len(c.InitialMint) != len(other.InitialMint) {
		return false
	}

	for address, amount := range c.InitialMint {
		val, ok := other.InitialMint[address]
		if !ok {
			return false
		}
		bigIntAmount := (*big.Int)(amount)
		bigIntVal := (*big.Int)(val)
		if !utils.BigNumEqual(bigIntAmount, bigIntVal) {
			return false
		}
	}

	return true
}

func (c *Config) Verify() error {
	// ensure that all of the initial mint values in the map are non-nil positive values
	for addr, amount := range c.InitialMint {
		if amount == nil {
			return fmt.Errorf("initial mint cannot contain nil amount for address %s", addr)
		}
		bigIntAmount := (*big.Int)(amount)
		if bigIntAmount.Sign() < 1 {
			return fmt.Errorf("initial mint cannot contain invalid amount %v for address %s", bigIntAmount, addr)
		}
	}
	return c.AllowListConfig.Verify()
}
