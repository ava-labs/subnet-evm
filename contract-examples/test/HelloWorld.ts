// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

import { expect } from "chai";
import { ethers } from "hardhat"
import {
    Contract,
    ContractFactory,
} from "ethers"


const HELLO_WORLD_ADDRESS = "0x0200000000000000000000000000000000000004";
const modifiedGreeting = "What is up!"

describe.only("HelloWorld", function () {
    let helloWorldContract: Contract;
    let helloWorldPrecompile: Contract;
    let account0;
    let account1;

    before(async function () {
        // Set up accounts 
        let accounts = await ethers.getSigners();
        account0 = accounts[0];
        account1 = accounts[1];
        console.log(`Account 0 Address: ${account0.address}`);
        console.log(`Account 1 Address: ${account1.address}`);

        // Deploy Hello World Contract
        const ContractF: ContractFactory = await ethers.getContractFactory("HelloWorld");
        helloWorldContract = await ContractF.deploy();
        await helloWorldContract.deployed();
        const helloWorldContractAddress: string = helloWorldContract.address;
        console.log(`Contract deployed to: ${helloWorldContractAddress}`);

        //Set up precompile
        // helloWorldPrecompile = await ethers.getContractAt("IHelloWorld", HELLO_WORLD_ADDRESS, account0);
    });

    it("should sayHello properly", async function () {
        let result = await helloWorldContract.callStatic.sayHello();
        expect(result).to.equal("Hello World!");
    });

    it("should setGreeting and sayHello", async function () {
        const modifiedGreeting = "What's up";
        let tx = await helloWorldContract.setGreeting(modifiedGreeting);
        await tx.wait();

        expect(await helloWorldContract.callStatic.sayHello()).to.be.equal(modifiedGreeting);
    });
});
