// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txallowlist

import (
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestVerifyTxAllowlistConfig(t *testing.T) {
	admins := []common.Address{{1}}
	enableds := []common.Address{{2}}
	tests := []struct {
		name          string
		config        precompile.StatefulPrecompileConfig
		expectedError string
	}{
		{
			name:          "invalid allow list config in tx allowlist",
			config:        NewTxAllowListConfig(big.NewInt(3), admins, admins),
			expectedError: "cannot set address",
		},
		{
			name:          "nil member allow list config in tx allowlist",
			config:        NewTxAllowListConfig(big.NewInt(3), nil, nil),
			expectedError: "",
		},
		{
			name:          "empty member allow list config in tx allowlist",
			config:        NewTxAllowListConfig(big.NewInt(3), []common.Address{}, []common.Address{}),
			expectedError: "",
		},
		{
			name:          "valid allow list config in tx allowlist",
			config:        NewTxAllowListConfig(big.NewInt(3), admins, enableds),
			expectedError: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)

			err := tt.config.Verify()
			if tt.expectedError == "" {
				require.NoError(err)
			} else {
				require.ErrorContains(err, tt.expectedError)
			}
		})
	}
}

func TestEqualTxAllowListConfig(t *testing.T) {
	admins := []common.Address{{1}}
	enableds := []common.Address{{2}}
	tests := []struct {
		name     string
		config   precompile.StatefulPrecompileConfig
		other    precompile.StatefulPrecompileConfig
		expected bool
	}{
		{
			name:     "non-nil config and nil other",
			config:   NewTxAllowListConfig(big.NewInt(3), admins, enableds),
			other:    nil,
			expected: false,
		},
		{
			name:     "different admin",
			config:   NewTxAllowListConfig(big.NewInt(3), admins, enableds),
			other:    NewTxAllowListConfig(big.NewInt(3), []common.Address{{3}}, enableds),
			expected: false,
		},
		{
			name:     "different enabled",
			config:   NewTxAllowListConfig(big.NewInt(3), admins, enableds),
			other:    NewTxAllowListConfig(big.NewInt(3), admins, []common.Address{{3}}),
			expected: false,
		},
		{
			name:     "different timestamp",
			config:   NewTxAllowListConfig(big.NewInt(3), admins, enableds),
			other:    NewTxAllowListConfig(big.NewInt(4), admins, enableds),
			expected: false,
		},
		{
			name:     "same config",
			config:   NewTxAllowListConfig(big.NewInt(3), admins, enableds),
			other:    NewTxAllowListConfig(big.NewInt(3), admins, enableds),
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)

			require.Equal(tt.expected, tt.config.Equal(tt.other))
		})
	}
}
