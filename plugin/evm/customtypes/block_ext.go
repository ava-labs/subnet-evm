// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package customtypes

import (
	"math/big"

	ethtypes "github.com/ava-labs/libevm/core/types"
)

func BlockExtDataGasUsed(b *ethtypes.Block) *big.Int {
	used := GetHeaderExtra(b.Header()).ExtDataGasUsed
	if used == nil {
		return nil
	}
	return new(big.Int).Set(used)
}

func BlockGasCost(b *ethtypes.Block) *big.Int {
	cost := GetHeaderExtra(b.Header()).BlockGasCost
	if cost == nil {
		return nil
	}
	return new(big.Int).Set(cost)
}
