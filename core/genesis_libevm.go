// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package core

import (
	"encoding/json"

	"github.com/ava-labs/subnet-evm/params"
	gethparams "github.com/ethereum/go-ethereum/params"
)

type ChainConfig params.ChainConfig

func (c *ChainConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal((*params.ChainConfig)(c))
}

func (c *ChainConfig) UnmarshalJSON(input []byte) error {
	var tmp gethparams.ChainConfig
	if err := json.Unmarshal(input, &tmp); err != nil {
		return err
	}
	*c = (ChainConfig)(*params.FromChainConfig(&tmp))
	return nil
}
