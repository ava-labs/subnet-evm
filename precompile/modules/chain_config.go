package modules

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
)

func GetActivatingPrecompileConfigs(c *params.ChainConfig, address common.Address, from *uint64, to uint64, upgrades []params.PrecompileUpgrade) []precompileconfig.Config {
	// Get key from address.
	module, ok := GetPrecompileModuleByAddress(address)
	if !ok {
		return nil
	}
	var configs []precompileconfig.Config
	key := module.ConfigKey
	// First check the embedded [upgrade] for precompiles configured
	// in the genesis chain config.
	if config, ok := c.GenesisPrecompiles[key]; ok {
		if utils.IsForkTransition(config.Timestamp(), from, to) {
			configs = append(configs, config)
		}
	}
	// Loop over all upgrades checking for the requested precompile config.
	for _, upgrade := range upgrades {
		if upgrade.Key() == key {
			// Check if the precompile activates in the specified range.
			if utils.IsForkTransition(upgrade.Timestamp(), from, to) {
				configs = append(configs, upgrade.Config)
			}
		}
	}
	return configs
}

func GetActivePrecompileConfig(c *params.ChainConfig, address common.Address, timestamp uint64) precompileconfig.Config {
	configs := GetActivatingPrecompileConfigs(c, address, nil, timestamp, c.PrecompileUpgrades)
	if len(configs) == 0 {
		return nil
	}
	return configs[len(configs)-1] // return the most recent config
}

// CheckPrecompilesCompatible checks if [precompileUpgrades] are compatible with [c] at [headTimestamp].
// Returns a ConfigCompatError if upgrades already activated at [headTimestamp] are missing from
// [precompileUpgrades]. Upgrades not already activated may be modified or absent from [precompileUpgrades].
// Returns nil if [precompileUpgrades] is compatible with [c].
// Assumes given timestamp is the last accepted block timestamp.
// This ensures that as long as the node has not accepted a block with a different rule set it will allow a
// new upgrade to be applied as long as it activates after the last accepted block.
func CheckPrecompilesCompatible(c *params.ChainConfig, precompileUpgrades []params.PrecompileUpgrade, time uint64) *params.ConfigCompatError {
	for _, module := range RegisteredModules() {
		if err := checkPrecompileCompatible(c, module.Address, precompileUpgrades, time); err != nil {
			return err
		}
	}

	return nil
}

// checkPrecompileCompatible verifies that the precompile specified by [address] is compatible between [c]
// and [precompileUpgrades] at [headTimestamp].
// Returns an error if upgrades already activated at [headTimestamp] are missing from [precompileUpgrades].
// Upgrades that have already gone into effect cannot be modified or absent from [precompileUpgrades].
func checkPrecompileCompatible(c *params.ChainConfig, address common.Address, precompileUpgrades []params.PrecompileUpgrade, time uint64) *params.ConfigCompatError {
	// All active upgrades (from nil to [lastTimestamp]) must match.
	activeUpgrades := GetActivatingPrecompileConfigs(c, address, nil, time, c.PrecompileUpgrades)
	newUpgrades := GetActivatingPrecompileConfigs(c, address, nil, time, precompileUpgrades)

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

func newTimestampCompatError(what string, storedtime, newtime *uint64) *params.ConfigCompatError {
	return nil // TODO: DO NOT MERGE
}

// EnabledStatefulPrecompiles returns current stateful precompile configs that are enabled at [blockTimestamp].
func EnabledStatefulPrecompiles(c *params.ChainConfig, blockTimestamp uint64) params.Precompiles {
	statefulPrecompileConfigs := make(params.Precompiles)
	for _, module := range RegisteredModules() {
		if config := GetActivePrecompileConfig(c, module.Address, blockTimestamp); config != nil && !config.IsDisabled() {
			statefulPrecompileConfigs[module.ConfigKey] = config
		}
	}

	return statefulPrecompileConfigs
}

// IsPrecompileEnabled returns whether precompile with [address] is enabled at [timestamp].
func IsPrecompileEnabled(c *params.ChainConfig, address common.Address, timestamp uint64) bool {
	config := GetActivePrecompileConfig(c, address, timestamp)
	return config != nil && !config.IsDisabled()
}
