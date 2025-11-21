// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nativemintertest

// Step 1: Compile Solidity contracts to generate ABI and bin files
//go:generate solc-v0.8.30 -o artifacts --overwrite --abi --bin --base-path ../../../.. @openzeppelin/contracts/=contracts/lib/openzeppelin-contracts/contracts/ precompile/=precompile/ --evm-version paris solidity/ERC20NativeMinter.sol solidity/Minter.sol
// Step 2: Generate Go bindings from the compiled artifacts
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg nativemintertest --type INativeMinter --abi artifacts/INativeMinter.abi --bin artifacts/INativeMinter.bin --out gen_inativeminter_binding.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg nativemintertest --type ERC20NativeMinter --abi artifacts/ERC20NativeMinter.abi --bin artifacts/ERC20NativeMinter.bin --out gen_erc20nativeminter_binding.go
//go:generate go run github.com/ava-labs/libevm/cmd/abigen --pkg nativemintertest --type Minter --abi artifacts/Minter.abi --bin artifacts/Minter.bin --out gen_minter_binding.go
// Step 3: Replace import paths in generated binding to use subnet-evm instead of libevm
//go:generate sh -c "sed -i.bak -e 's|github.com/ava-labs/libevm/accounts/abi|github.com/ava-labs/subnet-evm/accounts/abi|g' -e 's|github.com/ava-labs/libevm/accounts/abi/bind|github.com/ava-labs/subnet-evm/accounts/abi/bind|g' gen_inativeminter_binding.go gen_erc20nativeminter_binding.go gen_minter_binding.go && rm -f gen_inativeminter_binding.go.bak gen_erc20nativeminter_binding.go.bak gen_minter_binding.go.bak"
