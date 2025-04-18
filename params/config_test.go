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
	tests := map[string]struct {
		param any
		want  any
	}{
		"GasLimitBoundDivisor":               {param: ethparams.GasLimitBoundDivisor, want: uint64(1024)},
		"MinGasLimit":                        {param: ethparams.MinGasLimit, want: uint64(5000)},
		"MaxGasLimit":                        {param: ethparams.MaxGasLimit, want: uint64(0x7fffffffffffffff)},
		"GenesisGasLimit":                    {param: ethparams.GenesisGasLimit, want: uint64(4712388)},
		"MaximumExtraDataSize":               {param: ethparams.MaximumExtraDataSize, want: uint64(32)},
		"ExpByteGas":                         {param: ethparams.ExpByteGas, want: uint64(10)},
		"SloadGas":                           {param: ethparams.SloadGas, want: uint64(50)},
		"CallValueTransferGas":               {param: ethparams.CallValueTransferGas, want: uint64(9000)},
		"CallNewAccountGas":                  {param: ethparams.CallNewAccountGas, want: uint64(25000)},
		"TxGas":                              {param: ethparams.TxGas, want: uint64(21000)},
		"TxGasContractCreation":              {param: ethparams.TxGasContractCreation, want: uint64(53000)},
		"TxDataZeroGas":                      {param: ethparams.TxDataZeroGas, want: uint64(4)},
		"QuadCoeffDiv":                       {param: ethparams.QuadCoeffDiv, want: uint64(512)},
		"LogDataGas":                         {param: ethparams.LogDataGas, want: uint64(8)},
		"CallStipend":                        {param: ethparams.CallStipend, want: uint64(2300)},
		"Keccak256Gas":                       {param: ethparams.Keccak256Gas, want: uint64(30)},
		"Keccak256WordGas":                   {param: ethparams.Keccak256WordGas, want: uint64(6)},
		"InitCodeWordGas":                    {param: ethparams.InitCodeWordGas, want: uint64(2)},
		"SstoreSetGas":                       {param: ethparams.SstoreSetGas, want: uint64(20000)},
		"SstoreResetGas":                     {param: ethparams.SstoreResetGas, want: uint64(5000)},
		"SstoreClearGas":                     {param: ethparams.SstoreClearGas, want: uint64(5000)},
		"SstoreRefundGas":                    {param: ethparams.SstoreRefundGas, want: uint64(15000)},
		"NetSstoreNoopGas":                   {param: ethparams.NetSstoreNoopGas, want: uint64(200)},
		"NetSstoreInitGas":                   {param: ethparams.NetSstoreInitGas, want: uint64(20000)},
		"NetSstoreCleanGas":                  {param: ethparams.NetSstoreCleanGas, want: uint64(5000)},
		"NetSstoreDirtyGas":                  {param: ethparams.NetSstoreDirtyGas, want: uint64(200)},
		"NetSstoreClearRefund":               {param: ethparams.NetSstoreClearRefund, want: uint64(15000)},
		"NetSstoreResetRefund":               {param: ethparams.NetSstoreResetRefund, want: uint64(4800)},
		"NetSstoreResetClearRefund":          {param: ethparams.NetSstoreResetClearRefund, want: uint64(19800)},
		"SstoreSentryGasEIP2200":             {param: ethparams.SstoreSentryGasEIP2200, want: uint64(2300)},
		"SstoreSetGasEIP2200":                {param: ethparams.SstoreSetGasEIP2200, want: uint64(20000)},
		"SstoreResetGasEIP2200":              {param: ethparams.SstoreResetGasEIP2200, want: uint64(5000)},
		"SstoreClearsScheduleRefundEIP2200":  {param: ethparams.SstoreClearsScheduleRefundEIP2200, want: uint64(15000)},
		"ColdAccountAccessCostEIP2929":       {param: ethparams.ColdAccountAccessCostEIP2929, want: uint64(2600)},
		"ColdSloadCostEIP2929":               {param: ethparams.ColdSloadCostEIP2929, want: uint64(2100)},
		"WarmStorageReadCostEIP2929":         {param: ethparams.WarmStorageReadCostEIP2929, want: uint64(100)},
		"SstoreClearsScheduleRefundEIP3529":  {param: ethparams.SstoreClearsScheduleRefundEIP3529, want: uint64(5000 - 2100 + 1900)},
		"JumpdestGas":                        {param: ethparams.JumpdestGas, want: uint64(1)},
		"EpochDuration":                      {param: ethparams.EpochDuration, want: uint64(30000)},
		"CreateDataGas":                      {param: ethparams.CreateDataGas, want: uint64(200)},
		"CallCreateDepth":                    {param: ethparams.CallCreateDepth, want: uint64(1024)},
		"ExpGas":                             {param: ethparams.ExpGas, want: uint64(10)},
		"LogGas":                             {param: ethparams.LogGas, want: uint64(375)},
		"CopyGas":                            {param: ethparams.CopyGas, want: uint64(3)},
		"StackLimit":                         {param: ethparams.StackLimit, want: uint64(1024)},
		"TierStepGas":                        {param: ethparams.TierStepGas, want: uint64(0)},
		"LogTopicGas":                        {param: ethparams.LogTopicGas, want: uint64(375)},
		"CreateGas":                          {param: ethparams.CreateGas, want: uint64(32000)},
		"Create2Gas":                         {param: ethparams.Create2Gas, want: uint64(32000)},
		"SelfdestructRefundGas":              {param: ethparams.SelfdestructRefundGas, want: uint64(24000)},
		"MemoryGas":                          {param: ethparams.MemoryGas, want: uint64(3)},
		"TxDataNonZeroGasFrontier":           {param: ethparams.TxDataNonZeroGasFrontier, want: uint64(68)},
		"TxDataNonZeroGasEIP2028":            {param: ethparams.TxDataNonZeroGasEIP2028, want: uint64(16)},
		"TxAccessListAddressGas":             {param: ethparams.TxAccessListAddressGas, want: uint64(2400)},
		"TxAccessListStorageKeyGas":          {param: ethparams.TxAccessListStorageKeyGas, want: uint64(1900)},
		"CallGasFrontier":                    {param: ethparams.CallGasFrontier, want: uint64(40)},
		"CallGasEIP150":                      {param: ethparams.CallGasEIP150, want: uint64(700)},
		"BalanceGasFrontier":                 {param: ethparams.BalanceGasFrontier, want: uint64(20)},
		"BalanceGasEIP150":                   {param: ethparams.BalanceGasEIP150, want: uint64(400)},
		"BalanceGasEIP1884":                  {param: ethparams.BalanceGasEIP1884, want: uint64(700)},
		"ExtcodeSizeGasFrontier":             {param: ethparams.ExtcodeSizeGasFrontier, want: uint64(20)},
		"ExtcodeSizeGasEIP150":               {param: ethparams.ExtcodeSizeGasEIP150, want: uint64(700)},
		"SloadGasFrontier":                   {param: ethparams.SloadGasFrontier, want: uint64(50)},
		"SloadGasEIP150":                     {param: ethparams.SloadGasEIP150, want: uint64(200)},
		"SloadGasEIP1884":                    {param: ethparams.SloadGasEIP1884, want: uint64(800)},
		"SloadGasEIP2200":                    {param: ethparams.SloadGasEIP2200, want: uint64(800)},
		"ExtcodeHashGasConstantinople":       {param: ethparams.ExtcodeHashGasConstantinople, want: uint64(400)},
		"ExtcodeHashGasEIP1884":              {param: ethparams.ExtcodeHashGasEIP1884, want: uint64(700)},
		"SelfdestructGasEIP150":              {param: ethparams.SelfdestructGasEIP150, want: uint64(5000)},
		"ExpByteFrontier":                    {param: ethparams.ExpByteFrontier, want: uint64(10)},
		"ExpByteEIP158":                      {param: ethparams.ExpByteEIP158, want: uint64(50)},
		"ExtcodeCopyBaseFrontier":            {param: ethparams.ExtcodeCopyBaseFrontier, want: uint64(20)},
		"ExtcodeCopyBaseEIP150":              {param: ethparams.ExtcodeCopyBaseEIP150, want: uint64(700)},
		"CreateBySelfdestructGas":            {param: ethparams.CreateBySelfdestructGas, want: uint64(25000)},
		"DefaultBaseFeeChangeDenominator":    {param: ethparams.DefaultBaseFeeChangeDenominator, want: 8},
		"DefaultElasticityMultiplier":        {param: ethparams.DefaultElasticityMultiplier, want: 2},
		"InitialBaseFee":                     {param: ethparams.InitialBaseFee, want: 1000000000},
		"MaxCodeSize":                        {param: ethparams.MaxCodeSize, want: 24576},
		"MaxInitCodeSize":                    {param: ethparams.MaxInitCodeSize, want: 2 * 24576},
		"EcrecoverGas":                       {param: ethparams.EcrecoverGas, want: uint64(3000)},
		"Sha256BaseGas":                      {param: ethparams.Sha256BaseGas, want: uint64(60)},
		"Sha256PerWordGas":                   {param: ethparams.Sha256PerWordGas, want: uint64(12)},
		"Ripemd160BaseGas":                   {param: ethparams.Ripemd160BaseGas, want: uint64(600)},
		"Ripemd160PerWordGas":                {param: ethparams.Ripemd160PerWordGas, want: uint64(120)},
		"IdentityBaseGas":                    {param: ethparams.IdentityBaseGas, want: uint64(15)},
		"IdentityPerWordGas":                 {param: ethparams.IdentityPerWordGas, want: uint64(3)},
		"Bn256AddGasByzantium":               {param: ethparams.Bn256AddGasByzantium, want: uint64(500)},
		"Bn256AddGasIstanbul":                {param: ethparams.Bn256AddGasIstanbul, want: uint64(150)},
		"Bn256ScalarMulGasByzantium":         {param: ethparams.Bn256ScalarMulGasByzantium, want: uint64(40000)},
		"Bn256ScalarMulGasIstanbul":          {param: ethparams.Bn256ScalarMulGasIstanbul, want: uint64(6000)},
		"Bn256PairingBaseGasByzantium":       {param: ethparams.Bn256PairingBaseGasByzantium, want: uint64(100000)},
		"Bn256PairingBaseGasIstanbul":        {param: ethparams.Bn256PairingBaseGasIstanbul, want: uint64(45000)},
		"Bn256PairingPerPointGasByzantium":   {param: ethparams.Bn256PairingPerPointGasByzantium, want: uint64(80000)},
		"Bn256PairingPerPointGasIstanbul":    {param: ethparams.Bn256PairingPerPointGasIstanbul, want: uint64(34000)},
		"Bls12381G1AddGas":                   {param: ethparams.Bls12381G1AddGas, want: uint64(600)},
		"Bls12381G1MulGas":                   {param: ethparams.Bls12381G1MulGas, want: uint64(12000)},
		"Bls12381G2AddGas":                   {param: ethparams.Bls12381G2AddGas, want: uint64(4500)},
		"Bls12381G2MulGas":                   {param: ethparams.Bls12381G2MulGas, want: uint64(55000)},
		"Bls12381PairingBaseGas":             {param: ethparams.Bls12381PairingBaseGas, want: uint64(115000)},
		"Bls12381PairingPerPairGas":          {param: ethparams.Bls12381PairingPerPairGas, want: uint64(23000)},
		"Bls12381MapG1Gas":                   {param: ethparams.Bls12381MapG1Gas, want: uint64(5500)},
		"Bls12381MapG2Gas":                   {param: ethparams.Bls12381MapG2Gas, want: uint64(110000)},
		"RefundQuotient":                     {param: ethparams.RefundQuotient, want: uint64(2)},
		"RefundQuotientEIP3529":              {param: ethparams.RefundQuotientEIP3529, want: uint64(5)},
		"BlobTxBytesPerFieldElement":         {param: ethparams.BlobTxBytesPerFieldElement, want: 32},
		"BlobTxFieldElementsPerBlob":         {param: ethparams.BlobTxFieldElementsPerBlob, want: 4096},
		"BlobTxBlobGasPerBlob":               {param: ethparams.BlobTxBlobGasPerBlob, want: 1 << 17},
		"BlobTxMinBlobGasprice":              {param: ethparams.BlobTxMinBlobGasprice, want: 1},
		"BlobTxBlobGaspriceUpdateFraction":   {param: ethparams.BlobTxBlobGaspriceUpdateFraction, want: 3338477},
		"BlobTxPointEvaluationPrecompileGas": {param: ethparams.BlobTxPointEvaluationPrecompileGas, want: 50000},
		"BlobTxTargetBlobGasPerBlock":        {param: ethparams.BlobTxTargetBlobGasPerBlock, want: 3 * 131072},
		"MaxBlobGasPerBlock":                 {param: ethparams.MaxBlobGasPerBlock, want: 6 * 131072},
		"GenesisDifficulty":                  {param: ethparams.GenesisDifficulty.Int64(), want: int64(131072)},
		"BeaconRootsStorageAddress":          {param: ethparams.BeaconRootsStorageAddress, want: common.HexToAddress("0x000F3df6D732807Ef1319fB7B8bB8522d0Beac02")},
		"SystemAddress":                      {param: ethparams.SystemAddress, want: common.HexToAddress("0xfffffffffffffffffffffffffffffffffffffffe")},
	}

	for name, test := range tests {
		assert.Equal(t, test.want, test.param, name)
	}
}
