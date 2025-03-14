// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

import {DSTest} from "ds-test/src/test.sol";

contract FakeTest is DSTest {
    event NotFromDSTest();

    function logNonDSTest() external {
        emit NotFromDSTest();
    }

    function logString(string memory s) external {
        emit log(s);
    }

    function logNamedAddress(string memory name, address addr) external {
        emit log_named_address(name, addr);
    }
}
