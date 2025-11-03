// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package contracts

// Step 1: Compile Solidity contracts to generate ABI and bin files
//go:generate solc-v0.8.30 -o ../artifacts --overwrite --abi --bin --base-path . --evm-version paris @openzeppelin/contracts/=../node_modules/@openzeppelin/contracts/ AllowList.sol ERC20NativeMinter.sol ExampleDeployerList.sol ExampleFeeManager.sol ExampleRewardManager.sol ExampleTxAllowList.sol ExampleWarp.sol interfaces/IAllowList.sol interfaces/IFeeManager.sol interfaces/INativeMinter.sol interfaces/IRewardManager.sol interfaces/IWarpMessenger.sol

// Step 2: Generate Go bindings from the compiled artifacts
// Main contracts
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type AllowList --abi ../artifacts/AllowList.abi --bin ../artifacts/AllowList.bin --out ../bindings/allowlist.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type ERC20NativeMinter --abi ../artifacts/ERC20NativeMinter.abi --bin ../artifacts/ERC20NativeMinter.bin --out ../bindings/erc20nativeminter.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type ExampleDeployerList --abi ../artifacts/ExampleDeployerList.abi --bin ../artifacts/ExampleDeployerList.bin --out ../bindings/exampledeployerlist.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type ExampleFeeManager --abi ../artifacts/ExampleFeeManager.abi --bin ../artifacts/ExampleFeeManager.bin --out ../bindings/examplefeemanager.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type ExampleRewardManager --abi ../artifacts/ExampleRewardManager.abi --bin ../artifacts/ExampleRewardManager.bin --out ../bindings/examplerewardmanager.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type ExampleTxAllowList --abi ../artifacts/ExampleTxAllowList.abi --bin ../artifacts/ExampleTxAllowList.bin --out ../bindings/exampletxallowlist.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type ExampleWarp --abi ../artifacts/ExampleWarp.abi --bin ../artifacts/ExampleWarp.bin --out ../bindings/examplewarp.go

// Interface contracts (for interacting with precompiles)
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type IAllowList --abi ../artifacts/IAllowList.abi --bin ../artifacts/IAllowList.bin --out ../bindings/iallowlist.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type IFeeManager --abi ../artifacts/IFeeManager.abi --bin ../artifacts/IFeeManager.bin --out ../bindings/ifeemanager.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type INativeMinter --abi ../artifacts/INativeMinter.abi --bin ../artifacts/INativeMinter.bin --out ../bindings/inativeminter.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type IRewardManager --abi ../artifacts/IRewardManager.abi --bin ../artifacts/IRewardManager.bin --out ../bindings/irewardmanager.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg bindings --type IWarpMessenger --abi ../artifacts/IWarpMessenger.abi --bin ../artifacts/IWarpMessenger.bin --out ../bindings/iwarpmessenger.go
