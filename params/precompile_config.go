// (c) 2022 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

var errMultipleKeys = errors.New("PrecompileUpgrade must have exactly one key")

// PrecompileUpgrade is a helper struct embedded in UpgradeConfig, representing
// each of the possible stateful precompile types that can be activated
// as a network upgrade.
type PrecompileUpgrade struct {
	precompile.StatefulPrecompileConfig
}

// UnmarshalJSON unmarshals the json into the correct precompile config type
// based on the key. Keys are defined in each precompile module, and registered in
// params/precompile_modules.go.
// Ex: {"feeManagerConfig": {...}} where "feeManagerConfig" is the key
func (u *PrecompileUpgrade) UnmarshalJSON(data []byte) error {
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if len(raw) != 1 {
		return errMultipleKeys
	}
	for key, value := range raw {
		module, ok := precompile.GetPrecompileModule(key)
		if !ok {
			return fmt.Errorf("unknown precompile module: %s", key)
		}
		conf := module.New()
		err := json.Unmarshal(value, conf)
		if err != nil {
			return err
		}
		u.StatefulPrecompileConfig = conf
	}
	return nil
}

// MarshalJSON marshal the precompile config into json based on the precompile key.
// Ex: {"feeManagerConfig": {...}} where "feeManagerConfig" is the key
func (u *PrecompileUpgrade) MarshalJSON() ([]byte, error) {
	res := make(ChainConfigPrecompiles)
	res[u.Key()] = u.StatefulPrecompileConfig
	return json.Marshal(res)
}

// verifyPrecompileUpgrades checks [c.PrecompileUpgrades] is well formed:
// - [upgrades] must specify exactly one key per PrecompileUpgrade
// - the specified blockTimestamps must monotonically increase
// - the specified blockTimestamps must be compatible with those
//   specified in the chainConfig by genesis.
// - check a precompile is disabled before it is re-enabled
func (c *ChainConfig) verifyPrecompileUpgrades() error {
	// Store this struct to keep track of the last upgrade for each precompile key.
	// Required for timestamp and disabled checks.
	type lastUpgradeData struct {
		lastUpgraded *big.Int
		disabled     bool
	}

	lastUpgradeMap := make(map[string]lastUpgradeData)

	// verify genesis precompiles
	for key, config := range c.Precompiles {
		if err := config.Verify(); err != nil {
			return err
		}
		// check the genesis chain config for any enabled upgrade
		lastUpgradeMap[key] = lastUpgradeData{
			disabled:     false,
			lastUpgraded: config.Timestamp(),
		}
	}

	// next range over upgrades to verify correct use of disabled and blockTimestamps.
	var lastBlockTimestamp *big.Int
	for i, upgrade := range c.PrecompileUpgrades {
		key := upgrade.Key()

		lastUpgrade, ok := lastUpgradeMap[key]
		var (
			disabled     bool
			lastUpgraded *big.Int
		)
		if !ok {
			disabled = true
			lastUpgraded = nil
		} else {
			disabled = lastUpgrade.disabled
			lastUpgraded = lastUpgrade.lastUpgraded
		}
		upgradeTimestamp := upgrade.Timestamp()

		if upgradeTimestamp == nil {
			return fmt.Errorf("PrecompileUpgrades[%d] cannot have a nil timestamp", i)
		}
		// Verify specified timestamps are monotonically increasing across all precompile keys.
		// Note: It is OK for multiple configs of different keys to specify the same timestamp.
		if lastBlockTimestamp != nil && upgradeTimestamp.Cmp(lastBlockTimestamp) < 0 {
			return fmt.Errorf("PrecompileUpgrades[%d] config timestamp (%v) < previous timestamp (%v)", i, upgradeTimestamp, lastBlockTimestamp)
		}

		if disabled == upgrade.IsDisabled() {
			return fmt.Errorf("PrecompileUpgrades[%d] disable should be [%v]", i, !disabled)
		}
		if lastUpgraded != nil && (upgradeTimestamp.Cmp(lastUpgraded) <= 0) {
			return fmt.Errorf("PrecompileUpgrades[%d] config timestamp (%v) <= previous timestamp (%v)", i, upgradeTimestamp, lastUpgraded)
		}

		if err := upgrade.Verify(); err != nil {
			return err
		}

		lastUpgradeMap[key] = lastUpgradeData{
			disabled:     upgrade.IsDisabled(),
			lastUpgraded: upgradeTimestamp,
		}

		lastBlockTimestamp = upgradeTimestamp
	}

	return nil
}

// GetActivePrecompileConfig returns the most recent precompile config corresponding to [address].
// If none have occurred, returns nil.
func (c *ChainConfig) GetActivePrecompileConfig(address common.Address, blockTimestamp *big.Int) precompile.StatefulPrecompileConfig {
	configs := c.getActivatingPrecompileConfigs(address, nil, blockTimestamp, c.PrecompileUpgrades)
	if len(configs) == 0 {
		return nil
	}
	return configs[len(configs)-1] // return the most recent config
}

// getActivatingPrecompileConfigs returns all upgrades configured to activate during the state transition from a block with timestamp [from]
// to a block with timestamp [to].
func (c *ChainConfig) getActivatingPrecompileConfigs(address common.Address, from *big.Int, to *big.Int, upgrades []PrecompileUpgrade) []precompile.StatefulPrecompileConfig {
	configs := make([]precompile.StatefulPrecompileConfig, 0)
	// Get key from address.
	module, ok := precompile.GetPrecompileModuleByAddress(address)
	if !ok {
		return configs
	}

	key := module.Key()

	// First check the embedded [upgrade] for precompiles configured
	// in the genesis chain config.
	if config, ok := c.Precompiles[key]; ok {
		if utils.IsForkTransition(config.Timestamp(), from, to) {
			configs = append(configs, config)
		}
	}
	// Loop over all upgrades checking for the requested precompile config.
	for _, upgrade := range upgrades {
		if upgrade.Key() == key {
			// Check if the precompile activates in the specified range.
			if utils.IsForkTransition(upgrade.Timestamp(), from, to) {
				configs = append(configs, upgrade.StatefulPrecompileConfig)
			}
		}
	}
	return configs
}

// CheckPrecompilesCompatible checks if [precompileUpgrades] are compatible with [c] at [headTimestamp].
// Returns a ConfigCompatError if upgrades already forked at [headTimestamp] are missing from
// [precompileUpgrades]. Upgrades not already forked may be modified or absent from [precompileUpgrades].
// Returns nil if [precompileUpgrades] is compatible with [c].
// Assumes given timestamp is the last accepted block timestamp.
// This ensures that as long as the node has not accepted a block with a different rule set it will allow a new upgrade to be applied as long as it activates after the last accepted block.
func (c *ChainConfig) CheckPrecompilesCompatible(precompileUpgrades []PrecompileUpgrade, lastTimestamp *big.Int) *ConfigCompatError {
	for _, module := range precompile.RegisteredModules() {
		if err := c.checkPrecompileCompatible(module.Address(), precompileUpgrades, lastTimestamp); err != nil {
			return err
		}
	}

	return nil
}

// checkPrecompileCompatible verifies that the precompile specified by [address] is compatible between [c] and [precompileUpgrades] at [headTimestamp].
// Returns an error if upgrades already forked at [headTimestamp] are missing from [precompileUpgrades].
// Upgrades that have already gone into effect cannot be modified or absent from [precompileUpgrades].
func (c *ChainConfig) checkPrecompileCompatible(address common.Address, precompileUpgrades []PrecompileUpgrade, lastTimestamp *big.Int) *ConfigCompatError {
	// all active upgrades must match
	activeUpgrades := c.getActivatingPrecompileConfigs(address, nil, lastTimestamp, c.PrecompileUpgrades)
	newUpgrades := c.getActivatingPrecompileConfigs(address, nil, lastTimestamp, precompileUpgrades)

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
		if !upgrade.Equal(newUpgrades[i]) {
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

	return nil
}

// EnabledStatefulPrecompiles returns a slice of stateful precompile configs that
// have been activated through an upgrade.
func (c *ChainConfig) EnabledStatefulPrecompiles(blockTimestamp *big.Int) []precompile.StatefulPrecompileConfig {
	statefulPrecompileConfigs := make([]precompile.StatefulPrecompileConfig, 0)
	for _, module := range precompile.RegisteredModules() {
		if config := c.GetActivePrecompileConfig(module.Address(), blockTimestamp); config != nil {
			statefulPrecompileConfigs = append(statefulPrecompileConfigs, config)
		}
	}

	return statefulPrecompileConfigs
}

// ConfigurePrecompiles checks if any of the precompiles specified by the chain config are enabled or disabled by the block
// transition from [parentTimestamp] to the timestamp set in [blockContext]. If this is the case, it calls [Configure]
// or [Deconfigure] to apply the necessary state transitions for the upgrade.
// This function is called:
// - within genesis setup to configure the starting state for precompiles enabled at genesis,
// - during block processing to update the state before processing the given block.
func (c *ChainConfig) ConfigurePrecompiles(parentTimestamp *big.Int, blockContext precompile.BlockContext, statedb precompile.StateDB) error {
	blockTimestamp := blockContext.Timestamp()
	for _, module := range precompile.RegisteredModules() { // Note: configure precompiles in a deterministic order.
		key := module.Key()
		for _, config := range c.getActivatingPrecompileConfigs(module.Address(), parentTimestamp, blockTimestamp, c.PrecompileUpgrades) {
			// If this transition activates the upgrade, configure the stateful precompile.
			// (or deconfigure it if it is being disabled.)
			if config.IsDisabled() {
				log.Info("Disabling precompile", "name", key)
				statedb.Suicide(module.Address())
				// Calling Finalise here effectively commits Suicide call and wipes the contract state.
				// This enables re-configuration of the same contract state in the same block.
				// Without an immediate Finalise call after the Suicide, a reconfigured precompiled state can be wiped out
				// since Suicide will be committed after the reconfiguration.
				statedb.Finalise(true)
			} else {
				log.Info("Activating new precompile", "name", key, "config", config)
				if err := precompile.Configure(c, blockContext, config, statedb); err != nil {
					return fmt.Errorf("could not configure precompile, name: %s, reason: %w", key, err)
				}
			}
		}
	}
	return nil
}
