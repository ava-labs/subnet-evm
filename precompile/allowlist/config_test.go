// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package allowlist

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestVerifyAllowlistConfig(t *testing.T) {
	admins := []common.Address{{1}}
	enableds := []common.Address{{2}}
	tests := []struct {
		name          string
		config        Config
		expectedError string
	}{
		{
			name:          "invalid allow list config in allowlist",
			config:        Config{admins, admins},
			expectedError: "cannot set address",
		},
		{
			name:          "nil member allow list config in allowlist",
			config:        Config{nil, nil},
			expectedError: "",
		},
		{
			name:          "empty member allow list config in allowlist",
			config:        Config{[]common.Address{}, []common.Address{}},
			expectedError: "",
		},
		{
			name:          "valid allow list config in allowlist",
			config:        Config{admins, enableds},
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

func TestEqualAllowListConfig(t *testing.T) {
	admins := []common.Address{{1}}
	enableds := []common.Address{{2}}
	tests := []struct {
		name     string
		config   *Config
		other    *Config
		expected bool
	}{
		{
			name:     "non-nil config and nil other",
			config:   &Config{admins, enableds},
			other:    nil,
			expected: false,
		},
		{
			name:     "different admin",
			config:   &Config{admins, enableds},
			other:    &Config{[]common.Address{{3}}, enableds},
			expected: false,
		},
		{
			name:     "different enabled",
			config:   &Config{admins, enableds},
			other:    &Config{admins, []common.Address{{3}}},
			expected: false,
		},
		{
			name:     "same config",
			config:   &Config{admins, enableds},
			other:    &Config{admins, enableds},
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
