//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./IAllowList.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract ExampleFeeDistributor is Ownable {
  mapping(address => bool) rewardAddresses;
  uint256 rewardAddressCount;

  constructor() Ownable() {}

  function addRewardAddress(address addr) public onlyOwner {
    require(!rewardAddresses[addr], "Already a reward address");

    rewardAddresses[addr] = true;
    rewardAddressCount++;
  }

  function revoke(address addr) public onlyOwner {
    require(rewardAddresses[addr], "Not a reward address");

    rewardAddresses[addr] = false;
    rewardAddressCount--;
  }

  function claim() public {
    require(rewardAddresses[msg.sender], "Not a reward address");
  }
}
