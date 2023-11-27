// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"github.com/ava-labs/subnet-evm/utils"
)

var (
	LocalNetworkUpgrades = MandatoryNetworkUpgrades{
		SubnetEVMTimestamp: utils.NewUint64(0),
		DUpgradeTimestamp:  utils.NewUint64(0),
	}

	FujiNetworkUpgrades = MandatoryNetworkUpgrades{
		SubnetEVMTimestamp: utils.NewUint64(0),
		// DUpgradeTimestamp: utils.NewUint64(0), // TODO: Uncomment and set this to the correct value
	}

	MainnetNetworkUpgrades = MandatoryNetworkUpgrades{
		SubnetEVMTimestamp: utils.NewUint64(0),
		// DUpgradeTimestamp: utils.NewUint64(0), // TODO: Uncomment and set this to the correct value
	}

	UnitTestNetworkUpgrades = MandatoryNetworkUpgrades{
		SubnetEVMTimestamp: utils.NewUint64(0),
		DUpgradeTimestamp:  utils.NewUint64(0),
	}
)

// MandatoryNetworkUpgrades contains timestamps that enable mandatory network upgrades.
// These upgrades are mandatory, meaning that if a node does not upgrade by the
// specified timestamp, it will be unable to participate in consensus.
// Avalanche specific network upgrades are also included here.
type MandatoryNetworkUpgrades struct {
	// SubnetEVMTimestamp is a placeholder that activates Avalanche Upgrades prior to ApricotPhase6 (nil = no fork, 0 = already activated)
	SubnetEVMTimestamp *uint64 `json:"subnetEVMTimestamp,omitempty"`
	// DUpgrade activates the Shanghai Execution Spec Upgrade from Ethereum (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/shanghai.md#included-eips)
	// and Avalanche Warp Messaging. (nil = no fork, 0 = already activated)
	// Note: EIP-4895 is excluded since withdrawals are not relevant to the Avalanche C-Chain or Subnets running the EVM.
	DUpgradeTimestamp *uint64 `json:"dUpgradeTimestamp,omitempty"`
	// Cancun activates the Cancun upgrade from Ethereum. (nil = no fork, 0 = already activated)
	CancunTime *uint64 `json:"cancunTime,omitempty"`
}

func (m *MandatoryNetworkUpgrades) CheckMandatoryCompatible(newcfg *MandatoryNetworkUpgrades, time uint64) *ConfigCompatError {
	if isForkTimestampIncompatible(m.SubnetEVMTimestamp, newcfg.SubnetEVMTimestamp, time) {
		return newTimestampCompatError("SubnetEVM fork block timestamp", m.SubnetEVMTimestamp, newcfg.SubnetEVMTimestamp)
	}
	if isForkTimestampIncompatible(m.DUpgradeTimestamp, newcfg.DUpgradeTimestamp, time) {
		return newTimestampCompatError("DUpgrade fork block timestamp", m.DUpgradeTimestamp, newcfg.DUpgradeTimestamp)
	}
	if isForkTimestampIncompatible(m.CancunTime, newcfg.CancunTime, time) {
		return newTimestampCompatError("Cancun fork block timestamp", m.CancunTime, m.CancunTime)
	}
	return nil
}

func (m *MandatoryNetworkUpgrades) mandatoryForkOrder() []fork {
	return []fork{
		{name: "subnetEVMTimestamp", timestamp: m.SubnetEVMTimestamp},
		{name: "dUpgradeTimestamp", timestamp: m.DUpgradeTimestamp},
	}
}

type OptionalFork struct {
	Timestamp *uint64 `json:"timestamp" serialize:"true"`
}

func (m *OptionalFork) ToFork(name string) fork {
	var timestamp *uint64
	if m.Timestamp != nil {
		timestamp = m.Timestamp
	}
	return fork{
		name:      name,
		timestamp: timestamp,
		block:     nil,
		optional:  true,
	}
}

type OptionalNetworkUpgrades struct {
	// This is an example of a configuration.
	//FeatureConfig *OptionalFork `json:"test,omitempty" serialize:"true,nullable"`
}

func (n *OptionalNetworkUpgrades) CheckOptionalCompatible(newcfg *OptionalNetworkUpgrades, time uint64) *ConfigCompatError {
	return nil
}

func (n *OptionalNetworkUpgrades) optionalForkOrder() []fork {
	forks := make([]fork, 0)
	/*
		// This block of code should be added for each property inside OptionalNetworkUpgrades
		// Each property should have a test like https://github.com/ava-labs/subnet-evm/blob/c261ce270d16db5351b1e1ec5e128a3d23fe6a9c/params/config_test.go#L357-L360
		if n.FeatureConfig != nil {
			forks = append(forks, FeatureConfig.ToFork("featureConfig"))
		}

	*/
	return forks
}
