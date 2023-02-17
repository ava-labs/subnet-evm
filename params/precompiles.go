// (c) 2023 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"encoding/json"

	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ava-labs/subnet-evm/precompile/modules"
)

type Precompiles map[string]config.StatefulPrecompileConfig

// UnmarshalJSON parses the JSON-encoded data into the ChainConfigPrecompiles.
// ChainConfigPrecompiles is a map of precompile module keys to their
// configuration.
func (ccp *Precompiles) UnmarshalJSON(data []byte) error {
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*ccp = make(Precompiles)
	for _, module := range modules.RegisteredModules() {
		key := module.ConfigKey
		if value, ok := raw[key]; ok {
			conf := module.NewConfig()
			if err := json.Unmarshal(value, conf); err != nil {
				return err
			}
			(*ccp)[key] = conf
		}
	}
	return nil
}
