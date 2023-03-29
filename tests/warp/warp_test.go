// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements solidity tests.
package solidity

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"

	"github.com/ava-labs/avalanche-network-runner/rpcpb"
	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/set"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/interfaces"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/plugin/evm"
	"github.com/ava-labs/subnet-evm/precompile/contracts/warp"
	"github.com/ava-labs/subnet-evm/tests/utils"
	"github.com/ava-labs/subnet-evm/tests/utils/runner"
	warpBackend "github.com/ava-labs/subnet-evm/warp"
	warpTransaction "github.com/ava-labs/subnet-evm/warp/transaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var (
	config  = runner.NewDefaultANRConfig()
	manager = runner.NewNetworkManager(config)
)

func TestE2E(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "subnet-evm warp e2e test")
}

// BeforeSuite starts the default network and adds 10 new nodes as validators with BLS keys
// registered on the P-Chain.
// Adds two disjoin sets of 5 of the new validator nodes to validate two new subnets with a
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

	// Construct the network using the avalanche-network-runner
	_, err = manager.StartDefaultNetwork(ctx)
	gomega.Expect(err).Should(gomega.BeNil())
	err = manager.SetupNetwork(
		ctx,
		config.AvalancheGoExecPath,
		[]*rpcpb.BlockchainSpec{
			{
				VmName:       evm.IDStr,
				Genesis:      "./tests/precompile/genesis/warp.json",
				ChainConfig:  "",
				SubnetConfig: "",
				Participants: subnetANodeNames,
			},
			{
				VmName:       evm.IDStr,
				Genesis:      "./tests/precompile/genesis/warp.json",
				ChainConfig:  "",
				SubnetConfig: "",
				Participants: subnetBNodeNames,
			},
		},
	)
	gomega.Expect(err).Should(gomega.BeNil())
})

var _ = ginkgo.AfterSuite(func() {
	gomega.Expect(manager).ShouldNot(gomega.BeNil())
	gomega.Expect(manager.TeardownNetwork()).Should(gomega.BeNil())
	// TODO: bootstrap an additional node to ensure that we can bootstrap the test data correctly
})

var _ = ginkgo.Describe("[Warp]", ginkgo.Ordered, func() {
	var (
		unsignedWarpMsg                *avalancheWarp.UnsignedMessage
		unsignedWarpMessageID          ids.ID
		signedWarpMsg                  *avalancheWarp.Message
		blockchainIDA, blockchainIDB   ids.ID
		chainAURIs, chainBURIs         []string
		chainAWSClient, chainBWSClient ethclient.Client
		chainID                        = big.NewInt(99999)
		fundedKey                      *ecdsa.PrivateKey
		fundedAddress                  common.Address
		payload                        = []byte{1, 2, 3}
		txSigner                       = types.LatestSignerForChainID(chainID)
	)

	var err error
	fundedKey, err = crypto.HexToECDSA("56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027")
	if err != nil {
		panic(err)
	}
	fundedAddress = crypto.PubkeyToAddress(fundedKey.PublicKey)

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

		log.Info("Created URIs for both subnets", "ChainAURIs", chainAURIs, "ChainBURIs", chainBURIs, "blockchainIDA", blockchainIDA.String(), "blockchainIDB", blockchainIDB)

		chainAWSURI := utils.ToWebsocketURI(chainAURIs[0], blockchainIDA.String())
		log.Info("Creating ethclient for blockchainA", "wsURI", chainAWSURI)
		chainAWSClient, err = ethclient.Dial(chainAWSURI)
		gomega.Expect(err).Should(gomega.BeNil())

		chainBWSURI := utils.ToWebsocketURI(chainBURIs[0], blockchainIDB.String())
		log.Info("Creating ethclient for blockchainB", "wsURI", chainBWSURI)
		chainBWSClient, err = ethclient.Dial(chainBWSURI)
		gomega.Expect(err).Should(gomega.BeNil())

	})

	// Send a transaction to Subnet A to issue a Warp Message to Subnet B
	ginkgo.It("Send Message from A to B", ginkgo.Label("Warp", "SendWarp"), func() {
		ctx := context.Background()

		gomega.Expect(err).Should(gomega.BeNil())

		log.Info("Subscribing to new heads")
		newHeads := make(chan *types.Header, 10)
		sub, err := chainAWSClient.SubscribeNewHead(ctx, newHeads)
		gomega.Expect(err).Should(gomega.BeNil())
		defer sub.Unsubscribe()

		// TODO: switch to using SubscribeFilterLogs
		// Create a channel to handle accepted warp logs
		// warpLogsChan := make(chan types.Log, 10)
		// warpLogsSub, err := chainBWSClient.SubscribeFilterLogs(
		// 	ctx,
		// 	interfaces.FilterQuery{Addresses: []common.Address{warp.Module.Address}},
		// 	warpLogsChan,
		// )
		// gomega.Expect(err).Should(gomega.BeNil())
		// defer warpLogsSub.Unsubscribe()

		packedInput, err := warp.PackSendWarpMessage(warp.SendWarpMessageInput{
			DestinationChainID: blockchainIDB,
			DestinationAddress: fundedAddress.Hash(),
			Payload:            payload,
		})
		gomega.Expect(err).Should(gomega.BeNil())
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     0,
			To:        &warp.Module.Address,
			Gas:       200_000,
			GasFeeCap: big.NewInt(225 * params.GWei),
			GasTipCap: big.NewInt(params.GWei),
			Value:     common.Big0,
			Data:      packedInput,
		})
		signedTx, err := types.SignTx(tx, txSigner, fundedKey)
		gomega.Expect(err).Should(gomega.BeNil())
		log.Info("Sending sendWarpMessage transaction", "txHash", signedTx.Hash())
		err = chainAWSClient.SendTransaction(ctx, signedTx)
		gomega.Expect(err).Should(gomega.BeNil())

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
		// select {
		// case warpLog := <-warpLogsChan:
		// 	gomega.Expect(warpLog).Should(gomega.Equal(txLog))
		// case <-time.After(5 * time.Second):
		// 	gomega.Expect(fmt.Errorf("failed to fetch txLog (%v) via subscription", txLog)).Should(gomega.BeNil())
		// }
		// warpLogsSub.Unsubscribe()

		log.Info("Parsing logData as unsigned warp message")
		unsignedMsg, err := avalancheWarp.ParseUnsignedMessage(txLog.Data)
		gomega.Expect(err).Should(gomega.BeNil())

		// Set local variables for the duration of the test
		unsignedWarpMessageID = unsignedMsg.ID()
		unsignedWarpMsg = unsignedMsg
		log.Info("Parsed unsignedWarpMsg", "unsignedWarpMessageID", unsignedWarpMessageID, "unsignedWarpMessage", unsignedWarpMsg)
	})

	// Aggregate a Warp Signature by sending an API request to each node requesting its signature and manually
	// constructing a valid Avalanche Warp Message
	ginkgo.It("Aggregate Warp Signature via API", ginkgo.Label("Warp", "ReceiveWarp", "AggregateWarpManually"), func() {
		ctx := context.Background()

		blsSignatures := make([]*bls.Signature, 0, len(chainAURIs))
		for i, uri := range chainAURIs {
			warpClient, err := warpBackend.NewWarpClient(uri, blockchainIDA.String())
			gomega.Expect(err).Should(gomega.BeNil())
			log.Info("Fetching warp signature from node")
			rawSignatureBytes, err := warpClient.GetSignature(ctx, unsignedWarpMessageID)
			gomega.Expect(err).Should(gomega.BeNil())

			blsSignature, err := bls.SignatureFromBytes(rawSignatureBytes)
			gomega.Expect(err).Should(gomega.BeNil())

			infoClient := info.NewClient(uri)
			nodeID, blsSigner, err := infoClient.GetNodeID(ctx)
			gomega.Expect(err).Should(gomega.BeNil())

			blsSignatures = append(blsSignatures, blsSignature)

			blsPublicKey := blsSigner.Key()
			log.Info("Verifying BLS Signature from node", "nodeID", nodeID, "nodeIndex", i)
			gomega.Expect(bls.Verify(blsPublicKey, blsSignature, unsignedWarpMsg.Bytes())).Should(gomega.BeTrue())
		}

		blsAggregatedSignature, err := bls.AggregateSignatures(blsSignatures)
		gomega.Expect(err).Should(gomega.BeNil())

		signersBitSet := set.NewBits()
		for i := 0; i < len(blsSignatures); i++ {
			signersBitSet.Add(i)
		}
		warpSignature := &avalancheWarp.BitSetSignature{
			Signers: signersBitSet.Bytes(),
		}

		blsAggregatedSignatureBytes := bls.SignatureToBytes(blsAggregatedSignature)
		copy(warpSignature.Signature[:], blsAggregatedSignatureBytes)

		warpMsg, err := avalancheWarp.NewMessage(
			unsignedWarpMsg,
			warpSignature,
		)
		gomega.Expect(err).Should(gomega.BeNil())
		signedWarpMsg = warpMsg
	})

	// Aggregate a Warp Signature using the node's Signature Aggregation API call and verifying that its output matches the
	// the manual construction
	ginkgo.It("Aggregate Warp Signature via Aggregator", ginkgo.Label("Warp", "ReceiveWarp", "AggregatorWarp"), func() {
		ctx := context.Background()

		// Verify that the signature aggregation matches the results of manually constructing the warp message
		warpClient, err := warpBackend.NewWarpClient(chainAURIs[0], blockchainIDA.String())
		gomega.Expect(err).Should(gomega.BeNil())

		signedWarpMessageBytes, err := warpClient.GetAggregateSignature(ctx, unsignedWarpMessageID, 100)
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(signedWarpMessageBytes).Should(gomega.Equal(signedWarpMsg.Bytes()))
	})

	// Verify successful delivery of the Avalanche Warp Message from Chain A to Chain B
	ginkgo.It("Verify Message from A to B", ginkgo.Label("Warp", "VerifyMessage"), func() {
		ctx := context.Background()

		log.Info("Subscribing to new heads")
		newHeads := make(chan *types.Header, 10)
		sub, err := chainBWSClient.SubscribeNewHead(ctx, newHeads)
		gomega.Expect(err).Should(gomega.BeNil())
		defer sub.Unsubscribe()

		// Trigger building of a new block at the current timestamp.
		// This timestamp should be after the ProposerVM activation time or ApricotPhase4 block timestamp.
		// This should generate a PostForkBlock because its parent block (genesis) has a timestamp (0) that is greater than or equal
		// to the fork activation time of 0.
		// Therefore, when we build a subsequent block it should be built with BuildBlockWithContext
		triggerTx, err := types.SignTx(types.NewTransaction(0, fundedAddress, common.Big1, 21_000, big.NewInt(225*params.GWei), nil), txSigner, fundedKey)
		gomega.Expect(err).Should(gomega.BeNil())

		err = chainBWSClient.SendTransaction(ctx, triggerTx)
		gomega.Expect(err).Should(gomega.BeNil())
		newHead := <-newHeads
		log.Info("Transaction triggered new block", "blockHash", newHead.Hash())

		// Try building another block to see if that one ends up as a PostForkBlock
		triggerTx2, err := types.SignTx(types.NewTransaction(1, fundedAddress, common.Big1, 21_000, big.NewInt(225*params.GWei), nil), txSigner, fundedKey)
		gomega.Expect(err).Should(gomega.BeNil())

		err = chainBWSClient.SendTransaction(ctx, triggerTx2)
		gomega.Expect(err).Should(gomega.BeNil())
		newHead = <-newHeads
		log.Info("Transaction2 triggered new block", "blockHash", newHead.Hash())

		packedInput, err := warp.PackGetVerifiedWarpMessage()
		gomega.Expect(err).Should(gomega.BeNil())
		tx := warpTransaction.NewWarpTx(
			chainID,
			2,
			&warp.Module.Address,
			5_000_000,
			big.NewInt(225*params.GWei),
			big.NewInt(params.GWei),
			common.Big0,
			packedInput,
			types.AccessList{},
			signedWarpMsg,
		)
		signedTx, err := types.SignTx(tx, txSigner, fundedKey)
		gomega.Expect(err).Should(gomega.BeNil())
		txBytes, err := signedTx.MarshalBinary()
		gomega.Expect(err).Should(gomega.BeNil())
		log.Info("Sending getVerifiedWarpMessage transaction", "txHash", signedTx.Hash(), "txBytes", common.Bytes2Hex(txBytes))
		err = chainBWSClient.SendTransaction(ctx, signedTx)
		gomega.Expect(err).Should(gomega.BeNil())

		log.Info("Waiting for new block confirmation")
		newHead = <-newHeads
		blockHash := newHead.Hash()
		log.Info("Fetching relevant warp logs and receipts from new block")
		logs, err := chainBWSClient.FilterLogs(ctx, interfaces.FilterQuery{
			BlockHash: &blockHash,
			Addresses: []common.Address{warp.Module.Address},
		})
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(len(logs)).Should(gomega.Equal(0))
		receipt, err := chainBWSClient.TransactionReceipt(ctx, signedTx.Hash())
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(receipt.Status).Should(gomega.Equal(types.ReceiptStatusSuccessful))
	})
})
