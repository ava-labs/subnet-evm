// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements counter tests, requires network-runner cluster.
package counter

import (
	"os"
	"strings"
	"time"

	"github.com/ava-labs/subnet-evm/tests/e2e/utils"
	"github.com/ethereum/go-ethereum/common/math"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = utils.DescribeLocal("[Solidity Counter]", func() {
	ginkgo.It("can deploy counter contract", func() {
		ci := utils.GetClusterInfo()
		gomega.Expect(len(ci.URIs) > 0).Should(gomega.BeTrue())

		contractsFoundryDir := utils.GetContractsFoundryDir()
		utils.Outf("{{green}}testing contracts '%s' to:{{/}} %q\n", contractsFoundryDir, ci.SubnetEVMRPCEndpoints)
		gomega.Expect(os.Chdir(contractsFoundryDir)).Should(gomega.BeNil())

		s, err := utils.RunCommand(2*time.Minute, "forge", "test", "-vvv")
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(s.Complete && s.Exit == 0 && s.Error == nil).Should(gomega.BeTrue())

		utils.Outf("{{green}}deploying counter contract using foundry to:{{/}} %q\n", ci.SubnetEVMRPCEndpoints)
		s, err = utils.RunCommand(
			2*time.Minute,
			"forge",
			"create",
			"src/Counter.sol:Counter",
			"--private-key=56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027", // ewoq key
			"--rpc-url="+ci.SubnetEVMRPCEndpoints[0],
		)
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(s.Complete && s.Exit == 0 && s.Error == nil).Should(gomega.BeTrue())
		utils.Outf("{{green}}command output:{{/}}\n\n%s\n\n", strings.Join(s.Stdout, "\n"))

		// "Deployed to:" is the contract address
		contractAddr := ""
		for _, line := range s.Stdout {
			if strings.HasPrefix(line, "Deployed to: ") {
				contractAddr = strings.Replace(line, "Deployed to: ", "", 1)
				break
			}
		}
		gomega.Expect(contractAddr).ShouldNot(gomega.BeEmpty())
		utils.Outf("{{green}}counter contract address:{{/}} %q\n", contractAddr)

		utils.Outf("{{green}}set the current counter number{{/}}\n")
		s, err = utils.RunCommand(
			2*time.Minute,
			"cast",
			"send",
			"--private-key=56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027", // ewoq key
			"--rpc-url="+ci.SubnetEVMRPCEndpoints[0],
			contractAddr,
			"setNumber(uint256)",
			"100",
		)
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(s.Complete && s.Exit == 0 && s.Error == nil).Should(gomega.BeTrue())

		utils.Outf("{{green}}fetching the current counter number{{/}}\n")
		s, err = utils.RunCommand(
			2*time.Minute,
			"cast",
			"call",
			"--rpc-url="+ci.SubnetEVMRPCEndpoints[0],
			contractAddr,
			"number()",
		)
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(s.Complete && s.Exit == 0 && s.Error == nil).Should(gomega.BeTrue())

		bigNum, _ := math.ParseBig256(strings.TrimSpace(strings.Join(s.Stdout, "")))
		curCnt := bigNum.Uint64()
		gomega.Expect(curCnt).Should(gomega.BeNumerically("==", 100))
	})
})
