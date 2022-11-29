package message

import (
	"encoding/json"

	"github.com/ava-labs/subnet-evm/internal/ethapi"
)


func SendEthCallCrossChainRequest(networkClient peer.NetworkClient, chainID id.ID) error {
	// TO BE FILLED
	ethCallArgs := &ethapi.TransactionArgs{
		To: 
		Data: 
	} 

	ethCallBytes, err := json.Marshal(ethCallArgs)
	if err != nil {
		return err 
	}

	ethCallRequest := CrossChainCodec.Marshal(Version, EthCallRequest{RequestArgs: ethCallBytes})

	response, err := networkClient.CrossChainRequest(chainID, ethCethCallRequest)
	if err != nil {
		return err
	}

	return nil 
}