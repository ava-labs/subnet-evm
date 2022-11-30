package message

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/internal/ethapi"
	"github.com/ava-labs/subnet-evm/peer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func SendEthCallCrossChainRequest(networkClient peer.NetworkClient, chainID ids.ID) error {
	abiJSON, err := os.ReadFile("./abi/ERC20NativeMinterABI.json")
	if err != nil {
		return err
	}

	abi, err := abi.JSON(strings.NewReader(string(abiJSON)))
	if err != nil {
		return err
	}

	dataBytes, err := abi.Pack("decimals", "")
	data := hexutil.Bytes(dataBytes)

	address := common.HexToAddress("0XFAKE")

	//TO BE FILLED
	ethCallArgs := &ethapi.TransactionArgs{
		To:   &address,
		Data: &data,
	}

	ethCallBytes, err := json.Marshal(ethCallArgs)
	if err != nil {
		return err
	}

	ethCallRequest, err := CrossChainCodec.Marshal(Version, EthCallRequest{RequestArgs: ethCallBytes})
	if err != nil {
		return err
	}

	response, err := networkClient.CrossChainRequest(chainID, ethCallRequest)
	if err != nil {
		return err
	}

	return nil
}
