// (c) 2019-2020, Ava Labs, Inc.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

import (
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/commontype"
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
func assertConversions(t *testing.T, originalConfig *UpgradeConfig) {
	bytes, err := originalConfig.MarshalBinary()
	require.NoError(t, err)

	deserializedConfig := UpgradeConfig{}
	require.NoError(t, deserializedConfig.UnmarshalBinary(bytes))

	twiceDeserialized := UpgradeConfig{}
	newBytes, err := deserializedConfig.MarshalBinary()
	require.NoError(t, err)
	require.NoError(t, twiceDeserialized.UnmarshalBinary(newBytes))

	hash1, err := originalConfig.Hash()
	require.NoError(t, err)
	hash2, err := deserializedConfig.Hash()
	require.NoError(t, err)
	hash3, err := twiceDeserialized.Hash()
	require.NoError(t, err)

	require.Equal(t, deserializedConfig, twiceDeserialized)
	require.Equal(t, hash1, hash2)
	require.Equal(t, hash2, hash3)
}

func TestSerialize(t *testing.T) {
	var t0 uint64 = 0
	var t1 uint64 = 1
	config := UpgradeConfig{
		PrecompileUpgrades: []PrecompileUpgrade{
			{
				Config: nativeminter.NewConfig(&t0, nil, nil, nil, nil), // enable at genesis
			},
			{
				Config: nativeminter.NewDisableConfig(&t1), // disable at timestamp 1
			},
		},
	}
	assertConversions(t, &config)
}

func TestWithAddress(t *testing.T) {
	var t0 uint64 = 1
	var t1 uint64 = 11
	config := UpgradeConfig{
		PrecompileUpgrades: []PrecompileUpgrade{
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
	}
	assertConversions(t, &config)
}

func TestWithAddressAndMint(t *testing.T) {
	var t0 uint64 = 2
	var t1 uint64 = 1001
	config := UpgradeConfig{
		PrecompileUpgrades: []PrecompileUpgrade{
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
	}
	assertConversions(t, &config)
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

	config := UpgradeConfig{
		PrecompileUpgrades: []PrecompileUpgrade{
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
	}
	assertConversions(t, &config)
}

func TestWithDepoyerAllowList(t *testing.T) {
	var t0 uint64 = 2
	var t1 uint64 = 1001

	config := UpgradeConfig{
		PrecompileUpgrades: []PrecompileUpgrade{
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
	}
	assertConversions(t, &config)
}

func TestWithRewardManager(t *testing.T) {
	var t0 uint64 = 2
	var t1 uint64 = 1001

	config := UpgradeConfig{
		PrecompileUpgrades: []PrecompileUpgrade{
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
	}
	assertConversions(t, &config)
}

func TestWithRewardManagerWithNil(t *testing.T) {
	var t0 uint64 = 2
	var t1 uint64 = 1001

	config := UpgradeConfig{
		PrecompileUpgrades: []PrecompileUpgrade{
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
	}
	assertConversions(t, &config)
}

func TestStateUpgrades(t *testing.T) {
	var t0 uint64 = 2
	var t1 uint64 = 1001
	config := UpgradeConfig{
		StateUpgrades: []StateUpgrade{
			{
				BlockTimestamp: &t0,
				StateUpgradeAccounts: map[common.Address]StateUpgradeAccount{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000050")): StateUpgradeAccount{
						Code:          []byte{1, 2, 3, 4, 5, 6},
						BalanceChange: math.NewHexOrDecimal256(99),
						Storage: map[common.Hash]common.Hash{
							common.BytesToHash([]byte{1, 2, 4, 5}): common.BytesToHash([]byte{1, 2, 3}),
						},
					},
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000001000")): StateUpgradeAccount{
						Code:          []byte{1, 2, 9, 93, 4, 5, 6},
						BalanceChange: math.NewHexOrDecimal256(92312319),
						Storage: map[common.Hash]common.Hash{
							common.BytesToHash([]byte{11, 21, 99, 5}): common.BytesToHash([]byte{1, 2, 3}),
							common.BytesToHash([]byte{1, 21, 99, 5}):  common.BytesToHash([]byte{1, 2, 3}),
							common.BytesToHash([]byte{1, 2, 99, 5}):   common.BytesToHash([]byte{1, 2, 3}),
							common.BytesToHash([]byte{1, 2, 4, 5}):    common.BytesToHash([]byte{1, 2, 3}),
						},
					},
				},
			},
			{
				BlockTimestamp: &t1,
				StateUpgradeAccounts: map[common.Address]StateUpgradeAccount{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000001000")): StateUpgradeAccount{},
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000050")): StateUpgradeAccount{},
				},
			},
		},
	}
	assertConversions(t, &config)
}
