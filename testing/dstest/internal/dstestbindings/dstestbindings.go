// Package dstestbindings contains generated Solidity bindings for the
// DSTest testing contract.
package dstestbindings

import (
	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ethereum/go-ethereum/common"
)

//go:generate sh -c "go run $(git rev-parse --show-toplevel)/scripts/abigen --solc.base-path=../ds-test/src --solc.output=abi --abigen.pkg=dstestbindings ../ds-test/src/test.sol > generated.go"

var (
	parsed *abi.ABI
	bound  *bind.BoundContract
)

func init() {
	a, err := DSTestMetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	parsed = a
	bound = bind.NewBoundContract(common.Address{}, *parsed, nil, nil, nil)
}

// EventByID is a convenience wrapper around [abi.ABI.EventByID], returning the
// DSTest event with the corresponding topic.
func EventByID(topic common.Hash) (*abi.Event, error) {
	return parsed.EventByID(topic)
}

// UnpackLogIntoMap is a convenience wrapper around
// [bind.BoundContract.UnpackLogIntoMap], unpacking the log data of the named
// event.
func UnpackLogIntoMap(out map[string]any, event string, log types.Log) error {
	return bound.UnpackLogIntoMap(out, event, log)
}
