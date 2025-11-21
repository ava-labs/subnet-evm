//SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "./ERC20NativeMinterTest.sol";

// Helper contract to test minting from another contract
contract Minter {
  ERC20NativeMinterTest token;

  constructor(address tokenAddress) {
    token = ERC20NativeMinterTest(tokenAddress);
  }

  function mintdraw(uint amount) external {
    token.mintdraw(amount);
  }

  function deposit(uint value) external {
    token.deposit{value: value}();
  }
  
  // Allow the contract to receive ETH
  receive() external payable {}
}

