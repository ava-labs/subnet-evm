# Shared Memory Precompile

## Goal

We want to provide a precompile that will replace the functionality of import/export transactions for both AVAX and for smart contracts to export and import an ANT.

## Calculating Asset ID

Smart contracts should have the ability to export an Avalanche Native Token directly into shared memory and to import a UTXO that exists in shared memory. A contract should only be able to interact with the `assetID` that it corresponds to.

We will calculate the `assetID` by taking the `sha256(blockchainID + contract Address)`.

## Export Functionality

### Solidity Interface for Exports

We will provide the ability for contracts to export either a generic message, a UTXO, or AVAX to another chain within the same subnet. To do this, we will expose the following interface to perform an export to a specified destination chain:

```sol
export(destinationChain bytes32, msg []byte)
```

Export will be used to perform arbitrary message passing from one chain to another. The assetID of the message will be derived from the blockchainID and the contract address.

We need to ensure that an arbitrary message cannot pretend to be an actual UTXO.

TODO:
- How do we want to ensure that an arbitrary message cannot pretend to be an actual UTXO? Add an ID or use shared memory traits? This is handled by the codec
- How should we expose the ability to add traits for a given generic message?
- We should learn from the past and not implement this until someone asks for it explicitly in case it introduces a vulnerability

```sol
exportUTXO(amount uint64, locktime uint64, threshold uint64, addrs []address)
```

ExportUTXO will allow a contract to export a UTXO with an assetID derived from the contract address. The UTXO will be placed into shared memory when the block is accepted.

TODO:
- decide what denomination to use here. If we use uint64 some contracts may need to perform the conversion.

```sol
exportAVAX(amount uint256, locktime uint64, threshold uint64, addrs []address)
```

ExportAVAX will allow a contract to export a UTXO containing AVAX from the contract address. This will require a balance check internal to the EVM and will place a UTXO into shared memory containing AVAX on accept.

TODO
- Do we want to use an amount of type uint64 or uint256? Where do we want to perform the conversion between the two denominations? 

### Exports Under the Hood: Interacting with Shared Memory

Under the hood, an export may need to modify the `StateDB` when it is executed within the EVM and it may need to perform an operation on shared memory when the block is accepted. This will consist of the following steps:

When the EVM is executed, the precompile will emit a structured log (TODO: define structured log format). For export and exportUTXO, there are no other changes needed on the `StateDB`.

For `exportAVAX`, the precompile will emit a structured log, perform a required balance check, and decrease the balance of the caller to consume the requested amount of AVAX.

Outside of the EVM's execution we need to perform a Put operation into the shared memory of the recipient chain when the block gets accepted.

When we create the UTXO, we need to ensure that the UTXO ID will be unique (invariant of shared memory that was recently documented).

## Import Functionality

Smart contracts should have the ability to import an Avalanche Native Token from shared memory that was sent by another source blockchain. This ability will enable the contract to import a generic message, import a UTXO with an assetID that corresponds to the contract, or import an AVAX UTXO from shared memory.

To add this functionality, the precompile will expose the following functions:

```sol
import(sourceChain bytes32, msg []byte) (msg []byte, present bool)
```

Import will be used to acknowledge, process, and remove a generic message from shared memory. This means that the message will be consumed by a contract by its ID and the message will be removed from shared memory when the block gets accepted.

```sol
importUTXO(utxoID bytes32) (amount uin64, locktime uint64, threshold uint64, addrs []address)
```

TODO
- there's only one caller, so there's a question of if we should only have a single address in the returned set of addresses and require a threshold of 1
- Should the caller specify the amount? Probably not

ImportUTXO will be used to acknowledge, process, and remove a UTXO from shared memory. The contract calling this precompile will be responsible for processing the details of the UTXO and crediting the owner of the UTXO with the funds held by the UTXO.

```sol
importAVAX(utxoID bytes32)
```

ImportAVAX will allow a contract to import a UTXO with AVAX as the assetID. It will verify the caller address is the sole owner of the UTXO and then increment hte AVAX balance for either the caller or a requested address. The UTXO will be removed from shared memory when the block is accepted.

TODO
- should this perform evm.Call to send funds to a requested address or just add the balance
- should the user be allowed to specify an address to send the funds to

## Special Case: Handling AVAX Import/Exports

This is a special case because normally we are going to require that the only caller address that can perform an import/export of a specific assetID, is the contract itself whose address corresponds to that assetID.

For import/export of AVAX from the precompile, we will need to enforce the requirement that the caller is specified as the only owner of the UTXO.

In contrast, for an Avalanche Native Token, we will instead return the owner details of the UTXO to the caller, which will be responsible for performing the necessary verification.

This means that it is valid for a contract calling importUTXO (not importAVAX) to import a UTXO whose owner details does not include the contract's address. Instead a wrapper contract or any other contract is responsible for handling the imported UTXO and modifying balances as required.
