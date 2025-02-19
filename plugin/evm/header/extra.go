// (c) 2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package header

import (
	"errors"
	"fmt"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
)

var errInvalidExtraLength = errors.New("invalid header.Extra length")

// ExtraPrefix takes the previous header and the timestamp of its child
// block and calculates the expected extra prefix for the child block.
func ExtraPrefix(
	config *params.ChainConfig,
	feeConfig commontype.FeeConfig,
	parent *types.Header,
	timestamp uint64,
) ([]byte, error) {
	switch {
	case config.IsSubnetEVM(timestamp):
		window, err := feeWindow(config, feeConfig, parent, timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate fee window: %w", err)
		}
		return feeWindowBytes(window), nil
	default:
		// Prior to SubnetEVM there was no expected extra prefix.
		return nil, nil
	}
}

// VerifyExtra verifies that the header's Extra field is correctly formatted for
// [rules].
func VerifyExtra(rules params.AvalancheRules, extra []byte) error {
	extraLen := len(extra)
	switch {
	case rules.IsDurango:
		if extraLen < FeeWindowSize {
			return fmt.Errorf(
				"%w: expected >= %d but got %d",
				errInvalidExtraLength,
				FeeWindowSize,
				extraLen,
			)
		}
	case rules.IsSubnetEVM:
		if extraLen != FeeWindowSize {
			return fmt.Errorf(
				"%w: expected %d but got %d",
				errInvalidExtraLength,
				FeeWindowSize,
				extraLen,
			)
		}
	default:
		if uint64(extraLen) > params.MaximumExtraDataSize {
			return fmt.Errorf(
				"%w: expected <= %d but got %d",
				errInvalidExtraLength,
				params.MaximumExtraDataSize,
				extraLen,
			)
		}
	}
	return nil
}
