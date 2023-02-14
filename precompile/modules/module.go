// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package modules

import (
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
)

type Module struct {
	// ConfigKey is the key used in json config files to specify this precompile config.
	ConfigKey string
	// Address returns the address where the stateful precompile is accessible.
	Address common.Address
	// Contract returns a thread-safe singleton that can be used as the StatefulPrecompiledContract when
	// this config is enabled.
	Contract contract.StatefulPrecompiledContract
	contract.Configurator
}
