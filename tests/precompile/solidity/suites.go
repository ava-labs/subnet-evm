// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements solidity tests.
package solidity

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/ava-labs/avalanchego/api/health"
	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/avalanchego/genesis"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	wallet "github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/plugin/evm"
	"github.com/ava-labs/subnet-evm/tests/precompile/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("[Precompiles]", ginkgo.Ordered, func() {
	ginkgo.It("ping the network", ginkgo.Label("setup"), func() {
		client := health.NewClient(utils.DefaultLocalNodeURI)
		healthy, err := client.Readiness(context.Background())
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(healthy.Healthy).Should(gomega.BeTrue())
	})
})

func runHardhatTests(test string, rpcURI string) {
	log.Info("Sleeping to wait for test ping", "rpcURI", rpcURI)
	client, err := utils.NewEvmClient(rpcURI, 225, 2) // TODO this is failing because the Avalanche engine does not start bootstrapping of subnets when staking is disabled
	gomega.Expect(err).Should(gomega.BeNil())

	bal, err := client.FetchBalance(context.Background(), common.HexToAddress(""))
	gomega.Expect(err).Should(gomega.BeNil())
	gomega.Expect(bal.Cmp(common.Big0)).Should(gomega.Equal(0))

	err = os.Setenv("RPC_URI", rpcURI)
	gomega.Expect(err).Should(gomega.BeNil())
	cmd := exec.Command("npx", "hardhat", "test", fmt.Sprintf("./test/%s.ts", test), "--network", "local")
	cmd.Dir = "./contract-examples"
	log.Info("Running hardhat command", "cmd", cmd.String())

	out, err := cmd.CombinedOutput()
	fmt.Printf("\nCombined output:\n\n%s\n", string(out))
	gomega.Expect(err).Should(gomega.BeNil())
}

func createNewSubnet(ctx context.Context, genesisFilePath string) string {
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
	log.Info("Reading genesis file", "filePath", genesisFilePath, "pwd", wd)
	genesisBytes, err := os.ReadFile(genesisFilePath)
	gomega.Expect(err).Should(gomega.BeNil())

	log.Info("Creating new subnet")
	createSubnetTxID, err := pWallet.IssueCreateSubnetTx(owner)
	gomega.Expect(err).Should(gomega.BeNil())

	genesis := &core.Genesis{}
	err = json.Unmarshal(genesisBytes, genesis)
	gomega.Expect(err).Should(gomega.BeNil())

	log.Info("Creating new Subnet-EVM blockchain", "genesis", genesis)
	createChainTxID, err := pWallet.IssueCreateChainTx(
		createSubnetTxID,
		genesisBytes,
		evm.ID,
		nil,
		"testChain",
	)
	gomega.Expect(err).Should(gomega.BeNil())

	// Confirm the new blockchain is ready by waiting for the readiness endpoint
	infoClient := info.NewClient(utils.DefaultLocalNodeURI)
	bootstrapped, err := info.AwaitBootstrapped(ctx, infoClient, createChainTxID.String(), 5*time.Second)
	gomega.Expect(err).Should(gomega.BeNil())
	gomega.Expect(bootstrapped).Should(gomega.BeTrue())

	// Return the RPC Endpoint for the new blockchain
	return fmt.Sprintf("%s/ext/bc/%s/rpc", utils.DefaultLocalNodeURI, createChainTxID.String())
}

func executeHardHatTestOnNewBlockchain(ctx context.Context, test string) {
	log.Info("Executing HardHat tests on a new blockchain", "test", test)

	genesisFilePath := fmt.Sprintf("./tests/precompile/genesis/%s.json", test)

	createSubnetCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	chainURI := createNewSubnet(createSubnetCtx, genesisFilePath)

	log.Info("Created subnet successfully", "ChainURI", chainURI)
	runHardhatTests(test, chainURI)
}

// TODO: can we move where we register the precompile e2e tests, so that they stay within their package
var _ = ginkgo.Describe("[Precompiles]", ginkgo.Ordered, func() {
	// Each ginkgo It node specifies the name of the genesis file (in ./tests/precompile/genesis/)
	//to use to launch the subnet and the name of the TS test file to run on the subnet (in ./contract-examples/tests/)
	ginkgo.It("contract native minter", ginkgo.Label("solidity-with-npx"), func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		executeHardHatTestOnNewBlockchain(ctx, "contract_native_minter")
	})

	// ginkgo.It("tx allow list", ginkgo.Label("solidity-with-npx"), func() {
	// 	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	// 	defer cancel()

	// 	executeHardHatTestOnNewBlockchain(ctx, "tx_allow_list")
	// })

	// ginkgo.It("contract deployer allow list", ginkgo.Label("solidity-with-npx"), func() {
	// 	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	// 	defer cancel()

	// 	executeHardHatTestOnNewBlockchain(ctx, "contract_deployer_allow_list")
	// })

	// ginkgo.It("fee manager", ginkgo.Label("solidity-with-npx"), func() {
	// 	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	// 	defer cancel()

	// 	executeHardHatTestOnNewBlockchain(ctx, "fee_manager")
	// })

	// ginkgo.It("reward manager", ginkgo.Label("solidity-with-npx"), func() {
	// 	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	// 	defer cancel()

	// 	executeHardHatTestOnNewBlockchain(ctx, "reward_manager")
	// })

	// ADD YOUR PRECOMPILE HERE
	/*
		ginkgo.It("your precompile", ginkgo.Label("solidity-with-npx"), func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			// Specify the name shared by the genesis file in ./tests/precompile/genesis/{your_precompile}.json
			// and the test file in ./contract-examples/tests/{your_precompile}.ts
			executeHardHatTestOnNewBlockchain(ctx, "your_precompile")
		})
	*/
})
