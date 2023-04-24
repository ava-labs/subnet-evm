// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

import { ethers } from "hardhat"
import { test } from "./utils"

// make sure this is always an admin for minter precompile
const ADMIN_ADDRESS: string = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
const DEPLOYER_ALLOWLIST_ADDRESS = "0x0200000000000000000000000000000000000000";

describe("ExampleDeployerList", function () {
  beforeEach('Setup DS-Test contract', async function () {
    const signer = await ethers.getSigner(ADMIN_ADDRESS)
    const allowListPromise = ethers.getContractAt("IAllowList", DEPLOYER_ALLOWLIST_ADDRESS, signer)

    return ethers.getContractFactory("ExampleDeployerListTest", { signer })
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

  test("precompile should see owner address has admin role", "verifySenderIsAdmin")

  it.skip("precompile should see owner address has admin role", async function () {});

  test("precompile should see test address has no role", "newAddressHasNoRole")

  it.skip("precompile should see test address has no role", async function () {});

  test("contract should report test address has no admin role", "noRoleIsNotAdmin")

  it.skip("contract should report test address has no admin role", async function () {});

  test("contract should report owner address has admin role", "ownerIsAdmin")

  it.skip("contract should report owner address has admin role", async function () {});

  test("should not let test address deploy", "noRoleCannotDeploy")

  it.skip("should not let test address to deploy", async function () {});

  test("should allow admin to add contract as admin", "adminAddContractAsAdmin")

  it.skip("should allow admin to add contract as admin", async function () {});

  test("should allow admin to add deployer address as deployer through contract", "addDeployerThroughContract")

  it.skip("should allow admin to add deployer address as deployer through contract", async function () {});

  test("should let deployer address to deploy", "deployerCanDeploy")

  it.skip("should let deployer address to deploy", async function () {});

  test("should let admin revoke deployer", "adminCanRevokeDeployer")
  
  it.skip("should let admin revoke deployer", async function () {});
})
