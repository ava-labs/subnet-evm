// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package deployerallowlist

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/precompile/execution"
)

type Executor struct{}

// Configure configures [state] with the desired admins based on [c].
func (Executor) Configure(_ execution.ChainConfig, precompileConfig config.Config, state *state.StateDB, _ execution.BlockContext) error {
	deployerConfig, ok := precompileConfig.(*ContractDeployerAllowListConfig)
	if !ok {
		return fmt.Errorf("incorrect config %T: %v", deployerConfig, deployerConfig)
	}
	// Configure(chainConfig ChainConfig, precompileConfig precompileConfig.Config, state *state.StateDB, blockContext BlockContext) error
	return deployerConfig.AllowListConfig.Configure(state, ContractAddress)
}

func (Executor) Contract() execution.Contract {
	return ContractDeployerAllowListPrecompile
}
