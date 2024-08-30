// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package modules

import (
	"bytes"

	"github.com/ava-labs/subnet-evm/precompile/interfaces"
	"github.com/ethereum/go-ethereum/common"
)

type Module struct {
	// ConfigKey is the key used in json config files to specify this precompile config.
	ConfigKey string
	// Address returns the address where the stateful precompile is accessible.
	Address common.Address
	// Contract returns a thread-safe singleton that can be used as the StatefulPrecompiledContract when
	// this config is enabled.
	Contract interfaces.StatefulPrecompiledContract
	// Configurator is used to configure the stateful precompile when the config is enabled.
	interfaces.Configurator
}

type moduleArray []Module

func (u moduleArray) Len() int {
	return len(u)
}

func (u moduleArray) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (m moduleArray) Less(i, j int) bool {
	return bytes.Compare(m[i].Address.Bytes(), m[j].Address.Bytes()) < 0
}
