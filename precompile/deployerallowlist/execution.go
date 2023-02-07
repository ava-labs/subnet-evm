// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package deployerallowlist

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/precompile/execution"
)

type Executor struct{}

// Configure configures [state] with the desired admins based on [cfg].
func (Executor) Configure(_ execution.ChainConfig, cfg config.Config, state execution.StateDB, _ execution.BlockContext) error {
	config, ok := cfg.(*ContractDeployerAllowListConfig)
	if !ok {
		return fmt.Errorf("incorrect config %T: %v", config, config)
	}
	return config.AllowListConfig.Configure(state, ContractAddress)
}

func (Executor) Contract() execution.Contract {
	return ContractDeployerAllowListPrecompile
}
