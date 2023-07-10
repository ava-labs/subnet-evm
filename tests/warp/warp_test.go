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
	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/interfaces"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/plugin/evm"
	"github.com/ava-labs/subnet-evm/precompile/contract"
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

var (
	WarpExampleRawABI string = `[{"inputs":[{"internalType":"bytes32","name":"destinationChainID","type":"bytes32"},{"internalType":"bytes32","name":"destinationAddress","type":"bytes32"},{"internalType":"bytes","name":"payload","type":"bytes"}],"name":"sendWarpMessage","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"expectedBlockchainID","type":"bytes32"}],"name":"validateGetBlockchainID","outputs":[],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"originChainID","type":"bytes32"},{"internalType":"bytes32","name":"originSenderAddress","type":"bytes32"},{"internalType":"bytes32","name":"destinationChainID","type":"bytes32"},{"internalType":"bytes32","name":"destinationAddress","type":"bytes32"},{"internalType":"bytes","name":"payload","type":"bytes"}],"name":"validateWarpMessage","outputs":[],"stateMutability":"view","type":"function"}]`
	WarpExampleABI           = contract.ParseABI(WarpExampleRawABI)
	WarpExampleBin    []byte = common.FromHex("60806040527302000000000000000000000000000000000000056000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555034801561006457600080fd5b506108db806100746000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806315f0c959146100465780631d9a3246146100625780636a0b64dd1461007e575b600080fd5b610060600480360381019061005b919061033c565b61009a565b005b61007c600480360381019061007791906103ce565b61013e565b005b61009860048036038101906100939190610468565b61025b565b005b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16634213cf786040518163ffffffff1660e01b8152600401602060405180830381865afa158015610108573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061012c91906104f1565b905081811461013a57600080fd5b5050565b60008060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16631f9b40ec6040518163ffffffff1660e01b8152600401600060405180830381865afa1580156101ac573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f820116820180604052508101906101d59190610764565b91509150806101e357600080fd5b878260000151146101f357600080fd5b8682602001511461020357600080fd5b8582604001511461021357600080fd5b8482606001511461022357600080fd5b83836040516102339291906107ff565b60405180910390208260800151805190602001201461025157600080fd5b5050505050505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16636a0b64dd858585856040518563ffffffff1660e01b81526004016102ba9493929190610865565b600060405180830381600087803b1580156102d457600080fd5b505af11580156102e8573d6000803e3d6000fd5b5050505050505050565b6000604051905090565b600080fd5b600080fd5b6000819050919050565b61031981610306565b811461032457600080fd5b50565b60008135905061033681610310565b92915050565b600060208284031215610352576103516102fc565b5b600061036084828501610327565b91505092915050565b600080fd5b600080fd5b600080fd5b60008083601f84011261038e5761038d610369565b5b8235905067ffffffffffffffff8111156103ab576103aa61036e565b5b6020830191508360018202830111156103c7576103c6610373565b5b9250929050565b60008060008060008060a087890312156103eb576103ea6102fc565b5b60006103f989828a01610327565b965050602061040a89828a01610327565b955050604061041b89828a01610327565b945050606061042c89828a01610327565b935050608087013567ffffffffffffffff81111561044d5761044c610301565b5b61045989828a01610378565b92509250509295509295509295565b60008060008060608587031215610482576104816102fc565b5b600061049087828801610327565b94505060206104a187828801610327565b935050604085013567ffffffffffffffff8111156104c2576104c1610301565b5b6104ce87828801610378565b925092505092959194509250565b6000815190506104eb81610310565b92915050565b600060208284031215610507576105066102fc565b5b6000610515848285016104dc565b91505092915050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61056c82610523565b810181811067ffffffffffffffff8211171561058b5761058a610534565b5b80604052505050565b600061059e6102f2565b90506105aa8282610563565b919050565b600080fd5b600080fd5b600067ffffffffffffffff8211156105d4576105d3610534565b5b6105dd82610523565b9050602081019050919050565b60005b838110156106085780820151818401526020810190506105ed565b60008484015250505050565b6000610627610622846105b9565b610594565b905082815260208101848484011115610643576106426105b4565b5b61064e8482856105ea565b509392505050565b600082601f83011261066b5761066a610369565b5b815161067b848260208601610614565b91505092915050565b600060a0828403121561069a5761069961051e565b5b6106a460a0610594565b905060006106b4848285016104dc565b60008301525060206106c8848285016104dc565b60208301525060406106dc848285016104dc565b60408301525060606106f0848285016104dc565b606083015250608082015167ffffffffffffffff811115610714576107136105af565b5b61072084828501610656565b60808301525092915050565b60008115159050919050565b6107418161072c565b811461074c57600080fd5b50565b60008151905061075e81610738565b92915050565b6000806040838503121561077b5761077a6102fc565b5b600083015167ffffffffffffffff81111561079957610798610301565b5b6107a585828601610684565b92505060206107b68582860161074f565b9150509250929050565b600081905092915050565b82818337600083830152505050565b60006107e683856107c0565b93506107f38385846107cb565b82840190509392505050565b600061080c8284866107da565b91508190509392505050565b61082181610306565b82525050565b600082825260208201905092915050565b60006108448385610827565b93506108518385846107cb565b61085a83610523565b840190509392505050565b600060608201905061087a6000830187610818565b6108876020830186610818565b818103604083015261089a818486610838565b90509594505050505056fea2646970667358221220b2cefd832fa6ca66aff8dba2e446584042974aa24a3f54e0f99e83fc9cf34edf64736f6c63430008140033")
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
				VmName:      evm.IDStr,
				Genesis:     "./tests/precompile/genesis/warp.json",
				ChainConfig: "",
				SubnetSpec: &rpcpb.SubnetSpec{
					SubnetConfig: "",
					Participants: subnetANodeNames,
				},
			},
			{
				VmName:      evm.IDStr,
				Genesis:     "./tests/precompile/genesis/warp.json",
				ChainConfig: "",
				SubnetSpec: &rpcpb.SubnetSpec{
					SubnetConfig: "",
					Participants: subnetBNodeNames,
				},
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
		chainANonce, chainBNonce       = uint64(0), uint64(0)
		chainAExAddr, chainBExAddr     common.Address
		chainID                        = big.NewInt(99999)
		fundedKey                      *ecdsa.PrivateKey
		fundedAddress                  common.Address
		payload                        = []byte{1, 2, 3}
		sendToAddress                  common.Address
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

	// Deploy the test warp contracts to subnets A & B
	ginkgo.It("Deploy Test Warp Contracts", ginkgo.Label("Warp", "DeployWarp"), func() {
		ctx := context.Background()

		auth, err := bind.NewKeyedTransactorWithChainID(fundedKey, chainID)
		gomega.Expect(err).Should(gomega.BeNil())

		//deploy on subnetA
		newHeads := make(chan *types.Header, 10)
		sub, err := chainAWSClient.SubscribeNewHead(ctx, newHeads)
		gomega.Expect(err).Should(gomega.BeNil())

		addr, tx, _, err := bind.DeployContract(auth, WarpExampleABI, WarpExampleBin, chainAWSClient)
		gomega.Expect(err).Should(gomega.BeNil())
		chainANonce++
		chainAExAddr = addr
		signed, err := auth.Signer(auth.From, tx)
		gomega.Expect(err).Should(gomega.BeNil())

		log.Info("Waiting for contract creation on chain A to be accepted")
		<-newHeads
		recp, err := chainAWSClient.TransactionReceipt(ctx, signed.Hash())
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(recp.Status).Should(gomega.Equal(types.ReceiptStatusSuccessful)) //make sure status code is 1, contract deployed successfully
		sub.Unsubscribe()

		//deploy on subnetB
		newHeads = make(chan *types.Header, 10)
		sub, err = chainBWSClient.SubscribeNewHead(ctx, newHeads)
		gomega.Expect(err).Should(gomega.BeNil())

		addr, tx, _, err = bind.DeployContract(auth, WarpExampleABI, WarpExampleBin, chainBWSClient)
		gomega.Expect(err).Should(gomega.BeNil())
		chainBNonce++
		chainBExAddr = addr
		signed, err = auth.Signer(auth.From, tx)
		gomega.Expect(err).Should(gomega.BeNil())

		log.Info("waiting for contract creation on chain B to be accepted")
		<-newHeads
		recp, err = chainBWSClient.TransactionReceipt(ctx, signed.Hash())
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(recp.Status).Should(gomega.Equal(types.ReceiptStatusSuccessful)) //make sure status code is 1, contract deployed successfully
		sub.Unsubscribe()
		log.Info("contracts deployed")
	})

	// Send a transaction to Subnet A to issue a Warp Message to Subnet B
	ginkgo.It("Send Message from A to B", ginkgo.Label("Warp", "SendWarp"), func() {
		ctx := context.Background()

		gomega.Expect(err).Should(gomega.BeNil())

		// TODO: switch to using SubscribeFilterLogs to retrieve the warp log
		log.Info("Subscribing to new heads")
		newHeads := make(chan *types.Header, 10)
		sub, err := chainAWSClient.SubscribeNewHead(ctx, newHeads)
		gomega.Expect(err).Should(gomega.BeNil())
		defer sub.Unsubscribe()

		key2, err := crypto.GenerateKey()
		gomega.Expect(err).Should(gomega.BeNil())
		sendToAddress = crypto.PubkeyToAddress(key2.PublicKey)

		packedInput, err := warp.PackSendWarpMessage(warp.SendWarpMessageInput{
			DestinationChainID: blockchainIDB,
			DestinationAddress: sendToAddress.Hash(),
			Payload:            payload,
		})
		gomega.Expect(err).Should(gomega.BeNil())

		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     1,
			To:        &chainAExAddr,
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
		chainANonce++

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
		unsignedMsg, err := avalancheWarp.ParseUnsignedMessage(txLog.Data)
		gomega.Expect(err).Should(gomega.BeNil())

		// Set local variables for the duration of the test
		unsignedWarpMessageID = unsignedMsg.ID()
		unsignedWarpMsg = unsignedMsg
		log.Info("Parsed unsignedWarpMsg", "unsignedWarpMessageID", unsignedWarpMessageID, "unsignedWarpMessage", unsignedWarpMsg)

		// Loop over each client on chain A to ensure they all have time to accept the block.
		// Note: if we did not confirm this here, the next stage could be racy since it assumes every node
		// has accepted the block.
		for i, uri := range chainAURIs {
			chainAWSURI := utils.ToWebsocketURI(uri, blockchainIDA.String())
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
		nonce, err := chainBWSClient.NonceAt(ctx, fundedAddress, nil)
		gomega.Expect(err).Should(gomega.BeNil())

		triggerTx, err := types.SignTx(types.NewTransaction(nonce, fundedAddress, common.Big1, 21_000, big.NewInt(225*params.GWei), nil), txSigner, fundedKey)
		gomega.Expect(err).Should(gomega.BeNil())

		err = chainBWSClient.SendTransaction(ctx, triggerTx)
		gomega.Expect(err).Should(gomega.BeNil())
		newHead := <-newHeads
		log.Info("Transaction triggered new block", "blockHash", newHead.Hash())

		nonce++
		// Try building another block to see if that one ends up as a PostForkBlock
		triggerTx2, err := types.SignTx(types.NewTransaction(nonce, fundedAddress, common.Big1, 21_000, big.NewInt(225*params.GWei), nil), txSigner, fundedKey)
		gomega.Expect(err).Should(gomega.BeNil())

		err = chainBWSClient.SendTransaction(ctx, triggerTx2)
		gomega.Expect(err).Should(gomega.BeNil())
		newHead = <-newHeads
		log.Info("Transaction2 triggered new block", "blockHash", newHead.Hash())
		nonce++

		//attempt to send a warp message with no sigs, which should fail
		packedInput, err := WarpExampleABI.Pack(
			"validateWarpMessage",
			blockchainIDA,
			chainAExAddr.Hash(), //calling address is the deployed exampleWarpTx, not the original funded address
			blockchainIDB,
			sendToAddress.Hash(),
			payload,
		)
		gomega.Expect(err).Should(gomega.BeNil())

		tx := warpTransaction.NewWarpTx(
			chainID,
			nonce,
			&chainBExAddr,
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
