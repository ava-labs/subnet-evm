// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements solidity tests.
package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ava-labs/avalanche-network-runner/rpcpb"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/plugin/evm"
	"github.com/ava-labs/subnet-evm/tests/utils"
	"github.com/ava-labs/subnet-evm/tests/utils/runner"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

const fundedKeyStr = "56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027"

var (
	config              = runner.NewDefaultANRConfig()
	manager             = runner.NewNetworkManager(config)
	warpChainConfigPath string
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
	warpSetup()
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

func warpSetup() {
	var (
		chainAURIs, chainBURIs []string
	)

	fundedKey, err := crypto.HexToECDSA(fundedKeyStr)
	if err != nil {
		panic(err)
	}
	fundedAddress := crypto.PubkeyToAddress(fundedKey.PublicKey)

	log.Info("Funded address", "address", fundedAddress.String(), "fundedKey", fundedKeyStr)

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

	subnetB := subnetIDs[1]
	subnetBDetails, ok := manager.GetSubnet(subnetB)
	if !ok {
		panic("subnetB not found")
	}
	blockchainIDB := subnetBDetails.BlockchainID
	if len(subnetBDetails.ValidatorURIs) != 5 {
		panic("expected 5 validators in subnetB")
	}
	chainBURIs = append(chainBURIs, subnetBDetails.ValidatorURIs...)

	log.Info("Created URIs for both subnets", "ChainAURIs", chainAURIs, "ChainBURIs", chainBURIs, "blockchainIDA", blockchainIDA, "blockchainIDB", blockchainIDB)

	chainAWSURI := toWebsocketURI(chainAURIs[0], blockchainIDA.String())
	log.Info("Creating ethclient for blockchainA", "wsURI", chainAWSURI)

	chainBWSURI := toWebsocketURI(chainBURIs[0], blockchainIDB.String())
	log.Info("Creating ethclient for blockchainB", "wsURI", chainBWSURI)
}
