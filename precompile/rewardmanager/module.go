// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package rewardmanager

import (
	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ethereum/go-ethereum/common"
)

// ConfigKey is the key used in json config files to specify this precompile config.
// must be unique across all precompiles.
const ConfigKey = "rewardManagerConfig"

var ContractAddress = common.HexToAddress("0x0200000000000000000000000000000000000004")

type Module struct{}

func (Module) Key() string {
	return ConfigKey
}

// Address returns the address of the reward manager.
func (Module) Address() common.Address {
	return ContractAddress
}

func (Module) NewConfig() config.Config {
	return &RewardManagerConfig{}
}

func (Module) Executor() contract.Execution {
	return &Executor{}
}
