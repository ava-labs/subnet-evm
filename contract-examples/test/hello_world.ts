// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

import { expect } from "chai"
import { ethers } from "hardhat"
import { Contract, ContractFactory } from "ethers"

describe("ExampleHelloWorld", function () {
  let helloWorldContract: Contract

  before(async function () {
    // Deploy Hello World Contract
    const ContractF: ContractFactory = await ethers.getContractFactory(
      "ExampleHelloWorld"
    )
    helloWorldContract = await ContractF.deploy()
    await helloWorldContract.deployed()
    const helloWorldContractAddress: string = helloWorldContract.address
    console.log(`Contract deployed to: ${helloWorldContractAddress}`)
  })

  it("should getHello properly", async function () {
    let result = await helloWorldContract.callStatic.getHello()
    expect(result).to.equal("Hello World!")
  })

  it("should setGreeting and getHello", async function () {
    const modifiedGreeting = "What's up"
    let tx = await helloWorldContract.setGreeting(modifiedGreeting)
    await tx.wait()

    expect(await helloWorldContract.callStatic.getHello()).to.be.equal(
      modifiedGreeting
    )
  })
})
