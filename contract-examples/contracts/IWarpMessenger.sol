// (c) 2022-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

struct WarpMessage {
    bytes32 originChainID;
    bytes32 originSenderAddress;
    bytes32 destinationChainID;
    bytes32 destinationAddress;
    bytes payload;
}

interface WarpMessenger {
    event SendWarpMessage(
        bytes32 indexed destinationChainID,
        bytes32 indexed destinationAddress,
        bytes32 indexed sender,
        bytes message
    );

    // sendWarpMessage emits a request for the subnet to send a warp message from [msg.sender]
    // with the specified parameters.
    // This emits a SendWarpMessage log, which will be picked up by validators to queue the signing of
    // a Warp message if/when the block is accepted.
    function sendWarpMessage(
        bytes32 destinationChainID,
        bytes32 destinationAddress,
        bytes calldata payload
    ) external;

    // getVerifiedWarpMessage parses the pre-verified warp message in the
    // predicate storage slots as a WarpMessage and returns it to the caller.
    // Returns false if no such predicate exists.
    function getVerifiedWarpMessage()
        external view
        returns (WarpMessage calldata message, bool exists);

    // Note: getVerifiedWarpMessage takes no arguments because it returns a single verified
    // message that is encoded in the predicate (inside the tx access list) of the transaction.
    // The alternative design to this is to verify messages during the EVM's execution in which
    // case there would be no predicate and the block would encode the hits/misses that occur
    // throughout its execution.
    // This would result in the following alternative function signature:
    // function verifyMessage(bytes calldata signedWarpMsg) external returns (WarpMessage calldata message);

    // getBlockchainID returns the snow.Context BlockchainID of this chain.
    // This blockchainID is the hash of the transaction that created this blockchain on the P-Chain
    // and is not related to the Ethereum ChainID.
    function getBlockchainID() external view returns (bytes32 blockchainID);
}
