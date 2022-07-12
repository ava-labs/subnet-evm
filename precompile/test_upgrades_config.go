// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import "math/big"

// This file contains methods used to add network upgrades that enable and disable
// stateful precompiles. These methods should only be called from tests.
// Non-test code should configure the Updates struct by JSON using the chain config
// or upgrade bytes instead.

// AddContractDeployerAllowListUpgrade adds a ContractDeployerAllowListConfig
// upgrade at [blockTimestamp].
func (c *UpgradesConfig) AddContractDeployerAllowListUpgrade(blockTimestamp *big.Int, config *ContractDeployerAllowListConfig) {
	c.Upgrades = append(c.Upgrades, Upgrade{
		BlockTimestamp: blockTimestamp,
		enable: enable{
			ContractDeployerAllowListConfig: config,
		},
	})
}

// AddContractNativeMinterUpgrade adds a ContractNativeMinterConfig upgrade at
// [blockTimestamp].
func (c *UpgradesConfig) AddContractNativeMinterUpgrade(blockTimestamp *big.Int, config *ContractNativeMinterConfig) {
	c.Upgrades = append(c.Upgrades, Upgrade{
		BlockTimestamp: blockTimestamp,
		enable: enable{
			ContractNativeMinterConfig: config,
		},
	})
}

// AddTxAllowListUpgrade adds a TxAllowListConfig upgrade at [blockTimestamp].
func (c *UpgradesConfig) AddTxAllowListUpgrade(blockTimestamp *big.Int, config *TxAllowListConfig) {
	c.Upgrades = append(c.Upgrades, Upgrade{
		BlockTimestamp: blockTimestamp,
		enable: enable{
			TxAllowListConfig: config,
		},
	})
}

// AddFeeManagerUpgrade adds a FeeConfigManagerConfig upgrade at [blockTimestamp].
func (c *UpgradesConfig) AddFeeManagerUpgrade(blockTimestamp *big.Int, config *FeeConfigManagerConfig) {
	c.Upgrades = append(c.Upgrades, Upgrade{
		BlockTimestamp: blockTimestamp,
		enable: enable{
			FeeManagerConfig: config,
		},
	})
}

// DisableContractDeployerAllowListUpgrade disables a previously added
// ContractDeployerAllowListConfig upgrade at [blockTimestamp].
func (c *UpgradesConfig) DisableContractDeployerAllowListUpgrade(blockTimestamp *big.Int) {
	c.Upgrades = append(c.Upgrades, Upgrade{
		BlockTimestamp: blockTimestamp,
		disable: disable{
			DisableTxAllowList: &struct{}{},
		},
	})
}

// DisableContractNativeMinterUpgrade disables a previously added
// ContractDeployerAllowListConfig upgrade at [blockTimestamp].
func (c *UpgradesConfig) DisableContractNativeMinterUpgrade(blockTimestamp *big.Int) {
	c.Upgrades = append(c.Upgrades, Upgrade{
		BlockTimestamp: blockTimestamp,
		disable: disable{
			DisableTxAllowList: &struct{}{},
		},
	})
}

// DisableTxAllowListUpgrade disables a previously added TxAllowListConfig
// upgrade at [blockTimestamp].
func (c *UpgradesConfig) DisableTxAllowListUpgrade(blockTimestamp *big.Int) {
	c.Upgrades = append(c.Upgrades, Upgrade{
		BlockTimestamp: blockTimestamp,
		disable: disable{
			DisableTxAllowList: &struct{}{},
		},
	})
}

// DisableFeeManagerUpgrade disables a previously added FeeConfigManagerConfig
// upgrade at [blockTimestamp].
func (c *UpgradesConfig) DisableFeeManagerUpgrade(blockTimestamp *big.Int) {
	c.Upgrades = append(c.Upgrades, Upgrade{
		BlockTimestamp: blockTimestamp,
		disable: disable{
			DisableTxAllowList: &struct{}{},
		},
	})
}
