// Code generated
// This file is a generated precompile contract config with stubbed abstract functions.
// The file is generated by a template. Please inspect every code and comment in this file before use.

package allowlist

import (
	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// PackSetAdminEventEvent packs the event into the appropriate arguments for setAdminEvent.
// It returns topic hashes and the encoded non-indexed data.
func PackSetAdminEventEvent(contractAbi abi.ABI, sender common.Address, admin common.Address) ([]common.Hash, []byte, error) {
	return contractAbi.PackEvent("setAdminEvent", sender, admin)
}

// UnpackSetAdminEventEvent won't be generated because the event does not have any non-indexed data.

// PackSetEnabledEventEvent packs the event into the appropriate arguments for setEnabledEvent.
// It returns topic hashes and the encoded non-indexed data.
func PackSetEnabledEventEvent(contractAbi abi.ABI, sender common.Address, enabledAddr common.Address) ([]common.Hash, []byte, error) {
	return contractAbi.PackEvent("setEnabledEvent", sender, enabledAddr)
}

// UnpackSetEnabledEventEvent won't be generated because the event does not have any non-indexed data.

// PackSetManagerEventEvent packs the event into the appropriate arguments for setManagerEvent.
// It returns topic hashes and the encoded non-indexed data.
func PackSetManagerEventEvent(contractAbi abi.ABI, sender common.Address, manager common.Address) ([]common.Hash, []byte, error) {
	return contractAbi.PackEvent("setManagerEvent", sender, manager)
}

// UnpackSetManagerEventEvent won't be generated because the event does not have any non-indexed data.

// PackSetNoneEventEvent packs the event into the appropriate arguments for setNoneEvent.
// It returns topic hashes and the encoded non-indexed data.
func PackSetNoneEventEvent(contractAbi abi.ABI, sender common.Address, none common.Address) ([]common.Hash, []byte, error) {
	return contractAbi.PackEvent("setNoneEvent", sender, none)
}

// UnpackSetNoneEventEvent won't be generated because the event does not have any non-indexed data.
