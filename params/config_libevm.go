// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ethereum/go-ethereum/common"
	gethparams "github.com/ethereum/go-ethereum/params"
)

func do_init() any {
	getter = gethparams.RegisterExtras(gethparams.Extras[ChainConfigExtra, RulesExtra]{
		ReuseJSONRoot: true, // Reuse the root JSON input when unmarshalling the extra payload.
		NewRules:      constructRulesExtra,
	})
	return nil
}

var getter gethparams.ExtraPayloadGetter[ChainConfigExtra, RulesExtra]

// constructRulesExtra acts as an adjunct to the [params.ChainConfig.Rules]
// method. Its primary purpose is to construct the extra payload for the
// [params.Rules] but it MAY also modify the [params.Rules].
func constructRulesExtra(c *gethparams.ChainConfig, r *gethparams.Rules, cEx *ChainConfigExtra, blockNum *big.Int, isMerge bool, timestamp uint64) *RulesExtra {
	if cEx == nil {
		return &RulesExtra{}
	}
	var rules RulesExtra
	rules.AvalancheRules = cEx.GetAvalancheRules(timestamp)

	// Initialize the stateful precompiles that should be enabled at [blockTimestamp].
	rules.ActivePrecompiles = make(map[common.Address]precompileconfig.Config)
	rules.Predicaters = make(map[common.Address]precompileconfig.Predicater)
	rules.AccepterPrecompiles = make(map[common.Address]precompileconfig.Accepter)
	for _, module := range modules.RegisteredModules() {
		if config := cEx.getActivePrecompileConfig(module.Address, timestamp); config != nil && !config.IsDisabled() {
			rules.ActivePrecompiles[module.Address] = config
			if predicater, ok := config.(precompileconfig.Predicater); ok {
				rules.Predicaters[module.Address] = predicater
			}
			if precompileAccepter, ok := config.(precompileconfig.Accepter); ok {
				rules.AccepterPrecompiles[module.Address] = precompileAccepter
			}
		}
	}

	return &rules
}

// FromChainConfig returns the extra payload carried by the ChainConfig.
func FromChainConfig(c *gethparams.ChainConfig) *ChainConfigExtra {
	return getter.FromChainConfig(c)
}

// FromRules returns the extra payload carried by the Rules.
func FromRules(r *gethparams.Rules) *RulesExtra {
	return getter.FromRules(r)
}
