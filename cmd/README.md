# Precompile Generation Tutorial

We can now generate a stateful precompile with the Precompile gen tool!

### Assumption of Knowledge
Before starting this tutorial it would be helpful if you had some context on the EVM, precompiles, and stateful precompiles. 
Here are some resources to get started put together. 

- [EVM Handbook](https://noxx3xxon.notion.site/noxx3xxon/The-EVM-Handbook-bb38e175cc404111a391907c4975426d)

- [Precompiles in Solidity](https://medium.com/@rbkhmrcr/precompiles-solidity-e5d29bd428c4)

- [Customizing the EVM with Stateful Precompiles](https://medium.com/avalancheavax/customizing-the-evm-with-stateful-precompiles-f44a34f39efd)

 
## Tutorial

Let's start by creating the Solidity interface that we want to implement. We can put this  in `./contract-examples/contracts`

```
// (c) 2022-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// SPDX-License-Identifier: MIT

pragma solidity >=0.8.0;

interface IHelloWorld {
  function sayHello() external returns (string calldata);

  // SetGreeting
  function setGreeting(string calldata recipient) external;
}
```

Now we have an interface that our precompile can implement!
Let's create an [abi](https://docs.soliditylang.org/en/v0.8.13/abi-spec.html#:~:text=Contract%20ABI%20Specification-,Basic%20Design,as%20described%20in%20this%20specification.) of our solidity code.

In the same `./contract-examples/contracts` directory, let's run

```
solcjs --abi IHelloWorld.sol
```

This spits out the abi code. Let's move it into a brand new folder in 
`./contract-examples/contracts/contract-abis`. 

IHelloWorld.abi

```
[{"inputs":[],"name":"sayHello","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"string","name":"recipient","type":"string"}],"name":"setGreeting","outputs":[],"stateMutability":"nonpayable","type":"function"}]
```


## Some facts about the precompile tool

The precompile tool takes in 4 arguments. 
```
	abiFlag = &cli.StringFlag{
		Name:  "abi",
		Usage: "Path to the Ethereum contract ABI json to bind, - for STDIN",
	}
	typeFlag = &cli.StringFlag{
		Name:  "type",
		Usage: "Struct name for the precompile (default = ABI name)",
	}
	pkgFlag = &cli.StringFlag{
		Name:  "pkg",
		Usage: "Package name to generate the precompile into (default = precompile)",
	}
	outFlag = &cli.StringFlag{
		Name:  "out",
		Usage: "Output file for the generated precompile (default = STDOUT)",
	}
```

Currently it can only generate precompiles in Golang only.  


## Generating the precompile 

Now that we have an abi for the precompile gen tool to interact with. We can run the following command to generate our HelloWorld precompile!


In the root of the repo run 
```
go run ./cmd/precompilegen/main.go --abi ./contract-examples/contracts/contract-abis/IHelloWorld.abi --out ./precompile/hello_world.go
```


Wow! We just got a precompile template that's 80% complete located at `./precompile/hello_world.go`. Let's fill out the rest!






