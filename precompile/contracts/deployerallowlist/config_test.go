// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package deployerallowlist

import (
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestVerifyContractDeployerConfig(t *testing.T) {
	admins := []common.Address{{1}}
	tests := []struct {
		name          string
		config        config.Config
		expectedError string
	}{
		{
			name:          "invalid allow list config in deployer allowlist",
			config:        NewContractDeployerAllowListConfig(big.NewInt(3), admins, admins),
			expectedError: "cannot set address",
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

func TestEqualContractDeployerAllowListConfig(t *testing.T) {
	admins := []common.Address{{1}}
	enableds := []common.Address{{2}}
	tests := []struct {
		name     string
		config   config.Config
		other    config.Config
		expected bool
	}{
		{
			name:     "non-nil config and nil other",
			config:   NewContractDeployerAllowListConfig(big.NewInt(3), admins, enableds),
			other:    nil,
			expected: false,
		},
		{
			name:     "different type",
			config:   NewContractDeployerAllowListConfig(big.NewInt(3), admins, enableds),
			other:    config.NewNoopStatefulPrecompileConfig(),
			expected: false,
		},
		{
			name:     "different admin",
			config:   NewContractDeployerAllowListConfig(big.NewInt(3), admins, enableds),
			other:    NewContractDeployerAllowListConfig(big.NewInt(3), []common.Address{{3}}, enableds),
			expected: false,
		},
		{
			name:     "different enabled",
			config:   NewContractDeployerAllowListConfig(big.NewInt(3), admins, enableds),
			other:    NewContractDeployerAllowListConfig(big.NewInt(3), admins, []common.Address{{3}}),
			expected: false,
		},
		{
			name:     "different timestamp",
			config:   NewContractDeployerAllowListConfig(big.NewInt(3), admins, enableds),
			other:    NewContractDeployerAllowListConfig(big.NewInt(4), admins, enableds),
			expected: false,
		},
		{
			name:     "same config",
			config:   NewContractDeployerAllowListConfig(big.NewInt(3), admins, enableds),
			other:    NewContractDeployerAllowListConfig(big.NewInt(3), admins, enableds),
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
