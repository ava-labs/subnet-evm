// (c) 2023 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"encoding/json"
	"fmt"

	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
)

// PrecompileUpgrade is a helper struct embedded in UpgradeConfig.
// It is used to unmarshal the json into the correct precompile config type
// based on the key. Keys are defined in each precompile module, and registered in
// precompile/registry/registry.go.
type PrecompileUpgrade struct {
	precompileconfig.Config
}

// MarshalJSON marshal the precompile config into json based on the precompile key.
// Ex: {"feeManagerConfig": {...}} where "feeManagerConfig" is the key
func (u *PrecompileUpgrade) MarshalJSON() ([]byte, error) {
	res := make(map[string]precompileconfig.Config)
	res[u.Key()] = u.Config
	return json.Marshal(res)
}

// verifyPrecompileUpgrades checks [c.PrecompileUpgrades] is well formed:
//   - [upgrades] must specify exactly one key per PrecompileUpgrade
//   - the specified blockTimestamps must monotonically increase
//   - the specified blockTimestamps must be compatible with those
//     specified in the chainConfig by genesis.
//   - check a precompile is disabled before it is re-enabled
func (c *ChainConfig) verifyPrecompileUpgrades() error {
	// Store this struct to keep track of the last upgrade for each precompile key.
	// Required for timestamp and disabled checks.
	type lastUpgradeData struct {
		blockTimestamp uint64
		disabled       bool
	}

	lastPrecompileUpgrades := make(map[string]lastUpgradeData)

	// verify genesis precompiles
	for key, config := range c.GenesisPrecompiles {
		if err := config.Verify(c); err != nil {
			return err
		}
		// if the precompile is disabled at genesis, skip it.
		if config.Timestamp() == nil {
			continue
		}
		// check the genesis chain config for any enabled upgrade
		lastPrecompileUpgrades[key] = lastUpgradeData{
			disabled:       false,
			blockTimestamp: *config.Timestamp(),
		}
	}

	// next range over upgrades to verify correct use of disabled and blockTimestamps.
	// previousUpgradeTimestamp is used to verify monotonically increasing timestamps.
	var previousUpgradeTimestamp *uint64
	for i, upgrade := range c.PrecompileUpgrades {
		key := upgrade.Key()

		// lastUpgradeByKey is the previous processed upgrade for this precompile key.
		lastUpgradeByKey, ok := lastPrecompileUpgrades[key]
		var (
			disabled      bool
			lastTimestamp *uint64
		)
		if !ok {
			disabled = true
			lastTimestamp = nil
		} else {
			disabled = lastUpgradeByKey.disabled
			lastTimestamp = utils.NewUint64(lastUpgradeByKey.blockTimestamp)
		}
		upgradeTimestamp := upgrade.Timestamp()

		if upgradeTimestamp == nil {
			return fmt.Errorf("PrecompileUpgrade (%s) at [%d]: block timestamp cannot be nil ", key, i)
		}
		// Verify specified timestamps are monotonically increasing across all precompile keys.
		// Note: It is OK for multiple configs of DIFFERENT keys to specify the same timestamp.
		if previousUpgradeTimestamp != nil && *upgradeTimestamp < *previousUpgradeTimestamp {
			return fmt.Errorf("PrecompileUpgrade (%s) at [%d]: config block timestamp (%v) < previous timestamp (%v)", key, i, *upgradeTimestamp, *previousUpgradeTimestamp)
		}

		if disabled == upgrade.IsDisabled() {
			return fmt.Errorf("PrecompileUpgrade (%s) at [%d]: disable should be [%v]", key, i, !disabled)
		}
		// Verify specified timestamps are monotonically increasing across same precompile keys.
		// Note: It is NOT OK for multiple configs of the SAME key to specify the same timestamp.
		if lastTimestamp != nil && *upgradeTimestamp <= *lastTimestamp {
			return fmt.Errorf("PrecompileUpgrade (%s) at [%d]: config block timestamp (%v) <= previous timestamp (%v) of same key", key, i, *upgradeTimestamp, *lastTimestamp)
		}

		if err := upgrade.Verify(c); err != nil {
			return err
		}

		lastPrecompileUpgrades[key] = lastUpgradeData{
			disabled:       upgrade.IsDisabled(),
			blockTimestamp: *upgradeTimestamp,
		}

		previousUpgradeTimestamp = upgradeTimestamp
	}

	return nil
}

// getActivePrecompileConfig returns the most recent precompile config corresponding to [address].
// If none have occurred, returns nil.
func (c *ChainConfig) getActivePrecompileConfig(address common.Address, timestamp uint64) precompileconfig.Config {
	configs := c.GetActivatingPrecompileConfigs(address, nil, timestamp, c.PrecompileUpgrades)
	if len(configs) == 0 {
		return nil
	}
	return configs[len(configs)-1] // return the most recent config
}

// GetActivatingPrecompileConfigs returns all precompile upgrades configured to activate during the
// state transition from a block with timestamp [from] to a block with timestamp [to].
func (c *ChainConfig) GetActivatingPrecompileConfigs(address common.Address, from *uint64, to uint64, upgrades []PrecompileUpgrade) []precompileconfig.Config {
	var configs []precompileconfig.Config
	maybeAppend := func(pc precompileconfig.Config) {
		if pc.Address() == address && utils.IsForkTransition(pc.Timestamp(), from, to) {
			configs = append(configs, pc)
		}
	}
	for _, p := range c.GenesisPrecompiles {
		maybeAppend(p)
	}
	for _, upgrade := range upgrades {
		maybeAppend(upgrade.Config)
	}
	return configs
}

// CheckPrecompilesCompatible checks if [precompileUpgrades] are compatible with [c] at [headTimestamp].
// Returns a ConfigCompatError if upgrades already activated at [headTimestamp] are missing from
// [precompileUpgrades]. Upgrades not already activated may be modified or absent from [precompileUpgrades].
// Returns nil if [precompileUpgrades] is compatible with [c].
// Assumes given timestamp is the last accepted block timestamp.
// This ensures that as long as the node has not accepted a block with a different rule set it will allow a
// new upgrade to be applied as long as it activates after the last accepted block.
func (c *ChainConfig) CheckPrecompilesCompatible(precompileUpgrades []PrecompileUpgrade, time uint64) *ConfigCompatError {
	addrs := c.allPrecompileAddresses(precompileUpgrades...)
	for _, a := range addrs {
		if err := c.checkPrecompileCompatible(a, precompileUpgrades, time); err != nil {
			return err
		}
	}
	return nil
}

// checkPrecompileCompatible verifies that the precompile specified by [address] is compatible between [c]
// and [precompileUpgrades] at [headTimestamp].
// Returns an error if upgrades already activated at [headTimestamp] are missing from [precompileUpgrades].
// Upgrades that have already gone into effect cannot be modified or absent from [precompileUpgrades].
func (c *ChainConfig) checkPrecompileCompatible(address common.Address, precompileUpgrades []PrecompileUpgrade, time uint64) *ConfigCompatError {
	// All active upgrades (from nil to [lastTimestamp]) must match.
	activeUpgrades := c.GetActivatingPrecompileConfigs(address, nil, time, c.PrecompileUpgrades)
	newUpgrades := c.GetActivatingPrecompileConfigs(address, nil, time, precompileUpgrades)

	// Check activated upgrades are still present.
	for i, upgrade := range activeUpgrades {
		if len(newUpgrades) <= i {
			// missing upgrade
			return newTimestampCompatError(
				fmt.Sprintf("missing PrecompileUpgrade[%d]", i),
				upgrade.Timestamp(),
				nil,
			)
		}
		// All upgrades that have activated must be identical.
		if !upgrade.Equal(newUpgrades[i]) {
			return newTimestampCompatError(
				fmt.Sprintf("PrecompileUpgrade[%d]", i),
				upgrade.Timestamp(),
				newUpgrades[i].Timestamp(),
			)
		}
	}
	// then, make sure newUpgrades does not have additional upgrades
	// that are already activated. (cannot perform retroactive upgrade)
	if len(newUpgrades) > len(activeUpgrades) {
		return newTimestampCompatError(
			fmt.Sprintf("cannot retroactively enable PrecompileUpgrade[%d]", len(activeUpgrades)),
			nil,
			newUpgrades[len(activeUpgrades)].Timestamp(), // this indexes to the first element in newUpgrades after the end of activeUpgrades
		)
	}

	return nil
}

// EnabledStatefulPrecompiles returns current stateful precompile configs that are enabled at [blockTimestamp].
func (c *ChainConfig) EnabledStatefulPrecompiles(blockTimestamp uint64) Precompiles {
	statefulPrecompileConfigs := make(Precompiles)
	for key, addr := range c.allPrecompileAddresses() {
		if config := c.getActivePrecompileConfig(addr, blockTimestamp); config != nil && !config.IsDisabled() {
			statefulPrecompileConfigs[key] = config
		}
	}
	return statefulPrecompileConfigs
}

// allPrecompileAddresses returns a mapping from precompile config key to
// address for all precompiles defined in [ChainConfig.GenesisPrecompiles],
// [ChainConfig.UpgradeConfig.PrecompileUpgrades], and the `extra` upgrades.
func (c *ChainConfig) allPrecompileAddresses(extra ...PrecompileUpgrade) map[string]common.Address {
	all := make(map[string]common.Address)
	add := func(pc precompileconfig.Config) {
		if a, ok := all[pc.Key()]; ok && a != pc.Address() {
			panic("DO NOT MERGE")
		}
		all[pc.Key()] = pc.Address()
	}

	for _, p := range c.GenesisPrecompiles {
		add(p)
	}
	for _, p := range c.UpgradeConfig.PrecompileUpgrades {
		add(p.Config)
	}
	for _, p := range extra {
		add(p.Config)
	}
	return all
}
