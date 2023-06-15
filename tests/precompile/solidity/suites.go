// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements solidity tests.
package solidity

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/ava-labs/subnet-evm/tests/utils"
	"github.com/ethereum/go-ethereum/log"
	ginkgo "github.com/onsi/ginkgo/v2"
)

func runHardhatTests(ctx context.Context, test string) {
	log.Info("Executing HardHat tests on a new blockchain", "test", test)

	genesisFilePath := fmt.Sprintf("./tests/precompile/genesis/%s.json", test)

	blockchainID := utils.CreateNewSubnet(ctx, genesisFilePath)
	chainURI := utils.GetDefaultChainURI(blockchainID)
	log.Info("Created subnet successfully", "ChainURI", chainURI)

	cmdPath := "./contracts"
	// test path is relative to the cmd path
	testPath := fmt.Sprintf("./test/%s.ts", test)
	cmd := exec.Command("npx", "hardhat", "test", testPath, "--network", "local")
	cmd.Dir = cmdPath

	utils.RunTestCMD(cmd, chainURI)
}

var _ = ginkgo.Describe("[Precompiles]", ginkgo.Ordered, func() {
	// Register the ping test first
	utils.RegisterPingTest()

	// Each ginkgo It node specifies the name of the genesis file (in ./tests/precompile/genesis/)
	// to use to launch the subnet and the name of the TS test file to run on the subnet (in ./contracts/tests/)
	ginkgo.It("contract native minter", ginkgo.Label("Precompile"), ginkgo.Label("ContractNativeMinter"), func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		runHardhatTests(ctx, "contract_native_minter")
	})

	ginkgo.It("tx allow list", ginkgo.Label("Precompile"), ginkgo.Label("TxAllowList"), func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		runHardhatTests(ctx, "tx_allow_list")
	})

	ginkgo.It("contract deployer allow list", ginkgo.Label("Precompile"), ginkgo.Label("ContractDeployerAllowList"), func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		runHardhatTests(ctx, "contract_deployer_allow_list")
	})

	ginkgo.It("fee manager", ginkgo.Label("Precompile"), ginkgo.Label("FeeManager"), func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		runHardhatTests(ctx, "fee_manager")
	})

	ginkgo.It("reward manager", ginkgo.Label("Precompile"), ginkgo.Label("RewardManager"), func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		runHardhatTests(ctx, "reward_manager")
	})

	// and then runs the hardhat tests for each one without forcing precompile developers to modify this file.
	// ADD YOUR PRECOMPILE HERE
	/*
		ginkgo.It("your precompile", ginkgo.Label("Precompile"), ginkgo.Label("YourPrecompile"), func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			// Specify the name shared by the genesis file in ./tests/precompile/genesis/{your_precompile}.json
			// and the test file in ./contracts/tests/{your_precompile}.ts
			// If you want to use a different test command and genesis path than the defaults, you can
			// use the utils.RunTestCMD. See runHardhatTests for an example.
			runHardhatTests(ctx, "your_precompile")
		})
	*/
})
