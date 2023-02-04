// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package modules

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/precompile/execution"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
)

type Module interface {
	// Key returns the unique key for the stateful precompile.
	Key() string
	// Address returns the address where the stateful precompile is accessible.
	Address() common.Address
	// Contract returns a thread-safe singleton that can be used as the StatefulPrecompiledContract when
	// this config is enabled.
	Executor() execution.Execution
	// NewConfig returns a new instance of the stateful precompile config.
	NewConfig() config.Config
}

var (
	// registeredModulesIndex is a map of key to Module
	// used to quickly look up a module by key
	registeredModulesIndex = make(map[common.Address]int, 0)
	// registeredModules is a list of Module to preserve order
	// for deterministic iteration
	registeredModules = make([]Module, 0)

	reservedRanges = []utils.AddressRange{
		{
			Start: common.HexToAddress("0x0100000000000000000000000000000000000000"),
			End:   common.HexToAddress("0x01000000000000000000000000000000000000ff"),
		},
		{
			Start: common.HexToAddress("0x0200000000000000000000000000000000000000"),
			End:   common.HexToAddress("0x02000000000000000000000000000000000000ff"),
		},
		{
			Start: common.HexToAddress("0x0300000000000000000000000000000000000000"),
			End:   common.HexToAddress("0x03000000000000000000000000000000000000ff"),
		},
	}
)

// ReservedAddress returns true if [addr] is in a reserved range for custom precompiles
func ReservedAddress(addr common.Address) bool {
	for _, reservedRange := range reservedRanges {
		if reservedRange.Contains(addr) {
			return true
		}
	}

	return false
}

func RegisterModule(stm Module) error {
	address := stm.Address()
	key := stm.Key()
	if !ReservedAddress(address) {
		return fmt.Errorf("address %s not in a reserved range", address)
	}

	for _, module := range registeredModules {
		if module.Key() == key {
			return fmt.Errorf("name %s already used by a stateful precompile", key)
		}
		if module.Address() == address {
			return fmt.Errorf("address %s already used by a stateful precompile", address)
		}
	}

	registeredModulesIndex[address] = len(registeredModules)
	registeredModules = append(registeredModules, stm)

	return nil
}

func GetPrecompileModuleByAddress(address common.Address) (Module, bool) {
	index, ok := registeredModulesIndex[address]
	if !ok {
		return nil, false
	}

	return registeredModules[index], true
}

func GetPrecompileModule(key string) (Module, bool) {
	for _, stm := range registeredModules {
		if stm.Key() == key {
			return stm, true
		}
	}

	return nil, false
}

func RegisteredModules() []Module {
	return registeredModules
}
