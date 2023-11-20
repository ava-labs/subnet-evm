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
	message, err := NewUpgradeConfigMessage(&params.UpgradeConfig{
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

	config, err := NewUpgradeConfigMessageFromBytes(message.Bytes())
	require.NoError(t, err)

	message2, err := NewUpgradeConfigMessage(config)
	require.NoError(t, err)

	config3, err := NewUpgradeConfigMessageFromBytes(message2.Bytes())
	require.NoError(t, err)

	message3, err := NewUpgradeConfigMessage(config3)
	require.NoError(t, err)

	require.Equal(t, config, config3)
	require.Equal(t, message.hash, message2.hash)
	require.Equal(t, message2.hash, message3.hash)
}
