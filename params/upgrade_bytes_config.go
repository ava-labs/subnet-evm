// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"encoding/json"
	"math/big"
)

// ApplyUpgradeBytes applies modifications from upgradeBytes to chainConfig
// if upgradeBytes is compatible with activated forks.
func (c *ChainConfig) ApplyUpgradeBytes(upgradeBytes []byte, headTimestamp *big.Int) error {
	var upgradeConfig UpgradeConfig

	// Note: passing an empty slice is considered equivalent to an empty upgradeBytesConfig
	// we will still verify the empty config against the existing chainConfig, to ensure
	// activated upgrades are not removed.
	if len(upgradeBytes) > 0 {
		if err := json.Unmarshal(upgradeBytes, &upgradeConfig); err != nil {
			return err
		}
	}

	// Verify the precompile upgrades are internally consistent given the existing chainConfig.
	if err := c.ValidatePrecompileUpgrades(upgradeConfig.PrecompileUpgrades); err != nil {
		return err
	}

	// Verify newly specified NewtorkUpgrades is compatible with the existing chainConfig.
	if c.UpgradeConfig.NetworkUpgrades != nil && upgradeConfig.NetworkUpgrades == nil {
		// Note: if we have previously applied persisted upgrade bytes,
		// missing "networkUpgrades" will be treated as intention to
		// abort forks. Initialize NetworkUpgrades here.
		upgradeConfig.NetworkUpgrades = &NetworkUpgrades{}
	}
	if networkUpgrades := upgradeConfig.NetworkUpgrades; networkUpgrades != nil {
		if err := c.getNetworkUpgrades().CheckCompatible(networkUpgrades, headTimestamp); err != nil {
			return err
		}
	}

	// Verify newly specified precompiles are compatible with the existing chainConfig.
	if err := c.CheckPrecompilesCompatible(upgradeConfig.PrecompileUpgrades, headTimestamp); err != nil {
		return err
	}

	// Apply upgrades to chainConfig.
	c.UpgradeConfig = upgradeConfig
	return nil
}
