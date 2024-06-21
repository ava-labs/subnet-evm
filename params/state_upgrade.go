// (c) 2023 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"fmt"
	"reflect"

	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
)

// StateUpgrade describes the modifications to be made to the state during
// a state upgrade.
type StateUpgrade struct {
	BlockTimestamp *uint64 `json:"blockTimestamp,omitempty"`

	// map from account address to the modification to be made to the account.
	StateUpgradeAccounts map[common.Address]StateUpgradeAccount `json:"accounts"`
}

// StateUpgradeAccount describes the modifications to be made to an account during
// a state upgrade.
type StateUpgradeAccount struct {
	Code          hexutil.Bytes               `json:"code,omitempty"`
	Storage       map[common.Hash]common.Hash `json:"storage,omitempty"`
	BalanceChange *math.HexOrDecimal256       `json:"balanceChange,omitempty"`
}

func (s *StateUpgrade) Equal(other *StateUpgrade) bool {
	return reflect.DeepEqual(s, other)
}

// verifyStateUpgrades checks [c.StateUpgrades] is well formed:
// - the specified blockTimestamps must monotonically increase
func (c *ChainConfig) verifyStateUpgrades() error {
	var previousUpgradeTimestamp *uint64
	for i, upgrade := range c.StateUpgrades {
		upgradeTimestamp := upgrade.BlockTimestamp
		if upgradeTimestamp == nil {
			return fmt.Errorf("StateUpgrade[%d]: config block timestamp cannot be nil ", i)
		}
		// Verify the upgrade's timestamp is equal 0 (to avoid confusion with genesis).
		if *upgradeTimestamp == 0 {
			return fmt.Errorf("StateUpgrade[%d]: config block timestamp (%v) must be greater than 0", i, *upgradeTimestamp)
		}

		// Verify specified timestamps are strictly monotonically increasing.
		if previousUpgradeTimestamp != nil && *upgradeTimestamp <= *previousUpgradeTimestamp {
			return fmt.Errorf("StateUpgrade[%d]: config block timestamp (%v) <= previous timestamp (%v)", i, *upgradeTimestamp, *previousUpgradeTimestamp)
		}
		previousUpgradeTimestamp = upgradeTimestamp
	}
	return nil
}

// GetActivatingStateUpgrades returns all state upgrades configured to activate during the
// state transition from a block with timestamp [from] to a block with timestamp [to].
func (c *ChainConfig) GetActivatingStateUpgrades(from *uint64, to uint64, upgrades []StateUpgrade) []StateUpgrade {
	activating := make([]StateUpgrade, 0)
	for _, upgrade := range upgrades {
		if utils.IsForkTransition(upgrade.BlockTimestamp, from, to) {
			activating = append(activating, upgrade)
		}
	}
	return activating
}

// CheckStateUpgradesCompatible checks if [stateUpgrades] are compatible with [c] at [headTimestamp].
func (c *ChainConfig) CheckStateUpgradesCompatible(stateUpgrades []StateUpgrade, lastTimestamp uint64) *ConfigCompatError {
	// All active upgrades (from nil to [lastTimestamp]) must match.
	activeUpgrades := c.GetActivatingStateUpgrades(nil, lastTimestamp, c.StateUpgrades)
	newUpgrades := c.GetActivatingStateUpgrades(nil, lastTimestamp, stateUpgrades)

	// Check activated upgrades are still present.
	for i, upgrade := range activeUpgrades {
		if len(newUpgrades) <= i {
			// missing upgrade
			return newTimestampCompatError(
				fmt.Sprintf("missing StateUpgrade[%d]", i),
				upgrade.BlockTimestamp,
				nil,
			)
		}
		// All upgrades that have activated must be identical.
		if !upgrade.Equal(&newUpgrades[i]) {
			return newTimestampCompatError(
				fmt.Sprintf("StateUpgrade[%d]", i),
				upgrade.BlockTimestamp,
				newUpgrades[i].BlockTimestamp,
			)
		}
	}
	// then, make sure newUpgrades does not have additional upgrades
	// that are already activated. (cannot perform retroactive upgrade)
	if len(newUpgrades) > len(activeUpgrades) {
		return newTimestampCompatError(
			fmt.Sprintf("cannot retroactively enable StateUpgrade[%d]", len(activeUpgrades)),
			nil,
			newUpgrades[len(activeUpgrades)].BlockTimestamp, // this indexes to the first element in newUpgrades after the end of activeUpgrades
		)
	}

	return nil
}
