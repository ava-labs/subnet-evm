// (c) 2023 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"encoding/json"

	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ethereum/go-ethereum/common"
)

type ChainConfigPrecompiles map[string]config.Config

func (ccp *ChainConfigPrecompiles) GetConfigByAddress(address common.Address) (config.Config, bool) {
	module, ok := modules.GetPrecompileModuleByAddress(address)
	if !ok {
		return nil, false
	}
	key := module.ConfigKey
	config, ok := (*ccp)[key]
	return config, ok
}

// UnmarshalJSON parses the JSON-encoded data into the ChainConfigPrecompiles.
// ChainConfigPrecompiles is a map of precompile module keys to their
// configuration.
func (ccp *ChainConfigPrecompiles) UnmarshalJSON(data []byte) error {
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*ccp = make(ChainConfigPrecompiles)
	for _, module := range modules.RegisteredModules() {
		key := module.ConfigKey
		if value, ok := raw[key]; ok {
			conf := module.NewConfig()
			err := json.Unmarshal(value, conf)
			if err != nil {
				return err
			}
			(*ccp)[key] = conf
		}
	}
	return nil
}
