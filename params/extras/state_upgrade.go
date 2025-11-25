// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package extras

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/stateupgrade"

	ethparams "github.com/ava-labs/libevm/params"
)

// GetActivatingStateUpgrades returns all state upgrades configured to activate during the
// state transition from a block with timestamp [from] to a block with timestamp [to].
func (*ChainConfig) GetActivatingStateUpgrades(from *uint64, to uint64, upgrades []stateupgrade.StateUpgrade) []stateupgrade.StateUpgrade {
	activating := make([]stateupgrade.StateUpgrade, 0)
	for _, upgrade := range upgrades {
		if IsForkTransition(upgrade.BlockTimestamp, from, to) {
			activating = append(activating, upgrade)
		}
	}
	return activating
}

// checkStateUpgradesCompatible checks if [stateUpgrades] are compatible with [c] at [headTimestamp].
func (c *ChainConfig) checkStateUpgradesCompatible(stateUpgrades []stateupgrade.StateUpgrade, lastTimestamp uint64) *ethparams.ConfigCompatError {
	// All active upgrades (from nil to [lastTimestamp]) must match.
	activeUpgrades := c.GetActivatingStateUpgrades(nil, lastTimestamp, c.StateUpgrades)
	newUpgrades := c.GetActivatingStateUpgrades(nil, lastTimestamp, stateUpgrades)

	// Check activated upgrades are still present.
	for i, upgrade := range activeUpgrades {
		if len(newUpgrades) <= i {
			// missing upgrade
			return ethparams.NewTimestampCompatError(
				fmt.Sprintf("missing StateUpgrade[%d]", i),
				upgrade.BlockTimestamp,
				nil,
			)
		}
		// All upgrades that have activated must be identical.
		if !upgrade.Equal(&newUpgrades[i]) {
			return ethparams.NewTimestampCompatError(
				fmt.Sprintf("StateUpgrade[%d]", i),
				upgrade.BlockTimestamp,
				newUpgrades[i].BlockTimestamp,
			)
		}
	}
	// then, make sure newUpgrades does not have additional upgrades
	// that are already activated. (cannot perform retroactive upgrade)
	if len(newUpgrades) > len(activeUpgrades) {
		return ethparams.NewTimestampCompatError(
			fmt.Sprintf("cannot retroactively enable StateUpgrade[%d]", len(activeUpgrades)),
			nil,
			newUpgrades[len(activeUpgrades)].BlockTimestamp, // this indexes to the first element in newUpgrades after the end of activeUpgrades
		)
	}

	return nil
}
