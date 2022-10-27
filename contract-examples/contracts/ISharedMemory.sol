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
// When a block g
interface ISharedMemory {
    // TransferOutput is emitted by the contract when amount of assetID was exported from the chain
    // to be owned by addrs/threshold/locktime.
    event TransferOutput(bytes32 otherChain, bytes32 assetID, uint64 amount, uint64 locktime, uint64 threshold, address[] addrs);

    // TransferInput is emitted by the contract when it spends utxoID.
    event TransferInput(bytes32 otherChain, bytes32 assetID, bytes32 utxoID);

    // getNativeTokenAssetID returns the assetID that corresponds to the specified caller.
    // The returned assetID is sha256(caller, blockchainID) where the blockchainID is the Avalanche blockchainID as opposed to the EVM
    // ChainID.
    function getNativeTokenAssetID(address caller) external view returns (bytes32 assetID);

    // exportUTXO generates a UTXO with a unique UTXOID with the specified amount, locktime, threshold, and set of addresses
    // and an assetID derived from msg.sender and the blockchainID
    function exportUTXO(uint64 amount, bytes32 desinationChainID, uint64 locktime, uint64 threshold, address[] calldata addrs) external;

    // exportAVAX generates an AVAX UTXO with a unique UTXOID with the specified locktime, threshold, and set of addresses
    // and the AVAX assetID.
    // In order to ensure that we do not break any EVM invariants, we require that the caller specify the amount of AVAX to
    // use for the export as msg.value
    // XXX we do not include an amount paremter to avoid specifying it twice. We may consider specifying the parameter and verifying
    // that the amount matches msg.value
    function exportAVAX(uint64 locktime, uint64 threshold, address[] calldata addrs) external;

    // importUTXO attempts to import the UTXO specified by UTXOID. If the UTXO is still present, then it returns the UTXO details to the caller
    // for the caller to credit any balance changes as a result of importing the UTXO.
    // how do we verify msg.sender
    // when we look up the UTXO, we want to verify that the sender can actually use it
    // The precompile will perform the following steps:
    // 1. Verify that the UTXO is available within the TxContext (the TxContext will verify that the UTXOs as specfied in the transaction are available in shared memory)
    // 2. Verify that 
    // The precompile will look up the UTXO by sourceChain and UTXOID from shared memory.

    // If the atomic input does not match
    function importUTXO(bytes32 sourceChain, bytes32 utxoID) external returns (uint64 amount, uint64 locktime, uint64 threshold, address[] calldata addrs);
}

// SharedMemory requires additional steps for block verification in order to ensure that atomic UTXO conflicts are handled correctly across a chain of
// processing blocks.
// Currently, atomic UTXOs are only spent by import transactions processed at the end of the block, which specify the exact UTXOs that they spend. If any
// of these transactions are invalid or include a double spend, the block is considered invalid.
// 

// TODO
// write out block verification invariants
// write out the steps the precompile will actually take for each of these functions
