// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nativeminter

import (
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"
)

func TestVerifyContractNativeMinterConfig(t *testing.T) {
	admins := []common.Address{{1}}
	enableds := []common.Address{{2}}
	tests := []struct {
		name          string
		config        precompile.StatefulPrecompileConfig
		expectedError string
	}{
		{
			name:          "invalid allow list config in native minter allowlist",
			config:        NewContractNativeMinterConfig(big.NewInt(3), admins, admins, nil),
			expectedError: "cannot set address",
		},
		{
			name:          "duplicate admins in config in native minter allowlist",
			config:        NewContractNativeMinterConfig(big.NewInt(3), append(admins, admins[0]), enableds, nil),
			expectedError: "duplicate address",
		},
		{
			name:          "duplicate enableds in config in native minter allowlist",
			config:        NewContractNativeMinterConfig(big.NewInt(3), admins, append(enableds, enableds[0]), nil),
			expectedError: "duplicate address",
		},
		{
			name: "nil amount in native minter config",
			config: NewContractNativeMinterConfig(big.NewInt(3), admins, nil,
				map[common.Address]*math.HexOrDecimal256{
					common.HexToAddress("0x01"): math.NewHexOrDecimal256(123),
					common.HexToAddress("0x02"): nil,
				}),
			expectedError: "initial mint cannot contain nil",
		},
		{
			name: "negative amount in native minter config",
			config: NewContractNativeMinterConfig(big.NewInt(3), admins, nil,
				map[common.Address]*math.HexOrDecimal256{
					common.HexToAddress("0x01"): math.NewHexOrDecimal256(123),
					common.HexToAddress("0x02"): math.NewHexOrDecimal256(-1),
				}),
			expectedError: "initial mint cannot contain invalid amount",
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

func TestEqualContractNativeMinterConfig(t *testing.T) {
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
			config:   NewContractNativeMinterConfig(big.NewInt(3), admins, enableds, nil),
			other:    nil,
			expected: false,
		},
		{
			name:     "different type",
			config:   NewContractNativeMinterConfig(big.NewInt(3), admins, enableds, nil),
			other:    precompile.NewNoopStatefulPrecompileConfig(),
			expected: false,
		},
		{
			name:     "different timestamps",
			config:   NewContractNativeMinterConfig(big.NewInt(3), admins, nil, nil),
			other:    NewContractNativeMinterConfig(big.NewInt(4), admins, nil, nil),
			expected: false,
		},
		{
			name:     "different enabled",
			config:   NewContractNativeMinterConfig(big.NewInt(3), admins, nil, nil),
			other:    NewContractNativeMinterConfig(big.NewInt(3), admins, enableds, nil),
			expected: false,
		},
		{
			name: "different initial mint amounts",
			config: NewContractNativeMinterConfig(big.NewInt(3), admins, nil,
				map[common.Address]*math.HexOrDecimal256{
					common.HexToAddress("0x01"): math.NewHexOrDecimal256(1),
				}),
			other: NewContractNativeMinterConfig(big.NewInt(3), admins, nil,
				map[common.Address]*math.HexOrDecimal256{
					common.HexToAddress("0x01"): math.NewHexOrDecimal256(2),
				}),
			expected: false,
		},
		{
			name: "different initial mint addresses",
			config: NewContractNativeMinterConfig(big.NewInt(3), admins, nil,
				map[common.Address]*math.HexOrDecimal256{
					common.HexToAddress("0x01"): math.NewHexOrDecimal256(1),
				}),
			other: NewContractNativeMinterConfig(big.NewInt(3), admins, nil,
				map[common.Address]*math.HexOrDecimal256{
					common.HexToAddress("0x02"): math.NewHexOrDecimal256(1),
				}),
			expected: false,
		},
		{
			name: "same config",
			config: NewContractNativeMinterConfig(big.NewInt(3), admins, nil,
				map[common.Address]*math.HexOrDecimal256{
					common.HexToAddress("0x01"): math.NewHexOrDecimal256(1),
				}),
			other: NewContractNativeMinterConfig(big.NewInt(3), admins, nil,
				map[common.Address]*math.HexOrDecimal256{
					common.HexToAddress("0x01"): math.NewHexOrDecimal256(1),
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
