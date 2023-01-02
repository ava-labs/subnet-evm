// (c) 2022 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ava-labs/subnet-evm/precompile/deployerallowlist"
	"github.com/ava-labs/subnet-evm/precompile/feemanager"
	"github.com/ava-labs/subnet-evm/precompile/nativeminter"
	"github.com/ava-labs/subnet-evm/precompile/rewardmanager"
	"github.com/ava-labs/subnet-evm/precompile/txallowlist"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

// PrecompileUpgrade is a helper struct embedded in UpgradeConfig, representing
// each of the possible stateful precompile types that can be activated
// as a network upgrade.
type PrecompileUpgrade struct {
	ContractDeployerAllowListConfig *deployerallowlist.ContractDeployerAllowListConfig `json:"contractDeployerAllowListConfig,omitempty"` // Config for the contract deployer allow list precompile
	ContractNativeMinterConfig      *nativeminter.ContractNativeMinterConfig           `json:"contractNativeMinterConfig,omitempty"`      // Config for the native minter precompile
	TxAllowListConfig               *txallowlist.TxAllowListConfig                     `json:"txAllowListConfig,omitempty"`               // Config for the tx allow list precompile
	FeeManagerConfig                *feemanager.FeeConfigManagerConfig                 `json:"feeManagerConfig,omitempty"`                // Config for the fee manager precompile
	RewardManagerConfig             *rewardmanager.RewardManagerConfig                 `json:"rewardManagerConfig,omitempty"`             // Config for the reward manager precompile
	// ADD YOUR PRECOMPILE HERE
	// {YourPrecompile}Config  *precompile.{YourPrecompile}Config `json:"{yourPrecompile}Config,omitempty"`
}

// getByAddress returns the precompile config for the given address.
func (p *PrecompileUpgrade) getByAddress(address common.Address) (precompile.StatefulPrecompileConfig, bool) {
	switch address {
	case deployerallowlist.Address:
		return p.ContractDeployerAllowListConfig, p.ContractDeployerAllowListConfig != nil
	case nativeminter.Address:
		return p.ContractNativeMinterConfig, p.ContractNativeMinterConfig != nil
	case txallowlist.Address:
		return p.TxAllowListConfig, p.TxAllowListConfig != nil
	case feemanager.Address:
		return p.FeeManagerConfig, p.FeeManagerConfig != nil
	case rewardmanager.Address:
		return p.RewardManagerConfig, p.RewardManagerConfig != nil
	// ADD YOUR PRECOMPILE HERE
	/*
		case precompile.{YourPrecompile}Address:
			return p.{YourPrecompile}Config, p.{YourPrecompile}Config != nil
	*/
	default:
		panic(fmt.Sprintf("unknown precompile address: %v", address))
	}
}

// verifyPrecompileUpgrades checks [c.PrecompileUpgrades] is well formed:
// - [upgrades] must specify exactly one key per PrecompileUpgrade
// - the specified blockTimestamps must monotonically increase
// - the specified blockTimestamps must be compatible with those
//   specified in the chainConfig by genesis.
// - check a precompile is disabled before it is re-enabled
func (c *ChainConfig) verifyPrecompileUpgrades() error {
	var lastBlockTimestamp *big.Int
	for i, upgrade := range c.PrecompileUpgrades {
		hasKey := false // used to verify if there is only one key per Upgrade

		for _, module := range precompile.RegisteredModules {
			address := module.Address()
			config, ok := upgrade.getByAddress(address)
			if !ok {
				continue
			}
			if hasKey {
				return fmt.Errorf("PrecompileUpgrades[%d] has more than one key set", i)
			}
			configTimestamp := config.Timestamp()
			if configTimestamp == nil {
				return fmt.Errorf("PrecompileUpgrades[%d] cannot have a nil timestamp", i)
			}
			// Verify specified timestamps are monotonically increasing across all precompile keys.
			// Note: It is OK for multiple configs of different keys to specify the same timestamp.
			if lastBlockTimestamp != nil && configTimestamp.Cmp(lastBlockTimestamp) < 0 {
				return fmt.Errorf("PrecompileUpgrades[%d] config timestamp (%v) < previous timestamp (%v)", i, configTimestamp, lastBlockTimestamp)
			}
			lastBlockTimestamp = configTimestamp
			hasKey = true
		}
		if !hasKey {
			return fmt.Errorf("empty precompile upgrade at index %d", i)
		}
	}

	for _, module := range precompile.RegisteredModules {
		var (
			lastUpgraded *big.Int
			disabled     bool
		)
		address := module.Address()
		// check the genesis chain config for any enabled upgrade
		if config, ok := c.PrecompileUpgrade.getByAddress(address); ok {
			if err := config.Verify(); err != nil {
				return err
			}
			disabled = false
			lastUpgraded = config.Timestamp()
		} else {
			disabled = true
		}
		// next range over upgrades to verify correct use of disabled and blockTimestamps.
		for i, upgrade := range c.PrecompileUpgrades {
			config, ok := upgrade.getByAddress(address)
			// Skip the upgrade if it's not relevant to [address].
			if !ok {
				continue
			}

			if disabled == config.IsDisabled() {
				return fmt.Errorf("PrecompileUpgrades[%d] disable should be [%v]", i, !disabled)
			}
			if lastUpgraded != nil && (config.Timestamp().Cmp(lastUpgraded) <= 0) {
				return fmt.Errorf("PrecompileUpgrades[%d] config timestamp (%v) <= previous timestamp (%v)", i, config.Timestamp(), lastUpgraded)
			}

			if err := config.Verify(); err != nil {
				return err
			}

			disabled = config.IsDisabled()
			lastUpgraded = config.Timestamp()
		}
	}

	return nil
}

// getActivePrecompileConfig returns the most recent precompile config corresponding to [address].
// If none have occurred, returns nil.
func (c *ChainConfig) getActivePrecompileConfig(blockTimestamp *big.Int, address common.Address, upgrades []PrecompileUpgrade) precompile.StatefulPrecompileConfig {
	configs := c.getActivatingPrecompileConfigs(nil, blockTimestamp, address, upgrades)
	if len(configs) == 0 {
		return nil
	}
	return configs[len(configs)-1] // return the most recent config
}

// getActivatingPrecompileConfigs returns all forks configured to activate during the state transition from a block with timestamp [from]
// to a block with timestamp [to].
func (c *ChainConfig) getActivatingPrecompileConfigs(from *big.Int, to *big.Int, address common.Address, upgrades []PrecompileUpgrade) []precompile.StatefulPrecompileConfig {
	configs := make([]precompile.StatefulPrecompileConfig, 0)
	// First check the embedded [upgrade] for precompiles configured
	// in the genesis chain config.
	if config, ok := c.PrecompileUpgrade.getByAddress(address); ok {
		if utils.IsForkTransition(config.Timestamp(), from, to) {
			configs = append(configs, config)
		}
	}
	// Loop over all upgrades checking for the requested precompile config.
	for _, upgrade := range upgrades {
		if config, ok := upgrade.getByAddress(address); ok {
			// Check if the precompile activates in the specified range.
			if utils.IsForkTransition(config.Timestamp(), from, to) {
				configs = append(configs, config)
			}
		}
	}
	return configs
}

func (c *ChainConfig) GetPrecompileConfig(address common.Address, blockTimestamp *big.Int) precompile.StatefulPrecompileConfig {
	if val := c.getActivePrecompileConfig(blockTimestamp, address, c.PrecompileUpgrades); val != nil {
		return val
	}
	return nil
}

// CheckPrecompilesCompatible checks if [precompileUpgrades] are compatible with [c] at [headTimestamp].
// Returns a ConfigCompatError if upgrades already forked at [headTimestamp] are missing from
// [precompileUpgrades]. Upgrades not already forked may be modified or absent from [precompileUpgrades].
// Returns nil if [precompileUpgrades] is compatible with [c].
// Assumes given timestamp is the last accepted block timestamp.
// This ensures that as long as the node has not accepted a block with a different rule set it will allow a new upgrade to be applied as long as it activates after the last accepted block.
func (c *ChainConfig) CheckPrecompilesCompatible(precompileUpgrades []PrecompileUpgrade, lastTimestamp *big.Int) *ConfigCompatError {
	for _, module := range precompile.RegisteredModules {
		address := module.Address()
		if err := c.checkPrecompileCompatible(address, precompileUpgrades, lastTimestamp); err != nil {
			return err
		}
	}

	return nil
}

// checkPrecompileCompatible verifies that the precompile specified by [address] is compatible between [c] and [precompileUpgrades] at [headTimestamp].
// Returns an error if upgrades already forked at [headTimestamp] are missing from [precompileUpgrades].
// Upgrades that have already gone into effect cannot be modified or absent from [precompileUpgrades].
func (c *ChainConfig) checkPrecompileCompatible(address common.Address, precompileUpgrades []PrecompileUpgrade, lastTimestamp *big.Int) *ConfigCompatError {
	// all active upgrades must match
	activeUpgrades := c.getActivatingPrecompileConfigs(nil, lastTimestamp, address, c.PrecompileUpgrades)
	newUpgrades := c.getActivatingPrecompileConfigs(nil, lastTimestamp, address, precompileUpgrades)

	// first, check existing upgrades are there
	for i, upgrade := range activeUpgrades {
		if len(newUpgrades) <= i {
			// missing upgrade
			return newCompatError(
				fmt.Sprintf("missing PrecompileUpgrade[%d]", i),
				upgrade.Timestamp(),
				nil,
			)
		}
		// All upgrades that have forked must be identical.
		if !upgrade.Equal(newUpgrades[i]) {
			return newCompatError(
				fmt.Sprintf("PrecompileUpgrade[%d]", i),
				upgrade.Timestamp(),
				newUpgrades[i].Timestamp(),
			)
		}
	}
	// then, make sure newUpgrades does not have additional upgrades
	// that are already activated. (cannot perform retroactive upgrade)
	if len(newUpgrades) > len(activeUpgrades) {
		return newCompatError(
			fmt.Sprintf("cannot retroactively enable PrecompileUpgrade[%d]", len(activeUpgrades)),
			nil,
			newUpgrades[len(activeUpgrades)].Timestamp(), // this indexes to the first element in newUpgrades after the end of activeUpgrades
		)
	}

	return nil
}

// EnabledStatefulPrecompiles returns a slice of stateful precompile configs that
// have been activated through an upgrade.
func (c *ChainConfig) EnabledStatefulPrecompiles(blockTimestamp *big.Int) []precompile.StatefulPrecompileConfig {
	statefulPrecompileConfigs := make([]precompile.StatefulPrecompileConfig, 0)
	for _, module := range precompile.RegisteredModules {
		address := module.Address()
		if config := c.getActivePrecompileConfig(blockTimestamp, address, c.PrecompileUpgrades); config != nil {
			statefulPrecompileConfigs = append(statefulPrecompileConfigs, config)
		}
	}

	return statefulPrecompileConfigs
}

// ConfigurePrecompiles checks if any of the precompiles specified by the chain config are enabled or disabled by the block
// transition from [parentTimestamp] to the timestamp set in [blockContext]. If this is the case, it calls [Configure]
// or [Deconfigure] to apply the necessary state transitions for the upgrade.
// This function is called:
// - within genesis setup to configure the starting state for precompiles enabled at genesis,
// - during block processing to update the state before processing the given block.
func (c *ChainConfig) ConfigurePrecompiles(parentTimestamp *big.Int, blockContext precompile.BlockContext, statedb precompile.StateDB) error {
	blockTimestamp := blockContext.Timestamp()
	for _, module := range precompile.RegisteredModules { // Note: configure precompiles in a deterministic order.
		address := module.Address()
		for _, config := range c.getActivatingPrecompileConfigs(parentTimestamp, blockTimestamp, address, c.PrecompileUpgrades) {
			// If this transition activates the upgrade, configure the stateful precompile.
			// (or deconfigure it if it is being disabled.)
			if config.IsDisabled() {
				log.Info("Disabling precompile", "precompileAddress", address) // TODO: use proper names for precompiles
				statedb.Suicide(config.Address())
				// Calling Finalise here effectively commits Suicide call and wipes the contract state.
				// This enables re-configuration of the same contract state in the same block.
				// Without an immediate Finalise call after the Suicide, a reconfigured precompiled state can be wiped out
				// since Suicide will be committed after the reconfiguration.
				statedb.Finalise(true)
			} else {
				log.Info("Activating new precompile", "precompileAddress", address, "config", config)
				if err := precompile.Configure(c, blockContext, config, statedb); err != nil {
					return fmt.Errorf("could not configure precompile, precompileAddress: %s, reason: %w", address, err)
				}
			}
		}
	}
	return nil
}
