// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// e2e implements the e2e tests.
package load

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ava-labs/avalanche-network-runner/client"
	"github.com/ava-labs/avalanche-network-runner/rpcpb"
	"github.com/ava-labs/avalanche-network-runner/utils/constants"
	"github.com/ava-labs/avalanche-network-runner/ux"
	"github.com/ava-labs/avalanchego/utils/logging"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestE2E(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "load test")
}

// What do I need to be present here?
// AvalancheGo
// Plugin directory with plugin binary


var (
	logLevel      string
	logDir        string
	gRPCEp        string
	gRPCGatewayEp string
	execPath1     string
	execPath2     string
	subnetEvmPath string

	newNodeName       = "test-add-node"
	customNodeConfigs = map[string]string{
		"node1": `{"api-admin-enabled":true}`,
		"node2": `{"api-admin-enabled":true}`,
		"node3": `{"api-admin-enabled":true}`,
		"node4": `{"api-admin-enabled":false}`,
		"node5": `{"api-admin-enabled":false}`,
		"node6": `{"api-admin-enabled":false}`,
		"node7": `{"api-admin-enabled":false}`,
	}
	numNodes = uint32(5)
)

func init() {
	flag.StringVar(
		&logLevel,
		"log-level",
		logging.Info.String(),
		"log level",
	)
	flag.StringVar(
		&logDir,
		"log-dir",
		"",
		"log directory",
	)
	flag.StringVar(
		&gRPCEp,
		"grpc-endpoint",
		"0.0.0.0:8080",
		"gRPC server endpoint",
	)
	flag.StringVar(
		&gRPCGatewayEp,
		"grpc-gateway-endpoint",
		"0.0.0.0:8081",
		"gRPC gateway endpoint",
	)
	flag.StringVar(
		&execPath1,
		"avalanchego-path-1",
		"",
		"avalanchego executable path (to upgrade from)",
	)
	flag.StringVar(
		&execPath2,
		"avalanchego-path-2",
		"",
		"avalanchego executable path (to upgrade to)",
	)
	flag.StringVar(
		&subnetEvmPath,
		"subnet-evm-path",
		"",
		"path to subnet-evm binary",
	)
}

var (
	cli client.Client
	log logging.Logger
)

var _ = ginkgo.BeforeSuite(func() {
	var err error
	logDir, err = os.MkdirTemp("", fmt.Sprintf("subnet-evm-load-test-%d", time.Now().Unix()))
	gomega.Ω(err).Should(gomega.BeNil())
	lvl, err := logging.ToLevel(logLevel)
	gomega.Ω(err).Should(gomega.BeNil())
	lcfg := logging.Config{
		DisplayLevel: lvl,
	}
	logFactory := logging.NewFactory(lcfg)
	log, err = logFactory.Make(constants.LogNameTest)
	gomega.Ω(err).Should(gomega.BeNil())

	cli, err = client.New(client.Config{
		Endpoint:    gRPCEp,
		DialTimeout: 10 * time.Second,
	}, log)
	gomega.Ω(err).Should(gomega.BeNil())
})

var _ = ginkgo.AfterSuite(func() {
	ux.Print(log, logging.Red.Wrap("shutting down cluster"))
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	_, err := cli.Stop(ctx)
	cancel()
	gomega.Ω(err).Should(gomega.BeNil())

	ux.Print(log, logging.Red.Wrap("shutting down client"))
	err = cli.Close()
	gomega.Ω(err).Should(gomega.BeNil())
})

var _ = ginkgo.Describe("[Start/Remove/Restart/Add/Stop]", func() {
	ginkgo.It("test", func() {
		ginkgo.By("start with blockchain specs", func() {
			ux.Print(log, logging.Green.Wrap("sending 'start' with the valid binary path: %s"), execPath1)
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			resp, err := cli.Start(ctx, execPath1,
				client.WithBlockchainSpecs([]*rpcpb.BlockchainSpec{
					{
						VmName:  "subnetevm",
						Genesis: "tests/e2e/subnet-evm-genesis.json",
					},
				}),
			)
			cancel()
			gomega.Ω(err).Should(gomega.BeNil())
			ux.Print(log, logging.Green.Wrap("successfully started, node-names: %s"), resp.ClusterInfo.NodeNames)
		})
	})
})

func waitForCustomChainsHealthy() string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	var created bool
	continueLoop := true
	for continueLoop {
		select {
		case <-ctx.Done():
			continueLoop = false
		case <-time.After(5 * time.Second):
			cctx, ccancel := context.WithTimeout(context.Background(), 15*time.Second)
			status, err := cli.Status(cctx)
			ccancel()
			gomega.Ω(err).Should(gomega.BeNil())
			created = status.ClusterInfo.CustomChainsHealthy
			if created {
				existingSubnetID := status.ClusterInfo.GetSubnets()[0]
				gomega.Ω(existingSubnetID).Should(gomega.Not(gomega.BeNil()))
				cancel()
				return existingSubnetID
			}
		}
	}
	cancel()
	gomega.Ω(created).Should(gomega.Equal(true))
	return ""
}
