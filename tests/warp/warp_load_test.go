// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Load test for AWM. Shares the setup with the warp test.
package warp

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/ava-labs/subnet-evm/tests/utils/runner"
	"github.com/ethereum/go-ethereum/log"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("[AWM Load Simulator]", ginkgo.Ordered, func() {
	ginkgo.It("Run AWM load simulator against local nodes", ginkgo.Label("Warp", "load"), func() {
		subnetIDs := manager.GetSubnets()
		gomega.Expect(len(subnetIDs)).Should(gomega.Equal(2))

		// Create a temporary file to store the subnet information
		// so that the load simulator can read it.
		subnets := make([]*runner.Subnet, 0, len(subnetIDs))
		for _, subnetID := range subnetIDs {
			subnetDetails, ok := manager.GetSubnet(subnetID)
			gomega.Expect(ok).Should(gomega.BeTrue())
			subnets = append(subnets, subnetDetails)
		}
		jsonBytes, err := json.MarshalIndent(subnets, "", "  ")
		gomega.Expect(err).Should(gomega.BeNil())
		tmpFile, err := os.CreateTemp("", "subnet-info-*.json")
		gomega.Expect(err).Should(gomega.BeNil())
		tmpFileName := tmpFile.Name()
		_, err = tmpFile.Write(jsonBytes)
		gomega.Expect(err).Should(gomega.BeNil())
		err = tmpFile.Close()
		gomega.Expect(err).Should(gomega.BeNil())

		log.Info(
			"Running load simulator...",
			"rpcEndpointsFile", tmpFileName,
		)
		cmd := exec.Command("./scripts/run_simulator.sh")
		additionalEnv := []string{"RPC_ENDPOINTS_FILE=" + tmpFileName}

		cmd.Env = os.Environ() // Inherit environment variables from parent
		cmd.Env = append(cmd.Env, additionalEnv...)
		log.Info("Running load simulator script", "env", additionalEnv, "cmd", cmd.String())

		out, err := cmd.CombinedOutput()
		fmt.Printf("\nCombined output:\n\n%s\n", string(out))
		gomega.Expect(err).Should(gomega.BeNil())
	})
})
