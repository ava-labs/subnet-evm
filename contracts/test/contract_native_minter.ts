// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

import { ethers } from "hardhat"
import { test } from "./utils"
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { expect } from "chai";
import { Contract } from "ethers"

const ADMIN_ADDRESS: string = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
const MINT_PRECOMPILE_ADDRESS = "0x0200000000000000000000000000000000000001"

describe("ERC20NativeMinter", function () {
  beforeEach('Setup DS-Test contract', async function () {
    const signer = await ethers.getSigner(ADMIN_ADDRESS)
    const nativeMinterPromise = ethers.getContractAt("INativeMinter", MINT_PRECOMPILE_ADDRESS, signer)

    return ethers.getContractFactory("ERC20NativeMinterTest", { signer })
      .then(factory => factory.deploy())
      .then(contract => {
        this.testContract = contract
        return contract.deployed().then(() => contract)
      })
      .then(contract => contract.setUp())
      .then(tx => Promise.all([nativeMinterPromise, tx.wait()]))
      .then(([nativeMinter]) => nativeMinter.setAdmin(this.testContract.address))
      .then(tx => tx.wait())
  })

  test("contract should not be able to mintdraw", "step_mintdrawFailure")

  test("should be added to minter list", "step_addMinter")

  test("admin should mintdraw", "step_adminMintdraw")

  test("minter should not mintdraw", "step_minterMintdrawFailure")

  test("should deposit for minter", "step_minterDeposit")

  test("minter should mintdraw", "step_mintdraw")
})


describe("INativeMinter", function () {
  let owner: SignerWithAddress
  let contract: Contract
  before(async function () {
    owner = await ethers.getSigner(ADMIN_ADDRESS);
    contract = await ethers.getContractAt("INativeMinter", MINT_PRECOMPILE_ADDRESS, owner)
  });

  it("should emit NativeCoinMinted event", async function () {
    let testAddress = "0x0444400000000000000000000000000000000004"
    let amount = 1000
    await expect(contract.mintNativeCoin(testAddress, amount))
      .to.emit(contract, 'NativeCoinMinted')
      .withArgs(owner.address, testAddress, amount)
  })
})
