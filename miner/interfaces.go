// (c) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package miner

import (
	"math/big"

	"github.com/ava-labs/coreth/core/txpool"
	"github.com/ethereum/go-ethereum/common"
)

type TxPool interface {
	Locals() []common.Address
	PendingWithBaseFee(enforceTips bool, baseFee *big.Int) map[common.Address][]*txpool.LazyTransaction
}
