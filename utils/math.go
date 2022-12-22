// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
)

// SafeSumUint256 adds two big.Ints and returns the result and whether the result
// is overflowed.
// If the result is overflowed, the result will be set to math.MaxBig256.
func SafeSumUint256(a, b *big.Int) (*big.Int, bool) {
	sum := new(big.Int).Add(a, b)
	if sum.Cmp(math.MaxBig256) > 0 {
		return math.MaxBig256, true
	}
	return sum, false
}
