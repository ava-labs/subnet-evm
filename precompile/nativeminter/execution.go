// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nativeminter

import (
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/precompile/execution"
)

type Executor struct{}

// Configure configures [state] with the desired admins based on [cfg].
func (Executor) Configure(_ execution.ChainConfig, cfg config.Config, state execution.StateDB, _ execution.BlockContext) error {
	config, ok := cfg.(*ContractNativeMinterConfig)
	if !ok {
		return fmt.Errorf("incorrect config %T: %v", config, config)
	}
	for to, amount := range config.InitialMint {
		if amount != nil {
			bigIntAmount := (*big.Int)(amount)
			state.AddBalance(to, bigIntAmount)
		}
	}

	return config.AllowListConfig.Configure(state, ContractAddress)
}

func (Executor) Contract() execution.Contract {
	return ContractNativeMinterPrecompile
}
