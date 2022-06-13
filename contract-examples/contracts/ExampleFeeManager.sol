//SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "./IAllowList.sol";
import "./IFeeManager.sol";

// ExampleDeployerList shows how ContractDeployerAllowList precompile can be used in a smart contract
// All methods of [allowList] can be directly called. There are example calls as tasks in hardhat.config.ts file.
contract ExampleFeeManager is Ownable {
  // Precompiled Fee Manager Contract Address
  address constant FEE_MANAGER = 0x0200000000000000000000000000000000000003;
  IFeeManager feeManager = IFeeManager(FEE_MANAGER);

  uint256 constant STATUS_NONE = 0;
  uint256 constant STATUS_ENABLED = 1;
  uint256 constant STATUS_ADMIN = 2;

  struct fee_config {
    uint256 gas_limit;
    uint256 target_block_rate;
    uint256 min_base_fee;
    uint256 target_gas;
    uint256 base_fee_change_denominator;
    uint256 min_block_gas_cost;
    uint256 max_block_gas_cost;
    uint256 block_gas_cost_step;
  }

  constructor() Ownable() {}

  function isAdmin(address addr) public view returns (bool) {
    uint256 result = feeManager.readAllowList(addr);
    return result == STATUS_ADMIN;
  }

  function isAllowed(address addr) public view returns (bool) {
    uint256 result = feeManager.readAllowList(addr);
    // if address is ENABLED or ADMIN, it can change the fee
    return result != STATUS_NONE;
  }

  function addAdmin(address addr) public onlyOwner {
    feeManager.setAdmin(addr);
  }

  function addAllowed(address addr) public onlyOwner {
    feeManager.setEnabled(addr);
  }

  function revoke(address addr) public onlyOwner {
    require(_msgSender() != addr, "cannot revoke own role");
    feeManager.setNone(addr);
  }
}
