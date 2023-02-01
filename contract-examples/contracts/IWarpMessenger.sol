// (c) 2022-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// SPDX-License-Identifier: MIT

pragma solidity >=0.8.16;

struct CrossChainMessage {
    bytes32 sourceChainID;
    bytes32 senderAddress;
    bytes32 destinationChainID;
    bytes32 destinationAddress;
    bytes payload;
}

interface WarpMessenger {
    event SendCrossChainMessage(
        bytes32 indexed destinationChainID,
        bytes32 indexed destinationAddress,
        bytes32 indexed sender,
        bytes message
    );

    // sendCrossChainMessage emits a request for the Subnet to send a cross subnet message from [msg.sender]
    // with the specified parameters.
    // This emits a SendCrossChainMessage log, which will be picked up by validators to queue the signing of
    // a Warp message if/when the block is accepted.
    function sendCrossChainMessage(
        bytes32 destinationChainID,
        bytes32 destinationAddress,
        bytes calldata payload
    ) external;

    // getVerifiedCrossChainMessage parses the message in the predicate storage slots as a Warp Message.
    // This message is then delivered on chain by performing evm.Call with the Warp precompile as the caller,
    // the destinationAddress as the receiver.
    // The full message and a boolean indicating if the operation executed successfully is returned to the caller.
    function getVerifiedCrossChainMessage(uint256 messageIndex)
        external
        returns (CrossChainMessage calldata message, bool success);

    function getBlockchainId() external view returns (bytes32 chainID);
}