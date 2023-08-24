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
<<<<<<< HEAD
				expectedOutput, err := PackGetBlockchainIDOutput(blockchainID)
=======
				expectedOutput, err := PackGetBlockchainIDOutput(common.Hash(blockchainID))
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
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
<<<<<<< HEAD
				expectedOutput, err := PackGetBlockchainIDOutput(blockchainID)
=======
				expectedOutput, err := PackGetBlockchainIDOutput(common.Hash(blockchainID))
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
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
<<<<<<< HEAD
		DestinationChainID: destinationChainID,
		DestinationAddress: receiverAddr.Hash(),
=======
		DestinationChainID: common.Hash(destinationChainID),
		DestinationAddress: receiverAddr,
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
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
<<<<<<< HEAD
				mutatedSendWarpMessageInput := common.CopyBytes(sendWarpMessageInput)
				return mutatedSendWarpMessageInput[:4] // Include only the function selector, so that the input is invalid
=======
				return sendWarpMessageInput[:4] // Include only the function selector, so that the input is invalid
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
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

<<<<<<< HEAD
				require.Equal(t, unsignedWarpMsg.DestinationChainID, destinationChainID)
				require.Equal(t, unsignedWarpMsg.SourceChainID, blockchainID)
				require.Equal(t, addressedPayload.DestinationAddress, ids.ID(receiverAddr.Hash()))
				require.Equal(t, addressedPayload.SourceAddress, ids.ID(callerAddr.Hash()))
=======
				require.Equal(t, addressedPayload.SourceAddress, callerAddr)
				require.Equal(t, unsignedWarpMsg.SourceChainID, blockchainID)
				require.Equal(t, addressedPayload.DestinationChainID, common.Hash(destinationChainID))
				require.Equal(t, addressedPayload.DestinationAddress, receiverAddr)
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
				require.Equal(t, addressedPayload.Payload, sendWarpMessagePayload)
			},
		},
	}

	testutils.RunPrecompileTests(t, Module, state.NewTestStateDB, tests)
}

func TestGetVerifiedWarpMessage(t *testing.T) {
<<<<<<< HEAD
=======
	networkID := uint32(54321)
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
	callerAddr := common.HexToAddress("0x0123")
	sourceAddress := common.HexToAddress("0x456789")
	destinationAddress := common.HexToAddress("0x987654")
	sourceChainID := ids.GenerateTestID()
	packagedPayloadBytes := []byte("mcsorley")
	addressedPayload, err := warpPayload.NewAddressedPayload(
<<<<<<< HEAD
		ids.ID(sourceAddress.Hash()),
		ids.ID(destinationAddress.Hash()),
		packagedPayloadBytes,
	)
	require.NoError(t, err)
	unsignedWarpMsg, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, addressedPayload.Bytes())
=======
		sourceAddress,
		common.Hash(destinationChainID),
		destinationAddress,
		packagedPayloadBytes,
	)
	require.NoError(t, err)
	unsignedWarpMsg, err := avalancheWarp.NewUnsignedMessage(networkID, sourceChainID, addressedPayload.Bytes())
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
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
<<<<<<< HEAD
			SuppliedGas: GasCostPerWarpMessageBytes * uint64(len(warpMessagePredicateBytes)),
=======
			SuppliedGas: GetVerifiedWarpMessageBaseCost + GasCostPerWarpMessageBytes*uint64(len(warpMessagePredicateBytes)),
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
			ReadOnly:    false,
			ExpectedRes: func() []byte {
				res, err := PackGetVerifiedWarpMessageOutput(GetVerifiedWarpMessageOutput{
					Message: WarpMessage{
<<<<<<< HEAD
						OriginChainID:       sourceChainID,
						OriginSenderAddress: sourceAddress.Hash(),
						DestinationChainID:  destinationChainID,
						DestinationAddress:  destinationAddress.Hash(),
=======
						OriginChainID:       common.Hash(sourceChainID),
						OriginSenderAddress: sourceAddress,
						DestinationChainID:  common.Hash(destinationChainID),
						DestinationAddress:  destinationAddress,
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
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
<<<<<<< HEAD
			SuppliedGas: 0,
=======
			SuppliedGas: GetVerifiedWarpMessageBaseCost,
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
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
<<<<<<< HEAD
			SuppliedGas: GasCostPerWarpMessageBytes * uint64(len(warpMessagePredicateBytes)),
=======
			SuppliedGas: GetVerifiedWarpMessageBaseCost + GasCostPerWarpMessageBytes*uint64(len(warpMessagePredicateBytes)),
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
			ReadOnly:    true,
			ExpectedRes: func() []byte {
				res, err := PackGetVerifiedWarpMessageOutput(GetVerifiedWarpMessageOutput{
					Message: WarpMessage{
<<<<<<< HEAD
						OriginChainID:       sourceChainID,
						OriginSenderAddress: sourceAddress.Hash(),
						DestinationChainID:  destinationChainID,
						DestinationAddress:  destinationAddress.Hash(),
=======
						OriginChainID:       common.Hash(sourceChainID),
						OriginSenderAddress: sourceAddress,
						DestinationChainID:  common.Hash(destinationChainID),
						DestinationAddress:  destinationAddress,
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
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
<<<<<<< HEAD
			SuppliedGas: 0,
=======
			SuppliedGas: GetVerifiedWarpMessageBaseCost,
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
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
<<<<<<< HEAD
			SuppliedGas: GasCostPerWarpMessageBytes*uint64(len(warpMessagePredicateBytes)) - 1,
=======
			SuppliedGas: GetVerifiedWarpMessageBaseCost + GasCostPerWarpMessageBytes*uint64(len(warpMessagePredicateBytes)) - 1,
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
			ReadOnly:    false,
			ExpectedErr: vmerrs.ErrOutOfGas.Error(),
		},
		"get message invalid predicate packing": {
			Caller:  callerAddr,
			InputFn: func(t testing.TB) []byte { return getVerifiedWarpMsg },
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				state.SetPredicateStorageSlots(ContractAddress, warpMessage.Bytes())
			},
<<<<<<< HEAD
			SuppliedGas: GasCostPerWarpMessageBytes * uint64(len(warpMessage.Bytes())),
=======
			SuppliedGas: GetVerifiedWarpMessageBaseCost + GasCostPerWarpMessageBytes*uint64(len(warpMessage.Bytes())),
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
			ReadOnly:    false,
			ExpectedErr: errInvalidPredicateBytes.Error(),
		},
		"get message invalid warp message": {
			Caller:  callerAddr,
			InputFn: func(t testing.TB) []byte { return getVerifiedWarpMsg },
			BeforeHook: func(t testing.TB, state contract.StateDB) {
				state.SetPredicateStorageSlots(ContractAddress, predicateutils.PackPredicate([]byte{1, 2, 3}))
			},
<<<<<<< HEAD
			SuppliedGas: GasCostPerWarpMessageBytes * uint64(32),
=======
			SuppliedGas: GetVerifiedWarpMessageBaseCost + GasCostPerWarpMessageBytes*uint64(32),
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
			ReadOnly:    false,
			ExpectedErr: errInvalidWarpMsg.Error(),
		},
		"get message invalid addressed payload": {
			Caller:  callerAddr,
			InputFn: func(t testing.TB) []byte { return getVerifiedWarpMsg },
			BeforeHook: func(t testing.TB, state contract.StateDB) {
<<<<<<< HEAD
				unsignedMessage, err := avalancheWarp.NewUnsignedMessage(sourceChainID, destinationChainID, []byte{1, 2, 3}) // Invalid addressed payload
=======
				unsignedMessage, err := avalancheWarp.NewUnsignedMessage(networkID, sourceChainID, []byte{1, 2, 3}) // Invalid addressed payload
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
				require.NoError(t, err)
				warpMessage, err := avalancheWarp.NewMessage(unsignedMessage, &avalancheWarp.BitSetSignature{})
				require.NoError(t, err)

				state.SetPredicateStorageSlots(ContractAddress, predicateutils.PackPredicate(warpMessage.Bytes()))
			},
<<<<<<< HEAD
			SuppliedGas: GasCostPerWarpMessageBytes * uint64(192),
=======
			SuppliedGas: GetVerifiedWarpMessageBaseCost + GasCostPerWarpMessageBytes*uint64(160),
>>>>>>> c56d42d51da4d5423aa192d99e33a85c2b82747d
			ReadOnly:    false,
			ExpectedErr: errInvalidAddressedPayload.Error(),
		},
	}

	testutils.RunPrecompileTests(t, Module, state.NewTestStateDB, tests)
}
