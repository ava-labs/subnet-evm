// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

// NetworkUpgrades contains timestamps that enable avalanche network upgrades.
type NetworkUpgrades struct {
	// SubnetEVMTimestamp is a placeholder that activates Avalanche Upgrades prior to ApricotPhase6 (nil = no fork, 0 = already activated)
	SubnetEVMTimestamp *uint64 `json:"subnetEVMTimestamp,omitempty"`
	// DUpgrade activates the Shanghai upgrade from Ethereum. (nil = no fork, 0 = already activated)
	DUpgradeBlockTimestamp *uint64 `json:"dUpgradeBlockTimestamp,omitempty"`
}

func (n *NetworkUpgrades) CheckCompatible(newcfg *NetworkUpgrades, time uint64) *ConfigCompatError {
	// Check subnet-evm specific activations
	if isForkTimestampIncompatible(n.SubnetEVMTimestamp, newcfg.SubnetEVMTimestamp, time) {
		return newTimestampCompatError("SubnetEVM fork block timestamp", n.SubnetEVMTimestamp, newcfg.SubnetEVMTimestamp)
	}
	if isForkTimestampIncompatible(n.DUpgradeBlockTimestamp, newcfg.DUpgradeBlockTimestamp, time) {
		return newTimestampCompatError("DUpgrade fork block timestamp", n.DUpgradeBlockTimestamp, newcfg.DUpgradeBlockTimestamp)
	}

	return nil
}
