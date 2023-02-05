// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Defines the stateless interface for unmarshalling an arbitrary config of a precompile
package config

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// params package
// - IMPORTANT - cannot import precompile or statedb directly
// imports precompile/config to unmarshal relevant configs and nothing else
// this means that we also need to stop using AvalancheRules to hold the precompile contracts
// switch instead to map address -> config
// ActivePrecompiles was previously used to check if a precompile was enabled at that address, so we can still use this with the new map
// also used in core/vm/evm.go to get the precompile contract
// we will need to import the precompile/exec package to look up the precompile at that address now.

// update package
// import params package and precompile/exec
// used to Configure precompiles

var (
	registry               = make(map[string]Factory)
	addressToConfigFactory = make(map[common.Address]Config)
	addresses              = make([]common.Address, 0)
	configs                = make([]Config, 0)
)

func RegisterConfig(name string, addr common.Address, configFactory Factory) error {
	_, exists := registry[name]
	if exists {
		return fmt.Errorf("cannot register duplicate config with the name: %s", name)
	}
	_, exists = addressToConfigFactory[addr]
	if exists {
		return fmt.Errorf("cannot register duplicate config with address: %s", addr)
	}
	registry[name] = configFactory
	addressToConfigFactory[addr] = configFactory.NewConfig()
	addresses = append(addresses, addr)
	configs = append(configs, configFactory.NewConfig())
	return nil
}

func GetNewConfig(name string) (Config, bool) {
	config, ok := registry[name]
	if !ok {
		return nil, false
	}
	return config.NewConfig(), true
}

func GetAddresses() []common.Address { return addresses }

func GetConfigs() []Config { return configs }

// TODO: add lookups as needed
