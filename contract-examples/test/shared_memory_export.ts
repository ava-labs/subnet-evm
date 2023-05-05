// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { expect } from "chai";
import {
  BigNumber,
  Contract,
  ContractFactory,
  Event,
} from "ethers"
import { ethers } from "hardhat"
import ts = require("typescript");

const FUNDED_ADDRESS: string = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
const SHARED_MEMORY_ADDRESS = "0x0200000000000000000000000000000000000005";

enum BlockchainName {A, B}
const getBlockchainId = (name: BlockchainName) => {
  switch (name) {
    case BlockchainName.A: 
      return "0x" + process.env.BLOCKCHAIN_ID_A
    case BlockchainName.B: 
      return "0x" + process.env.BLOCKCHAIN_ID_B
  }
}

describe("SharedMemoryExport", function () {
  this.timeout("30s")

  // Populate blockchainIDs from the environment variables
  const [blockchainIDA, blockchainIDB] = [
    getBlockchainId(BlockchainName.A),
    getBlockchainId(BlockchainName.B),
  ]
  let fundedSigner: SignerWithAddress
  let signer1: SignerWithAddress
  let signer2: SignerWithAddress
  let precompile: Contract
  let contract: Contract

  // TODO: see why the nonce is not matching the import test if I do beforeEach
  // ends up with different contract address for the 2nd test.
  // Maybe better to leave this as before anyway.
  before('Setup DS-Test contract', async function () {
    console.log("blockchainIDA %s, blockchainIDB: %s", blockchainIDA, blockchainIDB)

    const signers = await ethers.getSigners();
    [fundedSigner, signer1, signer2] = signers

    const sharedMemory = await ethers.getContractAt(
       "ISharedMemory", SHARED_MEMORY_ADDRESS, fundedSigner)

    return  ethers.getContractFactory(
        "ERC20SharedMemoryTest", { signer: fundedSigner })
      .then(factory => factory.deploy())
      .then(contract => {
        this.testContract = contract
        return contract.deployed().then(() => contract)
      })
      .then(contract => contract.setUp())
      .then(tx => tx.wait())
  })


  it("exportAVAX via contract", async function () {
    let testContract: Contract = this.testContract;
    console.log("testContract", testContract.address)
    
    // Fund the contract
    // Note 1 gwei is the minimum amount of AVAX that can be sent due to
    // denomination adjustment in exported UTXOs.
    let amount = ethers.utils.parseUnits("1", "gwei") 
    let tx = await fundedSigner.sendTransaction({
      to: testContract.address,
      value: amount,
    })
    let receipt = await tx.wait()
    expect(receipt.status == 1).to.be.true

    // ExportAVAX
    // Note we export AVAX to testContract.address, which is the contract we
    // just deployed. This is because the import test will also deploy a
    // contract from the same account with the same nonce on the other
    // blockchain.
    let unsignedTx = await testContract.populateTransaction.test_exportAVAX(
       amount, blockchainIDB, testContract.address)
    let signedTx = await fundedSigner.sendTransaction(unsignedTx);
    let txReceipt = await signedTx.wait()
    console.log("txReceipt", txReceipt.status)
    expect(await testContract.callStatic.failed()).to.be.false

    // Verify logs were emitted as expected
    let foundLog = txReceipt.logs.find(
      (log: Event, _: any, __: any) => 
        log.address === SHARED_MEMORY_ADDRESS &&
        log.topics.length === 2 && // TODO: review the indexed vs. non-indexed log data
        // TODO: get the string from the contract abi
        log.topics[0] === ethers.utils.id("ExportAVAX(uint64,bytes32,uint64,uint64,address[])") &&
        log.topics[1] == blockchainIDB // destination
    )
    // TODO: consider verifying more about the logs
    expect(foundLog).to.exist;
  })


  it("exportUTXO via contract", async function () {
    let testContract: Contract = this.testContract;
    console.log("testContract", testContract.address)

    // Allow the ERC20 contract to spend tokens on behalf of the test contract
    let amount = 1_000_000_000;
    // let tx = await testContract.test_approveERC20(amount);
    // let receipt = await tx.wait();
    // expect(receipt.status == 1).to.be.true;

    let approvalAmount = await testContract.callStatic.approvalAmount();
    console.log("approvalAmount", approvalAmount);

    // ExportERC20
    // Note we export ERC20 to testContract.address, which is the contract we
    // just deployed. This is because the import test will also deploy a
    // contract from the same account with the same nonce.
   let unsignedTx = await testContract.populateTransaction.test_exportERC20(
     amount, blockchainIDB, testContract.address)
   let signedTx = await fundedSigner.sendTransaction(unsignedTx);
   let txReceipt = await signedTx.wait()
   console.log("txReceipt", txReceipt.status)
   expect(await testContract.callStatic.failed()).to.be.false

   // Verify logs were emitted as expected
   let foundLog = txReceipt.logs.find(
     (log: Event, _: any, __: any) => 
       log.address === SHARED_MEMORY_ADDRESS &&
       log.topics.length === 3 && // TODO: review the indexed vs. non-indexed log data
       // TODO: get the string from the contract abi
       log.topics[0] === ethers.utils.id("ExportUTXO(uint64,bytes32,bytes32,uint64,uint64,address[])") &&
       log.topics[1] == blockchainIDB // destination
       // TODO: verify the assetID in log.topics[2]
   )
   // TODO: consider verifying more about the logs
   expect(foundLog).to.exist;
  })
});
