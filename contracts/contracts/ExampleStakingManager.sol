// SPDX-License-Identifier: MIT
// from https://solidity-by-example.org/defi/staking-rewards
pragma solidity ^0.8.24;

import "./interfaces/IWarpMessenger.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract ExampleStakingManager is Ownable {
  IWarpMessenger public constant WARP_MESSENGER = IWarpMessenger(0x0200000000000000000000000000000000000005);

  IERC20 public stakingToken; // The stakingToken users will stake.
  uint256 public stakingDuration; // Duration for staking in seconds.
  uint256 public rewardAmount; // Reward amount users will receive.
  uint256 public startTime; // Start time of the staking period.

  mapping(address => uint256) public stakedBalance;
  mapping(address => uint256) public lastUpdateTime;

  // Validator Manager
  // subnetID + nodeID (stakingID) => messageID (RegisterValidatorMessage hash)
  mapping(bytes32 => bytes32) public activeValidators;
  // messageID => Validator
  mapping(bytes32 => Validator) public registeredValidatorMessages;

  struct Validator {
    bytes32 subnetID;
    bytes32 nodeID;
    uint64 weight;
    address rewardAddress;
  }

  struct RegisterValidatorMessage {
    bytes32 subnetID;
    bytes32 nodeID;
    uint64 weight;
    uint64 expiryTimestamp;
    bytes signature;
  }

  struct ValidatorRegisteredMessage {
    bytes32 messageID;
  }

  constructor(address _stakingToken, uint256 _stakingDuration, uint256 _rewardAmount) {
    stakingToken = IERC20(_stakingToken);
    stakingDuration = _stakingDuration;
    rewardAmount = _rewardAmount;
  }

  function registerValidator(
    bytes32 subnetID,
    bytes32 nodeID,
    uint64 amount,
    uint64 expiryTimestamp,
    bytes memory signature
  ) external {
    // Q: check if the subnetID is this subnet?
    require(amount > 0, "ExampleValidatorManager: amount must be greater than 0");
    require(block.timestamp >= startTime, "Staking period has not started");
    // Q: does below check make sense?
    require(expiryTimestamp > block.timestamp, "ExampleValidatorManager: expiry timestamp must be in the future");
    require(nodeID != bytes32(0), "ExampleValidatorManager: nodeID must not be zero");
    require(signature.length == 64, "ExampleValidatorManager: invalid signature length, must be 64");

    bytes32 stakingIDHash = keccak256(abi.encode(subnetID, nodeID));
    require(activeValidators[stakingIDHash] == bytes32(0), "ExampleValidatorManager: validator already exists");

    // TODO: do we need to specify allowed relayers?
    RegisterValidatorMessage memory message = RegisterValidatorMessage(
      subnetID,
      nodeID,
      amount,
      expiryTimestamp,
      signature
    );

    bytes memory messageBytes = abi.encode(message);
    bytes32 messageID = sha256(messageBytes);
    // This requires the message ID on P-Chain to be same as this message ID.
    require(
      registeredValidatorMessages[messageID].weight == 0,
      "ExampleValidatorManager: pending message already exists"
    );

    // TODO: decide on relayer fee info
    WARP_MESSENGER.sendWarpMessage(messageBytes);

    stakingToken.transferFrom(msg.sender, address(this), amount);
    stakedBalance[msg.sender] += amount;

    // TODO: see if we need to store the whole message vs receive it from P-Chain (in )
    registeredValidatorMessages[messageID] = Validator(subnetID, nodeID, amount, msg.sender);
  }

  function receiveRegisterValidator(uint32 messageIndex) public {
    (WarpMessage memory warpMessage, bool success) = WARP_MESSENGER.getVerifiedWarpMessage(messageIndex);
    require(success, "ExampleValidatorManager: invalid warp message");

    // TODO: check if the sender is P-Chain
    // require(warpMessage.sourceChainID == P_CHAIN_ID, "ExampleValidatorManager: invalid source chain ID");
    // require(warpMessage.originSenderAddress == address(this), "ExampleValidatorManager: invalid origin sender address");

    // Parse the payload of the message.
    ValidatorRegisteredMessage memory registeredMessage = abi.decode(warpMessage.payload, (ValidatorRegisteredMessage));

    bytes32 messageID = registeredMessage.messageID;
    require(messageID != bytes32(0), "ExampleValidatorManager: invalid messageID");

    Validator memory pendingValidator = registeredValidatorMessages[messageID];
    require(pendingValidator.weight != 0, "ExampleValidatorManager: pending message does not exist");

    bytes32 stakingID = keccak256(abi.encode(pendingValidator.subnetID, pendingValidator.nodeID));
    require(activeValidators[stakingID] == bytes32(0), "ExampleValidatorManager: validator already exists");

    activeValidators[stakingID] = messageID;
    lastUpdateTime[pendingValidator.rewardAddress] = block.timestamp;
  }

  // TODO: add delegation

  // Function to calculate the user's reward.
  // TODO: this is totally broken, replace this a proper reward calculation
  // add MinStakeDuration, reward based on staked amount, etc.
  function calculateReward(address user) internal view returns (uint256) {
    uint256 lastUpdate = lastUpdateTime[user];
    if (lastUpdate == 0 || lastUpdate >= startTime + stakingDuration) {
      return 0;
    }

    uint256 stakingTime = block.timestamp - lastUpdate;
    return (stakingTime * rewardAmount) / stakingDuration;
  }

  // Owner function to start the staking period.
  function startStaking() external onlyOwner {
    require(startTime == 0, "Staking has already started");
    startTime = block.timestamp;
  }
}

interface IERC20 {
  function totalSupply() external view returns (uint256);
  function balanceOf(address account) external view returns (uint256);
  function transfer(address recipient, uint256 amount) external returns (bool);
  function allowance(address owner, address spender) external view returns (uint256);
  function approve(address spender, uint256 amount) external returns (bool);
  function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);
}
