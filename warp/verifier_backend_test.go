// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/cache"
	"github.com/ava-labs/avalanchego/cache/lru"
	"github.com/ava-labs/avalanchego/database/memdb"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/network/p2p/acp118"
	"github.com/ava-labs/avalanchego/proto/pb/sdk"
	"github.com/ava-labs/avalanchego/snow/engine/common"
	"github.com/ava-labs/avalanchego/utils/timer/mockable"
	"github.com/ava-labs/avalanchego/vms/evm/metrics/metricstest"
	avalancheWarp "github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/payload"
	"github.com/ava-labs/subnet-evm/plugin/evm/validators"
	stateinterfaces "github.com/ava-labs/subnet-evm/plugin/evm/validators/state/interfaces"
	"github.com/ava-labs/subnet-evm/utils/utilstest"
	"github.com/ava-labs/subnet-evm/warp/messages"
	"github.com/ava-labs/subnet-evm/warp/warptest"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestAddressedCallSignatures(t *testing.T) {
	metricstest.WithMetrics(t)

	database := memdb.New()
	snowCtx := utilstest.NewTestSnowContext(t)

	offChainPayload, err := payload.NewAddressedCall([]byte{1, 2, 3}, []byte{1, 2, 3})
	require.NoError(t, err)
	offchainMessage, err := avalancheWarp.NewUnsignedMessage(snowCtx.NetworkID, snowCtx.ChainID, offChainPayload.Bytes())
	require.NoError(t, err)
	offchainSignature, err := snowCtx.WarpSigner.Sign(offchainMessage)
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
				signature, err := snowCtx.WarpSigner.Sign(msg)
				require.NoError(t, err)

				backend.AddMessage(msg)
				return msg.Bytes(), signature[:]
			},
			verifyStats: func(t *testing.T, stats *verifierStats) {
				require.EqualValues(t, 0, stats.messageParseFail.Snapshot().Count())
				require.EqualValues(t, 0, stats.blockValidationFail.Snapshot().Count())
			},
		},
		"offchain message": {
			setup: func(_ Backend) (request []byte, expectedResponse []byte) {
				return offchainMessage.Bytes(), offchainSignature[:]
			},
			verifyStats: func(t *testing.T, stats *verifierStats) {
				require.EqualValues(t, 0, stats.messageParseFail.Snapshot().Count())
				require.EqualValues(t, 0, stats.blockValidationFail.Snapshot().Count())
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
				require.EqualValues(t, 1, stats.messageParseFail.Snapshot().Count())
				require.EqualValues(t, 0, stats.blockValidationFail.Snapshot().Count())
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
					sigCache = lru.NewCache[ids.ID, []byte](100)
				} else {
					sigCache = &cache.Empty[ids.ID, []byte]{}
				}
				warpBackend, err := NewBackend(
					snowCtx.NetworkID,
					snowCtx.ChainID,
					snowCtx.WarpSigner,
					warptest.EmptyBlockClient,
					nil,
					database,
					sigCache,
					[][]byte{offchainMessage.Bytes()},
				)
				require.NoError(t, err)
				handler := acp118.NewCachedHandler(sigCache, warpBackend, snowCtx.WarpSigner)

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
	metricstest.WithMetrics(t)

	database := memdb.New()
	snowCtx := utilstest.NewTestSnowContext(t)

	knownBlkID := ids.GenerateTestID()
	blockClient := warptest.MakeBlockClient(knownBlkID)

	toMessageBytes := func(id ids.ID) []byte {
		idPayload, err := payload.NewHash(id)
		if err != nil {
			panic(err)
		}

		msg, err := avalancheWarp.NewUnsignedMessage(snowCtx.NetworkID, snowCtx.ChainID, idPayload.Bytes())
		if err != nil {
			panic(err)
		}

		return msg.Bytes()
	}

	tests := map[string]struct {
		setup       func() (request []byte, expectedResponse []byte)
		verifyStats func(t *testing.T, stats *verifierStats)
		err         error
	}{
		"known block": {
			setup: func() (request []byte, expectedResponse []byte) {
				hashPayload, err := payload.NewHash(knownBlkID)
				require.NoError(t, err)
				unsignedMessage, err := avalancheWarp.NewUnsignedMessage(snowCtx.NetworkID, snowCtx.ChainID, hashPayload.Bytes())
				require.NoError(t, err)
				signature, err := snowCtx.WarpSigner.Sign(unsignedMessage)
				require.NoError(t, err)
				return toMessageBytes(knownBlkID), signature[:]
			},
			verifyStats: func(t *testing.T, stats *verifierStats) {
				require.EqualValues(t, 0, stats.blockValidationFail.Snapshot().Count())
				require.EqualValues(t, 0, stats.messageParseFail.Snapshot().Count())
			},
		},
		"unknown block": {
			setup: func() (request []byte, expectedResponse []byte) {
				unknownBlockID := ids.GenerateTestID()
				return toMessageBytes(unknownBlockID), nil
			},
			verifyStats: func(t *testing.T, stats *verifierStats) {
				require.EqualValues(t, 1, stats.blockValidationFail.Snapshot().Count())
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
					sigCache = lru.NewCache[ids.ID, []byte](100)
				} else {
					sigCache = &cache.Empty[ids.ID, []byte]{}
				}
				warpBackend, err := NewBackend(
					snowCtx.NetworkID,
					snowCtx.ChainID,
					snowCtx.WarpSigner,
					blockClient,
					warptest.NoOpValidatorReader{},
					database,
					sigCache,
					nil,
				)
				require.NoError(t, err)
				handler := acp118.NewCachedHandler(sigCache, warpBackend, snowCtx.WarpSigner)

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
	snowCtx := utilstest.NewTestSnowContext(t)

	getUptimeMessageBytes := func(sourceAddress []byte, vID ids.ID, totalUptime uint64) ([]byte, *avalancheWarp.UnsignedMessage) {
		uptimePayload, err := messages.NewValidatorUptime(vID, 80)
		require.NoError(t, err)
		addressedCall, err := payload.NewAddressedCall(sourceAddress, uptimePayload.Bytes())
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
			sigCache = lru.NewCache[ids.ID, []byte](100)
		} else {
			sigCache = &cache.Empty[ids.ID, []byte]{}
		}
		chainCtx := utilstest.NewTestSnowContext(t)
		clk := &mockable.Clock{}
		validatorsManager, err := validators.NewManager(chainCtx, memdb.New(), clk)
		require.NoError(t, err)
		lock := &sync.RWMutex{}
		newLockedValidatorManager := validators.NewLockedValidatorReader(validatorsManager, lock)
		validatorsManager.StartTracking([]ids.NodeID{})
		warpBackend, err := NewBackend(
			snowCtx.NetworkID,
			snowCtx.ChainID,
			snowCtx.WarpSigner,
			warptest.EmptyBlockClient,
			newLockedValidatorManager,
			database,
			sigCache,
			nil,
		)
		require.NoError(t, err)
		handler := acp118.NewCachedHandler(sigCache, warpBackend, snowCtx.WarpSigner)

		// sourceAddress nonZero
		protoBytes, _ := getUptimeMessageBytes([]byte{1, 2, 3}, ids.GenerateTestID(), 80)
		_, appErr := handler.AppRequest(context.Background(), ids.GenerateTestNodeID(), time.Time{}, protoBytes)
		require.ErrorIs(t, appErr, &common.AppError{Code: VerifyErrCode})
		require.Contains(t, appErr.Error(), "source address should be empty")

		// not existing validationID
		vID := ids.GenerateTestID()
		protoBytes, _ = getUptimeMessageBytes([]byte{}, vID, 80)
		_, appErr = handler.AppRequest(context.Background(), ids.GenerateTestNodeID(), time.Time{}, protoBytes)
		require.ErrorIs(t, appErr, &common.AppError{Code: VerifyErrCode})
		require.Contains(t, appErr.Error(), "failed to get validator")

		// uptime is less than requested (not connected)
		validationID := ids.GenerateTestID()
		nodeID := ids.GenerateTestNodeID()
		require.NoError(t, validatorsManager.AddValidator(stateinterfaces.Validator{
			ValidationID:   validationID,
			NodeID:         nodeID,
			Weight:         1,
			StartTimestamp: clk.Unix(),
			IsActive:       true,
			IsL1Validator:  true,
		}))
		protoBytes, _ = getUptimeMessageBytes([]byte{}, validationID, 80)
		_, appErr = handler.AppRequest(context.Background(), nodeID, time.Time{}, protoBytes)
		require.ErrorIs(t, appErr, &common.AppError{Code: VerifyErrCode})
		require.Contains(t, appErr.Error(), "current uptime 0 is less than queried uptime 80")

		// uptime is less than requested (not enough)
		require.NoError(t, validatorsManager.Connect(nodeID))
		clk.Set(clk.Time().Add(40 * time.Second))
		protoBytes, _ = getUptimeMessageBytes([]byte{}, validationID, 80)
		_, appErr = handler.AppRequest(context.Background(), nodeID, time.Time{}, protoBytes)
		require.ErrorIs(t, appErr, &common.AppError{Code: VerifyErrCode})
		require.Contains(t, appErr.Error(), "current uptime 40 is less than queried uptime 80")

		// valid uptime
		clk.Set(clk.Time().Add(40 * time.Second))
		protoBytes, msg := getUptimeMessageBytes([]byte{}, validationID, 80)
		responseBytes, appErr := handler.AppRequest(context.Background(), nodeID, time.Time{}, protoBytes)
		require.Nil(t, appErr)
		expectedSignature, err := snowCtx.WarpSigner.Sign(msg)
		require.NoError(t, err)
		response := &sdk.SignatureResponse{}
		require.NoError(t, proto.Unmarshal(responseBytes, response))
		require.Equal(t, expectedSignature[:], response.Signature)
	}
}
