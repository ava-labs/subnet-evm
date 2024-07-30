// (c) 2023 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/utils"
)

// PrecompileUpgrade is a helper struct embedded in UpgradeConfig.
// It is used to unmarshal the json into the correct precompile config type
// based on the key. Keys are defined in each precompile module, and registered in
// precompile/registry/registry.go.
type PrecompileUpgrade struct {
	precompileconfig.Config
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
