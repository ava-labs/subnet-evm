//SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "precompile/contracts/warp/warpbindings/IWarpMessenger.sol";


// Assuming there is some IPrecompile.sol and a respective ABI, there are two scenarios for a Go TestPrecompile:
// 1. TestPrecompile -> abigen bindings -> precompile.
// 2. TestPrecompile -> intermediary contract1 -> IPrecompile -> precompile.
//
// Since we expect an end user to consume IPrecompile in their own smart contract, we MUST be including it as 
// part of the SUT (i.e. IPrecompile + precompile). Achieving this with (2) is self-evident. 
// To achieve this with (1) we MUST have the ABI+bindings generated from IPrecompile and the results confirmed by CI2.
// 
// It is not sufficient to say that the precompile was generated from some ABI of unknown provenance and that it's the 
// source of truth because then the tests exclude a key element on which our users depend.
//
// Note that the intermediary contract can be called any way, e.g. with its own abigen bindings, but for tests 
// that don't require dynamic arguments it can be no more than a constructor with revert()s, and then successful 
//deployment is the test.

contract WarpTest {
  IWarpMessenger private warp;

  constructor(address warpPrecompile) {
    warp = IWarpMessenger(warpPrecompile);
  }

  // Calls the getBlockchainID function on the precompile
  function getBlockchainID() external view returns (bytes32) {
    return warp.getBlockchainID();
  }

  // Calls the sendWarpMessage function on the precompile
  function sendWarpMessage(bytes calldata payload) external returns (bytes32 messageID) {
    return warp.sendWarpMessage(payload);
  }

  // Calls the getVerifiedWarpMessage function on the precompile
  function getVerifiedWarpMessage(
    uint32 index
  ) external view returns (WarpMessage memory message, bool valid) {
    return warp.getVerifiedWarpMessage(index);
  }

  // Calls the getVerifiedWarpBlockHash function on the precompile
  function getVerifiedWarpBlockHash(
    uint32 index
  ) external view returns (WarpBlockHash memory warpBlockHash, bool valid) {
    return warp.getVerifiedWarpBlockHash(index);
  }
}

