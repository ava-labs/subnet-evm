// SPDX-License-Identifier: BUSL-1.1

pragma solidity 0.8.9;

import { IOrderBook } from "./IOrderBook.sol";

interface IClearingHouse {
    enum Mode { Maintenance_Margin, Min_Allowable_Margin }
    function openPosition(IOrderBook.Order memory order, int256 fillAmount, uint256 fulfillPrice, bool isMakerOrder) external;
    function settleFunding() external;
    function getTotalNotionalPositionAndUnrealizedPnl(address trader, int256 margin, Mode mode)
        external
        view
        returns(uint256 notionalPosition, int256 unrealizedPnl);
    function isAboveMaintenanceMargin(address trader) external view returns(bool);
    function assertMarginRequirement(address trader) external view;
    function updatePositions(address trader) external;
    function getMarginFraction(address trader) external view returns(int256);
    function getTotalFunding(address trader) external view returns(int256 totalFunding);
    function getAmmsLength() external view returns(uint);
    function amms(uint idx) external view returns(address); // is returns(IAMM) in protocol repo IClearingHouse
    function maintenanceMargin() external view returns(int256);
    function minAllowableMargin() external view returns(int256);
    function takerFee() external view returns(uint256);
    function makerFee() external view returns(uint256);
    function liquidationPenalty() external view returns(uint256);
    function getNotionalPositionAndMargin(address trader, bool includeFundingPayments, Mode mode)
        external
        view
        returns(uint256 notionalPosition, int256 margin);
    function liquidate(address trader, uint ammIdx, uint price, int toLiquidate) external;
    function feeSink() external view returns(address);
    function calcMarginFraction(address trader, bool includeFundingPayments, Mode mode) external view returns(int256);
}
