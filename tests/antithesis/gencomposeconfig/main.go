// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ava-labs/avalanchego/tests/antithesis"
	"github.com/ava-labs/avalanchego/tests/fixture/tmpnet"

	"github.com/ava-labs/subnet-evm/tests/utils"
)

const baseImageName = "antithesis-subnet-evm"

// Creates docker-compose.yml and its associated volumes in the target path.
func main() {
	avalancheGoPath := os.Getenv("AVALANCHEGO_PATH")
	if len(avalancheGoPath) == 0 {
		log.Fatal("AVALANCHEGO_PATH environment variable not set")
	}

	pluginDir := os.Getenv("AVALANCHEGO_PLUGIN_DIR")
	if len(pluginDir) == 0 {
		log.Fatal("AVALANCHEGO_PLUGIN_DIR environment variable not set")
	}

	targetPath := os.Getenv("TARGET_PATH")
	if len(targetPath) == 0 {
		log.Fatal("TARGET_PATH environment variable not set")
	}

	imageTag := os.Getenv("IMAGE_TAG")
	if len(imageTag) == 0 {
		log.Fatal("IMAGE_TAG environment variable not set")
	}

	nodeImageName := fmt.Sprintf("%s-node:%s", baseImageName, imageTag)
	workloadImageName := fmt.Sprintf("%s-workload:%s", baseImageName, imageTag)

	// Assume the working directory is the root of the repository
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory: %s", err)
	}
	genesisPath := filepath.Join(cwd, "tests/load/genesis/genesis.json")

	// Create a network with an xsvm subnet
	network := tmpnet.LocalNetworkOrDie()
	network.Subnets = []*tmpnet.Subnet{
		utils.NewTmpnetSubnet("subnet-evm", genesisPath, utils.DefaultChainConfig, network.Nodes...),
	}

	if err := antithesis.InitDBVolumes(network, avalancheGoPath, pluginDir, targetPath); err != nil {
		log.Fatalf("failed to initialize db volumes: %s", err)
	}

	if err := antithesis.GenerateComposeConfig(network, nodeImageName, workloadImageName, targetPath); err != nil {
		log.Fatalf("failed to generate config for docker-compose: %s", err)
	}
}
