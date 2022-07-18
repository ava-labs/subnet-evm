// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// This file contains methods used to add network upgrades that enable and disable
// stateful precompiles. These methods should only be called from tests.
// Non-test code should configure the Updates struct by JSON using the chain config
// or upgrade bytes instead.

// AddContractDeployerAllowListUpgrade adds a ContractDeployerAllowList upgrade at [blockTimestamp].
func (c *UpgradesConfig) AddContractDeployerAllowListUpgrade(blockTimestamp *big.Int, admins []common.Address) {
	c.Upgrades = append(c.Upgrades, upgrade{
		ContractDeployerAllowListConfig: &ContractDeployerAllowListConfig{
			UpgradeableConfig: UpgradeableConfig{BlockTimestamp: blockTimestamp},
			AllowListConfig:   AllowListConfig{AllowListAdmins: admins},
		},
	})
}

// DisableContractDeployerAllowListUpgrade disables a previously added
// ContractDeployerAllowList upgrade at [blockTimestamp].
func (c *UpgradesConfig) DisableContractDeployerAllowListUpgrade(blockTimestamp *big.Int) {
	c.Upgrades = append(c.Upgrades, upgrade{
		ContractDeployerAllowListConfig: &ContractDeployerAllowListConfig{
			UpgradeableConfig: UpgradeableConfig{
				BlockTimestamp: blockTimestamp,
				Disable:        true,
			},
		},
	})
}

// AddContractNativeMinterUpgrade adds a ContractNativeMinter upgrade at [blockTimestamp].
func (c *UpgradesConfig) AddContractNativeMinterUpgrade(blockTimestamp *big.Int, admins []common.Address) {
	c.Upgrades = append(c.Upgrades, upgrade{
		ContractNativeMinterConfig: &ContractNativeMinterConfig{
			UpgradeableConfig: UpgradeableConfig{BlockTimestamp: blockTimestamp},
			AllowListConfig:   AllowListConfig{AllowListAdmins: admins},
		},
	})
}

// DisableContractNativeMinterUpgrade disables a previously added
// ContractDeployerAllowList upgrade at [blockTimestamp].
func (c *UpgradesConfig) DisableContractNativeMinterUpgrade(blockTimestamp *big.Int) {
	c.Upgrades = append(c.Upgrades, upgrade{
		ContractNativeMinterConfig: &ContractNativeMinterConfig{
			UpgradeableConfig: UpgradeableConfig{
				BlockTimestamp: blockTimestamp,
				Disable:        true,
			},
		},
	})
}

// AddTxAllowListUpgrade adds a TxAllowList upgrade at [blockTimestamp].
func (c *UpgradesConfig) AddTxAllowListUpgrade(blockTimestamp *big.Int, admins []common.Address) {
	c.Upgrades = append(c.Upgrades, upgrade{
		TxAllowListConfig: &TxAllowListConfig{
			UpgradeableConfig: UpgradeableConfig{BlockTimestamp: blockTimestamp},
			AllowListConfig:   AllowListConfig{AllowListAdmins: admins},
		},
	})
}

// DisableTxAllowListUpgrade disables a previously added TxAllowList
// upgrade at [blockTimestamp].
func (c *UpgradesConfig) DisableTxAllowListUpgrade(blockTimestamp *big.Int) {
	c.Upgrades = append(c.Upgrades, upgrade{
		TxAllowListConfig: &TxAllowListConfig{
			UpgradeableConfig: UpgradeableConfig{
				BlockTimestamp: blockTimestamp,
				Disable:        true,
			},
		},
	})
}

// AddFeeManagerUpgrade adds a FeeConfigManager upgrade at [blockTimestamp].
func (c *UpgradesConfig) AddFeeManagerUpgrade(blockTimestamp *big.Int, admins []common.Address) {
	c.Upgrades = append(c.Upgrades, upgrade{
		FeeManagerConfig: &FeeConfigManagerConfig{
			UpgradeableConfig: UpgradeableConfig{BlockTimestamp: blockTimestamp},
			AllowListConfig:   AllowListConfig{AllowListAdmins: admins},
		},
	})
}

// DisableFeeManagerUpgrade disables a previously added FeeConfigManager
// upgrade at [blockTimestamp].
func (c *UpgradesConfig) DisableFeeManagerUpgrade(blockTimestamp *big.Int) {
	c.Upgrades = append(c.Upgrades, upgrade{
		FeeManagerConfig: &FeeConfigManagerConfig{
			UpgradeableConfig: UpgradeableConfig{
				BlockTimestamp: blockTimestamp,
				Disable:        true,
			},
		},
	})
}
