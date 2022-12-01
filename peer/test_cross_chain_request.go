package peer

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/subnet-evm/accounts/abi"
	"github.com/ava-labs/subnet-evm/internal/ethapi"
	"github.com/ava-labs/subnet-evm/plugin/evm/message"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
)

func SendEthCallCrossChainRequest(networkClient NetworkClient, chainID ids.ID) error {
	abiJSON, err := os.ReadFile("./abi/ERC20NativeMinterABI.json")
	if err != nil {
		log.Error("failed to read file", "error", err)
		return err
	}

	abi, err := abi.JSON(strings.NewReader(string(abiJSON)))
	if err != nil {
		log.Error("failed to parse ABI interface ", "error", err)

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
		log.Error("failed to marshal into JSON encoding ", "error", err)

		return err
	}

	ethCallRequest, err := message.CrossChainCodec.Marshal(message.Version, message.EthCallRequest{RequestArgs: ethCallBytes})
	if err != nil {
		log.Error("failed to marshal into codec encoding", "error", err)
		return err
	}

	response, err := networkClient.CrossChainRequest(chainID, ethCallRequest)
	if err != nil {
		log.Error("failed to send CCR", "error", err)
		return err
	}

	log.Info("Success ! Response from CCR", "response", response)
	return nil
}
