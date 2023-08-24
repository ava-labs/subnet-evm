//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
pragma experimental ABIEncoderV2;

<<<<<<< HEAD
import "./IWarpMessenger.sol";
=======
import "./interfaces/IWarpMessenger.sol";
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d

contract ExampleWarp {
    address constant WARP_ADDRESS = 0x0200000000000000000000000000000000000005;
    WarpMessenger warp = WarpMessenger(WARP_ADDRESS);

    // sendWarpMessage sends a warp message to the specified destination chain and address pair containing the payload
    function sendWarpMessage(
        bytes32 destinationChainID,
<<<<<<< HEAD
        bytes32 destinationAddress,
=======
        address destinationAddress,
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
        bytes calldata payload
    ) external {
        warp.sendWarpMessage(destinationChainID, destinationAddress, payload);
    }


    // validateWarpMessage retrieves the warp message attached to the transaction and verifies all of its attributes.
    function validateWarpMessage(
        bytes32 originChainID,
<<<<<<< HEAD
        bytes32 originSenderAddress,
        bytes32 destinationChainID,
        bytes32 destinationAddress,
=======
        address originSenderAddress,
        bytes32 destinationChainID,
        address destinationAddress,
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
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

    // validateGetBlockchainID checks that the blockchainID returned by warp matches the argument
    function validateGetBlockchainID(bytes32 blockchainID) external view {
        require(blockchainID == warp.getBlockchainID());
    }
}
