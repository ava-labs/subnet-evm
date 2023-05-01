// (c) 2022 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/precompile/contracts/deployerallowlist"
	"github.com/ava-labs/subnet-evm/precompile/contracts/feemanager"
	"github.com/ava-labs/subnet-evm/precompile/contracts/nativeminter"
	"github.com/ava-labs/subnet-evm/precompile/contracts/rewardmanager"
	"github.com/ava-labs/subnet-evm/precompile/contracts/txallowlist"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestVerifyWithChainConfig(t *testing.T) {
	admins := []common.Address{{1}}
	baseConfig := *SubnetEVMDefaultChainConfig
	config := &baseConfig
	config.GenesisPrecompiles = Precompiles{
		txallowlist.ConfigKey: txallowlist.NewConfig(big.NewInt(2), nil, nil),
	}
	config.PrecompileUpgrades = []PrecompileUpgrade{
		{
			// disable TxAllowList at timestamp 4
			txallowlist.NewDisableConfig(big.NewInt(4)),
		},
		{
			// re-enable TxAllowList at timestamp 5
			txallowlist.NewConfig(big.NewInt(5), admins, nil),
		},
	}

	// check this config is valid
	err := config.Verify()
	require.NoError(t, err)

	// same precompile cannot be configured twice for the same timestamp
	badConfig := *config
	badConfig.PrecompileUpgrades = append(
		badConfig.PrecompileUpgrades,
		PrecompileUpgrade{
			Config: txallowlist.NewDisableConfig(big.NewInt(5)),
		},
	)
	err = badConfig.Verify()
	require.ErrorContains(t, err, "config block timestamp (5) <= previous timestamp (5) of same key")

	// cannot enable a precompile without disabling it first.
	badConfig = *config
	badConfig.PrecompileUpgrades = append(
		badConfig.PrecompileUpgrades,
		PrecompileUpgrade{
			Config: txallowlist.NewConfig(big.NewInt(5), admins, nil),
		},
	)
	err = badConfig.Verify()
	require.ErrorContains(t, err, "disable should be [true]")
}

func TestVerifyWithChainConfigAtNilTimestamp(t *testing.T) {
	admins := []common.Address{{0}}
	baseConfig := *SubnetEVMDefaultChainConfig
	config := &baseConfig
	config.PrecompileUpgrades = []PrecompileUpgrade{
		// this does NOT enable the precompile, so it should be upgradeable.
		{Config: txallowlist.NewConfig(nil, nil, nil)},
	}
	require.False(t, config.IsPrecompileEnabled(txallowlist.ContractAddress, common.Big0)) // check the precompile is not enabled.
	config.PrecompileUpgrades = []PrecompileUpgrade{
		{
			// enable TxAllowList at timestamp 5
			Config: txallowlist.NewConfig(big.NewInt(5), admins, nil),
		},
	}

	// check this config is valid
	require.NoError(t, config.Verify())
}

func TestVerifyPrecompileUpgrades(t *testing.T) {
	admins := []common.Address{{1}}
	tests := []struct {
		name          string
		upgrades      []PrecompileUpgrade
		expectedError string
	}{
		{
			name: "enable and disable tx allow list",
			upgrades: []PrecompileUpgrade{
				{
					Config: txallowlist.NewConfig(big.NewInt(1), admins, nil),
				},
				{
					Config: txallowlist.NewDisableConfig(big.NewInt(2)),
				},
			},
			expectedError: "",
		},
		{
			name: "invalid allow list config in tx allowlist",
			upgrades: []PrecompileUpgrade{
				{
					Config: txallowlist.NewConfig(big.NewInt(1), admins, nil),
				},
				{
					Config: txallowlist.NewDisableConfig(big.NewInt(2)),
				},
				{
					Config: txallowlist.NewConfig(big.NewInt(3), admins, admins),
				},
			},
			expectedError: "cannot set address",
		},
		{
			name: "invalid initial fee manager config",
			upgrades: []PrecompileUpgrade{
				{
					Config: feemanager.NewConfig(big.NewInt(3), admins, nil,
						func() *commontype.FeeConfig {
							feeConfig := DefaultFeeConfig
							feeConfig.GasLimit = big.NewInt(-1)
							return &feeConfig
						}()),
				},
			},
			expectedError: "gasLimit = -1 cannot be less than or equal to 0",
		},
		{
			name: "invalid initial fee manager config gas limit 0",
			upgrades: []PrecompileUpgrade{
				{
					Config: feemanager.NewConfig(big.NewInt(3), admins, nil,
						func() *commontype.FeeConfig {
							feeConfig := DefaultFeeConfig
							feeConfig.GasLimit = common.Big0
							return &feeConfig
						}()),
				},
			},
			expectedError: "gasLimit = 0 cannot be less than or equal to 0",
		},
		{
			name: "different upgrades are allowed to configure same timestamp for different precompiles",
			upgrades: []PrecompileUpgrade{
				{
					Config: txallowlist.NewConfig(big.NewInt(1), admins, nil),
				},
				{
					Config: feemanager.NewConfig(big.NewInt(1), admins, nil, nil),
				},
			},
			expectedError: "",
		},
		{
			name: "different upgrades must be monotonically increasing",
			upgrades: []PrecompileUpgrade{
				{
					Config: txallowlist.NewConfig(big.NewInt(2), admins, nil),
				},
				{
					Config: feemanager.NewConfig(big.NewInt(1), admins, nil, nil),
				},
			},
			expectedError: "config block timestamp (1) < previous timestamp (2)",
		},
		{
			name: "upgrades with same keys are not allowed to configure same timestamp for same precompiles",
			upgrades: []PrecompileUpgrade{
				{
					Config: txallowlist.NewConfig(big.NewInt(1), admins, nil),
				},
				{
					Config: txallowlist.NewDisableConfig(big.NewInt(1)),
				},
			},
			expectedError: "config block timestamp (1) <= previous timestamp (1) of same key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			baseConfig := *SubnetEVMDefaultChainConfig
			config := &baseConfig
			config.PrecompileUpgrades = tt.upgrades

			err := config.Verify()
			if tt.expectedError == "" {
				require.NoError(err)
			} else {
				require.ErrorContains(err, tt.expectedError)
			}
		})
	}
}

func TestVerifyPrecompiles(t *testing.T) {
	admins := []common.Address{{1}}
	tests := []struct {
		name          string
		precompiles   Precompiles
		expectedError string
	}{
		{
			name: "invalid allow list config in tx allowlist",
			precompiles: Precompiles{
				txallowlist.ConfigKey: txallowlist.NewConfig(big.NewInt(3), admins, admins),
			},
			expectedError: "cannot set address",
		},
		{
			name: "invalid initial fee manager config",
			precompiles: Precompiles{
				feemanager.ConfigKey: feemanager.NewConfig(big.NewInt(3), admins, nil,
					func() *commontype.FeeConfig {
						feeConfig := DefaultFeeConfig
						feeConfig.GasLimit = big.NewInt(-1)
						return &feeConfig
					}()),
			},
			expectedError: "gasLimit = -1 cannot be less than or equal to 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			baseConfig := *SubnetEVMDefaultChainConfig
			config := &baseConfig
			config.GenesisPrecompiles = tt.precompiles

			err := config.Verify()
			if tt.expectedError == "" {
				require.NoError(err)
			} else {
				require.ErrorContains(err, tt.expectedError)
			}
		})
	}
}

func TestVerifyRequiresSortedTimestamps(t *testing.T) {
	admins := []common.Address{{1}}
	baseConfig := *SubnetEVMDefaultChainConfig
	config := &baseConfig
	config.PrecompileUpgrades = []PrecompileUpgrade{
		{
			Config: txallowlist.NewConfig(big.NewInt(2), admins, nil),
		},
		{
			Config: deployerallowlist.NewConfig(big.NewInt(1), admins, nil),
		},
	}

	// block timestamps must be monotonically increasing, so this config is invalid
	err := config.Verify()
	require.ErrorContains(t, err, "config block timestamp (1) < previous timestamp (2)")
}

func TestGetPrecompileConfig(t *testing.T) {
	require := require.New(t)
	baseConfig := *SubnetEVMDefaultChainConfig
	config := &baseConfig
	config.GenesisPrecompiles = Precompiles{
		deployerallowlist.ConfigKey: deployerallowlist.NewConfig(big.NewInt(10), nil, nil),
	}

	deployerConfig := config.getActivePrecompileConfig(deployerallowlist.ContractAddress, big.NewInt(0))
	require.Nil(deployerConfig)

	deployerConfig = config.getActivePrecompileConfig(deployerallowlist.ContractAddress, big.NewInt(10))
	require.NotNil(deployerConfig)

	deployerConfig = config.getActivePrecompileConfig(deployerallowlist.ContractAddress, big.NewInt(11))
	require.NotNil(deployerConfig)

	txAllowListConfig := config.getActivePrecompileConfig(txallowlist.ContractAddress, big.NewInt(0))
	require.Nil(txAllowListConfig)
}

func TestPrecompileUpgradeUnmarshalJSON(t *testing.T) {
	require := require.New(t)

	upgradeBytes := []byte(`
			{
				"precompileUpgrades": [
					{
						"rewardManagerConfig": {
							"blockTimestamp": 1671542573,
							"adminAddresses": [
								"0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
							],
							"initialRewardConfig": {
								"allowFeeRecipients": true
							}
						}
					},
					{
						"contractNativeMinterConfig": {
							"blockTimestamp": 1671543172,
							"disable": false
						}
					}
				]
			}
	`)

	var upgradeConfig UpgradeConfig
	err := json.Unmarshal(upgradeBytes, &upgradeConfig)
	require.NoError(err)

	require.Len(upgradeConfig.PrecompileUpgrades, 2)

	rewardManagerConf := upgradeConfig.PrecompileUpgrades[0]
	require.Equal(rewardManagerConf.Key(), rewardmanager.ConfigKey)
	testRewardManagerConfig := rewardmanager.NewConfig(
		big.NewInt(1671542573),
		[]common.Address{common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")},
		nil,
		&rewardmanager.InitialRewardConfig{
			AllowFeeRecipients: true,
		})
	require.True(rewardManagerConf.Equal(testRewardManagerConfig))

	nativeMinterConfig := upgradeConfig.PrecompileUpgrades[1]
	require.Equal(nativeMinterConfig.Key(), nativeminter.ConfigKey)
	expectedNativeMinterConfig := nativeminter.NewConfig(big.NewInt(1671543172), nil, nil, nil)
	require.True(nativeMinterConfig.Equal(expectedNativeMinterConfig))

	// Marshal and unmarshal again and check that the result is the same
	upgradeBytes2, err := json.Marshal(upgradeConfig)
	require.NoError(err)
	var upgradeConfig2 UpgradeConfig
	err = json.Unmarshal(upgradeBytes2, &upgradeConfig2)
	require.NoError(err)
	require.Equal(upgradeConfig, upgradeConfig2)
}
