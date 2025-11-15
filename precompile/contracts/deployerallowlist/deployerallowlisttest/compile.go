// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package deployerallowlisttest

// Step 1: Compile Solidity contracts to generate ABI and bin files
//go:generate solc -o artifacts --overwrite --abi --bin --base-path . precompile/=../../../ --evm-version paris DeployerListTest.sol
// Step 2: Generate Go bindings from the compiled artifacts
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg deployerallowlisttest --type DeployerListTest --abi artifacts/DeployerListTest.abi --bin artifacts/DeployerListTest.bin --out deployerlisttest_binding.go
// Step 3: Replace import paths in generated binding to use subnet-evm instead of libevm
//go:generate sh -c "sed -i '' -e 's|github.com/ava-labs/libevm/accounts/abi|github.com/ava-labs/subnet-evm/accounts/abi|g' -e 's|github.com/ava-labs/libevm/accounts/abi/bind|github.com/ava-labs/subnet-evm/accounts/abi/bind|g' deployerlisttest_binding.go"
