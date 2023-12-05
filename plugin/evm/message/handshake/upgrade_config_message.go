// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package handshake

import (
	"encoding/json"

	"github.com/ava-labs/subnet-evm/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type UpgradeConfigMessage struct {
	bytes []byte
	hash  common.Hash
}

func (u *UpgradeConfigMessage) Bytes() []byte {
	return u.bytes
}

func (u *UpgradeConfigMessage) ID() common.Hash {
	return u.hash
}

// Attempts to parse a params.UpgradeConfig from a []byte
//
// The function returns a reference of *params.UpgradeConfig
func UpgradeConfigFromBytes(bytes []byte) (*params.UpgradeConfig, error) {
	var upgradeConfig params.UpgradeConfig
	err := json.Unmarshal(bytes, &upgradeConfig)
	if err != nil {
		return nil, err
	}

	return &upgradeConfig, nil
}

// Encodes any object to JSON with a deterministic output. All the keys, even
// into inner objects, are sorted alphabetically.
func DeterministicJsonEncoding(object interface{}) ([]byte, error) {
	bytes, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	var temporaryObject interface{}
	err = json.Unmarshal(bytes, &temporaryObject)
	if err != nil {
		return nil, err
	}
	return json.Marshal(temporaryObject)
}

// Wraps an instance of *params.UpgradeConfig
//
// This function returns the serialized UpgradeConfig, ready to be send over to
// other peers. The struct also includes a hash of the content, ready to be used
// as part of the handshake protocol.
//
// Since params.UpgradeConfig should never change without a node reloading, it
// is safe to call this function once and store its output globally to re-use
// multiple times
func NewUpgradeConfigMessage(config *params.UpgradeConfig) (*UpgradeConfigMessage, error) {
	bytes, err := DeterministicJsonEncoding(config)
	if err != nil {
		return nil, err
	}

	hash := crypto.Keccak256Hash(bytes)
	return &UpgradeConfigMessage{
		bytes: bytes,
		hash:  hash,
	}, nil
}
