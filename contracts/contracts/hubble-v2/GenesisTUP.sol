// SPDX-License-Identifier: BUSL-1.1

pragma solidity 0.8.9;

// unused import; required for a forced contract compilation
import { ProxyAdmin } from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import { TransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

contract GenesisTUP is TransparentUpgradeableProxy {
    // since this goes as a genesis contract, we cannot pass vars in the constructor
    constructor() TransparentUpgradeableProxy(address(0), address(0), "") {}

    function setGenesisAdmin(address admin_) external {
        // it is a known issue that this can be frontran, but we do not worry about it at the moment
        require(_admin() == address(0), "already initialized");
        _changeAdmin(admin_);
    }
}
