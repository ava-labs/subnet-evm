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

    // IOC Orders
    function validatePlaceIOCOrders(IImmediateOrCancelOrders.Order[] memory orders, address sender) external view returns(bytes32[] memory orderHashes);
}

interface IImmediateOrCancelOrders {
    struct Order {
        uint8 orderType;
        uint256 ammIndex;
        address trader;
        int256 baseAssetQuantity;
        uint256 price;
        uint256 salt;
        bool reduceOnly;
    }
}
