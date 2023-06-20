// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"math/big"
)

// NetworkUpgrades contains timestamps that enable avalanche network upgrades.
type NetworkUpgrades struct {
	SubnetEVMTimestamp *uint64 `json:"subnetEVMTimestamp,omitempty"` // A placeholder for the latest avalanche forks (nil = no fork, 0 = already activated)
}

func (n *NetworkUpgrades) CheckCompatible(newcfg *NetworkUpgrades, headTimestamp *big.Int) *ConfigCompatError {
	// Check subnet-evm specific activations
	if isForkTimestampIncompatible(n.SubnetEVMTimestamp, newcfg.SubnetEVMTimestamp, headTimestamp) {
		return newTimestampCompatError("SubnetEVM fork block timestamp", n.SubnetEVMTimestamp, newcfg.SubnetEVMTimestamp)
	}

	return nil
}
