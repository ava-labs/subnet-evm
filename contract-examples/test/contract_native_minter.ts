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
const assert = require("assert")

// make sure this is always an admin for minter precompile
const ADMIN_ADDRESS: string = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
const MINT_PRECOMPILE_ADDRESS = "0x0200000000000000000000000000000000000001";
const mintValue = ethers.utils.parseEther("1")
const initialValue = ethers.utils.parseEther("10")

const ROLES = {
  NONE: 0,
  MINTER: 1,
  ADMIN: 2
};

const testFn = (fnNameOrNames: string | string[], overrides = {}, debug = false) => async function () {
  const fnNames: string[] = Array.isArray(fnNameOrNames) ? fnNameOrNames : [fnNameOrNames];

  return fnNames.reduce((p: Promise<undefined>, fnName) => p.then(async () => {
    const tx = await this.testContract['test_' + fnName](overrides)
    const txReceipt = await tx.wait().catch(err => err.receipt)

    const failed = txReceipt.status !== 0 ? await this.testContract.callStatic.failed() : true
    
    if (debug || failed) {
      console.log('')

      if (!txReceipt.events) console.warn(txReceipt);

      txReceipt
        .events
        ?.filter(event => debug || event.event?.startsWith('log'))
        .map(event => event.args?.forEach(arg => console.log(arg)))

      console.log('')
    }

    assert(!failed, `${fnName} failed`)
  }), Promise.resolve());
}

const test = (name, fnName, overrides = {}) => it(name, testFn(fnName, overrides));
test.only = (name, fnName, overrides = {}) => it.only(name, testFn(fnName, overrides));
test.debug = (name, fnName, overrides = {}) => it.only(name, testFn(fnName, overrides, true));
test.skip = (name, fnName, overrides = {}) => it.skip(name, testFn(fnName, overrides));

describe("ERC20NativeMinter", function () {
  let owner: SignerWithAddress
  let contract: Contract
  let minter: SignerWithAddress
  
  beforeEach('Setup DS-Test contract', async function () {
    const signer = await ethers.getSigner(ADMIN_ADDRESS)
    const nativeMinterPromise = ethers.getContractAt("INativeMinter", MINT_PRECOMPILE_ADDRESS, signer)

    return ethers.getContractFactory("ERC20NativeMinterTest", { signer })
      .then(factory => factory.deploy())
      .then(contract => {
        this.testContract = contract;
        return contract.deployed().then(() => contract)
      })
      .then(contract => contract.setUp())
      .then(tx => Promise.all([nativeMinterPromise, tx.wait()]))
      .then(([nativeMinter]) => nativeMinter.setAdmin(this.testContract.address))
      .then(tx => tx.wait())
  })

  before(async function () {
    owner = await ethers.getSigner(ADMIN_ADDRESS);
    const Token: ContractFactory = await ethers.getContractFactory("ERC20NativeMinter", { signer: owner })
    contract = await Token.deploy(initialValue)
    await contract.deployed()
    const contractAddress: string = contract.address
    console.log(`Contract deployed to: ${contractAddress}`)

    const name: string = await contract.name()
    console.log(`Name: ${name}`)

    const symbol: string = await contract.symbol()
    console.log(`Symbol: ${symbol}`)

    const decimals: string = await contract.decimals()
    console.log(`Decimals: ${decimals}`)

    const signers: SignerWithAddress[] = await ethers.getSigners()
    minter = signers.slice(-1)[0]

    // Fund minter address
    await owner.sendTransaction({
      to: minter.address,
      value: ethers.utils.parseEther("1")
    })
  });

  // this test doesn't test precompile logic  
  it.skip("should add contract deployer as owner", async function () {});

  test("contract should not be able to mintdraw", "mintdrawFailure")

  // this contract is not given minter permission yet, so should not mintdraw
  it.skip("contract should not be able to mintdraw", async function () {})

  test("should be added to minter list", "addMinter")

  it.skip("should be added to minter list", async function () {});

  test("admin should mintdraw", "adminMintdraw")

  // admin should mintdraw since it has ERC20 token initially.
  it.skip("admin should mintdraw", async function () {})

  test("minter should not mintdraw ", "minterMintdrawFailure")

  // minter should not mintdraw since it has no ERC20 token.
  it.skip("minter should not mintdraw ", async function () {})
  
  // minter should not mintdraw since it has no ERC20 token.
  it("minter should not mintdraw ", async function () {
    try {
      await contract.connect(minter).mintdraw(mintValue)
    }
    catch (err) {
      return
    }
    expect.fail("should have errored")
  })

  test("should deposit for minter", "minterDeposit")

  // minter should not mintdraw since it has no ERC20 token.
  it.skip("should deposit for minter", async function () {})

  test("minter should mintdraw", "mintdraw")

  // minter should mintdraw now since it has ERC20 token.
  it.skip("minter should mintdraw", async function () {})
})
