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
