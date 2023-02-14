// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nativeminter

import (
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

var _ config.Config = &ContractNativeMinterConfig{}

// ContractNativeMinterConfig wraps [AllowListConfig] and uses it to implement the StatefulPrecompileConfig
// interface while adding in the ContractNativeMinter specific precompile address.
type ContractNativeMinterConfig struct {
	allowlist.Config
	config.Uprade
	InitialMint map[common.Address]*math.HexOrDecimal256 `json:"initialMint,omitempty"` // initial mint config to be immediately minted
}

// NewContractNativeMinterConfig returns a config for a network upgrade at [blockTimestamp] that enables
// ContractNativeMinter with the given [admins] and [enableds] as members of the allowlist. Also mints balances according to [initialMint] when the upgrade activates.
func NewContractNativeMinterConfig(blockTimestamp *big.Int, admins []common.Address, enableds []common.Address, initialMint map[common.Address]*math.HexOrDecimal256) *ContractNativeMinterConfig {
	return &ContractNativeMinterConfig{
		Config: allowlist.Config{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
		},
		Uprade:      config.Uprade{BlockTimestamp: blockTimestamp},
		InitialMint: initialMint,
	}
}

// NewDisableContractNativeMinterConfig returns config for a network upgrade at [blockTimestamp]
// that disables ContractNativeMinter.
func NewDisableContractNativeMinterConfig(blockTimestamp *big.Int) *ContractNativeMinterConfig {
	return &ContractNativeMinterConfig{
		Uprade: config.Uprade{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}
func (ContractNativeMinterConfig) Key() string { return ConfigKey }

// Equal returns true if [cfg] is a [*ContractNativeMinterConfig] and it has been configured identical to [c].
func (c *ContractNativeMinterConfig) Equal(cfg config.Config) bool {
	// typecast before comparison
	other, ok := (cfg).(*ContractNativeMinterConfig)
	if !ok {
		return false
	}
	eq := c.Uprade.Equal(&other.Uprade) && c.Config.Equal(&other.Config)
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

func (c *ContractNativeMinterConfig) Verify() error {
	if err := c.Config.Verify(); err != nil {
		return err
	}
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
	return nil
}
