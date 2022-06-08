package commontype

import "math/big"

type FeeConfig struct {
	GasLimit        *big.Int `json:"gasLimit,omitempty"`
	TargetBlockRate uint64   `json:"targetBlockRate,omitempty"`

	MinBaseFee               *big.Int `json:"minBaseFee,omitempty"`
	TargetGas                *big.Int `json:"targetGas,omitempty"`
	BaseFeeChangeDenominator *big.Int `json:"baseFeeChangeDenominator,omitempty"`

	MinBlockGasCost  *big.Int `json:"minBlockGasCost,omitempty"`
	MaxBlockGasCost  *big.Int `json:"maxBlockGasCost,omitempty"`
	BlockGasCostStep *big.Int `json:"blockGasCostStep,omitempty"`
}

var EmptyFeeConfig = FeeConfig{}
