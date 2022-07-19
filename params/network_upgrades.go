// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/utils"
)

// NetworkUpgrades contains timestamps that enable avalanche network upgrades.
type NetworkUpgrades struct {
	SubnetEVMTimestamp *big.Int `json:"subnetEVMTimestamp,omitempty"` // A placeholder for the latest avalanche forks (nil = no fork, 0 = already activated)
}

func (c *NetworkUpgrades) CheckCompatible(newcfg *NetworkUpgrades, headHeight *big.Int, headTimestamp *big.Int) *utils.ConfigCompatError {
	// Check subnet-evm specific activations
	if isForkIncompatible(c.SubnetEVMTimestamp, newcfg.SubnetEVMTimestamp, headTimestamp) {
		return utils.NewCompatError("SubnetEVM fork block timestamp", c.SubnetEVMTimestamp, newcfg.SubnetEVMTimestamp)
	}

	return nil
}

// IsSubnetEVM returns whether [blockTimestamp] is either equal to the SubnetEVM fork block timestamp or greater.
func (n *NetworkUpgrades) IsSubnetEVM(blockTimestamp *big.Int) bool {
	return utils.IsForked(n.SubnetEVMTimestamp, blockTimestamp)
}
