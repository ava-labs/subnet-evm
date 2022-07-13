// (c) 2022-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// SPDX-License-Identifier: MIT

pragma solidity >=0.8.0;

interface HelloWorldInterface {
    function sayHello() external;

    // setRecipient
    function setReceipient(string calldata recipient) external;
}
