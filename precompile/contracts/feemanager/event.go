// Code generated
// This file is a generated precompile contract config with stubbed abstract functions.
// The file is generated by a template. Please inspect every code and comment in this file before use.

package feemanager

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ethereum/go-ethereum/common"
)

// ContractChangeFeeConfig represents a ChangeFeeConfig non-indexed event data raised by the Contract contract.
type ChangeFeeConfigEventData struct {
	GasLimit                 *big.Int
	TargetBlockRate          *big.Int
	MinBaseFee               *big.Int
	TargetGas                *big.Int
	BaseFeeChangeDenominator *big.Int
	MinBlockGasCost          *big.Int
	MaxBlockGasCost          *big.Int
	BlockGasCostStep         *big.Int
}

// PackChangeFeeConfigEvent packs the event into the appropriate arguments for changeFeeConfig.
// It returns topic hashes and the encoded non-indexed data.
func PackChangeFeeConfigEvent(oldConfig commontype.FeeConfig, newConfig commontype.FeeConfig) ([]common.Hash, []byte, error) {
	oldConfigC := convertFromCommonConfig(oldConfig)
	newConfigC := convertFromCommonConfig(newConfig)
	return FeeManagerABI.PackEvent("FeeConfigChanged", oldConfigC, newConfigC)
}

// UnpackChangeFeeConfigEventData attempts to unpack non-indexed [dataBytes].
func UnpackChangeFeeConfigEventData(dataBytes []byte) (ChangeFeeConfigEventData, ChangeFeeConfigEventData, error) {
	eventData := make([]ChangeFeeConfigEventData, 2)
	err := FeeManagerABI.UnpackIntoInterface(&eventData, "FeeConfigChanged", dataBytes)
	return eventData[0], eventData[1], err
}

func convertFromCommonConfig(config commontype.FeeConfig) ChangeFeeConfigEventData {
	return ChangeFeeConfigEventData{
		GasLimit:                 config.GasLimit,
		TargetBlockRate:          new(big.Int).SetUint64(config.TargetBlockRate),
		MinBaseFee:               config.MinBaseFee,
		TargetGas:                config.TargetGas,
		BaseFeeChangeDenominator: config.BaseFeeChangeDenominator,
		MinBlockGasCost:          config.MinBlockGasCost,
		MaxBlockGasCost:          config.MaxBlockGasCost,
		BlockGasCostStep:         config.BlockGasCostStep,
	}
}

func convertToCommonConfig(config ChangeFeeConfigEventData) commontype.FeeConfig {
	return commontype.FeeConfig{
		GasLimit:                 config.GasLimit,
		TargetBlockRate:          config.TargetBlockRate.Uint64(),
		MinBaseFee:               config.MinBaseFee,
		TargetGas:                config.TargetGas,
		BaseFeeChangeDenominator: config.BaseFeeChangeDenominator,
		MinBlockGasCost:          config.MinBlockGasCost,
		MaxBlockGasCost:          config.MaxBlockGasCost,
		BlockGasCostStep:         config.BlockGasCostStep,
	}
}
