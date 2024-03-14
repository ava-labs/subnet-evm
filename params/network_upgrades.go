// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"fmt"
	"reflect"

	"github.com/ava-labs/avalanchego/version"
	"github.com/ava-labs/subnet-evm/utils"
)

// NetworkUpgrades contains timestamps that enable network upgrades.
// Avalanche specific network upgrades are also included here.
type NetworkUpgrades struct {
	// SubnetEVMTimestamp is a placeholder that activates Avalanche Upgrades prior to ApricotPhase6 (nil = no fork, 0 = already activated)
	SubnetEVMTimestamp *uint64 `json:"subnetEVMTimestamp,omitempty"`
	// Durango activates the Shanghai Execution Spec Upgrade from Ethereum (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/shanghai.md#included-eips)
	// and Avalanche Warp Messaging. (nil = no fork, 0 = already activated)
	// Note: EIP-4895 is excluded since withdrawals are not relevant to the Avalanche C-Chain or Subnets running the EVM.
	DurangoTimestamp *uint64 `json:"durangoTimestamp,omitempty"`
}

func (s *NetworkUpgrades) Equal(other *NetworkUpgrades) bool {
	return reflect.DeepEqual(s, other)
}

func (m *NetworkUpgrades) CheckNetworkUpgradesCompatible(newcfg *NetworkUpgrades, time uint64) *ConfigCompatError {
	if isForkTimestampIncompatible(m.SubnetEVMTimestamp, newcfg.SubnetEVMTimestamp, time) {
		return newTimestampCompatError("SubnetEVM fork block timestamp", m.SubnetEVMTimestamp, newcfg.SubnetEVMTimestamp)
	}
	if isForkTimestampIncompatible(m.DurangoTimestamp, newcfg.DurangoTimestamp, time) {
		return newTimestampCompatError("Durango fork block timestamp", m.DurangoTimestamp, newcfg.DurangoTimestamp)
	}
	return nil
}

func (m *NetworkUpgrades) forkOrder() []fork {
	return []fork{
		{name: "subnetEVMTimestamp", timestamp: m.SubnetEVMTimestamp},
		{name: "durangoTimestamp", timestamp: m.DurangoTimestamp},
	}
}

// setDefaults sets the default values for the network upgrades.
// This overrides deactivating the network upgrade by providing a timestamp of nil value.
func (m *NetworkUpgrades) setDefaults(networkID uint32) {
	defaults := getDefaultNetworkUpgrades(networkID)
	if m.SubnetEVMTimestamp == nil {
		m.SubnetEVMTimestamp = defaults.SubnetEVMTimestamp
	}
	if m.DurangoTimestamp == nil {
		m.DurangoTimestamp = defaults.DurangoTimestamp
	}
}

// verify checks that the network upgrades are well formed.
func (m *NetworkUpgrades) VerifyNetworkUpgrades(networkID uint32) error {
	defaults := getDefaultNetworkUpgrades(networkID)
	if isNilOrSmaller(m.SubnetEVMTimestamp, *defaults.SubnetEVMTimestamp) {
		return fmt.Errorf("SubnetEVM fork block timestamp (%v) must be greater than or equal to %v", m.SubnetEVMTimestamp, *defaults.SubnetEVMTimestamp)
	}
	if isNilOrSmaller(m.DurangoTimestamp, *defaults.DurangoTimestamp) {
		return fmt.Errorf("Durango fork block timestamp (%v) must be greater than or equal to %v", m.DurangoTimestamp, *defaults.DurangoTimestamp)
	}
	return nil
}

func (m *NetworkUpgrades) Override(o *NetworkUpgrades) {
	if o.SubnetEVMTimestamp != nil {
		m.SubnetEVMTimestamp = o.SubnetEVMTimestamp
	}
	if o.DurangoTimestamp != nil {
		m.DurangoTimestamp = o.DurangoTimestamp
	}
}

// getDefaultNetworkUpgrades returns the network upgrades for the specified network ID.
func getDefaultNetworkUpgrades(networkID uint32) NetworkUpgrades {
	return NetworkUpgrades{
		SubnetEVMTimestamp: utils.NewUint64(0),
		DurangoTimestamp:   getUpgradeTime(networkID, version.DurangoTimes),
	}
}

func isNilOrSmaller(a *uint64, b uint64) bool {
	if a == nil {
		return true
	}
	return *a < b
}
