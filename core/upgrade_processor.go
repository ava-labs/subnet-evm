// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package core

import (
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ethereum/go-ethereum/log"
)

// ConfigurePrecompiles checks if any of the precompiles specified by the chain config are enabled or disabled by the block
// transition from [parentTimestamp] to the timestamp set in [blockContext]. If this is the case, it calls [Configure]
// or [Deconfigure] to apply the necessary state transitions for the upgrade.
// This function is called:
// - within genesis setup to configure the starting state for precompiles enabled at genesis,
// - during block processing to update the state before processing the given block.
func ConfigurePrecompiles(c *params.ChainConfig, parentTimestamp *big.Int, blockContext precompile.BlockContext, statedb *state.StateDB) error {
	blockTimestamp := blockContext.Timestamp()
	// Note: RegisteredModules returns precompiles in order they are registered.
	// This is important because we want to configure precompiles in the same order
	// so that the state is deterministic.
	for _, module := range precompile.RegisteredModules() {
		key := module.Key()
		for _, config := range c.GetActivatingPrecompileConfigs(module.Address(), parentTimestamp, blockTimestamp, c.PrecompileUpgrades) {
			// If this transition activates the upgrade, configure the stateful precompile.
			// (or deconfigure it if it is being disabled.)
			if config.IsDisabled() {
				log.Info("Disabling precompile", "name", key)
				statedb.Suicide(module.Address())
				// Calling Finalise here effectively commits Suicide call and wipes the contract state.
				// This enables re-configuration of the same contract state in the same block.
				// Without an immediate Finalise call after the Suicide, a reconfigured precompiled state can be wiped out
				// since Suicide will be committed after the reconfiguration.
				statedb.Finalise(true)
			} else {
				log.Info("Activating new precompile", "name", key, "config", config)
				if err := precompile.Configure(c, blockContext, config, statedb); err != nil {
					return fmt.Errorf("could not configure precompile, name: %s, reason: %w", key, err)
				}
			}
		}
	}
	return nil
}
