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
	time := GetHeaderExtra(b.Header()).TimeMilliseconds
	if time == nil {
		return nil
	}
	cp := *time
	return &cp
}

func BlockMinDelayExcess(b *ethtypes.Block) *uint64 {
	e := GetHeaderExtra(b.Header()).MinDelayExcess
	if e == nil {
		return nil
	}
	cp := *e
	return &cp
}
