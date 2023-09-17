// SPDX-License-Identifier: BUSL-1.1

pragma solidity ^0.8.0;

interface IOrderBook {
    enum OrderStatus {
        Invalid,
        Placed,
        Filled,
        Cancelled
    }

    enum OrderExecutionMode {
        Taker,
        Maker,
        SameBlock,
        Liquidation
    }

    struct Order {
        uint256 ammIndex;
        address trader;
        int256 baseAssetQuantity;
        uint256 price;
        uint256 salt;
        bool reduceOnly;
        bool postOnly;
    }

    struct MatchInfo {
        bytes32 orderHash;
        uint blockPlaced;
        OrderExecutionMode mode;
    }

    event OrderAccepted(address indexed trader, bytes32 indexed orderHash, Order order, uint timestamp);
    event OrderCancelled(address indexed trader, bytes32 indexed orderHash, uint timestamp);
    event OrdersMatched(bytes32 indexed orderHash0, bytes32 indexed orderHash1, uint256 fillAmount, uint price, uint openInterestNotional, address relayer, uint timestamp);
    event LiquidationOrderMatched(address indexed trader, bytes32 indexed orderHash, bytes signature, uint256 fillAmount, uint price, uint openInterestNotional, address relayer, uint timestamp);
    event OrderMatchingError(bytes32 indexed orderHash, string err);
    event LiquidationError(address indexed trader, bytes32 indexed orderHash, string err, uint256 toLiquidate);

    function executeMatchedOrders(Order[2] memory orders, int256 fillAmount) external;
    function settleFunding() external;
    function liquidateAndExecuteOrder(address trader, Order memory order, bytes memory signature, uint256 toLiquidate) external;
    function getLastTradePrices() external view returns(uint[] memory lastTradePrices);
}
