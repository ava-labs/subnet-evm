//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./IAllowList.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

// ExampleRewardDistributor is a sample contract to be used in conjunction
// with the RewardManager precompile. This contract allows its owner to
// add reward addresses, each of which can claim fees that accumulate to the
// contract (up to rewardRate per block). The owner can also adjust rewardRate.
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

  // setRewardRate sets the reward rate per block
  function setRewardRate(uint256 rewardRate) public onlyOwner {
    rewardRatePerBlk = rewardRate;
  }

  // addRewardAddress adds an address to the reward list
  // and sets the last claimed block to the current block
  function addRewardAddress(address addr) public onlyOwner {
    require(rewardAddresses[addr] == 0, "Already a reward address");

    rewardAddresses[addr] = block.number;
  }

  // revoke removes address from reward addresses and transfers any remaining rewards to the address
  function revoke(address addr) public onlyOwner _isRewardAddress(addr) {
    uint256 reward = estimateReward(addr);
    delete rewardAddresses[addr];

    // reward only if there is any to reward
    if (reward > 0 && reward <= address(this).balance) {
      payable(addr).transfer(reward);
    }
  }

  // claim transfers any rewards to the address
  function claim() public _isRewardAddress(msg.sender) {
    uint256 reward = estimateReward(addr);
    require(reward <= address(this).balance, "Not enough collected balance");

    rewardAddresses[addr] = block.number;
    payable(addr).transfer(reward);
  }

  // estimateReward returns the estimated reward for the address
  function estimateReward(address addr) public view returns (uint256) {
    uint256 lastClaimedBlk = rewardAddresses[addr];
    if (lastClaimedBlk == 0) {
      return 0;
    }
    uint256 currentBlk = block.number;
    uint256 reward = (currentBlk - lastClaimedBlk) * rewardRatePerBlk;
    return reward;
  }
}
