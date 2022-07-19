// (c) 2022 Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestValidateWithChainConfig(t *testing.T) {
	admins := []common.Address{{1}}
	config := &UpgradesConfig{
		Upgrade: Upgrade{
			TxAllowListConfig: &TxAllowListConfig{
				UpgradeableConfig: UpgradeableConfig{
					BlockTimestamp: big.NewInt(2),
				},
			},
		},
	}
	config.DisableTxAllowListUpgrade(big.NewInt(4))
	config.AddTxAllowListUpgrade(big.NewInt(5), admins)

	// check this config is valid
	err := config.Validate()
	assert.NoError(t, err)

	// entries must be monotonically increasing
	badConfig := *config
	badConfig.DisableTxAllowListUpgrade(big.NewInt(1))
	err = badConfig.Validate()
	assert.ErrorContains(t, err, "timestamp should not be less than [5]")

	// cannot enable a precompile without disabling it first.
	badConfig = *config
	badConfig.AddTxAllowListUpgrade(big.NewInt(5), admins)
	err = badConfig.Validate()
	assert.ErrorContains(t, err, "disable should be [true]")
}

func TestValidate(t *testing.T) {
	admins := []common.Address{common.Address{1}}
	config := &UpgradesConfig{}
	config.AddTxAllowListUpgrade(big.NewInt(1), admins)
	config.DisableTxAllowListUpgrade(big.NewInt(2))

	// check this config is valid
	err := config.Validate()
	assert.NoError(t, err)
}
