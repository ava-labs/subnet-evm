package precompileconfig

import (
	"fmt"
	"sort"

	"github.com/ethereum/go-ethereum/common"
)

var (
	// registeredModules is a list of Module to preserve order
	// for deterministic iteration
	registeredModules = make([]Module, 0)
)

// RegisterModule registers a stateful precompile module
func RegisterModule(stm Module) error {
	address := stm.Address
	key := stm.ConfigKey

	for _, registeredModule := range registeredModules {
		if registeredModule.ConfigKey == key {
			return fmt.Errorf("name %s already used by a stateful precompile", key)
		}
		if registeredModule.Address == address {
			return fmt.Errorf("address %s already used by a stateful precompile", address)
		}
	}
	// sort by address to ensure deterministic iteration
	registeredModules = insertSortedByAddress(registeredModules, stm)
	return nil
}

func GetPrecompileModuleByAddress(address common.Address) (Module, bool) {
	for _, stm := range registeredModules {
		if stm.Address == address {
			return stm, true
		}
	}
	return Module{}, false
}

func GetPrecompileModule(key string) (Module, bool) {
	for _, stm := range registeredModules {
		if stm.ConfigKey == key {
			return stm, true
		}
	}
	return Module{}, false
}

func RegisteredModules() []Module {
	return registeredModules
}

func insertSortedByAddress(data []Module, stm Module) []Module {
	data = append(data, stm)
	sort.Sort(moduleArray(data))
	return data
}
