// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
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
				{BlockTimestamp: common.Big1, Accounts: modifiedAccounts},
				{BlockTimestamp: common.Big2, Accounts: modifiedAccounts},
			},
		},
		{
			name: "upgrade block timestamp is not strictly increasing",
			upgrades: []StateUpgrade{
				{BlockTimestamp: common.Big1, Accounts: modifiedAccounts},
				{BlockTimestamp: common.Big1, Accounts: modifiedAccounts},
			},
			expectedError: "config block timestamp (1) <= previous timestamp (1)",
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
	stateUpgrade := map[common.Address]StateUpgradeAccount{
		{1}: {BalanceChange: common.Big1},
	}
	differentStateUpgrade := map[common.Address]StateUpgradeAccount{
		{2}: {BalanceChange: common.Big1},
	}

	tests := map[string]upgradeCompatibilityTest{
		"reschedule upgrade before it happens": {
			startTimestamps: []*big.Int{big.NewInt(5), big.NewInt(6)},
			configs: []*UpgradeConfig{
				{
					StateUpgrades: []StateUpgrade{
						{BlockTimestamp: big.NewInt(6), Accounts: stateUpgrade},
					},
				},
				{
					StateUpgrades: []StateUpgrade{
						{BlockTimestamp: big.NewInt(6), Accounts: stateUpgrade},
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
						{BlockTimestamp: big.NewInt(6), Accounts: stateUpgrade},
						{BlockTimestamp: big.NewInt(7), Accounts: stateUpgrade},
					},
				},
				{
					StateUpgrades: []StateUpgrade{
						{BlockTimestamp: big.NewInt(6), Accounts: stateUpgrade},
						{BlockTimestamp: big.NewInt(7), Accounts: differentStateUpgrade},
					},
				},
			},
		},
		"cancel upgrade before it happens": {
			startTimestamps: []*big.Int{big.NewInt(5), big.NewInt(6)},
			configs: []*UpgradeConfig{
				{
					StateUpgrades: []StateUpgrade{
						{BlockTimestamp: big.NewInt(6), Accounts: stateUpgrade},
						{BlockTimestamp: big.NewInt(7), Accounts: stateUpgrade},
					},
				},
				{
					StateUpgrades: []StateUpgrade{
						{BlockTimestamp: big.NewInt(6), Accounts: stateUpgrade},
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
						{BlockTimestamp: big.NewInt(5), Accounts: stateUpgrade},
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
