// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ava-labs/avalanchego/api/health"
	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/avalanchego/genesis"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	wallet "github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/plugin/evm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/go-cmd/cmd"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

type SubnetSuite struct {
	GetBlockchainID func(alias string) string
}

// TODO: find a better way rather than using a global variable
// This is used to pass the blockchain IDs from the SynchronizedBeforeSuite() to the tests
var globalSuite SubnetSuite

// CreateSubnetsSuite creates subnets for given [genesisFiles], and registers a before suite that starts an AvalancheGo process to use for the e2e tests.
// genesisFiles is a map of test aliases to genesis file paths.
func CreateSubnetsSuite(genesisFiles map[string]string) *SubnetSuite {
	// Keep track of the AvalancheGo external bash script, it is null for most
	// processes except the first process that starts AvalancheGo
	var startCmd *cmd.Cmd

	// Our test suite runs in separate processes, ginkgo has
	// SynchronizedBeforeSuite() which runs once, and its return value is passed
	// over to each worker.
	//
	// Here an AvalancheGo node instance is started, and subnets are created for
	// each test case. Each test case has its own subnet, therefore all tests
	// can run in parallel without any issue.
	//
	var _ = ginkgo.SynchronizedBeforeSuite(func() []byte {
		ctx, cancel := context.WithTimeout(context.Background(), bootAvalancheNodeTimeout)
		defer cancel()

		wd, err := os.Getwd()
		gomega.Expect(err).Should(gomega.BeNil())
		log.Info("Starting AvalancheGo node", "wd", wd)
		cmd, err := RunCommand("./scripts/run.sh")
		startCmd = cmd
		gomega.Expect(err).Should(gomega.BeNil())

		// Assumes that startCmd will launch a node with HTTP Port at [utils.DefaultLocalNodeURI]
		healthClient := health.NewClient(DefaultLocalNodeURI)
		healthy, err := health.AwaitReady(ctx, healthClient, healthCheckTimeout, nil)
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(healthy).Should(gomega.BeTrue())
		log.Info("AvalancheGo node is healthy")

		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		blockchainIDs := make(map[string]string)
		for alias, file := range genesisFiles {
			blockchainIDs[alias] = CreateNewSubnet(ctx, file)
		}

		blockchainIDsBytes, err := json.Marshal(blockchainIDs)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		return blockchainIDsBytes
	}, func(ctx ginkgo.SpecContext, data []byte) {
		blockchainIDs := make(map[string]string)
		err := json.Unmarshal(data, &blockchainIDs)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		globalSuite.GetBlockchainID = func(alias string) string {
			return blockchainIDs[alias]
		}
	})

	// SynchronizedAfterSuite() takes two functions, the first runs after each test suite is done and the second
	// function is executed once when all the tests are done. This function is used
	// to gracefully shutdown the AvalancheGo node.
	var _ = ginkgo.SynchronizedAfterSuite(func() {}, func() {
		gomega.Expect(startCmd).ShouldNot(gomega.BeNil())
		gomega.Expect(startCmd.Stop()).Should(gomega.BeNil())
	})

	return &globalSuite
}

// RunDefaultHardhatTests runs the hardhat tests on a given blockchain ID
// with default parameters. Default parameters are:
// 1. Hardhat contract environment is located at ./contracts
// 2. Hardhat test file is located at ./contracts/test/<test>.ts
// 3. npx is available in the ./contracts directory
// 4. CreateSubnetsSynchronized() called before this function and with the [test] aliased to genesis file
func (s *SubnetSuite) RunHardhatTests(ctx context.Context, test string) {
	blockchainID := s.GetBlockchainID(test)
	runHardhatTests(ctx, blockchainID, test)
}

// CreateNewSubnet creates a new subnet and Subnet-EVM blockchain with the given genesis file.
// returns the ID of the new created blockchain.
func CreateNewSubnet(ctx context.Context, genesisFilePath string) string {
	kc := secp256k1fx.NewKeychain(genesis.EWOQKey)

	// NewWalletFromURI fetches the available UTXOs owned by [kc] on the network
	// that [LocalAPIURI] is hosting.
	wallet, err := wallet.NewWalletFromURI(ctx, DefaultLocalNodeURI, kc)
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
	log.Info("Reading genesis file", "filePath", genesisFilePath, "wd", wd)
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
	infoClient := info.NewClient(DefaultLocalNodeURI)
	bootstrapped, err := info.AwaitBootstrapped(ctx, infoClient, createChainTxID.String(), 2*time.Second)
	gomega.Expect(err).Should(gomega.BeNil())
	gomega.Expect(bootstrapped).Should(gomega.BeTrue())

	// Return the blockchainID of the newly created blockchain
	return createChainTxID.String()
}

// GetDefaultChainURI returns the default chain URI for a given blockchainID
func GetDefaultChainURI(blockchainID string) string {
	return fmt.Sprintf("%s/ext/bc/%s/rpc", DefaultLocalNodeURI, blockchainID)
}

// CreateAndRunHardhatTests creates a subnet and blockchain and then
// runs the hardhat tests on the new blockchain with default parameters.
//
//	Default parameters are:
//
// 1. Genesis file is located at ./tests/precompile/genesis/<test>.json
// 2. Hardhat contract environment is located at ./contracts
// 3. Hardhat test file is located at ./contracts/test/<test>.ts
// 4. npx is available in the ./contracts directory
func CreateAndRunHardhatTests(ctx context.Context, test string) {
	genesisFilePath := fmt.Sprintf("./tests/precompile/genesis/%s.json", test)

	blockchainID := CreateNewSubnet(ctx, genesisFilePath)
	log.Info("Created subnet successfully", "blockchainID", blockchainID)

	runHardhatTests(ctx, blockchainID, test)
}

// GetFilesAndAliases returns a map of aliases to file paths in given [dir].
func GetFilesAndAliases(dir string) (map[string]string, error) {
	files, err := filepath.Glob(dir)
	if err != nil {
		return nil, err
	}
	aliasesToFiles := make(map[string]string)
	for _, file := range files {
		alias := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		aliasesToFiles[alias] = file
	}
	return aliasesToFiles, nil
}

func runHardhatTests(ctx context.Context, blockchainID string, test string) {
	chainURI := GetDefaultChainURI(blockchainID)
	log.Info(
		"Executing HardHat tests on blockchain",
		"blockchainID", blockchainID,
		"test", test,
		"ChainURI", chainURI,
	)

	cmdPath := "./contracts"
	// test path is relative to the cmd path
	testPath := fmt.Sprintf("./test/%s.ts", test)
	cmd := exec.Command("npx", "hardhat", "test", testPath, "--network", "local")
	cmd.Dir = cmdPath

	RunTestCMD(cmd, chainURI)
}
