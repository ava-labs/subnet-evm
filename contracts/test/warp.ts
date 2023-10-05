// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { expect } from "chai";
import {
  Contract,
} from "ethers"
import { ethers } from "hardhat"

const WARP_ADDRESS = "0x0200000000000000000000000000000000000005";
var senderAddress: string
// Expected to be a hex string
var payload: string
var expectedUnsignedMessage: string
var sourceID: string

describe("IWarpMessenger", function () {
  this.timeout("30s")

  let owner: SignerWithAddress
  let contract: Contract
  before(async function () {
    senderAddress = process.env["SENDER_ADDRESS"];
    owner = await ethers.getSigner(senderAddress);
    contract = await ethers.getContractAt("IWarpMessenger", WARP_ADDRESS, owner)

    payload = process.env["PAYLOAD"];
    expectedUnsignedMessage = process.env["EXPECTED_UNSIGNED_MESSAGE"];
    sourceID = process.env["SOURCE_CHAIN_ID"];
  });

  it("contract should be to send warp message", async function () {
    expect(ethers.utils.isHexString(payload)).to.be.true;
    expect(ethers.utils.isHexString(expectedUnsignedMessage)).to.be.true;

    console.log(`Sending warp message with payload ${payload}, expected unsigned message ${expectedUnsignedMessage}`);

    await expect(contract.sendWarpMessage(payload))
      .to.emit(contract, 'SendWarpMessage')
      .withArgs(senderAddress, expectedUnsignedMessage);
  })

  it("should be able to fetch correct blockchain ID", async function () {
    expect(ethers.utils.isHexString(sourceID)).to.be.true;
    let blockchainID = await contract.getBlockchainID();
    expect(blockchainID).to.be.equal(sourceID);
  })
})
