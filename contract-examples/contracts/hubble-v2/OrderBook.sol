// SPDX-License-Identifier: BUSL-1.1

pragma solidity 0.8.9;

import { SafeCast } from "@openzeppelin/contracts/utils/math/SafeCast.sol";
import { ECDSAUpgradeable } from "@openzeppelin/contracts-upgradeable/utils/cryptography/ECDSAUpgradeable.sol";
import { EIP712Upgradeable } from "@openzeppelin/contracts-upgradeable/utils/cryptography/draft-EIP712Upgradeable.sol";

import { IOrderBook } from "./interfaces/IOrderBook.sol";

contract OrderBook is IOrderBook, EIP712Upgradeable {
    using SafeCast for uint256;
    using SafeCast for int256;

    // keccak256("Order(uint256 ammIndex,address trader,int256 baseAssetQuantity,uint256 price,uint256 salt,bool reduceOnly)");
    bytes32 public constant ORDER_TYPEHASH = 0x0a2e4d36552888a97d5a8975ad22b04e90efe5ea0a8abb97691b63b431eb25d2;

    struct OrderInfo {
        uint blockPlaced;
        int256 filledAmount;
        OrderStatus status;
    }
    mapping(bytes32 => OrderInfo) public orderInfo;

    struct Position {
        int256 size;
        uint256 openNotional;
    }

    // following vars are used to mock clearingHouse
    // ammIndex => address => Position
    mapping(uint => mapping(address => Position)) public positions;
    mapping(uint => uint) public lastPrices;
    uint public numAmms;

    function initialize(string memory name, string memory version) initializer public {
        __EIP712_init(name, version);
        setNumAMMs(1);
    }

    /**
     * Execute matched orders
     * @param orders It is required that orders[0] is a LONG and orders[1] is a short
     * @param fillAmount Amount to be filled for each order. This is to support partial fills.
     *        Should be > 0 and min(unfilled amount in both orders)
    */
    function executeMatchedOrders(
        Order[2] memory orders,
        int256 fillAmount
    )   external
    {
        // Checks and Effects
        require(orders[0].baseAssetQuantity > 0, "OB_order_0_is_not_long");
        require(orders[1].baseAssetQuantity < 0, "OB_order_1_is_not_short");
        require(fillAmount > 0, "OB_fillAmount_is_neg");
        require(orders[0].price /* buy */ >= orders[1].price /* sell */, "OB_orders_do_not_match");

        bytes32 orderHash0 = getOrderHash(orders[0]);
        bytes32 orderHash1 = getOrderHash(orders[1]);
        // // Effects
        _updateOrder(orderHash0, fillAmount, orders[0].baseAssetQuantity);
        _updateOrder(orderHash1, -fillAmount, orders[1].baseAssetQuantity);

        // // Interactions
        uint fulfillPrice = orders[0].price;
        _openPosition(orders[0], fillAmount, fulfillPrice);
        _openPosition(orders[1], -fillAmount, fulfillPrice);

        emit OrdersMatched(orderHash0, orderHash1, fillAmount.toUint256(), fulfillPrice, fillAmount.toUint256() * fulfillPrice, msg.sender, block.timestamp);
    }

    /**
     * @dev mocked version of clearingHouse.openPosition
    */
    function _openPosition(Order memory order, int fillAmount, uint fulfillPrice) internal {
        // update open notional
        uint delta = abs(fillAmount).toUint256() * fulfillPrice / 1e18;
        address trader = order.trader;
        uint ammIndex = order.ammIndex;
        require(ammIndex < numAmms, "OB_please_whitelist_new_amm");
        if (positions[ammIndex][trader].size * fillAmount >= 0) { // increase position
            positions[ammIndex][trader].openNotional += delta;
        } else { // reduce position
            if (positions[ammIndex][trader].openNotional >= delta) {
                positions[ammIndex][trader].openNotional -= delta; // position reduced
            } else { // open reverse position
                positions[ammIndex][trader].openNotional = (delta - positions[ammIndex][trader].openNotional);
            }
        }
        // update position size
        positions[ammIndex][trader].size += fillAmount;
        // update latest price
        lastPrices[ammIndex] = fulfillPrice;
    }

    function placeOrder(Order memory order) external {
        bytes32 orderHash = getOrderHash(order);
        // order should not exist in the orderStatus map already
        // require(orderInfo[orderHash].status == OrderStatus.Invalid, "OB_Order_already_exists");
        orderInfo[orderHash] = OrderInfo(block.number, 0, OrderStatus.Placed);
        // @todo assert margin requirements for placing the order
        // @todo min size requirement while placing order

        emit OrderPlaced(order.trader, orderHash, order, block.timestamp);
    }

    function cancelOrder(Order memory order) external {
        require(msg.sender == order.trader, "OB_sender_is_not_trader");
        bytes32 orderHash = getOrderHash(order);
        // order status should be placed
        require(orderInfo[orderHash].status == OrderStatus.Placed, "OB_Order_does_not_exist");
        orderInfo[orderHash].status = OrderStatus.Cancelled;

        emit OrderCancelled(order.trader, orderHash, block.timestamp);
    }

    /**
     * @dev is a no-op here but works in the implementation in the protocol repo
    */
    function settleFunding() external {}

    /**
    @dev assuming one order is in liquidation zone and other is out of it
    @notice liquidate trader
    @param trader trader to liquidate
    @param order order to match when liuidating for a particular amm
    @param signature signature corresponding to order
    @param toLiquidate baseAsset amount being traded/liquidated. -ve if short position is being liquidated, +ve if long
    */
    function liquidateAndExecuteOrder(address trader, Order memory order, bytes memory signature, uint256 toLiquidate) external {
        // liquidate
        positions[order.ammIndex][trader].openNotional -= (order.price * toLiquidate / 1e18);
        positions[order.ammIndex][trader].size -= toLiquidate.toInt256();

        (bytes32 orderHash,) = _verifyOrder(order, signature, toLiquidate.toInt256());
        _updateOrder(orderHash, toLiquidate.toInt256(), order.baseAssetQuantity);
        _openPosition(order, toLiquidate.toInt256(), order.price);
        emit LiquidationOrderMatched(trader, orderHash, signature, toLiquidate, order.price, order.price * toLiquidate, msg.sender, block.timestamp);
    }

    /* ****************** */
    /*      View      */
    /* ****************** */

    function getLastTradePrices() external view returns(uint[] memory lastTradePrices) {
        lastTradePrices = new uint[](numAmms);
        for (uint i; i < numAmms; i++) {
            lastTradePrices[i] = lastPrices[i];
        }
    }

    function verifySigner(Order memory order, bytes memory signature) public view returns (address, bytes32) {
        bytes32 orderHash = getOrderHash(order);

        // removed because verification is not required
        // address signer = ECDSAUpgradeable.recover(orderHash, signature);
        // OB_SINT: Signer Is Not Trader
        // require(signer == order.trader, "OB_SINT");

        return (order.trader, orderHash);
    }

    function getOrderHash(Order memory order) public view returns (bytes32) {
        return _hashTypedDataV4(keccak256(abi.encode(ORDER_TYPEHASH, order)));
    }

    /* ****************** */
    /*      Internal      */
    /* ****************** */

    function _verifyOrder(Order memory order, bytes memory signature, int256 fillAmount)
        internal
        view
        returns (bytes32 /* orderHash */, uint /* blockPlaced */)
    {
        (, bytes32 orderHash) = verifySigner(order, signature);
        // order should be in placed status
        require(orderInfo[orderHash].status == OrderStatus.Placed, "OB_invalid_order");
        // order.baseAssetQuantity and fillAmount should have same sign
        require(order.baseAssetQuantity * fillAmount > 0, "OB_fill_and_base_sign_not_match");
        // fillAmount[orderHash] should be strictly increasing or strictly decreasing
        require(orderInfo[orderHash].filledAmount * fillAmount >= 0, "OB_invalid_fillAmount");
        require(abs(orderInfo[orderHash].filledAmount) <= abs(order.baseAssetQuantity), "OB_filled_amount_higher_than_order_base");
        return (orderHash, orderInfo[orderHash].blockPlaced);
    }

    function _updateOrder(bytes32 orderHash, int256 fillAmount, int256 baseAssetQuantity) internal {
        orderInfo[orderHash].filledAmount += fillAmount;
        // update order status if filled
        if (orderInfo[orderHash].filledAmount == baseAssetQuantity) {
            orderInfo[orderHash].status = OrderStatus.Filled;
        }
    }

    /* ****************** */
    /*        Pure        */
    /* ****************** */

    function abs(int x) internal pure returns (int) {
        return x >= 0 ? x : -x;
    }

    /* ****************** */
    /*        Mocks       */
    /* ****************** */

    /**
    * @dev only for testing with evm
    */
    function executeTestOrder(Order memory order, bytes memory signature) external {
        (bytes32 orderHash0,) = _verifyOrder(order, signature, order.baseAssetQuantity);
        _updateOrder(orderHash0, order.baseAssetQuantity, order.baseAssetQuantity);
        _openPosition(order, order.baseAssetQuantity, order.price);
    }

    function setNumAMMs(uint _num) public {
        numAmms = _num;
    }
}
