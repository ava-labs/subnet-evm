// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"math/big"

	"github.com/holiman/uint256"
)

// SafeSumUint256 returns the sum of a and b and a boolean indicating if the sum
// operation was successful. If the sum operation overflows, the boolean will be
// false and the sum will be 0.
func SafeSumUint256(a, b *big.Int) (*big.Int, bool) {
	sum := new(big.Int).Add(a, b)
	if _, overflow := uint256.FromBig(sum); overflow {
		return big.NewInt(0), false
	}
	return sum, true
}
