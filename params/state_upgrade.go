// (c) 2023 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

type StateUpgrade struct {
	blockTimestamp *big.Int

	// Adds the specified amount to the balance of the specified address
	AddToBalance map[common.Address]*math.HexOrDecimal256 `json:"addToBalance,omitempty"`

	// Sets the specified storage slots of the specified addresses
	// to the given values. Note that the value of common.Hash{} will
	// remove the storage key
	SetStorage map[common.Address]map[common.Hash]common.Hash `json:"setStorage,omitempty"`

	// Sets the code of the specified contract to the given value
	SetCode map[common.Address][]byte `json:"setCode,omitempty"`

	// Deploys contracts with the specified creation bytecode to the
	// specified addresses, instead of the normal rules for deriving
	// the address of a created contract.
	DeployContractTo []ContractDeploy `json:"deployContractTo,omitempty"`
}

type ContractDeploy struct {
	DeployTo common.Address `json:"deployTo,omitempty"` // The address to deploy the contract to
	Caller   common.Address `json:"caller,omitempty"`   // The address of the caller
	Input    []byte         `json:"input,omitempty"`    // The input bytecode to create the contract
	Gas      uint64         `json:"gas,omitempty"`      // The gas to use when creating the contract
	Value    *big.Int       `json:"value,omitempty"`    // The value to send when creating the contract
}

func (s *StateUpgrade) Equal(other *StateUpgrade) bool {
	return reflect.DeepEqual(s, other)
}
