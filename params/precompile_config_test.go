// (c) 2022 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/precompile/deployerallowlist"
	"github.com/ava-labs/subnet-evm/precompile/feemanager"
	"github.com/ava-labs/subnet-evm/precompile/nativeminter"
	"github.com/ava-labs/subnet-evm/precompile/rewardmanager"
	"github.com/ava-labs/subnet-evm/precompile/txallowlist"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyWithChainConfig(t *testing.T) {
	admins := []common.Address{{1}}
	baseConfig := *SubnetEVMDefaultChainConfig
	config := &baseConfig
	config.Precompiles = ChainConfigPrecompiles{
		txallowlist.ConfigKey: txallowlist.NewTxAllowListConfig(big.NewInt(2), nil, nil),
	}
	config.PrecompileUpgrades = []PrecompileUpgrade{
		{
			// disable TxAllowList at timestamp 4
			txallowlist.NewDisableTxAllowListConfig(big.NewInt(4)),
		},
		{
			// re-enable TxAllowList at timestamp 5
			txallowlist.NewTxAllowListConfig(big.NewInt(5), admins, nil),
		},
	}

	// check this config is valid
	err := config.Verify()
	assert.NoError(t, err)

	// same precompile cannot be configured twice for the same timestamp
	badConfig := *config
	badConfig.PrecompileUpgrades = append(
		badConfig.PrecompileUpgrades,
		PrecompileUpgrade{
			txallowlist.NewDisableTxAllowListConfig(big.NewInt(5)),
		},
	)
	err = badConfig.Verify()
	assert.ErrorContains(t, err, "config timestamp (5) <= previous timestamp (5)")

	// cannot enable a precompile without disabling it first.
	badConfig = *config
	badConfig.PrecompileUpgrades = append(
		badConfig.PrecompileUpgrades,
		PrecompileUpgrade{
			txallowlist.NewTxAllowListConfig(big.NewInt(5), admins, nil),
		},
	)
	err = badConfig.Verify()
	assert.ErrorContains(t, err, "disable should be [true]")
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
					txallowlist.NewTxAllowListConfig(big.NewInt(1), admins, nil),
				},
				{
					txallowlist.NewDisableTxAllowListConfig(big.NewInt(2)),
				},
			},
			expectedError: "",
		},
		{
			name: "invalid allow list config in tx allowlist",
			upgrades: []PrecompileUpgrade{
				{
					txallowlist.NewTxAllowListConfig(big.NewInt(1), admins, nil),
				},
				{
					txallowlist.NewDisableTxAllowListConfig(big.NewInt(2)),
				},
				{
					txallowlist.NewTxAllowListConfig(big.NewInt(3), admins, admins),
				},
			},
			expectedError: "cannot set address",
		},
		{
			name: "invalid initial fee manager config",
			upgrades: []PrecompileUpgrade{
				{
					feemanager.NewFeeManagerConfig(big.NewInt(3), admins, nil,
						&commontype.FeeConfig{
							GasLimit: big.NewInt(-1),
						}),
				},
			},
			expectedError: "gasLimit = -1 cannot be less than or equal to 0",
		},
		{
			name: "invalid initial fee manager config gas limit 0",
			upgrades: []PrecompileUpgrade{
				{
					feemanager.NewFeeManagerConfig(big.NewInt(3), admins, nil,
						&commontype.FeeConfig{
							GasLimit: big.NewInt(0),
						}),
				},
			},
			expectedError: "gasLimit = 0 cannot be less than or equal to 0",
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
		upgrade       ChainConfigPrecompiles
		expectedError string
	}{
		{
			name: "invalid allow list config in tx allowlist",
			upgrade: ChainConfigPrecompiles{
				txallowlist.ConfigKey: txallowlist.NewTxAllowListConfig(big.NewInt(3), admins, admins),
			},
			expectedError: "cannot set address",
		},
		{
			name: "invalid initial fee manager config",
			upgrade: ChainConfigPrecompiles{
				feemanager.ConfigKey: feemanager.NewFeeManagerConfig(big.NewInt(3), admins, nil,
					&commontype.FeeConfig{
						GasLimit: big.NewInt(-1),
					}),
			},
			expectedError: "gasLimit = -1 cannot be less than or equal to 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			baseConfig := *SubnetEVMDefaultChainConfig
			config := &baseConfig
			config.Precompiles = tt.upgrade

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
			txallowlist.NewTxAllowListConfig(big.NewInt(2), admins, nil),
		},
		{
			deployerallowlist.NewContractDeployerAllowListConfig(big.NewInt(1), admins, nil),
		},
	}

	// block timestamps must be monotonically increasing, so this config is invalid
	err := config.Verify()
	assert.ErrorContains(t, err, "config timestamp (1) < previous timestamp (2)")
}

func TestGetPrecompileConfig(t *testing.T) {
	assert := assert.New(t)
	baseConfig := *SubnetEVMDefaultChainConfig
	config := &baseConfig
	config.Precompiles = ChainConfigPrecompiles{
		deployerallowlist.ConfigKey: deployerallowlist.NewContractDeployerAllowListConfig(big.NewInt(10), nil, nil),
	}

	deployerConfig := config.GetActivePrecompileConfig(deployerallowlist.ContractAddress, big.NewInt(0))
	assert.Nil(deployerConfig)

	deployerConfig = config.GetActivePrecompileConfig(deployerallowlist.ContractAddress, big.NewInt(10))
	assert.NotNil(deployerConfig)

	deployerConfig = config.GetActivePrecompileConfig(deployerallowlist.ContractAddress, big.NewInt(11))
	assert.NotNil(deployerConfig)

	txAllowListConfig := config.GetActivePrecompileConfig(txallowlist.ContractAddress, big.NewInt(0))
	assert.Nil(txAllowListConfig)
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
	testRewardManagerConfig := rewardmanager.NewRewardManagerConfig(
		big.NewInt(1671542573),
		[]common.Address{common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")},
		nil,
		&rewardmanager.InitialRewardConfig{
			AllowFeeRecipients: true,
		})
	require.True(rewardManagerConf.Equal(testRewardManagerConfig))

	contractNativeMinterConf := upgradeConfig.PrecompileUpgrades[1]
	require.Equal(contractNativeMinterConf.Key(), nativeminter.ConfigKey)
	testContractNativeMinterConfig := nativeminter.NewContractNativeMinterConfig(big.NewInt(1671543172), nil, nil, nil)
	require.True(contractNativeMinterConf.Equal(testContractNativeMinterConfig))

	// Marshal and unmarshal again and check that the result is the same
	upgradeBytes2, err := json.Marshal(upgradeConfig)
	require.NoError(err)
	var upgradeConfig2 UpgradeConfig
	err = json.Unmarshal(upgradeBytes2, &upgradeConfig2)
	require.NoError(err)
	require.Equal(upgradeConfig, upgradeConfig2)
}
