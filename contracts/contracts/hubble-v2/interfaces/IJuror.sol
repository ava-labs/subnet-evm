// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { IClearingHouse } from "./IClearingHouse.sol";

interface IJuror {
    function validateOrdersAndDetermineFillPrice(
        bytes[2] calldata data,
        int256 fillAmount
    )   external
        view
        returns(
            IClearingHouse.Instruction[2] memory instructions,
            uint8[2] memory orderTypes,
            bytes[2] memory encodedOrders,
            uint256 fillPrice
        );

    function validateLiquidationOrderAndDetermineFillPrice(bytes calldata data, uint256 liquidationAmount)
        external
        view
        returns(
            IClearingHouse.Instruction memory instruction,
            uint8 orderType,
            bytes memory encodedOrder,
            uint256 fillPrice,
            int256 fillAmount
        );


    // Limit Orders
    function validatePlaceLimitOrder(ILimitOrderBook.OrderV2 calldata order, address trader)
        external
        view
        returns (string memory errs, bytes32 orderhash, IOrderHandler.PlaceOrderRes memory res);

    function validateCancelLimitOrder(ILimitOrderBook.OrderV2 memory order, address trader, bool assertLowMargin)
        external
        view
        returns (string memory err, bytes32 orderHash, IOrderHandler.CancelOrderRes memory res);

    // IOC Orders
    function validatePlaceIOCOrders(IImmediateOrCancelOrders.Order[] memory orders, address sender) external view returns(bytes32[] memory orderHashes);

    // ticks
    function getPrevTick(address amm, bool isBid, uint tick) external view returns (uint prevTick);
    function sampleImpactBid(address amm) external view returns (uint impactBid);
    function sampleImpactAsk(address amm) external view returns (uint impactAsk);
    function getQuote(address amm, int256 baseAssetQuantity) external view returns (uint256 rate);
    function getBaseQuote(address amm, int256 quoteQuantity) external view returns (uint256 rate);
}

interface ILimitOrderBook {
    struct OrderV2 {
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

interface IOrderHandler {
    enum OrderStatus {
        Invalid,
        Placed,
        Filled,
        Cancelled
    }

    struct PlaceOrderRes {
        uint reserveAmount;
        address amm;
    }

    struct CancelOrderRes {
        int unfilledAmount;
        address amm;
    }
}
