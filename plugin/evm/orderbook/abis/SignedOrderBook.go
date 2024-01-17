package abis

var SignedOrderBookAbi = []byte(`{"abi": [
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_defaultOrderBook",
          "type": "address"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "constructor"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "uint8",
          "name": "version",
          "type": "uint8"
        }
      ],
      "name": "Initialized",
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
          "indexed": false,
          "internalType": "address",
          "name": "account",
          "type": "address"
        }
      ],
      "name": "Paused",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "address",
          "name": "account",
          "type": "address"
        }
      ],
      "name": "Unpaused",
      "type": "event"
    },
    {
      "inputs": [],
      "name": "ORDER_TYPEHASH",
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
            },
            {
              "internalType": "bool",
              "name": "postOnly",
              "type": "bool"
            }
          ],
          "internalType": "struct ISignedOrderBook.Order[]",
          "name": "orders",
          "type": "tuple[]"
        },
        {
          "internalType": "bytes[]",
          "name": "signatures",
          "type": "bytes[]"
        }
      ],
      "name": "cancelOrders",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "defaultOrderBook",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
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
            },
            {
              "internalType": "bool",
              "name": "postOnly",
              "type": "bool"
            }
          ],
          "internalType": "struct ISignedOrderBook.Order",
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
      "inputs": [],
      "name": "governance",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_governance",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "_juror",
          "type": "address"
        }
      ],
      "name": "initialize",
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
          "internalType": "address",
          "name": "sender",
          "type": "address"
        }
      ],
      "name": "isTradingAuthority",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "juror",
      "outputs": [
        {
          "internalType": "contract IJuror",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "",
          "type": "bytes32"
        }
      ],
      "name": "orderInfo",
      "outputs": [
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
          "internalType": "struct ISignedOrderBook.OrderInfo",
          "name": "",
          "type": "tuple"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "paused",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "__governance",
          "type": "address"
        }
      ],
      "name": "setGovernace",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes",
          "name": "encodedOrder",
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
            },
            {
              "internalType": "bool",
              "name": "postOnly",
              "type": "bool"
            }
          ],
          "internalType": "struct ISignedOrderBook.Order",
          "name": "order",
          "type": "tuple"
        },
        {
          "internalType": "bytes",
          "name": "signature",
          "type": "bytes"
        }
      ],
      "name": "verifySigner",
      "outputs": [
        {
          "internalType": "address",
          "name": "signer",
          "type": "address"
        },
        {
          "internalType": "bytes32",
          "name": "orderHash",
          "type": "bytes32"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    }
  ]}`)
