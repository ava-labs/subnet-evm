// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"math/big"

	"github.com/ava-labs/avalanchego/utils/constants"
)

var (
	LocalNetworkUpgrades = NetworkUpgrades{
		SubnetEVMTimestamp: big.NewInt(0),
		DUpgradeTimestamp:  big.NewInt(0),
	}

	FujiNetworkUpgrades = NetworkUpgrades{
		SubnetEVMTimestamp: big.NewInt(0),
		//	DUpgradeTimestamp:           big.NewInt(0), // TODO: Uncomment and set this to the correct value
	}

	MainnetNetworkUpgrades = NetworkUpgrades{
		SubnetEVMTimestamp: big.NewInt(0),
		//	DUpgradeTimestamp:           big.NewInt(0), // TODO: Uncomment and set this to the correct value
	}
)

// NetworkUpgrades contains timestamps that enable avalanche network upgrades.
type NetworkUpgrades struct {
	SubnetEVMTimestamp *big.Int `json:"subnetEVMTimestamp,omitempty"` // initial subnet-evm upgrade (nil = no fork, 0 = already activated)
	DUpgradeTimestamp  *big.Int `json:"dUpgradeTimestamp,omitempty"`  // A placeholder for the latest avalanche forks (nil = no fork, 0 = already activated)
}

func (n *NetworkUpgrades) CheckCompatible(newcfg NetworkUpgrades, headTimestamp *big.Int) *ConfigCompatError {
	// Check subnet-evm specific activations
	if isForkIncompatible(n.SubnetEVMTimestamp, newcfg.SubnetEVMTimestamp, headTimestamp) {
		return newCompatError("SubnetEVM fork block timestamp", n.SubnetEVMTimestamp, newcfg.SubnetEVMTimestamp)
	}
	if isForkIncompatible(n.DUpgradeTimestamp, newcfg.DUpgradeTimestamp, headTimestamp) {
		return newCompatError("DUpgrade fork block timestamp", n.DUpgradeTimestamp, newcfg.DUpgradeTimestamp)
	}

	return nil
}

func GetNetworkUpgrades(networkID uint32) NetworkUpgrades {
	switch networkID {
	case constants.FujiID:
		return FujiNetworkUpgrades
	case constants.MainnetID:
		return MainnetNetworkUpgrades
	default:
		return LocalNetworkUpgrades
	}
}
