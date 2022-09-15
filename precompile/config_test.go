// (c) 2022 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"
)

func TestVerifyPrecompileUpgrades(t *testing.T) {
	admins := []common.Address{{1}}
	enableds := []common.Address{{2}}
	tests := []struct {
		name          string
		config        StatefulPrecompileConfig
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
		{
			name:          "invalid allow list config in deployer allowlist",
			config:        NewContractDeployerAllowListConfig(big.NewInt(3), admins, admins),
			expectedError: "cannot set address",
		},
		{
			name:          "invalid allow list config in native minter allowlist",
			config:        NewContractNativeMinterConfig(big.NewInt(3), admins, admins, nil),
			expectedError: "cannot set address",
		},
		{
			name:          "invalid allow list config in fee manager allowlist",
			config:        NewFeeManagerConfig(big.NewInt(3), admins, admins, nil),
			expectedError: "cannot set address",
		},
		{
			name: "invalid initial fee manager config",
			config: NewFeeManagerConfig(big.NewInt(3), admins, nil,
				&commontype.FeeConfig{
					GasLimit: big.NewInt(0),
				}),
			expectedError: "gasLimit = 0 cannot be less than or equal to 0",
		},
		{
			name: "nil amount in native minter config",
			config: NewContractNativeMinterConfig(big.NewInt(3), admins, nil,
				map[common.Address]*math.HexOrDecimal256{
					common.HexToAddress("0x01"): math.NewHexOrDecimal256(123),
					common.HexToAddress("0x02"): nil,
				}),
			expectedError: ErrAmountNil.Error(),
		},
		{
			name: "negative amount in native minter config",
			config: NewContractNativeMinterConfig(big.NewInt(3), admins, nil,
				map[common.Address]*math.HexOrDecimal256{
					common.HexToAddress("0x01"): math.NewHexOrDecimal256(123),
					common.HexToAddress("0x02"): math.NewHexOrDecimal256(-1),
				}),
			expectedError: ErrAmountNonPositive.Error(),
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
