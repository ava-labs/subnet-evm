// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package stateupgrade

import (
	"testing"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/libevm/common/math"
	"github.com/stretchr/testify/require"

	"github.com/ava-labs/subnet-evm/utils"
)

func TestVerifyStateUpgrades(t *testing.T) {
	modifiedAccounts := map[common.Address]StateUpgradeAccount{
		{1}: {
			BalanceChange: (*math.HexOrDecimal256)(common.Big1),
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
				{BlockTimestamp: utils.NewUint64(1), StateUpgradeAccounts: modifiedAccounts},
				{BlockTimestamp: utils.NewUint64(2), StateUpgradeAccounts: modifiedAccounts},
			},
		},
		{
			name: "upgrade block timestamp is not strictly increasing",
			upgrades: []StateUpgrade{
				{BlockTimestamp: utils.NewUint64(1), StateUpgradeAccounts: modifiedAccounts},
				{BlockTimestamp: utils.NewUint64(1), StateUpgradeAccounts: modifiedAccounts},
			},
			expectedError: "config block timestamp (1) <= previous timestamp (1)",
		},
		{
			name: "upgrade block timestamp decreases",
			upgrades: []StateUpgrade{
				{BlockTimestamp: utils.NewUint64(2), StateUpgradeAccounts: modifiedAccounts},
				{BlockTimestamp: utils.NewUint64(1), StateUpgradeAccounts: modifiedAccounts},
			},
			expectedError: "config block timestamp (1) <= previous timestamp (2)",
		},
		{
			name: "upgrade block timestamp is zero",
			upgrades: []StateUpgrade{
				{BlockTimestamp: utils.NewUint64(0), StateUpgradeAccounts: modifiedAccounts},
			},
			expectedError: "config block timestamp (0) must be greater than 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)

			err := VerifyStateUpgrades(tt.upgrades)
			if tt.expectedError == "" {
				require.NoError(err)
			} else {
				require.ErrorContains(err, tt.expectedError)
			}
		})
	}
}
