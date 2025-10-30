// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Package bindings provides Go bindings for Solidity contracts.
// This file contains go:generate directives to generate Go bindings using abigen.
package bindings

//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type AllowList --abi ../artifacts/AllowList.abi --bin ../artifacts/AllowList.bin --out allowlist.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type ERC20NativeMinter --abi ../artifacts/ERC20NativeMinter.abi --bin ../artifacts/ERC20NativeMinter.bin --out erc20nativeminter.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type ExampleDeployerList --abi ../artifacts/ExampleDeployerList.abi --bin ../artifacts/ExampleDeployerList.bin --out exampledeployerlist.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type ExampleFeeManager --abi ../artifacts/ExampleFeeManager.abi --bin ../artifacts/ExampleFeeManager.bin --out examplefeemanager.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type ExampleRewardManager --abi ../artifacts/ExampleRewardManager.abi --bin ../artifacts/ExampleRewardManager.bin --out examplerewardmanager.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type ExampleTxAllowList --abi ../artifacts/ExampleTxAllowList.abi --bin ../artifacts/ExampleTxAllowList.bin --out exampletxallowlist.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type ExampleWarp --abi ../artifacts/ExampleWarp.abi --bin ../artifacts/ExampleWarp.bin --out examplewarp.go
