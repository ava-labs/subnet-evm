//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "./ISharedMemory.sol";
import "ds-test/src/test.sol";

contract ERC20SharedMemory is ERC20 {
  // Precompiled Shared Memory Contract Address
  address constant SHARED_MEMORY_ADDRESS = 0x0200000000000000000000000000000000000006;
  ISharedMemory sharedMemory = ISharedMemory(SHARED_MEMORY_ADDRESS);
  string private constant TOKEN_NAME = "ERC20SharedMemoryExampleToken";
  string private constant TOKEN_SYMBOL = "XMPL";

  bytes32 public importableAssetID;
  address public owner;

  event Deposit(address indexed dst, uint256 wad);
  event Mintdrawal(address indexed src, uint256 wad);

  constructor(uint256 initSupply) ERC20(TOKEN_NAME, TOKEN_SYMBOL) {
    // Mints INIT_SUPPLY to owner
    _mint(_msgSender(), initSupply);
    owner = _msgSender();
  }

  function setImportableAssetID(bytes32 assetID) public {
    require(
      msg.sender == owner,
      "ERC20SharedMemory: only owner can set importableAssetID");
    importableAssetID = assetID;
  }

  function exportUTXO(
      bytes32 destinationChainID, uint64 amount, address addr) public {
    // Transfer tokens to this contract
    // require(
    //   transferFrom(_msgSender(), address(this), uint(amount)),
    //   "ERC20SharedMemory: transferFrom failed"
    // );
    // TODO: figure out why the allowance is not working

    _transfer(_msgSender(), address(this), amount);


    // TODO: burn the tokens

    // Only support exporting to one address for now
    address[] memory addrs = new address[](1);
    addrs[0] = addr;

    // Export tokens to destination chain
    sharedMemory.exportUTXO(amount, destinationChainID, 0, 1, addrs);
  }

  function importUTXO(bytes32 sourceChain, bytes32 utxoID) public {
    // Import tokens from source chain
    // TODO: I don't think we need to expose all the utxo details,
    // the precompile already verifies timelock and threshold.
    (uint64 amount, bytes32 assetID, , , address[] memory addrs) = sharedMemory.importUTXO(
      sourceChain, utxoID);

    // Require the UTXO only has one spender and it matches the caller
    require(addrs.length == 1, "ERC20SharedMemory: invalid UTXO");
    require(addrs[0] == _msgSender(), "ERC20SharedMemory: invalid spender");
    require(assetID == importableAssetID, "ERC20SharedMemory: invalid assetID");

    // Mint tokens to the caller
    _mint(_msgSender(), amount);
  }

  function decimals() public view virtual override returns (uint8) {
    return 18;
  }

  function approvalAmount() public view returns (uint) {
    return allowance(_msgSender(), address(this));
  }
}

contract ERC20SharedMemoryTest is DSTest {
  address constant SHARED_MEMORY_ADDRESS = 0x0200000000000000000000000000000000000006;
  ISharedMemory sharedMemory = ISharedMemory(SHARED_MEMORY_ADDRESS);
  ERC20SharedMemory private erc20;

  function setUp() public {
    erc20 = new ERC20SharedMemory(1_000_000_000);
  }

  function approvalAmount() public view returns (uint) {
    return erc20.approvalAmount();
  }

  function test_exportAVAX(
      uint amount, bytes32 destinationChainID, address addr) public {
    address testAddress = address(this);
    uint balanceNative = testAddress.balance;
    assertEq(balanceNative, amount);

    address[] memory addrs = new address[](1);
    addrs[0] = addr;

    sharedMemory.exportAVAX{value: amount}(destinationChainID, 0, 1, addrs);

    // balance should decrease after a successful call to exportAVAX
    balanceNative = testAddress.balance;
    assertEq(balanceNative, 0);
  }

  function test_approveERC20(uint amount) public { 
    erc20.approve(address(erc20), amount);
  }

  function test_exportERC20(
      uint amount, bytes32 destinationChainID, address addr) public {
    address testAddress = address(this);
    uint balance = erc20.balanceOf(testAddress);
    assertEq(balance, amount);
    
    erc20.exportUTXO(destinationChainID, uint64(amount), addr);

    // balance should decrease after a successful call to exportUTXO
    balance = erc20.balanceOf(testAddress);
    assertEq(balance, 0);
  }
  
  function test_importAVAX(
      bytes32 sourceChain, bytes32 utxoID, uint expectedBalance) public {
    address testAddress = address(this);
    uint balanceNative = testAddress.balance;
    assertEq(balanceNative, 0);

    sharedMemory.importAVAX(sourceChain, utxoID);

    // balance should increase after a successful call to importAVAX
    balanceNative = testAddress.balance;
    assertEq(balanceNative, expectedBalance);
  }

  function test_importERC20(
      bytes32 sourceChain, bytes32 utxoID, uint expectedBalance) public {
    address testAddress = address(this);
    uint balance = erc20.balanceOf(testAddress);
    assertEq(balance, 1_000_000_000); // TODO: fix this

    // The test expects to import assets from sourceChain exported from a
    // contract with the same address as this contract.
    bytes32 importAssetID = keccak256(
      abi.encodePacked(sourceChain, address(erc20)));
    erc20.setImportableAssetID(importAssetID);
    erc20.importUTXO(sourceChain, utxoID);

    // balance should increase after a successful call to importUTXO
    balance = erc20.balanceOf(testAddress);
    assertEq(balance, 1_000_000_000+expectedBalance);
  }

  receive() external payable {}
}