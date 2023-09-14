// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package handshake

import (
	"testing"

	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contracts/nativeminter"
	"github.com/stretchr/testify/require"
)

func TestSerialize(t *testing.T) {
	var t0 uint64 = 0
	var t1 uint64 = 1
	config, err := NewUpgradeConfig(params.UpgradeConfig{
		PrecompileUpgrades: []params.PrecompileUpgrade{
			{
				Config: nativeminter.NewConfig(&t0, nil, nil, nil, nil), // enable at genesis
			},
			{
				Config: nativeminter.NewDisableConfig(&t1), // disable at timestamp 1
			},
		},
	})
	require.NoError(t, err)

	config2, err := ParseUpgradeConfig(config.Bytes())
	require.NoError(t, err)

	config3, err := NewUpgradeConfig(config2.Config())
	require.NoError(t, err)

	require.Equal(t, config2, config3)
	require.Equal(t, config.Hash(), config2.Hash())
	require.Equal(t, config.Hash(), config3.Hash())
}
