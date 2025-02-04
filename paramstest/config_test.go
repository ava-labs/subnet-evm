package paramstest

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contracts/nativeminter"
	"github.com/ava-labs/subnet-evm/precompile/contracts/rewardmanager"
	"github.com/ava-labs/subnet-evm/precompile/contracts/txallowlist"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestConfigUnmarshalJSON(t *testing.T) {
	require := require.New(t)

	testRewardManagerConfig := rewardmanager.NewConfig(
		utils.NewUint64(1671542573),
		[]common.Address{common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")},
		nil,
		nil,
		&rewardmanager.InitialRewardConfig{
			AllowFeeRecipients: true,
		})

	testContractNativeMinterConfig := nativeminter.NewConfig(
		utils.NewUint64(0),
		[]common.Address{common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")},
		nil,
		nil,
		nil,
	)

	config := []byte(`
	{
		"chainId": 43214,
		"allowFeeRecipients": true,
		"rewardManagerConfig": {
			"blockTimestamp": 1671542573,
			"adminAddresses": [
				"0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
			],
			"initialRewardConfig": {
				"allowFeeRecipients": true
			}
		},
		"contractNativeMinterConfig": {
			"blockTimestamp": 0,
			"adminAddresses": [
				"0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
			]
		}
	}
	`)
	c := params.ChainConfig{}
	err := json.Unmarshal(config, &c)
	require.NoError(err)

	require.Equal(c.ChainID, big.NewInt(43214))
	require.Equal(c.AllowFeeRecipients, true)

	rewardManagerConfig, ok := c.GenesisPrecompiles[rewardmanager.ConfigKey]
	require.True(ok)
	require.Equal(rewardManagerConfig.Key(), rewardmanager.ConfigKey)
	require.True(rewardManagerConfig.Equal(testRewardManagerConfig))

	nativeMinterConfig := c.GenesisPrecompiles[nativeminter.ConfigKey]
	require.Equal(nativeMinterConfig.Key(), nativeminter.ConfigKey)
	require.True(nativeMinterConfig.Equal(testContractNativeMinterConfig))

	// Marshal and unmarshal again and check that the result is the same
	marshaled, err := json.Marshal(c)
	require.NoError(err)
	c2 := params.ChainConfig{}
	err = json.Unmarshal(marshaled, &c2)
	require.NoError(err)
	require.Equal(c, c2)
}

func TestActivePrecompiles(t *testing.T) {
	config := params.ChainConfig{
		UpgradeConfig: params.UpgradeConfig{
			PrecompileUpgrades: []params.PrecompileUpgrade{
				{
					Config: nativeminter.NewConfig(utils.NewUint64(0), nil, nil, nil, nil), // enable at genesis
				},
				{
					Config: nativeminter.NewDisableConfig(utils.NewUint64(1)), // disable at timestamp 1
				},
			},
		},
	}

	rules0 := config.Rules(common.Big0, 0)
	require.True(t, rules0.IsPrecompileEnabled(nativeminter.Module.Address))

	rules1 := config.Rules(common.Big0, 1)
	require.False(t, rules1.IsPrecompileEnabled(nativeminter.Module.Address))
}

func TestChainConfigMarshalWithUpgrades(t *testing.T) {
	config := params.ChainConfigWithUpgradesJSON{
		ChainConfig: params.ChainConfig{
			ChainID:             big.NewInt(1),
			FeeConfig:           params.DefaultFeeConfig,
			AllowFeeRecipients:  false,
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
			MuirGlacierBlock:    big.NewInt(0),
			NetworkUpgrades: params.NetworkUpgrades{
				SubnetEVMTimestamp: utils.NewUint64(0),
				DurangoTimestamp:   utils.NewUint64(0),
			},
			GenesisPrecompiles: params.Precompiles{},
		},
		UpgradeConfig: params.UpgradeConfig{
			PrecompileUpgrades: []params.PrecompileUpgrade{
				{
					Config: txallowlist.NewConfig(utils.NewUint64(100), nil, nil, nil),
				},
			},
		},
	}
	result, err := json.Marshal(&config)
	require.NoError(t, err)
	expectedJSON := `{
		"chainId": 1,
		"feeConfig": {
			"gasLimit": 8000000,
			"targetBlockRate": 2,
			"minBaseFee": 25000000000,
			"targetGas": 15000000,
			"baseFeeChangeDenominator": 36,
			"minBlockGasCost": 0,
			"maxBlockGasCost": 1000000,
			"blockGasCostStep": 200000
		},
		"homesteadBlock": 0,
		"eip150Block": 0,
		"eip155Block": 0,
		"eip158Block": 0,
		"byzantiumBlock": 0,
		"constantinopleBlock": 0,
		"petersburgBlock": 0,
		"istanbulBlock": 0,
		"muirGlacierBlock": 0,
		"subnetEVMTimestamp": 0,
		"durangoTimestamp": 0,
		"upgrades": {
			"precompileUpgrades": [
				{
					"txAllowListConfig": {
						"blockTimestamp": 100
					}
				}
			]
		}
	}`
	require.JSONEq(t, expectedJSON, string(result))

	var unmarshalled params.ChainConfigWithUpgradesJSON
	err = json.Unmarshal(result, &unmarshalled)
	require.NoError(t, err)
	require.Equal(t, config, unmarshalled)
}
