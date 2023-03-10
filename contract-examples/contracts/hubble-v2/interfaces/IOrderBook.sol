// SPDX-License-Identifier: BUSL-1.1

pragma solidity 0.8.9;

interface IOrderBook {
    struct Order {
        uint256 ammIndex;
        address trader;
        int256 baseAssetQuantity;
        uint256 price;
        uint256 salt;
    }

    enum OrderStatus {
        Invalid,
        Placed,
        Filled,
        Cancelled
    }

    event OrderPlaced(address indexed trader, Order order, bytes signature);
    event OrderCancelled(address indexed trader, Order order);
    event OrdersMatched(Order[2] orders, bytes[2] signatures, uint256 fillAmount, address relayer);
    event LiquidationOrderMatched(address indexed trader, Order order, bytes signature, uint256 fillAmount, address relayer);

    function executeMatchedOrders(Order[2] memory orders, bytes[2] memory signatures, int256 fillAmount) external;
    function settleFunding() external;
    function getLastTradePrices() external view returns(uint[] memory lastTradePrices);
    function liquidateAndExecuteOrder(address trader, Order memory order, bytes memory signature, int toLiquidate) external;
}
