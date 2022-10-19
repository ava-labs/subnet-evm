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
    // getNativeTokenAssetID returns the assetID that corresponds to the specified caller.
    // The returned assetID is sha256(caller, blockchainID) where the blockchainID is the Avalanche blockchainID as opposed to the EVM
    // ChainID.
    function getNativeTokenAssetID(address caller) external view returns (bytes32 assetID);

    // exportUTXO generates a UTXO with a unique UTXOID with the specified amount, locktime, threshold, and set of addresses
    // and an assetID derived from msg.sender and the blockchainID
    function exportUTXO(uint64 amount, uint64 locktime, uint64 threshold, address[] calldata addrs) external;

    // exportAVAX generates an AVAX UTXO with a unique UTXOID with the specified locktime, threshold, and set of addresses
    // and the AVAX assetID.
    // In order to ensure that we do not break any EVM invariants, we require that the caller specify the amount of AVAX to
    // use for the export as msg.value
    // XXX we do not include an amount paremter to avoid specifying it twice. We may consider specifying the parameter and verifying
    // that the amount matches msg.value
    function exportAVAX(uint64 locktime, uint64 threshold, address[] calldata addrs) external;

    // importUTXO attempts to import the UTXO specified by UTXOID. If the UTXO is still present, then it returns the UTXO details to the caller
    // for the caller to credit any balance changes as a result of importing the UTXO.
    function importUTXO(bytes32 utxoID) external returns (uint64 amount, uint64 locktime, uint64 threshold, address[] calldata addrs);
}
