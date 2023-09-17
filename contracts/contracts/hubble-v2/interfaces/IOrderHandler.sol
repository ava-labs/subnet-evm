// SPDX-License-Identifier: BUSL-1.1

pragma solidity ^0.8.0;

import { IClearingHouse } from "./IClearingHouse.sol";

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

    struct MatchingValidationRes {
        IClearingHouse.Instruction[2] instructions;
        uint8[2] orderTypes;
        bytes[2] encodedOrders;
        uint256 fillPrice;
    }

    struct LiquidationMatchingValidationRes {
        IClearingHouse.Instruction instruction;
        uint8 orderType;
        bytes encodedOrder;
        uint256 fillPrice;
        int256 fillAmount;
    }
}
