// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Package contracts provides contract compilation directives.
// This file contains go:generate directives to compile Solidity contracts using solc.
package contracts

// Compile main contracts
//go:generate solc -o ./artifacts --overwrite --abi --bin --base-path . contracts/AllowList.sol @openzeppelin/contracts/=submodules/openzeppelin-contracts/contracts/
//go:generate solc -o ./artifacts --overwrite --abi --bin --base-path . contracts/ERC20NativeMinter.sol @openzeppelin/contracts/=submodules/openzeppelin-contracts/contracts/
//go:generate solc -o ./artifacts --overwrite --abi --bin --base-path . contracts/ExampleDeployerList.sol @openzeppelin/contracts/=submodules/openzeppelin-contracts/contracts/
//go:generate solc -o ./artifacts --overwrite --abi --bin --base-path . contracts/ExampleFeeManager.sol @openzeppelin/contracts/=submodules/openzeppelin-contracts/contracts/
//go:generate solc -o ./artifacts --overwrite --abi --bin --base-path . contracts/ExampleRewardManager.sol @openzeppelin/contracts/=submodules/openzeppelin-contracts/contracts/
//go:generate solc -o ./artifacts --overwrite --abi --bin --base-path . contracts/ExampleTxAllowList.sol @openzeppelin/contracts/=submodules/openzeppelin-contracts/contracts/
//go:generate solc -o ./artifacts --overwrite --abi --bin --base-path . contracts/ExampleWarp.sol @openzeppelin/contracts/=submodules/openzeppelin-contracts/contracts/

// Compile interface contracts
//go:generate solc -o ./artifacts --overwrite --abi --bin --base-path . contracts/interfaces/IAllowList.sol @openzeppelin/contracts/=submodules/openzeppelin-contracts/contracts/
//go:generate solc -o ./artifacts --overwrite --abi --bin --base-path . contracts/interfaces/IFeeManager.sol @openzeppelin/contracts/=submodules/openzeppelin-contracts/contracts/
//go:generate solc -o ./artifacts --overwrite --abi --bin --base-path . contracts/interfaces/INativeMinter.sol @openzeppelin/contracts/=submodules/openzeppelin-contracts/contracts/
//go:generate solc -o ./artifacts --overwrite --abi --bin --base-path . contracts/interfaces/IRewardManager.sol @openzeppelin/contracts/=submodules/openzeppelin-contracts/contracts/
//go:generate solc -o ./artifacts --overwrite --abi --bin --base-path . contracts/interfaces/IWarpMessenger.sol @openzeppelin/contracts/=submodules/openzeppelin-contracts/contracts/

// NOTE: Test contracts (contracts/test/*Test.sol) are NOT compiled here
// They use DS-Test which we're deprecating. These will be replaced with Go tests in Phase 3.
// The existing TypeScript tests will continue to use Hardhat-compiled versions until migration is complete.
