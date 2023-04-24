// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

import { ethers } from "hardhat"
import { test } from "./utils"

// make sure this is always an admin for minter precompile
const ADMIN_ADDRESS: string = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
const MINT_PRECOMPILE_ADDRESS = "0x0200000000000000000000000000000000000001";

describe("ERC20NativeMinter", function () {
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
  
  test("should deposit for minter", "minterDeposit")

  // minter should not mintdraw since it has no ERC20 token.
  it.skip("should deposit for minter", async function () {})

  test("minter should mintdraw", "mintdraw")

  // minter should mintdraw now since it has ERC20 token.
  it.skip("minter should mintdraw", async function () {})
})
