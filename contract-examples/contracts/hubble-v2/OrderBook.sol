// SPDX-License-Identifier: BUSL-1.1

pragma solidity 0.8.9;

import { SafeCast } from "@openzeppelin/contracts/utils/math/SafeCast.sol";
import { ECDSAUpgradeable } from "@openzeppelin/contracts-upgradeable/utils/cryptography/ECDSAUpgradeable.sol";
import { EIP712Upgradeable } from "@openzeppelin/contracts-upgradeable/utils/cryptography/draft-EIP712Upgradeable.sol";

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

    function executeMatchedOrders(Order[2] memory orders, bytes[2] memory signatures, int256 fillAmount) external;
    function executeFundingPayment() external;
    function getLastTradePrices() external view returns(uint[] memory lastTradePrices);
}

contract OrderBook is IOrderBook, EIP712Upgradeable {
    using SafeCast for uint256;
    using SafeCast for int256;

    // keccak256("Order(uint256 ammIndex,address trader,int256 baseAssetQuantity,uint256 price,uint256 salt)");
    bytes32 public constant ORDER_TYPEHASH = 0xba5bdc08c77846c2444ea7c84fcaf3479e3389b274ebc7ab59358538ca00dbe0;

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

    event OrderPlaced(address indexed trader, Order order, bytes signature);
    event OrderCancelled(address indexed trader, Order order);
    event OrdersMatched(Order[2] orders, bytes[2] signatures, uint256 fillAmount, address relayer);

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
     * @param signatures To verify authenticity of the order
     * @param fillAmount Amount to be filled for each order. This is to support partial fills.
     *        Should be > 0 and min(unfilled amount in both orders)
    */
    function executeMatchedOrders(
        Order[2] memory orders,
        bytes[2] memory signatures,
        int256 fillAmount
    )   external
        /* onlyValidator */
    {
        // Checks and Effects
        require(orders[0].baseAssetQuantity > 0, "OB_order_0_is_not_long");
        require(orders[1].baseAssetQuantity < 0, "OB_order_1_is_not_short");
        require(fillAmount > 0, "OB_fillAmount_is_neg");
        require(orders[0].price /* buy */ >= orders[1].price /* sell */, "OB_orders_do_not_match");
        (bytes32 orderHash0, uint blockPlaced0) = _verifyOrder(orders[0], signatures[0], fillAmount);
        (bytes32 orderHash1, uint blockPlaced1) = _verifyOrder(orders[1], signatures[1], -fillAmount);
        // @todo min fillAmount and min order.baseAsset check

        // Effects
        _updateOrder(orderHash0, fillAmount, orders[0].baseAssetQuantity);
        _updateOrder(orderHash1, -fillAmount, orders[1].baseAssetQuantity);

        // Interactions
        uint fulfillPrice = orders[0].price; // if prices are equal or long blockPlaced <= short blockPlaced
        if (orders[0].price != orders[1].price && blockPlaced0 > blockPlaced1) {
            fulfillPrice = orders[1].price;
        }

        _openPosition(orders[0], fillAmount, fulfillPrice);
        _openPosition(orders[1], -fillAmount, fulfillPrice);

        emit OrdersMatched(orders, signatures, fillAmount.toUint256(), msg.sender);
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

    function placeOrder(Order memory order, bytes memory signature) external {
        require(msg.sender == order.trader, "OB_sender_is_not_trader");
        // verifying signature here to avoid too many fake placeOrders
        (, bytes32 orderHash) = verifySigner(order, signature);
        // order should not exist in the orderStatus map already
        require(orderInfo[orderHash].status == OrderStatus.Invalid, "OB_Order_already_exists");
        orderInfo[orderHash] = OrderInfo(block.number, 0, OrderStatus.Placed);
        // @todo assert margin requirements for placing the order
        // @todo min size requirement while placing order

        emit OrderPlaced(order.trader, order, signature);
    }

    function cancelOrder(Order memory order) external {
        require(msg.sender == order.trader, "OB_sender_is_not_trader");
        bytes32 orderHash = getOrderHash(order);
        // order status should be placed
        require(orderInfo[orderHash].status == OrderStatus.Placed, "OB_Order_does_not_exist");
        orderInfo[orderHash].status = OrderStatus.Cancelled;

        emit OrderCancelled(order.trader, order);
    }

    function executeFundingPayment() external {}

    /**
    @dev assuming one order is in liquidation zone and other is out of it
    @notice liquidate trader
    @param trader trader to liquidate
    @param order order to match when liuidating for a particular amm
    @param signature signature corresponding to order
    @param toLiquidate baseAsset amount being traded/liquidated. -ve if short position is being liquidated, +ve if long
    */
    function liquidateAndExecuteOrder(address trader, Order memory order, bytes memory signature, int toLiquidate) external {
        // liquidate
        positions[order.ammIndex][trader].openNotional -= (order.price * abs(toLiquidate).toUint256() / 1e18);
        positions[order.ammIndex][trader].size -= toLiquidate;

        (bytes32 orderHash,) = _verifyOrder(order, signature, toLiquidate);
        _updateOrder(orderHash, toLiquidate, order.baseAssetQuantity);
        _openPosition(order, toLiquidate, order.price);
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
        address signer = ECDSAUpgradeable.recover(orderHash, signature);

        // OB_SINT: Signer Is Not Trader
        require(signer == order.trader, "OB_SINT");

        return (signer, orderHash);
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
