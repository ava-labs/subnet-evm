// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IHubbleBibliophile {
    struct Order {
        uint256 ammIndex;
        address trader;
        int256 baseAssetQuantity;
        uint256 price;
        uint256 salt;
        bool reduceOnly;
    }

    enum OrderExecutionMode {
        Taker,
        Maker,
        SameBlock,
        Liquidation
    }

    function getNotionalPositionAndMargin(address trader, bool includeFundingPayments, uint8 mode)
        external
        view
        returns(uint256 notionalPosition, int256 margin);

    function getPositionSizes(address trader) external view returns(int[] memory posSizes);

    function validateOrdersAndDetermineFillPrice(
        Order[2] memory orders,
        bytes32[2] memory orderHashes,
        int256 fillAmount
    ) external view returns(uint256 fillPrice, OrderExecutionMode mode0, OrderExecutionMode mode1);

    function validateLiquidationOrderAndDetermineFillPrice(
        Order memory order,
        int256 fillAmount
    ) external view returns(uint256 fillPrice);
}
