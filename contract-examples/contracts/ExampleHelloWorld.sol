//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./IHelloWorld.sol";

// ExampleHelloWorld shows how the HelloWorld precompile can be used in a smart conract
contract ExampleHelloWorld {
  address constant HELLO_WORLD_ADDRESS = 0x0200000000000000000000000000000000000004;
  IHelloWorld helloWorld = IHelloWorld(HELLO_WORLD_ADDRESS);

  function getHello() public returns (string memory) {
    return helloWorld.sayHello();
  }

  function setGreeting(string calldata greeting) public {
    helloWorld.setGreeting(greeting);
  }
}
