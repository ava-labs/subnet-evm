package commontype

import "math/big"

// This struct is used by params.Config and precompile.FeeConfigManager
// any modification of this struct has direct affect on the precompiled contract
// and changes should be carefully handled in the precompiled contract code.
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

var ZeroFeeConfig = FeeConfig{
	GasLimit:                 big.NewInt(0),
	TargetBlockRate:          0,
	MinBaseFee:               big.NewInt(0),
	TargetGas:                big.NewInt(0),
	BaseFeeChangeDenominator: big.NewInt(0),
	MinBlockGasCost:          big.NewInt(0),
	MaxBlockGasCost:          big.NewInt(0),
	BlockGasCostStep:         big.NewInt(0),
}
