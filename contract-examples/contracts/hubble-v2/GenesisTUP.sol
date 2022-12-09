// SPDX-License-Identifier: BUSL-1.1

pragma solidity 0.8.9;

// unused import; required for a forced contract compilation
import { ProxyAdmin } from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";

import { TransparentUpgradeableProxy } from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

contract GenesisTUP is TransparentUpgradeableProxy {

    constructor() TransparentUpgradeableProxy(address(0), address(0), "") {}

    // @todo initializer check
    function init(address admin_) public {
        _changeAdmin(admin_);
    }
}
