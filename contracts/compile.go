// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Package contracts provides contract compilation directives.
// This file contains go:generate directives to compile Solidity contracts using solc.
package contracts

//go:generate solc-v0.8.30 -o ./artifacts --overwrite --abi --bin --base-path . @openzeppelin/contracts/=lib/openzeppelin-contracts/contracts/ contracts/AllowList.sol contracts/ERC20NativeMinter.sol contracts/ExampleDeployerList.sol contracts/ExampleFeeManager.sol contracts/ExampleRewardManager.sol contracts/ExampleTxAllowList.sol contracts/ExampleWarp.sol

// Compile interface contracts
//go:generate solc-v0.8.30 -o ./artifacts --overwrite --abi --bin --base-path . contracts/interfaces/IAllowList.sol @openzeppelin/contracts/=lib/openzeppelin-contracts/contracts/
//go:generate solc-v0.8.30 -o ./artifacts --overwrite --abi --bin --base-path . contracts/interfaces/IFeeManager.sol @openzeppelin/contracts/=lib/openzeppelin-contracts/contracts/
//go:generate solc-v0.8.30 -o ./artifacts --overwrite --abi --bin --base-path . contracts/interfaces/INativeMinter.sol @openzeppelin/contracts/=lib/openzeppelin-contracts/contracts/
//go:generate solc-v0.8.30 -o ./artifacts --overwrite --abi --bin --base-path . contracts/interfaces/IRewardManager.sol @openzeppelin/contracts/=lib/openzeppelin-contracts/contracts/
//go:generate solc-v0.8.30 -o ./artifacts --overwrite --abi --bin --base-path . contracts/interfaces/IWarpMessenger.sol @openzeppelin/contracts/=lib/openzeppelin-contracts/contracts/

// NOTE: Test contracts (contracts/test/*Test.sol) are NOT compiled here
// They use DS-Test which we're deprecating. These will be replaced with Go tests in Phase 3.
// The existing TypeScript tests will continue to use Hardhat-compiled versions until migration is complete.
