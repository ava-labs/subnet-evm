// SPDX-License-Identifier: BUSL-1.1

pragma solidity 0.8.9;

import { ECDSA } from "../node_modules/@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import { EIP712 } from "../node_modules/@openzeppelin/contracts/utils/cryptography/draft-EIP712.sol";

contract OrderBook is EIP712 {
    struct Order {
        address trader;
        int256 baseAssetQuantity;
        uint256 price;
        uint256 salt;
    }

    enum OrderStatus {
        Unfilled,
        Filled,
        Cancelled
    }

    event OrderPlaced(address indexed trader, int256 baseAssetQuantity, uint256 price, address relayer);

    mapping(bytes32 => OrderStatus) public ordersStatus;
    Order[] public orders;

    // keccak256("Order(address trader,int256 baseAssetQuantity,uint256 price,uint256 salt)");
    bytes32 public constant LIMIT_ORDER_TYPEHASH = 0x4cab2d4fcf58d07df65ee3d9d1e6e3c407eae39d76ee15b247a025ab52e2c45d;

    constructor(string memory name, string memory version) EIP712(name, version) {}

    function placeOrder(Order memory order, bytes memory signature) external {
        (, bytes32 orderHash) = verifySigner(order, signature);

        // OB_OMBU: Order Must Be Unfilled
        require(ordersStatus[orderHash] == OrderStatus.Unfilled, "OB_OMBU");

        orders.push(order);
        emit OrderPlaced(order.trader, order.baseAssetQuantity, order.price, msg.sender);
    }

    function verifySigner(Order memory order, bytes memory signature) public view returns (address, bytes32) {
        bytes32 orderHash = getOrderHash(order);
        address signer = ECDSA.recover(orderHash, signature);

        // OB_SINT: Signer Is Not Trader
        require(signer == order.trader, "OB_SINT");

        return (signer, orderHash);
    }

    function getOrderHash(Order memory order) public view returns (bytes32) {
        return _hashTypedDataV4(keccak256(abi.encode(LIMIT_ORDER_TYPEHASH, order)));
    }

    function getAllOrders() external view returns (Order[] memory) {
        return orders;
    }

    function getOrdersLen() external view returns (uint) {
        return orders.length;
    }
}
