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
	message, err := UpgradeConfigToNetworkMessage(&params.UpgradeConfig{
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

	config, err := ParseUpgradeConfigMessage(message.Bytes)
	require.NoError(t, err)

	message2, err := UpgradeConfigToNetworkMessage(config)
	require.NoError(t, err)

	config3, err := ParseUpgradeConfigMessage(message2.Bytes)
	require.NoError(t, err)

	require.Equal(t, config, config3)
	require.Equal(t, message.Hash, message2.Hash)
}
