//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "./ISharedMemory.sol";

contract ERC20SharedMemory is ERC20 {
  // Precompiled Native Minter Contract Address
  address constant SHARED_MEMORY_ADDRESS = 0x0200000000000000000000000000000000000005;
  ISharedMemory sharedMemory = ISharedMemory(SHARED_MEMORY_ADDRESS);
  string private constant TOKEN_NAME = "ERC20NativeMinterToken";
  string private constant TOKEN_SYMBOL = "XMPL";

  event Deposit(address indexed dst, uint256 wad);
  event Mintdrawal(address indexed src, uint256 wad);

  constructor(uint256 initSupply) ERC20(TOKEN_NAME, TOKEN_SYMBOL) {
    // Mints INIT_SUPPLY to owner
    _mint(_msgSender(), initSupply);
  }

  function exportAVAX(bytes32 destinationChainID, uint64 locktime, uint64 threshold, address[] calldata addrs) external payable {
    sharedMemory.exportAVAX{value: msg.value}(destinationChainID, locktime, threshold, addrs);
  }

  function decimals() public view virtual override returns (uint8) {
    return 18;
  }
}
