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

	RegisteredModules = make([]StatefulPrecompileModule, 0)
	reservedRanges    = []AddressRange{
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

func RegisterPrecompile(stm StatefulPrecompileModule) error {
	address := stm.Address()
	name := stm.Name()
	if !ReservedAddress(address) {
		return fmt.Errorf("address %s not in a reserved range", address)
	}
	for _, precompile := range RegisteredModules {
		if precompile.Address() == address {
			return fmt.Errorf("address %s already used by a stateful precompile", address)
		}
		if precompile.Name() == name {
			return fmt.Errorf("name %s already used by a stateful precompile", name)
		}
	}
	RegisteredModules = append(RegisteredModules, stm)

	return nil
}
