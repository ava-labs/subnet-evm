// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package load

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/ava-labs/avalanche-network-runner/rpcpb"
	"github.com/ava-labs/subnet-evm/plugin/evm"
	"github.com/ava-labs/subnet-evm/tests/utils"
	"github.com/ava-labs/subnet-evm/tests/utils/runner"
	"github.com/ethereum/go-ethereum/log"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var (
	config   = runner.NewDefaultANRConfig()
	manager  = runner.NewNetworkManager(config)
	numNodes = 5
)

func TestE2E(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "subnet-evm small load simulator test suite")
}

// BeforeSuite starts an AvalancheGo process to use for the e2e tests
var _ = ginkgo.BeforeSuite(func() {
	// Name 10 new validators (which should have BLS key registered)
	subnetA := make([]string, 0)
	for i := 1; i <= numNodes; i++ {
		subnetA = append(subnetA, fmt.Sprintf("node%d-bls", i))
	}

	ctx := context.Background()
	var err error
	_, err = manager.StartDefaultNetwork(ctx)
	gomega.Expect(err).Should(gomega.BeNil())
	err = manager.SetupNetwork(
		ctx,
		config.AvalancheGoExecPath,
		[]*rpcpb.BlockchainSpec{
			{
				VmName:      evm.IDStr,
				Genesis:     "./tests/load/genesis/genesis.json",
				ChainConfig: "",
				SubnetSpec: &rpcpb.SubnetSpec{
					SubnetConfig: "",
					Participants: subnetA,
				},
			},
		},
	)
	gomega.Expect(err).Should(gomega.BeNil())
})

var _ = ginkgo.Describe("[Load Simulator]", ginkgo.Ordered, func() {
	ginkgo.It("basic subnet load test", ginkgo.Label("load"), func() {
		subnetIDs := manager.GetSubnets()
		gomega.Expect(len(subnetIDs)).Should(gomega.Equal(1))
		subnetID := subnetIDs[0]
		subnetDetails, ok := manager.GetSubnet(subnetID)
		gomega.Expect(ok).Should(gomega.BeTrue())
		blockchainID := subnetDetails.BlockchainID

		nodeURIs := subnetDetails.ValidatorURIs
		gomega.Expect(len(nodeURIs)).Should(gomega.Equal(numNodes))
		rpcEndpoints := make([]string, 0, len(nodeURIs))
		for _, uri := range nodeURIs {
			rpcEndpoints = append(rpcEndpoints, utils.ToRPCURI(uri, blockchainID.String()))
		}
		commaSeparatedRPCEndpoints := strings.Join(rpcEndpoints, ",")
		err := os.Setenv("RPC_ENDPOINTS", commaSeparatedRPCEndpoints)
		gomega.Expect(err).Should(gomega.BeNil())

		log.Info("Sleeping with network running", "rpcEndpoints", commaSeparatedRPCEndpoints)
		cmd := exec.Command("./scripts/run_simulator.sh")
		log.Info("Running load simulator script", "cmd", cmd.String())

		out, err := cmd.CombinedOutput()
		fmt.Printf("\nCombined output:\n\n%s\n", string(out))
		gomega.Expect(err).Should(gomega.BeNil())
	})
})

var _ = ginkgo.AfterSuite(func() {
	gomega.Expect(manager).ShouldNot(gomega.BeNil())
	gomega.Expect(manager.TeardownNetwork()).Should(gomega.BeNil())
	// TODO: bootstrap an additional node to ensure that we can bootstrap the test data correctly
})
