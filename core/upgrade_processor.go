// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package core

import (
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/core/vm"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/precompile/execution"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ethereum/go-ethereum/log"
)

// ConfigurePrecompiles checks if any of the precompiles specified by the chain config are enabled or disabled by the block
// transition from [parentTimestamp] to the timestamp set in [blockContext]. If this is the case, it calls [Configure]
// or [Deconfigure] to apply the necessary state transitions for the upgrade.
// This function is called:
// - within genesis setup to configure the starting state for precompiles enabled at genesis,
// - during block processing to update the state before processing the given block.
func ConfigurePrecompiles(c *params.ChainConfig, parentTimestamp *big.Int, blockContext execution.BlockContext, statedb vm.StateDB) error {
	blockTimestamp := blockContext.Timestamp()
	// Note: RegisteredModules returns precompiles in order they are registered.
	// This is important because we want to configure precompiles in the same order
	// so that the state is deterministic.
	for _, config := range config.GetConfigs() {
		key := config.Key()
		for _, activatingConfig := range c.GetActivatingPrecompileConfigs(config.Address(), parentTimestamp, blockTimestamp, c.PrecompileUpgrades) {
			// If this transition activates the upgrade, configure the stateful precompile.
			// (or deconfigure it if it is being disabled.)
			if activatingConfig.IsDisabled() {
				log.Info("Disabling precompile", "name", key)
				statedb.Suicide(config.Address())
				// Calling Finalise here effectively commits Suicide call and wipes the contract state.
				// This enables re-configuration of the same contract state in the same block.
				// Without an immediate Finalise call after the Suicide, a reconfigured precompiled state can be wiped out
				// since Suicide will be committed after the reconfiguration.
				statedb.Finalise(true)
			} else {
				module, ok := modules.GetPrecompileModule(key)
				if !ok {
					return fmt.Errorf("could not find module for activating precompile, name: %s", key)
				}
				log.Info("Activating new precompile", "name", key, "config", activatingConfig)
				if err := module.Executor().Configure(c, activatingConfig, statedb, blockContext); err != nil {
					return fmt.Errorf("could not configure precompile, name: %s, reason: %w", key, err)
				}
			}
		}
	}
	return nil
}
