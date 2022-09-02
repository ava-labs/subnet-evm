//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./IHelloWorld.sol";

contract MockHelloWorld is IHelloWorld {
  string originalGreeting = "Hello World!";

  function sayHello() public view override returns (string memory) {
    return originalGreeting;
  }

  function setGreeting(string calldata greeting) public override {
    originalGreeting = greeting;
  }
}
