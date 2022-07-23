import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import { expect } from "chai"
import {
  Contract,
  ContractFactory,
} from "ethers"
import { ethers } from "hardhat"

// make sure this is always an admin for minter precompile
const addressStr: string = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
const HELLO_WORLD_ADDRESS = "0x0300000000000000000000000000000000000000"
const modifiedGreeting = "Hey there sports fan!"

describe("ExampleHelloWorld", function () {
  let address: SignerWithAddress
  let contract: Contract
  before(async function () {
    address = await ethers.getSigner(addressStr);
    const contractF: ContractFactory = await ethers.getContractFactory("ExampleHelloWorld", { signer: address });
    contract = await contractF.deploy();
    await contract.deployed();
    const contractAddress: string = contract.address;
    console.log(`Contract deployed to: ${contractAddress}`);
  })

  it("sayHello", async function () {
    const greeting: string = await contract.sayHello();
    expect(greeting).to.be.equal("Hello World!");
  })

  it("set greeting and say hello", async function() {
    const helloWorld = await ethers.getContractAt("HelloWorld", HELLO_WORLD_ADDRESS, address);
    await helloWorld.setGreeting(modifiedGreeting)
    const greeting: string = await helloWorld.sayHello();
    expect(greeting).to.be.equal(modifiedGreeting)
  })
})
