// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

import { expect } from "chai";
import { ethers } from "hardhat";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { Contract, ContractFactory } from "ethers";

// make sure this is always an admin for the precompile
const adminAddress: string = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
const HELLO_WORLD_ADDRESS = "0x0300000000000000000000000000000000000000";

describe("ExampleHelloWorld", function () {
  let helloWorldContract: Contract;
  let adminSigner: SignerWithAddress;
  let adminSignerPrecompile: Contract;

  before(async function () {
    // Deploy Hello World Contract
    const ContractF: ContractFactory = await ethers.getContractFactory(
      "ExampleHelloWorld"
    );
    helloWorldContract = await ContractF.deploy();
    await helloWorldContract.deployed();
    const helloWorldContractAddress: string = helloWorldContract.address;
    console.log(`Contract deployed to: ${helloWorldContractAddress}`);

    adminSigner = await ethers.getSigner(adminAddress);
    adminSignerPrecompile = await ethers.getContractAt("IHelloWorld", HELLO_WORLD_ADDRESS, adminSigner);
  });

  it("should getHello properly", async function () {
    let result = await helloWorldContract.callStatic.getHello();
    expect(result).to.equal("Hello World!");
  });

  it("contract should not be able to set greeting without enabled", async function () {
    const modifiedGreeting = "What's up";
    let contractRole = await adminSignerPrecompile.readAllowList(helloWorldContract.address);
    expect(contractRole).to.be.equal(0); // 0 = NONE
    try {
      let tx = await helloWorldContract.setGreeting(modifiedGreeting)
      await tx.wait()
    }
    catch (err) {
      return
    }
    expect.fail("should have errored")
  })

  it("should be add contract to enabled list", async function () {
    let contractRole = await adminSignerPrecompile.readAllowList(helloWorldContract.address);
    expect(contractRole).to.be.equal(0)

    let enableTx = await adminSignerPrecompile.setEnabled(helloWorldContract.address);
    await enableTx.wait()
    contractRole = await adminSignerPrecompile.readAllowList(helloWorldContract.address);
    expect(contractRole).to.be.equal(1) // 1 = ENABLED
  });


  it("should setGreeting and getHello", async function () {
    const modifiedGreeting = "What's up";
    let tx = await helloWorldContract.setGreeting(modifiedGreeting);
    await tx.wait();

    expect(await helloWorldContract.callStatic.getHello()).to.be.equal(
      modifiedGreeting
    );
  });
});
