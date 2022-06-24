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

const ROLES = {
  NONE: 0,
  ENABLED: 1,
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
  blockGasCostStep: 100_000 // blockGasCostStep
}

const LOW_FEES = {
  gasLimit: 2_000_0000, // gasLimit
  targetBlockRate: 2, // targetBlockRate
  minBaseFee: 1_000_000_000, // minBaseFee
  targetGas: 100_000_000, // targetGas
  baseFeeChangeDenominator: 48, // baseFeeChangeDenominator
  minBlockGasCost: 0, // minBlockGasCost
  maxBlockGasCost: 10_000_000, // maxBlockGasCost
  blockGasCostStep: 100_000 // blockGasCostStep
}


describe("ExampleFeeManager", function () {
  this.timeout("30s")

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

    const feeManager = await ethers.getContractAt("IFeeManager", FEE_MANAGER, owner);
    let tx = await feeManager.setEnabled(manager.address);
    await tx.wait()

    tx = await owner.sendTransaction({
      to: manager.address,
      value: ethers.utils.parseEther("1")
    })
    await tx.wait()

    tx = await owner.sendTransaction({
      to: nonEnabled.address,
      value: ethers.utils.parseEther("1")
    })
    await tx.wait()
  });


  it("manager should be able to change fees from precompile", async function () {
    const feeManager = await ethers.getContractAt("IFeeManager", FEE_MANAGER, owner);
    let testMinBaseFee = LOW_FEES.minBaseFee + 123

    let enableTx = await feeManager.setFeeConfig(
      LOW_FEES.gasLimit,
      LOW_FEES.targetBlockRate,
      testMinBaseFee,
      LOW_FEES.targetGas,
      LOW_FEES.baseFeeChangeDenominator,
      LOW_FEES.minBlockGasCost,
      LOW_FEES.maxBlockGasCost,
      LOW_FEES.blockGasCostStep
    )
    let txRes = await enableTx.wait();

    var res = await contract.connect(manager).getCurrentFeeConfig()
    expect(res.minBaseFee).to.equal(testMinBaseFee)

    var res = await contract.getLastChangedAt()

    expect(res).to.equal(txRes.blockNumber)
  })

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
    const feeManager = await ethers.getContractAt("IFeeManager", FEE_MANAGER, owner);
    let adminRole = await feeManager.readAllowList(adminAddress);
    expect(adminRole).to.be.equal(ROLES.ADMIN)
    let contractRole = await feeManager.readAllowList(contract.address);
    expect(contractRole).to.be.equal(ROLES.NONE)

    let enableTx = await feeManager.setEnabled(contract.address);
    await enableTx.wait()
    contractRole = await feeManager.readAllowList(contract.address);
    expect(contractRole).to.be.equal(ROLES.ENABLED)
  });

  it("admin should be able to change fees through contract", async function () {
    let enableTx = await contract.enableCustomFees(LOW_FEES)
    let txRes = await enableTx.wait()

    var res = await contract.getCurrentFeeConfig()
    expect(res.gasLimit).to.equal(LOW_FEES.gasLimit)
    expect(res.minBaseFee).to.be.equal(LOW_FEES.minBaseFee)

    var res = await contract.getLastChangedAt()

    expect(res).to.equal(txRes.blockNumber)
  })

  it("should let low fee tx to be in mempool", async function () {
    var res = await contract.getCurrentFeeConfig()
    expect(res.minBaseFee).to.be.equal(LOW_FEES.minBaseFee)

    var testMaxFeePerGas = HIGH_FEES.minBaseFee - 10000

    let tx = await owner.sendTransaction({
      to: manager.address,
      value: ethers.utils.parseEther("0.1"),
      maxFeePerGas: testMaxFeePerGas,
      maxPriorityFeePerGas: 0
    });
    let confirmedTx = await tx.wait()
    expect(confirmedTx.confirmations).to.be.greaterThanOrEqual(1)
  })

  it("should not let low fee tx to be in mempool", async function () {
    var testMaxFeePerGas = HIGH_FEES.minBaseFee - 10000

    let enableTx = await contract.enableCustomFees(HIGH_FEES)
    await enableTx.wait()
    let getRes = await contract.getCurrentFeeConfig()
    expect(getRes.minBaseFee).to.equal(HIGH_FEES.minBaseFee)

    // send tx with lower han HIGH_FEES minBaseFee
    try {
      let tx = await owner.sendTransaction({
        to: manager.address,
        value: ethers.utils.parseEther("0.1"),
        maxFeePerGas: testMaxFeePerGas,
        maxPriorityFeePerGas: 0
      });
      let res = await tx.wait()
    }
    catch (err) {
      expect(err.toString()).to.include("transaction underpriced")
      return
    }
    expect.fail("should have errored")
  })

  it("should be able to get current fee config", async function () {
    let enableTx = await contract.enableCustomFees(HIGH_FEES)
    await enableTx.wait()

    var res = await contract.getCurrentFeeConfig()
    expect(res.gasLimit).to.equal(HIGH_FEES.gasLimit)

    var res = await contract.connect(manager).getCurrentFeeConfig()
    expect(res.gasLimit).to.equal(HIGH_FEES.gasLimit)

    var res = await contract.connect(nonEnabled).getCurrentFeeConfig()
    expect(res.gasLimit).to.equal(HIGH_FEES.gasLimit)
  });

  it("nonEnabled should not be able to set fee config", async function () {
    const feeManager = await ethers.getContractAt("IFeeManager", FEE_MANAGER, owner);

    let nonEnabledRole = await feeManager.readAllowList(nonEnabled.address);

    expect(nonEnabledRole).to.be.equal(ROLES.NONE)
    try {
      await contract.connect(nonEnabled).enableWAGMIFees()
    }
    catch (err) {
      return
    }
    expect.fail("should have errored")
  })

  it("manager should be able to change fees through contract", async function () {
    let enableTx = await contract.connect(manager).enableCustomFees(LOW_FEES)
    await enableTx.wait()

    var res = await contract.connect(manager).getCurrentFeeConfig()
    expect(res.minBaseFee).to.equal(LOW_FEES.minBaseFee)
  })


  it("non-enabled should not be able to change fees through contract", async function () {
    try {
      let enableTx = await contract.connect(nonEnabled).enableCustomFees(LOW_FEES)
      await enableTx.wait()
    }
    catch (err) {
      return
    }
    expect.fail("should have errored")
  })
})
