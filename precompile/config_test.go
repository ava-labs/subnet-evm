// (c) 2022 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestVerifyPrecompileUpgrades(t *testing.T) {
	admins := []common.Address{{1}}
	tests := []struct {
		name          string
		config        StatefulPrecompileConfig
		expectedError string
	}{
		{
			name: "invalid allow list config in tx allowlist",
			config: &TxAllowListConfig{
				AllowListConfig:   AllowListConfig{AllowListAdmins: admins, EnabledAddresses: admins},
				UpgradeableConfig: UpgradeableConfig{BlockTimestamp: big.NewInt(3)},
			},
			expectedError: "cannot set address",
		},
		{
			name: "invalid allow list config in deployer allowlist",
			config: &ContractDeployerAllowListConfig{
				AllowListConfig:   AllowListConfig{AllowListAdmins: admins, EnabledAddresses: admins},
				UpgradeableConfig: UpgradeableConfig{BlockTimestamp: big.NewInt(3)},
			},
			expectedError: "cannot set address",
		},
		{
			name: "invalid allow list config in native minter allowlist",
			config: &ContractNativeMinterConfig{
				AllowListConfig:   AllowListConfig{AllowListAdmins: admins, EnabledAddresses: admins},
				UpgradeableConfig: UpgradeableConfig{BlockTimestamp: big.NewInt(3)},
			},
			expectedError: "cannot set address",
		},
		{
			name: "invalid allow list config in fee manager allowlist",
			config: &FeeConfigManagerConfig{
				AllowListConfig:   AllowListConfig{AllowListAdmins: admins, EnabledAddresses: admins},
				UpgradeableConfig: UpgradeableConfig{BlockTimestamp: big.NewInt(3)},
			},
			expectedError: "cannot set address",
		},
		{
			name: "invalid initial fee manager config",
			config: &FeeConfigManagerConfig{
				AllowListConfig:   AllowListConfig{AllowListAdmins: admins},
				UpgradeableConfig: UpgradeableConfig{BlockTimestamp: big.NewInt(3)},
				InitialFeeConfig: &commontype.FeeConfig{
					GasLimit: big.NewInt(0),
				},
			},
			expectedError: "gasLimit = 0 cannot be less than or equal to 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)

			err := tt.config.Verify()
			if tt.expectedError == "" {
				require.NoError(err)
			} else {
				require.ErrorContains(err, tt.expectedError)
			}
		})
	}
}
