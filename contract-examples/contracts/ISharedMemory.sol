//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// ISharedMemory provides the interface for a shared memory precompile, which enables smart contracts to export UTXOs into Avalanche shared memory:
// https://github.com/ava-labs/avalanchego/blob/master/chains/atomic/README.md.
//
// This precompile provides the following functionality:
// 1. Allow a given address to export a UTXO into shared memory with an assetID derived from its address
// 2. Allow a given address to import a UTXO from shared memory when that UTXO's assetID matches the assetID derived from the caller address
// 3. Allow a caller to import an AVAX UTXO when that address is specified as the sole owner of the UTXO
// 4. Allow a caller to export an AVAX UTXO by sending AVAX through the exportAVAX function
//
// For each operation that is actually performed on shared memory, we cannot perform the operation until the block containing the action
// is accepted. Therefore, this precompile does not directly perform Put or Remove requests on shared memory, but instead emits log events
// and the details of the action that was performed.
//
// SharedMemory requires additional steps for block verification in order to ensure that atomic UTXO conflicts are handled correctly across a chain of
// processing blocks.
// Currently, atomic UTXOs are only spent by import transactions processed at the end of the block, which specify the exact UTXOs that they spend. If any
// of these transactions are invalid or include a double spend, the block is considered invalid.
// To include shared memory inside of the EVM, we need to provide a clear way to ensure that if we ever re-process a block, it observes a static state
// of shared memory, so that it can be deterministically re-executed.
// Therefore, we add the following steps to block verification prior to executing the block:
//
// 1. Iterate the EVM transactions in the block for all EVM transactions that include an import and create a set of UTXOs
// that may be imported (these transactions may contain duplicates, which will be handled within the EVM)
// 2. Verify that the full set of importUTXOs is available in the current state of shared memory
// 3. Verify that there are no conflicts with the specified importUTXOs set with the processing ancestry of the new block
// 4. Add the importedUTXO set into the BlockContext as the set of available UTXOs throughout the EVM's execution.
// 5. Add the importedUTXO set named by each transaction into its own TxContext
//
// The EVM is now responsible to ensure that any imported UTXO is available in both the BlockContext and TxContext when it is consumed
// and must handle conflicts within the block.
interface ISharedMemory {
    // IDEA:
    // should we provide an interface that exposes the ability for a contract to export multiple tokens
    // we could include an additional assetID parameter, which would just be a number where 0 specifies
    // AVAX and then 1...n specifies additional assets that could be managed by the contract.
    // This would mean changing the getNativeTokenAssetID to take an additional parameter.
    // This wold allow us to unify the interfaces for AVAX and UTXOs.

    // getNativeTokenAssetID returns the assetID that corresponds to the specified caller.
    // The returned assetID is sha256(caller, blockchainID) where the blockchainID is the Avalanche blockchainID as opposed to the EVM
    // ChainID.
    function getNativeTokenAssetID(address caller) external view returns (bytes32 assetID);

    // ExportUTXO is emitted by exportUTXO to indicate that the export operation has taken place.
    // When the block is accepted, the VM will parse the generated log (if it was not reverted) and perform the export
    // operation on shared memory.
    // TODO: should we use assetID or the contract address that corresponds to the assetID
    event ExportUTXO(uint64 amount, bytes32 indexed destinationChainID, bytes32 indexed assetID, uint64 locktime, uint64 threshold, address[] addrs);

    // exportUTXO generates a UTXO with a unique UTXOID with the specified amount, locktime, threshold, and set of addresses
    // and an assetID derived from msg.sender and the blockchainID
    // This emits an ExportUTXO event.
    function exportUTXO(uint64 amount, bytes32 desinationChainID, uint64 locktime, uint64 threshold, address[] calldata addrs) external;

    // ExportAVAX is emitted by exportAVAX to indicate that the export AVAX operation has taken place.
    // When the block is accepted, the VM will parse the generated log (if it was not reverted) and perform the export
    // operation on shared memory.
    // TODO: should we use assetID or the contract address that corresponds to the assetID
    event ExportAVAX(uint64 amount, bytes32 indexed destinationChainID, uint64 locktime, uint64 threshold, address[] addrs);

    // exportAVAX generates an AVAX UTXO with a unique UTXOID with the specified locktime, threshold, and set of addresses
    // and the AVAX assetID.
    // In order to ensure that we do not break any EVM invariants, we require that the caller specify the amount of AVAX to
    // use for the export as msg.value
    // XXX we do not include an amount paremter to avoid specifying it twice. We may consider specifying the parameter and verifying
    // that the amount matches msg.value
    function exportAVAX(bytes32 destinationChainID, uint64 locktime, uint64 threshold, address[] calldata addrs) external;

    // ImportUTXO is emitted by importUTXO to indicate that the import operation has taken place.
    // When the block is accepted, the VM will parse the generated log (if it was not reverted) and perform the import
    // operation on shared memory.
    // TODO: should we use assetID or the contract address that corresponds to the assetID
    event ImportUTXO(uint64 amount, bytes32 indexed sourceChainID, bytes32 indexed assetID, uint64 locktime, uint64 threshold, address[] addrs);

    // importUTXO attempts to import the UTXO specified by UTXOID. If the UTXO is available, then it returns the UTXO details to the caller
    // for the caller to credit any balance changes as a result of importing the UTXO.
    // importUTXO performs the following verification:
    // 1. Verify the UTXO is available in the BlockContext and TxContext
    // 2. Verify the UTXO has not already been spent at this point in the EVM's execution
    // 3. Verify that the assetID of the UTXO matches the assetID that corresponds to msg.sender
    //
    // If verification passes, then importUTXO emits an ImportUTXO event, which will trigger a Remove request on shared memory during block
    // acceptance.
    function importUTXO(bytes32 sourceChain, bytes32 utxoID) external returns (uint64 amount, uint64 locktime, uint64 threshold, address[] calldata addrs);

    // ImportAVAX is emitted by importAVAX to indicate that the import operation has taken place.
    // When the block is accepted, the VM will parse the generated log (if it was not reverted) and perform the import
    // operation on shared memory.
    event ImportAVAX(uint64 amount, bytes32 indexed sourceChainID, uint64 locktime, uint64 threshold, address[] addrs);

    // importAVAX attempts to import the AVAX UTXO specified by UTXOID.
    // importAVAX performs the following verification:
    // 1. Verify the UTXO is available in the BlockContext and TxContext
    // 2. Verify the UTXO has not already been spent at this point in the EVM's execution
    // 3. Verify that the UTXO's assetID is AVAX
    // 4. Verify that the multisig specifies a single address msg.sender with a threshold of 1.
    // 5. Verify that block.timestamp is after locktime
    //
    // If verification passes, then importAVAX emits an ImportAVAX event, which will trigger a Remove request on shared memory during block
    // acceptance.
    // Since this imports AVAX to the chain, this will additionally credit msg.sender with the imported amount of AVAX.
    // There are two options to do this:
    // 1. AddBalance directly (this breaks an invariant of the EVM that your balance cannot be modified during a call other than through a call operation)
    // 2. Call msg.sender with the amount of the UTXO and no calldata.
    // XXX
    function importAVAX(bytes32 sourceChain, bytes32 utxoID) external;
}
