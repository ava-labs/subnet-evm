// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Load test for AWM. Shares the setup with the warp test.
package warp

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ava-labs/subnet-evm/tests/utils/runner"
	"github.com/ethereum/go-ethereum/log"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func rpcEndpoints(subnetDetails *runner.Subnet) []string {
	nodeURIs := subnetDetails.ValidatorURIs
	rpcEndpoints := make([]string, 0, len(nodeURIs))
	for _, uri := range nodeURIs {
		rpcEndpoints = append(rpcEndpoints, fmt.Sprintf("%s/ext/bc/%s/rpc", uri, subnetDetails.BlockchainID))
	}
	return rpcEndpoints
}

var _ = ginkgo.Describe("[AWM Load Simulator]", ginkgo.Ordered, func() {
	ginkgo.It("Run AWM load simulator against local nodes", ginkgo.Label("Warp", "load"), func() {
		subnetIDs := manager.GetSubnets()
		gomega.Expect(len(subnetIDs)).Should(gomega.Equal(2))

		subnetA := subnetIDs[0]
		subnetADetails, ok := manager.GetSubnet(subnetA)
		gomega.Expect(ok).Should(gomega.BeTrue())

		subnetARPCEndpoints := rpcEndpoints(subnetADetails)
		commaSeparatedRPCEndpointsA := strings.Join(subnetARPCEndpoints, ",")

		subnetB := subnetIDs[1]
		subnetBDetails, ok := manager.GetSubnet(subnetB)
		gomega.Expect(ok).Should(gomega.BeTrue())

		subnetBRPCEndpoints := rpcEndpoints(subnetBDetails)
		commaSeparatedRPCEndpointsB := strings.Join(subnetBRPCEndpoints, ",")

		log.Info(
			"Running load simulator...",
			"rpcEndpointsSubnetA", commaSeparatedRPCEndpointsA,
			"rpcEndpointsSubnetB", commaSeparatedRPCEndpointsB,
		)
		cmd := exec.Command("./scripts/run_simulator.sh")
		additionalEnv := []string{
			"SUBNET_B=" + subnetB.String(),
			"RPC_ENDPOINTS=" + commaSeparatedRPCEndpointsA,
			"RPC_ENDPOINTS_SUBNET_A=" + commaSeparatedRPCEndpointsA,
			"RPC_ENDPOINTS_SUBNET_B=" + commaSeparatedRPCEndpointsB,
		}

		cmd.Env = os.Environ() // Inherit environment variables from parent
		cmd.Env = append(cmd.Env, additionalEnv...)
		log.Info("Running load simulator script", "env", additionalEnv, "cmd", cmd.String())

		out, err := cmd.CombinedOutput()
		fmt.Printf("\nCombined output:\n\n%s\n", string(out))
		gomega.Expect(err).Should(gomega.BeNil())
	})
})
