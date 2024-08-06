// (c) 2024 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
)

// UpgradeConfig includes the following configs that may be specified in upgradeBytes:
// - Timestamps that enable avalanche network upgrades,
// - Enabling or disabling precompiles as network upgrades.
type UpgradeConfig struct {
	// Config for timestamps that enable network upgrades.
	NetworkUpgradeOverrides *NetworkUpgrades `json:"networkUpgradeOverrides,omitempty"`

	// Config for modifying state as a network upgrade.
	StateUpgrades []StateUpgrade `json:"stateUpgrades,omitempty"`

	// Config for enabling and disabling precompiles as network upgrades.
	PrecompileUpgrades []PrecompileUpgrade `json:"precompileUpgrades,omitempty"`
}

// AvalancheContext provides Avalanche specific context directly into the EVM.
type AvalancheContext struct {
	SnowCtx *snow.Context
}

// MarshalJSON marshal the precompile config into json based on the precompile key.
// Ex: {"feeManagerConfig": {...}} where "feeManagerConfig" is the key
func (u *PrecompileUpgrade) MarshalJSON() ([]byte, error) {
	res := make(map[string]precompileconfig.Config)
	res[u.Key()] = u.Config
	return json.Marshal(res)
}

// MarshalJSON returns the JSON encoding of c.
// This is a custom marshaler to handle the Precompiles field.
func (c *ChainConfig) MarshalJSON() ([]byte, error) {
	// TODO refactor this (DO NOT MERGE)

	// Alias ChainConfig to avoid recursion
	type _ChainConfig ChainConfig
	tmp, err := json.Marshal((*_ChainConfig)(c))
	if err != nil {
		return nil, err
	}

	// To include PrecompileUpgrades, we unmarshal the json representing c
	// then directly add the corresponding keys to the json.
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(tmp, &raw); err != nil {
		return nil, err
	}

	for key, value := range c.GenesisPrecompiles {
		conf, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		raw[key] = conf
	}

	return json.Marshal(raw)
}

type ChainConfigWithUpgradesJSON struct {
	ChainConfig
	UpgradeConfig UpgradeConfig `json:"upgrades,omitempty"`
}

// MarshalJSON implements json.Marshaler. This is a workaround for the fact that
// the embedded ChainConfig struct has a MarshalJSON method, which prevents
// the default JSON marshalling from working for UpgradeConfig.
// TODO: consider removing this method by allowing external tag for the embedded
// ChainConfig struct.
func (cu ChainConfigWithUpgradesJSON) MarshalJSON() ([]byte, error) {
	// embed the ChainConfig struct into the response
	chainConfigJSON, err := json.Marshal(cu.ChainConfig)
	if err != nil {
		return nil, err
	}
	if len(chainConfigJSON) > maxJSONLen {
		return nil, errors.New("value too large")
	}

	type upgrades struct {
		UpgradeConfig UpgradeConfig `json:"upgrades"`
	}

	upgradeJSON, err := json.Marshal(upgrades{cu.UpgradeConfig})
	if err != nil {
		return nil, err
	}
	if len(upgradeJSON) > maxJSONLen {
		return nil, errors.New("value too large")
	}

	// merge the two JSON objects
	mergedJSON := make([]byte, 0, len(chainConfigJSON)+len(upgradeJSON)+1)
	mergedJSON = append(mergedJSON, chainConfigJSON[:len(chainConfigJSON)-1]...)
	mergedJSON = append(mergedJSON, ',')
	mergedJSON = append(mergedJSON, upgradeJSON[1:]...)
	return mergedJSON, nil
}

// ToWithUpgradesJSON converts the ChainConfig to ChainConfigWithUpgradesJSON with upgrades explicitly displayed.
// ChainConfig does not include upgrades in its JSON output.
// This is a workaround for showing upgrades in the JSON output.
func (c *ChainConfig) ToWithUpgradesJSON() *ChainConfigWithUpgradesJSON {
	return &ChainConfigWithUpgradesJSON{
		ChainConfig:   *c,
		UpgradeConfig: c.UpgradeConfig,
	}
}

func (c *ChainConfig) SetNetworkUpgradeDefaults() {
	if c.HomesteadBlock == nil {
		c.HomesteadBlock = big.NewInt(0)
	}
	if c.EIP150Block == nil {
		c.EIP150Block = big.NewInt(0)
	}
	if c.EIP155Block == nil {
		c.EIP155Block = big.NewInt(0)
	}
	if c.EIP158Block == nil {
		c.EIP158Block = big.NewInt(0)
	}
	if c.ByzantiumBlock == nil {
		c.ByzantiumBlock = big.NewInt(0)
	}
	if c.ConstantinopleBlock == nil {
		c.ConstantinopleBlock = big.NewInt(0)
	}
	if c.PetersburgBlock == nil {
		c.PetersburgBlock = big.NewInt(0)
	}
	if c.IstanbulBlock == nil {
		c.IstanbulBlock = big.NewInt(0)
	}
	if c.MuirGlacierBlock == nil {
		c.MuirGlacierBlock = big.NewInt(0)
	}

	c.NetworkUpgrades.setDefaults(c.SnowCtx.NetworkID)
}

// SetEVMUpgrades sets the mapped upgrades  Avalanche > EVM upgrades) for the chain config.
func (c *ChainConfig) SetEVMUpgrades(avalancheUpgrades NetworkUpgrades) {
	if avalancheUpgrades.EUpgradeTimestamp != nil {
		c.CancunTime = utils.NewUint64(*avalancheUpgrades.EUpgradeTimestamp)
	}
}

func (c *ChainConfig) GetActivePrecompileConfig(address common.Address, timestamp uint64) precompileconfig.Config {
	configs := c.GetActivatingPrecompileConfigs(address, nil, timestamp, c.PrecompileUpgrades)
	if len(configs) == 0 {
		return nil
	}
	return configs[len(configs)-1] // return the most recent config
}

func (c *ChainConfig) GetActivatingPrecompileConfigs(address common.Address, from *uint64, to uint64, upgrades []PrecompileUpgrade) []precompileconfig.Config {
	var configs []precompileconfig.Config
	maybeAppend := func(pc precompileconfig.Config) {
		if pc.Address() == address && utils.IsForkTransition(pc.Timestamp(), from, to) {
			configs = append(configs, pc)
		}
	}
	// First check the embedded [upgrade] for precompiles configured
	// in the genesis chain config.
	for _, p := range c.GenesisPrecompiles {
		maybeAppend(p)
	}
	// Loop over all upgrades checking for the requested precompile config.
	for _, upgrade := range upgrades {
		maybeAppend(upgrade.Config)
	}
	return configs
}

// IsPrecompileEnabled returns whether precompile with [address] is enabled at [timestamp].
func (c *ChainConfig) IsPrecompileEnabled(address common.Address, timestamp uint64) bool {
	config := c.GetActivePrecompileConfig(address, timestamp)
	return config != nil && !config.IsDisabled()
}

func (c *ChainConfig) allPrecompileAddresses(extra ...PrecompileUpgrade) map[string]common.Address {
	all := make(map[string]common.Address)
	add := func(pc precompileconfig.Config) {
		if a, ok := all[pc.Key()]; ok && a != pc.Address() {
			panic("DO NOT MERGE")
		}
		all[pc.Key()] = pc.Address()
	}

	for _, p := range c.GenesisPrecompiles {
		add(p)
	}
	for _, p := range c.PrecompileUpgrades {
		add(p)
	}
	for _, p := range extra {
		add(p)
	}
	return all
}

// CheckPrecompilesCompatible checks if [precompileUpgrades] are compatible with [c] at [headTimestamp].
// Returns a ConfigCompatError if upgrades already activated at [headTimestamp] are missing from
// [precompileUpgrades]. Upgrades not already activated may be modified or absent from [precompileUpgrades].
// Returns nil if [precompileUpgrades] is compatible with [c].
// Assumes given timestamp is the last accepted block timestamp.
// This ensures that as long as the node has not accepted a block with a different rule set it will allow a
// new upgrade to be applied as long as it activates after the last accepted block.
func (c *ChainConfig) CheckPrecompilesCompatible(precompileUpgrades []PrecompileUpgrade, time uint64) *ConfigCompatError {
	addrs := c.allPrecompileAddresses(precompileUpgrades...)
	for _, a := range addrs {
		if err := c.checkPrecompileCompatible(a, precompileUpgrades, time); err != nil {
			return err
		}
	}
	return nil
}

// checkPrecompileCompatible verifies that the precompile specified by [address] is compatible between [c]
// and [precompileUpgrades] at [headTimestamp].
// Returns an error if upgrades already activated at [headTimestamp] are missing from [precompileUpgrades].
// Upgrades that have already gone into effect cannot be modified or absent from [precompileUpgrades].
func (c *ChainConfig) checkPrecompileCompatible(address common.Address, precompileUpgrades []PrecompileUpgrade, time uint64) *ConfigCompatError {
	// All active upgrades (from nil to [lastTimestamp]) must match.
	activeUpgrades := c.GetActivatingPrecompileConfigs(address, nil, time, c.PrecompileUpgrades)
	newUpgrades := c.GetActivatingPrecompileConfigs(address, nil, time, precompileUpgrades)

	// Check activated upgrades are still present.
	for i, upgrade := range activeUpgrades {
		if len(newUpgrades) <= i {
			// missing upgrade
			return newTimestampCompatError(
				fmt.Sprintf("missing PrecompileUpgrade[%d]", i),
				upgrade.Timestamp(),
				nil,
			)
		}
		// All upgrades that have activated must be identical.
		if !upgrade.Equal(newUpgrades[i]) {
			return newTimestampCompatError(
				fmt.Sprintf("PrecompileUpgrade[%d]", i),
				upgrade.Timestamp(),
				newUpgrades[i].Timestamp(),
			)
		}
	}
	// then, make sure newUpgrades does not have additional upgrades
	// that are already activated. (cannot perform retroactive upgrade)
	if len(newUpgrades) > len(activeUpgrades) {
		return newTimestampCompatError(
			fmt.Sprintf("cannot retroactively enable PrecompileUpgrade[%d]", len(activeUpgrades)),
			nil,
			newUpgrades[len(activeUpgrades)].Timestamp(), // this indexes to the first element in newUpgrades after the end of activeUpgrades
		)
	}

	return nil
}

// EnabledStatefulPrecompiles returns current stateful precompile configs that are enabled at [blockTimestamp].
func (c *ChainConfig) EnabledStatefulPrecompiles(blockTimestamp uint64) Precompiles {
	statefulPrecompileConfigs := make(Precompiles)
	for key, addr := range c.allPrecompileAddresses() {
		if config := c.GetActivePrecompileConfig(addr, blockTimestamp); config != nil && !config.IsDisabled() {
			statefulPrecompileConfigs[key] = config
		}
	}
	return statefulPrecompileConfigs
}
