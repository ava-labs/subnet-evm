// (c) 2024 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package extras

import (
	"math/big"

	"github.com/ava-labs/avalanchego/upgrade"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/utils"
)

var (
	DefaultSubnetEVMChainID = big.NewInt(43214)
	DefaultFeeConfig        = commontype.FeeConfig{
		GasLimit:        big.NewInt(8_000_000),
		TargetBlockRate: 2, // in seconds

		MinBaseFee:               big.NewInt(25_000_000_000),
		TargetGas:                big.NewInt(15_000_000),
		BaseFeeChangeDenominator: big.NewInt(36),

		MinBlockGasCost:  big.NewInt(0),
		MaxBlockGasCost:  big.NewInt(1_000_000),
		BlockGasCostStep: big.NewInt(200_000),
	}

	SubnetEVMDefaultChainConfig = &ChainConfig{
		FeeConfig:          DefaultFeeConfig,
		NetworkUpgrades:    SubnetEVMDefaultNetworkUpgrades(upgrade.GetConfig(constants.MainnetID)),
		GenesisPrecompiles: Precompiles{},
	}

	TestChainConfig = &ChainConfig{
		AvalancheContext: AvalancheContext{SnowCtx: utils.TestSnowContext()},
		FeeConfig:        DefaultFeeConfig,
		NetworkUpgrades: NetworkUpgrades{
			ApricotPhase1BlockTimestamp:     utils.NewUint64(0),
			ApricotPhase2BlockTimestamp:     utils.NewUint64(0),
			ApricotPhase3BlockTimestamp:     utils.NewUint64(0),
			ApricotPhase4BlockTimestamp:     utils.NewUint64(0),
			ApricotPhase5BlockTimestamp:     utils.NewUint64(0),
			ApricotPhasePre6BlockTimestamp:  utils.NewUint64(0),
			ApricotPhase6BlockTimestamp:     utils.NewUint64(0),
			ApricotPhasePost6BlockTimestamp: utils.NewUint64(0),
			BanffBlockTimestamp:             utils.NewUint64(0),
			CortinaBlockTimestamp:           utils.NewUint64(0),
			SubnetEVMTimestamp:              utils.NewUint64(0),
			DurangoTimestamp:                utils.NewUint64(0),
			EtnaTimestamp:                   utils.NewUint64(0),
			FortunaTimestamp:                utils.NewUint64(0),
		},
		GenesisPrecompiles: Precompiles{},
	}

	TestCChainLaunchConfig = &ChainConfig{
		AvalancheContext: AvalancheContext{SnowCtx: utils.TestSnowContext()},
	}

	TestApricotPhase1Config = copyAndSet(TestCChainLaunchConfig, func(c *ChainConfig) {
		c.NetworkUpgrades.ApricotPhase1BlockTimestamp = utils.NewUint64(0)
	})

	TestApricotPhase2Config = copyAndSet(TestApricotPhase1Config, func(c *ChainConfig) {
		c.NetworkUpgrades.ApricotPhase2BlockTimestamp = utils.NewUint64(0)
	})

	TestApricotPhase3Config = copyAndSet(TestApricotPhase2Config, func(c *ChainConfig) {
		c.NetworkUpgrades.ApricotPhase3BlockTimestamp = utils.NewUint64(0)
	})

	TestApricotPhase4Config = copyAndSet(TestApricotPhase3Config, func(c *ChainConfig) {
		c.NetworkUpgrades.ApricotPhase4BlockTimestamp = utils.NewUint64(0)
	})

	TestApricotPhase5Config = copyAndSet(TestApricotPhase4Config, func(c *ChainConfig) {
		c.NetworkUpgrades.ApricotPhase5BlockTimestamp = utils.NewUint64(0)
	})

	TestApricotPhasePre6Config = copyAndSet(TestApricotPhase5Config, func(c *ChainConfig) {
		c.NetworkUpgrades.ApricotPhasePre6BlockTimestamp = utils.NewUint64(0)
	})

	TestApricotPhase6Config = copyAndSet(TestApricotPhasePre6Config, func(c *ChainConfig) {
		c.NetworkUpgrades.ApricotPhase6BlockTimestamp = utils.NewUint64(0)
	})

	TestApricotPhasePost6Config = copyAndSet(TestApricotPhase6Config, func(c *ChainConfig) {
		c.NetworkUpgrades.ApricotPhasePost6BlockTimestamp = utils.NewUint64(0)
	})

	TestBanffChainConfig = copyAndSet(TestApricotPhasePost6Config, func(c *ChainConfig) {
		c.NetworkUpgrades.BanffBlockTimestamp = utils.NewUint64(0)
	})

	TestCortinaChainConfig = copyAndSet(TestBanffChainConfig, func(c *ChainConfig) {
		c.NetworkUpgrades.CortinaBlockTimestamp = utils.NewUint64(0)
	})

	TestPreSubnetEVMChainConfig = copyAndSet(TestCortinaChainConfig, func(c *ChainConfig) {
		c.NetworkUpgrades.SubnetEVMTimestamp = nil
	})

	TestSubnetEVMChainConfig = copyAndSet(TestCortinaChainConfig, func(c *ChainConfig) {
		c.FeeConfig = DefaultFeeConfig
		c.NetworkUpgrades.SubnetEVMTimestamp = utils.NewUint64(0)
	})

	TestDurangoChainConfig = copyAndSet(TestSubnetEVMChainConfig, func(c *ChainConfig) {
		c.NetworkUpgrades.DurangoTimestamp = utils.NewUint64(0)
	})

	TestEtnaChainConfig = copyAndSet(TestDurangoChainConfig, func(c *ChainConfig) {
		c.NetworkUpgrades.EtnaTimestamp = utils.NewUint64(0)
	})

	TestFortunaChainConfig = copyAndSet(TestEtnaChainConfig, func(c *ChainConfig) {
		c.NetworkUpgrades.FortunaTimestamp = utils.NewUint64(0)
	})
)

func copyAndSet(c *ChainConfig, set func(*ChainConfig)) *ChainConfig {
	newConfig := *c
	set(&newConfig)
	return &newConfig
}
