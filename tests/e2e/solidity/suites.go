// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements solidity tests.
package solidity

import (
	"fmt"
	"os/exec"

	"github.com/ava-labs/subnet-evm/tests/e2e/utils"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

const vmName = "subnetevm"

func runHardhatTests(test string) {
	cmd := exec.Command("npx", "hardhat", "test", test, "--network", "e2e")
	cmd.Dir = "./contract-examples"
	out, err := cmd.Output()
	fmt.Println(string(out))
	gomega.Expect(err).Should(gomega.BeNil())
}

var _ = utils.DescribePrecompile(func() {
	ginkgo.It("tx allow list", ginkgo.Label("solidity-with-npx"), func() {
		runHardhatTests("./test/ExampleTxAllowList.ts")
	})

	ginkgo.It("deployer allow list", ginkgo.Label("solidity-with-npx"), func() {
		runHardhatTests("./test/ExampleDeployerList.ts")
	})

	ginkgo.It("contract native minter", ginkgo.Label("solidity-with-npx"), func() {
		runHardhatTests("./test/ERC20NativeMinter.ts")
	})

	ginkgo.It("fee manager", ginkgo.Label("solidity-with-npx"), func() {
		runHardhatTests("./test/ExampleFeeManager.ts")
	})

	// ADD YOUR PRECOMPILE HERE
	/*
			ginkgo.It("your precompile", ginkgo.Label("solidity-with-npx"), func() {
			runHardhatTests("./test/Example{YourPrecompile}Test.ts")
		})
	*/
})
