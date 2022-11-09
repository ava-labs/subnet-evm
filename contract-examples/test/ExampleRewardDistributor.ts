// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { expect } from "chai";
import {
  BigNumber,
  Contract,
  ContractFactory,
} from "ethers"
import { ethers } from "hardhat"
import ts = require("typescript");

// make sure this is always an admin for minter precompile
const adminAddress: string = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
const REWARD_MANAGER_ADDRESS = "0x0200000000000000000000000000000000000004";
const BLACKHOLE_ADDRESS = "0x0100000000000000000000000000000000000000";
const rewardRate = ethers.utils.parseEther("0.00001")

const ROLES = {
  NONE: 0,
  MINTER: 1,
  ADMIN: 2
};

describe("ExampleRewardDistributor", function () {
  this.timeout("30s")
  let owner: SignerWithAddress
  let contract: Contract
  let signer1: SignerWithAddress
  let signer2: SignerWithAddress
  before(async function () {
    owner = await ethers.getSigner(adminAddress);
    signer1 = (await ethers.getSigners())[1]
    signer2 = (await ethers.getSigners())[2]
    const Contract: ContractFactory = await ethers.getContractFactory("ExampleRewardDistributor", { signer: owner })
    contract = await Contract.deploy(rewardRate)
    await contract.deployed()
    const contractAddress: string = contract.address
    console.log(`Contract deployed to: ${contractAddress}`)

    // Send a transaction to mine a new block
    const tx = await owner.sendTransaction({
      to: signer1.address,
      value: ethers.utils.parseEther("10")
    })
    await tx.wait()
  });

  it("should add contract deployer as owner", async function () {
    const contractOwnerAddr: string = await contract.owner()
    expect(owner.address).to.equal(contractOwnerAddr)
  });

  // this contract is not selected as the reward address yet, so should not be able to receive fees
  it("contract should not be able to receive fees", async function () {
    const rewardManager = await ethers.getContractAt("IRewardManager", REWARD_MANAGER_ADDRESS, owner);
    let rewardAddress = await rewardManager.currentRewardAddress();
    expect(rewardAddress).to.be.equal(BLACKHOLE_ADDRESS)

    let balance = await ethers.provider.getBalance(contract.address)
    expect(balance).to.be.equal(0)

    let firstBHBalance = await ethers.provider.getBalance(BLACKHOLE_ADDRESS)

    // Send a transaction to mine a new block
    const tx = await owner.sendTransaction({
      to: signer1.address,
      value: ethers.utils.parseEther("0.0001")
    })
    await tx.wait()

    balance = await ethers.provider.getBalance(contract.address)
    expect(balance).to.be.equal(0)

    let secondBHBalance = await ethers.provider.getBalance(BLACKHOLE_ADDRESS)
    expect(secondBHBalance).to.be.greaterThan(firstBHBalance)
  })

  it("should be appointed as reward address", async function () {
    const rewardManager = await ethers.getContractAt("IRewardManager", REWARD_MANAGER_ADDRESS, owner);
    let adminRole = await rewardManager.readAllowList(adminAddress);
    expect(adminRole).to.be.equal(ROLES.ADMIN)

    let tx = await rewardManager.setRewardAddress(contract.address);
    await tx.wait()
    let rewardAddress = await rewardManager.currentRewardAddress();
    expect(rewardAddress).to.be.equal(contract.address)
  });

  it("contract should be able to receive fees", async function () {
    let previousBalance = await ethers.provider.getBalance(contract.address)

    // Send a transaction to mine a new block
    const tx = await owner.sendTransaction({
      to: signer1.address,
      value: ethers.utils.parseEther("0.0001")
    })
    await tx.wait()

    let balance = await ethers.provider.getBalance(contract.address)
    expect(balance.gt(previousBalance)).to.be.true
  })

  it("should not let non-added address to claim rewards", async function () {
    const nonRewardAddress = signer1
    try {
      await contract.connect(nonRewardAddress).claim();
    }
    catch (err) {
      expect(err.message).to.contains('Not a reward address')
      return
    }
    expect.fail("should have errored")
  })

  it("should not revoke non-added address", async function () {
    const nonRewardAddress = signer1
    try {
      await contract.revoke(nonRewardAddress.address);
    }
    catch (err) {
      expect(err.message).to.contains('Not a reward address')
      return
    }
    expect.fail("should have errored")
  })

  it("should return 0 estimate reward for non-added address", async function () {
    const nonRewardAddress = signer1
    const estimatedReward = await contract.estimateReward(nonRewardAddress.address)
    expect(estimatedReward).to.be.equal(0)
  })

  it("should return false for isRewardAddress for non-added address", async function () {
    const nonRewardAddress = signer1
    const isRewardAddress = await contract.isRewardAddress(nonRewardAddress.address)
    expect(isRewardAddress).to.be.equal(false)
  })

  it("should be able to add a reward address", async function () {
    const rewardAddress = signer1
    let tx = await contract.addRewardAddress(rewardAddress.address);
    await tx.wait()
    const isRewardAddress = await contract.isRewardAddress(rewardAddress.address)
    expect(isRewardAddress).to.be.equal(true)
  })

  it("should not be able to add a reward address twice", async function () {
    const rewardAddress = signer1
    try {
      await contract.addRewardAddress(rewardAddress.address);
    }
    catch (err) {
      expect(err.message).to.contains('Already a reward address')
      return
    }
    expect.fail("should have errored")
  })

  it("should distribute according to reward rate", async function () {
    const rewardAddress = signer1
    let previousBalance = await ethers.provider.getBalance(rewardAddress.address)
    let previousBlockNum = await ethers.provider.getBlockNumber()

    // send a transaction to mine a new block
    let transfer = await owner.sendTransaction({
      to: signer2.address,
      value: ethers.utils.parseEther("0.0001")
    })

    await transfer.wait()

    let tx = await contract.connect(rewardAddress).claim();
    await tx.wait()

    let balance = await ethers.provider.getBalance(rewardAddress.address)

    let txRec = await tx.wait()
    let gasUsed: BigNumber = txRec.cumulativeGasUsed
    let gasPrice: BigNumber = txRec.effectiveGasPrice
    let txFee = gasUsed.mul(gasPrice)
    let blockNum = await ethers.provider.getBlockNumber()

    let blockDiff = blockNum - previousBlockNum
    expect(balance).to.be.equal(previousBalance.add(rewardRate.mul(blockDiff)).sub(txFee))
  })

  it("should be able to revoke a reward address", async function () {
    const rewardAddress = signer1
    let tx = await contract.revoke(rewardAddress.address);
    await tx.wait()
    const isRewardAddress = await contract.isRewardAddress(rewardAddress.address)
    expect(isRewardAddress).to.be.equal(false)
  })

  it("should be able to claim after revoke", async function () {
    const nonRewardAddress = signer1
    const isRewardAddress = await contract.isRewardAddress(nonRewardAddress.address)
    expect(isRewardAddress).to.be.equal(false)

    try {
      await contract.connect(nonRewardAddress).claim();
    }
    catch (err) {
      expect(err.message).to.contains('Not a reward address')
      return
    }
    expect.fail("should have errored")
  })
});
