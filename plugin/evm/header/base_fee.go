// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package header

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ava-labs/libevm/core/types"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/params/extras"
)

var errEstimateBaseFeeWithoutActivation = errors.New("cannot estimate base fee for chain without activation scheduled")

// BaseFee takes the previous header and the timestamp of its child block and
// calculates the expected base fee for the child block.
//
// Prior to SubnetEVM, the returned base fee will be nil.
func BaseFee(
	config *extras.ChainConfig,
	feeConfig commontype.FeeConfig,
	acp224FeeConfig commontype.ACP224FeeConfig,
	parent *types.Header,
	timestamp uint64,
) (*big.Int, error) {
	// TODO: XXX Handle feeConfig with Fortuna here
	switch {
	case config.IsFortuna(timestamp):
		state, err := feeStateBeforeBlock(config, parent, timestamp)
		if err != nil {
			return nil, fmt.Errorf("calculating initial fee state: %w", err)
		}
		price := state.GasPrice()
		return new(big.Int).SetUint64(uint64(price)), nil
	case config.IsSubnetEVM(timestamp):
		return baseFeeFromWindow(config, feeConfig, parent, timestamp)
	default:
		// Prior to SubnetEVM the expected base fee is nil.
		return nil, nil
	}
}

// EstimateNextBaseFee attempts to estimate the base fee of a block built at
// `timestamp` on top of `parent`.
//
// If timestamp is before parent.Time or the SubnetEVM activation time, then timestamp
// is set to the maximum of parent.Time and the SubnetEVM activation time.
//
// Warning: This function should only be used in estimation and should not be
// used when calculating the canonical base fee for a block.
func EstimateNextBaseFee(
	config *extras.ChainConfig,
	feeConfig commontype.FeeConfig,
	acp224FeeConfig commontype.ACP224FeeConfig,
	parent *types.Header,
	timestamp uint64,
) (*big.Int, error) {
	if config.SubnetEVMTimestamp == nil {
		return nil, errEstimateBaseFeeWithoutActivation
	}

	timestamp = max(timestamp, parent.Time, *config.SubnetEVMTimestamp)
	return BaseFee(config, feeConfig, acp224FeeConfig, parent, timestamp)
}
