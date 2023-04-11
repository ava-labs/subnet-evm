//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "./AllowList.sol";
import "./INativeMinter.sol";
import "./AllowListTest.sol";

address constant MINTER_ADDRESS = 0x0200000000000000000000000000000000000001;
address constant BLACKHOLE_ADDRESS = 0x0100000000000000000000000000000000000000;

contract ERC20NativeMinter is ERC20, AllowList {
  // Precompiled Native Minter Contract Address
  // Designated Blackhole Address
  string private constant TOKEN_NAME = "ERC20NativeMinterToken";
  string private constant TOKEN_SYMBOL = "XMPL";

  INativeMinter nativeMinter = INativeMinter(MINTER_ADDRESS);

  event Deposit(address indexed dst, uint256 wad);
  event Mintdrawal(address indexed src, uint256 wad);

  constructor(uint256 initSupply) ERC20(TOKEN_NAME, TOKEN_SYMBOL) AllowList(MINTER_ADDRESS) {
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

  // Swaps [amount] number of ERC20 token for native coin.
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

  function decimals() public view virtual override returns (uint8) {
    return 18;
  }
}

// TODO:
// this contract adds another (unwanted) layer of indirection
// but it's the easiest way to match the previous HardHat testing functionality.
// Once we completely migrate to DS-test, we can simplify this set of tests.
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
}

contract ERC20NativeMinterTest is AllowListTest {
  function setUp() public {
    // noop
  }

  function test_mintdrawFailure() public {
    ERC20NativeMinter token = new ERC20NativeMinter(1000);
    address tokenAddress = address(token);

    INativeMinter nativeMinter = INativeMinter(MINTER_ADDRESS);

    assertRole(nativeMinter.readAllowList(tokenAddress), AllowList.Role.None);

    try token.mintdraw(100) {
      assertTrue(false, "mintdraw should fail");
    } catch {} // TODO should match on an error to make sure that this is failing in the way that's expected
  }

  function test_addMinter() public {
    ERC20NativeMinter token = new ERC20NativeMinter(1000);
    address tokenAddress = address(token);

    INativeMinter nativeMinter = INativeMinter(MINTER_ADDRESS);

    assertRole(nativeMinter.readAllowList(tokenAddress), AllowList.Role.None);

    nativeMinter.setEnabled(tokenAddress);

    assertRole(nativeMinter.readAllowList(tokenAddress), AllowList.Role.Enabled);
  }

  function test_adminMintdraw() public {
    ERC20NativeMinter token = new ERC20NativeMinter(1000);
    address tokenAddress = address(token);

    address testAddress = address(this);

    INativeMinter nativeMinter = INativeMinter(MINTER_ADDRESS);
    nativeMinter.setEnabled(tokenAddress);

    uint initialTokenBalance = token.balanceOf(testAddress);
    uint initialNativeBalance = testAddress.balance;

    uint amount = 100;

    token.mintdraw(amount);

    assertEq(token.balanceOf(testAddress), initialTokenBalance - amount);
    assertEq(testAddress.balance, initialNativeBalance + amount);
  }

  function test_minterMintdrawFailure() public {
    ERC20NativeMinter token = new ERC20NativeMinter(1000);
    address tokenAddress = address(token);

    Minter minter = new Minter(tokenAddress);
    address minterAddress = address(minter);

    INativeMinter nativeMinter = INativeMinter(MINTER_ADDRESS);
    nativeMinter.setEnabled(tokenAddress);

    uint initialTokenBalance = token.balanceOf(minterAddress);
    uint initialNativeBalance = minterAddress.balance;

    assertRole(initialTokenBalance, AllowList.Role.None);

    try minter.mintdraw(100) {
      assertTrue(false, "mintdraw should fail");
    } catch {} // TODO should match on an error to make sure that this is failing in the way that's expected

    assertEq(token.balanceOf(minterAddress), initialTokenBalance);
    assertEq(minterAddress.balance, initialNativeBalance);
  }

  function test_minterDeposit() public {
    ERC20NativeMinter token = new ERC20NativeMinter(1000);
    address tokenAddress = address(token);

    Minter minter = new Minter(tokenAddress);
    address minterAddress = address(minter);

    INativeMinter nativeMinter = INativeMinter(MINTER_ADDRESS);
    nativeMinter.setEnabled(tokenAddress);

    uint amount = 100;

    nativeMinter.mintNativeCoin(minterAddress, amount);

    uint initialTokenBalance = token.balanceOf(minterAddress);
    uint initialNativeBalance = minterAddress.balance;

    minter.deposit(amount);

    assertEq(token.balanceOf(minterAddress), initialTokenBalance + amount);
    assertEq(minterAddress.balance, initialNativeBalance - amount);
  }

  function test_mintdraw() public {
    ERC20NativeMinter token = new ERC20NativeMinter(1000);
    address tokenAddress = address(token);

    Minter minter = new Minter(tokenAddress);
    address minterAddress = address(minter);

    INativeMinter nativeMinter = INativeMinter(MINTER_ADDRESS);
    nativeMinter.setEnabled(tokenAddress);

    uint amount = 100;

    uint initialNativeBalance = minterAddress.balance;
    assertRole(initialNativeBalance, AllowList.Role.None);

    token.mint(minterAddress, amount);

    uint initialTokenBalance = token.balanceOf(minterAddress);
    assertEq(initialTokenBalance, amount);

    minter.mintdraw(amount);

    assertRole(token.balanceOf(minterAddress), AllowList.Role.None);
    assertEq(minterAddress.balance, amount);
  }
}
