// (c) 2022 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/config"
	precompileConfig "github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
)

var errMultipleKeys = errors.New("PrecompileUpgrade must have exactly one key")
var errNoKey = errors.New("PrecompileUpgrade cannot be empty")

// PrecompileUpgrade is a helper struct embedded in UpgradeConfig.
// It is used to unmarshal the json into the correct precompile config type
// based on the key. Keys are defined in each precompile module, and registered in
// params/precompile_modules.go.
type PrecompileUpgrade struct {
	config.Config
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
	if len(raw) == 0 {
		return errNoKey
	}
	if len(raw) > 1 {
		return errMultipleKeys
	}
	for key, value := range raw {
		config, ok := config.GetNewConfig(key)
		if !ok {
			return fmt.Errorf("unknown precompile config: %s", key)
		}
		err := json.Unmarshal(value, config)
		if err != nil {
			return err
		}
		u.Config = config
	}
	return nil
}

// MarshalJSON marshal the precompile config into json based on the precompile key.
// Ex: {"feeManagerConfig": {...}} where "feeManagerConfig" is the key
func (u *PrecompileUpgrade) MarshalJSON() ([]byte, error) {
	res := make(precompileConfig.Configs)
	res[u.Key()] = u.Config
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
		blockTimestamp *big.Int
		disabled       bool
	}

	lastPrecompileUpgrades := make(map[string]lastUpgradeData)

	// verify genesis precompiles
	for key, config := range c.GenesisPrecompiles {
		if err := config.Verify(); err != nil {
			return err
		}
		// if the precompile is disabled at genesis, skip it.
		if config.Timestamp() == nil {
			continue
		}
		// check the genesis chain config for any enabled upgrade
		lastPrecompileUpgrades[key] = lastUpgradeData{
			disabled:       false,
			blockTimestamp: config.Timestamp(),
		}
	}

	// next range over upgrades to verify correct use of disabled and blockTimestamps.
	// previousUpgradeTimestamp is used to verify monotonically increasing timestamps.
	var previousUpgradeTimestamp *big.Int
	for i, upgrade := range c.PrecompileUpgrades {
		key := upgrade.Key()

		// lastUpgradeByKey is the previous processed upgrade for this precompile key.
		lastUpgradeByKey, ok := lastPrecompileUpgrades[key]
		var (
			disabled      bool
			lastTimestamp *big.Int
		)
		if !ok {
			disabled = true
			lastTimestamp = nil
		} else {
			disabled = lastUpgradeByKey.disabled
			lastTimestamp = lastUpgradeByKey.blockTimestamp
		}
		upgradeTimestamp := upgrade.Timestamp()

		if upgradeTimestamp == nil {
			return fmt.Errorf("PrecompileUpgrade (%s) at [%d]: block timestamp cannot be nil ", key, i)
		}
		// Verify specified timestamps are monotonically increasing across all precompile keys.
		// Note: It is OK for multiple configs of DIFFERENT keys to specify the same timestamp.
		if previousUpgradeTimestamp != nil && upgradeTimestamp.Cmp(previousUpgradeTimestamp) < 0 {
			return fmt.Errorf("PrecompileUpgrade (%s) at [%d]: config block timestamp (%v) < previous timestamp (%v)", key, i, upgradeTimestamp, previousUpgradeTimestamp)
		}

		if disabled == upgrade.IsDisabled() {
			return fmt.Errorf("PrecompileUpgrade (%s) at [%d]: disable should be [%v]", key, i, !disabled)
		}
		// Verify specified timestamps are monotonically increasing across same precompile keys.
		// Note: It is NOT OK for multiple configs of the SAME key to specify the same timestamp.
		if lastTimestamp != nil && (upgradeTimestamp.Cmp(lastTimestamp) <= 0) {
			return fmt.Errorf("PrecompileUpgrade (%s) at [%d]: config block timestamp (%v) <= previous timestamp of same key (%v)", key, i, upgradeTimestamp, lastTimestamp)
		}

		if err := upgrade.Verify(); err != nil {
			return err
		}

		lastPrecompileUpgrades[key] = lastUpgradeData{
			disabled:       upgrade.IsDisabled(),
			blockTimestamp: upgradeTimestamp,
		}

		previousUpgradeTimestamp = upgradeTimestamp
	}

	return nil
}

// GetActivePrecompileConfig returns the most recent precompile config corresponding to [address].
// If none have occurred, returns nil.
func (c *ChainConfig) GetActivePrecompileConfig(address common.Address, blockTimestamp *big.Int) config.Config {
	configs := c.GetActivatingPrecompileConfigs(address, nil, blockTimestamp, c.PrecompileUpgrades)
	if len(configs) == 0 {
		return nil
	}
	return configs[len(configs)-1] // return the most recent config
}

// GetActivatingPrecompileConfigs returns all upgrades configured to activate during the state transition from a block with timestamp [from]
// to a block with timestamp [to].
func (c *ChainConfig) GetActivatingPrecompileConfigs(address common.Address, from *big.Int, to *big.Int, upgrades []PrecompileUpgrade) []config.Config {
	configs := make([]config.Config, 0)

	// First check the embedded [upgrade] for precompiles configured
	// in the genesis chain config.
	if config, ok := c.GenesisPrecompiles.GetConfigByAddress(address); ok {
		if utils.IsForkTransition(config.Timestamp(), from, to) {
			configs = append(configs, config)
		}
	}
	// Loop over all upgrades checking for the requested precompile config.
	for _, upgrade := range upgrades {
		if upgrade.Address() == address {
			// Check if the precompile activates in the specified range.
			if utils.IsForkTransition(upgrade.Timestamp(), from, to) {
				configs = append(configs, upgrade.Config)
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
	for _, address := range config.GetAddresses() {
		if err := c.checkPrecompileCompatible(address, precompileUpgrades, lastTimestamp); err != nil {
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
	activeUpgrades := c.GetActivatingPrecompileConfigs(address, nil, lastTimestamp, c.PrecompileUpgrades)
	newUpgrades := c.GetActivatingPrecompileConfigs(address, nil, lastTimestamp, precompileUpgrades)

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
func (c *ChainConfig) EnabledStatefulPrecompiles(blockTimestamp *big.Int) []config.Config {
	statefulPrecompileConfigs := make([]config.Config, 0)
	for _, address := range config.GetAddresses() {
		if config := c.GetActivePrecompileConfig(address, blockTimestamp); config != nil {
			statefulPrecompileConfigs = append(statefulPrecompileConfigs, config)
		}
	}

	return statefulPrecompileConfigs
}
