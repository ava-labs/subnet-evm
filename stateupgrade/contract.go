// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package stateupgrade

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// BlockContext defines an interface that provides information to a stateful stateupgrade
// about the block that activates the upgrade. The stateupgrade can access this information
// to initialize its state.
type BlockContext interface {
	Number() *big.Int
	Timestamp() *big.Int
}

// ChainConfig defines an interface that provides information to a stateful stateupgrade
// about the chain configuration. The stateupgrade can access this information to initialize
// its state.
type ChainConfig interface {

}

// StateDB is the interface for accessing EVM state
type StateDB interface {
	GetState(common.Address, common.Hash) common.Hash
	SetState(common.Address, common.Hash, common.Hash)

	SetCode(common.Address, []byte)
	GetCode(common.Address) []byte

	SetNonce(common.Address, uint64)
	GetNonce(common.Address) uint64

	GetBalance(common.Address) *big.Int
	AddBalance(common.Address, *big.Int)
	SubBalance(common.Address, *big.Int)

	CreateAccount(common.Address)
	Exist(common.Address) bool

	AddLog(addr common.Address, topics []common.Hash, data []byte, blockNumber uint64)

	Suicide(common.Address) bool
	Finalise(deleteEmptyObjects bool)
}
