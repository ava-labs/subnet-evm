// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestApplyUpgradeBytes(t *testing.T) {
	admins := []common.Address{{1}}
	chainConfig := *TestChainConfig
	chainConfig.TxAllowListConfig = precompile.NewTxAllowListConfig(big.NewInt(1), admins)
	chainConfig.ContractDeployerAllowListConfig = precompile.NewContractDeployerAllowListConfig(big.NewInt(10), admins)

	type test struct {
		configs             []*UpgradeConfig
		startTimestamps     []*big.Int
		expectedErrorString string
	}

	tests := map[string]test{
		"upgrade bytes conflicts with genesis (re-enable without disable)": {
			expectedErrorString: "disable should be [true]",
			startTimestamps:     []*big.Int{big.NewInt(5)},
			configs: []*UpgradeConfig{
				{
					PrecompileUpgrades: []PrecompileUpgrade{
						{
							TxAllowListConfig: precompile.NewTxAllowListConfig(big.NewInt(1), admins),
						},
					},
				},
			},
		},
		"upgrade bytes conflicts with genesis (disable before enable)": {
			expectedErrorString: "timestamp should not be less than [1]",
			startTimestamps:     []*big.Int{big.NewInt(5)},
			configs: []*UpgradeConfig{
				{
					PrecompileUpgrades: []PrecompileUpgrade{
						{
							TxAllowListConfig: precompile.NewDisableTxAllowListConfig(big.NewInt(0)),
						},
					},
				},
			},
		},
		"disable and re-enable": {
			startTimestamps: []*big.Int{big.NewInt(5)},
			configs: []*UpgradeConfig{
				{
					PrecompileUpgrades: []PrecompileUpgrade{
						{
							TxAllowListConfig: precompile.NewDisableTxAllowListConfig(big.NewInt(6)),
						},
						{
							TxAllowListConfig: precompile.NewTxAllowListConfig(big.NewInt(7), admins),
						},
					},
				},
			},
		},
		"disable and re-enable, reschedule upgrade before it happens": {
			startTimestamps: []*big.Int{big.NewInt(5), big.NewInt(6)},
			configs: []*UpgradeConfig{
				{
					PrecompileUpgrades: []PrecompileUpgrade{
						{
							TxAllowListConfig: precompile.NewDisableTxAllowListConfig(big.NewInt(6)),
						},
						{
							TxAllowListConfig: precompile.NewTxAllowListConfig(big.NewInt(7), admins),
						},
					},
				},
				{
					PrecompileUpgrades: []PrecompileUpgrade{
						{
							TxAllowListConfig: precompile.NewDisableTxAllowListConfig(big.NewInt(6)),
						},
						{
							TxAllowListConfig: precompile.NewTxAllowListConfig(big.NewInt(8), admins),
						},
					},
				},
			},
		},
		"disable and re-enable, reschedule upgrade after it happens": {
			expectedErrorString: "mismatching PrecompileUpgrade",
			startTimestamps:     []*big.Int{big.NewInt(5), big.NewInt(8)},
			configs: []*UpgradeConfig{
				{
					PrecompileUpgrades: []PrecompileUpgrade{
						{
							TxAllowListConfig: precompile.NewDisableTxAllowListConfig(big.NewInt(6)),
						},
						{
							TxAllowListConfig: precompile.NewTxAllowListConfig(big.NewInt(7), admins),
						},
					},
				},
				{
					PrecompileUpgrades: []PrecompileUpgrade{
						{
							TxAllowListConfig: precompile.NewDisableTxAllowListConfig(big.NewInt(6)),
						},
						{
							TxAllowListConfig: precompile.NewTxAllowListConfig(big.NewInt(8), admins),
						},
					},
				},
			},
		},
		"disable and re-enable, cancel upgrade before it happens": {
			startTimestamps: []*big.Int{big.NewInt(5), big.NewInt(6)},
			configs: []*UpgradeConfig{
				{
					PrecompileUpgrades: []PrecompileUpgrade{
						{
							TxAllowListConfig: precompile.NewDisableTxAllowListConfig(big.NewInt(6)),
						},
						{
							TxAllowListConfig: precompile.NewTxAllowListConfig(big.NewInt(7), admins),
						},
					},
				},
				{
					PrecompileUpgrades: []PrecompileUpgrade{
						{
							TxAllowListConfig: precompile.NewDisableTxAllowListConfig(big.NewInt(6)),
						},
					},
				},
			},
		},
		"disable and re-enable, cancel upgrade after it happens": {
			expectedErrorString: "mismatching missing PrecompileUpgrade",
			startTimestamps:     []*big.Int{big.NewInt(5), big.NewInt(8)},
			configs: []*UpgradeConfig{
				{
					PrecompileUpgrades: []PrecompileUpgrade{
						{
							TxAllowListConfig: &precompile.TxAllowListConfig{
								UpgradeableConfig: precompile.UpgradeableConfig{
									BlockTimestamp: big.NewInt(6),
									Disable:        true,
								},
							},
						},
						{
							TxAllowListConfig: &precompile.TxAllowListConfig{
								UpgradeableConfig: precompile.UpgradeableConfig{
									BlockTimestamp: big.NewInt(7),
								},
								AllowListConfig: precompile.AllowListConfig{
									AllowListAdmins: admins,
								},
							},
						},
					},
				},
				{
					PrecompileUpgrades: []PrecompileUpgrade{
						{
							TxAllowListConfig: &precompile.TxAllowListConfig{
								UpgradeableConfig: precompile.UpgradeableConfig{
									BlockTimestamp: big.NewInt(6),
									Disable:        true,
								},
							},
						},
					},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// make a local copy of the chainConfig
			chainConfig := chainConfig

			// apply all the upgrade bytes specified in order
			for i, upgrade := range tt.configs {
				newCfg := chainConfig
				newCfg.UpgradeConfig = *upgrade

				// TODO: split tests for verify vs. checkCompatible
				err := newCfg.Verify()
				if err == nil {
					if compatErr := chainConfig.checkCompatible(&newCfg, nil, tt.startTimestamps[i]); compatErr != nil {
						err = compatErr
					}
				}

				// if this is not the final upgradeBytes, continue applying
				// the next upgradeBytes. (only check the result on the last apply)
				if i != len(tt.configs)-1 {
					if err != nil {
						t.Fatalf("expecting ApplyUpgradeBytes call %d to return nil, got %s", i+1, err)
					}
					chainConfig = newCfg
					continue
				}

				if tt.expectedErrorString != "" {
					assert.ErrorContains(t, err, tt.expectedErrorString)
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}
