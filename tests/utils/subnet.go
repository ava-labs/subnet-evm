// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/avalanchego/genesis"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	wallet "github.com/ava-labs/avalanchego/wallet/subnet/primary"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/plugin/evm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/onsi/gomega"
	"golang.org/x/sync/errgroup"
)

func RunHardhatTests(test string, rpcURI string) {
	log.Info("Sleeping to wait for test ping", "rpcURI", rpcURI)
	client, err := NewEvmClient(rpcURI, 225, 2)
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
	if err != nil {
		fmt.Printf("\nErr: %s\n", err.Error())
	}
	gomega.Expect(err).Should(gomega.BeNil())
}

// CreateNewSubnet creates subnets with the specified genesisFilePaths
// using the P chain wallet [wallet] and returns the IDs of the newly created
// blockchains along with the AVAX asset ID.
func CreateNewSubnet(ctx context.Context, genesisFilePaths []string) ([]ids.ID, ids.ID) {
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

	genesisBytesArr := make([][]byte, 0, len(genesisFilePaths))
	wd, err := os.Getwd()
	gomega.Expect(err).Should(gomega.BeNil())
	log.Info("Creating new subnet with specified blockchains", "wd", wd)

	for _, genesisFilePath := range genesisFilePaths {
		log.Info("Reading genesis file", "filePath", genesisFilePath)
		genesisBytes, err := os.ReadFile(genesisFilePath)
		gomega.Expect(err).Should(gomega.BeNil())
		genesisBytesArr = append(genesisBytesArr, genesisBytes)
	}

	log.Info("Creating new subnet")
	createSubnetTxID, err := pWallet.IssueCreateSubnetTx(owner)
	gomega.Expect(err).Should(gomega.BeNil())

	blockchainIDs := make([]ids.ID, 0, len(genesisBytesArr))
	for _, genesisBytes := range genesisBytesArr {
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
		blockchainIDs = append(blockchainIDs, createChainTxID)
	}

	eg, egCtx := errgroup.WithContext(ctx)
	for _, blockchainID := range blockchainIDs {
		blockchainID := blockchainID
		eg.Go(func() error {
			// Confirm the new blockchain is ready by waiting for the readiness endpoint
			infoClient := info.NewClient(DefaultLocalNodeURI)
			bootstrapped, err := info.AwaitBootstrapped(egCtx, infoClient, blockchainID.String(), 2*time.Second)
			if err != nil {
				return err
			}
			if !bootstrapped {
				return fmt.Errorf("blockchain %s not bootstrapped", blockchainID)
			}
			return nil
		})
	}
	// Check that all blockchains bootstrap correctly
	gomega.Expect(eg.Wait()).Should(gomega.BeNil())

	// Return the blockchainIDs of the newly created blockchains
	return blockchainIDs, pWallet.AVAXAssetID()
}

func ExecuteHardHatTestOnNewBlockchain(ctx context.Context, test string) {
	log.Info("Executing HardHat tests on a new blockchain", "test", test)

	genesisFilePath := fmt.Sprintf("./tests/precompile/genesis/%s.json", test)

	blockchainIDs, _ := CreateNewSubnet(ctx, []string{genesisFilePath})
	chainURI := fmt.Sprintf("%s/ext/bc/%s/rpc", DefaultLocalNodeURI, blockchainIDs[0])

	log.Info("Created subnet successfully", "ChainURI", chainURI)
	RunHardhatTests(test, chainURI)
}
