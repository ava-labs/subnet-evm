// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// Gas costs for stateful precompiles
const (
	WriteGasCostPerSlot = 20_000
	ReadGasCostPerSlot  = 5_000
)

// Designated addresses of stateful precompiles
// Note: it is important that none of these addresses conflict with each other or any other precompiles
// in core/vm/contracts.go.
// The first stateful precompiles were added in coreth to support nativeAssetCall and nativeAssetBalance. New stateful precompiles
// originating in coreth will continue at this prefix, so we reserve this range in subnet-evm so that they can be migrated into
// subnet-evm without issue.
// These start at the address: 0x0100000000000000000000000000000000000000 and will increment by 1.
// Optional precompiles implemented in subnet-evm start at 0x0200000000000000000000000000000000000000 and will increment by 1
// from here to reduce the risk of conflicts.
// For forks of subnet-evm, users should start at 0x0300000000000000000000000000000000000000 to ensure
// that their own modifications do not conflict with stateful precompiles that may be added to subnet-evm
// in the future.
var (
	// This list is kept just for reference. The actual addresses defined in respective packages of precompiles.
	// ContractDeployerAllowListAddress = common.HexToAddress("0x0200000000000000000000000000000000000000")
	// ContractNativeMinterAddress      = common.HexToAddress("0x0200000000000000000000000000000000000001")
	// TxAllowListAddress               = common.HexToAddress("0x0200000000000000000000000000000000000002")
	// FeeManagerAddress         				= common.HexToAddress("0x0200000000000000000000000000000000000003")
	// RewardManagerAddress             = common.HexToAddress("0x0200000000000000000000000000000000000004")
	// ADD YOUR PRECOMPILE HERE
	// {YourPrecompile}Address       = common.HexToAddress("0x03000000000000000000000000000000000000??")

	registeredModules = make(map[string]StatefulPrecompileModule, 0)
	// Sorted with the addresses in ascending order
	sortedModules = make([]StatefulPrecompileModule, 0)

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

// UsedAddress returns true if [addr] is in a reserved range for custom precompiles
func ReservedAddress(addr common.Address) bool {
	for _, reservedRange := range reservedRanges {
		if reservedRange.Contains(addr) {
			return true
		}
	}

	return false
}

func RegisterModule(stm StatefulPrecompileModule) error {
	address := stm.Address()
	key := stm.Key()
	if !ReservedAddress(address) {
		return fmt.Errorf("address %s not in a reserved range", address)
	}
	_, ok := registeredModules[key]
	if ok {
		return fmt.Errorf("name %s already used by a stateful precompile", key)
	}

	for _, precompile := range registeredModules {
		if precompile.Address() == address {
			return fmt.Errorf("address %s already used by a stateful precompile", address)
		}
	}
	registeredModules[key] = stm
	// keep the list sorted
	insertSortedModules(stm)
	fmt.Println(sortedModules)
	return nil
}

// TODO: if not used remove this function
func GetPrecompileModule(name string) (StatefulPrecompileModule, bool) {
	stm, ok := registeredModules[name]
	return stm, ok
}

// insertSortedModules inserts the module into the sorted list of modules
func insertSortedModules(stm StatefulPrecompileModule) {
	for i := 0; i < len(sortedModules); i++ {
		if stm.Address().Hash().Big().Cmp(sortedModules[i].Address().Hash().Big()) < 0 {
			sortedModules = append(sortedModules[:i], append([]StatefulPrecompileModule{stm}, sortedModules[i:]...)...)
			return
		}
	}
	// if we get here, the module should be appended to the end of the list
	sortedModules = append(sortedModules, stm)
}

func RegisteredModules() []StatefulPrecompileModule {
	return sortedModules
}
