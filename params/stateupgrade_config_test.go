// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestVerifyStateUpgrades(t *testing.T) {
	modifiedAccounts := map[common.Address]StateUpgradeAccount{
		{1}: {
			BalanceChange: common.Big1,
		},
	}
	tests := []struct {
		name          string
		upgrades      []StateUpgrade
		expectedError string
	}{
		{
			name: "valid upgrade",
			upgrades: []StateUpgrade{
				{BlockTimestamp: common.Big1, StateUpgradeAccounts: modifiedAccounts},
				{BlockTimestamp: common.Big2, StateUpgradeAccounts: modifiedAccounts},
			},
		},
		{
			name: "upgrade block timestamp is not strictly increasing",
			upgrades: []StateUpgrade{
				{BlockTimestamp: common.Big1, StateUpgradeAccounts: modifiedAccounts},
				{BlockTimestamp: common.Big1, StateUpgradeAccounts: modifiedAccounts},
			},
			expectedError: "config block timestamp (1) <= previous timestamp (1)",
		},
		{
			name: "upgrade block timestamp is decreases",
			upgrades: []StateUpgrade{
				{BlockTimestamp: common.Big2, StateUpgradeAccounts: modifiedAccounts},
				{BlockTimestamp: common.Big1, StateUpgradeAccounts: modifiedAccounts},
			},
			expectedError: "config block timestamp (1) <= previous timestamp (2)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			baseConfig := *SubnetEVMDefaultChainConfig
			config := &baseConfig
			config.StateUpgrades = tt.upgrades

			err := config.Verify()
			if tt.expectedError == "" {
				require.NoError(err)
			} else {
				require.ErrorContains(err, tt.expectedError)
			}
		})
	}
}

func TestCheckCompatibleStateUpgradeConfigs(t *testing.T) {
	chainConfig := *TestChainConfig
	stateUpgrade := StateUpgradeAccounts{
		{1}: {BalanceChange: common.Big1},
	}
	differentStateUpgrade := StateUpgradeAccounts{
		{2}: {BalanceChange: common.Big1},
	}

	tests := map[string]upgradeCompatibilityTest{
		"reschedule upgrade before it happens": {
			startTimestamps: []*big.Int{big.NewInt(5), big.NewInt(6)},
			configs: []*UpgradeConfig{
				{
					StateUpgrades: []StateUpgrade{
						{BlockTimestamp: big.NewInt(6), StateUpgradeAccounts: stateUpgrade},
					},
				},
				{
					StateUpgrades: []StateUpgrade{
						{BlockTimestamp: big.NewInt(6), StateUpgradeAccounts: stateUpgrade},
					},
				},
			},
		},
		"modify upgrade after it happens not allowed": {
			expectedErrorString: "mismatching StateUpgrade",
			startTimestamps:     []*big.Int{big.NewInt(5), big.NewInt(8)},
			configs: []*UpgradeConfig{
				{
					StateUpgrades: []StateUpgrade{
						{BlockTimestamp: big.NewInt(6), StateUpgradeAccounts: stateUpgrade},
						{BlockTimestamp: big.NewInt(7), StateUpgradeAccounts: stateUpgrade},
					},
				},
				{
					StateUpgrades: []StateUpgrade{
						{BlockTimestamp: big.NewInt(6), StateUpgradeAccounts: stateUpgrade},
						{BlockTimestamp: big.NewInt(7), StateUpgradeAccounts: differentStateUpgrade},
					},
				},
			},
		},
		"cancel upgrade before it happens": {
			startTimestamps: []*big.Int{big.NewInt(5), big.NewInt(6)},
			configs: []*UpgradeConfig{
				{
					StateUpgrades: []StateUpgrade{
						{BlockTimestamp: big.NewInt(6), StateUpgradeAccounts: stateUpgrade},
						{BlockTimestamp: big.NewInt(7), StateUpgradeAccounts: stateUpgrade},
					},
				},
				{
					StateUpgrades: []StateUpgrade{
						{BlockTimestamp: big.NewInt(6), StateUpgradeAccounts: stateUpgrade},
					},
				},
			},
		},
		"retroactively enabling upgrades is not allowed": {
			expectedErrorString: "cannot retroactively enable StateUpgrade",
			startTimestamps:     []*big.Int{big.NewInt(6)},
			configs: []*UpgradeConfig{
				{
					StateUpgrades: []StateUpgrade{
						{BlockTimestamp: big.NewInt(5), StateUpgradeAccounts: stateUpgrade},
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

func TestUnmarshalJSON(t *testing.T) {
	jsonBytes := []byte(
		`{
			"stateUpgrades": [
				{
					"blockTimestamp": 1677608400,
					"accounts": {
						"8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC": {
							"balanceChange": "0x52B7D2DCC80CD2E4000000"
						}
					}
				}
			]
		}`,
	)

	var config UpgradeConfig
	err := json.Unmarshal(jsonBytes, &config)
	require.NoError(t, err)
}
