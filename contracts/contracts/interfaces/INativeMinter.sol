//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "./IAllowList.sol";

interface INativeMinter is IAllowList {
  event MintNativeCoin(address indexed sender, address indexed addr, uint256 indexed amount);

  // Mint [amount] number of native coins and send to [addr]
  function mintNativeCoin(address addr, uint256 amount) external;
}
