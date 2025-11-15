// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package extras

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/common/math"
	"github.com/stretchr/testify/require"

	"github.com/ava-labs/subnet-evm/stateupgrade"
	"github.com/ava-labs/subnet-evm/utils"
)

func TestCheckCompatibleStateUpgrades(t *testing.T) {
	chainConfig := *TestChainConfig
	stateUpgrade := map[common.Address]stateupgrade.StateUpgradeAccount{
		{1}: {BalanceChange: (*math.HexOrDecimal256)(common.Big1)},
	}
	differentStateUpgrade := map[common.Address]stateupgrade.StateUpgradeAccount{
		{2}: {BalanceChange: (*math.HexOrDecimal256)(common.Big1)},
	}

	tests := map[string]upgradeCompatibilityTest{
		"reschedule upgrade before it happens": {
			startTimestamps: []uint64{5, 6},
			configs: []*UpgradeConfig{
				{
					StateUpgrades: []stateupgrade.StateUpgrade{
						{BlockTimestamp: utils.NewUint64(6), StateUpgradeAccounts: stateUpgrade},
					},
				},
				{
					StateUpgrades: []stateupgrade.StateUpgrade{
						{BlockTimestamp: utils.NewUint64(6), StateUpgradeAccounts: stateUpgrade},
					},
				},
			},
		},
		"modify upgrade after it happens not allowed": {
			expectedErrorString: "mismatching StateUpgrade",
			startTimestamps:     []uint64{5, 8},
			configs: []*UpgradeConfig{
				{
					StateUpgrades: []stateupgrade.StateUpgrade{
						{BlockTimestamp: utils.NewUint64(6), StateUpgradeAccounts: stateUpgrade},
						{BlockTimestamp: utils.NewUint64(7), StateUpgradeAccounts: stateUpgrade},
					},
				},
				{
					StateUpgrades: []stateupgrade.StateUpgrade{
						{BlockTimestamp: utils.NewUint64(6), StateUpgradeAccounts: stateUpgrade},
						{BlockTimestamp: utils.NewUint64(7), StateUpgradeAccounts: differentStateUpgrade},
					},
				},
			},
		},
		"cancel upgrade before it happens": {
			startTimestamps: []uint64{5, 6},
			configs: []*UpgradeConfig{
				{
					StateUpgrades: []stateupgrade.StateUpgrade{
						{BlockTimestamp: utils.NewUint64(6), StateUpgradeAccounts: stateUpgrade},
						{BlockTimestamp: utils.NewUint64(7), StateUpgradeAccounts: stateUpgrade},
					},
				},
				{
					StateUpgrades: []stateupgrade.StateUpgrade{
						{BlockTimestamp: utils.NewUint64(6), StateUpgradeAccounts: stateUpgrade},
					},
				},
			},
		},
		"retroactively enabling upgrades is not allowed": {
			expectedErrorString: "cannot retroactively enable StateUpgrade[0] in database (have timestamp nil, want timestamp 5, rewindto timestamp 4)",
			startTimestamps:     []uint64{6},
			configs: []*UpgradeConfig{
				{
					StateUpgrades: []stateupgrade.StateUpgrade{
						{BlockTimestamp: utils.NewUint64(5), StateUpgradeAccounts: stateUpgrade},
					},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.run(t, chainConfig)
		})
	}
}

func TestUnmarshalStateUpgradeJSON(t *testing.T) {
	jsonBytes := []byte(
		`{
			"stateUpgrades": [
				{
					"blockTimestamp": 1677608400,
					"accounts": {
						"0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC": {
							"balanceChange": "100"
						}
					}
				}
			]
		}`,
	)

	upgradeConfig := UpgradeConfig{
		StateUpgrades: []stateupgrade.StateUpgrade{
			{
				BlockTimestamp: utils.NewUint64(1677608400),
				StateUpgradeAccounts: map[common.Address]stateupgrade.StateUpgradeAccount{
					common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"): {
						BalanceChange: (*math.HexOrDecimal256)(big.NewInt(100)),
					},
				},
			},
		},
	}
	var unmarshaledConfig UpgradeConfig
	require.NoError(t, json.Unmarshal(jsonBytes, &unmarshaledConfig))
	require.Equal(t, upgradeConfig, unmarshaledConfig)
}
