//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./IHelloWorld.sol";

// ExampleDeployerList shows how ContractDeployerAllowList precompile can be used in a smart contract
// All methods of [allowList] can be directly called. There are example calls as tasks in hardhat.config.ts file.
contract ExampleHelloWorld {
  // Precompiled Allow List Contract Address
  address constant HELLO_WORLD_ADDRESS = 0x0300000000000000000000000000000000000000;
  HelloWorld helloWorld = HelloWorld(HELLO_WORLD_ADDRESS);

  function sayHello() public returns (string memory) {
    return helloWorld.sayHello();
  }

  function setGreeting(string calldata greeting) public {
    helloWorld.setGreeting(greeting);
  }
}
