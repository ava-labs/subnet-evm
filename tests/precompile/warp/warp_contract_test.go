// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements solidity tests.
package warp

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/ava-labs/avalanche-network-runner/rpcpb"
	"github.com/ava-labs/avalanchego/ids"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/interfaces"
	"github.com/ava-labs/subnet-evm/plugin/evm"
	"github.com/ava-labs/subnet-evm/tests/utils"
	"github.com/ava-labs/subnet-evm/tests/utils/runner"
	warpPayload "github.com/ava-labs/subnet-evm/warp/payload"
	"github.com/ava-labs/subnet-evm/x/warp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

const fundedKeyStr = "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027"

var (
	config              = runner.NewDefaultANRConfig()
	manager             = runner.NewNetworkManager(config)
	warpChainConfigPath string
)

func TestE2E(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "subnet-evm warp contract e2e test")
}

func toWebsocketURI(uri string, blockchainID string) string {
	return fmt.Sprintf("ws://%s/ext/bc/%s/ws", strings.TrimPrefix(uri, "http://"), blockchainID)
}

func toRPCURI(uri string, blockchainID string) string {
	return fmt.Sprintf("%s/ext/bc/%s/rpc", uri, blockchainID)
}

// BeforeSuite starts the default network and adds 10 new nodes as validators with BLS keys
// registered on the P-Chain.
// Adds two disjoint sets of 5 of the new validator nodes to validate two new subnets with a
// a single Subnet-EVM blockchain.
var _ = ginkgo.BeforeSuite(func() {
	ctx := context.Background()
	var err error
	// Name 10 new validators (which should have BLS key registered)
	subnetANodeNames := make([]string, 0)
	subnetBNodeNames := []string{}
	for i := 1; i <= 10; i++ {
		n := fmt.Sprintf("node%d-bls", i)
		if i <= 5 {
			subnetANodeNames = append(subnetANodeNames, n)
		} else {
			subnetBNodeNames = append(subnetBNodeNames, n)
		}
	}
	f, err := os.CreateTemp(os.TempDir(), "config.json")
	gomega.Expect(err).Should(gomega.BeNil())
	_, err = f.Write([]byte(`{"warp-api-enabled": true}`))
	gomega.Expect(err).Should(gomega.BeNil())
	warpChainConfigPath = f.Name()

	// Construct the network using the avalanche-network-runner
	_, err = manager.StartDefaultNetwork(ctx)
	gomega.Expect(err).Should(gomega.BeNil())
	err = manager.SetupNetwork(
		ctx,
		config.AvalancheGoExecPath,
		[]*rpcpb.BlockchainSpec{
			{
				VmName:      evm.IDStr,
				Genesis:     "./tests/precompile/genesis/warp.json",
				ChainConfig: warpChainConfigPath,
				SubnetSpec: &rpcpb.SubnetSpec{
					SubnetConfig: "",
					Participants: subnetANodeNames,
				},
			},
			{
				VmName:      evm.IDStr,
				Genesis:     "./tests/precompile/genesis/warp.json",
				ChainConfig: warpChainConfigPath,
				SubnetSpec: &rpcpb.SubnetSpec{
					SubnetConfig: "",
					Participants: subnetBNodeNames,
				},
			},
		},
	)
	gomega.Expect(err).Should(gomega.BeNil())

	// Issue transactions to activate the proposerVM fork on the receiving chain
	chainID := big.NewInt(99999)
	fundedKey, err := crypto.HexToECDSA(fundedKeyStr)
	gomega.Expect(err).Should(gomega.BeNil())
	subnetB := manager.GetSubnets()[1]
	subnetBDetails, ok := manager.GetSubnet(subnetB)
	gomega.Expect(ok).Should(gomega.BeTrue())

	chainBID := subnetBDetails.BlockchainID
	uri := toWebsocketURI(subnetBDetails.ValidatorURIs[0], chainBID.String())
	client, err := ethclient.Dial(uri)
	gomega.Expect(err).Should(gomega.BeNil())

	err = utils.IssueTxsToActivateProposerVMFork(ctx, chainID, fundedKey, client)
	gomega.Expect(err).Should(gomega.BeNil())
})

var _ = ginkgo.AfterSuite(func() {
	gomega.Expect(manager).ShouldNot(gomega.BeNil())
	gomega.Expect(manager.TeardownNetwork()).Should(gomega.BeNil())
	gomega.Expect(os.Remove(warpChainConfigPath)).Should(gomega.BeNil())
	// TODO: bootstrap an additional node to ensure that we can bootstrap the test data correctly
})

var _ = ginkgo.Describe("[Warp]", ginkgo.Ordered, func() {
	var (
		unsignedWarpMsg              *avalancheWarp.UnsignedMessage
		unsignedWarpMessageID        ids.ID
		blockchainIDA, blockchainIDB ids.ID
		chainAURIs, chainBURIs       []string
		chainAWSClient               ethclient.Client
		payload                      = []byte{1, 2, 3}
		err                          error
	)

	ginkgo.It("Setup URIs", ginkgo.Label("Warp", "SetupWarp"), func() {
		subnetIDs := manager.GetSubnets()
		gomega.Expect(len(subnetIDs)).Should(gomega.Equal(2))

		subnetA := subnetIDs[0]
		subnetADetails, ok := manager.GetSubnet(subnetA)
		gomega.Expect(ok).Should(gomega.BeTrue())
		blockchainIDA = subnetADetails.BlockchainID
		gomega.Expect(len(subnetADetails.ValidatorURIs)).Should(gomega.Equal(5))
		chainAURIs = append(chainAURIs, subnetADetails.ValidatorURIs...)

		subnetB := subnetIDs[1]
		subnetBDetails, ok := manager.GetSubnet(subnetB)
		gomega.Expect(ok).Should(gomega.BeTrue())
		blockchainIDB := subnetBDetails.BlockchainID
		gomega.Expect(len(subnetBDetails.ValidatorURIs)).Should(gomega.Equal(5))
		chainBURIs = append(chainBURIs, subnetBDetails.ValidatorURIs...)

		log.Info("Created URIs for both subnets", "ChainAURIs", chainAURIs, "ChainBURIs", chainBURIs, "blockchainIDA", blockchainIDA, "blockchainIDB", blockchainIDB)

		chainAWSURI := toWebsocketURI(chainAURIs[0], blockchainIDA.String())
		log.Info("Creating ethclient for blockchainA", "wsURI", chainAWSURI)
		chainAWSClient, err = ethclient.Dial(chainAWSURI)
		gomega.Expect(err).Should(gomega.BeNil())

		chainBWSURI := toWebsocketURI(chainBURIs[0], blockchainIDB.String())
		log.Info("Creating ethclient for blockchainB", "wsURI", chainBWSURI)
		gomega.Expect(err).Should(gomega.BeNil())

	})

	// Send a transaction to Subnet A to issue a Warp Message to Subnet B
	ginkgo.It("Send Message from A to B", ginkgo.Label("WarpContract", "SendWarpContract"), func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		log.Info("Subscribing to new heads")
		newHeads := make(chan *types.Header, 10)
		sub, err := chainAWSClient.SubscribeNewHead(ctx, newHeads)
		gomega.Expect(err).Should(gomega.BeNil())
		defer sub.Unsubscribe()

		cmdPath := "./contracts"
		// test path is relative to the cmd path
		testPath := "./test/warp.ts"

		rpcURI := toRPCURI(chainAURIs[0], blockchainIDA.String())
		senderAddress := common.HexToAddress("0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC")
		destinationAddress := common.HexToAddress("0x0550000000000000000000000000000000000000")
		addressedPayload, err := warpPayload.NewAddressedPayload(
			senderAddress,
			common.Hash(blockchainIDB),
			destinationAddress,
			payload,
		)
		gomega.Expect(err).Should(gomega.BeNil())
		expectedUnsignedMessage, err := avalancheWarp.NewUnsignedMessage(
			1337,
			blockchainIDA,
			addressedPayload.Bytes(),
		)
		gomega.Expect(err).Should(gomega.BeNil())

		os.Setenv("SENDER_ADDRESS", senderAddress.Hex())
		os.Setenv("SOURCE_CHAIN_ID", blockchainIDA.Hex())
		os.Setenv("DESTINATION_CHAIN_ID", blockchainIDB.Hex())
		os.Setenv("PAYLOAD", common.Bytes2Hex(payload))
		os.Setenv("DESTINATION_ADDRESS", destinationAddress.Hex())
		os.Setenv("EXPECTED_UNSIGNED_MESSAGE", hex.EncodeToString(expectedUnsignedMessage.Bytes()))
		runWarpHardhatTests(ctx, rpcURI, cmdPath, testPath)

		log.Info("Waiting for new block confirmation")
		newHead := <-newHeads
		blockHash := newHead.Hash()

		log.Info("Fetching relevant warp logs from the newly produced block")
		logs, err := chainAWSClient.FilterLogs(ctx, interfaces.FilterQuery{
			BlockHash: &blockHash,
			Addresses: []common.Address{warp.Module.Address},
		})
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(len(logs)).Should(gomega.Equal(1))

		// Check for relevant warp log from subscription and ensure that it matches
		// the log extracted from the last block.
		txLog := logs[0]
		log.Info("Parsing logData as unsigned warp message")
		unsignedMsg, err := warp.UnpackSendWarpEventDataToMessage(txLog.Data)
		gomega.Expect(err).Should(gomega.BeNil())

		// Set local variables for the duration of the test
		unsignedWarpMessageID = unsignedMsg.ID()
		unsignedWarpMsg = unsignedMsg
		log.Info("Parsed unsignedWarpMsg", "unsignedWarpMessageID", unsignedWarpMessageID, "unsignedWarpMessage", unsignedWarpMsg)

		// Loop over each client on chain A to ensure they all have time to accept the block.
		// Note: if we did not confirm this here, the next stage could be racy since it assumes every node
		// has accepted the block.
		for i, uri := range chainAURIs {
			chainAWSURI := toWebsocketURI(uri, blockchainIDA.String())
			log.Info("Creating ethclient for blockchainA", "wsURI", chainAWSURI)
			client, err := ethclient.Dial(chainAWSURI)
			gomega.Expect(err).Should(gomega.BeNil())

			// Loop until each node has advanced to >= the height of the block that emitted the warp log
			for {
				block, err := client.BlockByNumber(ctx, nil)
				gomega.Expect(err).Should(gomega.BeNil())
				if block.NumberU64() >= newHead.Number.Uint64() {
					log.Info("client accepted the block containing SendWarpMessage", "client", i, "height", block.NumberU64())
					break
				}
			}
		}
	})
})

func runWarpHardhatTests(ctx context.Context, chainURI string, execPath string, testPath string) {
	log.Info(
		"Executing HardHat tests on warp blockchain",
		"testPath", testPath,
		"ChainURI", chainURI,
	)

	cmd := exec.Command("npx", "hardhat", "test", testPath, "--network", "local")
	cmd.Dir = execPath

	log.Info("Sleeping to wait for test ping", "rpcURI", chainURI)
	err := os.Setenv("RPC_URI", chainURI)
	gomega.Expect(err).Should(gomega.BeNil())
	log.Info("Running test command", "cmd", cmd.String())

	out, err := cmd.CombinedOutput()
	fmt.Printf("\nCombined output:\n\n%s\n", string(out))
	gomega.Expect(err).Should(gomega.BeNil())
}
