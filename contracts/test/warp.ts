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
var destAddress: string

describe("IWarpMessenger", function () {
  this.timeout("30s")

  let owner: SignerWithAddress
  let contract: Contract
  before(async function () {
    senderAddress = process.env["SENDER_ADDRESS"];
    destAddress = process.env["DESTINATION_ADDRESS"];
    owner = await ethers.getSigner(senderAddress);
    contract = await ethers.getContractAt("IWarpMessenger", WARP_ADDRESS, owner)
  });

  it("contract should be to send warp message", async function () {
    let destId = process.env["DESTINATION_CHAIN_ID"];
    let payload = process.env["PAYLOAD"];
    let expectedUnsignedMessage = process.env["EXPECTED_UNSIGNED_MESSAGE"];
    let destIdHex = "0x" + destId.toString().padStart(32, "0");
    expect(ethers.utils.isHexString(destIdHex)).to.be.true;
    let payloadHex = "0x" + payload.toString()
    expect(ethers.utils.isHexString(payloadHex)).to.be.true;
    let unsignedMessageHex = "0x" + expectedUnsignedMessage.toString()
    expect(ethers.utils.isHexString(unsignedMessageHex)).to.be.true;

    console.log(`Sending warp message to chain: ${destIdHex}, address: ${destAddress} with payload ${payloadHex}`);

    await expect(contract.sendWarpMessage(destIdHex, destAddress, payloadHex))
    .to.emit(contract, 'SendWarpMessage')
    .withArgs(destIdHex, destAddress, senderAddress, unsignedMessageHex);
  })

  it("should be able to fetch correct blockchain ID", async function () {
    let sourceID = process.env["SOURCE_CHAIN_ID"];
    let sourceIDHex = "0x" + sourceID.toString().padStart(32, "0");
    expect(ethers.utils.isHexString(sourceIDHex)).to.be.true;
  })
})
