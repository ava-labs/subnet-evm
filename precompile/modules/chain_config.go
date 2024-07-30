package modules

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
)

func InitChainConfig(c *params.ChainConfig) error {
	if c.GenesisPrecompiles == nil {
		c.GenesisPrecompiles = make(params.Precompiles)
	}
	for key, raw := range c.LazyUnmarshalData {
		mod, ok := GetPrecompileModule(key)
		if !ok {
			continue
		}

		conf := mod.MakeConfig()
		if err := json.Unmarshal(raw, conf); err != nil {
			return fmt.Errorf("unmarshal %T: %v", conf, err)
		}
		c.GenesisPrecompiles[key] = conf
	}
	return nil
}

func InitChainRules(c *params.ChainConfig, r *params.Rules, timestamp uint64) {
	// Initialize the stateful precompiles that should be enabled at [blockTimestamp].
	r.ActivePrecompiles = make(map[common.Address]precompileconfig.Config)
	r.Predicaters = make(map[common.Address]precompileconfig.Predicater)
	r.AccepterPrecompiles = make(map[common.Address]precompileconfig.Accepter)

	for _, module := range RegisteredModules() {
		if config := getActivePrecompileConfig(c, module.Address, timestamp); config != nil && !config.IsDisabled() {
			r.ActivePrecompiles[module.Address] = config
			if predicater, ok := config.(precompileconfig.Predicater); ok {
				r.Predicaters[module.Address] = predicater
			}
			if precompileAccepter, ok := config.(precompileconfig.Accepter); ok {
				r.AccepterPrecompiles[module.Address] = precompileAccepter
			}
		}
	}
}

func getActivePrecompileConfig(c *params.ChainConfig, address common.Address, timestamp uint64) precompileconfig.Config {
	configs := GetActivatingPrecompileConfigs(c, address, nil, timestamp, c.PrecompileUpgrades)
	if len(configs) == 0 {
		return nil
	}
	return configs[len(configs)-1] // return the most recent config
}

func GetActivatingPrecompileConfigs(c *params.ChainConfig, address common.Address, from *uint64, to uint64, upgrades []params.PrecompileUpgrade) []precompileconfig.Config {
	// Get key from address.
	module, ok := GetPrecompileModuleByAddress(address)
	if !ok {
		return nil
	}
	configs := make([]precompileconfig.Config, 0)
	key := module.ConfigKey
	// First check the embedded [upgrade] for precompiles configured
	// in the genesis chain config.
	if config, ok := c.GenesisPrecompiles[key]; ok {
		if utils.IsForkTransition(config.Timestamp(), from, to) {
			configs = append(configs, config)
		}
	}
	// Loop over all upgrades checking for the requested precompile config.
	for _, upgrade := range upgrades {
		if upgrade.Key() == key {
			// Check if the precompile activates in the specified range.
			if utils.IsForkTransition(upgrade.Timestamp(), from, to) {
				configs = append(configs, upgrade.Config)
			}
		}
	}
	return configs
}

// CheckPrecompilesCompatible checks if [precompileUpgrades] are compatible with [c] at [headTimestamp].
// Returns a ConfigCompatError if upgrades already activated at [headTimestamp] are missing from
// [precompileUpgrades]. Upgrades not already activated may be modified or absent from [precompileUpgrades].
// Returns nil if [precompileUpgrades] is compatible with [c].
// Assumes given timestamp is the last accepted block timestamp.
// This ensures that as long as the node has not accepted a block with a different rule set it will allow a
// new upgrade to be applied as long as it activates after the last accepted block.
func CheckPrecompilesCompatible(c *params.ChainConfig, precompileUpgrades []params.PrecompileUpgrade, time uint64) *params.ConfigCompatError {
	for _, module := range RegisteredModules() {
		if err := checkPrecompileCompatible(c, module.Address, precompileUpgrades, time); err != nil {
			return err
		}
	}

	return nil
}

// checkPrecompileCompatible verifies that the precompile specified by [address] is compatible between [c]
// and [precompileUpgrades] at [headTimestamp].
// Returns an error if upgrades already activated at [headTimestamp] are missing from [precompileUpgrades].
// Upgrades that have already gone into effect cannot be modified or absent from [precompileUpgrades].
func checkPrecompileCompatible(c *params.ChainConfig, address common.Address, precompileUpgrades []params.PrecompileUpgrade, time uint64) *params.ConfigCompatError {
	// All active upgrades (from nil to [lastTimestamp]) must match.
	activeUpgrades := GetActivatingPrecompileConfigs(c, address, nil, time, c.PrecompileUpgrades)
	newUpgrades := GetActivatingPrecompileConfigs(c, address, nil, time, precompileUpgrades)

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

func newTimestampCompatError(what string, storedtime, newtime *uint64) *params.ConfigCompatError {
	return nil // TODO: DO NOT MERGE
}

// EnabledStatefulPrecompiles returns current stateful precompile configs that are enabled at [blockTimestamp].
func EnabledStatefulPrecompiles(c *params.ChainConfig, blockTimestamp uint64) params.Precompiles {
	statefulPrecompileConfigs := make(params.Precompiles)
	for _, module := range RegisteredModules() {
		if config := getActivePrecompileConfig(c, module.Address, blockTimestamp); config != nil && !config.IsDisabled() {
			statefulPrecompileConfigs[module.ConfigKey] = config
		}
	}

	return statefulPrecompileConfigs
}

type PrecompileUpgrade struct {
	precompileconfig.Config // TODO: DO NOT MERGE
}

var errNoKey = errors.New("PrecompileUpgrade cannot be empty")

// UnmarshalJSON unmarshals the json into the correct precompile config type
// based on the key. Keys are defined in each precompile module, and registered in
// precompile/registry/registry.go.
// Ex: {"feeManagerConfig": {...}} where "feeManagerConfig" is the key
func (u *PrecompileUpgrade) UnmarshalJSON(data []byte) error {
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if len(raw) == 0 {
		return errNoKey
	}
	if len(raw) > 1 {
		return fmt.Errorf("PrecompileUpgrade must have exactly one key, got %d", len(raw))
	}
	for key, value := range raw {
		module, ok := GetPrecompileModule(key)
		if !ok {
			return fmt.Errorf("unknown precompile config: %s", key)
		}
		config := module.MakeConfig()
		if err := json.Unmarshal(value, config); err != nil {
			return err
		}
		u.Config = config
	}
	return nil
}

// MarshalJSON marshal the precompile config into json based on the precompile key.
// Ex: {"feeManagerConfig": {...}} where "feeManagerConfig" is the key
func (u *PrecompileUpgrade) MarshalJSON() ([]byte, error) {
	res := make(map[string]precompileconfig.Config)
	res[u.Key()] = u.Config
	return json.Marshal(res)
}

// IsPrecompileEnabled returns whether precompile with [address] is enabled at [timestamp].
func IsPrecompileEnabled(c *params.ChainConfig, address common.Address, timestamp uint64) bool {
	config := getActivePrecompileConfig(c, address, timestamp)
	return config != nil && !config.IsDisabled()
}
