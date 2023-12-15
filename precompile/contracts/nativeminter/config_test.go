// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nativeminter

import (
	"testing"

	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"go.uber.org/mock/gomock"
)

func TestVerify(t *testing.T) {
	admins := []common.Address{allowlist.TestAdminAddr}
	enableds := []common.Address{allowlist.TestEnabledAddr}
	managers := []common.Address{allowlist.TestManagerAddr}
	tests := map[string]testutils.ConfigVerifyTest{
		"invalid allow list config in native minter allowlist": {
			Config:        NewConfig(utils.NewUint64(3), admins, admins, nil, nil),
			ExpectedError: "cannot set address",
		},
		"duplicate admins in config in native minter allowlist": {
			Config:        NewConfig(utils.NewUint64(3), append(admins, admins[0]), enableds, managers, nil),
			ExpectedError: "duplicate address",
		},
		"duplicate enableds in config in native minter allowlist": {
			Config:        NewConfig(utils.NewUint64(3), admins, append(enableds, enableds[0]), managers, nil),
			ExpectedError: "duplicate address",
		},
		"nil amount in native minter config": {
			Config: NewConfig(utils.NewUint64(3), admins, nil, nil,
				map[common.Address]*math.HexOrDecimal256{
					common.HexToAddress("0x01"): math.NewHexOrDecimal256(123),
					common.HexToAddress("0x02"): nil,
				}),
			ExpectedError: "initial mint cannot contain nil",
		},
		"negative amount in native minter config": {
			Config: NewConfig(utils.NewUint64(3), admins, nil, nil,
				map[common.Address]*math.HexOrDecimal256{
					common.HexToAddress("0x01"): math.NewHexOrDecimal256(123),
					common.HexToAddress("0x02"): math.NewHexOrDecimal256(-1),
				}),
			ExpectedError: "initial mint cannot contain invalid amount",
		},
	}
	allowlist.VerifyPrecompileWithAllowListTests(t, Module, tests)
}

func TestSerialize(t *testing.T) {
	var t0 uint64 = 0
	var t1 uint64 = 1
	config := params.UpgradeConfig{
		PrecompileUpgrades: []params.PrecompileUpgrade{
			{
				Config: NewConfig(&t0, nil, nil, nil, nil), // enable at genesis
			},
			{
				Config: NewDisableConfig(&t1), // disable at timestamp 1
			},
		},
	}
	params.AssertConfigHashesAndSerialization(t, &config)
}
func TestSerializeWithAddresses(t *testing.T) {
	var t0 uint64 = 1
	var t1 uint64 = 11
	config := params.UpgradeConfig{
		PrecompileUpgrades: []params.PrecompileUpgrade{
			{
				Config: NewConfig(&t0, []common.Address{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000020")),
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000030")),
				}, []common.Address{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000040")),
				}, []common.Address{
					common.BytesToAddress(common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000050")),
				}, nil), // enable at genesis
			},
			{
				Config: NewDisableConfig(&t1), // disable at timestamp 1
			},
		},
	}
	params.AssertConfigHashesAndSerialization(t, &config)
}

func TestSerializeWithAddressAndMint(t *testing.T) {
	var t0 uint64 = 2
	var t1 uint64 = 1001
	config := params.UpgradeConfig{
		PrecompileUpgrades: []params.PrecompileUpgrade{
			{
				Config: NewConfig(&t0, []common.Address{
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
				Config: NewDisableConfig(&t1), // disable at timestamp 1
			},
		},
	}
	params.AssertConfigHashesAndSerialization(t, &config)
}

func TestEqual(t *testing.T) {
	admins := []common.Address{allowlist.TestAdminAddr}
	enableds := []common.Address{allowlist.TestEnabledAddr}
	managers := []common.Address{allowlist.TestManagerAddr}
	tests := map[string]testutils.ConfigEqualTest{
		"non-nil config and nil other": {
			Config:   NewConfig(utils.NewUint64(3), admins, enableds, managers, nil),
			Other:    nil,
			Expected: false,
		},
		"different type": {
			Config:   NewConfig(utils.NewUint64(3), admins, enableds, managers, nil),
			Other:    precompileconfig.NewMockConfig(gomock.NewController(t)),
			Expected: false,
		},
		"different timestamps": {
			Config:   NewConfig(utils.NewUint64(3), admins, nil, nil, nil),
			Other:    NewConfig(utils.NewUint64(4), admins, nil, nil, nil),
			Expected: false,
		},
		"different initial mint amounts": {
			Config: NewConfig(utils.NewUint64(3), admins, nil, nil,
				map[common.Address]*math.HexOrDecimal256{
					common.HexToAddress("0x01"): math.NewHexOrDecimal256(1),
				}),
			Other: NewConfig(utils.NewUint64(3), admins, nil, nil,
				map[common.Address]*math.HexOrDecimal256{
					common.HexToAddress("0x01"): math.NewHexOrDecimal256(2),
				}),
			Expected: false,
		},
		"different initial mint addresses": {
			Config: NewConfig(utils.NewUint64(3), admins, nil, nil,
				map[common.Address]*math.HexOrDecimal256{
					common.HexToAddress("0x01"): math.NewHexOrDecimal256(1),
				}),
			Other: NewConfig(utils.NewUint64(3), admins, nil, nil,
				map[common.Address]*math.HexOrDecimal256{
					common.HexToAddress("0x02"): math.NewHexOrDecimal256(1),
				}),
			Expected: false,
		},
		"same config": {
			Config: NewConfig(utils.NewUint64(3), admins, nil, nil,
				map[common.Address]*math.HexOrDecimal256{
					common.HexToAddress("0x01"): math.NewHexOrDecimal256(1),
				}),
			Other: NewConfig(utils.NewUint64(3), admins, nil, nil,
				map[common.Address]*math.HexOrDecimal256{
					common.HexToAddress("0x01"): math.NewHexOrDecimal256(1),
				}),
			Expected: true,
		},
	}
	allowlist.EqualPrecompileWithAllowListTests(t, Module, tests)
}
