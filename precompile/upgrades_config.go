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
	// upgrade is embedded here to deserialize precompiles configured
	// in the genesis chain config.
	upgrade

	Upgrades []upgrade
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

type getterFn func(*upgrade) StatefulPrecompileConfig

// getters contains a list of functions that return each of the stateful precompile
// types that can be retrieved from an upgrade.
var getters = []getterFn{
	func(u *upgrade) StatefulPrecompileConfig { return u.ContractDeployerAllowListConfig },
	func(u *upgrade) StatefulPrecompileConfig { return u.ContractNativeMinterConfig },
	func(u *upgrade) StatefulPrecompileConfig { return u.TxAllowListConfig },
	func(u *upgrade) StatefulPrecompileConfig { return u.FeeManagerConfig },
}

// TODO: Validate the config

// getActiveConfig returns the most recent config that has
// already forked, or nil if none have been configured or
// if none have forked yet.
func (c *UpgradesConfig) getActiveConfig(blockTimestamp *big.Int, getter getterFn) StatefulPrecompileConfig {
	configs := c.getActiveConfigs(nil, blockTimestamp, getter)
	if len(configs) == 0 {
		return nil
	}
	return configs[len(configs)-1] // return the most recent config
}

// getActiveConfigs returns all forks configured to activate in the
// specified [from, to] range, in the order of activation.
func (c *UpgradesConfig) getActiveConfigs(from *big.Int, to *big.Int, getter getterFn) []StatefulPrecompileConfig {
	configs := make([]StatefulPrecompileConfig, 0)
	// first check the embedded [upgrade] for precompiles configured
	// in the genesis chain config.
	if config := getter(&c.upgrade); config != nil {
		if utils.IsForkTransition(config.Timestamp(), from, to) {
			configs = append(configs, config)
		}
	}
	// loop on all upgrades
	for _, upgrade := range c.Upgrades {
		if config := getter(&upgrade); config != nil {
			// check if fork is activating in the specified range
			if utils.IsForkTransition(config.Timestamp(), from, to) {
				configs = append(configs, config)
			}
		}
	}
	return configs
}

// GetContractDeployerAllowListConfig returns the latest forked ContractDeployerAllowListConfig
// specified by [c] or nil if it was never enabled.
func (c *UpgradesConfig) GetContractDeployerAllowListConfig(blockTimestamp *big.Int) *ContractDeployerAllowListConfig {
	getter := func(u *upgrade) StatefulPrecompileConfig { return u.ContractDeployerAllowListConfig }
	return c.getActiveConfig(blockTimestamp, getter).(*ContractDeployerAllowListConfig)
}

// GetContractNativeMinterConfig returns the latest forked ContractNativeMinterConfig
// specified by [c] or nil if it was never enabled.
func (c *UpgradesConfig) GetContractNativeMinterConfig(blockTimestamp *big.Int) *ContractNativeMinterConfig {
	getter := func(u *upgrade) StatefulPrecompileConfig { return u.ContractNativeMinterConfig }
	return c.getActiveConfig(blockTimestamp, getter).(*ContractNativeMinterConfig)
}

// GetTxAllowListConfig returns the latest forked TxAllowListConfig
// specified by [c] or nil if it was never enabled.
func (c *UpgradesConfig) GetTxAllowListConfig(blockTimestamp *big.Int) *TxAllowListConfig {
	getter := func(u *upgrade) StatefulPrecompileConfig { return u.TxAllowListConfig }
	return c.getActiveConfig(blockTimestamp, getter).(*TxAllowListConfig)
}

// GetFeeConfigManagerConfig returns the latest forked FeeManagerConfig
// specified by [c] or nil if it was never enabled.
func (c *UpgradesConfig) GetFeeConfigManagerConfig(blockTimestamp *big.Int) *FeeConfigManagerConfig {
	getter := func(u *upgrade) StatefulPrecompileConfig { return u.FeeManagerConfig }
	return c.getActiveConfig(blockTimestamp, getter).(*FeeConfigManagerConfig)
}

// CheckCompatible checks if [newcfg] is compatible with [c] at [headTimestamp].
// Returns a ConfigCompatError if upgrades that have forked at [headTimestamp]
// are missing from [newcfg]. Upgrades that have not forked yet may be modified
// or absent from [newcfg]. Returns nil if [newcfg] is compatible with [c].
func (c *UpgradesConfig) CheckCompatible(newcfg *UpgradesConfig, headTimestamp *big.Int) *utils.ConfigCompatError {
	for _, getter := range getters {
		newUpgrades := newcfg.getActiveConfigs(nil, headTimestamp, getter)
		for i, upgrade := range c.getActiveConfigs(nil, headTimestamp, getter) {
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
	}
	return nil // newcfg is compatible
}

// EnabledStatefulPrecompiles returns a slice of stateful precompile configs that
// have been activated through an upgrade.
func (c *UpgradesConfig) EnabledStatefulPrecompiles(blockTimestamp *big.Int) []StatefulPrecompileConfig {
	statefulPrecompileConfigs := make([]StatefulPrecompileConfig, 0)
	for _, getter := range getters {
		if config := c.getActiveConfig(blockTimestamp, getter); config != nil {
			statefulPrecompileConfigs = append(statefulPrecompileConfigs, config)
		}
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
	for _, getter := range getters {
		for _, config := range c.getActiveConfigs(parentTimestamp, blockTimestamp, getter) {
			// If this transition activates the upgrade, configure the stateful precompile.
			// (or deconfigure it if it is being disabled.)
			if config.IsDisabled() {
				Deconfigure(config.Address(), statedb)
			} else {
				Configure(chainConfig, blockContext, config, statedb)
			}
		}
	}
}
