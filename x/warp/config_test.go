// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"fmt"
<<<<<<< HEAD
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/stretchr/testify/require"
)

func TestVerifyWarpconfig(t *testing.T) {
	tests := []struct {
		name          string
		config        precompileconfig.Config
		ExpectedError string
	}{
		{
			name:          "quorum numerator less than minimum",
			config:        NewConfig(big.NewInt(3), MinQuorumNumerator-1),
			ExpectedError: fmt.Sprintf("cannot specify quorum numerator (%d) < min quorum numerator (%d)", MinQuorumNumerator-1, MinQuorumNumerator),
		},
		{
			name:          "quorum numerator > quorum denominator",
			config:        NewConfig(big.NewInt(3), QuorumDenominator+1),
			ExpectedError: fmt.Sprintf("cannot specify quorum numerator (%d) > quorum denominator (%d)", QuorumDenominator+1, QuorumDenominator),
		},
		{
			name:   "default quorum numerator",
			config: NewDefaultConfig(big.NewInt(3)),
		},
		{
			name:   "valid quorum numerator 1 less than denominator",
			config: NewConfig(big.NewInt(3), QuorumDenominator-1),
		},
		{
			name:   "valid quorum numerator 1 more than minimum",
			config: NewConfig(big.NewInt(3), MinQuorumNumerator+1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)

			err := tt.config.Verify()
			if tt.ExpectedError == "" {
				require.NoError(err)
			} else {
				require.ErrorContains(err, tt.ExpectedError)
			}
		})
	}
}

func TestEqualWarpConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   precompileconfig.Config
		other    precompileconfig.Config
		expected bool
	}{
		{
			name:     "non-nil config and nil other",
			config:   NewDefaultConfig(big.NewInt(3)),
			other:    nil,
			expected: false,
		},
		{
			name:     "different type",
			config:   NewDefaultConfig(big.NewInt(3)),
			other:    precompileconfig.NewNoopStatefulPrecompileConfig(),
			expected: false,
		},
		{
			name:     "different timestamp",
			config:   NewDefaultConfig(big.NewInt(3)),
			other:    NewDefaultConfig(big.NewInt(4)),
			expected: false,
		},
		{
			name:     "different quorum numerator",
			config:   NewConfig(big.NewInt(3), MinQuorumNumerator+1),
			other:    NewConfig(big.NewInt(3), MinQuorumNumerator+2),
			expected: false,
		},
		{
			name:     "same default config",
			config:   NewDefaultConfig(big.NewInt(3)),
			other:    NewDefaultConfig(big.NewInt(3)),
			expected: true,
		},
		{
			name:     "same non-default config",
			config:   NewConfig(big.NewInt(3), MinQuorumNumerator+5),
			other:    NewConfig(big.NewInt(3), MinQuorumNumerator+5),
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)

			require.Equal(tt.expected, tt.config.Equal(tt.other))
		})
	}
=======
	"testing"

	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	"github.com/ava-labs/subnet-evm/utils"
)

func TestVerify(t *testing.T) {
	tests := map[string]testutils.ConfigVerifyTest{
		"quorum numerator less than minimum": {
			Config:        NewConfig(utils.NewUint64(3), params.WarpQuorumNumeratorMinimum-1),
			ExpectedError: fmt.Sprintf("cannot specify quorum numerator (%d) < min quorum numerator (%d)", params.WarpQuorumNumeratorMinimum-1, params.WarpQuorumNumeratorMinimum),
		},
		"quorum numerator greater than quorum denominator": {
			Config:        NewConfig(utils.NewUint64(3), params.WarpQuorumDenominator+1),
			ExpectedError: fmt.Sprintf("cannot specify quorum numerator (%d) > quorum denominator (%d)", params.WarpQuorumDenominator+1, params.WarpQuorumDenominator),
		},
		"default quorum numerator": {
			Config: NewDefaultConfig(utils.NewUint64(3)),
		},
		"valid quorum numerator 1 less than denominator": {
			Config: NewConfig(utils.NewUint64(3), params.WarpQuorumDenominator-1),
		},
		"valid quorum numerator 1 more than minimum": {
			Config: NewConfig(utils.NewUint64(3), params.WarpQuorumNumeratorMinimum+1),
		},
	}
	testutils.RunVerifyTests(t, tests)
}

func TestEqualWarpConfig(t *testing.T) {
	tests := map[string]testutils.ConfigEqualTest{
		"non-nil config and nil other": {
			Config:   NewDefaultConfig(utils.NewUint64(3)),
			Other:    nil,
			Expected: false,
		},

		"different type": {
			Config:   NewDefaultConfig(utils.NewUint64(3)),
			Other:    precompileconfig.NewNoopStatefulPrecompileConfig(),
			Expected: false,
		},

		"different timestamp": {
			Config:   NewDefaultConfig(utils.NewUint64(3)),
			Other:    NewDefaultConfig(utils.NewUint64(4)),
			Expected: false,
		},

		"different quorum numerator": {
			Config:   NewConfig(utils.NewUint64(3), params.WarpQuorumNumeratorMinimum+1),
			Other:    NewConfig(utils.NewUint64(3), params.WarpQuorumNumeratorMinimum+2),
			Expected: false,
		},

		"same default config": {
			Config:   NewDefaultConfig(utils.NewUint64(3)),
			Other:    NewDefaultConfig(utils.NewUint64(3)),
			Expected: true,
		},

		"same non-default config": {
			Config:   NewConfig(utils.NewUint64(3), params.WarpQuorumNumeratorMinimum+5),
			Other:    NewConfig(utils.NewUint64(3), params.WarpQuorumNumeratorMinimum+5),
			Expected: true,
		},
	}
	testutils.RunEqualTests(t, tests)
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
}
