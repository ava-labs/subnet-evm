// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package deployerallowlist

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ethereum/go-ethereum/common"
)

var _ config.Config = &ContractDeployerAllowListConfig{}

// ContractDeployerAllowListConfig wraps [AllowListConfig] and uses it to implement the StatefulPrecompileConfig
// interface while adding in the contract deployer specific precompile address.
type ContractDeployerAllowListConfig struct {
	allowlist.Config
	config.Uprade
}

// NewContractDeployerAllowListConfig returns a config for a network upgrade at [blockTimestamp] that enables
// ContractDeployerAllowList with [admins] and [enableds] as members of the allowlist.
func NewContractDeployerAllowListConfig(blockTimestamp *big.Int, admins []common.Address, enableds []common.Address) *ContractDeployerAllowListConfig {
	return &ContractDeployerAllowListConfig{
		Config: allowlist.Config{
			AdminAddresses:   admins,
			EnabledAddresses: enableds,
		},
		Uprade: config.Uprade{BlockTimestamp: blockTimestamp},
	}
}

// NewDisableContractDeployerAllowListConfig returns config for a network upgrade at [blockTimestamp]
// that disables ContractDeployerAllowList.
func NewDisableContractDeployerAllowListConfig(blockTimestamp *big.Int) *ContractDeployerAllowListConfig {
	return &ContractDeployerAllowListConfig{
		Uprade: config.Uprade{
			BlockTimestamp: blockTimestamp,
			Disable:        true,
		},
	}
}

func (ContractDeployerAllowListConfig) Key() string { return ConfigKey }

// Equal returns true if [cfg] is a [*ContractDeployerAllowListConfig] and it has been configured identical to [c].
func (c *ContractDeployerAllowListConfig) Equal(cfg config.Config) bool {
	// typecast before comparison
	other, ok := (cfg).(*ContractDeployerAllowListConfig)
	if !ok {
		return false
	}
	return c.Uprade.Equal(&other.Uprade) && c.Config.Equal(&other.Config)
}

func (c *ContractDeployerAllowListConfig) Verify() error { return c.Config.Verify() }
