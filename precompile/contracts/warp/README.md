# Avalanche Warp Messaging

Avalanche Warp Messaging offers a basic primitive to enable Cross-Subnet communication on the Avalanche Network.

It is intended to allow communication between arbitrary Custom Virtual Machines (including, but not limited to Subnet-EVM).

## How does Avalanche Warp Messaging Work

Avalanche Warp Messaging relies on the organization of Avalanche Subnets through the P-Chain.

- how are Avalanche Subnets organized
- P-Chain organizes Avalanche Subnets
- Every subnet validator has read-access to the P-Chain state
- the P-Chain registers BLS Public Keys for every validator

- that's where we are, what can we do with it?

## Subnet to Subnet

- the validator set of subent A can send a message on behalf of any blockchain on its Subnet
- note: AvalancheGo will only sign a message on behalf of a VM if the SourceChainID matches the blockchainID of that blockchain (add link)
- 

## Subnet to C-Chain (Primary Network)

## C-Chain to Subnet

## Warp Precompile

## Guarantees Offered by Precompile

## Building on the Warp Precompile