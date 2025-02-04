# Subnet EVM

[![CI](https://github.com/ava-labs/subnet-evm/actions/workflows/ci.yml/badge.svg)](https://github.com/ava-labs/subnet-evm/actions/workflows/ci.yml)
[![CodeQL](https://github.com/ava-labs/subnet-evm/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/ava-labs/subnet-evm/actions/workflows/codeql-analysis.yml)

[Avalanche](https://docs.avax.network/avalanche-l1s) is a network composed of multiple blockchains.
Each blockchain is an instance of a Virtual Machine (VM), much like an object in an object-oriented language is an instance of a class.
That is, the VM defines the behavior of the blockchain.

Subnet EVM is the [Virtual Machine (VM)](https://docs.avax.network/learn/virtual-machines) that defines the Subnet Contract Chains. Subnet EVM is a simplified version of [Coreth VM (C-Chain)](https://github.com/ava-labs/coreth).

This chain implements the Ethereum Virtual Machine and supports Solidity smart contracts as well as most other Ethereum client functionality.

## Building

The Subnet EVM runs in a separate process from the main AvalancheGo process and communicates with it over a local gRPC connection.

### AvalancheGo Compatibility

```text
[v0.7.0] AvalancheGo@v1.12.0-v1.12.1 (Protocol Version: 38)
[v0.7.1] AvalancheGo@v1.12.2 (Protocol Version: 39)
```

## API

The Subnet EVM supports the following API namespaces:

- `eth`
- `personal`
- `txpool`
- `debug`

Only the `eth` namespace is enabled by default.
Subnet EVM is a simplified version of [Coreth VM (C-Chain)](https://github.com/ava-labs/coreth).
Full documentation for the C-Chain's API can be found [here](https://docs.avax.network/apis/avalanchego/apis/c-chain).

## Compatibility

The Subnet EVM is compatible with almost all Ethereum tooling, including [Remix](https://docs.avax.network/build/dapp/smart-contracts/remix-deploy), [Metamask](https://docs.avax.network/build/dapp/chain-settings), and [Foundry](https://docs.avax.network/build/dapp/smart-contracts/toolchains/foundry).

## Differences Between Subnet EVM and Coreth

- Added configurable fees and gas limits in genesis
- Merged Avalanche hardforks into the single "Subnet EVM" hardfork
- Removed Atomic Txs and Shared Memory
- Removed Multicoin Contract and State

## Block Format

To support these changes, there have been a number of changes to the SubnetEVM block format compared to what exists on the C-Chain and Ethereum. Here we list the changes to the block format as compared to Ethereum.

### Block Header

- `BaseFee`: Added by EIP-1559 to represent the base fee of the block (present in Ethereum as of EIP-1559)
- `BlockGasCost`: surcharge for producing a block faster than the target rate

## Create an EVM Subnet on a Local Network

### Clone Subnet-evm

First install Go 1.22.8 or later. Follow the instructions [here](https://go.dev/doc/install). You can verify by running `go version`.

Set `$GOPATH` environment variable properly for Go to look for Go Workspaces. Please read [this](https://go.dev/doc/code) for details. You can verify by running `echo $GOPATH`.

As a few software will be installed into `$GOPATH/bin`, please make sure that `$GOPATH/bin` is in your `$PATH`, otherwise, you may get error running the commands below.

Download the `subnet-evm` repository into your `$GOPATH`:

```sh
cd $GOPATH
mkdir -p src/github.com/ava-labs
cd src/github.com/ava-labs
git clone git@github.com:ava-labs/subnet-evm.git
cd subnet-evm
```

This will clone and checkout to `master` branch.

### Run Local Network

To run a local network, it is recommended to use the [avalanche-cli](https://github.com/ava-labs/avalanche-cli#avalanche-cli) to set up an instance of Subnet-EVM on a local Avalanche Network.

There are two options when using the Avalanche-CLI:

1. Use an official Subnet-EVM release: https://docs.avax.network/subnets/build-first-subnet
2. Build and deploy a locally built (and optionally modified) version of Subnet-EVM: https://docs.avax.network/subnets/create-custom-subnet

## Run in Docker

The `subnet-evm` Docker image comes with AvalancheGo pre-installed, making it easy to run a node. You can find the latest image tags on [Docker Hub](https://hub.docker.com/r/avaplatform/subnet-evm/tags).

### Configuration

You can configure the `subnet-evm` Docker container using environment variables, flags or standard AvalancheGo config files.

- **Environment Variables**: Use uppercase variables prefixed with `AVAGO_`. For example, `AVAGO_NETWORK_ID` corresponds to the `--network-id` [flag in AvalancheGo](https://docs.avax.network/nodes/configure/configs-flags#--network-id-string).
- **Flags**: Please refer to the [AvalancheGo documentation](https://docs.avax.network/nodes/configure/configs-flags) for a list of available flags.
- **Config Files**: Configure as you would with the regular AvalancheGo binary using config files. Mount the config directory with `-v ~/.avalanchego:/root/.avalanchego`.

### Data Persistence

To persist data across container restarts, you need to mount the `/root/.avalanchego` directory. For example, `-v ~/.avalanchego:/root/.avalanchego`. The container runs as root by default.

### Updating

Run `docker stop avago; docker rm avago;` then start a new container with the latest version tag in your `docker run` command.

### Networking

Using `--network host` is recommended to avoid any issues.
If you know what you are doing, you will need port `AVAGO_STAKING_PORT` (default `9651`) open for the validator to connect to the subnet. For the RPC server, open `AVAGO_HTTP_PORT` (default `9650`). Do not attempt to remap `AVAGO_STAKING_PORT` using the Docker `-p` flag (e.g., `-p 9651:1234`); it will not work. Instead, set `AVAGO_STAKING_PORT=1234` and then use `-p 1234:1234`.

This is because the staking port is used for peer-to-peer communication between validators, and the port number is part of the node's identity in the network. If you remap it using Docker's port mapping, other nodes will still try to connect to the original port number advertised by your node, leading to connection failures. The HTTP port can be remapped freely since it's only used for API calls.

### Example Configs

#### Fuji Subnet Validator

```bash
docker run -it -d \
  --name avago \
  --network host \
  -v ~/.avalanchego:/home/avalanche/.avalanchego \
  -e AVAGO_NETWORK_ID=fuji \
  -e AVAGO_PARTIAL_SYNC_PRIMARY_NETWORK=true \
  -e AVAGO_TRACK_SUBNETS=REPLACE_THIS_WITH_YOUR_SUBNET_ID \
  -e AVAGO_PUBLIC_IP_RESOLUTION_SERVICE=ifconfigme \
  avaplatform/subnet-evm:v0.7.1-rc.0
```

- `AVAGO_PARTIAL_SYNC_PRIMARY_NETWORK`: Ensures you don't sync the X and C-Chains.
- `AVAGO_TRACK_SUBNETS`: Sets the subnet ID to track. It will track all chains in the subnet.
- `AVAGO_NETWORK_ID=fuji`: Sets the network ID to Fuji. Remove to sync Mainnet.
- `AVAGO_PUBLIC_IP_RESOLUTION_SERVICE=ifconfigme`: Required for AWS EC2 instances to be accessed from outside AWS.

#### Fuji Subnet RPC

```bash
docker run -it -d \
  --name rpc \
  --network host \
  -v ~/.avalanchego_rpc/:/home/avalanche/.avalanchego \
  -e AVAGO_NETWORK_ID=fuji \
  -e AVAGO_PARTIAL_SYNC_PRIMARY_NETWORK=true \
  -e AVAGO_TRACK_SUBNETS=hk755meusfKqBb9C9RfzzCxZFkdSXhaFHTcbtycMUSQ11o2cd \
  -e AVAGO_HTTP_PORT=8080 \
  -e AVAGO_STAKING_PORT=9653 \
  -e AVAGO_HTTP_ALLOWED_HOSTS="*" \
  -e AVAGO_HTTP_HOST=0.0.0.0 \
  -e AVAGO_PUBLIC_IP_RESOLUTION_SERVICE=ifconfigme \
  avaplatform/subnet-evm:v0.7.1-rc.0
```

- `AVAGO_STAKING_PORT` is set to `9653` in case you want to run this on the same machine as the validator. Remove this to set to the default `9651`.
- `AVAGO_HTTP_PORT` is set to `8080` instead of `9650` to avoid conflicts with the validator.
- `AVAGO_HTTP_ALLOWED_HOSTS` and `AVAGO_HTTP_HOST` are required to allow the RPC server to be accessed from outside. You'll need to secure it with HTTPS; Caddy is recommended.

RPC example uses another folder `~/.avalanchego_rpc` to avoid conflicts with the validator if you want to run both on the same machine.
