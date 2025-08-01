// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package params

import (
	"math/big"
	"testing"

	"github.com/ava-labs/avalanchego/upgrade/upgradetest"
	"github.com/ava-labs/subnet-evm/params/extras"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/stretchr/testify/require"
)

func TestSetEthUpgrades(t *testing.T) {
	genesisBlock := big.NewInt(0)
	genesisTimestamp := utils.NewUint64(initiallyActive)
	tests := []struct {
		fork     upgradetest.Fork
		expected *ChainConfig
	}{
		{
			fork: upgradetest.NoUpgrades,
			expected: &ChainConfig{
				HomesteadBlock:      genesisBlock,
				EIP150Block:         genesisBlock,
				EIP155Block:         genesisBlock,
				EIP158Block:         genesisBlock,
				ByzantiumBlock:      genesisBlock,
				ConstantinopleBlock: genesisBlock,
				PetersburgBlock:     genesisBlock,
				IstanbulBlock:       genesisBlock,
				MuirGlacierBlock:    genesisBlock,
				BerlinBlock:         genesisBlock,
				LondonBlock:         genesisBlock,
				ShanghaiTime:        nil,
				CancunTime:          nil,
			},
		},
		{
			fork: upgradetest.Durango,
			expected: &ChainConfig{
				HomesteadBlock:      genesisBlock,
				EIP150Block:         genesisBlock,
				EIP155Block:         genesisBlock,
				EIP158Block:         genesisBlock,
				ByzantiumBlock:      genesisBlock,
				ConstantinopleBlock: genesisBlock,
				PetersburgBlock:     genesisBlock,
				IstanbulBlock:       genesisBlock,
				MuirGlacierBlock:    genesisBlock,
				BerlinBlock:         genesisBlock,
				LondonBlock:         genesisBlock,
				ShanghaiTime:        genesisTimestamp,
				CancunTime:          nil,
			},
		},
		{
			fork: upgradetest.Etna,
			expected: &ChainConfig{
				HomesteadBlock:      genesisBlock,
				EIP150Block:         genesisBlock,
				EIP155Block:         genesisBlock,
				EIP158Block:         genesisBlock,
				ByzantiumBlock:      genesisBlock,
				ConstantinopleBlock: genesisBlock,
				PetersburgBlock:     genesisBlock,
				IstanbulBlock:       genesisBlock,
				MuirGlacierBlock:    genesisBlock,
				BerlinBlock:         genesisBlock,
				LondonBlock:         genesisBlock,
				ShanghaiTime:        genesisTimestamp,
				CancunTime:          genesisTimestamp,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.fork.String(), func(t *testing.T) {
			require := require.New(t)

			extraConfig := &extras.ChainConfig{
				NetworkUpgrades: extras.GetNetworkUpgrades(upgradetest.GetConfig(test.fork)),
			}
			actual := WithExtra(
				&ChainConfig{},
				extraConfig,
			)
			require.NoError(SetEthUpgrades(actual))

			expected := WithExtra(
				test.expected,
				extraConfig,
			)
			require.Equal(expected, actual)
		})
	}
}
