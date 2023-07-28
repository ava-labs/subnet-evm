package handshake

import (
	"testing"

	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contracts/nativeminter"
	"github.com/stretchr/testify/assert"
)

func TestSerialize(t *testing.T) {
	var t0 uint64 = 0
	var t1 uint64 = 1
	config, err := NewUpgradeConfig(params.UpgradeConfig{
		PrecompileUpgrades: []params.PrecompileUpgrade{
			{
				Config: nativeminter.NewConfig(&t0, nil, nil, nil), // enable at genesis
			},
			{
				Config: nativeminter.NewDisableConfig(&t1), // disable at timestamp 1
			},
		},
	})
	assert.NoError(t, err)

	config2, err := ParseUpgradeConfig(config.Bytes())
	assert.NoError(t, err)

	config3, err := NewUpgradeConfig(config2.Config())
	assert.NoError(t, err)

	assert.Equal(t, config2, config3)
	assert.Equal(t, config.Hash(), config2.Hash())
	assert.Equal(t, config.Hash(), config3.Hash())
}
