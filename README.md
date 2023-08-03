# Subnet EVM

[![Build + Test + Release](https://github.com/ava-labs/subnet-evm/actions/workflows/lint-tests-release.yml/badge.svg)](https://github.com/ava-labs/subnet-evm/actions/workflows/lint-tests-release.yml)
[![CodeQL](https://github.com/ava-labs/subnet-evm/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/ava-labs/subnet-evm/actions/workflows/codeql-analysis.yml)

[Avalanche](https://docs.avax.network/overview/getting-started/avalanche-platform) is a network composed of multiple blockchains.
Each blockchain is an instance of a Virtual Machine (VM), much like an object in an object-oriented language is an instance of a class.
That is, the VM defines the behavior of the blockchain.

Subnet EVM is the [Virtual Machine (VM)](https://docs.avax.network/learn/avalanche/virtual-machines) that defines the Subnet Contract Chains. Subnet EVM is a simplified version of [Coreth VM (C-Chain)](https://github.com/ava-labs/coreth).

This chain implements the Ethereum Virtual Machine and supports Solidity smart contracts as well as most other Ethereum client functionality.

## Building

The Subnet EVM runs in a separate process from the main AvalancheGo process and communicates with it over a local gRPC connection.

### AvalancheGo Compatibility

```text
[v0.1.0] AvalancheGo@v1.7.0-v1.7.4 (Protocol Version: 9)
[v0.1.1-v0.1.2] AvalancheGo@v1.7.5-v1.7.6 (Protocol Version: 10)
[v0.2.0] AvalancheGo@v1.7.7-v1.7.9 (Protocol Version: 11)
[v0.2.1] AvalancheGo@v1.7.10 (Protocol Version: 12)
[v0.2.2] AvalancheGo@v1.7.11-v1.7.12 (Protocol Version: 14)
[v0.2.3] AvalancheGo@v1.7.13-v1.7.16 (Protocol Version: 15)
[v0.2.4] AvalancheGo@v1.7.13-v1.7.16 (Protocol Version: 15)
[v0.2.5] AvalancheGo@v1.7.13-v1.7.16 (Protocol Version: 15)
[v0.2.6] AvalancheGo@v1.7.13-v1.7.16 (Protocol Version: 15)
[v0.2.7] AvalancheGo@v1.7.13-v1.7.16 (Protocol Version: 15)
[v0.2.8] AvalancheGo@v1.7.13-v1.7.18 (Protocol Version: 15)
[v0.2.9] AvalancheGo@v1.7.13-v1.7.18 (Protocol Version: 15)
[v0.3.0] AvalancheGo@v1.8.0-v1.8.6 (Protocol Version: 16)
[v0.4.0] AvalancheGo@v1.9.0 (Protocol Version: 17)
[v0.4.1] AvalancheGo@v1.9.1 (Protocol Version: 18)
[v0.4.2] AvalancheGo@v1.9.1 (Protocol Version: 18)
[v0.4.3] AvalancheGo@v1.9.2-v1.9.3 (Protocol Version: 19)
[v0.4.4] AvalancheGo@v1.9.2-v1.9.3 (Protocol Version: 19)
[v0.4.5] AvalancheGo@v1.9.4 (Protocol Version: 20)
[v0.4.6] AvalancheGo@v1.9.4 (Protocol Version: 20)
[v0.4.7] AvalancheGo@v1.9.5 (Protocol Version: 21)
[v0.4.8] AvalancheGo@v1.9.6-v1.9.8 (Protocol Version: 22)
[v0.4.9] AvalancheGo@v1.9.9 (Protocol Version: 23)
[v0.4.10] AvalancheGo@v1.9.9 (Protocol Version: 23)
[v0.4.11] AvalancheGo@v1.9.10-v1.9.16 (Protocol Version: 24)
[v0.4.12] AvalancheGo@v1.9.10-v1.9.16 (Protocol Version: 24)
[v0.5.0] AvalancheGo@v1.10.0 (Protocol Version: 25)
[v0.5.1] AvalancheGo@v1.10.1-v1.10.4 (Protocol Version: 26)
[v0.5.2] AvalancheGo@v1.10.1-v1.10.4 (Protocol Version: 26)
[v0.5.3] AvalancheGo@v1.10.5-v1.10.6 (Protocol Version: 27)
[v0.5.4] AvalancheGo@v1.10.5-v1.10.6 (Protocol Version: 27)
```

## API

The Subnet EVM supports the following API namespaces:

- `eth`
- `personal`
- `txpool`
- `debug`

Only the `eth` namespace is enabled by default.
Full documentation for the C-Chain's API can be found [here.](https://docs.avax.network/apis/avalanchego/apis/c-chain)

## Compatibility

The Subnet EVM is compatible with almost all Ethereum tooling, including [Remix](https://docs.avax.network/dapps/smart-contracts/deploy-a-smart-contract-on-avalanche-using-remix-and-metamask/), [Metamask](https://docs.avax.network/dapps/smart-contracts/deploy-a-smart-contract-on-avalanche-using-remix-and-metamask/) and [Truffle](https://docs.avax.network/dapps/smart-contracts/using-truffle-with-the-avalanche-c-chain/).

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

First install Go 1.19.6 or later. Follow the instructions [here](https://golang.org/doc/install). You can verify by running `go version`.

Set `$GOPATH` environment variable properly for Go to look for Go Workspaces. Please read [this](https://go.dev/doc/gopath_code) for details. You can verify by running `echo $GOPATH`.

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

To run a local network, it is recommended to use the [avalanche-cli](https://github.com/ava-labs/avalanche-cli#avalanche-cli) to set up an instance of Subnet-EVM on an local Avalanche Network.

There are two options when using the Avalanche-CLI:

1. Use an official Subnet-EVM release: https://docs.avax.network/subnets/build-first-subnet
2. Build and deploy a locally built (and optionally modified) version of Subnet-EVM: https://docs.avax.network/subnets/create-custom-subnet


### Run subnet-evm on AWS

[avalanche-ops](https://github.com/ava-labs/avalanche-ops) is a suite of tools to automate deploying a custom subnet on AWS. Launching a subnet, and a chain running subnet-evm, is straightforward using avalanche-ops. 

1. Clone [avalanche-ops](https://github.com/ava-labs/avalanche-ops) and build the project (using `--release`), or download the `avalancheup-aws` binary for your specific platform from the releases page. Using a release version is recommended.
2. Login to AWS locally, via the `aws` CLI tool. See [this doc](https://github.com/ava-labs/avalanche-ops#permissions) to login if using an authentication provider. In particular, your AWS session should have a "profile name", which is then passed directly to avalanche-ops.
3. Navigate to the `avalancheup-aws` binary, which is either in the local installation directory or `avalanche-ops/target/release` if you build the project locally.
4. Run the following command to generate a spec (a template for the cloud resources that will be created for you by `avalanche-ops`)
```bash
$ avalancheup-aws default-spec \
--arch-type amd64 \
--os-type ubuntu20.04 \
--anchor-nodes 1 \
--non-anchor-nodes 1 \
--instance-mode=on-demand \
--instance-size=2xlarge \
--ip-mode=elastic \
--ingress-ipv4-cidr 0.0.0.0/0 \
--network-name custom \
--keys-to-generate 10 \
--profile-name <YOUR AWS PROFILE NAME>
```
You can modify this spec as you like by adding optional flags and selecting different values.
The number of `anchor-nodes` is the number of seed/bootstrapper nodes.
The number of `non-anchor-nodes` is the number of existing anchor nodes for peer discovery.
The default value the optional flag `regions` is `us-west-2`. 
5. Once the spec is ready, run the `apply` command to spin up a subnet. 
```bash
$ avalancheup-aws apply \
--spec-file-path <SPEC FILE PATH>
```

6. Once the subnet has successfully installed, create a new chain using the subnet-evm.

First, generate the necessary subnet-evm config files locally:
```bash
$ avalancheup-aws subnet-evm chain-config -s /tmp/subnet-evm-genesis.json
``` 
```bash
$ avalancheup-aws subnet-evm genesis -s /tmp/subnet-evm-genesis.json
```
You can modify these files as well per the desired subnet-evm configuration. 

7. Install the chain via the following command. Substitute the required parameters as needed. This command is also outputted at the end of the `apply` command, so copy/paste if possible.

The only required fields to update are the `profile-name` and the `vm-binary-local-path`. The VM binary local path is the path to the local subnet-evm binary. The subnet-evm binary must be compiled for linux amd64. Do not rename the subnet-evm binary (use the original hash name).
```yaml
$ avalancheup-aws install-subnet-chain \
--log-level info \
--profile-name <AWS PROFILE NAME> 
--s3-region us-west-2 \
--s3-bucket avalanche-ops-202307-4lw5xroxbc-us-west-2 \
--s3-key-prefix aops-custom-202307-2fpVNd/install-subnet-chain \
--chain-rpc-url http://54.203.8.171:9650 \
--key 0x56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027 \
--subnet-config-local-path /tmp/subnet-config.json \
--subnet-config-remote-dir /data/avalanche-configs/subnets \
--vm-binary-local-path <PATH TO LOCAL SUBNET-EVM COMPILED FOR linux-amd64> \
--vm-binary-remote-dir /data/avalanche-plugins \
--chain-name subnetevm \
--chain-genesis-path /tmp/subnet-evm-genesis.json \
--chain-config-local-path /tmp/subnet-evm-chain-config.json \
--chain-config-remote-dir /data/avalanche-configs/chains \
--avalanchego-config-remote-path /data/avalanche-configs/config.json \
--ssm-docs '{"us-west-2":"aops-custom-202307-2fpVNd-ssm-install-subnet-chain"}' \
--target-nodes '{"NodeID-6RCgWDrGDWt8vPFb8iTw2qp9287ae2XjD":{"region":"us-west-2","machine_id":"i-058653a153723d60a"},"NodeID-2bVWnx4sFGfzCTi8ntTHcdTKc7KJbGUKi":{"region":"us-west-2","machine_id":"i-08417d553a2f18491"}}'
```

You could define these optional flags for these values: `s3-upload-timeout = 30`, `primary-network-validate-period-in-days = 16`, `subnet-validate-period-in-days = 14`, and
`staking-amount-in-avax = 2000`.
8. Once the `install-subnet-chain` commands returns successfully, log in to the AWS console and you should be able to see the cloud resources that were created. You can use the JSON-RPC APIs and query nodes as expected.