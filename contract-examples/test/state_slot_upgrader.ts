// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

import { expect } from "chai"
import { ethers } from "hardhat"
import { Contract, ContractFactory } from "ethers"

describe("ExampleStateSlotUpgrade", function () {
    let capIncreaserContract: Contract

    before(async function () {
        // Deploy ExampleStateSlotUpgrade Contract
        const ContractF: ContractFactory = await ethers.getContractFactory(
            "ExampleStateSlotUpgrade"
        )

        capIncreaserContract = await ContractF.deploy()
        await capIncreaserContract.deployTransaction.wait(1)
        console.log(`Contract deployed to: ${capIncreaserContract.address}`)

        console.log("Using hardcoded contract address: 0x52C84043CD9c865236f11d9Fc9F56aa003c1f922")
        capIncreaserContract = ContractF.attach("0x52C84043CD9c865236f11d9Fc9F56aa003c1f922")
    })

    it("cap was updated to 250 million", async function () {
        expect(await capIncreaserContract.cap()).to.be.equal(
            ethers.BigNumber.from(ethers.utils.parseEther("250000000"))
        )
    })
})