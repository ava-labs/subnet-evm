package solidity

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"

	_ "embed"

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
	"github.com/ava-labs/subnet-evm/tests/utils/runner"
	predicateutils "github.com/ava-labs/subnet-evm/utils/predicate"
	warpBackend "github.com/ava-labs/subnet-evm/warp"
	"github.com/ava-labs/subnet-evm/x/warp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func toWebsocketURI(uri string, blockchainID string) string {
	return fmt.Sprintf("ws://%s/ext/bc/%s/ws", strings.TrimPrefix(uri, "http://"), blockchainID)
}

var (
	config              = runner.NewDefaultANRConfig()
	manager             = runner.NewNetworkManager(config)
	warpChainConfigPath string
	//go:embed warp.abi
	WarpExampleRawABI string
	//go:embed warp.bin
	WarpExampleRawBin string
)

// used for indirect tests
type e2eContext struct {
	chainWSClient ethclient.Client
	addr          common.Address
	name          string
}

func newe2eContext(chainWSClient ethclient.Client, name string) *e2eContext {
	return &e2eContext{
		chainWSClient: chainWSClient,
		name:          name,
	}
}

// global variables used across multiple e2e tests
// allows tests to reuse a network
type testContext struct {
	blockchainIDA, blockchainIDB   ids.ID
	chainAURIs, chainBURIs         []string
	chainAWSClient, chainBWSClient ethclient.Client
	chainID                        big.Int
	fundedKey                      *ecdsa.PrivateKey
	fundedAddress                  common.Address
	txSigner                       types.Signer
}

func newTestContext() *testContext {
	chainID := big.NewInt(99999)
	return &testContext{
		chainID:  *chainID,
		txSigner: types.LatestSignerForChainID(chainID),
	}
}

func TestE2E(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "subnet-evm precompile ginkgo test suite")
}

var _ = ginkgo.BeforeSuite(func() {
	ctx := context.Background()
	var err error
	config = runner.NewDefaultANRConfig()
	manager = runner.NewNetworkManager(config)
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
})

var _ = ginkgo.AfterSuite(func() {
	gomega.Expect(manager).ShouldNot(gomega.BeNil())
	gomega.Expect(manager.TeardownNetwork()).Should(gomega.BeNil())
	gomega.Expect(os.Remove(warpChainConfigPath)).Should(gomega.BeNil())
	// TODO: bootstrap an additional node to ensure that we can bootstrap the test data correctly
})

var _ = ginkgo.Describe("[Warp]", ginkgo.Ordered, func() {
	testCtx := newTestContext()

	var err error
	testCtx.fundedKey, err = crypto.HexToECDSA("56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027")
	if err != nil {
		panic(err)
	}
	testCtx.fundedAddress = crypto.PubkeyToAddress(testCtx.fundedKey.PublicKey)

	//the same URIs are setup for both tests
	ginkgo.It("Setup URIs", ginkgo.Label("Warp", "SetupWarp"), func() {
		var err error
		subnetIDs := manager.GetSubnets()
		gomega.Expect(len(subnetIDs)).Should(gomega.Equal(2))

		subnetA := subnetIDs[0]
		subnetADetails, ok := manager.GetSubnet(subnetA)
		gomega.Expect(ok).Should(gomega.BeTrue())
		testCtx.blockchainIDA = subnetADetails.BlockchainID
		gomega.Expect(len(subnetADetails.ValidatorURIs)).Should(gomega.Equal(5))
		testCtx.chainAURIs = append(testCtx.chainAURIs, subnetADetails.ValidatorURIs...)

		subnetB := subnetIDs[1]
		subnetBDetails, ok := manager.GetSubnet(subnetB)
		gomega.Expect(ok).Should(gomega.BeTrue())
		testCtx.blockchainIDB = subnetBDetails.BlockchainID
		gomega.Expect(len(subnetBDetails.ValidatorURIs)).Should(gomega.Equal(5))
		testCtx.chainBURIs = append(testCtx.chainBURIs, subnetBDetails.ValidatorURIs...)

		log.Info("Created URIs for both subnets", "ChainAURIs", testCtx.chainAURIs, "ChainBURIs", testCtx.chainBURIs, "blockchainIDA", testCtx.blockchainIDA.String(), "blockchainIDB", testCtx.blockchainIDB)

		chainAWSURI := toWebsocketURI(testCtx.chainAURIs[0], testCtx.blockchainIDA.String())
		log.Info("Creating ethclient for blockchainA", "wsURI", chainAWSURI)
		testCtx.chainAWSClient, err = ethclient.Dial(chainAWSURI)
		gomega.Expect(err).Should(gomega.BeNil())

		chainBWSURI := toWebsocketURI(testCtx.chainBURIs[0], testCtx.blockchainIDB.String())
		log.Info("Creating ethclient for blockchainB", "wsURI", chainBWSURI)
		testCtx.chainBWSClient, err = ethclient.Dial(chainBWSURI)
		gomega.Expect(err).Should(gomega.BeNil())
	})

	// test that interacts with the warp precompile directly
	ginkgo.Describe("[Warp_Direct]", ginkgo.Ordered, func() {
		var (
			payload               = []byte("warp_direct")
			unsignedWarpMsg       *avalancheWarp.UnsignedMessage
			unsignedWarpMessageID ids.ID
			signedWarpMsg         *avalancheWarp.Message
		)

		ginkgo.It("Send Message from A to B", ginkgo.Label("Warp", "SendWarp"), func() {
			unsignedWarpMsg, unsignedWarpMessageID = sendMessage(testCtx, payload, warp.Module.Address, testCtx.fundedAddress)
		})

		// Aggregate a Warp Signature by sending an API request to each node requesting its signature and manually
		// constructing a valid Avalanche Warp Message
		ginkgo.It("Aggregate Warp Signature via API", ginkgo.Label("Warp", "ReceiveWarp", "AggregateWarpManually"), func() {
			signedWarpMsg = aggregateWarpSigsApi(testCtx, *unsignedWarpMsg, unsignedWarpMessageID)
		})

		// Aggregate a Warp Signature using the node's Signature Aggregation API call and verifying that its output matches the
		// the manual construction
		ginkgo.It("Aggregate Warp Signature via Aggregator", ginkgo.Label("Warp", "ReceiveWarp", "AggregatorWarp"), func() {
			aggregateWarpSigsAggreagtor(testCtx, *unsignedWarpMsg, unsignedWarpMessageID, *signedWarpMsg)
		})

		// Verify successful delivery of the Avalanche Warp Message from Chain A to Chain B
		ginkgo.It("Verify Message from A to B", ginkgo.Label("Warp", "VerifyMessage"), func() {
			packedInput, err := warp.PackGetVerifiedWarpMessage()
			gomega.Expect(err).Should(gomega.BeNil())

			verifyWarpMsg(testCtx, packedInput, warp.Module.Address, signedWarpMsg)
		})
	})

	// Test that interacts with the warp precompile via an intermediary contract
	ginkgo.Describe("Warp_Indirect", ginkgo.Ordered, func() {
		var (
			WarpExampleABI                    = contract.ParseABI(WarpExampleRawABI)
			WarpExampleBin             []byte = common.FromHex(WarpExampleRawBin)
			chainAExAddr, chainBExAddr common.Address
			unsignedWarpMsg            *avalancheWarp.UnsignedMessage
			unsignedWarpMessageID      ids.ID
			signedWarpMsg              *avalancheWarp.Message
			sendToAddress              = common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87")
			payload                    = []byte("warp_indirect")
		)
		// Deploy the test warp contracts to subnets A & B
		ginkgo.It("Deploy Test Warp Contracts", ginkgo.Label("Warp", "DeployWarp"), func() {
			ctx := context.Background()
			auth, err := bind.NewKeyedTransactorWithChainID(testCtx.fundedKey, &testCtx.chainID)
			gomega.Expect(err).Should(gomega.BeNil())

			AContext := newe2eContext(testCtx.chainAWSClient, "A")
			BContext := newe2eContext(testCtx.chainBWSClient, "B")
			e2eCtxs := []*e2eContext{AContext, BContext}

			for i := 0; i < len(e2eCtxs); i++ {
				e2eCtx := e2eCtxs[i]

				// deploy on subnetA
				newHeads := make(chan *types.Header, 1)
				sub, err := e2eCtx.chainWSClient.SubscribeNewHead(ctx, newHeads)
				gomega.Expect(err).Should(gomega.BeNil())

				addr, tx, _, err := bind.DeployContract(auth, WarpExampleABI, WarpExampleBin, e2eCtx.chainWSClient)
				e2eCtx.addr = addr
				gomega.Expect(err).Should(gomega.BeNil())
				signed, err := auth.Signer(auth.From, tx)
				gomega.Expect(err).Should(gomega.BeNil())
				log.Info(fmt.Sprintf("Waiting for contract creation on chain %s to be accepted", e2eCtx.name))
				<-newHeads
				txRecp, err := e2eCtx.chainWSClient.TransactionReceipt(ctx, signed.Hash())
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(txRecp.Status).Should(gomega.Equal(types.ReceiptStatusSuccessful)) // make sure status code is 1, contract deployed successfully
				sub.Unsubscribe()
			}
			chainAExAddr = AContext.addr
			chainBExAddr = BContext.addr
		})

		ginkgo.It("Send Message from A to B", ginkgo.Label("Warp", "SendWarp"), func() {
			unsignedWarpMsg, unsignedWarpMessageID = sendMessage(testCtx, payload, chainAExAddr, sendToAddress)
		})

		// Aggregate a Warp Signature by sending an API request to each node requesting its signature and manually
		// constructing a valid Avalanche Warp Message
		ginkgo.It("Aggregate Warp Signature via API", ginkgo.Label("Warp", "ReceiveWarp", "AggregateWarpManually"), func() {
			signedWarpMsg = aggregateWarpSigsApi(testCtx, *unsignedWarpMsg, unsignedWarpMessageID)
		})

		// Aggregate a Warp Signature using the node's Signature Aggregation API call and verifying that its output matches the
		// the manual construction
		ginkgo.It("Aggregate Warp Signature via Aggregator", ginkgo.Label("Warp", "ReceiveWarp", "AggregatorWarp"), func() {
			aggregateWarpSigsAggreagtor(testCtx, *unsignedWarpMsg, unsignedWarpMessageID, *signedWarpMsg)
		})

		// Verify successful delivery of the Avalanche Warp Message from Chain A to Chain B
		ginkgo.It("Verify Message from A to B", ginkgo.Label("Warp", "VerifyMessage"), func() {
			packedInput, err := WarpExampleABI.Pack(
				"validateWarpMessage",
				testCtx.blockchainIDA,
				chainAExAddr.Hash(), // calling address is the deployed exampleWarpTx, not the original funded address
				testCtx.blockchainIDB,
				sendToAddress.Hash(),
				payload,
			)
			gomega.Expect(err).Should(gomega.BeNil())
			verifyWarpMsg(testCtx, packedInput, chainBExAddr, signedWarpMsg)
		})
	})
})

func sendMessage(
	testCtx *testContext,
	payload []byte,
	contractAddress common.Address,
	destinationAddress common.Address,
) (*avalancheWarp.UnsignedMessage, ids.ID) {
	ctx := context.Background()
	log.Info("Subscribing to new heads")
	newHeads := make(chan *types.Header, 10)
	sub, err := testCtx.chainAWSClient.SubscribeNewHead(ctx, newHeads)
	gomega.Expect(err).Should(gomega.BeNil())
	defer sub.Unsubscribe()

	packedInput, err := warp.PackSendWarpMessage(warp.SendWarpMessageInput{
		DestinationChainID: common.BytesToHash(testCtx.blockchainIDB[:]),
		DestinationAddress: destinationAddress,
		Payload:            payload,
	})
	gomega.Expect(err).Should(gomega.BeNil())
	nonce, err := testCtx.chainAWSClient.NonceAt(ctx, testCtx.fundedAddress, nil)
	gomega.Expect(err).Should(gomega.BeNil())
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   &testCtx.chainID,
		Nonce:     nonce,
		To:        &contractAddress,
		Gas:       200_000,
		GasFeeCap: big.NewInt(225 * params.GWei),
		GasTipCap: big.NewInt(params.GWei),
		Value:     common.Big0,
		Data:      packedInput,
	})
	signedTx, err := types.SignTx(tx, testCtx.txSigner, testCtx.fundedKey)
	gomega.Expect(err).Should(gomega.BeNil())
	log.Info("Sending sendWarpMessage transaction", "txHash", signedTx.Hash())
	err = testCtx.chainAWSClient.SendTransaction(ctx, signedTx)
	gomega.Expect(err).Should(gomega.BeNil())

	log.Info("Waiting for new block confirmation")
	newHead := <-newHeads
	blockHash := newHead.Hash()

	log.Info("Fetching relevant warp logs from the newly produced block")
	logs, err := testCtx.chainAWSClient.FilterLogs(ctx, interfaces.FilterQuery{
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
	unsignedWarpMessageID := unsignedMsg.ID()
	unsignedWarpMsg := unsignedMsg
	log.Info("Parsed unsignedWarpMsg", "unsignedWarpMessageID", unsignedWarpMessageID, "unsignedWarpMessage", unsignedWarpMsg)

	// Loop over each client on chain A to ensure they all have time to accept the block.
	// Note: if we did not confirm this here, the next stage could be racy since it assumes every node
	// has accepted the block.
	for i, uri := range testCtx.chainAURIs {
		chainAWSURI := toWebsocketURI(uri, testCtx.blockchainIDA.String())
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
	return unsignedWarpMsg, unsignedWarpMessageID
}

// Aggregate the warp signatures via direct API calls
func aggregateWarpSigsApi(
	testCtx *testContext,
	unsignedWarpMsg avalancheWarp.UnsignedMessage,
	unsignedWarpMessageID ids.ID,
) *avalancheWarp.Message {
	ctx := context.Background()

	blsSignatures := make([]*bls.Signature, 0, len(testCtx.chainAURIs))
	for i, uri := range testCtx.chainAURIs {
		warpClient, err := warpBackend.NewWarpClient(uri, testCtx.blockchainIDA.String())
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
		&unsignedWarpMsg,
		warpSignature,
	)
	gomega.Expect(err).Should(gomega.BeNil())
	return warpMsg
}

// Aggregate the warp sigs via the aggregator, check that they match
func aggregateWarpSigsAggreagtor(
	testCtx *testContext,
	unsignedWarpMsg avalancheWarp.UnsignedMessage,
	unsignedWarpMessageID ids.ID,
	signedWarpMsg avalancheWarp.Message,
) {
	ctx := context.Background()

	// Verify that the signature aggregation matches the results of manually constructing the warp message
	warpClient, err := warpBackend.NewWarpClient(testCtx.chainAURIs[0], testCtx.blockchainIDA.String())
	gomega.Expect(err).Should(gomega.BeNil())

	signedWarpMessageBytes, err := warpClient.GetAggregateSignature(ctx, unsignedWarpMessageID, 100)
	gomega.Expect(err).Should(gomega.BeNil())
	gomega.Expect(signedWarpMessageBytes).Should(gomega.Equal(signedWarpMsg.Bytes()))
}

func verifyWarpMsg(
	testCtx *testContext,
	packedInput []byte,
	contractAddress common.Address,
	signedWarpMsg *avalancheWarp.Message,
) {
	ctx := context.Background()

	log.Info("Subscribing to new heads")
	newHeads := make(chan *types.Header, 10)
	sub, err := testCtx.chainBWSClient.SubscribeNewHead(ctx, newHeads)
	gomega.Expect(err).Should(gomega.BeNil())
	defer sub.Unsubscribe()

	// Trigger building of a new block at the current timestamp.
	// This timestamp should be after the ProposerVM activation time or ApricotPhase4 block timestamp.
	// This should generate a PostForkBlock because its parent block (genesis) has a timestamp (0) that is greater than or equal
	// to the fork activation time of 0.
	// Therefore, when we build a subsequent block it should be built with BuildBlockWithContext
	nonce, err := testCtx.chainBWSClient.NonceAt(ctx, testCtx.fundedAddress, nil)
	gomega.Expect(err).Should(gomega.BeNil())

	triggerTx, err := types.SignTx(types.NewTransaction(nonce, testCtx.fundedAddress, common.Big1, 21_000, big.NewInt(225*params.GWei), nil), testCtx.txSigner, testCtx.fundedKey)
	gomega.Expect(err).Should(gomega.BeNil())

	err = testCtx.chainBWSClient.SendTransaction(ctx, triggerTx)
	gomega.Expect(err).Should(gomega.BeNil())
	newHead := <-newHeads
	log.Info("Transaction triggered new block", "blockHash", newHead.Hash())

	nonce++
	// Try building another block to see if that one ends up as a PostForkBlock
	triggerTx2, err := types.SignTx(types.NewTransaction(nonce, testCtx.fundedAddress, common.Big1, 21_000, big.NewInt(225*params.GWei), nil), testCtx.txSigner, testCtx.fundedKey)
	gomega.Expect(err).Should(gomega.BeNil())

	err = testCtx.chainBWSClient.SendTransaction(ctx, triggerTx2)
	gomega.Expect(err).Should(gomega.BeNil())
	newHead = <-newHeads
	log.Info("Transaction2 triggered new block", "blockHash", newHead.Hash())
	nonce++
	tx := predicateutils.NewPredicateTx(
		&testCtx.chainID,
		nonce,
		&contractAddress,
		5_000_000,
		big.NewInt(225*params.GWei),
		big.NewInt(params.GWei),
		common.Big0,
		packedInput,
		types.AccessList{},
		warp.ContractAddress,
		signedWarpMsg.Bytes(),
	)
	signedTx, err := types.SignTx(tx, testCtx.txSigner, testCtx.fundedKey)
	gomega.Expect(err).Should(gomega.BeNil())
	txBytes, err := signedTx.MarshalBinary()
	gomega.Expect(err).Should(gomega.BeNil())
	log.Info("Sending getVerifiedWarpMessage transaction", "txHash", signedTx.Hash(), "txBytes", common.Bytes2Hex(txBytes))
	err = testCtx.chainBWSClient.SendTransaction(ctx, signedTx)
	gomega.Expect(err).Should(gomega.BeNil())

	log.Info("Waiting for new block confirmation")
	newHead = <-newHeads
	blockHash := newHead.Hash()
	log.Info("Fetching relevant warp logs and receipts from new block")
	logs, err := testCtx.chainBWSClient.FilterLogs(ctx, interfaces.FilterQuery{
		BlockHash: &blockHash,
		Addresses: []common.Address{warp.Module.Address},
	})
	gomega.Expect(err).Should(gomega.BeNil())
	gomega.Expect(len(logs)).Should(gomega.Equal(0))
	receipt, err := testCtx.chainBWSClient.TransactionReceipt(ctx, signedTx.Hash())
	gomega.Expect(err).Should(gomega.BeNil())
	gomega.Expect(receipt.Status).Should(gomega.Equal(types.ReceiptStatusSuccessful))
}
