// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package load

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	ginkgo "github.com/onsi/ginkgo/v2"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/log"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/tests/fixture/e2e"
	"github.com/ava-labs/avalanchego/tests/fixture/tmpnet"
	"github.com/ava-labs/avalanchego/utils/set"

	"github.com/ethereum/go-ethereum/tests"
	"github.com/ethereum/go-ethereum/tests/utils"
)

const (
	// The load test requires 5 nodes
	nodeCount = 5

	subnetAName = "load-subnet-a"
)

var (
	flagVars     *e2e.FlagVars
	repoRootPath = tests.GetRepoRootPath("tests/load")
)

func init() {
	// Configures flags used to configure tmpnet
	flagVars = e2e.RegisterFlags()
}

func TestE2E(t *testing.T) {
	ginkgo.RunSpecs(t, "subnet-evm small load simulator test suite")
}

var _ = ginkgo.Describe("[Load Simulator]", ginkgo.Ordered, func() {
	require := require.New(ginkgo.GinkgoT())

	var env *e2e.TestEnvironment

	ginkgo.BeforeAll(func() {
		genesisPath := filepath.Join(repoRootPath, "tests/load/genesis/genesis.json")

		nodes := utils.NewTmpnetNodes(nodeCount)
		env = e2e.NewTestEnvironment(
			flagVars,
			utils.NewTmpnetNetwork(
				"subnet-evm-small-load",
				nodes,
				tmpnet.FlagsMap{},
				utils.NewTmpnetSubnet(subnetAName, genesisPath, utils.DefaultChainConfig, nodes...),
			),
		)
	})

	ginkgo.It("basic subnet load test", ginkgo.Label("load"), func() {
		network := env.GetNetwork()

		subnet := network.GetSubnet(subnetAName)
		require.NotNil(subnet)
		blockchainID := subnet.Chains[0].ChainID

		nodeURIs := tmpnet.GetNodeURIs(network.Nodes)
		validatorIDs := set.NewSet[ids.NodeID](len(subnet.ValidatorIDs))
		validatorIDs.Add(subnet.ValidatorIDs...)
		rpcEndpoints := make([]string, 0, len(nodeURIs))
		for _, nodeURI := range nodeURIs {
			if !validatorIDs.Contains(nodeURI.NodeID) {
				continue
			}
			rpcEndpoints = append(rpcEndpoints, fmt.Sprintf("%s/ext/bc/%s/rpc", nodeURI.URI, blockchainID))
		}
		commaSeparatedRPCEndpoints := strings.Join(rpcEndpoints, ",")
		err := os.Setenv("RPC_ENDPOINTS", commaSeparatedRPCEndpoints)
		require.NoError(err)

		log.Info("Running load simulator...", "rpcEndpoints", commaSeparatedRPCEndpoints)
		cmd := exec.Command("./scripts/run_simulator.sh")
		cmd.Dir = repoRootPath
		log.Info("Running load simulator script", "cmd", cmd.String())

		out, err := cmd.CombinedOutput()
		fmt.Printf("\nCombined output:\n\n%s\n", string(out))
		require.NoError(err)
	})
})
