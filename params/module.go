// (c) 2023 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ethereum/go-ethereum/common"
)

var RegisteredModules = func() []Module { return nil }

type Module struct {
	ConfigKey  string
	Address    common.Address
	MakeConfig func() precompileconfig.Config
}

func getPrecompileModule(key string) (Module, bool) {
	for _, module := range RegisteredModules() {
		if module.ConfigKey == key {
			return module, true
		}
	}
	return Module{}, false
}

func getPrecompileModuleByAddress(addr common.Address) (Module, bool) {
	for _, module := range RegisteredModules() {
		if module.Address == addr {
			return module, true
		}
	}
	return Module{}, false
}
