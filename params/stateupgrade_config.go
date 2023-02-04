// (c) 2022 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/stateupgrade"
	"github.com/ava-labs/subnet-evm/utils"
)

// stateUpgradeKey is a helper type used to reference each of the
// possible state upgrade types that can be activated
// as a network upgrade.
type stateUpgradeKey int

const (
	stateUpgradeCodeKey stateUpgradeKey = iota + 1
	stateUpgradeStateKey
	stateUpgradeBalanceKey
	stateUpgradeDeployKey
)
var stateUpgradeKeys = []stateUpgradeKey{stateUpgradeCodeKey, stateUpgradeStateKey, stateUpgradeBalanceKey, stateUpgradeDeployKey}

// StateUpgrade is a helper struct embedded in UpgradeConfig, representing
// each of the possible stateful state upgrade types that can be activated
// as a network upgrade.
type StateUpgrade struct {
	StateUpgradeCodeConfig 			*stateupgrade.StateUpgradeCodeConfig		`json:"stateUpgradeCodeConfig,omitempty"` 		   // Config for the state upgrade code upgrader
	StateUpgradeStateConfig 		*stateupgrade.StateUpgradeStateConfig		`json:"stateUpgradeStateConfig,omitempty"` 		   // Config for the state upgrade state upgrader
	//StateUpgradeBalanceConfig 		*stateupgrade.StateUpgradeBalanceConfig		`json:"stateUpgradeBalanceConfig,omitempty"` 	   // Config for the state upgrade balance upgrader
	//StateUpgradeDeployConfig 		*stateupgrade.StateUpgradeDeployConfig		`json:"stateUpgradeDeployConfig,omitempty"` 	   // Config for the state upgrade deploy upgrader
}

func (p *StateUpgrade) getByKey(key stateUpgradeKey) (stateupgrade.StateUpgradeConfig, bool) {
	switch key {
	case stateUpgradeCodeKey:
		return p.StateUpgradeCodeConfig, p.StateUpgradeCodeConfig != nil
	default:
		panic(fmt.Sprintf("unknown state upgrade key: %v", key))
	}
}

// verifyStateUpgrades checks [c.StateUpgrades] is well formed:
// - [upgrades] must specify exactly one key per StateUpgrade
// - the specified blockTimestamps must monotonically increase
// - the specified blockTimestamps must be compatible with those
//   specified in the chainConfig by genesis.
// - check a state upgrade is disabled before it is re-enabled
func (c *ChainConfig) verifyStateUpgrades() error {
	var lastBlockTimestamp *big.Int
	for i, upgrade := range c.StateUpgrades {
		hasKey := false // used to verify if there is only one key per Upgrade

		for _, key := range stateUpgradeKeys {
			config, ok := upgrade.getByKey(key)
			if !ok {
				continue
			}
			if hasKey {
				return fmt.Errorf("StateUpgrades[%d] has more than one key set", i)
			}
			configTimestamp := config.Timestamp()
			if configTimestamp == nil {
				return fmt.Errorf("StateUpgrades[%d] cannot have a nil timestamp", i)
			}
			// Verify specified timestamps are monotonically increasing across all state upgrade keys.
			// Note: It is OK for multiple configs of different keys to specify the same timestamp.
			if lastBlockTimestamp != nil && configTimestamp.Cmp(lastBlockTimestamp) < 0 {
				return fmt.Errorf("StateUpgrades[%d] config timestamp (%v) < previous timestamp (%v)", i, configTimestamp, lastBlockTimestamp)
			}
			lastBlockTimestamp = configTimestamp
			hasKey = true
		}
		if !hasKey {
			return fmt.Errorf("empty state upgrade at index %d", i)
		}
	}

	for _, key := range stateUpgradeKeys {
		var (
			lastUpgraded *big.Int
			disabled     bool
		)
		disabled = true
		// next range over upgrades to verify correct use of disabled and blockTimestamps.
		for i, upgrade := range c.StateUpgrades {
			config, ok := upgrade.getByKey(key)
			// Skip the upgrade if it's not relevant to [key].
			if !ok {
				continue
			}

			if disabled == config.IsDisabled() {
				return fmt.Errorf("StateUpgrades[%d] disable should be [%v]", i, !disabled)
			}
			if lastUpgraded != nil && (config.Timestamp().Cmp(lastUpgraded) <= 0) {
				return fmt.Errorf("StateUpgrades[%d] config timestamp (%v) <= previous timestamp (%v)", i, config.Timestamp(), lastUpgraded)
			}

			if err := config.Verify(); err != nil {
				return err
			}

			disabled = config.IsDisabled()
			lastUpgraded = config.Timestamp()
		}
	}

	return nil
}

func (c *ChainConfig) GetActiveStateUpgrades(blockTimestamp *big.Int) StateUpgrade {
	u := StateUpgrade{}
	if config := c.GetStateUpgradeCodeConfig(blockTimestamp); config != nil && !config.Disable {
		u.StateUpgradeCodeConfig = config
	}

	return u
}

// getActiveStateUpgradeConfig returns the most recent state upgrade config corresponding to [key].
// If none have occurred, returns nil.
func (c *ChainConfig) getActiveStateUpgradeConfig(blockTimestamp *big.Int, key stateUpgradeKey, upgrades []StateUpgrade) stateupgrade.StateUpgradeConfig {
	configs := c.getActivatingStateUpgradeConfigs(nil, blockTimestamp, key, upgrades)
	if len(configs) == 0 {
		return nil
	}
	return configs[len(configs)-1] // return the most recent config
}

// getActivatingStateUpgradeConfigs returns all forks configured to activate during the state transition from a block with timestamp [from]
// to a block with timestamp [to].
func (c *ChainConfig) getActivatingStateUpgradeConfigs(from *big.Int, to *big.Int, key stateUpgradeKey, upgrades []StateUpgrade) []stateupgrade.StateUpgradeConfig {
	configs := make([]stateupgrade.StateUpgradeConfig, 0)
	// Loop over all upgrades checking for the requested state upgrade config.
	for _, upgrade := range upgrades {
		if config, ok := upgrade.getByKey(key); ok {
			// Check if the state upgrade activates in the specified range.
			if utils.IsForkTransition(config.Timestamp(), from, to) {
				configs = append(configs, config)
			}
		}
	}
	return configs
}

// GetStateUpgradeCodeConfig returns the latest forked StateUpgradeCodeConfig
// specified by [c] or nil if it was never enabled.
func (c *ChainConfig) GetStateUpgradeCodeConfig(blockTimestamp *big.Int) *stateupgrade.StateUpgradeCodeConfig {
	if val := c.getActiveStateUpgradeConfig(blockTimestamp, stateUpgradeCodeKey, c.StateUpgrades); val != nil {
		return val.(*stateupgrade.StateUpgradeCodeConfig)
	}
	return nil
}

// CheckStateUpgradesCompatible checks if [StateUpgrades] are compatible with [c] at [headTimestamp].
// Returns a ConfigCompatError if upgrades already forked at [headTimestamp] are missing from
// [stateUpgrades]. Upgrades not already forked may be modified or absent from [stateUpgrades].
// Returns nil if [stateUpgrades] is compatible with [c].
// Assumes given timestamp is the last accepted block timestamp.
// This ensures that as long as the node has not accepted a block with a different rule set it will allow a new upgrade to be applied as long as it activates after the last accepted block.
func (c *ChainConfig) CheckStateUpgradesCompatible(stateUpgrades []StateUpgrade, lastTimestamp *big.Int) *ConfigCompatError {
	for _, key := range stateUpgradeKeys {
		if err := c.checkStateUpgradeCompatible(key, stateUpgrades, lastTimestamp); err != nil {
			return err
		}
	}

	return nil
}

// checkStateUpgradeCompatible verifies that the state upgrade specified by [key] is compatible between [c] and [stateUpgrades] at [headTimestamp].
// Returns an error if upgrades already forked at [headTimestamp] are missing from [stateUpgrades].
// Upgrades that have already gone into effect cannot be modified or absent from [stateUpgrades].
func (c *ChainConfig) checkStateUpgradeCompatible(key stateUpgradeKey, stateUpgrades []StateUpgrade, lastTimestamp *big.Int) *ConfigCompatError {
	// all active upgrades must match
	activeUpgrades := c.getActivatingStateUpgradeConfigs(nil, lastTimestamp, key, c.StateUpgrades)
	newUpgrades := c.getActivatingStateUpgradeConfigs(nil, lastTimestamp, key, stateUpgrades)

	// first, check existing upgrades are there
	for i, upgrade := range activeUpgrades {
		if len(newUpgrades) <= i {
			// missing upgrade
			return newCompatError(
				fmt.Sprintf("missing StateUpgrade[%d]", i),
				upgrade.Timestamp(),
				nil,
			)
		}
		// All upgrades that have forked must be identical.
		if !upgrade.Equal(newUpgrades[i]) {
			return newCompatError(
				fmt.Sprintf("StateUpgrade[%d]", i),
				upgrade.Timestamp(),
				newUpgrades[i].Timestamp(),
			)
		}
	}
	// then, make sure newUpgrades does not have additional upgrades
	// that are already activated. (cannot perform retroactive upgrade)
	if len(newUpgrades) > len(activeUpgrades) {
		return newCompatError(
			fmt.Sprintf("cannot retroactively enable StateUpgrade[%d]", len(activeUpgrades)),
			nil,
			newUpgrades[len(activeUpgrades)].Timestamp(), // this indexes to the first element in newUpgrades after the end of activeUpgrades
		)
	}

	return nil
}

// EnabledStateUpgrades returns a slice of stateful state upgrade configs that
// have been activated through an upgrade.
func (c *ChainConfig) EnabledStateUpgrades(blockTimestamp *big.Int) []stateupgrade.StateUpgradeConfig {
	stateUpgradeConfigs := make([]stateupgrade.StateUpgradeConfig, 0)
	for _, key := range stateUpgradeKeys {
		if config := c.getActiveStateUpgradeConfig(blockTimestamp, key, c.StateUpgrades); config != nil {
			stateUpgradeConfigs = append(stateUpgradeConfigs, config)
		}
	}

	return stateUpgradeConfigs
}