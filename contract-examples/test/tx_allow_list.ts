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

  test("should add contract deployer as admin", "contractOwnerIsAdmin")

  // this is testing something that is setup in the before hook 
  // I don't really understand why... 
  it.skip("should add contract deployer as admin", async function () {})

  test("precompile should see admin address has admin role", "precompileHasDeployerAsAdmin")

  // this isn't really testing the functionality of the precompile but instead testing if the genesis config 
  // set things up as expected.
  it.skip("precompile should see admin address has admin role", async function () {})

  test("precompile should see test address has no role", "newAddressHasNoRole")

  it.skip("precompile should see test address has no role", async function () {})

  test("contract should report test address has on admin role", "noRoleIsNotAdmin")

  it.skip("contract should report test address has no admin role", async function () {})

  test("contract should report admin address has admin role", "exmapleAllowListReturnsTestIsAdmin")

  it.skip("contract should report admin address has admin role", async function () {})

  // does not test the VM, but rather isolates the contract functionality
  test("should not let test address submit txs", "cantDeployFromNoRole");

  // currently, this not actually testing the precompile logic, but instead the tx-pool/state-transition logic
  // (I'm guessing the former)
  //
  // if st.evm.ChainConfig().IsPrecompileEnabled(txallowlist.ContractAddress, st.evm.Context.Time) {
  //     txAllowListRole := txallowlist.GetTxAllowListStatus(st.state, st.msg.From())
  //     if !txAllowListRole.IsEnabled() {
  //         return fmt.Errorf("%w: %s", vmerrs.ErrSenderAddressNotAllowListed, st.msg.From())
  //     }
  // }
  it.skip("should not let test address submit txs", async function () {})

  test("should not allow noRole to enable itself", "noRoleCannotEnableItself")

  // addDeployer is not a function... this test passes for the wrong reasons
  it.skip("should not allow noRole to enable itself", async function () {})

  test("should allow admin to add contract as admin", "addContractAsAdmin")

  it.skip("should allow admin to add contract as admin", async function () {})

  test("should allow admin to add allowed address as allowed through contract", "enableThroughContract")

  it.skip("should allow admin to add allowed address as allowed through contract", async function () {})

  // seems like ethers can't estimate the gas properly on this one
  test("should let allowed address deploy", "canDeploy", { gasLimit: "20000000" })

  it.skip("should let allowed address deploy", async function () {})

  test("should not let allowed add another allowed", "onlyAdminCanEnable")

  it.skip("should not let allowed add another allowed", async function () {})

  test("should not let allowed to revoke admin", "onlyAdminCanRevoke") 

  it.skip("should not let allowed to revoke admin", async function () {})

  test("should let admin to revoke allowed", "adminCanRevoke")

  it.skip("should let admin to revoke allowed", async function () {})
})
