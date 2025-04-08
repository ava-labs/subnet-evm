// (c) 2025 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package extras

import (
	"math/big"
	"testing"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/stretchr/testify/assert"
)

func pointer[T any](v T) *T { return &v }

func TestChainConfigDescription(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		config *ChainConfig
		want   string
	}{
		"nil": {},
		"empty": {
			config: &ChainConfig{},
			want: `Avalanche Upgrades (timestamp based):
 - SubnetEVM Timestamp:          @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.10.0)
 - Durango Timestamp:            @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.11.0)
 - Etna Timestamp:               @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.12.0)
 - Fortuna Timestamp:            @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.13.0)

Upgrade Config: {}
Fee Config: {}
Allow Fee Recipients: false
`,
		},
		"set": {
			config: &ChainConfig{
				NetworkUpgrades: NetworkUpgrades{
					SubnetEVMTimestamp: pointer(uint64(1)),
					DurangoTimestamp:   pointer(uint64(2)),
					EtnaTimestamp:      pointer(uint64(3)),
					FortunaTimestamp:   pointer(uint64(4)),
				},
				FeeConfig: commontype.FeeConfig{
					GasLimit:                 big.NewInt(5),
					TargetBlockRate:          6,
					MinBaseFee:               big.NewInt(7),
					TargetGas:                big.NewInt(8),
					BaseFeeChangeDenominator: big.NewInt(9),
					MinBlockGasCost:          big.NewInt(10),
					MaxBlockGasCost:          big.NewInt(11),
					BlockGasCostStep:         big.NewInt(12),
				},
				AllowFeeRecipients: true,
				UpgradeConfig: UpgradeConfig{
					NetworkUpgradeOverrides: &NetworkUpgrades{
						SubnetEVMTimestamp: pointer(uint64(13)),
					},
					StateUpgrades: []StateUpgrade{
						{
							BlockTimestamp: pointer(uint64(14)),
							StateUpgradeAccounts: map[common.Address]StateUpgradeAccount{
								common.Address{15}: {
									Code: []byte{16},
								},
							},
						},
					},
				},
			},
			want: `Avalanche Upgrades (timestamp based):
 - SubnetEVM Timestamp:          @1          (https://github.com/ava-labs/avalanchego/releases/tag/v1.10.0)
 - Durango Timestamp:            @2          (https://github.com/ava-labs/avalanchego/releases/tag/v1.11.0)
 - Etna Timestamp:               @3          (https://github.com/ava-labs/avalanchego/releases/tag/v1.12.0)
 - Fortuna Timestamp:            @4          (https://github.com/ava-labs/avalanchego/releases/tag/v1.13.0)

Upgrade Config: {"networkUpgradeOverrides":{"subnetEVMTimestamp":13},"stateUpgrades":[{"blockTimestamp":14,"accounts":{"0x0f00000000000000000000000000000000000000":{"code":"0x10"}}}]}
Fee Config: {"gasLimit":5,"targetBlockRate":6,"minBaseFee":7,"targetGas":8,"baseFeeChangeDenominator":9,"minBlockGasCost":10,"maxBlockGasCost":11,"blockGasCostStep":12}
Allow Fee Recipients: true
`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := test.config.Description()
			assert.Equal(t, test.want, got, "config description mismatch")
		})
	}
}
