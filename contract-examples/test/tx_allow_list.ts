// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

import { ethers } from "hardhat"
import { test } from "./utils"

// make sure this is always an admin for minter precompile
const ADMIN_ADDRESS = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
const TX_ALLOW_LIST_ADDRESS = "0x0200000000000000000000000000000000000002"

describe("ExampleTxAllowList", function () {
  beforeEach('Setup DS-Test contract', async function () {
    const signer = await ethers.getSigner(ADMIN_ADDRESS)
    const allowListPromise = ethers.getContractAt("IAllowList", TX_ALLOW_LIST_ADDRESS, signer);

    return ethers.getContractFactory("ExampleTxAllowListTest", { signer })
      .then(factory => factory.deploy())
      .then(contract => {
        this.testContract = contract;
        return contract.deployed().then(() => contract)
      })
      .then(contract => contract.setUp())
      .then(tx => Promise.all([allowListPromise, tx.wait()]))
      .then(([allowList]) => allowList.setAdmin(this.testContract.address))
      .then(tx => tx.wait())
  })

  test("should add contract deployer as admin", "test_contractOwnerIsAdmin")

  test("precompile should see admin address has admin role", "test_precompileHasDeployerAsAdmin")

  test("precompile should see test address has no role", "test_newAddressHasNoRole")

  test("contract should report test address has on admin role", "test_noRoleIsNotAdmin")

  test("contract should report admin address has admin role", "test_exmapleAllowListReturnsTestIsAdmin")

  // does not test the VM, but rather isolates the contract functionality
  test("should not let test address submit txs", "test_cantDeployFromNoRole");

  test("should not allow noRole to enable itself", "test_noRoleCannotEnableItself")

  test("should allow admin to add contract as admin", "test_addContractAsAdmin")

  test("should allow admin to add allowed address as allowed through contract", "test_enableThroughContract")

  // seems like ethers can't estimate the gas properly on this one
  test("should let allowed address deploy", "test_canDeploy", { gasLimit: "20000000" })

  test("should not let allowed add another allowed", "test_onlyAdminCanEnable")

  test("should not let allowed to revoke admin", "test_onlyAdminCanRevoke") 

  test("should let admin to revoke allowed", "test_adminCanRevoke")
})
