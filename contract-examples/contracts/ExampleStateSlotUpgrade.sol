//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// ExampleStateSlotUpgrade shows how the Slot state upgrader can be used.
contract ExampleStateSlotUpgrade {
    // Make a 6 gap variable.
    uint256[6] private ______gap;

    // Cap is now at slot 7.
    uint256 public cap = 125000000 ether; // We initialize the cap to 125 million ether.
}