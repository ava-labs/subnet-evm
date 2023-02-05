//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// ExampleContractAfter shows how the CodeUpgrader precompile can be used to upgrade the code of a smart contract.
contract ExampleCodeUpgraderAfter {
    // Make a 6 gap variable.
    uint256[7] private ______gap;

    // Cap is now at slot 0x7.
    uint256 public cap = 125000000 ether; // We initialize the cap to 125 million ether.

    event CapUpdated(uint256 newCap);

    // Add a method to update the cap.
    function updateCap(uint256 _cap) public {
        cap = _cap;
        emit CapUpdated(_cap);
    }
}