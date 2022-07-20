// (c) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ava-labs/subnet-evm/precompile"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestApplyUpgradeBytes(t *testing.T) {
	admins := []common.Address{{1}}
	chainConfig := &ChainConfig{
		UpgradesConfig: UpgradesConfig{
			Upgrade: Upgrade{
				TxAllowListConfig: &precompile.TxAllowListConfig{
					UpgradeableConfig: precompile.UpgradeableConfig{
						BlockTimestamp: big.NewInt(1),
					},
					AllowListConfig: precompile.AllowListConfig{
						AllowListAdmins: admins,
					},
				},

				ContractDeployerAllowListConfig: &precompile.ContractDeployerAllowListConfig{
					UpgradeableConfig: precompile.UpgradeableConfig{
						BlockTimestamp: big.NewInt(10),
					},
					AllowListConfig: precompile.AllowListConfig{
						AllowListAdmins: admins,
					},
				},
			},
		},
	}

	type test struct {
		configs             []*UpgradeBytesConfig
		startTimestamps     []*big.Int
		expectedErrorString string
	}

	tests := map[string]test{
		"upgrade bytes conflicts with genesis (re-enable without disable)": {
			expectedErrorString: "disable should be [true]",
			startTimestamps:     []*big.Int{big.NewInt(5)},
			configs: []*UpgradeBytesConfig{
				{
					PrecompileUpgrades: []Upgrade{
						{
							TxAllowListConfig: &precompile.TxAllowListConfig{
								UpgradeableConfig: precompile.UpgradeableConfig{
									BlockTimestamp: big.NewInt(1),
								},
								AllowListConfig: precompile.AllowListConfig{
									AllowListAdmins: admins,
								},
							},
						},
					},
				},
			},
		},
		"upgrade bytes conflicts with genesis (disable before enable)": {
			expectedErrorString: "timestamp should not be less than [1]",
			startTimestamps:     []*big.Int{big.NewInt(5)},
			configs: []*UpgradeBytesConfig{
				{
					PrecompileUpgrades: []Upgrade{
						{
							TxAllowListConfig: &precompile.TxAllowListConfig{
								UpgradeableConfig: precompile.UpgradeableConfig{
									BlockTimestamp: big.NewInt(0),
									Disable:        true,
								},
							},
						},
					},
				},
			},
		},
		"disable and re-enable": {
			startTimestamps: []*big.Int{big.NewInt(5)},
			configs: []*UpgradeBytesConfig{
				{
					PrecompileUpgrades: []Upgrade{
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
			},
		},
		"disable and re-enable, reschedule upgrade before it happens": {
			startTimestamps: []*big.Int{big.NewInt(5), big.NewInt(6)},
			configs: []*UpgradeBytesConfig{
				{
					PrecompileUpgrades: []Upgrade{
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
					PrecompileUpgrades: []Upgrade{
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
									BlockTimestamp: big.NewInt(8),
								},
								AllowListConfig: precompile.AllowListConfig{
									AllowListAdmins: admins,
								},
							},
						},
					},
				},
			},
		},
		"disable and re-enable, reschedule upgrade after it happens": {
			expectedErrorString: "mismatching PrecompileUpgrade",
			startTimestamps:     []*big.Int{big.NewInt(5), big.NewInt(8)},
			configs: []*UpgradeBytesConfig{
				{
					PrecompileUpgrades: []Upgrade{
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
					PrecompileUpgrades: []Upgrade{
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
									BlockTimestamp: big.NewInt(8),
								},
								AllowListConfig: precompile.AllowListConfig{
									AllowListAdmins: admins,
								},
							},
						},
					},
				},
			},
		},
		"disable and re-enable, cancel upgrade before it happens": {
			startTimestamps: []*big.Int{big.NewInt(5), big.NewInt(6)},
			configs: []*UpgradeBytesConfig{
				{
					PrecompileUpgrades: []Upgrade{
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
					PrecompileUpgrades: []Upgrade{
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
		"disable and re-enable, cancel upgrade after it happens": {
			expectedErrorString: "mismatching missing PrecompileUpgrade",
			startTimestamps:     []*big.Int{big.NewInt(5), big.NewInt(8)},
			configs: []*UpgradeBytesConfig{
				{
					PrecompileUpgrades: []Upgrade{
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
					PrecompileUpgrades: []Upgrade{
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
			chainConfig := *chainConfig

			// apply all the upgrade bytes specified in order
			for i, upgrade := range tt.configs {
				upgradeBytes, err := json.Marshal(upgrade)
				if err != nil {
					t.Fatal(err)
				}

				err = chainConfig.ApplyUpgradeBytes(upgradeBytes, tt.startTimestamps[i])
				// if this is not the final upgradeBytes, continue applying
				// the next upgradeBytes. (only check the result on the last apply)
				if i != len(tt.configs)-1 {
					if err != nil {
						t.Fatalf("expecting ApplyUpgradeBytes call %d to return nil, got %s", i+1, err)
					}
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
