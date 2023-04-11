//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./IRewardManager.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "./AllowListTest.sol";

address constant REWARD_MANAGER_ADDRESS = 0x0200000000000000000000000000000000000004;
address constant BLACKHOLE_ADDRESS = 0x0100000000000000000000000000000000000000;

// ExampleRewardManager is a sample wrapper contract for RewardManager precompile.
contract ExampleRewardManager is Ownable {
  IRewardManager rewardManager = IRewardManager(REWARD_MANAGER_ADDRESS);

  constructor() Ownable() {}

  function currentRewardAddress() public view returns (address) {
    return rewardManager.currentRewardAddress();
  }

  function setRewardAddress(address addr) public onlyOwner {
    rewardManager.setRewardAddress(addr);
  }

  function allowFeeRecipients() public onlyOwner {
    rewardManager.allowFeeRecipients();
  }

  function disableRewards() public onlyOwner {
    rewardManager.disableRewards();
  }

  function areFeeRecipientsAllowed() public view returns (bool) {
    return rewardManager.areFeeRecipientsAllowed();
  }
}

contract ExampleRewardManagerTest is AllowListTest {
  uint blackholeBalance;
  ExampleRewardManager exampleReceiveFees;
  uint exampleBalance;

  function setUp() public {
    blackholeBalance = BLACKHOLE_ADDRESS.balance;
  }

  function test_captureBlackholeBalance() public {
    blackholeBalance = BLACKHOLE_ADDRESS.balance;
  }

  function test_checkSendFeesToBlackhole() public {
    assertGt(BLACKHOLE_ADDRESS.balance, blackholeBalance);
  }

  function test_doesNotSetRewardAddressBeforeEnabled() public {
    ExampleRewardManager example = new ExampleRewardManager();
    address exampleAddress = address(example);
    IRewardManager rewardManager = IRewardManager(REWARD_MANAGER_ADDRESS);

    assertRole(rewardManager.readAllowList(exampleAddress), AllowList.Role.None);

    try example.setRewardAddress(exampleAddress) {
      assertTrue(false, "setRewardAddress should fail");
    } catch {} // TODO should match on an error to make sure that this is failing in the way that's expected
  }

  function test_setEnabled() public {
    ExampleRewardManager example = new ExampleRewardManager();
    address exampleAddress = address(example);
    IRewardManager rewardManager = IRewardManager(REWARD_MANAGER_ADDRESS);

    assertRole(rewardManager.readAllowList(exampleAddress), AllowList.Role.None);
    rewardManager.setEnabled(exampleAddress);
    assertRole(rewardManager.readAllowList(exampleAddress), AllowList.Role.Enabled);
  }

  function test_setRewardAddress() public {
    ExampleRewardManager example = new ExampleRewardManager();
    address exampleAddress = address(example);
    IRewardManager rewardManager = IRewardManager(REWARD_MANAGER_ADDRESS);

    rewardManager.setEnabled(exampleAddress);
    example.setRewardAddress(exampleAddress);

    assertEq(example.currentRewardAddress(), exampleAddress);
  }

  function test_setupReceiveFees() public {
    ExampleRewardManager example = new ExampleRewardManager();
    address exampleAddress = address(example);

    IRewardManager rewardManager = IRewardManager(REWARD_MANAGER_ADDRESS);

    rewardManager.setEnabled(exampleAddress);
    example.setRewardAddress(exampleAddress);

    exampleReceiveFees = example;
    exampleBalance = exampleAddress.balance;
  }

  function test_receiveFees() public {
    // used as a noop to test if the correct address receives fees
  }

  function test_checkReceiveFees() public {
    assertGt(address(exampleReceiveFees).balance, exampleBalance);
  }

  function test_areFeeRecipientsAllowed() public {
    ExampleRewardManager example = new ExampleRewardManager();
    assertTrue(!example.areFeeRecipientsAllowed());
  }

  function test_allowFeeRecipients() public {
    ExampleRewardManager example = new ExampleRewardManager();
    address exampleAddress = address(example);

    IRewardManager rewardManager = IRewardManager(REWARD_MANAGER_ADDRESS);
    rewardManager.setEnabled(exampleAddress);

    example.allowFeeRecipients();
    assertTrue(example.areFeeRecipientsAllowed());
  }

  function test_disableRewardAddress() public {
    ExampleRewardManager example = new ExampleRewardManager();
    address exampleAddress = address(example);

    IRewardManager rewardManager = IRewardManager(REWARD_MANAGER_ADDRESS);
    rewardManager.setEnabled(exampleAddress);

    example.setRewardAddress(exampleAddress);

    assertEq(example.currentRewardAddress(), exampleAddress);

    example.disableRewards();

    assertEq(example.currentRewardAddress(), BLACKHOLE_ADDRESS);
  }
}
