//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./interfaces/IHelloWorld.sol";

address constant HELLO_WORLD_ADDRESS = 0x0300000000000000000000000000000000000000;

// ExampleHelloWorld shows how the HelloWorld precompile can be used in a smart contract.
contract ExampleHelloWorld {
  IHelloWorld helloWorld = IHelloWorld(HELLO_WORLD_ADDRESS);

  function sayHello() public view returns (string memory) {
    return helloWorld.sayHello();
  }

  function setGreeting(string calldata greeting) public {
    helloWorld.setGreeting(greeting);
  }
}
