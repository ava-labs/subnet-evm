// SPDX-License-Identifier: BUSL-1.1

pragma solidity 0.8.9;

import { IERC20 } from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

interface IMarginAccount {
    struct Collateral {
        IERC20 token;
        uint weight;
        uint8 decimals;
    }

    enum LiquidationStatus {
        IS_LIQUIDATABLE,
        OPEN_POSITIONS,
        NO_DEBT,
        ABOVE_THRESHOLD
    }

    function addMargin(uint idx, uint amount) external;
    function addMarginFor(uint idx, uint amount, address to) external;
    function removeMargin(uint idx, uint256 amount) external;
    function getSpotCollateralValue(address trader) external view returns(int256 spot);
    function weightedAndSpotCollateral(address trader) external view returns(int256, int256);
    function getNormalizedMargin(address trader) external view returns(int256);
    function realizePnL(address trader, int256 realizedPnl) external;
    function isLiquidatable(address trader, bool includeFunding) external view returns(LiquidationStatus, uint, uint);
    function supportedAssetsLen() external view returns(uint);
    function supportedAssets() external view returns (Collateral[] memory);
    function margin(uint idx, address trader) external view returns(int256);
    function transferOutVusd(address recipient, uint amount) external;
    function liquidateExactRepay(address trader, uint repay, uint idx, uint minSeizeAmount) external;
    function oracle() external view returns(address); // interface in the protocol repo returns IOracle
    function removeMarginFor(address trader, uint idx, uint256 amount) external;
}
