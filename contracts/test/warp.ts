// (c) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { expect } from "chai";
import {
  Contract,
  ContractFactory,
} from "ethers"
import { ethers } from "hardhat"

const fundedAddress: string = "0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
const WARP_ADDRESS = "0x0200000000000000000000000000000000000005";

let fundedSigner: SignerWithAddress
let warpPrecompile: Contract

let warpExample: Contract

var blockchainIDA = ""
var blockchainIDB = ""
var payload = ""
var task = process.env.WARP_TEST_TASK || "send"

describe("Warp", function () {
  before(async function () {
    fundedSigner = await ethers.getSigner(fundedAddress);
    warpPrecompile = await ethers.getContractAt("WarpMessenger", WARP_ADDRESS, fundedSigner);
    console.log(`Precompile contract handle at address: ${warpPrecompile.address}`);

    const contractF: ContractFactory = await ethers.getContractFactory("ExampleWarp", { signer: fundedSigner });
    warpExample = await contractF.deploy();
    console.log(`WarpExample deployed at address: ${warpExample.address}`);
  })

  it("warp action", async function () {
    // TODO: implement warp HardHat test on a single node staking network
    // send a warp message to a random destination
    // check that delivering the message unsigned fails
    // input the signed message (this can either be a constant or an environment variable populated from the precompile test) and deliver it
    // validate the contents using ExampleWarp.sol which will simply take in arguments and assert that it is correct
  })
})
