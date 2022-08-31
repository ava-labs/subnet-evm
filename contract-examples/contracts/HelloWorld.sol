//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./IHelloWorld.sol";

// HelloWorld shows how the HelloWorld precompile can be used in a smart conract
contract HelloWorld {
  string modifiedGreeting = "Hello World!";
  address constant HELLO_WORLD_ADDRESS = 0x0200000000000000000000000000000000000004;

  IHelloWorld helloWorld = IHelloWorld(HELLO_WORLD_ADDRESS);

  function sayHello() public returns (string memory) {
    return helloWorld.sayHello();
  }

  function setGreeting(string calldata greeting) public {
    helloWorld.setGreeting(greeting);
  }
}
