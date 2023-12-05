// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package handshake

import (
	"testing"

	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contracts/nativeminter"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"
)

// Checks if messages have the same hash
//
// `message` is the simulation of a configuration being parsed from the local
// config. `message2` is parsing a message being exchanged through the network
// (a foreign config), and `message3` is the the deserialization and
// serialization of the foreign config. All 3 instances should have the same
// hashing, depite maybe not being identical (some configurations may be in a
// different order, but our hashing algorithm is resilient to those changes,
// thanks for our serialization library, which produces always the same output.
func assertConversions(t *testing.T, message *UpgradeConfigMessage, err error) {
	require.NoError(t, err)

	config, err := UpgradeConfigFromBytes(message.Bytes())
	require.NoError(t, err)

	message2, err := NewUpgradeConfigMessage(config)
	require.NoError(t, err)

	config3, err := UpgradeConfigFromBytes(message2.Bytes())
	require.NoError(t, err)

	message3, err := NewUpgradeConfigMessage(config3)
	require.NoError(t, err)

	require.Equal(t, config, config3)
	require.Equal(t, message.hash, message2.hash)
	require.Equal(t, message2.hash, message3.hash)
}

func TestSerialize(t *testing.T) {
	var t0 uint64 = 0
	var t1 uint64 = 1
	message, err := NewUpgradeConfigMessage(&params.UpgradeConfig{
		PrecompileUpgrades: []params.PrecompileUpgrade{
			{
				Config: nativeminter.NewConfig(&t0, nil, nil, nil, nil), // enable at genesis
			},
			{
				Config: nativeminter.NewDisableConfig(&t1), // disable at timestamp 1
			},
		},
	})
	assertConversions(t, message, err)
}

func TestWithAddress(t *testing.T) {
	var t0 uint64 = 1
	var t1 uint64 = 11
	message, err := NewUpgradeConfigMessage(&params.UpgradeConfig{
		PrecompileUpgrades: []params.PrecompileUpgrade{
			{
				Config: nativeminter.NewConfig(&t0, []common.Address{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000020")),
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000030")),
				}, []common.Address{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000040")),
				}, []common.Address{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000050")),
				}, nil), // enable at genesis
			},
			{
				Config: nativeminter.NewDisableConfig(&t1), // disable at timestamp 1
			},
		},
	})
	assertConversions(t, message, err)
}

func TestWithAddressAndMint(t *testing.T) {
	var t0 uint64 = 2
	var t1 uint64 = 1001
	message, err := NewUpgradeConfigMessage(&params.UpgradeConfig{
		PrecompileUpgrades: []params.PrecompileUpgrade{
			{
				Config: nativeminter.NewConfig(&t0, []common.Address{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000020")),
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000030")),
				}, []common.Address{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000040")),
				}, []common.Address{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000050")),
				}, map[common.Address]*math.HexOrDecimal256{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000010")): math.NewHexOrDecimal256(64),
				}), // enable at genesis
			},
			{
				Config: nativeminter.NewDisableConfig(&t1), // disable at timestamp 1
			},
		},
	})
	assertConversions(t, message, err)
}

func Test(t *testing.T) {
	var t0 uint64 = 2
	var t1 uint64 = 1001
	config := params.UpgradeConfig{
		StateUpgrades: []params.StateUpgrade{
			{
				BlockTimestamp:       &t0,
				StateUpgradeAccounts: map[common.Address]params.StateUpgradeAccount{},
			},
			{
				BlockTimestamp: &t1,
				StateUpgradeAccounts: map[common.Address]params.StateUpgradeAccount{
					common.HexToAddress("00c1f1a2"): params.StateUpgradeAccount{
						Code:          common.Hex2Bytes("0fc01e"),
						BalanceChange: math.NewHexOrDecimal256(643),
					},
				},
			},
		},
	}
	message, err := NewUpgradeConfigMessage(&config)
	require.NoError(t, err)
	configFromBytes, err := UpgradeConfigFromBytes(message.Bytes())
	require.NoError(t, err)
	require.Equal(t, &config, configFromBytes)
}
