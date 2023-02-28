// (c) 2023 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package stateupgrade

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// StateDB is the interface for accessing EVM state in state upgrades
type StateDB interface {
	SetState(common.Address, common.Hash, common.Hash)
	SetCode(common.Address, []byte)
	AddBalance(common.Address, *big.Int)

	CreateAccount(common.Address)
	Exist(common.Address) bool

	Snapshot() int
	RevertToSnapshot(int)
}

// BlockContext defines an interface that provides information about the
// block that activates the state upgrade.
type BlockContext interface {
	Number() *big.Int
	Timestamp() *big.Int
}

// AccessibleState defines the interface exposed to state upgrades
type AccessibleState interface {
	CreateAt(contractAddr common.Address, callerAddr common.Address, code []byte, gas uint64, value *big.Int) ([]byte, common.Address, uint64, error)
}
