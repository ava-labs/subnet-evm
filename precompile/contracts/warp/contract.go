// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ava-labs/subnet-evm/vmerrs"
	warpPayload "github.com/ava-labs/subnet-evm/warp/payload"

	_ "embed"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/log"
)

const (
	GetBlockchainIDGasCost uint64 = 2 // Based on GasQuickStep used in existing EVM instructions
	// Cost of fixed log size and writing log to database (Note: this calculation includes buffer room because a warp message
	// is not stored in the trie and can be safely deleted within a reasonable period of time)
	SendWarpMessageGasCost uint64 = params.LogGas + 4*params.LogTopicGas + contract.WriteGasCostPerSlot
	// TODO: SendWarpMessageGasCostPerByte cost accounts for producing the log itself as well as the added work of signing
	// and storing the message.
	SendWarpMessageGasCostPerByte uint64 = params.TxDataNonZeroGasFrontier

	GasCostPerWarpSigner       uint64 = 300_000                        // TODO: charge for number of signers and the cost of the validator set lookup on the P-Chain
	GasCostPerWarpMessageBytes uint64 = params.TxDataNonZeroGasEIP2028 // TODO: charge O(n) cost for decoding predicate of input size n
	// GasCostPerSourceSubnetValidator uint64 = 1 // TODO: charge O(n) cost for subnet validator set lookup
)

// Singleton StatefulPrecompiledContract and signatures.
var (

	// WarpRawABI contains the raw ABI of Warp contract.
	//go:embed contract.abi
	WarpRawABI string

	WarpABI = contract.ParseABI(WarpRawABI)

	WarpPrecompile = createWarpPrecompile()
)

// WarpMessage is an auto generated low-level Go binding around an user-defined struct.
type WarpMessage struct {
	OriginChainID       [32]byte
	OriginSenderAddress [32]byte
	DestinationChainID  [32]byte
	DestinationAddress  [32]byte
	Payload             []byte
}

type GetVerifiedWarpMessageOutput struct {
	Message WarpMessage
	Exists  bool
}

type SendWarpMessageInput struct {
	DestinationChainID [32]byte
	DestinationAddress [32]byte
	Payload            []byte
}

// PackGetBlockchainID packs the include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackGetBlockchainID() ([]byte, error) {
	return WarpABI.Pack("getBlockchainID")
}

// PackGetBlockchainIDOutput attempts to pack given blockchainID of type [32]byte
// to conform the ABI outputs.
func PackGetBlockchainIDOutput(blockchainID [32]byte) ([]byte, error) {
	return WarpABI.PackOutput("getBlockchainID", blockchainID)
}

// getBlockchainID returns the snow Chain Context ChainID of this blockchain.
func getBlockchainID(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, GetBlockchainIDGasCost); err != nil {
		return nil, 0, err
	}
	packedOutput, err := PackGetBlockchainIDOutput(accessibleState.GetSnowContext().ChainID)
	if err != nil {
		return nil, remainingGas, err
	}

	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// PackGetVerifiedWarpMessage packs the include selector (first 4 func signature bytes).
// This function is mostly used for tests.
func PackGetVerifiedWarpMessage() ([]byte, error) {
	return WarpABI.Pack("getVerifiedWarpMessage")
}

// PackGetVerifiedWarpMessageOutput attempts to pack given [outputStruct] of type GetVerifiedWarpMessageOutput
// to conform the ABI outputs.
func PackGetVerifiedWarpMessageOutput(outputStruct GetVerifiedWarpMessageOutput) ([]byte, error) {
	return WarpABI.PackOutput("getVerifiedWarpMessage",
		outputStruct.Message,
		outputStruct.Exists,
	)
}

// getVerifiedWarpMessage retrieves the pre-verified warp message from the predicate storage slots and returns
// the expected ABI encoding of the message to the caller.
func getVerifiedWarpMessage(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, _ []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	remainingGas = suppliedGas
	// XXX Note: there is no base cost for retrieving a verified warp message. Instead, we charge for each piece of gas,
	// prior to each execution step.
	// Ignore input since there are no arguments
	predicateBytes, exists := accessibleState.GetStateDB().GetPredicateStorageSlots(ContractAddress)
	// If there is no such value, return false to the caller.
	// Note: decoding errors will return an error instead.
	if !exists {
		packedOutput, err := PackGetVerifiedWarpMessageOutput(GetVerifiedWarpMessageOutput{
			Exists: false,
		})
		if err != nil {
			return nil, remainingGas, err
		}
		return packedOutput, remainingGas, nil
	}

	msgBytesGas, overflow := math.SafeMul(GasCostPerWarpMessageBytes, uint64(len(predicateBytes)))
	if overflow {
		return nil, remainingGas, vmerrs.ErrOutOfGas
	}
	if remainingGas, err = contract.DeductGas(remainingGas, msgBytesGas); err != nil {
		return nil, 0, err
	}
	unpackedPredicateBytes, err := utils.UnpackPredicate(predicateBytes)
	if err != nil {
		return nil, remainingGas, err
	}
	warpMessage, err := warp.ParseMessage(unpackedPredicateBytes)
	if err != nil {
		return nil, remainingGas, err
	}

	addressedPayload, err := warpPayload.ParseAddressedPayload(warpMessage.UnsignedMessage.Payload)
	if err != nil {
		return nil, remainingGas, err
	}
	packedOutput, err := PackGetVerifiedWarpMessageOutput(GetVerifiedWarpMessageOutput{
		Message: WarpMessage{
			OriginChainID:       warpMessage.SourceChainID,
			OriginSenderAddress: addressedPayload.SourceAddress,
			DestinationChainID:  warpMessage.DestinationChainID,
			DestinationAddress:  addressedPayload.DestinationAddress,
			Payload:             addressedPayload.Payload,
		},
		Exists: true,
	})
	if err != nil {
		return nil, remainingGas, err
	}

	// TODO: remove
	log.Info("getVerifiedWarpMessage successful", "originaChainID", warpMessage.SourceChainID, "originSenderAddress", addressedPayload.SourceAddress, "destinationChainID", warpMessage.DestinationChainID, "destinationAddress", addressedPayload.DestinationAddress, "payload", common.Bytes2Hex(addressedPayload.Payload))
	// Return the packed output and the remaining gas
	return packedOutput, remainingGas, nil
}

// UnpackSendWarpMessageInput attempts to unpack [input] as SendWarpMessageInput
// assumes that [input] does not include selector (omits first 4 func signature bytes)
func UnpackSendWarpMessageInput(input []byte) (SendWarpMessageInput, error) {
	inputStruct := SendWarpMessageInput{}
	err := WarpABI.UnpackInputIntoInterface(&inputStruct, "sendWarpMessage", input)

	return inputStruct, err
}

// PackSendWarpMessage packs [inputStruct] of type SendWarpMessageInput into the appropriate arguments for sendWarpMessage.
func PackSendWarpMessage(inputStruct SendWarpMessageInput) ([]byte, error) {
	return WarpABI.Pack("sendWarpMessage", inputStruct.DestinationChainID, inputStruct.DestinationAddress, inputStruct.Payload)
}

// sendWarpMessage constructs an Avalanche Warp Message containing an AddressedPayload and emits a log to signal validators that they should
// be willing to sign this message.
func sendWarpMessage(accessibleState contract.AccessibleState, caller common.Address, addr common.Address, input []byte, suppliedGas uint64, readOnly bool) (ret []byte, remainingGas uint64, err error) {
	if remainingGas, err = contract.DeductGas(suppliedGas, SendWarpMessageGasCost); err != nil {
		return nil, 0, err
	}
	// This gas cost includes buffer room because it is based off of the total size of the input instead of the produced payload.
	// This ensures that we charge gas before we unpack the variable sized input.
	payloadGas, overflow := math.SafeMul(SendWarpMessageGasCostPerByte, uint64(len(input)))
	if overflow {
		return nil, 0, vmerrs.ErrOutOfGas
	}
	if remainingGas, err = contract.DeductGas(remainingGas, payloadGas); err != nil {
		return nil, 0, err
	}
	if readOnly {
		return nil, remainingGas, vmerrs.ErrWriteProtection
	}
	// unpack the arguments
	inputStruct, err := UnpackSendWarpMessageInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	var (
		sourceChainID      = accessibleState.GetSnowContext().ChainID
		destinationChainID = inputStruct.DestinationChainID
		sourceAddress      = caller.Hash()
		destinationAddress = inputStruct.DestinationAddress
		payload            = inputStruct.Payload
	)

	addressedPayload, err := warpPayload.NewAddressedPayload(
		ids.ID(sourceAddress),
		destinationAddress,
		payload,
	)
	if err != nil {
		return nil, remainingGas, err
	}
	warpMessage, err := warp.NewUnsignedMessage(
		sourceChainID,
		destinationChainID,
		addressedPayload.Bytes(),
	)
	if err != nil {
		return nil, remainingGas, err
	}

	// Add a log to be handled if this action is finalized.
	accessibleState.GetStateDB().AddLog(
		ContractAddress,
		[]common.Hash{
			WarpABI.Events["SendWarpMessage"].ID,
			destinationChainID,
			destinationAddress,
			sourceAddress,
		},
		warpMessage.Bytes(),
		accessibleState.GetBlockContext().Number().Uint64(),
	)

	// Return an empty output and the remaining gas
	return []byte{}, remainingGas, nil
}

// createWarpPrecompile returns a StatefulPrecompiledContract with getters and setters for the precompile.
func createWarpPrecompile() contract.StatefulPrecompiledContract {
	var functions []*contract.StatefulPrecompileFunction

	abiFunctionMap := map[string]contract.RunStatefulPrecompileFunc{
		"getBlockchainID":        getBlockchainID,
		"getVerifiedWarpMessage": getVerifiedWarpMessage,
		"sendWarpMessage":        sendWarpMessage,
	}

	for name, function := range abiFunctionMap {
		method, ok := WarpABI.Methods[name]
		if !ok {
			panic(fmt.Errorf("given method (%s) does not exist in the ABI", name))
		}
		functions = append(functions, contract.NewStatefulPrecompileFunction(method.ID, function))
	}
	// Construct the contract with no fallback function.
	statefulContract, err := contract.NewStatefulPrecompileContract(nil, functions)
	if err != nil {
		panic(err)
	}
	return statefulContract
}
