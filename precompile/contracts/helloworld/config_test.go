// (c) 2022 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package helloworld

import (
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestVerifyHelloWorldConfig(t *testing.T) {
	admins := []common.Address{{1}}
	tests := []struct {
		name          string
		config        precompileconfig.Config
		ExpectedError string
	}{
		{
			name:          "invalid allow list config in hello world allowlist",
			config:        NewConfig(big.NewInt(3), admins, admins),
			ExpectedError: "cannot set address",
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

func TestEqualHelloWorldConfig(t *testing.T) {
	admins := []common.Address{{1}}
	enableds := []common.Address{{2}}
	tests := []struct {
		name     string
		config   precompileconfig.Config
		other    precompileconfig.Config
		expected bool
	}{
		{
			name:     "non-nil config and nil other",
			config:   NewConfig(big.NewInt(3), admins, enableds),
			other:    nil,
			expected: false,
		},
		{
			name:     "different type",
			config:   NewConfig(big.NewInt(3), admins, enableds),
			other:    precompileconfig.NewNoopStatefulPrecompileConfig(),
			expected: false,
		},
		{
			name:     "different timestamp",
			config:   NewConfig(big.NewInt(3), admins, nil),
			other:    NewConfig(big.NewInt(4), admins, nil),
			expected: false,
		},
		{
			name:     "different enabled",
			config:   NewConfig(big.NewInt(3), admins, nil),
			other:    NewConfig(big.NewInt(3), admins, enableds),
			expected: false,
		},
		{
			name:     "same config",
			config:   NewConfig(big.NewInt(3), admins, nil),
			other:    NewConfig(big.NewInt(3), admins, nil),
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
