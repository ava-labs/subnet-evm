// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"encoding/json"
	"os"

	"github.com/ava-labs/avalanchego/config"
	"github.com/ava-labs/avalanchego/tests/fixture/tmpnet"

	"github.com/ava-labs/subnet-evm/plugin/evm"
)

func NewTmpnetNetwork(subnets ...*tmpnet.Subnet) *tmpnet.Network {
	return &tmpnet.Network{
		DefaultFlags: tmpnet.FlagsMap{
			config.ProposerVMUseCurrentHeightKey: true,
		},
		Subnets: subnets,
	}
}

// Create the configuration that will enable creation and access to a
// subnet created on a temporary network.
func NewTmpnetSubnet(name string, genesisPath string) *tmpnet.Subnet {
	genesisBytes, err := os.ReadFile(genesisPath)
	if err != nil {
		panic(err)
	}

	configBytes, err := json.Marshal(tmpnet.FlagsMap{
		"log-level":        "debug",
		"warp-api-enabled": true,
	})
	if err != nil {
		panic(err)
	}

	return &tmpnet.Subnet{
		Name: name,
		Chains: []*tmpnet.Chain{
			{
				VMID:         evm.ID,
				Genesis:      genesisBytes,
				Config:       string(configBytes),
				PreFundedKey: tmpnet.HardhatKey,
			},
		},
	}
}
