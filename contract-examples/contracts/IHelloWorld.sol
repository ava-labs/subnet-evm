// (c) 2022-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// SPDX-License-Identifier: MIT

pragma solidity >=0.8.0;

interface HelloWorld {
    function sayHello() external returns (string calldata);

    // SetGreeting
    function setGreeting(string calldata recipient) external;
}
