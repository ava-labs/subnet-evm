//SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;
import "./IAllowList.sol";

interface IStakingManager {
  function registerValidator(
    bytes32 subnetID,
    bytes32 nodeID,
    uint64 amount,
    uint64 expiryTimestamp,
    bytes memory signature
  ) external;

  function receiveValidatorRegistered(uint32 messageIndex) external;

  function removeValidator(bytes32 subnetID, bytes32 nodeID) external;

  function receiveRegisterMessageInvalid(uint32 messageIndex) external;

  function increaseValidatorWeight(bytes32 subnetID, bytes32 nodeID, uint64 amount) external;

  function decreaseValidatorWeight(bytes32 subnetID, bytes32 nodeID, uint64 amount) external;

  function receiveValidatorWeightChanged(uint32 messageIndex) external;
}
