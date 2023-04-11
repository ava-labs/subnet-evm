// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

import { expect } from "chai"
import { ethers } from "hardhat"
import { test } from "./utils"

const ADMIN_ADDRESS: string = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
const FEE_MANAGER = "0x0200000000000000000000000000000000000003"

const GENESIS_CONFIG = require('../../tests/precompile/genesis/fee_manager.json')

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

  test("should add contract deployer as owner", "test_addContractDeployerAsOwner")

  test("contract should not be able to change fee without enabled", "test_enableWAGMIFeesFailure")

  test("contract should be added to manager list", "test_addContractToManagerList")

  test("admin should be able to enable change fees", "test_changeFees")

  test("should confirm min-fee transaction", "test_minFeeTransaction", {
    maxFeePerGas: GENESIS_CONFIG.config.feeConfig.minBaseFee,
    maxPriorityFeePerGas: 0,
  })

  test("should reject a transaction below the minimum", [
    "test_raiseMinFeeByOne",
    {
      method: "test_minFeeTransaction",
      shouldFail: true,
    },
    "test_lowerMinFeeByOne",
  ])
})
