// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package precompile

import (
	"context"
	"testing"
	"time"

	"github.com/ava-labs/subnet-evm/tests/utils"

	"github.com/ava-labs/avalanchego/api/health"
	"github.com/ethereum/go-ethereum/log"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	// Import the solidity package, so that ginkgo maps out the tests declared within the package
	_ "github.com/ava-labs/subnet-evm/tests/precompile/solidity"
)

var tearDown func() error

func TestE2E(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "subnet-evm precompile ginkgo test suite")
}

// BeforeSuite starts an AvalancheGo process to use for the e2e tests
var _ = ginkgo.BeforeSuite(func() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	log.Info("Starting AvalancheGo node")
	nodeUris, tearDownFunc, err := utils.SpinupAvalancheNodes(utils.DefaultNumberOfNodesToSpinUp)
	gomega.Expect(err).Should(gomega.BeNil())
	gomega.Expect(tearDownFunc).ShouldNot(gomega.BeNil())
	gomega.Expect(nodeUris).ShouldNot(gomega.BeEmpty())
	firstNodeUri := nodeUris[0]
	tearDown = tearDownFunc

	// confirm that Kurtosis started the node on the expected url
	gomega.Expect(firstNodeUri).Should(gomega.Equal(utils.DefaultLocalNodeURI))

	// Assumes that startCmd will launch a node with HTTP Port at [utils.DefaultLocalNodeURI]
	healthClient := health.NewClient(firstNodeUri)
	healthy, err := health.AwaitReady(ctx, healthClient, 5*time.Second, nil)
	gomega.Expect(err).Should(gomega.BeNil())
	gomega.Expect(healthy).Should(gomega.BeTrue())
	log.Info("AvalancheGo node is healthy")
})

var _ = ginkgo.AfterSuite(func() {
	gomega.Expect(tearDown).ShouldNot(gomega.BeNil())
	err := tearDown()
	gomega.Expect(err).Should(gomega.BeNil())
	// TODO add a new node to bootstrap off of the existing node and ensure it can bootstrap all subnets
	// created during the test
})
