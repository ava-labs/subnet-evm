// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"fmt"

	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	predicateutils "github.com/ava-labs/subnet-evm/utils/predicate"
	"github.com/ava-labs/subnet-evm/vmerrs"
	warpPayload "github.com/ava-labs/subnet-evm/warp/payload"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

var (
	_ messageHandler = addressedPayloadHandler{}
	_ messageHandler = blockHashHandler{}
)

var (
	getVerifiedWarpMessageInvalidOutput   []byte
	getVerifiedWarpBlockHashInvalidOutput []byte
)

func init() {
	res, err := PackGetVerifiedWarpMessageOutput(GetVerifiedWarpMessageOutput{Valid: false})
	if err != nil {
		panic(err)
	}
	getVerifiedWarpMessageInvalidOutput = res

	res, err = PackGetVerifiedWarpBlockHashOutput(GetVerifiedWarpBlockHashOutput{Valid: false})
	if err != nil {
		panic(err)
	}
	getVerifiedWarpBlockHashInvalidOutput = res
}

type messageHandler interface {
	packFailed() []byte
	handleMessage(msg *warp.Message) ([]byte, error)
}

func handleWarpMessage(accessibleState contract.AccessibleState, input []byte, remainingGas uint64, handler messageHandler) ([]byte, uint64, error) {
	warpIndex, err := UnpackGetVerifiedWarpMessageInput(input)
	if err != nil {
		return nil, remainingGas, fmt.Errorf("%w: %s", errInvalidIndexInput, err)
	}
	state := accessibleState.GetStateDB()
	predicateBytes, exists := state.GetPredicateStorageSlots(ContractAddress, warpIndex)
	predicateResults := accessibleState.GetBlockContext().GetPredicateResults(state.GetTxHash(), ContractAddress)
	valid := exists && set.BitsFromBytes(predicateResults).Contains(int(warpIndex))
	if !valid {
		return handler.packFailed(), remainingGas, nil
	}

	// Note: we charge for the size of the message during both predicate verification and each time the message is read during
	// EVM execution because each execution incurs an additional read cost.
	msgBytesGas, overflow := math.SafeMul(GasCostPerWarpMessageBytes, uint64(len(predicateBytes)))
	if overflow {
		return nil, 0, vmerrs.ErrOutOfGas
	}
	if remainingGas, err = contract.DeductGas(remainingGas, msgBytesGas); err != nil {
		return nil, 0, err
	}
	// Note: since the predicate is verified in advance of execution, the precompile should not
	// hit an error during execution.
	unpackedPredicateBytes, err := predicateutils.UnpackPredicate(predicateBytes)
	if err != nil {
		return nil, remainingGas, fmt.Errorf("%w: %s", errInvalidPredicateBytes, err)
	}
	warpMessage, err := warp.ParseMessage(unpackedPredicateBytes)
	if err != nil {
		return nil, remainingGas, fmt.Errorf("%w: %s", errInvalidWarpMsg, err)
	}
	res, err := handler.handleMessage(warpMessage)
	if err != nil {
		return nil, remainingGas, err
	}
	return res, remainingGas, nil
}

type addressedPayloadHandler struct{}

func (addressedPayloadHandler) packFailed() []byte {
	return getVerifiedWarpMessageInvalidOutput
}

func (addressedPayloadHandler) handleMessage(warpMessage *warp.Message) ([]byte, error) {
	addressedPayload, err := warpPayload.ParseAddressedPayload(warpMessage.UnsignedMessage.Payload)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errInvalidAddressedPayload, err)
	}
	return PackGetVerifiedWarpMessageOutput(GetVerifiedWarpMessageOutput{
		Message: WarpMessage{
			SourceChainID:       common.Hash(warpMessage.SourceChainID),
			OriginSenderAddress: addressedPayload.SourceAddress,
			DestinationChainID:  addressedPayload.DestinationChainID,
			DestinationAddress:  addressedPayload.DestinationAddress,
			Payload:             addressedPayload.Payload,
		},
		Valid: true,
	})
}

type blockHashHandler struct{}

func (blockHashHandler) packFailed() []byte {
	return getVerifiedWarpBlockHashInvalidOutput
}

func (blockHashHandler) handleMessage(warpMessage *warp.Message) ([]byte, error) {
	blockHashPayload, err := warpPayload.ParseBlockHashPayload(warpMessage.UnsignedMessage.Payload)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errInvalidBlockHashPayload, err)
	}
	return PackGetVerifiedWarpBlockHashOutput(GetVerifiedWarpBlockHashOutput{
		WarpBlockHash: WarpBlockHash{
			SourceChainID: common.Hash(warpMessage.SourceChainID),
			BlockHash:     blockHashPayload.BlockHash,
		},
		Valid: true,
	})
}
