// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import { expect } from "chai"
import {
  Contract,
  ContractFactory,
} from "ethers"
import { ethers } from "hardhat"
const assert = require("assert")

// make sure this is always an admin for minter precompile
const ADMIN_ADDRESS = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
const TX_ALLOW_LIST_ADDRESS = "0x0200000000000000000000000000000000000002"

const ROLES = {
  NONE: 0,
  ALLOWED: 1,
  ADMIN: 2
}

const testFn = (fnName, overrides = {}, debug = false) => async function () {
  const tx = await this.testContract['test_' + fnName](overrides)
  const txReceipt = await tx.wait().catch(err => err.receipt)

  const failed = txReceipt.status !== 0 ? await this.testContract.callStatic.failed() : true
  
  if (debug || failed) {
    console.log('')

    txReceipt
      .events
      ?.filter(event => debug || event.event?.startsWith('log'))
      .map(event => event.args?.forEach(arg => console.log(arg)))

    console.log('')
  }

  assert(!failed, `${fnName} failed`)
}

const test = (name, fnName, overrides = {}) => it(name, testFn(fnName, overrides));
test.only = (name, fnName, overrides = {}) => it.only(name, testFn(fnName, overrides));
test.debug = (name, fnName, overrides = {}) => it.only(name, testFn(fnName, overrides, true));

describe("ExampleTxAllowList", function () {
  let admin: SignerWithAddress
  let contract: Contract
  let allowed: SignerWithAddress
  let noRole: SignerWithAddress

  beforeEach('Setup DS-Test contract', async function () {
    const signer = await ethers.getSigner(ADMIN_ADDRESS)
    const allowListPromise = ethers.getContractAt("IAllowList", TX_ALLOW_LIST_ADDRESS, signer);

    return ethers.getSigner(ADMIN_ADDRESS)
      .then(signer => ethers.getContractFactory("ExampleTxAllowListTest", signer))
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

  before(async function () {
    admin = await ethers.getSigner(ADMIN_ADDRESS)
    const contractF: ContractFactory = await ethers.getContractFactory("ExampleTxAllowList", { signer: admin })
    contract = await contractF.deploy()
    await contract.deployed()
    const contractAddress: string = contract.address
    console.log(`Contract deployed to: ${contractAddress}`)

      ;[, allowed, noRole] = await ethers.getSigners()

    // Fund allowed address
    await admin.sendTransaction({
      to: allowed.address,
      value: ethers.utils.parseEther("10")
    })

    // Fund no role address
    let tx = await admin.sendTransaction({
      to: noRole.address,
      value: ethers.utils.parseEther("10")
    })
    await tx.wait()
  })

  test("should add contract deployer as admin", "contractOwnerIsAdmin")

  // this is testing something that is setup in the before hook 
  // I don't really understand why... 
  it("should add contract deployer as admin", async function () {
    const contractOwnerAdmin: string = await contract.isAdmin(contract.owner())
    expect(contractOwnerAdmin).to.be.true
  })

  test("precompile should see admin address has admin role", "precompileHasDeployerAsAdmin")

  // this isn't really testing the functionality of the precompile but instead testing if the genesis config 
  // set things up as expected.
  it("precompile should see admin address has admin role", async function () {
    // test precompile first
    const allowList = await ethers.getContractAt("IAllowList", TX_ALLOW_LIST_ADDRESS, admin)
    let adminRole = await allowList.readAllowList(admin.address)
    expect(adminRole).to.be.equal(ROLES.ADMIN)
  })

  test("precompile should see test address has no role", "newAddressHasNoRole")

  it("precompile should see test address has no role", async function () {
    // test precompile first
    const allowList = await ethers.getContractAt("IAllowList", TX_ALLOW_LIST_ADDRESS, admin)
    let role = await allowList.readAllowList(noRole.address)
    expect(role).to.be.equal(ROLES.NONE)
  })

  test("contract should report test address has on admin role", "noRoleIsNotAdmin")

  it("contract should report test address has no admin role", async function () {
    const result = await contract.isAdmin(noRole.address)
    expect(result).to.be.false
  })

  test("contract should report admin address has admin role", "exmapleAllowListReturnsTestIsAdmin")

  it("contract should report admin address has admin role", async function () {
    const result = await contract.isAdmin(admin.address)
    expect(result).to.be.true
  })

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
  it("should not let test address submit txs", async function () {
    const Token: ContractFactory = await ethers.getContractFactory("ERC20NativeMinter", { signer: noRole })
    let token: Contract
    try {
      token = await Token.deploy(11111)
    }
    catch (err) {
      expect(err.message).contains("cannot issue transaction from non-allow listed address")
      return
    }
    expect.fail("should have errored")
  })

  test("should not allow noRole to enable itself", "noRoleCannotEnableItself")

  // addDeployer is not a function... this test passes for the wrong reasons
  it("should not allow noRole to enable itself", async function () {
    try {
      await contract.connect(noRole).addDeployer(noRole.address)
    }
    catch (err) {
      return
    }
    expect.fail("should have errored")
  })

  // this test is redundant with the one above
  it("should not allow admin to enable noRole without enabling contract", async function () {
    const allowList = await ethers.getContractAt("IAllowList", TX_ALLOW_LIST_ADDRESS, admin)
    let role = await allowList.readAllowList(contract.address)
    expect(role).to.be.equal(ROLES.NONE)
    const result = await contract.isEnabled(contract.address)
    expect(result).to.be.false
    try {
      await contract.setEnabled(noRole.address)
    }
    catch (err) {
      return
    }
    expect.fail("should have errored")
  })

  test("should allow admin to add contract as admin", "addContractAsAdmin")

  it.only("should allow admin to add contract as admin", async function () {
    const allowList = await ethers.getContractAt("IAllowList", TX_ALLOW_LIST_ADDRESS, admin)
    let role = await allowList.readAllowList(contract.address)
    expect(role).to.be.equal(ROLES.NONE)
    let tx = await allowList.setAdmin(contract.address)
    await tx.wait()
    role = await allowList.readAllowList(contract.address)
    expect(role).to.be.equal(ROLES.ADMIN)
    const result = await contract.isAdmin(contract.address)
    expect(result).to.be.true
  })

  test("should allow admin to add allowed address as allowed through contract", "enableThroughContract")

  it("should allow admin to add allowed address as allowed through contract", async function () {
    let result = await contract.isEnabled(allowed.address)
    expect(result).to.be.false
    let tx = await contract.setEnabled(allowed.address)
    await tx.wait()
    result = await contract.isEnabled(allowed.address)
    expect(result).to.be.true
  })

  test("should let allowed address deploy", "canDeploy")

  it("should let allowed address deploy", async function () {
    const Token: ContractFactory = await ethers.getContractFactory("ERC20NativeMinter", { signer: allowed })
    let token: Contract
    token = await Token.deploy(11111)
    await token.deployed()
    expect(token.address).not.null
  })

  test("should not let allowed add another allowed", "onlyAdminCanEnable")

  it("should not let allowed add another allowed", async function () {
    try {
      const signers: SignerWithAddress[] = await ethers.getSigners()
      const testAddress = signers.slice(-2)[0]
      await contract.connect(allowed).setEnabled(noRole)
    }
    catch (err) {
      return
    }
    expect.fail("should have errored")
  })

  test("should not let allowed to revoke admin", "onlyAdminCanRevoke") 

  it("should not let allowed to revoke admin", async function () {
    try {
      await contract.connect(allowed).revoke(admin.address)
    }
    catch (err) {
      return
    }
    expect.fail("should have errored")
  })

  // this test is redundant, allowed can't revoke anyone including itself 
  it("should not let allowed to revoke itself", async function () {
    try {
      await contract.connect(allowed).revoke(allowed.address)
    }
    catch (err) {
      return
    }
    expect.fail("should have errored")
  })

  test("should let admin to revoke allowed", "adminCanRevoke")

  it("should let admin to revoke allowed", async function () {
    let tx = await contract.revoke(allowed.address)
    await tx.wait()
    const allowList = await ethers.getContractAt("IAllowList", TX_ALLOW_LIST_ADDRESS, admin)
    let noRole = await allowList.readAllowList(allowed.address)
    expect(noRole).to.be.equal(ROLES.NONE)
  })

  // this is only testing a require statement in the example contract
  it("should not let admin to revoke itself", async function () {
    try {
      await contract.revoke(admin.address)
    }
    catch (err) {
      console.log(err)
      return
    }
    expect.fail("should have errored")
  })
})
