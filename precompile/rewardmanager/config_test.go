// (c) 2022 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package rewardmanager

import (
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/precompile/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestVerifyRewardManagerConfig(t *testing.T) {
	admins := []common.Address{{1}}
	enableds := []common.Address{{2}}
	tests := []struct {
		name          string
		config        config.Config
		expectedError string
	}{
		{
			name:          "duplicate enableds in config in reward manager allowlist",
			config:        NewRewardManagerConfig(big.NewInt(3), admins, append(enableds, enableds[0]), nil),
			expectedError: "duplicate address",
		},
		{
			name: "both reward mechanisms should not be activated at the same time in reward manager",
			config: NewRewardManagerConfig(big.NewInt(3), admins, enableds, &InitialRewardConfig{
				AllowFeeRecipients: true,
				RewardAddress:      common.HexToAddress("0x01"),
			}),
			expectedError: ErrCannotEnableBothRewards.Error(),
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

func TestEqualRewardManagerConfig(t *testing.T) {
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
			config:   NewRewardManagerConfig(big.NewInt(3), admins, enableds, nil),
			other:    nil,
			expected: false,
		},
		{
			name:     "different type",
			config:   NewRewardManagerConfig(big.NewInt(3), admins, enableds, nil),
			other:    config.NewNoopStatefulPrecompileConfig(),
			expected: false,
		},
		{
			name:     "different timestamp",
			config:   NewRewardManagerConfig(big.NewInt(3), admins, nil, nil),
			other:    NewRewardManagerConfig(big.NewInt(4), admins, nil, nil),
			expected: false,
		},
		{
			name:     "different enabled",
			config:   NewRewardManagerConfig(big.NewInt(3), admins, nil, nil),
			other:    NewRewardManagerConfig(big.NewInt(3), admins, enableds, nil),
			expected: false,
		},
		{
			name: "non-nil initial config and nil initial config",
			config: NewRewardManagerConfig(big.NewInt(3), admins, nil, &InitialRewardConfig{
				AllowFeeRecipients: true,
			}),
			other:    NewRewardManagerConfig(big.NewInt(3), admins, nil, nil),
			expected: false,
		},
		{
			name: "different initial config",
			config: NewRewardManagerConfig(big.NewInt(3), admins, nil, &InitialRewardConfig{
				RewardAddress: common.HexToAddress("0x01"),
			}),
			other: NewRewardManagerConfig(big.NewInt(3), admins, nil,
				&InitialRewardConfig{
					RewardAddress: common.HexToAddress("0x02"),
				}),
			expected: false,
		},
		{
			name: "same config",
			config: NewRewardManagerConfig(big.NewInt(3), admins, nil, &InitialRewardConfig{
				RewardAddress: common.HexToAddress("0x01"),
			}),
			other: NewRewardManagerConfig(big.NewInt(3), admins, nil, &InitialRewardConfig{
				RewardAddress: common.HexToAddress("0x01"),
			}),
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
