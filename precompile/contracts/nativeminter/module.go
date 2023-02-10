// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nativeminter

import (
	"fmt"
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
)

var _ contract.Module = &Module{}

// ConfigKey is the key used in json config files to specify this precompile config.
// must be unique across all precompiles.
const ConfigKey = "contractNativeMinterConfig"

var ContractAddress = common.HexToAddress("0x0200000000000000000000000000000000000001")

type Module struct{}

// Address returns the address of the native minter.
func (Module) Address() common.Address {
	return ContractAddress
}

func (Module) NewConfig() config.Config {
	return &ContractNativeMinterConfig{}
}

// Configure configures [state] with the desired admins based on [cfg].
func (Module) Configure(_ contract.ChainConfig, cfg config.Config, state contract.StateDB, _ contract.BlockContext) error {
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

func (Module) Contract() contract.StatefulPrecompiledContract {
	return ContractNativeMinterPrecompile
}
