//SPDX-License-Identifier: MIT
pragma solidity >=0.6.2;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "./INativeMinter.sol";

contract ERC20NativeMinter is ERC20, Ownable {
  // Precompiled Native Minter Contract Address
  address constant MINTER_ADDRESS = 0x0200000000000000000000000000000000000001;
  // Designated Blackhole Address
  address constant BLACKHOLE_ADDRESS = 0x0100000000000000000000000000000000000000;
  INativeMinter nativeMinter = INativeMinter(MINTER_ADDRESS);
  string private TOKEN_NAME = "ERC20NativeMinterToken";
  string private TOKEN_SYMBOL = "XMPL";

  event Deposit(address indexed dst, uint256 wad);
  event Mintdrawal(address indexed src, uint256 wad);

  constructor(uint256 initSupply) ERC20(TOKEN_NAME, TOKEN_SYMBOL) {
    // Mints INIT_SUPPLY to owner
    _mint(_msgSender(), initSupply);
  }

  // Mints [amount] number of ERC20 token to [to] address.
  function mint(address to, uint256 amount) external onlyOwner {
    _mint(to, amount);
  }

  // Burns [amount] number of ERC20 token from [from] address.
  function burn(address from, uint256 amount) external onlyOwner {
    _burn(from, amount);
  }

  // Swaps [amount] number of ERC20 token for native gas coins.
  function mintdraw(uint256 wad) external {
    // Burn ERC20 token first.
    _burn(_msgSender(), wad);
    // Mints [amount] number of native coins (gas coin) to [msg.sender] address.
    // Calls NativeMinter precompile through INativeMinter interface.
    nativeMinter.mintNativeCoin(_msgSender(), wad);
    emit Mintdrawal(_msgSender(), wad);
  }

  // Swaps [amount] number of native gas coins for ERC20 tokens.
  function deposit() external payable {
    // Burn native token by sending to BLACKHOLE_ADDRESS
    payable(BLACKHOLE_ADDRESS).transfer(msg.value);
    // Mint ERC20 token.
    _mint(_msgSender(), msg.value);
    emit Deposit(_msgSender(), msg.value);
  }
}
