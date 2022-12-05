package peer

import (
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/subnet-evm/plugin/evm/message"
	"github.com/ethereum/go-ethereum/log"
)

func SendPongCrossChainRequest(networkClient NetworkClient, chainID ids.ID) error {
	var pongRequest message.CrossChainRequest = message.PongRequest{
		RequestBytes: []byte("ping"),
	}

	crossChainCodec := message.CrossChainCodec
	crossChainRequest, err := crossChainCodec.Marshal(message.Version, &pongRequest)
	if err != nil {
		log.Error("failed to marshal into codec encoding", "error", err)
		return err
	}

	responseBytes, err := networkClient.CrossChainRequest(chainID, crossChainRequest)
	if err != nil {
		log.Error("error with CCR", "error", err)
		return err
	}

	crossChainResponse := &message.PongResponse{}
	_, err = crossChainCodec.Unmarshal(responseBytes, crossChainResponse)
	if err != nil {
		log.Error("error with CCR", "error", err)
		return err
	}

	result := string(crossChainResponse.ResponseBytes[:])

	log.Info("Success ! Response from CCR", "response", result)
	return nil

}

// func SendEthCallCrossChainRequest(networkClient NetworkClient, chainID ids.ID) error {
// 	log.Info("_--------------------------------- SLEEP CCR  -------------------------------------")
// 	time.Sleep(60 * time.Second)
// 	log.Info("_--------------------------------- START CCR  -------------------------------------")

// 	abi, err := getABI()
// 	if err != nil {
// 		log.Error("failed to parse ABI interface ", "error", err)
// 		return err
// 	}

// 	dataBytes, err := abi.Pack("decimals")
// 	data := hexutil.Bytes(dataBytes)
// 	address := common.HexToAddress("0x52C84043CD9c865236f11d9Fc9F56aa003c1f922")

// 	ethCallArgs := &ethapi.TransactionArgs{
// 		To:   &address,
// 		Data: &data,
// 	}

// 	ethCallBytes, err := json.Marshal(ethCallArgs)
// 	if err != nil {
// 		log.Error("failed to marshal into JSON encoding ", "error", err)

// 		return err
// 	}

// 	var ethCallRequest message.CrossChainRequest = message.EthCallRequest{
// 		RequestArgs: ethCallBytes,
// 	}

// 	crossChainCodec := message.CrossChainCodec
// 	crossChainRequest, err := crossChainCodec.Marshal(message.Version, &ethCallRequest)
// 	if err != nil {
// 		log.Error("failed to marshal into codec encoding", "error", err)
// 		return err
// 	}

// 	responseBytes, err := networkClient.CrossChainRequest(chainID, crossChainRequest)
// 	if err != nil {
// 		log.Error("error with CCR", "error", err)
// 		return err
// 	}

// 	crossChainResponse := &message.EthCallResponse{}
// 	_, err = crossChainCodec.Unmarshal(responseBytes, crossChainResponse)
// 	if err != nil {
// 		log.Error("error with CCR", "error", err)
// 		return err
// 	}

// 	executionResult := &core.ExecutionResult{}
// 	json.Unmarshal(crossChainResponse.ExecutionResult, executionResult)
// 	log.Info("Success ! Response from CCR", "response", executionResult)
// 	return nil

// }

// func getABI() (abi.ABI, error) {
// 	buf := `[{
// 		"inputs": [
// 			{
// 				"internalType": "uint256",
// 				"name": "initSupply",
// 				"type": "uint256"
// 			}
// 		],
// 		"stateMutability": "nonpayable",
// 		"type": "constructor"
// 	},
// 	{
// 		"anonymous": false,
// 		"inputs": [
// 			{
// 				"indexed": true,
// 				"internalType": "address",
// 				"name": "owner",
// 				"type": "address"
// 			},
// 			{
// 				"indexed": true,
// 				"internalType": "address",
// 				"name": "spender",
// 				"type": "address"
// 			},
// 			{
// 				"indexed": false,
// 				"internalType": "uint256",
// 				"name": "value",
// 				"type": "uint256"
// 			}
// 		],
// 		"name": "Approval",
// 		"type": "event"
// 	},
// 	{
// 		"anonymous": false,
// 		"inputs": [
// 			{
// 				"indexed": true,
// 				"internalType": "address",
// 				"name": "dst",
// 				"type": "address"
// 			},
// 			{
// 				"indexed": false,
// 				"internalType": "uint256",
// 				"name": "wad",
// 				"type": "uint256"
// 			}
// 		],
// 		"name": "Deposit",
// 		"type": "event"
// 	},
// 	{
// 		"anonymous": false,
// 		"inputs": [
// 			{
// 				"indexed": true,
// 				"internalType": "address",
// 				"name": "src",
// 				"type": "address"
// 			},
// 			{
// 				"indexed": false,
// 				"internalType": "uint256",
// 				"name": "wad",
// 				"type": "uint256"
// 			}
// 		],
// 		"name": "Mintdrawal",
// 		"type": "event"
// 	},
// 	{
// 		"anonymous": false,
// 		"inputs": [
// 			{
// 				"indexed": true,
// 				"internalType": "address",
// 				"name": "previousOwner",
// 				"type": "address"
// 			},
// 			{
// 				"indexed": true,
// 				"internalType": "address",
// 				"name": "newOwner",
// 				"type": "address"
// 			}
// 		],
// 		"name": "OwnershipTransferred",
// 		"type": "event"
// 	},
// 	{
// 		"anonymous": false,
// 		"inputs": [
// 			{
// 				"indexed": true,
// 				"internalType": "address",
// 				"name": "from",
// 				"type": "address"
// 			},
// 			{
// 				"indexed": true,
// 				"internalType": "address",
// 				"name": "to",
// 				"type": "address"
// 			},
// 			{
// 				"indexed": false,
// 				"internalType": "uint256",
// 				"name": "value",
// 				"type": "uint256"
// 			}
// 		],
// 		"name": "Transfer",
// 		"type": "event"
// 	},
// 	{
// 		"inputs": [
// 			{
// 				"internalType": "address",
// 				"name": "owner",
// 				"type": "address"
// 			},
// 			{
// 				"internalType": "address",
// 				"name": "spender",
// 				"type": "address"
// 			}
// 		],
// 		"name": "allowance",
// 		"outputs": [
// 			{
// 				"internalType": "uint256",
// 				"name": "",
// 				"type": "uint256"
// 			}
// 		],
// 		"stateMutability": "view",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [
// 			{
// 				"internalType": "address",
// 				"name": "spender",
// 				"type": "address"
// 			},
// 			{
// 				"internalType": "uint256",
// 				"name": "amount",
// 				"type": "uint256"
// 			}
// 		],
// 		"name": "approve",
// 		"outputs": [
// 			{
// 				"internalType": "bool",
// 				"name": "",
// 				"type": "bool"
// 			}
// 		],
// 		"stateMutability": "nonpayable",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [
// 			{
// 				"internalType": "address",
// 				"name": "account",
// 				"type": "address"
// 			}
// 		],
// 		"name": "balanceOf",
// 		"outputs": [
// 			{
// 				"internalType": "uint256",
// 				"name": "",
// 				"type": "uint256"
// 			}
// 		],
// 		"stateMutability": "view",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [
// 			{
// 				"internalType": "address",
// 				"name": "from",
// 				"type": "address"
// 			},
// 			{
// 				"internalType": "uint256",
// 				"name": "amount",
// 				"type": "uint256"
// 			}
// 		],
// 		"name": "burn",
// 		"outputs": [],
// 		"stateMutability": "nonpayable",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [],
// 		"name": "decimals",
// 		"outputs": [
// 			{
// 				"internalType": "uint8",
// 				"name": "",
// 				"type": "uint8"
// 			}
// 		],
// 		"stateMutability": "view",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [
// 			{
// 				"internalType": "address",
// 				"name": "spender",
// 				"type": "address"
// 			},
// 			{
// 				"internalType": "uint256",
// 				"name": "subtractedValue",
// 				"type": "uint256"
// 			}
// 		],
// 		"name": "decreaseAllowance",
// 		"outputs": [
// 			{
// 				"internalType": "bool",
// 				"name": "",
// 				"type": "bool"
// 			}
// 		],
// 		"stateMutability": "nonpayable",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [],
// 		"name": "deposit",
// 		"outputs": [],
// 		"stateMutability": "payable",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [
// 			{
// 				"internalType": "address",
// 				"name": "spender",
// 				"type": "address"
// 			},
// 			{
// 				"internalType": "uint256",
// 				"name": "addedValue",
// 				"type": "uint256"
// 			}
// 		],
// 		"name": "increaseAllowance",
// 		"outputs": [
// 			{
// 				"internalType": "bool",
// 				"name": "",
// 				"type": "bool"
// 			}
// 		],
// 		"stateMutability": "nonpayable",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [
// 			{
// 				"internalType": "address",
// 				"name": "addr",
// 				"type": "address"
// 			}
// 		],
// 		"name": "isAdmin",
// 		"outputs": [
// 			{
// 				"internalType": "bool",
// 				"name": "",
// 				"type": "bool"
// 			}
// 		],
// 		"stateMutability": "view",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [
// 			{
// 				"internalType": "address",
// 				"name": "addr",
// 				"type": "address"
// 			}
// 		],
// 		"name": "isEnabled",
// 		"outputs": [
// 			{
// 				"internalType": "bool",
// 				"name": "",
// 				"type": "bool"
// 			}
// 		],
// 		"stateMutability": "view",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [
// 			{
// 				"internalType": "address",
// 				"name": "to",
// 				"type": "address"
// 			},
// 			{
// 				"internalType": "uint256",
// 				"name": "amount",
// 				"type": "uint256"
// 			}
// 		],
// 		"name": "mint",
// 		"outputs": [],
// 		"stateMutability": "nonpayable",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [
// 			{
// 				"internalType": "uint256",
// 				"name": "wad",
// 				"type": "uint256"
// 			}
// 		],
// 		"name": "mintdraw",
// 		"outputs": [],
// 		"stateMutability": "nonpayable",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [],
// 		"name": "name",
// 		"outputs": [
// 			{
// 				"internalType": "string",
// 				"name": "",
// 				"type": "string"
// 			}
// 		],
// 		"stateMutability": "view",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [],
// 		"name": "owner",
// 		"outputs": [
// 			{
// 				"internalType": "address",
// 				"name": "",
// 				"type": "address"
// 			}
// 		],
// 		"stateMutability": "view",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [],
// 		"name": "renounceOwnership",
// 		"outputs": [],
// 		"stateMutability": "nonpayable",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [
// 			{
// 				"internalType": "address",
// 				"name": "addr",
// 				"type": "address"
// 			}
// 		],
// 		"name": "revoke",
// 		"outputs": [],
// 		"stateMutability": "nonpayable",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [
// 			{
// 				"internalType": "address",
// 				"name": "addr",
// 				"type": "address"
// 			}
// 		],
// 		"name": "setAdmin",
// 		"outputs": [],
// 		"stateMutability": "nonpayable",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [
// 			{
// 				"internalType": "address",
// 				"name": "addr",
// 				"type": "address"
// 			}
// 		],
// 		"name": "setEnabled",
// 		"outputs": [],
// 		"stateMutability": "nonpayable",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [],
// 		"name": "symbol",
// 		"outputs": [
// 			{
// 				"internalType": "string",
// 				"name": "",
// 				"type": "string"
// 			}
// 		],
// 		"stateMutability": "view",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [],
// 		"name": "totalSupply",
// 		"outputs": [
// 			{
// 				"internalType": "uint256",
// 				"name": "",
// 				"type": "uint256"
// 			}
// 		],
// 		"stateMutability": "view",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [
// 			{
// 				"internalType": "address",
// 				"name": "to",
// 				"type": "address"
// 			},
// 			{
// 				"internalType": "uint256",
// 				"name": "amount",
// 				"type": "uint256"
// 			}
// 		],
// 		"name": "transfer",
// 		"outputs": [
// 			{
// 				"internalType": "bool",
// 				"name": "",
// 				"type": "bool"
// 			}
// 		],
// 		"stateMutability": "nonpayable",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [
// 			{
// 				"internalType": "address",
// 				"name": "from",
// 				"type": "address"
// 			},
// 			{
// 				"internalType": "address",
// 				"name": "to",
// 				"type": "address"
// 			},
// 			{
// 				"internalType": "uint256",
// 				"name": "amount",
// 				"type": "uint256"
// 			}
// 		],
// 		"name": "transferFrom",
// 		"outputs": [
// 			{
// 				"internalType": "bool",
// 				"name": "",
// 				"type": "bool"
// 			}
// 		],
// 		"stateMutability": "nonpayable",
// 		"type": "function"
// 	},
// 	{
// 		"inputs": [
// 			{
// 				"internalType": "address",
// 				"name": "newOwner",
// 				"type": "address"
// 			}
// 		],
// 		"name": "transferOwnership",
// 		"outputs": [],
// 		"stateMutability": "nonpayable",
// 		"type": "function"
// 	}
// ]`

// 	return abi.JSON(strings.NewReader(string(buf)))
// }
