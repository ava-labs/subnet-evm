// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements solidity tests.
package warp

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	ginkgo "github.com/onsi/ginkgo/v2"

	"github.com/onsi/gomega"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/validators"
	"github.com/ava-labs/avalanchego/tests/fixture/e2e"
	"github.com/ava-labs/avalanchego/tests/fixture/tmpnet"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/vms/platformvm"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/payload"

	"github.com/ava-labs/subnet-evm/cmd/simulator/key"
	"github.com/ava-labs/subnet-evm/cmd/simulator/load"
	"github.com/ava-labs/subnet-evm/cmd/simulator/metrics"
	"github.com/ava-labs/subnet-evm/cmd/simulator/txs"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/interfaces"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contracts/warp"
	"github.com/ava-labs/subnet-evm/predicate"
	"github.com/ava-labs/subnet-evm/tests"
	"github.com/ava-labs/subnet-evm/tests/utils"
	warpBackend "github.com/ava-labs/subnet-evm/warp"
	"github.com/ava-labs/subnet-evm/warp/aggregator"
)

const (
	subnetAName = "warp-subnet-a"
	subnetBName = "warp-subnet-b"
)

var (
	flagVars *e2e.FlagVars

	repoRootPath = tests.GetRepoRootPath("tests/warp")

	genesisPath = filepath.Join(repoRootPath, "tests/precompile/genesis/warp.json")

	subnetA, subnetB, cChainSubnetDetails *Subnet

	testPayload = []byte{1, 2, 3}
)

func init() {
	// Configures flags used to configure tmpnet (via SynchronizedBeforeSuite)
	flagVars = e2e.RegisterFlags()
}

// Subnet provides the basic details of a created subnet
type Subnet struct {
	// SubnetID is the txID of the transaction that created the subnet
	SubnetID ids.ID
	// For simplicity assume a single blockchain per subnet
	BlockchainID ids.ID
	// Key funded in the genesis of the blockchain
	PreFundedKey *ecdsa.PrivateKey
	// ValidatorURIs are the base URIs for each participant of the Subnet
	ValidatorURIs []string
}

func TestE2E(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "subnet-evm warp e2e test")
}

var _ = ginkgo.SynchronizedBeforeSuite(func() []byte {
	// Run only once in the first ginkgo process
	env := e2e.NewTestEnvironment(
		flagVars,
		utils.NewTmpnetNetwork(
			tmpnet.DefaultNodeCount,
			utils.NewTmpnetSubnet(subnetAName, genesisPath),
			utils.NewTmpnetSubnet(subnetBName, genesisPath),
		),
	)

	return env.Marshal()
}, func(envBytes []byte) {
	// Run in every ginkgo process

	require := require.New(ginkgo.GinkgoT())

	// Initialize the local test environment from the global state
	if len(envBytes) > 0 {
		e2e.InitSharedTestEnvironment(envBytes)
	}

	network := e2e.Env.GetNetwork()

	// By default all nodes are validating all subnets
	validatorURIs := make([]string, len(network.Nodes))
	for i, node := range network.Nodes {
		validatorURIs[i] = node.URI
	}

	tmpnetSubnetA := network.GetSubnet(subnetAName)
	require.NotNil(tmpnetSubnetA)
	subnetA = &Subnet{
		SubnetID:      tmpnetSubnetA.SubnetID,
		BlockchainID:  tmpnetSubnetA.Chains[0].ChainID,
		PreFundedKey:  tmpnetSubnetA.Chains[0].PreFundedKey.ToECDSA(),
		ValidatorURIs: validatorURIs,
	}

	tmpnetSubnetB := network.GetSubnet(subnetBName)
	require.NotNil(tmpnetSubnetB)
	subnetB = &Subnet{
		SubnetID:      tmpnetSubnetB.SubnetID,
		BlockchainID:  tmpnetSubnetB.Chains[0].ChainID,
		PreFundedKey:  tmpnetSubnetB.Chains[0].PreFundedKey.ToECDSA(),
		ValidatorURIs: validatorURIs,
	}

	infoClient := info.NewClient(network.Nodes[0].URI)
	cChainBlockchainID, err := infoClient.GetBlockchainID(e2e.DefaultContext(), "C")
	require.NoError(err)

	cChainSubnetDetails = &Subnet{
		SubnetID:      constants.PrimaryNetworkID,
		BlockchainID:  cChainBlockchainID,
		PreFundedKey:  tmpnet.HardhatKey.ToECDSA(),
		ValidatorURIs: validatorURIs,
	}
})

var _ = ginkgo.Describe("[Warp]", func() {
	testFunc := func(sendingSubnet *Subnet, receivingSubnet *Subnet) {
		w := newWarpTest(e2e.DefaultContext(), sendingSubnet, receivingSubnet)

		log.Info("Sending message from A to B")
		w.sendMessageFromSendingSubnet()

		log.Info("Aggregating signatures via API")
		w.aggregateSignaturesViaAPI()

		log.Info("Aggregating signatures via p2p aggregator")
		w.aggregateSignatures()

		log.Info("Delivering addressed call payload to receiving subnet")
		w.deliverAddressedCallToReceivingSubnet()

		log.Info("Delivering block hash payload to receiving subnet")
		w.deliverBlockHashPayload()

		log.Info("Executing HardHat test")
		w.executeHardHatTest()

		log.Info("Executing warp load test")
		w.warpLoad()
	}
	ginkgo.It("SubnetA -> SubnetB", func() { testFunc(subnetA, subnetB) })
	ginkgo.It("SubnetA -> SubnetA", func() { testFunc(subnetA, subnetA) })
	ginkgo.It("SubnetA -> C-Chain", func() { testFunc(subnetA, cChainSubnetDetails) })
	ginkgo.It("C-Chain -> SubnetA", func() { testFunc(cChainSubnetDetails, subnetA) })
	ginkgo.It("C-Chain -> C-Chain", func() { testFunc(cChainSubnetDetails, cChainSubnetDetails) })
})

type warpTest struct {
	// network-wide fields set in the constructor
	networkID uint32

	// sendingSubnet fields set in the constructor
	sendingSubnet              *Subnet
	sendingSubnetURIs          []string
	sendingSubnetClients       []ethclient.Client
	sendingSubnetFundedKey     *ecdsa.PrivateKey
	sendingSubnetFundedAddress common.Address
	sendingSubnetChainID       *big.Int
	sendingSubnetSigner        types.Signer

	// receivingSubnet fields set in the constructor
	receivingSubnet              *Subnet
	receivingSubnetURIs          []string
	receivingSubnetClients       []ethclient.Client
	receivingSubnetFundedKey     *ecdsa.PrivateKey
	receivingSubnetFundedAddress common.Address
	receivingSubnetChainID       *big.Int
	receivingSubnetSigner        types.Signer

	// Fields set throughout test execution
	blockID                     ids.ID
	blockPayload                *payload.Hash
	blockPayloadUnsignedMessage *avalancheWarp.UnsignedMessage
	blockPayloadSignedMessage   *avalancheWarp.Message

	addressedCallUnsignedMessage *avalancheWarp.UnsignedMessage
	addressedCallSignedMessage   *avalancheWarp.Message
}

func newWarpTest(ctx context.Context, sendingSubnet *Subnet, receivingSubnet *Subnet) *warpTest {
	require := require.New(ginkgo.GinkgoT())

	sendingSubnetFundedKey := sendingSubnet.PreFundedKey
	receivingSubnetFundedKey := receivingSubnet.PreFundedKey

	warpTest := &warpTest{
		sendingSubnet:                sendingSubnet,
		sendingSubnetURIs:            sendingSubnet.ValidatorURIs,
		receivingSubnet:              receivingSubnet,
		receivingSubnetURIs:          receivingSubnet.ValidatorURIs,
		sendingSubnetFundedKey:       sendingSubnetFundedKey,
		sendingSubnetFundedAddress:   crypto.PubkeyToAddress(sendingSubnetFundedKey.PublicKey),
		receivingSubnetFundedKey:     receivingSubnetFundedKey,
		receivingSubnetFundedAddress: crypto.PubkeyToAddress(receivingSubnetFundedKey.PublicKey),
	}
	infoClient := info.NewClient(sendingSubnet.ValidatorURIs[0])
	networkID, err := infoClient.GetNetworkID(ctx)
	require.NoError(err)
	warpTest.networkID = networkID

	warpTest.initClients()

	sendingClient := warpTest.sendingSubnetClients[0]
	sendingSubnetChainID, err := sendingClient.ChainID(ctx)
	require.NoError(err)
	warpTest.sendingSubnetChainID = sendingSubnetChainID
	warpTest.sendingSubnetSigner = types.LatestSignerForChainID(sendingSubnetChainID)

	receivingClient := warpTest.receivingSubnetClients[0]
	receivingChainID, err := receivingClient.ChainID(ctx)
	require.NoError(err)
	// Issue transactions to activate ProposerVM on the receiving chain
	require.NoError(utils.IssueTxsToActivateProposerVMFork(ctx, receivingChainID, receivingSubnetFundedKey, receivingClient))
	warpTest.receivingSubnetChainID = receivingChainID
	warpTest.receivingSubnetSigner = types.LatestSignerForChainID(receivingChainID)

	return warpTest
}

func (w *warpTest) initClients() {
	require := require.New(ginkgo.GinkgoT())

	w.sendingSubnetClients = make([]ethclient.Client, 0, len(w.sendingSubnetClients))
	for _, uri := range w.sendingSubnet.ValidatorURIs {
		wsURI := toWebsocketURI(uri, w.sendingSubnet.BlockchainID.String())
		log.Info("Creating ethclient for blockchain A", "blockchainID", w.sendingSubnet.BlockchainID)
		client, err := ethclient.Dial(wsURI)
		require.NoError(err)
		w.sendingSubnetClients = append(w.sendingSubnetClients, client)
	}

	w.receivingSubnetClients = make([]ethclient.Client, 0, len(w.receivingSubnetClients))
	for _, uri := range w.receivingSubnet.ValidatorURIs {
		wsURI := toWebsocketURI(uri, w.receivingSubnet.BlockchainID.String())
		log.Info("Creating ethclient for blockchain B", "blockchainID", w.receivingSubnet.BlockchainID)
		client, err := ethclient.Dial(wsURI)
		require.NoError(err)
		w.receivingSubnetClients = append(w.receivingSubnetClients, client)
	}
}

func (w *warpTest) getBlockHashAndNumberFromTxReceipt(ctx context.Context, client ethclient.Client, tx *types.Transaction) (common.Hash, uint64) {
	// This uses the Subnet-EVM client to fetch a block from Coreth (when testing the C-Chain), so we use this
	// workaround to get the correct block hash. Note the client recalculates the block hash locally, which results
	// in a different block hash due to small differences in the block format.
	require := require.New(ginkgo.GinkgoT())
	for {
		require.NoError(ctx.Err())
		receipt, err := client.TransactionReceipt(ctx, tx.Hash())
		if err == nil {
			return receipt.BlockHash, receipt.BlockNumber.Uint64()
		}
	}
}

func (w *warpTest) sendMessageFromSendingSubnet() {
	ctx := e2e.DefaultContext()
	require := require.New(ginkgo.GinkgoT())

	client := w.sendingSubnetClients[0]
	log.Info("Subscribing to new heads")
	newHeads := make(chan *types.Header, 10)
	sub, err := client.SubscribeNewHead(ctx, newHeads)
	require.NoError(err)
	defer sub.Unsubscribe()

	startingNonce, err := client.NonceAt(ctx, w.sendingSubnetFundedAddress, nil)
	require.NoError(err)

	packedInput, err := warp.PackSendWarpMessage(testPayload)
	require.NoError(err)
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   w.sendingSubnetChainID,
		Nonce:     startingNonce,
		To:        &warp.Module.Address,
		Gas:       200_000,
		GasFeeCap: big.NewInt(225 * params.GWei),
		GasTipCap: big.NewInt(params.GWei),
		Value:     common.Big0,
		Data:      packedInput,
	})
	signedTx, err := types.SignTx(tx, w.sendingSubnetSigner, w.sendingSubnetFundedKey)
	require.NoError(err)
	log.Info("Sending sendWarpMessage transaction", "txHash", signedTx.Hash())
	err = client.SendTransaction(ctx, signedTx)
	require.NoError(err)

	log.Info("Waiting for new block confirmation")
	<-newHeads
	receiptCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	blockHash, blockNumber := w.getBlockHashAndNumberFromTxReceipt(receiptCtx, client, signedTx)

	log.Info("Constructing warp block hash unsigned message", "blockHash", blockHash)
	w.blockID = ids.ID(blockHash) // Set blockID to construct a warp message containing a block hash payload later
	w.blockPayload, err = payload.NewHash(w.blockID)
	require.NoError(err)
	w.blockPayloadUnsignedMessage, err = avalancheWarp.NewUnsignedMessage(w.networkID, w.sendingSubnet.BlockchainID, w.blockPayload.Bytes())
	require.NoError(err)

	log.Info("Fetching relevant warp logs from the newly produced block")
	logs, err := client.FilterLogs(ctx, interfaces.FilterQuery{
		BlockHash: &blockHash,
		Addresses: []common.Address{warp.Module.Address},
	})
	require.NoError(err)
	require.Len(logs, 1)

	// Check for relevant warp log from subscription and ensure that it matches
	// the log extracted from the last block.
	txLog := logs[0]
	log.Info("Parsing logData as unsigned warp message")
	unsignedMsg, err := warp.UnpackSendWarpEventDataToMessage(txLog.Data)
	require.NoError(err)

	// Set local variables for the duration of the test
	w.addressedCallUnsignedMessage = unsignedMsg
	log.Info("Parsed unsignedWarpMsg", "unsignedWarpMessageID", w.addressedCallUnsignedMessage.ID(), "unsignedWarpMessage", w.addressedCallUnsignedMessage)

	// Loop over each client on chain A to ensure they all have time to accept the block.
	// Note: if we did not confirm this here, the next stage could be racy since it assumes every node
	// has accepted the block.
	for i, client := range w.sendingSubnetClients {
		// Loop until each node has advanced to >= the height of the block that emitted the warp log
		for {
			block, err := client.BlockByNumber(ctx, nil)
			require.NoError(err)
			if block.NumberU64() >= blockNumber {
				log.Info("client accepted the block containing SendWarpMessage", "client", i, "height", block.NumberU64())
				break
			}
		}
	}
}

func (w *warpTest) aggregateSignaturesViaAPI() {
	require := require.New(ginkgo.GinkgoT())
	ctx := e2e.DefaultContext()

	warpAPIs := make(map[ids.NodeID]warpBackend.Client, len(w.sendingSubnetURIs))
	for _, uri := range w.sendingSubnetURIs {
		client, err := warpBackend.NewClient(uri, w.sendingSubnet.BlockchainID.String())
		require.NoError(err)

		infoClient := info.NewClient(uri)
		nodeID, _, err := infoClient.GetNodeID(ctx)
		require.NoError(err)
		warpAPIs[nodeID] = client
	}

	pChainClient := platformvm.NewClient(w.sendingSubnetURIs[0])
	pChainHeight, err := pChainClient.GetHeight(ctx)
	require.NoError(err)
	// If the source subnet is the Primary Network, then we only need to aggregate signatures from the receiving
	// subnet's validator set instead of the entire Primary Network.
	// If the destination turns out to be the Primary Network as well, then this is a no-op.
	var validators map[ids.NodeID]*validators.GetValidatorOutput
	if w.sendingSubnet.SubnetID == constants.PrimaryNetworkID {
		validators, err = pChainClient.GetValidatorsAt(ctx, w.receivingSubnet.SubnetID, pChainHeight)
	} else {
		validators, err = pChainClient.GetValidatorsAt(ctx, w.sendingSubnet.SubnetID, pChainHeight)
	}
	require.NoError(err)
	require.NotZero(len(validators))

	totalWeight := uint64(0)
	warpValidators := make([]*avalancheWarp.Validator, 0, len(validators))
	for nodeID, validator := range validators {
		warpValidators = append(warpValidators, &avalancheWarp.Validator{
			PublicKey: validator.PublicKey,
			Weight:    validator.Weight,
			NodeIDs:   []ids.NodeID{nodeID},
		})
		totalWeight += validator.Weight
	}

	log.Info("Aggregating signatures from validator set", "numValidators", len(warpValidators), "totalWeight", totalWeight)
	apiSignatureGetter := warpBackend.NewAPIFetcher(warpAPIs)
	signatureResult, err := aggregator.New(apiSignatureGetter, warpValidators, totalWeight).AggregateSignatures(ctx, w.addressedCallUnsignedMessage, 100)
	require.NoError(err)
	require.Equal(signatureResult.SignatureWeight, signatureResult.TotalWeight)
	require.Equal(signatureResult.SignatureWeight, totalWeight)

	w.addressedCallSignedMessage = signatureResult.Message

	signatureResult, err = aggregator.New(apiSignatureGetter, warpValidators, totalWeight).AggregateSignatures(ctx, w.blockPayloadUnsignedMessage, 100)
	require.NoError(err)
	require.Equal(signatureResult.SignatureWeight, signatureResult.TotalWeight)
	require.Equal(signatureResult.SignatureWeight, totalWeight)
	w.blockPayloadSignedMessage = signatureResult.Message

	log.Info("Aggregated signatures for warp messages", "addressedCallMessage", common.Bytes2Hex(w.addressedCallSignedMessage.Bytes()), "blockPayloadMessage", common.Bytes2Hex(w.blockPayloadSignedMessage.Bytes()))
}

func (w *warpTest) aggregateSignatures() {
	require := require.New(ginkgo.GinkgoT())
	ctx := e2e.DefaultContext()

	// Verify that the signature aggregation matches the results of manually constructing the warp message
	client, err := warpBackend.NewClient(w.sendingSubnetURIs[0], w.sendingSubnet.BlockchainID.String())
	require.NoError(err)

	log.Info("Fetching addressed call aggregate signature via p2p API")
	subnetIDStr := ""
	if w.sendingSubnet.SubnetID == constants.PrimaryNetworkID {
		subnetIDStr = w.receivingSubnet.SubnetID.String()
	}
	signedWarpMessageBytes, err := client.GetMessageAggregateSignature(ctx, w.addressedCallSignedMessage.ID(), warp.WarpQuorumDenominator, subnetIDStr)
	require.NoError(err)
	require.Equal(w.addressedCallSignedMessage.Bytes(), signedWarpMessageBytes)

	log.Info("Fetching block payload aggregate signature via p2p API")
	signedWarpBlockBytes, err := client.GetBlockAggregateSignature(ctx, w.blockID, warp.WarpQuorumDenominator, subnetIDStr)
	require.NoError(err)
	require.Equal(w.blockPayloadSignedMessage.Bytes(), signedWarpBlockBytes)
}

func (w *warpTest) deliverAddressedCallToReceivingSubnet() {
	require := require.New(ginkgo.GinkgoT())
	ctx := e2e.DefaultContext()

	client := w.receivingSubnetClients[0]
	log.Info("Subscribing to new heads")
	newHeads := make(chan *types.Header, 10)
	sub, err := client.SubscribeNewHead(ctx, newHeads)
	require.NoError(err)
	defer sub.Unsubscribe()

	nonce, err := client.NonceAt(ctx, w.receivingSubnetFundedAddress, nil)
	require.NoError(err)

	packedInput, err := warp.PackGetVerifiedWarpMessage(0)
	require.NoError(err)
	tx := predicate.NewPredicateTx(
		w.receivingSubnetChainID,
		nonce,
		&warp.Module.Address,
		5_000_000,
		big.NewInt(225*params.GWei),
		big.NewInt(params.GWei),
		common.Big0,
		packedInput,
		types.AccessList{},
		warp.ContractAddress,
		w.addressedCallSignedMessage.Bytes(),
	)
	signedTx, err := types.SignTx(tx, w.receivingSubnetSigner, w.receivingSubnetFundedKey)
	require.NoError(err)
	txBytes, err := signedTx.MarshalBinary()
	require.NoError(err)
	log.Info("Sending getVerifiedWarpMessage transaction", "txHash", signedTx.Hash(), "txBytes", common.Bytes2Hex(txBytes))
	require.NoError(client.SendTransaction(ctx, signedTx))

	log.Info("Waiting for new block confirmation")
	<-newHeads
	receiptCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	blockHash, _ := w.getBlockHashAndNumberFromTxReceipt(receiptCtx, client, signedTx)

	log.Info("Fetching relevant warp logs and receipts from new block")
	logs, err := client.FilterLogs(ctx, interfaces.FilterQuery{
		BlockHash: &blockHash,
		Addresses: []common.Address{warp.Module.Address},
	})
	require.NoError(err)
	require.Len(logs, 0)
	receipt, err := client.TransactionReceipt(ctx, signedTx.Hash())
	require.NoError(err)
	require.Equal(receipt.Status, types.ReceiptStatusSuccessful)
}

func (w *warpTest) deliverBlockHashPayload() {
	require := require.New(ginkgo.GinkgoT())
	ctx := e2e.DefaultContext()

	client := w.receivingSubnetClients[0]
	log.Info("Subscribing to new heads")
	newHeads := make(chan *types.Header, 10)
	sub, err := client.SubscribeNewHead(ctx, newHeads)
	require.NoError(err)
	defer sub.Unsubscribe()

	nonce, err := client.NonceAt(ctx, w.receivingSubnetFundedAddress, nil)
	require.NoError(err)

	packedInput, err := warp.PackGetVerifiedWarpBlockHash(0)
	require.NoError(err)
	tx := predicate.NewPredicateTx(
		w.receivingSubnetChainID,
		nonce,
		&warp.Module.Address,
		5_000_000,
		big.NewInt(225*params.GWei),
		big.NewInt(params.GWei),
		common.Big0,
		packedInput,
		types.AccessList{},
		warp.ContractAddress,
		w.blockPayloadSignedMessage.Bytes(),
	)
	signedTx, err := types.SignTx(tx, w.receivingSubnetSigner, w.receivingSubnetFundedKey)
	require.NoError(err)
	txBytes, err := signedTx.MarshalBinary()
	require.NoError(err)
	log.Info("Sending getVerifiedWarpBlockHash transaction", "txHash", signedTx.Hash(), "txBytes", common.Bytes2Hex(txBytes))
	err = client.SendTransaction(ctx, signedTx)
	require.NoError(err)

	log.Info("Waiting for new block confirmation")
	<-newHeads
	receiptCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	blockHash, _ := w.getBlockHashAndNumberFromTxReceipt(receiptCtx, client, signedTx)
	log.Info("Fetching relevant warp logs and receipts from new block")
	logs, err := client.FilterLogs(ctx, interfaces.FilterQuery{
		BlockHash: &blockHash,
		Addresses: []common.Address{warp.Module.Address},
	})
	require.NoError(err)
	require.Len(logs, 0)
	receipt, err := client.TransactionReceipt(ctx, signedTx.Hash())
	require.NoError(err)
	require.Equal(receipt.Status, types.ReceiptStatusSuccessful)
}

func (w *warpTest) executeHardHatTest() {
	require := require.New(ginkgo.GinkgoT())
	ctx := e2e.DefaultContext()

	client := w.sendingSubnetClients[0]
	log.Info("Subscribing to new heads")
	newHeads := make(chan *types.Header, 10)
	sub, err := client.SubscribeNewHead(ctx, newHeads)
	require.NoError(err)
	defer sub.Unsubscribe()

	chainID, err := client.ChainID(ctx)
	require.NoError(err)

	rpcURI := toRPCURI(w.sendingSubnetURIs[0], w.sendingSubnet.BlockchainID.String())

	os.Setenv("SENDER_ADDRESS", crypto.PubkeyToAddress(w.sendingSubnetFundedKey.PublicKey).Hex())
	os.Setenv("SOURCE_CHAIN_ID", "0x"+w.sendingSubnet.BlockchainID.Hex())
	os.Setenv("PAYLOAD", "0x"+common.Bytes2Hex(testPayload))
	os.Setenv("EXPECTED_UNSIGNED_MESSAGE", "0x"+hex.EncodeToString(w.addressedCallUnsignedMessage.Bytes()))
	os.Setenv("CHAIN_ID", fmt.Sprintf("%d", chainID.Uint64()))

	cmdPath := filepath.Join(repoRootPath, "contracts")
	// test path is relative to the cmd path
	testPath := "./test/warp.ts"
	utils.RunHardhatTestsCustomURI(ctx, rpcURI, cmdPath, testPath)
}

func (w *warpTest) warpLoad() {
	require := require.New(ginkgo.GinkgoT())
	ctx := e2e.DefaultContext()

	var (
		numWorkers           = len(w.sendingSubnetClients)
		txsPerWorker  uint64 = 10
		batchSize     uint64 = 10
		sendingClient        = w.sendingSubnetClients[0]
	)

	chainAKeys, chainAPrivateKeys := generateKeys(w.sendingSubnetFundedKey, numWorkers)
	chainBKeys, chainBPrivateKeys := generateKeys(w.receivingSubnetFundedKey, numWorkers)

	loadMetrics := metrics.NewDefaultMetrics()

	log.Info("Distributing funds on sending subnet", "numKeys", len(chainAKeys))
	chainAKeys, err := load.DistributeFunds(ctx, sendingClient, chainAKeys, len(chainAKeys), new(big.Int).Mul(big.NewInt(100), big.NewInt(params.Ether)), loadMetrics)
	require.NoError(err)

	log.Info("Distributing funds on receiving subnet", "numKeys", len(chainBKeys))
	_, err = load.DistributeFunds(ctx, w.receivingSubnetClients[0], chainBKeys, len(chainBKeys), new(big.Int).Mul(big.NewInt(100), big.NewInt(params.Ether)), loadMetrics)
	require.NoError(err)

	log.Info("Creating workers for each subnet...")
	chainAWorkers := make([]txs.Worker[*types.Transaction], 0, len(chainAKeys))
	for i := range chainAKeys {
		chainAWorkers = append(chainAWorkers, load.NewTxReceiptWorker(ctx, w.sendingSubnetClients[i]))
	}
	chainBWorkers := make([]txs.Worker[*types.Transaction], 0, len(chainBKeys))
	for i := range chainBKeys {
		chainBWorkers = append(chainBWorkers, load.NewTxReceiptWorker(ctx, w.receivingSubnetClients[i]))
	}

	log.Info("Subscribing to warp send events on sending subnet")
	logs := make(chan types.Log, numWorkers*int(txsPerWorker))
	sub, err := sendingClient.SubscribeFilterLogs(ctx, interfaces.FilterQuery{
		Addresses: []common.Address{warp.Module.Address},
	}, logs)
	require.NoError(err)
	defer func() {
		sub.Unsubscribe()
		err := <-sub.Err()
		require.NoError(err)
	}()

	log.Info("Generating tx sequence to send warp messages...")
	warpSendSequences, err := txs.GenerateTxSequences(ctx, func(key *ecdsa.PrivateKey, nonce uint64) (*types.Transaction, error) {
		data, err := warp.PackSendWarpMessage([]byte(fmt.Sprintf("Jets %d-%d Dolphins", key.X.Int64(), nonce)))
		if err != nil {
			return nil, err
		}
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   w.sendingSubnetChainID,
			Nonce:     nonce,
			To:        &warp.Module.Address,
			Gas:       200_000,
			GasFeeCap: big.NewInt(225 * params.GWei),
			GasTipCap: big.NewInt(params.GWei),
			Value:     common.Big0,
			Data:      data,
		})
		return types.SignTx(tx, w.sendingSubnetSigner, key)
	}, w.sendingSubnetClients[0], chainAPrivateKeys, txsPerWorker, false)
	require.NoError(err)
	log.Info("Executing warp send loader...")
	warpSendLoader := load.New(chainAWorkers, warpSendSequences, batchSize, loadMetrics)
	// TODO: execute send and receive loaders concurrently.
	require.NoError(warpSendLoader.Execute(ctx))
	require.NoError(warpSendLoader.ConfirmReachedTip(ctx))

	warpClient, err := warpBackend.NewClient(w.sendingSubnetURIs[0], w.sendingSubnet.BlockchainID.String())
	require.NoError(err)
	subnetIDStr := ""
	if w.sendingSubnet.SubnetID == constants.PrimaryNetworkID {
		subnetIDStr = w.receivingSubnet.SubnetID.String()
	}

	log.Info("Executing warp delivery sequences...")
	warpDeliverSequences, err := txs.GenerateTxSequences(ctx, func(key *ecdsa.PrivateKey, nonce uint64) (*types.Transaction, error) {
		// Wait for the next warp send log
		warpLog := <-logs

		unsignedMessage, err := warp.UnpackSendWarpEventDataToMessage(warpLog.Data)
		if err != nil {
			return nil, err
		}
		log.Info("Fetching addressed call aggregate signature via p2p API")

		signedWarpMessageBytes, err := warpClient.GetMessageAggregateSignature(ctx, unsignedMessage.ID(), warp.WarpDefaultQuorumNumerator, subnetIDStr)
		if err != nil {
			return nil, err
		}

		packedInput, err := warp.PackGetVerifiedWarpMessage(0)
		if err != nil {
			return nil, err
		}
		tx := predicate.NewPredicateTx(
			w.receivingSubnetChainID,
			nonce,
			&warp.Module.Address,
			5_000_000,
			big.NewInt(225*params.GWei),
			big.NewInt(params.GWei),
			common.Big0,
			packedInput,
			types.AccessList{},
			warp.ContractAddress,
			signedWarpMessageBytes,
		)
		return types.SignTx(tx, w.receivingSubnetSigner, key)
	}, w.receivingSubnetClients[0], chainBPrivateKeys, txsPerWorker, true)
	require.NoError(err)

	log.Info("Executing warp delivery...")
	warpDeliverLoader := load.New(chainBWorkers, warpDeliverSequences, batchSize, loadMetrics)
	require.NoError(warpDeliverLoader.Execute(ctx))
	require.NoError(warpSendLoader.ConfirmReachedTip(ctx))
	log.Info("Completed warp delivery successfully.")
}

func generateKeys(preFundedKey *ecdsa.PrivateKey, numWorkers int) ([]*key.Key, []*ecdsa.PrivateKey) {
	keys := []*key.Key{
		key.CreateKey(preFundedKey),
	}
	privateKeys := []*ecdsa.PrivateKey{
		preFundedKey,
	}
	for i := 1; i < numWorkers; i++ {
		newKey, err := key.Generate()
		require.NoError(ginkgo.GinkgoT(), err)
		keys = append(keys, newKey)
		privateKeys = append(privateKeys, newKey.PrivKey)
	}
	return keys, privateKeys
}

func toWebsocketURI(uri string, blockchainID string) string {
	return fmt.Sprintf("ws://%s/ext/bc/%s/ws", strings.TrimPrefix(uri, "http://"), blockchainID)
}

func toRPCURI(uri string, blockchainID string) string {
	return fmt.Sprintf("%s/ext/bc/%s/rpc", uri, blockchainID)
}
