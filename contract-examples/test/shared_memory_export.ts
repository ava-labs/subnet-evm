// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { expect } from "chai";
import {
  BigNumber,
  Contract,
  ContractFactory,
} from "ethers"
import { ethers } from "hardhat"
import ts = require("typescript");

// make sure this is always an admin for reward manager precompile
const fundedAddr: string = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
const SHARED_MEMORY_ADDRESS = "0x0200000000000000000000000000000000000005";

describe("SharedMemoryExport", function () {
  this.timeout("30s")
  let fundedSigner: SignerWithAddress
  let contract: Contract
  let signer1: SignerWithAddress
  let signer2: SignerWithAddress
  let precompile: Contract
  let blockchainIDA: string
  let blockchainIDB: string

  before(async function () {
    // Populate blockchainIDs from the environment variables
    blockchainIDA = "0x" + process.env.BLOCKCHAIN_ID_A
    blockchainIDB = "0x" + process.env.BLOCKCHAIN_ID_B
    console.log("blockchainIDA %s, blockchainIDB: %s", blockchainIDA, blockchainIDB)

    

    fundedSigner = await ethers.getSigner(fundedAddr);
    signer1 = (await ethers.getSigners())[1]
    signer2 = (await ethers.getSigners())[2]
    const Contract: ContractFactory = await ethers.getContractFactory("ERC20SharedMemory", { signer: fundedSigner })
    contract = await Contract.deploy(1000000000)
    await contract.deployed()
    const contractAddress: string = contract.address
    console.log(`Contract deployed to: ${contractAddress}`)

    precompile = await ethers.getContractAt("ISharedMemory", SHARED_MEMORY_ADDRESS, fundedSigner);
  });

  it("exportAVAX via EOA", async function () {
    let startingBalance: BigNumber = await ethers.provider.getBalance(fundedAddr)
    console.log("Starting balance of %d", startingBalance)
    
    // call exportAVAX via EOA
    let value = ethers.utils.parseUnits("1", "ether")
    let exportAVAXTx = await precompile.exportAVAX(blockchainIDB, 0, 1, [fundedAddr], { value: value })
    let receipt = await exportAVAXTx.wait()

    // find the tx fee
    let gasUsed: BigNumber = receipt.cumulativeGasUsed
    let gasPrice: BigNumber = receipt.effectiveGasPrice
    let txFee = gasUsed.mul(gasPrice)
    
    // verify balance update
    let updatedBalance = await ethers.provider.getBalance(fundedAddr)
    console.log("Starting balance: %d, Updated balance: %d", startingBalance, updatedBalance)
    expect(updatedBalance.lt(startingBalance)).to.be.true
    let expectedBalance = startingBalance.sub(value).sub(txFee)
    console.log("expected balance %d", expectedBalance)
    expect(expectedBalance.eq(updatedBalance)).to.be.true
  });

  it("exportAVAX via contract", async function () {
    // TODO: implement
    // let fundTx = 
    // let startingBalance: BigNumber = await ethers.provider.getBalance(fundedAddr)
    // console.log("Starting balance of %d", startingBalance)
    
    // // call exportAVAX via EOA
    // let value = ethers.utils.parseUnits("1", "ether")
    // let exportAVAXTx = await precompile.exportAVAX(blockchainIDB, 0, 1, [fundedAddr], { value: value })
    // await exportAVAXTx.wait()
    
    // // verify balance update
    // let updatedBalance = await ethers.provider.getBalance(fundedAddr)
    // console.log("Starting balance: %d, Updated balance: %d", startingBalance, updatedBalance)
    // expect(updatedBalance.lt(startingBalance)).to.be.true
    // let expectedBalance = startingBalance.sub(value)
    // console.log("expected balance %d", expectedBalance)
    // expect(expectedBalance.lt(updatedBalance)).to.be.true
  });

  // TODO: export non-AVAX asset from EOA and contract
});
