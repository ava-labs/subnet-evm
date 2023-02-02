// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements solidity tests.
package solidity

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ava-labs/subnet-evm/tests/utils"
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
	ginkgo.It("ExportAVAX", ginkgo.Label("Precompile"), ginkgo.Label("ExportAVAX"), ginkgo.Label("SharedMemory"), func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		log.Info("Executing HardHat tests on a new blockchain", "test", "exportAVAX")

		genesisFilePath := "./tests/precompile/genesis/shared_memory.json"

		// CreateNewSubnet returns the same length of blockchainIDs as genesisFilePaths passed in or errors
		// so we do not check the length here
		blockchainIDs := utils.CreateNewSubnet(ctx, []string{genesisFilePath, genesisFilePath})
		blockchainA := blockchainIDs[0]
		blockchainB := blockchainIDs[1]
		err := os.Setenv("BLOCKCHAIN_ID_A", blockchainA.Hex())
		gomega.Expect(err).Should(gomega.BeNil())
		err = os.Setenv("BLOCKCHAIN_ID_B", blockchainB.Hex())
		gomega.Expect(err).Should(gomega.BeNil())
		chainURI := fmt.Sprintf("%s/ext/bc/%s/rpc", utils.DefaultLocalNodeURI, blockchainA.String())

		// Execute export tests
		utils.RunHardhatTests("shared_memory_export", chainURI) // TODO: complete TODOs in this test file

		// TODO: Verify eth_logs for shared memory are correct

		// Confirm via API that the UTXOs were created on blockchainB
		utils.RunHardhatTests("shared_memory_import", chainURI) // TODO: implement this test file
	})
})
