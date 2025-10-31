// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Package contracts provides contract compilation directives.
// This file contains go:generate directives to compile Solidity contracts using solc.
package contracts

// Step 1: Compile Solidity contracts to generate ABI and bin files
//go:generate solc-v0.8.30 -o ../artifacts --overwrite --abi --bin --base-path . @openzeppelin/contracts/=../node_modules/@openzeppelin/contracts/ AllowList.sol ERC20NativeMinter.sol ExampleDeployerList.sol ExampleFeeManager.sol ExampleRewardManager.sol ExampleTxAllowList.sol ExampleWarp.sol

// Step 2: Generate Go bindings from the compiled artifacts
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type AllowList --abi ../artifacts/AllowList.abi --bin ../artifacts/AllowList.bin --out ../bindings/allowlist.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type ERC20NativeMinter --abi ../artifacts/ERC20NativeMinter.abi --bin ../artifacts/ERC20NativeMinter.bin --out ../bindings/erc20nativeminter.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type ExampleDeployerList --abi ../artifacts/ExampleDeployerList.abi --bin ../artifacts/ExampleDeployerList.bin --out ../bindings/exampledeployerlist.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type ExampleFeeManager --abi ../artifacts/ExampleFeeManager.abi --bin ../artifacts/ExampleFeeManager.bin --out ../bindings/examplefeemanager.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type ExampleRewardManager --abi ../artifacts/ExampleRewardManager.abi --bin ../artifacts/ExampleRewardManager.bin --out ../bindings/examplerewardmanager.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type ExampleTxAllowList --abi ../artifacts/ExampleTxAllowList.abi --bin ../artifacts/ExampleTxAllowList.bin --out ../bindings/exampletxallowlist.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type ExampleWarp --abi ../artifacts/ExampleWarp.abi --bin ../artifacts/ExampleWarp.bin --out ../bindings/examplewarp.go
