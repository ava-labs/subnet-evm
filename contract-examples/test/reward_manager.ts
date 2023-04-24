// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

import { ethers } from "hardhat"
import { test } from "./utils"

// make sure this is always an admin for reward manager precompile
const ADMIN_ADDRESS = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
const REWARD_MANAGER_ADDRESS = "0x0200000000000000000000000000000000000004";

describe("ExampleRewardManager", function () {
  this.timeout("30s")

  beforeEach('Setup DS-Test contract', async function () {
    const signer = await ethers.getSigner(ADMIN_ADDRESS)
    const rewardManagerPromise = ethers.getContractAt("IRewardManager", REWARD_MANAGER_ADDRESS, signer);

    return ethers.getContractFactory("ExampleRewardManagerTest", { signer })
      .then(factory => factory.deploy())
      .then(contract => {
        this.testContract = contract
        return contract.deployed().then(() => contract)
      })
      .then(contract => contract.setUp())
      .then(tx => Promise.all([rewardManagerPromise, tx.wait()]))
      .then(([rewardManager]) => rewardManager.setAdmin(this.testContract.address))
      .then(tx => tx.wait())
  })

  test("should send fees to blackhole", ["sendFeesToBlackhole", "checkSendFeesToBlackhole"])

  // this contract is not selected as the reward address yet, so should not be able to receive fees
  it.skip("should send fees to blackhole", async function () {})

  test("should not appoint reward address before enabled", "doesNotSetRewardAddressBeforeEnabled")

  it.skip("should not appoint reward address before enabled", async function () {});

  test("contract should be added to enabled list", "setEnabled")

  it.skip("contract should be added to enabled list", async function () {});

  test("should be appointed as reward address", "setRewardAddress")

  it.skip("should be appointed as reward address", async function () {});

  // we need to change the fee receiver, send a transaction for the new receiver to receive fees, then check the balance change. 
  // the new fee receiver won't receive fees in the same block where it was set.
  test("should be able to receive fees", ["setupReceiveFees", "receiveFees", "checkReceiveFees"])

  // I don't think it's necessary to test with an EOA since the logic is the same
  it.skip("should be able to receive fees", async function () {})

  test("should return false for allowFeeRecipients check", "areFeeRecipientsAllowed")

  it.skip("should return false for allowFeeRecipients check", async function () {})

  test("should enable allowFeeRecipients", "allowFeeRecipients")

  it.skip("should enable allowFeeRecipients", async function () {})

  test("should disable reward address", "disableRewardAddress")

  it.skip("should disable reward address", async function () {})
});
