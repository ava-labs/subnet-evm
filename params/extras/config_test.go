// (c) 2025 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package extras

import (
	"math/big"
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/upgrade"
	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/precompile/contracts/txallowlist"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
 - Apricot Phase 1 Timestamp:        @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.3.0)
 - Apricot Phase 2 Timestamp:        @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.4.0)
 - Apricot Phase 3 Timestamp:        @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.5.0)
 - Apricot Phase 4 Timestamp:        @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.6.0)
 - Apricot Phase 5 Timestamp:        @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.7.0)
 - Apricot Phase P6 Timestamp:        @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.8.0)
 - Apricot Phase 6 Timestamp:        @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.8.0)
 - Apricot Phase Post-6 Timestamp:   @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.8.0
 - Banff Timestamp:                  @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.9.0)
 - Cortina Timestamp:                @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.10.0)
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
					SubnetEVMTimestamp: utils.NewUint64(uint64(1)),
					DurangoTimestamp:   utils.NewUint64(uint64(2)),
					EtnaTimestamp:      utils.NewUint64(uint64(3)),
					FortunaTimestamp:   utils.NewUint64(uint64(4)),
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
						SubnetEVMTimestamp: utils.NewUint64(uint64(13)),
					},
					StateUpgrades: []StateUpgrade{
						{
							BlockTimestamp: utils.NewUint64(uint64(14)),
							StateUpgradeAccounts: map[common.Address]StateUpgradeAccount{
								{15}: {
									Code: []byte{16},
								},
							},
						},
					},
				},
			},
			want: `Avalanche Upgrades (timestamp based):
 - Apricot Phase 1 Timestamp:        @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.3.0)
 - Apricot Phase 2 Timestamp:        @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.4.0)
 - Apricot Phase 3 Timestamp:        @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.5.0)
 - Apricot Phase 4 Timestamp:        @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.6.0)
 - Apricot Phase 5 Timestamp:        @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.7.0)
 - Apricot Phase P6 Timestamp:        @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.8.0)
 - Apricot Phase 6 Timestamp:        @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.8.0)
 - Apricot Phase Post-6 Timestamp:   @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.8.0
 - Banff Timestamp:                  @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.9.0)
 - Cortina Timestamp:                @nil        (https://github.com/ava-labs/avalanchego/releases/tag/v1.10.0)
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
			assert.Equal(t, test.want, got)
		})
	}
}

func TestChainConfigVerify(t *testing.T) {
	t.Parallel()

	validFeeConfig := commontype.FeeConfig{
		GasLimit:                 big.NewInt(1),
		TargetBlockRate:          1,
		MinBaseFee:               big.NewInt(1),
		TargetGas:                big.NewInt(1),
		BaseFeeChangeDenominator: big.NewInt(1),
		MinBlockGasCost:          big.NewInt(1),
		MaxBlockGasCost:          big.NewInt(1),
		BlockGasCostStep:         big.NewInt(1),
	}

	tests := map[string]struct {
		config   ChainConfig
		errRegex string
	}{
		"invalid_feeconfig": {
			config: ChainConfig{
				FeeConfig: commontype.FeeConfig{
					GasLimit: nil,
				},
			},
			errRegex: "^invalid fee config: ",
		},
		"invalid_precompile_upgrades": {
			// Also see precompile_config_test.go TestVerifyWithChainConfig* tests
			config: ChainConfig{
				FeeConfig: validFeeConfig,
				UpgradeConfig: UpgradeConfig{
					PrecompileUpgrades: []PrecompileUpgrade{
						// same precompile cannot be configured twice for the same timestamp
						{Config: txallowlist.NewDisableConfig(utils.NewUint64(uint64(1)))},
						{Config: txallowlist.NewDisableConfig(utils.NewUint64(uint64(1)))},
					},
				},
			},
			errRegex: "^invalid precompile upgrades: ",
		},
		"invalid_state_upgrades": {
			config: ChainConfig{
				FeeConfig: validFeeConfig,
				UpgradeConfig: UpgradeConfig{
					StateUpgrades: []StateUpgrade{
						{BlockTimestamp: nil},
					},
				},
			},
			errRegex: "^invalid state upgrades: ",
		},
		"invalid_network_upgrades": {
			config: ChainConfig{
				FeeConfig: validFeeConfig,
				NetworkUpgrades: NetworkUpgrades{
					SubnetEVMTimestamp: utils.NewUint64(1),
				},
				AvalancheContext: AvalancheContext{SnowCtx: &snow.Context{
					NetworkUpgrades: upgrade.Config{
						DurangoTime: time.Unix(2, 0),
						EtnaTime:    time.Unix(3, 0),
						FortunaTime: time.Unix(4, 0),
					},
				}},
			},
			errRegex: "^invalid network upgrades: ",
		},
		"valid": {
			config: ChainConfig{
				FeeConfig: validFeeConfig,
				NetworkUpgrades: NetworkUpgrades{
					SubnetEVMTimestamp: utils.NewUint64(uint64(0)),
					DurangoTimestamp:   utils.NewUint64(uint64(2)),
					EtnaTimestamp:      utils.NewUint64(uint64(3)),
					FortunaTimestamp:   utils.NewUint64(uint64(4)),
				},
				AvalancheContext: AvalancheContext{SnowCtx: &snow.Context{
					NetworkUpgrades: upgrade.Config{
						DurangoTime: time.Unix(2, 0),
						EtnaTime:    time.Unix(3, 0),
						FortunaTime: time.Unix(4, 0),
					},
				}},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.config.Verify()
			if test.errRegex == "" {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Regexp(t, test.errRegex, err.Error())
			}
		})
	}
}
