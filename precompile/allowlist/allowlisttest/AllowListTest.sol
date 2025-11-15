//SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "precompile/allowlist/IAllowList.sol";
import "precompile/allowlist/allowlisttest/AllowList.sol";
import "precompile/precompiletest/DSTest.sol";

contract AllowListTestHelper is DSTest {
    function assertRole(uint result, AllowList.Role role) internal {
        assertEq(result, uint(role));
    }
}

contract AllowListTest is AllowList {
    // Precompiled Allow List Contract Address
    constructor(address precompileAddr) AllowList(precompileAddr) {}

    function deployContract() public {
        new Example();
    }
}

// This is an empty contract that can be used to test contract deployment
contract Example {}
