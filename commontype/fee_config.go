// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package commontype

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/avalanchego/vms/components/gas"
	"github.com/ava-labs/avalanchego/vms/evm/upgrade/acp176"
	"github.com/ava-labs/libevm/common"

	"github.com/ava-labs/subnet-evm/utils"
)

type ACP224FeeConfig struct {
	TargetGas          *big.Int `json:"targetGas,omitempty"`
	MinGasPrice        *big.Int `json:"minGasPrice,omitempty"`
	TimeToFillCapacity *big.Int `json:"timeToFillCapacity,omitempty"`
	TimeToDouble       *big.Int `json:"timeToDouble,omitempty"`
}

// represents an empty ACP-224 fee config without any field
var EmptyACP224FeeConfig = ACP224FeeConfig{}

func (f *ACP224FeeConfig) Verify() error {
	switch {
	case f.TargetGas == nil:
		return errors.New("targetGas cannot be nil")
	case f.MinGasPrice == nil:
		return errors.New("minGasPrice cannot be nil")
	case f.TimeToFillCapacity == nil:
		return errors.New("timeToFillCapacity cannot be nil")
	case f.TimeToDouble == nil:
		return errors.New("timeToDouble cannot be nil")
	}

	switch {
	case f.TargetGas.Cmp(common.Big0) != 1:
		return fmt.Errorf("targetGas = %d cannot be less than or equal to 0", f.TargetGas)
	case f.MinGasPrice.Cmp(common.Big0) != 1:
		return fmt.Errorf("minGasPrice = %d cannot be less than or equal to 0", f.MinGasPrice)
	case f.TimeToFillCapacity.Cmp(common.Big0) == -1:
		return fmt.Errorf("timeToFillCapacity = %d cannot be less than 0", f.TimeToFillCapacity)
	case f.TimeToDouble.Cmp(common.Big0) == -1:
		return fmt.Errorf("timeToDouble = %d cannot be less than 0", f.TimeToDouble)
	}

	switch {
	case f.TimeToFillCapacity.Cmp(big.NewInt(acp176.MaxTimeToFillCapacity)) == 1:
		return fmt.Errorf("timeToFillCapacity = %d cannot be greater than %d", f.TimeToFillCapacity, acp176.MaxTimeToFillCapacity)
	case f.TimeToDouble.Cmp(big.NewInt(acp176.MaxTimeToDouble)) == 1:
		return fmt.Errorf("timeToDouble = %d cannot be greater than %d", f.TimeToDouble, acp176.MaxTimeToDouble)
	}
	return f.checkByteLens()
}

func (f *ACP224FeeConfig) Equal(other *ACP224FeeConfig) bool {
	if other == nil {
		return false
	}

	return utils.BigNumEqual(f.TargetGas, other.TargetGas) &&
		utils.BigNumEqual(f.MinGasPrice, other.MinGasPrice) &&
		utils.BigNumEqual(f.TimeToFillCapacity, other.TimeToFillCapacity) &&
		utils.BigNumEqual(f.TimeToDouble, other.TimeToDouble)
}

func (f *ACP224FeeConfig) ToACP176Config() (acp176.Config, error) {
	if err := f.Verify(); err != nil {
		return acp176.Config{}, err
	}

	config := acp176.Config{
		MinGasPrice:        gas.Price(f.MinGasPrice.Uint64()),
		TimeToFillCapacity: gas.Gas(f.TimeToFillCapacity.Uint64()),
		TimeToDouble:       f.TimeToDouble.Uint64(),
	}
	return config, config.Verify()
}

// checkByteLens checks byte lengths against common.HashLen (32 bytes) and returns error
func (f *ACP224FeeConfig) checkByteLens() error {
	if isBiggerThanHashLen(f.TargetGas) {
		return fmt.Errorf("targetGas exceeds %d bytes", common.HashLength)
	}
	if isBiggerThanHashLen(f.MinGasPrice) {
		return fmt.Errorf("minGasPrice exceeds %d bytes", common.HashLength)
	}
	if isBiggerThanHashLen(f.TimeToFillCapacity) {
		return fmt.Errorf("timeToFillCapacity exceeds %d bytes", common.HashLength)
	}
	if isBiggerThanHashLen(f.TimeToDouble) {
		return fmt.Errorf("timeToDouble exceeds %d bytes", common.HashLength)
	}
	return nil
}

// FeeConfig specifies the parameters for the dynamic fee algorithm, which determines the gas limit, base fee, and block gas cost of blocks
// on the network.
//
// The dynamic fee algorithm simply increases fees when the network is operating at a utilization level above the target and decreases fees
// when the network is operating at a utilization level below the target.
// This struct is used by Genesis and Fee Manager precompile.
// Any modification of this struct has direct affect on the precompiled contract
// and changes should be carefully handled in the precompiled contract code.
type FeeConfig struct {
	// GasLimit sets the max amount of gas consumed per block.
	GasLimit *big.Int `json:"gasLimit,omitempty"`

	// TargetBlockRate sets the target rate of block production in seconds.
	// A target of 2 will target producing a block every 2 seconds.
	TargetBlockRate uint64 `json:"targetBlockRate,omitempty"`

	// The minimum base fee sets a lower bound on the EIP-1559 base fee of a block.
	// Since the block's base fee sets the minimum gas price for any transaction included in that block, this effectively sets a minimum
	// gas price for any transaction.
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

// represents an empty fee config without any field
var EmptyFeeConfig = FeeConfig{}

// Verify checks fields of this config to ensure a valid fee configuration is provided.
func (f *FeeConfig) Verify() error {
	switch {
	case f.GasLimit == nil:
		return errors.New("gasLimit cannot be nil")
	case f.MinBaseFee == nil:
		return errors.New("minBaseFee cannot be nil")
	case f.TargetGas == nil:
		return errors.New("targetGas cannot be nil")
	case f.BaseFeeChangeDenominator == nil:
		return errors.New("baseFeeChangeDenominator cannot be nil")
	case f.MinBlockGasCost == nil:
		return errors.New("minBlockGasCost cannot be nil")
	case f.MaxBlockGasCost == nil:
		return errors.New("maxBlockGasCost cannot be nil")
	case f.BlockGasCostStep == nil:
		return errors.New("blockGasCostStep cannot be nil")
	}

	switch {
	case f.GasLimit.Cmp(common.Big0) != 1:
		return fmt.Errorf("gasLimit = %d cannot be less than or equal to 0", f.GasLimit)
	case f.TargetBlockRate <= 0:
		return fmt.Errorf("targetBlockRate = %d cannot be less than or equal to 0", f.TargetBlockRate)
	case f.MinBaseFee.Cmp(common.Big0) == -1:
		return fmt.Errorf("minBaseFee = %d cannot be less than 0", f.MinBaseFee)
	case f.TargetGas.Cmp(common.Big0) != 1:
		return fmt.Errorf("targetGas = %d cannot be less than or equal to 0", f.TargetGas)
	case f.BaseFeeChangeDenominator.Cmp(common.Big0) != 1:
		return fmt.Errorf("baseFeeChangeDenominator = %d cannot be less than or equal to 0", f.BaseFeeChangeDenominator)
	case f.MinBlockGasCost.Cmp(common.Big0) == -1:
		return fmt.Errorf("minBlockGasCost = %d cannot be less than 0", f.MinBlockGasCost)
	case f.MinBlockGasCost.Cmp(f.MaxBlockGasCost) == 1:
		return fmt.Errorf("minBlockGasCost = %d cannot be greater than maxBlockGasCost = %d", f.MinBlockGasCost, f.MaxBlockGasCost)
	case f.BlockGasCostStep.Cmp(common.Big0) == -1:
		return fmt.Errorf("blockGasCostStep = %d cannot be less than 0", f.BlockGasCostStep)
	case !f.MaxBlockGasCost.IsUint64():
		return fmt.Errorf("maxBlockGasCost = %d is not a valid uint64", f.MaxBlockGasCost)
	}
	return f.checkByteLens()
}

// Equal checks if given [other] is same with this FeeConfig.
func (f *FeeConfig) Equal(other *FeeConfig) bool {
	if other == nil {
		return false
	}

	return utils.BigNumEqual(f.GasLimit, other.GasLimit) &&
		f.TargetBlockRate == other.TargetBlockRate &&
		utils.BigNumEqual(f.MinBaseFee, other.MinBaseFee) &&
		utils.BigNumEqual(f.TargetGas, other.TargetGas) &&
		utils.BigNumEqual(f.BaseFeeChangeDenominator, other.BaseFeeChangeDenominator) &&
		utils.BigNumEqual(f.MinBlockGasCost, other.MinBlockGasCost) &&
		utils.BigNumEqual(f.MaxBlockGasCost, other.MaxBlockGasCost) &&
		utils.BigNumEqual(f.BlockGasCostStep, other.BlockGasCostStep)
}

// checkByteLens checks byte lengths against common.HashLen (32 bytes) and returns error
func (f *FeeConfig) checkByteLens() error {
	if isBiggerThanHashLen(f.GasLimit) {
		return fmt.Errorf("gasLimit exceeds %d bytes", common.HashLength)
	}
	if isBiggerThanHashLen(new(big.Int).SetUint64(f.TargetBlockRate)) {
		return fmt.Errorf("targetBlockRate exceeds %d bytes", common.HashLength)
	}
	if isBiggerThanHashLen(f.MinBaseFee) {
		return fmt.Errorf("minBaseFee exceeds %d bytes", common.HashLength)
	}
	if isBiggerThanHashLen(f.TargetGas) {
		return fmt.Errorf("targetGas exceeds %d bytes", common.HashLength)
	}
	if isBiggerThanHashLen(f.BaseFeeChangeDenominator) {
		return fmt.Errorf("baseFeeChangeDenominator exceeds %d bytes", common.HashLength)
	}
	if isBiggerThanHashLen(f.MinBlockGasCost) {
		return fmt.Errorf("minBlockGasCost exceeds %d bytes", common.HashLength)
	}
	if isBiggerThanHashLen(f.MaxBlockGasCost) {
		return fmt.Errorf("maxBlockGasCost exceeds %d bytes", common.HashLength)
	}
	if isBiggerThanHashLen(f.BlockGasCostStep) {
		return fmt.Errorf("blockGasCostStep exceeds %d bytes", common.HashLength)
	}
	return nil
}

func isBiggerThanHashLen(bigint *big.Int) bool {
	buf := bigint.Bytes()
	isBigger := len(buf) > common.HashLength
	return isBigger
}
