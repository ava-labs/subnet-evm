package limitorders

var orderBookAbi = []byte(`{"abi": [
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
            }
          ],
          "indexed": false,
          "internalType": "struct IOrderBook.Order",
          "name": "order",
          "type": "tuple"
        },
        {
          "indexed": false,
          "internalType": "bytes",
          "name": "signature",
          "type": "bytes"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "fillAmount",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "address",
          "name": "relayer",
          "type": "address"
        }
      ],
      "name": "LiquidationOrderMatched",
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
            }
          ],
          "indexed": false,
          "internalType": "struct IOrderBook.Order",
          "name": "order",
          "type": "tuple"
        }
      ],
      "name": "OrderCancelled",
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
            }
          ],
          "indexed": false,
          "internalType": "struct IOrderBook.Order",
          "name": "order",
          "type": "tuple"
        },
        {
          "indexed": false,
          "internalType": "bytes",
          "name": "signature",
          "type": "bytes"
        }
      ],
      "name": "OrderPlaced",
      "type": "event"
    },
    {
      "anonymous": false,
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
            }
          ],
          "indexed": false,
          "internalType": "struct IOrderBook.Order[2]",
          "name": "orders",
          "type": "tuple[2]"
        },
        {
          "indexed": false,
          "internalType": "bytes[2]",
          "name": "signatures",
          "type": "bytes[2]"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "fillAmount",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "address",
          "name": "relayer",
          "type": "address"
        }
      ],
      "name": "OrdersMatched",
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
            }
          ],
          "internalType": "struct IOrderBook.Order[2]",
          "name": "orders",
          "type": "tuple[2]"
        },
        {
          "internalType": "bytes[2]",
          "name": "signatures",
          "type": "bytes[2]"
        },
        {
          "internalType": "int256",
          "name": "fillAmount",
          "type": "int256"
        }
      ],
      "name": "executeMatchedOrders",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getLastTradePrices",
      "outputs": [
        {
          "internalType": "uint256[]",
          "name": "lastTradePrices",
          "type": "uint256[]"
        }
      ],
      "stateMutability": "view",
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
            }
          ],
          "internalType": "struct IOrderBook.Order",
          "name": "order",
          "type": "tuple"
        },
        {
          "internalType": "bytes",
          "name": "signature",
          "type": "bytes"
        },
        {
          "internalType": "int256",
          "name": "toLiquidate",
          "type": "int256"
        }
      ],
      "name": "liquidateAndExecuteOrder",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "settleFunding",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]}`)

var marginAccountAbi = []byte(`{"abi": [
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "idx",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "addMargin",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "idx",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        },
        {
          "internalType": "address",
          "name": "to",
          "type": "address"
        }
      ],
      "name": "addMarginFor",
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
        }
      ],
      "name": "getNormalizedMargin",
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
          "internalType": "address",
          "name": "trader",
          "type": "address"
        }
      ],
      "name": "getSpotCollateralValue",
      "outputs": [
        {
          "internalType": "int256",
          "name": "spot",
          "type": "int256"
        }
      ],
      "stateMutability": "view",
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
          "internalType": "bool",
          "name": "includeFunding",
          "type": "bool"
        }
      ],
      "name": "isLiquidatable",
      "outputs": [
        {
          "internalType": "enum IMarginAccount.LiquidationStatus",
          "name": "",
          "type": "uint8"
        },
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        },
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
          "internalType": "address",
          "name": "trader",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "repay",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "idx",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "minSeizeAmount",
          "type": "uint256"
        }
      ],
      "name": "liquidateExactRepay",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "idx",
          "type": "uint256"
        },
        {
          "internalType": "address",
          "name": "trader",
          "type": "address"
        }
      ],
      "name": "margin",
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
      "inputs": [],
      "name": "oracle",
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
          "name": "trader",
          "type": "address"
        },
        {
          "internalType": "int256",
          "name": "realizedPnl",
          "type": "int256"
        }
      ],
      "name": "realizePnL",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "idx",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "removeMargin",
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
          "name": "idx",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "removeMarginFor",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "supportedAssets",
      "outputs": [
        {
          "components": [
            {
              "internalType": "contract IERC20",
              "name": "token",
              "type": "address"
            },
            {
              "internalType": "uint256",
              "name": "weight",
              "type": "uint256"
            },
            {
              "internalType": "uint8",
              "name": "decimals",
              "type": "uint8"
            }
          ],
          "internalType": "struct IMarginAccount.Collateral[]",
          "name": "",
          "type": "tuple[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "supportedAssetsLen",
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
          "internalType": "address",
          "name": "recipient",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "transferOutVusd",
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
        }
      ],
      "name": "weightedAndSpotCollateral",
      "outputs": [
        {
          "internalType": "int256",
          "name": "",
          "type": "int256"
        },
        {
          "internalType": "int256",
          "name": "",
          "type": "int256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    }
  ]}`)

var clearingHouseAbi = []byte(`{"abi": [
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "idx",
          "type": "uint256"
        }
      ],
      "name": "amms",
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
          "name": "trader",
          "type": "address"
        }
      ],
      "name": "assertMarginRequirement",
      "outputs": [],
      "stateMutability": "view",
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
          "internalType": "bool",
          "name": "includeFundingPayments",
          "type": "bool"
        },
        {
          "internalType": "enum IClearingHouse.Mode",
          "name": "mode",
          "type": "uint8"
        }
      ],
      "name": "calcMarginFraction",
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
      "inputs": [],
      "name": "feeSink",
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
      "inputs": [],
      "name": "getAmmsLength",
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
          "internalType": "address",
          "name": "trader",
          "type": "address"
        }
      ],
      "name": "getMarginFraction",
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
          "internalType": "address",
          "name": "trader",
          "type": "address"
        },
        {
          "internalType": "bool",
          "name": "includeFundingPayments",
          "type": "bool"
        },
        {
          "internalType": "enum IClearingHouse.Mode",
          "name": "mode",
          "type": "uint8"
        }
      ],
      "name": "getNotionalPositionAndMargin",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "notionalPosition",
          "type": "uint256"
        },
        {
          "internalType": "int256",
          "name": "margin",
          "type": "int256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "trader",
          "type": "address"
        }
      ],
      "name": "getTotalFunding",
      "outputs": [
        {
          "internalType": "int256",
          "name": "totalFunding",
          "type": "int256"
        }
      ],
      "stateMutability": "view",
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
          "internalType": "int256",
          "name": "margin",
          "type": "int256"
        },
        {
          "internalType": "enum IClearingHouse.Mode",
          "name": "mode",
          "type": "uint8"
        }
      ],
      "name": "getTotalNotionalPositionAndUnrealizedPnl",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "notionalPosition",
          "type": "uint256"
        },
        {
          "internalType": "int256",
          "name": "unrealizedPnl",
          "type": "int256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "trader",
          "type": "address"
        }
      ],
      "name": "isAboveMaintenanceMargin",
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
          "name": "trader",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "ammIdx",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "price",
          "type": "uint256"
        },
        {
          "internalType": "int256",
          "name": "toLiquidate",
          "type": "int256"
        }
      ],
      "name": "liquidate",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "liquidationPenalty",
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
      "inputs": [],
      "name": "maintenanceMargin",
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
      "inputs": [],
      "name": "makerFee",
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
      "inputs": [],
      "name": "minAllowableMargin",
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
            }
          ],
          "internalType": "struct IOrderBook.Order",
          "name": "order",
          "type": "tuple"
        },
        {
          "internalType": "int256",
          "name": "fillAmount",
          "type": "int256"
        },
        {
          "internalType": "uint256",
          "name": "fulfillPrice",
          "type": "uint256"
        },
        {
          "internalType": "bool",
          "name": "isMakerOrder",
          "type": "bool"
        }
      ],
      "name": "openPosition",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "settleFunding",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "takerFee",
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
          "internalType": "address",
          "name": "trader",
          "type": "address"
        }
      ],
      "name": "updatePositions",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]}`)
