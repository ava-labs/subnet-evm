// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Defines the stateless interface for unmarshalling an arbitrary config of a precompile
package config

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Factory interface {
	NewConfig() Config
}

// StatefulPrecompileConfig defines the interface for a stateful precompile to
// be enabled via a network upgrade.
type Config interface {
	Key() string
	Address() common.Address
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

func (ccp Configs) GetConfigByAddress(addr common.Address) (Config, bool) {
	config, ok := addressToConfigFactory[addr]
	if !ok {
		return nil, false
	}
	config, exists := ccp[config.Key()]
	return config, exists
}

// UnmarshalJSON parses the JSON-encoded data into the ChainConfigPrecompiles.
// ChainConfigPrecompiles is a map of precompile module keys to their
// configuration.
func (ccp *Configs) UnmarshalJSON(data []byte) error {
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*ccp = make(Configs)
	for name, configFactory := range registry {
		if value, ok := raw[name]; ok {
			conf := configFactory.NewConfig()
			err := json.Unmarshal(value, conf)
			if err != nil {
				return err
			}
			(*ccp)[name] = conf
		}
	}
	return nil
}
