// (c) 2022 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import "math/big"

// UpgradeableConfig contains the timestamp for the upgrade along with
// a boolean [Disable]. If [Disable] is set, the upgrade deactivates
// the precompile and resets its storage.
type UpgradeableConfig struct {
	BlockTimestamp *big.Int `json:"blockTimestamp"`
	Disable        bool     `json:"disable,omitempty"`
}

// Timestamp returns the timestamp this network upgrade goes into effect.
func (c *UpgradeableConfig) Timestamp() *big.Int {
	return c.BlockTimestamp
}

// IsDisabled returns true if the network upgrade deactivates the precompile.
func (c *UpgradeableConfig) IsDisabled() bool {
	return c.Disable
}
