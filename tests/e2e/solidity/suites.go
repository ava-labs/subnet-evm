// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements solidity tests.
package solidity

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ava-labs/avalanchego/api/health"
	"github.com/ava-labs/avalanchego/genesis"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	wallet "github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/ava-labs/subnet-evm/plugin/evm"
	"github.com/ava-labs/subnet-evm/tests/e2e/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var localURI = "http://127.0.0.1:9650"

var _ = ginkgo.Describe("[Precompiles]", ginkgo.Ordered, func() {
	ginkgo.It("ping the network", ginkgo.Label("setup"), func() {
		client := health.NewClient(localURI)
		healthy, err := client.Readiness(context.Background())
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(healthy.Healthy).Should(gomega.BeTrue())
	})
})

func runHardhatTests(test string, rpcURI string) {
	log.Info("Sleeping to wait for test ping", "rpcURI", rpcURI)
	time.Sleep(time.Minute)
	client, err := utils.NewEvmClient(rpcURI, 225, 2) // TODO this is failing because the Avalanche engine does not start bootstrapping of subnets when staking is disabled
	gomega.Expect(err).Should(gomega.BeNil())

	bal, err := client.FetchBalance(context.Background(), common.HexToAddress(""))
	gomega.Expect(err).Should(gomega.BeNil())
	gomega.Expect(bal).Should(gomega.Equal(common.Big0))

	// err := os.Setenv("RPC_URI", rpcURI)
	// gomega.Expect(err).Should(gomega.BeNil())

	// utils.RunCommand(fmt.Sprintf("npx hardhat test %s", "--network=local"))
	// cmd := exec.Command("npx", "hardhat", "test", test, "--network", "local")
	// cmd.Env = append(cmd.Env, fmt.Sprintf("RPC_URI=%s", rpcURI))
	// cmd.Dir = "./contract-examples"
	// out, err := cmd.Output()
	// if err != nil {
	// 	fmt.Println(string(out))
	// 	fmt.Println(err)
	// }
	// gomega.Expect(err).Should(gomega.BeNil())
}

func executeHardHatTestOnNewBlockchain(ctx context.Context, test string) {
	log.Info("Executing HardHat tests on a new blockchain", "test", test)
	kc := secp256k1fx.NewKeychain(genesis.EWOQKey)

	// NewWalletFromURI fetches the available UTXOs owned by [kc] on the network
	// that [LocalAPIURI] is hosting.
	wallet, err := wallet.NewWalletFromURI(ctx, utils.DefaultLocalNodeURI, kc)
	gomega.Expect(err).Should(gomega.BeNil())

	pWallet := wallet.P()

	owner := &secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs: []ids.ShortID{
			genesis.EWOQKey.PublicKey().Address(),
		},
	}

	wd, err := os.Getwd()
	gomega.Expect(err).Should(gomega.BeNil())
	genesisFilePath := fmt.Sprintf("./tests/e2e/genesis/%s.json", test)
	log.Info("Reading genesis file", "filePath", genesisFilePath, "pwd", wd)
	genesisBytes, err := os.ReadFile(genesisFilePath)
	gomega.Expect(err).Should(gomega.BeNil())

	log.Info("Creating new subnet")
	createSubnetTxID, err := pWallet.IssueCreateSubnetTx(owner)
	gomega.Expect(err).Should(gomega.BeNil())

	log.Info("Creating new blockchain", "genesis", genesisBytes)
	createChainTx, err := pWallet.IssueCreateChainTx(
		createSubnetTxID,
		genesisBytes,
		evm.ID,
		nil,
		strings.ReplaceAll(test, "_", ""),
	)
	gomega.Expect(err).Should(gomega.BeNil())

	// Confirm the new blockchain is ready by waiting for the readiness endpoint
	healthClient := health.NewClient(utils.DefaultLocalNodeURI)
	healthy, err := healthClient.AwaitReady(ctx, 5*time.Second)
	gomega.Expect(err).Should(gomega.BeNil())
	gomega.Expect(healthy).Should(gomega.BeTrue())

	// Confirm the new blockchain is up
	chainURI := fmt.Sprintf("%s/ext/bc/%s/rpc", utils.DefaultLocalNodeURI, createChainTx.String())

	runHardhatTests(test, chainURI)
}

var _ = ginkgo.Describe("[Precompiles]", ginkgo.Ordered, func() {
	ginkgo.It("create blockchain", ginkgo.Label("test"), func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		executeHardHatTestOnNewBlockchain(ctx, "contract_native_minter")
	})

	// ginkgo.It("tx allow list", ginkgo.Label("solidity-with-npx"), func() {
	// 	err := startSubnet("./tests/e2e/genesis/tx_allow_list.json")
	// 	gomega.Expect(err).Should(gomega.BeNil())
	// 	running := runner.IsRunnerUp(grpcEp)
	// 	gomega.Expect(running).Should(gomega.BeTrue())
	// 	runHardhatTests("./test/ExampleTxAllowList.ts")
	// 	stopSubnet()
	// 	running = runner.IsRunnerUp(grpcEp)
	// 	gomega.Expect(running).Should(gomega.BeFalse())
	// })

	// ginkgo.It("deployer allow list", ginkgo.Label("solidity-with-npx"), func() {
	// 	err := startSubnet("./tests/e2e/genesis/deployer_allow_list.json")
	// 	gomega.Expect(err).Should(gomega.BeNil())
	// 	running := runner.IsRunnerUp(grpcEp)
	// 	gomega.Expect(running).Should(gomega.BeTrue())
	// 	runHardhatTests("./test/ExampleDeployerList.ts")
	// 	stopSubnet()
	// 	running = runner.IsRunnerUp(grpcEp)
	// 	gomega.Expect(running).Should(gomega.BeFalse())
	// })

	// ginkgo.It("contract native minter", ginkgo.Label("solidity-with-npx"), func() {
	// 	err := startSubnet("./tests/e2e/genesis/contract_native_minter.json")
	// 	gomega.Expect(err).Should(gomega.BeNil())
	// 	running := runner.IsRunnerUp(grpcEp)
	// 	gomega.Expect(running).Should(gomega.BeTrue())
	// 	runHardhatTests("./test/ERC20NativeMinter.ts")
	// 	stopSubnet()
	// 	running = runner.IsRunnerUp(grpcEp)
	// 	gomega.Expect(running).Should(gomega.BeFalse())
	// })

	// ginkgo.It("fee manager", ginkgo.Label("solidity-with-npx"), func() {
	// 	err := startSubnet("./tests/e2e/genesis/fee_manager.json")
	// 	gomega.Expect(err).Should(gomega.BeNil())
	// 	running := runner.IsRunnerUp(grpcEp)
	// 	gomega.Expect(running).Should(gomega.BeTrue())
	// 	runHardhatTests("./test/ExampleFeeManager.ts")
	// 	stopSubnet()
	// 	running = runner.IsRunnerUp(grpcEp)
	// 	gomega.Expect(running).Should(gomega.BeFalse())
	// })

	// ginkgo.It("reward manager", ginkgo.Label("solidity-with-npx"), func() {
	// 	err := startSubnet("./tests/e2e/genesis/reward_manager.json")
	// 	gomega.Expect(err).Should(gomega.BeNil())
	// 	running := runner.IsRunnerUp(grpcEp)
	// 	gomega.Expect(running).Should(gomega.BeTrue())
	// 	runHardhatTests("./test/ExampleRewardManager.ts")
	// 	stopSubnet()
	// 	running = runner.IsRunnerUp(grpcEp)
	// 	gomega.Expect(running).Should(gomega.BeFalse())
	// })

	// ADD YOUR PRECOMPILE HERE
	/*
			ginkgo.It("your precompile", ginkgo.Label("solidity-with-npx"), func() {
			err := startSubnet("./tests/e2e/genesis/{your_precompile}.json")
			gomega.Expect(err).Should(gomega.BeNil())
			running := runner.IsRunnerUp(grpcEp)
			gomega.Expect(running).Should(gomega.BeTrue())
			runHardhatTests("./test/{YourPrecompileTest}.ts")
			stopSubnet()
			running = runner.IsRunnerUp(grpcEp)
			gomega.Expect(running).Should(gomega.BeFalse())
		})
	*/
})
