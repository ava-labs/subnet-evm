// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

import { expect } from "chai";
import { ethers } from "hardhat"
const assert = require("assert");

// make sure this is always an admin for the precompile
const ADMIN_ADDRESS: string = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
const FEE_MANAGER = "0x0200000000000000000000000000000000000003";

const GENESIS_CONFIG = require('../../tests/precompile/genesis/fee_manager.json');

const testFn = (fnName, overrides = {}) => async function () {
  const tx = await this.testContract["test_" + fnName](overrides)
  const txReceipt = await tx.wait()
  const failed = await this.testContract.callStatic.failed()

  if (failed) {
    console.log('')

    txReceipt
      .events
      ?.filter(event => event.event?.startsWith('log'))
      .map(event => event.args?.forEach(arg => console.log(arg)))

    console.log('')
  }

  assert(!failed, `${fnName} failed`)
}

const test = (name, fnName, overrides = {}) => it(name, testFn(fnName, overrides)); 
test.only = (name, fnName, overrides = {}) => it.only(name, testFn(fnName, overrides));

// TODO: These tests keep state to the next state. It means that some tests cases assumes some preconditions
// set by previous test cases. We should make these tests stateless.
describe("ExampleFeeManager", function () {
  this.timeout("30s")

  beforeEach("setup DS-Test contract", async function () {
    const signer = await ethers.getSigner(ADMIN_ADDRESS)
    const feeManagerPromise = ethers.getContractAt("IFeeManager", FEE_MANAGER, signer)

    return ethers.getContractFactory("ExampleFeeManagerTest", { signer })
      .then(factory => factory.deploy())
      .then(contract => {
        this.testContract = contract
        return contract.deployed().then(() => contract)
      })
      .then(contract => contract.setUp())
      .then(tx => Promise.all([feeManagerPromise, tx.wait()]))
      .then(([feeManager]) => feeManager.setAdmin(this.testContract.address))
      .then(tx => tx.wait())
  })

  test("should add contract deployer as owner", "addContractDeployerAsOwner");

  test("contract should not be able to change fee without enabled", "enableWAGMIFeesFailure"); 

  test("contract should be added to manager list", "addContractToManagerList");

  test("admin should be able to enable change fees", "changeFees");

  test("should confirm min-fee transaction", "minFeeTransaction", {
    maxFeePerGas: GENESIS_CONFIG.config.feeConfig.minBaseFee,
    maxPriorityFeePerGas: 0,
  })

  // TODO: I should be able to test inside the contract by manipulating gas mid-call
  it("should reject a transaction below the minimum", async function() {
    const maxFeePerGas = GENESIS_CONFIG.config.feeConfig.minBaseFee;
    const maxPriorityFeePerGas = 0;
    const gasLimit = GENESIS_CONFIG.config.feeConfig.gasLimit;

    await this.testContract.raiseMinFeeByOne({ maxFeePerGas }).then(tx => tx.wait());

    try {
      await this.testContract.lowerMinFeeByOne({ maxFeePerGas }).then(tx => tx.wait());
    } catch (err) {
      expect(err.toString()).to.include("max fee per gas less than block base fee");

      await this.testContract.lowerMinFeeByOne({ maxFeePerGas: maxFeePerGas + 1 }).then(tx => tx.wait());

      return;
    }

    expect.fail("should have errored")
  })
})
