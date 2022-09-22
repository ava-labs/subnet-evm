// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"math/big"

	"github.com/holiman/uint256"
)

func SafeSumUint256(a, b *big.Int) (*big.Int, bool) {
	sum := new(big.Int).Add(a, b)
	if _, overflow := uint256.FromBig(sum); overflow {
		return big.NewInt(0), false
	}
	return sum, true
}
