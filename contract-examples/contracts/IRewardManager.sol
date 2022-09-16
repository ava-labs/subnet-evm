//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "./IAllowList.sol";

interface IRewardManager is IAllowList {
  function setRewardAddress(address addr) external;

  function allowFeeRecipients() external;

  function disableRewards() external;

  function currentRewardAddress() external returns (address rewardAddress);

  function areFeeRecipientsAllowed() external returns (bool isAllowed);
}
