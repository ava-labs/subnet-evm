// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Package contracts provides contract compilation directives.
// This file contains go:generate directives to compile Solidity contracts using solc.
package contracts

//go:generate solc-v0.8.30 -o ../artifacts --overwrite --abi --bin --base-path . @openzeppelin/contracts/=../lib/openzeppelin-contracts/contracts/ AllowList.sol ERC20NativeMinter.sol ExampleDeployerList.sol ExampleFeeManager.sol ExampleRewardManager.sol ExampleTxAllowList.sol ExampleWarp.sol
