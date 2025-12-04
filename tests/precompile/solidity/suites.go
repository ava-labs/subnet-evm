// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements solidity tests.
package solidity

import (
	"context"
	"fmt"

	"github.com/ava-labs/subnet-evm/tests/utils"

	ginkgo "github.com/onsi/ginkgo/v2"
)

// Registers the Asynchronized Precompile Tests
// Before running the tests, this function creates all subnets given in the genesis files
// and then runs the hardhat tests for each one asynchronously if called with `ginkgo run -procs=`.
func RegisterAsyncTests() {
	// Tests here assumes that the genesis files are in ./tests/precompile/genesis/
	// with the name {precompile_name}.json
	genesisFiles, err := utils.GetFilesAndAliases("./tests/precompile/genesis/*.json")
	if err != nil {
		ginkgo.AbortSuite("Failed to get genesis files: " + err.Error())
	}
	if len(genesisFiles) == 0 {
		ginkgo.AbortSuite("No genesis files found")
	}
}

//	Default parameters are:
//
// 1. Hardhat contract environment is located at ./contracts
// 2. Hardhat test file is located at ./contracts/test/<test>.ts
// 3. npx is available in the ./contracts directory
func runDefaultHardhatTests(ctx context.Context, blockchainID, testName string) {
	cmdPath := "./contracts"
	// test path is relative to the cmd path
	testPath := fmt.Sprintf("./test/%s.ts", testName)
	utils.RunHardhatTests(ctx, blockchainID, cmdPath, testPath)
}
