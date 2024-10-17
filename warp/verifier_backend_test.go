// (c) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"context"
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/cache"
	"github.com/ava-labs/avalanchego/database/memdb"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/network/p2p/acp118"
	"github.com/ava-labs/avalanchego/proto/pb/sdk"
	"github.com/ava-labs/avalanchego/snow/engine/common"
	"github.com/ava-labs/avalanchego/snow/uptime"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/timer/mockable"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/payload"
	"github.com/ava-labs/subnet-evm/plugin/evm/validators"
	"github.com/ava-labs/subnet-evm/utils"
	"github.com/ava-labs/subnet-evm/warp/messages"
	"github.com/ava-labs/subnet-evm/warp/warptest"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestAddressedCallSignatures(t *testing.T) {
	database := memdb.New()
	snowCtx := utils.TestSnowContext()
	blsSecretKey, err := bls.NewSecretKey()
	require.NoError(t, err)
	warpSigner := avalancheWarp.NewSigner(blsSecretKey, snowCtx.NetworkID, snowCtx.ChainID)

	offChainPayload, err := payload.NewAddressedCall([]byte{1, 2, 3}, []byte{1, 2, 3})
	require.NoError(t, err)
	offchainMessage, err := avalancheWarp.NewUnsignedMessage(snowCtx.NetworkID, snowCtx.ChainID, offChainPayload.Bytes())
	require.NoError(t, err)
	offchainSignature, err := warpSigner.Sign(offchainMessage)
	require.NoError(t, err)

	tests := map[string]struct {
		setup       func(backend Backend) (request []byte, expectedResponse []byte)
		verifyStats func(t *testing.T, stats *verifierStats)
		err         error
	}{
		"known message": {
			setup: func(backend Backend) (request []byte, expectedResponse []byte) {
				knownPayload, err := payload.NewAddressedCall([]byte{0, 0, 0}, []byte("test"))
				require.NoError(t, err)
				msg, err := avalancheWarp.NewUnsignedMessage(snowCtx.NetworkID, snowCtx.ChainID, knownPayload.Bytes())
				require.NoError(t, err)
				knownSignature, err := warpSigner.Sign(msg)
				require.NoError(t, err)

				backend.AddMessage(msg)
				return msg.Bytes(), knownSignature[:]
			},
			verifyStats: func(t *testing.T, stats *verifierStats) {
				require.EqualValues(t, 0, stats.messageParseFail.Snapshot().Count())
				require.EqualValues(t, 0, stats.addressedCallSignatureValidationFail.Snapshot().Count())
				require.EqualValues(t, 0, stats.blockSignatureValidationFail.Snapshot().Count())
			},
		},
		"offchain message": {
			setup: func(_ Backend) (request []byte, expectedResponse []byte) {
				return offchainMessage.Bytes(), offchainSignature[:]
			},
			verifyStats: func(t *testing.T, stats *verifierStats) {
				require.EqualValues(t, 0, stats.messageParseFail.Snapshot().Count())
				require.EqualValues(t, 0, stats.addressedCallSignatureValidationFail.Snapshot().Count())
				require.EqualValues(t, 0, stats.blockSignatureValidationFail.Snapshot().Count())
			},
		},
		"unknown message": {
			setup: func(_ Backend) (request []byte, expectedResponse []byte) {
				unknownPayload, err := payload.NewAddressedCall([]byte{0, 0, 0}, []byte("unknown message"))
				require.NoError(t, err)
				unknownMessage, err := avalancheWarp.NewUnsignedMessage(snowCtx.NetworkID, snowCtx.ChainID, unknownPayload.Bytes())
				require.NoError(t, err)
				return unknownMessage.Bytes(), nil
			},
			verifyStats: func(t *testing.T, stats *verifierStats) {
				require.EqualValues(t, 0, stats.messageParseFail.Snapshot().Count())
				require.EqualValues(t, 1, stats.addressedCallSignatureValidationFail.Snapshot().Count())
				require.EqualValues(t, 0, stats.blockSignatureValidationFail.Snapshot().Count())
			},
			err: &common.AppError{Code: ParseErrCode},
		},
	}

	for name, test := range tests {
		for _, withCache := range []bool{true, false} {
			if withCache {
				name += "_with_cache"
			} else {
				name += "_no_cache"
			}
			t.Run(name, func(t *testing.T) {
				var sigCache cache.Cacher[ids.ID, []byte]
				if withCache {
					sigCache = &cache.LRU[ids.ID, []byte]{Size: 100}
				} else {
					sigCache = &cache.Empty[ids.ID, []byte]{}
				}
				warpBackend, err := NewBackend(snowCtx.NetworkID, snowCtx.ChainID, warpSigner, warptest.EmptyBlockClient, uptime.NoOpCalculator, validators.NoOpState, database, sigCache, [][]byte{offchainMessage.Bytes()})
				require.NoError(t, err)
				handler := acp118.NewCachedHandler(sigCache, warpBackend, warpSigner)

				requestBytes, expectedResponse := test.setup(warpBackend)
				protoMsg := &sdk.SignatureRequest{Message: requestBytes}
				protoBytes, err := proto.Marshal(protoMsg)
				require.NoError(t, err)
				responseBytes, appErr := handler.AppRequest(context.Background(), ids.GenerateTestNodeID(), time.Time{}, protoBytes)
				if test.err != nil {
					require.Error(t, appErr)
					require.ErrorIs(t, appErr, test.err)
				} else {
					require.Nil(t, appErr)
				}

				test.verifyStats(t, warpBackend.(*backend).stats)

				// If the expected response is empty, assert that the handler returns an empty response and return early.
				if len(expectedResponse) == 0 {
					require.Len(t, responseBytes, 0, "expected response to be empty")
					return
				}
				// check cache is populated
				if withCache {
					require.NotZero(t, warpBackend.(*backend).signatureCache.Len())
				} else {
					require.Zero(t, warpBackend.(*backend).signatureCache.Len())
				}
				response := &sdk.SignatureResponse{}
				require.NoError(t, proto.Unmarshal(responseBytes, response))
				require.NoError(t, err, "error unmarshalling SignatureResponse")

				require.Equal(t, expectedResponse, response.Signature)
			})
		}
	}
}

func TestBlockSignatures(t *testing.T) {
	database := memdb.New()
	snowCtx := utils.TestSnowContext()
	blsSecretKey, err := bls.NewSecretKey()
	require.NoError(t, err)

	warpSigner := avalancheWarp.NewSigner(blsSecretKey, snowCtx.NetworkID, snowCtx.ChainID)
	blkID := ids.GenerateTestID()
	blockClient := warptest.MakeBlockClient(blkID)

	unknownBlockID := ids.GenerateTestID()

	toMessageBytes := func(id ids.ID) []byte {
		idPayload, err := payload.NewHash(id)
		require.NoError(t, err)

		msg, err := avalancheWarp.NewUnsignedMessage(snowCtx.NetworkID, snowCtx.ChainID, idPayload.Bytes())
		require.NoError(t, err)

		return msg.Bytes()
	}

	tests := map[string]struct {
		setup       func() (request []byte, expectedResponse []byte)
		verifyStats func(t *testing.T, stats *verifierStats)
		err         error
	}{
		"known block": {
			setup: func() (request []byte, expectedResponse []byte) {
				hashPayload, err := payload.NewHash(blkID)
				require.NoError(t, err)
				unsignedMessage, err := avalancheWarp.NewUnsignedMessage(snowCtx.NetworkID, snowCtx.ChainID, hashPayload.Bytes())
				require.NoError(t, err)
				signature, err := warpSigner.Sign(unsignedMessage)
				require.NoError(t, err)
				return toMessageBytes(blkID), signature[:]
			},
			verifyStats: func(t *testing.T, stats *verifierStats) {
				require.EqualValues(t, 0, stats.addressedCallSignatureValidationFail.Snapshot().Count())
				require.EqualValues(t, 0, stats.blockSignatureValidationFail.Snapshot().Count())
				require.EqualValues(t, 0, stats.messageParseFail.Snapshot().Count())
			},
		},
		"unknown block": {
			setup: func() (request []byte, expectedResponse []byte) {
				return toMessageBytes(unknownBlockID), nil
			},
			verifyStats: func(t *testing.T, stats *verifierStats) {
				require.EqualValues(t, 0, stats.addressedCallSignatureValidationFail.Snapshot().Count())
				require.EqualValues(t, 1, stats.blockSignatureValidationFail.Snapshot().Count())
				require.EqualValues(t, 0, stats.messageParseFail.Snapshot().Count())
			},
			err: &common.AppError{Code: VerifyErrCode},
		},
	}

	for name, test := range tests {
		for _, withCache := range []bool{true, false} {
			if withCache {
				name += "_with_cache"
			} else {
				name += "_no_cache"
			}
			t.Run(name, func(t *testing.T) {
				var sigCache cache.Cacher[ids.ID, []byte]
				if withCache {
					sigCache = &cache.LRU[ids.ID, []byte]{Size: 100}
				} else {
					sigCache = &cache.Empty[ids.ID, []byte]{}
				}
				warpBackend, err := NewBackend(
					snowCtx.NetworkID,
					snowCtx.ChainID,
					warpSigner,
					blockClient,
					uptime.NoOpCalculator,
					validators.NoOpState,
					database,
					sigCache,
					nil,
				)
				require.NoError(t, err)
				handler := acp118.NewCachedHandler(sigCache, warpBackend, warpSigner)

				requestBytes, expectedResponse := test.setup()
				protoMsg := &sdk.SignatureRequest{Message: requestBytes}
				protoBytes, err := proto.Marshal(protoMsg)
				require.NoError(t, err)
				responseBytes, appErr := handler.AppRequest(context.Background(), ids.GenerateTestNodeID(), time.Time{}, protoBytes)
				if test.err != nil {
					require.NotNil(t, appErr)
					require.ErrorIs(t, test.err, appErr)
				} else {
					require.Nil(t, appErr)
				}

				test.verifyStats(t, warpBackend.(*backend).stats)

				// If the expected response is empty, assert that the handler returns an empty response and return early.
				if len(expectedResponse) == 0 {
					require.Len(t, responseBytes, 0, "expected response to be empty")
					return
				}
				// check cache is populated
				if withCache {
					require.NotZero(t, warpBackend.(*backend).signatureCache.Len())
				} else {
					require.Zero(t, warpBackend.(*backend).signatureCache.Len())
				}
				var response sdk.SignatureResponse
				err = proto.Unmarshal(responseBytes, &response)
				require.NoError(t, err, "error unmarshalling SignatureResponse")
				require.Equal(t, expectedResponse, response.Signature)
			})
		}
	}
}

func TestUptimeSignatures(t *testing.T) {
	database := memdb.New()
	snowCtx := utils.TestSnowContext()
	blsSecretKey, err := bls.NewSecretKey()
	require.NoError(t, err)
	warpSigner := avalancheWarp.NewSigner(blsSecretKey, snowCtx.NetworkID, snowCtx.ChainID)

	getUptimeMessageBytes := func(vID ids.ID, totalUptime uint64) ([]byte, *avalancheWarp.UnsignedMessage) {
		uptimePayload, err := messages.NewValidatorUptime(vID, 80)
		require.NoError(t, err)
		addressedCall, err := payload.NewAddressedCall([]byte{1, 2, 3}, uptimePayload.Bytes())
		require.NoError(t, err)
		unsignedMessage, err := avalancheWarp.NewUnsignedMessage(snowCtx.NetworkID, snowCtx.ChainID, addressedCall.Bytes())
		require.NoError(t, err)

		protoMsg := &sdk.SignatureRequest{Message: unsignedMessage.Bytes()}
		protoBytes, err := proto.Marshal(protoMsg)
		require.NoError(t, err)
		return protoBytes, unsignedMessage
	}

	for _, withCache := range []bool{true, false} {
		var sigCache cache.Cacher[ids.ID, []byte]
		if withCache {
			sigCache = &cache.LRU[ids.ID, []byte]{Size: 100}
		} else {
			sigCache = &cache.Empty[ids.ID, []byte]{}
		}
		state, err := validators.NewState(memdb.New())
		require.NoError(t, err)
		clk := &mockable.Clock{}
		uptimeManager := uptime.NewManager(state, clk)
		uptimeManager.StartTracking([]ids.NodeID{})
		warpBackend, err := NewBackend(snowCtx.NetworkID, snowCtx.ChainID, warpSigner, warptest.EmptyBlockClient, uptimeManager, state, database, sigCache, nil)
		require.NoError(t, err)
		handler := acp118.NewCachedHandler(sigCache, warpBackend, warpSigner)

		// not existing validationID
		vID := ids.GenerateTestID()
		protoBytes, _ := getUptimeMessageBytes(vID, 80)
		_, appErr := handler.AppRequest(context.Background(), ids.GenerateTestNodeID(), time.Time{}, protoBytes)
		require.ErrorIs(t, appErr, &common.AppError{Code: VerifyErrCode})
		require.Contains(t, appErr.Error(), "failed to get validator")

		// uptime is less than requested (not connected)
		validationID := ids.GenerateTestID()
		nodeID := ids.GenerateTestNodeID()
		require.NoError(t, state.AddValidator(validationID, nodeID, clk.Unix(), true))
		protoBytes, _ = getUptimeMessageBytes(validationID, 80)
		_, appErr = handler.AppRequest(context.Background(), nodeID, time.Time{}, protoBytes)
		require.ErrorIs(t, appErr, &common.AppError{Code: VerifyErrCode})
		require.Contains(t, appErr.Error(), "current uptime 0 is less than queried uptime 80")

		// uptime is less than requested (not enough)
		require.NoError(t, uptimeManager.Connect(nodeID))
		clk.Set(clk.Time().Add(40 * time.Second))
		protoBytes, _ = getUptimeMessageBytes(validationID, 80)
		_, appErr = handler.AppRequest(context.Background(), nodeID, time.Time{}, protoBytes)
		require.ErrorIs(t, appErr, &common.AppError{Code: VerifyErrCode})
		require.Contains(t, appErr.Error(), "current uptime 40 is less than queried uptime 80")

		// valid uptime
		clk.Set(clk.Time().Add(40 * time.Second))
		protoBytes, msg := getUptimeMessageBytes(validationID, 80)
		responseBytes, appErr := handler.AppRequest(context.Background(), nodeID, time.Time{}, protoBytes)
		require.Nil(t, appErr)
		expectedSignature, err := warpSigner.Sign(msg)
		require.NoError(t, err)
		response := &sdk.SignatureResponse{}
		require.NoError(t, proto.Unmarshal(responseBytes, response))
		require.Equal(t, expectedSignature[:], response.Signature)
	}
}
