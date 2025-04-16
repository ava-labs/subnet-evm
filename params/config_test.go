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
	"encoding/json"
	"math"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/ava-labs/libevm/common"
	ethparams "github.com/ava-labs/libevm/params"
	"github.com/ava-labs/subnet-evm/params/extras"
	"github.com/ava-labs/subnet-evm/precompile/contracts/nativeminter"
	"github.com/ava-labs/subnet-evm/precompile/contracts/rewardmanager"
	"github.com/ava-labs/subnet-evm/precompile/contracts/txallowlist"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckCompatible(t *testing.T) {
	type test struct {
		stored, new   *ChainConfig
		headBlock     uint64
		headTimestamp uint64
		wantErr       *ethparams.ConfigCompatError
	}
	tests := []test{
		{stored: TestChainConfig, new: TestChainConfig, headBlock: 0, headTimestamp: 0, wantErr: nil},
		{stored: TestChainConfig, new: TestChainConfig, headBlock: 0, headTimestamp: uint64(time.Now().Unix()), wantErr: nil},
		{stored: TestChainConfig, new: TestChainConfig, headBlock: 100, wantErr: nil},
		{
			stored:        &ChainConfig{EIP150Block: big.NewInt(10)},
			new:           &ChainConfig{EIP150Block: big.NewInt(20)},
			headBlock:     9,
			headTimestamp: 90,
			wantErr:       nil,
		},
		{
			stored:        TestChainConfig,
			new:           &ChainConfig{HomesteadBlock: nil},
			headBlock:     3,
			headTimestamp: 30,
			wantErr: &ethparams.ConfigCompatError{
				What:          "Homestead fork block",
				StoredBlock:   big.NewInt(0),
				NewBlock:      nil,
				RewindToBlock: 0,
			},
		},
		{
			stored:        TestChainConfig,
			new:           &ChainConfig{HomesteadBlock: big.NewInt(1)},
			headBlock:     3,
			headTimestamp: 30,
			wantErr: &ethparams.ConfigCompatError{
				What:          "Homestead fork block",
				StoredBlock:   big.NewInt(0),
				NewBlock:      big.NewInt(1),
				RewindToBlock: 0,
			},
		},
		{
			stored:        &ChainConfig{HomesteadBlock: big.NewInt(30), EIP150Block: big.NewInt(10)},
			new:           &ChainConfig{HomesteadBlock: big.NewInt(25), EIP150Block: big.NewInt(20)},
			headBlock:     25,
			headTimestamp: 250,
			wantErr: &ethparams.ConfigCompatError{
				What:          "EIP150 fork block",
				StoredBlock:   big.NewInt(10),
				NewBlock:      big.NewInt(20),
				RewindToBlock: 9,
			},
		},
		{
			stored:        &ChainConfig{ConstantinopleBlock: big.NewInt(30)},
			new:           &ChainConfig{ConstantinopleBlock: big.NewInt(30), PetersburgBlock: big.NewInt(30)},
			headBlock:     40,
			headTimestamp: 400,
			wantErr:       nil,
		},
		{
			stored:        &ChainConfig{ConstantinopleBlock: big.NewInt(30)},
			new:           &ChainConfig{ConstantinopleBlock: big.NewInt(30), PetersburgBlock: big.NewInt(31)},
			headBlock:     40,
			headTimestamp: 400,
			wantErr: &ethparams.ConfigCompatError{
				What:          "Petersburg fork block",
				StoredBlock:   nil,
				NewBlock:      big.NewInt(31),
				RewindToBlock: 30,
			},
		},
		{
			stored:        TestChainConfig,
			new:           TestPreSubnetEVMChainConfig,
			headBlock:     0,
			headTimestamp: 0,
			wantErr: &ethparams.ConfigCompatError{
				What:         "SubnetEVM fork block timestamp",
				StoredTime:   utils.NewUint64(0),
				NewTime:      GetExtra(TestPreSubnetEVMChainConfig).NetworkUpgrades.SubnetEVMTimestamp,
				RewindToTime: 0,
			},
		},
		{
			stored:        TestChainConfig,
			new:           TestPreSubnetEVMChainConfig,
			headBlock:     10,
			headTimestamp: 100,
			wantErr: &ethparams.ConfigCompatError{
				What:         "SubnetEVM fork block timestamp",
				StoredTime:   utils.NewUint64(0),
				NewTime:      GetExtra(TestPreSubnetEVMChainConfig).NetworkUpgrades.SubnetEVMTimestamp,
				RewindToTime: 0,
			},
		},
	}

	for _, test := range tests {
		err := test.stored.CheckCompatible(test.new, test.headBlock, test.headTimestamp)
		if !reflect.DeepEqual(err, test.wantErr) {
			t.Errorf("error mismatch:\nstored: %v\nnew: %v\nblockHeight: %v\nerr: %v\nwant: %v", test.stored, test.new, test.headBlock, err, test.wantErr)
		}
	}
}

func TestConfigRules(t *testing.T) {
	c := WithExtra(
		&ChainConfig{},
		&extras.ChainConfig{
			NetworkUpgrades: extras.NetworkUpgrades{
				SubnetEVMTimestamp: utils.NewUint64(500),
			},
		},
	)

	var stamp uint64
	if r := c.Rules(big.NewInt(0), IsMergeTODO, stamp); GetRulesExtra(r).IsSubnetEVM {
		t.Errorf("expected %v to not be subnet-evm", stamp)
	}
	stamp = 500
	if r := c.Rules(big.NewInt(0), IsMergeTODO, stamp); !GetRulesExtra(r).IsSubnetEVM {
		t.Errorf("expected %v to be subnet-evm", stamp)
	}
	stamp = math.MaxInt64
	if r := c.Rules(big.NewInt(0), IsMergeTODO, stamp); !GetRulesExtra(r).IsSubnetEVM {
		t.Errorf("expected %v to be subnet-evm", stamp)
	}
}

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
	c := ChainConfig{}
	err := json.Unmarshal(config, &c)
	require.NoError(err)

	require.Equal(c.ChainID, big.NewInt(43214))
	require.Equal(GetExtra(&c).AllowFeeRecipients, true)

	rewardManagerConfig, ok := GetExtra(&c).GenesisPrecompiles[rewardmanager.ConfigKey]
	require.True(ok)
	require.Equal(rewardManagerConfig.Key(), rewardmanager.ConfigKey)
	require.True(rewardManagerConfig.Equal(testRewardManagerConfig))

	nativeMinterConfig := GetExtra(&c).GenesisPrecompiles[nativeminter.ConfigKey]
	require.Equal(nativeMinterConfig.Key(), nativeminter.ConfigKey)
	require.True(nativeMinterConfig.Equal(testContractNativeMinterConfig))

	// Marshal and unmarshal again and check that the result is the same
	marshaled, err := json.Marshal(&c)
	require.NoError(err)
	c2 := ChainConfig{}
	err = json.Unmarshal(marshaled, &c2)
	require.NoError(err)
	require.Equal(c, c2)
}

func TestActivePrecompiles(t *testing.T) {
	config := *WithExtra(
		&ChainConfig{},
		&extras.ChainConfig{
			UpgradeConfig: extras.UpgradeConfig{
				PrecompileUpgrades: []extras.PrecompileUpgrade{
					{
						Config: nativeminter.NewConfig(utils.NewUint64(0), nil, nil, nil, nil), // enable at genesis
					},
					{
						Config: nativeminter.NewDisableConfig(utils.NewUint64(1)), // disable at timestamp 1
					},
				},
			},
		},
	)

	rules0 := config.Rules(common.Big0, IsMergeTODO, 0)
	require.True(t, GetRulesExtra(rules0).IsPrecompileEnabled(nativeminter.Module.Address))

	rules1 := config.Rules(common.Big0, IsMergeTODO, 1)
	require.False(t, GetRulesExtra(rules1).IsPrecompileEnabled(nativeminter.Module.Address))
}

func TestChainConfigMarshalWithUpgrades(t *testing.T) {
	config := ChainConfigWithUpgradesJSON{
		ChainConfig: *WithExtra(
			&ChainConfig{
				ChainID:             big.NewInt(1),
				HomesteadBlock:      big.NewInt(0),
				EIP150Block:         big.NewInt(0),
				EIP155Block:         big.NewInt(0),
				EIP158Block:         big.NewInt(0),
				ByzantiumBlock:      big.NewInt(0),
				ConstantinopleBlock: big.NewInt(0),
				PetersburgBlock:     big.NewInt(0),
				IstanbulBlock:       big.NewInt(0),
				MuirGlacierBlock:    big.NewInt(0),
			},
			&extras.ChainConfig{
				FeeConfig:          DefaultFeeConfig,
				AllowFeeRecipients: false,
				NetworkUpgrades: extras.NetworkUpgrades{
					SubnetEVMTimestamp: utils.NewUint64(0),
					DurangoTimestamp:   utils.NewUint64(0),
				},
				GenesisPrecompiles: extras.Precompiles{},
			},
		),
		UpgradeConfig: extras.UpgradeConfig{
			PrecompileUpgrades: []extras.PrecompileUpgrade{
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

	var unmarshalled ChainConfigWithUpgradesJSON
	err = json.Unmarshal(result, &unmarshalled)
	require.NoError(t, err)
	require.Equal(t, config, unmarshalled)
}

// TestUpstreamParamsValues detects when a params value changes upstream to prevent a subtle change
// to one of the values to have an unpredicted impact in the libevm consumer.
// Values should be updated to newer upstream values once the consumer is updated to handle the
// updated value(s).
func TestUpstreamParamsValues(t *testing.T) {
	assert.Equal(t, uint64(1024), ethparams.GasLimitBoundDivisor, "GasLimitBoundDivisor")
	assert.Equal(t, uint64(5000), ethparams.MinGasLimit, "MinGasLimit")
	assert.Equal(t, uint64(0x7fffffffffffffff), ethparams.MaxGasLimit, "MaxGasLimit")
	assert.Equal(t, uint64(4712388), ethparams.GenesisGasLimit, "GenesisGasLimit")
	assert.Equal(t, uint64(32), ethparams.MaximumExtraDataSize, "MaximumExtraDataSize")
	assert.Equal(t, uint64(10), ethparams.ExpByteGas, "ExpByteGas")
	assert.Equal(t, uint64(50), ethparams.SloadGas, "SloadGas")
	assert.Equal(t, uint64(9000), ethparams.CallValueTransferGas, "CallValueTransferGas")
	assert.Equal(t, uint64(25000), ethparams.CallNewAccountGas, "CallNewAccountGas")
	assert.Equal(t, uint64(21000), ethparams.TxGas, "TxGas")
	assert.Equal(t, uint64(53000), ethparams.TxGasContractCreation, "TxGasContractCreation")
	assert.Equal(t, uint64(4), ethparams.TxDataZeroGas, "TxDataZeroGas")
	assert.Equal(t, uint64(512), ethparams.QuadCoeffDiv, "QuadCoeffDiv")
	assert.Equal(t, uint64(8), ethparams.LogDataGas, "LogDataGas")
	assert.Equal(t, uint64(2300), ethparams.CallStipend, "CallStipend")
	assert.Equal(t, uint64(30), ethparams.Keccak256Gas, "Keccak256Gas")
	assert.Equal(t, uint64(6), ethparams.Keccak256WordGas, "Keccak256WordGas")
	assert.Equal(t, uint64(2), ethparams.InitCodeWordGas, "InitCodeWordGas")
	assert.Equal(t, uint64(20000), ethparams.SstoreSetGas, "SstoreSetGas")
	assert.Equal(t, uint64(5000), ethparams.SstoreResetGas, "SstoreResetGas")
	assert.Equal(t, uint64(5000), ethparams.SstoreClearGas, "SstoreClearGas")
	assert.Equal(t, uint64(15000), ethparams.SstoreRefundGas, "SstoreRefundGas")
	assert.Equal(t, uint64(200), ethparams.NetSstoreNoopGas, "NetSstoreNoopGas")
	assert.Equal(t, uint64(20000), ethparams.NetSstoreInitGas, "NetSstoreInitGas")
	assert.Equal(t, uint64(5000), ethparams.NetSstoreCleanGas, "NetSstoreCleanGas")
	assert.Equal(t, uint64(200), ethparams.NetSstoreDirtyGas, "NetSstoreDirtyGas")
	assert.Equal(t, uint64(15000), ethparams.NetSstoreClearRefund, "NetSstoreClearRefund")
	assert.Equal(t, uint64(4800), ethparams.NetSstoreResetRefund, "NetSstoreResetRefund")
	assert.Equal(t, uint64(19800), ethparams.NetSstoreResetClearRefund, "NetSstoreResetClearRefund")
	assert.Equal(t, uint64(2300), ethparams.SstoreSentryGasEIP2200, "SstoreSentryGasEIP2200")
	assert.Equal(t, uint64(20000), ethparams.SstoreSetGasEIP2200, "SstoreSetGasEIP2200")
	assert.Equal(t, uint64(5000), ethparams.SstoreResetGasEIP2200, "SstoreResetGasEIP2200")
	assert.Equal(t, uint64(15000), ethparams.SstoreClearsScheduleRefundEIP2200, "SstoreClearsScheduleRefundEIP2200")
	assert.Equal(t, uint64(2600), ethparams.ColdAccountAccessCostEIP2929, "ColdAccountAccessCostEIP2929")
	assert.Equal(t, uint64(2100), ethparams.ColdSloadCostEIP2929, "ColdSloadCostEIP2929")
	assert.Equal(t, uint64(100), ethparams.WarmStorageReadCostEIP2929, "WarmStorageReadCostEIP2929")
	assert.Equal(t, uint64(5000-2100+1900), ethparams.SstoreClearsScheduleRefundEIP3529, "SstoreClearsScheduleRefundEIP3529")
	assert.Equal(t, uint64(1), ethparams.JumpdestGas, "JumpdestGas")
	assert.Equal(t, uint64(30000), ethparams.EpochDuration, "EpochDuration")
	assert.Equal(t, uint64(200), ethparams.CreateDataGas, "CreateDataGas")
	assert.Equal(t, uint64(1024), ethparams.CallCreateDepth, "CallCreateDepth")
	assert.Equal(t, uint64(10), ethparams.ExpGas, "ExpGas")
	assert.Equal(t, uint64(375), ethparams.LogGas, "LogGas")
	assert.Equal(t, uint64(3), ethparams.CopyGas, "CopyGas")
	assert.Equal(t, uint64(1024), ethparams.StackLimit, "StackLimit")
	assert.Equal(t, uint64(0), ethparams.TierStepGas, "TierStepGas")
	assert.Equal(t, uint64(375), ethparams.LogTopicGas, "LogTopicGas")
	assert.Equal(t, uint64(32000), ethparams.CreateGas, "CreateGas")
	assert.Equal(t, uint64(32000), ethparams.Create2Gas, "Create2Gas")
	assert.Equal(t, uint64(24000), ethparams.SelfdestructRefundGas, "SelfdestructRefundGas")
	assert.Equal(t, uint64(3), ethparams.MemoryGas, "MemoryGas")
	assert.Equal(t, uint64(68), ethparams.TxDataNonZeroGasFrontier, "TxDataNonZeroGasFrontier")
	assert.Equal(t, uint64(16), ethparams.TxDataNonZeroGasEIP2028, "TxDataNonZeroGasEIP2028")
	assert.Equal(t, uint64(2400), ethparams.TxAccessListAddressGas, "TxAccessListAddressGas")
	assert.Equal(t, uint64(1900), ethparams.TxAccessListStorageKeyGas, "TxAccessListStorageKeyGas")
	assert.Equal(t, uint64(40), ethparams.CallGasFrontier, "CallGasFrontier")
	assert.Equal(t, uint64(700), ethparams.CallGasEIP150, "CallGasEIP150")
	assert.Equal(t, uint64(20), ethparams.BalanceGasFrontier, "BalanceGasFrontier")
	assert.Equal(t, uint64(400), ethparams.BalanceGasEIP150, "BalanceGasEIP150")
	assert.Equal(t, uint64(700), ethparams.BalanceGasEIP1884, "BalanceGasEIP1884")
	assert.Equal(t, uint64(20), ethparams.ExtcodeSizeGasFrontier, "ExtcodeSizeGasFrontier")
	assert.Equal(t, uint64(700), ethparams.ExtcodeSizeGasEIP150, "ExtcodeSizeGasEIP150")
	assert.Equal(t, uint64(50), ethparams.SloadGasFrontier, "SloadGasFrontier")
	assert.Equal(t, uint64(200), ethparams.SloadGasEIP150, "SloadGasEIP150")
	assert.Equal(t, uint64(800), ethparams.SloadGasEIP1884, "SloadGasEIP1884")
	assert.Equal(t, uint64(800), ethparams.SloadGasEIP2200, "SloadGasEIP2200")
	assert.Equal(t, uint64(400), ethparams.ExtcodeHashGasConstantinople, "ExtcodeHashGasConstantinople")
	assert.Equal(t, uint64(700), ethparams.ExtcodeHashGasEIP1884, "ExtcodeHashGasEIP1884")
	assert.Equal(t, uint64(5000), ethparams.SelfdestructGasEIP150, "SelfdestructGasEIP150")
	assert.Equal(t, uint64(10), ethparams.ExpByteFrontier, "ExpByteFrontier")
	assert.Equal(t, uint64(50), ethparams.ExpByteEIP158, "ExpByteEIP158")
	assert.Equal(t, uint64(20), ethparams.ExtcodeCopyBaseFrontier, "ExtcodeCopyBaseFrontier")
	assert.Equal(t, uint64(700), ethparams.ExtcodeCopyBaseEIP150, "ExtcodeCopyBaseEIP150")
	assert.Equal(t, uint64(25000), ethparams.CreateBySelfdestructGas, "CreateBySelfdestructGas")
	assert.Equal(t, 8, ethparams.DefaultBaseFeeChangeDenominator, "DefaultBaseFeeChangeDenominator")
	assert.Equal(t, 2, ethparams.DefaultElasticityMultiplier, "DefaultElasticityMultiplier")
	assert.Equal(t, 1000000000, ethparams.InitialBaseFee, "InitialBaseFee")
	assert.Equal(t, 24576, ethparams.MaxCodeSize, "MaxCodeSize")
	assert.Equal(t, 2*24576, ethparams.MaxInitCodeSize, "MaxInitCodeSize")
	assert.Equal(t, uint64(3000), ethparams.EcrecoverGas, "EcrecoverGas")
	assert.Equal(t, uint64(60), ethparams.Sha256BaseGas, "Sha256BaseGas")
	assert.Equal(t, uint64(12), ethparams.Sha256PerWordGas, "Sha256PerWordGas")
	assert.Equal(t, uint64(600), ethparams.Ripemd160BaseGas, "Ripemd160BaseGas")
	assert.Equal(t, uint64(120), ethparams.Ripemd160PerWordGas, "Ripemd160PerWordGas")
	assert.Equal(t, uint64(15), ethparams.IdentityBaseGas, "IdentityBaseGas")
	assert.Equal(t, uint64(3), ethparams.IdentityPerWordGas, "IdentityPerWordGas")
	assert.Equal(t, uint64(500), ethparams.Bn256AddGasByzantium, "Bn256AddGasByzantium")
	assert.Equal(t, uint64(150), ethparams.Bn256AddGasIstanbul, "Bn256AddGasIstanbul")
	assert.Equal(t, uint64(40000), ethparams.Bn256ScalarMulGasByzantium, "Bn256ScalarMulGasByzantium")
	assert.Equal(t, uint64(6000), ethparams.Bn256ScalarMulGasIstanbul, "Bn256ScalarMulGasIstanbul")
	assert.Equal(t, uint64(100000), ethparams.Bn256PairingBaseGasByzantium, "Bn256PairingBaseGasByzantium")
	assert.Equal(t, uint64(45000), ethparams.Bn256PairingBaseGasIstanbul, "Bn256PairingBaseGasIstanbul")
	assert.Equal(t, uint64(80000), ethparams.Bn256PairingPerPointGasByzantium, "Bn256PairingPerPointGasByzantium")
	assert.Equal(t, uint64(34000), ethparams.Bn256PairingPerPointGasIstanbul, "Bn256PairingPerPointGasIstanbul")
	assert.Equal(t, uint64(600), ethparams.Bls12381G1AddGas, "Bls12381G1AddGas")
	assert.Equal(t, uint64(12000), ethparams.Bls12381G1MulGas, "Bls12381G1MulGas")
	assert.Equal(t, uint64(4500), ethparams.Bls12381G2AddGas, "Bls12381G2AddGas")
	assert.Equal(t, uint64(55000), ethparams.Bls12381G2MulGas, "Bls12381G2MulGas")
	assert.Equal(t, uint64(115000), ethparams.Bls12381PairingBaseGas, "Bls12381PairingBaseGas")
	assert.Equal(t, uint64(23000), ethparams.Bls12381PairingPerPairGas, "Bls12381PairingPerPairGas")
	assert.Equal(t, uint64(5500), ethparams.Bls12381MapG1Gas, "Bls12381MapG1Gas")
	assert.Equal(t, uint64(110000), ethparams.Bls12381MapG2Gas, "Bls12381MapG2Gas")
	assert.Equal(t, uint64(2), ethparams.RefundQuotient, "RefundQuotient")
	assert.Equal(t, uint64(5), ethparams.RefundQuotientEIP3529, "RefundQuotientEIP3529")
	assert.Equal(t, 32, ethparams.BlobTxBytesPerFieldElement, "BlobTxBytesPerFieldElement")
	assert.Equal(t, 4096, ethparams.BlobTxFieldElementsPerBlob, "BlobTxFieldElementsPerBlob")
	assert.Equal(t, 1<<17, ethparams.BlobTxBlobGasPerBlob, "BlobTxBlobGasPerBlob")
	assert.Equal(t, 1, ethparams.BlobTxMinBlobGasprice, "BlobTxMinBlobGasprice")
	assert.Equal(t, 3338477, ethparams.BlobTxBlobGaspriceUpdateFraction, "BlobTxBlobGaspriceUpdateFraction")
	assert.Equal(t, 50000, ethparams.BlobTxPointEvaluationPrecompileGas, "BlobTxPointEvaluationPrecompileGas")
	assert.Equal(t, 3*131072, ethparams.BlobTxTargetBlobGasPerBlock, "BlobTxTargetBlobGasPerBlock")
	assert.Equal(t, 6*131072, ethparams.MaxBlobGasPerBlock, "MaxBlobGasPerBlock")
	assert.Equal(t, int64(131072), ethparams.GenesisDifficulty.Int64(), "GenesisDifficulty")
	assert.Equal(t, common.HexToAddress("0x000F3df6D732807Ef1319fB7B8bB8522d0Beac02"), ethparams.BeaconRootsStorageAddress, "BeaconRootsStorageAddress")
	assert.Equal(t, common.HexToAddress("0xfffffffffffffffffffffffffffffffffffffffe"), ethparams.SystemAddress, "SystemAddress")
}
