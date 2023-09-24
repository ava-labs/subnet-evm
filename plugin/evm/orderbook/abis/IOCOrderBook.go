package abis

var IOCOrderBookAbi = []byte(`{"abi": [
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "trader",
        "type": "address"
      },
      {
        "indexed": true,
        "internalType": "bytes32",
        "name": "orderHash",
        "type": "bytes32"
      },
      {
        "components": [
          {
            "internalType": "uint8",
            "name": "orderType",
            "type": "uint8"
          },
          {
            "internalType": "uint256",
            "name": "expireAt",
            "type": "uint256"
          },
          {
            "internalType": "uint256",
            "name": "ammIndex",
            "type": "uint256"
          },
          {
            "internalType": "address",
            "name": "trader",
            "type": "address"
          },
          {
            "internalType": "int256",
            "name": "baseAssetQuantity",
            "type": "int256"
          },
          {
            "internalType": "uint256",
            "name": "price",
            "type": "uint256"
          },
          {
            "internalType": "uint256",
            "name": "salt",
            "type": "uint256"
          },
          {
            "internalType": "bool",
            "name": "reduceOnly",
            "type": "bool"
          }
        ],
        "indexed": false,
        "internalType": "struct IImmediateOrCancelOrders.Order",
        "name": "order",
        "type": "tuple"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "timestamp",
        "type": "uint256"
      }
    ],
    "name": "OrderAccepted",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "trader",
        "type": "address"
      },
      {
        "indexed": true,
        "internalType": "bytes32",
        "name": "orderHash",
        "type": "bytes32"
      },
      {
        "components": [
          {
            "internalType": "uint8",
            "name": "orderType",
            "type": "uint8"
          },
          {
            "internalType": "uint256",
            "name": "expireAt",
            "type": "uint256"
          },
          {
            "internalType": "uint256",
            "name": "ammIndex",
            "type": "uint256"
          },
          {
            "internalType": "address",
            "name": "trader",
            "type": "address"
          },
          {
            "internalType": "int256",
            "name": "baseAssetQuantity",
            "type": "int256"
          },
          {
            "internalType": "uint256",
            "name": "price",
            "type": "uint256"
          },
          {
            "internalType": "uint256",
            "name": "salt",
            "type": "uint256"
          },
          {
            "internalType": "bool",
            "name": "reduceOnly",
            "type": "bool"
          }
        ],
        "indexed": false,
        "internalType": "struct IImmediateOrCancelOrders.Order",
        "name": "order",
        "type": "tuple"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "timestamp",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "string",
        "name": "err",
        "type": "string"
      }
    ],
    "name": "OrderRejected",
    "type": "event"
  },
  {
    "inputs": [],
    "name": "expirationCap",
    "outputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "components": [
          {
            "internalType": "uint8",
            "name": "orderType",
            "type": "uint8"
          },
          {
            "internalType": "uint256",
            "name": "expireAt",
            "type": "uint256"
          },
          {
            "internalType": "uint256",
            "name": "ammIndex",
            "type": "uint256"
          },
          {
            "internalType": "address",
            "name": "trader",
            "type": "address"
          },
          {
            "internalType": "int256",
            "name": "baseAssetQuantity",
            "type": "int256"
          },
          {
            "internalType": "uint256",
            "name": "price",
            "type": "uint256"
          },
          {
            "internalType": "uint256",
            "name": "salt",
            "type": "uint256"
          },
          {
            "internalType": "bool",
            "name": "reduceOnly",
            "type": "bool"
          }
        ],
        "internalType": "struct IImmediateOrCancelOrders.Order",
        "name": "order",
        "type": "tuple"
      }
    ],
    "name": "getOrderHash",
    "outputs": [
      {
        "internalType": "bytes32",
        "name": "",
        "type": "bytes32"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "bytes32",
        "name": "orderHash",
        "type": "bytes32"
      }
    ],
    "name": "orderStatus",
    "outputs": [
      {
        "components": [
          {
            "internalType": "uint256",
            "name": "blockPlaced",
            "type": "uint256"
          },
          {
            "internalType": "int256",
            "name": "filledAmount",
            "type": "int256"
          },
          {
            "internalType": "enum IOrderHandler.OrderStatus",
            "name": "status",
            "type": "uint8"
          }
        ],
        "internalType": "struct IImmediateOrCancelOrders.OrderInfo",
        "name": "",
        "type": "tuple"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "components": [
          {
            "internalType": "uint8",
            "name": "orderType",
            "type": "uint8"
          },
          {
            "internalType": "uint256",
            "name": "expireAt",
            "type": "uint256"
          },
          {
            "internalType": "uint256",
            "name": "ammIndex",
            "type": "uint256"
          },
          {
            "internalType": "address",
            "name": "trader",
            "type": "address"
          },
          {
            "internalType": "int256",
            "name": "baseAssetQuantity",
            "type": "int256"
          },
          {
            "internalType": "uint256",
            "name": "price",
            "type": "uint256"
          },
          {
            "internalType": "uint256",
            "name": "salt",
            "type": "uint256"
          },
          {
            "internalType": "bool",
            "name": "reduceOnly",
            "type": "bool"
          }
        ],
        "internalType": "struct IImmediateOrCancelOrders.Order[]",
        "name": "orders",
        "type": "tuple[]"
      }
    ],
    "name": "placeOrders",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "bytes",
        "name": "data",
        "type": "bytes"
      },
      {
        "internalType": "bytes",
        "name": "metadata",
        "type": "bytes"
      }
    ],
    "name": "updateOrder",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  }
]}`)
