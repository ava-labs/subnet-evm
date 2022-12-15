// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
)

// SafeSumUint256 adds two big.Ints and returns the result and whether the result
// is less than or equal to 2^256-1 (MaxBig256).
func SafeSumUint256(a, b *big.Int) (*big.Int, bool) {
	sum := new(big.Int).Add(a, b)
	return sum, sum.Cmp(math.MaxBig256) <= 0
}
