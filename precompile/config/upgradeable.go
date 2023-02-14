// (c) 2022 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package config

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/utils"
)

// Uprade contains the timestamp for the upgrade along with
// a boolean [Disable]. If [Disable] is set, the upgrade deactivates
// the precompile and resets its storage.
type Uprade struct {
	BlockTimestamp *big.Int `json:"blockTimestamp"`
	Disable        bool     `json:"disable,omitempty"`
}

// Timestamp returns the timestamp this network upgrade goes into effect.
func (c *Uprade) Timestamp() *big.Int {
	return c.BlockTimestamp
}

// IsDisabled returns true if the network upgrade deactivates the precompile.
func (c *Uprade) IsDisabled() bool {
	return c.Disable
}

// Equal returns true iff [other] has the same blockTimestamp and has the
// same on value for the Disable flag.
func (c *Uprade) Equal(other *Uprade) bool {
	if other == nil {
		return false
	}
	return c.Disable == other.Disable && utils.BigNumEqual(c.BlockTimestamp, other.BlockTimestamp)
}
