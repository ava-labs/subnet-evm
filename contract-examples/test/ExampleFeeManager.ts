import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { expect } from "chai";
import {
  BigNumber,
  Contract,
  ContractFactory,
} from "ethers"
import { ethers } from "hardhat"

// make sure this is always an admin for minter precompile
const adminAddress: string = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
const FEE_MANAGER = "0x0200000000000000000000000000000000000003";
const initialValue = ethers.utils.parseEther("10")

const ROLES = {
  NONE: 0,
  MANAGER: 1,
  ADMIN: 2
};

const HIGH_FEES = {
  gasLimit: 8_000_000, // gasLimit
  targetBlockRate: 2, // targetBlockRate
  minBaseFee: 25_000_000_000, // minBaseFee
  targetGas: 15_000_000, // targetGas
  baseFeeChangeDenominator: 36, // baseFeeChangeDenominator
  minBlockGasCost: 0, // minBlockGasCost
  maxBlockGasCost: 1_000_000, // maxBlockGasCost
  blockGasCostStep: 0 // blockGasCostStep
}

const LOW_FEES = {
  gasLimit: 2_000_0000, // gasLimit
  targetBlockRate: 2, // targetBlockRate
  minBaseFee: 1_000_000_000, // minBaseFee
  targetGas: 100_000_000, // targetGas
  baseFeeChangeDenominator: 48, // baseFeeChangeDenominator
  minBlockGasCost: 0, // minBlockGasCost
  maxBlockGasCost: 10_000_000, // maxBlockGasCost
  blockGasCostStep: 0 // blockGasCostStep
}


describe("ExampleFeeManager", function () {
  let owner: SignerWithAddress
  let contract: Contract
  let manager: SignerWithAddress
  let nonEnabled: SignerWithAddress
  before(async function () {
    owner = await ethers.getSigner(adminAddress);
    const FeeManager: ContractFactory = await ethers.getContractFactory("ExampleFeeManager", { signer: owner })
    contract = await FeeManager.deploy()
    await contract.deployed()
    const contractAddress: string = contract.address
    console.log(`Contract deployed to: ${contractAddress}`)

    const signers: SignerWithAddress[] = await ethers.getSigners()
    manager = signers.slice(-1)[0]
    nonEnabled = signers.slice(-2)[0]
  });

  it("should add contract deployer as owner", async function () {
    const contractOwnerAddr: string = await contract.owner()
    expect(owner.address).to.equal(contractOwnerAddr)
  });

  // this contract is not given minter permission yet, so should not mintdraw
  it("contract should not be able to change fee without enabled", async function () {
    const managerPrecompile = await ethers.getContractAt("IFeeManager", FEE_MANAGER, owner);
    let contractRole = await managerPrecompile.readAllowList(contract.address);
    expect(contractRole).to.be.equal(ROLES.NONE)
    try {
      await contract.enableWAGMIFees()
    }
    catch (err) {
      return
    }
    expect.fail("should have errored")
  })

  it("contract should be added to manager list", async function () {
    const managerList = await ethers.getContractAt("IFeeManager", FEE_MANAGER, owner);
    let adminRole = await managerList.readAllowList(adminAddress);
    expect(adminRole).to.be.equal(ROLES.ADMIN)
    let contractRole = await managerList.readAllowList(contract.address);
    expect(contractRole).to.be.equal(ROLES.NONE)

    let enableTx = await managerList.setEnabled(contract.address);
    await enableTx.wait()
    contractRole = await managerList.readAllowList(contract.address);
    expect(contractRole).to.be.equal(ROLES.MANAGER)
  });

  it("admin should be able to change fees through contract", async function () {
    var res = await contract.getCurrentFeeConfig()
    expect(res.gasLimit).to.equal(0)

    let enableTx = await contract.enableCustomFees(LOW_FEES)
    await enableTx.wait()

    var res = await contract.getCurrentFeeConfig()
    expect(res.gasLimit).to.equal(LOW_FEES.gasLimit)
    expect(res.minBaseFee).to.be.equal(LOW_FEES.minBaseFee)
  })

  it("should not let low fee tx to be in mempool", async function () {
    var res = await contract.getCurrentFeeConfig()
    expect(res.minBaseFee).to.be.equal(LOW_FEES.minBaseFee)
    // Send 1 ether to an ens name.
    let tx = await owner.sendTransaction({
      to: manager.address,
      value: ethers.utils.parseEther("0.1"),
      maxFeePerGas: LOW_FEES.minBaseFee * 2 // safe value
    });
    let confirmedTx = await tx.wait()
    expect(confirmedTx.confirmations).to.be.greaterThanOrEqual(1)

    let enableTx = await contract.enableCustomFees(HIGH_FEES)
    await enableTx.wait()
    let getRes = await contract.getCurrentFeeConfig()
    expect(getRes.minBaseFee).to.equal(HIGH_FEES.minBaseFee)

    // send tx with WAGMI_FEES minBaseFee
    try {
      let tx = await owner.sendTransaction({
        to: manager.address,
        value: ethers.utils.parseEther("0.1"),
        maxFeePerGas: LOW_FEES.minBaseFee * 2, // safe value
      });

      await tx.wait()
    }
    catch (err) {
      expect(err).to.include("max fee per gas less than block base fee")
      return
    }
    expect.fail("should have errored")
  })

  it("admin should be able to reset fees", async function () {
    let enableTx = await contract.resetFeeConfig()
    await enableTx.wait()

    var res = await contract.getCurrentFeeConfig()
    expect(res.gasLimit).to.equal(0)
  })

  it("should be able to get current fee config", async function () {
    var res = await contract.getCurrentFeeConfig()
    expect(res.gasLimit).to.equal(0)

    var res = await contract.connect(manager).getCurrentFeeConfig()
    expect(res.gasLimit).to.equal(0)

    var res = await contract.connect(nonEnabled).getCurrentFeeConfig()
    expect(res.gasLimit).to.equal(0)
  });

  it("manager should not be able to set fee config before enabled", async function () {
    const managerPrecompile = await ethers.getContractAt("IFeeManager", FEE_MANAGER, owner);
    let contractRole = await managerPrecompile.readAllowList(manager.address);
    expect(contractRole).to.be.equal(ROLES.NONE)
    try {
      await contract.connect(manager).enableWAGMIFees()
    }
    catch (err) {
      return
    }
    expect.fail("should have errored")
  })

  it("manager should be enabled", async function () {
    let enableTx = await contract.setEnabled(manager.address);
    await enableTx.wait()
    let contractRole = await contract.isEnabled(manager.address);
    expect(contractRole).to.be.true
  });

  it("manager should be able to change fees through contract", async function () {
    var res = await contract.connect(manager).getCurrentFeeConfig()
    expect(res.gasLimit).to.equal(0)

    let enableTx = await contract.connect(manager).enableWAGMIFees()
    await enableTx.wait()

    var res = await contract.connect(manager).getCurrentFeeConfig()
    expect(res.gasLimit).to.equal(2_000_0000)
  })

  it("manager should be able to reset fees", async function () {
    let enableTx = await contract.connect(manager).resetFeeConfig()
    await enableTx.wait()

    var res = await contract.connect(manager).getCurrentFeeConfig()
    expect(res.gasLimit).to.equal(0)
  })

  it("non-enabled should not be able to change fees through contract", async function () {
    var res = await contract.connect(nonEnabled).getCurrentFeeConfig()
    expect(res.gasLimit).to.equal(0)

    try {
      let enableTx = await contract.connect(nonEnabled).enableWAGMIFees()
      await enableTx.wait()
    }
    catch (err) {
      return
    }
    expect.fail("should have errored")
  })
})
