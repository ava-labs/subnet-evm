// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

import { expect } from "chai";
import { ethers } from "hardhat";

import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { Contract, ContractFactory } from "ethers";

// make sure this is always an admin for the precompile
const adminAddress: string = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC";
const COUNTER_ADDRESS = "0x0300000000000000000000000000000000000099";

describe("ExampleCounter", function () {
  let counterContract: Contract;
  let adminSigner: SignerWithAddress;
  let adminSignerPrecompile: Contract;

  before(async function () {
    // Deploy Counter Contract
    const ContractF: ContractFactory = await ethers.getContractFactory(
      "ExampleCounter"
    );
    counterContract = await ContractF.deploy();
    await counterContract.deployed();
    const counterContractAddress: string = counterContract.address;
    console.log(`Contract deployed to: ${counterContractAddress}`);

    adminSigner = await ethers.getSigner(adminAddress);
    adminSignerPrecompile = await ethers.getContractAt(
      "ICounter",
      COUNTER_ADDRESS,
      adminSigner
    );
  });

  it("should getCounter properly", async function () {
    let result = await counterContract.callStatic.getCounter();
    expect(result).to.equal(0);
  });

  it("should incrementByOne and getCounter", async function () {
    let tx = await counterContract.incrementByOne();
    await tx.wait();

    expect(await counterContract.callStatic.getCounter()).to.be.equal(1);
  });

  it("should incrementByX and getCounter", async function () {
    let tx = await counterContract.incrementByX(5);
    await tx.wait();

    expect(await counterContract.callStatic.getCounter()).to.be.equal(6);
  });
});
