// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nativeminter

import (
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

var _ precompile.StatefulPrecompileConfig = &ContractNativeMinterConfig{}

// ConfigKey is the key used in json config files to specify this precompile config.
// must be unique across all precompiles.
const ConfigKey = "contractNativeMinterConfig"

// ContractNativeMinterConfig wraps [AllowListConfig] and uses it to implement the StatefulPrecompileConfig
// interface while adding in the ContractNativeMinter specific precompile address.
type ContractNativeMinterConfig struct {
	allowlist.AllowListConfig
	precompile.UpgradeableConfig
	InitialMint map[common.Address]*math.HexOrDecimal256 `json:"initialMint,omitempty"` // initial mint config to be immediately minted
}

// NewContractNativeMinterConfig returns a config for a network upgrade at [blockTimestamp] that enables
// ContractNativeMinter with the given [admins] and [enableds] as members of the allowlist. Also mints balances according to [initialMint] when the upgrade activates.
func NewContractNativeMinterConfig(blockTimestamp *big.Int, admins []common.Address, enableds []common.Address, initialMint map[common.Address]*math.HexOrDecimal256) *ContractNativeMinterConfig {
	return &ContractNativeMinterConfig{
		AllowListConfig: allowlist.AllowListConfig{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
		},
		UpgradeableConfig: precompile.UpgradeableConfig{BlockTimestamp: blockTimestamp},
		InitialMint:       initialMint,
	}
}

// NewDisableContractNativeMinterConfig returns config for a network upgrade at [blockTimestamp]
// that disables ContractNativeMinter.
func NewDisableContractNativeMinterConfig(blockTimestamp *big.Int) *ContractNativeMinterConfig {
	return &ContractNativeMinterConfig{
		UpgradeableConfig: precompile.UpgradeableConfig{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

// Address returns the address of the native minter contract.
func (c ContractNativeMinterConfig) Address() common.Address {
	return ContractAddress
}

// Configure configures [state] with the desired admins based on [c].
func (c *ContractNativeMinterConfig) Configure(_ precompile.ChainConfig, state precompile.StateDB, _ precompile.BlockContext) error {
	for to, amount := range c.InitialMint {
		if amount != nil {
			bigIntAmount := (*big.Int)(amount)
			state.AddBalance(to, bigIntAmount)
		}
	}

	return c.AllowListConfig.Configure(state, ContractAddress)
}

// Contract returns the singleton stateful precompiled contract to be used for the native minter.
func (ContractNativeMinterConfig) Contract() precompile.StatefulPrecompiledContract {
	return ContractNativeMinterPrecompile
}

func (c *ContractNativeMinterConfig) Verify() error {
	if err := c.AllowListConfig.Verify(); err != nil {
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

// Equal returns true if [s] is a [*ContractNativeMinterConfig] and it has been configured identical to [c].
func (c *ContractNativeMinterConfig) Equal(s precompile.StatefulPrecompileConfig) bool {
	// typecast before comparison
	other, ok := (s).(*ContractNativeMinterConfig)
	if !ok {
		return false
	}
	eq := c.UpgradeableConfig.Equal(&other.UpgradeableConfig) && c.AllowListConfig.Equal(&other.AllowListConfig)
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

func (ContractNativeMinterConfig) Key() string {
	return ConfigKey
}

func (ContractNativeMinterConfig) New() precompile.StatefulPrecompileConfig {
	return new(ContractNativeMinterConfig)
}
