//SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;
import "./IAllowList.sol";

interface IStakingManager {
  // RegisterValidatorMessage is the message that is sent by the staking manager contract to the P-Chain as a warp message
  // to register a validator on a subnet.
  // A messageID can be calculated by hashing (sha256) the RegisterValidatorMessage.
  struct RegisterValidatorMessage {
    bytes32 subnetID;
    bytes32 nodeID;
    uint64 weight;
    uint64 expiryTimestamp;
    bytes signature;
  }

  // ValidatorRegisteredMessage is the message that is sentto the staking manager contract as a warp message
  // that confirms the registration of a validator on P-chain. The messageID is the sha256 of the RegisterValidatorMessage.
  struct ValidatorRegisteredMessage {
    // The messageID is the sha256 of the related RegisterValidatorMessage.
    bytes32 messageID;
  }

  // InvalidValidatorRegisterMessage is the message that is sent to the staking manager contract as a warp message
  // that indicates that the registration of a validator on P-chain was invalid or will forever be invalid (validation finished).
  struct InvalidValidatorRegisterMessage {
    // The messageID is the sha256 of the related RegisterValidatorMessage.
    bytes32 messageID;
  }

  // SetSubnetValidatorWeightMessage is the message that is sent by the staking manager contract to the P-Chain as a warp message
  // to set the weight of a validator on a subnet.
  struct SetSubnetValidatorWeightMessage {
    // The messageID is the sha256 of the related RegisterValidatorMessage.
    bytes32 messageID;
    uint64 nonce;
    uint64 weight;
  }

  // ValidatorWeightChangedMessage is the message that is sent to the staking manager contract as a warp message
  // that confirms the change of weight of a validator on P-chain.
  struct ValidatorWeightChangedMessage {
    // The messageID is the sha256 of the related RegisterValidatorMessage.
    bytes32 messageID;
    uint64 nonce;
    uint64 weight;
  }

  // UptimeMessage is the warp message that can be received by the staking manager contract to evaluate the uptime of a validator.
  // This can be used to calculate the reward for a validator.
  struct UptimeMessage {
    // The messageID is the sha256 of the related RegisterValidatorMessage.
    bytes32 messageID;
    uint256 from;
    uint256 to;
    uint256 uptime;
  }

  // SetSubnetValidatorManager is the warp message that can be sent by the staking manager contract to the P-Chain
  // to change the validator manager of a subnet.
  // Note: this can only be done if the existing validator manager is the contract that is sending the message.
  // The first validator manager must be registered manually (not via a contract) on the P-chain.
  struct SetSubnetValidatorManager {
    bytes32 subnetID;
    bytes32 chainID;
    address validatorManager;
  }

  // registerValidator creates and sends a warp message for P-Chain to register a validator on a subnet.
  // A successfull transaction will not imply that the validator is registered. See receiveValidatorRegistered.
  // The contract sending the message must be a registered staking contract on the P-chain via SetSubnetValidatorManagerTx.
  function registerValidator(
    bytes32 subnetID,
    bytes32 nodeID,
    uint64 amount,
    uint64 expiryTimestamp,
    bytes memory signature
  ) external;

  // removeValidator creates and sends a warp message for P-Chain to remove a validator from a subnet.
  // A successfull transaction will not imply that the validator is removed. See receiveRegisterMessageInvalid.
  function removeValidator(bytes32 subnetID, bytes32 nodeID) external;

  // increaseValidatorWeight creates and sends a warp message for P-Chain to increase the weight of a validator on a subnet.
  // A successfull transaction will not imply that the weight is increased. See receiveValidatorWeightChanged.
  function increaseValidatorWeight(bytes32 subnetID, bytes32 nodeID, uint64 amount) external;

  // decreaseValidatorWeight creates and sends a warp message for P-Chain to decrease the weight of a validator on a subnet.
  // A successfull transaction will not imply that the weight is decreased. See receiveValidatorWeightChanged.
  function decreaseValidatorWeight(bytes32 subnetID, bytes32 nodeID, uint64 amount) external;

  // receiveValidatorRegistered is called with a verified warp messageIndex to confirm the registration of a validator on P-chain.
  // The warp payload corresponds to messageIndex must be a ValidatorRegisteredMessage.
  function receiveValidatorRegistered(uint32 messageIndex) external;

  // receiveValidatorWeightChanged is called with a verified warp messageIndex to confirm the change of weight of a validator on P-chain.
  // The warp payload corresponds to messageIndex must be a ValidatorWeightChangedMessage.
  function receiveValidatorWeightChanged(uint32 messageIndex) external;

  // receiveRegisterMessageInvalid is called with a verified warp messageIndex to confirm that the registration of a validator on P-chain was invalid
  // or will forever be invalid (validation finished).
  function receiveRegisterMessageInvalid(uint32 messageIndex) external;

  // receiveUptimeMessage is called with a verified warp messageIndex to report the confirmed uptime of a validator between two timestamps.
  // The warp payload corresponds to messageIndex must be a UptimeMessage.
  function receiveUptimeMessage(uint32 messageIndex) external;

  // setSubnetValidatorManager creates and sends a warp message for P-Chain to change the validator manager of a subnet.
  // This can only be done if the existing validator manager is the contract that is sending the message.
  function setSubnetValidatorManager(bytes32 subnetID, bytes32 chainID, address validatorManager) external;

  // TODO: add delegation
}
