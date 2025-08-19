// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

import { expect } from "chai"
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { Contract } from "ethers"
import { ethers } from "hardhat"
import { test } from "./utils"

// make sure this is always an admin for hello world precompile
const ADMIN_ADDRESS = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
const HELLO_WORLD_ADDRESS = "0x0300000000000000000000000000000000000000"

describe("ExampleHelloWorldTest", function () {
  this.timeout("30s")

  beforeEach('Setup DS-Test contract', async function () {
    const signer = await ethers.getSigner(ADMIN_ADDRESS)
    const helloWorldPromise = ethers.getContractAt("IHelloWorld", HELLO_WORLD_ADDRESS, signer)

    return ethers.getContractFactory("ExampleHelloWorldTest", { signer })
      .then(factory => factory.deploy())
      .then(contract => {
        this.testContract = contract
        return contract.deployed().then(() => contract)
      })
      .then(() => Promise.all([helloWorldPromise]))
      .then(([helloWorld]) => helloWorld.setAdmin(this.testContract.address))
      .then(tx => tx.wait())
  })

  test("should gets default hello world", ["step_getDefaultHelloWorld"])

  test("should not set greeting before enabled", "step_doesNotSetGreetingBeforeEnabled")

  test("should set and get greeting with enabled account", "step_setAndGetGreeting")
});

describe("IHelloWorld events", function () {
  let owner: SignerWithAddress
  let contract: Contract
  let defaultGreeting = "Hello, World!"
  before(async function () {
    owner = await ethers.getSigner(ADMIN_ADDRESS);
    contract = await ethers.getContractAt("IHelloWorld", HELLO_WORLD_ADDRESS, owner)

    // reset greeting
    let tx = await contract.setGreeting(defaultGreeting)
    await tx.wait()
  });

  it("should emit GreetingChanged event", async function () {
    let newGreeting = "helloprecompile"
    await expect(contract.setGreeting(newGreeting)
    )
    .to.emit(contract, "GreetingChanged").withArgs(owner.address,
      // old greeting
      defaultGreeting,
      // new greeting
      newGreeting
    )
  })
})
