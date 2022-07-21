// (c) 2022 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/utils"
)

// precompileKey is a helper type used to reference each of the
// possible stateful precompile types that can be activated
// as a network upgrade.
type precompileKey int

const (
	contractDeployerAllowListKey precompileKey = iota + 1
	contractNativeMinterKey
	txAllowListKey
	feeManagerKey
)

var (
	precompileKeys = []precompileKey{contractDeployerAllowListKey, contractNativeMinterKey, txAllowListKey, feeManagerKey}
)

// PrecompileUpgrade is a helper struct embedded in UpgradeConfig, representing
// each of the possible stateful precompile types that can be activated
// as a network upgrade.
type PrecompileUpgrade struct {
	ContractDeployerAllowListConfig *precompile.ContractDeployerAllowListConfig `json:"contractDeployerAllowListConfig,omitempty"` // Config for the contract deployer allow list precompile
	ContractNativeMinterConfig      *precompile.ContractNativeMinterConfig      `json:"contractNativeMinterConfig,omitempty"`      // Config for the native minter precompile
	TxAllowListConfig               *precompile.TxAllowListConfig               `json:"txAllowListConfig,omitempty"`               // Config for the tx allow list precompile
	FeeManagerConfig                *precompile.FeeConfigManagerConfig          `json:"feeManagerConfig,omitempty"`                // Config for the fee manager precompile
}

func (p *PrecompileUpgrade) getByKey(key precompileKey) (precompile.StatefulPrecompileConfig, bool) {
	switch key {
	case contractDeployerAllowListKey:
		return p.ContractDeployerAllowListConfig, p.ContractDeployerAllowListConfig != nil
	case contractNativeMinterKey:
		return p.ContractNativeMinterConfig, p.ContractNativeMinterConfig != nil
	case txAllowListKey:
		return p.TxAllowListConfig, p.TxAllowListConfig != nil
	case feeManagerKey:
		return p.FeeManagerConfig, p.FeeManagerConfig != nil
	}
	panic("unknown upgrade key")
}

// ValidatePrecompileUpgrades checks [upgrades] is well formed:
// - [upgrades] must specify only one key per PrecompileUpgrade
// - the specified blockTimestamps must monotonically increase
// - the specified blockTimestamps must be compatible with those
//   specified in the chainConfig by genesis.
// - check a precompile is disabled before it is re-enabled
func (c *ChainConfig) ValidatePrecompileUpgrades(upgrades []PrecompileUpgrade) error {
	for i, upgrade := range upgrades {
		hasKey := false // used to verify if there is only one key per Upgrade

		for _, key := range precompileKeys {
			if _, ok := upgrade.getByKey(key); !ok {
				continue
			}
			if hasKey {
				return fmt.Errorf("PrecompileUpgrades[%d] has more than one key set", i)
			}
			hasKey = true
		}
	}

	for _, key := range precompileKeys {
		var (
			lastUpgraded *big.Int
			disabled     bool
		)
		// check the genesis chain config for any enabled upgrade
		if config, ok := c.PrecompileUpgrade.getByKey(key); ok {
			disabled = false
			lastUpgraded = config.Timestamp()
		} else {
			disabled = true
		}
		// next range over upgrades to verify correct use of disabled and blockTimestamps.
		for i, upgrade := range upgrades {
			if config, ok := upgrade.getByKey(key); ok {
				if disabled == config.IsDisabled() {
					return fmt.Errorf("PrecompileUpgrades[%d] disable should be [%v]", i, !disabled)
				}
				if lastUpgraded != nil && (config.Timestamp() == nil || lastUpgraded.Cmp(config.Timestamp()) > 0) {
					return fmt.Errorf("PrecompileUpgrades[%d] timestamp should not be less than [%v]", i, lastUpgraded)
				}

				disabled = config.IsDisabled()
				lastUpgraded = config.Timestamp()
			}
		}
	}

	return nil // successfully verified all conditions
}

// getActiveConfig returns the most recent config that has
// already forked, or nil if none have been configured or
// if none have forked yet.
func (c *ChainConfig) getActiveConfig(blockTimestamp *big.Int, key precompileKey, upgrades []PrecompileUpgrade) precompile.StatefulPrecompileConfig {
	configs := c.getActiveConfigs(nil, blockTimestamp, key, upgrades)
	if len(configs) == 0 {
		return nil
	}
	return configs[len(configs)-1] // return the most recent config
}

// getActiveConfigs returns all forks configured to activate in the
// specified [from, to] range, in the order of activation.
func (c *ChainConfig) getActiveConfigs(from *big.Int, to *big.Int, key precompileKey, upgrades []PrecompileUpgrade) []precompile.StatefulPrecompileConfig {
	configs := make([]precompile.StatefulPrecompileConfig, 0)
	// first check the embedded [upgrade] for precompiles configured
	// in the genesis chain config.
	if config, ok := c.PrecompileUpgrade.getByKey(key); ok {
		if utils.IsForkTransition(config.Timestamp(), from, to) {
			configs = append(configs, config)
		}
	}
	// loop on all upgrades
	for _, upgrade := range upgrades {
		if config, ok := upgrade.getByKey(key); ok {
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
func (c *ChainConfig) GetContractDeployerAllowListConfig(blockTimestamp *big.Int) *precompile.ContractDeployerAllowListConfig {
	if val := c.getActiveConfig(blockTimestamp, contractDeployerAllowListKey, c.PrecompileUpgrades); val != nil {
		return val.(*precompile.ContractDeployerAllowListConfig)
	}
	return nil
}

// GetContractNativeMinterConfig returns the latest forked ContractNativeMinterConfig
// specified by [c] or nil if it was never enabled.
func (c *ChainConfig) GetContractNativeMinterConfig(blockTimestamp *big.Int) *precompile.ContractNativeMinterConfig {
	if val := c.getActiveConfig(blockTimestamp, contractNativeMinterKey, c.PrecompileUpgrades); val != nil {
		return val.(*precompile.ContractNativeMinterConfig)
	}
	return nil
}

// GetTxAllowListConfig returns the latest forked TxAllowListConfig
// specified by [c] or nil if it was never enabled.
func (c *ChainConfig) GetTxAllowListConfig(blockTimestamp *big.Int) *precompile.TxAllowListConfig {
	if val := c.getActiveConfig(blockTimestamp, txAllowListKey, c.PrecompileUpgrades); val != nil {
		return val.(*precompile.TxAllowListConfig)
	}
	return nil
}

// GetFeeConfigManagerConfig returns the latest forked FeeManagerConfig
// specified by [c] or nil if it was never enabled.
func (c *ChainConfig) GetFeeConfigManagerConfig(blockTimestamp *big.Int) *precompile.FeeConfigManagerConfig {
	if val := c.getActiveConfig(blockTimestamp, feeManagerKey, c.PrecompileUpgrades); val != nil {
		return val.(*precompile.FeeConfigManagerConfig)
	}
	return nil
}

// CheckPrecompilesCompatible checks if [precompileUpgrades] are compatible with [c] at [headTimestamp].
// Returns a ConfigCompatError if upgrades already forked at [headTimestamp] are missing from
// [precompileUpgrades]. Upgrades not already forked may be modified or absent from [precompileUpgrades].
// Returns nil if [precompileUpgrades] is compatible with [c].
func (c *ChainConfig) CheckPrecompilesCompatible(precompileUpgrades []PrecompileUpgrade, headTimestamp *big.Int) *ConfigCompatError {
	for _, key := range precompileKeys {
		// all active upgrades must match
		activeUpgrades := c.getActiveConfigs(nil, headTimestamp, key, c.PrecompileUpgrades)
		newUpgrades := c.getActiveConfigs(nil, headTimestamp, key, precompileUpgrades)
		// first, check existing upgrades are there
		for i, upgrade := range activeUpgrades {
			if len(newUpgrades) <= i {
				// missing upgrade
				return newCompatError(
					fmt.Sprintf("missing PrecompileUpgrade[%d]", i),
					upgrade.Timestamp(),
					nil,
				)
			}
			// All upgrades that have forked must be identical.
			// TODO: verify config?
			if upgrade.Timestamp().Cmp(newUpgrades[i].Timestamp()) != 0 ||
				upgrade.IsDisabled() != newUpgrades[i].IsDisabled() {
				return newCompatError(
					fmt.Sprintf("PrecompileUpgrade[%d]", i),
					upgrade.Timestamp(),
					newUpgrades[i].Timestamp(),
				)
			}
		}
		// then, make sure newUpgrades does not have additional upgrades
		// that are already activated. (cannot perform retroactive upgrade)
		if len(newUpgrades) > len(activeUpgrades) {
			return newCompatError(
				fmt.Sprintf("cannot retroactively enable PrecompileUpgrade[%d]", len(activeUpgrades)),
				nil,
				newUpgrades[len(activeUpgrades)].Timestamp(), // this indexes to the first element in newUpgrades after the end of activeUpgrades
			)
		}
	}
	return nil // precompileUpgrades is compatible
}

// EnabledStatefulPrecompiles returns a slice of stateful precompile configs that
// have been activated through an upgrade.
func (c *ChainConfig) EnabledStatefulPrecompiles(blockTimestamp *big.Int) []precompile.StatefulPrecompileConfig {
	statefulPrecompileConfigs := make([]precompile.StatefulPrecompileConfig, 0)
	for _, key := range precompileKeys {
		if config := c.getActiveConfig(blockTimestamp, key, c.PrecompileUpgrades); config != nil {
			statefulPrecompileConfigs = append(statefulPrecompileConfigs, config)
		}
	}

	return statefulPrecompileConfigs
}

// CheckConfigurePrecompiles checks if any of the precompiles specified by [c] is enabled or disabled by the block
// transition from [parentTimestamp] to the timestamp set in [blockContext]. If this is the case, it calls [Configure]
// or [Deconfigure] to apply the necessary state transitions for the upgrade.
// This function is called:
// - within genesis setup to configure the starting state for precompiles enabled at genesis,
// - during block processing to update the state before processing the given block.
func (c *ChainConfig) CheckConfigurePrecompiles(parentTimestamp *big.Int, blockContext precompile.BlockContext, statedb precompile.StateDB) {
	blockTimestamp := blockContext.Timestamp()
	for _, key := range precompileKeys {
		for _, config := range c.getActiveConfigs(parentTimestamp, blockTimestamp, key, c.PrecompileUpgrades) {
			// If this transition activates the upgrade, configure the stateful precompile.
			// (or deconfigure it if it is being disabled.)
			if config.IsDisabled() {
				precompile.Deconfigure(config.Address(), statedb)
			} else {
				precompile.Configure(c, blockContext, config, statedb)
			}
		}
	}
}
