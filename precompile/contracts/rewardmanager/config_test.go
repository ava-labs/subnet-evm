// (c) 2022 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package rewardmanager

import (
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/precompile/allowlist"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestVerifyRewardManagerConfig(t *testing.T) {
	admins := []common.Address{{1}}
	enableds := []common.Address{{2}}
	tests := map[string]testutils.ConfigVerifyTest{
		"both reward mechanisms should not be activated at the same time in reward manager": {
			Config: NewConfig(big.NewInt(3), admins, enableds, &InitialRewardConfig{
				AllowFeeRecipients: true,
				RewardAddress:      common.HexToAddress("0x01"),
			}),
			ExpectedError: ErrCannotEnableBothRewards.Error(),
		},
	}
	allowlist.VerifyPrecompileWithAllowListTests(t, Module, tests)
}

func TestEqualRewardManagerConfig(t *testing.T) {
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
			config:   NewConfig(big.NewInt(3), admins, enableds, nil),
			other:    nil,
			expected: false,
		},
		{
			name:     "different type",
			config:   NewConfig(big.NewInt(3), admins, enableds, nil),
			other:    precompileconfig.NewNoopStatefulPrecompileConfig(),
			expected: false,
		},
		{
			name:     "different timestamp",
			config:   NewConfig(big.NewInt(3), admins, nil, nil),
			other:    NewConfig(big.NewInt(4), admins, nil, nil),
			expected: false,
		},
		{
			name:     "different enabled",
			config:   NewConfig(big.NewInt(3), admins, nil, nil),
			other:    NewConfig(big.NewInt(3), admins, enableds, nil),
			expected: false,
		},
		{
			name: "non-nil initial config and nil initial config",
			config: NewConfig(big.NewInt(3), admins, nil, &InitialRewardConfig{
				AllowFeeRecipients: true,
			}),
			other:    NewConfig(big.NewInt(3), admins, nil, nil),
			expected: false,
		},
		{
			name: "different initial config",
			config: NewConfig(big.NewInt(3), admins, nil, &InitialRewardConfig{
				RewardAddress: common.HexToAddress("0x01"),
			}),
			other: NewConfig(big.NewInt(3), admins, nil,
				&InitialRewardConfig{
					RewardAddress: common.HexToAddress("0x02"),
				}),
			expected: false,
		},
		{
			name: "same config",
			config: NewConfig(big.NewInt(3), admins, nil, &InitialRewardConfig{
				RewardAddress: common.HexToAddress("0x01"),
			}),
			other: NewConfig(big.NewInt(3), admins, nil, &InitialRewardConfig{
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
