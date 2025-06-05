//SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

interface IBlockList {
  event AddressBlocked(address indexed account, string reason);
  event AddressUnblocked(address indexed account, string reason);
  event AdminChanged(address indexed account);

  // Set [addr] to be the admin of the blocklist.
  function changeAdmin(address addr) external;

  // Set [addr] to be added to the blocklist.
  function blockAddress(address addr, string calldata reason) external;

  // Set [addr] to be removed from the blocklist.
  function unblockAddress(address addr, string calldata reason) external;

  // Read the blocklist status of [addr].
  function readBlockList(address addr) external view returns (uint256 role);

  // Read the admin address of the blocklist.
  function admin() external view returns (address addr);
}
