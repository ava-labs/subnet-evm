//SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "./ERC20NativeMinter.sol";

// Helper contract to test minting from another contract
contract Minter {
  ERC20NativeMinter token;

  constructor(address tokenAddress) {
    token = ERC20NativeMinter(tokenAddress);
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

