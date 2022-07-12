// (c) 2022 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/ava-labs/subnet-evm/utils"
)

// UpgradesConfig includes a list of network upgrades.
// Upgrades must be sorted in increasing order of BlockTimestamp.
type UpgradesConfig struct {
	Upgrades []Upgrade `json:"upgrades,omitempty"`
}

// Upgrade specifies a single network upgrade that may
// enable or disable a set of precompiles.
type Upgrade struct {
	BlockTimestamp *big.Int     `json:"blockTimestamp"`
	Enable         *upgradeable `json:"enable,omitempty"`
	Disable        *upgradeable `json:"disable,omitempty"`
}

type upgradeable struct {
	ContractDeployerAllowListConfig *ContractDeployerAllowListConfig `json:"contractDeployerAllowListConfig,omitempty"` // Config for the contract deployer allow list precompile
	ContractNativeMinterConfig      *ContractNativeMinterConfig      `json:"contractNativeMinterConfig,omitempty"`      // Config for the native minter precompile
	TxAllowListConfig               *TxAllowListConfig               `json:"txAllowListConfig,omitempty"`               // Config for the tx allow list precompile
	FeeManagerConfig                *FeeConfigManagerConfig          `json:"feeManagerConfig,omitempty"`                // Config for the fee manager precompile
}

// TODO: Validate the config

// AddContractDeployerAllowListUpgrade adds a ContractDeployerAllowListConfig upgrade at [blockTimestamp].
func (c *UpgradesConfig) AddContractDeployerAllowListUpgrade(blockTimestamp *big.Int, config *ContractDeployerAllowListConfig) {
	c.Upgrades = append(c.Upgrades, Upgrade{
		BlockTimestamp: blockTimestamp,
		Enable: &upgradeable{
			ContractDeployerAllowListConfig: config,
		},
	})
}

// AddContractNativeMinterUpgrade adds a ContractNativeMinterConfig upgrade at [blockTimestamp].
func (c *UpgradesConfig) AddContractNativeMinterUpgrade(blockTimestamp *big.Int, config *ContractNativeMinterConfig) {
	c.Upgrades = append(c.Upgrades, Upgrade{
		BlockTimestamp: blockTimestamp,
		Enable: &upgradeable{
			ContractNativeMinterConfig: config,
		},
	})
}

// AddTxAllowListUpgrade adds a TxAllowListConfig upgrade at [blockTimestamp].
func (c *UpgradesConfig) AddTxAllowListUpgrade(blockTimestamp *big.Int, config *TxAllowListConfig) {
	c.Upgrades = append(c.Upgrades, Upgrade{
		BlockTimestamp: blockTimestamp,
		Enable: &upgradeable{
			TxAllowListConfig: config,
		},
	})
}

// AddFeeManagerUpgrade adds a FeeConfigManagerConfig upgrade at [blockTimestamp].
func (c *UpgradesConfig) AddFeeManagerUpgrade(blockTimestamp *big.Int, config *FeeConfigManagerConfig) {
	c.Upgrades = append(c.Upgrades, Upgrade{
		BlockTimestamp: blockTimestamp,
		Enable: &upgradeable{
			FeeManagerConfig: config,
		},
	})
}

// GetContractDeployerAllowListConfig returns the latest ContractDeployerAllowListConfig specified by [c] or nil if
// ContractDeployerAllowListConfig was disabled or never enabled.
func (c *UpgradesConfig) GetContractDeployerAllowListConfig(blockTimestamp *big.Int) *ContractDeployerAllowListConfig {
	// loop backwards on all upgrades
	for i := len(c.Upgrades) - 1; i >= 0; i-- {
		upgrade := c.Upgrades[i]
		// check disable first
		if upgrade.Disable != nil {
			if upgrade.Disable.ContractDeployerAllowListConfig != nil && utils.IsForked(upgrade.BlockTimestamp, blockTimestamp) {
				return nil
			}
		}

		// then check enables
		if upgrade.Enable != nil {
			if upgrade.Enable.ContractDeployerAllowListConfig != nil && utils.IsForked(upgrade.BlockTimestamp, blockTimestamp) {
				return upgrade.Enable.ContractDeployerAllowListConfig
			}
		}
	}
	return nil
}

// GetContractNativeMinterConfig returns the latest ContractNativeMinterConfig specified by [c] or nil if
// ContractNativeMinterConfig was disabled or never enabled.
func (c *UpgradesConfig) GetContractNativeMinterConfig(blockTimestamp *big.Int) *ContractNativeMinterConfig {
	// loop backwards on all upgrades
	for i := len(c.Upgrades) - 1; i >= 0; i-- {
		upgrade := c.Upgrades[i]
		// check disable first
		if upgrade.Disable != nil {
			if upgrade.Disable.ContractNativeMinterConfig != nil && utils.IsForked(upgrade.BlockTimestamp, blockTimestamp) {
				return nil
			}
		}

		// then check enables
		if upgrade.Enable != nil {
			if upgrade.Enable.ContractNativeMinterConfig != nil && utils.IsForked(upgrade.BlockTimestamp, blockTimestamp) {
				return upgrade.Enable.ContractNativeMinterConfig
			}
		}
	}
	return nil
}

// GetTxAllowListConfig returns the latest TxAllowListConfig specified by [c] or nil if
// TxAllowListConfig was disabled or never enabled.
func (c *UpgradesConfig) GetTxAllowListConfig(blockTimestamp *big.Int) *TxAllowListConfig {
	// loop backwards on all upgrades
	for i := len(c.Upgrades) - 1; i >= 0; i-- {
		upgrade := c.Upgrades[i]
		// check disable first
		if upgrade.Disable != nil {
			if upgrade.Disable.TxAllowListConfig != nil && utils.IsForked(upgrade.BlockTimestamp, blockTimestamp) {
				return nil
			}
		}

		// then check enables
		if upgrade.Enable != nil {
			if upgrade.Enable.TxAllowListConfig != nil && utils.IsForked(upgrade.BlockTimestamp, blockTimestamp) {
				return upgrade.Enable.TxAllowListConfig
			}
		}
	}
	return nil
}

// GetFeeConfigManagerConfig returns the latest FeeManagerConfig specified by [c] or nil if
// FeeManagerConfig was disabled or never enabled.
func (c *UpgradesConfig) GetFeeConfigManagerConfig(blockTimestamp *big.Int) *FeeConfigManagerConfig {
	// loop backwards on all upgrades
	for i := len(c.Upgrades) - 1; i >= 0; i-- {
		upgrade := c.Upgrades[i]
		// check disable first
		if upgrade.Disable != nil {
			if upgrade.Disable.FeeManagerConfig != nil && utils.IsForked(upgrade.BlockTimestamp, blockTimestamp) {
				return nil
			}
		}

		// then check enables
		if upgrade.Enable != nil {
			if upgrade.Enable.FeeManagerConfig != nil && utils.IsForked(upgrade.BlockTimestamp, blockTimestamp) {
				return upgrade.Enable.FeeManagerConfig
			}
		}
	}
	return nil
}

// CheckCompatible checks if [newcfg] is compatible with [c] at [headTimestamp].
// Returns a ConfigCompatError if upgrades that have forked at [headTimestamp]
// are missing from [newcfg]. Upgrades that have not forked yet may be modified
// or absent from [newcfg]. Returns nil if [newcfg] is compatible with [c].
func (c *UpgradesConfig) CheckCompatible(newcfg *UpgradesConfig, headTimestamp *big.Int) *utils.ConfigCompatError {
	for i, upgrade := range c.Upgrades {
		if !utils.IsForked(upgrade.BlockTimestamp, headTimestamp) {
			// we have checked all the forked upgrades, so we can break here
			// to allow modifying upgrades that have not forked yet.
			break
		}
		if len(newcfg.Upgrades) <= i {
			// missing upgrade
			return utils.NewCompatError(
				fmt.Sprintf("missing PrecompileUpgradeConfig[%d]", i),
				upgrade.BlockTimestamp,
				nil,
			)
		}

		// All upgrades that have forked must be identical.
		// TODO: return error w/ details instead of this
		if !reflect.DeepEqual(upgrade, newcfg.Upgrades[i]) {
			return utils.NewCompatError(
				fmt.Sprintf("PrecompileUpgradeConfig[%d]", i),
				upgrade.BlockTimestamp,
				newcfg.Upgrades[i].BlockTimestamp,
			)
		}
	}
	return nil // newcfg is compatible
}

// ActiveStatefulPrecompiles returns a slice of stateful precompiles that have been
// activated through an upgrade but not deactivated yet.
func (c *UpgradesConfig) ActiveStatefulPrecompiles(blockTimestamp *big.Int) []StatefulPrecompileConfig {
	statefulPrecompileConfigs := make([]StatefulPrecompileConfig, 0)
	if config := c.GetContractDeployerAllowListConfig(blockTimestamp); config != nil {
		statefulPrecompileConfigs = append(statefulPrecompileConfigs, config)
	}
	if config := c.GetContractNativeMinterConfig(blockTimestamp); config != nil {
		statefulPrecompileConfigs = append(statefulPrecompileConfigs, config)
	}
	if config := c.GetTxAllowListConfig(blockTimestamp); config != nil {
		statefulPrecompileConfigs = append(statefulPrecompileConfigs, config)
	}
	if config := c.GetFeeConfigManagerConfig(blockTimestamp); config != nil {
		statefulPrecompileConfigs = append(statefulPrecompileConfigs, config)
	}

	return statefulPrecompileConfigs
}

// CheckConfigure checks if any of the precompiles specified by [c] is enabled or disabled by the block transition
// from [parentTimestamp] to [currentTimestamp]. If this is the case, it calls [Configure] or [Deconfigure] to apply
// the necessary state transitions for the upgrade.
// Note: this function is called within genesis to configure the starting state if it [config] specifies that it should be
// configured at genesis, or happens during block processing to update the state before processing the given block.
func (c *UpgradesConfig) CheckConfigure(chainConfig ChainConfig, parentTimestamp *big.Int, blockContext BlockContext, statedb StateDB) {
	for i := len(c.Upgrades) - 1; i >= 0; i-- {
		upgrade := c.Upgrades[i]
		if utils.IsForked(upgrade.BlockTimestamp, parentTimestamp) {
			// all forks already applied
			break
		}

		// If [upgrade] goes into effect within this transition, configure the stateful precompile
		if utils.IsForkTransition(upgrade.BlockTimestamp, parentTimestamp, blockContext.Timestamp()) {
			// handle disables first (in case an upgrade is disabled and enabled in the same fork)
			if upgrade.Disable != nil {
				if upgrade.Disable.ContractDeployerAllowListConfig != nil {
					Deconfigure(upgrade.Disable.ContractDeployerAllowListConfig, statedb)
				}
				if upgrade.Disable.ContractNativeMinterConfig != nil {
					Deconfigure(upgrade.Disable.ContractNativeMinterConfig, statedb)
				}
				if upgrade.Disable.TxAllowListConfig != nil {
					Deconfigure(upgrade.Disable.TxAllowListConfig, statedb)
				}
				if upgrade.Disable.FeeManagerConfig != nil {
					Deconfigure(upgrade.Disable.FeeManagerConfig, statedb)
				}
			}

			// handle upgrades that are enabled
			if upgrade.Enable != nil {
				if upgrade.Enable.ContractDeployerAllowListConfig != nil {
					Configure(chainConfig, blockContext, upgrade.Enable.ContractDeployerAllowListConfig, statedb)
				}
				if upgrade.Enable.ContractNativeMinterConfig != nil {
					Configure(chainConfig, blockContext, upgrade.Enable.ContractNativeMinterConfig, statedb)
				}
				if upgrade.Enable.TxAllowListConfig != nil {
					Configure(chainConfig, blockContext, upgrade.Enable.TxAllowListConfig, statedb)
				}
				if upgrade.Enable.FeeManagerConfig != nil {
					Configure(chainConfig, blockContext, upgrade.Enable.FeeManagerConfig, statedb)
				}
			}
		}
	}
}
