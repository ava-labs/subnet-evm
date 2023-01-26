// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txallowlist

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ethereum/go-ethereum/common"
)

var _ precompile.StatefulPrecompileConfig = &TxAllowListConfig{}

// ConfigKey is the key used in json config files to specify this precompile config.
// must be unique across all precompiles.
const ConfigKey = "txAllowListConfig"

// TxAllowListConfig wraps [AllowListConfig] and uses it to implement the StatefulPrecompileConfig
// interface while adding in the TxAllowList specific precompile address.
type TxAllowListConfig struct {
	allowlist.AllowListConfig
	precompile.UpgradeableConfig
}

// NewModule returns a new module for TxAllowList.
func NewModule() precompile.StatefulPrecompileModule {
	return &TxAllowListConfig{}
}

// NewTxAllowListConfig returns a config for a network upgrade at [blockTimestamp] that enables
// TxAllowList with the given [admins] and [enableds] as members of the allowlist.
func NewTxAllowListConfig(blockTimestamp *big.Int, admins []common.Address, enableds []common.Address) *TxAllowListConfig {
	return &TxAllowListConfig{
		AllowListConfig: allowlist.AllowListConfig{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
		},
		UpgradeableConfig: precompile.UpgradeableConfig{BlockTimestamp: blockTimestamp},
	}
}

// NewDisableTxAllowListConfig returns config for a network upgrade at [blockTimestamp]
// that disables TxAllowList.
func NewDisableTxAllowListConfig(blockTimestamp *big.Int) *TxAllowListConfig {
	return &TxAllowListConfig{
		UpgradeableConfig: precompile.UpgradeableConfig{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

// Address returns the address of the contract deployer allow list.
func (_ *TxAllowListConfig) Address() common.Address {
	return ContractAddress
}

// Contract returns the singleton stateful precompiled contract to be used for the allow list.
func (_ *TxAllowListConfig) Contract() precompile.StatefulPrecompiledContract {
	return TxAllowListPrecompile
}

// Key returns the key used in json config files to specify this precompile config.
func (c *TxAllowListConfig) Key() string {
	return ConfigKey
}

// NewConfig returns a new instance of TxAllowListConfig.
func (_ *TxAllowListConfig) NewConfig() precompile.StatefulPrecompileConfig {
	return new(TxAllowListConfig)
}

// Configure configures [state] with the desired admins based on [c].
func (c *TxAllowListConfig) Configure(_ precompile.ChainConfig, state precompile.StateDB, _ precompile.BlockContext) error {
	return c.AllowListConfig.Configure(state, c.Address())
}

// Equal returns true if [s] is a [*TxAllowListConfig] and it has been configured identical to [c].
func (c *TxAllowListConfig) Equal(s precompile.StatefulPrecompileConfig) bool {
	// typecast before comparison
	other, ok := (s).(*TxAllowListConfig)
	if !ok {
		return false
	}
	return c.UpgradeableConfig.Equal(&other.UpgradeableConfig) && c.AllowListConfig.Equal(&other.AllowListConfig)
}
