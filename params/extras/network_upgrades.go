// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package extras

import (
	"fmt"
	"reflect"

	"github.com/ava-labs/avalanchego/upgrade"
	ethparams "github.com/ava-labs/libevm/params"
	"github.com/ava-labs/subnet-evm/utils"
)

var errCannotBeNil = fmt.Errorf("timestamp cannot be nil")

// NetworkUpgrades contains timestamps that enable network upgrades.
// Avalanche specific network upgrades are also included here.
// (nil = no fork, 0 = already activated)
type NetworkUpgrades struct {
	ApricotPhase1BlockTimestamp *uint64 `json:"apricotPhase1BlockTimestamp,omitempty"` // Apricot Phase 1 Block Timestamp (nil = no fork, 0 = already activated)
	// Apricot Phase 2 Block Timestamp (nil = no fork, 0 = already activated)
	// Apricot Phase 2 includes a modified version of the Berlin Hard Fork from Ethereum
	ApricotPhase2BlockTimestamp *uint64 `json:"apricotPhase2BlockTimestamp,omitempty"`
	// Apricot Phase 3 introduces dynamic fees and a modified version of the London Hard Fork from Ethereum (nil = no fork, 0 = already activated)
	ApricotPhase3BlockTimestamp *uint64 `json:"apricotPhase3BlockTimestamp,omitempty"`
	// Apricot Phase 4 introduces the notion of a block fee to the dynamic fee algorithm (nil = no fork, 0 = already activated)
	ApricotPhase4BlockTimestamp *uint64 `json:"apricotPhase4BlockTimestamp,omitempty"`
	// Apricot Phase 5 introduces a batch of atomic transactions with a maximum atomic gas limit per block. (nil = no fork, 0 = already activated)
	ApricotPhase5BlockTimestamp *uint64 `json:"apricotPhase5BlockTimestamp,omitempty"`
	// Apricot Phase Pre-6 deprecates the NativeAssetCall precompile (soft). (nil = no fork, 0 = already activated)
	ApricotPhasePre6BlockTimestamp *uint64 `json:"apricotPhasePre6BlockTimestamp,omitempty"`
	// Apricot Phase 6 deprecates the NativeAssetBalance and NativeAssetCall precompiles. (nil = no fork, 0 = already activated)
	ApricotPhase6BlockTimestamp *uint64 `json:"apricotPhase6BlockTimestamp,omitempty"`
	// Apricot Phase Post-6 deprecates the NativeAssetCall precompile (soft). (nil = no fork, 0 = already activated)
	ApricotPhasePost6BlockTimestamp *uint64 `json:"apricotPhasePost6BlockTimestamp,omitempty"`
	// Banff restricts import/export transactions to AVAX. (nil = no fork, 0 = already activated)
	BanffBlockTimestamp *uint64 `json:"banffBlockTimestamp,omitempty"`
	// Cortina increases the block gas limit to 15M. (nil = no fork, 0 = already activated)
	CortinaBlockTimestamp *uint64 `json:"cortinaBlockTimestamp,omitempty"`
	// SubnetEVMTimestamp is a placeholder that activates Avalanche Upgrades prior to Durango
	SubnetEVMTimestamp *uint64 `json:"subnetEVMTimestamp,omitempty"`
	// Durango activates the Shanghai Execution Spec Upgrade from Ethereum (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/shanghai.md#included-eips)
	// and Avalanche Warp Messaging. (nil = no fork, 0 = already activated)
	// Note: EIP-4895 is excluded since withdrawals are not relevant to the Avalanche C-Chain or Subnets running the EVM.
	DurangoTimestamp *uint64 `json:"durangoTimestamp,omitempty"`
	// Etna activates Cancun from Ethereum (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/cancun.md#included-eips)
	// and reduces the min base fee. (nil = no fork, 0 = already activated)
	// Note: EIP-4844 BlobTxs are not enabled in the mempool and blocks are not
	// allowed to contain them. For details see https://github.com/avalanche-foundation/ACPs/pull/131
	EtnaTimestamp *uint64 `json:"etnaTimestamp,omitempty"`
	// Fortuna is a placeholder for the next upgrade.
	// (nil = no fork, 0 = already activated)
	FortunaTimestamp *uint64 `json:"fortunaTimestamp,omitempty"`
}

func (n *NetworkUpgrades) Equal(other *NetworkUpgrades) bool {
	return reflect.DeepEqual(n, other)
}

func (n *NetworkUpgrades) checkNetworkUpgradesCompatible(newcfg *NetworkUpgrades, time uint64) *ethparams.ConfigCompatError {
	if isForkTimestampIncompatible(n.ApricotPhase1BlockTimestamp, newcfg.ApricotPhase1BlockTimestamp, time) {
		return ethparams.NewTimestampCompatError("ApricotPhase1 fork block timestamp", n.ApricotPhase1BlockTimestamp, newcfg.ApricotPhase1BlockTimestamp)
	}
	if isForkTimestampIncompatible(n.ApricotPhase2BlockTimestamp, newcfg.ApricotPhase2BlockTimestamp, time) {
		return ethparams.NewTimestampCompatError("ApricotPhase2 fork block timestamp", n.ApricotPhase2BlockTimestamp, newcfg.ApricotPhase2BlockTimestamp)
	}
	if isForkTimestampIncompatible(n.ApricotPhase3BlockTimestamp, newcfg.ApricotPhase3BlockTimestamp, time) {
		return ethparams.NewTimestampCompatError("ApricotPhase3 fork block timestamp", n.ApricotPhase3BlockTimestamp, newcfg.ApricotPhase3BlockTimestamp)
	}
	if isForkTimestampIncompatible(n.ApricotPhase4BlockTimestamp, newcfg.ApricotPhase4BlockTimestamp, time) {
		return ethparams.NewTimestampCompatError("ApricotPhase4 fork block timestamp", n.ApricotPhase4BlockTimestamp, newcfg.ApricotPhase4BlockTimestamp)
	}
	if isForkTimestampIncompatible(n.ApricotPhase5BlockTimestamp, newcfg.ApricotPhase5BlockTimestamp, time) {
		return ethparams.NewTimestampCompatError("ApricotPhase5 fork block timestamp", n.ApricotPhase5BlockTimestamp, newcfg.ApricotPhase5BlockTimestamp)
	}
	if isForkTimestampIncompatible(n.ApricotPhasePre6BlockTimestamp, newcfg.ApricotPhasePre6BlockTimestamp, time) {
		return ethparams.NewTimestampCompatError("ApricotPhasePre6 fork block timestamp", n.ApricotPhasePre6BlockTimestamp, newcfg.ApricotPhasePre6BlockTimestamp)
	}
	if isForkTimestampIncompatible(n.ApricotPhase6BlockTimestamp, newcfg.ApricotPhase6BlockTimestamp, time) {
		return ethparams.NewTimestampCompatError("ApricotPhase6 fork block timestamp", n.ApricotPhase6BlockTimestamp, newcfg.ApricotPhase6BlockTimestamp)
	}
	if isForkTimestampIncompatible(n.ApricotPhasePost6BlockTimestamp, newcfg.ApricotPhasePost6BlockTimestamp, time) {
		return ethparams.NewTimestampCompatError("ApricotPhasePost6 fork block timestamp", n.ApricotPhasePost6BlockTimestamp, newcfg.ApricotPhasePost6BlockTimestamp)
	}
	if isForkTimestampIncompatible(n.BanffBlockTimestamp, newcfg.BanffBlockTimestamp, time) {
		return ethparams.NewTimestampCompatError("Banff fork block timestamp", n.BanffBlockTimestamp, newcfg.BanffBlockTimestamp)
	}
	if isForkTimestampIncompatible(n.CortinaBlockTimestamp, newcfg.CortinaBlockTimestamp, time) {
		return ethparams.NewTimestampCompatError("Cortina fork block timestamp", n.CortinaBlockTimestamp, newcfg.CortinaBlockTimestamp)
	}
	if isForkTimestampIncompatible(n.SubnetEVMTimestamp, newcfg.SubnetEVMTimestamp, time) {
		return ethparams.NewTimestampCompatError("SubnetEVM fork block timestamp", n.SubnetEVMTimestamp, newcfg.SubnetEVMTimestamp)
	}
	if isForkTimestampIncompatible(n.DurangoTimestamp, newcfg.DurangoTimestamp, time) {
		return ethparams.NewTimestampCompatError("Durango fork block timestamp", n.DurangoTimestamp, newcfg.DurangoTimestamp)
	}
	if isForkTimestampIncompatible(n.EtnaTimestamp, newcfg.EtnaTimestamp, time) {
		return ethparams.NewTimestampCompatError("Etna fork block timestamp", n.EtnaTimestamp, newcfg.EtnaTimestamp)
	}
	if isForkTimestampIncompatible(n.FortunaTimestamp, newcfg.FortunaTimestamp, time) {
		return ethparams.NewTimestampCompatError("Fortuna fork block timestamp", n.FortunaTimestamp, newcfg.FortunaTimestamp)
	}

	return nil
}

func (n *NetworkUpgrades) forkOrder() []fork {
	return []fork{
		{name: "apricotPhase1BlockTimestamp", timestamp: n.ApricotPhase1BlockTimestamp},
		{name: "apricotPhase2BlockTimestamp", timestamp: n.ApricotPhase2BlockTimestamp},
		{name: "apricotPhase3BlockTimestamp", timestamp: n.ApricotPhase3BlockTimestamp},
		{name: "apricotPhase4BlockTimestamp", timestamp: n.ApricotPhase4BlockTimestamp},
		{name: "apricotPhase5BlockTimestamp", timestamp: n.ApricotPhase5BlockTimestamp},
		{name: "apricotPhasePre6BlockTimestamp", timestamp: n.ApricotPhasePre6BlockTimestamp},
		{name: "apricotPhase6BlockTimestamp", timestamp: n.ApricotPhase6BlockTimestamp},
		{name: "apricotPhasePost6BlockTimestamp", timestamp: n.ApricotPhasePost6BlockTimestamp},
		{name: "banffBlockTimestamp", timestamp: n.BanffBlockTimestamp},
		{name: "cortinaBlockTimestamp", timestamp: n.CortinaBlockTimestamp},
		{name: "subnetEVMTimestamp", timestamp: n.SubnetEVMTimestamp, optional: true},
		{name: "durangoTimestamp", timestamp: n.DurangoTimestamp},
		{name: "etnaTimestamp", timestamp: n.EtnaTimestamp},
		{name: "fortunaTimestamp", timestamp: n.FortunaTimestamp, optional: true},
	}
}

// IsApricotPhase1 returns whether [time] represents a block
// with a timestamp after the Apricot Phase 1 upgrade time.
func (n NetworkUpgrades) IsApricotPhase1(time uint64) bool {
	return isTimestampForked(n.ApricotPhase1BlockTimestamp, time)
}

// IsApricotPhase2 returns whether [time] represents a block
// with a timestamp after the Apricot Phase 2 upgrade time.
func (n NetworkUpgrades) IsApricotPhase2(time uint64) bool {
	return isTimestampForked(n.ApricotPhase2BlockTimestamp, time)
}

// IsApricotPhase3 returns whether [time] represents a block
// with a timestamp after the Apricot Phase 3 upgrade time.
func (n *NetworkUpgrades) IsApricotPhase3(time uint64) bool {
	return isTimestampForked(n.ApricotPhase3BlockTimestamp, time)
}

// IsApricotPhase4 returns whether [time] represents a block
// with a timestamp after the Apricot Phase 4 upgrade time.
func (n NetworkUpgrades) IsApricotPhase4(time uint64) bool {
	return isTimestampForked(n.ApricotPhase4BlockTimestamp, time)
}

// IsApricotPhase5 returns whether [time] represents a block
// with a timestamp after the Apricot Phase 5 upgrade time.
func (n NetworkUpgrades) IsApricotPhase5(time uint64) bool {
	return isTimestampForked(n.ApricotPhase5BlockTimestamp, time)
}

// IsApricotPhasePre6 returns whether [time] represents a block
// with a timestamp after the Apricot Phase Pre 6 upgrade time.
func (n NetworkUpgrades) IsApricotPhasePre6(time uint64) bool {
	return isTimestampForked(n.ApricotPhasePre6BlockTimestamp, time)
}

// IsApricotPhase6 returns whether [time] represents a block
// with a timestamp after the Apricot Phase 6 upgrade time.
func (n NetworkUpgrades) IsApricotPhase6(time uint64) bool {
	return isTimestampForked(n.ApricotPhase6BlockTimestamp, time)
}

// IsApricotPhasePost6 returns whether [time] represents a block
// with a timestamp after the Apricot Phase 6 Post upgrade time.
func (n NetworkUpgrades) IsApricotPhasePost6(time uint64) bool {
	return isTimestampForked(n.ApricotPhasePost6BlockTimestamp, time)
}

// IsBanff returns whether [time] represents a block
// with a timestamp after the Banff upgrade time.
func (n NetworkUpgrades) IsBanff(time uint64) bool {
	return isTimestampForked(n.BanffBlockTimestamp, time)
}

// IsCortina returns whether [time] represents a block
// with a timestamp after the Cortina upgrade time.
func (n NetworkUpgrades) IsCortina(time uint64) bool {
	return isTimestampForked(n.CortinaBlockTimestamp, time)
}

// IsSubnetEVM returns whether [time] represents a block
// with a timestamp after the SubnetEVM upgrade time.
func (n NetworkUpgrades) IsSubnetEVM(time uint64) bool {
	return isTimestampForked(n.SubnetEVMTimestamp, time)
}

// IsDurango returns whether [time] represents a block
// with a timestamp after the Durango upgrade time.
func (n NetworkUpgrades) IsDurango(time uint64) bool {
	return isTimestampForked(n.DurangoTimestamp, time)
}

// IsEtna returns whether [time] represents a block
// with a timestamp after the Etna upgrade time.
func (n NetworkUpgrades) IsEtna(time uint64) bool {
	return isTimestampForked(n.EtnaTimestamp, time)
}

// IsFortuna returns whether [time] represents a block
// with a timestamp after the Fortuna upgrade time.
func (n *NetworkUpgrades) IsFortuna(time uint64) bool {
	return isTimestampForked(n.FortunaTimestamp, time)
}

func (n *NetworkUpgrades) Description() string {
	var banner string
	banner += fmt.Sprintf(" - Apricot Phase 1 Timestamp:        @%-10v (https://github.com/ava-labs/avalanchego/releases/tag/v1.3.0)\n", ptrToString(n.ApricotPhase1BlockTimestamp))
	banner += fmt.Sprintf(" - Apricot Phase 2 Timestamp:        @%-10v (https://github.com/ava-labs/avalanchego/releases/tag/v1.4.0)\n", ptrToString(n.ApricotPhase2BlockTimestamp))
	banner += fmt.Sprintf(" - Apricot Phase 3 Timestamp:        @%-10v (https://github.com/ava-labs/avalanchego/releases/tag/v1.5.0)\n", ptrToString(n.ApricotPhase3BlockTimestamp))
	banner += fmt.Sprintf(" - Apricot Phase 4 Timestamp:        @%-10v (https://github.com/ava-labs/avalanchego/releases/tag/v1.6.0)\n", ptrToString(n.ApricotPhase4BlockTimestamp))
	banner += fmt.Sprintf(" - Apricot Phase 5 Timestamp:        @%-10v (https://github.com/ava-labs/avalanchego/releases/tag/v1.7.0)\n", ptrToString(n.ApricotPhase5BlockTimestamp))
	banner += fmt.Sprintf(" - Apricot Phase P6 Timestamp        @%-10v (https://github.com/ava-labs/avalanchego/releases/tag/v1.8.0)\n", ptrToString(n.ApricotPhasePre6BlockTimestamp))
	banner += fmt.Sprintf(" - Apricot Phase 6 Timestamp:        @%-10v (https://github.com/ava-labs/avalanchego/releases/tag/v1.8.0)\n", ptrToString(n.ApricotPhase6BlockTimestamp))
	banner += fmt.Sprintf(" - Apricot Phase Post-6 Timestamp:   @%-10v (https://github.com/ava-labs/avalanchego/releases/tag/v1.8.0\n", ptrToString(n.ApricotPhasePost6BlockTimestamp))
	banner += fmt.Sprintf(" - Banff Timestamp:                  @%-10v (https://github.com/ava-labs/avalanchego/releases/tag/v1.9.0)\n", ptrToString(n.BanffBlockTimestamp))
	banner += fmt.Sprintf(" - Cortina Timestamp:                @%-10v (https://github.com/ava-labs/avalanchego/releases/tag/v1.10.0)\n", ptrToString(n.CortinaBlockTimestamp))
	banner += fmt.Sprintf(" - SubnetEVM Timestamp:          @%-10v (https://github.com/ava-labs/avalanchego/releases/tag/v1.10.0)\n", ptrToString(n.SubnetEVMTimestamp))
	banner += fmt.Sprintf(" - Durango Timestamp:            @%-10v (https://github.com/ava-labs/avalanchego/releases/tag/v1.11.0)\n", ptrToString(n.DurangoTimestamp))
	banner += fmt.Sprintf(" - Etna Timestamp:               @%-10v (https://github.com/ava-labs/avalanchego/releases/tag/v1.12.0)\n", ptrToString(n.EtnaTimestamp))
	banner += fmt.Sprintf(" - Fortuna Timestamp:            @%-10v (https://github.com/ava-labs/avalanchego/releases/tag/v1.13.0)\n", ptrToString(n.FortunaTimestamp))
	return banner
}

func CorethDefaultNetworkUpgrades(agoUpgrade upgrade.Config) NetworkUpgrades {
	return NetworkUpgrades{
		ApricotPhase1BlockTimestamp:     utils.TimeToNewUint64(agoUpgrade.ApricotPhase1Time),
		ApricotPhase2BlockTimestamp:     utils.TimeToNewUint64(agoUpgrade.ApricotPhase2Time),
		ApricotPhase3BlockTimestamp:     utils.TimeToNewUint64(agoUpgrade.ApricotPhase3Time),
		ApricotPhase4BlockTimestamp:     utils.TimeToNewUint64(agoUpgrade.ApricotPhase4Time),
		ApricotPhase5BlockTimestamp:     utils.TimeToNewUint64(agoUpgrade.ApricotPhase5Time),
		ApricotPhasePre6BlockTimestamp:  utils.TimeToNewUint64(agoUpgrade.ApricotPhasePre6Time),
		ApricotPhase6BlockTimestamp:     utils.TimeToNewUint64(agoUpgrade.ApricotPhase6Time),
		ApricotPhasePost6BlockTimestamp: utils.TimeToNewUint64(agoUpgrade.ApricotPhasePost6Time),
		BanffBlockTimestamp:             utils.TimeToNewUint64(agoUpgrade.BanffTime),
		CortinaBlockTimestamp:           utils.TimeToNewUint64(agoUpgrade.CortinaTime),
		DurangoTimestamp:                utils.TimeToNewUint64(agoUpgrade.DurangoTime),
		EtnaTimestamp:                   utils.TimeToNewUint64(agoUpgrade.EtnaTime),
		FortunaTimestamp:                utils.TimeToNewUint64(agoUpgrade.FortunaTime),
	}
}

func SubnetEVMDefaultNetworkUpgrades(agoUpgrade upgrade.Config) NetworkUpgrades {
	defaultSubnetEVMTimestamp := uint64(0)
	return NetworkUpgrades{
		// All these timetamps are required to be set to SubnetEVM default timestamp
		ApricotPhase1BlockTimestamp:     utils.NewUint64(defaultSubnetEVMTimestamp),
		ApricotPhase2BlockTimestamp:     utils.NewUint64(defaultSubnetEVMTimestamp),
		ApricotPhase3BlockTimestamp:     utils.NewUint64(defaultSubnetEVMTimestamp),
		ApricotPhase4BlockTimestamp:     utils.NewUint64(defaultSubnetEVMTimestamp),
		ApricotPhase5BlockTimestamp:     utils.NewUint64(defaultSubnetEVMTimestamp),
		ApricotPhasePre6BlockTimestamp:  utils.NewUint64(defaultSubnetEVMTimestamp),
		ApricotPhase6BlockTimestamp:     utils.NewUint64(defaultSubnetEVMTimestamp),
		ApricotPhasePost6BlockTimestamp: utils.NewUint64(defaultSubnetEVMTimestamp),
		BanffBlockTimestamp:             utils.NewUint64(defaultSubnetEVMTimestamp),
		CortinaBlockTimestamp:           utils.NewUint64(defaultSubnetEVMTimestamp),
		SubnetEVMTimestamp:              utils.NewUint64(defaultSubnetEVMTimestamp),
		DurangoTimestamp:                utils.TimeToNewUint64(agoUpgrade.DurangoTime),
		EtnaTimestamp:                   utils.TimeToNewUint64(agoUpgrade.EtnaTime),
		FortunaTimestamp:                nil, // Fortuna is optional and has no effect on Subnet-EVM
	}
}

// SetSubnetEVMDefaults sets the default values for the network upgrades.
// This overrides deactivating the network upgrade by providing a timestamp of nil value.
// If the network upgrade is not set, sets it to the default value.
// If the network upgrade is set to 0, we also treat it as nil and set it default.
// This is because in prior versions, upgrades were not modifiable and were directly set to their default values.
// Most of the tools and configurations just provide these as 0, so it is safer to treat 0 as nil and set to default
// to prevent premature activations of the network upgrades for live networks.
func SetSubnetEVMDefaults(n *NetworkUpgrades, agoUpgrade upgrade.Config) {
	defaults := SubnetEVMDefaultNetworkUpgrades(agoUpgrade)

	n.ApricotPhase1BlockTimestamp = defaults.ApricotPhase1BlockTimestamp
	n.ApricotPhase2BlockTimestamp = defaults.ApricotPhase2BlockTimestamp
	n.ApricotPhase3BlockTimestamp = defaults.ApricotPhase3BlockTimestamp
	n.ApricotPhase4BlockTimestamp = defaults.ApricotPhase4BlockTimestamp
	n.ApricotPhase5BlockTimestamp = defaults.ApricotPhase5BlockTimestamp
	n.ApricotPhasePre6BlockTimestamp = defaults.ApricotPhasePre6BlockTimestamp
	n.ApricotPhase6BlockTimestamp = defaults.ApricotPhase6BlockTimestamp
	n.ApricotPhasePost6BlockTimestamp = defaults.ApricotPhasePost6BlockTimestamp
	n.BanffBlockTimestamp = defaults.BanffBlockTimestamp
	n.CortinaBlockTimestamp = defaults.CortinaBlockTimestamp
	n.SubnetEVMTimestamp = defaults.SubnetEVMTimestamp

	if n.DurangoTimestamp == nil || *n.DurangoTimestamp == 0 {
		n.DurangoTimestamp = defaults.DurangoTimestamp
	}
	if n.EtnaTimestamp == nil || *n.EtnaTimestamp == 0 {
		n.EtnaTimestamp = defaults.EtnaTimestamp
	}
	if n.FortunaTimestamp == nil || *n.FortunaTimestamp == 0 {
		n.FortunaTimestamp = defaults.FortunaTimestamp
	}
}

// verifyNetworkUpgrades checks that the network upgrades are well formed.
func (n *NetworkUpgrades) verifyNetworkUpgrades(agoUpgrades upgrade.Config) error {
	var defaults NetworkUpgrades
	if n.SubnetEVMTimestamp != nil {
		defaults = SubnetEVMDefaultNetworkUpgrades(agoUpgrades)
		if *n.SubnetEVMTimestamp != *defaults.SubnetEVMTimestamp {
			return fmt.Errorf("SubnetEVM fork block timestamp is invalid: %d != %d", *n.SubnetEVMTimestamp, *defaults.SubnetEVMTimestamp)
		}
	} else {
		defaults = CorethDefaultNetworkUpgrades(agoUpgrades)
	}
	if err := verifyWithDefault(n.DurangoTimestamp, defaults.DurangoTimestamp); err != nil {
		return fmt.Errorf("Durango fork block timestamp is invalid: %w", err)
	}
	if err := verifyWithDefault(n.EtnaTimestamp, defaults.EtnaTimestamp); err != nil {
		return fmt.Errorf("Etna fork block timestamp is invalid: %w", err)
	}
	if err := verifyWithDefault(n.FortunaTimestamp, defaults.FortunaTimestamp); err != nil {
		return fmt.Errorf("Fortuna fork block timestamp is invalid: %w", err)
	}
	return nil
}

func (n *NetworkUpgrades) Override(o *NetworkUpgrades) {
	if o.SubnetEVMTimestamp != nil {
		n.SubnetEVMTimestamp = o.SubnetEVMTimestamp
	}
	if o.DurangoTimestamp != nil {
		n.DurangoTimestamp = o.DurangoTimestamp
	}
	if o.EtnaTimestamp != nil {
		n.EtnaTimestamp = o.EtnaTimestamp
	}
	if o.FortunaTimestamp != nil {
		n.FortunaTimestamp = o.FortunaTimestamp
	}
}

type AvalancheRules struct {
	IsApricotPhase1, IsApricotPhase2, IsApricotPhase3, IsApricotPhase4, IsApricotPhase5 bool
	IsApricotPhasePre6, IsApricotPhase6, IsApricotPhasePost6, IsBanff, IsCortina        bool
	IsSubnetEVM                                                                         bool
	IsDurango                                                                           bool
	IsEtna                                                                              bool
	IsFortuna                                                                           bool
}

func (n *NetworkUpgrades) GetAvalancheRules(timestamp uint64) AvalancheRules {
	return AvalancheRules{
		IsApricotPhase1:     n.IsApricotPhase1(timestamp),
		IsApricotPhase2:     n.IsApricotPhase2(timestamp),
		IsApricotPhase3:     n.IsApricotPhase3(timestamp),
		IsApricotPhase4:     n.IsApricotPhase4(timestamp),
		IsApricotPhase5:     n.IsApricotPhase5(timestamp),
		IsApricotPhasePre6:  n.IsApricotPhasePre6(timestamp),
		IsApricotPhase6:     n.IsApricotPhase6(timestamp),
		IsApricotPhasePost6: n.IsApricotPhasePost6(timestamp),
		IsBanff:             n.IsBanff(timestamp),
		IsCortina:           n.IsCortina(timestamp),
		IsSubnetEVM:         n.IsSubnetEVM(timestamp),
		IsDurango:           n.IsDurango(timestamp),
		IsEtna:              n.IsEtna(timestamp),
		IsFortuna:           n.IsFortuna(timestamp),
	}
}

// verifyWithDefault checks that the provided timestamp is greater than or equal to the default timestamp.
func verifyWithDefault(configTimestamp *uint64, defaultTimestamp *uint64) error {
	if defaultTimestamp == nil {
		return nil
	}

	if configTimestamp == nil {
		return errCannotBeNil
	}

	if *configTimestamp < *defaultTimestamp {
		return fmt.Errorf("provided timestamp (%d) must be greater than or equal to the default timestamp (%d)", *configTimestamp, *defaultTimestamp)
	}
	return nil
}

func ptrToString(val *uint64) string {
	if val == nil {
		return "nil"
	}
	return fmt.Sprintf("%d", *val)
}
