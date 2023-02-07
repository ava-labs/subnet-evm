// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nativeminter

import (
	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/precompile/execution"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ethereum/go-ethereum/common"
)

// ConfigKey is the key used in json config files to specify this precompile config.
// must be unique across all precompiles.
const ConfigKey = "contractNativeMinterConfig"

var ContractAddress = common.HexToAddress("0x0200000000000000000000000000000000000001")

func init() {
	modules.RegisterModule(&Module{})
}

type Module struct{}

func (Module) Key() string {
	return ConfigKey
}

// Address returns the address of the native minter.
func (Module) Address() common.Address {
	return ContractAddress
}

func (Module) NewConfig() config.Config {
	return &ContractNativeMinterConfig{}
}

func (Module) Executor() execution.Execution {
	return &Executor{}
}
