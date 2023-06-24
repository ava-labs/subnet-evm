//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./interfaces/ICounter.sol";

contract ExampleCounter {
  address constant COUNTER_ADDRESS = 0x0300000000000000000000000000000000000099;
  ICounter counter = ICounter(COUNTER_ADDRESS);

  function getCounter() public view returns (uint256) {
    return counter.getCounter();
  }

  function incrementByOne() public {
    counter.incrementByOne();
  }

  function incrementByX(uint256 x) public {
    counter.incrementByX(x);
  }
}