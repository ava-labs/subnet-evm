// (c) 2022-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// SPDX-License-Identifier: MIT

pragma solidity >=0.8.0;

interface IHelloWorld {
  // sayHello returns the string located at [key]
  function sayHello() external returns (string calldata);

  // setGreeting sets the string located at [key]
  function setGreeting(string calldata response) external;
}
