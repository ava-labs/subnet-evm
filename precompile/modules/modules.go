// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package modules

import (
	"fmt"

	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
)

var (
	// registeredModulesIndex is a map of key to Module
	// used to quickly look up a module by key
	registeredModulesIndex = make(map[common.Address]int, 0)
	// registeredModules is a list of Module to preserve order
	// for deterministic iteration
	registeredModules = make([]contract.Module, 0)

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

// RegisterModule registers a stateful precompile module
// This function should be called in the init function of the module
// to ensure that the module is registered before the node starts
// and the module is available for use.
// This function will panic if the module cannot be registered.
func RegisterModule(stm contract.Module) {
	if err := registerModule(stm); err != nil {
		panic(err)
	}
}

func registerModule(stm contract.Module) error {
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
	return config.RegisterConfig(key, address, stm)
}

func GetPrecompileModuleByAddress(address common.Address) (contract.Module, bool) {
	index, ok := registeredModulesIndex[address]
	if !ok {
		return nil, false
	}

	return registeredModules[index], true
}

func GetPrecompileModule(key string) (contract.Module, bool) {
	for _, stm := range registeredModules {
		if stm.Key() == key {
			return stm, true
		}
	}

	return nil, false
}

func RegisteredModules() []contract.Module {
	return registeredModules
}
