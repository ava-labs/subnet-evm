// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type StatefulPrecompileModule struct {
	// Address returns the address where the stateful precompile is accessible.
	Address common.Address
	// Contract returns a thread-safe singleton that can be used as the StatefulPrecompiledContract when
	// this config is enabled.
	Contract StatefulPrecompiledContract
	// Key returns the unique key for the stateful precompile.
	Key string
	// NewConfig returns a new instance of the stateful precompile config.
	NewConfigFn func() StatefulPrecompileConfig
}

var (
	// registeredModulesIndex is a map of key to StatefulPrecompileModule
	// used to quickly look up a module by key
	registeredModulesIndex = make(map[common.Address]int, 0)
	// registeredModules is a list of StatefulPrecompileModule to preserve order
	// for deterministic iteration
	registeredModules = make([]StatefulPrecompileModule, 0)

	reservedRanges = []AddressRange{
		{
			common.HexToAddress("0x0100000000000000000000000000000000000000"),
			common.HexToAddress("0x01000000000000000000000000000000000000ff"),
		},
		{
			common.HexToAddress("0x0200000000000000000000000000000000000000"),
			common.HexToAddress("0x02000000000000000000000000000000000000ff"),
		},
		{
			common.HexToAddress("0x0300000000000000000000000000000000000000"),
			common.HexToAddress("0x03000000000000000000000000000000000000ff"),
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

func RegisterModule(stm StatefulPrecompileModule) error {
	address := stm.Address
	key := stm.Key
	if !ReservedAddress(address) {
		return fmt.Errorf("address %s not in a reserved range", address)
	}

	for _, registeredModule := range registeredModules {
		if registeredModule.Key == key {
			return fmt.Errorf("name %s already used by a stateful precompile", key)
		}
		if registeredModule.Address == address {
			return fmt.Errorf("address %s already used by a stateful precompile", address)
		}
	}

	registeredModulesIndex[address] = len(registeredModules)
	registeredModules = append(registeredModules, stm)

	return nil
}

func GetPrecompileModuleByAddress(address common.Address) (StatefulPrecompileModule, bool) {
	index, ok := registeredModulesIndex[address]
	if !ok {
		return StatefulPrecompileModule{}, false
	}

	return registeredModules[index], true
}

func GetPrecompileModule(key string) (StatefulPrecompileModule, bool) {
	for _, stm := range registeredModules {
		if stm.Key == key {
			return stm, true
		}
	}

	return StatefulPrecompileModule{}, false
}

func RegisteredModules() []StatefulPrecompileModule {
	return registeredModules
}
