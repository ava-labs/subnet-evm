package abis

var ClearingHouseAbi = []byte(`{"abi": [
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
        "internalType": "uint256",
        "name": "idx",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "int256",
        "name": "takerFundingPayment",
        "type": "int256"
      },
      {
        "indexed": false,
        "internalType": "int256",
        "name": "cumulativePremiumFraction",
        "type": "int256"
      }
    ],
    "name": "FundingPaid",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "uint256",
        "name": "idx",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "int256",
        "name": "premiumFraction",
        "type": "int256"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "underlyingPrice",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "int256",
        "name": "cumulativePremiumFraction",
        "type": "int256"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "nextFundingTime",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "timestamp",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "blockNumber",
        "type": "uint256"
      }
    ],
    "name": "FundingRateUpdated",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "uint256",
        "name": "idx",
        "type": "uint256"
      },
      {
        "indexed": true,
        "internalType": "address",
        "name": "amm",
        "type": "address"
      }
    ],
    "name": "MarketAdded",
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
        "internalType": "uint256",
        "name": "idx",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "int256",
        "name": "baseAsset",
        "type": "int256"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "price",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "int256",
        "name": "realizedPnl",
        "type": "int256"
      },
      {
        "indexed": false,
        "internalType": "int256",
        "name": "size",
        "type": "int256"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "openNotional",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "int256",
        "name": "fee",
        "type": "int256"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "timestamp",
        "type": "uint256"
      }
    ],
    "name": "PositionLiquidated",
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
        "internalType": "uint256",
        "name": "idx",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "int256",
        "name": "baseAsset",
        "type": "int256"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "price",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "int256",
        "name": "realizedPnl",
        "type": "int256"
      },
      {
        "indexed": false,
        "internalType": "int256",
        "name": "size",
        "type": "int256"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "openNotional",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "int256",
        "name": "fee",
        "type": "int256"
      },
      {
        "indexed": false,
        "internalType": "enum IOrderBook.OrderExecutionMode",
        "name": "mode",
        "type": "uint8"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "timestamp",
        "type": "uint256"
      }
    ],
    "name": "PositionModified",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "referrer",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "referralBonus",
        "type": "uint256"
      }
    ],
    "name": "ReferralBonusAdded",
    "type": "event"
  },
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
        "internalType": "contract IAMM",
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
    "inputs": [],
    "name": "getUnderlyingPrice",
    "outputs": [
      {
        "internalType": "uint256[]",
        "name": "prices",
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
          }
        ],
        "internalType": "struct IOrderBook.Order",
        "name": "order",
        "type": "tuple"
      },
      {
        "components": [
          {
            "internalType": "bytes32",
            "name": "orderHash",
            "type": "bytes32"
          },
          {
            "internalType": "uint256",
            "name": "blockPlaced",
            "type": "uint256"
          },
          {
            "internalType": "enum IOrderBook.OrderExecutionMode",
            "name": "mode",
            "type": "uint8"
          }
        ],
        "internalType": "struct IOrderBook.MatchInfo",
        "name": "matchInfo",
        "type": "tuple"
      },
      {
        "internalType": "int256",
        "name": "liquidationAmount",
        "type": "int256"
      },
      {
        "internalType": "uint256",
        "name": "price",
        "type": "uint256"
      },
      {
        "internalType": "address",
        "name": "trader",
        "type": "address"
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
          },
          {
            "internalType": "bool",
            "name": "reduceOnly",
            "type": "bool"
          }
        ],
        "internalType": "struct IOrderBook.Order[2]",
        "name": "orders",
        "type": "tuple[2]"
      },
      {
        "components": [
          {
            "internalType": "bytes32",
            "name": "orderHash",
            "type": "bytes32"
          },
          {
            "internalType": "uint256",
            "name": "blockPlaced",
            "type": "uint256"
          },
          {
            "internalType": "enum IOrderBook.OrderExecutionMode",
            "name": "mode",
            "type": "uint8"
          }
        ],
        "internalType": "struct IOrderBook.MatchInfo[2]",
        "name": "matchInfo",
        "type": "tuple[2]"
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
      }
    ],
    "name": "openComplementaryPositions",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [],
    "name": "orderBook",
    "outputs": [
      {
        "internalType": "contract IOrderBook",
        "name": "",
        "type": "address"
      }
    ],
    "stateMutability": "view",
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
    "name": "updatePositions",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  }
]}`)
