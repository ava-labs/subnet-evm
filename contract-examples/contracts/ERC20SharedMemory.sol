//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "./ISharedMemory.sol";
import "ds-test/src/test.sol";

contract ERC20SharedMemory is ERC20 {
  // Precompiled Shared Memory Contract Address
  address constant SHARED_MEMORY_ADDRESS = 0x0200000000000000000000000000000000000005;
  ISharedMemory sharedMemory = ISharedMemory(SHARED_MEMORY_ADDRESS);
  string private constant TOKEN_NAME = "ERC20SharedMemoryExampleToken";
  string private constant TOKEN_SYMBOL = "XMPL";

  event Deposit(address indexed dst, uint256 wad);
  event Mintdrawal(address indexed src, uint256 wad);

  constructor(uint256 initSupply) ERC20(TOKEN_NAME, TOKEN_SYMBOL) {
    // Mints INIT_SUPPLY to owner
    _mint(_msgSender(), initSupply);
  }

  function decimals() public view virtual override returns (uint8) {
    return 18;
  }

  receive() external payable {}
}

contract ERC20SharedMemoryTest is DSTest {
  address constant SHARED_MEMORY_ADDRESS = 0x0200000000000000000000000000000000000005;
  ISharedMemory sharedMemory = ISharedMemory(SHARED_MEMORY_ADDRESS);
  ERC20SharedMemory private erc20;

  function setUp() public {
    erc20 = new ERC20SharedMemory(100);
  }

  function test_exportAVAX(
      uint value, bytes32 destinationChainID, address addr) public {
    address testAddress = address(this);
    uint balanceNative = testAddress.balance;
    assertEq(balanceNative, value);

    address[] memory addrs = new address[](1);
    addrs[0] = addr;

    sharedMemory.exportAVAX{value: value}(destinationChainID, 0, 1, addrs);

    // balance should decrease after a successful call to exportAVAX
    balanceNative = testAddress.balance;
    assertEq(balanceNative, 0);
  }
  
  function test_importAVAX(
      bytes32 sourceChain, bytes32 utxoID, uint expectedValue) public {
    address testAddress = address(this);
    uint balanceNative = testAddress.balance;
    assertEq(balanceNative, 0);

    sharedMemory.importAVAX(sourceChain, utxoID);

    // balance should increase after a successful call to importAVAX
    balanceNative = testAddress.balance;
    assertEq(balanceNative, expectedValue);
  }

  receive() external payable {}
}