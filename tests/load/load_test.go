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
	"time"

	"github.com/ava-labs/avalanchego/api/health"
	"github.com/ava-labs/subnet-evm/tests/utils"
	"github.com/ethereum/go-ethereum/log"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var tearDown func() error

func TestE2E(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "subnet-evm small load simulator test suite")
}

// BeforeSuite starts an AvalancheGo process to use for the e2e tests
var _ = ginkgo.BeforeSuite(func() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	log.Info("Starting AvalancheGo node")
	nodeUris, tearDownFunc, err := utils.SpinupAvalancheNodes(utils.NumNodesToSpinUpForMultiNodeTest)
	gomega.Expect(err).Should(gomega.BeNil())
	gomega.Expect(tearDownFunc).ShouldNot(gomega.BeNil())
	tearDown = tearDownFunc
	gomega.Expect(nodeUris).Should(gomega.HaveLen(utils.NumNodesToSpinUpForMultiNodeTest))
	for index, nodeUri := range nodeUris {
		// confirm that Kurtosis started the node on the expected url
		gomega.Expect(nodeUri).Should(gomega.Equal(utils.NodeURIs[index]))
		healthClient := health.NewClient(nodeUri)
		healthy, err := health.AwaitReady(ctx, healthClient, 5*time.Second, nil)
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(healthy).Should(gomega.BeTrue())
		log.Info("AvalancheGo node is healthy")
	}
})

var _ = ginkgo.Describe("[Load Simulator]", ginkgo.Ordered, func() {
	ginkgo.It("basic subnet load test", ginkgo.Label("load"), func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		blockchainID := utils.CreateNewSubnet(ctx, "./tests/load/genesis/genesis.json")

		rpcEndpoints := make([]string, 0, len(utils.NodeURIs))
		for _, uri := range utils.NodeURIs {
			rpcEndpoints = append(rpcEndpoints, fmt.Sprintf("%s/ext/bc/%s/rpc", uri, blockchainID))
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
	gomega.Expect(tearDown).ShouldNot(gomega.BeNil())
	err := tearDown()
	gomega.Expect(err).Should(gomega.BeNil())
	// TODO add a new node to bootstrap off of the existing node and ensure it can bootstrap all subnets
	// created during the test
})
