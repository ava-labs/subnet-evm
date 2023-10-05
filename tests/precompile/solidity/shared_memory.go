// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements solidity tests.
package solidity

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/components/avax"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/interfaces"
	"github.com/ava-labs/subnet-evm/precompile/contracts/sharedmemory"
	"github.com/ava-labs/subnet-evm/tests/utils"
	"github.com/ava-labs/subnet-evm/utils/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// TODO: restructure to follow good ginkgo style, Before, It export, It Import
var _ = ginkgo.Describe("[Shared Memory]", ginkgo.Ordered, func() {
	// Each ginkgo It node specifies the name of the genesis file (in ./tests/precompile/genesis/)
	// to use to launch the subnet and the name of the TS test file to run on the subnet (in ./contract-examples/tests/)
	// Steps:
	// 1. Set up two blockchains with shared memory enabled on the same subnet
	// 2. Export AVAX (and other assets in the future) from blockchain A to blockchain B - verify logs, created UTXOs, and balance updates on blockchain A
	// 3. Import AVAX (and other assets in the future) onto subnet B - verify logs, verify UTXOs are consumed, and balance updated on subnet B
	ginkgo.It("Can import exported assets", ginkgo.Label("Precompile"), ginkgo.Label("SharedMemory"), func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		log.Info("Executing HardHat tests on a new blockchain", "test", "exportAVAX")

		genesisFilePath := "./tests/precompile/genesis/shared_memory.json"

		// CreateNewSubnet returns the same length of blockchainIDs as genesisFilePaths passed in or errors
		// so we do not check the length here
		subnetSuite := utils.CreateSubnetsSuite(map[string]string{
			"subnet_a": genesisFilePath,
			"subnet_b": genesisFilePath,
		})
		blockchainA := subnetSuite.GetBlockchainID("subnet_a")
		blockchainB := subnetSuite.GetBlockchainID("subnet_b")
		err := os.Setenv("BLOCKCHAIN_ID_A", blockchainA)
		gomega.Expect(err).Should(gomega.BeNil())
		err = os.Setenv("BLOCKCHAIN_ID_B", blockchainB)
		gomega.Expect(err).Should(gomega.BeNil())

		testDir := "./contracts"
		// Execute export tests
		utils.RunHardhatTests(ctx, blockchainA, testDir, "shared_memory_export")

		// Dial RPC to blockchainA to fetch logs
		uriChainA := fmt.Sprintf(
			"%s/ext/bc/%s/rpc", utils.DefaultLocalNodeURI, blockchainA)
		client, err := ethclient.Dial(uriChainA)
		gomega.Expect(err).Should(gomega.BeNil())
		defer client.Close()

		// Fetch all logs from shared memory contract so we
		// can derive the exported UTXO IDs from logs and provide
		// them to the import tests.
		logs, err := client.FilterLogs(
			ctx,
			interfaces.FilterQuery{
				Addresses: []common.Address{sharedmemory.ContractAddress},
			},
		)
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(logs).Should(gomega.HaveLen(2))

		avaxAssetID := utils.MakeDefaultWallet(ctx).P().AVAXAssetID()
		for idx, log := range logs {
			// TODO: I am going to calculate the predicate bytes here now to
			// close the loop on testing. We should have a design that does not
			// require the importer to use the codec to issue a transaction.
			var (
				parsedUTXO *avax.UTXO
				err        error
			)
			if idx == 0 {
				parsedUTXO, err = sharedmemory.ExportAVAXEventToUTXO(
					avaxAssetID,
					log.TxHash,
					int(log.Index),
					log.Topics,
					log.Data,
				)
			} else {
				parsedUTXO, err = sharedmemory.ExportUTXOEventToUTXO(
					log.TxHash,
					int(log.Index),
					log.Topics,
					log.Data,
				)
			}
			gomega.Expect(err).Should(gomega.BeNil())

			utxoID := parsedUTXO.InputID()
			err = os.Setenv(
				fmt.Sprintf("UTXO_ID_%d", idx),
				common.Bytes2Hex(utxoID[:]))
			gomega.Expect(err).Should(gomega.BeNil())

			blockchainAID, err := ids.FromString(blockchainA)
			gomega.Expect(err).Should(gomega.BeNil())
			predicate := &sharedmemory.AtomicPredicate{
				SourceChain:   blockchainAID,
				ImportedUTXOs: []*avax.UTXO{parsedUTXO},
			}
			predicateBytes, err := codec.Codec.Marshal(
				codec.CodecVersion, predicate)
			gomega.Expect(err).Should(gomega.BeNil())
			err = os.Setenv(
				fmt.Sprintf("PREDICATE_BYTES_%d", idx),
				common.Bytes2Hex(predicateBytes))
			gomega.Expect(err).Should(gomega.BeNil())

		}

		// Import the UTXOs on blockchainB
		utils.RunHardhatTests(ctx, blockchainB, testDir, "shared_memory_import")
	})
})
