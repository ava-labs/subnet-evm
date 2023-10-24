// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements solidity tests.
package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ava-labs/avalanche-network-runner/rpcpb"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/interfaces"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/plugin/evm"
	"github.com/ava-labs/subnet-evm/tests/utils"
	"github.com/ava-labs/subnet-evm/tests/utils/runner"
	"github.com/ava-labs/subnet-evm/x/warp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

const fundedKeyStr = "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027"

var (
	config              = runner.NewDefaultANRConfig()
	manager             = runner.NewNetworkManager(config)
	warpChainConfigPath string

	chainID     = big.NewInt(99999)
	testPayload = []byte{1, 2, 3}
	txSigner    = types.LatestSignerForChainID(chainID)
)

func toWebsocketURI(uri string, blockchainID string) string {
	return fmt.Sprintf("ws://%s/ext/bc/%s/ws", strings.TrimPrefix(uri, "http://"), blockchainID)
}

// Starts the default network and adds 10 new nodes as validators with BLS keys
// registered on the P-Chain.
// Adds two disjoint sets of 5 of the new validator nodes to validate two new subnets with a
// a single Subnet-EVM blockchain.
func main() {
	ctx := context.Background()

	// Create buffered sigChan to receive SIGINT notifications
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	// Create context with cancel
	ctx, cancel := context.WithCancel(ctx)

	setup(ctx)
	cancelChan := make(chan os.Signal, 1)
	// catch SIGETRM or SIGINTERRUPT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-cancelChan
	log.Info("Exitting...", "caught signal", sig)
	cancel()
	shutdown()
}
func setup(ctx context.Context) {
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
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	_, err = f.Write([]byte(`{"warp-api-enabled": true}`))
	if err != nil {
		panic(err)
	}
	warpChainConfigPath = f.Name()

	// Construct the network using the avalanche-network-runner
	_, err = manager.StartDefaultNetwork(ctx)
	if err != nil {
		panic(err)
	}
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
	if err != nil {
		panic(err)
	}

	// Issue transactions to activate the proposerVM fork on the receiving chain
	chainID := big.NewInt(99999)
	fundedKey, err := crypto.HexToECDSA(fundedKeyStr)
	if err != nil {
		panic(err)
	}

	fundedAddress := crypto.PubkeyToAddress(fundedKey.PublicKey)

	log.Info("Funded address", "address", fundedAddress.String(), "fundedKey", fundedKeyStr)
	subnetB := manager.GetSubnets()[1]
	subnetBDetails, ok := manager.GetSubnet(subnetB)
	if !ok {
		panic("subnetB not found")
	}

	chainBID := subnetBDetails.BlockchainID
	uri := toWebsocketURI(subnetBDetails.ValidatorURIs[0], chainBID.String())
	client, err := ethclient.Dial(uri)
	if err != nil {
		panic(err)
	}

	err = utils.IssueTxsToActivateProposerVMFork(ctx, chainID, fundedKey, client)
	if err != nil {
		panic(err)
	}
	var (
		chainAURIs, chainBURIs []string
	)

	subnetIDs := manager.GetSubnets()
	if len(subnetIDs) != 2 {
		panic("expected 2 subnets")
	}

	subnetA := subnetIDs[0]
	subnetADetails, ok := manager.GetSubnet(subnetA)
	if !ok {
		panic("subnetA not found")
	}

	blockchainIDA := subnetADetails.BlockchainID
	if len(subnetADetails.ValidatorURIs) != 5 {
		panic("expected 5 validators in subnetA")
	}
	chainAURIs = append(chainAURIs, subnetADetails.ValidatorURIs...)

	blockchainIDB := subnetBDetails.BlockchainID
	if len(subnetBDetails.ValidatorURIs) != 5 {
		panic("expected 5 validators in subnetB")
	}
	chainBURIs = append(chainBURIs, subnetBDetails.ValidatorURIs...)

	// print out full URIs for both chains
	for i, uri := range chainAURIs {
		log.Info("Printing full chain URI for Subnet A", "nodeName", subnetANodeNames[i], "uri", fmt.Sprintf("%s/ext/bc/%s", uri, blockchainIDA.String()))
	}

	for i, uri := range chainBURIs {
		log.Info("Printing full chain URI for Subnet B", "nodeName", subnetBNodeNames[i], "uri", fmt.Sprintf("%s/ext/bc/%s", uri, blockchainIDB.String()))
	}

	chainAWSURI := toWebsocketURI(chainAURIs[0], blockchainIDA.String())
	log.Info("Creating ethclient for blockchainA", "wsURI", chainAWSURI)
	chainAWSClient, err := ethclient.Dial(chainAWSURI)
	if err != nil {
		panic(err)
	}

	log.Info("Subscribing to new heads")
	newHeads := make(chan *types.Header, 10)
	sub, err := chainAWSClient.SubscribeNewHead(ctx, newHeads)
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	log.Info("Subscribing to pending txs")
	newPendingTxs := make(chan *common.Hash, 10)
	sub, err = chainAWSClient.SubscribeNewPendingTransactions(ctx, newPendingTxs)
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	startingNonce, err := chainAWSClient.NonceAt(ctx, fundedAddress, nil)
	if err != nil {
		panic(err)
	}

	// Create 7 accounts
	keys := make([]*ecdsa.PrivateKey, 7)
	accs := make([]common.Address, len(keys))

	for i := 0; i < len(keys); i++ {
		keys[i], err = crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		accs[i] = crypto.PubkeyToAddress(keys[i].PublicKey)
	}

	// Fund those accounts
	for i := 0; i < len(keys); i++ {
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     startingNonce + uint64(i),
			To:        &accs[i],
			Gas:       200_000,
			GasFeeCap: big.NewInt(225 * params.GWei),
			GasTipCap: big.NewInt(params.GWei),
			Value:     new(big.Int).Mul(big.NewInt(100_000), big.NewInt(params.Ether)),
		})
		signedTx, err := types.SignTx(tx, txSigner, fundedKey)
		if err != nil {
			panic(err)
		}
		log.Info("Funding account", "account", accs[i], "txHash", signedTx.Hash())
		err = chainAWSClient.SendTransaction(ctx, signedTx)
		if err != nil {
			panic(err)
		}
	}

	targetNonce := startingNonce + uint64(len(keys))
	for targetNonce != startingNonce {
		log.Info("Blocking until all accounts are funded")
		startingNonce, err = chainAWSClient.NonceAt(ctx, fundedAddress, nil)
		if err != nil {
			panic(err)
		}
		time.Sleep(2 * time.Second)
	}

	for i := 0; i < len(keys); i++ {
		bal, err := chainAWSClient.BalanceAt(ctx, accs[i], nil)
		if err != nil {
			panic(err)
		}
		log.Info("Account balance", "account", accs[i], "balance", bal)
	}

	go listenForLogs(ctx, newHeads, chainAWSClient)
	go listenForPendingTxs(ctx, newPendingTxs, chainAWSClient)

	time.Sleep(5 * time.Second)

	go spamWarpMessages(ctx, chainAWSClient, fundedKey, fundedAddress, 50)

	for i := 0; i < len(keys); i++ {
		go spamWarpMessages(ctx, chainAWSClient, keys[i], accs[i], 50)
	}
}

func spamWarpMessages(ctx context.Context, client ethclient.Client, key *ecdsa.PrivateKey, addr common.Address, numTxs uint64) {
	log.Info("spamming network with warp messages", "addr", addr, "numTxs", numTxs)

	startingNonce, err := client.NonceAt(ctx, addr, nil)
	if err != nil {
		panic(err)
	}

	packedInput, err := warp.PackSendWarpMessage(testPayload)
	if err != nil {
		panic(err)
	}

	for i := uint64(0); i < numTxs; i++ {
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     startingNonce + i,
			To:        &warp.Module.Address,
			Gas:       200_000,
			GasFeeCap: big.NewInt(225 * params.GWei),
			GasTipCap: big.NewInt(params.GWei),
			Value:     common.Big0,
			Data:      packedInput,
		})
		signedTx, err := types.SignTx(tx, txSigner, key)
		if err != nil {
			panic(err)
		}
		err = client.SendTransaction(ctx, signedTx)
		if err != nil {
			panic(err)
		}
	}
}

func listenForPendingTxs(ctx context.Context, newPendingTxs chan *common.Hash, client ethclient.Client) {
	for {
		newPendingTx := <-newPendingTxs
		log.Info("new pending tx", "txHash", newPendingTx)
	}
}

func listenForLogs(ctx context.Context, newHeads chan *types.Header, chainAWSClient ethclient.Client) {
	for {
		newHead := <-newHeads
		blockHash := newHead.Hash()

		log.Info("new block", "blockHash", blockHash)

		logs, err := chainAWSClient.FilterLogs(ctx, interfaces.FilterQuery{
			BlockHash: &blockHash,
			Addresses: []common.Address{warp.Module.Address},
		})
		if err != nil {
			panic(err)
		}

		for _, ethLog := range logs {
			log.Info("received ethLog",
				"address", ethLog.Address,
				"blockNumber", ethLog.BlockNumber,
				"txIndex", ethLog.TxIndex,
			)
		}
	}
}

func shutdown() {
	if manager == nil {
		return
	}

	if err := manager.TeardownNetwork(); err != nil {
		panic(err)
	}

	if err := os.Remove(warpChainConfigPath); err != nil {
		panic(err)
	}
}
