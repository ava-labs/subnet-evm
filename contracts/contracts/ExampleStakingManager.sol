// SPDX-License-Identifier: MIT
// from https://solidity-by-example.org/defi/staking-rewards
pragma solidity ^0.8.24;

import "./interfaces/IWarpMessenger.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract ExampleStakingManager is Ownable {
  IWarpMessenger public constant WARP_MESSENGER = IWarpMessenger(0x0200000000000000000000000000000000000005);
  uint64 constant MAX_UINT64 = type(uint64).max;

  IERC20 public immutable stakingToken;
  IERC20 public immutable rewardsToken;

  uint256 public minStakingDuration; // Minimun duration for staking in seconds.
  uint256 public minStakingAmount; // Minimum amount for staking.
  bool public stakingEnabled; // Start time of the staking period.s
  uint256 public rewardRate; // Reward rate per second.
  // User address => rewardPerTokenStored
  mapping(address => uint256) public userRewardPerTokenPaid;
  // User address => rewards to be claimed
  mapping(address => uint256) public rewards;
  // Total staked
  uint256 public totalSupply;
  // User address => staked amount
  mapping(address => uint256) public balanceOf;
  mapping(address => uint256) public unlockedBalanceOf;

  // Validator Manager
  // subnetID + nodeID (stakingID) => messageID (RegisterValidatorMessage hash)
  mapping(bytes32 => bytes32) public activeValidators;
  // messageID => Validator
  mapping(bytes32 => Validator) public registeredValidatorMessages;

  struct Validator {
    bytes32 subnetID;
    bytes32 nodeID;
    uint64 weight;
    uint256 startedAt;
    uint256 redeemedAt;
    uint64 nonce;
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
    // invalid is true if the message represents an invalid registration
    // if a registered message considered invalid, it means the message cannot be used
    // to register a validator, including the case that the registration is already finished/expired.
    // it is required to prevent replay attacks
    bool invalid;
  }

  struct SetSubnetValidatorWeightMessage {
    bytes32 messageID;
    uint64 nonce;
    uint64 weight;
  }

  constructor(address _stakingToken, address _rewardToken, uint256 _minStakingDuration, uint256 _minStakingAmount) {
    stakingToken = IERC20(_stakingToken);
    rewardsToken = IERC20(_rewardToken);
    minStakingDuration = _minStakingDuration;
    minStakingAmount = _minStakingAmount;
  }

  function registerValidator(
    bytes32 subnetID,
    bytes32 nodeID,
    uint64 amount,
    uint64 expiryTimestamp,
    bytes memory signature
  ) external {
    // Q: check if the subnetID is this subnet?
    require(amount > minStakingAmount, "ExampleValidatorManager: amount must be greater than minStakingAmount");
    require(stakingEnabled, "Staking period has not started");
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
    require(registeredValidatorMessages[messageID].weight == 0, "ExampleValidatorManager: message already exists");

    // TODO: decide on relayer fee info
    WARP_MESSENGER.sendWarpMessage(messageBytes);

    stakingToken.transferFrom(msg.sender, address(this), amount);
    balanceOf[msg.sender] += amount;
    totalSupply += amount;

    // TODO: see if we need to store the whole message vs receive it from P-Chain (in receiveRegisterValidator)
    registeredValidatorMessages[messageID] = Validator(subnetID, nodeID, amount, 0, 0, 0, msg.sender);
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
    require(!registeredMessage.invalid, "ExampleValidatorManager: invalid message");

    // TODO: maybe we want to minimize errors here?
    Validator memory pendingValidator = registeredValidatorMessages[messageID];
    require(pendingValidator.weight != 0, "ExampleValidatorManager: pending message does not exist");

    require(pendingValidator.startedAt == 0, "ExampleValidatorManager: register message already consumed");

    bytes32 stakingID = keccak256(abi.encode(pendingValidator.subnetID, pendingValidator.nodeID));
    require(activeValidators[stakingID] == bytes32(0), "ExampleValidatorManager: validator already exists");

    activeValidators[stakingID] = messageID;
    registeredValidatorMessages[messageID].startedAt = block.timestamp;
  }

  // TODO: review error messages

  // TODO: add cooldown period for withdraw/redeem
  function removeValidator(bytes32 subnetID, bytes32 nodeID) external {
    bytes32 stakingID = keccak256(abi.encode(subnetID, nodeID));
    bytes32 messageID = activeValidators[stakingID];
    require(messageID != bytes32(0), "Validator not found");

    Validator memory validator = registeredValidatorMessages[messageID];
    require(validator.rewardAddress == msg.sender, "Only the validator can remove itself");

    require(
      block.timestamp >= validator.startedAt + minStakingDuration,
      "Cannot remove validator before min staking duration"
    );

    require(validator.redeemedAt == 0, "Validator already redeemed");

    SetSubnetValidatorWeightMessage memory message = SetSubnetValidatorWeightMessage(messageID, MAX_UINT64, 0);
    bytes memory messageBytes = abi.encode(message);
    WARP_MESSENGER.sendWarpMessage(messageBytes);
    registeredValidatorMessages[messageID].redeemedAt = block.timestamp;
  }

  function receiveRegisterMessageInvalid(uint32 messageIndex) public {
    (WarpMessage memory warpMessage, bool success) = WARP_MESSENGER.getVerifiedWarpMessage(messageIndex);
    require(success, "ExampleValidatorManager: invalid warp message");

    ValidatorRegisteredMessage memory registeredMessage = abi.decode(warpMessage.payload, (ValidatorRegisteredMessage));

    bytes32 messageID = registeredMessage.messageID;
    require(messageID != bytes32(0), "ExampleValidatorManager: invalid messageID");
    require(registeredMessage.invalid, "ExampleValidatorManager: expected an invalid message");

    Validator memory pendingValidator = registeredValidatorMessages[messageID];
    bytes32 stakingID = keccak256(abi.encode(pendingValidator.subnetID, pendingValidator.nodeID));
    delete activeValidators[stakingID];
    // if redeemedAt is 0, this was not a graceful exit.
    // TODO: do we need to check balanceOf[pendingValidator.rewardAddress] > weight?
    uint256 totalAmount = pendingValidator.weight;
    if (pendingValidator.redeemedAt != 0) {
      uint256 reward = calculateReward(
        pendingValidator.weight,
        pendingValidator.startedAt,
        pendingValidator.redeemedAt
      );
      totalAmount += reward;
    }
    unlockedBalanceOf[pendingValidator.rewardAddress] += totalAmount;
  }

  // TODO: add uptime tracking/rewards based on uptimes

  // TODO: add partial withdraw + increase stake

  // TODO: add delegation

  // Owner function to start the staking period.
  function startStaking() external onlyOwner {
    require(!stakingEnabled, "Staking has already started");
    stakingEnabled = true;
  }

  // TODO: validation vs delegation should be weighted differently for rewards
  // Function to calculate the user's reward.
  function calculateReward(uint256 amount, uint256 startedAt, uint256 finishedAt) internal view returns (uint256) {
    if (finishedAt <= startedAt || !stakingEnabled || finishedAt - startedAt < minStakingDuration) {
      return 0;
    }

    uint256 stakingTime = finishedAt - startedAt;
    return (stakingTime * rewardRate * amount) / (1 ether); // Assuming rewardRate is scaled appropriately
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
