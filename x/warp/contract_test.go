// (c) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"testing"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/utils"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/subnet-evm/core/state"
	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/testutils"
	predicateutils "github.com/ava-labs/subnet-evm/utils/predicate"
	"github.com/ava-labs/subnet-evm/vmerrs"
	warpPayload "github.com/ava-labs/subnet-evm/warp/payload"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestGetBlockchainID(t *testing.T) {
	callerAddr := common.HexToAddress("0x0123")

	defaultSnowCtx := snow.DefaultContextTest()
	blockchainID := defaultSnowCtx.ChainID

	tests := map[string]testutils.PrecompileTest{
		"getBlockchainID success": {
			Caller: callerAddr,
			InputFn: func(t testing.TB) []byte {
				input, err := PackGetBlockchainID()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: GetBlockchainIDGasCost,
			ReadOnly:    false,
			ExpectedRes: func() []byte {
				expectedOutput, err := PackGetBlockchainIDOutput(common.Hash(blockchainID))
				require.NoError(t, err)

				return expectedOutput
			}(),
		},
		"getBlockchainID readOnly": {
			Caller: callerAddr,
			InputFn: func(t testing.TB) []byte {
				input, err := PackGetBlockchainID()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: GetBlockchainIDGasCost,
			ReadOnly:    true,
			ExpectedRes: func() []byte {
				expectedOutput, err := PackGetBlockchainIDOutput(common.Hash(blockchainID))
				require.NoError(t, err)

				return expectedOutput
			}(),
		},
		"getBlockchainID insufficient gas": {
			Caller: callerAddr,
			InputFn: func(t testing.TB) []byte {
				input, err := PackGetBlockchainID()
				require.NoError(t, err)

				return input
			},
			SuppliedGas: GetBlockchainIDGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
	}

	testutils.RunPrecompileTests(t, Module, state.NewTestStateDB, tests)
}

func TestSendWarpMessage(t *testing.T) {
	callerAddr := common.HexToAddress("0x0123")
	receiverAddr := common.HexToAddress("0x456789")

	defaultSnowCtx := snow.DefaultContextTest()
	blockchainID := defaultSnowCtx.ChainID
	destinationChainID := ids.GenerateTestID()
	sendWarpMessagePayload := utils.RandomBytes(100)

	sendWarpMessageInput, err := PackSendWarpMessage(SendWarpMessageInput{
		DestinationChainID: common.Hash(destinationChainID),
		DestinationAddress: receiverAddr.Hash(),
		Payload:            sendWarpMessagePayload,
	})
	require.NoError(t, err)

	tests := map[string]testutils.PrecompileTest{
		"send warp message readOnly": {
			Caller:      callerAddr,
			InputFn:     func(t testing.TB) []byte { return sendWarpMessageInput },
			SuppliedGas: SendWarpMessageGasCost + uint64(len(sendWarpMessageInput[4:])*int(SendWarpMessageGasCostPerByte)),
			ReadOnly:    true,
			ExpectedErr: vmerrs.ErrWriteProtection.Error(),
		},
		"send warp message insufficient gas for first step": {
			Caller:      callerAddr,
			InputFn:     func(t testing.TB) []byte { return sendWarpMessageInput },
			SuppliedGas: SendWarpMessageGasCost - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"send warp message insufficient gas for payload bytes": {
			Caller:      callerAddr,
			InputFn:     func(t testing.TB) []byte { return sendWarpMessageInput },
			SuppliedGas: SendWarpMessageGasCost + uint64(len(sendWarpMessageInput[4:])*int(SendWarpMessageGasCostPerByte)) - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"send warp message invalid input": {
			Caller: callerAddr,
			InputFn: func(t testing.TB) []byte {
				return sendWarpMessageInput[:4] // Include only the function selector, so that the input is invalid
			},
			SuppliedGas: SendWarpMessageGasCost,
			ReadOnly:    false,
			ExpectedErr: errInvalidSendInput.Error(),
		},
		"send warp message success": {
			Caller:      callerAddr,
			InputFn:     func(t testing.TB) []byte { return sendWarpMessageInput },
			SuppliedGas: SendWarpMessageGasCost + uint64(len(sendWarpMessageInput[4:])*int(SendWarpMessageGasCostPerByte)),
			ReadOnly:    false,
			ExpectedRes: []byte{},
			AfterHook: func(t testing.TB, state contract.StateDB) {
				logsData := state.GetLogData()
				require.Len(t, logsData, 1)
				logData := logsData[0]

				unsignedWarpMsg, err := avalancheWarp.ParseUnsignedMessage(logData)
				require.NoError(t, err)
				addressedPayload, err := warpPayload.ParseAddressedPayload(unsignedWarpMsg.Payload)
				require.NoError(t, err)

				require.Equal(t, unsignedWarpMsg.DestinationChainID, destinationChainID)
				require.Equal(t, unsignedWarpMsg.SourceChainID, common.Hash(blockchainID))
				require.Equal(t, addressedPayload.DestinationAddress, ids.ID(receiverAddr.Hash()))
				require.Equal(t, addressedPayload.SourceAddress, ids.ID(callerAddr.Hash()))
				require.Equal(t, addressedPayload.Payload, sendWarpMessagePayload)
			},
		},
	}

	testutils.RunPrecompileTests(t, Module, state.NewTestStateDB, tests)
}

func TestGetVerifiedWarpMessage(t *testing.T) {
	callerAddr := common.HexToAddress("0x0123")
	sourceAddress := common.HexToAddress("0x456789")
	destinationAddress := common.HexToAddress("0x987654")
	sourceChainID := ids.GenerateTestID()
	packagedPayloadBytes := []byte("mcsorley")
	addressedPayload, err := warpPayload.NewAddressedPayload(
		ids.ID(sourceAddress.Hash()),
		ids.ID(destinationAddress.Hash()),
		packagedPayloadBytes,
	)
	require.NoError(t, err)
	unsignedWarpMsg, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, addressedPayload.Bytes())
	require.NoError(t, err)
	warpMessage, err := avalancheWarp.NewMessage(unsignedWarpMsg, &avalancheWarp.BitSetSignature{}) // Create message with empty signature for testing
	require.NoError(t, err)
	warpMessagePredicateBytes := predicateutils.PackPredicate(warpMessage.Bytes())
	getVerifiedWarpMsg, err := PackGetVerifiedWarpMessage()
	require.NoError(t, err)

	tests := map[string]testutils.PrecompileTest{
		"get message success": {
			Caller:  callerAddr,
			InputFn: func(t testing.TB) []byte { return getVerifiedWarpMsg },
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				state.SetPredicateStorageSlots(ContractAddress, warpMessagePredicateBytes)
			},
			SuppliedGas: GasCostPerWarpMessageBytes * uint64(len(warpMessagePredicateBytes)),
			ReadOnly:    false,
			ExpectedRes: func() []byte {
				res, err := PackGetVerifiedWarpMessageOutput(GetVerifiedWarpMessageOutput{
					Message: WarpMessage{
						OriginChainID:       common.Hash(sourceChainID),
						OriginSenderAddress: sourceAddress.Hash(),
						DestinationChainID:  common.Hash(destinationChainID),
						DestinationAddress:  destinationAddress.Hash(),
						Payload:             packagedPayloadBytes,
					},
					Exists: true,
				})
				if err != nil {
					panic(err)
				}
				return res
			}(),
		},
		"get non-existent message": {
			Caller:      callerAddr,
			InputFn:     func(t testing.TB) []byte { return getVerifiedWarpMsg },
			SuppliedGas: 0,
			ReadOnly:    false,
			ExpectedRes: func() []byte {
				res, err := PackGetVerifiedWarpMessageOutput(GetVerifiedWarpMessageOutput{Exists: false})
				if err != nil {
					panic(err)
				}
				return res
			}(),
		},
		"get message success readOnly": {
			Caller:  callerAddr,
			InputFn: func(t testing.TB) []byte { return getVerifiedWarpMsg },
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				state.SetPredicateStorageSlots(ContractAddress, warpMessagePredicateBytes)
			},
			SuppliedGas: GasCostPerWarpMessageBytes * uint64(len(warpMessagePredicateBytes)),
			ReadOnly:    true,
			ExpectedRes: func() []byte {
				res, err := PackGetVerifiedWarpMessageOutput(GetVerifiedWarpMessageOutput{
					Message: WarpMessage{
						OriginChainID:       common.Hash(sourceChainID),
						OriginSenderAddress: sourceAddress.Hash(),
						DestinationChainID:  common.Hash(destinationChainID),
						DestinationAddress:  destinationAddress.Hash(),
						Payload:             packagedPayloadBytes,
					},
					Exists: true,
				})
				if err != nil {
					panic(err)
				}
				return res
			}(),
		},
		"get non-existent message readOnly": {
			Caller:      callerAddr,
			InputFn:     func(t testing.TB) []byte { return getVerifiedWarpMsg },
			SuppliedGas: 0,
			ReadOnly:    true,
			ExpectedRes: func() []byte {
				res, err := PackGetVerifiedWarpMessageOutput(GetVerifiedWarpMessageOutput{Exists: false})
				if err != nil {
					panic(err)
				}
				return res
			}(),
		},
		"get message out of gas": {
			Caller:  callerAddr,
			InputFn: func(t testing.TB) []byte { return getVerifiedWarpMsg },
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				state.SetPredicateStorageSlots(ContractAddress, warpMessagePredicateBytes)
			},
			SuppliedGas: GasCostPerWarpMessageBytes*uint64(len(warpMessagePredicateBytes)) - 1,
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"get message invalid predicate packing": {
			Caller:  callerAddr,
			InputFn: func(t testing.TB) []byte { return getVerifiedWarpMsg },
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				state.SetPredicateStorageSlots(ContractAddress, warpMessage.Bytes())
			},
			SuppliedGas: GasCostPerWarpMessageBytes * uint64(len(warpMessage.Bytes())),
			ReadOnly:    false,
			ExpectedErr: errInvalidPredicateBytes.Error(),
		},
		"get message invalid warp message": {
			Caller:  callerAddr,
			InputFn: func(t testing.TB) []byte { return getVerifiedWarpMsg },
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				state.SetPredicateStorageSlots(ContractAddress, predicateutils.PackPredicate([]byte{1, 2, 3}))
			},
			SuppliedGas: GasCostPerWarpMessageBytes * uint64(32),
			ReadOnly:    false,
			ExpectedErr: errInvalidWarpMsg.Error(),
		},
		"get message invalid addressed payload": {
			Caller:  callerAddr,
			InputFn: func(t testing.TB) []byte { return getVerifiedWarpMsg },
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				unsignedMessage, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, []byte{1, 2, 3}) // Invalid addressed payload
				require.NoError(t, err)
				warpMessage, err := avalancheWarp.NewMessage(unsignedMessage, &avalancheWarp.BitSetSignature{})
				require.NoError(t, err)

				state.SetPredicateStorageSlots(ContractAddress, predicateutils.PackPredicate(warpMessage.Bytes()))
			},
			SuppliedGas: GasCostPerWarpMessageBytes * uint64(192),
			ReadOnly:    false,
			ExpectedErr: errInvalidAddressedPayload.Error(),
		},
	}

	testutils.RunPrecompileTests(t, Module, state.NewTestStateDB, tests)
}
