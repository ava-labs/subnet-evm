// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// e2e implements the e2e tests.
package e2e

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/api/health"
	"github.com/ava-labs/avalanchego/config"
	"github.com/ava-labs/subnet-evm/tests/e2e/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/go-cmd/cmd"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	_ "github.com/ava-labs/subnet-evm/tests/e2e/solidity"
)

var (
	// sets the "avalanchego" exec path
	avalanchegoExecPath string
	dataDir             string
	configFilePath      string

	setupTimeout time.Duration

	startCmd *cmd.Cmd

	defaultConfigJSON = `{
		"network-id": "local",
		"staking-enabled": false
	  }`
)

func init() {
	// Assumes that the plugin directory will be found in a default location, so we do not set it here.
	flag.StringVar(
		&avalanchegoExecPath,
		"avalanchego-path",
		os.ExpandEnv("$GOPATH/src/github.com/ava-labs/avalanchego/build/avalanchego"),
		"avalanchego executable path",
	)
	flag.StringVar(
		&dataDir,
		config.DataDirKey,
		fmt.Sprintf("/tmp/subnet-evm-e2e-test/%v", time.Now().Unix()),
		"Data directory",
	)
	flag.StringVar(
		&configFilePath,
		config.ConfigFileKey,
		"",
		"Path to specify a config file",
	)
	flag.DurationVar(
		&setupTimeout,
		"setup-timeout",
		time.Minute,
		"Timeout for setting up the node for the e2e test (timeout for BeforeSuite to complete)",
	)
}

func TestE2E(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "subnet-evm e2e test suites")
}

// BeforeSuite starts an AvalancheGo process to use for the e2e tests
var _ = ginkgo.BeforeSuite(func() {
	ctx, cancel := context.WithTimeout(context.Background(), setupTimeout)
	defer cancel()

	var err error
	log.Info("Starting AvalancheGo node")
	startCmd, err = utils.RunCommand("./scripts/run_single_node.sh")
	gomega.Expect(err).Should(gomega.BeNil())
	healthClient := health.NewClient(utils.DefaultLocalNodeURI)
	healthy, err := health.AwaitReady(ctx, healthClient, 5*time.Second)
	gomega.Expect(err).Should(gomega.BeNil())
	gomega.Expect(healthy).Should(gomega.BeTrue())
	log.Info("AvalancheGo node is healthy")
})

var _ = ginkgo.AfterSuite(func() {
	// TODO add a new node to bootstrap off of the existing node and make sure we can bootstrap all of the data created in the test.
	gomega.Expect(startCmd).ShouldNot(gomega.BeNil())
	gomega.Expect(startCmd.Stop()).Should(gomega.BeNil())
})
