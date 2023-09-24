package abis

var LimitOrderBookAbi = []byte(`{"abi": [
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
            },
            {
              "internalType": "bool",
              "name": "postOnly",
              "type": "bool"
            }
          ],
          "indexed": false,
          "internalType": "struct ILimitOrderBook.Order",
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
          "indexed": false,
          "internalType": "uint256",
          "name": "timestamp",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "bool",
          "name": "isAutoCancelled",
          "type": "bool"
        }
      ],
      "name": "OrderCancelAccepted",
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
      "name": "OrderCancelRejected",
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
            },
            {
              "internalType": "bool",
              "name": "postOnly",
              "type": "bool"
            }
          ],
          "indexed": false,
          "internalType": "struct ILimitOrderBook.Order",
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
      "inputs": [
        {
          "components": [
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
            },
            {
              "internalType": "bool",
              "name": "postOnly",
              "type": "bool"
            }
          ],
          "internalType": "struct ILimitOrderBook.Order[]",
          "name": "orders",
          "type": "tuple[]"
        }
      ],
      "name": "cancelOrders",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "components": [
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
            },
            {
              "internalType": "bool",
              "name": "postOnly",
              "type": "bool"
            }
          ],
          "internalType": "struct ILimitOrderBook.Order[]",
          "name": "orders",
          "type": "tuple[]"
        }
      ],
      "name": "cancelOrdersWithLowMargin",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "components": [
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
            },
            {
              "internalType": "bool",
              "name": "postOnly",
              "type": "bool"
            }
          ],
          "internalType": "struct ILimitOrderBook.Order",
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
              "internalType": "uint256",
              "name": "reservedMargin",
              "type": "uint256"
            },
            {
              "internalType": "enum IOrderHandler.OrderStatus",
              "name": "status",
              "type": "uint8"
            }
          ],
          "internalType": "struct ILimitOrderBook.OrderInfo",
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
            },
            {
              "internalType": "bool",
              "name": "postOnly",
              "type": "bool"
            }
          ],
          "internalType": "struct ILimitOrderBook.Order[]",
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
          "internalType": "address",
          "name": "trader",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "ammIndex",
          "type": "uint256"
        }
      ],
      "name": "reduceOnlyAmount",
      "outputs": [
        {
          "internalType": "int256",
          "name": "",
          "type": "int256"
        }
      ],
      "stateMutability": "view",
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
