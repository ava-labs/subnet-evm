package modules

import (
	"encoding/json"
	"fmt"

	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ethereum/go-ethereum/common"
)

func InitChainConfig(c *params.ChainConfig) error {
	if c.GenesisPrecompiles == nil {
		c.GenesisPrecompiles = make(params.Precompiles)
	}
	for key, raw := range c.LazyUnmarshalData {
		mod, ok := GetPrecompileModule(key)
		if !ok {
			continue
		}

		conf := mod.MakeConfig()
		if err := json.Unmarshal(raw, conf); err != nil {
			return fmt.Errorf("unmarshal %T: %v", conf, err)
		}
		c.GenesisPrecompiles[key] = conf
	}
	return nil
}

func InitChainRules(r *params.Rules) {
	// Initialize the stateful precompiles that should be enabled at [blockTimestamp].
	r.ActivePrecompiles = make(map[common.Address]precompileconfig.Config)
	r.Predicaters = make(map[common.Address]precompileconfig.Predicater)
	r.AccepterPrecompiles = make(map[common.Address]precompileconfig.Accepter)

	// TODO(arr4n): range over config keys instead
	for _, module := range RegisteredModules() {
		if config := c.getActivePrecompileConfig(module.Address, timestamp); config != nil && !config.IsDisabled() {
			r.ActivePrecompiles[module.Address] = config
			if predicater, ok := config.(precompileconfig.Predicater); ok {
				r.Predicaters[module.Address] = predicater
			}
			if precompileAccepter, ok := config.(precompileconfig.Accepter); ok {
				r.AccepterPrecompiles[module.Address] = precompileAccepter
			}
		}
	}
}
