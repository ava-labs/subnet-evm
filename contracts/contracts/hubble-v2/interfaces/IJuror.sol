// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { IOrderHandler } from "./IOrderHandler.sol";

interface IJuror {
    enum BadElement { Order0, Order1, Generic }

    // Order Matching
    function validateOrdersAndDetermineFillPrice(
        bytes[2] calldata data,
        int256 fillAmount
    )   external
        view
        returns(string memory err, BadElement reason, IOrderHandler.MatchingValidationRes memory res);

    function validateLiquidationOrderAndDetermineFillPrice(bytes calldata data, uint256 liquidationAmount)
        external
        view
        returns(string memory err, IOrderHandler.LiquidationMatchingValidationRes memory res);

    // Limit Orders
    function validatePlaceLimitOrder(ILimitOrderBook.Order calldata order, address sender)
        external
        view
        returns (string memory err, bytes32 orderhash, IOrderHandler.PlaceOrderRes memory res);

    function validateCancelLimitOrder(ILimitOrderBook.Order memory order, address sender, bool assertLowMargin)
        external
        view
        returns (string memory err, bytes32 orderHash, IOrderHandler.CancelOrderRes memory res);

    // IOC Orders
    function validatePlaceIOCOrder(IImmediateOrCancelOrders.Order memory order, address sender) external view returns(string memory err, bytes32 orderHash);

    // other methods
    function getNotionalPositionAndMargin(address trader, bool includeFundingPayments, uint8 mode) external view returns(uint256 notionalPosition, int256 margin);
}

interface ILimitOrderBook {
    struct Order {
        uint256 ammIndex;
        address trader;
        int256 baseAssetQuantity;
        uint256 price;
        uint256 salt;
        bool reduceOnly;
        bool postOnly;
    }
}

interface IImmediateOrCancelOrders {
    struct Order {
        uint8 orderType;
        uint256 expireAt;
        uint256 ammIndex;
        address trader;
        int256 baseAssetQuantity;
        uint256 price;
        uint256 salt;
        bool reduceOnly;
    }
}
