// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"context"
	"encoding/binary"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/database/memdb"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/network/p2p"
	"github.com/ava-labs/avalanchego/network/p2p/gossip"
	"github.com/ava-labs/avalanchego/proto/pb/sdk"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/engine/common"
	"github.com/ava-labs/avalanchego/snow/validators"
	agoUtils "github.com/ava-labs/avalanchego/utils"
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"

	"google.golang.org/protobuf/proto"

	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/utils"
)

func TestEthTxGossip(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()
	snowCtx := utils.TestSnowContext()
	validatorState := &validators.TestState{}
	snowCtx.ValidatorState = validatorState

	responseSender := &common.FakeSender{
		SentAppResponse: make(chan []byte, 1),
	}
	vm := &VM{
		p2pSender: responseSender,
	}

	require.NoError(vm.Initialize(
		ctx,
		snowCtx,
		memdb.New(),
		[]byte(genesisJSONLatest),
		nil,
		nil,
		make(chan common.Message),
		nil,
		&common.SenderTest{},
	))
	require.NoError(vm.SetState(ctx, snow.NormalOp))

	defer func() {
		require.NoError(vm.Shutdown(ctx))
	}()

	// sender for the peer requesting gossip from [vm]
	peerSender := &common.FakeSender{
		SentAppRequest: make(chan []byte, 1),
	}

	network, err := p2p.NewNetwork(logging.NoLog{}, peerSender, prometheus.NewRegistry(), "")
	require.NoError(err)
	client := network.NewClient(ethTxGossipProtocol)

	// we only accept gossip requests from validators
	requestingNodeID := ids.GenerateTestNodeID()
	require.NoError(vm.Network.Connected(ctx, requestingNodeID, nil))
	validatorState.GetCurrentHeightF = func(context.Context) (uint64, error) {
		return 0, nil
	}
	validatorState.GetValidatorSetF = func(context.Context, uint64, ids.ID) (map[ids.NodeID]*validators.GetValidatorOutput, error) {
		return map[ids.NodeID]*validators.GetValidatorOutput{requestingNodeID: nil}, nil
	}

	// Ask the VM for any new transactions. We should get nothing at first.
	emptyBloomFilter, err := gossip.NewBloomFilter(prometheus.NewRegistry(), "", txGossipBloomMinTargetElements, txGossipBloomTargetFalsePositiveRate, txGossipBloomResetFalsePositiveRate)
	require.NoError(err)
	emptyBloomFilterBytes, _ := emptyBloomFilter.Marshal()
	request := &sdk.PullGossipRequest{
		Filter: emptyBloomFilterBytes,
		Salt:   agoUtils.RandomBytes(32),
	}

	requestBytes, err := proto.Marshal(request)
	require.NoError(err)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	onResponse := func(_ context.Context, nodeID ids.NodeID, responseBytes []byte, err error) {
		require.NoError(err)

		response := &sdk.PullGossipResponse{}
		require.NoError(proto.Unmarshal(responseBytes, response))
		require.Empty(response.Gossip)
		wg.Done()
	}
	require.NoError(client.AppRequest(ctx, set.Of(vm.ctx.NodeID), requestBytes, onResponse))
	require.NoError(vm.AppRequest(ctx, requestingNodeID, 1, time.Time{}, <-peerSender.SentAppRequest))
	require.NoError(network.AppResponse(ctx, snowCtx.NodeID, 1, <-responseSender.SentAppResponse))
	wg.Wait()

	// Issue a tx to the VM
	address := testEthAddrs[0]
	key := testKeys[0]
	tx := types.NewTransaction(0, address, big.NewInt(10), 21000, big.NewInt(testMinGasPrice), nil)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(vm.chainConfig.ChainID), key)
	require.NoError(err)

	errs := vm.txPool.AddLocals([]*types.Transaction{signedTx})
	require.Len(errs, 1)
	require.Nil(errs[0])

	// wait so we aren't throttled by the vm
	time.Sleep(5 * time.Second)

	marshaller := GossipEthTxMarshaller{}
	// Ask the VM for new transactions. We should get the newly issued tx.
	wg.Add(1)
	onResponse = func(_ context.Context, nodeID ids.NodeID, responseBytes []byte, err error) {
		require.NoError(err)

		response := &sdk.PullGossipResponse{}
		require.NoError(proto.Unmarshal(responseBytes, response))
		require.Len(response.Gossip, 1)

		gotTx, err := marshaller.UnmarshalGossip(response.Gossip[0])
		require.NoError(err)
		require.Equal(signedTx.Hash(), gotTx.Tx.Hash())

		wg.Done()
	}
	require.NoError(client.AppRequest(ctx, set.Of(vm.ctx.NodeID), requestBytes, onResponse))
	require.NoError(vm.AppRequest(ctx, requestingNodeID, 3, time.Time{}, <-peerSender.SentAppRequest))
	require.NoError(network.AppResponse(ctx, snowCtx.NodeID, 3, <-responseSender.SentAppResponse))
	wg.Wait()
}

func TestEthTxPushGossipOutbound(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()
	snowCtx := utils.TestSnowContext()
	sender := &common.FakeSender{
		SentAppGossip: make(chan []byte, 1),
	}

	vm := &VM{
		p2pSender:         sender,
		ethTxPullGossiper: gossip.NoOpGossiper{},
	}

	require.NoError(vm.Initialize(
		ctx,
		snowCtx,
		memdb.New(),
		[]byte(genesisJSONLatest),
		nil,
		nil,
		make(chan common.Message),
		nil,
		&common.FakeSender{},
	))
	require.NoError(vm.SetState(ctx, snow.NormalOp))

	address := testEthAddrs[0]
	key := testKeys[0]
	tx := types.NewTransaction(0, address, big.NewInt(10), 21000, big.NewInt(testMinGasPrice), nil)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(vm.chainConfig.ChainID), key)
	require.NoError(err)

	// issue a tx
	require.NoError(vm.txPool.AddLocal(signedTx))

	sent := <-sender.SentAppGossip
	got := &sdk.PushGossip{}

	// we should get a message that has the protocol prefix and the gossip
	// message
	require.Equal(byte(ethTxGossipProtocol), sent[0])
	require.NoError(proto.Unmarshal(sent[1:], got))

	marshaller := GossipEthTxMarshaller{}
	require.Len(got.Gossip, 1)
	gossipedTx, err := marshaller.UnmarshalGossip(got.Gossip[0])
	require.NoError(err)
	require.Equal(ids.ID(signedTx.Hash()), gossipedTx.GossipID())
}

// Tests that a gossiped tx is added to the mempool and forwarded
func TestEthTxPushGossipInbound(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()
	snowCtx := utils.TestSnowContext()

	sender := &common.FakeSender{
		SentAppGossip: make(chan []byte, 1),
	}
	vm := &VM{
		p2pSender:         sender,
		ethTxPullGossiper: gossip.NoOpGossiper{},
	}

	require.NoError(vm.Initialize(
		ctx,
		snowCtx,
		memdb.New(),
		[]byte(genesisJSONLatest),
		nil,
		nil,
		make(chan common.Message),
		nil,
		&common.FakeSender{},
	))
	require.NoError(vm.SetState(ctx, snow.NormalOp))

	address := testEthAddrs[0]
	key := testKeys[0]
	tx := types.NewTransaction(0, address, big.NewInt(10), 21000, big.NewInt(testMinGasPrice), nil)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(vm.chainConfig.ChainID), key)
	require.NoError(err)

	marshaller := GossipEthTxMarshaller{}
	gossipedTx := &GossipEthTx{
		Tx: signedTx,
	}
	gossipedTxBytes, err := marshaller.MarshalGossip(gossipedTx)
	require.NoError(err)

	inboundGossip := &sdk.PushGossip{
		Gossip: [][]byte{gossipedTxBytes},
	}

	inboundGossipBytes, err := proto.Marshal(inboundGossip)
	require.NoError(err)

	inboundGossipMsg := append(binary.AppendUvarint(nil, ethTxGossipProtocol), inboundGossipBytes...)
	require.NoError(vm.AppGossip(ctx, ids.EmptyNodeID, inboundGossipMsg))

	forwardedMsg := &sdk.PushGossip{}
	outboundGossipBytes := <-sender.SentAppGossip

	require.Equal(byte(ethTxGossipProtocol), outboundGossipBytes[0])
	require.NoError(proto.Unmarshal(outboundGossipBytes[1:], forwardedMsg))
	require.Len(forwardedMsg.Gossip, 1)

	forwardedTx, err := marshaller.UnmarshalGossip(forwardedMsg.Gossip[0])
	require.NoError(err)
	require.Equal(gossipedTx.GossipID(), forwardedTx.GossipID())
	require.True(vm.txPool.Has(signedTx.Hash()))
}
