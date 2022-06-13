package commontype

import "math/big"

// FeeConfig specifies the parameters for the dynamic fee algorithm, which determines the gas limit, base fee, and block gas cost of blocks
// on the network.
//
// The dynamic fee algorithm simply increases fees when the network is operating at a utilization level above the target and decreases fees
// when the network is operating at a utilization level below the target.
// This struct is used by params.Config and precompile.FeeConfigManager
// any modification of this struct has direct affect on the precompiled contract
// and changes should be carefully handled in the precompiled contract code.
type FeeConfig struct {
	// GasLimit sets the max amount of gas consumed per block.
	GasLimit *big.Int `json:"gasLimit,omitempty"`

	// TargetBlockRate sets the target rate of block production in seconds.
	// A target of 2 will target producing a block every 2 seconds.
	TargetBlockRate uint64 `json:"targetBlockRate,omitempty"`

	// The minimum base fee sets a lower bound on the EIP-1559 base fee of a block.
	// Since the block's base fee sets the minimum gas price for any transaction included in that block, this effectively sets a minimum
	// gas price for any tranasction.
	MinBaseFee *big.Int `json:"minBaseFee,omitempty"`

	// When the dynamic fee algorithm observes that network activity is above/below the [TargetGas], it increases/decreases the base fee proportionally to
	// how far above/below the target actual network activity is.

	// TargetGas specifies the targeted amount of gas (including block gas cost) to consume within a rolling 10s window.
	TargetGas *big.Int `json:"targetGas,omitempty"`
	// The BaseFeeChangeDenominator divides the difference between actual and target utilization to determine how much to increase/decrease the base fee.
	// This means that a larger denominator indicates a slower changing, stickier base fee, while a lower denominator will allow the base fee to adjust
	// more quickly.
	BaseFeeChangeDenominator *big.Int `json:"baseFeeChangeDenominator,omitempty"`

	// MinBlockGasCost sets the minimum amount of gas to charge for the production of a block.
	MinBlockGasCost *big.Int `json:"minBlockGasCost,omitempty"`
	// MaxBlockGasCost sets the maximum amount of gas to charge for the production of a block.
	MaxBlockGasCost *big.Int `json:"maxBlockGasCost,omitempty"`
	// BlockGasCostStep determines how much to increase/decrease the block gas cost depending on the amount of time elapsed since the previous block.
	// If the block is produced at the target rate, the block gas cost will stay the same as the block gas cost for the parent block.
	// If it is produced faster/slower, the block gas cost will be increased/decreased by the step value for each second faster/slower than the target
	// block rate accordingly.
	// Note: if the BlockGasCostStep is set to a very large number, it effectively requires block production to go no faster than the TargetBlockRate.
	//
	// Ex: if a block is produced two seconds faster than the target block rate, the block gas cost will increase by 2 * BlockGasCostStep.
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
