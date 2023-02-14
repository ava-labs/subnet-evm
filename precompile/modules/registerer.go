// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package modules

import (
	"bytes"
	"fmt"

	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
)

var (
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

// RegisterModule registers a stateful precompile module
func RegisterModule(stm Module) error {
	address := stm.Address
	key := stm.ConfigKey
	if !ReservedAddress(address) {
		return fmt.Errorf("address %s not in a reserved range", address)
	}

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
	// sort by address to ensure deterministic iteration
	// start at the end of the list and work backwards
	// this is faster than sorting the list every time
	// since we expect sorted inserts
	index := 0
	for i := len(data) - 1; i >= 0; i-- {
		if bytes.Compare(stm.Address.Bytes(), data[i].Address.Bytes()) > 0 {
			index = i + 1
			break
		}
	}
	return insertAt(data, index, stm)
}

func insertAt(data []Module, index int, stm Module) []Module {
	// if the index is out of bounds, append the module
	if index >= len(data) {
		data = append(data, stm)
		return data
	}
	// shift the slice to the right and leave a space for the new element
	data = append(data[:index+1], data[index:]...)
	// Insert the new element.
	data[index] = stm
	return data
}
