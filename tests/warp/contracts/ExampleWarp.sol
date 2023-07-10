//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
pragma experimental ABIEncoderV2;

import "./IWarpMessenger.sol";

address constant WARP_ADDRESS = 0x0200000000000000000000000000000000000005;

contract ExampleWarp {

    IWarpMessenger warp = IWarpMessenger(WARP_ADDRESS);

    // sendWarpMessage sends a warp message to the specified destination chain and address pair containing the payload
    function sendWarpMessage(
        bytes32 destinationChainID,
        bytes32 destinationAddress,
        bytes calldata payload
    ) external {
        warp.sendWarpMessage(destinationChainID, destinationAddress, payload);
    }

    // validateWarpMessage retrieves the warp message attached to the transaction and verifies all of its attributes.
    function validateWarpMessage(
        bytes32 originChainID,
        bytes32 originSenderAddress,
        bytes32 destinationChainID,
        bytes32 destinationAddress,
        bytes calldata payload
    ) external view {
        (WarpMessage memory message, bool exists) = warp.getVerifiedWarpMessage();
        require(exists);
        require(message.originChainID == originChainID);
        require(message.originSenderAddress == originSenderAddress);
        require(message.destinationChainID == destinationChainID);
        require(message.destinationAddress == destinationAddress);
        require(keccak256(message.payload) == keccak256(payload));
    }

    function validateGetBlockchainID(
        bytes32 expectedBlockchainID
    ) external view {
        bytes32 blockchainID = warp.getBlockchainID();
        require(blockchainID==expectedBlockchainID);
    }
}
