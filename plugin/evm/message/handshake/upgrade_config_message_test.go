// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package handshake

import (
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contracts/deployerallowlist"
	"github.com/ava-labs/subnet-evm/precompile/contracts/feemanager"
	"github.com/ava-labs/subnet-evm/precompile/contracts/nativeminter"
	"github.com/ava-labs/subnet-evm/precompile/contracts/rewardmanager"
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

	config, err := NewUpgradeConfigMessageFromBytes(message.Bytes())
	require.NoError(t, err)

	message2, err := NewUpgradeConfigMessage(config)
	require.NoError(t, err)

	config3, err := NewUpgradeConfigMessageFromBytes(message2.Bytes())
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
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000090")): math.NewHexOrDecimal256(6402100201021),
				}), // enable at genesis
			},
			{
				Config: nativeminter.NewDisableConfig(&t1), // disable at timestamp 1
			},
		},
	})
	assertConversions(t, message, err)
}

func TestWithAddressFeeMinter(t *testing.T) {
	var t0 uint64 = 2
	var t1 uint64 = 1001
	var validFeeConfig = commontype.FeeConfig{
		GasLimit:        big.NewInt(8_000_000),
		TargetBlockRate: 2, // in seconds

		MinBaseFee:               big.NewInt(25_000_000_000),
		TargetGas:                big.NewInt(15_000_000),
		BaseFeeChangeDenominator: big.NewInt(36),

		MinBlockGasCost:  big.NewInt(0),
		MaxBlockGasCost:  big.NewInt(1_000_000),
		BlockGasCostStep: big.NewInt(200_000),
	}

	message, err := NewUpgradeConfigMessage(&params.UpgradeConfig{
		PrecompileUpgrades: []params.PrecompileUpgrade{
			{
				Config: feemanager.NewConfig(&t0, []common.Address{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000020")),
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000030")),
				}, []common.Address{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000040")),
				}, []common.Address{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000050")),
				}, &validFeeConfig), // enable at genesis
			},
			{
				Config: feemanager.NewDisableConfig(&t1), // disable at timestamp 1
			},
		},
	})
	assertConversions(t, message, err)
}

func TestWithDepoyerAllowList(t *testing.T) {
	var t0 uint64 = 2
	var t1 uint64 = 1001

	message, err := NewUpgradeConfigMessage(&params.UpgradeConfig{
		PrecompileUpgrades: []params.PrecompileUpgrade{
			{
				Config: deployerallowlist.NewConfig(&t0, []common.Address{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000020")),
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000030")),
				}, []common.Address{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000040")),
				}, []common.Address{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000050")),
				}), // enable at genesis
			},
			{
				Config: feemanager.NewDisableConfig(&t1), // disable at timestamp 1
			},
		},
	})
	assertConversions(t, message, err)
}

func TestWithRewardManager(t *testing.T) {
	var t0 uint64 = 2
	var t1 uint64 = 1001

	message, err := NewUpgradeConfigMessage(&params.UpgradeConfig{
		PrecompileUpgrades: []params.PrecompileUpgrade{
			{
				Config: rewardmanager.NewConfig(&t0, []common.Address{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000020")),
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000030")),
				}, []common.Address{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000040")),
				}, []common.Address{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000050")),
				}, nil), // enable at genesis
			},
			{
				Config: feemanager.NewDisableConfig(&t1), // disable at timestamp 1
			},
		},
	})
	assertConversions(t, message, err)
}

func TestWithRewardManagerWithNil(t *testing.T) {
	var t0 uint64 = 2
	var t1 uint64 = 1001

	message, err := NewUpgradeConfigMessage(&params.UpgradeConfig{
		PrecompileUpgrades: []params.PrecompileUpgrade{
			{
				Config: rewardmanager.NewConfig(&t0, []common.Address{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000020")),
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000030")),
				}, []common.Address{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000040")),
				}, []common.Address{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000050")),
				}, &rewardmanager.InitialRewardConfig{
					AllowFeeRecipients: true,
				}),
			},
			{
				Config: feemanager.NewDisableConfig(&t1), // disable at timestamp 1
			},
		},
	})
	assertConversions(t, message, err)
}
