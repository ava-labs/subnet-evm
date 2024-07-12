// Package dstestbindings contains generated Solidity bindings for the
// DSTest testing contract.
package dstestbindings

import (
	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ethereum/go-ethereum/common"
)

//go:generate sh -c "solc --evm-version=paris --base-path=../ds-test/src --combined-json=abi ../ds-test/src/test.sol | abigen --pkg dstestbindings --combined-json=- | sed -E 's,github.com/ethereum/go-ethereum/(accounts|core)/,github.com/ava-labs/subnet-evm/\\1/,' > generated.go"

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
