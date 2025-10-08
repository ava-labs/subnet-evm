// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package customtypes

import (
	"math/big"

	ethtypes "github.com/ava-labs/libevm/core/types"
)

func BlockGasCost(b *ethtypes.Block) *big.Int {
	cost := GetHeaderExtra(b.Header()).BlockGasCost
	if cost == nil {
		return nil
	}
	return new(big.Int).Set(cost)
}

func BlockTimeMilliseconds(b *ethtypes.Block) *uint64 {
	var time *uint64
	if t := GetHeaderExtra(b.Header()).TimeMilliseconds; t != nil {
		time = new(uint64)
		*time = *t
	}
	return time
}

func BlockMinDelayExcess(b *ethtypes.Block) *uint64 {
	var excess *uint64
	if e := GetHeaderExtra(b.Header()).MinDelayExcess; e != nil {
		excess = new(uint64)
		*excess = *e
	}
	return excess
}
