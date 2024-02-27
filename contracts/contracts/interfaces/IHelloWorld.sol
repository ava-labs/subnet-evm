// (c) 2022-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// SPDX-License-Identifier: MIT

pragma solidity >=0.8.0;
import "./IAllowList.sol";

interface IHelloWorld is IAllowList {
  event GreetingChanged(address indexed sender, string oldGreeting, string newGreeting);
  // sayHello returns the stored greeting string
  function sayHello() external view returns (string calldata result);

  // setGreeting  stores the greeting string
  function setGreeting(string calldata response) external;
}
