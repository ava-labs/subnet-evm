// (c) 2022 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/utils"
)

// UpgradesConfig includes a list of network upgrades.
// Upgrades must be sorted in increasing order of BlockTimestamp.
type UpgradesConfig struct {
	Upgrades []upgrade `json:"precompileUpgrades,omitempty"`

	// support top-level stateful precompile configs by embedding
	// an Upgrade struct here. this style of configuration is
	// deprecated and is for backwards-compatibility.
	upgrade
}

// upgrade is a helper struct embedded in UpgradesConfig, representing
// each of the possible stateful precompile types that can be activated
// through UpgradesConfig.
type upgrade struct {
	ContractDeployerAllowListConfig *ContractDeployerAllowListConfig `json:"contractDeployerAllowListConfig,omitempty"` // Config for the contract deployer allow list precompile
	ContractNativeMinterConfig      *ContractNativeMinterConfig      `json:"contractNativeMinterConfig,omitempty"`      // Config for the native minter precompile
	TxAllowListConfig               *TxAllowListConfig               `json:"txAllowListConfig,omitempty"`               // Config for the tx allow list precompile
	FeeManagerConfig                *FeeConfigManagerConfig          `json:"feeManagerConfig,omitempty"`                // Config for the fee manager precompile
}

// TODO: Validate the config

// GetContractDeployerAllowListConfig returns the latest forked ContractDeployerAllowListConfig
// specified by [c] or nil if it was never enabled.
func (c *UpgradesConfig) GetContractDeployerAllowListConfig(blockTimestamp *big.Int) *ContractDeployerAllowListConfig {
	// loop backwards on all upgrades
	for i := len(c.Upgrades) - 1; i >= 0; i-- {
		upgrade := c.Upgrades[i]
		if config := upgrade.ContractDeployerAllowListConfig; config != nil {
			if !utils.IsForked(config.BlockTimestamp, blockTimestamp) {
				// this fork has not happened yet
				continue
			}
			return config
		}

	}
	// fallback to the embedded [Upgrade] for backwards compatibility.
	if config := c.upgrade.ContractDeployerAllowListConfig; config != nil {
		if utils.IsForked(config.BlockTimestamp, blockTimestamp) {
			return config
		}
	}
	return nil
}

// GetContractNativeMinterConfig returns the latest forked ContractNativeMinterConfig
// specified by [c] or nil if it was never enabled.
func (c *UpgradesConfig) GetContractNativeMinterConfig(blockTimestamp *big.Int) *ContractNativeMinterConfig {
	// loop backwards on all upgrades
	for i := len(c.Upgrades) - 1; i >= 0; i-- {
		upgrade := c.Upgrades[i]
		if config := upgrade.ContractNativeMinterConfig; config != nil {
			if !utils.IsForked(config.BlockTimestamp, blockTimestamp) {
				// this fork has not happened yet
				continue
			}
			return config
		}

	}
	// fallback to the embedded [Upgrade] for backwards compatibility.
	if config := c.upgrade.ContractNativeMinterConfig; config != nil {
		if utils.IsForked(config.BlockTimestamp, blockTimestamp) {
			return config
		}
	}
	return nil
}

// GetTxAllowListConfig returns the latest forked TxAllowListConfig
// specified by [c] or nil if it was never enabled.
func (c *UpgradesConfig) GetTxAllowListConfig(blockTimestamp *big.Int) *TxAllowListConfig {
	// loop backwards on all upgrades
	for i := len(c.Upgrades) - 1; i >= 0; i-- {
		upgrade := c.Upgrades[i]
		if config := upgrade.TxAllowListConfig; config != nil {
			if !utils.IsForked(config.BlockTimestamp, blockTimestamp) {
				// this fork has not happened yet
				continue
			}
			return config
		}

	}
	// fallback to the embedded [Upgrade] for backwards compatibility.
	if config := c.upgrade.TxAllowListConfig; config != nil {
		if utils.IsForked(config.BlockTimestamp, blockTimestamp) {
			return config
		}
	}
	return nil
}

// GetFeeConfigManagerConfig returns the latest forked FeeManagerConfig
// specified by [c] or nil if it was never enabled.
func (c *UpgradesConfig) GetFeeConfigManagerConfig(blockTimestamp *big.Int) *FeeConfigManagerConfig {
	// loop backwards on all upgrades
	for i := len(c.Upgrades) - 1; i >= 0; i-- {
		upgrade := c.Upgrades[i]
		if config := upgrade.FeeManagerConfig; config != nil {
			if !utils.IsForked(config.BlockTimestamp, blockTimestamp) {
				// this fork has not happened yet
				continue
			}
			return config
		}

	}
	// fallback to the embedded [Upgrade] for backwards compatibility.
	if config := c.upgrade.FeeManagerConfig; config != nil {
		if utils.IsForked(config.BlockTimestamp, blockTimestamp) {
			return config
		}
	}
	return nil
}

// CheckCompatible checks if [newcfg] is compatible with [c] at [headTimestamp].
// Returns a ConfigCompatError if upgrades that have forked at [headTimestamp]
// are missing from [newcfg]. Upgrades that have not forked yet may be modified
// or absent from [newcfg]. Returns nil if [newcfg] is compatible with [c].
func (c *UpgradesConfig) CheckCompatible(newcfg *UpgradesConfig, headTimestamp *big.Int) *utils.ConfigCompatError {
	newUpgrades := newcfg.EnabledStatefulPrecompiles(headTimestamp)
	for i, upgrade := range c.EnabledStatefulPrecompiles(headTimestamp) {
		if len(newUpgrades) <= i {
			// missing upgrade
			return utils.NewCompatError(
				fmt.Sprintf("missing PrecompileUpgradeConfig[%d]", i),
				upgrade.Timestamp(),
				nil,
			)
		}
		// All upgrades that have forked must be identical.
		// TODO: verify config?
		if upgrade.Timestamp().Cmp(newUpgrades[i].Timestamp()) != 0 ||
			upgrade.IsDisabled() != newUpgrades[i].IsDisabled() {
			return utils.NewCompatError(
				fmt.Sprintf("PrecompileUpgradeConfig[%d]", i),
				upgrade.Timestamp(),
				newUpgrades[i].Timestamp(),
			)
		}
	}
	return nil // newcfg is compatible
}

// EnabledStatefulPrecompiles returns a slice of stateful precompile configss that
// have been activated through an upgrade.
func (c *UpgradesConfig) EnabledStatefulPrecompiles(blockTimestamp *big.Int) []StatefulPrecompileConfig {
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
// from [parentTimestamp] to the timestamp set in [blockContext]. If this is the case, it calls [Configure] or
// [Deconfigure] to apply the necessary state transitions for the upgrade.
// This function is called:
// - within genesis setup to configure the starting state for precompiles enabled at genesis,
// - during block processing to update the state before processing the given block.
func (c *UpgradesConfig) CheckConfigure(chainConfig ChainConfig, parentTimestamp *big.Int, blockContext BlockContext, statedb StateDB) {
	blockTimestamp := blockContext.Timestamp()
	for _, config := range c.EnabledStatefulPrecompiles(blockTimestamp) {
		// If [config] this transition activates the upgrade, configure the stateful precompile.
		// (or deconfigure it if it is being disabled.)
		if utils.IsForkTransition(config.Timestamp(), parentTimestamp, blockTimestamp) {
			if config.IsDisabled() {
				Deconfigure(config.Address(), statedb)
			} else {
				Configure(chainConfig, blockContext, config, statedb)
			}
		}
	}
}
