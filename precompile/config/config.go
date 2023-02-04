// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Defines the stateless interface for unmarshalling an arbitrary config of a precompile
package config

import (
	"encoding/json"
	"fmt"
	"math/big"
)

// new breakdown:
// precompile/config - does not need to import anything sits on an island

// precompile/exec - imports statedb and should be able to import from params and core/types
// precompile/module <- config and exec packages

// precompile/txallowlist defines a module that can create a new config, Configure, and return the Contract instance
// imports from both config and exec package
// it should be alright for this to import the params package and the statedb package
// registers itself by calling precompileModule.Register

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

var registeredConfigs = make(map[string]ConfigFactory)

func RegisterConfig(name string, configFactory ConfigFactory) error {
	_, exists := registeredConfigs[name]
	if exists {
		return fmt.Errorf("cannot register duplicate config with the name: %s", name)
	}
	registeredConfigs[name] = configFactory
	return nil
}

type ConfigFactory interface {
	New() Config
}

// StatefulPrecompileConfig defines the interface for a stateful precompile to
// be enabled via a network upgrade.
type Config interface {
	// Timestamp returns the timestamp at which this stateful precompile should be enabled.
	// 1) 0 indicates that the precompile should be enabled from genesis.
	// 2) n indicates that the precompile should be enabled in the first block with timestamp >= [n].
	// 3) nil indicates that the precompile is never enabled.
	Timestamp() *big.Int
	// IsDisabled returns true if this network upgrade should disable the precompile.
	IsDisabled() bool
	// Equal returns true if the provided argument configures the same precompile with the same parameters.
	Equal(Config) bool
	// Verify is called on startup and an error is treated as fatal. Configure can assume the Config has passed verification.
	Verify() error
}

type Configs map[string]Config

// UnmarshalJSON parses the JSON-encoded data into the ChainConfigPrecompiles.
// ChainConfigPrecompiles is a map of precompile module keys to their
// configuration.
func (ccp *Configs) UnmarshalJSON(data []byte) error {
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*ccp = make(Configs)
	for name, configFactory := range registeredConfigs {
		if value, ok := raw[name]; ok {
			conf := configFactory.New()
			err := json.Unmarshal(value, conf)
			if err != nil {
				return err
			}
			(*ccp)[name] = conf
		}
	}
	return nil
}
