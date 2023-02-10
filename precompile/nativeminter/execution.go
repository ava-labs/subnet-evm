// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nativeminter

import (
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/precompile/contract"
)

type Executor struct{}

// Configure configures [state] with the desired admins based on [cfg].
func (Executor) Configure(_ contract.ChainConfig, cfg config.Config, state contract.StateDB, _ contract.BlockContext) error {
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

func (Executor) Contract() contract.StatefulPrecompiledContract {
	return ContractNativeMinterPrecompile
}
