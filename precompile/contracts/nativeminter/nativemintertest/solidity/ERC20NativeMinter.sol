//SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "precompile/allowlist/allowlisttest/solidity/AllowList.sol";
import "./INativeMinter.sol";

// Designated Blackhole Address
address constant BLACKHOLE_ADDRESS = 0x0100000000000000000000000000000000000000;

contract ERC20NativeMinter is ERC20, Ownable, AllowList {
  string private constant TOKEN_NAME = "ERC20NativeMinterToken";
  string private constant TOKEN_SYMBOL = "XMPL";

  INativeMinter nativeMinter;

  event Deposit(address indexed dst, uint256 wad);
  event Mintdrawal(address indexed src, uint256 wad);

  constructor(address nativeMinterPrecompile, uint256 initSupply) 
    ERC20(TOKEN_NAME, TOKEN_SYMBOL) 
    Ownable(msg.sender)
    AllowList(nativeMinterPrecompile) 
  {
    nativeMinter = INativeMinter(nativeMinterPrecompile);
    // Mints init supply to msg.sender
    _mint(msg.sender, initSupply);
  }

  // Mints [amount] number of ERC20 token to [to] address.
  function mint(address to, uint256 amount) external onlyOwner {
    _mint(to, amount);
  }

  // Burns [amount] number of ERC20 token from [from] address.
  function burn(address from, uint256 amount) external onlyOwner {
    _burn(from, amount);
  }

  // Swaps [amount] number of ERC20 token for native coin.
  function mintdraw(uint256 wad) external {
    // Burn ERC20 token first.
    _burn(msg.sender, wad);
    // Mints [amount] number of native coins (gas coin) to [msg.sender] address.
    // Calls NativeMinter precompile through INativeMinter interface.
    nativeMinter.mintNativeCoin(msg.sender, wad);
    emit Mintdrawal(msg.sender, wad);
  }

  // Swaps [amount] number of native gas coins for ERC20 tokens.
  function deposit() external payable {
    // Burn native token by sending to BLACKHOLE_ADDRESS
    payable(BLACKHOLE_ADDRESS).transfer(msg.value);
    // Mint ERC20 token.
    _mint(msg.sender, msg.value);
    emit Deposit(msg.sender, msg.value);
  }

  function decimals() public view virtual override returns (uint8) {
    return 18;
  }
}
