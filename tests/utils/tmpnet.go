// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"encoding/json"
	"os"

	"github.com/ava-labs/avalanchego/config"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/tests/fixture/tmpnet"

	"github.com/ava-labs/subnet-evm/plugin/evm"
)

func NewTmpnetNodes(count int) []*tmpnet.Node {
	nodes := make([]*tmpnet.Node, count)
	for i := range nodes {
		node := tmpnet.NewNode("")
		node.EnsureKeys()
		nodes[i] = node
	}
	return nodes
}

func NewTmpnetNetwork(nodes []*tmpnet.Node, flags tmpnet.FlagsMap, subnets ...*tmpnet.Subnet) *tmpnet.Network {
	defaultFlags := tmpnet.FlagsMap{}
	defaultFlags.SetDefaults(flags)
	defaultFlags.SetDefaults(tmpnet.FlagsMap{
		// Remove when vendored tmpnet default is `off`. tmpnet nodes are run headless so stdout logging is unnecessary.
		config.LogDisplayLevelKey:            "off",
		config.ProposerVMUseCurrentHeightKey: true,
	})
	return &tmpnet.Network{
		DefaultFlags: defaultFlags,
		Nodes:        nodes,
		Subnets:      subnets,
	}
}

// Create the configuration that will enable creation and access to a
// subnet created on a temporary network.
func NewTmpnetSubnet(name string, genesisPath string, chainConfig tmpnet.FlagsMap, nodes ...*tmpnet.Node) *tmpnet.Subnet {
	if len(nodes) == 0 {
		panic("a subnet must be validated by at least one node")
	}

	validatorIDs := make([]ids.NodeID, len(nodes))
	for i, node := range nodes {
		validatorIDs[i] = node.NodeID
	}

	genesisBytes, err := os.ReadFile(genesisPath)
	if err != nil {
		panic(err)
	}

	chainConfigBytes, err := json.Marshal(chainConfig)
	if err != nil {
		panic(err)
	}

	return &tmpnet.Subnet{
		Name: name,
		Chains: []*tmpnet.Chain{
			{
				VMID:         evm.ID,
				Genesis:      genesisBytes,
				Config:       string(chainConfigBytes),
				PreFundedKey: tmpnet.HardhatKey,
			},
		},
		ValidatorIDs: validatorIDs,
	}
}
