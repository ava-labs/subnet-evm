//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./IAllowList.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract ExampleRewardDistributor is Ownable {
  mapping(address => uint256) rewardAddresses; // address to last claimed reward block number

  uint256 rewardRatePerBlk; // reward rate per block

  constructor(uint256 rewardRate) Ownable() {
    rewardRatePerBlk = rewardRate;
  }

  modifier _isRewardAddress(address addr) {
    require(rewardAddresses[addr] > 0, "Not a reward address");
    _;
  }

  function isRewardAddress(address addr) public view returns (bool) {
    return rewardAddresses[addr] != 0;
  }

  function getRewardRate() public view returns (uint256) {
    return rewardRatePerBlk;
  }

  function setRewardRate(uint256 rewardRate) public onlyOwner {
    rewardRatePerBlk = rewardRate;
  }

  function addRewardAddress(address addr) public onlyOwner {
    require(rewardAddresses[addr] == 0, "Already a reward address");

    rewardAddresses[addr] = block.number;
  }

  function revoke(address addr) public onlyOwner _isRewardAddress(addr) {
    uint256 reward = _estimateReward(addr);
    delete rewardAddresses[addr];

    payable(addr).transfer(reward);
  }

  function claim() public _isRewardAddress(msg.sender) {
    _reward(msg.sender);
  }

  function estimateReward(address addr) public view returns (uint256) {
    return _estimateReward(addr);
  }

  function _estimateReward(address addr) private view returns (uint256) {
    uint256 lastClaimedBlk = rewardAddresses[addr];
    if (lastClaimedBlk == 0) {
      return 0;
    }
    uint256 currentBlk = block.number;
    uint256 reward = (currentBlk - lastClaimedBlk) * rewardRatePerBlk;
    return reward;
  }

  function _reward(address addr) private {
    uint256 reward = _estimateReward(addr);
    require(reward > 0, "Nothing to claim");
    require(reward <= address(this).balance, "Not enough collected balance");

    rewardAddresses[addr] = block.number;
    payable(addr).transfer(reward);
  }
}
