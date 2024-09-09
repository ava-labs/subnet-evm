// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"math/big"

	gethparams "github.com/ethereum/go-ethereum/params"
)

func init() {
	getter = gethparams.RegisterExtras(gethparams.Extras[ChainConfig, RulesExtra]{
		NewRules: constructRulesExtra,
	})
}

var getter gethparams.ExtraPayloadGetter[ChainConfig, RulesExtra]

// constructRulesExtra acts as an adjunct to the [params.ChainConfig.Rules]
// method. Its primary purpose is to construct the extra payload for the
// [params.Rules] but it MAY also modify the [params.Rules].
func constructRulesExtra(c *gethparams.ChainConfig, r *gethparams.Rules, cEx *ChainConfig, blockNum *big.Int, isMerge bool, timestamp uint64) *RulesExtra {
	return &RulesExtra{}
}

// FromChainConfig returns the extra payload carried by the ChainConfig.
func FromChainConfig(c *gethparams.ChainConfig) *ChainConfig {
	return getter.FromChainConfig(c)
}

// FromRules returns the extra payload carried by the Rules.
func FromRules(r *gethparams.Rules) *RulesExtra {
	return getter.FromRules(r)
}
