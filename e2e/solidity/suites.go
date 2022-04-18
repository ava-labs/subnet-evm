// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements ping tests, requires network-runner cluster.
package solidity

import (
	"context"
	"time"

	ginkgo "github.com/onsi/ginkgo/v2"

	runner_client "github.com/ava-labs/avalanche-network-runner/client"

	"github.com/onsi/gomega"

	"github.com/ava-labs/avalanchego/tests"
	"github.com/ava-labs/subnet-evm/e2e"
)

var _ = e2e.DescribePrecompile("[Solidity]", func() {
	ginkgo.BeforeEach(func() {
		if e2e.GetRunnerGRPCEndpoint() == "" {
			ginkgo.Skip("no local network-runner, failing")
		}
	})

	ginkgo.AfterEach(func() {
		// if e2e.GetRunnerGRPCEndpoint() == "" {
		// 	ginkgo.Fail("no local network-runner, skipping")
		// }
		// runnerCli := e2e.GetRunnerClient()
		// ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		// runnerCli.Stop(ctx)
		// cancel()

		if e2e.GetRunnerGRPCEndpoint() != "" {
			runnerCli := e2e.GetRunnerClient()
			gomega.Expect(runnerCli).ShouldNot(gomega.BeNil())

			tests.Outf("{{red}}shutting down network-runner cluster{{/}}\n")
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			_, err := runnerCli.Stop(ctx)
			cancel()
			gomega.Expect(err).Should(gomega.BeNil())

			tests.Outf("{{red}}shutting down network-runner client{{/}}\n")
			err = e2e.CloseRunnerClient()
			gomega.Expect(err).Should(gomega.BeNil())
		}
	})

	ginkgo.It("can ping network-runner RPC server", func() {

		runnerCli := e2e.GetRunnerClient()
		gomega.Expect(runnerCli).ShouldNot(gomega.BeNil())

		execPath := e2e.GetExecPath()
		logLevel := e2e.GetLogLevel()

		tests.Outf("{{magenta}}starting network-runner with %q{{/}}\n", execPath)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		resp, err := runnerCli.Start(ctx, execPath, runner_client.WithLogLevel(logLevel))
		cancel()
		gomega.Expect(err).Should(gomega.BeNil())
		tests.Outf("{{green}}successfully started network-runner :{{/}} %+v\n", resp.ClusterInfo.NodeNames)

		// start is async, so wait some time for cluster health
		time.Sleep(time.Minute)

		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Minute)
		_, err = runnerCli.Health(ctx)
		cancel()
		gomega.Expect(err).Should(gomega.BeNil())

		var uriSlice []string
		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Minute)
		uriSlice, err = runnerCli.URIs(ctx)
		cancel()
		gomega.Expect(err).Should(gomega.BeNil())
		e2e.SetURIs(uriSlice)

		gomega.Expect(err).Should(gomega.BeNil())
	})
})
