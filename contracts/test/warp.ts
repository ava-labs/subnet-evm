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

const WARP_ADDRESS = "0x0200000000000000000000000000000000000005";
const senderAddress: string = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
const destAddress: string = "0x0550000000000000000000000000000000000000"

describe("ExampleWarp", function () {
  this.timeout("30s")

  let owner: SignerWithAddress
  let contract: Contract
  before(async function () {
    owner = await ethers.getSigner(senderAddress);
    contract = await ethers.getContractAt("WarpMessenger", WARP_ADDRESS, owner)
  });

  it("contract should be to send warp message", async function () {
    let destId = process.env["DESTINATION_CHAIN_ID"];
    let payload = process.env["PAYLOAD"];
    let expectedUnsignedMessage = process.env["EXPECTED_UNSIGNED_MESSAGE"];
    console.log(`Sending warp message to chain: ${destId}, address: ${destAddress} with payload ${payload}`);
    let destIdHex = "0x" + destId.toString().padStart(32, "0");
    expect(ethers.utils.isHexString(destIdHex)).to.be.true;
    let payloadHex = "0x" + payload.toString()
    expect(ethers.utils.isHexString(payloadHex)).to.be.true;
    let unsignedMessageHex = "0x" + expectedUnsignedMessage.toString()
    expect(ethers.utils.isHexString(unsignedMessageHex)).to.be.true;

    await expect(contract.sendWarpMessage(destIdHex, destAddress, payloadHex))
    .to.emit(contract, 'SendWarpMessage')
    .withArgs(destIdHex, destAddress, senderAddress, unsignedMessageHex);
  })

})
