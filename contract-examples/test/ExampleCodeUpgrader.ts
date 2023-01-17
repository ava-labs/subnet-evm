// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

import { expect } from "chai"
import { ethers } from "hardhat"
import { Contract, ContractFactory } from "ethers"

describe("ExampleCodeUpgrader", function () {
    let beforeContract: Contract
    let afterContract: Contract
    let ExampleCodeUpgraderBeforeFactory: ContractFactory
    let ExampleCodeUpgraderAfterFactory: ContractFactory

    before(async function () {
        // Deploy the two contracts, before and after.
        ExampleCodeUpgraderBeforeFactory = await ethers.getContractFactory(
            "ExampleCodeUpgraderBefore"
        )
        beforeContract = await ExampleCodeUpgraderBeforeFactory.deploy()
        await beforeContract.deployTransaction.wait(1)
        console.log(`Before contract deployed to: ${beforeContract.address}`)

        ExampleCodeUpgraderAfterFactory = await ethers.getContractFactory(
            "ExampleCodeUpgraderAfter"
        )
        afterContract = await ExampleCodeUpgraderAfterFactory.deploy()
        await afterContract.deployTransaction.wait(1)
        console.log(`After contract deployed to: ${afterContract.address}`)
    })

    it("update cap method fails on before contract, but not after", async function () {
        // expect(await beforeContract.updateCap()).to.be.equal(
        //     ethers.BigNumber.from(ethers.utils.parseEther("250000000"))
        // )

        expect(await afterContract.updateCap("250000000")).to.emit(
            afterContract, "CapUpdated"
        )
    })

    it("after network upgrade, before contract now has after contract's code", async function () {
        console.log("Using hardcoded before contract addresses: 0x52C84043CD9c865236f11d9Fc9F56aa003c1f922")
        const newBeforeContract = ExampleCodeUpgraderAfterFactory.attach("0x52C84043CD9c865236f11d9Fc9F56aa003c1f922")

        expect(await newBeforeContract.updateCap(ethers.BigNumber.from("250000000"))).to.emit(
            newBeforeContract, "CapUpdated"
        )
    })
})